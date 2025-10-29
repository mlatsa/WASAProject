<template>
  <main style="padding:16px;font-family:system-ui">
    <h1>WebUI is running ✅</h1>
    <div style="margin:12px 0">
      <button @click="pingHealth">Ping /v1/health</button>
      <button @click="pingExt" style="margin-left:8px">GET /ext/todos/1</button>
    </div>
    <pre style="background:#f6f6f6;padding:12px;border-radius:8px">{{ out }}</pre>
  </main>
</template>

<script setup>
import { ref } from 'vue'
const out = ref('')
async function pingHealth() {
  const r = await fetch('http://localhost:3000/v1/health')
  out.value = await r.text()
}
async function pingExt() {
  const r = await fetch('http://localhost:3000/ext/todos/1')
  out.value = await r.text()
}
</script>
