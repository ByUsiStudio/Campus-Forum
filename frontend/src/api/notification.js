import api from './http'

// 閫氱煡鐩稿叧API
export const notificationApi = {
  // 鑾峰彇閫氱煡鍒楄〃
  getNotifications: () => api.get('/notifications'),

  // 鑾峰彇鏈閫氱煡鏁伴噺
  getUnreadCount: () => api.get('/notifications/unread-count'),

  // 鏍囪閫氱煡涓哄凡璇?  markNotificationRead: (id) => api.post(`/notifications/${id}/read`),

  // 鏍囪鎵€鏈夐€氱煡涓哄凡璇?  markAllNotificationsRead: () => api.post('/notifications/read-all'),

  // 鑾峰彇璇勮鍥炲閫氱煡
  getCommentReplyNotifications: () => api.get('/comment-reply-notifications'),

  // 鏍囪璇勮鍥炲閫氱煡涓哄凡璇?  markCommentReplyNotificationRead: (id) => api.post(`/comment-reply-notifications/${id}/read`),

  // 鏍囪鎵€鏈夎瘎璁哄洖澶嶉€氱煡涓哄凡璇?  markAllCommentReplyNotificationsRead: () => api.post('/comment-reply-notifications/read-all'),
}

// 鐢ㄦ埛涓汉閫氱煡API
export const userNotificationApi = {
  // 鑾峰彇鐢ㄦ埛閫氱煡鍒楄〃
  getUserNotifications: () => api.get('/user-notifications'),

  // 鑾峰彇鍗曚釜閫氱煡璇︽儏
  getNotification: (id) => api.get(`/user-notifications/${id}`),

  // 鏍囪閫氱煡涓哄凡璇?  markAsRead: (id) => api.post(`/user-notifications/${id}/read`),

  // 鏍囪鎵€鏈夐€氱煡涓哄凡璇?  markAllAsRead: () => api.post('/user-notifications/read-all'),

  // 鍒犻櫎閫氱煡
  deleteNotification: (id) => api.delete(`/user-notifications/${id}`),

  // 娓呯┖鎵€鏈夐€氱煡
  clearAll: () => api.delete('/user-notifications/clear'),
}

// 绠＄悊鍛樼敤鎴烽€氱煡API
export const adminUserNotificationApi = {
  // 鍙戦€佸崟涓€氱煡
  sendNotification: (data) => api.post('/user-notifications/send', {
    user_id: data.user_id,
    title: data.title,
    content: data.content,
    type: data.type
  }),

  // 鎵归噺鍙戦€侀€氱煡
  sendBatchNotifications: (data) => api.post('/user-notifications/send-batch', {
    user_ids: data.user_ids,
    title: data.title,
    content: data.content,
    type: data.type
  }),

  // 鑾峰彇鐢ㄦ埛閫氱煡锛堢鐞嗗憳锛?  getUserNotifications: (userId) => api.get(`/admin/user-notifications/${userId}`),
}
