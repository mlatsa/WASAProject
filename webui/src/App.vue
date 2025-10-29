<template>
  <main style="padding:16px;font-family:system-ui">
    <h1>WebUI is running ✅</h1>
    <p>Click to call the API at <code>/v1/health</code>.</p>
    <button @click="ping" style="padding:8px 12px;border:1px solid #ccc;border-radius:8px">Ping /v1/health</button>
    <pre style="background:#f6f6f6;padding:12px;border-radius:8px;margin-top:12px">{{ out }}</pre>
  </main>
</template>

<script setup>
import { ref } from 'vue'
const out = ref('')

async function ping() {
  try {
    const r = await fetch('http://localhost:3000/v1/health')
    const j = await r.json()
    out.value = JSON.stringify(j, null, 2)
  } catch (e) {
    out.value = String(e)
  }
}
</script>
