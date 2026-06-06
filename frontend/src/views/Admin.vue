<template>
  <v-app v-if="isInitialized && isAdmin">
    <!-- 左侧导航栏 - 桌面端 -->
    <v-navigation-drawer
      v-if="!isMobile"
      v-model="drawerOpen"
      :rail="sidebarCollapsed"
      permanent
      :width="sidebarCollapsed ? 80 : 260"
      color="surface"
      border
    >
      <!-- 头部 -->
      <v-sheet color="primary" class="pa-4 d-flex align-center">
        <v-icon size="24" color="white">mdi-shield-crown</v-icon>
        <span v-if="!sidebarCollapsed" class="text-subtitle-1 font-weight-bold ml-2 text-white">管理后台</span>
        <v-btn
          icon
          variant="text"
          size="x-small"
          @click="sidebarCollapsed = !sidebarCollapsed"
          color="white"
          class="ml-auto"
        >
          <v-icon size="16">{{ sidebarCollapsed ? 'mdi-chevron-right' : 'mdi-chevron-left' }}</v-icon>
        </v-btn>
      </v-sheet>

      <v-divider />

      <!-- 导航菜单 -->
      <v-list nav density="compact" class="pa-2">
        <v-list-item
          v-for="item in adminItems"
          :key="item.route"
          :to="{ name: item.route }"
          color="primary"
          rounded="lg"
          class="mb-1"
        >
          <template v-slot:prepend>
            <v-icon :size="sidebarCollapsed ? 22 : 20">{{ item.icon }}</v-icon>
          </template>
          <v-list-item-title v-if="!sidebarCollapsed">{{ item.title }}</v-list-item-title>
          <v-badge v-if="item.badge && item.badge() > 0 && sidebarCollapsed" :content="item.badge()" color="error" floating />
          <template v-slot:append v-if="item.badge && item.badge() > 0 && !sidebarCollapsed">
            <v-badge :content="item.badge()" color="error" inline />
          </template>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

    <!-- 顶部工具栏 -->
    <v-app-bar flat color="surface" border height="56">
      <!-- 移动端菜单按钮 -->
      <v-app-bar-nav-icon v-if="isMobile" @click="drawerOpen = !drawerOpen" />
      
      <v-app-bar-title class="text-subtitle-1 font-weight-medium">
        {{ currentPageTitle }}
      </v-app-bar-title>
      
      <v-spacer />
      
      <v-btn icon variant="text" size="small" to="/">
        <v-icon size="20">mdi-home</v-icon>
        <v-tooltip activator="parent" location="bottom">返回首页</v-tooltip>
      </v-btn>
      
      <v-btn icon variant="text" size="small" @click="handleRefresh" :loading="isRefreshing">
        <v-icon size="20">mdi-refresh</v-icon>
        <v-tooltip activator="parent" location="bottom">刷新</v-tooltip>
      </v-btn>
    </v-app-bar>

    <!-- 移动端侧边栏 -->
    <v-navigation-drawer
      v-if="isMobile"
      v-model="drawerOpen"
      temporary
      width="280"
      color="surface"
    >
      <v-sheet color="primary" class="pa-4">
        <div class="d-flex align-center">
          <v-icon size="24" color="white">mdi-shield-crown</v-icon>
          <span class="text-subtitle-1 font-weight-bold ml-2 text-white">管理后台</span>
        </div>
      </v-sheet>

      <v-divider />

      <v-list nav density="compact" class="pa-2">
        <v-list-item
          v-for="item in adminItems"
          :key="item.route"
          :to="{ name: item.route }"
          color="primary"
          rounded="lg"
          @click="drawerOpen = false"
          class="mb-1"
        >
          <template v-slot:prepend>
            <v-icon size="20">{{ item.icon }}</v-icon>
          </template>
          <v-list-item-title>{{ item.title }}</v-list-item-title>
          <template v-slot:append v-if="item.badge && item.badge() > 0">
            <v-badge :content="item.badge()" color="error" inline />
          </template>
        </v-list-item>
      </v-list>

      <v-divider class="mt-4" />
      
      <v-list nav density="compact" class="pa-2">
        <v-list-item to="/" @click="drawerOpen = false" rounded="lg">
          <template v-slot:prepend>
            <v-icon size="20">mdi-home</v-icon>
          </template>
          <v-list-item-title>返回首页</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

    <!-- 主要内容区域 -->
    <v-main class="bg-grey-lighten-4">
      <v-container fluid class="pa-4 pa-sm-6">
        <router-view />
      </v-container>
    </v-main>
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
const isRefreshing = ref(false)
const isMobile = ref(false)
const deletionCount = ref(0)
const sidebarCollapsed = ref(false)
const drawerOpen = ref(false)

const adminItems = [
  { route: 'AdminIndex', title: '数据概览', icon: 'mdi-view-dashboard' },
  { route: 'AdminUsers', title: '用户管理', icon: 'mdi-account-group' },
  { route: 'AdminArticles', title: '文章管理', icon: 'mdi-file-document' },
  { route: 'AdminComments', title: '评论管理', icon: 'mdi-comment-text' },
  { route: 'AdminCategories', title: '分区管理', icon: 'mdi-shape' },
  { route: 'AdminTitles', title: '头衔管理', icon: 'mdi-medal' },
  { route: 'AdminSidebar', title: '侧边栏', icon: 'mdi-web' },
  { route: 'AdminDeletions', title: '删除申请', icon: 'mdi-delete', badge: () => deletionCount.value },
  { route: 'AdminAnnouncement', title: '公告管理', icon: 'mdi-bullhorn' },
  { route: 'AdminSiteConfig', title: '网站配置', icon: 'mdi-cog' },
  { route: 'AdminSMTPConfig', title: '邮件配置', icon: 'mdi-email' },
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

const currentPageTitle = computed(() => pageTitles[route.name] || '管理后台')

const checkMobile = () => {
  isMobile.value = window.innerWidth < 768
  if (isMobile.value) {
    drawerOpen.value = false
    sidebarCollapsed.value = false
  } else {
    const saved = localStorage.getItem('adminSidebarCollapsed')
    if (saved !== null) sidebarCollapsed.value = JSON.parse(saved)
    drawerOpen.value = true
  }
}

const handleRefresh = () => {
  isRefreshing.value = true
  window.location.reload()
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

watch(sidebarCollapsed, (val) => {
  localStorage.setItem('adminSidebarCollapsed', val)
})

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  checkAdmin()
  loadDeletionCount()
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})

watch(() => route.path, loadDeletionCount)
</script>

<style scoped>
:deep(.v-navigation-drawer) {
  transition: width 0.2s ease !important;
}
</style>