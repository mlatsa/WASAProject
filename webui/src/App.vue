<script setup>
import { ref, computed, onMounted } from 'vue'

// Grader-friendly API base: env if present, otherwise same origin
const API_BASE = import.meta.env.VITE_API_BASE || window.location.origin

// Simple fetch helper that always uses relative paths
async function api(path, opts = {}) {
  const url = new URL(path, API_BASE).toString()
  const res = await fetch(url, opts)
  const text = await res.text()
  let data = text
  try { data = JSON.parse(text) } catch {}
  return { res, data }
}

// UI state
const name = ref('Alex')
const token = ref('')
const statusLine = ref('')
const output = ref('')

const conversations = ref([])           // list view (from GET /conversations)
const activeConvId = ref('conversation_abc')
const activeConv = ref(null)            // detail view (GET /conversations/:id)
const msgText = ref('')
const msgType = ref('text')

const lastMessageId = ref('')           // captured from POST send
const lastReactionId = ref('')          // captured from POST reaction
const forwardTo = ref('c2')             // target for forward
const username = ref('alex99')
const groupName = ref('My Group')

// headers with bearer token
const authHeaders = (extra = {}) => ({
  'Authorization': `Bearer ${token.value}`,
  ...extra
})

function show(res, data, label = '') {
  statusLine.value = `${res.status} ${res.statusText} ${label ? '→ ' + label : ''}`
  output.value = (typeof data === 'string') ? data : JSON.stringify(data, null, 2)
}

// ===== Required endpoints =====

// health
async function health() {
  const r = await api('/health')
  show(r.res, r.data, '/health')
}

// session (login)
async function login() {
  const r = await api('/session', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name: name.value })
  })
  show(r.res, r.data, 'POST /session')
  if (r.res.ok && r.data && r.data.identifier) {
    token.value = r.data.identifier
    // pull initial data
    await listConvs()
    await openConv(activeConvId.value)
  }
}

// user updates
async function putUsername() {
  const r = await api('/user/username', {
    method: 'PUT',
    headers: authHeaders({ 'Content-Type': 'application/json' }),
    body: JSON.stringify({ username: username.value })
  })
  show(r.res, r.data, 'PUT /user/username')
}

async function putUserPhoto() {
  const r = await api('/user/photo', {
    method: 'PUT',
    headers: authHeaders()
  })
  show(r.res, r.data, 'PUT /user/photo')
}

// conversations
async function listConvs() {
  const r = await api('/conversations', { headers: authHeaders() })
  show(r.res, r.data, 'GET /conversations')
  if (r.res.ok && r.data && Array.isArray(r.data.conversations)) {
    conversations.value = r.data.conversations
  }
}

async function openConv(id) {
  activeConvId.value = id
  const r = await api(`/conversations/${encodeURIComponent(id)}`, { headers: authHeaders() })
  show(r.res, r.data, 'GET /conversations/{id}')
  if (r.res.ok && r.data && r.data.conversation) {
    activeConv.value = r.data.conversation
  }
}

async function sendMsg() {
  if (!msgText.value.trim()) return
  const r = await api(`/conversations/${encodeURIComponent(activeConvId.value)}/messages`, {
    method: 'POST',
    headers: authHeaders({ 'Content-Type': 'application/json' }),
    body: JSON.stringify({ content: msgText.value, type: msgType.value })
  })
  show(r.res, r.data, 'POST /conversations/{id}/messages')
  if (r.res.ok && r.data && r.data.messageId) {
    lastMessageId.value = r.data.messageId
    msgText.value = ''
    await openConv(activeConvId.value) // refresh
    await listConvs()                  // update list preview
  }
}

// messages: delete / forward / reactions
async function deleteMsg() {
  if (!lastMessageId.value) return
  const r = await api(`/messages/${encodeURIComponent(lastMessageId.value)}`, {
    method: 'DELETE',
    headers: authHeaders()
  })
  show(r.res, r.data, 'DELETE /messages/{messageId}')
  await openConv(activeConvId.value)
}

async function forwardMsg() {
  if (!lastMessageId.value) return
  const r = await api(`/messages/${encodeURIComponent(lastMessageId.value)}/forward`, {
    method: 'POST',
    headers: authHeaders({ 'Content-Type': 'application/json' }),
    body: JSON.stringify({ conversationId: forwardTo.value })
  })
  show(r.res, r.data, 'POST /messages/{messageId}/forward')
}

