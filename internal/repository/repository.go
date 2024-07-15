package repository

import (
	"database/sql"
	"excercise_library/internal/entity"
	"excercise_library/internal/errors"
	"log"
	"net/http"
	"time"
)

type Repository interface {
	GetAllAuthors() ([]entity.Author, error)
	GetAuthor(authorID int) (entity.Author, error)
	CreateAuthor(firstName, lastName, biography string, birthDate time.Time) (int, error)
	UpdateAuthor(authorID int, firstName, lastName, biography string, birthDate time.Time) error
	DeleteAuthor(authorID int) error
	GetAllBooks() ([]entity.Book, error)
	GetBooksByAuthor(authorID int) ([]entity.Book, error)
	GetBook(bookID int) (entity.Book, error)
	CreateBook(title string, year int, isbn string, authorID int) (int, error)
	UpdateBook(bookID int, title string, year int, isbn string) error
	DeleteBook(bookID int) error
	UpdateBookAndAuthor(bookID int, newTitle string, newYear int, newISBN string, authorID int, newFirstName string, newLastName string, newBiography string, newBirthDate time.Time) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetAllAuthors() ([]entity.Author, error) {
	rows, err := r.db.Query("SELECT id, first_name, last_name, biography, birth_date FROM authors")
	if err != nil {
		return nil, errors.MapErrorToHTTP(err)
	}
	defer rows.Close()

	var authors []entity.Author
	for rows.Next() {
		var author entity.Author
		if err := rows.Scan(&author.ID, &author.FirstName, &author.LastName, &author.Biography, &author.BirthDate); err != nil {
			return nil, errors.MapErrorToHTTP(err)
		}
		authors = append(authors, author)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.MapErrorToHTTP(err)
	}
	return authors, nil
}

func (r *repository) GetAuthor(authorID int) (entity.Author, error) {
	var author entity.Author
	err := r.db.QueryRow("SELECT id, first_name, last_name, biography, birth_date FROM authors WHERE id = $1", authorID).Scan(&author.ID, &author.FirstName, &author.LastName, &author.Biography, &author.BirthDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return author, errors.ErrNotFound
		}
		return author, errors.ErrDB
	}
	return author, nil
}

func (r *repository) CreateAuthor(firstName, lastName, biography string, birthDate time.Time) (int, error) {
	var authorID int
	err := r.db.QueryRow("INSERT INTO authors (first_name, last_name, biography, birth_date) VALUES ($1, $2, $3, $4) RETURNING id", firstName, lastName, biography, birthDate).Scan(&authorID)
	log.Printf("error=%+v", err)
	if err != nil {
		return 0, errors.MapErrorToHTTP(err)
	}
	return authorID, nil
}

func (r *repository) UpdateAuthor(authorID int, firstName, lastName, biography string, birthDate time.Time) error {
	result, err := r.db.Exec("UPDATE authors SET first_name = $1, last_name = $2, biography = $3, birth_date = $4 WHERE id = $5", firstName, lastName, biography, birthDate, authorID)
	rows, _ := result.RowsAffected()
	if err != nil {
		return errors.MapErrorToHTTP(err)
	} else if rows == 0 {
		return errors.NewHTTPError(http.StatusNotFound, "author not found", "UpdateAuthor")
	}

	return nil
}

func (r *repository) DeleteAuthor(authorID int) error {
	result, err := r.db.Exec("DELETE FROM authors WHERE id = $1", authorID)
	rows, _ := result.RowsAffected()
	if err != nil {
		return errors.MapErrorToHTTP(err)
	} else if rows == 0 {
		return errors.NewHTTPError(http.StatusNotFound, "author not found", "UpdateAuthor")
	}

	books, err := r.GetBooksByAuthor(authorID)
	if err != nil {
		return errors.MapErrorToHTTP(err)
	}
	for _, book := range books {
		if err := r.DeleteBook(book.ID); err != nil {
			return errors.MapErrorToHTTP(err)
		}
	}
	return err
}

