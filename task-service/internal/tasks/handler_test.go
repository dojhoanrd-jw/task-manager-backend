package tasks

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/task-manager/task-service/pkg/models"
)

// mockServiceForHandler implements ServiceInterface for handler testing
type mockServiceForHandler struct{}

func (m *mockServiceForHandler) GetByProject(ctx context.Context, projectID string, userID string, limit int, lastID string) ([]models.Task, error) {
	return []models.Task{{ID: "t1", Title: "Task 1", ProjectID: projectID}}, nil
}

func (m *mockServiceForHandler) GetByID(ctx context.Context, taskID string) (*models.Task, error) {
	return &models.Task{ID: taskID, Title: "Task 1"}, nil
}

func (m *mockServiceForHandler) Create(ctx context.Context, req CreateTaskRequest, projectID string, userID string) (*models.Task, error) {
	return &models.Task{ID: "new-task", Title: req.Title, ProjectID: projectID}, nil
}

func (m *mockServiceForHandler) Update(ctx context.Context, taskID string, req UpdateTaskRequest) (*models.Task, error) {
	return &models.Task{ID: taskID, Title: "Updated"}, nil
}

func (m *mockServiceForHandler) Delete(ctx context.Context, taskID string) error {
	return nil
}

func TestHandlerGetByProject(t *testing.T) {
	handler := NewHandler(&mockServiceForHandler{})

	req := httptest.NewRequest(http.MethodGet, "/projects/proj1/tasks?limit=10", nil)
	req.SetPathValue("projectId", "proj1")
	req.Header.Set("X-User-ID", "user1")
	rec := httptest.NewRecorder()

	handler.GetByProject(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestHandlerGetByProjectMissingID(t *testing.T) {
	handler := NewHandler(&mockServiceForHandler{})

	req := httptest.NewRequest(http.MethodGet, "/projects//tasks", nil)
	rec := httptest.NewRecorder()

	handler.GetByProject(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestHandlerCreate(t *testing.T) {
	handler := NewHandler(&mockServiceForHandler{})

	body, _ := json.Marshal(CreateTaskRequest{Title: "New Task", Description: "Desc"})
	req := httptest.NewRequest(http.MethodPost, "/projects/proj1/tasks", bytes.NewReader(body))
	req.SetPathValue("projectId", "proj1")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "user1")
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", rec.Code)
	}
}

func TestHandlerDelete(t *testing.T) {
	handler := NewHandler(&mockServiceForHandler{})

	req := httptest.NewRequest(http.MethodDelete, "/projects/proj1/tasks/t1", nil)
	req.SetPathValue("taskId", "t1")
	rec := httptest.NewRecorder()

	handler.Delete(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}
