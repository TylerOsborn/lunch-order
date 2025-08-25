import { ref, computed } from 'vue';
import api from '../axios/axios';

const AUTH_STORAGE_KEY = 'lunch-order-admin-auth';

// Global state for authentication
const isAuthenticated = ref(false);

// Initialize from localStorage
const initializeAuth = () => {
  const stored = localStorage.getItem(AUTH_STORAGE_KEY);
  if (stored === 'authenticated') {
    isAuthenticated.value = true;
  }
};

// Initialize on module load
initializeAuth();

export const useAuth = () => {
  const login = async (password: string): Promise<boolean> => {
    try {
      const response = await api.post('/Api/Admin/Login', { password });

      if (response.status === 200 && response.data?.data?.authenticated) {
        isAuthenticated.value = true;
        localStorage.setItem(AUTH_STORAGE_KEY, 'authenticated');
        return true;
      }
      
      return false;
    } catch (error) {
      console.error('Authentication error:', error);
      return false;
    }
  };

  const logout = () => {
    isAuthenticated.value = false;
    localStorage.removeItem(AUTH_STORAGE_KEY);
  };

  const checkAuth = computed(() => isAuthenticated.value);

  return {
    isAuthenticated: checkAuth,
    login,
    logout,
  };
};