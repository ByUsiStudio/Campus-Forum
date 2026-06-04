<template>
  <v-app v-if="isInitialized && isAdmin" class="admin-page">
    <v-navigation-drawer
      v-model="drawerOpen"
      :rail="sidebarCollapsed"
      :permanent="!isMobile"
      :temporary="isMobile"
      :width="240"
      class="admin-sidebar"
    >
      <div class="sidebar-inner">
        <div class="sidebar-header" @click="toggleSidebar()">
          <v-icon size="24" color="primary">mdi-shield-crown</v-icon>
          <span v-if="!sidebarCollapsed" class="header-title">管理后台</span>
          <v-spacer />
          <v-btn
            icon
            variant="text"
            size="x-small"
            @click.stop="toggleSidebar"
          >
            <v-icon size="18">
              {{ sidebarCollapsed ? 'mdi-chevron-right' : 'mdi-chevron-left' }}
            </v-icon>
          </v-btn>
        </div>

        <v-divider class="sidebar-divider" />

        <nav class="sidebar-nav">
          <div v-if="sidebarItems.length > 0">
            <router-link
              v-for="item in sidebarItems"
              :key="item.link"
              :to="item.link"
              class="nav-item"
              :class="{ 'active': isActive(item.link) }"
              @click="drawerOpen = false"
            >
              <v-icon size="22" class="nav-icon">{{ getIcon(item.icon) }}</v-icon>
              <span v-if="!sidebarCollapsed" class="nav-text">{{ item.title }}</span>
            </router-link>
          </div>

          <v-divider class="sidebar-divider mt-4" />

          <div class="admin-section-title" v-if="!sidebarCollapsed">
            管理功能
          </div>

          <router-link
            v-for="item in adminItems"
            :key="item.route"
            :to="{ name: item.route }"
            class="nav-item"
            :class="{ 'active': route.name === item.route }"
            @click="drawerOpen = false"
          >
            <v-icon size="22" class="nav-icon">{{ item.icon }}</v-icon>
            <span v-if="!sidebarCollapsed" class="nav-text">{{ item.title }}</span>
            <v-chip
              v-if="item.badge && item.badge() > 0 && !sidebarCollapsed"
              size="x-small"
              color="error"
              class="nav-badge"
            >
              {{ item.badge() }}
            </v-chip>
          </router-link>
        </nav>

        <div class="sidebar-footer">
          <v-divider class="sidebar-divider" />
          <div class="footer-info" v-if="!sidebarCollapsed">
            <div class="text-caption text-medium-emphasis">管理员面板</div>
          </div>
        </div>
      </div>
    </v-navigation-drawer>

    <v-app-bar flat class="admin-header">
      <div class="header-content">
        <div class="header-brand">
          <v-icon size="24" color="primary">mdi-shield-crown</v-icon>
          <div class="brand-text">
            <div class="brand-title">{{ currentPageTitle }}</div>
          </div>
        </div>
        <v-spacer />
        <div class="header-actions">
          <v-btn icon variant="text" size="small" @click="goToHome" title="返回首页" class="mr-2">
            <v-icon>mdi-home</v-icon>
          </v-btn>
          <v-btn icon variant="text" size="small" @click="handleRefresh" :class="{ 'spin': isRefreshing }" title="刷新">
            <v-icon>mdi-refresh</v-icon>
          </v-btn>
        </div>
      </div>
    </v-app-bar>

    <v-main class="admin-main">
      <router-view />
    </v-main>

    <!-- 底部快速操作栏 -->
    <v-footer class="admin-footer" app>
      <div class="footer-content">
        <div class="footer-info">
          <span class="text-caption text-medium-emphasis">{{ currentPageTitle }}</span>
        </div>
        <v-spacer />
        <div class="footer-actions">
          <v-chip
            v-if="deletionCount > 0"
            color="error"
            size="small"
            class="mr-2"
          >
            <v-icon start size="16">mdi-trash-can</v-icon>
            {{ deletionCount }} 个删除请求
          </v-chip>
          <v-tooltip text="刷新数据">
            <template #activator="{ props }">
              <v-btn
                v-bind="props"
                icon="mdi-refresh"
                variant="text"
                size="small"
                @click="handleRefresh"
                :loading="isRefreshing"
              />
            </template>
          </v-tooltip>
          <v-divider vertical class="mx-2" />
          <v-tooltip text="返回首页">
            <template #activator="{ props }">
              <v-btn
                v-bind="props"
                icon="mdi-home"
                variant="text"
                size="small"
                @click="goToHome"
              />
            </template>
          </v-tooltip>
        </div>
      </div>
    </v-footer>

    <div class="floating-ball" @click="toggleDrawer">
      <v-btn
        icon
        size="large"
        color="primary"
        elevation="4"
        class="floating-btn"
      >
        <v-icon size="28">{{ drawerOpen ? 'mdi-close' : 'mdi-menu' }}</v-icon>
      </v-btn>
    </div>

    <div
      v-if="drawerOpen"
      class="drawer-overlay"
      @click="drawerOpen = false"
    />
  </v-app>

  <div v-else-if="isInitialized && !isAdmin" class="error-page">
    <v-card class="error-card">
      <v-icon size="64" color="error" class="mb-4">mdi-shield-alert</v-icon>
      <h2 class="text-h5 mb-2">权限不足</h2>
      <p class="text-body-2 text-medium-emphasis mb-4">您没有访问管理后台的权限</p>
      <v-btn color="primary" variant="flat" @click="goToHome">返回首页</v-btn>
    </v-card>
  </div>

  <div v-else class="loading-page">
    <v-progress-circular indeterminate color="primary" size="48" />
    <div class="mt-4 text-body-2 text-medium-emphasis">加载中...</div>
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

