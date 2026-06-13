import api from './index'

// 用户相关API
export const userApi = {
  // 获取用户公开信息
  getUserByID: (id) => api.get(`/users/${id}`),

  // 获取用户文章列表
  getUserArticles: (id, params = {}) => api.get(`/users/${id}/articles`, { params }),

  // 获取用户关注列表
  getUserFollowing: (id) => api.get(`/users/${id}/following`),

  // 获取用户粉丝列表
  getUserFollowers: (id) => api.get(`/users/${id}/followers`),

  // 获取用户头衔
  getUserTitles: (id) => api.get(`/users/${id}/titles`),
}

// 签到相关API
export const signinApi = {
  // 用户签到
  signIn: () => api.post('/signin'),

  // 获取签到状态
  getSignInStatus: () => api.get('/signin/status'),

  // 获取签到历史
  getSignInHistory: (params = {}) => api.get('/signin/history', { params }),

  // 获取签到排行榜
  getSignInRankings: (params = {}) => api.get('/signin/rankings', { params }),

  // 获取签到配置
  getSignInConfig: () => api.get('/signin/config'),

  // 更新签到配置（管理员）
  updateSignInConfig: (data) => api.put('/signin/config', data),
}
