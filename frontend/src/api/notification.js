import api from './http'

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
