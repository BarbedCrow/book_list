package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BarbedCrow/book_list/internal/handler"
)

type mockTokenValidator struct {
	validate func(token string) (string, error)
}

func (m *mockTokenValidator) Validate(token string) (string, error) {
	return m.validate(token)
}

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name       string
		authHeader string
		validator  *mockTokenValidator
		wantStatus int
	}{
		{
			name:       "valid token",
			authHeader: "Bearer valid-token",
			validator: &mockTokenValidator{
				validate: func(_ string) (string, error) { return "user-1", nil },
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing header",
			authHeader: "",
			validator:  &mockTokenValidator{},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "invalid token",
			authHeader: "Bearer bad-token",
			validator: &mockTokenValidator{
				validate: func(_ string) (string, error) { return "", errors.New("invalid") },
			},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			mw := handler.AuthMiddleware(tt.validator)(next)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rec := httptest.NewRecorder()
			mw.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("want status %d, got %d", tt.wantStatus, rec.Code)
			}
		})
	}
}
