<template>
  <div class="notifications-page">
    <div class="notifications-container">
      <div class="notifications-card">
        <div class="card-header">
          <div class="header-left">
            <button class="back-btn" @click="router.back()">
              <i class="fa-solid fa-arrow-left"></i>
            </button>
            <h1 class="page-title">我的通知</h1>
          </div>
          <button 
            v-if="unreadCount > 0"
            class="layui-btn layui-btn-normal"
            @click="markAllRead"
          >
            <i class="fa-solid fa-check-all mr-2"></i>
            全部标记已读
          </button>
        </div>

        <div class="divider"></div>

        <div v-if="notifications.length > 0" class="notifications-list">
          <div 
            v-for="notification in notifications"
            :key="notification.id"
            :class="['notification-item', { unread: !notification.is_read }]"
          >
            <div class="notification-avatar" :style="{ background: getTypeColor(notification.type) }">
              <i :class="getTypeIcon(notification.type)" class="text-white"></i>
            </div>
            
            <div class="notification-content">
              <div class="notification-title">
                <span :class="['type-tag', getTypeColor(notification.type)]">
                  {{ getTypeText(notification.type) }}
                </span>
                {{ notification.title }}
              </div>
              <div class="notification-desc">{{ notification.content }}</div>
              <div class="notification-time">
                <i class="fa-solid fa-clock mr-1"></i>
                {{ formatDate(notification.created_at) }}
              </div>
            </div>

            <button 
              v-if="!notification.is_read"
              class="mark-read-btn"
              @click="markRead(notification)"
            >
              标记已读
            </button>
          </div>
        </div>

        <div v-else class="empty-state">
          <i class="fa-regular fa-bell-slash"></i>
          <div class="empty-title">暂无通知</div>
          <div class="empty-desc">暂无新的通知消息</div>
        </div>

        <div v-if="totalPages > 1" class="pagination">
          <button 
            class="page-btn" 
            :disabled="page <= 1"
            @click="page--"
          >
            <i class="fa-solid fa-chevron-left"></i>
          </button>
          <span class="page-info">{{ page }} / {{ totalPages }}</span>
          <button 
            class="page-btn" 
            :disabled="page >= totalPages"
            @click="page++"
          >
            <i class="fa-solid fa-chevron-right"></i>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import api from '../api'

export default {
  name: 'Notifications',
  setup() {
    const router = useRouter()
    const notifications = ref([])
    const unreadCount = ref(0)
    const page = ref(1)
    const totalPages = ref(1)
    const pageSize = 20

    const loadNotifications = async () => {
      try {
        const response = await api.get('/notifications')
        const allNotifications = response.data.notifications || []
        
        const total = allNotifications.length
        totalPages.value = Math.ceil(total / pageSize)
        
        const start = (page.value - 1) * pageSize
        notifications.value = allNotifications.slice(start, start + pageSize)
      } catch (error) {
        console.error('加载通知失败', error)
      }
    }

    const loadUnreadCount = async () => {
      try {
        const response = await api.get('/notifications/unread-count')
        unreadCount.value = response.data.unread_count || 0
      } catch (error) {
        console.error('加载未读数量失败', error)
      }
    }

    const markRead = async (notification) => {
      try {
        await api.post(`/notifications/${notification.id}/read`)
        notification.is_read = true
        unreadCount.value = Math.max(0, unreadCount.value - 1)
      } catch (error) {
        console.error('标记已读失败', error)
      }
    }

    const markAllRead = async () => {
      try {
        await api.post('/notifications/read-all')
        notifications.value.forEach(n => n.is_read = true)
        unreadCount.value = 0
      } catch (error) {
        console.error('标记全部已读失败', error)
      }
    }

    const getTypeColor = (type) => {
      const colors = {
        system: '#1E9FFF',
        activity: '#52C41A',
        update: '#13C2C2',
        warning: '#FAAD14'
      }
      return colors[type] || '#999'
    }

    const getTypeIcon = (type) => {
      const icons = {
        system: 'fa-solid fa-circle-info',
        activity: 'fa-solid fa-calendar-star',
        update: 'fa-solid fa-rotate-right',
        warning: 'fa-solid fa-triangle-exclamation'
      }
      return icons[type] || 'fa-solid fa-bell'
    }

    const getTypeText = (type) => {
      const texts = {
        system: '系统通知',
        activity: '活动公告',
        update: '更新通知',
        warning: '警告通知'
      }
      return texts[type] || type
    }

    const formatDate = (date) => {
      return new Date(date).toLocaleString('zh-CN')
    }

    watch(page, () => {
      loadNotifications()
    })

    onMounted(() => {
      loadNotifications()
      loadUnreadCount()
    })

    return {
      notifications,
      unreadCount,
      page,
      totalPages,
      markRead,
      markAllRead,
      getTypeColor,
      getTypeIcon,
      getTypeText,
      formatDate
    }
  }
}
</script>

