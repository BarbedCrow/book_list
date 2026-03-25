package list

import (
	"context"
	"fmt"

	"github.com/BarbedCrow/book_list/internal/domain"
)

type RemoveBookFromList struct {
	repo ListRepository
}

func NewRemoveBookFromList(repo ListRepository) *RemoveBookFromList {
	return &RemoveBookFromList{repo: repo}
}

func (uc *RemoveBookFromList) Execute(ctx context.Context, userID, listID, bookID string) error {
	l, err := uc.repo.FindByID(ctx, listID)
	if err != nil {
		return fmt.Errorf("remove book from list: %w", err)
	}

	if l.OwnerID != userID {
		return domain.ErrNotOwner
	}

	if err := uc.repo.RemoveBook(ctx, listID, bookID); err != nil {
		return fmt.Errorf("remove book from list: %w", err)
	}

	return nil
}
