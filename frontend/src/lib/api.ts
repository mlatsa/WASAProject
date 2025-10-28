// frontend/src/lib/api.ts
import axios from "axios";

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || "",
  timeout: 15000,
});

// Optional request logger
api.interceptors.request.use((req) => {
  try {
    const url = (req.baseURL || "") + (req.url || "");
    // eslint-disable-next-line no-console
    console.log("[API req]", (req.method || "GET").toUpperCase(), url, req.data ?? null);
  } catch {}
  return req;
});

export type Message = {
  id: string;
  messageId?: string;
  sender: string;
  content: string;
  type: string;
  createdAt: string;
  reactions?: string[];
};

export type Conversation = {
  id: string;
  title: string;
  lastMessage: string;
  messages: Message[];
};

// ---- Conversations
export async function listConversations(): Promise<Conversation[]> {
  const { data } = await api.get("/api/conversations");
  return data;
}

export async function getConversation(id: string): Promise<Conversation> {
  const { data } = await api.get(`/api/conversations/${id}`);
  return data;
}

export async function sendMessage(id: string, content: string): Promise<Message> {
  const { data } = await api.post(
    `/api/conversations/${id}/messages`,
    { content },
    { headers: { Authorization: "Bearer bearer-demo-token" } }
  );
  return data;
}

export async function patchConversationTitle(id: string, title: string): Promise<Conversation> {
  const { data } = await api.patch(`/api/conversations/${id}`, { title });
  return data;
}

export async function deleteConversation(id: string): Promise<void> {
  await api.delete(`/api/conversations/${id}`);
}

// ---- Message actions
export async function reactToMessage(messageId: string, emoji: string) {
  const { data } = await api.post(`/api/messages/${messageId}/react`, { emoji });
  return data; // typically the updated reactions array
}

export async function deleteMessageById(messageId: string) {
  await api.delete(`/api/messages/${messageId}`);
}

export async function forwardMessage(messageId: string, targetId: string) {
  const { data } = await api.post(`/api/forward`, { messageId, targetId });
  return data; // forwarded message object in the target chat
}

export default api;
