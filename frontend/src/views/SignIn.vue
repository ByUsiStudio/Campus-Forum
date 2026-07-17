<template>
  <div class="sign-in-page">
    <div class="sign-in-card">
      <div class="card-header">
        <button class="back-btn" @click="$router.back()">
          <i class="fa-solid fa-arrow-left"></i>
        </button>
        <h2 class="page-title">每日签到</h2>
      </div>

      <div v-if="loading" class="loading-state">
        <div class="loading-spinner"></div>
      </div>

      <template v-else>
        <div class="sign-in-content">
          <div v-if="status.has_signed_in" class="signed-in-section">
            <div class="avatar success">
              <i class="fa-solid fa-check-circle"></i>
            </div>
            <h3 class="text-h5 mb-2">今日已签到</h3>
            <p class="text-body-2 text-grey">明天再来领取更多积分吧~</p>
          </div>

          <div v-else class="unsigned-section">
            <div class="avatar primary">
              <i class="fa-solid fa-calendar-check"></i>
            </div>
            <h3 class="text-h5 mb-2">今日可签到</h3>
            <p class="text-body-2 text-grey mb-4">
              连续签到 <span class="text-primary font-weight-bold">{{ status.sign_in_days }}</span> 天，获得 <span class="text-primary font-weight-bold">{{ status.total_coins || status.total_points || 0 }}</span> 币
            </p>
            <button
              class="layui-btn layui-btn-normal layui-btn-lg"
              :disabled="signing"
              @click="handleSignIn"
            >
              <i class="fa-solid fa-pencil mr-2"></i>
              {{ signing ? '签到中...' : '立即签到' }}
            </button>
          </div>

          <div class="stats-row">
            <div class="stat-item">
              <div class="stat-value text-primary">{{ status.total_coins || status.total_points || 0 }}</div>
              <div class="stat-label">累计币</div>
            </div>
            <div class="stat-item">
              <div class="stat-value text-primary">{{ status.sign_in_days || 0 }}</div>
              <div class="stat-label">连续天数</div>
            </div>
            <div class="stat-item">
              <div class="stat-value text-primary">{{ status.total_sign_ins || 0 }}</div>
              <div class="stat-label">累计次数</div>
            </div>
          </div>
        </div>

        <div class="stats-cards">
          <div class="stat-card">
            <i class="fa-solid fa-calendar-week text-orange"></i>
            <div class="stat-num">{{ status.week_sign_in_count || 0 }}</div>
            <div class="stat-desc">本周签到</div>
          </div>
          <div class="stat-card">
            <i class="fa-solid fa-calendar-days text-blue"></i>
            <div class="stat-num">{{ status.month_sign_in_count || 0 }}</div>
            <div class="stat-desc">本月签到</div>
          </div>
        </div>

        <div class="reward-card">
          <div class="card-title">
            <i class="fa-solid fa-gift mr-2"></i>
            连续签到奖励
          </div>
          <div class="reward-list">
            <div class="reward-item">
              <i class="fa-solid fa-circle text-amber"></i>
              <span>连续 <strong class="text-primary">7</strong> 天</span>
              <span class="text-success">+5 币</span>
            </div>
            <div class="reward-item">
              <i class="fa-solid fa-circle text-amber"></i>
              <span>连续 <strong class="text-primary">30</strong> 天</span>
              <span class="text-success">+15 币</span>
            </div>
            <div class="reward-item">
              <i class="fa-solid fa-circle text-amber"></i>
              <span>连续 <strong class="text-primary">365</strong> 天</span>
              <span class="text-success">+50 币</span>
            </div>
          </div>
        </div>

        <div class="rank-card">
          <div class="card-title">
            <i class="fa-solid fa-trophy mr-2"></i>
            签到排行榜
            <div class="rank-tabs">
              <button 
                class="rank-tab" 
                :class="{ active: rankType === 'continuous' }"
                @click="rankType = 'continuous'"
              >
                连续
              </button>
              <button 
                class="rank-tab" 
                :class="{ active: rankType === 'points' }"
                @click="rankType = 'points'"
              >
                币
              </button>
            </div>
          </div>
          <div class="rank-list">
            <div 
              v-for="(user, index) in (rankType === 'continuous' ? rankings.continuous_rankings : rankings.points_rankings)" 
              :key="user.id"
              class="rank-item"
            >
              <div class="rank-avatar" :class="getRankClass(index)">
                {{ index + 1 }}
              </div>
              <div class="rank-info">
                <div class="rank-name">{{ user.display_name || user.username }}</div>
              </div>
              <div class="rank-score text-primary">
                {{ rankType === 'continuous' ? user.sign_in_days + ' 天' : (user.total_coins || user.total_points) + ' 币' }}
              </div>
            </div>
            <div v-if="(!rankings.continuous_rankings?.length && rankType === 'continuous') || (!rankings.points_rankings?.length && rankType === 'points')" class="empty-rank">
              暂无数据
            </div>
          </div>
        </div>

        <div class="history-card">
          <div class="card-title">
            <i class="fa-solid fa-history mr-2"></i>
            签到记录
          </div>
          <div class="history-list">
            <div 
              v-for="record in history.records" 
              :key="record.id"
              class="history-item"
            >
              <div class="history-icon">
                <i class="fa-solid fa-check text-white"></i>
              </div>
              <div class="history-info">
                <div class="history-date">{{ record.sign_in_date }}</div>
                <div class="history-desc">连续 {{ record.continuous_day }} 天</div>
              </div>
              <div class="history-reward text-success">+{{ record.reward_points }} 币</div>
            </div>
            <div v-if="!history.records?.length" class="empty-history">
              暂无签到记录
            </div>
          </div>

          <div v-if="history.total_pages > 1" class="history-pagination">
            <button 
              class="page-btn" 
              :disabled="page <= 1"
              @click="page--; fetchHistory()"
            >
              <i class="fa-solid fa-chevron-left"></i>
            </button>
            <span class="page-info">{{ page }} / {{ history.total_pages }}</span>
            <button 
              class="page-btn" 
              :disabled="page >= history.total_pages"
              @click="page++; fetchHistory()"
            >
              <i class="fa-solid fa-chevron-right"></i>
            </button>
          </div>
        </div>
      </template>

      <div v-if="snackbar.show" class="snackbar" :class="snackbar.color">
        {{ snackbar.text }}
        <button class="snackbar-close" @click="snackbar.show = false">
          <i class="fa-solid fa-xmark"></i>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { signinApi } from '@/api'

