package middleware

import (
	"net/http"
	"time"

	"github.com/task-manager/task-service/pkg/logger"
)

// responseWriter wraps http.ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Logger logs each HTTP request with structured JSON output
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)

		requestID := GetRequestID(r.Context())
		logger.Request(requestID, r.Method, r.URL.Path, rw.statusCode, time.Since(start))
	})
}
