<template>
  <Toast />
  <router-view></router-view>
</template>

<script lang="ts">
  import Toast from 'primevue/toast';

  export default {
    components: {
      Toast,
    },
  };
</script>
