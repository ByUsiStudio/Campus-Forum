import api from './index'

// 数据统计与分析API
export const statisticsApi = {
  // 获取用户统计数据
  getUserStatistics() {
    return api.get('/statistics')
  },

  // 获取每日统计数据
  getDailyStatistics(startDate, endDate) {
    return api.get('/statistics/daily', { params: { start_date: startDate, end_date: endDate } })
  },

  // 获取系统概览数据
  getSystemOverview() {
    return api.get('/statistics/overview')
  },

  // 获取用户活跃度数据
  getUserActivity(startDate, endDate) {
    return api.get('/statistics/activity', { params: { start_date: startDate, end_date: endDate } })
  },

  // 获取统计仪表板数据（管理员）
  getStatisticsDashboard() {
    return api.get('/statistics/dashboard')
  },

  // 获取文章统计数据
  getArticleStatistics(articleId) {
    return api.get(`/articles/${articleId}/statistics`)
  }
}