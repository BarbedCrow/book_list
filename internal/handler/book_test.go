package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BarbedCrow/book_list/internal/domain"
	"github.com/BarbedCrow/book_list/internal/handler"
)

type mockBookSearcher struct {
	execute func(ctx context.Context, title string, limit, offset int) ([]domain.Book, error)
}

func (m *mockBookSearcher) Execute(ctx context.Context, title string, limit, offset int) ([]domain.Book, error) {
	return m.execute(ctx, title, limit, offset)
}

type mockBookDetailer struct {
	execute func(ctx context.Context, id string) (domain.Book, error)
}

func (m *mockBookDetailer) Execute(ctx context.Context, id string) (domain.Book, error) {
	return m.execute(ctx, id)
}

func TestBookHandler_Search(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		searcher   *mockBookSearcher
		wantStatus int
		wantCount  int
	}{
		{
			name:  "success",
			query: "?title=Go",
			searcher: &mockBookSearcher{
				execute: func(_ context.Context, _ string, _, _ int) ([]domain.Book, error) {
					return []domain.Book{{ID: "1", Title: "Go Programming"}}, nil
				},
			},
			wantStatus: http.StatusOK,
			wantCount:  1,
		},
		{
			name:  "empty result",
			query: "?title=nothing",
			searcher: &mockBookSearcher{
				execute: func(_ context.Context, _ string, _, _ int) ([]domain.Book, error) {
					return nil, nil
				},
			},
			wantStatus: http.StatusOK,
			wantCount:  0,
		},
		{
			name:       "missing title",
			query:      "",
			searcher:   &mockBookSearcher{},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handler.NewBookHandler(tt.searcher, nil)
			req := httptest.NewRequest(http.MethodGet, "/books"+tt.query, nil)
			rec := httptest.NewRecorder()
			h.Search(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("want status %d, got %d", tt.wantStatus, rec.Code)
			}
			if tt.wantStatus == http.StatusOK {
				var books []domain.Book
				json.NewDecoder(rec.Body).Decode(&books)
				if len(books) != tt.wantCount {
					t.Fatalf("want %d books, got %d", tt.wantCount, len(books))
				}
			}
		})
	}
}

func TestBookHandler_GetDetails(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		detailer   *mockBookDetailer
		wantStatus int
	}{
		{
			name: "success",
			id:   "1",
			detailer: &mockBookDetailer{
				execute: func(_ context.Context, id string) (domain.Book, error) {
					return domain.Book{ID: id, Title: "Go Programming"}, nil
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "not found",
			id:   "999",
			detailer: &mockBookDetailer{
				execute: func(_ context.Context, _ string) (domain.Book, error) {
					return domain.Book{}, domain.ErrBookNotFound
				},
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handler.NewBookHandler(nil, tt.detailer)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /books/{id}", h.GetDetails)

			req := httptest.NewRequest(http.MethodGet, "/books/"+tt.id, nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("want status %d, got %d", tt.wantStatus, rec.Code)
			}
		})
	}
}
