import axios from 'axios'
import { toast as jcToast } from '@/utils/message'

/**
 * 统一 HTTP 层
 * - axios 实例，baseURL = /api
 * - 自动携带 / 刷新 token（并发去重）
 * - 统一错误解析与提示（JCuPupw Toast）
 */
const http = axios.create({
  baseURL: '/api',
  timeout: 600000
})

let isRefreshing = false
let refreshQueue = []

function getToken() {
  return localStorage.getItem('token') || ''
}

function clearAuth() {
  localStorage.removeItem('token')
  localStorage.removeItem('refresh_token')
  localStorage.removeItem('user')
  localStorage.removeItem('token_expires_at')
}

async function refreshToken() {
  const refreshToken = localStorage.getItem('refresh_token')
  if (!refreshToken) return false
  try {
    const { data } = await axios.post('/api/auth/refresh-token', { refresh_token: refreshToken })
    if (!data.token) return false
    localStorage.setItem('token', data.token)
    if (data.refresh_token) localStorage.setItem('refresh_token', data.refresh_token)
    localStorage.setItem('token_expires_at', Date.now() + 3600 * 1000)
    return true
  } catch (e) {
    clearAuth()
    return false
  }
}

function attachToken(config) {
  const t = getToken()
  if (t) config.headers.Authorization = `Bearer ${t}`
  return config
}

function redirectToLogin() {
  if (window.location.pathname !== '/login') window.location.href = '/login'
}

http.interceptors.request.use(
  async config => {
    attachToken(config)
    const expiresAt = localStorage.getItem('token_expires_at')
    const expiresSoon = expiresAt && (parseInt(expiresAt, 10) - Date.now()) < 5 * 60 * 1000

    if (config.url !== '/auth/refresh-token' && expiresSoon) {
      if (isRefreshing) {
        return new Promise(resolve => refreshQueue.push(() => resolve(attachToken(config))))
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
    jcToast(error?.message || '请求配置错误', { type: 'error' })
    return Promise.reject(error)
  }
)

http.interceptors.response.use(
  response => response,
  async error => {
    const info = parseError(error)

    if (info.code === 401 && error.config && error.config.url !== '/auth/refresh-token') {
      if (isRefreshing) await new Promise(resolve => refreshQueue.push(resolve))
      const ok = await refreshToken()
      if (ok) {
        error.config.headers.Authorization = `Bearer ${getToken()}`
        return axios.request(error.config)
      }
      clearAuth()
      redirectToLogin()
      return Promise.reject(error)
    }

    switch (info.code) {
      case 403: jcToast(info.message || '您没有权限执行此操作', { type: 'error' }); break
      case 404: jcToast(info.message || '请求的资源不存在', { type: 'warning' }); break
      case 429: jcToast(info.message || '请求过于频繁，请稍后再试', { type: 'warning' }); break
      case 409: jcToast(info.message || '资源冲突', { type: 'error' }); break
      default:
        if (info.code !== 401) jcToast(info.message || '操作失败', { type: 'error' })
    }

    return Promise.reject(error)
  }
)

function parseError(error) {
  const result = { code: 500, message: '服务器内部错误' }
  if (!error.response) {
    const msg = error?.message || ''
    if (/Network Error/i.test(msg)) {
      result.code = 0
      result.message = '网络连接失败，请检查网络设置'
    } else if (/timeout|ECONNABORTED/i.test(msg)) {
      result.message = '请求超时，请稍后重试'
    } else {
      result.code = 0
      result.message = '网络错误，请稍后重试'
    }
    return result
  }
  const res = error.response
  result.code = res.status
  const data = res.data
  if (data && typeof data === 'object') {
    if (data.message) result.message = data.message
    else if (data.error) result.message = data.error
  }
  if (!result.message || result.message === '服务器内部错误') {
    const map = {
      400: '请求参数错误', 401: '登录已过期，请重新登录', 403: '您没有权限执行此操作',
      404: '请求的资源不存在', 408: '请求超时', 409: '资源冲突', 429: '请求过于频繁',
      500: '服务器内部错误', 502: '网关错误', 503: '服务暂时不可用', 504: '网关超时'
    }
    result.message = map[res.status] || `请求失败，状态码: ${res.status}`
  }
  return result
}

export default http
