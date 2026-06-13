<template>
  <v-app v-if="isInitialized && isAdmin">
    <!-- 侧边栏导航 -->
    <v-navigation-drawer
      v-model="drawerOpen"
      :rail="sidebarCollapsed && !isMobile"
      :permanent="!isMobile"
      :temporary="isMobile"
      :width="260"
      color="primary"
      class="admin-drawer"
    >
      <!-- Logo 区域 -->
      <div class="drawer-header pa-4">
        <div class="d-flex align-center">
          <v-icon size="32" color="white">mdi-shield-crown</v-icon>
          <transition name="fade">
            <span v-if="!sidebarCollapsed" class="ml-3 text-h6 font-weight-bold text-white">
              校园论坛
            </span>
          </transition>
        </div>
        <v-btn
          icon
          size="small"
          variant="text"
          @click="sidebarCollapsed = !sidebarCollapsed"
          color="white"
          class="collapse-btn"
        >
          <v-icon>{{ sidebarCollapsed ? 'mdi-chevron-right' : 'mdi-chevron-left' }}</v-icon>
        </v-btn>
      </div>

      <v-divider color="white" opacity="0.2" />

      <!-- 导航菜单 -->
      <v-list nav density="compact" class="pa-2">
        <v-list-item
          v-for="item in adminItems"
          :key="item.route"
          :to="{ name: item.route }"
          color="white"
          rounded="lg"
          class="mb-1 nav-item"
          active-class="nav-item-active"
        >
          <template v-slot:prepend>
            <v-icon :size="24">{{ item.icon }}</v-icon>
          </template>
          <v-list-item-title class="font-weight-medium">{{ item.title }}</v-list-item-title>
          
          <template v-slot:append v-if="item.badge && item.badge() > 0">
            <v-badge
              :content="item.badge()"
              color="error"
              inline
            />
          </template>
        </v-list-item>
      </v-list>

      <template v-slot:append>
        <v-divider color="white" opacity="0.2" />
        <!-- 版本信息 -->
        <div class="version-info pa-3">
          <div class="d-flex align-center mb-1">
            <v-icon size="14" color="white" class="mr-1">mdi-information</v-icon>
            <span class="text-caption text-white text-truncate">版本信息</span>
          </div>
          <div class="text-caption text-white opacity-70">前端: {{ frontendVersion }}</div>
          <div class="text-caption text-white opacity-70">后端: {{ backendVersion }}</div>
        </div>
        <v-list nav density="compact" class="pa-2">
          <v-list-item
            to="/"
            color="white"
            rounded="lg"
            class="nav-item"
          >
            <template v-slot:prepend>
              <v-icon size="24">mdi-home</v-icon>
            </template>
            <v-list-item-title class="font-weight-medium">返回首页</v-list-item-title>
          </v-list-item>
        </v-list>
      </template>
    </v-navigation-drawer>

    <!-- 顶部栏 -->
    <v-app-bar flat color="white" border="b" height="64" class="admin-header">
      <v-app-bar-nav-icon @click="drawerOpen = !drawerOpen" class="d-lg-none" />
      
      <v-breadcrumbs :items="breadcrumbs" class="pa-0">
        <template v-slot:divider>
          <v-icon size="small">mdi-chevron-right</v-icon>
        </template>
      </v-breadcrumbs>

      <v-spacer />

      <!-- 搜索框 -->
      <v-text-field
        v-model="searchQuery"
        placeholder="搜索..."
        variant="outlined"
        density="compact"
        prepend-inner-icon="mdi-magnify"
        hide-details
        style="max-width: 300px;"
        class="mx-4 d-none d-sm-block"
      />

      <!-- 通知按钮 -->
      <v-btn icon variant="text" size="small" class="mx-1">
        <v-badge :content="notificationCount" color="error" floating>
          <v-icon>mdi-bell-outline</v-icon>
        </v-badge>
      </v-btn>

      <!-- 用户菜单 -->
      <v-menu offset-y>
        <template v-slot:activator="{ props }">
          <v-btn
            v-bind="props"
            variant="text"
            class="ml-2"
          >
            <v-avatar size="32" color="primary">
              <v-img v-if="currentUser?.avatar" :src="currentUser.avatar" />
              <span v-else class="text-white text-body-2">
                {{ currentUser?.display_name?.[0] || currentUser?.username?.[0] || 'A' }}
              </span>
            </v-avatar>
            <span class="ml-2 d-none d-sm-inline">{{ currentUser?.display_name || currentUser?.username }}</span>
            <v-icon size="small" class="ml-1">mdi-chevron-down</v-icon>
          </v-btn>
        </template>
        <v-list density="compact" min-width="200">
          <v-list-item to="/profile" prepend-icon="mdi-account">
            <v-list-item-title>个人资料</v-list-item-title>
          </v-list-item>
          <v-divider />
          <v-list-item @click="handleLogout" prepend-icon="mdi-logout">
            <v-list-item-title>退出登录</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>
    </v-app-bar>

    <!-- 主要内容区域 -->
    <v-main class="admin-main">
      <v-container fluid class="pa-6">
        <transition name="fade-slide" mode="out-in">
          <router-view />
        </transition>
      </v-container>
    </v-main>

    <!-- 移动端底部导航 -->
    <v-bottom-navigation grow color="primary" class="d-lg-none admin-bottom-nav">
      <v-btn
        v-for="item in bottomNavItems"
        :key="item.route"
        :to="{ name: item.route }"
      >
        <v-icon>{{ item.icon }}</v-icon>
        <span class="text-caption">{{ item.title }}</span>
      </v-btn>
    </v-bottom-navigation>
  </v-app>

  <!-- 权限不足 -->
  <v-container v-else-if="isInitialized && !isAdmin" fluid class="fill-height">
    <v-row justify="center" align="center">
      <v-col cols="12" sm="8" md="6" lg="4">
        <v-card class="text-center pa-8" elevation="2">
          <v-icon size="80" color="error" class="mb-4">mdi-shield-alert</v-icon>
          <div class="text-h5 font-weight-bold mb-2">权限不足</div>
          <div class="text-body-2 text-medium-emphasis mb-6">您没有访问管理后台的权限</div>
          <v-btn color="primary" to="/" prepend-icon="mdi-home">返回首页</v-btn>
        </v-card>
      </v-col>
    </v-row>
  </v-container>

  <!-- 加载状态 -->
  <v-container v-else fluid class="fill-height">
    <v-row justify="center" align="center">
      <v-progress-circular indeterminate color="primary" size="48" />
    </v-row>
  </v-container>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import api from '../api'

