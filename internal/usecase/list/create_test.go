package list_test

import (
	"context"
	"errors"
	"testing"

	"github.com/BarbedCrow/book_list/internal/domain"
	"github.com/BarbedCrow/book_list/internal/usecase/list"
)

type mockListRepo struct {
	findByOwner func(ctx context.Context, ownerID string) ([]domain.List, error)
	findByID    func(ctx context.Context, id string) (domain.List, error)
	save        func(ctx context.Context, l domain.List) error
	addBook     func(ctx context.Context, listID, bookID string) error
	removeBook  func(ctx context.Context, listID, bookID string) error
}

func (m *mockListRepo) FindByOwner(ctx context.Context, ownerID string) ([]domain.List, error) {
	return m.findByOwner(ctx, ownerID)
}
func (m *mockListRepo) FindByID(ctx context.Context, id string) (domain.List, error) {
	return m.findByID(ctx, id)
}
func (m *mockListRepo) Save(ctx context.Context, l domain.List) error {
	return m.save(ctx, l)
}
func (m *mockListRepo) AddBook(ctx context.Context, listID, bookID string) error {
	return m.addBook(ctx, listID, bookID)
}
func (m *mockListRepo) RemoveBook(ctx context.Context, listID, bookID string) error {
	return m.removeBook(ctx, listID, bookID)
}

func TestCreateCustomList(t *testing.T) {
	tests := []struct {
		name    string
		ownerID string
		lName   string
		repo    *mockListRepo
		wantErr bool
	}{
		{
			name:    "success",
			ownerID: "u1",
			lName:   "My List",
			repo: &mockListRepo{
				save: func(_ context.Context, _ domain.List) error { return nil },
			},
		},
		{
			name:    "repo error",
			ownerID: "u1",
			lName:   "My List",
			repo: &mockListRepo{
				save: func(_ context.Context, _ domain.List) error { return errors.New("db down") },
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := list.NewCreateCustomList(tt.repo, func() string { return "list-1" })
			got, err := uc.Execute(context.Background(), tt.ownerID, tt.lName)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Name != tt.lName {
				t.Fatalf("want name %s, got %s", tt.lName, got.Name)
			}
			if got.Type != domain.ListTypeCustom {
				t.Fatalf("want type custom, got %s", got.Type)
			}
		})
	}
}
