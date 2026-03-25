//go:build integration

package postgres_test

import (
	"context"
	"errors"
	"testing"

	"github.com/BarbedCrow/book_list/internal/postgres"
	authoruc "github.com/BarbedCrow/book_list/internal/usecase/author"
)

func TestAuthorRepo_FindByID(t *testing.T) {
	pool := newTestPool(t)
	cleanTables(t, pool)
	ctx := context.Background()

	pool.Exec(ctx, `INSERT INTO authors (id, name) VALUES ('a1', 'J.R.R. Tolkien')`)

	repo := postgres.NewAuthorRepo(pool)

	t.Run("found", func(t *testing.T) {
		a, err := repo.FindByID(ctx, "a1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if a.Name != "J.R.R. Tolkien" {
			t.Fatalf("want name 'J.R.R. Tolkien', got %s", a.Name)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.FindByID(ctx, "nonexistent")
		if !errors.Is(err, authoruc.ErrAuthorNotFound) {
			t.Fatalf("want ErrAuthorNotFound, got %v", err)
		}
	})
}

func TestAuthorRepo_FindByName(t *testing.T) {
	pool := newTestPool(t)
	cleanTables(t, pool)
	ctx := context.Background()

	pool.Exec(ctx, `INSERT INTO authors (id, name) VALUES ('a1', 'J.R.R. Tolkien')`)

	repo := postgres.NewAuthorRepo(pool)

	authors, err := repo.FindByName(ctx, "Tolkien")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(authors) != 1 {
		t.Fatalf("want 1 author, got %d", len(authors))
	}
}

func TestAuthorRepo_FindBooksByAuthorID(t *testing.T) {
	pool := newTestPool(t)
	cleanTables(t, pool)
	ctx := context.Background()

	pool.Exec(ctx, `INSERT INTO authors (id, name) VALUES ('a1', 'Tolkien')`)
	pool.Exec(ctx, `INSERT INTO books (id, title) VALUES ('b1', 'The Hobbit')`)
	pool.Exec(ctx, `INSERT INTO book_authors (book_id, author_id) VALUES ('b1', 'a1')`)

	repo := postgres.NewAuthorRepo(pool)

	books, err := repo.FindBooksByAuthorID(ctx, "a1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(books) != 1 {
		t.Fatalf("want 1 book, got %d", len(books))
	}
	if books[0].Title != "The Hobbit" {
		t.Fatalf("want title 'The Hobbit', got %s", books[0].Title)
	}
}
