package auth

import (
	"encoding/json"
	"net/http"

	"github.com/task-manager/task-service/pkg/apperror"
	"github.com/task-manager/task-service/pkg/response"
)

// Handler handles HTTP requests for authentication
type Handler struct {
	service ServiceInterface
}

// NewHandler creates a new auth handler
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

// Register handles POST /auth/register
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.service.Register(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, result)
}

// Login handles POST /auth/login
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.service.Login(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}
