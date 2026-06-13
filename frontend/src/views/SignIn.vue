<template>
  <v-container class="pa-4">
    <!-- 页面标题 -->
    <div class="d-flex align-center mb-4">
      <v-btn icon variant="text" @click="$router.back()" class="mr-2">
        <v-icon>mdi-arrow-left</v-icon>
      </v-btn>
      <h2 class="text-h5 font-weight-bold">每日签到</h2>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="d-flex justify-center align-center py-10">
      <v-progress-circular indeterminate color="primary"></v-progress-circular>
    </div>

    <template v-else>
      <!-- 签到卡片 -->
      <v-card class="mb-4" elevation="2">
        <v-card-text class="text-center pa-6">
          <!-- 签到状态显示 -->
          <div v-if="status.has_signed_in" class="signed-in-section">
            <v-avatar size="80" color="success" class="mb-4">
              <v-icon size="48" color="white">mdi-check-circle</v-icon>
            </v-avatar>
            <h3 class="text-h5 mb-2">今日已签到</h3>
            <p class="text-body-2 text-grey">明天再来领取更多积分吧~</p>
          </div>

          <!-- 未签到状态 -->
          <div v-else class="unsigned-section">
            <v-avatar size="80" color="primary" class="mb-4">
              <v-icon size="48" color="white">mdi-calendar-check</v-icon>
            </v-avatar>
            <h3 class="text-h5 mb-2">今日可签到</h3>
            <p class="text-body-2 text-grey mb-4">
              连续签到 <span class="text-primary font-weight-bold">{{ status.sign_in_days }}</span> 天
            </p>
            <v-btn
              color="primary"
              size="large"
              :loading="signing"
              @click="handleSignIn"
            >
              <v-icon left>mdi-pencil</v-icon>
              立即签到
            </v-btn>
          </div>

          <!-- 签到统计 -->
          <v-row class="mt-6" dense>
            <v-col cols="4">
              <div class="stat-item">
                <div class="stat-value text-h6 font-weight-bold text-primary">
                  {{ status.total_points || 0 }}
                </div>
                <div class="stat-label text-caption text-grey">累计积分</div>
              </div>
            </v-col>
            <v-col cols="4">
              <div class="stat-item">
                <div class="stat-value text-h6 font-weight-bold text-primary">
                  {{ status.sign_in_days || 0 }}
                </div>
                <div class="stat-label text-caption text-grey">连续天数</div>
              </div>
            </v-col>
            <v-col cols="4">
              <div class="stat-item">
                <div class="stat-value text-h6 font-weight-bold text-primary">
                  {{ status.total_sign_ins || 0 }}
                </div>
                <div class="stat-label text-caption text-grey">累计次数</div>
              </div>
            </v-col>
          </v-row>
        </v-card-text>
      </v-card>

      <!-- 本周本月统计 -->
      <v-row dense class="mb-4">
        <v-col cols="6">
          <v-card elevation="1">
            <v-card-text class="text-center pa-3">
              <v-icon color="orange" class="mb-1">mdi-calendar-week</v-icon>
              <div class="text-h6 font-weight-bold">{{ status.week_sign_in_count || 0 }}</div>
              <div class="text-caption text-grey">本周签到</div>
            </v-card-text>
          </v-card>
        </v-col>
        <v-col cols="6">
          <v-card elevation="1">
            <v-card-text class="text-center pa-3">
              <v-icon color="blue" class="mb-1">mdi-calendar-month</v-icon>
              <div class="text-h6 font-weight-bold">{{ status.month_sign_in_count || 0 }}</div>
              <div class="text-caption text-grey">本月签到</div>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <!-- 连续签到奖励说明 -->
      <v-card class="mb-4" elevation="1">
        <v-card-title class="text-subtitle-1 pb-0">
          <v-icon left size="20">mdi-gift</v-icon>
          连续签到奖励
        </v-card-title>
        <v-card-text>
          <v-list density="compact" class="py-0">
            <v-list-item class="px-0">
              <template v-slot:prepend>
                <v-icon color="amber" size="20">mdi-circle</v-icon>
              </template>
              <v-list-item-title class="text-body-2">
                连续 <span class="text-primary font-weight-bold">7</span> 天
              </v-list-item-title>
              <template v-slot:append>
                <span class="text-success text-body-2">+5 积分</span>
              </template>
            </v-list-item>
            <v-list-item class="px-0">
              <template v-slot:prepend>
                <v-icon color="amber" size="20">mdi-circle</v-icon>
              </template>
              <v-list-item-title class="text-body-2">
                连续 <span class="text-primary font-weight-bold">30</span> 天
              </v-list-item-title>
              <template v-slot:append>
                <span class="text-success text-body-2">+15 积分</span>
              </template>
            </v-list-item>
            <v-list-item class="px-0">
              <template v-slot:prepend>
                <v-icon color="amber" size="20">mdi-circle</v-icon>
              </template>
              <v-list-item-title class="text-body-2">
                连续 <span class="text-primary font-weight-bold">365</span> 天
              </v-list-item-title>
              <template v-slot:append>
                <span class="text-success text-body-2">+50 积分</span>
              </template>
            </v-list-item>
          </v-list>
        </v-card-text>
      </v-card>

      <!-- 排行榜 -->
      <v-card class="mb-4" elevation="1">
        <v-card-title class="d-flex justify-space-between align-center">
          <span class="text-subtitle-1">
            <v-icon left size="20">mdi-trophy</v-icon>
            签到排行榜
          </span>
          <v-btn-toggle v-model="rankType" mandatory density="compact" variant="outlined" divided>
            <v-btn value="continuous" size="small">连续</v-btn>
            <v-btn value="points" size="small">积分</v-btn>
          </v-btn-toggle>
        </v-card-title>
        <v-card-text class="pa-0">
          <v-list density="compact" v-if="rankType === 'continuous'">
            <v-list-item
              v-for="(user, index) in rankings.continuous_rankings"
              :key="user.id"
              class="ranking-item"
            >
              <template v-slot:prepend>
                <v-avatar
                  size="32"
                  :color="index < 3 ? ['amber', 'grey', 'brown'][index] : 'grey-lighten-2'"
                >
                  <span class="text-body-2 font-weight-bold">{{ index + 1 }}</span>
                </v-avatar>
              </template>
              <v-list-item-title class="text-body-2">
                {{ user.display_name || user.username }}
              </v-list-item-title>
              <template v-slot:append>
                <span class="text-primary text-body-2 font-weight-bold">
                  {{ user.sign_in_days }} 天
                </span>
              </template>
            </v-list-item>
            <v-list-item v-if="!rankings.continuous_rankings?.length">
              <v-list-item-title class="text-center text-grey">暂无数据</v-list-item-title>
            </v-list-item>
          </v-list>

          <v-list density="compact" v-else>
            <v-list-item
              v-for="(user, index) in rankings.points_rankings"
              :key="user.id"
              class="ranking-item"
            >
              <template v-slot:prepend>
                <v-avatar
                  size="32"
                  :color="index < 3 ? ['amber', 'grey', 'brown'][index] : 'grey-lighten-2'"
                >
                  <span class="text-body-2 font-weight-bold">{{ index + 1 }}</span>
                </v-avatar>
              </template>
              <v-list-item-title class="text-body-2">
                {{ user.display_name || user.username }}
              </v-list-item-title>
              <template v-slot:append>
                <span class="text-primary text-body-2 font-weight-bold">
                  {{ user.total_points }} 分
                </span>
              </template>
            </v-list-item>
            <v-list-item v-if="!rankings.points_rankings?.length">
              <v-list-item-title class="text-center text-grey">暂无数据</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-card-text>
      </v-card>

      <!-- 签到历史 -->
      <v-card elevation="1">
        <v-card-title class="text-subtitle-1">
          <v-icon left size="20">mdi-history</v-icon>
          签到记录
        </v-card-title>
        <v-card-text class="pa-0">
          <v-list density="compact">
            <v-list-item
              v-for="record in history.records"
              :key="record.id"
              class="history-item"
            >
              <template v-slot:prepend>
                <v-avatar size="36" color="success" class="mr-2">
                  <v-icon size="20" color="white">mdi-check</v-icon>
                </v-avatar>
              </template>
              <v-list-item-title class="text-body-2">
                {{ record.sign_in_date }}
              </v-list-item-title>
              <v-list-item-subtitle class="text-caption">
                连续 {{ record.continuous_day }} 天
              </v-list-item-subtitle>
              <template v-slot:append>
                <span class="text-success text-body-2">
                  +{{ record.reward_points }}
                </span>
              </template>
            </v-list-item>
            <v-list-item v-if="!history.records?.length">
              <v-list-item-title class="text-center text-grey">暂无签到记录</v-list-item-title>
            </v-list-item>
          </v-list>

          <!-- 分页 -->
          <v-pagination
            v-if="history.total_pages > 1"
            v-model="page"
            :length="history.total_pages"
            :total-visible="5"
            density="compact"
            class="my-2"
            @update:model-value="fetchHistory"
          ></v-pagination>
        </v-card-text>
      </v-card>
    </template>

    <!-- 签到成功提示 -->
    <v-snackbar v-model="snackbar.show" :color="snackbar.color" location="top">
      {{ snackbar.text }}
      <template v-slot:actions>
        <v-btn variant="text" @click="snackbar.show = false">关闭</v-btn>
      </template>
    </v-snackbar>
  </v-container>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { signinApi } from '@/api'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()

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
    snackbar.text = res.data?.message || '签到成功'
    snackbar.color = 'success'
    snackbar.show = true

    // 更新状态
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
.stat-item {
  text-align: center;
}

.ranking-item,
.history-item {
  border-bottom: 1px solid #f0f0f0;
}

.ranking-item:last-child,
.history-item:last-child {
  border-bottom: none;
}
</style>
