package users

import (
	"context"

	"github.com/task-manager/task-service/pkg/apperror"
	"github.com/task-manager/task-service/pkg/models"
)

// ServiceInterface defines the contract for user business logic
type ServiceInterface interface {
	GetAll(ctx context.Context) ([]models.User, error)
	UpdateRole(ctx context.Context, userID string, roleStr string) error
}

// Service handles user business logic
type Service struct {
	repo RepositoryInterface
}

// NewService creates a new users service
func NewService(repo RepositoryInterface) *Service {
	return &Service{repo: repo}
}

// GetAll returns all users (admin only)
func (s *Service) GetAll(ctx context.Context) ([]models.User, error) {
	return s.repo.GetAll(ctx)
}

// UpdateRole changes a user's role with validation
func (s *Service) UpdateRole(ctx context.Context, userID string, roleStr string) error {
	role := models.Role(roleStr)

	switch role {
	case models.RoleAdmin, models.RoleMember, models.RoleViewer:
		// valid role
	default:
		return apperror.BadRequest("invalid role, must be: admin, member or viewer")
	}

	return s.repo.UpdateRole(ctx, userID, role)
}
