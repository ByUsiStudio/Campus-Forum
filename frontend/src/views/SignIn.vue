<template>
  <div class="page-container">
    <!-- 页面标题 -->
    <div class="header-row mb-4">
      <el-button circle :icon="ArrowLeft" class="back-btn" @click="$router.back()" />
      <h2 class="page-title">每日签到</h2>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-wrap">
      <div class="loading"></div>
    </div>

    <template v-else>
      <!-- 签到卡片 -->
      <el-card class="mb-4 card-surface" shadow="never">
        <div class="text-center card-body-pad">
          <!-- 签到状态显示 -->
          <div v-if="status.has_signed_in" class="text-center pad-16">
            <el-avatar :size="80" class="signed-avatar mb-4">
              <el-icon :size="48"><CircleCheckFilled /></el-icon>
            </el-avatar>
            <h3 class="stat-title mb-2">今日已签到</h3>
            <p class="text-secondary sub-hint">明天再来领取更多积分吧~</p>
          </div>

          <!-- 未签到状态 -->
          <div v-else class="unsigned-section">
            <el-avatar :size="80" class="unsigned-avatar mb-4">
              <el-icon :size="48"><Calendar /></el-icon>
            </el-avatar>
            <h3 class="stat-title mb-2">今日可签到</h3>
            <p class="text-secondary mb-4">
              连续签到 <span class="text-primary stat-strong">{{ status.sign_in_days }}</span> 天，获得 <span class="text-primary stat-strong">{{ status.total_coins || status.total_points || 0 }}</span> 币
            </p>
            <el-button
              type="primary"
              size="large"
              :loading="signing"
              @click="handleSignIn"
            >
              <el-icon class="btn-icon"><EditPen /></el-icon>
              立即签到
            </el-button>
          </div>

          <!-- 签到统计 -->
          <el-row :gutter="16" class="mt-4 stats-row">
            <el-col :xs="8" :sm="8" :md="8" :lg="8">
              <div class="stat-item">
                <div class="stat-value text-primary">{{ status.total_coins || status.total_points || 0 }}</div>
                <div class="stat-label text-secondary">累计币</div>
              </div>
            </el-col>
            <el-col :xs="8" :sm="8" :md="8" :lg="8">
              <div class="stat-item">
                <div class="stat-value text-primary">{{ status.sign_in_days || 0 }}</div>
                <div class="stat-label text-secondary">连续天数</div>
              </div>
            </el-col>
            <el-col :xs="8" :sm="8" :md="8" :lg="8">
              <div class="stat-item">
                <div class="stat-value text-primary">{{ status.total_sign_ins || 0 }}</div>
                <div class="stat-label text-secondary">累计次数</div>
              </div>
            </el-col>
          </el-row>
        </div>
      </el-card>

      <!-- 本周本月统计 -->
      <el-row :gutter="16" class="mb-4">
        <el-col :xs="24" :sm="12" :md="12" :lg="12">
          <el-card class="card-surface" shadow="never">
            <div class="text-center small-card-pad">
              <div class="week-icon mb-4"><el-icon :size="20"><Calendar /></el-icon></div>
              <div class="mini-value text-primary">{{ status.week_sign_in_count || 0 }}</div>
              <div class="text-secondary mini-label">本周签到</div>
            </div>
          </el-card>
        </el-col>
        <el-col :xs="24" :sm="12" :md="12" :lg="12">
          <el-card class="card-surface" shadow="never">
            <div class="text-center small-card-pad">
              <div class="month-icon mb-4"><el-icon :size="20"><Calendar /></el-icon></div>
              <div class="mini-value text-primary">{{ status.month_sign_in_count || 0 }}</div>
              <div class="text-secondary mini-label">本月签到</div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 连续签到奖励说明 -->
      <el-card class="mb-4 card-surface" shadow="never">
        <template #header>
          <div class="card-title">
            <el-icon :size="20" class="card-title-icon"><Present /></el-icon>
            连续签到奖励
          </div>
        </template>
        <div class="reward-item">
          <el-tag type="warning" effect="plain" class="reward-tag">7</el-tag>
          <span class="reward-text">
            连续 <span class="text-primary stat-strong">7</span> 天
          </span>
          <span class="reward-coin success-text">+5 币</span>
        </div>
        <div class="reward-item">
          <el-tag type="warning" effect="plain" class="reward-tag">30</el-tag>
          <span class="reward-text">
            连续 <span class="text-primary stat-strong">30</span> 天
          </span>
          <span class="reward-coin success-text">+15 币</span>
        </div>
        <div class="reward-item">
          <el-tag type="warning" effect="plain" class="reward-tag">365</el-tag>
          <span class="reward-text">
            连续 <span class="text-primary stat-strong">365</span> 天
          </span>
          <span class="reward-coin success-text">+50 币</span>
        </div>
      </el-card>

      <!-- 排行榜 -->
      <el-card class="mb-4 card-surface" shadow="never">
        <template #header>
          <div class="rank-header">
            <span class="card-title">
              <el-icon :size="20" class="card-title-icon"><Trophy /></el-icon>
              签到排行榜
            </span>
            <el-radio-group v-model="rankType" size="small">
              <el-radio-button value="continuous">连续</el-radio-button>
              <el-radio-button value="points">币</el-radio-button>
            </el-radio-group>
          </div>
        </template>
        <div v-if="rankType === 'continuous'" class="list-pad">
          <el-list v-if="rankings.continuous_rankings?.length">
            <el-list-item
              v-for="(user, index) in rankings.continuous_rankings"
              :key="user.id"
              class="ranking-item"
            >
              <el-avatar
                :size="32"
                class="rank-avatar"
                :class="rankAvatarClass(index)"
              >
                <span class="rank-number">{{ index + 1 }}</span>
              </el-avatar>
              <span class="rank-name">{{ user.display_name || user.username }}</span>
              <span class="rank-value text-primary">
                {{ user.sign_in_days }} 天
              </span>
            </el-list-item>
          </el-list>
          <div v-else class="empty-text">暂无数据</div>
        </div>

        <div v-else class="list-pad">
          <el-list v-if="rankings.points_rankings?.length">
            <el-list-item
              v-for="(user, index) in rankings.points_rankings"
              :key="user.id"
              class="ranking-item"
            >
              <el-avatar
                :size="32"
                class="rank-avatar"
                :class="rankAvatarClass(index)"
              >
                <span class="rank-number">{{ index + 1 }}</span>
              </el-avatar>
              <span class="rank-name">{{ user.display_name || user.username }}</span>
              <span class="rank-value text-primary">
                {{ user.total_coins || user.total_points }} 币
              </span>
            </el-list-item>
          </el-list>
          <div v-else class="empty-text">暂无数据</div>
        </div>
      </el-card>

      <!-- 签到历史 -->
      <el-card class="card-surface" shadow="never">
        <template #header>
          <div class="card-title">
            <el-icon :size="20" class="card-title-icon"><Clock /></el-icon>
            签到记录
          </div>
        </template>
        <div class="list-pad">
          <el-list v-if="history.records?.length">
            <el-list-item
              v-for="record in history.records"
              :key="record.id"
              class="history-item"
            >
              <el-avatar :size="36" class="history-avatar">
                <el-icon :size="20"><Check /></el-icon>
              </el-avatar>
              <div class="history-info">
                <span class="history-date">{{ record.sign_in_date }}</span>
                <span class="text-secondary history-sub">连续 {{ record.continuous_day }} 天</span>
              </div>
              <span class="success-text history-coin">
                +{{ record.reward_points }} 币
              </span>
            </el-list-item>
          </el-list>
          <div v-else class="empty-text">暂无签到记录</div>

          <!-- 分页 -->
          <el-pagination
            v-if="history.total_pages > 1"
            v-model:current-page="page"
            :page-size="history.page_size || 30"
            :total="history.total"
            layout="prev, pager, next"
            :pager-count="5"
            class="my-4 pagination"
            @current-change="fetchHistory"
          ></el-pagination>
        </div>
      </el-card>
    </template>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
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

