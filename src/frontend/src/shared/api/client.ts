const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

interface RequestOptions extends RequestInit {
  params?: Record<string, string | number>;
}

interface ApiError {
  message: string;
  status: number;
}

class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  private buildUrl(path: string, params?: Record<string, string | number>): string {
    const url = new URL(path, this.baseUrl);
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        url.searchParams.append(key, String(value));
      });
    }
    return url.toString();
  }

  private async handleResponse<T>(response: Response): Promise<T> {
    if (!response.ok) {
      const error = new Error(`HTTP error: ${response.statusText}`) as Error & { status: number };
      error.status = response.status;
      throw error;
    }
    return response.json() as Promise<T>;
  }

  async get<T>(path: string, options?: RequestOptions): Promise<T> {
    const url = this.buildUrl(path, options?.params);
    const response = await fetch(url, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
      ...options,
    });
    return this.handleResponse<T>(response);
  }

  async post<T, D = unknown>(path: string, data?: D, options?: RequestOptions): Promise<T> {
    const url = this.buildUrl(path, options?.params);
    const fetchOptions: RequestInit = {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
      ...options,
    };
    if (data !== undefined) {
      fetchOptions.body = JSON.stringify(data);
    }
    const response = await fetch(url, fetchOptions);
    return this.handleResponse<T>(response);
  }

  async put<T, D = unknown>(path: string, data?: D, options?: RequestOptions): Promise<T> {
    const url = this.buildUrl(path, options?.params);
    const fetchOptions: RequestInit = {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
      ...options,
    };
    if (data !== undefined) {
      fetchOptions.body = JSON.stringify(data);
    }
    const response = await fetch(url, fetchOptions);
    return this.handleResponse<T>(response);
  }

  async delete<T>(path: string, options?: RequestOptions): Promise<T> {
    const url = this.buildUrl(path, options?.params);
    const response = await fetch(url, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
      ...options,
    });
    return this.handleResponse<T>(response);
  }
}

export const apiClient = new ApiClient(API_BASE_URL);
export type { ApiError, RequestOptions };
