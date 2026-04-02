package middleware

import (
	"fmt"
	"net/http"

	"github.com/task-manager/task-service/pkg/logger"
	"github.com/task-manager/task-service/pkg/response"
)

// Recovery catches panics and returns a 500 error instead of crashing the server
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error(fmt.Sprintf("panic recovered: %v", err))
				response.Error(w, http.StatusInternalServerError, "internal server error")
			}
		}()

		next.ServeHTTP(w, r)
	})
}
