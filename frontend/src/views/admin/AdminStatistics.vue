<template>
  <div>
    <el-card class="page-container" shadow="never">
      <!-- 系统概览 -->
      <template #header>
        <div class="card-header">
          <el-icon :size="18"><DataAnalysis /></el-icon>
          <span>系统概览</span>
        </div>
      </template>

      <el-row :gutter="16">
        <el-col :xs="24" :sm="12" :md="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-card-inner">
              <el-icon :size="40" color="#409EFF"><UserFilled /></el-icon>
              <div class="stat-value mt-2">{{ overview.total_users }}</div>
              <div class="stat-label">总用户数</div>
            </div>
          </el-card>
        </el-col>
        <el-col :xs="24" :sm="12" :md="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-card-inner">
              <el-icon :size="40" color="#67C23A"><Document /></el-icon>
              <div class="stat-value mt-2">{{ overview.total_articles }}</div>
              <div class="stat-label">总文章数</div>
            </div>
          </el-card>
        </el-col>
        <el-col :xs="24" :sm="12" :md="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-card-inner">
              <el-icon :size="40" color="#909399"><ChatDotRound /></el-icon>
              <div class="stat-value mt-2">{{ overview.total_comments }}</div>
              <div class="stat-label">总评论数</div>
            </div>
          </el-card>
        </el-col>
        <el-col :xs="24" :sm="12" :md="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-card-inner">
              <el-icon :size="40" color="#E6A23C"><CircleCheck /></el-icon>
              <div class="stat-value mt-2">{{ overview.online_users }}</div>
              <div class="stat-label">在线用户</div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </el-card>

    <el-row :gutter="16" class="mt-4">
      <!-- 最近7天统计 -->
      <el-col :xs="24" :md="12">
        <el-card class="page-container" shadow="never">
          <template #header>
            <div class="card-header">
              <el-icon :size="18"><TrendCharts /></el-icon>
              <span>最近7天活跃度</span>
            </div>
          </template>
          <el-table :data="recentStats" stripe>
            <el-table-column prop="date" label="日期" />
            <el-table-column prop="new_users" label="新增用户" />
            <el-table-column prop="active_users" label="活跃用户" />
            <el-table-column prop="new_articles" label="新增文章" />
            <el-table-column prop="new_comments" label="新增评论" />
          </el-table>
        </el-card>
      </el-col>

      <!-- 热门文章 -->
      <el-col :xs="24" :md="12">
        <el-card class="page-container" shadow="never">
          <template #header>
            <div class="card-header">
              <el-icon :size="18"><HotWater /></el-icon>
              <span>热门文章</span>
            </div>
          </template>
          <el-empty v-if="!hotArticles.length" description="暂无热门文章" />
          <div v-for="article in hotArticles" :key="article.id" class="hot-article">
            <div class="hot-article-info">
              <div class="hot-article-title">{{ article.title }}</div>
              <div class="hot-article-meta">
                <el-icon :size="14"><View /></el-icon> {{ article.view_count }}
                <el-icon :size="14" class="ml-3"><Star /></el-icon> {{ article.like_count }}
                <el-icon :size="14" class="ml-3"><ChatDotSquare /></el-icon> {{ article.comment_count }}
              </div>
            </div>
            <el-button type="primary" text size="small" @click="goArticle(article.id)">查看</el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 活跃用户 -->
    <el-card class="page-container mt-4" shadow="never">
      <template #header>
        <div class="card-header">
          <el-icon :size="18"><User /></el-icon>
          <span>活跃用户</span>
        </div>
      </template>
      <el-row :gutter="16">
        <el-col
          v-for="user in activeUsers"
          :key="user.id"
          :xs="24"
          :sm="12"
          :md="8"
          :lg="6"
        >
          <el-card shadow="hover" class="stat-card">
            <div class="stat-card-inner">
              <el-avatar :size="60" :src="user.avatar || '/default-avatar.png'" />
              <div class="mt-2">{{ user.display_name || user.username }}</div>
              <el-tag size="small" type="success" class="mt-1">
                <el-icon :size="12" class="mr-1"><CircleCheckFilled /></el-icon>
                在线
              </el-tag>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  DataAnalysis,
  UserFilled,
  Document,
  ChatDotRound,
  CircleCheck,
  TrendCharts,
  HotWater,
  View,
  Star,
  ChatDotSquare,
  User,
  CircleCheckFilled
} from '@element-plus/icons-vue'
import { statisticsApi } from '@/api'

const router = useRouter()

const overview = ref({
  total_users: 0,
  total_articles: 0,
  total_comments: 0,
  online_users: 0
})
const recentStats = ref([])
const hotArticles = ref([])
const activeUsers = ref([])

const loadDashboard = async () => {
  try {
    const res = await statisticsApi.getStatisticsDashboard()
    if (res.data.success) {
      overview.value = res.data.data.overview
      recentStats.value = res.data.data.recent_stats
      hotArticles.value = res.data.data.hot_articles
      activeUsers.value = res.data.data.active_users
    }
  } catch (err) {
    console.error('加载仪表板数据失败:', err)
  }
}

const goArticle = (id) => {
  router.push(`/article/${id}`)
}

onMounted(() => {
  loadDashboard()
})
</script>

<style scoped>
.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.stat-card {
  margin-bottom: 16px;
}

.stat-card-inner {
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.stat-value {
  font-size: 22px;
  font-weight: 600;
  line-height: 1.4;
}

.stat-label {
  color: #909399;
  font-size: 13px;
}

.hot-article {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 4px;
  border-bottom: 1px solid #f0f2f5;
}

.hot-article:last-child {
  border-bottom: none;
}

.hot-article-title {
  font-size: 14px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 320px;
}

.hot-article-meta {
  margin-top: 4px;
  font-size: 12px;
  color: #909399;
  display: flex;
  align-items: center;
  gap: 4px;
}

.ml-3 {
  margin-left: 12px;
}

.mr-1 {
  margin-right: 4px;
}

.mt-1 {
  margin-top: 4px;
}

.mt-2 {
  margin-top: 8px;
}

.mt-4 {
  margin-top: 16px;
}
</style>
