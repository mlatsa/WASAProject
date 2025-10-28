<template>
  <div class="min-h-screen bg-slate-50">
    <header class="border-b border-slate-200 bg-white shadow-sm">
      <div class="mx-auto flex max-w-6xl items-center justify-between px-6 py-4">
        <RouterLink to="/" class="text-xl font-semibold text-slate-900">WASA</RouterLink>
        <nav class="flex items-center gap-4">
          <RouterLink
            v-if="isAuthenticated"
            to="/dashboard"
            class="text-sm font-medium text-slate-600 hover:text-slate-900"
          >
            Dashboard
          </RouterLink>
          <RouterLink
            v-if="isAuthenticated"
            to="/conversations"
            class="text-sm font-medium text-slate-600 hover:text-slate-900"
          >
            Conversations
          </RouterLink>
          <button
            v-if="isAuthenticated"
            type="button"
            class="rounded-md border border-transparent bg-slate-900 px-4 py-2 text-sm font-semibold text-white hover:bg-slate-800"
            @click="logout"
          >
            Logout
          </button>
          <RouterLink
            v-else
            to="/login"
            class="text-sm font-medium text-slate-600 hover:text-slate-900"
          >
            Login
          </RouterLink>
        </nav>
      </div>
    </header>
    <main class="mx-auto min-h-[calc(100vh-80px)] max-w-6xl px-6 py-8">
      <RouterView />
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { RouterLink, RouterView, useRouter } from 'vue-router';

import { useAuthStore } from './stores/auth';

const auth = useAuthStore();
const router = useRouter();

const isAuthenticated = computed(() => auth.isAuthenticated);

const logout = () => {
  auth.logout();
  router.push({ name: 'login' });
};
</script>
