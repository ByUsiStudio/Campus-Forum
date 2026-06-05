<template>
  <v-app v-if="isInitialized && isAdmin">
    <!-- 左侧导航栏 -->
    <v-navigation-drawer
      v-model="drawerOpen"
      :rail="sidebarCollapsed"
      :permanent="!isMobile"
      :temporary="isMobile"
      :width="drawerWidth"
      :temporary-width="drawerWidth"
      color="surface"
      border
      elevation="1"
      touchless
      @touchstart="handleTouchStart"
      @touchmove="handleTouchMove"
      @touchend="handleTouchEnd"
    >
      <!-- 头部品牌区域 -->
      <v-sheet color="primary" class="pa-4">
        <div class="d-flex align-center justify-space-between">
          <div class="d-flex align-center">
            <v-icon size="28" color="white">mdi-shield-crown</v-icon>
            <span v-if="!sidebarCollapsed" class="text-h6 font-weight-bold ml-3 text-white">管理后台</span>
          </div>
          <v-btn
            v-if="!isMobile"
            icon
            variant="text"
            size="x-small"
            @click="toggleSidebar"
            color="white"
          >
            <v-icon size="18">
              {{ sidebarCollapsed ? 'mdi-chevron-right' : 'mdi-chevron-left' }}
            </v-icon>
          </v-btn>
        </div>
      </v-sheet>

      <v-divider />

      <!-- 导航菜单 -->
      <v-list nav density="compact" class="pa-2">
        <!-- 动态链接项 -->
        <v-list-subheader v-if="sidebarItems.length > 0 && !sidebarCollapsed">快捷链接</v-list-subheader>
        <v-list-item
          v-for="item in sidebarItems"
          :key="item.link"
          :to="item.link"
          :active="isActive(item.link)"
          color="primary"
          rounded="lg"
          @click="drawerOpen = false"
        >
          <template v-slot:prepend>
            <v-icon size="20">{{ getIcon(item.icon) }}</v-icon>
          </template>
          <v-list-item-title v-if="!sidebarCollapsed">{{ item.title }}</v-list-item-title>
        </v-list-item>

        <!-- 管理功能 -->
        <v-divider class="my-2" />
        <v-list-subheader v-if="!sidebarCollapsed">管理功能</v-list-subheader>
        
        <v-list-item
          v-for="item in adminItems"
          :key="item.route"
          :to="{ name: item.route }"
          :active="route.name === item.route"
          color="primary"
          rounded="lg"
          @click="drawerOpen = false"
        >
          <template v-slot:prepend>
            <v-icon size="20">{{ item.icon }}</v-icon>
          </template>
          <v-list-item-title v-if="!sidebarCollapsed">{{ item.title }}</v-list-item-title>
          <template v-slot:append v-if="item.badge && item.badge() > 0 && !sidebarCollapsed">
            <v-badge :content="item.badge()" color="error" inline />
          </template>
        </v-list-item>
      </v-list>

      <!-- 底部信息 -->
      <template v-slot:append>
        <v-divider />
        <v-sheet color="grey-lighten-4" class="pa-2 text-center" v-if="!sidebarCollapsed">
          <v-chip size="x-small" variant="tonal" color="primary">
            v{{ version || '1.0' }}
          </v-chip>
        </v-sheet>
      </template>
    </v-navigation-drawer>

    <!-- 顶部工具栏 -->
    <v-app-bar flat color="surface" border elevation="1">
      <v-app-bar-title>
        <v-breadcrumbs :items="breadcrumbs">
          <template v-slot:divider>
            <v-icon icon="mdi-chevron-right" size="small"></v-icon>
          </template>
        </v-breadcrumbs>
      </v-app-bar-title>
      
      <v-spacer />
      
      <v-tooltip text="返回首页" location="bottom">
        <template v-slot:activator="{ props }">
          <v-btn
            v-bind="props"
            icon
            variant="text"
            size="small"
            @click="goToHome"
          >
            <v-icon size="20">mdi-home-outline</v-icon>
          </v-btn>
        </template>
      </v-tooltip>
      
      <v-tooltip text="刷新页面" location="bottom">
        <template v-slot:activator="{ props }">
          <v-btn
            v-bind="props"
            icon
            variant="text"
            size="small"
            @click="handleRefresh"
            :class="{ 'animate-spin': isRefreshing }"
          >
            <v-icon size="20">mdi-refresh</v-icon>
          </v-btn>
        </template>
      </v-tooltip>
    </v-app-bar>

    <!-- 主要内容区域 -->
    <v-main class="bg-grey-lighten-4">
      <v-container fluid class="pa-6">
        <router-view />
      </v-container>
    </v-main>

    <!-- 悬浮按钮 - 打开/关闭侧边栏 -->
    <v-tooltip :text="drawerOpen ? '关闭侧边栏' : '打开侧边栏'" location="left">
      <template v-slot:activator="{ props }">
        <v-fab
          v-show="!drawerOpen || isMobile"
          v-bind="props"
          @click="toggleDrawer"
          :color="drawerOpen ? 'grey' : 'primary'"
          :icon="drawerOpen ? 'mdi-close' : 'mdi-menu'"
          location="bottom right"
          app
          appear
          size="large"
          elevation="4"
          class="ma-4"
        />
      </template>
    </v-tooltip>

    <!-- 移动端遮罩层 -->
    <v-overlay
      v-model="drawerOpen"
      v-if="isMobile"
      scrim
      persistent
      z-index="999"
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

