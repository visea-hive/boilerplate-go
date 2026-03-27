package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// GenerateRandomBytes returns n cryptographically secure random bytes.
// Use this whenever you need raw entropy (tokens, salts, nonces, etc.).
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("crypto: failed to generate %d random bytes: %w", n, err)
	}
	return b, nil
}

// SHA256Hex returns the SHA-256 hex digest of the given byte slice.
// Use this to hash any arbitrary value for safe storage or comparison.
func SHA256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// NewToken generates a verification token using a double-hash chain:
//
//		random bytes → SHA-256 → emailToken (hex) → SHA-256 → storedHash (hex)
//
//	  - emailToken  is embedded in the link sent to the user.
//	  - storedHash  is stored in the DB (e.g. email_verifications.token_hash).
//
// On verification, call HashToken(tokenFromURL) to reproduce storedHash for the DB lookup.
func NewToken() (emailToken string, storedHash string, err error) {
	raw, err := GenerateRandomBytes(32)
	if err != nil {
		return "", "", err
	}

	emailToken = SHA256Hex(raw)
	storedHash = HashToken(emailToken)
	return emailToken, storedHash, nil
}

// HashToken returns the SHA-256 hex digest of an email token string.
// Call this on the token received from the URL to produce the storedHash for DB lookup.
func HashToken(emailToken string) string {
	return SHA256Hex([]byte(emailToken))
}
