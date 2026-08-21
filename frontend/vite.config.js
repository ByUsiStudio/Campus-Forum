import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import fs from 'fs'
import path from 'path'

// 优先从命令行参数获取版本号
let version = 'unknown'

// 从命令行参数解析版本号 (e.g., --version=1.3.5)
const args = process.argv.slice(2)
for (const arg of args) {
  if (arg.startsWith('--version=')) {
    version = arg.split('=')[1]
    break
  }
}

// 如果命令行没有传入，尝试读取 version.json
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
