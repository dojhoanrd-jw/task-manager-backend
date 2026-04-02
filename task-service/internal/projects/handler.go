package projects

import (
	"encoding/json"
	"net/http"

	"github.com/task-manager/task-service/pkg/models"
	"github.com/task-manager/task-service/pkg/response"
)

const (
	headerUserID         = "X-User-ID"
	errProjectIDRequired = "project ID is required"
	errInvalidBody       = "invalid request body"
	errUnauthorized      = "unauthorized"
)

// Handler handles HTTP requests for projects
type Handler struct {
	service ServiceInterface
}

// NewHandler creates a new project handler
func NewHandler(service ServiceInterface) *Handler {
	return &Handler{service: service}
}

// GetByUser handles GET /projects
func (h *Handler) GetByUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get(headerUserID)
	if userID == "" {
		response.Error(w, http.StatusUnauthorized, errUnauthorized)
		return
	}

	projects, err := h.service.GetByUser(r.Context(), userID)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	if projects == nil {
		projects = []models.Project{}
	}

	response.JSON(w, http.StatusOK, projects)
}

// GetByID handles GET /projects/{projectId}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	if projectID == "" {
		response.Error(w, http.StatusBadRequest, errProjectIDRequired)
		return
	}

	project, err := h.service.GetByID(r.Context(), projectID)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, project)
}

// Create handles POST /projects
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get(headerUserID)
	if userID == "" {
		response.Error(w, http.StatusUnauthorized, errUnauthorized)
		return
	}

	var req CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, errInvalidBody)
		return
	}

	project, err := h.service.Create(r.Context(), req, userID)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, project)
}

// Update handles PUT /projects/{projectId}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	userID := r.Header.Get(headerUserID)

	if projectID == "" {
		response.Error(w, http.StatusBadRequest, errProjectIDRequired)
		return
	}

	var req UpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, errInvalidBody)
		return
	}

	project, err := h.service.Update(r.Context(), projectID, req, userID)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, project)
}

// Delete handles DELETE /projects/{projectId}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	userID := r.Header.Get(headerUserID)

	if projectID == "" {
		response.Error(w, http.StatusBadRequest, errProjectIDRequired)
		return
	}

	if err := h.service.Delete(r.Context(), projectID, userID); err != nil {
		response.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "project deleted successfully"})
}

// AddMember handles POST /projects/{projectId}/members
func (h *Handler) AddMember(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	userID := r.Header.Get(headerUserID)

	if projectID == "" {
		response.Error(w, http.StatusBadRequest, errProjectIDRequired)
		return
	}

	var req AddMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, errInvalidBody)
		return
	}

	if err := h.service.AddMember(r.Context(), projectID, req.UserID, userID); err != nil {
		response.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "member added successfully"})
}

// RemoveMember handles DELETE /projects/{projectId}/members/{userId}
func (h *Handler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	memberID := r.PathValue("userId")
	userID := r.Header.Get(headerUserID)

	if projectID == "" {
		response.Error(w, http.StatusBadRequest, errProjectIDRequired)
		return
	}

	if err := h.service.RemoveMember(r.Context(), projectID, memberID, userID); err != nil {
		response.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "member removed successfully"})
}
