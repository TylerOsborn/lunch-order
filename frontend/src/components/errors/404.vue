<script setup lang="ts">
import excuses from '../../excuses/excuses_404.json';
import ErrorPage from "./ErrorPage.vue";
</script>

<template>
  <ErrorPage title="404 Page Not Found 🚫" :excuses="excuses"/>
</template>