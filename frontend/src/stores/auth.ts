import { computed, ref } from 'vue';
import { defineStore } from 'pinia';

import { login as loginRequest } from '@/api';
import { onUnauthorized, setAuthToken } from '@/services/http';

const TOKEN_STORAGE_KEY = 'wasa_token';

const readToken = () => {
  if (typeof window === 'undefined') return null;
  return window.localStorage.getItem(TOKEN_STORAGE_KEY);
};

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(readToken());
  const displayName = ref<string | null>(null);
  const loading = ref(false);
  const error = ref<string | null>(null);

  const isAuthenticated = computed(() => Boolean(token.value));

  const persistToken = (value: string | null) => {
    token.value = value;
    setAuthToken(value);
    if (typeof window === 'undefined') {
      return;
    }
    if (value) {
      window.localStorage.setItem(TOKEN_STORAGE_KEY, value);
    } else {
      window.localStorage.removeItem(TOKEN_STORAGE_KEY);
    }
  };

  if (token.value) {
    setAuthToken(token.value);
  }

  const logout = () => {
    displayName.value = null;
    error.value = null;
    persistToken(null);
  };

  onUnauthorized(() => {
    logout();
  });

  const login = async (name: string) => {
    loading.value = true;
    error.value = null;
    try {
      const response = await loginRequest({ name });
      persistToken(response.identifier);
      displayName.value = name;
      return true;
    } catch (err) {
      console.error(err);
      error.value = 'Unable to log in. Please try again.';
      persistToken(null);
      return false;
    } finally {
      loading.value = false;
    }
  };

  return {
    token,
    displayName,
    loading,
    error,
    isAuthenticated,
    login,
    logout
  };
});
