import api from './index'

// 关注相关API
export const followApi = {
  // 关注用户
  followUser: (id) => api.post(`/follow/${id}`),

  // 取消关注
  unfollowUser: (id) => api.delete(`/follow/${id}`),

  // 获取当前用户关注列表
  getFollowingList: () => api.get('/following'),

  // 获取当前用户粉丝列表
  getFollowerList: () => api.get('/followers'),

  // 检查关注状态
  checkFollowStatus: (id) => api.get(`/follow/status/${id}`),

  // 获取互相关注好友
  getMutualFriends: () => api.get('/mutual'),
}
