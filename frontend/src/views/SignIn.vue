<template>
  <div class="page-container signin-page">
    <!-- 页面标题 -->
    <div class="header-row">
      <button class="back-btn" @click="$router.back()" aria-label="返回">
        <el-icon :size="20"><ArrowLeft /></el-icon>
      </button>
      <h2 class="page-title">每日签到</h2>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-wrap">
      <div class="loading"></div>
    </div>

    <template v-else>
      <div class="signin-layout">
        <!-- 左侧：签到主卡片 -->
        <div class="signin-main">
          <div class="card-surface signin-card">
            <!-- 状态图 -->
            <div v-if="status.has_signed_in" class="visual">
              <div class="streak-ring done">
                <el-icon :size="56"><CircleCheckFilled /></el-icon>
              </div>
              <h3 class="visual-title">今日已签到</h3>
              <p class="visual-sub">明天再来领取更多积分吧~</p>
            </div>
            <div v-else class="visual">
              <div class="streak-ring">
                <span class="streak-day">{{ status.sign_in_days || 0 }}</span>
                <span class="streak-unit">天</span>
              </div>
              <h3 class="visual-title">连续签到</h3>
              <p class="visual-sub">
                已连续签到
                <span class="text-primary stat-strong">{{ status.sign_in_days || 0 }}</span> 天
              </p>
              <button class="btn btn-primary signin-btn" :disabled="signing" @click="handleSignIn">
                <el-icon v-if="!signing" style="margin-right:6px"><EditPen /></el-icon>
                {{ signing ? '签到中…' : '立即签到' }}
              </button>
            </div>

            <!-- 签到统计 -->
            <div class="stats-row">
              <div class="stat-cell">
                <div class="stat-value text-primary">{{ status.total_coins || status.total_points || 0 }}</div>
                <div class="stat-label">累计币</div>
              </div>
              <div class="stat-cell">
                <div class="stat-value text-primary">{{ status.sign_in_days || 0 }}</div>
                <div class="stat-label">连续天数</div>
              </div>
              <div class="stat-cell">
                <div class="stat-value text-primary">{{ status.total_sign_ins || 0 }}</div>
                <div class="stat-label">累计次数</div>
              </div>
            </div>

            <!-- 本周 / 本月 -->
            <div class="period-row">
              <div class="period-cell">
                <div class="period-icon week"><el-icon :size="18"><Calendar /></el-icon></div>
                <div class="period-value">{{ status.week_sign_in_count || 0 }}</div>
                <div class="period-label">本周签到</div>
              </div>
              <div class="period-cell">
                <div class="period-icon month"><el-icon :size="18"><Calendar /></el-icon></div>
                <div class="period-value">{{ status.month_sign_in_count || 0 }}</div>
                <div class="period-label">本月签到</div>
              </div>
            </div>
          </div>

          <!-- 连续签到奖励说明 -->
          <div class="card-surface reward-card">
            <div class="section-title">
              <el-icon :size="18" class="title-icon"><Present /></el-icon>
              连续签到奖励
            </div>
            <div class="reward-item">
              <span class="reward-days">7 天</span>
              <span class="reward-bar"><i style="width:25%"></i></span>
              <span class="reward-coin">+5 币</span>
            </div>
            <div class="reward-item">
              <span class="reward-days">30 天</span>
              <span class="reward-bar"><i style="width:70%"></i></span>
              <span class="reward-coin">+15 币</span>
            </div>
            <div class="reward-item">
              <span class="reward-days">365 天</span>
              <span class="reward-bar"><i style="width:100%"></i></span>
              <span class="reward-coin">+50 币</span>
            </div>
          </div>
        </div>

        <!-- 右侧：排行榜 -->
        <div class="signin-side">
          <!-- 排行榜 -->
          <div class="card-surface rank-card">
            <div class="rank-head">
              <span class="section-title">
                <el-icon :size="18" class="title-icon"><Trophy /></el-icon>
                签到排行榜
              </span>
              <div class="rank-switch">
                <button
                  class="rank-switch-btn"
                  :class="{ active: rankType === 'continuous' }"
                  @click="rankType = 'continuous'"
                >连续</button>
                <button
                  class="rank-switch-btn"
                  :class="{ active: rankType === 'points' }"
                  @click="rankType = 'points'"
                >币</button>
              </div>
            </div>

            <template v-if="rankType === 'continuous'">
              <div v-if="rankings.continuous_rankings?.length" class="rank-list">
                <div
                  v-for="(user, index) in rankings.continuous_rankings"
                  :key="user.id"
                  class="rank-row"
                >
                  <span class="rank-no" :class="'rank-no-' + (index + 1)">{{ index + 1 }}</span>
                  <span class="rank-name">{{ user.display_name || user.username }}</span>
                  <span class="rank-value text-primary">{{ user.sign_in_days }} 天</span>
                </div>
              </div>
              <div v-else class="empty-text">暂无数据</div>
            </template>
            <template v-else>
              <div v-if="rankings.points_rankings?.length" class="rank-list">
                <div
                  v-for="(user, index) in rankings.points_rankings"
                  :key="user.id"
                  class="rank-row"
                >
                  <span class="rank-no" :class="'rank-no-' + (index + 1)">{{ index + 1 }}</span>
                  <span class="rank-name">{{ user.display_name || user.username }}</span>
                  <span class="rank-value text-primary">{{ user.total_coins || user.total_points }} 币</span>
                </div>
              </div>
              <div v-else class="empty-text">暂无数据</div>
            </template>
          </div>
        </div>
      </div>

      <!-- 签到历史 -->
      <div class="card-surface history-card">
        <div class="section-title history-title">
          <el-icon :size="18" class="title-icon"><Clock /></el-icon>
          签到记录
        </div>
        <div v-if="history.records?.length" class="history-list">
          <div v-for="record in history.records" :key="record.id" class="history-item">
            <div class="history-check"><el-icon :size="16"><Check /></el-icon></div>
            <div class="history-info">
              <span class="history-date">{{ record.sign_in_date }}</span>
              <span class="history-sub">连续 {{ record.continuous_day }} 天</span>
            </div>
            <span class="history-coin">+{{ record.reward_points }} 币</span>
          </div>
        </div>
        <div v-else class="empty-text">暂无签到记录</div>

        <el-pagination
          v-if="history.total_pages > 1"
          v-model:current-page="page"
          :page-size="history.page_size || 30"
          :total="history.total"
          layout="prev, pager, next"
          :pager-count="5"
          class="pagination"
          @current-change="fetchHistory"
        ></el-pagination>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import {
  ArrowLeft,
  Calendar,
  Check,
  CircleCheckFilled,
  Clock,
  EditPen,
  Present,
  Trophy
} from '@element-plus/icons-vue'
import { signinApi } from '@/api'
import { success, error } from '@/utils/message'

