package user

import (
	"context"
	"fmt"

	"github.com/BarbedCrow/book_list/internal/domain"
)

type AuthenticateUser struct {
	repo   UserRepository
	hasher PasswordHasher
	tokens TokenProvider
}

func NewAuthenticateUser(repo UserRepository, hasher PasswordHasher, tokens TokenProvider) *AuthenticateUser {
	return &AuthenticateUser{repo: repo, hasher: hasher, tokens: tokens}
}

func (uc *AuthenticateUser) Execute(ctx context.Context, email, password string) (string, error) {
	u, err := uc.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("authenticate user: %w", err)
	}

	if err := uc.hasher.Verify(u.PasswordHash, password); err != nil {
		return "", fmt.Errorf("authenticate user: %w", domain.ErrWrongPassword)
	}

	token, err := uc.tokens.Generate(u.ID)
	if err != nil {
		return "", fmt.Errorf("authenticate user: %w", err)
	}

	return token, nil
}
