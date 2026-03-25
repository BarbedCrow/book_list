package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/BarbedCrow/book_list/internal/domain"
	bookuc "github.com/BarbedCrow/book_list/internal/usecase/book"
)

type BookSearcher interface {
	Execute(ctx context.Context, title string) ([]domain.Book, error)
}

type BookDetailer interface {
	Execute(ctx context.Context, id string) (domain.Book, error)
}

type BookHandler struct {
	search  BookSearcher
	details BookDetailer
}

func NewBookHandler(search BookSearcher, details BookDetailer) *BookHandler {
	return &BookHandler{search: search, details: details}
}

func (h *BookHandler) Search(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	if title == "" {
		writeError(w, http.StatusBadRequest, "missing title query parameter")
		return
	}

	books, err := h.search.Execute(r.Context(), title)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if books == nil {
		books = []domain.Book{}
	}
	writeJSON(w, http.StatusOK, books)
}

func (h *BookHandler) GetDetails(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing book id")
		return
	}

	b, err := h.details.Execute(r.Context(), id)
	if err != nil {
		if errors.Is(err, bookuc.ErrBookNotFound) {
			writeError(w, http.StatusNotFound, "book not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, b)
}
