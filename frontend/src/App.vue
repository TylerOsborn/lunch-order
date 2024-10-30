<template>
  <Toast />
  <router-view></router-view>
</template>

<script>
  import Toast from 'primevue/toast';

  export default {
    components: {
      Toast,
    },
  };
</script>
