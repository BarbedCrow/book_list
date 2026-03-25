//go:build integration

package postgres_test

import (
	"context"
	"errors"
	"testing"

	"github.com/BarbedCrow/book_list/internal/domain"
	"github.com/BarbedCrow/book_list/internal/postgres"
	listuc "github.com/BarbedCrow/book_list/internal/usecase/list"
)

func TestListRepo(t *testing.T) {
	pool := newTestPool(t)
	cleanTables(t, pool)
	ctx := context.Background()

	// Setup user and book
	pool.Exec(ctx, `INSERT INTO users (id, email, password_hash) VALUES ('u1', 'test@example.com', 'hash')`)
	pool.Exec(ctx, `INSERT INTO books (id, title) VALUES ('b1', 'Go Programming')`)

	repo := postgres.NewListRepo(pool)

	t.Run("save", func(t *testing.T) {
		l := domain.List{ID: "l1", OwnerID: "u1", Name: "My List", Type: domain.ListTypeCustom}
		if err := repo.Save(ctx, l); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("find by owner", func(t *testing.T) {
		lists, err := repo.FindByOwner(ctx, "u1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(lists) != 1 {
			t.Fatalf("want 1 list, got %d", len(lists))
		}
	})

	t.Run("find by id", func(t *testing.T) {
		l, err := repo.FindByID(ctx, "l1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if l.Name != "My List" {
			t.Fatalf("want name 'My List', got %s", l.Name)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.FindByID(ctx, "nonexistent")
		if !errors.Is(err, listuc.ErrListNotFound) {
			t.Fatalf("want ErrListNotFound, got %v", err)
		}
	})

	t.Run("add and list books", func(t *testing.T) {
		if err := repo.AddBook(ctx, "l1", "b1"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		l, err := repo.FindByID(ctx, "l1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(l.Books) != 1 || l.Books[0] != "b1" {
			t.Fatalf("want books [b1], got %v", l.Books)
		}
	})

	t.Run("remove book", func(t *testing.T) {
		if err := repo.RemoveBook(ctx, "l1", "b1"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		l, err := repo.FindByID(ctx, "l1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(l.Books) != 0 {
			t.Fatalf("want 0 books, got %d", len(l.Books))
		}
	})
}
