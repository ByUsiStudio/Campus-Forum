<template>
  <div class="notification-wrapper" @click="toggleMenu">
    <button class="notification-btn">
      <span class="badge" v-if="unreadCount > 0">{{ unreadCount }}</span>
      <i class="fa-solid fa-bell"></i>
    </button>

    <div v-if="menuOpen" class="notification-dropdown">
      <div class="dropdown-header">
        <span class="dropdown-title">通知</span>
        <button v-if="unreadCount > 0" class="mark-all-btn" @click.stop="markAllRead">
          全部已读
        </button>
      </div>

      <div class="dropdown-divider"></div>

      <div v-if="notifications.length > 0" class="notification-list">
        <div
          v-for="notification in notifications"
          :key="notification.id"
          class="notification-item"
          :class="{ 'is-unread': !notification.is_read }"
          @click.stop="handleClick(notification)"
        >
          <div class="notification-icon" :style="{ background: getTypeColor(notification.type) + '20', color: getTypeColor(notification.type) }">
            <i :class="getTypeIcon(notification.type)"></i>
          </div>
          <div class="notification-content">
            <div class="notification-title">{{ notification.title }}</div>
            <div class="notification-body">{{ notification.content }}</div>
            <div class="notification-time">{{ formatDate(notification.created_at) }}</div>
          </div>
        </div>
      </div>

      <div v-else class="empty-state">
        <i class="fa-solid fa-inbox"></i>
        <div>暂无通知</div>
      </div>

      <div class="dropdown-divider" v-if="notifications.length > 0"></div>

      <div v-if="notifications.length > 0" class="dropdown-footer">
        <button class="view-all-btn" @click.stop="goToNotifications">查看全部</button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '../api'

export default {
  name: 'NotificationBell',
  setup() {
    const router = useRouter()
    const menuOpen = ref(false)
    const notifications = ref([])
    const unreadCount = ref(0)
    let pollInterval = null

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
      menuOpen.value = false
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

    const goToNotifications = () => {
      menuOpen.value = false
      router.push('/notifications')
    }

    const toggleMenu = () => {
      menuOpen.value = !menuOpen.value
      if (menuOpen.value) {
        loadNotifications()
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
        system: 'fa-solid fa-info-circle',
        activity: 'fa-solid fa-calendar-star',
        update: 'fa-solid fa-refresh',
        warning: 'fa-solid fa-exclamation-triangle'
      }
      return icons[type] || 'fa-solid fa-bell'
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

    const handleClickOutside = (event) => {
      const wrapper = document.querySelector('.notification-wrapper')
      if (wrapper && !wrapper.contains(event.target)) {
        menuOpen.value = false
      }
    }

    onMounted(() => {
      const token = localStorage.getItem('token')
      if (token) {
        loadNotifications()
        loadUnreadCount()
        pollInterval = setInterval(() => {
          loadUnreadCount()
        }, 60000)
      }
      document.addEventListener('click', handleClickOutside)
    })

    onUnmounted(() => {
      if (pollInterval) {
        clearInterval(pollInterval)
      }
      document.removeEventListener('click', handleClickOutside)
    })

    return {
      menuOpen,
      notifications,
      unreadCount,
      handleClick,
      markAllRead,
      goToNotifications,
      toggleMenu,
      getTypeColor,
      getTypeIcon,
      formatDate
    }
  }
}
</script>

<style scoped>
.notification-wrapper {
  position: relative;
}

.notification-btn {
  position: relative;
  background: none;
  border: none;
  font-size: 18px;
  color: #666;
  cursor: pointer;
  padding: 8px;
  border-radius: 6px;
  transition: all 0.2s ease;

  &:hover {
    background: rgba(30, 159, 255, 0.1);
    color: #1E9FFF;
  }
}

.badge {
  position: absolute;
  top: 2px;
  right: 2px;
  background: #FF5722;
  color: #fff;
  font-size: 10px;
  padding: 1px 5px;
  border-radius: 10px;
  min-width: 16px;
  text-align: center;
}

.notification-dropdown {
  position: absolute;
  top: 100%;
  right: 0;
  width: 360px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  z-index: 1000;
  overflow: hidden;
}

.dropdown-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid #f0f0f0;
}

.dropdown-title {
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.mark-all-btn {
  background: none;
  border: none;
  color: #1E9FFF;
  font-size: 12px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;

  &:hover {
    background: rgba(30, 159, 255, 0.1);
  }
}

.dropdown-divider {
  height: 1px;
  background: #f0f0f0;
}

.notification-list {
  max-height: 400px;
  overflow-y: auto;
}

.notification-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px 16px;
  cursor: pointer;
  transition: background 0.2s ease;

  &:hover {
    background: #f8f9fa;
  }

  &.is-unread {
    background: rgba(30, 159, 255, 0.04);
  }
}

.notification-icon {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  font-size: 14px;
}

.notification-content {
  flex: 1;
  min-width: 0;
}

.notification-title {
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin-bottom: 4px;
}

.notification-body {
  font-size: 13px;
  color: #666;
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.notification-time {
  font-size: 12px;
  color: #999;
}

.empty-state {
  padding: 32px 16px;
  text-align: center;
  color: #999;

  i {
    font-size: 32px;
    margin-bottom: 8px;
    display: block;
  }

  div {
    font-size: 14px;
  }
}

.dropdown-footer {
  padding: 12px 16px;
  border-top: 1px solid #f0f0f0;
  text-align: right;
}

.view-all-btn {
  background: none;
  border: none;
  color: #1E9FFF;
  font-size: 13px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;

  &:hover {
    background: rgba(30, 159, 255, 0.1);
  }
}
</style>