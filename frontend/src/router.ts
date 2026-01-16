import { createRouter, createWebHistory } from 'vue-router';
import HomeScreen from './components/HomeScreen.vue';
import GiveMealScreen from './components/GiveMealScreen.vue';
import ReceiveMealScreen from './components/ReceiveMealScreen.vue';
import AdminScreen from './components/AdminScreen.vue';
import DonationRequestScreen from './components/DonationRequestScreen.vue';
import LoginScreen from './components/LoginScreen.vue';
import NotFound from './components/errors/404.vue';
import Unauthorized from './components/errors/401.vue';
import Forbidden from './components/errors/403.vue';
import { userStore } from './store/user';

const routes = [
  { path: '/login', component: LoginScreen },
  { path: '/', component: HomeScreen, meta: { requiresAuth: true } },
  { path: '/give-meal', component: GiveMealScreen, meta: { requiresAuth: true } },
  { path: '/receive-meal', component: ReceiveMealScreen, meta: { requiresAuth: true } },
  { path: '/donation-request', component: DonationRequestScreen, meta: { requiresAuth: true } },
  { path: '/admin', component: AdminScreen, meta: { requiresAuth: true, requiresAdmin: true } },
  { path: '/401', component: Unauthorized },
  { path: '/403', component: Forbidden },
  { path: '/:pathMatch(.*)*', component: NotFound },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach(async (to, _from, next) => {
  if (!userStore.isAuthenticated) {
      await userStore.fetchUser();
  }

  if (to.meta.requiresAuth && !userStore.isAuthenticated) {
    next('/401');
  } else if (to.path === '/login' && userStore.isAuthenticated) {
    next('/');
  } else if (to.meta.requiresAdmin && !userStore.user?.isAdmin) {
    next('/403');
  } else {
    next();
  }
});

export default router;
