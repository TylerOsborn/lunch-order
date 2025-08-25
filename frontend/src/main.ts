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
import LoginScreen from './components/LoginScreen.vue';
import { useAuth } from './composables/useAuth';

const routes = [
  { path: '/', component: HomeScreen, meta: { requiresAuth: false } },
  { path: '/login', component: LoginScreen, meta: { requiresAuth: false } },
  { path: '/give-meal', component: GiveMealScreen, meta: { requiresAuth: true } },
  { path: '/receive-meal', component: ReceiveMealScreen, meta: { requiresAuth: true } },
  { path: '/donation-request', component: DonationRequestScreen, meta: { requiresAuth: true } },
  { path: '/admin', component: AdminScreen, meta: { requiresAuth: true, requiresAdmin: true } },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// Route guards
router.beforeEach(async (to, _from, next) => {
  const { state, checkAuth, isAdmin } = useAuth();
  
  // Check authentication status
  if (state.isLoading) {
    await checkAuth();
  }

  const requiresAuth = to.meta.requiresAuth;
  const requiresAdmin = to.meta.requiresAdmin;

  if (requiresAuth && !state.isAuthenticated) {
    next('/login');
  } else if (requiresAdmin && !isAdmin()) {
    next('/'); // Redirect to home if not admin
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
