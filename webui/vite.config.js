import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    proxy: {
      // proxy API calls during dev to your Go server on 3000
      '^/(session|conversations|liveness)': {
        target: 'http://localhost:3000',
        changeOrigin: true
      }
    }
  }
})