const router = useRouter()
const route = useRoute()

const isInitialized = ref(false)
const isAdmin = ref(false)
const isMobile = ref(false)
const deletionCount = ref(0)
const notificationCount = ref(3)
const sidebarCollapsed = ref(false)
const drawerOpen = ref(true)
const searchQuery = ref('')
const currentUser = ref(null)
const backendVersion = ref('加载中...')
const frontendVersion = ref(typeof __FRONTEND_VERSION__ !== 'undefined' ? __FRONTEND_VERSION__ : 'unknown')

const adminItems = [
  { route: 'AdminIndex', title: '数据概览', icon: 'mdi-view-dashboard' },
  { route: 'AdminUsers', title: '用户管理', icon: 'mdi-account-group', badge: null },
  { route: 'AdminArticles', title: '文章管理', icon: 'mdi-file-document' },
  { route: 'AdminComments', title: '评论管理', icon: 'mdi-comment-text' },
  { route: 'AdminCategories', title: '分区管理', icon: 'mdi-shape' },
  { route: 'AdminTitles', title: '头衔管理', icon: 'mdi-medal' },
  { route: 'AdminSidebar', title: '侧边栏', icon: 'mdi-web' },
  { route: 'AdminDeletions', title: '删除申请', icon: 'mdi-delete', badge: () => deletionCount.value },
  { route: 'AdminReports', title: '举报管理', icon: 'mdi-flag' },
  { route: 'AdminAnnouncement', title: '公告管理', icon: 'mdi-bullhorn' },
  { route: 'AdminUserNotifications', title: '用户通知', icon: 'mdi-bell' },
  { route: 'AdminSystemLogs', title: '系统日志', icon: 'mdi-text-box-multiple' },
  { route: 'AdminSiteConfig', title: '网站配置', icon: 'mdi-cog' },
  { route: 'AdminSMTPConfig', title: '邮件配置', icon: 'mdi-email' },
]

