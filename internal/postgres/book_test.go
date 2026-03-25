//go:build integration

package postgres_test

import (
	"context"
	"errors"
	"testing"

	"github.com/BarbedCrow/book_list/internal/postgres"
	bookuc "github.com/BarbedCrow/book_list/internal/usecase/book"
)

func TestBookRepo_FindByTitle(t *testing.T) {
	pool := newTestPool(t)
	cleanTables(t, pool)
	ctx := context.Background()

	pool.Exec(ctx, `INSERT INTO books (id, title) VALUES ('b1', 'Go Programming')`)
	pool.Exec(ctx, `INSERT INTO authors (id, name) VALUES ('a1', 'John Doe')`)
	pool.Exec(ctx, `INSERT INTO book_authors (book_id, author_id) VALUES ('b1', 'a1')`)

	repo := postgres.NewBookRepo(pool)

	t.Run("found", func(t *testing.T) {
		books, err := repo.FindByTitle(ctx, "Go")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(books) != 1 {
			t.Fatalf("want 1 book, got %d", len(books))
		}
		if books[0].Title != "Go Programming" {
			t.Fatalf("want title 'Go Programming', got %s", books[0].Title)
		}
		if len(books[0].Authors) != 1 || books[0].Authors[0] != "John Doe" {
			t.Fatalf("want author 'John Doe', got %v", books[0].Authors)
		}
	})

	t.Run("not found", func(t *testing.T) {
		books, err := repo.FindByTitle(ctx, "nonexistent")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(books) != 0 {
			t.Fatalf("want 0 books, got %d", len(books))
		}
	})
}

func TestBookRepo_FindByID(t *testing.T) {
	pool := newTestPool(t)
	cleanTables(t, pool)
	ctx := context.Background()

	pool.Exec(ctx, `INSERT INTO books (id, title) VALUES ('b1', 'Go Programming')`)

	repo := postgres.NewBookRepo(pool)

	t.Run("found", func(t *testing.T) {
		book, err := repo.FindByID(ctx, "b1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if book.ID != "b1" {
			t.Fatalf("want id 'b1', got %s", book.ID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.FindByID(ctx, "nonexistent")
		if !errors.Is(err, bookuc.ErrBookNotFound) {
			t.Fatalf("want ErrBookNotFound, got %v", err)
		}
	})
}
