package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/BarbedCrow/book_list/internal/domain"
	"github.com/BarbedCrow/book_list/internal/usecase/user"
)

func TestAuthenticateUser(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		pass    string
		repo    *mockUserRepo
		hasher  *mockHasher
		tokens  *mockTokenProvider
		wantErr bool
	}{
		{
			name:  "success",
			email: "a@b.com",
			pass:  "secret",
			repo: &mockUserRepo{
				findByEmail: func(_ context.Context, _ string) (domain.User, error) {
					return domain.User{ID: "1", Email: "a@b.com", PasswordHash: "hashed"}, nil
				},
			},
			hasher: &mockHasher{
				verify: func(_, _ string) error { return nil },
			},
			tokens: &mockTokenProvider{
				generate: func(_ string) (string, error) { return "token-123", nil },
			},
		},
		{
			name:  "wrong email",
			email: "unknown@b.com",
			pass:  "secret",
			repo: &mockUserRepo{
				findByEmail: func(_ context.Context, _ string) (domain.User, error) {
					return domain.User{}, domain.ErrUserNotFound
				},
			},
			hasher:  &mockHasher{},
			tokens:  &mockTokenProvider{},
			wantErr: true,
		},
		{
			name:  "wrong password",
			email: "a@b.com",
			pass:  "wrong",
			repo: &mockUserRepo{
				findByEmail: func(_ context.Context, _ string) (domain.User, error) {
					return domain.User{ID: "1", Email: "a@b.com", PasswordHash: "hashed"}, nil
				},
			},
			hasher: &mockHasher{
				verify: func(_, _ string) error { return errors.New("mismatch") },
			},
			tokens:  &mockTokenProvider{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := user.NewAuthenticateUser(tt.repo, tt.hasher, tt.tokens)
			token, err := uc.Execute(context.Background(), tt.email, tt.pass)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if token == "" {
				t.Fatal("expected non-empty token")
			}
		})
	}
}
