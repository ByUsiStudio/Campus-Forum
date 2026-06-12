<template>
  <div class="sidebar-component">
    <!-- 用户信息卡片 -->
    <v-card class="user-card mb-4">
      <v-card-text v-if="user" class="text-center pa-5">
        <UserAvatar :user="user" :size="72" class="mb-3" />
        <div class="text-subtitle-1 font-weight-bold text-grey-darken-3 mb-1">
          {{ user.display_name || user.username }}
        </div>
        <div class="text-caption text-grey mb-4">@{{ user.username }}</div>
        <v-btn
          variant="tonal"
          size="small"
          to="/profile"
          color="primary"
          block
        >
          <v-icon start size="16">mdi-account-outline</v-icon>
          个人中心
        </v-btn>
      </v-card-text>
      <v-card-text v-else class="text-center pa-5">
        <v-icon size="48" color="primary-lighten-2" class="mb-3">mdi-account-circle-outline</v-icon>
        <div class="text-body-2 text-grey mb-4">登录后享受更多功能</div>
        <div class="d-flex gap-2 justify-center">
          <v-btn variant="flat" color="primary" to="/login" size="small">
            登录
          </v-btn>
          <v-btn variant="outlined" color="primary" to="/register" size="small">
            注册
          </v-btn>
        </div>
      </v-card-text>
    </v-card>

    <!-- 快捷导航 -->
    <v-card class="nav-card mb-4">
      <v-card-title class="text-subtitle-2 pa-4 pb-2">
        <v-icon size="18" class="mr-2">mdi-compass-outline</v-icon>
        快捷导航
      </v-card-title>
      <v-list density="compact" class="pt-0 pb-2" bg-color="transparent">
        <v-list-item
          v-for="item in sidebarItems"
          :key="item.link"
          :to="item.link"
          color="primary"
          rounded="lg"
          class="mx-2"
        >
          <template #prepend>
            <v-icon size="18">{{ getIcon(item.icon) }}</v-icon>
          </template>
          <v-list-item-title class="text-body-2">{{ item.title }}</v-list-item-title>
          <template #append v-if="item.badge">
            <v-badge
              :content="item.badge"
              color="error"
              inline
            ></v-badge>
          </template>
        </v-list-item>
      </v-list>
    </v-card>

    <!-- 社交互动 -->
    <v-card v-if="user" class="nav-card mb-4">
      <v-card-title class="text-subtitle-2 pa-4 pb-2">
        <v-icon size="18" class="mr-2">mdi-account-group-outline</v-icon>
        社交互动
      </v-card-title>
      <v-list density="compact" class="pt-0 pb-2" bg-color="transparent">
        <v-list-item
          to="/friends"
          color="primary"
          rounded="lg"
          class="mx-2"
        >
          <template #prepend>
            <v-icon size="18">mdi-account-multiple-outline</v-icon>
          </template>
          <v-list-item-title class="text-body-2">好友</v-list-item-title>
        </v-list-item>
        <v-list-item
          to="/chat"
          color="primary"
          rounded="lg"
          class="mx-2"
        >
          <template #prepend>
            <v-icon size="18">mdi-chat-outline</v-icon>
          </template>
          <v-list-item-title class="text-body-2">聊天</v-list-item-title>
          <template #append v-if="unreadCount > 0">
            <v-badge
              :content="unreadCount"
              color="error"
              inline
            ></v-badge>
          </template>
        </v-list-item>
        <v-list-item
          to="/notifications"
          color="primary"
          rounded="lg"
          class="mx-2"
        >
          <template #prepend>
            <v-icon size="18">mdi-bell-outline</v-icon>
          </template>
          <v-list-item-title class="text-body-2">通知</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-card>

    <!-- 论坛统计 -->
    <v-card class="nav-card">
      <v-card-title class="text-subtitle-2 pa-4 pb-2">
        <v-icon size="18" class="mr-2">mdi-chart-line</v-icon>
        论坛统计
      </v-card-title>
      <v-list density="compact" class="pt-0 pb-3" bg-color="transparent">
        <v-list-item rounded="lg" class="mx-2">
          <template #prepend>
            <v-icon size="18" color="primary">mdi-file-document-outline</v-icon>
          </template>
          <v-list-item-title class="text-body-2">文章总数</v-list-item-title>
          <template #append>
            <v-chip size="x-small" color="primary" variant="tonal">
              {{ stats.totalArticles }}
            </v-chip>
          </template>
        </v-list-item>
        <v-list-item rounded="lg" class="mx-2">
          <template #prepend>
            <v-icon size="18" color="success">mdi-account-group-outline</v-icon>
          </template>
          <v-list-item-title class="text-body-2">用户总数</v-list-item-title>
          <template #append>
            <v-chip size="x-small" color="success" variant="tonal">
              {{ stats.totalUsers }}
            </v-chip>
          </template>
        </v-list-item>
      </v-list>
    </v-card>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import api from '../api'
