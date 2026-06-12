import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

// 响应拦截器
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// ============ 公开接口 ============

// 认证相关（公开）
export const authApi = {
  register: (data) => api.post('/auth/register', data),
  login: (data) => api.post('/auth/login', data),
  initAdmin: (data) => api.post('/auth/init-admin', data),
  checkInit: () => api.get('/auth/check-init')
}

// 文章相关（公开）
export const articleApi = {
  // 公开
  getArticles: (params) => api.get('/articles', { params }),
  getArticle: (id) => api.get(`/articles/${id}`),
  searchArticles: (params) => api.get('/articles/search', { params }),
  // 需认证
  create: (data) => api.post('/articles', data),
  update: (id, data) => api.put(`/articles/${id}`, data),
  delete: (id) => api.delete(`/articles/${id}`),
  restore: (id) => api.post(`/articles/${id}/restore`),
  like: (id) => api.post(`/articles/${id}/like`),
  unlike: (id) => api.delete(`/articles/${id}/like`),
  getMyArticles: (params) => api.get('/my/articles', { params }),
  share: (id) => api.post(`/articles/${id}/share`),
  getDrafts: (params) => api.get('/my/drafts', { params }),
  publishDraft: (id) => api.post(`/articles/${id}/publish`),
  pin: (id) => api.post(`/articles/${id}/pin`),
  unpin: (id) => api.delete(`/articles/${id}/pin`)
}

// 分类相关（公开）
export const categoryApi = {
  getCategories: () => api.get('/categories'),
  create: (data) => api.post('/categories', data),
  update: (id, data) => api.put(`/categories/${id}`, data),
  delete: (id) => api.delete(`/categories/${id}`)
}

// 用户相关（公开）
export const userApi = {
  getUserByID: (id) => api.get(`/users/${id}`),
  getUserArticles: (id, params) => api.get(`/users/${id}/articles`, { params })
}

// 站点配置（公开）
export const siteApi = {
  getAnnouncement: () => api.get('/announcement'),
  getSidebarConfig: () => api.get('/sidebar-config'),
  getSiteConfig: () => api.get('/site-config'),
  getVersion: () => api.get('/version')
}

// 密码重置（公开）
export const passwordApi = {
  sendResetCode: (data) => api.post('/password/reset-code', data),
  resetPassword: (data) => api.post('/password/reset', data)
}

// ============ 需认证接口 ============

