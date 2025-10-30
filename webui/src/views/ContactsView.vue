<script setup>
import { ref, onMounted } from 'vue'
import axios from '../services/axios'
import jwt_decode from 'jwt-decode'

const contacts = ref([])
const query = ref('')
const errorMsg = ref(null)

function myId(){
  const uid = Number(localStorage.getItem('userId'))
  if (uid) return uid
  const t = localStorage.getItem('authToken') || localStorage.getItem('identifier')
  try { return jwt_decode(t).user_id } catch { return null }
}

async function load(){
  errorMsg.value = null
  try{
    const resp = await axios.get(`/users/${myId()}/contacts`)
    contacts.value = resp.data || []
  }catch(e){
    contacts.value = []
  }
}

async function add(){
  errorMsg.value = null
  try{
    await axios.post(`/users/${myId()}/contacts`, { username: query.value.trim() })
    query.value = ''
    await load()
  }catch(e){
    errorMsg.value = 'Could not add'
  }
}

async function remove(name){
  try{
    await axios.delete(`/users/${myId()}/contacts`, { data:{ username:name } })
    await load()
  }catch(e){}
}

onMounted(load)
</script>

<template>
  <div class="p-6 max-w-3xl mx-auto">
    <h1 class="text-3xl font-bold mb-4">Your Contacts</h1>
    <div class="flex gap-2 mb-4">
      <input v-model="query" class="border p-2 rounded flex-1" placeholder="username"/>
      <button @click="add" class="bg-blue-600 text-white px-4 py-2 rounded">Add</button>
    </div>
    <p v-if="errorMsg" class="text-red-600">{{ errorMsg }}</p>
    <ul class="space-y-2">
      <li v-for="c in contacts" :key="c" class="flex justify-between border rounded p-2">
        <span>{{ c }}</span>
        <button @click="remove(c)" class="bg-red-600 text-white px-3 py-1 rounded">Remove</button>
      </li>
    </ul>
  </div>
</template>