import { chatApi } from '../api'
import UserAvatar from './UserAvatar.vue'

const user = ref(null)
const sidebarItems = ref([])
const stats = ref({
  totalArticles: 0,
  totalUsers: 0
})
const unreadCount = ref(0)

const loadUser = () => {
  const userStr = localStorage.getItem('user')
  if (userStr) {
    user.value = JSON.parse(userStr)
  }
}

const loadSidebarConfig = async () => {
  try {
    const response = await api.get('/sidebar-config')
    if (response.data && response.data.items) {
      sidebarItems.value = response.data.items
    } else {
      sidebarItems.value = [
        { title: '首页', link: '/', icon: 'mdi-home-outline' },
        { title: '发布', link: '/create', icon: 'mdi-pencil-outline' },
        { title: '搜索', link: '/search', icon: 'mdi-magnify' }
      ]
    }
  } catch (error) {
    console.error('加载侧边栏配置失败', error)
    sidebarItems.value = [
      { title: '首页', link: '/', icon: 'mdi-home-outline' },
      { title: '发布', link: '/create', icon: 'mdi-pencil-outline' },
      { title: '搜索', link: '/search', icon: 'mdi-magnify' }
    ]
  }
}

const loadStats = async () => {
  try {
    const response = await api.get('/articles', { params: { page: 1, page_size: 1 } })
    stats.value.totalArticles = response.data.total || 0
    stats.value.totalUsers = response.data.total || 1
  } catch (error) {
    console.error('加载统计失败', error)
  }
}

const loadUnreadCount = async () => {
  if (!user.value) return
  try {
    const response = await chatApi.getUnreadCount()
    unreadCount.value = response.data.unread_count || 0
  } catch (error) {
    console.error('加载未读消息数失败:', error)
  }
}

const getIcon = (icon) => {
  if (!icon) return 'mdi-link'
  if (icon.match(/[\u{1F600}-\u{1F64F}]/u)) return icon
  return icon.startsWith('mdi-') ? icon : `mdi-${icon}`
}

let refreshInterval = null

onMounted(() => {
  loadUser()
  loadSidebarConfig()
  loadStats()
  loadUnreadCount()

  refreshInterval = setInterval(() => {
    loadUnreadCount()
  }, 5000)
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})
</script>

<style scoped>
.sidebar-component {
  position: sticky;
  top: 88px;
}

.user-card {
  border: 1px solid rgba(149, 117, 205, 0.12) !important;
  border-radius: 16px !important;
}

.nav-card {
  border: 1px solid rgba(149, 117, 205, 0.08) !important;
  border-radius: 16px !important;
}

:deep(.v-list-item) {
  margin: 2px 0;
  border-radius: 10px;
}

:deep(.v-list-item--active) {
  background-color: rgba(149, 117, 205, 0.1) !important;
}
</style>
