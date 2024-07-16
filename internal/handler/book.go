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

type BookHandler struct {
	service usecase.Service
}

func NewBookHandler(service usecase.Service) *BookHandler {
	return &BookHandler{service: service}
}

func (h *BookHandler) HandleBooks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getBooks(w, r)
	case http.MethodPost:
		h.createBook(w, r)
	default:
		h.sendHTTPError(w, errors.NewHTTPMethodError("Method not supported", nil))
	}
}

func (h *BookHandler) HandleBook(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, "/")
	if len(urlPathSegments) < 3 {
		h.sendHTTPError(w, errors.NewValidationError("Invalid URL", nil))
		return
	}

	id, err := strconv.Atoi(urlPathSegments[2])
	if err != nil {
		h.sendHTTPError(w, errors.NewValidationError("Invalid ID", err))
		return
	}

	if len(urlPathSegments) == 3 {
		switch r.Method {
		case http.MethodGet:
			h.getBookByID(w, r, id)
		case http.MethodPut:
			h.updateBook(w, r, id)
		case http.MethodDelete:
			h.deleteBook(w, r, id)
		default:
			h.sendHTTPError(w, errors.NewHTTPMethodError("Method not supported", nil))
		}
	} else if len(urlPathSegments) == 5 && urlPathSegments[3] == "authors" {
		authorID, err := strconv.Atoi(urlPathSegments[4])
		if err != nil {
			h.sendHTTPError(w, errors.NewValidationError("Invalid author ID", err))
			return
		}
		h.handleBookAndAuthorUpdate(w, r, id, authorID)
	} else {
		h.sendHTTPError(w, errors.NewValidationError("Invalid URL", nil))
	}
}

func (h *BookHandler) handleBookAndAuthorUpdate(w http.ResponseWriter, r *http.Request, bookID, authorID int) {
	var payload entity.BookAuthorPayload

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		h.sendHTTPError(w, errors.NewValidationError("Invalid request payload", err))
		return
	}

	payload.Book.ID = bookID
	payload.Author.ID = authorID

	if err := h.service.UpdateBookWithAuthor(payload.Book, payload.Author); err != nil {
		h.sendHTTPError(w, errors.MapErrorToHTTP(err))
		return
	}

	h.sendJSONResponse(w, http.StatusOK, "Book and author updated successfully")
}

func (h *BookHandler) getBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.service.GetAllBooks()
	if err != nil {
		h.sendHTTPError(w, errors.MapErrorToHTTP(err))
		return
	}
	h.sendJSONResponse(w, http.StatusOK, books)
}

func (h *BookHandler) getBookByID(w http.ResponseWriter, r *http.Request, bookID int) {
	book, err := h.service.GetBook(bookID)
	if err != nil {
		h.sendHTTPError(w, errors.NewNotFoundError("Book", bookID, err))
		return
	}
	h.sendJSONResponse(w, http.StatusOK, book)
}

func (h *BookHandler) createBook(w http.ResponseWriter, r *http.Request) {
	var book entity.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		h.sendHTTPError(w, errors.NewValidationError("Invalid data format", err))
		return
	}

	bookID, err := h.service.CreateBook(*book.Title, *book.Year, *book.ISBN, *book.AuthorID)
	if err != nil {
		h.sendHTTPError(w, errors.MapErrorToHTTP(err))
		return
	}

	h.sendJSONResponse(w, http.StatusCreated, map[string]int{"book_id": bookID})
}

func (h *BookHandler) updateBook(w http.ResponseWriter, r *http.Request, bookID int) {
	var book entity.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		h.sendHTTPError(w, errors.NewValidationError("Invalid data format", err))
		return
	}

	book.ID = bookID

	if err := h.service.UpdateBook(book); err != nil {
		h.sendHTTPError(w, errors.MapErrorToHTTP(err))
		return
	}

	h.sendJSONResponse(w, http.StatusOK, map[string]int{"book_id": bookID})
}

func (h *BookHandler) deleteBook(w http.ResponseWriter, r *http.Request, bookID int) {
	if err := h.service.DeleteBook(bookID); err != nil {
		h.sendHTTPError(w, errors.MapErrorToHTTP(err))
		return
	}

	h.sendJSONResponse(w, http.StatusOK, map[string]int{"book_id": bookID})
}

func (h *BookHandler) sendHTTPError(w http.ResponseWriter, httpErr *errors.AppError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpErr.Code)
	json.NewEncoder(w).Encode(errors.CreateErrorResponse(httpErr))
}

func (h *BookHandler) sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
