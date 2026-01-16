<script setup lang="ts">
import { ref, onMounted } from 'vue';

const props = defineProps<{
  title: string;
  excuses: string[]
}>();

const excuse = ref('');
onMounted(() => {
  const randomIndex = Math.floor(Math.random() * props.excuses.length);
  excuse.value = props.excuses[randomIndex];
});
</script>

<template>
  <div class="flex flex-col items-center justify-center min-h-screen bg-gray-100">
    <div class="p-8 bg-white rounded-lg shadow-md text-center">
      <h1 class="mb-6 text-2xl font-bold text-gray-800">{{ title }}</h1>
      <p>
        {{ excuse }}
      </p>
    </div>
  </div>
</template>

<style scoped>

</style>