const loading = ref(true)
const signing = ref(false)
const rankType = ref('continuous')

const status = ref({
  has_signed_in: false,
  sign_in_days: 0,
  total_sign_ins: 0,
  max_continuous_days: 0,
  total_points: 0,
  month_sign_in_count: 0,
  week_sign_in_count: 0
})

const rankings = ref({
  continuous_rankings: [],
  points_rankings: []
})

const history = ref({
  records: [],
  total: 0,
  page: 1,
  page_size: 30,
  total_pages: 0
})

const page = ref(1)

const snackbar = reactive({
  show: false,
  text: '',
  color: 'success'
})

const getRankClass = (index) => {
  if (index === 0) return 'gold'
  if (index === 1) return 'silver'
  if (index === 2) return 'bronze'
  return 'grey'
}

const fetchStatus = async () => {
  try {
    const res = await signinApi.getSignInStatus()
    if (res.data) {
      status.value = res.data
    }
  } catch (error) {
    console.error('获取签到状态失败:', error)
  }
}

const fetchRankings = async () => {
  try {
    const res = await signinApi.getSignInRankings({ limit: 10 })
    if (res.data) {
      rankings.value = res.data
    }
  } catch (error) {
    console.error('获取排行榜失败:', error)
  }
}

const fetchHistory = async () => {
  try {
    const res = await signinApi.getSignInHistory({
      page: page.value,
      page_size: 30
    })
    if (res.data) {
      history.value = res.data
    }
  } catch (error) {
    console.error('获取签到历史失败:', error)
  }
}

const handleSignIn = async () => {
  signing.value = true
  try {
    const res = await signinApi.signIn()
    snackbar.text = res.data?.message || '签到成功'
    snackbar.color = 'success'
    snackbar.show = true

    await fetchStatus()
    await fetchRankings()
  } catch (error) {
    snackbar.text = error.response?.data?.error || '签到失败'
    snackbar.color = 'error'
    snackbar.show = true
  } finally {
    signing.value = false
  }
}

onMounted(async () => {
  await Promise.all([
    fetchStatus(),
    fetchRankings(),
    fetchHistory()
  ])
  loading.value = false
})
</script>

<style scoped>
.sign-in-page {
  padding: 24px;
  max-width: 800px;
  margin: 0 auto;
}

.sign-in-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
  overflow: hidden;
}

.card-header {
  display: flex;
  align-items: center;
  padding: 16px 24px;
  border-bottom: 1px solid #f0f0f0;
}

.back-btn {
  background: none;
  border: none;
  font-size: 20px;
  color: #333;
  cursor: pointer;
  margin-right: 16px;
  padding: 8px;
  border-radius: 6px;
  
  &:hover {
    background: #f5f5f5;
  }
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin: 0;
}

