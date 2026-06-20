import api from './index'

// 排行榜与徽章API
export const leaderboardApi = {
  // 获取排行榜
  getLeaderboard(type = 'experience', period = 'all_time', limit = 50) {
    return api.get('/leaderboard', { params: { type, period, limit } })
  },

  // 获取用户排名
  getUserRank(type = 'experience', period = 'all_time') {
    return api.get('/leaderboard/rank', { params: { type, period } })
  },

  // 获取用户徽章
  getUserBadges() {
    return api.get('/badges')
  },

  // 更新徽章显示状态
  updateBadgeDisplay(badgeId, isDisplayed) {
    return api.put(`/badges/${badgeId}/display`, { is_displayed: isDisplayed })
  },

  // 授予用户徽章（管理员）
  grantBadge(data) {
    return api.post('/badges/grant', data)
  },

  // 撤销用户徽章（管理员）
  revokeBadge(badgeId) {
    return api.delete(`/badges/${badgeId}`)
  }
}