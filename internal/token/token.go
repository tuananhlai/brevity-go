package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	DefaultAccessTokenExpiry = time.Hour * 24 * 30
)

// AccessTokenIssuer handles JWT token generation and verification with a shared secret.
type AccessTokenIssuer struct {
	secret string
	expiry time.Duration
}

// NewIssuer creates a new token manager with the given secret and default expiry.
func NewIssuer(secret string) *AccessTokenIssuer {
	return &AccessTokenIssuer{
		secret: secret,
		expiry: DefaultAccessTokenExpiry,
	}
}

// Issue creates a signed JWT token for the given user ID.
func (m *AccessTokenIssuer) Issue(userID string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(m.expiry).Unix(),
	})

	return t.SignedString([]byte(m.secret))
}

// Verify parses and validates a JWT token, returning the user ID (subject).
func (m *AccessTokenIssuer) Verify(accessToken string) (string, error) {
	t, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
		return []byte(m.secret), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("missing token subject")
	}

	return subject, nil
}
