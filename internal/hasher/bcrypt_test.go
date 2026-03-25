package hasher_test

import (
	"testing"

	"github.com/BarbedCrow/book_list/internal/hasher"
)

func TestBcryptHasher(t *testing.T) {
	h := hasher.NewBcryptHasher(4) // low cost for fast tests

	t.Run("hash and verify success", func(t *testing.T) {
		hash, err := h.Hash("my-password")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := h.Verify(hash, "my-password"); err != nil {
			t.Fatalf("verify should succeed: %v", err)
		}
	})

	t.Run("verify wrong password", func(t *testing.T) {
		hash, err := h.Hash("my-password")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := h.Verify(hash, "wrong-password"); err == nil {
			t.Fatal("verify should fail for wrong password")
		}
	})
}
