// Auth Store - Pinia store for authentication state
import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import type { LoginRequest, RegisterRequest, AuthUser } from '../types/auth';
import * as authApi from '../api/authApi';
import { tokenService } from '../services/tokenService';

export const useAuthStore = defineStore('auth', () => {
  const router = useRouter();

  // State
  const user = ref<AuthUser | null>(null);
  const isLoading = ref(false);
  const error = ref<string | null>(null);
  const isInitialized = ref(false);

  // Getters
  const isAuthenticated = computed(() => !!user.value);

  // Initialize auth state from stored tokens
  async function initialize(): Promise<void> {
    if (isInitialized.value) return;

    const storedUser = tokenService.getUser();
    if (storedUser && tokenService.hasValidSession()) {
      user.value = storedUser;

      // If access token is expired, try to refresh it
      if (tokenService.isTokenExpired()) {
        try {
          await refreshTokens();
        } catch {
          // Refresh failed, clear session
          await logout();
        }
      }
    }

    isInitialized.value = true;
  }

  // Login
  async function login(credentials: LoginRequest): Promise<boolean> {
    isLoading.value = true;
    error.value = null;

    try {
      const response = await authApi.login(credentials);

      // Store tokens
      tokenService.setTokens(
        response.accessToken,
        response.refreshToken,
        response.expiresIn,
      );

      // Store user info (extract from email)
      const userName = credentials.email.split('@')[0];
      tokenService.setUser(credentials.email, userName);
      user.value = { email: credentials.email, ...(userName ? { name: userName } : {}) };

      return true;
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Login failed';
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  // Register
  async function register(data: RegisterRequest): Promise<boolean> {
    isLoading.value = true;
    error.value = null;

    try {
      const response = await authApi.register(data);

      // Store tokens
      tokenService.setTokens(
        response.accessToken,
        response.refreshToken,
        response.expiresIn,
      );

      // Store user info
      const userName = data.name ?? data.email.split('@')[0];
      tokenService.setUser(data.email, userName);
      user.value = { email: data.email, ...(userName ? { name: userName } : {}) };

      return true;
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Registration failed';
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  // Refresh tokens
  async function refreshTokens(): Promise<void> {
    const refreshToken = tokenService.getRefreshToken();
    if (!refreshToken) {
      throw new Error('No refresh token available');
    }

    const response = await authApi.refreshToken(refreshToken);

    tokenService.setTokens(
      response.accessToken,
      response.refreshToken,
      response.expiresIn,
    );
  }

  // Logout
  async function logout(): Promise<void> {
    const refreshToken = tokenService.getRefreshToken();

    // Try to notify server, but don't block on failure
    if (refreshToken) {
      try {
        await authApi.logout(refreshToken);
      } catch {
        // Ignore logout errors
      }
    }

    // Clear local state
    tokenService.clearTokens();
    user.value = null;
    error.value = null;

    // Redirect to login
    void router.push({ name: 'login' });
  }

  // Clear error
  function clearError(): void {
    error.value = null;
  }

  return {
    // State
    user,
    isLoading,
    error,
    isInitialized,
    // Getters
    isAuthenticated,
    // Actions
    initialize,
    login,
    register,
    refreshTokens,
    logout,
    clearError,
  };
});
