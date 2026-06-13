import api from './index'

// 认证相关API
export const authApi = {
  // 用户注册
  register: (data) => api.post('/auth/register', data),

  // 用户登录
  login: (data) => api.post('/auth/login', data),

  // 初始化管理员
  initAdmin: (data) => api.post('/auth/init-admin', data),

  // 检查系统是否已初始化
  checkInit: () => api.get('/auth/check-init'),

  // 发送密码重置验证码
  sendResetCode: (data) => api.post('/password/reset-code', data),

  // 重置密码
  resetPassword: (data) => api.post('/password/reset', data),

  // 获取当前用户信息
  getProfile: () => api.get('/profile'),

  // 更新当前用户信息
  updateProfile: (data) => api.put('/profile', data),
}