const bottomNavItems = [
  { route: 'AdminIndex', title: '概览', icon: 'mdi-view-dashboard' },
  { route: 'AdminUsers', title: '用户', icon: 'mdi-account-group' },
  { route: 'AdminArticles', title: '文章', icon: 'mdi-file-document' },
  { route: 'AdminDeletions', title: '删除', icon: 'mdi-delete' },
]

const breadcrumbs = computed(() => {
  const crumbs = [{ title: '首页', to: { name: 'AdminIndex' }, disabled: false }]
  if (route.name !== 'AdminIndex') {
    const pageTitle = pageTitles[route.name] || route.name
    crumbs.push({ title: pageTitle, disabled: true })
  }
  return crumbs
})

const pageTitles = {
  AdminIndex: '数据概览',
  AdminUsers: '用户管理',
  AdminArticles: '文章管理',
  AdminComments: '评论管理',
  AdminCategories: '分区管理',
  AdminTitles: '头衔管理',
  AdminSidebar: '侧边栏配置',
  AdminDeletions: '删除申请',
  AdminReports: '举报管理',
  AdminAnnouncement: '公告管理',
  AdminUserNotifications: '用户通知与权限',
  AdminSystemLogs: '系统操作日志',
  AdminSiteConfig: '网站配置',
  AdminSMTPConfig: '邮件配置',
}

const checkMobile = () => {
  const wasMobile = isMobile.value
  isMobile.value = window.innerWidth < 1024
  
  if (isMobile.value) {
    // 移动端：默认关闭侧边栏
    drawerOpen.value = false
    sidebarCollapsed.value = true
  } else {
    // 桌面端：恢复侧边栏状态
    const saved = localStorage.getItem('adminSidebarCollapsed')
    if (saved !== null) sidebarCollapsed.value = JSON.parse(saved)
    drawerOpen.value = true
  }
}

const handleLogout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('user')
  router.push('/login')
}

const checkAdmin = async () => {
  try {
    const response = await api.get('/admin/check')
    isAdmin.value = response.data.is_admin
  } catch (error) {
    isAdmin.value = false
  } finally {
    isInitialized.value = true
  }
}

const loadDeletionCount = async () => {
  try {
    const response = await api.get('/deletion-requests')
    deletionCount.value = response.data.requests?.length || 0
  } catch (error) {
    deletionCount.value = 0
  }
}

const loadCurrentUser = () => {
  const storedUser = localStorage.getItem('user')
  if (storedUser) {
    currentUser.value = JSON.parse(storedUser)
  }
}

const loadVersion = async () => {
  try {
    const response = await api.get('/version')
    backendVersion.value = response.data.backend?.version || response.data.backend_version || response.data.version || 'unknown'
  } catch (error) {
    backendVersion.value = '获取失败'
    console.error('加载后端版本失败', error)
  }
}

watch(sidebarCollapsed, (val) => {
  localStorage.setItem('adminSidebarCollapsed', val)
})

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  checkAdmin()
  loadDeletionCount()
  loadCurrentUser()
  loadVersion()
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})

watch(() => route.path, loadDeletionCount)
</script>

<style scoped>
.admin-drawer {
  border: none;
  box-shadow: 4px 0 20px rgba(0, 0, 0, 0.08);
}

.drawer-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  min-height: 64px;
}

.collapse-btn {
  opacity: 0.8;
  transition: opacity 0.2s;
}

.collapse-btn:hover {
  opacity: 1;
}

.nav-item {
  transition: all 0.2s;
  margin-bottom: 4px;
}

.nav-item:hover {
  background: rgba(255, 255, 255, 0.1);
}

.nav-item-active {
  background: rgba(255, 255, 255, 0.2) !important;
}

.admin-header {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.admin-main {
  background: #f5f5f5;
  min-height: 100vh;
}

.version-info {
  font-size: 0.75rem;
  opacity: 0.9;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  margin: 8px;
}

.admin-bottom-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
}

/* 过渡动画 */
.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.3s ease;
}

.fade-slide-enter-from {
  opacity: 0;
  transform: translateX(20px);
}

.fade-slide-leave-to {
  opacity: 0;
  transform: translateX(-20px);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

@media (max-width: 1024px) {
  .admin-main {
    padding-bottom: 56px !important;
  }
}
</style>
