package repository

import (
	"api_library/internal/entity"
	"api_library/internal/errors"
	"database/sql"
	"net/http"
	"strconv"
	"strings"
)

// Объявление интерфейса Repository
type Repository interface {
	GetAllAuthors() ([]entity.Author, error)
	GetAuthor(authorID int) (entity.Author, error)
	CreateAuthor(firstName, lastName, biography string, birthDate entity.Date) (int, error)
	UpdateAuthor(author entity.Author) error
	DeleteAuthor(authorID int) error
	GetAllBooks() ([]entity.Book, error)
	GetBooksByAuthor(authorID int) ([]entity.Book, error)
	GetBook(bookID int) (entity.Book, error)
	CreateBook(title string, year int, isbn string, authorID int) (int, error)
	UpdateBook(book entity.Book) error
	DeleteBook(bookID int) error
	UpdateBookAndAuthor(book entity.Book, author entity.Author) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
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
		var birthDate sql.NullTime
		if err := rows.Scan(&author.ID, &author.FirstName, &author.LastName, &author.Biography, &birthDate); err != nil {
			return nil, errors.MapErrorToHTTP(err)
		}
		if birthDate.Valid {
			author.BirthDate = &entity.Date{Time: birthDate.Time}
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
	var birthDate sql.NullTime
	err := r.db.QueryRow("SELECT id, first_name, last_name, biography, birth_date FROM authors WHERE id = $1", authorID).Scan(&author.ID, &author.FirstName, &author.LastName, &author.Biography, &birthDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return author, errors.ErrNotFound
		}
		return author, errors.ErrDB
	}
	if birthDate.Valid {
		author.BirthDate = &entity.Date{Time: birthDate.Time}
	}
	return author, nil
}

func (r *repository) CreateAuthor(firstName, lastName, biography string, birthDate entity.Date) (int, error) {
	var authorID int
	err := r.db.QueryRow("INSERT INTO authors (first_name, last_name, biography, birth_date) VALUES ($1, $2, $3, $4) RETURNING id", firstName, lastName, biography, birthDate.Time).Scan(&authorID)
	if err != nil {
		return 0, errors.MapErrorToHTTP(err)
	}
	return authorID, nil
}

func (r *repository) UpdateAuthor(author entity.Author) error {
	authorFields := extractAuthorFields(author)

	if len(authorFields) == 0 {
		return nil // Нет полей для обновления
	}

	query, args := createUpdateQuery("authors", authorFields)
	query += " WHERE id = $" + strconv.Itoa(len(args)+1)
	args = append(args, author.ID)

	result, err := r.db.Exec(query, args...)
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
		return errors.NewHTTPError(http.StatusNotFound, "author not found", "DeleteAuthor")
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
	return nil
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

func (r *repository) UpdateBook(book entity.Book) error {
	bookFields := extractBookFields(book)

	if len(bookFields) == 0 {
		return nil // Нет полей для обновления
	}

	query, args := createUpdateQuery("books", bookFields)
	query += " WHERE id = $" + strconv.Itoa(len(args)+1)
	args = append(args, book.ID)

	result, err := r.db.Exec(query, args...)
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

func (r *repository) UpdateBookAndAuthor(book entity.Book, author entity.Author) error {
	tx, err := r.db.Begin()
	if err != nil {
		return errors.MapErrorToHTTP(err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// Обновление книги
	bookFields := extractBookFields(book)

	if len(bookFields) > 0 {
		bookQuery, bookArgs := createUpdateQuery("books", bookFields)
		bookQuery += " WHERE id = $" + strconv.Itoa(len(bookArgs)+1)
		bookArgs = append(bookArgs, book.ID)

		if _, err = tx.Exec(bookQuery, bookArgs...); err != nil {
			return errors.MapErrorToHTTP(err)
		}
	}

	// Обновление автора
	authorFields := extractAuthorFields(author)

	if len(authorFields) > 0 {
		authorQuery, authorArgs := createUpdateQuery("authors", authorFields)
		authorQuery += " WHERE id = $" + strconv.Itoa(len(authorArgs)+1)
		authorArgs = append(authorArgs, author.ID)

		if _, err = tx.Exec(authorQuery, authorArgs...); err != nil {
			return errors.MapErrorToHTTP(err)
		}
	}

	return nil
}

// Приватная вспомогательная функция для создания SQL-запроса обновления
func createUpdateQuery(tableName string, fields map[string]interface{}) (string, []interface{}) {
	setClauses := []string{}
	args := []interface{}{}
	argIdx := 1

	for field, value := range fields {
		setClauses = append(setClauses, field+" = $"+strconv.Itoa(argIdx))
		args = append(args, value)
		argIdx++
	}

	query := "UPDATE " + tableName + " SET " + strings.Join(setClauses, ", ")

	return query, args
}

// Приватная вспомогательная функция для извлечения полей книги
func extractBookFields(book entity.Book) map[string]interface{} {
	fields := map[string]interface{}{}
	if book.Title != nil {
		fields["title"] = *book.Title
	}
	if book.AuthorID != nil {
		fields["author_id"] = *book.AuthorID
	}
	if book.Year != nil {
		fields["year"] = *book.Year
	}
	if book.ISBN != nil {
		fields["isbn"] = *book.ISBN
	}
	return fields
}

// Приватная вспомогательная функция для извлечения полей автора
func extractAuthorFields(author entity.Author) map[string]interface{} {
	fields := map[string]interface{}{}
	if author.FirstName != nil {
		fields["first_name"] = *author.FirstName
	}
	if author.LastName != nil {
		fields["last_name"] = *author.LastName
	}
	if author.Biography != nil {
		fields["biography"] = *author.Biography
	}
	if author.BirthDate != nil {
		fields["birth_date"] = author.BirthDate.Time.Format("2006-01-02")
	}
	return fields
}
