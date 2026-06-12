import request from './common'

// 聊天相关 API（通过论坛端调用IM服务）

// 获取会话列表
export const getConversations = () => {
  return request({
    url: '/api/chat/conversations',
    method: 'get'
  })
}

// 获取聊天消息
export const getMessages = (conversationId, limit = 20, offset = 0) => {
  return request({
    url: '/api/chat/messages',
    method: 'get',
    params: {
      conversation_id: conversationId,
      limit,
      offset
    }
  })
}

// 发送消息
export const sendMessage = (conversationId, content, type = 'text') => {
  return request({
    url: '/api/chat/messages',
    method: 'post',
    data: {
      conversation_id: conversationId,
      content,
      type
    }
  })
}

// 创建私聊会话
export const createPrivateConversation = (targetUserId) => {
  return request({
    url: '/api/chat/conversations/private',
    method: 'post',
    data: {
      target_user_id: targetUserId
    }
  })
}

// 获取未读消息数
export const getUnreadCount = () => {
  return request({
    url: '/api/chat/unread',
    method: 'get'
  })
}

// 标记会话已读
export const markConversationRead = (conversationId) => {
  return request({
    url: `/api/chat/conversations/${conversationId}/read`,
    method: 'post'
  })
}