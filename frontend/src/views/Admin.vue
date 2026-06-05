<template>
  <v-app v-if="isInitialized && isAdmin" class="admin-page">
    <!-- 左侧导航栏 -->
    <v-navigation-drawer
      v-model="drawerOpen"
      :rail="sidebarCollapsed"
      :permanent="!isMobile"
      :temporary="isMobile"
      :width="260"
      class="admin-sidebar"
      color="surface"
      border
    >
      <div class="sidebar-inner">
        <!-- 头部品牌区域 -->
        <div class="sidebar-header">
          <div class="brand-container">
            <v-icon size="28" color="primary" class="brand-icon">mdi-shield-crown</v-icon>
            <span v-if="!sidebarCollapsed" class="brand-title">管理后台</span>
          </div>
          <v-btn
            v-if="!isMobile"
            icon
            variant="text"
            size="x-small"
            @click="toggleSidebar"
            class="collapse-btn"
          >
            <v-icon size="18">
              {{ sidebarCollapsed ? 'mdi-chevron-right' : 'mdi-chevron-left' }}
            </v-icon>
          </v-btn>
        </div>

        <v-divider class="sidebar-divider" />

        <!-- 导航菜单 -->
        <nav class="sidebar-nav">
          <!-- 动态链接项 -->
          <div v-if="sidebarItems.length > 0" class="nav-section">
            <div v-if="!sidebarCollapsed" class="section-title">快捷链接</div>
            <router-link
              v-for="item in sidebarItems"
              :key="item.link"
              :to="item.link"
              class="nav-item"
              :class="{ 'active': isActive(item.link) }"
              @click="drawerOpen = false"
            >
              <v-icon size="20" class="nav-icon">{{ getIcon(item.icon) }}</v-icon>
              <span v-if="!sidebarCollapsed" class="nav-text">{{ item.title }}</span>
            </router-link>
          </div>

          <!-- 管理功能 -->
          <div class="nav-section mt-4">
            <div v-if="!sidebarCollapsed" class="section-title">管理功能</div>
            
            <router-link
              v-for="item in adminItems"
              :key="item.route"
              :to="{ name: item.route }"
              class="nav-item"
              :class="{ 'active': route.name === item.route }"
              @click="drawerOpen = false"
            >
              <v-icon size="20" class="nav-icon">{{ item.icon }}</v-icon>
              <span v-if="!sidebarCollapsed" class="nav-text">{{ item.title }}</span>
              <v-chip
                v-if="item.badge && item.badge() > 0 && !sidebarCollapsed"
                size="x-small"
                color="error"
                class="nav-badge"
                variant="tonal"
              >
                {{ item.badge() }}
              </v-chip>
            </router-link>
          </div>
        </nav>

        <!-- 底部信息 -->
        <div class="sidebar-footer">
          <v-divider class="sidebar-divider" />
          <div v-if="!sidebarCollapsed" class="footer-info">
            <div class="version-info">
              <v-chip size="x-small" variant="tonal" color="primary">
                v{{ version || '1.0' }}
              </v-chip>
            </div>
          </div>
        </div>
      </div>
    </v-navigation-drawer>

    <!-- 顶部工具栏 -->
    <v-app-bar flat class="admin-header" color="surface" border>
      <div class="header-content">
        <div class="header-brand">
          <v-breadcrumbs :items="breadcrumbs" divider="/">
            <template #divider>
              <v-icon size="16">mdi-chevron-right</v-icon>
            </template>
          </v-breadcrumbs>
        </div>
        
        <v-spacer />
        
        <div class="header-actions">
          <v-tooltip text="返回首页" location="bottom">
            <template #activator="{ props }">
              <v-btn
                v-bind="props"
                icon
                variant="text"
                size="small"
                @click="goToHome"
                class="action-btn"
              >
                <v-icon size="20">mdi-home-outline</v-icon>
              </v-btn>
            </template>
          </v-tooltip>
          
          <v-tooltip text="刷新页面" location="bottom">
            <template #activator="{ props }">
              <v-btn
                v-bind="props"
                icon
                variant="text"
                size="small"
                @click="handleRefresh"
                :class="{ 'refreshing': isRefreshing }"
                class="action-btn"
              >
                <v-icon size="20">mdi-refresh</v-icon>
              </v-btn>
            </template>
          </v-tooltip>
        </div>
      </div>
    </v-app-bar>

    <!-- 主要内容区域 -->
    <v-main class="admin-main">
      <div class="main-container">
        <router-view />
      </div>
    </v-main>

    <!-- 移动端浮动按钮 -->
    <div v-if="isMobile" class="floating-ball" @click="toggleDrawer">
      <v-btn
        icon
        size="large"
        color="primary"
        elevation="3"
        class="floating-btn"
      >
        <v-icon size="24">{{ drawerOpen ? 'mdi-close' : 'mdi-menu' }}</v-icon>
      </v-btn>
    </div>

    <!-- 移动端遮罩层 -->
    <div
      v-if="drawerOpen && isMobile"
      class="drawer-overlay"
      @click="drawerOpen = false"
    />
  </v-app>

  <!-- 权限不足提示 -->
  <div v-else-if="isInitialized && !isAdmin" class="error-page">
    <v-card class="error-card" elevation="2">
      <div class="error-content">
        <v-icon size="72" color="error" class="mb-4">mdi-shield-alert-outline</v-icon>
        <h2 class="text-h5 mb-2 font-weight-bold">权限不足</h2>
        <p class="text-body-1 text-medium-emphasis mb-6">您没有访问管理后台的权限</p>
        <v-btn color="primary" variant="flat" size="large" @click="goToHome" prepend-icon="mdi-home">
          返回首页
        </v-btn>
      </div>
    </v-card>
  </div>

  <!-- 加载状态 -->
  <div v-else class="loading-page">
    <v-progress-circular indeterminate color="primary" size="56" width="4" />
    <div class="mt-6 text-body-1 text-medium-emphasis">加载中...</div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import api from '../api'

