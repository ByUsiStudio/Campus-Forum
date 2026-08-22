<template>
  <el-container v-if="isInitialized && isAdmin" class="admin-shell">
  <div v-if="isMobile && drawerOpen" class="drawer-mask" @click="drawerOpen = false"></div>

  <el-aside
    :width="(isMobile ? 260 : (sidebarCollapsed ? 64 : 260)) + 'px'"
    class="admin-aside"
    :class="{ 'is-open': drawerOpen, 'is-mobile': isMobile }"
  >
    <div class="aside-inner">
        <div class="drawer-header">
          <div class="d-flex align-center">
            <el-icon :size="30" class="logo-icon" color="#fff">
              <Medal />
            </el-icon>
            <span
              v-if="!sidebarCollapsed"
              class="logo-title"
            >校园论坛</span>
          </div>
          <el-button
            v-if="!isMobile"
            text
            circle
            class="collapse-btn"
            :aria-label="sidebarCollapsed ? '展开侧边栏' : '收起侧边栏'"
            @click="sidebarCollapsed = !sidebarCollapsed"
          >
            <el-icon :size="16" color="#fff">
              <ArrowLeft v-if="!sidebarCollapsed" />
              <ArrowRight v-else />
            </el-icon>
          </el-button>
        </div>

        <el-divider class="aside-divider" />

        <el-menu
          :collapse="sidebarCollapsed && !isMobile"
          :default-active="activeMenu"
          background-color="transparent"
          text-color="#fff"
          active-text-color="#fff"
          class="admin-menu"
          :collapse-transition="false"
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
          <el-divider class="aside-divider" />
          <div v-if="!sidebarCollapsed || isMobile" class="version-info">
            <div class="d-flex align-center mb-1">
              <el-icon :size="14" class="mr-1"><InfoFilled /></el-icon>
              <span class="version-label">版本信息</span>
            </div>
            <div class="version-row">
              <span class="version-key">前端</span>
              <el-tag size="small" effect="dark" class="version-tag">{{ frontendVersion }}</el-tag>
            </div>
            <div class="version-row">
              <span class="version-key">后端</span>
              <el-tag size="small" effect="dark" class="version-tag backend">{{ backendVersion }}</el-tag>
            </div>
          </div>
          <el-menu
            :collapse="sidebarCollapsed && !isMobile"
            background-color="transparent"
            text-color="#fff"
            active-text-color="#fff"
            class="admin-menu footer-menu"
            :collapse-transition="false"
            router
          >
            <el-menu-item index="/">
              <el-icon>
                <HomeFilled />
              </el-icon>
              <template #title>
                <span class="menu-title">返回首页</span>
              </template>
            </el-menu-item>
          </el-menu>
        </div>
      </div>
    </el-aside>

    <el-container class="admin-body">
      <el-header height="64px" class="admin-header">
        <el-button
          class="d-lg-none menu-toggle"
          text
          @click="drawerOpen = !drawerOpen"
        >
          <el-icon :size="20"><Expand /></el-icon>
        </el-button>

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

        <el-input
          v-model="searchQuery"
          placeholder="搜索..."
          clearable
          class="admin-search d-none d-sm-block"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>

        <el-badge
          :value="notificationCount"
          :max="99"
          type="danger"
          class="notification-badge"
        >
          <el-button text circle class="icon-btn">
            <el-icon :size="18"><Bell /></el-icon>
          </el-button>
        </el-badge>

        <el-dropdown trigger="click" class="admin-user-dropdown">
          <div class="user-trigger">
            <el-avatar :size="32" class="user-avatar">
              <el-image
                v-if="currentUser?.avatar"
                :src="currentUser.avatar"
                fit="cover"
              />
              <template v-else>
                {{ currentUser?.display_name?.[0] || currentUser?.username?.[0] || 'A' }}
              </template>
            </el-avatar>
            <span class="user-name d-none d-sm-inline">
              {{ currentUser?.display_name || currentUser?.username }}
            </span>
            <el-icon :size="14" class="user-caret">
              <ArrowDown />
            </el-icon>
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
      </el-header>

      <!-- 主要内容区域 -->
      <el-main class="admin-main">
        <div class="admin-container">
          <transition name="fade-slide" mode="out-in">
            <router-view />
          </transition>
        </div>
      </el-main>
    </el-container>

    <!-- 移动端底部导航 -->
    <div class="bottom-navigation d-lg-none">
      <el-button
        v-for="item in bottomNavItems"
        :key="item.route"
        text
        class="bottom-nav-item"
        @click="router.push({ name: item.route })"
      >
        <el-icon :size="20">
          <component :is="item.icon" />
        </el-icon>
        <span class="bottom-nav-title">{{ item.title }}</span>
      </el-button>
    </div>
  </el-container>

  <!-- 权限不足 -->
  <el-container v-else-if="isInitialized && !isAdmin" class="fill-height">
    <el-row justify="center" align="middle" class="fill-height">
      <el-col :xs="22" :sm="16" :md="12" :lg="8">
        <el-card class="text-center">
          <el-icon :size="80" color="#f56c6c" class="mb-4">
            <WarningFilled />
          </el-icon>
          <div class="forbidden-title">权限不足</div>
          <div class="forbidden-desc">您没有访问管理后台的权限</div>
          <el-button type="primary" @click="router.push('/')">
            <el-icon class="mr-1"><HomeFilled /></el-icon>返回首页
          </el-button>
        </el-card>
      </el-col>
    </el-row>
  </el-container>

  <!-- 加载状态 -->
  <el-container v-else class="fill-height">
    <el-row justify="center" align="middle" class="fill-height">
      <el-icon class="is-loading" :size="40" color="#409eff">
        <Loading />
      </el-icon>
    </el-row>
  </el-container>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import api from '../api'
