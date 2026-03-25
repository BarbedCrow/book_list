package user

import (
	"context"

	"github.com/BarbedCrow/book_list/internal/domain"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	Save(ctx context.Context, u domain.User) error
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(hash, password string) error
}

type TokenProvider interface {
	Generate(userID string) (string, error)
	Validate(token string) (string, error)
}
