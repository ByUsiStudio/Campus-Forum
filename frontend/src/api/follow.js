import request from './common'

// 好友系统 API

// 发送好友请求
export const sendFriendRequest = (friendId, message = '') => {
  return request({
    url: '/api/friends/request',
    method: 'post',
    data: {
      friend_id: friendId,
      message
    }
  })
}

// 同意好友请求
export const acceptFriendRequest = (requestId) => {
  return request({
    url: `/api/friends/request/${requestId}/accept`,
    method: 'post'
  })
}

// 拒绝好友请求
export const rejectFriendRequest = (requestId) => {
  return request({
    url: `/api/friends/request/${requestId}/reject`,
    method: 'post'
  })
}

// 删除好友
export const deleteFriend = (friendId) => {
  return request({
    url: `/api/friends/${friendId}`,
    method: 'delete'
  })
}

// 获取好友列表
export const getFriendList = () => {
  return request({
    url: '/api/friends',
    method: 'get'
  })
}

// 获取好友请求列表
export const getFriendRequests = () => {
  return request({
    url: '/api/friends/requests',
    method: 'get'
  })
}

// 获取已发送的好友请求列表
export const getSentFriendRequests = () => {
  return request({
    url: '/api/friends/sent-requests',
    method: 'get'
  })
}

// 更新好友备注名
export const updateFriendDisplayName = (friendId, displayName) => {
  return request({
    url: `/api/friends/${friendId}/display-name`,
    method: 'put',
    data: {
      display_name: displayName
    }
  })
}

// 检查好友状态
export const checkFriendStatus = (userId) => {
  return request({
    url: `/api/friends/status/${userId}`,
    method: 'get'
  })
}

// 获取共同好友
export const getMutualFriends = (userId) => {
  return request({
    url: `/api/friends/mutual/${userId}`,
    method: 'get'
  })
}