import { createRouter, createWebHistory } from 'vue-router';

import { useAuthStore } from '@/stores/auth';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: { name: 'dashboard' }
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/pages/LoginView.vue')
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: () => import('@/pages/DashboardView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/conversations',
      name: 'conversations',
      component: () => import('@/pages/ConversationsView.vue'),
      meta: { requiresAuth: true }
    }
  ]
});

router.beforeEach((to) => {
  const auth = useAuthStore();
  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    return { name: 'login', query: { redirect: to.fullPath } };
  }

  if (to.name === 'login' && auth.isAuthenticated) {
    return { name: 'dashboard' };
  }

  return true;
});

export default router;
