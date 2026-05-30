<template>
  <div class="overview-panel">
    <div class="panel-header">
      <div class="header-left">
        <h2 class="panel-title">数据概览</h2>
        <p class="panel-subtitle">查看论坛核心数据统计</p>
      </div>
      <v-btn variant="outlined" color="primary" @click="$emit('refresh')" :loading="loading">
        <v-icon start>mdi-refresh</v-icon>
        刷新
      </v-btn>
    </div>

    <v-row class="stats-grid">
      <v-col cols="12" sm="6" lg="3">
        <div class="stat-card stat-users">
          <div class="stat-icon-wrapper">
            <v-icon size="28" color="white">mdi-account-group</v-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ statistics.user_count || 0 }}</div>
            <div class="stat-label">用户总数</div>
          </div>
        </div>
      </v-col>

      <v-col cols="12" sm="6" lg="3">
        <div class="stat-card stat-articles">
          <div class="stat-icon-wrapper">
            <v-icon size="28" color="white">mdi-file-document</v-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ statistics.article_count || 0 }}</div>
            <div class="stat-label">文章总数</div>
          </div>
        </div>
      </v-col>

      <v-col cols="12" sm="6" lg="3">
        <div class="stat-card stat-comments">
          <div class="stat-icon-wrapper">
            <v-icon size="28" color="white">mdi-comment-text</v-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ statistics.comment_count || 0 }}</div>
            <div class="stat-label">评论总数</div>
          </div>
        </div>
      </v-col>

      <v-col cols="12" sm="6" lg="3">
        <div class="stat-card stat-views">
          <div class="stat-icon-wrapper">
            <v-icon size="28" color="white">mdi-eye</v-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ formatNumber(statistics.view_count || 0) }}</div>
            <div class="stat-label">总浏览量</div>
          </div>
        </div>
      </v-col>
    </v-row>

    <v-row class="mt-6">
      <v-col cols="12" lg="8">
        <v-card class="quick-actions-card">
          <v-card-title class="card-title">
            <v-icon class="title-icon">mdi-lightning-bolt</v-icon>
            快捷操作
          </v-card-title>
          <v-card-text>
            <div class="quick-actions-grid">
              <v-btn color="primary" variant="tonal" size="large" class="action-btn" to="/">
                <v-icon start>mdi-home</v-icon>
                返回首页
              </v-btn>
              <v-btn color="success" variant="tonal" size="large" class="action-btn">
                <v-icon start>mdi-plus</v-icon>
                创建文章
              </v-btn>
              <v-btn color="warning" variant="tonal" size="large" class="action-btn">
                <v-icon start>mdi-bell</v-icon>
                发送通知
              </v-btn>
              <v-btn color="info" variant="tonal" size="large" class="action-btn">
                <v-icon start>mdi-cog</v-icon>
                系统设置
              </v-btn>
            </div>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" lg="4">
        <v-card class="system-info-card">
          <v-card-title class="card-title">
            <v-icon class="title-icon">mdi-information</v-icon>
            系统信息
          </v-card-title>
          <v-card-text>
            <div class="info-list">
              <div class="info-item">
                <span class="info-label">系统版本</span>
                <span class="info-value">v1.0.0</span>
              </div>
              <div class="info-item">
                <span class="info-label">数据库状态</span>
                <v-chip size="small" color="success">正常</v-chip>
              </div>
              <div class="info-item">
                <span class="info-label">服务器状态</span>
                <v-chip size="small" color="success">在线</v-chip>
              </div>
              <div class="info-item">
                <span class="info-label">当前时间</span>
                <span class="info-value">{{ currentTime }}</span>
              </div>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted } from 'vue'

export default {
  name: 'OverviewPanel',
  props: {
    statistics: {
      type: Object,
      default: () => ({
        user_count: 0,
        article_count: 0,
        comment_count: 0,
        view_count: 0
      })
    },
    loading: {
      type: Boolean,
      default: false
    }
  },
  emits: ['refresh'],
  setup() {
    const currentTime = ref('')
    let timeInterval = null

    const updateTime = () => {
      const now = new Date()
      currentTime.value = now.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
      })
    }

    const formatNumber = (num) => {
      if (num >= 10000) {
        return (num / 10000).toFixed(1) + 'w'
      }
      if (num >= 1000) {
        return (num / 1000).toFixed(1) + 'k'
      }
      return num.toString()
    }

    onMounted(() => {
      updateTime()
      timeInterval = setInterval(updateTime, 1000)
    })

    onUnmounted(() => {
      if (timeInterval) {
        clearInterval(timeInterval)
      }
    })

    return {
      currentTime,
      formatNumber
    }
  }
}
</script>

