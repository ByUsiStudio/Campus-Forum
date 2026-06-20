import api from './index'

// 话题管理API
export const topicApi = {
  // 获取话题列表
  getTopics(page = 1, limit = 20, isHot = false) {
    return api.get('/topics', { params: { page, limit, is_hot: isHot } })
  },

  // 获取话题详情
  getTopic(topicId, page = 1, limit = 20) {
    return api.get(`/topics/${topicId}`, { params: { page, limit } })
  },

  // 创建话题（管理员）
  createTopic(data) {
    return api.post('/topics', data)
  },

  // 更新话题（管理员）
  updateTopic(topicId, data) {
    return api.put(`/topics/${topicId}`, data)
  },

  // 删除话题（管理员）
  deleteTopic(topicId) {
    return api.delete(`/topics/${topicId}`)
  },

  // 关注话题
  followTopic(topicId) {
    return api.post(`/topics/${topicId}/follow`)
  },

  // 取消关注话题
  unfollowTopic(topicId) {
    return api.delete(`/topics/${topicId}/follow`)
  },

  // 获取用户关注的话题
  getFollowedTopics() {
    return api.get('/topics/followed')
  },

  // 为文章添加话题
  addTopicToArticle(articleId, topicId) {
    return api.post(`/articles/${articleId}/topics`, { topic_id: topicId })
  },

  // 从文章移除话题
  removeTopicFromArticle(articleId, topicId) {
    return api.delete(`/articles/${articleId}/topics/${topicId}`)
  },

  // 获取热门话题
  getHotTopics(period = 'daily', limit = 10) {
    return api.get('/topics/hot', { params: { period, limit } })
  }
}