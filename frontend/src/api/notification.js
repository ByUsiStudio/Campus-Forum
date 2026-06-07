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
