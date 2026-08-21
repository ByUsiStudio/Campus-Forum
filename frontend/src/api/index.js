import axios from 'axios'
import { error as showError, warning as showWarning } from '../utils/modal'

const api = axios.create({
  baseURL: '/api',
  timeout: 600000
})

let isRefreshing = false
let refreshQueue = []

async function refreshToken() {
  const stored = localStorage.getItem('refresh_token')
  if (!stored) return false

  try {
    const response = await axios.post('/api/auth/refresh-token', {
      refresh_token: stored
    })

    const { token, refresh_token } = response.data
    if (!token) return false

    localStorage.setItem('token', token)
    if (refresh_token) localStorage.setItem('refresh_token', refresh_token)
    localStorage.setItem('token_expires_at', Date.now() + 3600 * 1000)

    return true
  } catch (error) {
    clearAuth()
    return false
  }
}

function clearAuth() {
  localStorage.removeItem('token')
  localStorage.removeItem('refresh_token')
  localStorage.removeItem('user')
  localStorage.removeItem('token_expires_at')
}

function isTokenExpiringSoon() {
  const expiresAt = localStorage.getItem('token_expires_at')
  if (!expiresAt) return false
  const expires = parseInt(expiresAt, 10)
  if (isNaN(expires)) return false
  return expires - Date.now() < 5 * 60 * 1000
}

function redirectToLogin() {
  if (window.location.pathname !== '/login') {
    window.location.href = '/login'
  }
}

function attachToken(config) {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
}

api.interceptors.request.use(
  async config => {
    attachToken(config)
    if (config.url !== '/auth/refresh-token' && isTokenExpiringSoon()) {
      if (isRefreshing) {
        return new Promise(resolve => {
          refreshQueue.push(() => resolve(attachToken(config)))
        })
      }

      isRefreshing = true
      const ok = await refreshToken()
      isRefreshing = false

      refreshQueue.forEach(cb => cb())
      refreshQueue = []

      if (!ok) {
        redirectToLogin()
        return Promise.reject(new Error('Token expired'))
      }

      attachToken(config)
    }

    return config
  },
  error => {
    const message = error?.message || '请求配置错误'
    showError(message, { title: '请求错误' })
    return Promise.reject(error)
  }
)

api.interceptors.response.use(
  response => response,
  async error => {
    const errorInfo = parseErrorResponse(error)

    if (errorInfo.code === 401 && error.config && error.config.url !== '/auth/refresh-token') {
      if (isRefreshing) {
        await new Promise(resolve => refreshQueue.push(resolve))
      }

      const ok = await refreshToken()
      if (ok) {
        error.config.headers.Authorization = `Bearer ${localStorage.getItem('token')}`
        return axios.request(error.config)
      }

      clearAuth()
      redirectToLogin()
      return Promise.reject(error)
    }

    switch (errorInfo.code) {
      case 403:
        showError(errorInfo.message, { title: '权限不足', detail: errorInfo.detail })
        break
      case 404:
        showWarning(errorInfo.message, { title: '资源未找到', detail: errorInfo.detail })
        break
      case 429:
        showWarning(errorInfo.message || '请求过于频繁，请稍后再试', { title: '请求限流', detail: errorInfo.detail })
        break
      case 409:
        showError(errorInfo.message, { title: '资源冲突', detail: errorInfo.detail })
        break
      default:
        if (errorInfo.code !== 401) {
          showError(errorInfo.message, { title: '操作失败', detail: errorInfo.detail })
        }
    }

    return Promise.reject(error)
  }
)

function parseErrorResponse(error) {
  const result = {
    code: 500,
    error: 'INTERNAL_ERROR',
    message: '服务器内部错误',
    detail: ''
  }

  if (!error.response) {
    if (error.code === 'ECONNREFUSED') {
      result.message = '无法连接到服务器，请稍后重试'
      result.detail = '服务器可能未启动或网络不可达'
    } else if (error.code === 'ECONNABORTED' || /timeout/i.test(error.message || '')) {
      result.message = '请求超时，请稍后重试'
      result.detail = '服务器响应时间过长'
    } else if (/Network Error/i.test(error.message || '')) {
      result.message = '网络连接失败，请检查网络设置'
      result.detail = '请确保您的设备已连接到网络'
    } else {
      result.message = '网络错误，请稍后重试'
    }
    return result
  }

  const response = error.response
  result.code = response.status

  const data = response.data
  if (data && typeof data === 'object') {
    if (data.error) result.error = data.error
    if (data.message) result.message = data.message
    if (data.detail) result.detail = data.detail
  }

  if (!result.message || result.message === '服务器内部错误') {
    const statusMessages = {
      400: '请求参数错误',
      401: '登录已过期，请重新登录',
      403: '您没有权限执行此操作',
      404: '请求的资源不存在',
      408: '请求超时',
      409: '资源冲突',
      429: '请求过于频繁',
      500: '服务器内部错误',
      502: '网关错误',
      503: '服务暂时不可用',
      504: '网关超时'
    }
    result.message = statusMessages[response.status] || `请求失败，状态码: ${response.status}`
  }

  return result
}

export default api

export { authApi } from './auth'
export { articleApi, reportApi } from './article'
export { commentApi } from './comment'
export { userApi, signinApi } from './user'
export { friendApi } from './friend'
export { notificationApi, userNotificationApi, adminUserNotificationApi } from './notification'
export { commonApi } from './common'
export { uploadApi } from './upload'
export { default as adminApi, adminUserApi, adminArticleApi, adminCommentApi, adminStatsApi, adminSidebarApi, adminDeletionApi, adminAnnouncementApi, adminSiteConfigApi, adminNotificationApi, adminTitleApi, adminCategoryApi, adminAuthApi, permissionGroupApi } from './admin'
export { levelApi } from './level'
export { statisticsApi } from './statistics'
export { collectionApi } from './collection'
export { topicApi } from './topic'
export { leaderboardApi } from './leaderboard'
