package book

import (
	"context"
	"fmt"

	"github.com/BarbedCrow/book_list/internal/domain"
)

type GetBookDetails struct {
	repo BookRepository
}

func NewGetBookDetails(repo BookRepository) *GetBookDetails {
	return &GetBookDetails{repo: repo}
}

func (uc *GetBookDetails) Execute(ctx context.Context, id string) (domain.Book, error) {
	book, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return domain.Book{}, fmt.Errorf("get book details: %w", err)
	}
	return book, nil
}
