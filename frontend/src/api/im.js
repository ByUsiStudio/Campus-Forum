import api from './index'

// 野火IM配置获取
export const getIMConfig = () => api.get('/im/config')

// 获取IM Token
export const getIMToken = (deviceId = 'web') => api.get('/im/token', { params: { device_id: deviceId } })

// 创建IM用户
export const createIMUser = (data) => api.post('/im/users', data)

// 兼容旧接口的API封装
export const imApi = {
  // 获取IM配置
  getConfig: () => getIMConfig(),

  // 获取Token
  getToken: (deviceId) => getIMToken(deviceId),

  // 创建用户
  createUser: (data) => createIMUser(data),

  // 获取会话列表（通过后端代理）
  getConversations: () => api.get('/im/conversations'),

  // 获取消息历史（通过后端代理）
  getMessages: (params) => api.get('/im/messages', { params }),

  // 发送消息（通过后端代理）
  sendMessage: (data) => api.post('/im/messages', data),
}
