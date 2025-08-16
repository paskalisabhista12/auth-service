package exception

import "net/http"

type AppError struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	Detail     string `json:"detail,omitempty"`
}

var (
	ErrBadRequest   = &AppError{StatusCode: 400, Code: "BAD_REQUEST", Message: "Invalid request"}
	ErrUnauthorized = &AppError{StatusCode: 401, Code: "UNAUTHORIZED", Message: "Authentication required"}
	ErrForbidden    = &AppError{StatusCode: 403, Code: "FORBIDDEN", Message: "Access denied"}
	ErrNotFound     = &AppError{StatusCode: 404, Code: "NOT_FOUND", Message: "Resource not found"}
	ErrConflict     = &AppError{StatusCode: 409, Code: "CONFLICT", Message: "Resource conflict"}
	ErrInternal     = &AppError{StatusCode: 500, Code: "INTERNAL_ERROR", Message: "Internal server error"}
)

// Implement Go's error interface
func (e *AppError) Error() string {
	return e.Message
}

func NewBadRequest(msg string) *AppError {
	return &AppError{StatusCode: http.StatusBadRequest, Code: "BAD_REQUEST", Message: msg}
}

func NewNotFound(msg string) *AppError {
	return &AppError{StatusCode: http.StatusNotFound, Code: "NOT_FOUND", Message: msg}
}

func NewInternal(msg string) *AppError {
	return &AppError{StatusCode: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: msg}
}

func NewConflictBusinessException(msg string) *AppError {
	return &AppError{StatusCode: http.StatusConflict, Code: "CONFLICT", Message: msg}
}

func NewUnauthorizedBusinessException(msg string) *AppError {
	return &AppError{StatusCode: http.StatusUnauthorized, Code: "UNAUTHORIZED", Message: msg}
}
