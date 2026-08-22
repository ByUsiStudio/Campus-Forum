<template>
  <div class="overview">
    <!-- 欢迎横幅 -->
    <div class="welcome-bar">
      <div>
        <h1 class="welcome-title">欢迎回来 👋</h1>
        <p class="welcome-sub">这里是校园论坛管理后台数据概览</p>
      </div>
      <el-button type="primary" round :loading="loading" @click="$emit('refresh')">
        <el-icon class="mr-1"><Refresh /></el-icon>
        刷新数据
      </el-button>
    </div>

    <!-- 统计卡片 -->
    <div class="stat-grid">
      <div class="stat-card" v-for="s in statCards" :key="s.key">
        <div class="stat-card-top">
          <div class="stat-icon" :class="'tone-' + s.tone">
            <el-icon :size="22"><component :is="s.icon" /></el-icon>
          </div>
          <span class="stat-change" :class="'tone-text-' + s.tone">{{ s.trend }}</span>
        </div>
        <div class="stat-value">{{ formatNumber(statistics[s.field] || 0) }}</div>
        <div class="stat-label">{{ s.label }}</div>
      </div>
    </div>

    <!-- 快捷操作 -->
    <div class="card-surface section-block">
      <div class="card-header">
        <el-icon :size="20" class="card-header-icon"><Lightning /></el-icon>
        <span>快捷操作</span>
      </div>
      <div class="quick-grid">
        <el-button type="primary" plain class="quick-btn" @click="$router.push('/')">
          <el-icon><HomeFilled /></el-icon>返回首页
        </el-button>
        <el-button type="success" plain class="quick-btn" @click="$router.push('/create')">
          <el-icon><Plus /></el-icon>创建文章
        </el-button>
        <el-button type="info" plain class="quick-btn" @click="$router.push({ name: 'AdminUsers' })">
          <el-icon><User /></el-icon>用户管理
        </el-button>
        <el-button type="warning" plain class="quick-btn" @click="$router.push({ name: 'AdminArticles' })">
          <el-icon><Document /></el-icon>文章管理
        </el-button>
      </div>
    </div>

    <!-- 系统信息 -->
    <div class="card-surface">
      <div class="card-header">
        <el-icon :size="20" class="card-header-icon"><InfoFilled /></el-icon>
        <span>系统信息</span>
      </div>
      <div class="sys-info">
        <div class="sys-info-item">
          <el-icon class="sys-icon tone-text-info"><PriceTag /></el-icon>
          <span class="sys-label">系统版本</span>
          <span class="sys-value">v1.3.14</span>
        </div>
        <div class="sys-info-item">
          <el-icon class="sys-icon tone-text-success"><Coin /></el-icon>
          <span class="sys-label">数据库状态</span>
          <span class="sys-value sys-success">正常</span>
        </div>
        <div class="sys-info-item">
          <el-icon class="sys-icon tone-text-primary"><Monitor /></el-icon>
          <span class="sys-label">服务器状态</span>
          <span class="sys-value sys-info-color">运行中</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import {
  UserFilled,
  Document,
  ChatDotRound,
  View,
  Lightning,
  HomeFilled,
  Plus,
  User,
  InfoFilled,
  PriceTag,
  Coin,
  Monitor,
  Refresh
} from '@element-plus/icons-vue'

export default {
  name: 'OverviewPanel',
  emits: ['refresh'],
  components: {
    UserFilled,
    Document,
    ChatDotRound,
    View,
    Lightning,
    HomeFilled,
    Plus,
    User,
    InfoFilled,
    PriceTag,
    Coin,
    Monitor,
    Refresh
  },
  props: {
    statistics: { type: Object, default: () => ({}) },
    loading: { type: Boolean, default: false }
  },
  data() {
    return {
      statCards: [
        { key: 'user', field: 'user_count', label: '用户总数', icon: 'UserFilled', tone: 'primary', trend: '+ 今日' },
        { key: 'article', field: 'article_count', label: '文章总数', icon: 'Document', tone: 'violet', trend: '+ 今日' },
        { key: 'comment', field: 'comment_count', label: '评论总数', icon: 'ChatDotRound', tone: 'cyan', trend: '+ 今日' },
        { key: 'view', field: 'view_count', label: '总浏览量', icon: 'View', tone: 'green', trend: '+ 今日' }
      ]
    }
  },
  methods: {
    formatNumber(num) {
      return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',')
    }
  }
}
</script>

