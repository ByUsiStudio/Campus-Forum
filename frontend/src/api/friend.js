import api from './index'

// 好友相关API
export const friendApi = {
  // 发送好友请求
  sendFriendRequest: (data) => api.post('/friends/request', {
    user_id: data.user_id,
    message: data.message
  }),

  // 同意好友请求
  acceptFriendRequest: (requestId) => api.post(`/friends/request/${requestId}/accept`),

  // 拒绝好友请求
  rejectFriendRequest: (requestId) => api.post(`/friends/request/${requestId}/reject`),

  // 删除好友
  deleteFriend: (friendId) => api.delete(`/friends/${friendId}`),

  // 获取好友列表
  getFriendList: () => api.get('/friends'),

  // 获取收到的好友请求
  getFriendRequests: () => api.get('/friends/requests'),

  // 获取发送的好友请求
  getSentFriendRequests: () => api.get('/friends/sent-requests'),

  // 更新好友备注名
  updateFriendDisplayName: (friendId, displayName) => api.put(`/friends/${friendId}/display-name`, {
    display_name: displayName
  }),

  // 检查好友状态
  checkFriendStatus: (userId) => api.get(`/friends/status/${userId}`),

  // 获取共同好友
  getMutualFriends: (userId) => api.get(`/friends/mutual/${userId}`)
}