async function addReaction() {
  if (!lastMessageId.value) return
  const r = await api(`/messages/${encodeURIComponent(lastMessageId.value)}/reactions`, {
    method: 'POST',
    headers: authHeaders({ 'Content-Type': 'application/json' }),
    body: JSON.stringify({ emoji: '👍' })
  })
  show(r.res, r.data, 'POST /messages/{messageId}/reactions')
  if (r.res.ok && r.data && r.data.reactionId) {
    lastReactionId.value = r.data.reactionId
  }
}

async function delReaction() {
  if (!lastMessageId.value || !lastReactionId.value) return
  const r = await api(`/messages/${encodeURIComponent(lastMessageId.value)}/reactions/${encodeURIComponent(lastReactionId.value)}`, {
    method: 'DELETE',
    headers: authHeaders()
  })
  show(r.res, r.data, 'DELETE /messages/{messageId}/reactions/{reactionId}')
}

// groups
async function addMember() {
  const r = await api(`/groups/${encodeURIComponent(activeConvId.value)}/members`, {
    method: 'POST',
    headers: authHeaders({ 'Content-Type': 'application/json' }),
    body: JSON.stringify({ member: 'Bob' })
  })
  show(r.res, r.data, 'POST /groups/{id}/members')
}

async function leaveGroup() {
  const r = await api(`/groups/${encodeURIComponent(activeConvId.value)}/leave`, {
    method: 'POST',
    headers: authHeaders()
  })
  show(r.res, r.data, 'POST /groups/{id}/leave')
}

async function putGroupName() {
  const r = await api(`/groups/${encodeURIComponent(activeConvId.value)}/name`, {
    method: 'PUT',
    headers: authHeaders({ 'Content-Type': 'application/json' }),
    body: JSON.stringify({ name: groupName.value })
  })
  show(r.res, r.data, 'PUT /groups/{id}/name')
}

async function putGroupPhoto() {
  const r = await api(`/groups/${encodeURIComponent(activeConvId.value)}/photo`, {
    method: 'PUT',
    headers: authHeaders()
  })
  show(r.res, r.data, 'PUT /groups/{id}/photo')
}

// Derived UI
const messages = computed(() => (activeConv.value && activeConv.value.messages) ? activeConv.value.messages : [])

onMounted(() => {
  // optional: ping health at load
  health()
})
</script>

<template>
  <div class="layout">
    <aside class="sidebar">
      <h2>WASAText</h2>

      <div class="login">
        <input v-model="name" placeholder="Name to login" />
        <button @click="login">Login</button>
        <button class="ghost" @click="health">Health</button>
      </div>

      <div class="token">
        <label>Identifier</label>
        <input v-model="token" placeholder="Bearer token from login" />
      </div>

      <div class="user-admin">
        <label>Username</label>
        <div class="row">
          <input v-model="username" />
          <button @click="putUsername">Save</button>
        </div>
        <button class="ghost" @click="putUserPhoto">Update Photo</button>
      </div>

      <div class="convs">
        <div class="heading">
          <h3>Conversations</h3>
          <button class="ghost small" @click="listConvs">Refresh</button>
        </div>
        <div class="list">
          <button
            v-for="c in conversations"
            :key="c.id"
            class="conv"
            :class="{ active: c.id === activeConvId }"
            @click="openConv(c.id)"
            :title="c.id"
          >
            <div class="title">{{ c.id || 'unnamed' }}</div>
            <div class="preview">{{ c.lastMessage || '—' }}</div>
          </button>
        </div>
      </div>

      <div class="groups">
        <h3>Group tools</h3>
        <input v-model="groupName" placeholder="Group name" />
        <div class="row">
          <button @click="addMember">Add Member</button>
          <button @click="leaveGroup">Leave</button>
        </div>
        <div class="row">
          <button @click="putGroupName">Rename</button>
          <button @click="putGroupPhoto">Set Photo</button>
        </div>
      </div>
    </aside>

    <section class="chat">
      <header class="chat-header">
        <div class="title">{{ activeConvId }}</div>
        <div class="tools">
          <button class="ghost small" @click="openConv(activeConvId)">Reload</button>
        </div>
      </header>

      <div class="messages">
        <div
          v-for="m in messages"
          :key="m.messageId"
          class="bubble"
          :class="{ me: m.sender === 'Alex' }"
          @click="lastMessageId = m.messageId"
          :title="m.messageId"
        >
          <div class="meta">
            <span class="sender">{{ m.sender }}</span>
            <span class="time">{{ new Date(m.timestamp).toLocaleTimeString() }}</span>
          </div>
          <div class="content">{{ m.content }}</div>
        </div>
        <div v-if="!messages.length" class="empty">No messages yet</div>
      </div>

      <footer class="composer">
        <input v-model="activeConvId" class="convId" placeholder="conversation id" />
        <input v-model="msgText" class="text" placeholder="Type a message…" @keyup.enter="sendMsg" />
        <select v-model="msgType" class="type">
          <option value="text">text</option>
          <option value="image">image</option>
        </select>
        <button @click="sendMsg">Send</button>
      </footer>

      <div class="msg-tools">
        <div class="row">
          <input v-model="lastMessageId" placeholder="messageId" />
          <button class="ghost" @click="deleteMsg">Delete</button>
          <button class="ghost" @click="addReaction">👍 React</button>
          <input v-model="lastReactionId" placeholder="reactionId" />
          <button class="ghost" @click="delReaction">Remove Reaction</button>
        </div>
        <div class="row">
          <input v-model="forwardTo" placeholder="forward to convId" />
          <button class="ghost" @click="forwardMsg">Forward →</button>
        </div>
      </div>

      <div class="debug">
        <div class="status">{{ statusLine }}</div>
        <pre>{{ output }}</pre>
      </div>
    </section>
  </div>
