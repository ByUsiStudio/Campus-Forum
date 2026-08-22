import { defineStore } from 'pinia'
import http from '@/api/http'
import { ref, computed } from 'vue'

/**
 * 用户 / 认证状态
 */
export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref(null)
  const ready = ref(false)

  const isLoggedIn = computed(() => !!token.value)
  const isAdmin = computed(() => {
    const role = user.value?.role || ''
    return role === 'admin' || role === 'Admin'
  })

  function persist(data) {
    token.value = data.token || ''
    if (data.refresh_token) localStorage.setItem('refresh_token', data.refresh_token)
    if (data.expires_in) localStorage.setItem('token_expires_at', Date.now() + data.expires_in * 1000)
    localStorage.setItem('token', token.value)
    if (data.user) {
      user.value = data.user
      localStorage.setItem('user', JSON.stringify(data.user))
    }
  }

  function setUser(data) {
    user.value = data
    if (data) localStorage.setItem('user', JSON.stringify(data))
    else localStorage.removeItem('user')
  }

  /** 拉取最新完整资料并以之为准 */
  async function refreshProfile() {
    if (!token.value) return null
    const { data } = await http.get('/profile')
    setUser(data)
    return data
  }

  /** 初始化：从本地缓存恢复用户，再向后端拉取最新资料 */
  async function initUser() {
    const cached = localStorage.getItem('user')
    if (cached) {
      try { user.value = JSON.parse(cached) } catch (e) { /* ignore */ }
    }
    if (token.value) {
      try {
        const { data } = await http.get('/profile')
        setUser(data)
        ready.value = true
      } catch (e) {
        // 未授权：交给 http 层统一处理
        ready.value = true
      }
    } else {
      ready.value = true
    }
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('refresh_token')
    localStorage.removeItem('user')
    localStorage.removeItem('token_expires_at')
  }

  return { token, user, ready, isLoggedIn, isAdmin, persist, setUser, refreshProfile, initUser, logout }
})
