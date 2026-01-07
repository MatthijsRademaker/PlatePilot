// Token Service - Handles secure token storage and refresh
// Uses localStorage for persistent sessions (mobile-friendly)

const ACCESS_TOKEN_KEY = 'platepilot_access_token';
const REFRESH_TOKEN_KEY = 'platepilot_refresh_token';
const TOKEN_EXPIRY_KEY = 'platepilot_token_expiry';
const USER_KEY = 'platepilot_user';

export interface StoredTokens {
  accessToken: string;
  refreshToken: string;
  expiresAt: number;
}

export const tokenService = {
  // Store tokens after login/register
  setTokens(accessToken: string, refreshToken: string, expiresIn: number): void {
    const expiresAt = Date.now() + expiresIn * 1000;
    localStorage.setItem(ACCESS_TOKEN_KEY, accessToken);
    localStorage.setItem(REFRESH_TOKEN_KEY, refreshToken);
    localStorage.setItem(TOKEN_EXPIRY_KEY, expiresAt.toString());
  },

  // Get the current access token
  getAccessToken(): string | null {
    return localStorage.getItem(ACCESS_TOKEN_KEY);
  },

  // Get the refresh token
  getRefreshToken(): string | null {
    return localStorage.getItem(REFRESH_TOKEN_KEY);
  },

  // Get all stored tokens
  getTokens(): StoredTokens | null {
    const accessToken = localStorage.getItem(ACCESS_TOKEN_KEY);
    const refreshToken = localStorage.getItem(REFRESH_TOKEN_KEY);
    const expiresAt = localStorage.getItem(TOKEN_EXPIRY_KEY);

    if (!accessToken || !refreshToken || !expiresAt) {
      return null;
    }

    return {
      accessToken,
      refreshToken,
      expiresAt: parseInt(expiresAt, 10),
    };
  },

  // Check if the access token is expired (with 30s buffer)
  isTokenExpired(): boolean {
    const expiresAt = localStorage.getItem(TOKEN_EXPIRY_KEY);
    if (!expiresAt) return true;
    return Date.now() > parseInt(expiresAt, 10) - 30000;
  },

  // Check if we have a valid session
  hasValidSession(): boolean {
    const tokens = this.getTokens();
    if (!tokens) return false;
    // We have a session if we have a refresh token
    // Access token can be refreshed if expired
    return !!tokens.refreshToken;
  },

  // Store user info
  setUser(email: string, name?: string): void {
    localStorage.setItem(USER_KEY, JSON.stringify({ email, name }));
  },

  // Get stored user info
  getUser(): { email: string; name?: string } | null {
    const user = localStorage.getItem(USER_KEY);
    if (!user) return null;
    try {
      return JSON.parse(user);
    } catch {
      return null;
    }
  },

  // Clear all auth data on logout
  clearTokens(): void {
    localStorage.removeItem(ACCESS_TOKEN_KEY);
    localStorage.removeItem(REFRESH_TOKEN_KEY);
    localStorage.removeItem(TOKEN_EXPIRY_KEY);
    localStorage.removeItem(USER_KEY);
  },
};
