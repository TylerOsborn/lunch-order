import { ref, computed } from 'vue';

const ADMIN_PASSWORD = 'admin123'; // Simple password for basic auth
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
  const login = (password: string): boolean => {
    if (password === ADMIN_PASSWORD) {
      isAuthenticated.value = true;
      localStorage.setItem(AUTH_STORAGE_KEY, 'authenticated');
      return true;
    }
    return false;
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