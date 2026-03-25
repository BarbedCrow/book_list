package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/BarbedCrow/book_list/internal/domain"
	listuc "github.com/BarbedCrow/book_list/internal/usecase/list"
)

type ListCreator interface {
	Execute(ctx context.Context, ownerID, name string) (domain.List, error)
}

type ListGetter interface {
	Execute(ctx context.Context, ownerID string) ([]domain.List, error)
}

type ListBookAdder interface {
	Execute(ctx context.Context, userID, listID, bookID string) error
}

type ListBookRemover interface {
	Execute(ctx context.Context, userID, listID, bookID string) error
}

type ListHandler struct {
	create     ListCreator
	getLists   ListGetter
	addBook    ListBookAdder
	removeBook ListBookRemover
}

func NewListHandler(create ListCreator, getLists ListGetter, addBook ListBookAdder, removeBook ListBookRemover) *ListHandler {
	return &ListHandler{create: create, getLists: getLists, addBook: addBook, removeBook: removeBook}
}

func (h *ListHandler) GetUserLists(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	lists, err := h.getLists.Execute(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if lists == nil {
		lists = []domain.List{}
	}
	writeJSON(w, http.StatusOK, lists)
}

type createListRequest struct {
	Name string `json:"name"`
}

func (h *ListHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req createListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	l, err := h.create.Execute(r.Context(), userID, req.Name)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusCreated, l)
}

type listBookRequest struct {
	BookID string `json:"book_id"`
}

func (h *ListHandler) AddBook(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	listID := r.PathValue("id")
	if listID == "" {
		writeError(w, http.StatusBadRequest, "missing list id")
		return
	}

	var req listBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	err := h.addBook.Execute(r.Context(), userID, listID, req.BookID)
	if err != nil {
		if errors.Is(err, listuc.ErrNotOwner) {
			writeError(w, http.StatusForbidden, "not the list owner")
			return
		}
		if errors.Is(err, listuc.ErrListNotFound) {
			writeError(w, http.StatusNotFound, "list not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ListHandler) RemoveBook(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	listID := r.PathValue("id")
	bookID := r.PathValue("book_id")
	if listID == "" || bookID == "" {
		writeError(w, http.StatusBadRequest, "missing list or book id")
		return
	}

	err := h.removeBook.Execute(r.Context(), userID, listID, bookID)
	if err != nil {
		if errors.Is(err, listuc.ErrNotOwner) {
			writeError(w, http.StatusForbidden, "not the list owner")
			return
		}
		if errors.Is(err, listuc.ErrListNotFound) {
			writeError(w, http.StatusNotFound, "list not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
