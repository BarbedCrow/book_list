package list

import (
	"context"
	"fmt"

	"github.com/BarbedCrow/book_list/internal/domain"
)

type GetUserLists struct {
	repo ListRepository
}

func NewGetUserLists(repo ListRepository) *GetUserLists {
	return &GetUserLists{repo: repo}
}

func (uc *GetUserLists) Execute(ctx context.Context, ownerID string) ([]domain.List, error) {
	lists, err := uc.repo.FindByOwner(ctx, ownerID)
	if err != nil {
		return nil, fmt.Errorf("get user lists: %w", err)
	}
	return lists, nil
}
