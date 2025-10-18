import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// Dev server on :8000 and proxy /api and /v1 to the Go backend on :3000
export default defineConfig({
  plugins: [vue()],
  server: {
    port: 8000,
    strictPort: true,
    proxy: {
      '^/(api|v1)(/|$)': 'http://localhost:3000'
    }
  }
})
