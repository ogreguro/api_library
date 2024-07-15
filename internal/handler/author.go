package handler

import (
	"api_library/internal/entity"
	"api_library/internal/errors"
	"api_library/internal/usecase"
	"encoding/json"
	"fmt"
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
		h.sendHTTPError(w, errors.NewHTTPError(http.StatusMethodNotAllowed, "method not supported", "HandleAuthors"))
	}
}

func (h *AuthorHandler) HandleAuthor(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, "authors/")
	authorID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		h.sendHTTPError(w, errors.NewHTTPError(http.StatusBadRequest, "invalid author ID", "HandleAuthor"))
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
		h.sendHTTPError(w, errors.NewHTTPError(http.StatusMethodNotAllowed, "method not supported", "HandleAuthor"))
	}
}

func (h *AuthorHandler) getAuthors(w http.ResponseWriter, r *http.Request) {
	authors, err := h.service.GetAllAuthors()
	if err != nil {
		h.sendHTTPError(w, errors.MapErrorToHTTP(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authors)
	w.WriteHeader(http.StatusOK)
}

func (h *AuthorHandler) getAuthorByID(w http.ResponseWriter, r *http.Request, authorID int) {
	author, err := h.service.GetAuthor(authorID)
	if err != nil {
		h.sendHTTPError(w, errors.MapErrorToHTTP(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(author)
	w.WriteHeader(http.StatusOK)
}

func (h *AuthorHandler) createAuthor(w http.ResponseWriter, r *http.Request) {
	var author entity.Author
	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		h.sendHTTPError(w, errors.NewHTTPError(http.StatusBadRequest, err.Error(), "createAuthor"))
		return
	}

	authorID, err := h.service.CreateAuthor(*author.FirstName, *author.LastName, *author.Biography, author.BirthDate.Time)
	if err != nil {
		h.sendHTTPError(w, errors.MapErrorToHTTP(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Author created with ID: %d", authorID)
}

func (h *AuthorHandler) updateAuthor(w http.ResponseWriter, r *http.Request, authorID int) {
	var author entity.Author
	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		h.sendHTTPError(w, errors.NewHTTPError(http.StatusBadRequest, err.Error(), "updateAuthor"))
		return
	}

	author.ID = authorID

	err := h.service.UpdateAuthor(author)
	if err != nil {
		h.sendHTTPError(w, errors.MapErrorToHTTP(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Author updated with ID: %d", authorID)
}

func (h *AuthorHandler) deleteAuthor(w http.ResponseWriter, r *http.Request, authorID int) {
	err := h.service.DeleteAuthor(authorID)
	if err != nil {
		h.sendHTTPError(w, errors.MapErrorToHTTP(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Author deleted with ID: %d", authorID)
}

func (h *AuthorHandler) sendHTTPError(w http.ResponseWriter, httpErr *errors.HTTPError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpErr.Code)
	json.NewEncoder(w).Encode(httpErr)
}
