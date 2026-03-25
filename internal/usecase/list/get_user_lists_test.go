package list_test

import (
	"context"
	"errors"
	"testing"

	"github.com/BarbedCrow/book_list/internal/domain"
	"github.com/BarbedCrow/book_list/internal/usecase/list"
)

func TestGetUserLists(t *testing.T) {
	tests := []struct {
		name      string
		ownerID   string
		repo      *mockListRepo
		wantCount int
		wantErr   bool
	}{
		{
			name:    "found",
			ownerID: "u1",
			repo: &mockListRepo{
				findByOwner: func(_ context.Context, _ string) ([]domain.List, error) {
					return []domain.List{{ID: "l1", OwnerID: "u1", Name: "Favorites"}}, nil
				},
			},
			wantCount: 1,
		},
		{
			name:    "empty",
			ownerID: "u2",
			repo: &mockListRepo{
				findByOwner: func(_ context.Context, _ string) ([]domain.List, error) {
					return nil, nil
				},
			},
			wantCount: 0,
		},
		{
			name:    "repo error",
			ownerID: "u1",
			repo: &mockListRepo{
				findByOwner: func(_ context.Context, _ string) ([]domain.List, error) {
					return nil, errors.New("db down")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := list.NewGetUserLists(tt.repo)
			got, err := uc.Execute(context.Background(), tt.ownerID)
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
				t.Fatalf("want %d lists, got %d", tt.wantCount, len(got))
			}
		})
	}
}
