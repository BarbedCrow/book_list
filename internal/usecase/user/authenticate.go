package user

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrWrongPassword = errors.New("wrong password")
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
		return "", fmt.Errorf("authenticate user: %w", ErrWrongPassword)
	}

	token, err := uc.tokens.Generate(u.ID)
	if err != nil {
		return "", fmt.Errorf("authenticate user: %w", err)
	}

	return token, nil
}
