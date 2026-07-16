import api from './index'

// 通用API
export const commonApi = {
  // 获取版本信息
  getVersion: () => api.get('/version'),

  // 获取分区列表
  getCategories: () => api.get('/categories'),

  // 获取公告
  getAnnouncement: () => api.get('/announcement'),

  // 获取侧边栏配置
  getSidebarConfig: () => api.get('/sidebar-config'),

  // 获取网站配置
  getSiteConfig: () => api.get('/site-config'),

  // 获取所有头衔
  getTitles: () => api.get('/titles'),

  // 检查管理员权限
  checkAdmin: () => api.get('/admin/check'),

  // 用户登录
  login: (data) => api.post('/auth/login', data),

  // 检查系统是否已初始化
  checkInit: () => api.get('/auth/check-init'),

  // 用户注册
  register: (data) => api.post('/auth/register', data),

  // 获取用户信息
  getProfile: () => api.get('/profile'),

  // 更新用户信息
  updateProfile: (data) => api.put('/profile', data),
}
