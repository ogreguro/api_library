# Используем официальный образ Golang
FROM golang:1.21-alpine

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum для загрузки зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем все остальные файлы проекта
COPY . .

# Собираем приложение
RUN go build -o main ./cmd

# Определяем переменные окружения
ENV DB_HOST=db
ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_PASSWORD=postgres
ENV DB_NAME=library

# Запускаем приложение
CMD ["./main"]