import {
  ArrowDown,
  ArrowLeft,
  ArrowRight,
  Bell,
  Box,
  ChatLineRound,
  Delete,
  Document,
  DocumentCopy,
  Expand,
  Flag,
  Grid,
  HomeFilled,
  InfoFilled,
  Loading,
  Medal,
  Message,
  Monitor,
  Promotion,
  Search,
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
const sidebarCollapsed = ref(false)
const drawerOpen = ref(true)
const searchQuery = ref('')
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
  { route: 'AdminUserNotifications', path: 'notifications', title: '用户通知', icon: 'Bell' },
  { route: 'AdminSystemLogs', path: 'system-logs', title: '系统日志', icon: 'DocumentCopy' },
  { route: 'AdminSiteConfig', path: 'siteconfig', title: '网站配置', icon: 'Setting' },
  { route: 'AdminSMTPConfig', path: 'smtpconfig', title: '邮件配置', icon: 'Message' },
]

const bottomNavItems = [
  { route: 'AdminIndex', title: '概览', icon: 'Grid' },
  { route: 'AdminUsers', title: '用户', icon: 'UserFilled' },
  { route: 'AdminArticles', title: '文章', icon: 'Document' },
  { route: 'AdminDeletions', title: '删除', icon: 'Delete' },
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
    // 若前端未注入构建版本，则回退使用后端上报的前端版本
    if (frontendVersion.value === 'unknown' || frontendVersion.value === '') {
      frontendVersion.value = response.data.frontend?.version || 'unknown'
    }
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

// 路由变化时关闭移动端抽屉
watch(() => route.path, () => {
  if (isMobile.value) drawerOpen.value = false
})
</script>

<style scoped>
.admin-shell {
  min-height: 100vh;
}

.admin-aside {
  background: #409eff;
  transition: width 0.2s;
  overflow: hidden;
}

/* 移动端：侧边栏变为离屏抽屉 */
.admin-aside.is-mobile {
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  z-index: 1000;
  height: 100vh;
  transform: translateX(-100%);
  transition: transform 0.28s ease;
  box-shadow: 4px 0 16px rgba(0, 0, 0, 0.2);
}
.admin-aside.is-mobile.is-open {
  transform: translateX(0);
}
.drawer-mask {
  position: fixed;
  inset: 0;
  z-index: 999;
  background: rgba(0, 0, 0, 0.4);
}

.aside-inner {
  display: flex;
  flex-direction: column;
  height: 100vh;
}

.drawer-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  min-height: 64px;
  padding: 0 16px;
}

