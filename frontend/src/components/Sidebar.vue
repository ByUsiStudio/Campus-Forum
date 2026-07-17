<template>
  <div class="sidebar-component">
    <div v-if="user" class="user-card mb-4">
      <div class="user-card-body text-center pa-4">
        <UserAvatar :user="user" :size="80" class="mb-3" />
        <div class="text-body-1 font-weight-bold mb-1">{{ user.display_name }}</div>
        <div class="text-body-2 text-medium-emphasis mb-2">@{{ user.username }}</div>
        <button class="layui-btn layui-btn-normal layui-btn-sm w-full" @click="goToProfile">
          <i class="fa-solid fa-user mr-2"></i>
          个人中心
        </button>
      </div>
    </div>
    <div v-else class="user-card mb-4">
      <div class="user-card-body text-center pa-4">
        <div class="text-body-2 text-medium-emphasis mb-3">登录后享受更多功能</div>
        <div class="d-flex gap-2 justify-center">
          <button class="layui-btn layui-btn-primary layui-btn-sm" @click="goToLogin">
            登录
          </button>
          <button class="layui-btn layui-btn-outline layui-btn-sm" @click="goToRegister">
            注册
          </button>
        </div>
      </div>
    </div>

    <div class="menu-card mb-4">
      <div class="menu-header pa-3 pb-0">
        <div class="menu-title">
          <i class="fa-solid fa-compass mr-2"></i>
          快捷导航
        </div>
      </div>
      <ul class="menu-list pt-2 px-2">
        <li v-for="item in sidebarItems" :key="item.link">
          <router-link :to="item.link" class="menu-item">
            <i :class="getIconClass(item.icon)" class="menu-icon"></i>
            <span>{{ item.title }}</span>
          </router-link>
        </li>
      </ul>
    </div>

    <div class="stats-card">
      <div class="stats-header pa-3 pb-0">
        <div class="stats-title">
          <i class="fa-solid fa-chart-bar mr-2"></i>
          论坛统计
        </div>
      </div>
      <ul class="stats-list pt-2 px-2">
        <li class="stats-item">
          <i class="fa-solid fa-file-lines text-primary"></i>
          <span class="stats-label">文章总数</span>
          <span class="stats-value">{{ stats.totalArticles }}</span>
        </li>
        <li class="stats-item">
          <i class="fa-solid fa-users text-success"></i>
          <span class="stats-label">用户总数</span>
          <span class="stats-value">{{ stats.totalUsers }}</span>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '../api'
import UserAvatar from './UserAvatar.vue'

export default {
  name: 'Sidebar',
  components: {
    UserAvatar
  },
  setup() {
    const router = useRouter()
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
        if (response.data && response.data.items) {
          sidebarItems.value = response.data.items
        } else {
          sidebarItems.value = [
            { title: '首页', link: '/', icon: 'home' }
          ]
        }
      } catch (error) {
        console.error('加载侧边栏配置失败', error)
        sidebarItems.value = [
          { title: '首页', link: '/', icon: 'home' }
        ]
      }
    }

    const loadStats = async () => {
      try {
        const response = await api.get('/statistics/overview')
        if (response.data.success) {
          stats.value.totalArticles = response.data.data.total_articles || 0
          stats.value.totalUsers = response.data.data.total_users || 0
        }
      } catch (error) {
        console.error('加载统计失败', error)
      }
    }

    const getIconClass = (icon) => {
      if (!icon) return 'fa-solid fa-link'
      if (icon.startsWith('mdi-')) {
        const mdiToFa = {
          'mdi-home': 'fa-solid fa-house',
          'mdi-account': 'fa-solid fa-user',
          'mdi-navigation': 'fa-solid fa-compass',
          'mdi-chart-bar': 'fa-solid fa-chart-bar',
          'mdi-file-document': 'fa-solid fa-file-lines',
          'mdi-account-group': 'fa-solid fa-users',
          'mdi-link': 'fa-solid fa-link',
          'mdi-bookmark': 'fa-solid fa-bookmark',
          'mdi-trophy': 'fa-solid fa-trophy',
          'mdi-calendar': 'fa-solid fa-calendar',
          'mdi-search': 'fa-solid fa-search',
          'mdi-bell': 'fa-solid fa-bell',
          'mdi-cog': 'fa-solid fa-gear',
          'mdi-logout': 'fa-solid fa-right-from-bracket',
          'mdi-pencil': 'fa-solid fa-pencil',
          'mdi-tag': 'fa-solid fa-tag',
          'mdi-folder': 'fa-solid fa-folder',
          'mdi-image': 'fa-solid fa-image',
          'mdi-video': 'fa-solid fa-video'
        }
        return mdiToFa[icon] || 'fa-solid fa-circle-question'
      }
      return `fa-solid fa-${icon}`
    }

    const goToProfile = () => {
      router.push('/profile')
    }

    const goToLogin = () => {
      router.push('/login')
    }

    const goToRegister = () => {
      router.push('/register')
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
      getIconClass,
      goToProfile,
      goToLogin,
      goToRegister
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

.user-card, .menu-card, .stats-card {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
}

.user-card-body {
  background: #f8f9fa;
  border-radius: 8px;
}

.menu-header, .stats-header {
  display: flex;
  align-items: center;
}

.menu-title, .stats-title {
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.menu-list, .stats-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  color: #666;
  text-decoration: none;
  font-size: 14px;
  border-radius: 6px;
  transition: all 0.2s ease;

  &:hover {
    background: rgba(30, 159, 255, 0.08);
    color: #1E9FFF;
  }
}

.menu-icon {
  font-size: 14px;
}

.stats-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  border-radius: 6px;

  &:hover {
    background: #f8f9fa;
  }
}

.stats-label {
  flex: 1;
  font-size: 14px;
  color: #666;
}

.stats-value {
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.text-primary {
  color: #1E9FFF;
}

.text-success {
  color: #52C41A;
}

.text-medium-emphasis {
  color: #999;
}

.text-body-1 {
  font-size: 16px;
}

.text-body-2 {
  font-size: 14px;
}

.font-weight-bold {
  font-weight: 600;
}

.w-full {
  width: 100%;
}

.d-flex {
  display: flex;
}

.gap-2 {
  gap: 8px;
}

.justify-center {
  justify-content: center;
}

.pa-3 {
  padding: 12px;
}

.pa-4 {
  padding: 16px;
}

.pb-0 {
  padding-bottom: 0;
}

.pt-2 {
  padding-top: 8px;
}

.px-2 {
  padding-left: 8px;
  padding-right: 8px;
}

.mb-2 {
  margin-bottom: 8px;
}

.mb-3 {
  margin-bottom: 12px;
}

.mb-4 {
  margin-bottom: 16px;
}

.mr-2 {
  margin-right: 8px;
}
</style>