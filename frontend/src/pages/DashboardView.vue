<template>
  <div class="space-y-8">
    <section class="rounded-lg border border-slate-200 bg-white p-6 shadow-sm">
      <h2 class="text-xl font-semibold text-slate-900">Welcome</h2>
      <p class="mt-2 text-sm text-slate-600">
        You're connected to the WASA backend. Use the tools below to explore the API.
      </p>
      <div class="mt-4 flex items-center gap-3">
        <span class="inline-flex items-center rounded-full bg-emerald-100 px-3 py-1 text-sm font-medium text-emerald-700">
          Status: {{ healthStatus ?? 'Loading…' }}
        </span>
      </div>
    </section>

    <section class="rounded-lg border border-slate-200 bg-white p-6 shadow-sm">
      <h3 class="text-lg font-semibold text-slate-900">Update username</h3>
      <p class="mt-1 text-sm text-slate-600">
        Set a username for your account. This uses the <code class="rounded bg-slate-100 px-1 py-0.5">POST /user/username</code>
        endpoint.
      </p>
      <form class="mt-4 flex flex-col gap-3 sm:flex-row" @submit.prevent="onSubmit">
        <input
          v-model="username"
          type="text"
          required
          minlength="3"
          maxlength="16"
          placeholder="Choose a username"
          class="flex-1 rounded-md border border-slate-300 px-3 py-2 text-sm focus:border-slate-900 focus:outline-none focus:ring-2 focus:ring-slate-200"
        />
        <button
          type="submit"
          class="rounded-md bg-slate-900 px-4 py-2 text-sm font-semibold text-white hover:bg-slate-800 disabled:cursor-not-allowed disabled:bg-slate-400"
          :disabled="isSubmitting"
        >
          <span v-if="isSubmitting">Saving…</span>
          <span v-else>Save username</span>
        </button>
      </form>
      <p v-if="message" class="mt-2 text-sm text-emerald-600">{{ message }}</p>
      <p v-if="error" class="mt-2 text-sm text-red-600">{{ error }}</p>
    </section>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';

import { fetchHealth, updateUsername } from '@/api';

const healthStatus = ref<string | null>(null);
const username = ref('');
const error = ref<string | null>(null);
const message = ref<string | null>(null);
const isSubmitting = ref(false);

onMounted(async () => {
  try {
    const data = await fetchHealth();
    healthStatus.value = data.status;
  } catch (err) {
    console.error(err);
    healthStatus.value = 'Unavailable';
  }
});

const onSubmit = async () => {
  if (!username.value) return;
  error.value = null;
  message.value = null;
  isSubmitting.value = true;
  try {
    const response = await updateUsername({ name: username.value });
    message.value = response.message ?? 'Username updated successfully.';
  } catch (err) {
    console.error(err);
    error.value = 'Unable to update username.';
  } finally {
    isSubmitting.value = false;
  }
};
</script>
