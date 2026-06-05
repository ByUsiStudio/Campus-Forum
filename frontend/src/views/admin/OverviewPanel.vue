<template>
  <div class="overview-panel">
    <!-- 面板头部 -->
    <div class="panel-header">
      <div class="header-left">
        <h2 class="panel-title">数据概览</h2>
        <v-chip size="small" variant="tonal" color="primary" class="ml-3">
          实时统计
        </v-chip>
      </div>
      <v-btn 
        variant="tonal" 
        color="primary" 
        @click="$emit('refresh')" 
        :loading="loading"
        prepend-icon="mdi-refresh"
      >
        刷新数据
      </v-btn>
    </div>

    <!-- 统计卡片网格 -->
    <v-row class="stats-grid" gutter="16">
      <v-col cols="12" sm="6" lg="3">
        <div class="stat-card stat-users">
          <div class="stat-header">
            <div class="stat-icon">
              <v-icon size="24" color="white">mdi-account-group</v-icon>
            </div>
            <v-chip size="x-small" color="success" variant="tonal">+12%</v-chip>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ statistics.user_count || 0 }}</div>
            <div class="stat-label">用户总数</div>
          </div>
        </div>
      </v-col>

      <v-col cols="12" sm="6" lg="3">
        <div class="stat-card stat-articles">
          <div class="stat-header">
            <div class="stat-icon">
              <v-icon size="24" color="white">mdi-file-document</v-icon>
            </div>
            <v-chip size="x-small" color="success" variant="tonal">+8%</v-chip>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ statistics.article_count || 0 }}</div>
            <div class="stat-label">文章总数</div>
          </div>
        </div>
      </v-col>

      <v-col cols="12" sm="6" lg="3">
        <div class="stat-card stat-comments">
          <div class="stat-header">
            <div class="stat-icon">
              <v-icon size="24" color="white">mdi-comment-text</v-icon>
            </div>
            <v-chip size="x-small" color="success" variant="tonal">+15%</v-chip>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ statistics.comment_count || 0 }}</div>
            <div class="stat-label">评论总数</div>
          </div>
        </div>
      </v-col>

      <v-col cols="12" sm="6" lg="3">
        <div class="stat-card stat-views">
          <div class="stat-header">
            <div class="stat-icon">
              <v-icon size="24" color="white">mdi-eye</v-icon>
            </div>
            <v-chip size="x-small" color="success" variant="tonal">+23%</v-chip>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ formatNumber(statistics.view_count || 0) }}</div>
            <div class="stat-label">总浏览量</div>
          </div>
        </div>
      </v-col>
    </v-row>

    <!-- 内容区域 -->
    <v-row class="content-grid" gutter="16">
      <v-col cols="12" lg="8">
        <v-card class="content-card" elevation="1">
          <v-card-title class="card-title">
            <v-icon start size="20" color="primary">mdi-flash</v-icon>
            快捷操作
          </v-card-title>
          <v-card-text>
            <div class="quick-actions">
              <v-btn 
                variant="tonal" 
                color="primary" 
                to="/"
                prepend-icon="mdi-home"
                class="action-btn"
              >
                返回首页
              </v-btn>
              <v-btn 
                variant="tonal" 
                color="success"
                prepend-icon="mdi-plus"
                class="action-btn"
              >
                创建文章
              </v-btn>
              <v-btn 
                variant="tonal" 
                color="warning"
                prepend-icon="mdi-bell"
                class="action-btn"
              >
                发送通知
              </v-btn>
              <v-btn 
                variant="tonal" 
                color="info"
                prepend-icon="mdi-cog"
                class="action-btn"
              >
                系统设置
              </v-btn>
            </div>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" lg="4">
        <v-card class="content-card" elevation="1">
          <v-card-title class="card-title">
            <v-icon start size="20" color="primary">mdi-information</v-icon>
            系统信息
          </v-card-title>
          <v-card-text>
            <div class="info-list">
              <div class="info-item">
                <span class="info-label">系统版本</span>
                <v-chip size="small" variant="tonal" color="primary">v1.0.0</v-chip>
              </div>
              <div class="info-item">
                <span class="info-label">数据库状态</span>
                <v-chip size="small" color="success" variant="tonal">
                  <v-icon start size="14">mdi-check-circle</v-icon>
                  正常
                </v-chip>
              </div>
              <div class="info-item">
                <span class="info-label">服务器状态</span>
                <v-chip size="small" color="success" variant="tonal">
                  <v-icon start size="14">mdi-server</v-icon>
                  在线
                </v-chip>
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
      if (num >= 10000) return (num / 10000).toFixed(1) + 'w'
      if (num >= 1000) return (num / 1000).toFixed(1) + 'k'
      return num.toString()
    }

    onMounted(() => {
      updateTime()
      timeInterval = setInterval(updateTime, 1000)
    })

    onUnmounted(() => {
      if (timeInterval) clearInterval(timeInterval)
    })

    return { currentTime, formatNumber }
  }
}
</script>

