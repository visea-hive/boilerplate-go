package crypto

import (
	"testing"
)

func TestNewToken_Format(t *testing.T) {
	emailToken, storedHash, err := NewToken()
	if err != nil {
		t.Fatalf("NewToken() error: %v", err)
	}

	// Both must be SHA-256 hex strings → 64 chars
	if len(emailToken) != 64 {
		t.Errorf("emailToken: expected length 64, got %d", len(emailToken))
	}
	if len(storedHash) != 64 {
		t.Errorf("storedHash: expected length 64, got %d", len(storedHash))
	}

	// They must be different (SHA-256 of different inputs)
	if emailToken == storedHash {
		t.Error("emailToken and storedHash must not be equal")
	}
}

func TestNewToken_IsUnique(t *testing.T) {
	t1, _, _ := NewToken()
	t2, _, _ := NewToken()
	if t1 == t2 {
		t.Error("expected two email tokens to be different")
	}
}

func TestHashToken_Deterministic(t *testing.T) {
	emailToken, storedHash, err := NewToken()
	if err != nil {
		t.Fatalf("NewToken() error: %v", err)
	}

	// Re-hashing the emailToken must reproduce the storedHash
	if got := HashToken(emailToken); got != storedHash {
		t.Errorf("HashToken(emailToken) = %s, want %s", got, storedHash)
	}
}

func TestHashToken_DifferentInput(t *testing.T) {
	h1 := HashToken("token-one")
	h2 := HashToken("token-two")
	if h1 == h2 {
		t.Error("different email tokens should produce different stored hashes")
	}
}
