package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/BarbedCrow/book_list/internal/domain"
	authoruc "github.com/BarbedCrow/book_list/internal/usecase/author"
)

type AuthorSearcher interface {
	Execute(ctx context.Context, name string) ([]domain.Author, error)
}

type AuthorDetailer interface {
	Execute(ctx context.Context, id string) (domain.Author, error)
}

type AuthorBooksGetter interface {
	Execute(ctx context.Context, authorID string) ([]domain.Book, error)
}

type AuthorHandler struct {
	search   AuthorSearcher
	details  AuthorDetailer
	getBooks AuthorBooksGetter
}

func NewAuthorHandler(search AuthorSearcher, details AuthorDetailer, getBooks AuthorBooksGetter) *AuthorHandler {
	return &AuthorHandler{search: search, details: details, getBooks: getBooks}
}

func (h *AuthorHandler) Search(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		writeError(w, http.StatusBadRequest, "missing name query parameter")
		return
	}

	authors, err := h.search.Execute(r.Context(), name)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if authors == nil {
		authors = []domain.Author{}
	}
	writeJSON(w, http.StatusOK, authors)
}

func (h *AuthorHandler) GetDetails(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing author id")
		return
	}

	a, err := h.details.Execute(r.Context(), id)
	if err != nil {
		if errors.Is(err, authoruc.ErrAuthorNotFound) {
			writeError(w, http.StatusNotFound, "author not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, a)
}

func (h *AuthorHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing author id")
		return
	}

	books, err := h.getBooks.Execute(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if books == nil {
		books = []domain.Book{}
	}
	writeJSON(w, http.StatusOK, books)
}
