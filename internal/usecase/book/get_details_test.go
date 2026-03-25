package book_test

import (
	"context"
	"errors"
	"testing"

	"github.com/BarbedCrow/book_list/internal/domain"
	"github.com/BarbedCrow/book_list/internal/usecase/book"
)

func TestGetBookDetails(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		repo    *mockBookRepo
		wantID  string
		wantErr bool
	}{
		{
			name: "found",
			id:   "1",
			repo: &mockBookRepo{
				findByID: func(_ context.Context, id string) (domain.Book, error) {
					return domain.Book{ID: id, Title: "Go Programming"}, nil
				},
			},
			wantID: "1",
		},
		{
			name: "not found",
			id:   "999",
			repo: &mockBookRepo{
				findByID: func(_ context.Context, _ string) (domain.Book, error) {
					return domain.Book{}, book.ErrBookNotFound
				},
			},
			wantErr: true,
		},
		{
			name: "repo error",
			id:   "1",
			repo: &mockBookRepo{
				findByID: func(_ context.Context, _ string) (domain.Book, error) {
					return domain.Book{}, errors.New("db down")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := book.NewGetBookDetails(tt.repo)
			got, err := uc.Execute(context.Background(), tt.id)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.ID != tt.wantID {
				t.Fatalf("want book ID %s, got %s", tt.wantID, got.ID)
			}
		})
	}
}
