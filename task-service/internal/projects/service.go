package projects

import (
	"context"
	"time"

	"github.com/task-manager/task-service/pkg/apperror"
	"github.com/task-manager/task-service/pkg/models"
)

// ServiceInterface defines the contract for project business logic
type ServiceInterface interface {
	GetByUser(ctx context.Context, userID string) ([]models.Project, error)
	GetByID(ctx context.Context, projectID string, userID string) (*models.Project, error)
	Create(ctx context.Context, req CreateProjectRequest, ownerID string) (*models.Project, error)
	Update(ctx context.Context, projectID string, req UpdateProjectRequest, userID string) (*models.Project, error)
	Delete(ctx context.Context, projectID string, userID string) error
	AddMember(ctx context.Context, projectID string, memberID string, userID string) error
	RemoveMember(ctx context.Context, projectID string, memberID string, userID string) error
}

// Service handles project business logic
type Service struct {
	repo RepositoryInterface
}

// NewService creates a new project service
func NewService(repo RepositoryInterface) *Service {
	return &Service{repo: repo}
}

// GetByUser returns all projects for a user
func (s *Service) GetByUser(ctx context.Context, userID string) ([]models.Project, error) {
	return s.repo.GetByUser(ctx, userID)
}

// GetByID returns a single project if user is a member
func (s *Service) GetByID(ctx context.Context, projectID string, userID string) (*models.Project, error) {
	project, err := s.repo.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	for _, m := range project.Members {
		if m == userID {
			return project, nil
		}
	}

	return nil, apperror.Forbidden("you are not a member of this project")
}

// Create validates and creates a new project
func (s *Service) Create(ctx context.Context, req CreateProjectRequest, ownerID string) (*models.Project, error) {
	if req.Name == "" {
		return nil, apperror.BadRequest("project name is required")
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
		return nil, apperror.Forbidden("only the project owner can update it")
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
		project.Name = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
		project.Description = *req.Description
	}

	if len(updates) == 0 {
		return nil, apperror.BadRequest("no fields to update")
	}

	if err := s.repo.Update(ctx, projectID, updates); err != nil {
		return nil, err
	}

	return project, nil
}

// Delete removes a project
func (s *Service) Delete(ctx context.Context, projectID string, userID string) error {
	project, err := s.repo.GetByID(ctx, projectID)
	if err != nil {
		return err
	}

	if project.OwnerID != userID {
		return apperror.Forbidden("only the project owner can delete it")
	}

	return s.repo.Delete(ctx, projectID)
}

// AddMember adds a user to a project using a transaction
func (s *Service) AddMember(ctx context.Context, projectID string, memberID string, userID string) error {
	project, err := s.repo.GetByID(ctx, projectID)
	if err != nil {
		return err
	}

	if project.OwnerID != userID {
		return apperror.Forbidden("only the project owner can add members")
	}

	return s.repo.AddMemberTx(ctx, projectID, memberID)
}

// RemoveMember removes a user from a project using a transaction
func (s *Service) RemoveMember(ctx context.Context, projectID string, memberID string, userID string) error {
	project, err := s.repo.GetByID(ctx, projectID)
	if err != nil {
		return err
	}

	if project.OwnerID != userID {
		return apperror.Forbidden("only the project owner can remove members")
	}

	return s.repo.RemoveMemberTx(ctx, projectID, memberID)
}
