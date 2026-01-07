package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TokenService issues and validates access tokens.
type TokenService struct {
	secret    []byte
	issuer    string
	accessTTL time.Duration
}

// NewTokenService creates a new TokenService.
func NewTokenService(secret, issuer string, accessTTL time.Duration) *TokenService {
	return &TokenService{
		secret:    []byte(secret),
		issuer:    issuer,
		accessTTL: accessTTL,
	}
}

// IssueAccessToken issues a signed JWT access token for the given user.
func (s *TokenService) IssueAccessToken(userID uuid.UUID) (string, time.Time, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(s.accessTTL)

	claims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		Issuer:    s.issuer,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("sign token: %w", err)
	}

	return signed, expiresAt, nil
}

// ParseAccessToken validates a JWT access token and returns the user ID.
func (s *TokenService) ParseAccessToken(tokenString string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	parsed, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %T", token.Method)
		}
		return s.secret, nil
	})
	if err != nil || !parsed.Valid {
		return uuid.Nil, ErrInvalidCredentials
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, ErrInvalidCredentials
	}

	return userID, nil
}

// GenerateRefreshToken returns a raw token and its hash for storage.
func GenerateRefreshToken() (string, string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", "", fmt.Errorf("generate refresh token: %w", err)
	}

	raw := base64.RawURLEncoding.EncodeToString(buf)
	hash := sha256.Sum256([]byte(raw))
	return raw, hex.EncodeToString(hash[:]), nil
}

// HashRefreshToken hashes a raw refresh token for lookup.
func HashRefreshToken(raw string) string {
	hash := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(hash[:])
}
