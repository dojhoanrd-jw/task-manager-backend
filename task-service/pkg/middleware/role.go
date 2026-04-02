package middleware

import (
	"net/http"

	"github.com/task-manager/task-service/pkg/models"
	"github.com/task-manager/task-service/pkg/response"
)

// RequireRole restricts access to users with the specified roles
func RequireRole(allowedRoles ...models.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole := r.Header.Get("X-User-Role")
			if userRole == "" {
				response.Error(w, http.StatusForbidden, "access denied")
				return
			}

			for _, role := range allowedRoles {
				if models.Role(userRole) == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			response.Error(w, http.StatusForbidden, "insufficient permissions")
		})
	}
}
