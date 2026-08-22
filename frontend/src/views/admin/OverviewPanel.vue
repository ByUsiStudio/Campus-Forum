<template>
  <div>
    <!-- 统计卡片 -->
    <el-row :gutter="16" class="mb-4">
      <el-col :xs="12" :sm="12" :md="6">
        <el-card class="stat-card stat-primary" shadow="never">
          <div class="stat-card-inner d-flex align-center mb-2">
            <el-icon :size="32"><UserFilled /></el-icon>
          </div>
          <div class="stat-value text-white">{{ statistics.user_count || 0 }}</div>
          <div class="stat-label text-white">用户总数</div>
        </el-card>
      </el-col>

      <el-col :xs="12" :sm="12" :md="6">
        <el-card class="stat-card stat-pink" shadow="never">
          <div class="stat-card-inner d-flex align-center mb-2">
            <el-icon :size="32"><Document /></el-icon>
          </div>
          <div class="stat-value text-white">{{ statistics.article_count || 0 }}</div>
          <div class="stat-label text-white">文章总数</div>
        </el-card>
      </el-col>

      <el-col :xs="12" :sm="12" :md="6">
        <el-card class="stat-card stat-cyan" shadow="never">
          <div class="stat-card-inner d-flex align-center mb-2">
            <el-icon :size="32"><ChatDotRound /></el-icon>
          </div>
          <div class="stat-value text-white">{{ statistics.comment_count || 0 }}</div>
          <div class="stat-label text-white">评论总数</div>
        </el-card>
      </el-col>

      <el-col :xs="12" :sm="12" :md="6">
        <el-card class="stat-card stat-green" shadow="never">
          <div class="stat-card-inner d-flex align-center mb-2">
            <el-icon :size="32"><View /></el-icon>
          </div>
          <div class="stat-value text-white">{{ formatNumber(statistics.view_count || 0) }}</div>
          <div class="stat-label text-white">总浏览量</div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 快捷操作 -->
    <el-card class="page-container mb-4" shadow="never">
      <template #header>
        <div class="card-header">
          <el-icon :size="20"><Lightning /></el-icon>
          <span>快捷操作</span>
        </div>
      </template>
      <el-row :gutter="16">
        <el-col :xs="12" :sm="6" :md="6">
          <el-button type="primary" plain class="quick-btn" @click="$router.push('/')">
            <el-icon><HomeFilled /></el-icon>返回首页
          </el-button>
        </el-col>
        <el-col :xs="12" :sm="6" :md="6">
          <el-button type="success" plain class="quick-btn" @click="$router.push('/create')">
            <el-icon><Plus /></el-icon>创建文章
          </el-button>
        </el-col>
        <el-col :xs="12" :sm="6" :md="6">
          <el-button type="info" plain class="quick-btn" @click="$router.push({ name: 'AdminUsers' })">
            <el-icon><User /></el-icon>用户管理
          </el-button>
        </el-col>
        <el-col :xs="12" :sm="6" :md="6">
          <el-button type="warning" plain class="quick-btn" @click="$router.push({ name: 'AdminArticles' })">
            <el-icon><Document /></el-icon>文章管理
          </el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- 系统信息 -->
    <el-card class="page-container" shadow="never">
      <template #header>
        <div class="card-header">
          <el-icon :size="20"><InfoFilled /></el-icon>
          <span>系统信息</span>
        </div>
      </template>
      <div class="sys-info">
        <div class="sys-info-item">
          <el-icon><PriceTag /></el-icon>
          <span class="sys-label">系统版本</span>
          <span class="sys-value">v1.3.14</span>
        </div>
        <div class="sys-info-item">
          <el-icon color="#67C23A"><Coin /></el-icon>
          <span class="sys-label">数据库状态</span>
          <span class="sys-value sys-success">正常运行</span>
        </div>
        <div class="sys-info-item">
          <el-icon color="#409EFF"><Monitor /></el-icon>
          <span class="sys-label">服务器状态</span>
          <span class="sys-value sys-info-color">运行中</span>
        </div>
      </div>
    </el-card>
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
  Monitor
} from '@element-plus/icons-vue'

export default {
  name: 'OverviewPanel',
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
    Monitor
  },
  props: {
    statistics: { type: Object, default: () => ({}) },
    loading: { type: Boolean, default: false }
  },
  methods: {
    formatNumber(num) {
      return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',')
    }
  }
}
</script>

<style scoped>
.mb-2 {
  margin-bottom: 8px;
}
.mb-4 {
  margin-bottom: 16px;
}

.stat-card {
  margin-bottom: 16px;
}

.stat-card-inner {
  display: flex;
  align-items: center;
}

.stat-primary :deep(.el-card__body) {
  background: #409EFF;
}
.stat-pink :deep(.el-card__body) {
  background: #E91E63;
}
.stat-cyan :deep(.el-card__body) {
  background: #00ACC1;
}
.stat-green :deep(.el-card__body) {
  background: #67C23A;
}

.stat-value {
  font-size: 22px;
  font-weight: 700;
  line-height: 1.4;
}

.stat-label {
  font-size: 13px;
}

.text-white {
  color: #fff;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.quick-btn {
  width: 100%;
  margin-left: 0;
  margin-bottom: 12px;
}

.quick-btn .el-icon {
  margin-right: 4px;
}

.sys-info {
  display: flex;
  flex-direction: column;
}

.sys-info-item {
  display: flex;
  align-items: center;
  padding: 12px 4px;
  border-bottom: 1px solid #f0f2f5;
}

.sys-info-item:last-child {
  border-bottom: none;
}

.sys-label {
  margin-left: 8px;
  color: #303133;
}

.sys-value {
  margin-left: auto;
  color: #606266;
}

.sys-success {
  color: #67C23A;
}

.sys-info-color {
  color: #409EFF;
}
</style>
