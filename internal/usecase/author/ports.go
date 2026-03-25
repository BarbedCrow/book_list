package author

import (
	"context"

	"github.com/BarbedCrow/book_list/internal/domain"
)

type AuthorRepository interface {
	FindByName(ctx context.Context, name string, limit, offset int) ([]domain.Author, error)
	FindByID(ctx context.Context, id string) (domain.Author, error)
	FindBooksByAuthorID(ctx context.Context, authorID string) ([]domain.Book, error)
}
