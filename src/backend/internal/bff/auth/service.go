package auth

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/platepilot/backend/internal/common/auth"
)

// Tokens contains issued access and refresh tokens.
type Tokens struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

// Service handles auth workflows.
type Service struct {
	repo       *Repository
	tokens     *TokenService
	refreshTTL time.Duration
}

// NewService creates a new auth service.
func NewService(repo *Repository, tokens *TokenService, refreshTTL time.Duration) *Service {
	return &Service{
		repo:       repo,
		tokens:     tokens,
		refreshTTL: refreshTTL,
	}
}

// Register creates a user and issues tokens.
func (s *Service) Register(ctx context.Context, email, password, displayName, userAgent, ipAddress string) (Tokens, error) {
	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return Tokens{}, err
	}

	user, err := s.repo.CreateUser(ctx, email, displayName, passwordHash)
	if err != nil {
		return Tokens{}, err
	}

	return s.issueTokens(ctx, user.ID, userAgent, ipAddress)
}

// Login verifies credentials and issues tokens.
func (s *Service) Login(ctx context.Context, email, password, userAgent, ipAddress string) (Tokens, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return Tokens{}, ErrInvalidCredentials
	}

	hash, err := s.repo.GetPasswordHash(ctx, user.ID)
	if err != nil {
		return Tokens{}, ErrInvalidCredentials
	}

	if err := auth.CheckPassword(hash, password); err != nil {
		return Tokens{}, ErrInvalidCredentials
	}

	return s.issueTokens(ctx, user.ID, userAgent, ipAddress)
}

// Refresh validates a refresh token and rotates it.
func (s *Service) Refresh(ctx context.Context, refreshToken, userAgent, ipAddress string) (Tokens, error) {
	tokenHash := HashRefreshToken(refreshToken)
	record, err := s.repo.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		return Tokens{}, ErrInvalidRefreshToken
	}

	now := time.Now().UTC()
	if record.RevokedAt != nil || record.ExpiresAt.Before(now) {
		return Tokens{}, ErrInvalidRefreshToken
	}

	newRaw, newHash, err := GenerateRefreshToken()
	if err != nil {
		return Tokens{}, err
	}

	expiresAt := now.Add(s.refreshTTL)
	if err := s.repo.RotateRefreshToken(ctx, record.ID, record.UserID, newHash, expiresAt, userAgent, ipAddress); err != nil {
		return Tokens{}, err
	}

	accessToken, accessExpires, err := s.tokens.IssueAccessToken(record.UserID)
	if err != nil {
		return Tokens{}, err
	}

	return Tokens{
		AccessToken:  accessToken,
		RefreshToken: newRaw,
		ExpiresIn:    int64(accessExpires.Sub(now).Seconds()),
	}, nil
}

// Logout revokes the provided refresh token.
func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	tokenHash := HashRefreshToken(refreshToken)
	record, err := s.repo.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		return ErrInvalidRefreshToken
	}

	return s.repo.RevokeRefreshToken(ctx, record.ID)
}

func (s *Service) issueTokens(ctx context.Context, userID uuid.UUID, userAgent, ipAddress string) (Tokens, error) {
	now := time.Now().UTC()

	accessToken, accessExpires, err := s.tokens.IssueAccessToken(userID)
	if err != nil {
		return Tokens{}, err
	}

	refreshToken, refreshHash, err := GenerateRefreshToken()
	if err != nil {
		return Tokens{}, err
	}

	expiresAt := now.Add(s.refreshTTL)
	if err := s.repo.CreateRefreshToken(ctx, userID, refreshHash, expiresAt, userAgent, ipAddress); err != nil {
		return Tokens{}, err
	}

	return Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(accessExpires.Sub(now).Seconds()),
	}, nil
}
