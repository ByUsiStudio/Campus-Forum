import api from './http'

// 鐢ㄦ埛鐩稿叧API
export const userApi = {
  // 鑾峰彇鐢ㄦ埛鍏紑淇℃伅
  getUserByID: (id) => api.get(`/users/${id}`),

  // 鑾峰彇鐢ㄦ埛鏂囩珷鍒楄〃
  getUserArticles: (id, params = {}) => api.get(`/users/${id}/articles`, { params }),

  // 鑾峰彇鐢ㄦ埛澶磋
  getUserTitles: (id) => api.get(`/users/${id}/titles`),
}

// 绛惧埌鐩稿叧API
export const signinApi = {
  // 鐢ㄦ埛绛惧埌
  signIn: () => api.post('/signin'),

  // 鑾峰彇绛惧埌鐘舵€?  getSignInStatus: () => api.get('/signin/status'),

  // 鑾峰彇绛惧埌鍘嗗彶
  getSignInHistory: (params = {}) => api.get('/signin/history', { params }),

  // 鑾峰彇绛惧埌鎺掕姒?  getSignInRankings: (params = {}) => api.get('/signin/rankings', { params }),

  // 鑾峰彇绛惧埌閰嶇疆
  getSignInConfig: () => api.get('/signin/config'),

  // 鏇存柊绛惧埌閰嶇疆锛堢鐞嗗憳锛?  updateSignInConfig: (data) => api.put('/signin/config', data),
}
