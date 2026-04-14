import axios, { type AxiosRequestConfig } from 'axios';
import router from '../router';
import type { ApiResult } from '../api/types';

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
  withCredentials: true,
});

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response) {
      switch (error.response.status) {
        case 401:
          router.push('/401');
          break;
        case 403:
          router.push('/403');
          break;
      }
    }
    return Promise.reject(error);
  }
);

export async function apiGet<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
  const response = await api.get<ApiResult<T>>(url, config);
  return response.data.data;
}

export async function apiPost<T = void>(url: string, body?: unknown, config?: AxiosRequestConfig): Promise<T> {
  const response = await api.post<ApiResult<T>>(url, body, config);
  return response.data.data;
}

export default api;