// 个人信息
export const profileApi = {
  getProfile: () => api.get('/profile'),
  updateProfile: (data) => api.put('/profile', data),
  uploadAvatar: (formData) => api.post('/upload/avatar', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

// 上传
export const uploadApi = {
  uploadImage: (formData) => api.post('/upload/image', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  }),
  uploadVideo: (formData) => api.post('/upload/video', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  }),
  uploadVoice: (formData) => api.post('/upload/voice', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

// 评论
export const commentApi = {
  create: (articleId, data) => api.post(`/articles/${articleId}/comments`, data),
  delete: (id) => api.delete(`/comments/${id}`),
  like: (id) => api.post(`/comments/${id}/like`),
  unlike: (id) => api.delete(`/comments/${id}/like`)
}

// 收藏
export const favoriteApi = {
  add: (id) => api.post(`/articles/${id}/favorite`),
  remove: (id) => api.delete(`/articles/${id}/favorite`),
  getFavorites: (params) => api.get('/favorites', { params }),
  check: (id) => api.get(`/articles/${id}/favorite/check`)
}

// 签到
export const signinApi = {
  signin: () => api.post('/signin'),
  getStatus: () => api.get('/signin/status'),
  getHistory: (params) => api.get('/signin/history', { params })
}

// 举报
export const reportApi = {
  create: (data) => api.post('/reports', data),
  getReports: () => api.get('/reports'),
  getReport: (id) => api.get(`/reports/${id}`),
  handle: (id, data) => api.put(`/reports/${id}/handle`, data)
}

// 通知
export const notificationApi = {
  getNotifications: () => api.get('/notifications'),
  getUnreadCount: () => api.get('/notifications/unread-count'),
  markRead: (id) => api.post(`/notifications/${id}/read`),
  markAllRead: () => api.post('/notifications/read-all'),
  delete: (id) => api.delete(`/notifications/${id}`),
  // 评论回复通知
  getCommentReplies: () => api.get('/comment-reply-notifications'),
  markCommentReplyRead: (id) => api.post(`/comment-reply-notifications/${id}/read`),
  markAllCommentRepliesRead: () => api.post('/comment-reply-notifications/read-all'),
  // 用户个人通知
  getUserNotifications: () => api.get('/user-notifications'),
  getUserNotification: (id) => api.get(`/user-notifications/${id}`),
  markUserNotificationRead: (id) => api.post(`/user-notifications/${id}/read`),
  markAllUserNotificationsRead: () => api.post('/user-notifications/read-all'),
  deleteUserNotification: (id) => api.delete(`/user-notifications/${id}`),
  clearUserNotifications: () => api.delete('/user-notifications/clear')
}

// 好友
export const friendApi = {
  sendRequest: (friendId, message) => api.post('/friends/request', { friend_id: friendId, message }),
  acceptRequest: (id) => api.post(`/friends/request/${id}/accept`),
  rejectRequest: (id) => api.post(`/friends/request/${id}/reject`),
  deleteFriend: (id) => api.delete(`/friends/${id}`),
  getFriends: () => api.get('/friends'),
  getRequests: () => api.get('/friends/requests'),
  getSentRequests: () => api.get('/friends/sent-requests'),
  updateDisplayName: (id, displayName) => api.put(`/friends/${id}/display-name`, { display_name: displayName }),
  checkStatus: (id) => api.get(`/friends/status/${id}`),
  getMutual: (id) => api.get(`/friends/mutual/${id}`)
}

// 聊天
export const chatApi = {
  getConversations: () => api.get('/chat/conversations'),
  getMessages: (params) => api.get('/chat/messages', { params }),
  sendMessage: (data) => api.post('/chat/messages', data),
  createPrivateConversation: (userId) => api.post('/chat/conversations/private', { user_id: userId }),
  getUnreadCount: () => api.get('/chat/unread'),
  markConversationRead: (id) => api.post(`/chat/conversations/${id}/read`)
}

// 头衔
export const titleApi = {
  getAll: () => api.get('/titles'),
  create: (data) => api.post('/titles', data),
  update: (id, data) => api.put(`/titles/${id}`, data),
  delete: (id) => api.delete(`/titles/${id}`),
  grant: (userId, titleId) => api.post('/titles/grant', { user_id: userId, title_id: titleId }),
  revoke: (userId, titleId) => api.post('/titles/revoke', { user_id: userId, title_id: titleId }),
  getUserTitles: (userId) => api.get(`/users/${userId}/titles`)
}

// 权限组
export const permissionGroupApi = {
  getGroups: () => api.get('/permission-groups'),
  getGroup: (id) => api.get(`/permission-groups/${id}`),
  create: (data) => api.post('/permission-groups', data),
  update: (id, data) => api.put(`/permission-groups/${id}`, data),
  delete: (id) => api.delete(`/permission-groups/${id}`),
  grant: (data) => api.post('/permission-groups/grant', data),
  revokeUser: (groupId, userId) => api.delete(`/permission-groups/${groupId}/revoke-user/${userId}`),
  getUserGroups: (userId) => api.get(`/users/${userId}/permission-groups`),
  checkPermissions: (permissions) => api.get('/permissions/check', { params: { permissions } }),
  initDefaults: () => api.post('/permission-groups/init')
}

// 用户状态
export const userStatusApi = {
  updateStatus: (data) => api.post('/user/status', data),
  getStatus: (id) => api.get(`/user/status/${id}`)
}

// 系统日志
export const logApi = {
  getMyLogs: () => api.get('/my-logs')
}

// ============ 管理员接口 ============

export const adminApi = {
  // 检查权限
  checkAdmin: () => api.get('/admin/check'),

  // 统计
  getStatistics: () => api.get('/admin/statistics'),

  // 用户管理
  getUsers: () => api.get('/admin/users'),
  updateUser: (id, data) => api.put(`/admin/users/${id}`, data),
  updateUserRole: (id, data) => api.put(`/admin/users/${id}/role`, data),
  banUser: (id, data) => api.post(`/admin/users/${id}/ban`, data),
  unbanUser: (id) => api.post(`/admin/users/${id}/unban`),
  deleteUser: (id) => api.delete(`/admin/users/${id}`),

  // 文章管理
  getArticles: () => api.get('/admin/articles'),
  updateArticleStatus: (id, data) => api.put(`/admin/articles/${id}/status`, data),

  // 评论管理
  getComments: () => api.get('/admin/comments'),
  deleteComment: (id) => api.delete(`/admin/comments/${id}`),

  // 通知管理
  createNotification: (data) => api.post('/notifications', data),
  getAdminNotifications: () => api.get('/notifications/admin'),
  sendUserNotification: (data) => api.post('/user-notifications/send', data),
  sendBatchNotification: (data) => api.post('/user-notifications/send-batch', data),
  getAdminUserNotifications: (userId, params) => api.get(`/admin/user-notifications/${userId}`, { params }),

  // 侧边栏配置
  updateSidebarConfig: (data) => api.put('/sidebar-config', data),

  // 删除审核
  getDeletionRequests: () => api.get('/deletion-requests'),
  approveDeletion: (id) => api.post(`/deletion-requests/${id}/approve`),
  rejectDeletion: (id) => api.post(`/deletion-requests/${id}/reject`),

  // 公告
  updateAnnouncement: (data) => api.put('/announcement', data),

  // 网站配置
  updateSiteConfig: (data) => api.put('/site-config', data),
  testSmtp: (data) => api.post('/site-config/test-smtp', data),

  // 用户状态管理
  getAllUserStatuses: () => api.get('/users/status'),
  getOnlineUsers: () => api.get('/users/online'),
  cleanupUserStatuses: () => api.post('/users/status/cleanup'),

  // 系统日志
  getSystemLogs: (params) => api.get('/system-logs', { params }),
  getLogModules: () => api.get('/system-logs/modules'),
  deleteOldLogs: () => api.delete('/system-logs/old')
}

// 导出所有API
export default {
  api,
  auth: authApi,
  article: articleApi,
  category: categoryApi,
  user: userApi,
  site: siteApi,
  password: passwordApi,
  profile: profileApi,
  upload: uploadApi,
  comment: commentApi,
  favorite: favoriteApi,
  signin: signinApi,
  report: reportApi,
  notification: notificationApi,
  friend: friendApi,
  chat: chatApi,
  title: titleApi,
  permissionGroup: permissionGroupApi,
  userStatus: userStatusApi,
  log: logApi,
  admin: adminApi
}
