<script setup lang="ts">
import { ref } from 'vue';
import { useAuth } from '../composables/useAuth';
import { useToast } from 'primevue/usetoast';

import Card from 'primevue/card';
import Password from 'primevue/password';
import Button from 'primevue/button';

const emit = defineEmits<{
  login: [];
}>();

const { login } = useAuth();
const toast = useToast();

const password = ref('');
const isLoading = ref(false);

const handleLogin = async () => {
  if (!password.value.trim()) {
    toast.add({ 
      severity: 'error', 
      summary: 'Error', 
      detail: 'Please enter a password',
      life: 3000 
    });
    return;
  }

  isLoading.value = true;
  
  const success = login(password.value);
  
  if (success) {
    toast.add({ 
      severity: 'success', 
      summary: 'Success', 
      detail: 'Login successful',
      life: 3000 
    });
    emit('login');
  } else {
    toast.add({ 
      severity: 'error', 
      summary: 'Error', 
      detail: 'Invalid password',
      life: 3000 
    });
  }
  
  isLoading.value = false;
};

const handleKeyPress = (event: KeyboardEvent) => {
  if (event.key === 'Enter') {
    handleLogin();
  }
};
</script>

<template>
  <div class="login-container">
    <Card class="login-card">
      <template #title>
        <h2>Admin Login</h2>
      </template>
      <template #content>
        <div class="login-form">
          <div class="field">
            <label for="password">Password</label>
            <Password 
              id="password"
              v-model="password" 
              :feedback="false"
              :toggleMask="true"
              placeholder="Enter admin password"
              class="full-width"
              @keypress="handleKeyPress"
              :disabled="isLoading"
            />
          </div>
          <Button 
            type="submit" 
            label="Login" 
            class="login-button"
            @click="handleLogin"
            :loading="isLoading"
            :disabled="isLoading"
          />
        </div>
      </template>
    </Card>
  </div>
</template>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: calc(100vh - 4rem);
  width: 100%;
}

.login-card {
  width: 400px;
  max-width: 90%;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.field label {
  font-weight: 600;
}

.full-width {
  width: 100%;
}

.login-button {
  width: 100%;
  margin-top: 0.5rem;
}
</style>