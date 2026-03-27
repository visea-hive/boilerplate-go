package crypto

import (
	"strings"
	"testing"
)

const testPepper = "test-pepper-secret"

func newTestHasher() *Hasher {
	return New(testPepper)
}

func TestHash_ProducesValidFormat(t *testing.T) {
	h := newTestHasher()
	encoded, err := h.Hash("MySecret123!")
	if err != nil {
		t.Fatalf("Hash() error: %v", err)
	}

	// Must start with $argon2id$
	if !strings.HasPrefix(encoded, "$argon2id$") {
		t.Errorf("expected encoded hash to start with $argon2id$, got: %s", encoded)
	}

	// Must have 6 parts separated by $
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 {
		t.Errorf("expected 6 parts in encoded hash, got %d: %s", len(parts), encoded)
	}
}

func TestHash_IsNonDeterministic(t *testing.T) {
	h := newTestHasher()
	password := "SamePassword99!"

	h1, err := h.Hash(password)
	if err != nil {
		t.Fatalf("Hash() first call error: %v", err)
	}
	h2, err := h.Hash(password)
	if err != nil {
		t.Fatalf("Hash() second call error: %v", err)
	}

	if h1 == h2 {
		t.Error("expected two hashes of the same password to differ (different salts)")
	}
}

func TestVerify_CorrectPassword(t *testing.T) {
	h := newTestHasher()
	password := "Correct$Horse9"

	encoded, err := h.Hash(password)
	if err != nil {
		t.Fatalf("Hash() error: %v", err)
	}

	ok, err := h.Verify(password, encoded)
	if err != nil {
		t.Fatalf("Verify() error: %v", err)
	}
	if !ok {
		t.Error("Verify() returned false for correct password")
	}
}

func TestVerify_WrongPassword(t *testing.T) {
	h := newTestHasher()

	encoded, err := h.Hash("RightPassword1!")
	if err != nil {
		t.Fatalf("Hash() error: %v", err)
	}

	ok, err := h.Verify("WrongPassword1!", encoded)
	if err != nil {
		t.Fatalf("Verify() unexpected error: %v", err)
	}
	if ok {
		t.Error("Verify() returned true for wrong password")
	}
}

func TestVerify_WrongPepper(t *testing.T) {
	h1 := New("pepper-one")
	h2 := New("pepper-two")

	encoded, err := h1.Hash("Password123!")
	if err != nil {
		t.Fatalf("Hash() error: %v", err)
	}

	ok, err := h2.Verify("Password123!", encoded)
	if err != nil {
		t.Fatalf("Verify() unexpected error: %v", err)
	}
	if ok {
		t.Error("Verify() returned true when verified with a different pepper")
	}
}

func TestVerify_TamperedHash(t *testing.T) {
	h := newTestHasher()

	_, err := h.Verify("password", "not-a-valid-hash")
	if err == nil {
		t.Error("Verify() expected error for malformed hash, got nil")
	}
}

func TestVerify_TamperedHashMissingParts(t *testing.T) {
	h := newTestHasher()

	_, err := h.Verify("password", "$argon2id$v=19$m=65536,t=3,p=4$onlysalt")
	if err == nil {
		t.Error("Verify() expected error for hash with missing parts, got nil")
	}
}
