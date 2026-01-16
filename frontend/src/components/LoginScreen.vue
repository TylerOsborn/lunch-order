<template>
  <div class="flex flex-col items-center justify-center min-h-screen bg-gray-100">
    <div class="p-8 bg-white rounded-lg shadow-md text-center">
      <h1 class="mb-6 text-2xl font-bold text-gray-800">Lunch Order Login</h1>
      <p class="mb-6 text-gray-600">Please sign in with your @impact.com email.</p>
      <Button asChild v-slot="slotProps">
        <a href="/auth/google/login" :class="slotProps.class" style="text-decoration: none">
          <i class="pi pi-google"></i>
          Sign in with Google
        </a>
      </Button>
    </div>
  </div>
</template>

<script setup lang="ts">
import Button from "primevue/button";
</script>
