import api from './index'

// 评论相关API
export const commentApi = {
  // 创建评论
  createComment: (articleId, data) => api.post(`/articles/${articleId}/comments`, data),

  // 删除评论
  deleteComment: (id) => api.delete(`/comments/${id}`),

  // 点赞评论
  likeComment: (id) => api.post(`/comments/${id}/like`),

  // 取消点赞
  unlikeComment: (id) => api.delete(`/comments/${id}/like`),
}
