package author_test

import (
	"context"
	"errors"
	"testing"

	"github.com/BarbedCrow/book_list/internal/domain"
	"github.com/BarbedCrow/book_list/internal/usecase/author"
)

func TestGetAuthorDetails(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		repo    *mockAuthorRepo
		wantID  string
		wantErr bool
	}{
		{
			name: "found",
			id:   "1",
			repo: &mockAuthorRepo{
				findByID: func(_ context.Context, id string) (domain.Author, error) {
					return domain.Author{ID: id, Name: "Tolkien"}, nil
				},
			},
			wantID: "1",
		},
		{
			name: "not found",
			id:   "999",
			repo: &mockAuthorRepo{
				findByID: func(_ context.Context, _ string) (domain.Author, error) {
					return domain.Author{}, author.ErrAuthorNotFound
				},
			},
			wantErr: true,
		},
		{
			name: "repo error",
			id:   "1",
			repo: &mockAuthorRepo{
				findByID: func(_ context.Context, _ string) (domain.Author, error) {
					return domain.Author{}, errors.New("db down")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := author.NewGetAuthorDetails(tt.repo)
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
				t.Fatalf("want author ID %s, got %s", tt.wantID, got.ID)
			}
		})
	}
}
