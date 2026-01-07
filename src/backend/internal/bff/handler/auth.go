package handler

import (
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"net/mail"
	"strings"

	"github.com/platepilot/backend/internal/bff/auth"
)

// AuthHandler handles authentication requests.
type AuthHandler struct {
	service *auth.Service
	logger  *slog.Logger
}

// NewAuthHandler creates a new auth handler.
func NewAuthHandler(service *auth.Service, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{service: service, logger: logger}
}

// Register handles POST /v1/auth/register.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	email, err := normalizeEmail(req.Email)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid email")
		return
	}

	if len(req.Password) < 8 {
		writeError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	tokens, err := h.service.Register(r.Context(), email, req.Password, strings.TrimSpace(req.DisplayName), userAgent(r), clientIP(r))
	if err != nil {
		switch err {
		case auth.ErrEmailAlreadyExists:
			writeError(w, http.StatusConflict, "email already registered")
		default:
			h.logger.Error("failed to register", "error", err)
			writeError(w, http.StatusInternalServerError, "failed to register")
		}
		return
	}

	writeJSON(w, http.StatusCreated, TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    tokens.ExpiresIn,
	})
}

// Login handles POST /v1/auth/login.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	email, err := normalizeEmail(req.Email)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid email")
		return
	}

	tokens, err := h.service.Login(r.Context(), email, req.Password, userAgent(r), clientIP(r))
	if err != nil {
		switch err {
		case auth.ErrInvalidCredentials:
			writeError(w, http.StatusUnauthorized, "invalid credentials")
		default:
			h.logger.Error("failed to login", "error", err)
			writeError(w, http.StatusInternalServerError, "failed to login")
		}
		return
	}

	writeJSON(w, http.StatusOK, TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    tokens.ExpiresIn,
	})
}

// Refresh handles POST /v1/auth/refresh.
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.RefreshToken) == "" {
		writeError(w, http.StatusBadRequest, "refresh token is required")
		return
	}

	tokens, err := h.service.Refresh(r.Context(), req.RefreshToken, userAgent(r), clientIP(r))
	if err != nil {
		switch err {
		case auth.ErrInvalidRefreshToken:
			writeError(w, http.StatusUnauthorized, "invalid refresh token")
		default:
			h.logger.Error("failed to refresh token", "error", err)
			writeError(w, http.StatusInternalServerError, "failed to refresh token")
		}
		return
	}

	writeJSON(w, http.StatusOK, TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    tokens.ExpiresIn,
	})
}

// Logout handles POST /v1/auth/logout.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.RefreshToken) == "" {
		writeError(w, http.StatusBadRequest, "refresh token is required")
		return
	}

	if err := h.service.Logout(r.Context(), req.RefreshToken); err != nil {
		switch err {
		case auth.ErrInvalidRefreshToken:
			writeError(w, http.StatusUnauthorized, "invalid refresh token")
		default:
			h.logger.Error("failed to logout", "error", err)
			writeError(w, http.StatusInternalServerError, "failed to logout")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RegisterRequest is the request body for registration.
type RegisterRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"displayName"`
}

// LoginRequest is the request body for login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RefreshRequest is the request body for token refresh/logout.
type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

// TokenResponse is the response body for token issuance.
type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	TokenType    string `json:"tokenType"`
	ExpiresIn    int64  `json:"expiresIn"`
}

func normalizeEmail(email string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if _, err := mail.ParseAddress(email); err != nil {
		return "", err
	}
	return email, nil
}

func userAgent(r *http.Request) string {
	return r.Header.Get("User-Agent")
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
