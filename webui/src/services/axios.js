import axios from "axios";
const instance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE,
  timeout: 8000,
});
instance.interceptors.request.use((config) => {
  const t = localStorage.getItem("authToken") || localStorage.getItem("identifier");
  if (t) config.headers.Authorization = `Bearer ${t}`;
  return config;
});
export default instance;
