package tasks

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/task-manager/task-service/pkg/apperror"
	"github.com/task-manager/task-service/pkg/models"
	"github.com/task-manager/task-service/pkg/response"
)

const (
	errTaskIDRequired    = "task ID is required"
	errProjectIDRequired = "project ID is required"
	errInvalidBody       = "invalid request body"
)

// Handler handles HTTP requests for tasks
type Handler struct {
	service ServiceInterface
}

// NewHandler creates a new task handler
func NewHandler(service ServiceInterface) *Handler {
	return &Handler{service: service}
}

// handleError writes the appropriate HTTP response based on error type
func handleError(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*apperror.AppError); ok {
		response.Error(w, appErr.Code, appErr.Message)
		return
	}
	response.Error(w, http.StatusInternalServerError, "internal server error")
}

// GetByProject handles GET /projects/{projectId}/tasks
func (h *Handler) GetByProject(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	if projectID == "" {
		response.Error(w, http.StatusBadRequest, errProjectIDRequired)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	lastID := r.URL.Query().Get("lastId")

	userID := r.Header.Get("X-User-ID")
	tasks, err := h.service.GetByProject(r.Context(), projectID, userID, limit, lastID)
	if err != nil {
		handleError(w, err)
		return
	}

	if tasks == nil {
		tasks = []models.Task{}
	}

	response.JSON(w, http.StatusOK, tasks)
}

// GetByID handles GET /projects/{projectId}/tasks/{taskId}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("taskId")
	if taskID == "" {
		response.Error(w, http.StatusBadRequest, errTaskIDRequired)
		return
	}

	task, err := h.service.GetByID(r.Context(), taskID)
	if err != nil {
		handleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, task)
}

// Create handles POST /projects/{projectId}/tasks
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	if projectID == "" {
		response.Error(w, http.StatusBadRequest, errProjectIDRequired)
		return
	}

	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, errInvalidBody)
		return
	}

	userID := r.Header.Get("X-User-ID")
	task, err := h.service.Create(r.Context(), req, projectID, userID)
	if err != nil {
		handleError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, task)
}

// Update handles PUT /projects/{projectId}/tasks/{taskId}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("taskId")
	if taskID == "" {
		response.Error(w, http.StatusBadRequest, errTaskIDRequired)
		return
	}

	var req UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, errInvalidBody)
		return
	}

	task, err := h.service.Update(r.Context(), taskID, req)
	if err != nil {
		handleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, task)
}

// Delete handles DELETE /projects/{projectId}/tasks/{taskId}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("taskId")
	if taskID == "" {
		response.Error(w, http.StatusBadRequest, errTaskIDRequired)
		return
	}

	if err := h.service.Delete(r.Context(), taskID); err != nil {
		handleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "task deleted successfully"})
}
