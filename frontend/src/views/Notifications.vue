<script setup>
import { ref, inject, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { notificationApi } from '../api'

const router = useRouter()
const user = inject('user')

const notifications = ref([])
const isLoading = ref(false)

const loadNotifications = async () => {
  isLoading.value = true
  try {
    const response = await notificationApi.getNotifications()
    notifications.value = response.data.notifications
  } catch (error) {
    console.error('加载通知失败:', error)
  } finally {
    isLoading.value = false
  }
}

const markAsRead = async (id) => {
  try {
    await notificationApi.markRead(id)
    const notification = notifications.value.find(n => n.id === id)
    if (notification) {
      notification.is_read = true
    }
  } catch (error) {
    console.error('标记已读失败:', error)
  }
}

const markAllAsRead = async () => {
  try {
    await notificationApi.markAllRead()
    notifications.value.forEach(n => n.is_read = true)
  } catch (error) {
    console.error('全部标记已读失败:', error)
  }
}

const formatTime = (timeStr) => {
  const date = new Date(timeStr)
  const now = new Date()
  const diff = now - date
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  
  if (days === 0) {
    const hours = Math.floor(diff / (1000 * 60 * 60))
    if (hours === 0) {
      const minutes = Math.floor(diff / (1000 * 60))
      return minutes <= 0 ? '刚刚' : `${minutes}分钟前`
    }
    return `${hours}小时前`
  } else if (days < 7) {
    return `${days}天前`
  } else {
    return date.toLocaleDateString('zh-CN')
  }
}

const getNotificationIcon = (type) => {
  switch (type) {
    case 'like':
      return 'mdi-heart'
    case 'comment':
      return 'mdi-message'
    case 'follow':
      return 'mdi-user-follow'
    case 'reply':
      return 'mdi-reply'
    default:
      return 'mdi-bell'
  }
}

onMounted(() => {
  if (!user.value) {
    router.push('/login')
    return
  }
  loadNotifications()
})
</script>

<template>
  <v-container class="max-w-2xl mx-auto px-4 py-8">
    <!-- 返回按钮 -->
    <v-btn 
      text 
      color="gray-600" 
      class="mb-6 hover:text-primary transition-colors"
      @click="router.back()"
    >
      <v-icon class="mr-2" size="20">mdi-arrow-left</v-icon>
      返回
    </v-btn>
    
    <v-card rounded="2xl" elevation="4" class="overflow-hidden">
      <!-- 标题栏 -->
      <v-card-title class="gradient-purple text-white py-6 px-8 flex items-center justify-between">
        <div class="flex items-center">
          <v-icon class="mr-3" size="24">mdi-bell</v-icon>
          <span class="font-bold text-xl">通知中心</span>
        </div>
        <v-btn 
          v-if="notifications.some(n => !n.is_read)"
          text 
          color="white" 
          class="hover:bg-white/20 rounded-lg transition-colors"
          @click="markAllAsRead"
        >
          <v-icon class="mr-1" size="18">mdi-check-all</v-icon>
          全部已读
        </v-btn>
      </v-card-title>
      
      <!-- 内容区域 -->
      <div class="p-6">
        <div v-if="isLoading" class="loading-center">
          <v-progress-circular indeterminate color="primary" :size="48" />
        </div>
        
        <v-list v-else-if="notifications.length > 0" class="space-y-3">
          <v-card 
            v-for="notification in notifications" 
            :key="notification.id"
            :class="[
              'card-hover cursor-pointer',
              notification.is_read ? 'bg-white' : 'bg-primary/5 border border-primary/20'
            ]"
            rounded="xl"
            @click="markAsRead(notification.id)"
          >
            <v-list-item class="px-4 py-3">
              <v-list-item-icon>
                <v-icon 
                  :color="notification.is_read ? 'gray' : 'primary'"
                  size="24"
                >
                  {{ getNotificationIcon(notification.type) }}
                </v-icon>
              </v-list-item-icon>
              <v-list-item-content>
                <v-list-item-title class="font-medium text-gray-800">{{ notification.content }}</v-list-item-title>
                <v-list-item-subtitle class="text-gray-400 text-sm">{{ formatTime(notification.created_at) }}</v-list-item-subtitle>
              </v-list-item-content>
              <v-list-item-action>
                <v-chip 
                  v-if="!notification.is_read" 
                  size="small" 
                  color="primary"
                  class="tag-purple"
                >
                  未读
                </v-chip>
              </v-list-item-action>
            </v-list-item>
          </v-card>
        </v-list>
        
        <div v-else class="empty-state">
          <v-icon size="96" color="gray-200" class="empty-state-icon">mdi-bell-off</v-icon>
          <p class="text-gray-400">暂无通知</p>
          <p class="text-gray-400 text-sm mt-1">有新消息时会在这里显示</p>
        </div>
      </div>
    </v-card>
  </v-container>
</template>