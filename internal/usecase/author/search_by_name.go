package author

import (
	"context"
	"fmt"

	"github.com/BarbedCrow/book_list/internal/domain"
)

type SearchAuthorsByName struct {
	repo AuthorRepository
}

func NewSearchAuthorsByName(repo AuthorRepository) *SearchAuthorsByName {
	return &SearchAuthorsByName{repo: repo}
}

func (uc *SearchAuthorsByName) Execute(ctx context.Context, name string, limit, offset int) ([]domain.Author, error) {
	authors, err := uc.repo.FindByName(ctx, name, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("search authors by name: %w", err)
	}
	return authors, nil
}
