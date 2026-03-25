package book

import (
	"context"

	"github.com/BarbedCrow/book_list/internal/domain"
)

type BookRepository interface {
	FindByTitle(ctx context.Context, title string, limit, offset int) ([]domain.Book, error)
	FindByID(ctx context.Context, id string) (domain.Book, error)
}
