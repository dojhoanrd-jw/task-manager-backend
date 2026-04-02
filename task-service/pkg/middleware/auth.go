package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/task-manager/task-service/pkg/response"
)

type contextKey string

const (
	UserIDKey contextKey = "userId"
	EmailKey  contextKey = "email"
	RoleKey   contextKey = "role"
)

// Auth validates the JWT token and injects user data into the request context
func Auth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Error(w, http.StatusUnauthorized, "authorization header is required")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				response.Error(w, http.StatusUnauthorized, "invalid authorization format")
				return
			}

			token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				response.Error(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				response.Error(w, http.StatusUnauthorized, "invalid token claims")
				return
			}

			// Inject user data into context and headers for downstream services
			ctx := context.WithValue(r.Context(), UserIDKey, claims["userId"])
			ctx = context.WithValue(ctx, EmailKey, claims["email"])
			ctx = context.WithValue(ctx, RoleKey, claims["role"])

			r.Header.Set("X-User-ID", claims["userId"].(string))
			r.Header.Set("X-User-Email", claims["email"].(string))
			r.Header.Set("X-User-Role", claims["role"].(string))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
