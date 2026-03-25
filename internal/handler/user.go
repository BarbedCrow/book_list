package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/BarbedCrow/book_list/internal/domain"
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

const maxPasswordBytes = 72

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req userRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	if !isValidEmail(req.Email) {
		writeError(w, http.StatusBadRequest, "invalid email format")
		return
	}

	if utf8.RuneCountInString(req.Password) < 8 {
		writeError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}
	if len(req.Password) > maxPasswordBytes {
		writeError(w, http.StatusBadRequest, "password must be at most 72 bytes")
		return
	}

	u, err := h.register.Execute(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrDuplicateEmail) {
			writeError(w, http.StatusConflict, "email already registered")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusCreated, registerResponse{ID: u.ID, Email: u.Email})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

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
		if errors.Is(err, domain.ErrUserNotFound) || errors.Is(err, domain.ErrWrongPassword) {
			writeError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, loginResponse{Token: token})
}

func isValidEmail(email string) bool {
	at := strings.LastIndex(email, "@")
	if at < 1 {
		return false
	}
	domainPart := email[at+1:]
	if len(domainPart) < 3 || !strings.Contains(domainPart, ".") {
		return false
	}
	return true
}
