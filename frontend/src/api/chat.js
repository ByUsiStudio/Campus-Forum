import axios from 'axios'

// 聊天相关 API（通过论坛端调用IM服务）

// 获取会话列表
export const getConversations = () => {
  return axios.get('/api/chat/conversations')
}

// 获取聊天消息
export const getMessages = (conversationId, limit = 20, offset = 0) => {
  return axios.get('/api/chat/messages', {
    params: {
      conversation_id: conversationId,
      limit,
      offset
    }
  })
}

// 发送消息
export const sendMessage = (conversationId, content, type = 'text') => {
  return axios.post('/api/chat/messages', {
    conversation_id: conversationId,
    content,
    type
  })
}

// 创建私聊会话
export const createPrivateConversation = (targetUserId) => {
  return axios.post('/api/chat/conversations/private', {
    target_user_id: targetUserId
  })
}

// 获取未读消息数
export const getUnreadCount = () => {
  return axios.get('/api/chat/unread')
}

// 标记会话已读
export const markConversationRead = (conversationId) => {
  return axios.post(`/api/chat/conversations/${conversationId}/read`)
}
