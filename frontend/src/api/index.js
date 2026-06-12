import axios from 'axios'

// 认证接口
export const authApi = {
  login: (data) => axios.post('/api/auth/login', data),
  register: (data) => axios.post('/api/auth/register', data),
  getProfile: () => axios.get('/api/profile'),
  updateProfile: (data) => axios.put('/api/profile', data),
  uploadAvatar: (formData) => axios.post('/api/upload/avatar', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  }),
  checkInit: () => axios.get('/api/auth/check-init'),
  initAdmin: (data) => axios.post('/api/auth/init-admin', data)
}

// 文章接口
export const articleApi = {
  getArticles: (params) => axios.get('/api/articles', { params }),
  getArticle: (id) => axios.get(`/api/articles/${id}`),
  createArticle: (data) => axios.post('/api/articles', data),
  updateArticle: (id, data) => axios.put(`/api/articles/${id}`, data),
  deleteArticle: (id) => axios.delete(`/api/articles/${id}`),
  likeArticle: (id) => axios.post(`/api/articles/${id}/like`),
  unlikeArticle: (id) => axios.delete(`/api/articles/${id}/like`),
  searchArticles: (params) => axios.get('/api/articles/search', { params }),
  shareArticle: (id) => axios.post(`/api/articles/${id}/share`),
  getMyArticles: (params) => axios.get('/api/my/articles', { params }),
  getMyDrafts: (params) => axios.get('/api/my/drafts', { params }),
  publishDraft: (id) => axios.post(`/api/articles/${id}/publish`),
  pinArticle: (id) => axios.post(`/api/articles/${id}/pin`),
  unpinArticle: (id) => axios.delete(`/api/articles/${id}/pin`),
  restoreArticle: (id) => axios.post(`/api/articles/${id}/restore`)
}

// 评论接口
export const commentApi = {
  createComment: (articleId, data) => axios.post(`/api/articles/${articleId}/comments`, data),
  deleteComment: (id) => axios.delete(`/api/comments/${id}`),
  likeComment: (id) => axios.post(`/api/comments/${id}/like`),
  unlikeComment: (id) => axios.delete(`/api/comments/${id}/like`)
}

// 用户接口
export const userApi = {
  getUser: (id) => axios.get(`/api/users/${id}`),
  follow: (id) => axios.post(`/api/follow/${id}`),
  unfollow: (id) => axios.delete(`/api/follow/${id}`),
  getFollowing: () => axios.get('/api/following'),
  getFollowers: () => axios.get('/api/followers'),
  checkFollowStatus: (id) => axios.get(`/api/follow/status/${id}`),
  getMutualFriends: () => axios.get('/api/mutual'),
  getUserFollowing: (id) => axios.get(`/api/users/${id}/following`),
  getUserFollowers: (id) => axios.get(`/api/users/${id}/followers`)
}

// 收藏接口
export const favoriteApi = {
  addFavorite: (id) => axios.post(`/api/articles/${id}/favorite`),
  removeFavorite: (id) => axios.delete(`/api/articles/${id}/favorite`),
  getFavorites: (params) => axios.get('/api/favorites', { params }),
  checkFavorite: (id) => axios.get(`/api/articles/${id}/favorite/check`)
}

// 签到接口
export const signinApi = {
  signin: () => axios.post('/api/signin'),
  getStatus: () => axios.get('/api/signin/status'),
  getHistory: (params) => axios.get('/api/signin/history', { params })
}

// 通知接口
export const notificationApi = {
  getNotifications: () => axios.get('/api/notifications'),
  getUnreadCount: () => axios.get('/api/notifications/unread-count'),
  markRead: (id) => axios.post(`/api/notifications/${id}/read`),
  markAllRead: () => axios.post('/api/notifications/read-all'),
  getCommentReplyNotifications: () => axios.get('/api/comment-reply-notifications'),
  getFollowNotifications: () => axios.get('/api/follow-notifications'),
  getUserNotifications: () => axios.get('/api/user-notifications')
}

// 分类接口
export const categoryApi = {
  getCategories: () => axios.get('/api/categories'),
  createCategory: (data) => axios.post('/api/categories', data),
  updateCategory: (id, data) => axios.put(`/api/categories/${id}`, data),
  deleteCategory: (id) => axios.delete(`/api/categories/${id}`)
}

// 举报接口
export const reportApi = {
  getReports: () => axios.get('/api/reports'),
  handleReport: (id, data) => axios.put(`/api/reports/${id}/handle`, data),
  createReport: (data) => axios.post('/api/reports', data)
}

// 管理员接口
export const adminApi = {
  checkAdmin: () => axios.get('/api/admin/check'),
  getStatistics: () => axios.get('/api/admin/statistics'),
  getUsers: () => axios.get('/api/admin/users'),
  updateUser: (id, data) => axios.put(`/api/admin/users/${id}`, data),
  updateUserRole: (id, data) => axios.put(`/api/admin/users/${id}/role`, data),
  banUser: (id, data) => axios.post(`/api/admin/users/${id}/ban`, data),
  unbanUser: (id) => axios.post(`/api/admin/users/${id}/unban`),
  deleteUser: (id) => axios.delete(`/api/admin/users/${id}`),
  getArticles: () => axios.get('/api/admin/articles'),
  updateArticleStatus: (id, data) => axios.put(`/api/admin/articles/${id}/status`, data),
  getComments: () => axios.get('/api/admin/comments'),
  deleteComment: (id) => axios.delete(`/api/admin/comments/${id}`),
  createCategory: (data) => axios.post('/api/categories', data),
  updateCategory: (id, data) => axios.put(`/api/categories/${id}`, data),
  deleteCategory: (id) => axios.delete(`/api/categories/${id}`),
  updateAnnouncement: (data) => axios.put('/api/announcement', data),
  updateSiteConfig: (data) => axios.put('/api/site-config', data),
  createNotification: (data) => axios.post('/api/notifications', data),
  sendUserNotification: (data) => axios.post('/api/user-notifications/send', data),
  sendBatchNotification: (data) => axios.post('/api/user-notifications/send-batch', data)
}

export default {
  auth: authApi,
  article: articleApi,
  comment: commentApi,
  user: userApi,
  favorite: favoriteApi,
  signin: signinApi,
  notification: notificationApi,
  category: categoryApi,
  report: reportApi,
  admin: adminApi
}
