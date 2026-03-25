package auth_test

import (
	"testing"
	"time"

	"github.com/BarbedCrow/book_list/internal/auth"
)

func TestJWTProvider(t *testing.T) {
	p := auth.NewJWTProvider("test-secret", time.Hour)

	t.Run("generate and validate round-trip", func(t *testing.T) {
		token, err := p.Generate("user-123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		userID, err := p.Validate(token)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if userID != "user-123" {
			t.Fatalf("want user-123, got %s", userID)
		}
	})

	t.Run("expired token", func(t *testing.T) {
		expired := auth.NewJWTProvider("test-secret", -time.Hour)
		token, err := expired.Generate("user-123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		_, err = p.Validate(token)
		if err == nil {
			t.Fatal("expected error for expired token")
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		_, err := p.Validate("not-a-valid-token")
		if err == nil {
			t.Fatal("expected error for invalid token")
		}
	})
}
