package projects

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/task-manager/task-service/pkg/models"
)

// Service handles project business logic
type Service struct {
	repo *Repository
}

// NewService creates a new project service
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// GetByUser returns all projects for a user
func (s *Service) GetByUser(ctx context.Context, userID string) ([]models.Project, error) {
	return s.repo.GetByUser(ctx, userID)
}

// GetByID returns a single project
func (s *Service) GetByID(ctx context.Context, projectID string) (*models.Project, error) {
	return s.repo.GetByID(ctx, projectID)
}

// Create validates and creates a new project
func (s *Service) Create(ctx context.Context, req CreateProjectRequest, ownerID string) (*models.Project, error) {
	if req.Name == "" {
		return nil, errors.New("project name is required")
	}

	project := &models.Project{
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     ownerID,
		Members:     []string{ownerID},
		CreatedAt:   time.Now(),
	}

	id, err := s.repo.Create(ctx, project)
	if err != nil {
		return nil, err
	}

	project.ID = id
	return project, nil
}

// Update modifies an existing project
func (s *Service) Update(ctx context.Context, projectID string, req UpdateProjectRequest, userID string) (*models.Project, error) {
	project, err := s.repo.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if project.OwnerID != userID {
		return nil, errors.New("only the project owner can update it")
	}

	var updates []firestore.Update
	if req.Name != nil {
		updates = append(updates, firestore.Update{Path: "name", Value: *req.Name})
	}
	if req.Description != nil {
		updates = append(updates, firestore.Update{Path: "description", Value: *req.Description})
	}

	if len(updates) == 0 {
		return nil, errors.New("no fields to update")
	}

	if err := s.repo.Update(ctx, projectID, updates); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, projectID)
}

// Delete removes a project
func (s *Service) Delete(ctx context.Context, projectID string, userID string) error {
	project, err := s.repo.GetByID(ctx, projectID)
	if err != nil {
		return err
	}

	if project.OwnerID != userID {
		return errors.New("only the project owner can delete it")
	}

	return s.repo.Delete(ctx, projectID)
}

// AddMember adds a user to a project
func (s *Service) AddMember(ctx context.Context, projectID string, memberID string, userID string) error {
	project, err := s.repo.GetByID(ctx, projectID)
	if err != nil {
		return err
	}

	if project.OwnerID != userID {
		return errors.New("only the project owner can add members")
	}

	// Check if already a member
	for _, m := range project.Members {
		if m == memberID {
			return errors.New("user is already a member")
		}
	}

	updates := []firestore.Update{
		{Path: "members", Value: firestore.ArrayUnion(memberID)},
	}
	return s.repo.Update(ctx, projectID, updates)
}

// RemoveMember removes a user from a project
func (s *Service) RemoveMember(ctx context.Context, projectID string, memberID string, userID string) error {
	project, err := s.repo.GetByID(ctx, projectID)
	if err != nil {
		return err
	}

	if project.OwnerID != userID {
		return errors.New("only the project owner can remove members")
	}

	if memberID == project.OwnerID {
		return errors.New("cannot remove the project owner")
	}

	updates := []firestore.Update{
		{Path: "members", Value: firestore.ArrayRemove(memberID)},
	}
	return s.repo.Update(ctx, projectID, updates)
}
