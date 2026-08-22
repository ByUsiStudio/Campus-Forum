<template>
  <el-container v-if="isInitialized && isAdmin" class="admin-shell">
    <!-- 顶部管理后台 Header -->
    <header class="admin-header">
      <div class="admin-header-inner">
        <div class="brand">
          <div class="brand-logo">
            <el-icon :size="22"><Medal /></el-icon>
          </div>
          <span class="brand-title">管理后台</span>
        </div>

        <el-breadcrumb separator="/" class="admin-breadcrumb">
          <el-breadcrumb-item
            v-for="(crumb, index) in breadcrumbs"
            :key="index"
            :to="crumb.to"
          >
            {{ crumb.title }}
          </el-breadcrumb-item>
        </el-breadcrumb>

        <div class="flex-spacer" />

        <div class="header-actions">
          <router-link to="/" class="header-link">
            <el-icon :size="16"><HomeFilled /></el-icon>
            <span class="header-link-text">返回站点</span>
          </router-link>

          <el-badge
            :value="notificationCount"
            :max="99"
            type="danger"
            class="notification-badge"
          >
            <el-button text circle class="icon-btn" @click="router.push({ name: 'AdminNotifications' })">
              <el-icon :size="18"><Bell /></el-icon>
            </el-button>
          </el-badge>

          <el-dropdown trigger="click" class="admin-user-dropdown">
            <div class="user-trigger">
              <el-avatar :size="30" class="user-avatar">
                <el-image
                  v-if="currentUser?.avatar"
                  :src="currentUser.avatar"
                  fit="cover"
                />
                <template v-else>
                  {{ currentUser?.display_name?.[0] || currentUser?.username?.[0] || 'A' }}
                </template>
              </el-avatar>
              <span class="user-name">
                {{ currentUser?.display_name || currentUser?.username }}
              </span>
              <el-icon :size="14" class="user-caret"><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="router.push('/profile')">
                  <el-icon><User /></el-icon>个人资料
                </el-dropdown-item>
                <el-dropdown-item divided @click="handleLogout">
                  <el-icon><SwitchButton /></el-icon>退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>
    </header>

    <!-- 移动端 / 平板：横向可滚动菜单栏 -->
    <nav class="mobile-nav">
      <a
        v-for="item in adminItems"
        :key="item.route"
        class="mobile-nav-item"
        :class="{ 'is-active': activeMenu === '/admin/' + item.path }"
        @click.prevent="router.push('/admin/' + item.path)"
      >
        <el-icon :size="16"><component :is="item.icon" /></el-icon>
        <span class="mobile-nav-title">{{ item.title }}</span>
        <span v-if="item.badge && item.badge() > 0" class="mobile-nav-badge">{{ item.badge() > 99 ? '99+' : item.badge() }}</span>
      </a>
    </nav>

    <el-container class="admin-body">
      <!-- 桌面端左侧导航 -->
      <aside class="admin-aside">
        <div class="aside-inner">
          <div class="menu-label">导航菜单</div>
          <el-menu
            :default-active="activeMenu"
            class="admin-menu"
            router
          >
            <el-menu-item
              v-for="item in adminItems"
              :key="item.route"
              :index="'/admin/' + item.path"
            >
              <el-icon>
                <component :is="item.icon" />
              </el-icon>
              <template #title>
                <span class="menu-title">{{ item.title }}</span>
                <el-badge
                  v-if="item.badge && item.badge() > 0"
                  :value="item.badge()"
                  :max="99"
                  type="danger"
                  class="menu-badge"
                />
              </template>
            </el-menu-item>
          </el-menu>

          <div class="aside-footer">
            <div class="version-block">
              <div class="version-row">
                <span class="version-key">前端版本</span>
                <span class="version-value">{{ frontendVersion }}</span>
              </div>
              <div class="version-row">
                <span class="version-key">后端版本</span>
                <span class="version-value">{{ backendVersion }}</span>
              </div>
            </div>
          </div>
        </div>
      </aside>

      <!-- 主要内容区域 -->
      <el-main class="admin-main">
        <div class="admin-container">
          <transition name="fade-slide" mode="out-in">
            <router-view />
          </transition>
        </div>
      </el-main>
    </el-container>
  </el-container>

  <!-- 权限不足 -->
  <el-container v-else-if="isInitialized && !isAdmin" class="fill-height">
    <el-row justify="center" align="middle" class="fill-height">
      <el-col :xs="22" :sm="16" :md="12" :lg="8">
        <el-card class="text-center forbidden-card">
          <el-icon :size="80" color="#ef4444" class="mb-4"><WarningFilled /></el-icon>
          <div class="forbidden-title">权限不足</div>
          <div class="forbidden-desc">您没有访问管理后台的权限</div>
          <el-button type="primary" round @click="router.push('/')">
            <el-icon class="mr-1"><HomeFilled /></el-icon>返回首页
          </el-button>
        </el-card>
      </el-col>
    </el-row>
  </el-container>

  <!-- 加载状态 -->
  <el-container v-else class="fill-height">
    <el-row justify="center" align="middle" class="fill-height">
      <el-icon class="is-loading" :size="40" color="#4f6ef7"><Loading /></el-icon>
    </el-row>
  </el-container>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import api from '../api'