.logo-icon {
  margin-right: 6px;
}

.logo-title {
  color: #fff;
  font-size: 18px;
  font-weight: 700;
  white-space: nowrap;
}

.collapse-btn {
  opacity: 0.8;
  transition: opacity 0.2s;
}

.collapse-btn:hover {
  opacity: 1;
}

.aside-divider {
  border-color: rgba(255, 255, 255, 0.2);
  margin: 8px 0;
}

.admin-menu {
  border-right: none;
  flex: 1;
  overflow-y: auto;
}

.admin-menu .el-menu-item {
  border-radius: 8px;
  margin: 0 8px 4px;
}

.admin-menu .el-menu-item:hover {
  background: rgba(255, 255, 255, 0.1) !important;
}

.admin-menu .el-menu-item.is-active {
  background: rgba(255, 255, 255, 0.2) !important;
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
}

.version-info {
  font-size: 0.75rem;
  opacity: 0.9;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  margin: 8px;
  padding: 10px 12px;
}

.version-label {
  color: #fff;
  font-size: 0.75rem;
}

.version-text {
  color: #fff;
  opacity: 0.7;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.version-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  margin-top: 6px;
}

.version-key {
  color: #fff;
  opacity: 0.75;
  font-size: 0.75rem;
}

.version-tag {
  font-weight: 600;
  border: none;
}

.version-tag.backend {
  background: #606266;
  color: #fff;
}

.footer-menu {
  flex: none;
}

.admin-body {
  min-width: 0;
}

.admin-header {
  display: flex;
  align-items: center;
  background: #fff;
  border-bottom: 1px solid #ebeef5;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  padding: 0 20px;
}

.menu-toggle {
  margin-right: 12px;
}

.admin-breadcrumb {
  font-size: 14px;
}

.flex-spacer {
  flex: 1;
}

.admin-search {
  max-width: 300px;
  margin: 0 16px;
}

.notification-badge {
  margin-right: 8px;
}

.icon-btn {
  color: #606266;
}

.admin-user-dropdown {
  cursor: pointer;
}

.user-trigger {
  display: flex;
  align-items: center;
  padding: 4px 8px;
  border-radius: 8px;
}

.user-trigger:hover {
  background: #f5f7fa;
}

.user-avatar {
  background: #409eff;
}

.user-name {
  margin-left: 8px;
  color: #303133;
  font-size: 14px;
}

.user-caret {
  margin-left: 4px;
  color: #909399;
}

.admin-main {
  background: #f5f5f5;
  min-height: calc(100vh - 64px);
  padding: 24px;
}

/* 移动端底部导航 */
.bottom-navigation {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  z-index: 100;
  display: flex;
  justify-content: space-around;
  background: #fff;
  border-top: 1px solid #ebeef5;
  padding-bottom: env(safe-area-inset-bottom);
}

.bottom-nav-item {
  display: flex;
  flex-direction: column;
  height: 56px;
  color: #606266;
}

.bottom-nav-title {
  font-size: 12px;
  line-height: 1;
}

.fill-height {
  min-height: 100vh;
}

.text-center {
  text-align: center;
}

.forbidden-title {
  font-size: 20px;
  font-weight: 700;
  margin-bottom: 8px;
  color: #303133;
}

.forbidden-desc {
  font-size: 14px;
  color: #909399;
  margin-bottom: 24px;
}

.mb-4 {
  margin-bottom: 16px;
}

.mr-1 {
  margin-right: 4px;
}

.mb-1 {
  margin-bottom: 4px;
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

@media (max-width: 991px) {
  .bottom-navigation {
    display: flex;
  }
  .admin-main {
    padding-bottom: 72px !important;
  }
}

@media (max-width: 575px) {
  .admin-main {
    padding: 12px;
  }
}
</style>
