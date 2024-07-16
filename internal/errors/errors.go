package errors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func NewValidationError(message string, err error) *AppError {
	return NewAppError(http.StatusBadRequest, fmt.Sprintf("%s: %v", message, err), err)
}

func NewNotFoundError(resource string, id int, err error) *AppError {
	return NewAppError(http.StatusNotFound, fmt.Sprintf("%s with ID %d not found", resource, id), err)
}

func NewDBError(message string, err error) *AppError {
	return NewAppError(http.StatusInternalServerError, message, err)
}

func NewHTTPMethodError(message string, err error) *AppError {
	return NewAppError(http.StatusMethodNotAllowed, message, err)
}

func MapErrorToHTTP(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return NewAppError(http.StatusInternalServerError, "Internal Server Error", err)
}

func CreateErrorResponse(appErr *AppError) map[string]interface{} {
	return map[string]interface{}{
		"code":    appErr.Code,
		"message": appErr.Message,
	}
}
