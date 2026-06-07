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
}
