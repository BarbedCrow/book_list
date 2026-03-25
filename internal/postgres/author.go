package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/BarbedCrow/book_list/internal/domain"
	authoruc "github.com/BarbedCrow/book_list/internal/usecase/author"
)

type AuthorRepo struct {
	pool *pgxpool.Pool
}

func NewAuthorRepo(pool *pgxpool.Pool) *AuthorRepo {
	return &AuthorRepo{pool: pool}
}

func (r *AuthorRepo) FindByName(ctx context.Context, name string) ([]domain.Author, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT a.id, a.name, COALESCE(array_agg(b.title) FILTER (WHERE b.title IS NOT NULL), '{}')
		FROM authors a
		LEFT JOIN book_authors ba ON ba.author_id = a.id
		LEFT JOIN books b ON b.id = ba.book_id
		WHERE a.name ILIKE '%' || $1 || '%'
		GROUP BY a.id, a.name
	`, name)
	if err != nil {
		return nil, fmt.Errorf("author find by name: %w", err)
	}
	defer rows.Close()

	var authors []domain.Author
	for rows.Next() {
		var a domain.Author
		if err := rows.Scan(&a.ID, &a.Name, &a.Books); err != nil {
			return nil, fmt.Errorf("author find by name scan: %w", err)
		}
		authors = append(authors, a)
	}
	return authors, rows.Err()
}

func (r *AuthorRepo) FindByID(ctx context.Context, id string) (domain.Author, error) {
	var a domain.Author
	err := r.pool.QueryRow(ctx, `
		SELECT a.id, a.name, COALESCE(array_agg(b.title) FILTER (WHERE b.title IS NOT NULL), '{}')
		FROM authors a
		LEFT JOIN book_authors ba ON ba.author_id = a.id
		LEFT JOIN books b ON b.id = ba.book_id
		WHERE a.id = $1
		GROUP BY a.id, a.name
	`, id).Scan(&a.ID, &a.Name, &a.Books)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Author{}, authoruc.ErrAuthorNotFound
	}
	if err != nil {
		return domain.Author{}, fmt.Errorf("author find by id: %w", err)
	}
	return a, nil
}

func (r *AuthorRepo) FindBooksByAuthorID(ctx context.Context, authorID string) ([]domain.Book, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT b.id, b.title, COALESCE(array_agg(a2.name) FILTER (WHERE a2.name IS NOT NULL), '{}')
		FROM books b
		JOIN book_authors ba ON ba.book_id = b.id
		LEFT JOIN book_authors ba2 ON ba2.book_id = b.id
		LEFT JOIN authors a2 ON a2.id = ba2.author_id
		WHERE ba.author_id = $1
		GROUP BY b.id, b.title
	`, authorID)
	if err != nil {
		return nil, fmt.Errorf("author find books: %w", err)
	}
	defer rows.Close()

	var books []domain.Book
	for rows.Next() {
		var b domain.Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Authors); err != nil {
			return nil, fmt.Errorf("author find books scan: %w", err)
		}
		books = append(books, b)
	}
	return books, rows.Err()
}
