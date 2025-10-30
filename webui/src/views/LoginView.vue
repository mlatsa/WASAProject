<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import axios from '../services/axios'
const router = useRouter()
const username = ref('')
const errorMsg = ref(null)
const loading = ref(false)
async function login(){
  errorMsg.value = null
  loading.value = true
  try{
    const resp = await axios.post('/session',{ username: username.value })
    const token = resp.data.identifier
    const userId = resp.data.userId
    localStorage.setItem('authToken', token)
    localStorage.setItem('identifier', token)
    localStorage.setItem('userId', String(userId))
    router.push({ name: 'ChatList' })
  }catch(e){
    errorMsg.value = 'Login failed'
  }finally{
    loading.value = false
  }
}
</script>

<template>
  <div class="p-6 max-w-xl mx-auto">
    <h1 class="text-3xl font-bold mb-6">Login</h1>
    <div class="space-y-4">
      <input v-model="username" class="w-full border rounded p-2" placeholder="Username"/>
      <button @click="login" class="bg-blue-600 text-white px-4 py-2 rounded" :disabled="loading">Login</button>
      <p v-if="errorMsg" class="text-red-600">{{ errorMsg }}</p>
    </div>
  </div>
</template>
