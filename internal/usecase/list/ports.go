package list

import (
	"context"

	"github.com/BarbedCrow/book_list/internal/domain"
)

type ListRepository interface {
	FindByOwner(ctx context.Context, ownerID string) ([]domain.List, error)
	FindByID(ctx context.Context, id string) (domain.List, error)
	Save(ctx context.Context, l domain.List) error
	AddBook(ctx context.Context, listID, bookID string) error
	RemoveBook(ctx context.Context, listID, bookID string) error
}
