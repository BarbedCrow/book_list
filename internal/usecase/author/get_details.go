package author

import (
	"context"
	"fmt"

	"github.com/BarbedCrow/book_list/internal/domain"
)

type GetAuthorDetails struct {
	repo AuthorRepository
}

func NewGetAuthorDetails(repo AuthorRepository) *GetAuthorDetails {
	return &GetAuthorDetails{repo: repo}
}

func (uc *GetAuthorDetails) Execute(ctx context.Context, id string) (domain.Author, error) {
	a, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return domain.Author{}, fmt.Errorf("get author details: %w", err)
	}
	return a, nil
}
