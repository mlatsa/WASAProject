<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import jwt_decode from 'jwt-decode'
import axios from '../services/axios'

const route = useRoute()
const conv = ref(null)
const messageText = ref('')
const errorMsg = ref(null)

function myId(){
  const uid = Number(localStorage.getItem('userId'))
  if (uid) return uid
  const t = localStorage.getItem('authToken') || localStorage.getItem('identifier')
  try { return jwt_decode(t).user_id } catch { return null }
}

async function loadConversation(){
  errorMsg.value = null
  const me = myId()
  if (!me) return
  try{
    const resp = await axios.get(`/users/${me}/conversations/${route.params.convId}`)
    conv.value = resp.data
  }catch(e){
    errorMsg.value = 'Failed to load conversation'
  }
}

async function sendMessage(){
  const me = myId()
  if (!me || !messageText.value.trim()) return
  try{
    await axios.post(`/users/${me}/conversations/${route.params.convId}/messages`, { content: messageText.value.trim() })
    messageText.value = ''
    await loadConversation()
  }catch(e){
    errorMsg.value = 'Failed to send'
  }
}

onMounted(loadConversation)
watch(() => route.params.convId, loadConversation)
</script>

<template>
  <div class="p-6 max-w-4xl mx-auto">
    <h2 class="text-2xl font-bold mb-4">{{ conv?.name || ('Chat #' + route.params.convId) }}</h2>
    <p v-if="errorMsg" class="text-red-600">{{ errorMsg }}</p>
    <div class="border rounded p-4 h-96 overflow-auto whitespace-pre-line">
      <template v-for="m in (conv?.messages || [])" :key="m.id">
        <div class="mb-2"><strong>{{ m.senderId===myId() ? 'You' : 'Them' }}:</strong> {{ m.content }}</div>
      </template>
    </div>
    <div class="flex gap-2 mt-3">
      <input v-model="messageText" @keyup.enter="sendMessage" class="flex-1 border rounded p-2" placeholder="Type a message..."/>
      <button @click="sendMessage" class="px-4 py-2 bg-blue-600 text-white rounded">Send</button>
    </div>
  </div>
</template>
