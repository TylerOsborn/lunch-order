import { createApp } from 'vue';
import './style.css';
import App from './App.vue';
import PrimeVue from 'primevue/config';
import Aura from '@primevue/themes/aura';
import ToastService from 'primevue/toastservice';
import { VueQueryPlugin } from '@tanstack/vue-query'
import router from './router';
import 'primeicons/primeicons.css';

const app = createApp(App);
app.use(router);
app.use(ToastService);
app.use(VueQueryPlugin);

app.use(PrimeVue, {
  theme: {
    preset: Aura,
  },
});

app.mount('#app');
