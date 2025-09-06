<script setup>
import { ref } from 'vue'

const api = import.meta.env.VITE_API_BASE || 'http://localhost:3000'

const name = ref('Alex')
const identifier = ref('')
const convId = ref('conversation_abc')
const content = ref('hey!')
const out = ref('')

async function call(path, opts={}) {
  out.value = '...loading'
  try {
    const r = await fetch(`${api}${path}`, opts)
    const txt = await r.text()
    try { out.value = JSON.stringify(JSON.parse(txt), null, 2) }
    catch { out.value = txt }
  } catch (e) {
    out.value = String(e)
  }
}

async function login () {
  await call('/session', {
    method:'POST',
    headers:{ 'Content-Type':'application/json' },
    body: JSON.stringify({ name: name.value })
  })
  try { identifier.value = JSON.parse(out.value).identifier || '' } catch {}
}

async function listConvs () {
  await call('/conversations', {
    headers:{ Authorization:`Bearer ${identifier.value}` }
  })
}

async function getConv () {
  await call(`/conversations/${encodeURIComponent(convId.value)}`, {
    headers:{ Authorization:`Bearer ${identifier.value}` }
  })
}

async function sendMsg () {
  await call(`/conversations/${encodeURIComponent(convId.value)}/messages`, {
    method:'POST',
    headers:{
      'Content-Type':'application/json',
      Authorization:`Bearer ${identifier.value}`
    },
    body: JSON.stringify({ content: content.value, type: 'text' })
  })
}
</script>

<template>
  <main style="max-width:760px;margin:2rem auto;font-family:system-ui, sans-serif">
    <h1>WASAText (Frontend)</h1>
    <p>API base: <code>{{ api }}</code></p>

    <section style="display:grid;gap:8px;grid-template-columns:1fr auto">
      <input v-model="name" placeholder="Name to login" />
      <button @click="login">POST /session (login)</button>

      <input v-model="identifier" placeholder="Identifier (from login)" />
      <button @click="listConvs">GET /conversations</button>

      <input v-model="convId" placeholder="Conversation ID" />
      <button @click="getConv">GET /conversations/{id}</button>

      <input v-model="content" placeholder="Message content" />
      <button @click="sendMsg">POST /conversations/{id}/messages</button>
    </section>

    <pre style="background:#111;color:#0f0;padding:1rem;overflow:auto;min-height:180px;margin-top:1rem">{{ out }}</pre>
  </main>
</template>

<style>
button { padding: 8px 12px; cursor: pointer; }
input  { padding: 8px 10px; }
</style>
