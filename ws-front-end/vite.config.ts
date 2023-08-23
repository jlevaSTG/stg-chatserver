import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/

export default defineConfig(({ mode }) => {
  console.log("mode =", mode)
  return {
    plugins: [react()],
    server: {
      proxy: mode === 'development' ? {

        '/api': {
          target: 'http://localhost:8080',
          changeOrigin: true,
          ws: false // You can explicitly set this to false for clarity, although it's the default
        },
        '/ws': {
          target: 'ws://localhost:8080',
          ws: true
        },
      } : {}
    },
  }
});