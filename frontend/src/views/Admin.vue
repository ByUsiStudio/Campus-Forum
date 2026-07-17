<template>
  <div v-if="isInitialized && isAdmin" class="admin-page">
    <div class="admin-layout">
      <aside :class="['admin-sidebar', { collapsed: sidebarCollapsed }, { 'temporary': isMobile }]" :style="{ width: sidebarCollapsed && !isMobile ? '64px' : '260px' }">
        <div class="sidebar-header">
          <div class="logo-container">
            <i class="fa-solid fa-shield-halved text-white" style="font-size: 32px;"></i>
            <transition name="fade">
              <span v-if="!sidebarCollapsed" class="logo-text">校园论坛</span>
            </transition>
          </div>
          <button 
            class="collapse-btn" 
            @click="sidebarCollapsed = !sidebarCollapsed"
            :class="{ 'rotated': !sidebarCollapsed }"
          >
            <i class="fa-solid fa-chevron-left text-white"></i>
          </button>
        </div>

        <div class="sidebar-divider"></div>

        <nav class="sidebar-nav">
          <router-link
            v-for="item in adminItems"
            :key="item.route"
            :to="{ name: item.route }"
            :class="['nav-item', { active: route.name === item.route }]"
          >
            <i :class="getIconClass(item.icon)"></i>
            <span v-if="!sidebarCollapsed" class="nav-text">{{ item.title }}</span>
            <span v-if="item.badge && item.badge() > 0" class="nav-badge">{{ item.badge() }}</span>
          </router-link>
        </nav>

        <div class="sidebar-footer">
          <div class="sidebar-divider"></div>
          <div v-if="!sidebarCollapsed" class="version-info">
            <div class="d-flex align-center mb-1">
              <i class="fa-solid fa-info text-white opacity-70" style="font-size: 14px;"></i>
              <span class="text-caption text-white text-truncate">版本信息</span>
            </div>
            <div class="text-caption text-white opacity-50">前端: {{ frontendVersion }}</div>
            <div class="text-caption text-white opacity-50">后端: {{ backendVersion }}</div>
          </div>
          <router-link to="/" class="nav-item">
            <i class="fa-solid fa-house"></i>
            <span v-if="!sidebarCollapsed" class="nav-text">返回首页</span>
          </router-link>
        </div>
      </aside>

      <div class="admin-main-content">
        <header class="admin-header">
          <button v-if="isMobile" class="mobile-menu-btn" @click="drawerOpen = !drawerOpen">
            <i class="fa-solid fa-bars"></i>
          </button>

          <div class="breadcrumbs">
            <a :to="{ name: 'AdminIndex' }">首页</a>
            <i class="fa-solid fa-chevron-right"></i>
            <span>{{ pageTitles[route.name] || route.name }}</span>
          </div>

          <div class="header-right">
            <div class="search-box">
              <i class="fa-solid fa-magnifying-glass"></i>
              <input 
                v-model="searchQuery" 
                type="text" 
                placeholder="搜索..."
                class="search-input"
              />
            </div>

            <button class="notification-btn relative">
              <i class="fa-regular fa-bell"></i>
              <span v-if="notificationCount > 0" class="notification-badge">{{ notificationCount }}</span>
            </button>

            <div class="user-menu">
              <button class="user-btn" @click="userMenuOpen = !userMenuOpen">
                <div class="user-avatar">
                  <img v-if="currentUser?.avatar" :src="currentUser.avatar" />
                  <span v-else>{{ currentUser?.display_name?.[0] || currentUser?.username?.[0] || 'A' }}</span>
                </div>
                <span v-if="!isMobile" class="user-name">{{ currentUser?.display_name || currentUser?.username }}</span>
                <i class="fa-solid fa-chevron-down" :class="{ 'rotated': userMenuOpen }"></i>
              </button>
              
              <div v-if="userMenuOpen" class="user-dropdown">
                <router-link to="/profile" @click="userMenuOpen = false">
                  <i class="fa-solid fa-user"></i>
                  <span>个人资料</span>
                </router-link>
                <div class="dropdown-divider"></div>
                <button @click="handleLogout">
                  <i class="fa-solid fa-right-from-bracket"></i>
                  <span>退出登录</span>
                </button>
              </div>
            </div>
          </div>
        </header>

        <main class="admin-content">
          <div class="content-wrapper">
            <router-view />
          </div>
        </main>

        <nav v-if="isMobile" class="admin-bottom-nav">
          <router-link
            v-for="item in bottomNavItems"
            :key="item.route"
            :to="{ name: item.route }"
            :class="['bottom-nav-item', { active: route.name === item.route }]"
          >
            <i :class="getIconClass(item.icon)"></i>
            <span>{{ item.title }}</span>
          </router-link>
        </nav>
      </div>
    </div>
  </div>

  <div v-else-if="isInitialized && !isAdmin" class="access-denied">
    <div class="denied-card">
      <i class="fa-solid fa-shield-exclamation" style="font-size: 80px; color: #FF4D4F;"></i>
      <div class="denied-title">权限不足</div>
      <div class="denied-desc">您没有访问管理后台的权限</div>
      <router-link to="/" class="layui-btn layui-btn-normal">
        <i class="fa-solid fa-house mr-2"></i>返回首页
      </router-link>
    </div>
  </div>

  <div v-else class="loading-state">
    <div class="loading-spinner"></div>
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
const isMobile = ref(false)
const deletionCount = ref(0)
const notificationCount = ref(3)
const sidebarCollapsed = ref(false)
const drawerOpen = ref(true)
const searchQuery = ref('')
const currentUser = ref(null)
const backendVersion = ref('加载中...')
const frontendVersion = ref(typeof __FRONTEND_VERSION__ !== 'undefined' ? __FRONTEND_VERSION__ : 'unknown')
const userMenuOpen = ref(false)

