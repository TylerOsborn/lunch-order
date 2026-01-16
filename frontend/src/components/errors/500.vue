<script setup lang="ts">
import excuses from '../../excuses/excuses_500.json';
import ErrorPage from "./ErrorPage.vue";
</script>

<template>
<ErrorPage title="500 Internal Server Error 😵‍💫" :excuses="excuses"/>
</template>