const adminItems = [
  { route: 'AdminIndex', title: '概览', icon: 'mdi-view-dashboard' },
  { route: 'AdminUsers', title: '用户管理', icon: 'mdi-account-group' },
  { route: 'AdminArticles', title: '文章管理', icon: 'mdi-file-document-edit' },
  { route: 'AdminComments', title: '评论管理', icon: 'mdi-comment-text-multiple' },
  { route: 'AdminCategories', title: '分区管理', icon: 'mdi-shape' },
  { route: 'AdminTitles', title: '头衔管理', icon: 'mdi-medal' },
  { route: 'AdminSidebar', title: '侧边栏配置', icon: 'mdi-web' },
  { route: 'AdminDeletions', title: '删除申请', icon: 'mdi-delete-forever', badge: () => deletionCount.value },
  { route: 'AdminAnnouncement', title: '公告管理', icon: 'mdi-bullhorn' },
  { route: 'AdminSiteConfig', title: '网站配置', icon: 'mdi-globe' },
  { route: 'AdminSMTPConfig', title: '邮件配置', icon: 'mdi-email-settings' },
  { route: 'AdminNotifications', title: '通知管理', icon: 'mdi-bell' },
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
    const response = await api.get('/sidebar-config')
    sidebarItems.value = response.data.items || []
  } catch (error) {
    console.error('加载侧边栏配置失败', error)
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

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  loadSidebarConfig()
  checkAdmin()
  loadDeletionCount()
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
  background: #fff;
  border-right: 1px solid #E5E5E5;
  height: calc(100vh - 64px);
  top: 64px;
  z-index: 1000 !important;
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
  gap: 12px;
  padding: 16px;
  cursor: pointer;
  transition: background 0.2s;
}

.sidebar-header:hover {
  background: #F5F5F5;
}

.header-title {
  font-size: 16px;
  font-weight: 600;
  color: #1A1A1A;
  flex: 1;
}

.sidebar-divider {
  margin: 0;
  border-color: #F0F0F0;
}

.sidebar-nav {
  flex: 1;
  padding: 8px;
  overflow-y: auto;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  border-radius: 8px;
  color: #666;
  text-decoration: none;
  transition: all 0.2s;
  margin-bottom: 2px;
  position: relative;
}

.nav-item:hover {
  background: #F5F5F5;
  color: #333;
}

.nav-item.active {
  background: #EEF2FF;
  color: #6750A4;
}

.nav-item.active .nav-icon {
  color: #6750A4;
}

.nav-icon {
  flex-shrink: 0;
}

.nav-text {
  font-size: 14px;
  font-weight: 500;
  white-space: nowrap;
}

.nav-badge {
  margin-left: auto;
}

.admin-section-title {
  font-size: 11px;
  font-weight: 600;
  color: #999;
  text-transform: uppercase;
  padding: 12px 12px 4px;
  letter-spacing: 0.5px;
}

.sidebar-footer {
  padding: 8px;
}

.admin-footer {
  background: #fff;
  border-top: 1px solid #E5E5E5;
  padding: 8px 16px;
  position: relative;
  z-index: 100;
}

.footer-content {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
}

.footer-info {
  flex: 0 0 auto;
}

.admin-header {
  background: #fff;
  border-bottom: 1px solid #E5E5E5;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 900;
}

.header-content {
  display: flex;
  align-items: center;
  width: 100%;
  padding: 0 24px;
}

.header-brand {
  display: flex;
  align-items: center;
  gap: 12px;
}

.brand-title {
  font-size: 18px;
  font-weight: 600;
  color: #1A1A1A;
}

.admin-main {
  background: #FAFAFA;
  min-height: calc(100vh - 128px);
  padding-top: 64px;
  padding-bottom: 56px;
}

.admin-main > div {
  margin: 0;
  padding: 0;
}

.loading-page,
.error-page {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: #FAFAFA;
}

.error-card {
  padding: 48px;
  text-align: center;
  border-radius: 16px;
}

.spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
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
  transition: transform 0.3s ease;
}

.floating-btn:hover {
  transform: scale(1.1);
}

.drawer-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 999;
}

@media (max-width: 960px) {
  .admin-sidebar {
    position: fixed;
    height: 100vh;
    top: 0;
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
    font-size: 24px;
  }

  .nav-text {
    font-size: 15px;
  }
}
</style>