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
		code, msg := errors.MapErrorToHTTP(errors.NewMethodNotAllowedError())
		h.sendHTTPError(w, code, msg)
	}
}

func (h *BookHandler) HandleBook(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, "/")

	if len(urlPathSegments) < 3 {
		code, msg := errors.MapErrorToHTTP(errors.NewEndpointNotFoundError())
		h.sendHTTPError(w, code, msg)
		return
	}

	id, err := strconv.Atoi(urlPathSegments[2])
	if err != nil {
		code, msg := errors.MapErrorToHTTP(errors.NewInvalidIDError(err))
		h.sendHTTPError(w, code, msg)
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
			code, msg := errors.MapErrorToHTTP(errors.NewMethodNotAllowedError())
			h.sendHTTPError(w, code, msg)
		}
	} else if len(urlPathSegments) == 5 && urlPathSegments[3] == "authors" {
		authorID, err := strconv.Atoi(urlPathSegments[4])
		if err != nil {
			code, msg := errors.MapErrorToHTTP(errors.NewInvalidIDError(err))
			h.sendHTTPError(w, code, msg)
			return
		}
		h.handleBookAndAuthorUpdate(w, r, id, authorID)
	} else {
		code, msg := errors.MapErrorToHTTP(errors.NewEndpointNotFoundError())
		h.sendHTTPError(w, code, msg)
	}
}

func (h *BookHandler) createBook(w http.ResponseWriter, r *http.Request) {
	var book entity.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		code, msg := errors.MapErrorToHTTP(errors.NewInvalidRequestPayloadError(err))
		h.sendHTTPError(w, code, msg)
		return
	}

	validationErrors := validateBookAttributes(book)
	if len(validationErrors) > 0 {
		code, msg := errors.MapErrorToHTTP(errors.NewValidationError(validationErrors))
		h.sendHTTPError(w, code, msg)
		return
	}

	bookID, err := h.service.CreateBook(*book.Title, *book.Year, *book.ISBN, *book.AuthorID)
	if err != nil {
		code, msg := errors.MapErrorToHTTP(err)
		h.sendHTTPError(w, code, msg)
		return
	}

	h.sendJSONResponse(w, http.StatusCreated, map[string]int{"book_id": bookID})
}

func (h *BookHandler) updateBook(w http.ResponseWriter, r *http.Request, bookID int) {
	var book entity.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		code, msg := errors.MapErrorToHTTP(errors.NewInvalidRequestPayloadError(err))
		h.sendHTTPError(w, code, msg)
		return
	}

	validationErrors := validateBookAttributes(book)
	if len(validationErrors) > 0 {
		code, msg := errors.MapErrorToHTTP(errors.NewValidationError(validationErrors))
		h.sendHTTPError(w, code, msg)
		return
	}

	book.ID = bookID

	if err := h.service.UpdateBook(book); err != nil {
		code, msg := errors.MapErrorToHTTP(err)
		h.sendHTTPError(w, code, msg)
		return
	}

	h.sendJSONResponse(w, http.StatusOK, map[string]int{"book_id": bookID})
}

func (h *BookHandler) handleBookAndAuthorUpdate(w http.ResponseWriter, r *http.Request, bookID, authorID int) {
	var payload entity.BookAuthorPayload

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		code, msg := errors.MapErrorToHTTP(errors.NewInvalidRequestPayloadError(err))
		h.sendHTTPError(w, code, msg)
		return
	}

	validationErrors := validateBookAndAuthorAttributes(payload.Book, payload.Author)
	if len(validationErrors) > 0 {
		code, msg := errors.MapErrorToHTTP(errors.NewValidationError(validationErrors))
		h.sendHTTPError(w, code, msg)
		return
	}

	payload.Book.ID = bookID
	payload.Author.ID = authorID

	if err := h.service.UpdateBookWithAuthor(payload.Book, payload.Author); err != nil {
		code, msg := errors.MapErrorToHTTP(err)
		h.sendHTTPError(w, code, msg)
		return
	}

	h.sendJSONResponse(w, http.StatusOK, map[string]string{"message": "Book and author updated successfully"})
}

func (h *BookHandler) getBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.service.GetAllBooks()
	if err != nil {
		code, msg := errors.MapErrorToHTTP(err)
		h.sendHTTPError(w, code, msg)
		return
	}
	h.sendJSONResponse(w, http.StatusOK, books)
}

func (h *BookHandler) getBookByID(w http.ResponseWriter, r *http.Request, bookID int) {
	book, err := h.service.GetBook(bookID)
	if err != nil {
		code, msg := errors.MapErrorToHTTP(err)
		h.sendHTTPError(w, code, msg)
		return
	}
	h.sendJSONResponse(w, http.StatusOK, book)
}

func (h *BookHandler) deleteBook(w http.ResponseWriter, r *http.Request, bookID int) {
	if err := h.service.DeleteBook(bookID); err != nil {
		code, msg := errors.MapErrorToHTTP(err)
		h.sendHTTPError(w, code, msg)
		return
	}

	h.sendJSONResponse(w, http.StatusOK, map[string]int{"book_id": bookID})
}

func validateBookAttributes(book entity.Book) []string {
	var validationErrors []string

	if book.Title == nil || *book.Title == "" {
		validationErrors = append(validationErrors, "title is required")
	}
	if book.Year == nil {
		validationErrors = append(validationErrors, "year is required")
	}
	if book.AuthorID == nil {
		validationErrors = append(validationErrors, "author_id is required")
	}
	if book.ISBN == nil || *book.ISBN == "" {
		validationErrors = append(validationErrors, "ISBN is required")
	}

	return validationErrors
}

func validateBookAndAuthorAttributes(book entity.Book, author entity.Author) []string {
	var validationErrors []string

	if book.Title == nil && book.Year == nil && book.AuthorID == nil && book.ISBN == nil {
		validationErrors = append(validationErrors, "at least one book attribute is required")
	}

	if author.FirstName == nil && author.LastName == nil && author.Biography == nil && author.BirthDate == nil {
		validationErrors = append(validationErrors, "at least one author attribute is required")
	}

	return validationErrors
}

func (h *BookHandler) sendHTTPError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{"error": message, "code": code})
}

func (h *BookHandler) sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
