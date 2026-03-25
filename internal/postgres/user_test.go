//go:build integration

package postgres_test

import (
	"context"
	"errors"
	"testing"

	"github.com/BarbedCrow/book_list/internal/domain"
	"github.com/BarbedCrow/book_list/internal/postgres"
	useruc "github.com/BarbedCrow/book_list/internal/usecase/user"
)

func TestUserRepo_SaveAndFindByEmail(t *testing.T) {
	pool := newTestPool(t)
	cleanTables(t, pool)
	ctx := context.Background()

	repo := postgres.NewUserRepo(pool)

	u := domain.User{ID: "u1", Email: "test@example.com", PasswordHash: "hashed"}

	t.Run("save", func(t *testing.T) {
		if err := repo.Save(ctx, u); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("find by email", func(t *testing.T) {
		found, err := repo.FindByEmail(ctx, "test@example.com")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if found.ID != "u1" {
			t.Fatalf("want id 'u1', got %s", found.ID)
		}
	})

	t.Run("duplicate email", func(t *testing.T) {
		dup := domain.User{ID: "u2", Email: "test@example.com", PasswordHash: "hashed2"}
		err := repo.Save(ctx, dup)
		if !errors.Is(err, useruc.ErrDuplicateEmail) {
			t.Fatalf("want ErrDuplicateEmail, got %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.FindByEmail(ctx, "nobody@example.com")
		if !errors.Is(err, useruc.ErrUserNotFound) {
			t.Fatalf("want ErrUserNotFound, got %v", err)
		}
	})
}
