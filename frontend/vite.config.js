import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import fs from 'fs'
import path from 'path'

let version = 'unknown'
const args = process.argv.slice(2)
for (const arg of args) {
  if (arg.startsWith('--version=')) {
    version = arg.split('=')[1]
    break
  }
}

if (version === 'unknown') {
  try {
    const versionPath = path.resolve(__dirname, '../version.json')
    const versionContent = fs.readFileSync(versionPath, 'utf-8')
    const versionData = JSON.parse(versionContent)
    version = versionData.frontend?.version || 'unknown'
  } catch (error) {
    console.error('读取 version.json 失败:', error)
  }
}

console.log(`Frontend version: ${version}`)

export default defineConfig({
  plugins: [vue()],
  define: {
    __FRONTEND_VERSION__: JSON.stringify(version)
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src'),
      '@styles': path.resolve(__dirname, 'src/styles')
    }
  },
  css: {
    preprocessorOptions: {
      less: {
        additionalData: `@import "${path.resolve(__dirname, '../t.less')}";`,
        javascriptEnabled: true
      }
    }
  },
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://0.0.0.0:3620',
        changeOrigin: true
      }
    }
  },
  build: {
    outDir: '../backend/dist',
    assetsDir: 'assets',
    rollupOptions: {
      output: {
        manualChunks: undefined
      }
    }
  }
})
