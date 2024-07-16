package errors

import (
	"fmt"
	"net/http"
)

const (
	ErrInvalidRequestPayload = "Invalid request payload"
	ErrInvalidDataFormat     = "Invalid data format"
	ErrInvalidID             = "Invalid ID"
	ErrMethodNotAllowed      = "Method not allowed"
	ErrEndpointNotFound      = "Endpoint not found"
	ErrInternalServerError   = "Internal Server Error"
	ErrValidationFailed      = "Validation failed"
	ErrDBError               = "Database error"
	ErrResourceNotFound      = "Resource not found"
)

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func NewInvalidRequestPayloadError(err error) *AppError {
	return NewAppError(http.StatusBadRequest, ErrInvalidRequestPayload, err)
}

func NewInvalidDataFormatError(err error) *AppError {
	return NewAppError(http.StatusBadRequest, ErrInvalidDataFormat, err)
}

func NewInvalidIDError(err error) *AppError {
	return NewAppError(http.StatusBadRequest, ErrInvalidID, err)
}

func NewMethodNotAllowedError() *AppError {
	return NewAppError(http.StatusMethodNotAllowed, ErrMethodNotAllowed, nil)
}

func NewEndpointNotFoundError() *AppError {
	return NewAppError(http.StatusNotFound, ErrEndpointNotFound, nil)
}

func NewInternalServerError(err error) *AppError {
	return NewAppError(http.StatusInternalServerError, ErrInternalServerError, err)
}

func NewValidationError(errors []string) *AppError {
	return NewAppError(http.StatusBadRequest, fmt.Sprintf("%s: %v", ErrValidationFailed, errors), nil)
}

func NewDBError(err error) *AppError {
	return NewAppError(http.StatusInternalServerError, ErrDBError, err)
}

func NewResourceNotFoundError(resource string, id *int) *AppError {
	if id != nil {
		return NewAppError(http.StatusNotFound, fmt.Sprintf("%s with ID %d not found", resource, *id), nil)
	}
	return NewAppError(http.StatusNotFound, fmt.Sprintf("%s not found", resource), nil)
}

func MapErrorToHTTP(err error) (int, string) {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code, appErr.Error()
	}
	return http.StatusInternalServerError, ErrInternalServerError
}
