package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/visea-hive/auth-core/pkg/messages"
)

// Claims is the full set of JWT claims for an access token.
// Standard fields (iss, sub, exp, iat) are embedded via RegisteredClaims
type Claims struct {
	jwt.RegisteredClaims

	// User identity
	UserUUID       string  `json:"user_uuid"`
	Name           string  `json:"name"`
	Email          string  `json:"email"`
	AvatarURL      *string `json:"avatar_url"`
	SessionID      string  `json:"sid"`
	ActiveOrgID    uint    `json:"active_org_id"`
	ActiveOrgName  string  `json:"active_org_name"`
	ActiveRoleID   uint    `json:"active_role_id"`
	ActiveRoleName string  `json:"active_role_name"`
}

// Manager mints and verifies JWT access tokens.
type Manager struct {
	secret []byte
	ttl    time.Duration
}

// New creates a Manager with the given HMAC secret and token TTL.
func New(secret string, ttl time.Duration) *Manager {
	return &Manager{
		secret: []byte(secret),
		ttl:    ttl,
	}
}

// Sign mints a signed HS256 access token from the given claims.
// iat and exp are set automatically; any value already in RegisteredClaims is respected.
func (m *Manager) Sign(claims Claims) (string, error) {
	now := time.Now()
	claims.IssuedAt = jwt.NewNumericDate(now)
	claims.ExpiresAt = jwt.NewNumericDate(now.Add(m.ttl))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("jwt: failed to sign token: %w", err)
	}
	return signed, nil
}

// Verify parses and validates a token string, returning the embedded Claims.
// Returns messages.ErrTokenExpired or messages.ErrTokenInvalid on failure.
func (m *Manager) Verify(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, messages.ErrTokenUnexpectedMethod
		}
		return m.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, messages.ErrTokenExpired
		}
		return nil, messages.ErrTokenInvalid
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, messages.ErrTokenInvalid
	}

	return claims, nil
}