<style scoped>
.overview {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* ---------------- 欢迎横幅 ---------------- */
.welcome-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 16px;
  padding: 20px;
  border-radius: var(--campus-radius-lg);
  background: linear-gradient(120deg, var(--campus-primary), var(--campus-primary-light));
  color: #fff;
  box-shadow: var(--campus-shadow);
}

.welcome-title {
  margin: 0 0 4px;
  font-size: 22px;
  font-weight: 700;
  letter-spacing: -0.01em;
}

.welcome-sub {
  margin: 0;
  font-size: 14px;
  opacity: 0.85;
}

.welcome-bar :deep(.el-button--primary) {
  background: rgba(255, 255, 255, 0.16);
  border-color: rgba(255, 255, 255, 0.4);
  color: #fff;
}

.welcome-bar :deep(.el-button--primary:hover) {
  background: rgba(255, 255, 255, 0.28);
  border-color: #fff;
}

/* ---------------- 统计卡片 ---------------- */
.stat-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.stat-card {
  background: var(--campus-surface);
  border: 1px solid var(--campus-border);
  border-radius: var(--campus-radius);
  padding: 18px;
  box-shadow: var(--campus-shadow-sm);
  transition: var(--campus-transition);
}

.stat-card:hover {
  transform: translateY(-3px);
  box-shadow: var(--campus-shadow);
}

.stat-card-top {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 14px;
}

.stat-icon {
  width: 44px;
  height: 44px;
  display: grid;
  place-items: center;
  border-radius: 12px;
}

.tone-primary { background: rgba(79, 110, 247, 0.1); color: #4f6ef7; }
.tone-violet { background: rgba(124, 58, 237, 0.12); color: #7c3aed; }
.tone-cyan { background: rgba(6, 182, 212, 0.12); color: #0891b2; }
.tone-green { background: rgba(34, 197, 94, 0.12); color: #22c55e; }

.stat-change {
  font-size: 12px;
  font-weight: 600;
  padding: 3px 8px;
  border-radius: 999px;
  background: var(--campus-surface-2);
}

.tone-text-primary { color: #4f6ef7; }
.tone-text-violet { color: #7c3aed; }
.tone-text-cyan { color: #0891b2; }
.tone-text-green { color: #22c55e; }
.tone-text-info { color: #64748b; }
.tone-text-success { color: #22c55e; }

.stat-value {
  font-size: 26px;
  font-weight: 800;
  letter-spacing: -0.02em;
  color: var(--campus-text);
  line-height: 1.2;
}

.stat-label {
  margin-top: 4px;
  font-size: 13px;
  color: var(--campus-text-secondary);
}

/* ---------------- 通用卡片 ---------------- */
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

/* ---------------- 快捷操作 ---------------- */
.quick-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
}

.quick-btn {
  width: 100%;
  margin-left: 0;
  padding: 18px 0;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 600;
}

.quick-btn .el-icon {
  margin-right: 6px;
}

/* ---------------- 系统信息 ---------------- */
.sys-info {
  display: flex;
  flex-direction: column;
}

.sys-info-item {
  display: flex;
  align-items: center;
  padding: 13px 4px;
  border-bottom: 1px solid var(--campus-border);
}

.sys-info-item:last-child {
  border-bottom: none;
}

.sys-icon {
  margin-right: 10px;
  font-size: 16px;
}

.sys-label {
  color: var(--campus-text);
  font-weight: 500;
  font-size: 14px;
}

.sys-value {
  margin-left: auto;
  color: var(--campus-text-secondary);
  font-weight: 600;
}

.sys-success {
  color: var(--campus-success);
}

.sys-info-color {
  color: var(--campus-primary);
}

/* ---------------- 响应式 ---------------- */
@media (max-width: 1200px) {
  .stat-grid { grid-template-columns: repeat(2, 1fr); }
  .quick-grid { grid-template-columns: repeat(2, 1fr); }
}

@media (max-width: 575px) {
  .stat-grid { grid-template-columns: 1fr; }
  .quick-grid { grid-template-columns: 1fr; }
  .welcome-title { font-size: 19px; }
  .stat-value { font-size: 22px; }
}
</style>