import {
  ArrowDown,
  Bell,
  Box,
  ChatLineRound,
  Delete,
  Document,
  DocumentCopy,
  Flag,
  Grid,
  HomeFilled,
  Loading,
  Medal,
  Message,
  Monitor,
  Notification,
  Promotion,
  Setting,
  SwitchButton,
  TrendCharts,
  User,
  UserFilled,
  WarningFilled
} from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()

const isInitialized = ref(false)
const isAdmin = ref(false)
const isMobile = ref(false)
const deletionCount = ref(0)
const notificationCount = ref(3)
const currentUser = ref(null)
const backendVersion = ref('加载中...')
const frontendVersion = ref(typeof __FRONTEND_VERSION__ !== 'undefined' ? __FRONTEND_VERSION__ : 'unknown')

const adminItems = [
  { route: 'AdminIndex', path: 'index', title: '数据概览', icon: 'Grid' },
  { route: 'AdminStatistics', path: 'statistics', title: '数据统计', icon: 'TrendCharts' },
  { route: 'AdminUsers', path: 'users', title: '用户管理', icon: 'UserFilled', badge: null },
  { route: 'AdminArticles', path: 'articles', title: '文章管理', icon: 'Document' },
  { route: 'AdminComments', path: 'comments', title: '评论管理', icon: 'ChatLineRound' },
  { route: 'AdminCategories', path: 'categories', title: '分区管理', icon: 'Box' },
  { route: 'AdminTitles', path: 'titles', title: '头衔管理', icon: 'Medal' },
  { route: 'AdminSidebar', path: 'sidebar', title: '侧边栏', icon: 'Monitor' },
  { route: 'AdminDeletions', path: 'deletions', title: '删除申请', icon: 'Delete', badge: () => deletionCount.value },
  { route: 'AdminReports', path: 'reports', title: '举报管理', icon: 'Flag' },
  { route: 'AdminAnnouncement', path: 'announcement', title: '公告管理', icon: 'Promotion' },
  { route: 'AdminNotifications', path: 'admin-notifications', title: '系统通知', icon: 'Notification' },
  { route: 'AdminUserNotifications', path: 'notifications', title: '用户通知', icon: 'Bell' },
  { route: 'AdminSystemLogs', path: 'system-logs', title: '系统日志', icon: 'DocumentCopy' },
  { route: 'AdminSiteConfig', path: 'siteconfig', title: '网站配置', icon: 'Setting' },
  { route: 'AdminSMTPConfig', path: 'smtpconfig', title: '邮件配置', icon: 'Message' },
]

