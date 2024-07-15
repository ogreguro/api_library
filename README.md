# api_library

## Описание

Проект для управления коллекцией книг и авторами, написанный на Go с использованием PostgreSQL в качестве базы данных. Приложение контейнеризировано с помощью Docker и docker-compose.

## Требования

- Docker
- Docker Compose

## Установка

1. Клонируйте репозиторий:
    ```sh
    git clone https://github.com/yourusername/api_library.git
    cd api_library
    ```

2. Запустите контейнеры:
    ```sh
    docker-compose up --build
    ```

Приложение будет доступно по адресу `http://localhost:8888`, а PostgreSQL по `localhost:5433`.

## Makefile
```makefile
.PHONY: build run stop

build:
    docker-compose build

run:
    docker-compose up -d

stop:
   docker-compose down