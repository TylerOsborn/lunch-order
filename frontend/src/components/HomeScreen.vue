<template>
  <div class="home-screen">
    <h1>Meal Sharing App</h1>
    <div class="flex">
      <Button @click="$router.push('/give-meal')">Give a Meal</Button>
      <Button @click="$router.push('/receive-meal')">Receive a Meal</Button>
    </div>
  </div>
</template>

<script lang="ts">
import Button from "primevue/button";

export default {
  name: "HomeScreen",
  components: {
    Button,
  },
};
</script>

<style scoped>
.flex {
  display: flex;
  flex-direction: row;
  gap: 1rem;
  justify-content: center;
}
</style>
