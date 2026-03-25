package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/BarbedCrow/book_list/internal/domain"
	useruc "github.com/BarbedCrow/book_list/internal/usecase/user"
)

type UserRegisterer interface {
	Execute(ctx context.Context, email, password string) (domain.User, error)
}

type UserAuthenticator interface {
	Execute(ctx context.Context, email, password string) (string, error)
}

type UserHandler struct {
	register     UserRegisterer
	authenticate UserAuthenticator
}

func NewUserHandler(register UserRegisterer, authenticate UserAuthenticator) *UserHandler {
	return &UserHandler{register: register, authenticate: authenticate}
}

type userRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req userRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	u, err := h.register.Execute(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, useruc.ErrDuplicateEmail) {
			writeError(w, http.StatusConflict, "email already registered")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusCreated, registerResponse{ID: u.ID, Email: u.Email})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req userRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	token, err := h.authenticate.Execute(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, useruc.ErrUserNotFound) || errors.Is(err, useruc.ErrWrongPassword) {
			writeError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, loginResponse{Token: token})
}
