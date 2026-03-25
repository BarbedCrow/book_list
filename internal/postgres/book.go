package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/BarbedCrow/book_list/internal/domain"
)

type BookRepo struct {
	pool *pgxpool.Pool
}

func NewBookRepo(pool *pgxpool.Pool) *BookRepo {
	return &BookRepo{pool: pool}
}

func (r *BookRepo) FindByTitle(ctx context.Context, title string, limit, offset int) ([]domain.Book, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT b.id, b.title, COALESCE(array_agg(a.name) FILTER (WHERE a.name IS NOT NULL), '{}')
		FROM books b
		LEFT JOIN book_authors ba ON ba.book_id = b.id
		LEFT JOIN authors a ON a.id = ba.author_id
		WHERE b.title ILIKE '%' || $1 || '%'
		GROUP BY b.id, b.title
		LIMIT $2 OFFSET $3
	`, title, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("book find by title: %w", err)
	}
	defer rows.Close()

	var books []domain.Book
	for rows.Next() {
		var b domain.Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Authors); err != nil {
			return nil, fmt.Errorf("book find by title scan: %w", err)
		}
		books = append(books, b)
	}
	return books, rows.Err()
}

func (r *BookRepo) FindByID(ctx context.Context, id string) (domain.Book, error) {
	var b domain.Book
	err := r.pool.QueryRow(ctx, `
		SELECT b.id, b.title, COALESCE(array_agg(a.name) FILTER (WHERE a.name IS NOT NULL), '{}')
		FROM books b
		LEFT JOIN book_authors ba ON ba.book_id = b.id
		LEFT JOIN authors a ON a.id = ba.author_id
		WHERE b.id = $1
		GROUP BY b.id, b.title
	`, id).Scan(&b.ID, &b.Title, &b.Authors)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Book{}, domain.ErrBookNotFound
	}
	if err != nil {
		return domain.Book{}, fmt.Errorf("book find by id: %w", err)
	}
	return b, nil
}
