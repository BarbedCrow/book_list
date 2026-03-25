-- +goose Up
CREATE TABLE books (
    id    TEXT PRIMARY KEY,
    title TEXT NOT NULL
);

-- +goose Down
DROP TABLE books;
