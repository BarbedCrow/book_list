# Book List

Backend service for tracking read books, searching a catalog, and managing personal reading lists.

## Stack

- Go 1.25+
- PostgreSQL 16
- Docker & Docker Compose

## Requirements

- [Go 1.25+](https://golang.org/dl/)
- [Docker](https://docs.docker.com/get-docker/)
- [goose](https://github.com/pressly/goose) (for local migrations)

## Quick Start (local)

```bash
# Start PostgreSQL
make docker-up

# Run migrations
make migrate-up

# Start the server
make run
```

The server starts on `http://localhost:8080`.

## Quick Start (Docker)

```bash
docker compose up --build
```

This starts PostgreSQL, runs migrations, and launches the app on `:8080`.

## Environment Variables

See [`.env.example`](.env.example) for all available variables:

| Variable       | Default                    | Description           |
|----------------|----------------------------|-----------------------|
| `SERVER_ADDR`  | `:8080`                    | Server listen address |
| `DB_HOST`      | `localhost`                | Database host         |
| `DB_PORT`      | `5432`                     | Database port         |
| `DB_USER`      | `postgres`                 | Database user         |
| `DB_PASSWORD`  | `postgres`                 | Database password     |
| `DB_NAME`      | `book_list`                | Database name         |
| `DB_SSLMODE`   | `disable`                  | SSL mode              |
| `JWT_SECRET`   | `change-me-in-production`  | JWT signing secret    |
| `JWT_TTL`      | `24h`                      | JWT token lifetime    |

## API Endpoints

| Method   | Path                             | Auth | Description              |
|----------|----------------------------------|------|--------------------------|
| `GET`    | `/health`                        | No   | Health check             |
| `POST`   | `/register`                      | No   | Register a new user      |
| `POST`   | `/login`                         | No   | Authenticate, get JWT    |
| `GET`    | `/books?title=...`               | No   | Search books by title    |
| `GET`    | `/books/{id}`                    | No   | Get book details         |
| `GET`    | `/authors?name=...`              | No   | Search authors by name   |
| `GET`    | `/authors/{id}`                  | No   | Get author details       |
| `GET`    | `/authors/{id}/books`            | No   | Get books by author      |
| `GET`    | `/lists`                         | Yes  | Get user's lists         |
| `POST`   | `/lists`                         | Yes  | Create a custom list     |
| `POST`   | `/lists/{id}/books`              | Yes  | Add book to list         |
| `DELETE`  | `/lists/{id}/books/{book_id}`   | Yes  | Remove book from list    |

Protected endpoints require `Authorization: Bearer <token>` header.

## Makefile Commands

| Command                | Description                        |
|------------------------|------------------------------------|
| `make run`             | Run the server locally             |
| `make test`            | Run unit tests                     |
| `make test-integration`| Run integration tests (needs DB)   |
| `make migrate-up`      | Apply database migrations          |
| `make migrate-down`    | Rollback last migration            |
| `make docker-up`       | Start PostgreSQL container         |
| `make docker-down`     | Stop PostgreSQL container          |

## Project Structure

```
cmd/server/             — Application entrypoint
internal/
  domain/               — Domain entities (Book, Author, User, List)
  usecase/
    book/               — Book search & details use cases
    author/             — Author search, details & books use cases
    user/               — Registration & authentication use cases
    list/               — List CRUD & book management use cases
  handler/              — HTTP handlers & middleware
  postgres/             — PostgreSQL repository implementations
  auth/                 — JWT token provider
  hasher/               — Bcrypt password hasher
migrations/             — SQL migration files (goose)
```
