package users

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/task-manager/task-service/pkg/models"
	"google.golang.org/api/iterator"
)

const collectionName = "users"

// Repository handles user data access in Firestore
type Repository struct {
	client *firestore.Client
}

// NewRepository creates a new users repository
func NewRepository(client *firestore.Client) *Repository {
	return &Repository{client: client}
}

// GetAll returns all users
func (r *Repository) GetAll(ctx context.Context) ([]models.User, error) {
	iter := r.client.Collection(collectionName).
		OrderBy("createdAt", firestore.Desc).
		Documents(ctx)
	defer iter.Stop()

	var users []models.User
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate users: %w", err)
		}

		var user models.User
		if err := doc.DataTo(&user); err != nil {
			return nil, fmt.Errorf("failed to parse user: %w", err)
		}
		user.ID = doc.Ref.ID
		users = append(users, user)
	}

	return users, nil
}

// UpdateRole changes a user's role
func (r *Repository) UpdateRole(ctx context.Context, userID string, role models.Role) error {
	_, err := r.client.Collection(collectionName).Doc(userID).Update(ctx, []firestore.Update{
		{Path: "role", Value: string(role)},
	})
	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}
	return nil
}
