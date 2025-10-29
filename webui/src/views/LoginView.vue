<template>
  <div class="p-8 max-w-xl mx-auto">
    <h1 class="text-3xl font-semibold mb-6">Login</h1>
    <label class="block text-sm mb-2">Username:</label>
    <input v-model="name" class="border px-3 py-2 w-full rounded mb-4" />
    <button @click="login" class="bg-blue-600 text-white px-4 py-2 rounded">Login</button>
  </div>
</template>
<script setup>
import { ref } from 'vue';
import axios from '../services/axios';
import { setAuth } from '../services/auth';
import { useRouter } from 'vue-router';
const router = useRouter();
const name = ref('');
async function login() {
  const r = await axios.post('/session', { name: name.value.trim() || 'mariam' });
  const { identifier, userId } = r.data;
  setAuth({ token: identifier, userId, username: name.value.trim() || 'mariam' });
  router.push('/chats');
}
</script>
