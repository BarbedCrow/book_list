package author_test

import (
	"context"
	"errors"
	"testing"

	"github.com/BarbedCrow/book_list/internal/domain"
	"github.com/BarbedCrow/book_list/internal/usecase/author"
)

func TestGetBooksByAuthor(t *testing.T) {
	tests := []struct {
		name      string
		authorID  string
		repo      *mockAuthorRepo
		wantCount int
		wantErr   bool
	}{
		{
			name:     "found",
			authorID: "1",
			repo: &mockAuthorRepo{
				findBooksByAuthorID: func(_ context.Context, _ string) ([]domain.Book, error) {
					return []domain.Book{
						{ID: "b1", Title: "The Hobbit"},
						{ID: "b2", Title: "The Lord of the Rings"},
					}, nil
				},
			},
			wantCount: 2,
		},
		{
			name:     "empty",
			authorID: "2",
			repo: &mockAuthorRepo{
				findBooksByAuthorID: func(_ context.Context, _ string) ([]domain.Book, error) {
					return nil, nil
				},
			},
			wantCount: 0,
		},
		{
			name:     "repo error",
			authorID: "1",
			repo: &mockAuthorRepo{
				findBooksByAuthorID: func(_ context.Context, _ string) ([]domain.Book, error) {
					return nil, errors.New("db down")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := author.NewGetBooksByAuthor(tt.repo)
			got, err := uc.Execute(context.Background(), tt.authorID)
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
