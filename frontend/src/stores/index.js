import { ref, computed } from 'vue'

const user = ref(null)
const token = ref(null)
const refreshToken = ref(null)
const tokenExpiresAt = ref(null)
const siteConfig = ref(null)
const unreadCount = ref(0)
const sidebarConfig = ref(null)
const announcement = ref(null)

const STORAGE_KEYS = {
  USER: 'user',
  TOKEN: 'token',
  REFRESH_TOKEN: 'refresh_token',
  TOKEN_EXPIRES_AT: 'token_expires_at'
}

function loadFromStorage() {
  user.value = loadJson(STORAGE_KEYS.USER)
  token.value = localStorage.getItem(STORAGE_KEYS.TOKEN)
  refreshToken.value = localStorage.getItem(STORAGE_KEYS.REFRESH_TOKEN)
  
  const expiresAt = localStorage.getItem(STORAGE_KEYS.TOKEN_EXPIRES_AT)
  tokenExpiresAt.value = expiresAt ? parseInt(expiresAt) : null
}

function loadJson(key) {
  const value = localStorage.getItem(key)
  if (!value) return null
  try {
    return JSON.parse(value)
  } catch {
    return null
  }
}

function saveToStorage() {
  if (user.value) {
    localStorage.setItem(STORAGE_KEYS.USER, JSON.stringify(user.value))
  } else {
    localStorage.removeItem(STORAGE_KEYS.USER)
  }
  
  if (token.value) {
    localStorage.setItem(STORAGE_KEYS.TOKEN, token.value)
  } else {
    localStorage.removeItem(STORAGE_KEYS.TOKEN)
  }
  
  if (refreshToken.value) {
    localStorage.setItem(STORAGE_KEYS.REFRESH_TOKEN, refreshToken.value)
  } else {
    localStorage.removeItem(STORAGE_KEYS.REFRESH_TOKEN)
  }
  
  if (tokenExpiresAt.value) {
    localStorage.setItem(STORAGE_KEYS.TOKEN_EXPIRES_AT, tokenExpiresAt.value.toString())
  } else {
    localStorage.removeItem(STORAGE_KEYS.TOKEN_EXPIRES_AT)
  }
}

export function useStore() {
  const isLoggedIn = computed(() => !!token.value)
  
  const isAdmin = computed(() => {
    if (!user.value) return false
    return user.value.role === 'admin' || user.value.role === 'system_admin'
  })
  
  const hasMinLevel = computed(() => (level) => {
    if (!user.value) return false
    return user.value.level >= level
  })
  
  const isTokenExpiringSoon = computed(() => {
    if (!tokenExpiresAt.value) return false
    return tokenExpiresAt.value - Date.now() < 5 * 60 * 1000
  })

  const login = (userData, accessToken, refreshTokenValue, expiresIn) => {
    user.value = userData
    token.value = accessToken
    refreshToken.value = refreshTokenValue
    tokenExpiresAt.value = Date.now() + expiresIn * 1000
    saveToStorage()
  }

  const logout = () => {
    user.value = null
    token.value = null
    refreshToken.value = null
    tokenExpiresAt.value = null
    unreadCount.value = 0
    saveToStorage()
  }

  const updateUser = (userData) => {
    user.value = { ...user.value, ...userData }
    saveToStorage()
  }

  const setToken = (newToken, newRefreshToken, expiresIn) => {
    token.value = newToken
    refreshToken.value = newRefreshToken
    tokenExpiresAt.value = Date.now() + expiresIn * 1000
    saveToStorage()
  }

  const setSiteConfig = (config) => {
    siteConfig.value = config
  }

  const setUnreadCount = (count) => {
    unreadCount.value = count
  }

  const incrementUnreadCount = () => {
    unreadCount.value++
  }

  const decrementUnreadCount = () => {
    if (unreadCount.value > 0) {
      unreadCount.value--
    }
  }

  const clearUnreadCount = () => {
    unreadCount.value = 0
  }

  const setSidebarConfig = (config) => {
    sidebarConfig.value = config
  }

  const setAnnouncement = (announce) => {
    announcement.value = announce
  }

  const init = () => {
    loadFromStorage()
  }

  return {
    user,
    token,
    refreshToken,
    tokenExpiresAt,
    siteConfig,
    unreadCount,
    sidebarConfig,
    announcement,
    isLoggedIn,
    isAdmin,
    hasMinLevel,
    isTokenExpiringSoon,
    login,
    logout,
    updateUser,
    setToken,
    setSiteConfig,
    setUnreadCount,
    incrementUnreadCount,
    decrementUnreadCount,
    clearUnreadCount,
    setSidebarConfig,
    setAnnouncement,
    init
  }
}