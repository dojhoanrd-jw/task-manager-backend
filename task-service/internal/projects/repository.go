package projects

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/task-manager/task-service/pkg/models"
	"google.golang.org/api/iterator"
)

const collectionName = "projects"

// Repository handles project data access in Firestore
type Repository struct {
	client *firestore.Client
}

// NewRepository creates a new project repository
func NewRepository(client *firestore.Client) *Repository {
	return &Repository{client: client}
}

// GetByUser returns all projects where the user is owner or member
func (r *Repository) GetByUser(ctx context.Context, userID string) ([]models.Project, error) {
	// Get projects where user is a member
	iter := r.client.Collection(collectionName).
		Where("members", "array-contains", userID).
		OrderBy("createdAt", firestore.Desc).
		Documents(ctx)
	defer iter.Stop()

	var projects []models.Project
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate projects: %w", err)
		}

		var project models.Project
		if err := doc.DataTo(&project); err != nil {
			return nil, fmt.Errorf("failed to parse project: %w", err)
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
		return nil, fmt.Errorf("project not found: %w", err)
	}

	var project models.Project
	if err := doc.DataTo(&project); err != nil {
		return nil, fmt.Errorf("failed to parse project: %w", err)
	}
	project.ID = doc.Ref.ID
	return &project, nil
}

// Create adds a new project to Firestore
func (r *Repository) Create(ctx context.Context, project *models.Project) (string, error) {
	ref, _, err := r.client.Collection(collectionName).Add(ctx, project)
	if err != nil {
		return "", fmt.Errorf("failed to create project: %w", err)
	}
	return ref.ID, nil
}

// Update modifies an existing project in Firestore
func (r *Repository) Update(ctx context.Context, projectID string, updates []firestore.Update) error {
	_, err := r.client.Collection(collectionName).Doc(projectID).Update(ctx, updates)
	if err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}
	return nil
}

// Delete removes a project from Firestore
func (r *Repository) Delete(ctx context.Context, projectID string) error {
	_, err := r.client.Collection(collectionName).Doc(projectID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}
	return nil
}
