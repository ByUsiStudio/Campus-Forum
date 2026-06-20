import api from './index'

// 用户等级与成就系统API
export const levelApi = {
  // 获取用户等级信息
  getLevel() {
    return api.get('/level')
  },

  // 获取用户经验记录
  getExperienceRecords(page = 1, limit = 20) {
    return api.get('/level/experience-records', { params: { page, limit } })
  },

  // 获取用户成就列表
  getUserAchievements() {
    return api.get('/achievements')
  },

  // 获取所有成就定义
  getAllAchievements() {
    return api.get('/achievements/all')
  },

  // 获取等级配置
  getLevelConfig() {
    return api.get('/level/config')
  },

  // 创建等级配置（管理员）
  createLevelConfig(data) {
    return api.post('/level/config', data)
  },

  // 更新等级配置（管理员）
  updateLevelConfig(id, data) {
    return api.put(`/level/config/${id}`, data)
  },

  // 创建成就（管理员）
  createAchievement(data) {
    return api.post('/achievements', data)
  },

  // 更新成就（管理员）
  updateAchievement(id, data) {
    return api.put(`/achievements/${id}`, data)
  },

  // 删除成就（管理员）
  deleteAchievement(id) {
    return api.delete(`/achievements/${id}`)
  }
}