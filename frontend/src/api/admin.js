import api from './index'

// ==================== 用户管理API ====================
export const adminUserApi = {
  // 获取所有用户
  getUsers: () => api.get('/admin/users'),
  
  // 更新用户信息
  updateUser: (userId, data) => api.put(`/admin/users/${userId}`, data),
  
  // 更新用户角色
  updateUserRole: (userId, role) => api.put(`/admin/users/${userId}/role`, { role }),
  
  // 封禁用户
  banUser: (userId, reason) => api.post(`/admin/users/${userId}/ban`, { reason }),
  
  // 解封用户
  unbanUser: (userId) => api.post(`/admin/users/${userId}/unban`),
  
  // 删除用户
  deleteUser: (userId) => api.delete(`/admin/users/${userId}`)
}

// ==================== 文章管理API ====================
export const adminArticleApi = {
  // 获取所有文章
  getArticles: (params = {}) => api.get('/admin/articles', { params }),
  
  // 更新文章状态
  updateArticleStatus: (articleId, status) => 
    api.put(`/admin/articles/${articleId}/status`, { status })
}

// ==================== 评论管理API ====================
export const adminCommentApi = {
  // 获取所有评论
  getComments: (params = {}) => api.get('/admin/comments', { params }),
  
  // 删除评论
  deleteComment: (commentId) => api.delete(`/admin/comments/${commentId}`)
}

// ==================== 统计数据API ====================
export const adminStatsApi = {
  // 获取统计数据
  getStatistics: () => api.get('/admin/statistics')
}

// ==================== 侧边栏配置API ====================
export const adminSidebarApi = {
  // 获取侧边栏配置
  getConfig: () => api.get('/sidebar-config'),
  
  // 更新侧边栏配置
  updateConfig: (config) => api.put('/sidebar-config', config)
}

// ==================== 删除申请API ====================
export const adminDeletionApi = {
  // 获取删除申请列表
  getRequests: () => api.get('/deletion-requests'),
  
  // 批准删除申请
  approveRequest: (requestId) => api.post(`/deletion-requests/${requestId}/approve`),
  
  // 拒绝删除申请
  rejectRequest: (requestId) => api.post(`/deletion-requests/${requestId}/reject`)
}

// ==================== 公告管理API ====================
export const adminAnnouncementApi = {
  // 获取公告
  getAnnouncement: () => api.get('/announcement'),
  
  // 更新公告
  updateAnnouncement: (content) => api.put('/announcement', { content })
}

// ==================== 网站配置API ====================
export const adminSiteConfigApi = {
  // 获取网站配置
  getConfig: () => api.get('/site-config'),
  
  // 更新网站配置
  updateConfig: (config) => api.put('/site-config', config),
  
  // 测试SMTP配置
  testSmtp: (config) => api.post('/site-config/test-smtp', config)
}

// ==================== 通知管理API ====================
export const adminNotificationApi = {
  // 获取管理员通知
  getNotifications: () => api.get('/notifications/admin'),
  
  // 创建通知
  createNotification: (data) => api.post('/notifications', data),
  
  // 删除通知
  deleteNotification: (notificationId) => api.delete(`/notifications/${notificationId}`)
}

// ==================== 头衔管理API ====================
export const adminTitleApi = {
  // 获取所有头衔
  getTitles: () => api.get('/titles'),
  
  // 创建头衔
  createTitle: (titleData) => api.post('/titles', titleData),
  
  // 更新头衔
  updateTitle: (titleId, titleData) => api.put(`/titles/${titleId}`, titleData),
  
  // 删除头衔
  deleteTitle: (titleId) => api.delete(`/titles/${titleId}`),
  
  // 授予头衔
  grantTitle: (data) => api.post('/titles/grant', data),
  
  // 撤销头衔
  revokeTitle: (data) => api.post('/titles/revoke', data),
  
  // 获取用户头衔
  getUserTitles: (userId) => api.get(`/users/${userId}/titles`)
}

// ==================== 分区管理API ====================
export const adminCategoryApi = {
  // 获取所有分区
  getCategories: () => api.get('/categories'),
  
  // 创建分区
  createCategory: (categoryData) => api.post('/categories', categoryData),
  
  // 更新分区
  updateCategory: (categoryId, categoryData) => api.put(`/categories/${categoryId}`, categoryData),
  
  // 删除分区
  deleteCategory: (categoryId) => api.delete(`/categories/${categoryId}`)
}

// ==================== 管理员权限检查API ====================
export const adminAuthApi = {
  // 检查管理员权限
  checkAdmin: () => api.get('/admin/check'),
  
  // 获取系统统计
  getSystemStats: () => api.get('/admin/statistics'),
  
  // 检查系统初始化状态
  checkInit: () => api.get('/auth/check-init')
}

// 导出所有API
export default {
  user: adminUserApi,
  article: adminArticleApi,
  comment: adminCommentApi,
  stats: adminStatsApi,
  sidebar: adminSidebarApi,
  deletion: adminDeletionApi,
  announcement: adminAnnouncementApi,
  siteConfig: adminSiteConfigApi,
  notification: adminNotificationApi,
  title: adminTitleApi,
  category: adminCategoryApi,
  auth: adminAuthApi
}
