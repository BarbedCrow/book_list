package list

import (
	"context"
	"fmt"

	"github.com/BarbedCrow/book_list/internal/domain"
)

type CreateCustomList struct {
	repo  ListRepository
	idGen func() string
}

func NewCreateCustomList(repo ListRepository, idGen func() string) *CreateCustomList {
	return &CreateCustomList{repo: repo, idGen: idGen}
}

func (uc *CreateCustomList) Execute(ctx context.Context, ownerID, name string) (domain.List, error) {
	l := domain.List{
		ID:      uc.idGen(),
		OwnerID: ownerID,
		Name:    name,
		Type:    domain.ListTypeCustom,
	}

	if err := uc.repo.Save(ctx, l); err != nil {
		return domain.List{}, fmt.Errorf("create custom list: %w", err)
	}

	return l, nil
}
