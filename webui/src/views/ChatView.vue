<template>
  <div class="card">
    <h2 class="h2">Chat {{ id }}</h2>
    <div class="list" style="margin-bottom:12px">
      <div v-for="m in msgs" :key="m.id" class="item">
        <div><b>{{ nameOf(m.senderId) }}</b></div>
        <div>{{ m.text }}</div>
      </div>
    </div>
    <div style="display:flex; gap:8px">
      <input v-model="text" class="input" placeholder="Type a message"/>
      <button class="btn" @click="send">Send</button>
    </div>
  </div>
</template>
<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import api from '../services/api'
import { getAuth } from '../services/auth'

const route = useRoute()
const id = computed(()=> Number(route.params.id))
const msgs = ref([])
const text = ref('')
const users = ref([])

function nameOf(uid){
  const u = users.value.find(x=>x.id===uid)
  return u ? u.username : ('User '+uid)
}

async function load(){
  const a = getAuth()
  users.value = await api.listUsers()
  msgs.value = await api.listMessages(a.userId, id.value)
}

async function send(){
  const a = getAuth()
  if(!text.value.trim()) return
  await api.sendMessage(a.userId, id.value, { senderId: a.userId, text: text.value.trim() })
  text.value = ''
  await load()
}

onMounted(load)
</script>
