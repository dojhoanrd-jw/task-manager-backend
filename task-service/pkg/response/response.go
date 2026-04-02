package response

import (
	"encoding/json"
	"net/http"

	"github.com/task-manager/task-service/pkg/apperror"
)

// JSON sends a JSON response with the given status code
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Error sends an error response
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, map[string]string{"error": message})
}

// HandleError writes the appropriate HTTP response based on error type
func HandleError(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*apperror.AppError); ok {
		Error(w, appErr.Code, appErr.Message)
		return
	}
	Error(w, http.StatusInternalServerError, "internal server error")
}
