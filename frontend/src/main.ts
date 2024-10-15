import {createApp} from 'vue'
import {createRouter, createWebHistory} from 'vue-router'
import './style.css'
import App from './App.vue'
import PrimeVue from 'primevue/config';
import Aura from '@primevue/themes/aura';
import ToastService from 'primevue/toastservice';
import HomeScreen from './components/HomeScreen.vue';
import GiveMealScreen from './components/GiveMealScreen.vue';
import ReceiveMealScreen from './components/ReceiveMealScreen.vue';
import AdminScreen from './components/AdminScreen.vue';

const routes = [
    {path: '/', component: HomeScreen},
    {path: '/give-meal', component: GiveMealScreen},
    {path: '/receive-meal', component: ReceiveMealScreen},
    {path: '/admin', component: AdminScreen}
]

const router = createRouter({
  history: createWebHistory(),
  routes,
});

const app = createApp(App);
app.use(router);
app.use(ToastService);

app.use(PrimeVue, {
  theme: {
    preset: Aura,
  },
});

app.mount('#app');
