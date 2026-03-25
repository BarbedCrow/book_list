package list_test

import (
	"context"
	"errors"
	"testing"

	"github.com/BarbedCrow/book_list/internal/domain"
	"github.com/BarbedCrow/book_list/internal/usecase/list"
)

func TestRemoveBookFromList(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		listID  string
		bookID  string
		repo    *mockListRepo
		wantErr error
	}{
		{
			name:   "success",
			userID: "u1",
			listID: "l1",
			bookID: "b1",
			repo: &mockListRepo{
				findByID: func(_ context.Context, _ string) (domain.List, error) {
					return domain.List{ID: "l1", OwnerID: "u1"}, nil
				},
				removeBook: func(_ context.Context, _, _ string) error { return nil },
			},
		},
		{
			name:   "not owner",
			userID: "u2",
			listID: "l1",
			bookID: "b1",
			repo: &mockListRepo{
				findByID: func(_ context.Context, _ string) (domain.List, error) {
					return domain.List{ID: "l1", OwnerID: "u1"}, nil
				},
			},
			wantErr: domain.ErrNotOwner,
		},
		{
			name:   "list not found",
			userID: "u1",
			listID: "l999",
			bookID: "b1",
			repo: &mockListRepo{
				findByID: func(_ context.Context, _ string) (domain.List, error) {
					return domain.List{}, domain.ErrListNotFound
				},
			},
			wantErr: domain.ErrListNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := list.NewRemoveBookFromList(tt.repo)
			err := uc.Execute(context.Background(), tt.userID, tt.listID, tt.bookID)
			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("want error %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