const mdiToFa = {
  'mdi-shield-crown': 'fa-solid fa-shield-halved',
  'mdi-chevron-right': 'fa-solid fa-chevron-right',
  'mdi-chevron-left': 'fa-solid fa-chevron-left',
  'mdi-view-dashboard': 'fa-solid fa-gauge',
  'mdi-chart-box': 'fa-solid fa-chart-column',
  'mdi-account-group': 'fa-solid fa-users',
  'mdi-file-document': 'fa-solid fa-file-lines',
  'mdi-comment-text': 'fa-solid fa-comment',
  'mdi-shape': 'fa-solid fa-shapes',
  'mdi-medal': 'fa-solid fa-medal',
  'mdi-web': 'fa-solid fa-globe',
  'mdi-delete': 'fa-solid fa-trash',
  'mdi-flag': 'fa-solid fa-flag',
  'mdi-bullhorn': 'fa-solid fa-bullhorn',
  'mdi-bell': 'fa-solid fa-bell',
  'mdi-text-box-multiple': 'fa-solid fa-file-lines',
  'mdi-cog': 'fa-solid fa-gear',
  'mdi-email': 'fa-solid fa-envelope',
  'mdi-information': 'fa-solid fa-info',
  'mdi-home': 'fa-solid fa-house',
  'mdi-magnify': 'fa-solid fa-magnifying-glass',
  'mdi-bell-outline': 'fa-regular fa-bell',
  'mdi-account': 'fa-solid fa-user',
  'mdi-logout': 'fa-solid fa-right-from-bracket',
  'mdi-shield-alert': 'fa-solid fa-shield-exclamation'
}

const getIconClass = (icon) => {
  return mdiToFa[icon] || 'fa-solid fa-circle-question'
}

