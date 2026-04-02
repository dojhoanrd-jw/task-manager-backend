package tasks

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/task-manager/task-service/pkg/models"
)

// Service handles task business logic
type Service struct {
	repo *Repository
}

// NewService creates a new task service
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// GetByProject returns paginated tasks for a project
func (s *Service) GetByProject(ctx context.Context, projectID string, limit int, lastID string) ([]models.Task, error) {
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
func (s *Service) Create(ctx context.Context, req CreateTaskRequest, projectID string) (*models.Task, error) {
	if req.Title == "" {
		return nil, errors.New("title is required")
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
	var updates []firestore.Update

	if req.Title != nil {
		updates = append(updates, firestore.Update{Path: "title", Value: *req.Title})
	}
	if req.Description != nil {
		updates = append(updates, firestore.Update{Path: "description", Value: *req.Description})
	}
	if req.Completed != nil {
		updates = append(updates, firestore.Update{Path: "completed", Value: *req.Completed})
	}
	if req.AssignedTo != nil {
		updates = append(updates, firestore.Update{Path: "assignedTo", Value: *req.AssignedTo})
	}

	if len(updates) == 0 {
		return nil, errors.New("no fields to update")
	}

	if err := s.repo.Update(ctx, taskID, updates); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, taskID)
}

// Delete removes a task
func (s *Service) Delete(ctx context.Context, taskID string) error {
	return s.repo.Delete(ctx, taskID)
}
