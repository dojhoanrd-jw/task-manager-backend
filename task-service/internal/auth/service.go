package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/task-manager/task-service/pkg/apperror"
	"github.com/task-manager/task-service/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

// ServiceInterface defines the contract for auth business logic
type ServiceInterface interface {
	Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
}

// Service handles authentication business logic
type Service struct {
	repo      RepositoryInterface
	jwtSecret string
	jwtExp    time.Duration
}

// NewService creates a new auth service
func NewService(repo RepositoryInterface, jwtSecret string, jwtExpiration string) *Service {
	exp, err := time.ParseDuration(jwtExpiration)
	if err != nil {
		exp = 24 * time.Hour
	}

	return &Service{
		repo:      repo,
		jwtSecret: jwtSecret,
		jwtExp:    exp,
	}
}

// Register creates a new user account
func (s *Service) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return nil, apperror.BadRequest("name, email and password are required")
	}

	// Check if user already exists
	existing, _ := s.repo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return nil, apperror.Conflict("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperror.Wrap(500, "failed to hash password", err)
	}

	user := &models.User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Role:      models.RoleMember,
		CreatedAt: time.Now(),
	}

	id, err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = id

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User: UserInfo{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  string(user.Role),
		},
	}, nil
}

// Login authenticates a user and returns a JWT token
func (s *Service) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, apperror.BadRequest("email and password are required")
	}

	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, apperror.New(401, "invalid credentials")
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, apperror.New(401, "invalid credentials")
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User: UserInfo{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  string(user.Role),
		},
	}, nil
}

// generateToken creates a JWT token for a user
func (s *Service) generateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"userId": user.ID,
		"email":  user.Email,
		"role":   string(user.Role),
		"exp":    time.Now().Add(s.jwtExp).Unix(),
		"iat":    time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", apperror.Wrap(500, "failed to generate token", err)
	}
	return tokenStr, nil
}
