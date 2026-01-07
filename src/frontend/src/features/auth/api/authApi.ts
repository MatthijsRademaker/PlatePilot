// Auth API for PlatePilot

import type { LoginRequest, RegisterRequest, AuthResponse } from '../types/auth';

const apiHost = import.meta.env.VITE_API_URL || 'http://localhost:8080';
const baseURL = `${apiHost}/v1/auth`;

async function authFetch<T>(
  endpoint: string,
  options: RequestInit = {},
): Promise<T> {
  const response = await fetch(`${baseURL}${endpoint}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
  });

  if (!response.ok) {
    const errorBody = await response.json().catch(() => ({ error: response.statusText }));
    throw new Error(errorBody.error || errorBody.message || 'Authentication failed');
  }

  return response.json();
}

export async function login(credentials: LoginRequest): Promise<AuthResponse> {
  return authFetch<AuthResponse>('/login', {
    method: 'POST',
    body: JSON.stringify(credentials),
  });
}

export async function register(data: RegisterRequest): Promise<AuthResponse> {
  return authFetch<AuthResponse>('/register', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export async function refreshToken(refreshToken: string): Promise<AuthResponse> {
  return authFetch<AuthResponse>('/refresh', {
    method: 'POST',
    body: JSON.stringify({ refreshToken }),
  });
}

export async function logout(refreshToken: string): Promise<void> {
  await authFetch<void>('/logout', {
    method: 'POST',
    body: JSON.stringify({ refreshToken }),
  });
}
