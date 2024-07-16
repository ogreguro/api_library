package entity

import (
	"fmt"
	"time"
)

type Date struct {
	time.Time
}

// Переопределяет поведение по умолчанию для десериализации
func (d *Date) UnmarshalJSON(b []byte) error {
	str := string(b)
	str = str[1 : len(str)-1] // убираем кавычки

	t, err := time.Parse("2006-01-02", str)
	if err != nil {
		return fmt.Errorf("Date parse error %q as \"2006-01-02\": %v", str, err)
	}
	d.Time = t
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", d.Time.Format("2006-01-02"))), nil
}

type Author struct {
	ID        int     `json:"id"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Biography *string `json:"biography,omitempty"`
	BirthDate *Date   `json:"birth_date,omitempty"`
}

type Book struct {
	ID       int     `json:"id"`
	Title    *string `json:"title,omitempty"`
	AuthorID *int    `json:"author_id,omitempty"`
	Year     *int    `json:"year,omitempty"`
	ISBN     *string `json:"isbn,omitempty"`
}

type BookAuthorPayload struct {
	Book   Book   `json:"book"`
	Author Author `json:"author"`
}