<style scoped>
.overview-panel {
  animation: panelFadeIn 0.4s ease;
}

@keyframes panelFadeIn {
  from {
    opacity: 0;
    transform: translateY(15px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.stats-grid {
  margin-bottom: 16px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 28px;
  border-radius: 20px;
  background: linear-gradient(145deg, #ffffff 0%, #f8fafc 100%);
  box-shadow: 
    0 4px 6px -1px rgba(0, 0, 0, 0.05),
    0 2px 4px -2px rgba(0, 0, 0, 0.05),
    inset 0 1px 0 rgba(255, 255, 255, 0.8);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  height: 100%;
  position: relative;
  overflow: hidden;
}

.stat-card::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -25%;
  width: 100px;
  height: 100px;
  background: rgba(103, 80, 164, 0.03);
  border-radius: 50%;
}

.stat-card:hover {
  transform: translateY(-6px);
  box-shadow: 
    0 20px 40px -10px rgba(0, 0, 0, 0.1),
    0 8px 16px -4px rgba(0, 0, 0, 0.06);
}

.stat-icon-wrapper {
  width: 64px;
  height: 64px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: transform 0.3s ease;
  position: relative;
  z-index: 1;
}

.stat-card:hover .stat-icon-wrapper {
  transform: scale(1.08);
}

.stat-users .stat-icon-wrapper {
  background: linear-gradient(135deg, #6366F1 0%, #8B5CF6 100%);
  box-shadow: 0 8px 24px -8px rgba(99, 102, 241, 0.4);
}

.stat-articles .stat-icon-wrapper {
  background: linear-gradient(135deg, #EC4899 0%, #F43F5E 100%);
  box-shadow: 0 8px 24px -8px rgba(236, 72, 153, 0.4);
}

.stat-comments .stat-icon-wrapper {
  background: linear-gradient(135deg, #06B6D4 0%, #3B82F6 100%);
  box-shadow: 0 8px 24px -8px rgba(6, 182, 212, 0.4);
}

.stat-views .stat-icon-wrapper {
  background: linear-gradient(135deg, #10B981 0%, #059669 100%);
  box-shadow: 0 8px 24px -8px rgba(16, 185, 129, 0.4);
}

.stat-content {
  flex: 1;
  min-width: 0;
  position: relative;
  z-index: 1;
}

.stat-value {
  font-size: 2.5rem;
  font-weight: 800;
  color: #1C1B1F;
  line-height: 1.15;
  letter-spacing: -1px;
}

.stat-label {
  font-size: 14px;
  color: #625B71;
  margin-top: 6px;
  font-weight: 500;
}

.quick-actions-card,
.system-info-card {
  border-radius: 20px;
  height: 100%;
  background: #fff;
  box-shadow: 0 4px 20px -4px rgba(0, 0, 0, 0.06);
  overflow: hidden;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 14px;
  font-size: 16px;
  font-weight: 600;
  padding: 22px 24px !important;
  border-bottom: 1px solid #F2F0F4;
  background: linear-gradient(90deg, #F8F7FF 0%, #ffffff 100%);
}

.title-icon {
  width: 40px;
  height: 40px;
  padding: 8px;
  border-radius: 12px;
  background: linear-gradient(135deg, rgba(103, 80, 164, 0.1) 0%, rgba(123, 58, 237, 0.1) 100%);
  color: #6750A4 !important;
  transition: transform 0.2s ease;
}

.card-title:hover .title-icon {
  transform: scale(1.1);
}

.quick-actions-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 14px;
  padding: 8px;
}

.action-btn {
  justify-content: flex-start;
  height: 56px;
  border-radius: 14px;
  font-weight: 500;
  font-size: 14px;
  transition: all 0.25s ease;
}

.action-btn:hover {
  transform: translateY(-2px);
}

.info-list {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  transition: background 0.2s ease;
}

.info-item:hover {
  background: #F8F7FF;
}

.info-item:not(:last-child) {
  border-bottom: 1px solid #F2F0F4;
}

.info-label {
  font-size: 14px;
  color: #625B71;
  font-weight: 500;
}

.info-value {
  font-size: 14px;
  font-weight: 600;
  color: #1C1B1F;
}

@media (max-width: 600px) {
  .stats-grid {
    margin-bottom: 0;
  }

  .stat-card {
    padding: 22px;
  }

  .stat-value {
    font-size: 2rem;
  }

  .quick-actions-grid {
    grid-template-columns: 1fr;
  }
  .stat-value {
    font-size: 1.5rem;
  }

  .quick-actions-grid {
    grid-template-columns: 1fr;
  }
}
</style>