const router = useRouter()
const route = useRoute()

const isInitialized = ref(false)
const isAdmin = ref(false)
const isRefreshing = ref(false)
const isMobile = ref(false)
const sidebarItems = ref([])
const deletionCount = ref(0)
const sidebarCollapsed = ref(false)
const drawerOpen = ref(false)
const version = ref('1.0')

const adminItems = [
  { route: 'AdminIndex', title: '数据概览', icon: 'mdi-view-dashboard-outline' },
  { route: 'AdminUsers', title: '用户管理', icon: 'mdi-account-group-outline' },
  { route: 'AdminArticles', title: '文章管理', icon: 'mdi-file-document-edit-outline' },
  { route: 'AdminComments', title: '评论管理', icon: 'mdi-comment-text-multiple-outline' },
  { route: 'AdminCategories', title: '分区管理', icon: 'mdi-shape-outline' },
  { route: 'AdminTitles', title: '头衔管理', icon: 'mdi-medal-outline' },
  { route: 'AdminSidebar', title: '侧边栏配置', icon: 'mdi-web' },
  { route: 'AdminDeletions', title: '删除申请', icon: 'mdi-delete-forever', badge: () => deletionCount.value },
  { route: 'AdminAnnouncement', title: '公告管理', icon: 'mdi-bullhorn-outline' },
  { route: 'AdminSiteConfig', title: '网站配置', icon: 'mdi-cog-outline' },
  { route: 'AdminSMTPConfig', title: '邮件配置', icon: 'mdi-email-settings-outline' },
  { route: 'AdminNotifications', title: '通知管理', icon: 'mdi-bell-outline' },
]

const pageTitles = {
  AdminIndex: '数据概览',
  AdminUsers: '用户管理',
  AdminArticles: '文章管理',
  AdminComments: '评论管理',
  AdminCategories: '分区管理',
  AdminTitles: '头衔管理',
  AdminSidebar: '侧边栏配置',
  AdminDeletions: '删除申请',
  AdminAnnouncement: '公告管理',
  AdminSiteConfig: '网站配置',
  AdminSMTPConfig: '邮件配置',
  AdminNotifications: '通知管理',
}

const currentPageTitle = computed(() => pageTitles[route.name] || '概览')

const breadcrumbs = computed(() => [
  { title: '首页', disabled: false, href: '/' },
  { title: '管理后台', disabled: true },
  { title: currentPageTitle.value, disabled: true }
])

const getIcon = (icon) => {
  if (!icon) return 'mdi-link'
  if (icon.match(/[\u{1F600}-\u{1F64F}]/u)) return icon
  return icon.startsWith('mdi-') ? icon : `mdi-${icon}`
}

const isActive = (link) => {
  return route.path === link
}

const toggleSidebar = () => {
  sidebarCollapsed.value = !sidebarCollapsed.value
  localStorage.setItem('adminSidebarCollapsed', sidebarCollapsed.value)
}

const toggleDrawer = () => {
  drawerOpen.value = !drawerOpen.value
}

const checkMobile = () => {
  isMobile.value = window.innerWidth < 960
}

const goToHome = () => {
  router.push('/')
}

const handleRefresh = () => {
  isRefreshing.value = true
  setTimeout(() => {
    window.location.reload()
  }, 800)
}

const loadSidebarConfig = async () => {
  try {
    console.log('开始加载管理员侧边栏配置...')
    const response = await api.get('/sidebar-config')
    console.log('管理员侧边栏配置响应:', response.data)
    
    if (response.data && response.data.items) {
      sidebarItems.value = response.data.items
      console.log('侧边栏项目数量:', sidebarItems.value.length)
    } else {
      console.warn('侧边栏配置数据格式异常:', response.data)
      sidebarItems.value = []
    }
  } catch (error) {
    console.error('加载侧边栏配置失败', error)
    console.error('错误详情:', error.response?.data || error.message)
    sidebarItems.value = []
  }
}

const checkAdmin = async () => {
  try {
    const response = await api.get('/admin/check')
    isAdmin.value = response.data.is_admin
  } catch (error) {
    console.error('检查权限失败', error)
    isAdmin.value = false
  } finally {
    isInitialized.value = true
  }
}

