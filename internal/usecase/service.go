package usecase

import (
	"api_library/internal/entity"
	"api_library/internal/repository"
	"time"
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
	return s.repo.CreateAuthor(firstName, lastName, biography, entity.Date{Time: birthDate})
}

func (s *service) UpdateAuthor(author entity.Author) error {
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
	return s.repo.CreateBook(title, year, isbn, authorID)
}

func (s *service) UpdateBook(book entity.Book) error {
	return s.repo.UpdateBook(book)
}

func (s *service) DeleteBook(id int) error {
	return s.repo.DeleteBook(id)
}

func (s *service) GetBooksByAuthor(id int) ([]entity.Book, error) {
	return s.repo.GetBooksByAuthor(id)
}

func (s *service) UpdateBookWithAuthor(book entity.Book, author entity.Author) error {
	return s.repo.UpdateBookAndAuthor(book, author)
}
