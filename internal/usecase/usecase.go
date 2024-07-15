package usecase

import (
	"excercise_library/internal/entity"
	"excercise_library/internal/repository"
	"time"
)

type Service interface {
	GetAllAuthors() ([]entity.Author, error)
	GetAuthor(id int) (entity.Author, error)
	CreateAuthor(firstName, lastName, biography string, birthDate time.Time) (int, error)
	UpdateAuthor(id int, firstName, lastName, biography string, birthDate time.Time) error
	DeleteAuthor(id int) error

	GetAllBooks() ([]entity.Book, error)
	GetBook(id int) (entity.Book, error)
	CreateBook(title string, year int, isbn string, authorID int) (int, error)
	UpdateBook(id int, title string, year int, isbn string, authorID int) error
	DeleteBook(id int) error
	UpdateBookWithAuthor(bookID int, newTitle string, newYear int, newISBN string, authorID int, newFirstName string, newLastName string, newBiography string, newBirthDate time.Time) error
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
	return s.repo.CreateAuthor(firstName, lastName, biography, birthDate)
}

func (s *service) UpdateAuthor(id int, firstName, lastName, biography string, birthDate time.Time) error {
	return s.repo.UpdateAuthor(id, firstName, lastName, biography, birthDate)
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
	return s.repo.CreateBook(title, year, isbn, authorID)
}

func (s *service) UpdateBook(id int, title string, year int, isbn string, authorID int) error {
	return s.repo.UpdateBook(id, title, year, isbn)
}

func (s *service) DeleteBook(id int) error {
	return s.repo.DeleteBook(id)
}

func (s *service) GetBooksByAuthor(id int) ([]entity.Book, error) {
	return s.repo.GetBooksByAuthor(id)
}

func (s *service) UpdateBookWithAuthor(bookID int, newTitle string, newYear int, newISBN string, authorID int, newFirstName string, newLastName string, newBiography string, newBirthDate time.Time) error {
	return s.repo.UpdateBookAndAuthor(bookID, newTitle, newYear, newISBN, authorID, newFirstName, newLastName, newBiography, newBirthDate)
}
