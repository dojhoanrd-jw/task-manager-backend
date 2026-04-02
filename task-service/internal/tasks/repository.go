package tasks

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/task-manager/task-service/pkg/models"
	"google.golang.org/api/iterator"
)

const collectionName = "tasks"

// Repository handles task data access in Firestore
type Repository struct {
	client *firestore.Client
}

// NewRepository creates a new task repository
func NewRepository(client *firestore.Client) *Repository {
	return &Repository{client: client}
}

// GetByProject returns paginated tasks for a project
func (r *Repository) GetByProject(ctx context.Context, projectID string, limit int, lastID string) ([]models.Task, error) {
	query := r.client.Collection(collectionName).
		Where("projectId", "==", projectID).
		OrderBy("createdAt", firestore.Desc).
		Limit(limit)

	// Cursor-based pagination
	if lastID != "" {
		lastDoc, err := r.client.Collection(collectionName).Doc(lastID).Get(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get cursor document: %w", err)
		}
		query = query.StartAfter(lastDoc.Data()["createdAt"])
	}

	iter := query.Documents(ctx)
	defer iter.Stop()

	var tasks []models.Task
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate tasks: %w", err)
		}

		var task models.Task
		if err := doc.DataTo(&task); err != nil {
			return nil, fmt.Errorf("failed to parse task: %w", err)
		}
		task.ID = doc.Ref.ID
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// GetByID returns a single task by ID
func (r *Repository) GetByID(ctx context.Context, taskID string) (*models.Task, error) {
	doc, err := r.client.Collection(collectionName).Doc(taskID).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	var task models.Task
	if err := doc.DataTo(&task); err != nil {
		return nil, fmt.Errorf("failed to parse task: %w", err)
	}
	task.ID = doc.Ref.ID
	return &task, nil
}

// Create adds a new task to Firestore
func (r *Repository) Create(ctx context.Context, task *models.Task) (string, error) {
	ref, _, err := r.client.Collection(collectionName).Add(ctx, task)
	if err != nil {
		return "", fmt.Errorf("failed to create task: %w", err)
	}
	return ref.ID, nil
}

// Update modifies an existing task in Firestore
func (r *Repository) Update(ctx context.Context, taskID string, updates []firestore.Update) error {
	_, err := r.client.Collection(collectionName).Doc(taskID).Update(ctx, updates)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}
	return nil
}

// Delete removes a task from Firestore
func (r *Repository) Delete(ctx context.Context, taskID string) error {
	_, err := r.client.Collection(collectionName).Doc(taskID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	return nil
}
