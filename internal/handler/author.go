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
		h.sendHTTPError(w, errors.NewHTTPMethodError("Method not supported", nil))
	}
}

func (h *AuthorHandler) HandleAuthor(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, "authors/")
	authorID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		h.sendHTTPError(w, errors.NewValidationError("Invalid author ID", err))
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
		h.sendHTTPError(w, errors.NewHTTPMethodError("Method not supported", nil))
	}
}

func (h *AuthorHandler) getAuthors(w http.ResponseWriter, r *http.Request) {
	authors, err := h.service.GetAllAuthors()
	if err != nil {
		h.sendHTTPError(w, errors.MapErrorToHTTP(err))
		return
	}
	h.sendJSONResponse(w, http.StatusOK, authors)
}

func (h *AuthorHandler) getAuthorByID(w http.ResponseWriter, r *http.Request, authorID int) {
	author, err := h.service.GetAuthor(authorID)
	if err != nil {
		h.sendHTTPError(w, errors.NewNotFoundError("Author", authorID, err))
		return
	}
	h.sendJSONResponse(w, http.StatusOK, author)
}

func (h *AuthorHandler) createAuthor(w http.ResponseWriter, r *http.Request) {
	var author entity.Author
	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		h.sendHTTPError(w, errors.NewValidationError("Invalid request payload", err))
		return
	}

	authorID, err := h.service.CreateAuthor(*author.FirstName, *author.LastName, *author.Biography, author.BirthDate.Time)
	if err != nil {
		h.sendHTTPError(w, errors.MapErrorToHTTP(err))
		return
	}

	h.sendJSONResponse(w, http.StatusCreated, map[string]int{"author_id": authorID})
}

func (h *AuthorHandler) updateAuthor(w http.ResponseWriter, r *http.Request, authorID int) {
	var author entity.Author
	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		h.sendHTTPError(w, errors.NewValidationError("Invalid request payload", err))
		return
	}

	author.ID = authorID

	err := h.service.UpdateAuthor(author)
	if err != nil {
		h.sendHTTPError(w, errors.MapErrorToHTTP(err))
		return
	}

	h.sendJSONResponse(w, http.StatusOK, map[string]int{"author_id": authorID})
}

func (h *AuthorHandler) deleteAuthor(w http.ResponseWriter, r *http.Request, authorID int) {
	err := h.service.DeleteAuthor(authorID)
	if err != nil {
		h.sendHTTPError(w, errors.MapErrorToHTTP(err))
		return
	}

	h.sendJSONResponse(w, http.StatusOK, map[string]int{"author_id": authorID})
}

func (h *AuthorHandler) sendHTTPError(w http.ResponseWriter, httpErr *errors.AppError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpErr.Code)
	json.NewEncoder(w).Encode(errors.CreateErrorResponse(httpErr))
}

func (h *AuthorHandler) sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
