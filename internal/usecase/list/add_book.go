package list

import (
	"context"
	"fmt"

	"github.com/BarbedCrow/book_list/internal/domain"
)

type AddBookToList struct {
	repo ListRepository
}

func NewAddBookToList(repo ListRepository) *AddBookToList {
	return &AddBookToList{repo: repo}
}

func (uc *AddBookToList) Execute(ctx context.Context, userID, listID, bookID string) error {
	l, err := uc.repo.FindByID(ctx, listID)
	if err != nil {
		return fmt.Errorf("add book to list: %w", err)
	}

	if l.OwnerID != userID {
		return domain.ErrNotOwner
	}

	if err := uc.repo.AddBook(ctx, listID, bookID); err != nil {
		return fmt.Errorf("add book to list: %w", err)
	}

	return nil
}
