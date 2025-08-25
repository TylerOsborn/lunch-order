<template>
  <div>
    <!-- Navigation Bar -->
    <div class="navbar" v-if="$route.path !== '/login'">
      <div class="nav-brand">
        <h3>Meal Share</h3>
      </div>
      <div class="nav-links">
        <router-link to="/" class="nav-link">Home</router-link>
        <router-link 
          v-if="authState.isAuthenticated" 
          to="/give-meal" 
          class="nav-link"
        >
          Give Meal
        </router-link>
        <router-link 
          v-if="authState.isAuthenticated" 
          to="/receive-meal" 
          class="nav-link"
        >
          Receive Meal
        </router-link>
        <router-link 
          v-if="authState.isAuthenticated && isAdmin()" 
          to="/admin" 
          class="nav-link"
        >
          Admin
        </router-link>
      </div>
      <div class="nav-auth">
        <div v-if="authState.isLoading" class="loading">Loading...</div>
        <div v-else-if="authState.isAuthenticated" class="user-info">
          <span class="user-name">{{ authState.user?.name }}</span>
          <Button @click="handleLogout" severity="secondary" size="small">
            Logout
          </Button>
        </div>
        <div v-else>
          <Button @click="$router.push('/login')" severity="info" size="small">
            Login
          </Button>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="main-content">
      <Toast />
      <router-view></router-view>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { useAuth } from './composables/useAuth';
import { useRouter } from 'vue-router';
import Toast from 'primevue/toast';
import Button from 'primevue/button';

const { state: authState, checkAuth, logout, isAdmin } = useAuth();
const router = useRouter();

onMounted(async () => {
  await checkAuth();
});

const handleLogout = async () => {
  await logout();
  router.push('/');
};
</script>

<style scoped>
.navbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 2rem;
  background-color: #fff;
  border-bottom: 1px solid #e0e0e0;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.nav-brand h3 {
  margin: 0;
  color: #2c3e50;
}

.nav-links {
  display: flex;
  gap: 1rem;
}

.nav-link {
  text-decoration: none;
  color: #666;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  transition: all 0.2s;
}

.nav-link:hover {
  background-color: #f5f5f5;
  color: #2c3e50;
}

.nav-link.router-link-active {
  background-color: #3498db;
  color: white;
}

.nav-auth {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.user-name {
  font-weight: 500;
  color: #2c3e50;
}

.main-content {
  min-height: calc(100vh - 80px);
}

.loading {
  color: #666;
}
</style>
