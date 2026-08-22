<template>
  <el-popover
    v-model:visible="menuOpen"
    :width="popoverWidth"
    placement="bottom-end"
    trigger="click"
    :show-arrow="false"
    popper-class="notification-popover"
  >
    <template #reference>
      <el-badge
        :value="unreadCount"
        :hidden="unreadCount <= 0"
        :max="99"
        class="notification-btn"
      >
        <el-button
          text
          circle
          class="bell-button"
          aria-label="通知"
        >
          <el-icon :size="20"><Bell /></el-icon>
        </el-button>
      </el-badge>
    </template>

    <div class="notification-card">
      <div class="notification-header">
        <span class="notification-title">通知</span>
        <el-button
          v-if="unreadCount > 0"
          text
          type="primary"
          size="small"
          class="mark-all-btn"
          @click="markAllRead"
        >
          <el-icon class="mark-icon"><Select /></el-icon>
          全部已读
        </el-button>
      </div>

      <el-scrollbar max-height="400px" class="notification-scroll">
        <ul v-if="notifications.length > 0" class="notification-list">
          <li
            v-for="notification in notifications"
            :key="notification.id"
            class="notification-item cursor-pointer"
            :class="{ 'notification-item-unread': !notification.is_read }"
            @click="handleClick(notification)"
          >
            <el-avatar
              :size="36"
              :style="{ backgroundColor: getTypeColor(notification.type) }"
              class="notification-avatar"
            >
              <el-icon :size="17" color="white">
                <component :is="getTypeIcon(notification.type)" />
              </el-icon>
            </el-avatar>
            <div class="notification-body">
              <div class="notification-item-top">
                <span class="notification-item-title">{{ notification.title }}</span>
                <i v-if="!notification.is_read" class="unread-dot" aria-hidden="true"></i>
              </div>
              <div class="notification-item-content">{{ notification.content }}</div>
              <div class="notification-item-time">{{ formatDate(notification.created_at) }}</div>
            </div>
          </li>
        </ul>
        <div v-else class="notification-empty">
          <el-icon :size="40" class="empty-icon"><Bell /></el-icon>
          <div class="empty-text">暂无通知</div>
        </div>
      </el-scrollbar>

      <div v-if="notifications.length > 0" class="notification-footer">
        <el-button text type="primary" size="small" class="view-all-btn" @click="goToNotifications">
          查看全部通知
          <el-icon class="arrow-icon"><ArrowRight /></el-icon>
        </el-button>
      </div>
    </div>
  </el-popover>
</template>

<script>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { Bell, InfoFilled, Calendar, Refresh, Warning, Select, ArrowRight } from '@element-plus/icons-vue'
import api from '../api'
import { success } from '../utils/message'