const adminItems = [
  { route: 'AdminIndex', title: '数据概览', icon: 'mdi-view-dashboard' },
  { route: 'AdminStatistics', title: '数据统计', icon: 'mdi-chart-box' },
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

const pageTitles = {
  AdminIndex: '数据概览',
  AdminStatistics: '数据统计',
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
    drawerOpen.value = false
    sidebarCollapsed.value = true
  } else {
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

watch(() => route.path, loadDeletionCount)

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
</script>

<style scoped>
.admin-page {
  min-height: 100vh;
  background: #f5f5f5;
}

.admin-layout {
  display: flex;
  min-height: 100vh;
}

.admin-sidebar {
  background: var(--primary);
  color: white;
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
  z-index: 1000;
  transition: width 0.3s ease;
  overflow: hidden;
}

.admin-sidebar.collapsed {
  width: 64px;
}

.admin-sidebar.temporary {
  transform: translateX(-100%);
}

.admin-sidebar.temporary.open {
  transform: translateX(0);
}

.sidebar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  min-height: 64px;
}

.logo-container {
  display: flex;
  align-items: center;
}

.logo-text {
  margin-left: 12px;
  font-size: 18px;
  font-weight: 700;
  white-space: nowrap;
}

.collapse-btn {
  background: rgba(255, 255, 255, 0.1);
  border: none;
  border-radius: 6px;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
  
  &:hover {
    background: rgba(255, 255, 255, 0.2);
  }
  
  &.rotated {
    transform: rotate(180deg);
  }
}

.sidebar-divider {
  height: 1px;
  background: rgba(255, 255, 255, 0.1);
  margin: 8px 0;
}

.sidebar-nav {
  padding: 8px;
  flex: 1;
}

.nav-item {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  color: rgba(255, 255, 255, 0.8);
  text-decoration: none;
  position: relative;
  
  &:hover {
    background: rgba(255, 255, 255, 0.1);
    color: white;
  }
  
  &.active {
    background: rgba(255, 255, 255, 0.2);
    color: white;
  }
  
  i {
    font-size: 20px;
    flex-shrink: 0;
  }
  
  .nav-text {
    margin-left: 12px;
    font-size: 14px;
    font-weight: 500;
    white-space: nowrap;
  }
  
  .nav-badge {
    margin-left: auto;
    background: #FF4D4F;
    color: white;
    font-size: 12px;
    padding: 2px 8px;
    border-radius: 10px;
    min-width: 20px;
    text-align: center;
  }
}

.sidebar-footer {
  padding: 8px;
}

.version-info {
  padding: 12px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 8px;
  font-size: 12px;
  margin-bottom: 8px;
}

.d-flex {
  display: flex;
}

.align-center {
  align-items: center;
}

.text-caption {
  font-size: 12px;
}

.text-truncate {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.admin-main-content {
  flex: 1;
  margin-left: 260px;
  display: flex;
  flex-direction: column;
  transition: margin-left 0.3s ease;
}

.admin-sidebar.collapsed ~ .admin-main-content {
  margin-left: 64px;
}

.admin-header {
  background: white;
  padding: 0 24px;
  height: 64px;
  display: flex;
  align-items: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  gap: 16px;
}

.mobile-menu-btn {
  background: none;
  border: none;
  font-size: 24px;
  color: #333;
  cursor: pointer;
  display: none;
}

.breadcrumbs {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  
  a {
    color: var(--primary);
    text-decoration: none;
    
    &:hover {
      text-decoration: underline;
    }
  }
  
  span {
    color: #666;
  }
  
  i {
    font-size: 12px;
    color: #999;
  }
}

.header-right {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 16px;
}

.search-box {
  position: relative;
  display: none;
  
  i {
    position: absolute;
    left: 12px;
    top: 50%;
    transform: translateY(-50%);
    font-size: 14px;
    color: #999;
  }
  
  .search-input {
    padding: 8px 12px 8px 36px;
    border: 1px solid #e8e8e8;
    border-radius: 8px;
    font-size: 14px;
    width: 200px;
    outline: none;
    
    &:focus {
      border-color: var(--primary);
    }
  }
}

@media (min-width: 576px) {
  .search-box {
    display: flex;
  }
}

.notification-btn {
  background: none;
  border: none;
  font-size: 20px;
  color: #666;
  cursor: pointer;
  padding: 8px;
  border-radius: 6px;
  
  &:hover {
    background: #f5f5f5;
    color: var(--primary);
  }
  
  .notification-badge {
    position: absolute;
    top: 4px;
    right: 4px;
    background: #FF4D4F;
    color: white;
    font-size: 10px;
    padding: 1px 5px;
    border-radius: 10px;
    min-width: 18px;
    text-align: center;
  }
}

.user-menu {
  position: relative;
}

.user-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  background: none;
  border: none;
  cursor: pointer;
  padding: 8px;
  border-radius: 8px;
  
  &:hover {
    background: #f5f5f5;
  }
  
  i {
    font-size: 12px;
    transition: transform 0.3s ease;
    
    &.rotated {
      transform: rotate(180deg);
    }
  }
}

.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: var(--primary);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 14px;
  overflow: hidden;
  
  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}

.user-name {
  font-size: 14px;
  color: #333;
}

.user-dropdown {
  position: absolute;
  right: 0;
  top: calc(100% + 8px);
  background: white;
  border-radius: 8px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
  min-width: 200px;
  z-index: 100;
  overflow: hidden;
  
  a, button {
    display: flex;
    align-items: center;
    padding: 12px 16px;
    width: 100%;
    border: none;
    background: none;
    cursor: pointer;
    font-size: 14px;
    color: #333;
    text-decoration: none;
    gap: 12px;
    
    &:hover {
      background: #f5f5f5;
    }
    
    i {
      font-size: 16px;
      color: #666;
    }
  }
}

.dropdown-divider {
  height: 1px;
  background: #f0f0f0;
}

.admin-content {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
}

.content-wrapper {
  max-width: 1400px;
}

.admin-bottom-nav {
  display: none;
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: var(--primary);
  padding: 8px 0;
  gap: 0;
}

@media (max-width: 1024px) {
  .admin-main-content {
    margin-left: 0;
    padding-bottom: 56px;
  }
  
  .admin-bottom-nav {
    display: flex;
    justify-content: space-around;
  }
  
  .mobile-menu-btn {
    display: block;
  }
}

.bottom-nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 8px 16px;
  color: rgba(255, 255, 255, 0.6);
  text-decoration: none;
  font-size: 12px;
  transition: all 0.2s ease;
  
  &:hover {
    color: white;
  }
  
  &.active {
    color: white;
    
    i {
      transform: scale(1.1);
    }
  }
  
  i {
    font-size: 20px;
    margin-bottom: 4px;
    transition: transform 0.2s ease;
  }
}

.access-denied {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: #f5f5f5;
  padding: 24px;
}

.denied-card {
  background: white;
  border-radius: 12px;
  padding: 48px;
  text-align: center;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  max-width: 400px;
}

.denied-title {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin-bottom: 8px;
}

.denied-desc {
  font-size: 14px;
  color: #999;
  margin-bottom: 24px;
}

.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: #f5f5f5;
}

.loading-spinner {
  width: 48px;
  height: 48px;
  border: 4px solid #f0f0f0;
  border-top-color: var(--primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>