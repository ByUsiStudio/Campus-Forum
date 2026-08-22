<template>
  <div class="sidebar-component">
    <!-- 用户信息卡片 -->
    <div class="sb-card sb-user">
      <div v-if="user" class="user-panel">
        <UserAvatar :user="user" :size="72" class="avatar-wrap" />
        <div class="user-name-title">{{ user.display_name }}</div>
        <div class="user-name-sub">@{{ user.username }}</div>
        <el-button type="primary" round size="small" class="profile-btn" @click="goTo('/profile')">
          <el-icon class="btn-icon"><User /></el-icon>
          个人中心
        </el-button>
      </div>
      <div v-else class="user-panel">
        <div class="user-guest-icon">
          <el-icon :size="28"><Avatar /></el-icon>
        </div>
        <div class="user-name-sub mb-3">登录后享受更多功能</div>
        <div class="auth-actions">
          <el-button type="primary" round size="small" @click="goTo('/login')">登录</el-button>
          <el-button round size="small" plain @click="goTo('/register')">注册</el-button>
        </div>
      </div>
    </div>

    <!-- 导航菜单 -->
    <div class="sb-card">
      <div class="card-header">
        <el-icon class="header-icon"><Compass /></el-icon>
        <span>快捷导航</span>
      </div>
      <div class="nav-list">
        <div
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
        </div>
      </div>
    </div>

    <!-- 统计信息 -->
    <div class="sb-card">
      <div class="card-header">
        <el-icon class="header-icon"><DataAnalysis /></el-icon>
        <span>论坛统计</span>
      </div>
      <div class="stat-list">
        <div class="stat-row">
          <div class="stat-label">
            <el-icon class="nav-icon primary"><Document /></el-icon>
            <span>文章总数</span>
          </div>
          <el-tag size="small" type="primary" effect="light" round>{{ stats.totalArticles }}</el-tag>
        </div>
        <div class="stat-row">
          <div class="stat-label">
            <el-icon class="nav-icon success"><User /></el-icon>
            <span>用户总数</span>
          </div>
          <el-tag size="small" type="success" effect="light" round>{{ stats.totalUsers }}</el-tag>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { User, Compass, Document, DataAnalysis, Avatar } from '@element-plus/icons-vue'
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
      DataAnalysis,
      Avatar
    }
  }
}
</script>

<style scoped>
.sidebar-component {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.sb-card {
  background: var(--campus-surface);
  border: 1px solid var(--campus-border);
  border-radius: var(--campus-radius);
  box-shadow: var(--campus-shadow-sm);
  padding: 16px;
}

.sb-user {
  background: linear-gradient(160deg, var(--campus-surface), var(--campus-surface-2));
}

.user-panel {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.user-name-title {
  font-size: 16px;
  font-weight: 700;
  color: var(--campus-text);
  line-height: 1.3;
  margin-top: 10px;
}

.user-name-sub {
  font-size: 13px;
  color: var(--campus-text-secondary);
  margin: 2px 0 10px;
}

.user-guest-icon {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--campus-primary-soft);
  color: var(--campus-primary);
  margin-bottom: 8px;
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
  gap: 10px;
  justify-content: center;
  width: 100%;
}
.auth-actions .el-button {
  flex: 1;
}

.card-header {
  display: flex;
  align-items: center;
  font-size: 13px;
  font-weight: 700;
  color: var(--campus-text);
  margin-bottom: 8px;
}

.header-icon {
  margin-right: 8px;
  color: var(--campus-primary);
}

.nav-list {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.nav-item {
  padding: 9px 10px;
  border-radius: 10px;
  transition: var(--campus-transition);
}

.nav-item:hover {
  background: var(--campus-primary-soft);
}

.nav-item-inner {
  display: flex;
  align-items: center;
  gap: 10px;
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
  color: var(--campus-success);
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

.stat-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.stat-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: 6px 2px;
}

.stat-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: var(--campus-text-secondary);
}
</style>
