package author

import (
	"context"
	"fmt"

	"github.com/BarbedCrow/book_list/internal/domain"
)

type GetBooksByAuthor struct {
	repo AuthorRepository
}

func NewGetBooksByAuthor(repo AuthorRepository) *GetBooksByAuthor {
	return &GetBooksByAuthor{repo: repo}
}

func (uc *GetBooksByAuthor) Execute(ctx context.Context, authorID string) ([]domain.Book, error) {
	books, err := uc.repo.FindBooksByAuthorID(ctx, authorID)
	if err != nil {
		return nil, fmt.Errorf("get books by author: %w", err)
	}
	return books, nil
}
