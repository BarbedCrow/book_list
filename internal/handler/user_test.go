package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/BarbedCrow/book_list/internal/domain"
	"github.com/BarbedCrow/book_list/internal/handler"
)

type mockUserRegisterer struct {
	execute func(ctx context.Context, email, password string) (domain.User, error)
}

func (m *mockUserRegisterer) Execute(ctx context.Context, email, password string) (domain.User, error) {
	return m.execute(ctx, email, password)
}

type mockUserAuthenticator struct {
	execute func(ctx context.Context, email, password string) (string, error)
}

func (m *mockUserAuthenticator) Execute(ctx context.Context, email, password string) (string, error) {
	return m.execute(ctx, email, password)
}

func TestUserHandler_Register(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		registerer *mockUserRegisterer
		wantStatus int
	}{
		{
			name: "success",
			body: `{"email":"a@b.com","password":"secretpass"}`,
			registerer: &mockUserRegisterer{
				execute: func(_ context.Context, email, _ string) (domain.User, error) {
					return domain.User{ID: "1", Email: email}, nil
				},
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "bad JSON",
			body:       `{invalid`,
			registerer: &mockUserRegisterer{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "duplicate email",
			body: `{"email":"a@b.com","password":"secretpass"}`,
			registerer: &mockUserRegisterer{
				execute: func(_ context.Context, _, _ string) (domain.User, error) {
					return domain.User{}, domain.ErrDuplicateEmail
				},
			},
			wantStatus: http.StatusConflict,
		},
		{
			name:       "short password",
			body:       `{"email":"a@b.com","password":"short"}`,
			registerer: &mockUserRegisterer{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid email",
			body:       `{"email":"notanemail","password":"secretpass"}`,
			registerer: &mockUserRegisterer{},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handler.NewUserHandler(tt.registerer, nil)
			req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()
			h.Register(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("want status %d, got %d", tt.wantStatus, rec.Code)
			}
		})
	}
}

func TestUserHandler_Login(t *testing.T) {
	tests := []struct {
		name          string
		body          string
		authenticator *mockUserAuthenticator
		wantStatus    int
	}{
		{
			name: "success",
			body: `{"email":"a@b.com","password":"secret"}`,
			authenticator: &mockUserAuthenticator{
				execute: func(_ context.Context, _, _ string) (string, error) {
					return "token-123", nil
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name:          "bad JSON",
			body:          `{invalid`,
			authenticator: &mockUserAuthenticator{},
			wantStatus:    http.StatusBadRequest,
		},
		{
			name: "wrong credentials",
			body: `{"email":"a@b.com","password":"wrong"}`,
			authenticator: &mockUserAuthenticator{
				execute: func(_ context.Context, _, _ string) (string, error) {
					return "", domain.ErrWrongPassword
				},
			},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handler.NewUserHandler(nil, tt.authenticator)
			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()
			h.Login(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("want status %d, got %d", tt.wantStatus, rec.Code)
			}
		})
	}
}
