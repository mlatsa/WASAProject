<template>
  <div>
    <header class="nav">
      <span class="brand">APP</span>
      <RouterLink to="/chats">Chats</RouterLink>
      <RouterLink to="/contacts">Contacts</RouterLink>
      <RouterLink to="/profile">Profile</RouterLink>
      <RouterLink to="/user">User</RouterLink>
      <span class="spacer"></span>
      <button v-if="authed" class="btn outline" @click="logout">Logout</button>
      <RouterLink v-else class="btn outline" to="/">Login</RouterLink>
    </header>
    <main class="container">
      <RouterView />
    </main>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import realtime from './services/realtime'
import './style.css'
import { getAuth, clearAuth } from './services/auth'

const router = useRouter()
const authed = ref(false)

function boot(){
  const a = getAuth()
  authed.value = a.isAuthed
  if(a.isAuthed){
    const url = `ws://localhost:3000/ws?token=${encodeURIComponent(a.token)}`
    realtime.connect(url)
  }else{
    realtime.disconnect()
    router.push('/')
  }
}

function logout(){
  clearAuth()
  boot()
}

onMounted(boot)
</script>

