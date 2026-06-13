import api from './index'

// 评论相关API
export const commentApi = {
  // 创建评论
  createComment: (articleId, data) => api.post(`/articles/${articleId}/comments`, {
    content: data.content,
    parent_id: data.parent_id || null,
    is_anonymous: data.is_anonymous || false
  }),

  // 删除评论
  deleteComment: (id) => api.delete(`/comments/${id}`),

  // 点赞评论
  likeComment: (id) => api.post(`/comments/${id}/like`),

  // 取消点赞
  unlikeComment: (id) => api.delete(`/comments/${id}/like`),
}
