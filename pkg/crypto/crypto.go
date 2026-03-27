package crypto

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/visea-hive/auth-core/pkg/messages"
	"golang.org/x/crypto/argon2"
)

// Argon2id parameters.
const (
	argonMemory  uint32 = 65536 // 64 MiB
	argonTime    uint32 = 3
	argonThreads uint8  = 4
	argonKeyLen  uint32 = 32
	saltLen             = 16
)

// Hasher hashes and verifies passwords using argon2id with a pepper.
// The pepper (secret key from env) is HMAC-SHA256'd onto the password
// before it is fed into argon2id, so a DB dump alone is never sufficient
// to crack passwords.
type Hasher struct {
	pepper []byte
}

// New returns a Hasher initialised with the given pepper string.
// An empty pepper is allowed but strongly discouraged in production.
func New(pepper string) *Hasher {
	return &Hasher{pepper: []byte(pepper)}
}

// Hash derives a secure hash for password and returns a self-contained
// encoded string in the format:
//
//	$argon2id$v=19$m=65536,t=3,p=4$<base64-salt>$<base64-hash>
//
// The salt is randomly generated on every call, so the same password
// produces a different encoded string each time.
func (h *Hasher) Hash(password string) (string, error) {
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("crypto: failed to generate salt: %w", err)
	}

	peppered := h.applyPepper(password)
	hash := argon2.IDKey(peppered, salt, argonTime, argonMemory, argonThreads, argonKeyLen)

	encoded := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		argonMemory,
		argonTime,
		argonThreads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)

	return encoded, nil
}

// Verify checks whether password matches the encoded hash.
// Returns (true, nil) on match, (false, nil) on mismatch, and
// (false, err) if the encoded hash is malformed.
func (h *Hasher) Verify(password, encodedHash string) (bool, error) {
	salt, hash, err := decode(encodedHash)
	if err != nil {
		return false, err
	}

	peppered := h.applyPepper(password)
	candidate := argon2.IDKey(peppered, salt, argonTime, argonMemory, argonThreads, argonKeyLen)

	if subtle.ConstantTimeCompare(hash, candidate) == 1 {
		return true, nil
	}

	return false, nil
}

// applyPepper applies HMAC-SHA256 with the configured pepper to the password.
func (h *Hasher) applyPepper(password string) []byte {
	mac := hmac.New(sha256.New, h.pepper)
	mac.Write([]byte(password))
	return mac.Sum(nil)
}

// decode parses an encoded argon2id hash string and returns the salt and hash bytes.
func decode(encodedHash string) (salt, hash []byte, err error) {
	parts := strings.Split(encodedHash, "$")
	// expected: ["", "argon2id", "v=19", "m=65536,t=3,p=4", "<salt>", "<hash>"]
	if len(parts) != 6 {
		return nil, nil, messages.ErrHashInvalidFormat
	}

	if parts[1] != "argon2id" {
		return nil, nil, messages.ErrHashInvalidFormat
	}

	var version int
	if _, err = fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return nil, nil, messages.ErrHashInvalidFormat
	}
	if version != argon2.Version {
		return nil, nil, messages.ErrHashIncompatibleVersion
	}

	salt, err = base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, messages.ErrHashInvalidFormat
	}

	hash, err = base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, messages.ErrHashInvalidFormat
	}

	return salt, hash, nil
}
