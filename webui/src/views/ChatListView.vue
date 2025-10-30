<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import jwt_decode from 'jwt-decode'
import axios from '../services/axios'

const router = useRouter()
const conversations = ref([])
const participantInput = ref('')
const chatName = ref('')
const errorMsg = ref(null)

function myId(){
  const uid = Number(localStorage.getItem('userId'))
  if (uid) return uid
  const t = localStorage.getItem('authToken') || localStorage.getItem('identifier')
  if (!t) return null
  try { return jwt_decode(t).user_id } catch { return null }
}

async function loadConversations(){
  errorMsg.value = null
  const me = myId()
  if (!me){ router.push({name:'Login'}); return }
  try{
    const resp = await axios.get(`/users/${me}/conversations`)
    conversations.value = resp.data || []
  }catch(e){
    errorMsg.value = 'Failed to load conversations'
  }
}

async function resolveUserId(input){
  const n = Number(input)
  if (!Number.isNaN(n) && n>0) return n
  try{
    await axios.post(`/users/${myId()}/contacts`, { username: input })
  }catch(e){}
  const users = await axios.get('/users')
  const u = (users.data||[]).find(x => (x.username||'').toLowerCase() === String(input).toLowerCase())
  return u ? u.id : null
}

async function createConversation(){
  errorMsg.value = null
  const me = myId()
  if (!me){ router.push({name:'Login'}); return }
  const other = await resolveUserId(participantInput.value.trim())
  if (!other || other===me){ errorMsg.value='User not found'; return }
  try{
    const body = { members:[other], name: chatName.value||'' }
    const resp = await axios.post(`/users/${me}/conversations`, body)
    const conv = resp.data
    router.push({ name:'Chat', params:{ convId: conv.id } })
  }catch(e){
    errorMsg.value = 'Could not create conversation'
  }
}

function openConversation(conv){
  router.push({ name:'Chat', params:{ convId: conv.id } })
}

onMounted(loadConversations)
</script>

<template>
  <div class="p-6 max-w-3xl mx-auto space-y-4">
    <h1 class="text-3xl font-bold">Chats</h1>
    <div class="flex gap-2">
      <input v-model="participantInput" class="border rounded p-2 flex-1" placeholder="Participant username or ID"/>
      <input v-model="chatName" class="border rounded p-2 flex-1" placeholder="Chat name (optional)"/>
      <button @click="createConversation" class="bg-blue-600 text-white px-4 py-2 rounded">Start chat</button>
    </div>
    <p v-if="errorMsg" class="text-red-600">{{ errorMsg }}</p>
    <div class="mt-4 space-y-2">
      <div v-if="(conversations||[]).length===0" class="border rounded p-3 text-gray-600">No conversations yet.</div>
      <button v-for="c in conversations" :key="c.id" @click="openConversation(c)"
              class="w-full text-left border rounded p-3 hover:bg-gray-50">
        <div class="font-semibold">{{ c.name || ('Chat #'+c.id) }}</div>
      </button>
    </div>
  </div>
</template>
