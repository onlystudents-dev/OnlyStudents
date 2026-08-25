import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'

export default defineConfig({
  plugins: [vue()],
  base: '/',
  server: { proxy: { '/api': 'http://localhost:8080' } },
})
