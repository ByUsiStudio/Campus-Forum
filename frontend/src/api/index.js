import axios from 'axios'

// 认证接口
export const authApi = {
  login: (data) => axios.post('/auth/login', data),
  register: (data) => axios.post('/auth/register', data),
  getProfile: () => axios.get('/profile'),
  updateProfile: (data) => axios.put('/profile', data),
  uploadAvatar: (formData) => axios.post('/upload/avatar', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

// 文章接口
export const articleApi = {
  getArticles: (params) => axios.get('/articles', { params }),
  getArticle: (id) => axios.get(`/articles/${id}`),
  createArticle: (data) => axios.post('/articles', data),
  updateArticle: (id, data) => axios.put(`/articles/${id}`, data),
  deleteArticle: (id) => axios.delete(`/articles/${id}`),
  likeArticle: (id) => axios.post(`/articles/${id}/like`),
  unlikeArticle: (id) => axios.delete(`/articles/${id}/like`),
  searchArticles: (params) => axios.get('/articles/search', { params }),
  shareArticle: (id) => axios.post(`/articles/${id}/share`)
}

// 评论接口
export const commentApi = {
  createComment: (articleId, data) => axios.post(`/articles/${articleId}/comments`, data),
  deleteComment: (id) => axios.delete(`/comments/${id}`),
  likeComment: (id) => axios.post(`/comments/${id}/like`),
  unlikeComment: (id) => axios.delete(`/comments/${id}/like`)
}

// 用户接口
export const userApi = {
  getUser: (id) => axios.get(`/users/${id}`),
  getUserArticles: (id, params) => axios.get(`/users/${id}/articles`, { params }),
  follow: (id) => axios.post(`/follow/${id}`),
  unfollow: (id) => axios.delete(`/follow/${id}`),
  getFollowing: () => axios.get('/following'),
  getFollowers: () => axios.get('/followers')
}

// 收藏接口
export const favoriteApi = {
  addFavorite: (id) => axios.post(`/articles/${id}/favorite`),
  removeFavorite: (id) => axios.delete(`/articles/${id}/favorite`),
  getFavorites: (params) => axios.get('/favorites', { params }),
  checkFavorite: (id) => axios.get(`/articles/${id}/favorite/check`)
}

// 签到接口
export const signinApi = {
  signin: () => axios.post('/signin'),
  getStatus: () => axios.get('/signin/status'),
  getHistory: (params) => axios.get('/signin/history', { params })
}

// 通知接口
export const notificationApi = {
  getNotifications: () => axios.get('/notifications'),
  getUnreadCount: () => axios.get('/notifications/unread-count'),
  markRead: (id) => axios.post(`/notifications/${id}/read`),
  markAllRead: () => axios.post('/notifications/read-all')
}

// 分类接口
export const categoryApi = {
  getCategories: () => axios.get('/categories')
}

// 管理员接口
export const adminApi = {
  getUsers: () => axios.get('/admin/users'),
  getArticles: () => axios.get('/admin/articles'),
  getComments: () => axios.get('/admin/comments'),
  getReports: () => axios.get('/reports'),
  handleReport: (id, data) => axios.put(`/reports/${id}/handle`, data),
  updateUserRole: (id, data) => axios.put(`/admin/users/${id}/role`, data),
  banUser: (id, data) => axios.post(`/admin/users/${id}/ban`, data),
  deleteUser: (id) => axios.delete(`/admin/users/${id}`),
  deleteArticle: (id) => axios.delete(`/articles/${id}`),
  deleteComment: (id) => axios.delete(`/admin/comments/${id}`),
  createCategory: (data) => axios.post('/categories', data),
  updateCategory: (id, data) => axios.put(`/categories/${id}`, data),
  deleteCategory: (id) => axios.delete(`/categories/${id}`),
  updateAnnouncement: (data) => axios.put('/announcement', data),
  updateSiteConfig: (data) => axios.put('/site-config', data)
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
  admin: adminApi
}