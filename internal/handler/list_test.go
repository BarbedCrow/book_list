package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/BarbedCrow/book_list/internal/domain"
	"github.com/BarbedCrow/book_list/internal/handler"
	listuc "github.com/BarbedCrow/book_list/internal/usecase/list"
)

type mockListCreator struct {
	execute func(ctx context.Context, ownerID, name string) (domain.List, error)
}

func (m *mockListCreator) Execute(ctx context.Context, ownerID, name string) (domain.List, error) {
	return m.execute(ctx, ownerID, name)
}

type mockListGetter struct {
	execute func(ctx context.Context, ownerID string) ([]domain.List, error)
}

func (m *mockListGetter) Execute(ctx context.Context, ownerID string) ([]domain.List, error) {
	return m.execute(ctx, ownerID)
}

type mockListBookAdder struct {
	execute func(ctx context.Context, userID, listID, bookID string) error
}

func (m *mockListBookAdder) Execute(ctx context.Context, userID, listID, bookID string) error {
	return m.execute(ctx, userID, listID, bookID)
}

type mockListBookRemover struct {
	execute func(ctx context.Context, userID, listID, bookID string) error
}

func (m *mockListBookRemover) Execute(ctx context.Context, userID, listID, bookID string) error {
	return m.execute(ctx, userID, listID, bookID)
}

func withUserID(r *http.Request, userID string) *http.Request {
	ctx := context.WithValue(r.Context(), handler.UserIDKey, userID)
	return r.WithContext(ctx)
}

func TestListHandler_GetUserLists(t *testing.T) {
	tests := []struct {
		name       string
		userID     string
		getter     *mockListGetter
		wantStatus int
	}{
		{
			name:   "success",
			userID: "u1",
			getter: &mockListGetter{
				execute: func(_ context.Context, _ string) ([]domain.List, error) {
					return []domain.List{{ID: "l1", Name: "Favorites"}}, nil
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "unauthorized",
			userID:     "",
			getter:     &mockListGetter{},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handler.NewListHandler(nil, tt.getter, nil, nil)
			req := httptest.NewRequest(http.MethodGet, "/lists", nil)
			if tt.userID != "" {
				req = withUserID(req, tt.userID)
			}
			rec := httptest.NewRecorder()
			h.GetUserLists(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("want status %d, got %d", tt.wantStatus, rec.Code)
			}
		})
	}
}

func TestListHandler_Create(t *testing.T) {
	h := handler.NewListHandler(
		&mockListCreator{
			execute: func(_ context.Context, ownerID, name string) (domain.List, error) {
				return domain.List{ID: "l1", OwnerID: ownerID, Name: name, Type: domain.ListTypeCustom}, nil
			},
		}, nil, nil, nil,
	)

	req := httptest.NewRequest(http.MethodPost, "/lists", strings.NewReader(`{"name":"My List"}`))
	req = withUserID(req, "u1")
	rec := httptest.NewRecorder()
	h.Create(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d", rec.Code)
	}
}

func TestListHandler_AddBook(t *testing.T) {
	tests := []struct {
		name       string
		userID     string
		adder      *mockListBookAdder
		wantStatus int
	}{
		{
			name:   "success",
			userID: "u1",
			adder: &mockListBookAdder{
				execute: func(_ context.Context, _, _, _ string) error { return nil },
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:   "not owner",
			userID: "u2",
			adder: &mockListBookAdder{
				execute: func(_ context.Context, _, _, _ string) error { return listuc.ErrNotOwner },
			},
			wantStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handler.NewListHandler(nil, nil, tt.adder, nil)

			mux := http.NewServeMux()
			mux.HandleFunc("POST /lists/{id}/books", h.AddBook)

			req := httptest.NewRequest(http.MethodPost, "/lists/l1/books", strings.NewReader(`{"book_id":"b1"}`))
			req = withUserID(req, tt.userID)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("want status %d, got %d", tt.wantStatus, rec.Code)
			}
		})
	}
}

func TestListHandler_RemoveBook(t *testing.T) {
	h := handler.NewListHandler(nil, nil, nil, &mockListBookRemover{
		execute: func(_ context.Context, _, _, _ string) error { return nil },
	})

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /lists/{id}/books/{book_id}", h.RemoveBook)

	req := httptest.NewRequest(http.MethodDelete, "/lists/l1/books/b1", nil)
	req = withUserID(req, "u1")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("want 204, got %d", rec.Code)
	}
}
