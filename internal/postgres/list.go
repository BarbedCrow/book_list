package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/BarbedCrow/book_list/internal/domain"
	listuc "github.com/BarbedCrow/book_list/internal/usecase/list"
)

type ListRepo struct {
	pool *pgxpool.Pool
}

func NewListRepo(pool *pgxpool.Pool) *ListRepo {
	return &ListRepo{pool: pool}
}

func (r *ListRepo) FindByOwner(ctx context.Context, ownerID string) ([]domain.List, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT l.id, l.owner_id, l.name, l.type,
		       COALESCE(array_agg(lb.book_id) FILTER (WHERE lb.book_id IS NOT NULL), '{}')
		FROM lists l
		LEFT JOIN list_books lb ON lb.list_id = l.id
		WHERE l.owner_id = $1
		GROUP BY l.id, l.owner_id, l.name, l.type
	`, ownerID)
	if err != nil {
		return nil, fmt.Errorf("list find by owner: %w", err)
	}
	defer rows.Close()

	var lists []domain.List
	for rows.Next() {
		var l domain.List
		if err := rows.Scan(&l.ID, &l.OwnerID, &l.Name, &l.Type, &l.Books); err != nil {
			return nil, fmt.Errorf("list find by owner scan: %w", err)
		}
		lists = append(lists, l)
	}
	return lists, rows.Err()
}

func (r *ListRepo) FindByID(ctx context.Context, id string) (domain.List, error) {
	var l domain.List
	err := r.pool.QueryRow(ctx, `
		SELECT l.id, l.owner_id, l.name, l.type,
		       COALESCE(array_agg(lb.book_id) FILTER (WHERE lb.book_id IS NOT NULL), '{}')
		FROM lists l
		LEFT JOIN list_books lb ON lb.list_id = l.id
		WHERE l.id = $1
		GROUP BY l.id, l.owner_id, l.name, l.type
	`, id).Scan(&l.ID, &l.OwnerID, &l.Name, &l.Type, &l.Books)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.List{}, listuc.ErrListNotFound
	}
	if err != nil {
		return domain.List{}, fmt.Errorf("list find by id: %w", err)
	}
	return l, nil
}

func (r *ListRepo) Save(ctx context.Context, l domain.List) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO lists (id, owner_id, name, type) VALUES ($1, $2, $3, $4)
		 ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, type = EXCLUDED.type`,
		l.ID, l.OwnerID, l.Name, string(l.Type),
	)
	if err != nil {
		return fmt.Errorf("list save: %w", err)
	}
	return nil
}

func (r *ListRepo) AddBook(ctx context.Context, listID, bookID string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO list_books (list_id, book_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		listID, bookID,
	)
	if err != nil {
		return fmt.Errorf("list add book: %w", err)
	}
	return nil
}

func (r *ListRepo) RemoveBook(ctx context.Context, listID, bookID string) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM list_books WHERE list_id = $1 AND book_id = $2`,
		listID, bookID,
	)
	if err != nil {
		return fmt.Errorf("list remove book: %w", err)
	}
	return nil
}
