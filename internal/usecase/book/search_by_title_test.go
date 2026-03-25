package book_test

import (
	"context"
	"errors"
	"testing"

	"github.com/BarbedCrow/book_list/internal/domain"
	"github.com/BarbedCrow/book_list/internal/usecase/book"
)

type mockBookRepo struct {
	findByTitle func(ctx context.Context, title string, limit, offset int) ([]domain.Book, error)
	findByID    func(ctx context.Context, id string) (domain.Book, error)
}

func (m *mockBookRepo) FindByTitle(ctx context.Context, title string, limit, offset int) ([]domain.Book, error) {
	return m.findByTitle(ctx, title, limit, offset)
}

func (m *mockBookRepo) FindByID(ctx context.Context, id string) (domain.Book, error) {
	return m.findByID(ctx, id)
}

func TestSearchBooksByTitle(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		repo      *mockBookRepo
		wantCount int
		wantErr   bool
	}{
		{
			name:  "found",
			title: "Go",
			repo: &mockBookRepo{
				findByTitle: func(_ context.Context, _ string, _, _ int) ([]domain.Book, error) {
					return []domain.Book{{ID: "1", Title: "Go Programming"}}, nil
				},
			},
			wantCount: 1,
		},
		{
			name:  "empty",
			title: "nonexistent",
			repo: &mockBookRepo{
				findByTitle: func(_ context.Context, _ string, _, _ int) ([]domain.Book, error) {
					return nil, nil
				},
			},
			wantCount: 0,
		},
		{
			name:  "repo error",
			title: "Go",
			repo: &mockBookRepo{
				findByTitle: func(_ context.Context, _ string, _, _ int) ([]domain.Book, error) {
					return nil, errors.New("db down")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := book.NewSearchBooksByTitle(tt.repo)
			got, err := uc.Execute(context.Background(), tt.title, 20, 0)
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
				t.Fatalf("want %d books, got %d", tt.wantCount, len(got))
			}
		})
	}
}
