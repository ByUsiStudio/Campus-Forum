import api from './index'

// 聊天相关API
export const chatApi = {
  // 获取会话列表
  getConversations: () => api.get('/chat/conversations'),

  // 获取消息历史
  getMessages: (params) => api.get('/chat/messages', { params }),

  // 发送消息
  sendMessage: (data) => api.post('/chat/messages', {
    conversation_id: data.conversation_id,
    content: data.content,
    type: data.type || 'text'
  }),

  // 创建私聊会话
  createPrivateConversation: (userId) => api.post('/chat/conversations/private', { user_id: userId }),

  // 获取未读消息数量
  getChatUnreadCount: () => api.get('/chat/unread'),

  // 标记会话已读
  markConversationRead: (conversationId) => api.post(`/chat/conversations/${conversationId}/read`)
}
