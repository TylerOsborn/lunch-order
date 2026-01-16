import axios from 'axios';
import router from '../router';

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
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

export default api;
