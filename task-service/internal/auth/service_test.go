package auth

import (
	"context"
	"testing"

	"github.com/task-manager/task-service/pkg/apperror"
	"github.com/task-manager/task-service/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

// mockRepository implements RepositoryInterface for testing
type mockRepository struct {
	users  map[string]*models.User
	lastID int
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		users: make(map[string]*models.User),
	}
}

func (m *mockRepository) Create(ctx context.Context, user *models.User) (string, error) {
	m.lastID++
	id := "user-" + string(rune('0'+m.lastID))
	m.users[id] = user
	return id, nil
}

func (m *mockRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, apperror.NotFound("user")
}

func (m *mockRepository) GetByID(ctx context.Context, userID string) (*models.User, error) {
	u, ok := m.users[userID]
	if !ok {
		return nil, apperror.NotFound("user")
	}
	return u, nil
}

func TestRegister_Success(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo, "test-secret", "24h")

	result, err := svc.Register(context.Background(), RegisterRequest{
		Name:     "Test",
		Email:    "test@test.com",
		Password: "123456",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Token == "" {
		t.Error("expected token to be set")
	}
	if result.User.Email != "test@test.com" {
		t.Errorf("expected email 'test@test.com', got '%s'", result.User.Email)
	}
	if result.User.Role != "member" {
		t.Errorf("expected role 'member', got '%s'", result.User.Role)
	}
}

func TestRegister_EmptyFields(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo, "test-secret", "24h")

	_, err := svc.Register(context.Background(), RegisterRequest{
		Name:  "",
		Email: "test@test.com",
	})

	if err == nil {
		t.Fatal("expected error for empty fields")
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo, "test-secret", "24h")

	// Register first user
	_, _ = svc.Register(context.Background(), RegisterRequest{
		Name:     "User 1",
		Email:    "dup@test.com",
		Password: "123456",
	})

	// Try duplicate
	_, err := svc.Register(context.Background(), RegisterRequest{
		Name:     "User 2",
		Email:    "dup@test.com",
		Password: "123456",
	})

	if err == nil {
		t.Fatal("expected error for duplicate email")
	}
}

func TestLogin_Success(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo, "test-secret", "24h")

	// Create user with hashed password
	hashed, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	repo.users["user-1"] = &models.User{
		ID:       "user-1",
		Name:     "Test",
		Email:    "login@test.com",
		Password: string(hashed),
		Role:     models.RoleMember,
	}

	result, err := svc.Login(context.Background(), LoginRequest{
		Email:    "login@test.com",
		Password: "123456",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Token == "" {
		t.Error("expected token to be set")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo, "test-secret", "24h")

	hashed, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	repo.users["user-1"] = &models.User{
		ID:       "user-1",
		Email:    "login@test.com",
		Password: string(hashed),
		Role:     models.RoleMember,
	}

	_, err := svc.Login(context.Background(), LoginRequest{
		Email:    "login@test.com",
		Password: "wrong",
	})

	if err == nil {
		t.Fatal("expected error for wrong password")
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo, "test-secret", "24h")

	_, err := svc.Login(context.Background(), LoginRequest{
		Email:    "noexist@test.com",
		Password: "123456",
	})

	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}