const activeMenu = computed(() => {
  const item = adminItems.find(i => i.route === route.name)
  return item ? '/admin/' + item.path : ''
})

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
  AdminNotifications: '系统通知',
  AdminSystemLogs: '系统操作日志',
  AdminSiteConfig: '网站配置',
  AdminSMTPConfig: '邮件配置',
}

const checkMobile = () => {
  isMobile.value = window.innerWidth < 1024
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
    if (frontendVersion.value === 'unknown' || frontendVersion.value === '') {
      frontendVersion.value = response.data.frontend?.version || 'unknown'
    }
  } catch (error) {
    backendVersion.value = '获取失败'
    console.error('加载后端版本失败', error)
  }
}

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
.admin-shell {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--campus-bg);
}

/* ---------------- 顶部 Header ---------------- */
.admin-header {
  position: sticky;
  top: 0;
  z-index: 100;
  background: var(--campus-surface);
  border-bottom: 1px solid var(--campus-border);
  box-shadow: var(--campus-shadow-sm);
}

.admin-header-inner {
  height: var(--campus-header-h);
  max-width: 1600px;
  margin: 0 auto;
  padding: 0 20px;
  display: flex;
  align-items: center;
  gap: 16px;
}

.brand {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-shrink: 0;
}

.brand-logo {
  width: 36px;
  height: 36px;
  display: grid;
  place-items: center;
  border-radius: 10px;
  color: #fff;
  background: linear-gradient(135deg, var(--campus-primary), var(--campus-primary-light));
  box-shadow: var(--campus-shadow-sm);
}

.brand-title {
  font-size: 18px;
  font-weight: 700;
  color: var(--campus-text);
  letter-spacing: 0.02em;
}

.admin-breadcrumb {
  font-size: 14px;
  flex-shrink: 0;
}

@media (max-width: 575px) {
  .admin-breadcrumb {
    display: none;
  }
  .header-link-text {
    display: none;
  }
}

@media (max-width: 767px) {
  .user-name {
    display: none;
  }
}

