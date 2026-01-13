import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

// https://vite.dev/config/
export default defineConfig(({ command }) => ({
  plugins: [vue(), command === 'serve' ? vueDevTools() : null].filter(Boolean),
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  build: {
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (!id.includes('node_modules')) return
          if (id.includes('element-plus') || id.includes('@element-plus')) return 'element-plus'
          if (id.includes('vue-router') || id.includes('/vue/') || id.includes('pinia')) return 'vue'
          if (id.includes('axios')) return 'axios'
          if (id.includes('dayjs')) return 'dayjs'
          return 'vendor'
        }
      }
    }
  },
  server: {
    proxy: {
      '/uploads': {
        target: 'http://118.89.187.189:8081',
        changeOrigin: true
      }
    }
  }
}))
