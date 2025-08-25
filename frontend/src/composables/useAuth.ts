import { reactive } from 'vue';
import api from '../axios/axios';

interface User {
  id: number;
  name: string;
  email: string;
  role: string;
}

interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
}

const state = reactive<AuthState>({
  user: null,
  isAuthenticated: false,
  isLoading: true,
});

export const useAuth = () => {
  const checkAuth = async () => {
    try {
      state.isLoading = true;
      const response = await api.get('/Api/Auth/Profile');
      if (response.data.user) {
        state.user = response.data.user;
        state.isAuthenticated = true;
      }
    } catch (error) {
      state.user = null;
      state.isAuthenticated = false;
    } finally {
      state.isLoading = false;
    }
  };

  const login = async () => {
    try {
      const response = await api.get('/Api/Auth/Login');
      if (response.data.auth_url) {
        window.location.href = response.data.auth_url;
      }
    } catch (error) {
      console.error('Login failed:', error);
    }
  };

  const logout = async () => {
    try {
      await api.post('/Api/Auth/Logout');
      state.user = null;
      state.isAuthenticated = false;
    } catch (error) {
      console.error('Logout failed:', error);
    }
  };

  const isAdmin = () => {
    return state.user?.role === 'admin';
  };

  return {
    state,
    checkAuth,
    login,
    logout,
    isAdmin,
  };
};