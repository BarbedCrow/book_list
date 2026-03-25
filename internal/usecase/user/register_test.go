package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/BarbedCrow/book_list/internal/domain"
	"github.com/BarbedCrow/book_list/internal/usecase/user"
)

type mockUserRepo struct {
	findByEmail func(ctx context.Context, email string) (domain.User, error)
	save        func(ctx context.Context, u domain.User) error
}

func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	return m.findByEmail(ctx, email)
}

func (m *mockUserRepo) Save(ctx context.Context, u domain.User) error {
	return m.save(ctx, u)
}

type mockHasher struct {
	hash   func(password string) (string, error)
	verify func(hash, password string) error
}

func (m *mockHasher) Hash(password string) (string, error) { return m.hash(password) }
func (m *mockHasher) Verify(hash, password string) error   { return m.verify(hash, password) }

type mockTokenProvider struct {
	generate func(userID string) (string, error)
}

func (m *mockTokenProvider) Generate(userID string) (string, error) { return m.generate(userID) }

func TestRegisterUser(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		pass    string
		repo    *mockUserRepo
		hasher  *mockHasher
		wantErr bool
	}{
		{
			name:  "success",
			email: "a@b.com",
			pass:  "secret",
			repo: &mockUserRepo{
				save: func(_ context.Context, _ domain.User) error { return nil },
			},
			hasher: &mockHasher{
				hash: func(_ string) (string, error) { return "hashed", nil },
			},
		},
		{
			name:  "duplicate email",
			email: "a@b.com",
			pass:  "secret",
			repo: &mockUserRepo{
				save: func(_ context.Context, _ domain.User) error { return domain.ErrDuplicateEmail },
			},
			hasher: &mockHasher{
				hash: func(_ string) (string, error) { return "hashed", nil },
			},
			wantErr: true,
		},
		{
			name:  "hasher error",
			email: "a@b.com",
			pass:  "secret",
			repo:  &mockUserRepo{},
			hasher: &mockHasher{
				hash: func(_ string) (string, error) { return "", errors.New("hash fail") },
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := user.NewRegisterUser(tt.repo, tt.hasher, func() string { return "uuid-1" })
			got, err := uc.Execute(context.Background(), tt.email, tt.pass)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Email != tt.email {
				t.Fatalf("want email %s, got %s", tt.email, got.Email)
			}
		})
	}
}
