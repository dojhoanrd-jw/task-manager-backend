package users

import (
	"encoding/json"
	"net/http"

	"github.com/task-manager/task-service/pkg/response"
)

// Handler handles HTTP requests for user management
type Handler struct {
	service ServiceInterface
}

// NewHandler creates a new users handler
func NewHandler(service ServiceInterface) *Handler {
	return &Handler{service: service}
}

// GetAll handles GET /users (admin only)
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAll(r.Context())
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, users)
}

// UpdateRole handles PUT /users/{userId}/role (admin only)
func (h *Handler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		response.Error(w, http.StatusBadRequest, "user ID is required")
		return
	}

	var req UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.service.UpdateRole(r.Context(), userID, req.Role); err != nil {
		response.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "role updated successfully"})
}
