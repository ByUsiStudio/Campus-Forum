import axios from 'axios'
import { error as showError, warning as showWarning } from '../utils/modal'
import { useStore } from '../stores'

const api = axios.create({
  baseURL: '/api',
  timeout: 600000
})

let isRefreshing = false
let refreshQueue = []

async function refreshToken() {
  const store = useStore()
  
  if (!store.refreshToken.value) {
    return false
  }

  try {
    const response = await axios.post('/api/auth/refresh-token', {
      refresh_token: store.refreshToken.value
    })

    const { token, refresh_token } = response.data
    store.setToken(token, refresh_token, 3600)

    return true
  } catch (error) {
    store.logout()
    return false
  }
}

api.interceptors.request.use(
  async config => {
    const store = useStore()
    
    if (store.token.value) {
      if (store.isTokenExpiringSoon.value && !isRefreshing && config.url !== '/auth/refresh-token') {
        isRefreshing = true
        const success = await refreshToken()
        isRefreshing = false

        if (!success) {
          window.location.href = '/login'
          return Promise.reject(new Error('Token refresh failed'))
        }

        refreshQueue.forEach(callback => callback())
        refreshQueue = []
      }

      if (isRefreshing && config.url !== '/auth/refresh-token') {
        return new Promise((resolve) => {
          refreshQueue.push(() => {
            config.headers.Authorization = `Bearer ${useStore().token.value}`
            resolve(config)
          })
        })
      }

      config.headers.Authorization = `Bearer ${store.token.value}`
    }
    return config
  },
  error => {
    const errorMsg = getErrorMessage(error)
    showError(errorMsg, { title: '请求错误' })
    return Promise.reject(error)
  }
)

api.interceptors.response.use(
  response => {
    return response
  },
  async error => {
    const errorInfo = parseErrorResponse(error)

    if (errorInfo.code === 401) {
      if (error.config.url !== '/auth/refresh-token') {
        const success = await refreshToken()
        if (success) {
          error.config.headers.Authorization = `Bearer ${useStore().token.value}`
          return axios.request(error.config)
        }
      }

      useStore().logout()
      window.location.href = '/login'
      return Promise.reject(error)
    }

    if (errorInfo.code === 403) {
      showError(errorInfo.message, { title: '权限不足', detail: errorInfo.detail })
      return Promise.reject(error)
    }

    if (errorInfo.code === 404) {
      showWarning(errorInfo.message, { title: '资源未找到', detail: errorInfo.detail })
      return Promise.reject(error)
    }

    if (errorInfo.code === 429) {
      showWarning(errorInfo.message || '请求过于频繁，请稍后再试', { title: '请求限流', detail: errorInfo.detail })
      return Promise.reject(error)
    }

    if (errorInfo.code === 409) {
      showError(errorInfo.message, { title: '资源冲突', detail: errorInfo.detail })
      return Promise.reject(error)
    }

    showError(errorInfo.message, { title: '操作失败', detail: errorInfo.detail })
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

  const response = error.response
  result.code = response.status

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

function getErrorMessage(error) {
  if (error.message) {
    return error.message
  }
  return '请求错误'
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