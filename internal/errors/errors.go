package errors

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrDB           = errors.New("database error")
)

type HTTPError struct {
	Code    int
	Message string
	Source  string
}

func NewHTTPError(code int, message string, source string) *HTTPError {
	return &HTTPError{Code: code, Message: message, Source: source}
}

func (e *HTTPError) Error() string {
	return e.Message + " " + e.Source
}

func MapErrorToHTTP(err error) *HTTPError {
	switch err {
	case ErrNotFound:
		return NewHTTPError(http.StatusNotFound, err.Error(), "")
	case ErrInvalidInput:
		return NewHTTPError(http.StatusBadRequest, err.Error(), "")
	case ErrDB:
		return NewHTTPError(http.StatusInternalServerError, err.Error(), "")
	default:
		return NewHTTPError(http.StatusInternalServerError, err.Error(), "")
	}
}
