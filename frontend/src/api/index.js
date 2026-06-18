import axios from 'axios'
import { error as showError, warning as showWarning } from '../utils/modal'

const api = axios.create({
  baseURL: '/api',
  timeout: 600000
})

// 是否正在刷新 token
let isRefreshing = false
// 等待刷新 token 的请求队列
let refreshQueue = []

// 刷新 token
async function refreshToken() {
  const refreshToken = localStorage.getItem('refresh_token')
  if (!refreshToken) {
    return false
  }

  try {
    const response = await axios.post('/api/auth/refresh-token', {
      refresh_token: refreshToken
    })

    const { token, refresh_token } = response.data
    localStorage.setItem('token', token)
    localStorage.setItem('refresh_token', refresh_token)
    localStorage.setItem('token_expires_at', Date.now() + 3600 * 1000)

    return true
  } catch (error) {
    localStorage.removeItem('token')
    localStorage.removeItem('refresh_token')
    localStorage.removeItem('user')
    localStorage.removeItem('token_expires_at')
    return false
  }
}

// 检查 token 是否即将过期（5分钟内）
function isTokenExpiringSoon() {
  const expiresAt = localStorage.getItem('token_expires_at')
  if (!expiresAt) return false

  const now = Date.now()
  const expires = parseInt(expiresAt)
  return expires - now < 5 * 60 * 1000 // 5分钟内过期
}

// 请求拦截器
api.interceptors.request.use(
  async config => {
    const token = localStorage.getItem('token')
    if (token) {
      // 检查 token 是否即将过期
      if (isTokenExpiringSoon() && !isRefreshing && config.url !== '/auth/refresh-token') {
        isRefreshing = true
        const success = await refreshToken()
        isRefreshing = false

        if (!success) {
          // 刷新失败，跳转到登录页
          window.location.href = '/login'
          return Promise.reject(new Error('Token refresh failed'))
        }

        // 处理等待队列
        refreshQueue.forEach(callback => callback())
        refreshQueue = []
      }

      // 如果正在刷新，等待刷新完成
      if (isRefreshing && config.url !== '/auth/refresh-token') {
        return new Promise((resolve) => {
          refreshQueue.push(() => {
            config.headers.Authorization = `Bearer ${localStorage.getItem('token')}`
            resolve(config)
          })
        })
      }

      config.headers.Authorization = `Bearer ${localStorage.getItem('token')}`
    }
    return config
  },
  error => {
    // 请求配置错误
    const errorMsg = getErrorMessage(error)
    showError(errorMsg, { title: '请求错误' })
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  response => {
    return response
  },
  async error => {
    const errorInfo = parseErrorResponse(error)

    // 401 未授权 - 尝试刷新 token
    if (errorInfo.code === 401) {
      // 如果不是刷新 token 的请求失败
      if (error.config.url !== '/auth/refresh-token') {
        // 尝试刷新 token
        const success = await refreshToken()
        if (success) {
          // 重新发送原请求
          error.config.headers.Authorization = `Bearer ${localStorage.getItem('token')}`
          return axios.request(error.config)
        }
      }

      // 刷新失败，跳转到登录页
      localStorage.removeItem('token')
      localStorage.removeItem('refresh_token')
      localStorage.removeItem('user')
      localStorage.removeItem('token_expires_at')
      window.location.href = '/login'
      return Promise.reject(error)
    }

    // 403 禁止访问
    if (errorInfo.code === 403) {
      showError(errorInfo.message, { title: '权限不足', detail: errorInfo.detail })
      return Promise.reject(error)
    }

    // 404 资源不存在
    if (errorInfo.code === 404) {
      showWarning(errorInfo.message, { title: '资源未找到', detail: errorInfo.detail })
      return Promise.reject(error)
    }

    // 429 限流
    if (errorInfo.code === 429) {
      showWarning(errorInfo.message || '请求过于频繁，请稍后再试', { title: '请求限流', detail: errorInfo.detail })
      return Promise.reject(error)
    }

    // 409 冲突
    if (errorInfo.code === 409) {
      showError(errorInfo.message, { title: '资源冲突', detail: errorInfo.detail })
      return Promise.reject(error)
    }

    // 其他错误
    showError(errorInfo.message, { title: '操作失败', detail: errorInfo.detail })
    return Promise.reject(error)
  }
)

// 解析错误响应
function parseErrorResponse(error) {
  const result = {
    code: 500,
    error: 'INTERNAL_ERROR',
    message: '服务器内部错误',
    detail: ''
  }

  // 网络错误
  if (!error.response) {
    if (error.code === 'ECONNREFUSED') {
      result.message = '无法连接到服务器，请稍后重试'
      result.detail = '服务器可能未启动或网络不可达'
    } else if (error.code === 'ECONNABORTED' || error.message.includes('timeout')) {
      result.message = '请求超时，请稍后重试'
      result.detail = '服务器响应时间过长'
    } else if (error.message.includes('Network Error')) {
      result.message = '网络连接失败，请检查网络设置'
      result.detail = '请确保您的设备已连接到网络'
    } else {
      result.message = '网络错误，请稍后重试'
    }
    return result
  }

  // 服务器返回的错误
  const response = error.response
  result.code = response.status

  // 尝试解析统一错误响应格式
  if (response.data && typeof response.data === 'object') {
    if (response.data.error) {
      result.error = response.data.error
    }
    if (response.data.message) {
      result.message = response.data.message
    }
    if (response.data.detail) {
      result.detail = response.data.detail
    }
  }

  // 如果没有获取到消息，使用状态码对应的默认消息
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

// 导出所有API模块
export { authApi } from './auth'
export { articleApi, reportApi } from './article'
export { commentApi } from './comment'
export { userApi, signinApi } from './user'
export { friendApi } from './friend'
export { notificationApi, userNotificationApi, adminUserNotificationApi } from './notification'
export { commonApi } from './common'
export { uploadApi } from './upload'
export { default as adminApi, adminUserApi, adminArticleApi, adminCommentApi, adminStatsApi, adminSidebarApi, adminDeletionApi, adminAnnouncementApi, adminSiteConfigApi, adminNotificationApi, adminTitleApi, adminCategoryApi, adminAuthApi, permissionGroupApi } from './admin'