<style scoped>
.notifications-page {
  padding: 24px 0;
  min-height: 100vh;
  background: #f5f5f5;
}

.notifications-container {
  max-width: 900px;
  margin: 0 auto;
  padding: 0 15px;
}

.notifications-card {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.back-btn {
  background: none;
  border: none;
  font-size: 20px;
  color: #666;
  cursor: pointer;
  padding: 8px;
  border-radius: 6px;
  
  &:hover {
    background: #f5f5f5;
    color: var(--primary);
  }
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  color: #333;
  margin: 0;
}

.divider {
  height: 1px;
  background: #f0f0f0;
  margin: 20px 0;
}

.notifications-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.notification-item {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  padding: 16px;
  background: #f8f9fa;
  border-radius: 8px;
  transition: all 0.2s;
  
  &:hover {
    background: #f0f0f0;
  }
  
  &.unread {
    background: rgba(30, 159, 255, 0.05);
    border-left: 4px solid var(--primary);
  }
}

.notification-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  
  i {
    font-size: 24px;
  }
}

.notification-content {
  flex: 1;
  min-width: 0;
}

.notification-title {
  font-size: 15px;
  font-weight: 600;
  color: #333;
  margin-bottom: 8px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.type-tag {
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 4px;
  color: white;
  
  &.system {
    background: #1E9FFF;
  }
  
  &.activity {
    background: #52C41A;
  }
  
  &.update {
    background: #13C2C2;
  }
  
  &.warning {
    background: #FAAD14;
  }
}

.notification-desc {
  font-size: 14px;
  color: #666;
  margin-bottom: 8px;
}

.notification-time {
  font-size: 12px;
  color: #999;
}

.mark-read-btn {
  padding: 6px 12px;
  font-size: 12px;
  color: var(--primary);
  background: rgba(30, 159, 255, 0.1);
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
  
  &:hover {
    background: rgba(30, 159, 255, 0.2);
  }
}

.empty-state {
  text-align: center;
  padding: 48px 24px;
  color: #999;
  
  i {
    font-size: 64px;
    margin-bottom: 16px;
    color: #e8e8e8;
  }
}

.empty-title {
  font-size: 16px;
  margin-bottom: 8px;
  color: #666;
}

.empty-desc {
  font-size: 14px;
  color: #999;
}

.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px solid #f0f0f0;
}

.page-btn {
  width: 32px;
  height: 32px;
  border: 1px solid #e8e8e8;
  border-radius: 6px;
  background: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  color: #666;
  transition: all 0.2s;
  
  &:hover:not(:disabled) {
    border-color: var(--primary);
    color: var(--primary);
  }
  
  &:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }
}

.page-info {
  font-size: 14px;
  color: #666;
}

.mr-1 {
  margin-right: 4px;
}

.mr-2 {
  margin-right: 8px;
}
</style>