const rankAvatarClass = (index) => {
  if (index === 0) return 'rank-first'
  if (index === 1) return 'rank-second'
  if (index === 2) return 'rank-third'
  return 'rank-default'
}

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
.header-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.back-btn {
  flex-shrink: 0;
}

.page-title {
  margin: 0;
  font-size: 20px;
  font-weight: 700;
}

.loading-wrap {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 48px 0;
}

.card-body-pad {
  padding: 24px;
}

.pad-16 {
  padding: 16px;
}

.stat-title {
  margin: 0;
  font-size: 18px;
  font-weight: 700;
  color: var(--campus-text);
}

.sub-hint {
  margin: 0;
  font-size: 14px;
}

.signed-avatar {
  background: #67c23a;
}

.unsigned-avatar {
  background: var(--campus-primary);
}

.stat-strong {
  font-weight: 700;
}

.stat-item {
  text-align: center;
}

.stat-value {
  font-size: 20px;
  font-weight: 700;
}

.stat-label {
  font-size: 12px;
}

.stats-row {
  row-gap: 16px;
}

.small-card-pad {
  padding: 12px;
}

.week-icon,
.month-icon {
  color: #ff9800;
  display: flex;
  justify-content: center;
}

.month-icon {
  color: #409eff;
}

.mini-value {
  font-size: 18px;
  font-weight: 700;
}

.mini-label {
  font-size: 12px;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: var(--campus-text);
}

.card-title-icon {
  color: var(--campus-primary);
}

.rank-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.reward-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 0;
  border-bottom: 1px solid #f0f0f0;
}

.reward-item:last-child {
  border-bottom: none;
}

.reward-tag {
  width: 48px;
  text-align: center;
}

.reward-text {
  flex: 1;
  font-size: 14px;
}

.reward-coin {
  font-size: 14px;
  font-weight: 600;
}

.success-text {
  color: #67c23a;
  font-weight: 600;
}

.list-pad {
  padding: 8px 0;
  width: 100%;
}

.ranking-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 4px;
  border-bottom: 1px solid #f0f0f0;
}

.ranking-item:last-child {
  border-bottom: none;
}

.rank-avatar {
  flex-shrink: 0;
}

.rank-first {
  background: #ffb300;
}

.rank-second {
  background: #9e9e9e;
}

.rank-third {
  background: #8d6e63;
}

.rank-default {
  background: #e0e0e0;
  color: #666;
}

.rank-number {
  font-weight: 700;
  font-size: 13px;
  color: #fff;
}

.rank-default .rank-number {
  color: #666;
}

.rank-name {
  flex: 1;
  font-size: 14px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.rank-value {
  font-size: 14px;
  font-weight: 700;
  white-space: nowrap;
}

.empty-text {
  text-align: center;
  color: var(--campus-text-secondary);
  padding: 24px 0;
  font-size: 14px;
}

.history-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 4px;
  border-bottom: 1px solid #f0f0f0;
}

.history-item:last-child {
  border-bottom: none;
}

.history-avatar {
  background: #67c23a;
  flex-shrink: 0;
}

.history-info {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.history-date {
  font-size: 14px;
}

.history-sub {
  font-size: 12px;
}

.history-coin {
  font-size: 14px;
  white-space: nowrap;
}

.pagination {
  justify-content: center;
}

@media (max-width: 768px) {
  .card-body-pad {
    padding: 16px;
  }
}
</style>
