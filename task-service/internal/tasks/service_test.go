package tasks

import (
	"context"
	"fmt"
	"testing"

	"github.com/task-manager/task-service/pkg/models"
)

// mockRepository implements RepositoryInterface for testing
type mockRepository struct {
	tasks   map[string]*models.Task
	lastID  int
	created *models.Task
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		tasks: make(map[string]*models.Task),
	}
}

func (m *mockRepository) GetByProject(ctx context.Context, projectID string, limit int, lastID string) ([]models.Task, error) {
	var result []models.Task
	for _, t := range m.tasks {
		if t.ProjectID == projectID {
			result = append(result, *t)
		}
	}
	return result, nil
}

func (m *mockRepository) GetByID(ctx context.Context, taskID string) (*models.Task, error) {
	t, ok := m.tasks[taskID]
	if !ok {
		return nil, fmt.Errorf("task not found")
	}
	return t, nil
}

func (m *mockRepository) Create(ctx context.Context, task *models.Task) (string, error) {
	m.lastID++
	id := "task-" + string(rune('0'+m.lastID))
	m.tasks[id] = task
	m.created = task
	return id, nil
}

func (m *mockRepository) Update(ctx context.Context, taskID string, updates map[string]interface{}) error {
	t, ok := m.tasks[taskID]
	if !ok {
		return nil
	}
	for key, val := range updates {
		switch key {
		case "title":
			t.Title = val.(string)
		case "completed":
			t.Completed = val.(bool)
		}
	}
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, taskID string) error {
	delete(m.tasks, taskID)
	return nil
}

func TestServiceCreate_Success(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo)

	task, err := svc.Create(context.Background(), CreateTaskRequest{
		Title:       "Test Task",
		Description: "Description",
	}, "project-1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if task.Title != "Test Task" {
		t.Errorf("expected title 'Test Task', got '%s'", task.Title)
	}
	if task.ProjectID != "project-1" {
		t.Errorf("expected projectID 'project-1', got '%s'", task.ProjectID)
	}
	if task.Completed {
		t.Error("expected completed to be false")
	}
	if task.ID == "" {
		t.Error("expected task ID to be set")
	}
}

func TestServiceCreate_EmptyTitle(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo)

	_, err := svc.Create(context.Background(), CreateTaskRequest{
		Title: "",
	}, "project-1")

	if err == nil {
		t.Fatal("expected error for empty title")
	}
}

func TestServiceUpdate_Success(t *testing.T) {
	repo := newMockRepository()
	repo.tasks["task-1"] = &models.Task{
		ID:        "task-1",
		Title:     "Original",
		Completed: false,
	}
	svc := NewService(repo)

	completed := true
	_, err := svc.Update(context.Background(), "task-1", UpdateTaskRequest{
		Completed: &completed,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !repo.tasks["task-1"].Completed {
		t.Error("expected task to be completed")
	}
}

func TestServiceUpdate_NoFields(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo)

	_, err := svc.Update(context.Background(), "task-1", UpdateTaskRequest{})

	if err == nil {
		t.Fatal("expected error for empty update")
	}
}

func TestServiceDelete(t *testing.T) {
	repo := newMockRepository()
	repo.tasks["task-1"] = &models.Task{ID: "task-1"}
	svc := NewService(repo)

	err := svc.Delete(context.Background(), "task-1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if _, ok := repo.tasks["task-1"]; ok {
		t.Error("expected task to be deleted")
	}
}

func TestServiceGetByProject_DefaultLimit(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo)

	// limit 0 should default to 20
	_, err := svc.GetByProject(context.Background(), "project-1", 0, "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
