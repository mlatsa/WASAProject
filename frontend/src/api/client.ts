import apiClient from '@/services/http';

import type { components, paths, ConversationEnvelope, ConversationsEnvelope } from './generated';

type LoginBody = components['schemas']['LoginBody'];
type IdentifierResponse = paths['/session']['post']['responses']['201']['content']['application/json'];
type HealthResponse = paths['/health']['get']['responses']['200']['content']['application/json'];
type SetNameBody = components['schemas']['SetNameBody'];
type UpdateUsernameResponse = paths['/user/username']['post']['responses']['200']['content']['application/json'];
type ConversationListPayload = ConversationsEnvelope;
type ConversationResponse = ConversationEnvelope;
type Message = components['schemas']['Message'];
type SendMessageInput = components['schemas']['SendMessageInput'];

export const login = async (payload: LoginBody) => {
  const { data } = await apiClient.post<IdentifierResponse>('/session', payload);
  return data;
};

export const fetchHealth = async () => {
  const { data } = await apiClient.get<HealthResponse>('/health');
  return data;
};

export const updateUsername = async (payload: SetNameBody) => {
  const { data } = await apiClient.post<UpdateUsernameResponse>('/user/username', payload);
  return data;
};

export const fetchConversations = async () => {
  const { data } = await apiClient.get<ConversationListPayload>('/conversations');
  return {
    conversations: data.conversations ?? []
  };
};

export const fetchConversation = async (conversationId: string) => {
  const { data } = await apiClient.get<ConversationResponse>(`/conversations/${conversationId}`);
  return data.conversation;
};

export const sendMessage = async (conversationId: string, payload: SendMessageInput) => {
  const { data } = await apiClient.post<Message>(`/conversations/${conversationId}/messages`, payload);
  return data;
};

export const deleteMessage = async (messageId: string) => {
  await apiClient.delete(`/messages/${messageId}`);
};
