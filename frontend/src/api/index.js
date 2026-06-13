import axios from 'axios'
import { error as showError, warning as showWarning } from '../utils/modal'

const api = axios.create({
  baseURL: '/api',
  timeout: 600000
})

// 请求拦截器
api.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
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
    // 检查响应数据是否正常
    if (!response.data || typeof response.data !== 'object') {
      showError('服务器返回的数据格式不正确', { title: '数据解析错误' })
      return Promise.reject(new Error('Invalid response data'))
    }

    // 检查是否有错误字段
    if (response.data.error) {
      showError(response.data.error, { title: '操作失败' })
      return Promise.reject(new Error(response.data.error))
    }

    return response
  },
  error => {
    const errorMsg = getErrorMessage(error)

    // 401 未授权 - 跳转到登录页
    if (error.response && error.response.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
      return Promise.reject(error)
    }

    // 403 禁止访问
    if (error.response && error.response.status === 403) {
      showError('您没有权限执行此操作', { title: '权限不足' })
      return Promise.reject(error)
    }

    // 404 资源不存在
    if (error.response && error.response.status === 404) {
      showWarning('请求的资源不存在', { title: '资源未找到' })
      return Promise.reject(error)
    }

    // 其他错误
    showError(errorMsg, { title: '操作失败' })
    return Promise.reject(error)
  }
)

// 获取错误信息
function getErrorMessage(error) {
  // 网络错误
  if (!error.response) {
    if (error.code === 'ECONNREFUSED') {
      return '无法连接到服务器，请稍后重试'
    }
    if (error.code === 'ECONNABORTED' || error.message.includes('timeout')) {
      return '请求超时，请稍后重试'
    }
    if (error.message.includes('Network Error')) {
      return '网络连接失败，请检查网络设置'
    }
    return '网络错误，请稍后重试'
  }

  // 服务器返回的错误
  const response = error.response

  // 尝试解析响应数据中的错误信息
  if (response.data) {
    if (typeof response.data === 'string') {
      return response.data
    }
    if (response.data.error) {
      return response.data.error
    }
    if (response.data.message) {
      return response.data.message
    }
  }

  // 根据状态码返回通用错误信息
  const statusText = {
    400: '请求参数错误',
    401: '登录已过期，请重新登录',
    403: '您没有权限执行此操作',
    404: '请求的资源不存在',
    408: '请求超时',
    409: '资源冲突，请稍后重试',
    500: '服务器内部错误',
    502: '网关错误',
    503: '服务暂时不可用',
    504: '网关超时'
  }

  return statusText[response.status] || `请求失败，状态码: ${response.status}`
}

export default api

// 导出所有API模块
export { authApi } from './auth'
export { articleApi, reportApi } from './article'
export { commentApi } from './comment'
export { userApi, signinApi } from './user'
export { friendApi } from './friend'
export { chatApi } from './chat'
export { notificationApi, userNotificationApi, adminNotificationApi } from './notification'
export { commonApi } from './common'
export { uploadApi } from './upload'
export { default as adminApi, adminUserApi, adminArticleApi, adminCommentApi, adminStatsApi, adminSidebarApi, adminDeletionApi, adminAnnouncementApi, adminSiteConfigApi, adminNotificationApi, adminTitleApi, adminCategoryApi, adminAuthApi, userNotificationApi, permissionGroupApi } from './admin'