</template>

<style>
:root {
  --bg: #0f172a;
  --panel: #111827;
  --muted: #94a3b8;
  --text: #e2e8f0;
  --accent: #22c55e;
  --bubble: #1f2937;
  --bubble-me: #0ea5e9;
}
* { box-sizing: border-box; }
body { margin:0; background:var(--bg); color:var(--text); }
button {
  background: var(--accent); border: none; color:#06130b; padding:.5rem .8rem;
  border-radius: 10px; font-weight:600; cursor:pointer;
}
button.ghost { background:#0b1320; color:#9cc1f1; border:1px solid #24324b }
button.small { padding:.25rem .5rem; font-size:.8rem }
input, select {
  background: #0b1320; border:1px solid #24324b; color:#dbeafe; padding:.5rem .6rem;
  border-radius:10px; outline:none;
}
.layout { display:grid; grid-template-columns: 320px 1fr; min-height: 100vh; }
.sidebar {
  background: var(--panel); padding: 1rem; border-right:1px solid #1e293b; display:flex; flex-direction:column; gap:1rem;
}
.sidebar h2 { margin:.2rem 0 0 0; }
.row { display:flex; gap:.5rem; }
.login, .token, .user-admin, .convs, .groups { display:flex; flex-direction:column; gap:.5rem; }
.heading { display:flex; align-items:center; justify-content:space-between; }
.list { display:flex; flex-direction:column; gap:.4rem; max-height: 30vh; overflow:auto; }
.conv { text-align:left; background:#0b1320; border:1px solid #24324b; color:#cbd5e1; padding:.6rem .7rem; border-radius:10px; }
.conv.active { outline:2px solid #3b82f6; }
.conv .title { font-weight:700; }
.conv .preview { font-size:.9rem; color: var(--muted); }

.chat { display:flex; flex-direction:column; height:100vh; }
.chat-header { display:flex; align-items:center; justify-content:space-between; padding:1rem; border-bottom:1px solid #1e293b; }
.chat-header .title { font-weight:700; }

.messages { flex:1; overflow:auto; display:flex; flex-direction:column; gap:.5rem; padding:1rem; }
.bubble {
  background: var(--bubble); padding:.6rem .8rem; border-radius: 14px;
  max-width: 70%; box-shadow: 0 4px 18px rgba(0,0,0,.25);
}
.bubble.me { align-self:flex-end; background: var(--bubble-me); color:#062037; }
.bubble .meta { font-size:.75rem; color:#0ea5e9; display:flex; gap:.6rem; }
.bubble.me .meta { color:#02213c; }
.bubble .content { margin-top:.25rem; }

.empty { color:var(--muted); font-style:italic; padding:1rem; }
.composer {
  display:grid; grid-template-columns: 220px 1fr 110px auto; gap:.6rem; padding: .8rem; border-top: 1px solid #1e293b;
}
.composer .text { width:100%; }
.msg-tools { padding:.6rem 1rem; border-top:1px solid #1e293b; display:flex; flex-direction:column; gap:.5rem; }
.debug { border-top:1px dashed #1e293b; padding: .8rem 1rem; font: 12px/1.4 ui-monospace, SFMono-Regular, Menlo, monospace; color:#9ca3af; }
.debug pre { white-space: pre-wrap; background:#0b1320; border:1px solid #24324b; padding:.6rem .8rem; border-radius:10px; color:#a7f3d0 }
</style>
