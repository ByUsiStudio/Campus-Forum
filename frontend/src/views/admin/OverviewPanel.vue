<template>
  <div class="overview-panel">
    <div class="panel-header">
      <h2 class="panel-title">数据概览</h2>
      <v-btn variant="text" color="primary" @click="$emit('refresh')" :loading="loading">
        <v-icon start size="18">mdi-refresh</v-icon>
        刷新
      </v-btn>
    </div>

    <v-row class="stats-grid">
      <v-col cols="12" sm="6" lg="3">
        <div class="stat-card">
          <div class="stat-icon stat-users">
            <v-icon size="24" color="white">mdi-account-group</v-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ statistics.user_count || 0 }}</div>
            <div class="stat-label">用户总数</div>
          </div>
        </div>
      </v-col>

      <v-col cols="12" sm="6" lg="3">
        <div class="stat-card">
          <div class="stat-icon stat-articles">
            <v-icon size="24" color="white">mdi-file-document</v-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ statistics.article_count || 0 }}</div>
            <div class="stat-label">文章总数</div>
          </div>
        </div>
      </v-col>

      <v-col cols="12" sm="6" lg="3">
        <div class="stat-card">
          <div class="stat-icon stat-comments">
            <v-icon size="24" color="white">mdi-comment-text</v-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ statistics.comment_count || 0 }}</div>
            <div class="stat-label">评论总数</div>
          </div>
        </div>
      </v-col>

      <v-col cols="12" sm="6" lg="3">
        <div class="stat-card">
          <div class="stat-icon stat-views">
            <v-icon size="24" color="white">mdi-eye</v-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ formatNumber(statistics.view_count || 0) }}</div>
            <div class="stat-label">总浏览量</div>
          </div>
        </div>
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="12" lg="8">
        <v-card class="content-card">
          <v-card-title class="card-title">快捷操作</v-card-title>
          <v-card-text>
            <div class="quick-actions">
              <v-btn variant="tonal" color="primary" to="/">
                <v-icon start size="18">mdi-home</v-icon>
                返回首页
              </v-btn>
              <v-btn variant="tonal" color="success">
                <v-icon start size="18">mdi-plus</v-icon>
                创建文章
              </v-btn>
              <v-btn variant="tonal" color="warning">
                <v-icon start size="18">mdi-bell</v-icon>
                发送通知
              </v-btn>
              <v-btn variant="tonal" color="info">
                <v-icon start size="18">mdi-cog</v-icon>
                系统设置
              </v-btn>
            </div>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" lg="4">
        <v-card class="content-card">
          <v-card-title class="card-title">系统信息</v-card-title>
          <v-card-text>
            <div class="info-list">
              <div class="info-item">
                <span class="info-label">系统版本</span>
                <span class="info-value">v1.0.0</span>
              </div>
              <div class="info-item">
                <span class="info-label">数据库状态</span>
                <v-chip size="small" color="success" variant="tonal">正常</v-chip>
              </div>
              <div class="info-item">
                <span class="info-label">服务器状态</span>
                <v-chip size="small" color="success" variant="tonal">在线</v-chip>
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
        minute: '2-digit'
      })
    }

    const formatNumber = (num) => {
      if (num >= 10000) return (num / 10000).toFixed(1) + 'w'
      if (num >= 1000) return (num / 1000).toFixed(1) + 'k'
      return num.toString()
    }

    onMounted(() => {
      updateTime()
      timeInterval = setInterval(updateTime, 60000)
    })

    onUnmounted(() => {
      if (timeInterval) clearInterval(timeInterval)
    })

    return { currentTime, formatNumber }
  }
}
</script>

<style scoped>
.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.panel-title {
  font-size: 20px;
  font-weight: 600;
  color: #1A1A1A;
  margin: 0;
}

.stats-grid {
  margin-bottom: 16px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: #fff;
  border-radius: 12px;
  border: 1px solid #F0F0F0;
  transition: border-color 0.2s;
}

.stat-card:hover {
  border-color: #E0E0E0;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-users { background: #6750A4; }
.stat-articles { background: #EC4899; }
.stat-comments { background: #06B6D4; }
.stat-views { background: #10B981; }

.stat-info {
  min-width: 0;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #1A1A1A;
  line-height: 1.2;
}

.stat-label {
  font-size: 13px;
  color: #888;
  margin-top: 2px;
}

.content-card {
  border-radius: 12px;
  border: 1px solid #F0F0F0;
}

.card-title {
  font-size: 15px;
  font-weight: 600;
  color: #1A1A1A;
  padding-bottom: 12px;
}

.quick-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.info-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.info-label {
  font-size: 14px;
  color: #666;
}

.info-value {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}
</style>