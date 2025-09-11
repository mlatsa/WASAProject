import { createApp, ref, computed, onMounted } from "vue";

/** Build API base from env for production (e.g., '/api') or empty for dev (proxy) */
const API_BASE = import.meta.env.VITE_API_BASE || "";

/** tiny fetch helper */
async function api(path, init = {}) {
  const url = API_BASE ? `${API_BASE}${path}` : path; // relative path (no absolute host)
  const res = await fetch(url, init);
  const text = await res.text();
  let data = text;
  try { data = JSON.parse(text); } catch {}
  return { res, data };
}

const App = {
  name: "App",
  setup() {
    const name = ref("Alex");
    const token = ref("");
    const status = ref("");
    const body = ref("");
    const convs = ref([]);
    const convId = ref("conversation_abc");
    const convo = ref(null);
    const msgText = ref("");
    const msgType = ref("text");
    const lastMsgId = ref("");
    const lastReactionId = ref("");
    const forwardTo = ref("c2");
    const username = ref("alex99");
    const groupName = ref("My Group");

    const auth = (extra = {}) => ({ Authorization: `Bearer ${token.value}`, ...extra });
    const setDbg = (r, d, route = "") => {
      status.value = `${r.status} ${r.statusText}${route ? " → " + route : ""}`;
      body.value = typeof d === "string" ? d : JSON.stringify(d, null, 2);
    };

    async function health() {
      const r = await api("/health");
      setDbg(r.res, r.data, "/health");
    }

    async function login() {
      const r = await api("/session", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name: name.value }),
      });
      setDbg(r.res, r.data, "POST /session");
      if (r.res.ok && r.data && r.data.identifier) {
        token.value = r.data.identifier;
        await listConvs();
        await getConv(convId.value);
      }
    }

    async function setUsername() {
      const r = await api("/user/username", {
        method: "PUT",
        headers: auth({ "Content-Type": "application/json" }),
        body: JSON.stringify({ username: username.value }),
      });
      setDbg(r.res, r.data, "PUT /user/username");
    }

    async function setUserPhoto() {
      const r = await api("/user/photo", { method: "PUT", headers: auth() });
      setDbg(r.res, r.data, "PUT /user/photo");
    }

    async function listConvs() {
      const r = await api("/conversations", { headers: auth() });
      setDbg(r.res, r.data, "GET /conversations");
      if (r.res.ok && r.data && Array.isArray(r.data.conversations)) {
        convs.value = r.data.conversations;
      }
    }

    async function getConv(id) {
      convId.value = id;
      const r = await api(`/conversations/${encodeURIComponent(id)}`, { headers: auth() });
      setDbg(r.res, r.data, "GET /conversations/{id}");
      if (r.res.ok && r.data && r.data.conversation) {
        convo.value = r.data.conversation;
      }
    }

    async function sendMsg() {
      if (!msgText.value.trim()) return;
      const r = await api(`/conversations/${encodeURIComponent(convId.value)}/messages`, {
        method: "POST",
        headers: auth({ "Content-Type": "application/json" }),
        body: JSON.stringify({ content: msgText.value, type: msgType.value }),
      });
      setDbg(r.res, r.data, "POST /conversations/{id}/messages");
      if (r.res.ok && r.data && r.data.messageId) {
        lastMsgId.value = r.data.messageId;
        msgText.value = "";
        await getConv(convId.value);
        await listConvs();
      }
    }

    async function delMsg() {
      if (!lastMsgId.value) return;
      const r = await api(`/messages/${encodeURIComponent(lastMsgId.value)}`, {
        method: "DELETE",
        headers: auth(),
      });
      setDbg(r.res, r.data, "DELETE /messages/{messageId}");
      await getConv(convId.value);
    }

    async function react() {
      if (!lastMsgId.value) return;
      const r = await api(`/messages/${encodeURIComponent(lastMsgId.value)}/reactions`, {
        method: "POST",
        headers: auth({ "Content-Type": "application/json" }),
        body: JSON.stringify({ emoji: "👍" }),
      });
      setDbg(r.res, r.data, "POST /messages/{messageId}/reactions");
      if (r.res.ok && r.data && r.data.reactionId) {
        lastReactionId.value = r.data.reactionId;
      }
    }

    async function removeReaction() {
      if (!lastMsgId.value || !lastReactionId.value) return;
      const r = await api(`/messages/${encodeURIComponent(lastMsgId.value)}/reactions/${encodeURIComponent(lastReactionId.value)}`, {
        method: "DELETE",
        headers: auth(),
      });
      setDbg(r.res, r.data, "DELETE /messages/{messageId}/reactions/{reactionId}");
    }

    async function fwd() {
      if (!lastMsgId.value) return;
      const r = await api(`/messages/${encodeURIComponent(lastMsgId.value)}/forward`, {
        method: "POST",
        headers: auth({ "Content-Type": "application/json" }),
        body: JSON.stringify({ conversationId: forwardTo.value }),
      });
      setDbg(r.res, r.data, "POST /messages/{messageId}/forward");
    }

    async function addMember() {
      const r = await api(`/groups/${encodeURIComponent(convId.value)}/members`, {
        method: "POST",
        headers: auth({ "Content-Type": "application/json" }),
        body: JSON.stringify({ member: "Bob" }),
      });
      setDbg(r.res, r.data, "POST /groups/{id}/members");
    }

    async function leaveGroup() {
      const r = await api(`/groups/${encodeURIComponent(convId.value)}/leave`, {
        method: "POST",
        headers: auth(),
      });
      setDbg(r.res, r.data, "POST /groups/{id}/leave");
    }

    async function renameGroup() {
      const r = await api(`/groups/${encodeURIComponent(convId.value)}/name`, {
        method: "PUT",
        headers: auth({ "Content-Type": "application/json" }),
        body: JSON.stringify({ name: groupName.value }),
      });
      setDbg(r.res, r.data, "PUT /groups/{id}/name");
    }

    async function setGroupPhoto() {
      const r = await api(`/groups/${encodeURIComponent(convId.value)}/photo`, {
        method: "PUT",
        headers: auth(),
      });
      setDbg(r.res, r.data, "PUT /groups/{id}/photo");
    }

    const messages = computed(() =>
      convo.value && Array.isArray(convo.value.messages) ? convo.value.messages : []
    );

    onMounted(() => { health(); });

    return {
      name, token, status, body,
      convs, convId, convo, messages,
      msgText, msgType, lastMsgId, lastReactionId, forwardTo,
      username, groupName,
      health, login, setUsername, setUserPhoto,
      listConvs, getConv, sendMsg, delMsg, react, removeReaction, fwd,
      addMember, leaveGroup, renameGroup, setGroupPhoto
    };
  },
};

createApp(App).mount("#app");
