import { defineConfig, resolveBaseUrl } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, "./src")
    }
  },
  devServer:{
    proxy: "http://localhost:8000/",
    host:"localhost",
    port:8000,
    client: {
      webSocketURL: 'ws://localhost:8000/monster_situation',
    },
    headers: {
      'Access-Control-Allow-Origin': '*',
    }
  }
})
