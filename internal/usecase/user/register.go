package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/BarbedCrow/book_list/internal/domain"
)

var ErrDuplicateEmail = errors.New("email already registered")

type RegisterUser struct {
	repo   UserRepository
	hasher PasswordHasher
	idGen  func() string
}

func NewRegisterUser(repo UserRepository, hasher PasswordHasher, idGen func() string) *RegisterUser {
	return &RegisterUser{repo: repo, hasher: hasher, idGen: idGen}
}

func (uc *RegisterUser) Execute(ctx context.Context, email, password string) (domain.User, error) {
	hash, err := uc.hasher.Hash(password)
	if err != nil {
		return domain.User{}, fmt.Errorf("register user: %w", err)
	}

	u := domain.User{
		ID:           uc.idGen(),
		Email:        email,
		PasswordHash: hash,
	}

	if err := uc.repo.Save(ctx, u); err != nil {
		return domain.User{}, fmt.Errorf("register user: %w", err)
	}

	return u, nil
}
