package list

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrListNotFound = errors.New("list not found")
	ErrNotOwner     = errors.New("user is not the list owner")
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
		return ErrNotOwner
	}

	if err := uc.repo.AddBook(ctx, listID, bookID); err != nil {
		return fmt.Errorf("add book to list: %w", err)
	}

	return nil
}
