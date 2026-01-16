<script setup lang="ts">
import excuses from '../../excuses/excuses_401.json';
import ErrorPage from "./ErrorPage.vue";
</script>

<template>
<ErrorPage title="401 Unauthorized 🤡" :excuses="excuses"/>
</template>