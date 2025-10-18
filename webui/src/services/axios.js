import axios from 'axios';

const base = import.meta.env?.VITE_API_BASE_URL || '/api';

const instance = axios.create({
  baseURL: base,
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: false,
});

export default instance;
