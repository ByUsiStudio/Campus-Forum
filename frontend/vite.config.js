import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import fs from 'fs'
import path from 'path'

let version = process.env.VITE_APP_VERSION || ''

if (!version || version === 'unknown') {
  try {
    const versionPath = path.resolve(__dirname, '../version.json')
    const versionContent = fs.readFileSync(versionPath, 'utf-8')
    const versionData = JSON.parse(versionContent)
    version = versionData.frontend?.version || 'unknown'
  } catch (error) {
    console.error('读取 version.json 失败:', error)
    version = 'unknown'
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
      '@': path.resolve(__dirname, 'src')
    }
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
