<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import axios from '../services/axios'
const route = useRoute()
const conversation = ref(null)
const messages = ref([])
const errorMsg = ref(null)
const loading = ref(false)
async function loadConversation() {
  const userId = localStorage.getItem("userId")
  const id = route.params.convId
  if (!id) return
  loading.value = true
  errorMsg.value = null
  try {
    const convResp = await axios.get(`/users/${userId}/conversations/${id}`)
    conversation.value = convResp.data
    const msgResp = await axios.get(`/users/${userId}/conversations/${id}/messages`)
    messages.value = msgResp.data || []
  } catch (e) {
    errorMsg.value = (e?.response?.data?.error || e?.message || 'Failed to load conversation')
  } finally {
    loading.value = false
  }
}
onMounted(loadConversation)
watch(() => route.params.convId, loadConversation)
</script>

<template>
  <main v-if="conversation" class="container mt-4">
    <h2>{{ conversation.name || ('Conversation ' + conversation.id) }}</h2>
    <div v-if="errorMsg" class="alert alert-danger">{{ errorMsg }}</div>
    <div v-else-if="loading">Loading…</div>
    <div v-else class="list-group">
      <div v-for="m in (messages || [])" :key="m.id" class="list-group-item">
        <strong>{{ m.author?.name || m.author || 'User' }}:</strong>
        <span class="ms-1">{{ m.text || m.body || m.content }}</span>
      </div>
    </div>
  </main>
</template>
