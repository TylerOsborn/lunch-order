<template>
  <div class="home-screen">
    <h1>Meal Sharing App</h1>
    
    <div v-if="authState.isLoading" class="loading">
      Loading...
    </div>
    
    <div v-else-if="!authState.isAuthenticated" class="welcome-section">
      <p>Welcome to Meal Share! Please sign in to start sharing meals.</p>
      <Button @click="$router.push('/login')" severity="info" size="large">
        Sign In
      </Button>
    </div>
    
    <div v-else class="authenticated-section">
      <p>Welcome back, {{ authState.user?.name }}!</p>
      <div class="flex">
        <Button @click="$router.push('/give-meal')">Give a Meal</Button>
        <Button @click="$router.push('/receive-meal')">Receive a Meal</Button>
        <Button 
          v-if="isAdmin()" 
          @click="$router.push('/admin')" 
          severity="secondary"
        >
          Admin Panel
        </Button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { useAuth } from '../composables/useAuth';
import Button from 'primevue/button';

const { state: authState, checkAuth, isAdmin } = useAuth();

onMounted(async () => {
  if (authState.isLoading) {
    await checkAuth();
  }
});
</script>

<style scoped>
.home-screen {
  padding: 2rem;
  text-align: center;
  max-width: 800px;
  margin: 0 auto;
}

.flex {
  display: flex;
  flex-direction: row;
  gap: 1rem;
  justify-content: center;
  margin-top: 2rem;
}

.welcome-section, .authenticated-section {
  margin-top: 2rem;
}

.welcome-section p, .authenticated-section p {
  font-size: 1.2rem;
  margin-bottom: 2rem;
  color: #666;
}

.loading {
  font-size: 1.2rem;
  color: #666;
  margin-top: 2rem;
}

@media (max-width: 768px) {
  .flex {
    flex-direction: column;
    align-items: center;
  }
}
</style>
