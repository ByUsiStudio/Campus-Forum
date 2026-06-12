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

onMounted(() => {
  if (!user.value) {
    router.push('/login')
    return
  }
  loadNotifications()
})
</script>

<template>
  <v-container class="max-w-4xl mx-auto py-8">
    <v-card rounded="xl" elevation="4">
      <v-card-title class="gradient-purple text-white flex items-center justify-between">
        <div class="flex items-center">
          <v-icon class="mr-2">mdi-bell</v-icon>
          <span class="font-bold text-xl">通知中心</span>
        </div>
        <v-btn 
          text 
          color="white" 
          @click="markAllAsRead"
          v-if="notifications.some(n => !n.is_read)"
        >
          全部已读
        </v-btn>
      </v-card-title>
      
      <v-card-text>
        <div v-if="isLoading" class="text-center py-8">
          <v-progress-circular indeterminate color="primary" />
        </div>
        
        <v-list v-else-if="notifications.length > 0">
          <v-list-item 
            v-for="notification in notifications" 
            :key="notification.id"
            :class="notification.is_read ? '' : 'bg-primary/5'"
            class="border-b border-gray-100 last:border-0 cursor-pointer"
            @click="markAsRead(notification.id)"
          >
            <v-list-item-icon>
              <v-icon 
                :color="notification.is_read ? 'gray' : 'primary'"
                size="24"
              >
                {{ notification.type === 'like' ? 'mdi-heart' : 'mdi-message' }}
              </v-icon>
            </v-list-item-icon>
            <v-list-item-content>
              <v-list-item-title>{{ notification.content }}</v-list-item-title>
              <v-list-item-subtitle>{{ formatTime(notification.created_at) }}</v-list-item-subtitle>
            </v-list-item-content>
          </v-list-item>
        </v-list>
        
        <div v-else class="text-center py-12">
          <v-icon size="64" color="gray" class="mx-auto mb-4">mdi-bell-off</v-icon>
          <p class="text-gray-500">暂无通知</p>
        </div>
      </v-card-text>
    </v-card>
  </v-container>
</template>