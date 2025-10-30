<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import axios from '../services/axios'
const router = useRouter()
const username = ref('')
const loading = ref(false)
const errorMsg = ref(null)
async function login() {
  loading.value = true
  errorMsg.value = null
  try {
    const resp = await axios.post('/session', { username: username.value })
    localStorage.setItem("authToken", (resp.data.token || resp.data.access_token || resp.data.jwt || resp.data.identifier || ""))
    router.push({ name: 'ChatList' })
  } catch (e) {
    errorMsg.value = (e?.response?.data?.error || e?.message || 'Login failed')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="container mt-5" style="max-width:640px">
    <h1>Login</h1>
    <div class="mb-3">
      <label class="form-label">Username:</label>
      <input class="form-control" v-model="username" placeholder="username" />
    </div>
    <button class="btn btn-primary" :disabled="loading" @click="login">Login</button>
    <div v-if="errorMsg" class="alert alert-danger mt-3">{{ errorMsg }}</div>
  </div>
</template>
