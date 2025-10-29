<template>
  <main style="padding:16px">
    <h1>WebUI is running</h1>
    <button @click="ping">Ping /v1/health</button>
    <pre>{{ out }}</pre>
  </main>
</template>

<script setup>
import axios from '../services/axios'
import { ref } from 'vue'
const out = ref('')
async function ping() {
  const r = await axios.get('/v1/health')
  out.value = JSON.stringify(r.data, null, 2)
}
</script>
