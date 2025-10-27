<template>
  <div class="mx-auto flex max-w-md flex-col gap-6 rounded-lg border border-slate-200 bg-white p-8 shadow-sm">
    <div class="space-y-2 text-center">
      <h1 class="text-2xl font-semibold text-slate-900">Welcome back</h1>
      <p class="text-sm text-slate-500">Enter a display name to log in or create an account.</p>
    </div>
    <form class="space-y-4" @submit.prevent="onSubmit">
      <div class="space-y-2 text-left">
        <label for="name" class="block text-sm font-medium text-slate-700">Display name</label>
        <input
          id="name"
          v-model="name"
          type="text"
          name="name"
          required
          minlength="3"
          maxlength="16"
          class="w-full rounded-md border border-slate-300 px-3 py-2 text-sm focus:border-slate-900 focus:outline-none focus:ring-2 focus:ring-slate-200"
        />
      </div>
      <button
        type="submit"
        class="flex w-full items-center justify-center rounded-md bg-slate-900 px-4 py-2 text-sm font-semibold text-white transition hover:bg-slate-800 disabled:cursor-not-allowed disabled:bg-slate-400"
        :disabled="isSubmitting"
      >
        <span v-if="isSubmitting" class="h-5 w-5 animate-spin rounded-full border-2 border-white border-b-transparent"></span>
        <span v-else>Continue</span>
      </button>
      <p v-if="error" class="text-sm text-red-600">{{ error }}</p>
    </form>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import { useAuthStore } from '@/stores/auth';

const auth = useAuthStore();
const router = useRouter();
const route = useRoute();

const name = ref('');
const isSubmitting = computed(() => auth.loading);
const error = computed(() => auth.error);

const onSubmit = async () => {
  if (!name.value) return;
  const redirect = (route.query.redirect as string | undefined) ?? '/dashboard';
  const success = await auth.login(name.value);
  if (success) {
    router.push(redirect);
  }
};
</script>
