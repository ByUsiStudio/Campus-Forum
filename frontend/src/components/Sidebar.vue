<template>
  <div class="sidebar-component">
    <!-- 用户信息卡片 -->
    <v-card class="mb-4" elevation="0" color="transparent">
      <v-card-text v-if="user" class="text-center pa-4 bg-surface rounded-lg">
        <UserAvatar :user="user" :size="80" class="mb-3" />
        <div class="text-body-1 font-weight-bold mb-1">{{ user.display_name }}</div>
        <div class="text-body-2 text-medium-emphasis mb-2">@{{ user.username }}</div>
        <v-btn variant="tonal" size="small" to="/profile" color="primary" block>
          <v-icon start size="small">mdi-account</v-icon>
          个人中心
        </v-btn>
      </v-card-text>
      <v-card-text v-else class="text-center pa-4 bg-surface rounded-lg">
        <div class="text-body-2 text-medium-emphasis mb-3">登录后享受更多功能</div>
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

    <!-- 导航菜单 -->
    <v-card class="mb-4" elevation="0" color="transparent">
      <v-card-item class="pa-3 pb-0">
        <v-card-title class="text-subtitle-2 px-2">
          <v-icon size="small" class="mr-2">mdi-navigation</v-icon>
          快捷导航
        </v-card-title>
      </v-card-item>
      <v-list density="compact" class="pt-2 px-2" bg-color="transparent">
        <v-list-item
          v-for="item in sidebarItems"
          :key="item.link"
          :to="item.link"
          color="primary"
          rounded="lg"
          class="mb-1"
        >
          <template v-slot:prepend>
            <v-icon size="small">{{ getIcon(item.icon) }}</v-icon>
          </template>
          <v-list-item-title class="text-body-2">{{ item.title }}</v-list-item-title>
          <template v-slot:append>
            <span v-if="item.badge" class="badge badge-pill badge-danger text-xs">{{ item.badge }}</span>
          </template>
        </v-list-item>
      </v-list>
    </v-card>

    <!-- 聊天入口 -->
    <v-card v-if="user" class="mb-4" elevation="0" color="transparent">
      <v-card-item class="pa-3 pb-0">
        <v-card-title class="text-subtitle-2 px-2">
          <v-icon size="small" class="mr-2">mdi-message-circle</v-icon>
          社交互动
        </v-card-title>
      </v-card-item>
      <v-list density="compact" class="pt-2 px-2" bg-color="transparent">
        <v-list-item
          to="/friends"
          color="primary"
          rounded="lg"
          class="mb-1"
        >
          <template v-slot:prepend>
            <v-icon size="small">mdi-users</v-icon>
          </template>
          <v-list-item-title class="text-body-2">好友</v-list-item-title>
        </v-list-item>
        <v-list-item
          to="/chat"
          color="primary"
          rounded="lg"
          class="mb-1"
        >
          <template v-slot:prepend>
            <v-icon size="small">mdi-chat</v-icon>
          </template>
          <v-list-item-title class="text-body-2">聊天</v-list-item-title>
          <template v-slot:append>
            <span v-if="unreadCount > 0" class="badge badge-pill badge-danger text-xs">{{ unreadCount }}</span>
          </template>
        </v-list-item>
        <v-list-item
          to="/notifications"
          color="primary"
          rounded="lg"
        >
          <template v-slot:prepend>
            <v-icon size="small">mdi-bell</v-icon>
          </template>
          <v-list-item-title class="text-body-2">通知</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-card>

    <!-- 统计信息 -->
    <v-card elevation="0" color="transparent">
      <v-card-item class="pa-3 pb-0">
        <v-card-title class="text-subtitle-2 px-2">
          <v-icon size="small" class="mr-2">mdi-chart-bar</v-icon>
          论坛统计
        </v-card-title>
      </v-card-item>
      <v-list density="compact" class="pt-2 px-2" bg-color="transparent">
        <v-list-item rounded="lg" class="mb-1">
          <template v-slot:prepend>
            <v-icon size="small" color="primary">mdi-file-document</v-icon>
          </template>
          <v-list-item-title class="text-body-2">文章总数</v-list-item-title>
          <template v-slot:append>
            <v-chip size="x-small" color="primary" variant="tonal">{{ stats.totalArticles }}</v-chip>
          </template>
        </v-list-item>
        <v-list-item rounded="lg">
          <template v-slot:prepend>
            <v-icon size="small" color="success">mdi-account-group</v-icon>
          </template>
          <v-list-item-title class="text-body-2">用户总数</v-list-item-title>
          <template v-slot:append>
            <v-chip size="x-small" color="success" variant="tonal">{{ stats.totalUsers }}</v-chip>
          </template>
        </v-list-item>
      </v-list>
    </v-card>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted } from 'vue'
import api from '../api'
import { getUnreadCount } from '../api/chat'
import UserAvatar from './UserAvatar.vue'

export default {
  name: 'Sidebar',
  components: {
    UserAvatar
  },
  setup() {
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
        console.log('开始加载侧边栏配置...')
        const response = await api.get('/sidebar-config')
        console.log('侧边栏配置响应:', response.data)
        
        if (response.data && response.data.items) {
          sidebarItems.value = response.data.items
          console.log('侧边栏项目数量:', sidebarItems.value.length)
        } else {
          console.warn('侧边栏配置数据格式异常:', response.data)
          sidebarItems.value = [
            { title: '首页', link: '/', icon: 'mdi-home' },
            { title: '发布', link: '/create', icon: 'mdi-pencil' },
            { title: '搜索', link: '/search', icon: 'mdi-magnify' }
          ]
        }
      } catch (error) {
        console.error('加载侧边栏配置失败', error)
        console.error('错误详情:', error.response?.data || error.message)
        sidebarItems.value = [
          { title: '首页', link: '/', icon: 'mdi-home' },
          { title: '发布', link: '/create', icon: 'mdi-pencil' },
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
        const response = await getUnreadCount()
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

    return {
      user,
      sidebarItems,
      stats,
      unreadCount,
      getIcon
    }
  }
}
</script>

<style scoped>
.sidebar-component {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.bg-surface {
  background-color: rgb(var(--v-theme-surface));
}

.badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0 6px;
  border-radius: 10px;
  font-size: 10px;
  font-weight: 600;
  min-width: 18px;
  height: 18px;
}

.badge-danger {
  background-color: #ef4444;
  color: white;
}
</style>
