<template>
  <div class="statistics-page">
    <!-- 页面标题 -->
    <div class="page-head">
      <div>
        <h1 class="page-title">数据统计</h1>
        <p class="page-desc">查看论坛运营的核心指标与趋势</p>
      </div>
    </div>

    <!-- 系统概览 -->
    <div class="card-surface section-block">
      <div class="card-header">
        <el-icon :size="18" class="card-header-icon"><DataAnalysis /></el-icon>
        <span>系统概览</span>
      </div>
      <div class="stat-grid">
        <div class="stat-card" v-for="c in overviewCards" :key="c.key">
          <div class="stat-icon" :class="'tone-' + c.tone">
            <el-icon :size="22"><component :is="c.icon" /></el-icon>
          </div>
          <div class="stat-value">{{ c.field }}</div>
          <div class="stat-label">{{ c.label }}</div>
        </div>
      </div>
    </div>

    <div class="panel-grid mt-4">
      <!-- 最近7天统计 -->
      <div class="card-surface section-block">
        <div class="card-header">
          <el-icon :size="18" class="card-header-icon"><TrendCharts /></el-icon>
          <span>最近7天活跃度</span>
        </div>
        <div class="table-wrap">
          <el-table :data="recentStats" stripe class="campus-table">
            <el-table-column prop="date" label="日期" min-width="100" />
            <el-table-column prop="new_users" label="新增用户" min-width="90" align="center" />
            <el-table-column prop="active_users" label="活跃用户" min-width="90" align="center" />
            <el-table-column prop="new_articles" label="新增文章" min-width="90" align="center" />
            <el-table-column prop="new_comments" label="新增评论" min-width="90" align="center" />
          </el-table>
        </div>
      </div>

      <!-- 热门文章 -->
      <div class="card-surface section-block">
        <div class="card-header">
          <el-icon :size="18" class="card-header-icon"><HotWater /></el-icon>
          <span>热门文章</span>
        </div>
        <el-empty v-if="!hotArticles.length" description="暂无热门文章" :image-size="64" />
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
      </div>
    </div>

    <!-- 活跃用户 -->
    <div class="card-surface section-block">
      <div class="card-header">
        <el-icon :size="18" class="card-header-icon"><User /></el-icon>
        <span>活跃用户</span>
      </div>
      <div class="user-grid">
        <div class="user-card" v-for="user in activeUsers" :key="user.id">
          <el-avatar :size="56" :src="user.avatar || '/default-avatar.png'" class="user-avatar" />
          <div class="user-name">{{ user.display_name || user.username }}</div>
          <el-tag size="small" type="success" round>
            <el-icon :size="12" class="mr-1"><CircleCheckFilled /></el-icon>
            在线
          </el-tag>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
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

const overviewCards = computed(() => [
  { key: 'users', field: overview.value.total_users, label: '总用户数', icon: 'UserFilled', tone: 'primary' },
  { key: 'articles', field: overview.value.total_articles, label: '总文章数', icon: 'Document', tone: 'green' },
  { key: 'comments', field: overview.value.total_comments, label: '总评论数', icon: 'ChatDotRound', tone: 'cyan' },
  { key: 'online', field: overview.value.online_users, label: '在线用户', icon: 'CircleCheck', tone: 'violet' }
])

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
.statistics-page {
  display: flex;
  flex-direction: column;
}

.page-head {
  margin-bottom: 20px;
}

.page-title {
  margin: 0 0 4px;
  font-size: 22px;
  font-weight: 700;
  letter-spacing: -0.01em;
  color: var(--campus-text);
}

.page-desc {
  margin: 0;
  font-size: 14px;
  color: var(--campus-text-secondary);
}

/* ---------------- 卡片 ---------------- */
.section-block {
  padding: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 700;
  font-size: 16px;
  color: var(--campus-text);
  margin-bottom: 16px;
}

.card-header-icon {
  color: var(--campus-primary);
}

/* ---------------- 概览统计 ---------------- */
.stat-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.stat-card {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 18px;
  border: 1px solid var(--campus-border);
  border-radius: var(--campus-radius);
  background: var(--campus-surface);
  transition: var(--campus-transition);
}

.stat-card:hover {
  box-shadow: var(--campus-shadow);
  transform: translateY(-2px);
}

.stat-icon {
  width: 42px;
  height: 42px;
  display: grid;
  place-items: center;
  border-radius: 12px;
  margin-bottom: 10px;
}

.tone-primary { background: rgba(79, 110, 247, 0.1); color: #4f6ef7; }
.tone-green { background: rgba(34, 197, 94, 0.12); color: #22c55e; }
.tone-cyan { background: rgba(6, 182, 212, 0.12); color: #0891b2; }
.tone-violet { background: rgba(124, 58, 237, 0.12); color: #7c3aed; }

.stat-value {
  font-size: 24px;
  font-weight: 800;
  letter-spacing: -0.02em;
  color: var(--campus-text);
  line-height: 1.2;
}

.stat-label {
  font-size: 13px;
  color: var(--campus-text-secondary);
}

/* ---------------- 数据面板网格 ---------------- */
.panel-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  align-items: start;
}

.mt-4 {
  margin-top: 16px;
}

/* ---------------- 表格 ---------------- */
.table-wrap {
  overflow-x: auto;
}

.campus-table {
  width: 100%;
}

.campus-table :deep(.el-table__row) {
  font-size: 13px;
}

/* ---------------- 热门文章 ---------------- */
.hot-article {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 12px 4px;
  border-bottom: 1px solid var(--campus-border);
}

.hot-article:last-child {
  border-bottom: none;
}

.hot-article-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--campus-text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 320px;
}

.hot-article-meta {
  margin-top: 4px;
  font-size: 12px;
  color: var(--campus-text-muted);
  display: flex;
  align-items: center;
  gap: 4px;
}

.ml-3 {
  margin-left: 12px;
}

/* ---------------- 活跃用户 ---------------- */
.user-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.user-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 20px 12px;
  border: 1px solid var(--campus-border);
  border-radius: var(--campus-radius);
  background: var(--campus-surface-2);
  transition: var(--campus-transition);
}

.user-card:hover {
  box-shadow: var(--campus-shadow);
  transform: translateY(-2px);
}

.user-avatar {
  background: var(--campus-primary-soft);
}

.user-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--campus-text);
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.mr-1 {
  margin-right: 4px;
}

/* ---------------- 响应式 ---------------- */
@media (max-width: 1200px) {
  .stat-grid { grid-template-columns: repeat(2, 1fr); }
  .user-grid { grid-template-columns: repeat(3, 1fr); }
}

@media (max-width: 768px) {
  .panel-grid { grid-template-columns: 1fr; }
  .user-grid { grid-template-columns: repeat(2, 1fr); }
}

@media (max-width: 575px) {
  .stat-grid { grid-template-columns: 1fr 1fr; }
  .user-grid { grid-template-columns: 1fr 1fr; }
  .section-block { padding: 16px; }
}
</style>
