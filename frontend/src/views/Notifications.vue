<template>
  <v-card class="pa-4">
    <v-card-title class="d-flex justify-space-between align-center">
      <span class="text-h5">我的通知</span>
      <v-btn
        v-if="unreadCount > 0"
        variant="outlined"
        color="primary"
        @click="markAllRead"
      >
        全部标记已读
      </v-btn>
    </v-card-title>

    <v-divider class="mb-4"></v-divider>

    <v-list lines="three" v-if="notifications.length > 0">
      <v-list-item
        v-for="notification in notifications"
        :key="notification.id"
        :class="{ 'bg-blue-lighten-5': !notification.is_read }"
        class="mb-2 rounded"
      >
        <template v-slot:prepend>
          <v-avatar :color="getTypeColor(notification.type)" size="40">
            <v-icon color="white">{{ getTypeIcon(notification.type) }}</v-icon>
          </v-avatar>
        </template>

        <v-list-item-title class="font-weight-bold mb-1">
          <v-chip size="small" :color="getTypeColor(notification.type)" class="mr-2">
            {{ getTypeText(notification.type) }}
          </v-chip>
          {{ notification.title }}
        </v-list-item-title>

        <v-list-item-subtitle class="text-body-2 mt-2">
          {{ notification.content }}
        </v-list-item-subtitle>

        <v-list-item-subtitle class="text-caption text-medium-emphasis mt-2">
          {{ formatDate(notification.created_at) }}
        </v-list-item-subtitle>

        <template v-slot:append>
          <v-btn
            v-if="!notification.is_read"
            variant="text"
            size="small"
            color="primary"
            @click="markRead(notification)"
          >
            标记已读
          </v-btn>
        </template>
      </v-list-item>
    </v-list>

    <v-card-text v-else class="text-center py-12">
      <v-icon size="64" color="grey-lighten-1">mdi-bell-off</v-icon>
      <div class="text-h6 text-medium-emphasis mt-4">暂无通知</div>
    </v-card-text>

    <v-card-actions v-if="notifications.length > 0" class="justify-center">
      <v-pagination
        v-model="page"
        :length="totalPages"
        :total-visible="5"
      ></v-pagination>
    </v-card-actions>
  </v-card>
</template>

<script>
import { ref, onMounted, watch } from 'vue'
import api from '../api'

export default {
  name: 'Notifications',
  setup() {
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
