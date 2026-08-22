import api from './http'

// 鏂囩珷鐩稿叧API
export const articleApi = {
  // 鑾峰彇鏂囩珷鍒楄〃
  getArticles: (params = {}) => api.get('/articles', { params }),

  // 鑾峰彇鏂囩珷璇︽儏
  getArticle: (id, params = {}) => api.get(`/articles/${id}`, { params }),

  // 鍒涘缓鏂囩珷
  createArticle: (data) => api.post('/articles', {
    title: data.title,
    content: data.content,
    category_id: data.category_id,
    voice_url: data.voice_url,
    is_anonymous: data.is_anonymous || false,
    status: data.status || 'published'
  }),

  // 鏇存柊鏂囩珷
  updateArticle: (id, data) => api.put(`/articles/${id}`, {
    title: data.title,
    content: data.content,
    category_id: data.category_id,
    voice_url: data.voice_url,
    is_anonymous: data.is_anonymous
  }),

  // 鍒犻櫎鏂囩珷锛堣蒋鍒犻櫎锛?  deleteArticle: (id, data) => api.delete(`/articles/${id}`, data),

  // 鎭㈠鏂囩珷锛堢鐞嗗憳锛?  restoreArticle: (id) => api.post(`/articles/${id}/restore`),

  // 鐐硅禐鏂囩珷
  likeArticle: (id) => api.post(`/articles/${id}/like`),

  // 鍙栨秷鐐硅禐
  unlikeArticle: (id) => api.delete(`/articles/${id}/like`),

  // 鏀惰棌鏂囩珷
  addFavorite: (id) => api.post(`/articles/${id}/favorite`),

  // 鍙栨秷鏀惰棌
  removeFavorite: (id) => api.delete(`/articles/${id}/favorite`),

  // 妫€鏌ユ槸鍚﹀凡鏀惰棌
  checkFavorite: (id) => api.get(`/articles/${id}/favorite/check`),

  // 鑾峰彇褰撳墠鐢ㄦ埛鏂囩珷
  getMyArticles: (params = {}) => api.get('/my/articles', { params }),

  // 鎼滅储鏂囩珷
  searchArticles: (params = {}) => api.get('/articles/search', { params }),

  // 鍒嗕韩鏂囩珷
  shareArticle: (id) => api.post(`/articles/${id}/share`),

  // 鑾峰彇鏀惰棌鍒楄〃
  getFavorites: () => api.get('/favorites'),

  // 鑾峰彇鑽夌鍒楄〃
  getDrafts: (params = {}) => api.get('/my/drafts', { params }),

  // 鍙戝竷鑽夌
  publishDraft: (id) => api.post(`/articles/${id}/publish`),

  // 缃《鏂囩珷锛堢鐞嗗憳锛?  pinArticle: (id) => api.post(`/articles/${id}/pin`),

  // 鍙栨秷缃《锛堢鐞嗗憳锛?  unpinArticle: (id) => api.delete(`/articles/${id}/pin`),
}

// 涓炬姤鐩稿叧API
export const reportApi = {
  // 鍒涘缓涓炬姤
  createReport: (data) => api.post('/reports', data),

  // 鑾峰彇涓炬姤鍒楄〃锛堢鐞嗗憳锛?  getReports: (params = {}) => api.get('/reports', { params }),

  // 鑾峰彇涓炬姤璇︽儏锛堢鐞嗗憳锛?  getReport: (id) => api.get(`/reports/${id}`),

  // 澶勭悊涓炬姤锛堢鐞嗗憳锛?  handleReport: (id, data) => api.put(`/reports/${id}/handle`, data),
}
