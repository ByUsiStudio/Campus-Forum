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
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 16px;
}

.header-left {
  flex: 1;
}

.panel-title {
  font-size: 1.5rem;
  font-weight: 700;
  color: #1a1a2e;
  margin: 0 0 4px 0;
}

.panel-subtitle {
  font-size: 0.9rem;
  color: #6b7280;
  margin: 0;
}

.stats-grid {
  margin-bottom: 12px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 24px;
  border-radius: 16px;
  background: white;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
  transition: transform 0.2s, box-shadow 0.2s;
  height: 100%;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.1);
}

.stat-icon-wrapper {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-users .stat-icon-wrapper {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stat-articles .stat-icon-wrapper {
  background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
}

.stat-comments .stat-icon-wrapper {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.stat-views .stat-icon-wrapper {
  background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);
}

.stat-content {
  flex: 1;
  min-width: 0;
}

.stat-value {
  font-size: 2rem;
  font-weight: 700;
  color: #1a1a2e;
  line-height: 1.2;
}

.stat-label {
  font-size: 0.9rem;
  color: #6b7280;
  margin-top: 4px;
}

.quick-actions-card,
.system-info-card {
  border-radius: 16px;
  height: 100%;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 1.1rem;
  font-weight: 700;
  padding: 20px 24px !important;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.title-icon {
  width: 36px;
  height: 36px;
  padding: 6px;
  border-radius: 8px;
  background: rgba(103, 80, 164, 0.1);
  color: #6750A4 !important;
}

.quick-actions-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.action-btn {
  justify-content: flex-start;
}

.info-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.info-label {
  font-size: 0.9rem;
  color: #6b7280;
}

.info-value {
  font-size: 0.9rem;
  font-weight: 600;
  color: #1a1a2e;
}

@media (max-width: 600px) {
  .stats-grid {
    margin-bottom: 0;
  }

  .stat-card {
    padding: 20px;
  }

  .stat-value {
    font-size: 1.5rem;
  }

  .quick-actions-grid {
    grid-template-columns: 1fr;
  }
}
</style>
