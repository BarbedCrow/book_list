package author_test

import (
	"context"
	"errors"
	"testing"

	"github.com/BarbedCrow/book_list/internal/domain"
	"github.com/BarbedCrow/book_list/internal/usecase/author"
)

type mockAuthorRepo struct {
	findByName          func(ctx context.Context, name string, limit, offset int) ([]domain.Author, error)
	findByID            func(ctx context.Context, id string) (domain.Author, error)
	findBooksByAuthorID func(ctx context.Context, authorID string) ([]domain.Book, error)
}

func (m *mockAuthorRepo) FindByName(ctx context.Context, name string, limit, offset int) ([]domain.Author, error) {
	return m.findByName(ctx, name, limit, offset)
}

func (m *mockAuthorRepo) FindByID(ctx context.Context, id string) (domain.Author, error) {
	return m.findByID(ctx, id)
}

func (m *mockAuthorRepo) FindBooksByAuthorID(ctx context.Context, authorID string) ([]domain.Book, error) {
	return m.findBooksByAuthorID(ctx, authorID)
}

func TestSearchAuthorsByName(t *testing.T) {
	tests := []struct {
		name      string
		query     string
		repo      *mockAuthorRepo
		wantCount int
		wantErr   bool
	}{
		{
			name:  "found",
			query: "Tolkien",
			repo: &mockAuthorRepo{
				findByName: func(_ context.Context, _ string, _, _ int) ([]domain.Author, error) {
					return []domain.Author{{ID: "1", Name: "J.R.R. Tolkien"}}, nil
				},
			},
			wantCount: 1,
		},
		{
			name:  "empty",
			query: "nobody",
			repo: &mockAuthorRepo{
				findByName: func(_ context.Context, _ string, _, _ int) ([]domain.Author, error) {
					return nil, nil
				},
			},
			wantCount: 0,
		},
		{
			name:  "repo error",
			query: "Tolkien",
			repo: &mockAuthorRepo{
				findByName: func(_ context.Context, _ string, _, _ int) ([]domain.Author, error) {
					return nil, errors.New("db down")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := author.NewSearchAuthorsByName(tt.repo)
			got, err := uc.Execute(context.Background(), tt.query, 20, 0)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != tt.wantCount {
				t.Fatalf("want %d authors, got %d", tt.wantCount, len(got))
			}
		})
	}
}
