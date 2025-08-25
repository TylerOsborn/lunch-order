import { createApp } from 'vue';
import { createRouter, createWebHistory } from 'vue-router';
import './style.css';
import App from './App.vue';
import PrimeVue from 'primevue/config';
import Aura from '@primevue/themes/aura';
import ToastService from 'primevue/toastservice';
import { VueQueryPlugin } from '@tanstack/vue-query'
import HomeScreen from './components/HomeScreen.vue';
import GiveMealScreen from './components/GiveMealScreen.vue';
import ReceiveMealScreen from './components/ReceiveMealScreen.vue';
import AdminScreen from './components/AdminScreen.vue';
import DonationRequestScreen from './components/DonationRequestScreen.vue';

const routes = [
  { path: '/', component: HomeScreen },
  { path: '/give-meal', component: GiveMealScreen },
  { path: '/receive-meal', component: ReceiveMealScreen },
  { path: '/donation-request', component: DonationRequestScreen },
  { 
    path: '/admin', 
    component: AdminScreen,
    meta: { requiresAuth: true }
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// Add route guard for admin authentication
router.beforeEach((to, _from, next) => {
  if (to.meta.requiresAuth) {
    // The admin screen will handle showing the login form
    // So we always allow navigation to admin route
    next();
  } else {
    next();
  }
});

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
