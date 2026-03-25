package book

import (
	"context"
	"fmt"

	"github.com/BarbedCrow/book_list/internal/domain"
)

type SearchBooksByTitle struct {
	repo BookRepository
}

func NewSearchBooksByTitle(repo BookRepository) *SearchBooksByTitle {
	return &SearchBooksByTitle{repo: repo}
}

func (uc *SearchBooksByTitle) Execute(ctx context.Context, title string, limit, offset int) ([]domain.Book, error) {
	books, err := uc.repo.FindByTitle(ctx, title, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("search books by title: %w", err)
	}
	return books, nil
}
