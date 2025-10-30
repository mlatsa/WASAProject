<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from '../services/axios'
const router = useRouter()
const chats = ref([])
const errorMsg = ref(null)
const loading = ref(false)
async function loadConversations() {
  const userId = localStorage.getItem("userId")
  loading.value = true
  errorMsg.value = null
  try {
    const resp = await axios.get(`/users/${userId}/conversations`)
    chats.value = resp.data || []
  } catch (e) {
    errorMsg.value = (e?.response?.data?.error || e?.message || 'Failed to load chats')
  } finally {
    loading.value = false
  }
}
function openConversation(conv) {
  if (!conv) return
  router.push({ name: 'Chat', params: { convId: conv.id } })
}
onMounted(loadConversations)
</script>

<template>
  <div class="container mt-4">
    <h2>Chats</h2>
    <div v-if="errorMsg" class="alert alert-danger">{{ errorMsg }}</div>
    <div v-else-if="loading">Loading…</div>
    <ul class="list-group" v-else>
      <li v-for="c in chats" :key="c.id" class="list-group-item d-flex justify-content-between align-items-center" @click="openConversation(c)" style="cursor:pointer">
        <span>{{ c.name || ('Conversation ' + c.id) }}</span>
        <span class="badge bg-secondary">{{ (c.participants||[]).length }}</span>
      </li>
    </ul>
  </div>
</template>
