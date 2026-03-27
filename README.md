# Book List

Backend service for tracking read books, searching a catalog, and managing personal reading lists.

## Stack

- Go 1.25+
- PostgreSQL 16
- Prometheus + Alertmanager (metrics & alerts)
- Grafana (dashboards)
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

The monitoring stack (Prometheus, Alertmanager, Grafana) starts alongside the app:

| Service        | URL                    | Credentials   |
|----------------|------------------------|---------------|
| App            | http://localhost:8080   | —             |
| Prometheus     | http://localhost:9090   | —             |
| Alertmanager   | http://localhost:9093   | —             |
| Grafana        | http://localhost:3000   | admin / admin |

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
| `GET`    | `/health`                        | No   | Health check (incl. DB)  |
| `GET`    | `/metrics`                       | No   | Prometheus metrics       |
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

## Monitoring & Alerts

### Collected Metrics

- **HTTP**: `http_requests_total`, `http_request_duration_seconds`, `http_requests_in_flight`
- **DB pool**: `pgxpool_total_conns`, `pgxpool_acquired_conns`, `pgxpool_idle_conns`, `pgxpool_acquire_count_total`, `pgxpool_acquire_duration_seconds_total`, `pgxpool_empty_acquire_count_total`, `pgxpool_max_conns`
- **Runtime**: Go runtime and process metrics (via `prometheus/client_golang`)

### Alert Rules

| Alert                  | Condition                                       | Severity |
|------------------------|-------------------------------------------------|----------|
| ServiceDown            | Target unreachable for > 1 min                  | critical |
| HighErrorRate          | 5xx rate > 5% of all requests for 5 min         | critical |
| HighLatencyP95         | p95 latency > 500ms for 5 min                   | warning  |
| PgxPoolExhausted       | No idle connections for > 2 min                 | warning  |
| HighUnauthorizedRate   | 401 rate > 1 req/s for 5 min                    | warning  |

Alertmanager receiver is configured as a placeholder — update `monitoring/alertmanager/alertmanager.yml` with Slack, email, or PagerDuty settings.

### Grafana

A pre-provisioned **Book List — Overview** dashboard is available out of the box with panels for request rate, error rate, latency percentiles, in-flight requests, DB pool connections, and HTTP status codes.

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
  handler/              — HTTP handlers, middleware & metrics
  postgres/             — PostgreSQL repository implementations
  auth/                 — JWT token provider
  hasher/               — Bcrypt password hasher
  monitor/              — Prometheus collectors (pgxpool)
migrations/             — SQL migration files (goose)
monitoring/
  prometheus/           — Prometheus config & alert rules
  alertmanager/         — Alertmanager config
  grafana/              — Grafana provisioning & dashboards
```