<style scoped>
.overview-panel {
  width: 100%;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  gap: 16px;
  flex-wrap: wrap;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.panel-title {
  font-size: 20px;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0;
}

.stats-grid {
  margin-bottom: 24px;
}

.stat-card {
  padding: 24px;
  background: linear-gradient(135deg, #ffffff 0%, #fafbfc 100%);
  border-radius: 16px;
  border: 1px solid #e8eaed;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  height: 100%;
  display: flex;
  flex-direction: column;
  position: relative;
  overflow: hidden;
}

.stat-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, transparent 0%, currentColor 50%, transparent 100%);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 32px rgba(103, 80, 164, 0.15);
  border-color: transparent;
}

.stat-card:hover::before {
  opacity: 0.6;
}

.stat-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 16px;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-users .stat-icon { 
  background: linear-gradient(135deg, #6750A4 0%, #9C27B0 100%);
  box-shadow: 0 4px 12px rgba(103, 80, 164, 0.3);
}
.stat-articles .stat-icon { 
  background: linear-gradient(135deg, #EC4899 0%, #F472B6 100%);
  box-shadow: 0 4px 12px rgba(236, 72, 153, 0.3);
}
.stat-comments .stat-icon { 
  background: linear-gradient(135deg, #06B6D4 0%, #22D3EE 100%);
  box-shadow: 0 4px 12px rgba(6, 182, 212, 0.3);
}
.stat-views .stat-icon { 
  background: linear-gradient(135deg, #10B981 0%, #34D399 100%);
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.3);
}

.stat-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #1a1a1a;
  line-height: 1.2;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: #666;
  font-weight: 500;
}

.content-grid {
  margin-bottom: 0;
}

.content-card {
  border-radius: 16px;
  border: 1px solid #e8eaed;
  height: 100%;
  background: linear-gradient(135deg, #ffffff 0%, #fafbfc 100%);
  transition: all 0.3s ease;
}

.content-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #1a1a1a;
  padding: 20px 20px 16px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.quick-actions {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: 12px;
  padding: 4px 0;
}

.action-btn {
  justify-content: flex-start;
  padding: 12px 16px;
  font-weight: 500;
}

.info-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 4px 0;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
}

.info-label {
  font-size: 14px;
  color: #666;
  font-weight: 500;
}

.info-value {
  font-size: 14px;
  color: #333;
  font-weight: 600;
}

@media (max-width: 768px) {
  .panel-header {
    flex-direction: column;
    align-items: stretch;
  }
  
  .header-left {
    justify-content: space-between;
  }
  
  .stats-grid {
    margin-bottom: 16px;
  }
  
  .stat-card {
    padding: 20px;
  }
  
  .stat-value {
    font-size: 24px;
  }
  
  .quick-actions {
    grid-template-columns: 1fr;
  }
  
  .content-grid {
    gap: 16px;
  }
}

@media (max-width: 480px) {
  .stat-card {
    padding: 16px;
  }
  
  .stat-header {
    margin-bottom: 12px;
  }
  
  .stat-icon {
    width: 40px;
    height: 40px;
  }
  
  .stat-value {
    font-size: 22px;
  }
  
  .stat-label {
    font-size: 13px;
  }
}
</style>