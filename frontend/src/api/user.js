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
