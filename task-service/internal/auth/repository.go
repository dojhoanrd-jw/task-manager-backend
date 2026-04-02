package auth

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/task-manager/task-service/pkg/apperror"
	"github.com/task-manager/task-service/pkg/models"
	"google.golang.org/api/iterator"
)

const collectionName = "users"

// RepositoryInterface defines the contract for auth data access
type RepositoryInterface interface {
	Create(ctx context.Context, user *models.User) (string, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, userID string) (*models.User, error)
}

// Repository handles user data access in Firestore
type Repository struct {
	client *firestore.Client
}

// NewRepository creates a new auth repository
func NewRepository(client *firestore.Client) *Repository {
	return &Repository{client: client}
}

// Create adds a new user to Firestore
func (r *Repository) Create(ctx context.Context, user *models.User) (string, error) {
	ref, _, err := r.client.Collection(collectionName).Add(ctx, user)
	if err != nil {
		return "", apperror.Wrap(500, "failed to create user", err)
	}
	return ref.ID, nil
}

// GetByEmail finds a user by email
func (r *Repository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	iter := r.client.Collection(collectionName).
		Where("email", "==", email).
		Limit(1).
		Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, apperror.NotFound("user")
	}
	if err != nil {
		return nil, apperror.Wrap(500, "failed to query user", err)
	}

	var user models.User
	if err := doc.DataTo(&user); err != nil {
		return nil, apperror.Wrap(500, "failed to parse user", err)
	}
	user.ID = doc.Ref.ID
	return &user, nil
}

// GetByID finds a user by ID
func (r *Repository) GetByID(ctx context.Context, userID string) (*models.User, error) {
	doc, err := r.client.Collection(collectionName).Doc(userID).Get(ctx)
	if err != nil {
		return nil, apperror.NotFound("user")
	}

	var user models.User
	if err := doc.DataTo(&user); err != nil {
		return nil, apperror.Wrap(500, "failed to parse user", err)
	}
	user.ID = doc.Ref.ID
	return &user, nil
}
