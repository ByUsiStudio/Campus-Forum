<template>
  <div class="sidebar-component">
    <v-card class="mb-4">
      <v-card-text v-if="user" class="text-center pa-4">
        <v-avatar size="80" class="mb-3">
          <v-img :src="user.avatar" :alt="user.display_name"></v-img>
        </v-avatar>
        <div class="text-h6 mb-1">{{ user.display_name }}</div>
        <v-btn variant="text" size="small" to="/profile" color="primary">
          个人中心
        </v-btn>
      </v-card-text>
      <v-card-text v-else class="text-center pa-4">
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
    
    <v-card class="mb-4">
      <v-card-title class="text-subtitle-2 pa-3 pb-0">导航</v-card-title>
      <v-list density="compact" class="pt-0">
        <v-list-item
          v-for="item in sidebarItems"
          :key="item.link"
          :href="item.link"
          :to="item.link"
          color="primary"
        >
          <template v-slot:prepend>
            <v-icon size="small">{{ getIcon(item.icon) }}</v-icon>
          </template>
          <v-list-item-title>{{ item.title }}</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-card>
    
    <v-card>
      <v-card-title class="text-subtitle-2 pa-3 pb-0">统计信息</v-card-title>
      <v-list density="compact" class="pt-0">
        <v-list-item>
          <template v-slot:prepend>
            <v-icon size="small" color="primary">mdi-file-document</v-icon>
          </template>
          <v-list-item-title>总文章数：{{ stats.totalArticles }}</v-list-item-title>
        </v-list-item>
        <v-list-item>
          <template v-slot:prepend>
            <v-icon size="small" color="primary">mdi-account-group</v-icon>
          </template>
          <v-list-item-title>总用户数：{{ stats.totalUsers }}</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-card>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import api from '../api'

export default {
  name: 'Sidebar',
  setup() {
    const user = ref(null)
    const sidebarItems = ref([])
    const stats = ref({
      totalArticles: 0,
      totalUsers: 0
    })
    
    const loadUser = () => {
      const userStr = localStorage.getItem('user')
      if (userStr) {
        user.value = JSON.parse(userStr)
      }
    }
    
    const loadSidebarConfig = async () => {
      try {
        const response = await api.get('/sidebar-config')
        sidebarItems.value = response.data.items || []
      } catch (error) {
        console.error('加载侧边栏配置失败', error)
        sidebarItems.value = [
          { title: '首页', link: '/', icon: 'mdi-home' }
        ]
      }
    }
    
    const loadStats = async () => {
      try {
        const response = await api.get('/articles', { params: { page: 1, page_size: 1 } })
        stats.value.totalArticles = response.data.total
        stats.value.totalUsers = 1
      } catch (error) {
        console.error('加载统计失败', error)
      }
    }
    
    const getIcon = (icon) => {
      if (!icon) return 'mdi-link'
      // 如果是emoji
      if (icon.match(/[\u{1F600}-\u{1F64F}]/u)) return icon
      // 否则当作MDI图标
      return icon.startsWith('mdi-') ? icon : `mdi-${icon}`
    }
    
    onMounted(() => {
      loadUser()
      loadSidebarConfig()
      loadStats()
    })
    
    return {
      user,
      sidebarItems,
      stats,
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
</style>
