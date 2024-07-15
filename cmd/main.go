package main

import (
	"api_library/internal/handler"
	"api_library/internal/repository"
	"api_library/internal/usecase"
	"log"
	"net/http"
)

func main() {
	db, err := repository.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Инициализация репозитория
	repo := repository.NewRepository(db)

	// Инициализация сервисa
	service := usecase.NewService(repo)

	// Инициализация обработчиков
	authorHandler := handler.NewAuthorHandler(service)
	bookHandler := handler.NewBookHandler(service)

	// Маршруты
	http.HandleFunc("/authors", authorHandler.HandleAuthors)
	http.HandleFunc("/authors/", authorHandler.HandleAuthor)
	http.HandleFunc("/books", bookHandler.HandleBooks)
	http.HandleFunc("/books/", bookHandler.HandleBook)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