// 触摸滑动相关
const touchStartX = ref(0)
const touchStartY = ref(0)
const touchDeltaX = ref(0)
const isSwiping = ref(false)

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

// 根据屏幕大小动态计算抽屉宽度
const drawerWidth = computed(() => {
  if (isMobile.value) return 280
  if (window.innerWidth < 1280) return 240
  if (window.innerWidth < 1920) return 260
  return 300
})

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
  const width = window.innerWidth
  isMobile.value = width < 960
  
  // 根据屏幕大小调整侧边栏状态
  if (width < 600) {
    // 小屏幕默认收起
    sidebarCollapsed.value = false
    drawerOpen.value = false
  } else if (width < 960) {
    // 中等屏幕
    sidebarCollapsed.value = false
  } else if (width < 1280) {
    // 较大屏幕
    sidebarCollapsed.value = false
  } else {
    // 大屏幕，恢复保存的状态
    const savedState = localStorage.getItem('adminSidebarCollapsed')
    if (savedState !== null) {
      sidebarCollapsed.value = JSON.parse(savedState)
    }
  }
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

// 触摸滑动处理
const handleTouchStart = (e) => {
  touchStartX.value = e.touches[0].clientX
  touchStartY.value = e.touches[0].clientY
  touchDeltaX.value = 0
  isSwiping.value = false
}

const handleTouchMove = (e) => {
  const deltaX = e.touches[0].clientX - touchStartX.value
  const deltaY = e.touches[0].clientY - touchStartY.value
  
  // 判断是否为水平滑动（横向移动大于纵向移动）
  if (Math.abs(deltaX) > Math.abs(deltaY) && Math.abs(deltaX) > 10) {
    isSwiping.value = true
    touchDeltaX.value = deltaX
    
    // 阻止默认的滚动行为
    if (Math.abs(deltaX) > 20) {
      e.preventDefault()
    }
  }
}

const handleTouchEnd = () => {
  if (!isSwiping.value) return
  
  const threshold = 80 // 滑动阈值
  
  if (touchDeltaX.value > threshold && !drawerOpen.value) {
    // 从左向右滑动，打开抽屉
    drawerOpen.value = true
  } else if (touchDeltaX.value < -threshold && drawerOpen.value) {
    // 从右向左滑动，关闭抽屉
    drawerOpen.value = false
  }
  
  isSwiping.value = false
  touchDeltaX.value = 0
}

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  
  // 桌面端默认打开抽屉
  if (!isMobile.value) {
    drawerOpen.value = true
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
/* 使用 Vuetify 内置样式，仅保留必要动画 */
.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* 侧边栏触摸优化 */
:deep(.v-navigation-drawer) {
  transition: width 0.3s cubic-bezier(0.4, 0, 0.2, 1) !important;
}

:deep(.v-navigation-drawer__content) {
  touch-action: pan-y;
}

/* 滑动时的视觉反馈 */
:deep(.v-navigation-drawer.temporary) {
  will-change: transform;
}

@media (max-width: 960px) {
  .v-main {
    padding-top: 64px !important;
  }
}

/* 小屏幕优化 */
@media (max-width: 600px) {
  :deep(.v-navigation-drawer) {
    width: 100% !important;
  }
}
</style>