//go:build integration

package postgres_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func newTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		envOrDefault("DB_HOST", "localhost"),
		envOrDefault("DB_PORT", "5432"),
		envOrDefault("DB_USER", "postgres"),
		envOrDefault("DB_PASSWORD", "postgres"),
		envOrDefault("DB_NAME", "book_list"),
		envOrDefault("DB_SSLMODE", "disable"),
	)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("connect to test db: %v", err)
	}

	t.Cleanup(func() { pool.Close() })
	return pool
}

func cleanTables(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	tables := []string{"list_books", "lists", "book_authors", "books", "authors", "users"}
	for _, table := range tables {
		_, err := pool.Exec(context.Background(), "DELETE FROM "+table)
		if err != nil {
			t.Fatalf("clean table %s: %v", table, err)
		}
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
