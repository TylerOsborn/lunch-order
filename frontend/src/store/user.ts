import { reactive } from 'vue';

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
      const response = await fetch('/Api/Me');
      if (response.ok) {
        this.user = await response.json();
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
          await fetch('/auth/logout', { method: 'POST' });
      } catch (e) {
          console.error(e);
      } finally {
          this.user = null;
          this.isAuthenticated = false;
          window.location.href = '/login';
      }
  }
});
