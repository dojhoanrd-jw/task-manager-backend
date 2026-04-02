package tasks

import (
	"context"
	"time"

	"github.com/task-manager/task-service/pkg/apperror"
	"github.com/task-manager/task-service/pkg/models"
)

// MembershipChecker verifies if a user belongs to a project
type MembershipChecker interface {
	GetByID(ctx context.Context, projectID string) (*models.Project, error)
}

// ServiceInterface defines the contract for task business logic
type ServiceInterface interface {
	GetByProject(ctx context.Context, projectID string, userID string, limit int, lastID string) ([]models.Task, error)
	GetByID(ctx context.Context, taskID string) (*models.Task, error)
	Create(ctx context.Context, req CreateTaskRequest, projectID string, userID string) (*models.Task, error)
	Update(ctx context.Context, taskID string, req UpdateTaskRequest) (*models.Task, error)
	Delete(ctx context.Context, taskID string) error
}

// Service handles task business logic
type Service struct {
	repo       RepositoryInterface
	membership MembershipChecker
}

// NewService creates a new task service
func NewService(repo RepositoryInterface, membership MembershipChecker) *Service {
	return &Service{repo: repo, membership: membership}
}

// validateMembership checks if the user is a member of the project
func (s *Service) validateMembership(ctx context.Context, projectID string, userID string) error {
	project, err := s.membership.GetByID(ctx, projectID)
	if err != nil {
		return apperror.NotFound("project")
	}

	for _, m := range project.Members {
		if m == userID {
			return nil
		}
	}

	return apperror.Forbidden("you are not a member of this project")
}

// GetByProject returns paginated tasks for a project
func (s *Service) GetByProject(ctx context.Context, projectID string, userID string, limit int, lastID string) ([]models.Task, error) {
	if err := s.validateMembership(ctx, projectID, userID); err != nil {
		return nil, err
	}

	if limit <= 0 || limit > 50 {
		limit = 20
	}
	return s.repo.GetByProject(ctx, projectID, limit, lastID)
}

// GetByID returns a single task
func (s *Service) GetByID(ctx context.Context, taskID string) (*models.Task, error) {
	return s.repo.GetByID(ctx, taskID)
}

// Create validates and creates a new task
func (s *Service) Create(ctx context.Context, req CreateTaskRequest, projectID string, userID string) (*models.Task, error) {
	if req.Title == "" {
		return nil, apperror.BadRequest("title is required")
	}

	if err := s.validateMembership(ctx, projectID, userID); err != nil {
		return nil, err
	}

	task := &models.Task{
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
		ProjectID:   projectID,
		AssignedTo:  req.AssignedTo,
		CreatedAt:   time.Now(),
	}

	id, err := s.repo.Create(ctx, task)
	if err != nil {
		return nil, err
	}

	task.ID = id
	return task, nil
}

// Update modifies an existing task
func (s *Service) Update(ctx context.Context, taskID string, req UpdateTaskRequest) (*models.Task, error) {
	// Verify task exists before updating
	task, err := s.repo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})

	if req.Title != nil {
		updates["title"] = *req.Title
		task.Title = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
		task.Description = *req.Description
	}
	if req.Completed != nil {
		updates["completed"] = *req.Completed
		task.Completed = *req.Completed
	}
	if req.AssignedTo != nil {
		updates["assignedTo"] = *req.AssignedTo
		task.AssignedTo = *req.AssignedTo
	}

	if len(updates) == 0 {
		return nil, apperror.BadRequest("no fields to update")
	}

	if err := s.repo.Update(ctx, taskID, updates); err != nil {
		return nil, err
	}

	return task, nil
}

// Delete removes a task
func (s *Service) Delete(ctx context.Context, taskID string) error {
	return s.repo.Delete(ctx, taskID)
}
