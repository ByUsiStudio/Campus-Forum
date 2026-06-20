import api from './index'

// 收藏夹管理API
export const collectionApi = {
  // 获取用户收藏夹列表
  getCollections() {
    return api.get('/collections')
  },

  // 获取单个收藏夹详情
  getCollection(collectionId) {
    return api.get(`/collections/${collectionId}`)
  },

  // 创建收藏夹
  createCollection(data) {
    return api.post('/collections', data)
  },

  // 更新收藏夹
  updateCollection(collectionId, data) {
    return api.put(`/collections/${collectionId}`, data)
  },

  // 删除收藏夹
  deleteCollection(collectionId) {
    return api.delete(`/collections/${collectionId}`)
  },

  // 将文章添加到收藏夹
  addArticleToCollection(collectionId, articleId, note = '') {
    return api.post(`/collections/${collectionId}/articles`, { article_id: articleId, note })
  },

  // 从收藏夹移除文章
  removeArticleFromCollection(collectionId, articleId) {
    return api.delete(`/collections/${collectionId}/articles/${articleId}`)
  },

  // 获取文章版本历史
  getArticleVersions(articleId) {
    return api.get(`/articles/${articleId}/versions`)
  },

  // 获取文章特定版本
  getArticleVersion(articleId, version) {
    return api.get(`/articles/${articleId}/versions/${version}`)
  },

  // 恢复文章到特定版本
  restoreArticleVersion(articleId, version) {
    return api.post(`/articles/${articleId}/versions/${version}/restore`)
  }
}