package usecase

import (
	"api_library/internal/entity"
	"api_library/internal/errors"
	"api_library/internal/repository"
	"regexp"
	"time"
)

var (
	isbn10Regex = regexp.MustCompile(`^\d{9}[0-9X]$`)
	isbn13Regex = regexp.MustCompile(`^\d{13}$`)
)

type Service interface {
	GetAllAuthors() ([]entity.Author, error)
	GetAuthor(id int) (entity.Author, error)
	CreateAuthor(firstName, lastName, biography string, birthDate time.Time) (int, error)
	UpdateAuthor(author entity.Author) error
	DeleteAuthor(id int) error

	GetAllBooks() ([]entity.Book, error)
	GetBook(id int) (entity.Book, error)
	CreateBook(title string, year int, isbn string, authorID int) (int, error)
	UpdateBook(book entity.Book) error
	DeleteBook(id int) error
	GetBooksByAuthor(id int) ([]entity.Book, error)
	UpdateBookWithAuthor(book entity.Book, author entity.Author) error
}

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAllAuthors() ([]entity.Author, error) {
	return s.repo.GetAllAuthors()
}

func (s *service) GetAuthor(id int) (entity.Author, error) {
	return s.repo.GetAuthor(id)
}

func (s *service) CreateAuthor(firstName, lastName, biography string, birthDate time.Time) (int, error) {
	validationErrors := s.validateAuthor(entity.Author{FirstName: &firstName, LastName: &lastName, Biography: &biography, BirthDate: &entity.Date{Time: birthDate}})
	if len(validationErrors) > 0 {
		return 0, errors.NewValidationError(validationErrors)
	}
	return s.repo.CreateAuthor(firstName, lastName, biography, entity.Date{Time: birthDate})
}

func (s *service) UpdateAuthor(author entity.Author) error {
	validationErrors := s.validateAuthor(author)
	if len(validationErrors) > 0 {
		return errors.NewValidationError(validationErrors)
	}
	return s.repo.UpdateAuthor(author)
}

func (s *service) DeleteAuthor(id int) error {
	return s.repo.DeleteAuthor(id)
}

func (s *service) GetAllBooks() ([]entity.Book, error) {
	return s.repo.GetAllBooks()
}

func (s *service) GetBook(id int) (entity.Book, error) {
	return s.repo.GetBook(id)
}

func (s *service) CreateBook(title string, year int, isbn string, authorID int) (int, error) {
	validationErrors := s.validateBook(entity.Book{Title: &title, Year: &year, ISBN: &isbn, AuthorID: &authorID})
	if len(validationErrors) > 0 {
		return 0, errors.NewValidationError(validationErrors)
	}
	return s.repo.CreateBook(title, year, isbn, authorID)
}

func (s *service) UpdateBook(book entity.Book) error {
	validationErrors := s.validateBook(book)
	if len(validationErrors) > 0 {
		return errors.NewValidationError(validationErrors)
	}
	return s.repo.UpdateBook(book)
}

func (s *service) DeleteBook(id int) error {
	return s.repo.DeleteBook(id)
}

func (s *service) GetBooksByAuthor(id int) ([]entity.Book, error) {
	return s.repo.GetBooksByAuthor(id)
}

func (s *service) UpdateBookWithAuthor(book entity.Book, author entity.Author) error {
	bookValidationErrors := s.validateBook(book)
	authorValidationErrors := s.validateAuthor(author)

	if len(bookValidationErrors) > 0 && len(authorValidationErrors) > 0 {
		return errors.NewValidationError(append(bookValidationErrors, authorValidationErrors...))
	}

	return s.repo.UpdateBookAndAuthor(book, author)
}

func (s *service) validateBook(book entity.Book) []string {
	var validationErrors []string

	if book.Title != nil && *book.Title == "" {
		validationErrors = append(validationErrors, "title cannot be empty")
	}
	if book.Year != nil && *book.Year == 0 {
		validationErrors = append(validationErrors, "year cannot be empty")
	}
	if book.AuthorID != nil && *book.AuthorID == 0 {
		validationErrors = append(validationErrors, "author_id cannot be empty")
	}
	if book.ISBN != nil {
		strippedISBN := removeHyphens(*book.ISBN)
		if strippedISBN == "" {
			validationErrors = append(validationErrors, "ISBN cannot be empty")
		} else if !isbn10Regex.MatchString(strippedISBN) && !isbn13Regex.MatchString(strippedISBN) {
			validationErrors = append(validationErrors, "ISBN format is invalid")
		}
	}

	return validationErrors
}

func (s *service) validateAuthor(author entity.Author) []string {
	var validationErrors []string

	if author.FirstName != nil && *author.FirstName == "" {
		validationErrors = append(validationErrors, "first_name cannot be empty")
	}
	if author.LastName != nil && *author.LastName == "" {
		validationErrors = append(validationErrors, "last_name cannot be empty")
	}
	if author.Biography != nil && *author.Biography == "" {
		validationErrors = append(validationErrors, "biography cannot be empty")
	}
	if author.BirthDate != nil && author.BirthDate.Time.IsZero() {
		validationErrors = append(validationErrors, "birth_date cannot be empty")
	}

	return validationErrors
}

func removeHyphens(input string) string {
	return regexp.MustCompile(`-`).ReplaceAllString(input, "")
}
