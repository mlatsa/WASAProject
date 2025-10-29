<template>
  <div class="card">
    <h2 class="h1">Your Conversations</h2>
    <div class="list">
      <RouterLink v-for="c in rows" :key="c.id" :to="`/chat/${c.id}`" class="item">
        {{ c.name || ('Conversation ' + c.id) }}
      </RouterLink>
    </div>
    <div style="margin-top:12px; display:flex; gap:8px">
      <select v-model.number="picked" class="input" style="max-width:240px">
        <option disabled :value="0">Select user to start chat</option>
        <option v-for="u in users" :key="u.id" :value="u.id">{{ u.username }}</option>
      </select>
      <button class="btn" @click="start">New Conversation</button>
    </div>
  </div>
</template>
<script setup>
import { ref, onMounted } from 'vue'
import api from '../services/api'
import { getAuth } from '../services/auth'

const rows = ref([])
const users = ref([])
const picked = ref(0)

async function load(){
  const a = getAuth()
  users.value = await api.listUsers()
  rows.value = await api.listConversations(a.userId)
}

async function start(){
  const a = getAuth()
  if(!picked.value) return
  await api.createConversation(a.userId, { name: 'Chat', memberIds: [picked.value] })
  await load()
}

onMounted(load)
</script>