const loadDeletionCount = async () => {
  try {
    const response = await api.get('/deletion-requests')
    deletionCount.value = response.data.requests ? response.data.requests.length : 0
  } catch (error) {
    deletionCount.value = 0
  }
}

const loadVersion = async () => {
  try {
    const response = await api.get('/version')
    version.value = response.data.backend?.version || response.data.version || '1.0'
  } catch (error) {
    console.error('加载版本信息失败', error)
  }
}

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  
  // 恢复侧边栏状态
  const savedState = localStorage.getItem('adminSidebarCollapsed')
  if (savedState !== null) {
    sidebarCollapsed.value = JSON.parse(savedState)
  }
  
  loadSidebarConfig()
  checkAdmin()
  loadDeletionCount()
  loadVersion()
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})

watch(() => route.path, () => {
  loadDeletionCount()
})
</script>

<style scoped>
.admin-sidebar {
  background: #ffffff;
  border-right: 1px solid #e0e0e0;
  height: calc(100vh - 64px);
  top: 64px;
  z-index: 1000 !important;
  transition: all 0.3s ease;
}

.sidebar-inner {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

.sidebar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 16px;
  min-height: 64px;
}

.brand-container {
  display: flex;
  align-items: center;
  gap: 12px;
}

.brand-icon {
  flex-shrink: 0;
}

.brand-title {
  font-size: 18px;
  font-weight: 700;
  color: #1a1a1a;
  white-space: nowrap;
}

.collapse-btn {
  opacity: 0.6;
  transition: opacity 0.2s;
}

.collapse-btn:hover {
  opacity: 1;
}

.sidebar-divider {
  margin: 0;
  border-color: #f0f0f0;
}

.sidebar-nav {
  flex: 1;
  padding: 12px 8px;
  overflow-y: auto;
  overflow-x: hidden;
}

.nav-section {
  margin-bottom: 8px;
}

.section-title {
  font-size: 11px;
  font-weight: 600;
  color: #999;
  text-transform: uppercase;
  padding: 8px 12px 4px;
  letter-spacing: 0.5px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  border-radius: 8px;
  color: #666;
  text-decoration: none;
  transition: all 0.2s ease;
  margin-bottom: 2px;
  position: relative;
  font-size: 14px;
}

.nav-item:hover {
  background: #f5f5f5;
  color: #333;
}

.nav-item.active {
  background: linear-gradient(135deg, #6750A4 0%, #7E6BC4 100%);
  color: #ffffff;
  box-shadow: 0 2px 8px rgba(103, 80, 164, 0.2);
}

.nav-item.active .nav-icon {
  color: #ffffff;
}

.nav-icon {
  flex-shrink: 0;
  transition: transform 0.2s;
}

.nav-item:hover .nav-icon {
  transform: scale(1.1);
}

.nav-text {
  font-weight: 500;
  white-space: nowrap;
  flex: 1;
}

.nav-badge {
  margin-left: auto;
}

.sidebar-footer {
  padding: 8px;
}

.footer-info {
  display: flex;
  justify-content: center;
  padding: 8px 0;
}

.version-info {
  text-align: center;
}

.admin-header {
  background: #ffffff;
  border-bottom: 1px solid #e0e0e0;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 900;
  height: 64px;
}

.header-content {
  display: flex;
  align-items: center;
  width: 100%;
  padding: 0 24px;
  height: 100%;
}

.header-brand {
  display: flex;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  transition: all 0.2s;
}

.action-btn:hover {
  background: #f5f5f5;
}

.refreshing {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.admin-main {
  background: #fafafa;
  min-height: calc(100vh - 64px);
  padding-top: 64px;
}

.main-container {
  padding: 24px;
  max-width: 100%;
}

.loading-page,
.error-page {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #fafafa 100%);
}

.error-card {
  padding: 48px;
  text-align: center;
  border-radius: 16px;
  max-width: 400px;
}

.error-content {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.floating-ball {
  position: fixed;
  bottom: 24px;
  right: 24px;
  z-index: 2000;
}

.floating-btn {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(103, 80, 164, 0.3);
}

.floating-btn:hover {
  transform: scale(1.1);
  box-shadow: 0 6px 16px rgba(103, 80, 164, 0.4);
}

.drawer-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 999;
  backdrop-filter: blur(2px);
}

@media (max-width: 960px) {
  .admin-sidebar {
    position: fixed;
    height: 100vh;
    top: 0;
  }
  
  .admin-main {
    padding-top: 64px;
  }
  
  .main-container {
    padding: 16px;
  }
}

@media (max-width: 600px) {
  .floating-ball {
    bottom: 20px;
    right: 20px;
  }
  
  .sidebar-nav {
    padding: 8px 4px;
  }
  
  .nav-item {
    padding: 12px 16px;
    border-radius: 12px;
    margin-bottom: 4px;
  }
  
  .nav-icon {
    font-size: 22px;
  }
  
  .nav-text {
    font-size: 15px;
  }
  
  .header-content {
    padding: 0 16px;
  }
  
  .main-container {
    padding: 12px;
  }
}
</style>