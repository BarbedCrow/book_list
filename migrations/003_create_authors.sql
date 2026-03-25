-- +goose Up
CREATE TABLE authors (
    id   TEXT PRIMARY KEY,
    name TEXT NOT NULL
);

-- +goose Down
DROP TABLE authors;
