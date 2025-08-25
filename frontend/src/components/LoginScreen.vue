<template>
  <div class="login-screen">
    <Card class="login-card">
      <template #title>
        <h2>Sign In to Meal Share</h2>
      </template>
      <template #content>
        <div class="login-content">
          <p>Please sign in with your Google account to continue</p>
          <Button @click="handleLogin" severity="info" class="google-login-btn">
            <i class="pi pi-google"></i>
            Sign in with Google
          </Button>
        </div>
      </template>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { useAuth } from '../composables/useAuth';
import Card from 'primevue/card';
import Button from 'primevue/button';

const { login } = useAuth();

const handleLogin = async () => {
  await login();
};
</script>

<style scoped>
.login-screen {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f5f5f5;
}

.login-card {
  width: 100%;
  max-width: 400px;
  margin: 2rem;
}

.login-content {
  text-align: center;
  padding: 1rem;
}

.google-login-btn {
  margin-top: 1rem;
  width: 100%;
}

.google-login-btn i {
  margin-right: 0.5rem;
}
</style>