.flex-spacer {
  flex: 1;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.header-link {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 7px 12px;
  border-radius: var(--campus-radius-sm);
  color: var(--campus-text-secondary);
  font-size: 14px;
  font-weight: 500;
  transition: var(--campus-transition);
}

.header-link:hover {
  background: var(--campus-primary-soft);
  color: var(--campus-primary);
}

.icon-btn {
  color: var(--campus-text-secondary);
}

.admin-user-dropdown {
  cursor: pointer;
}

.user-trigger {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 5px 10px;
  border-radius: 999px;
  border: 1px solid var(--campus-border);
  background: var(--campus-surface);
  transition: var(--campus-transition);
}

.user-trigger:hover {
  border-color: var(--campus-primary);
  background: var(--campus-surface-2);
}

.user-avatar {
  background: var(--campus-primary);
}

.user-name {
  color: var(--campus-text);
  font-size: 14px;
  max-width: 140px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.user-caret {
  color: var(--campus-text-muted);
}

/* ---------------- 移动端横向滚动菜单 ---------------- */
.mobile-nav {
  position: sticky;
  top: var(--campus-header-h);
  z-index: 90;
  display: flex;
  align-items: center;
  gap: 6px;
  overflow-x: auto;
  padding: 8px 12px;
  background: var(--campus-surface);
  border-bottom: 1px solid var(--campus-border);
  scrollbar-width: none;
  -webkit-overflow-scrolling: touch;
}

.mobile-nav::-webkit-scrollbar {
  display: none;
}

.mobile-nav-item {
  flex-shrink: 0;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  border-radius: 999px;
  color: var(--campus-text-secondary);
  font-size: 13px;
  font-weight: 500;
  white-space: nowrap;
  cursor: pointer;
  border: 1px solid transparent;
  transition: var(--campus-transition);
}

.mobile-nav-item:hover {
  color: var(--campus-primary);
  background: var(--campus-surface-2);
}

.mobile-nav-item.is-active {
  color: var(--campus-primary);
  background: var(--campus-primary-soft);
  border-color: rgba(79, 110, 247, 0.18);
}

.mobile-nav-badge {
  min-width: 16px;
  height: 16px;
  padding: 0 4px;
  display: inline-grid;
  place-items: center;
  border-radius: 999px;
  background: var(--campus-danger);
  color: #fff;
  font-size: 10px;
  font-weight: 700;
  line-height: 1;
}

/* ---------------- 主体布局 ---------------- */
.admin-body {
  flex: 1;
  min-height: 0;
  max-width: 1600px;
  margin: 0 auto;
  width: 100%;
}

/* ======== 响应式 ---------------- */
@media (min-width: 992px) {
  .mobile-nav {
    display: none !important;
  }
}

@media (max-width: 991px) {
  .admin-aside {
    display: none !important;
  }
  .mobile-nav {
    display: flex;
  }
  .admin-main {
    padding: 16px;
  }
}

@media (max-width: 575px) {
  .admin-main {
    padding: 12px;
  }
  .admin-header-inner {
    padding: 0 12px;
    gap: 8px;
  }
}

.admin-aside {
  width: 260px;
  flex-shrink: 0;
  display: flex;
  border-right: 1px solid var(--campus-border);
}

.aside-inner {
  position: sticky;
  top: calc(var(--campus-header-h) + 1px);
  height: calc(100vh - var(--campus-header-h) - 1px);
  display: flex;
  flex-direction: column;
  padding: 16px 12px;
  box-sizing: border-box;
  background: var(--campus-surface);
}

.menu-label {
  padding: 0 12px 8px;
  font-size: 12px;
  font-weight: 600;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  color: var(--campus-text-muted);
}

.admin-menu {
  flex: 1;
  overflow-y: auto;
  border-right: none;
  background: transparent;
}

.admin-menu .el-menu-item {
  height: 44px;
  line-height: 44px;
  margin: 2px 0;
  border-radius: 10px;
  color: var(--campus-text-secondary);
  transition: var(--campus-transition);
}

.admin-menu .el-menu-item:hover {
  background: var(--campus-surface-2);
  color: var(--campus-primary);
}

.admin-menu .el-menu-item.is-active {
  background: var(--campus-primary-soft);
  color: var(--campus-primary);
  font-weight: 600;
}

.admin-menu .el-menu-item .el-icon {
  color: inherit;
}

.menu-title {
  font-weight: 500;
}

.menu-badge {
  margin-left: 8px;
  vertical-align: middle;
}

.aside-footer {
  flex-shrink: 0;
  padding-top: 12px;
  margin-top: 12px;
  border-top: 1px solid var(--campus-border);
}

.version-block {
  padding: 12px;
  border-radius: 10px;
  background: var(--campus-surface-2);
}

.version-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 12px;
  padding: 2px 0;
}

.version-key {
  color: var(--campus-text-muted);
}

.version-value {
  color: var(--campus-text-secondary);
  font-weight: 600;
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* ---------------- 主内容 ---------------- */
.admin-main {
  min-height: calc(100vh - var(--campus-header-h));
  background: var(--campus-bg);
  padding: 24px;
}

.admin-container {
  max-width: 1200px;
  margin: 0 auto;
}

/* ---------------- 权限不足 / 加载 ---------------- */
.fill-height {
  min-height: 100vh;
}

.forbidden-card {
  padding: 32px 0;
  border-radius: var(--campus-radius-lg);
}

.text-center {
  text-align: center;
}

.forbidden-title {
  font-size: 20px;
  font-weight: 700;
  margin-bottom: 8px;
  color: var(--campus-text);
}

.forbidden-desc {
  font-size: 14px;
  color: var(--campus-text-secondary);
  margin-bottom: 24px;
}

.mb-4 {
  margin-bottom: 16px;
}

.mr-1 {
  margin-right: 4px;
}

/* ---------------- 过渡动画 ---------------- */
.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.3s ease;
}

.fade-slide-enter-from {
  opacity: 0;
  transform: translateY(12px);
}

.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>
