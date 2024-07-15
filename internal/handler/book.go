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
	bookService usecase.Service
}

func NewBookHandler(service usecase.Service) *BookHandler {
	return &BookHandler{bookService: service}
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
	urlPathSegments := strings.Split(r.URL.Path, "books/")
	bookID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		h.sendResponse(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getBookByID(w, r, bookID)
	case http.MethodPut:
		h.updateBook(w, r, bookID)
	case http.MethodDelete:
		h.deleteBook(w, r, bookID)
	default:
		h.sendResponse(w, http.StatusMethodNotAllowed, "Method not supported")
	}
}

func (h *BookHandler) getBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.bookService.GetAllBooks()
	if err != nil {
		h.sendResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error retrieving books: %v", err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func (h *BookHandler) getBookByID(w http.ResponseWriter, r *http.Request, bookID int) {
	book, err := h.bookService.GetBook(bookID)
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
		h.sendResponse(w, http.StatusBadRequest, "Invalid data format")
		return
	}

	bookID, err := h.bookService.CreateBook(book.Title, book.Year, book.ISBN, book.AuthorID)
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
		h.sendResponse(w, http.StatusBadRequest, "Invalid data format")
		return
	}

	err := h.bookService.UpdateBook(bookID, book.Title, book.Year, book.ISBN, book.AuthorID)
	if err != nil {
		h.sendResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error updating book: %v", err))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Book with ID %d updated")
}

func (h *BookHandler) deleteBook(w http.ResponseWriter, r *http.Request, bookID int) {
	err := h.bookService.DeleteBook(bookID)
	if err != nil {
		h.sendResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting book: %v", err))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Book with ID %d deleted")
}

func (h *BookHandler) sendResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}
