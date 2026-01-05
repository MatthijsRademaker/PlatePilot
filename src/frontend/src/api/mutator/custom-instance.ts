const apiHost = import.meta.env.VITE_API_URL || 'http://localhost:8080';
const baseURL = `${apiHost}/v1`;

export interface ApiError {
  message: string;
  status: number;
}

export const customInstance = async <T>(config: {
  url: string;
  method: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH';
  params?: Record<string, unknown> | undefined;
  data?: unknown;
  headers?: Record<string, string>;
}): Promise<T> => {
  const { url, method, params, data, headers } = config;

  let targetUrl = `${baseURL}${url}`;

  if (params) {
    const searchParams = new URLSearchParams();
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        if (typeof value === 'object') {
          searchParams.append(key, JSON.stringify(value));
        } else if (typeof value === 'string' || typeof value === 'number' || typeof value === 'boolean') {
          searchParams.append(key, String(value));
        }
      }
    });
    const queryString = searchParams.toString();
    if (queryString) {
      targetUrl += `?${queryString}`;
    }
  }

  const fetchOptions: RequestInit = {
    method,
    headers: {
      'Content-Type': 'application/json',
      ...headers,
    },
  };

  if (data !== undefined) {
    fetchOptions.body = JSON.stringify(data);
  }

  const response = await fetch(targetUrl, fetchOptions);

  if (!response.ok) {
    const errorBody = await response.json().catch(() => ({ error: response.statusText }));
    const apiError = new Error(errorBody.error || response.statusText) as Error & ApiError;
    apiError.status = response.status;
    throw apiError;
  }

  return response.json();
};

export default customInstance;