.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48px;
}

.loading-spinner {
  width: 48px;
  height: 48px;
  border: 4px solid #f0f0f0;
  border-top-color: var(--primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.sign-in-content {
  padding: 32px;
  text-align: center;
}

.avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 16px;
  
  &.success {
    background: #52C41A;
    i {
      font-size: 48px;
      color: white;
    }
  }
  
  &.primary {
    background: var(--primary);
    i {
      font-size: 48px;
      color: white;
    }
  }
}

.text-h5 {
  font-size: 20px;
  font-weight: 600;
  color: #333;
}

.text-body-2 {
  font-size: 14px;
}

.text-grey {
  color: #999;
}

.text-primary {
  color: var(--primary);
}

.text-success {
  color: #52C41A;
}

.font-weight-bold {
  font-weight: 600;
}

.mb-2 {
  margin-bottom: 8px;
}

.mb-4 {
  margin-bottom: 16px;
}

.stats-row {
  display: flex;
  justify-content: center;
  gap: 48px;
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid #f0f0f0;
}

.stat-item {
  text-align: center;
}

.stat-value {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 12px;
  color: #999;
}

.stats-cards {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  padding: 0 24px 24px;
}

.stat-card {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 16px;
  text-align: center;
}

.stat-card i {
  font-size: 24px;
  margin-bottom: 8px;
}

.text-orange {
  color: #FAAD14;
}

.text-blue {
  color: var(--primary);
}

.stat-num {
  font-size: 24px;
  font-weight: 600;
  color: #333;
  margin-bottom: 4px;
}

.stat-desc {
  font-size: 12px;
  color: #999;
}

.reward-card,
.rank-card,
.history-card {
  margin: 0 24px 24px;
  background: #f8f9fa;
  border-radius: 8px;
  padding: 20px;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.reward-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.reward-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 0;
  border-bottom: 1px solid #e8e8e8;
  
  &:last-child {
    border-bottom: none;
  }
}

.text-amber {
  color: #FAAD14;
}

.rank-tabs {
  display: flex;
  gap: 4px;
}

.rank-tab {
  padding: 4px 12px;
  border: 1px solid #e8e8e8;
  border-radius: 4px;
  background: white;
  font-size: 12px;
  cursor: pointer;
  color: #666;
  
  &.active {
    background: var(--primary);
    border-color: var(--primary);
    color: white;
  }
}

.rank-list {
  display: flex;
  flex-direction: column;
}

.rank-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 0;
  border-bottom: 1px solid #e8e8e8;
  
  &:last-child {
    border-bottom: none;
  }
}

.rank-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
  color: white;
  
  &.gold { background: #FFD700; }
  &.silver { background: #C0C0C0; }
  &.bronze { background: #CD7F32; }
  &.grey { background: #e8e8e8; color: #666; }
}

.rank-info {
  flex: 1;
}

.rank-name {
  font-size: 14px;
  color: #333;
}

.rank-score {
  font-size: 14px;
  font-weight: 600;
}

.empty-rank,
.empty-history {
  text-align: center;
  padding: 24px;
  color: #999;
}

.history-list {
  display: flex;
  flex-direction: column;
}

.history-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 0;
  border-bottom: 1px solid #e8e8e8;
  
  &:last-child {
    border-bottom: none;
  }
}

.history-icon {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: #52C41A;
  display: flex;
  align-items: center;
  justify-content: center;
  
  i {
    font-size: 16px;
  }
}

.history-info {
  flex: 1;
}

.history-date {
  font-size: 14px;
  color: #333;
}

.history-desc {
  font-size: 12px;
  color: #999;
}

.history-reward {
  font-size: 14px;
  font-weight: 600;
}

.history-pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  margin-top: 16px;
}

.page-btn {
  width: 32px;
  height: 32px;
  border: 1px solid #e8e8e8;
  border-radius: 4px;
  background: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  
  &:hover:not(:disabled) {
    border-color: var(--primary);
    color: var(--primary);
  }
  
  &:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }
}

.page-info {
  font-size: 14px;
  color: #666;
}

.snackbar {
  position: fixed;
  top: 24px;
  left: 50%;
  transform: translateX(-50%);
  padding: 12px 24px;
  border-radius: 8px;
  color: white;
  font-size: 14px;
  z-index: 1000;
  display: flex;
  align-items: center;
  gap: 12px;
  
  &.success {
    background: #52C41A;
  }
  
  &.error {
    background: #FF4D4F;
  }
}

.snackbar-close {
  background: none;
  border: none;
  color: white;
  cursor: pointer;
  font-size: 16px;
}

.mr-2 {
  margin-right: 8px;
}
</style>