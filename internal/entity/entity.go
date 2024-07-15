package entity

import (
	"fmt"
	"log"
	"time"
)

type Date struct {
	time.Time
}

// переопределяет поведение по умолчанию для десериализации
func (d *Date) UnmarshalJSON(b []byte) error {
	str := string(b)
	str = str[1 : len(str)-1] // убираем кавычки

	t, err := time.Parse("2006-01-02", str)
	if err != nil {
		log.Printf("Date parse error %q as \"2006-01-02\": %v", str, err)
		return fmt.Errorf("Date parse error %q as \"2006-01-02\": %v", str, err)
	}
	d.Time = t
	return nil
}

// переопределяет поведение по умолчанию для сериализации
func (d Date) MarshalJSON() (string, error) {
	return string(d.Time.Format("2006-01-02")), nil
}

type Author struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Biography string `json:"biography"`
	BirthDate Date   `json:"birth_date"`
}

type Book struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	AuthorID int    `json:"author_id"`
	Year     int    `json:"year"`
	ISBN     string `json:"isbn"`
}

type BookAuthorPayload struct {
	Book   Book   `json:"book"`
	Author Author `json:"author"`
}
