package users

import (
	"context"
	"errors"

	"github.com/task-manager/task-service/pkg/models"
)

// Service handles user business logic
type Service struct {
	repo *Repository
}

// NewService creates a new users service
func NewService(repo *Repository) *Service {
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
		return errors.New("invalid role, must be: admin, member or viewer")
	}

	return s.repo.UpdateRole(ctx, userID, role)
}
