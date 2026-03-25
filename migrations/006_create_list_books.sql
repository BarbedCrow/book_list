-- +goose Up
CREATE TABLE list_books (
    list_id TEXT NOT NULL REFERENCES lists(id) ON DELETE CASCADE,
    book_id TEXT NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    PRIMARY KEY (list_id, book_id)
);

-- +goose Down
DROP TABLE list_books;