export default {
  name: 'NotificationBell',
  setup() {
    const router = useRouter()
    const menuOpen = ref(false)
    const notifications = ref([])
    const unreadCount = ref(0)
    const popoverWidth = ref(400)
    let pollInterval = null

    const updatePopoverWidth = () => {
      popoverWidth.value = window.innerWidth < 440 ? Math.min(400, window.innerWidth - 24) : 400
    }

    const loadNotifications = async () => {
      try {
        const response = await api.get('/notifications')
        notifications.value = (response.data.notifications || []).slice(0, 5)
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

    const handleClick = async (notification) => {
      if (!notification.is_read) {
        try {
          await api.post(`/notifications/${notification.id}/read`)
          notification.is_read = true
          unreadCount.value = Math.max(0, unreadCount.value - 1)
        } catch (error) {
          console.error('标记已读失败', error)
        }
      }
    }

    const markAllRead = async () => {
      try {
        await api.post('/notifications/read-all')
        notifications.value.forEach(n => n.is_read = true)
        unreadCount.value = 0
        success('已全部标记为已读')
      } catch (error) {
        console.error('标记全部已读失败', error)
      }
    }

    const goToNotifications = () => {
      menuOpen.value = false
      router.push('/notifications')
    }

    const getTypeColor = (type) => {
      const colors = {
        system: 'var(--campus-primary)',
        activity: '#22c55e',
        update: '#409eff',
        warning: '#f59e0b'
      }
      return colors[type] || 'var(--campus-text-secondary)'
    }

    const getTypeIcon = (type) => {
      const icons = {
        system: InfoFilled,
        activity: Calendar,
        update: Refresh,
        warning: Warning
      }
      return icons[type] || Bell
    }

    const formatDate = (date) => {
      const d = new Date(date)
      const now = new Date()
      const diff = now - d
      const minutes = Math.floor(diff / 60000)
      const hours = Math.floor(diff / 3600000)
      const days = Math.floor(diff / 86400000)

      if (minutes < 1) return '刚刚'
      if (minutes < 60) return `${minutes}分钟前`
      if (hours < 24) return `${hours}小时前`
      if (days < 7) return `${days}天前`
      return d.toLocaleDateString('zh-CN')
    }

    onMounted(() => {
      updatePopoverWidth()
      window.addEventListener('resize', updatePopoverWidth)
      const token = localStorage.getItem('token')
      if (token) {
        loadNotifications()
        loadUnreadCount()
        pollInterval = setInterval(() => {
          loadUnreadCount()
        }, 60000)
      }
    })

    onUnmounted(() => {
      window.removeEventListener('resize', updatePopoverWidth)
      if (pollInterval) {
        clearInterval(pollInterval)
      }
    })

    return {
      menuOpen,
      notifications,
      unreadCount,
      popoverWidth,
      handleClick,
      markAllRead,
      goToNotifications,
      getTypeColor,
      getTypeIcon,
      formatDate
    }
  }
}
</script>

<style scoped>
.notification-btn {
  margin-left: 8px;
}

.bell-button {
  color: var(--campus-text);
}

.notification-card {
  padding: 8px;
}

.notification-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 6px 8px 10px;
}

.notification-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--campus-text);
  letter-spacing: -0.01em;
}

.mark-all-btn {
  gap: 4px;
}

.mark-icon {
  font-size: 14px;
}

.notification-scroll {
  border-radius: var(--campus-radius-sm);
}

.notification-list {
  list-style: none;
  margin: 0;
  padding: 0;
}

.notification-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px;
  border-radius: 10px;
  transition: background-color 0.2s ease;
}

.notification-item-unread {
  background: var(--campus-primary-soft);
}

.notification-item:hover {
  background: var(--campus-surface-2);
}

.notification-avatar {
  flex-shrink: 0;
  margin-top: 2px;
  border-radius: 10px;
}

.notification-body {
  flex: 1;
  min-width: 0;
}

.notification-item-top {
  display: flex;
  align-items: center;
  gap: 6px;
}

.notification-item-title {
  font-weight: 600;
  font-size: 14px;
  color: var(--campus-text);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.unread-dot {
  flex-shrink: 0;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--campus-primary);
  box-shadow: 0 0 0 3px rgba(79, 110, 247, 0.15);
}

.notification-item-content {
  font-size: 13px;
  color: var(--campus-text-secondary);
  margin-top: 2px;
  word-break: break-word;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.notification-item-time {
  font-size: 12px;
  color: var(--campus-text-muted);
  margin-top: 4px;
}

.notification-empty {
  padding: 36px 0;
  text-align: center;
}

.empty-icon {
  color: var(--campus-text-muted);
}

.empty-text {
  margin-top: 10px;
  color: var(--campus-text-secondary);
  font-size: 14px;
}

.notification-footer {
  display: flex;
  justify-content: center;
  padding: 10px 8px 4px;
  border-top: 1px solid var(--campus-border);
  margin-top: 4px;
}

.view-all-btn {
  gap: 4px;
  width: 100%;
  justify-content: center;
}

.arrow-icon {
  font-size: 14px;
}
</style>

<style>
/* 控制弹层自身圆角与阴影 */
.notification-popover {
  border-radius: var(--campus-radius) !important;
  box-shadow: var(--campus-shadow-lg) !important;
  border: 1px solid var(--campus-border) !important;
  padding: 8px !important;
  overflow: hidden;
}
</style>