func (r *repository) GetAllBooks() ([]entity.Book, error) {
	rows, err := r.db.Query("SELECT id, title, year, isbn, author_id FROM books")
	if err != nil {
		return nil, errors.MapErrorToHTTP(err)
	}
	defer rows.Close()

	var books []entity.Book
	for rows.Next() {
		var book entity.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Year, &book.ISBN, &book.AuthorID); err != nil {
			return nil, errors.MapErrorToHTTP(err)
		}
		books = append(books, book)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.MapErrorToHTTP(err)
	}
	return books, nil
}

func (r *repository) GetBooksByAuthor(authorID int) ([]entity.Book, error) {
	rows, err := r.db.Query("SELECT id, title, year, isbn, author_id FROM books WHERE author_id = $1", authorID)
	if err != nil {
		return nil, errors.MapErrorToHTTP(err)
	}
	defer rows.Close()

	var books []entity.Book
	for rows.Next() {
		var book entity.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Year, &book.ISBN, &book.AuthorID); err != nil {
			return nil, errors.MapErrorToHTTP(err)
		}
		books = append(books, book)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.MapErrorToHTTP(err)
	}
	return books, nil
}

func (r *repository) GetBook(bookID int) (entity.Book, error) {
	var book entity.Book
	err := r.db.QueryRow("SELECT id, title, year, isbn, author_id FROM books WHERE id = $1", bookID).Scan(&book.ID, &book.Title, &book.Year, &book.ISBN, &book.AuthorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return book, errors.ErrNotFound
		}
		return book, errors.ErrDB
	}
	return book, nil
}

func (r *repository) CreateBook(title string, year int, isbn string, authorID int) (int, error) {
	var bookID int
	err := r.db.QueryRow("INSERT INTO books (title, year, isbn, author_id) VALUES ($1, $2, $3, $4) RETURNING id", title, year, isbn, authorID).Scan(&bookID)
	if err != nil {
		return 0, errors.MapErrorToHTTP(err)
	}
	return bookID, nil
}

func (r *repository) UpdateBook(bookID int, title string, year int, isbn string) error {
	result, err := r.db.Exec("UPDATE books SET title = $1, year = $2, isbn = $3 WHERE id = $4", title, year, isbn, bookID)
	rows, _ := result.RowsAffected()
	if err != nil {
		return errors.MapErrorToHTTP(err)
	} else if rows == 0 {
		return errors.NewHTTPError(http.StatusNotFound, "book not found", "UpdateBook")
	}
	return nil
}

func (r *repository) DeleteBook(bookID int) error {
	result, err := r.db.Exec("DELETE FROM books WHERE id = $1", bookID)
	rows, _ := result.RowsAffected()
	if err != nil {
		return errors.MapErrorToHTTP(err)
	} else if rows == 0 {
		return errors.NewHTTPError(http.StatusNotFound, "book not found", "DeleteBook")
	}
	return nil
}

func (r *repository) UpdateBookAndAuthor(bookID int, newTitle string, newYear int, newISBN string, authorID int, newFirstName string, newLastName string, newBiography string, newBirthDate time.Time) error {
	tx, err := r.db.Begin()
	if err != nil {
		return errors.MapErrorToHTTP(err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			errors.MapErrorToHTTP(err)
		} else {
			err = tx.Commit()
		}
	}()

	var result sql.Result
	var rows int64

	result, err = tx.Exec("UPDATE books SET title = $1, year = $2, isbn = $3 WHERE id = $4", newTitle, newYear, newISBN, bookID)
	rows, _ = result.RowsAffected()
	if err != nil {
		return errors.MapErrorToHTTP(err)
	} else if rows == 0 {
		return errors.NewHTTPError(http.StatusNotFound, "book not found", "UpdateBookAndAuthor")
	}

	result, err = tx.Exec("UPDATE authors SET first_name = $1, last_name = $2, biography = $3, birth_date = $4 WHERE id = $5", newFirstName, newLastName, newBiography, newBirthDate.Format("2006-01-02"), authorID)
	rows, _ = result.RowsAffected()
	if err != nil {
		return errors.MapErrorToHTTP(err)
	} else if rows == 0 {
		return errors.NewHTTPError(http.StatusNotFound, "author not found", "UpdateBookAndAuthor")
	}

	return nil
}
