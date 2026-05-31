import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import fs from 'fs'
import path from 'path'

let version = 'unknown'
try {
  const versionPath = path.resolve(__dirname, '../version.json')
  const versionContent = fs.readFileSync(versionPath, 'utf-8')
  const versionData = JSON.parse(versionContent)
  version = versionData.frontend?.version || 'unknown'
} catch (error) {
  console.error('读取 version.json 失败:', error)
}

export default defineConfig({
  plugins: [vue()],
  define: {
    __FRONTEND_VERSION__: JSON.stringify(version)
  },
  server: {
    port: 3000
  },
  build: {
    assetsDir: 'assets',
    rollupOptions: {
      output: {
        manualChunks: undefined
      }
    }
  }
})
