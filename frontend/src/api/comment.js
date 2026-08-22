import api from './http'

// 璇勮鐩稿叧API
export const commentApi = {
  // 鍒涘缓璇勮
  createComment: (articleId, data) => api.post(`/articles/${articleId}/comments`, {
    content: data.content,
    parent_id: data.parent_id || null,
    is_anonymous: data.is_anonymous || false
  }),

  // 鍒犻櫎璇勮
  deleteComment: (id) => api.delete(`/comments/${id}`),

  // 鐐硅禐璇勮
  likeComment: (id) => api.post(`/comments/${id}/like`),

  // 鍙栨秷鐐硅禐
  unlikeComment: (id) => api.delete(`/comments/${id}/like`),
}
