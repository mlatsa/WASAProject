<template>
  <div class="grid gap-6 lg:grid-cols-[280px_1fr]">
    <aside class="space-y-4">
      <div class="flex items-center justify-between">
        <h2 class="text-lg font-semibold text-slate-900">Conversations</h2>
        <button
          type="button"
          class="text-sm font-medium text-slate-500 hover:text-slate-900"
          @click="loadConversations"
        >
          Refresh
        </button>
      </div>
      <ul class="space-y-2">
        <li v-for="conversation in conversations" :key="conversation.id">
          <button
            type="button"
            class="w-full rounded-md border px-3 py-2 text-left text-sm transition"
            :class="conversation.id === selectedId ? 'border-slate-900 bg-slate-100 text-slate-900' : 'border-slate-200 bg-white text-slate-600 hover:border-slate-300'"
            @click="selectConversation(conversation.id)"
          >
            <p class="font-semibold text-slate-900">{{ conversation.name ?? conversation.id }}</p>
            <p class="mt-1 truncate text-xs text-slate-500">{{ conversation.lastMessage ?? 'No messages yet' }}</p>
          </button>
        </li>
      </ul>
    </aside>

    <section v-if="selectedConversation" class="space-y-6">
      <header>
        <h3 class="text-xl font-semibold text-slate-900">
          {{ selectedConversation.name ?? selectedConversation.id }}
        </h3>
        <p class="text-sm text-slate-500">
          Participants: {{ selectedConversation.participants?.join(', ') ?? 'Unknown' }}
        </p>
      </header>

      <div class="space-y-3 rounded-lg border border-slate-200 bg-white p-4 shadow-sm">
        <h4 class="text-sm font-semibold uppercase tracking-wide text-slate-500">Messages</h4>
        <div v-if="!selectedConversation.messages?.length" class="py-6 text-center text-sm text-slate-500">
          No messages yet. Send one below!
        </div>
        <ul v-else class="space-y-3">
          <li
            v-for="message in selectedConversation.messages"
            :key="message.messageId"
            class="flex items-start justify-between rounded-md border border-slate-200 px-3 py-2"
          >
            <div>
              <p class="text-sm text-slate-900">{{ message.content }}</p>
              <p class="mt-1 text-xs text-slate-400">{{ formatTimestamp(message.timestamp) }}</p>
            </div>
            <button
              type="button"
              class="text-xs font-medium text-red-600 hover:text-red-700"
              @click="removeMessage(message.messageId)"
            >
              Delete
            </button>
          </li>
        </ul>
      </div>

      <form class="flex flex-col gap-3 rounded-lg border border-slate-200 bg-white p-4 shadow-sm" @submit.prevent="send">
        <label class="text-sm font-semibold text-slate-700" for="message">Send a message</label>
        <textarea
          id="message"
          v-model="newMessage"
          rows="3"
          class="w-full rounded-md border border-slate-300 px-3 py-2 text-sm focus:border-slate-900 focus:outline-none focus:ring-2 focus:ring-slate-200"
          placeholder="Type your message"
        ></textarea>
        <div class="flex justify-end gap-3">
          <button
            type="submit"
            class="rounded-md bg-slate-900 px-4 py-2 text-sm font-semibold text-white hover:bg-slate-800 disabled:cursor-not-allowed disabled:bg-slate-400"
            :disabled="!newMessage || isSubmitting"
          >
            <span v-if="isSubmitting">Sending…</span>
            <span v-else>Send message</span>
          </button>
        </div>
        <p v-if="error" class="text-sm text-red-600">{{ error }}</p>
      </form>
    </section>

    <section v-else class="flex items-center justify-center rounded-lg border border-dashed border-slate-300 p-8 text-center text-slate-500">
      Select a conversation to view its messages.
    </section>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue';

import type { Conversation, ConversationWithMessages } from '@/api';
import { deleteMessage, fetchConversation, fetchConversations, sendMessage } from '@/api';

const conversations = ref<Conversation[]>([]);
const selectedId = ref<string | null>(null);
const selectedConversation = ref<ConversationWithMessages | null>(null);
const newMessage = ref('');
const isSubmitting = ref(false);
const error = ref<string | null>(null);

const loadConversations = async () => {
  try {
    const data = await fetchConversations();
    conversations.value = data.conversations;
    if (!selectedId.value && data.conversations.length) {
      selectedId.value = data.conversations[0].id;
    }
  } catch (err) {
    console.error(err);
  }
};

const loadConversation = async (conversationId: string) => {
  try {
    const data = await fetchConversation(conversationId);
    selectedConversation.value = data ?? null;
  } catch (err) {
    console.error(err);
    error.value = 'Unable to load conversation.';
  }
};

const selectConversation = (conversationId: string) => {
  selectedId.value = conversationId;
};

watch(selectedId, (id) => {
  error.value = null;
  if (id) {
    loadConversation(id);
  } else {
    selectedConversation.value = null;
  }
});

const send = async () => {
  if (!selectedId.value || !newMessage.value) return;
  isSubmitting.value = true;
  error.value = null;
  try {
    await sendMessage(selectedId.value, { content: newMessage.value, type: 'text' });
    newMessage.value = '';
    await loadConversation(selectedId.value);
  } catch (err) {
    console.error(err);
    error.value = 'Unable to send message.';
  } finally {
    isSubmitting.value = false;
  }
};

const removeMessage = async (messageId: string) => {
  if (!selectedId.value) return;
  try {
    await deleteMessage(messageId);
    await loadConversation(selectedId.value);
  } catch (err) {
    console.error(err);
    error.value = 'Unable to delete message.';
  }
};

const formatTimestamp = (value?: string) => {
  if (!value) return 'Just now';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }
  return date.toLocaleString();
};

onMounted(() => {
  loadConversations();
});
</script>
