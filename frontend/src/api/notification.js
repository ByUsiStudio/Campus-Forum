import api from './index'

// 通知相关API
export const notificationApi = {
  // 获取通知列表
  getNotifications: () => api.get('/notifications'),

  // 获取未读通知数量
  getUnreadCount: () => api.get('/notifications/unread-count'),

  // 标记通知为已读
  markNotificationRead: (id) => api.post(`/notifications/${id}/read`),

  // 标记所有通知为已读
  markAllNotificationsRead: () => api.post('/notifications/read-all'),

  // 获取评论回复通知
  getCommentReplyNotifications: () => api.get('/comment-reply-notifications'),

  // 标记评论回复通知为已读
  markCommentReplyNotificationRead: (id) => api.post(`/comment-reply-notifications/${id}/read`),

  // 标记所有评论回复通知为已读
  markAllCommentReplyNotificationsRead: () => api.post('/comment-reply-notifications/read-all'),

  // 获取粉丝通知
  getFollowNotifications: () => api.get('/follow-notifications'),

  // 标记粉丝通知为已读
  markFollowNotificationRead: (id) => api.post(`/follow-notifications/${id}/read`),

  // 标记所有粉丝通知为已读
  markAllFollowNotificationsRead: () => api.post('/follow-notifications/read-all'),

  // 获取粉丝通知未读数量
  getFollowNotificationUnreadCount: () => api.get('/follow-notifications/unread-count'),
}

// 用户个人通知API
export const userNotificationApi = {
  // 获取用户通知列表
  getUserNotifications: () => api.get('/user-notifications'),

  // 获取单个通知详情
  getNotification: (id) => api.get(`/user-notifications/${id}`),

  // 标记通知为已读
  markAsRead: (id) => api.post(`/user-notifications/${id}/read`),

  // 标记所有通知为已读
  markAllAsRead: () => api.post('/user-notifications/read-all'),

  // 删除通知
  deleteNotification: (id) => api.delete(`/user-notifications/${id}`),

  // 清空所有通知
  clearAll: () => api.delete('/user-notifications/clear'),
}

// 管理员用户通知API
export const adminUserNotificationApi = {
  // 发送单个通知
  sendNotification: (data) => api.post('/user-notifications/send', {
    user_id: data.user_id,
    title: data.title,
    content: data.content,
    type: data.type
  }),

  // 批量发送通知
  sendBatchNotifications: (data) => api.post('/user-notifications/send-batch', {
    user_ids: data.user_ids,
    title: data.title,
    content: data.content,
    type: data.type
  }),

  // 获取用户通知（管理员）
  getUserNotifications: (userId) => api.get(`/admin/user-notifications/${userId}`),
}
