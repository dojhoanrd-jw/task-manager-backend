package apperror

import "fmt"

// AppError represents a structured application error
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

// New creates a new AppError
func New(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

// Wrap creates a new AppError wrapping an existing error
func Wrap(code int, message string, err error) *AppError {
	return &AppError{Code: code, Message: fmt.Sprintf("%s: %v", message, err)}
}

// Common errors
var (
	ErrNotFound       = New(404, "resource not found")
	ErrUnauthorized   = New(401, "unauthorized")
	ErrForbidden      = New(403, "insufficient permissions")
	ErrBadRequest     = New(400, "bad request")
	ErrAlreadyExists  = New(409, "resource already exists")
	ErrInternalServer = New(500, "internal server error")
)

// NotFound creates a not found error with custom message
func NotFound(resource string) *AppError {
	return New(404, fmt.Sprintf("%s not found", resource))
}

// BadRequest creates a bad request error with custom message
func BadRequest(message string) *AppError {
	return New(400, message)
}

// Forbidden creates a forbidden error with custom message
func Forbidden(message string) *AppError {
	return New(403, message)
}

// Conflict creates a conflict error with custom message
func Conflict(message string) *AppError {
	return New(409, message)
}
