package jwt

import (
	"testing"
	"time"

	"github.com/visea-hive/auth-core/pkg/messages"
)

func newTestManager() *Manager {
	return New("test-secret-key-32-bytes-minimum!", 15*time.Minute)
}

func testClaims() Claims {
	name := "John Doe"
	avatar := "https://example.com/avatar.jpg"
	orgID := uint(42)
	return Claims{
		UserUUID:       "01944b1f-0000-7000-8000-000000000001",
		Name:           name,
		Email:          "john@example.com",
		AvatarURL:      &avatar,
		SessionID:      "session-abc-123",
		ActiveOrgID:    orgID,
		ActiveOrgName:  "Acme Corp",
		ActiveRoleID:   10,
		ActiveRoleName: "admin",
	}
}

func TestSign_ProducesNonEmptyToken(t *testing.T) {
	m := newTestManager()
	token, err := m.Sign(testClaims())
	if err != nil {
		t.Fatalf("Sign() error: %v", err)
	}
	if token == "" {
		t.Error("expected non-empty token string")
	}
}

func TestVerify_ValidToken(t *testing.T) {
	m := newTestManager()
	original := testClaims()

	token, err := m.Sign(original)
	if err != nil {
		t.Fatalf("Sign() error: %v", err)
	}

	claims, err := m.Verify(token)
	if err != nil {
		t.Fatalf("Verify() error: %v", err)
	}

	if claims.UserUUID != original.UserUUID {
		t.Errorf("UserUUID: got %s, want %s", claims.UserUUID, original.UserUUID)
	}
	if claims.Email != original.Email {
		t.Errorf("Email: got %s, want %s", claims.Email, original.Email)
	}
	if claims.SessionID != original.SessionID {
		t.Errorf("SessionID: got %s, want %s", claims.SessionID, original.SessionID)
	}
	if claims.ActiveRoleName != original.ActiveRoleName {
		t.Errorf("ActiveRoleName: got %s, want %s", claims.ActiveRoleName, original.ActiveRoleName)
	}
}

func TestVerify_ExpiredToken(t *testing.T) {
	m := New("test-secret-key-32-bytes-minimum!", -1*time.Second) // already expired
	token, err := m.Sign(testClaims())
	if err != nil {
		t.Fatalf("Sign() error: %v", err)
	}

	_, err = m.Verify(token)
	if err != messages.ErrTokenExpired {
		t.Errorf("expected ErrTokenExpired, got: %v", err)
	}
}

func TestVerify_TamperedToken(t *testing.T) {
	m := newTestManager()
	token, _ := m.Sign(testClaims())

	_, err := m.Verify(token + "tampered")
	if err != messages.ErrTokenInvalid {
		t.Errorf("expected ErrTokenInvalid for tampered token, got: %v", err)
	}
}

func TestVerify_WrongSecret(t *testing.T) {
	signer := New("correct-secret-key-32-bytes-min!", 15*time.Minute)
	verifier := New("wrong-secret-key-32-bytes-minnnn", 15*time.Minute)

	token, _ := signer.Sign(testClaims())

	_, err := verifier.Verify(token)
	if err != messages.ErrTokenInvalid {
		t.Errorf("expected ErrTokenInvalid for wrong secret, got: %v", err)
	}
}
