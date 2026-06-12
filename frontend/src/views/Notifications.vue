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
    notifications.value = response.data.notifications || []
  } catch (error) {
    console.error('加载通知失败:', error)
  } finally {
    isLoading.value = false
  }
}

const markRead = async (id) => {
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

const markAllRead = async () => {
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
    return hours === 0 ? '刚刚' : `${hours}小时前`
  }
  return `${days}天前`
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
  <v-app>
    <v-app-bar app>
      <v-btn icon @click="router.push('/')">
        <v-icon>mdi-arrow-left</v-icon>
      </v-btn>
      <v-toolbar-title>通知</v-toolbar-title>
      <v-spacer></v-spacer>
      <v-btn text @click="markAllRead">全部已读</v-btn>
    </v-app-bar>
    
    <v-container class="py-6">
      <v-card v-if="notifications.length > 0">
        <v-list>
          <v-list-item
            v-for="notification in notifications"
            :key="notification.id"
            :class="notification.is_read ? '' : 'bg-grey-light'"
            @click="markRead(notification.id)"
          >
            <v-list-item-icon>
              <v-icon :color="notification.is_read ? 'grey' : 'primary'">
                {{ notification.type === 'comment' ? 'mdi-comment' : 'mdi-bell' }}
              </v-icon>
            </v-list-item-icon>
            <v-list-item-content>
              <v-list-item-title>{{ notification.content }}</v-list-item-title>
              <v-list-item-subtitle>{{ formatTime(notification.created_at) }}</v-list-item-subtitle>
            </v-list-item-content>
          </v-list-item>
        </v-list>
      </v-card>
      
      <v-card v-else class="text-center py-12">
        <v-icon size="64" color="grey">mdi-inbox</v-icon>
        <p class="mt-4 text-grey">暂无通知</p>
      </v-card>
    </v-container>
  </v-app>
</template>

<style scoped>
.bg-grey-light {
  background-color: rgba(0, 0, 0, 0.05);
}
</style>
