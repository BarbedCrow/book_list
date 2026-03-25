-- +goose Up
CREATE TABLE book_authors (
    book_id   TEXT NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    author_id TEXT NOT NULL REFERENCES authors(id) ON DELETE CASCADE,
    PRIMARY KEY (book_id, author_id)
);

-- +goose Down
DROP TABLE book_authors;
