<template>
  <v-menu
    v-model="menuOpen"
    :close-on-content-click="false"
    location="bottom end"
    max-width="400"
  >
    <template v-slot:activator="{ props }">
      <v-btn
        v-bind="props"
        icon
        variant="text"
        class="notification-btn"
      >
        <v-badge
          :content="unreadCount"
          :model-value="unreadCount > 0"
          color="error"
          overlap
        >
          <v-icon>mdi-bell</v-icon>
        </v-badge>
      </v-btn>
    </template>

    <v-card>
      <v-card-title class="d-flex justify-space-between align-center pa-3">
        <span class="text-subtitle-1">通知</span>
        <v-btn
          v-if="unreadCount > 0"
          variant="text"
          size="small"
          color="primary"
          @click="markAllRead"
        >
          全部已读
        </v-btn>
      </v-card-title>

      <v-divider></v-divider>

      <v-list v-if="notifications.length > 0" lines="three" max-height="400" class="overflow-auto">
        <v-list-item
          v-for="notification in notifications"
          :key="notification.id"
          :class="{ 'bg-grey-lighten-4': !notification.is_read }"
          @click="handleClick(notification)"
        >
          <template v-slot:prepend>
            <v-avatar :color="getTypeColor(notification.type)" size="32">
              <v-icon size="small" color="white">{{ getTypeIcon(notification.type) }}</v-icon>
            </v-avatar>
          </template>

          <v-list-item-title class="font-weight-medium">
            {{ notification.title }}
          </v-list-item-title>
          <v-list-item-subtitle class="text-caption">
            {{ notification.content }}
          </v-list-item-subtitle>
          <v-list-item-subtitle class="text-caption text-medium-emphasis mt-1">
            {{ formatDate(notification.created_at) }}
          </v-list-item-subtitle>
        </v-list-item>
      </v-list>

      <v-card-text v-else class="text-center py-8 text-medium-emphasis">
        暂无通知
      </v-card-text>

      <v-divider v-if="notifications.length > 0"></v-divider>

      <v-card-actions v-if="notifications.length > 0">
        <v-spacer></v-spacer>
        <v-btn variant="text" color="primary" @click="goToNotifications">
          查看全部
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-menu>
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

    const getTypeColor = (type) => {
      const colors = {
        system: 'primary',
        activity: 'success',
        update: 'info',
        warning: 'warning'
      }
      return colors[type] || 'default'
    }

    const getTypeIcon = (type) => {
      const icons = {
        system: 'mdi-information',
        activity: 'mdi-calendar-star',
        update: 'mdi-update',
        warning: 'mdi-alert'
      }
      return icons[type] || 'mdi-bell'
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
      if (pollInterval) {
        clearInterval(pollInterval)
      }
    })

    return {
      menuOpen,
      notifications,
      unreadCount,
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
</style>
