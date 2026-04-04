package errors

import "net/http"

type AppError struct {
	Code    int
	Message string
	Details string
}

func (e *AppError) Error() string {
	return e.Message
}

var (
	ErrInvalidCredentials = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "Invalid email or password",
	}

	ErrUserNotFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "User not found",
	}

	ErrUserAlreadyExists = &AppError{
		Code:    http.StatusConflict,
		Message: "User with this email already exists",
	}

	ErrInvalidToken = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "Invalid or expired token",
	}

	ErrTokenBlacklisted = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "Token has been revoked",
	}

	ErrMissingToken = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Authorization token is required",
	}

	ErrInvalidInput = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Invalid input data",
	}
)

func NewAppError(code int, message, details string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
	}
}
