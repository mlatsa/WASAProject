import axios from "axios";
const instance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE || "http://localhost:3000",
  timeout: 8000,
});
export default instance;
