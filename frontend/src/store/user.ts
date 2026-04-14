import { reactive } from 'vue';
import api from '../axios/axios';

interface User {
  id: number;
  name: string;
  email: string;
  firstName: string;
  lastName: string;
  avatarUrl: string;
  isAdmin: boolean;
}

export const userStore = reactive({
  user: null as User | null,
  isAuthenticated: false,
  async fetchUser() {
    try {
      const response = await api.get('/Api/Me', { validateStatus: () => true });
      if (response.status === 200) {
        this.user = response.data;
        this.isAuthenticated = true;
      } else {
        this.user = null;
        this.isAuthenticated = false;
      }
    } catch (error) {
      this.user = null;
      this.isAuthenticated = false;
    }
  },
  async logout() {
    try {
      await api.post('/auth/logout', undefined, { validateStatus: () => true });
    } catch (e) {
      console.error(e);
    } finally {
      this.user = null;
      this.isAuthenticated = false;
      window.location.href = '/login';
    }
  }
});
