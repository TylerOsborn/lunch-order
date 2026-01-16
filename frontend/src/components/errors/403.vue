<script setup lang="ts">
import excuses from '../../excuses/excuses_403.json';
import ErrorPage from "./ErrorPage.vue";
</script>

<template>
<ErrorPage title="403 Forbidden 🥶" :excuses="excuses"/>
</template>