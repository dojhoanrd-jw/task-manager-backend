package users

import (
	"context"
	"testing"

	"github.com/task-manager/task-service/pkg/models"
)

// mockRepository implements RepositoryInterface for testing
type mockRepository struct {
	users    map[string]*models.User
	lastRole models.Role
}

func newMockRepository() *mockRepository {
	return &mockRepository{users: make(map[string]*models.User)}
}

func (m *mockRepository) GetAll(ctx context.Context) ([]models.User, error) {
	var result []models.User
	for _, u := range m.users {
		result = append(result, *u)
	}
	return result, nil
}

func (m *mockRepository) UpdateRole(ctx context.Context, userID string, role models.Role) error {
	u, ok := m.users[userID]
	if !ok {
		return nil
	}
	u.Role = role
	m.lastRole = role
	return nil
}

func TestGetAllUsers(t *testing.T) {
	repo := newMockRepository()
	repo.users["u1"] = &models.User{ID: "u1", Name: "User 1", Role: models.RoleMember}
	repo.users["u2"] = &models.User{ID: "u2", Name: "User 2", Role: models.RoleAdmin}
	svc := NewService(repo)

	users, err := svc.GetAll(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}

func TestUpdateRoleValid(t *testing.T) {
	repo := newMockRepository()
	repo.users["u1"] = &models.User{ID: "u1", Role: models.RoleMember}
	svc := NewService(repo)

	err := svc.UpdateRole(context.Background(), "u1", "admin")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if repo.users["u1"].Role != models.RoleAdmin {
		t.Errorf("expected role 'admin', got '%s'", repo.users["u1"].Role)
	}
}

func TestUpdateRoleInvalid(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo)

	err := svc.UpdateRole(context.Background(), "u1", "superadmin")

	if err == nil {
		t.Fatal("expected error for invalid role")
	}
}

func TestUpdateRoleToViewer(t *testing.T) {
	repo := newMockRepository()
	repo.users["u1"] = &models.User{ID: "u1", Role: models.RoleMember}
	svc := NewService(repo)

	err := svc.UpdateRole(context.Background(), "u1", "viewer")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if repo.users["u1"].Role != models.RoleViewer {
		t.Errorf("expected role 'viewer', got '%s'", repo.users["u1"].Role)
	}
}
