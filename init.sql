-- Переключение на контекст базы данных "library"
\c library;

CREATE TABLE authors (
                         id SERIAL PRIMARY KEY,
                         first_name VARCHAR(100),
                         last_name VARCHAR(100),
                         biography TEXT,
                         birth_date DATE
);

CREATE TABLE books (
                       id SERIAL PRIMARY KEY,
                       title VARCHAR(255),
                       author_id INT,
                       year INT,
                       isbn VARCHAR(13)
);