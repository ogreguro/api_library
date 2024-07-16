package handler

import (
	"api_library/internal/entity"
	"api_library/internal/errors"
	"api_library/internal/usecase"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type AuthorHandler struct {
	service usecase.Service
}

func NewAuthorHandler(service usecase.Service) *AuthorHandler {
	return &AuthorHandler{service: service}
}

func (h *AuthorHandler) HandleAuthors(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getAuthors(w, r)
	case http.MethodPost:
		h.createAuthor(w, r)
	default:
		code, msg := errors.MapErrorToHTTP(errors.NewMethodNotAllowedError())
		h.sendHTTPError(w, code, msg)
	}
}

func (h *AuthorHandler) HandleAuthor(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, "/")
	if len(urlPathSegments) != 3 {
		code, msg := errors.MapErrorToHTTP(errors.NewEndpointNotFoundError())
		h.sendHTTPError(w, code, msg)
		return
	}
	authorID, err := strconv.Atoi(urlPathSegments[2])
	if err != nil {
		code, msg := errors.MapErrorToHTTP(errors.NewInvalidIDError(err))
		h.sendHTTPError(w, code, msg)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getAuthorByID(w, r, authorID)
	case http.MethodPut:
		h.updateAuthor(w, r, authorID)
	case http.MethodDelete:
		h.deleteAuthor(w, r, authorID)
	default:
		code, msg := errors.MapErrorToHTTP(errors.NewMethodNotAllowedError())
		h.sendHTTPError(w, code, msg)
	}
}

func (h *AuthorHandler) getAuthors(w http.ResponseWriter, r *http.Request) {
	authors, err := h.service.GetAllAuthors()
	if err != nil {
		code, msg := errors.MapErrorToHTTP(err)
		h.sendHTTPError(w, code, msg)
		return
	}
	h.sendJSONResponse(w, http.StatusOK, authors)
}

func (h *AuthorHandler) getAuthorByID(w http.ResponseWriter, r *http.Request, authorID int) {
	author, err := h.service.GetAuthor(authorID)
	if err != nil {
		code, msg := errors.MapErrorToHTTP(err)
		h.sendHTTPError(w, code, msg)
		return
	}
	h.sendJSONResponse(w, http.StatusOK, author)
}

func (h *AuthorHandler) createAuthor(w http.ResponseWriter, r *http.Request) {
	var author entity.Author
	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		code, msg := errors.MapErrorToHTTP(errors.NewInvalidRequestPayloadError(err))
		h.sendHTTPError(w, code, msg)
		return
	}

	validationErrors := validateAuthorAttributes(author)
	if len(validationErrors) > 0 {
		code, msg := errors.MapErrorToHTTP(errors.NewValidationError(validationErrors))
		h.sendHTTPError(w, code, msg)
		return
	}

	authorID, err := h.service.CreateAuthor(*author.FirstName, *author.LastName, *author.Biography, author.BirthDate.Time)
	if err != nil {
		code, msg := errors.MapErrorToHTTP(err)
		h.sendHTTPError(w, code, msg)
		return
	}

	h.sendJSONResponse(w, http.StatusCreated, map[string]int{"author_id": authorID})
}

func (h *AuthorHandler) updateAuthor(w http.ResponseWriter, r *http.Request, authorID int) {
	var author entity.Author
	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		code, msg := errors.MapErrorToHTTP(errors.NewInvalidRequestPayloadError(err))
		h.sendHTTPError(w, code, msg)
		return
	}

	validationErrors := validateAuthorAttributes(author)
	if len(validationErrors) > 0 {
		code, msg := errors.MapErrorToHTTP(errors.NewValidationError(validationErrors))
		h.sendHTTPError(w, code, msg)
		return
	}

	author.ID = authorID

	err := h.service.UpdateAuthor(author)
	if err != nil {
		code, msg := errors.MapErrorToHTTP(err)
		h.sendHTTPError(w, code, msg)
		return
	}

	h.sendJSONResponse(w, http.StatusOK, map[string]int{"author_id": authorID})
}

func (h *AuthorHandler) deleteAuthor(w http.ResponseWriter, r *http.Request, authorID int) {
	err := h.service.DeleteAuthor(authorID)
	if err != nil {
		code, msg := errors.MapErrorToHTTP(err)
		h.sendHTTPError(w, code, msg)
		return
	}

	h.sendJSONResponse(w, http.StatusOK, map[string]int{"author_id": authorID})
}

func validateAuthorAttributes(author entity.Author) []string {
	var validationErrors []string

	if author.FirstName == nil || *author.FirstName == "" {
		validationErrors = append(validationErrors, "first_name is required")
	}
	if author.LastName == nil || *author.LastName == "" {
		validationErrors = append(validationErrors, "last_name is required")
	}
	if author.Biography == nil || *author.Biography == "" {
		validationErrors = append(validationErrors, "biography is required")
	}
	if author.BirthDate == nil {
		validationErrors = append(validationErrors, "birth_date is required")
	}

	return validationErrors
}

func (h *AuthorHandler) sendHTTPError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{"error": message, "code": code})
}

func (h *AuthorHandler) sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
