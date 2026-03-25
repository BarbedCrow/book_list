-- +goose Up
CREATE TABLE lists (
    id       TEXT PRIMARY KEY,
    owner_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name     TEXT NOT NULL,
    type     TEXT NOT NULL DEFAULT 'custom'
);

-- +goose Down
DROP TABLE lists;
