import api from './index'

// 文章相关API
export const articleApi = {
  // 获取文章列表
  getArticles: (params = {}) => api.get('/articles', { params }),

  // 获取文章详情
  getArticle: (id, params = {}) => api.get(`/articles/${id}`, { params }),

  // 创建文章
  createArticle: (data) => api.post('/articles', {
    title: data.title,
    content: data.content,
    category_id: data.category_id,
    voice_url: data.voice_url,
    is_anonymous: data.is_anonymous || false,
    status: data.status || 'published'
  }),

  // 更新文章
  updateArticle: (id, data) => api.put(`/articles/${id}`, {
    title: data.title,
    content: data.content,
    category_id: data.category_id,
    voice_url: data.voice_url,
    is_anonymous: data.is_anonymous
  }),

  // 删除文章（软删除）
  deleteArticle: (id, data) => api.delete(`/articles/${id}`, data),

  // 恢复文章（管理员）
  restoreArticle: (id) => api.post(`/articles/${id}/restore`),

  // 点赞文章
  likeArticle: (id) => api.post(`/articles/${id}/like`),

  // 取消点赞
  unlikeArticle: (id) => api.delete(`/articles/${id}/like`),

  // 收藏文章
  addFavorite: (id) => api.post(`/articles/${id}/favorite`),

  // 取消收藏
  removeFavorite: (id) => api.delete(`/articles/${id}/favorite`),

  // 检查是否已收藏
  checkFavorite: (id) => api.get(`/articles/${id}/favorite/check`),

  // 获取当前用户文章
  getMyArticles: (params = {}) => api.get('/my/articles', { params }),

  // 搜索文章
  searchArticles: (params = {}) => api.get('/articles/search', { params }),

  // 分享文章
  shareArticle: (id) => api.post(`/articles/${id}/share`),

  // 获取收藏列表
  getFavorites: () => api.get('/favorites'),

  // 获取草稿列表
  getDrafts: (params = {}) => api.get('/my/drafts', { params }),

  // 发布草稿
  publishDraft: (id) => api.post(`/articles/${id}/publish`),

  // 置顶文章（管理员）
  pinArticle: (id) => api.post(`/articles/${id}/pin`),

  // 取消置顶（管理员）
  unpinArticle: (id) => api.delete(`/articles/${id}/pin`),
}

// 举报相关API
export const reportApi = {
  // 创建举报
  createReport: (data) => api.post('/reports', data),

  // 获取举报列表（管理员）
  getReports: (params = {}) => api.get('/reports', { params }),

  // 获取举报详情（管理员）
  getReport: (id) => api.get(`/reports/${id}`),

  // 处理举报（管理员）
  handleReport: (id, data) => api.put(`/reports/${id}/handle`, data),
}