const loading = ref(true)
const signing = ref(false)
const rankType = ref('continuous')

const status = ref({
  has_signed_in: false,
  sign_in_days: 0,
  total_sign_ins: 0,
  max_continuous_days: 0,
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

// 获取签到状态
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

// 获取排行榜
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

// 获取签到历史
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

// 签到
const handleSignIn = async () => {
  signing.value = true
  try {
    const res = await signinApi.signIn()
    success(res.data?.message || '签到成功')

    // 更新状态
    await fetchStatus()
    await fetchRankings()
    await fetchHistory()
  } catch (err) {
    error(err.response?.data?.error || '签到失败')
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
.signin-page {
  padding-bottom: 40px;
}

.header-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;
}

.back-btn {
  width: 38px;
  height: 38px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  border-radius: 12px;
  background: var(--campus-surface);
  color: var(--campus-text);
  cursor: pointer;
  box-shadow: var(--campus-shadow-sm);
  transition: var(--campus-transition);
}

.back-btn:hover {
  background: var(--campus-primary-soft);
  color: var(--campus-primary);
}

.page-title {
  margin: 0;
  font-size: 22px;
  font-weight: 800;
  color: var(--campus-text);
}

.loading-wrap {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 56px 0;
}

.loading {
  width: 40px;
  height: 40px;
  border: 4px solid var(--campus-primary-light);
  border-top-color: var(--campus-primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* ---------------- 布局 ---------------- */
.signin-layout {
  display: grid;
  grid-template-columns: 1fr 360px;
  gap: 20px;
  align-items: start;
}

.signin-main {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.signin-side {
  min-width: 0;
}

/* ---------------- 签到主卡片 ---------------- */
.signin-card {
  padding: 32px;
  border-radius: var(--campus-radius-lg);
  box-shadow: var(--campus-shadow);
  text-align: center;
  background: linear-gradient(160deg, var(--campus-surface) 0%, var(--campus-surface-2) 100%);
}

.visual {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.visual-title {
  margin: 16px 0 4px;
  font-size: 22px;
  font-weight: 800;
  color: var(--campus-text);
}

.visual-sub {
  margin: 0 0 4px;
  font-size: 14px;
  color: var(--campus-text-secondary);
}

.streak-ring {
  width: 120px;
  height: 120px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: conic-gradient(from 0deg, #4f6ef7, #9db4ff);
  color: #fff;
  box-shadow: 0 0 0 8px var(--campus-primary-soft);
}

.streak-ring.done {
  background: conic-gradient(from 0deg, #22c55e, #4ade80);
}

.streak-day {
  font-size: 44px;
  line-height: 1;
  font-weight: 800;
}

.streak-unit {
  font-size: 14px;
  opacity: 0.9;
  margin-top: 2px;
}

.signed-visual {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.stat-strong {
  font-weight: 800;
}

.signin-btn {
  margin-top: 18px;
  min-width: 160px;
  padding: 12px 32px;
  font-size: 16px;
  border-radius: 14px;
}

.signin-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

/* ---------------- 统计 ---------------- */
.stats-row {
  display: flex;
  margin-top: 28px;
  padding-top: 20px;
  border-top: 1px solid var(--campus-border);
}

.stat-cell {
  flex: 1;
  text-align: center;
}

.stat-value {
  font-size: 24px;
  font-weight: 800;
}

.stat-label {
  font-size: 12px;
  color: var(--campus-text-secondary);
  margin-top: 2px;
}

/* ---------------- 本周本月 ---------------- */
.period-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  margin-top: 20px;
}

.period-cell {
  background: var(--campus-surface);
  border: 1px solid var(--campus-border);
  border-radius: 14px;
  padding: 16px;
  text-align: center;
}

.period-icon {
  width: 34px;
  height: 34px;
  margin: 0 auto 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 10px;
}

.period-icon.week {
  background: rgba(245, 158, 11, 0.12);
  color: #f59e0b;
}

.period-icon.month {
  background: rgba(79, 110, 247, 0.12);
  color: var(--campus-primary);
}

.period-value {
  font-size: 20px;
  font-weight: 800;
  color: var(--campus-primary);
}

.period-label {
  font-size: 12px;
  color: var(--campus-text-secondary);
  margin-top: 2px;
}

/* ---------------- 通用卡片 ---------------- */
.reward-card,
.rank-card,
.history-card {
  border-radius: var(--campus-radius-lg);
  box-shadow: var(--campus-shadow-sm);
}

.reward-card,
.rank-card {
  padding: 20px;
}

.history-card {
  padding: 20px;
  margin-top: 20px;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 700;
  color: var(--campus-text);
  margin-bottom: 16px;
}

.title-icon {
  color: var(--campus-primary);
}

/* ---------------- 奖励 ---------------- */
.reward-item {
  display: grid;
  grid-template-columns: 64px 1fr auto;
  align-items: center;
  gap: 12px;
  padding: 10px 0;
}

.reward-days {
  font-size: 14px;
  font-weight: 600;
  color: var(--campus-warning);
}

.reward-bar {
  height: 8px;
  background: var(--campus-surface-2);
  border-radius: 999px;
  overflow: hidden;
}

.reward-bar i {
  display: block;
  height: 100%;
  border-radius: 999px;
  background: linear-gradient(90deg, var(--campus-warning), #f59e0b);
}

.reward-coin {
  font-size: 14px;
  font-weight: 700;
  color: var(--campus-success);
  white-space: nowrap;
}

/* ---------------- 排行榜 ---------------- */
.rank-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 8px;
}

.rank-switch {
  display: inline-flex;
  background: var(--campus-surface-2);
  border: 1px solid var(--campus-border);
  border-radius: 10px;
  padding: 2px;
}

.rank-switch-btn {
  border: none;
  background: transparent;
  padding: 4px 12px;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 600;
  color: var(--campus-text-secondary);
  cursor: pointer;
  transition: var(--campus-transition);
}

.rank-switch-btn.active {
  background: var(--campus-surface);
  color: var(--campus-primary);
  box-shadow: var(--campus-shadow-sm);
}

.rank-list {
  display: flex;
  flex-direction: column;
}

.rank-row {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 9px 4px;
  border-bottom: 1px solid var(--campus-border);
}

.rank-row:last-child {
  border-bottom: none;
}

.rank-no {
  width: 24px;
  height: 24px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 800;
  background: var(--campus-surface-2);
  color: var(--campus-text-secondary);
}

.rank-no-1 {
  background: #ffb300;
  color: #fff;
}

.rank-no-2 {
  background: #9e9e9e;
  color: #fff;
}

.rank-no-3 {
  background: #8d6e63;
  color: #fff;
}

.rank-name {
  flex: 1;
  font-size: 14px;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.rank-value {
  font-size: 14px;
  font-weight: 700;
  white-space: nowrap;
}

/* ---------------- 历史记录 ---------------- */
.history-list {
  display: flex;
  flex-direction: column;
}

.history-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 4px;
  border-bottom: 1px solid var(--campus-border);
}

.history-item:last-child {
  border-bottom: none;
}

.history-check {
  width: 32px;
  height: 32px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: rgba(34, 197, 94, 0.15);
  color: var(--campus-success);
}

.history-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.history-date {
  font-size: 14px;
  font-weight: 600;
  color: var(--campus-text);
}

.history-sub {
  font-size: 12px;
  color: var(--campus-text-secondary);
}

.history-coin {
  font-size: 14px;
  font-weight: 700;
  color: var(--campus-success);
  white-space: nowrap;
}

.empty-text {
  text-align: center;
  color: var(--campus-text-secondary);
  padding: 24px 0;
  font-size: 14px;
}

.pagination {
  justify-content: center;
  margin-top: 16px;
  border-top: 1px solid var(--campus-border);
  padding-top: 16px;
}

/* ---------------- 按钮 ---------------- */
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 9px 18px;
  border-radius: 12px;
  border: none;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: var(--campus-transition);
}

.btn-primary {
  background: var(--campus-primary);
  color: #fff;
}

.btn-primary:hover {
  background: var(--campus-primary-dark);
  transform: translateY(-1px);
  box-shadow: var(--campus-shadow);
}

/* ---------------- 响应式 ---------------- */
@media (max-width: 820px) {
  .signin-layout {
    grid-template-columns: 1fr;
  }

  .signin-side {
    order: -1;
  }

  .signin-card {
    padding: 24px 16px;
  }
}

@media (max-width: 480px) {
  .streak-ring {
    width: 100px;
    height: 100px;
  }

  .streak-day {
    font-size: 36px;
  }

  .period-row {
    grid-template-columns: 1fr;
  }
}
</style>
