.PHONY: run test migrate-up migrate-down docker-up docker-down

run:
	go run ./cmd/server

test:
	go test ./...

docker-up:
	docker compose up -d

docker-down:
	docker compose down

GOOSE := $(shell go env GOPATH)/bin/goose
DB_DSN := host=$${DB_HOST:-localhost} port=$${DB_PORT:-5432} user=$${DB_USER:-postgres} password=$${DB_PASSWORD:-postgres} dbname=$${DB_NAME:-book_list} sslmode=$${DB_SSLMODE:-disable}

migrate-up:
	$(GOOSE) -dir migrations postgres "$(DB_DSN)" up

migrate-down:
	$(GOOSE) -dir migrations postgres "$(DB_DSN)" down
