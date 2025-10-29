import axios from './axios'
function listUsers(){ return axios.get('/users').then(r=>r.data) }
function listConversations(userId){ return axios.get(`/users/${userId}/conversations`).then(r=>r.data) }
function createConversation(userId, payload){ return axios.post(`/users/${userId}/conversations`, payload).then(r=>r.data) }
function listMessages(userId, convId){ return axios.get(`/users/${userId}/conversations/${convId}/messages`).then(r=>r.data) }
function sendMessage(userId, convId, payload){ return axios.post(`/users/${userId}/conversations/${convId}/messages`, payload).then(r=>r.data) }
export default { listUsers, listConversations, createConversation, listMessages, sendMessage }
