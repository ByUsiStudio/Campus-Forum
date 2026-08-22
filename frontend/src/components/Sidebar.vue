<template>
  <div class="sidebar-component">
    <!-- 用户信息卡片 -->
    <el-card class="mb-4 sidebar-card" shadow="never">
      <div v-if="user" class="text-center pa-4 user-panel">
        <UserAvatar :user="user" :size="80" class="mb-3 avatar-wrap" />
        <div class="user-name-title mb-1">{{ user.display_name }}</div>
        <div class="user-name-sub mb-2">@{{ user.username }}</div>
        <el-button type="primary" size="small" plain class="profile-btn" @click="goTo('/profile')">
          <el-icon class="btn-icon"><User /></el-icon>
          个人中心
        </el-button>
      </div>
      <div v-else class="text-center pa-4 user-panel">
        <div class="user-name-sub mb-3">登录后享受更多功能</div>
        <div class="auth-actions">
          <el-button type="primary" size="small" @click="goTo('/login')">
            登录
          </el-button>
          <el-button size="small" plain type="primary" @click="goTo('/register')">
            注册
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- 导航菜单 -->
    <el-card class="mb-4 sidebar-card" shadow="never">
      <template #header>
        <div class="card-header">
          <el-icon class="header-icon"><Compass /></el-icon>
          <span>快捷导航</span>
        </div>
      </template>
      <el-list class="nav-list">
        <el-list-item
          v-for="item in sidebarItems"
          :key="item.link"
          class="nav-item cursor-pointer"
          @click="goTo(item.link)"
        >
          <div class="nav-item-inner">
            <el-icon v-if="!isEmoji(item.icon)" class="nav-icon primary">
              <component :is="resolveIcon(item.icon)" />
            </el-icon>
            <span v-else class="nav-emoji">{{ item.icon }}</span>
            <span class="nav-title">{{ item.title }}</span>
          </div>
        </el-list-item>
      </el-list>
    </el-card>

    <!-- 统计信息 -->
    <el-card class="sidebar-card" shadow="never">
      <template #header>
        <div class="card-header">
          <el-icon class="header-icon"><DataAnalysis /></el-icon>
          <span>论坛统计</span>
        </div>
      </template>
      <el-list class="nav-list">
        <el-list-item>
          <div class="stat-item">
            <div class="nav-item-inner">
              <el-icon class="nav-icon primary"><Document /></el-icon>
              <span class="nav-title">文章总数</span>
            </div>
            <el-tag size="small" type="primary" effect="plain">{{ stats.totalArticles }}</el-tag>
          </div>
        </el-list-item>
        <el-list-item>
          <div class="stat-item">
            <div class="nav-item-inner">
              <el-icon class="nav-icon success"><User /></el-icon>
              <span class="nav-title">用户总数</span>
            </div>
            <el-tag size="small" type="success" effect="plain">{{ stats.totalUsers }}</el-tag>
          </div>
        </el-list-item>
      </el-list>
    </el-card>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { User, Compass, Document, DataAnalysis } from '@element-plus/icons-vue'
import api from '../api'
import UserAvatar from './UserAvatar.vue'

// mdi 图标名 -> Element Plus 图标组件名（全部图标已全局注册为同名组件）
const ICON_MAP = {
  'mdi-home': 'Home',
  'mdi-account': 'User',
  'mdi-account-group': 'User',
  'mdi-file-document': 'Document',
  'mdi-navigation': 'Compass',
  'mdi-chart-bar': 'DataAnalysis',
  'mdi-link': 'Link',
  'mdi-video': 'VideoCamera',
  'mdi-bell': 'Bell',
  'mdi-magnify': 'Search',
  'mdi-star': 'Star',
  'mdi-emoticon': 'Avatar',
  'mdi-trophy': 'Trophy',
  'mdi-lock': 'Lock',
  'mdi-login': 'SwitchButton',
  'mdi-book': 'Reading',
  'mdi-tag': 'PriceTag',
  'mdi-folder': 'Folder',
  'mdi-school': 'School'
}

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
          console.warn('侧边栏配置数据格式异常:', response.data)
          sidebarItems.value = [
            { title: '首页', link: '/', icon: 'mdi-home' }
          ]
        }
      } catch (error) {
        console.error('加载侧边栏配置失败', error)
        sidebarItems.value = [
          { title: '首页', link: '/', icon: 'mdi-home' }
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

    // 判断是否为 emoji 图标
    const isEmoji = (icon) => {
      if (!icon) return false
      return icon.match(/[\u{1F300}-\u{1FAFF}]/u) !== null
    }

    // 解析后端图标名 -> Element Plus 全局图标组件名；未知回退 Link
    const resolveIcon = (icon) => {
      if (!icon) return 'Link'
      const name = icon.startsWith('mdi-') ? icon : `mdi-${icon}`
      return ICON_MAP[name] || 'Link'
    }

    const goTo = (link) => {
      if (link) router.push(link)
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
      isEmoji,
      resolveIcon,
      goTo,
      User,
      Compass,
      Document,
      DataAnalysis
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

.sidebar-card {
  border-radius: var(--campus-radius);
  border-color: var(--campus-border);
  box-shadow: var(--campus-shadow);
}

.user-panel {
  background: var(--campus-surface);
  border-radius: var(--campus-radius);
}

.avatar-wrap {
  display: inherit;
}

.user-name-title {
  font-size: 16px;
  font-weight: 700;
  line-height: 1.3;
}

.user-name-sub {
  font-size: 13px;
  color: var(--campus-text-secondary);
}

.profile-btn {
  width: 100%;
  margin-top: 4px;
}

.btn-icon {
  margin-right: 4px;
}

.auth-actions {
  display: flex;
  gap: 8px;
  justify-content: center;
}

.card-header {
  display: flex;
  align-items: center;
  font-size: 13px;
  font-weight: 600;
}

.header-icon {
  margin-right: 8px;
  color: var(--campus-primary);
}

.nav-list {
  padding: 4px 0;
}

.nav-item {
  margin-bottom: 4px;
  border-radius: 8px;
  transition: background-color 0.2s;
}

.nav-item:hover {
  background-color: rgba(103, 80, 164, 0.08);
}

.nav-item-inner {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.nav-icon {
  font-size: 16px;
  flex-shrink: 0;
}

.nav-icon.primary {
  color: var(--campus-primary);
}

.nav-icon.success {
  color: #67c23a;
}

.nav-emoji {
  font-size: 16px;
  line-height: 1;
}

.nav-title {
  font-size: 14px;
  color: var(--campus-text);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.stat-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}
</style>
