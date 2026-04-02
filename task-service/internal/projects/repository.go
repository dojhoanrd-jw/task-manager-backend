package projects

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/task-manager/task-service/pkg/apperror"
	"github.com/task-manager/task-service/pkg/models"
	"google.golang.org/api/iterator"
)

const (
	collectionName     = "projects"
	errFailedToParse   = "failed to parse project"
)

// RepositoryInterface defines the contract for project data access
type RepositoryInterface interface {
	GetByUser(ctx context.Context, userID string) ([]models.Project, error)
	GetByID(ctx context.Context, projectID string) (*models.Project, error)
	Create(ctx context.Context, project *models.Project) (string, error)
	Update(ctx context.Context, projectID string, updates map[string]interface{}) error
	Delete(ctx context.Context, projectID string) error
	AddMemberTx(ctx context.Context, projectID string, memberID string) error
	RemoveMemberTx(ctx context.Context, projectID string, memberID string) error
}

// Repository handles project data access in Firestore
type Repository struct {
	client *firestore.Client
}

// NewRepository creates a new project repository
func NewRepository(client *firestore.Client) *Repository {
	return &Repository{client: client}
}

const maxProjectsPerUser = 100

// GetByUser returns projects where the user is owner or member (limited)
func (r *Repository) GetByUser(ctx context.Context, userID string) ([]models.Project, error) {
	iter := r.client.Collection(collectionName).
		Where("members", "array-contains", userID).
		OrderBy("createdAt", firestore.Desc).
		Limit(maxProjectsPerUser).
		Documents(ctx)
	defer iter.Stop()

	var projects []models.Project
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, apperror.Wrap(500, "failed to iterate projects", err)
		}

		var project models.Project
		if err := doc.DataTo(&project); err != nil {
			return nil, apperror.Wrap(500, errFailedToParse, err)
		}
		project.ID = doc.Ref.ID
		projects = append(projects, project)
	}

	return projects, nil
}

// GetByID returns a single project by ID
func (r *Repository) GetByID(ctx context.Context, projectID string) (*models.Project, error) {
	doc, err := r.client.Collection(collectionName).Doc(projectID).Get(ctx)
	if err != nil {
		return nil, apperror.NotFound("project")
	}

	var project models.Project
	if err := doc.DataTo(&project); err != nil {
		return nil, apperror.Wrap(500, errFailedToParse, err)
	}
	project.ID = doc.Ref.ID
	return &project, nil
}

// Create adds a new project to Firestore
func (r *Repository) Create(ctx context.Context, project *models.Project) (string, error) {
	ref, _, err := r.client.Collection(collectionName).Add(ctx, project)
	if err != nil {
		return "", apperror.Wrap(500, "failed to create project", err)
	}
	return ref.ID, nil
}

// Update modifies an existing project in Firestore
func (r *Repository) Update(ctx context.Context, projectID string, updates map[string]interface{}) error {
	var fsUpdates []firestore.Update
	for path, value := range updates {
		fsUpdates = append(fsUpdates, firestore.Update{Path: path, Value: value})
	}

	_, err := r.client.Collection(collectionName).Doc(projectID).Update(ctx, fsUpdates)
	if err != nil {
		return apperror.Wrap(500, "failed to update project", err)
	}
	return nil
}

// Delete removes a project from Firestore
func (r *Repository) Delete(ctx context.Context, projectID string) error {
	_, err := r.client.Collection(collectionName).Doc(projectID).Delete(ctx)
	if err != nil {
		return apperror.Wrap(500, "failed to delete project", err)
	}
	return nil
}

// AddMemberTx adds a member using a Firestore transaction for atomicity
func (r *Repository) AddMemberTx(ctx context.Context, projectID string, memberID string) error {
	ref := r.client.Collection(collectionName).Doc(projectID)

	return r.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(ref)
		if err != nil {
			return apperror.NotFound("project")
		}

		var project models.Project
		if err := doc.DataTo(&project); err != nil {
			return apperror.Wrap(500, errFailedToParse, err)
		}

		// Check if already a member
		for _, m := range project.Members {
			if m == memberID {
				return apperror.Conflict("user is already a member")
			}
		}

		return tx.Update(ref, []firestore.Update{
			{Path: "members", Value: firestore.ArrayUnion(memberID)},
		})
	})
}

// RemoveMemberTx removes a member using a Firestore transaction for atomicity
func (r *Repository) RemoveMemberTx(ctx context.Context, projectID string, memberID string) error {
	ref := r.client.Collection(collectionName).Doc(projectID)

	return r.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(ref)
		if err != nil {
			return apperror.NotFound("project")
		}

		var project models.Project
		if err := doc.DataTo(&project); err != nil {
			return apperror.Wrap(500, errFailedToParse, err)
		}

		if memberID == project.OwnerID {
			return apperror.BadRequest("cannot remove the project owner")
		}

		return tx.Update(ref, []firestore.Update{
			{Path: "members", Value: firestore.ArrayRemove(memberID)},
		})
	})
}
