import axios from 'axios'

// 好友系统 API

// 发送好友请求
export const sendFriendRequest = (friendId, message = '') => {
  return axios.post('/api/friends/request', {
    friend_id: friendId,
    message
  })
}

// 同意好友请求
export const acceptFriendRequest = (requestId) => {
  return axios.post(`/api/friends/request/${requestId}/accept`)
}

// 拒绝好友请求
export const rejectFriendRequest = (requestId) => {
  return axios.post(`/api/friends/request/${requestId}/reject`)
}

// 删除好友
export const deleteFriend = (friendId) => {
  return axios.delete(`/api/friends/${friendId}`)
}

// 获取好友列表
export const getFriendList = () => {
  return axios.get('/api/friends')
}

// 获取好友请求列表
export const getFriendRequests = () => {
  return axios.get('/api/friends/requests')
}

// 获取已发送的好友请求列表
export const getSentFriendRequests = () => {
  return axios.get('/api/friends/sent-requests')
}

// 更新好友备注名
export const updateFriendDisplayName = (friendId, displayName) => {
  return axios.put(`/api/friends/${friendId}/display-name`, {
    display_name: displayName
  })
}

// 检查好友状态
export const checkFriendStatus = (userId) => {
  return axios.get(`/api/friends/status/${userId}`)
}

// 获取共同好友
export const getMutualFriends = (userId) => {
  return axios.get(`/api/friends/mutual/${userId}`)
}
