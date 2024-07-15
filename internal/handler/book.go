package handler

import (
	"api_library/internal/entity"
	"api_library/internal/usecase"
	"encoding/json"
	"fmt"
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
		h.sendResponse(w, http.StatusMethodNotAllowed, "Method not supported")
	}
}

func (h *BookHandler) HandleBook(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, "/")
	if len(urlPathSegments) < 3 {
		h.sendResponse(w, http.StatusBadRequest, "Invalid URL")
		return
	}

	id, err := strconv.Atoi(urlPathSegments[2])
	if err != nil {
		h.sendResponse(w, http.StatusBadRequest, "Invalid ID")
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
			h.sendResponse(w, http.StatusMethodNotAllowed, "Method not supported")
		}
	} else if len(urlPathSegments) == 5 && urlPathSegments[3] == "authors" {
		authorID, err := strconv.Atoi(urlPathSegments[4])
		if err != nil {
			h.sendResponse(w, http.StatusBadRequest, "Invalid author ID")
			return
		}
		h.handleBookAndAuthorUpdate(w, r, id, authorID)
	} else {
		h.sendResponse(w, http.StatusBadRequest, "Invalid URL")
	}
}

func (h *BookHandler) handleBookAndAuthorUpdate(w http.ResponseWriter, r *http.Request, bookID, authorID int) {
	var payload entity.BookAuthorPayload

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		h.sendResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	payload.Book.ID = bookID
	payload.Author.ID = authorID

	if err := h.service.UpdateBookWithAuthor(payload.Book, payload.Author); err != nil {
		h.sendResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error updating book and author: %v", err))
		return
	}

	h.sendResponse(w, http.StatusOK, "Book and author updated successfully")
}

func (h *BookHandler) getBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.service.GetAllBooks()
	if err != nil {
		h.sendResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error retrieving books: %v", err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func (h *BookHandler) getBookByID(w http.ResponseWriter, r *http.Request, bookID int) {
	book, err := h.service.GetBook(bookID)
	if err != nil {
		h.sendResponse(w, http.StatusNotFound, fmt.Sprintf("Book with ID %d not found", bookID))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) createBook(w http.ResponseWriter, r *http.Request) {
	var book entity.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		h.sendResponse(w, http.StatusBadRequest, fmt.Sprintf("Invalid data format: %v", err))
		return
	}

	bookID, err := h.service.CreateBook(*book.Title, *book.Year, *book.ISBN, *book.AuthorID)
	if err != nil {
		h.sendResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error creating book: %v", err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Book created with ID: %d", bookID)
}

func (h *BookHandler) updateBook(w http.ResponseWriter, r *http.Request, bookID int) {
	var book entity.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		h.sendResponse(w, http.StatusBadRequest, fmt.Sprintf("Invalid data format: %v", err))
		return
	}

	book.ID = bookID

	if err := h.service.UpdateBook(book); err != nil {
		h.sendResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error updating book: %v", err))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Book with ID %d updated successfully", bookID)
}

func (h *BookHandler) deleteBook(w http.ResponseWriter, r *http.Request, bookID int) {
	if err := h.service.DeleteBook(bookID); err != nil {
		h.sendResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting book: %v", err))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Book with ID %d deleted successfully", bookID)
}

func (h *BookHandler) sendResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}
