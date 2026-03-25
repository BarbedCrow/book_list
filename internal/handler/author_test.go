package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BarbedCrow/book_list/internal/domain"
	"github.com/BarbedCrow/book_list/internal/handler"
	authoruc "github.com/BarbedCrow/book_list/internal/usecase/author"
)

type mockAuthorSearcher struct {
	execute func(ctx context.Context, name string) ([]domain.Author, error)
}

func (m *mockAuthorSearcher) Execute(ctx context.Context, name string) ([]domain.Author, error) {
	return m.execute(ctx, name)
}

type mockAuthorDetailer struct {
	execute func(ctx context.Context, id string) (domain.Author, error)
}

func (m *mockAuthorDetailer) Execute(ctx context.Context, id string) (domain.Author, error) {
	return m.execute(ctx, id)
}

type mockAuthorBooksGetter struct {
	execute func(ctx context.Context, authorID string) ([]domain.Book, error)
}

func (m *mockAuthorBooksGetter) Execute(ctx context.Context, authorID string) ([]domain.Book, error) {
	return m.execute(ctx, authorID)
}

func TestAuthorHandler_Search(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		searcher   *mockAuthorSearcher
		wantStatus int
		wantCount  int
	}{
		{
			name:  "success",
			query: "?name=Tolkien",
			searcher: &mockAuthorSearcher{
				execute: func(_ context.Context, _ string) ([]domain.Author, error) {
					return []domain.Author{{ID: "1", Name: "Tolkien"}}, nil
				},
			},
			wantStatus: http.StatusOK,
			wantCount:  1,
		},
		{
			name:       "missing name",
			query:      "",
			searcher:   &mockAuthorSearcher{},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handler.NewAuthorHandler(tt.searcher, nil, nil)
			req := httptest.NewRequest(http.MethodGet, "/authors"+tt.query, nil)
			rec := httptest.NewRecorder()
			h.Search(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("want status %d, got %d", tt.wantStatus, rec.Code)
			}
			if tt.wantStatus == http.StatusOK {
				var authors []domain.Author
				json.NewDecoder(rec.Body).Decode(&authors)
				if len(authors) != tt.wantCount {
					t.Fatalf("want %d authors, got %d", tt.wantCount, len(authors))
				}
			}
		})
	}
}

func TestAuthorHandler_GetDetails(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		detailer   *mockAuthorDetailer
		wantStatus int
	}{
		{
			name: "success",
			id:   "1",
			detailer: &mockAuthorDetailer{
				execute: func(_ context.Context, id string) (domain.Author, error) {
					return domain.Author{ID: id, Name: "Tolkien"}, nil
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "not found",
			id:   "999",
			detailer: &mockAuthorDetailer{
				execute: func(_ context.Context, _ string) (domain.Author, error) {
					return domain.Author{}, authoruc.ErrAuthorNotFound
				},
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handler.NewAuthorHandler(nil, tt.detailer, nil)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /authors/{id}", h.GetDetails)

			req := httptest.NewRequest(http.MethodGet, "/authors/"+tt.id, nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("want status %d, got %d", tt.wantStatus, rec.Code)
			}
		})
	}
}

func TestAuthorHandler_GetBooks(t *testing.T) {
	h := handler.NewAuthorHandler(nil, nil, &mockAuthorBooksGetter{
		execute: func(_ context.Context, _ string) ([]domain.Book, error) {
			return []domain.Book{{ID: "b1", Title: "The Hobbit"}}, nil
		},
	})

	mux := http.NewServeMux()
	mux.HandleFunc("GET /authors/{id}/books", h.GetBooks)

	req := httptest.NewRequest(http.MethodGet, "/authors/1/books", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rec.Code)
	}
}
