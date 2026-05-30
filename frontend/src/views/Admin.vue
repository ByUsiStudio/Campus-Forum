<template>
  <v-app v-if="isInitialized && isAdmin">
    <v-app-bar flat class="admin-header" height="72">
      <div class="header-content">
        <div class="header-brand">
          <v-icon size="28" color="primary" class="brand-icon">mdi-shield-crown</v-icon>
          <div class="brand-text">
            <div class="brand-title">管理后台</div>
            <div class="brand-subtitle">校园论坛管理中心</div>
          </div>
        </div>

        <div class="header-actions">
          <v-btn icon variant="text" size="small" @click="goToHome">
            <v-icon>mdi-home</v-icon>
          </v-btn>
          <v-btn icon variant="text" size="small" @click="handleRefresh">
            <v-icon :class="{ 'spin': isRefreshing }">mdi-refresh</v-icon>
          </v-btn>
        </div>
      </div>
    </v-app-bar>

    <v-navigation-drawer
      permanent
      expand-on-hover
      width="250"
      class="admin-sidebar"
    >
      <v-list dense class="sidebar-menu">
        <v-list-item
          :to="{ name: 'AdminIndex' }"
          :active="route.name === 'AdminIndex'"
          class="sidebar-item"
        >
          <v-list-item-icon class="sidebar-icon">
            <v-icon>mdi-view-dashboard</v-icon>
          </v-list-item-icon>
          <v-list-item-title>概览</v-list-item-title>
        </v-list-item>

        <v-list-item
          :to="{ name: 'AdminUsers' }"
          :active="route.name === 'AdminUsers'"
          class="sidebar-item"
        >
          <v-list-item-icon class="sidebar-icon">
            <v-icon>mdi-account-group</v-icon>
          </v-list-item-icon>
          <v-list-item-title>用户管理</v-list-item-title>
        </v-list-item>

        <v-list-item
          :to="{ name: 'AdminArticles' }"
          :active="route.name === 'AdminArticles'"
          class="sidebar-item"
        >
          <v-list-item-icon class="sidebar-icon">
            <v-icon>mdi-file-document-edit</v-icon>
          </v-list-item-icon>
          <v-list-item-title>文章管理</v-list-item-title>
        </v-list-item>

        <v-list-item
          :to="{ name: 'AdminComments' }"
          :active="route.name === 'AdminComments'"
          class="sidebar-item"
        >
          <v-list-item-icon class="sidebar-icon">
            <v-icon>mdi-comment-text-multiple</v-icon>
          </v-list-item-icon>
          <v-list-item-title>评论管理</v-list-item-title>
        </v-list-item>

        <v-list-item
          :to="{ name: 'AdminCategories' }"
          :active="route.name === 'AdminCategories'"
          class="sidebar-item"
        >
          <v-list-item-icon class="sidebar-icon">
            <v-icon>mdi-shape</v-icon>
          </v-list-item-icon>
          <v-list-item-title>分区管理</v-list-item-title>
        </v-list-item>

        <v-list-item
          :to="{ name: 'AdminTitles' }"
          :active="route.name === 'AdminTitles'"
          class="sidebar-item"
        >
          <v-list-item-icon class="sidebar-icon">
            <v-icon>mdi-medal</v-icon>
          </v-list-item-icon>
          <v-list-item-title>头衔管理</v-list-item-title>
        </v-list-item>

        <v-list-item
          :to="{ name: 'AdminSidebar' }"
          :active="route.name === 'AdminSidebar'"
          class="sidebar-item"
        >
          <v-list-item-icon class="sidebar-icon">
            <v-icon>mdi-web</v-icon>
          </v-list-item-icon>
          <v-list-item-title>侧边栏配置</v-list-item-title>
        </v-list-item>

        <v-list-item
          :to="{ name: 'AdminDeletions' }"
          :active="route.name === 'AdminDeletions'"
          class="sidebar-item"
        >
          <v-list-item-icon class="sidebar-icon">
            <v-icon>mdi-delete-forever</v-icon>
          </v-list-item-icon>
          <v-list-item-title>删除申请</v-list-item-title>
          <v-chip v-if="deletionCount > 0" size="x-small" color="error" class="sidebar-badge">
            {{ deletionCount }}
          </v-chip>
        </v-list-item>

        <v-divider class="sidebar-divider"></v-divider>

        <v-list-item
          :to="{ name: 'AdminAnnouncement' }"
          :active="route.name === 'AdminAnnouncement'"
          class="sidebar-item"
        >
          <v-list-item-icon class="sidebar-icon">
            <v-icon>mdi-bullhorn</v-icon>
          </v-list-item-icon>
          <v-list-item-title>公告管理</v-list-item-title>
        </v-list-item>

        <v-list-item
          :to="{ name: 'AdminSiteConfig' }"
          :active="route.name === 'AdminSiteConfig'"
          class="sidebar-item"
        >
          <v-list-item-icon class="sidebar-icon">
            <v-icon>mdi-globe</v-icon>
          </v-list-item-icon>
          <v-list-item-title>网站配置</v-list-item-title>
        </v-list-item>

        <v-list-item
          :to="{ name: 'AdminSMTPConfig' }"
          :active="route.name === 'AdminSMTPConfig'"
          class="sidebar-item"
        >
          <v-list-item-icon class="sidebar-icon">
            <v-icon>mdi-email-settings</v-icon>
          </v-list-item-icon>
          <v-list-item-title>邮件配置</v-list-item-title>
        </v-list-item>

        <v-list-item
          :to="{ name: 'AdminNotifications' }"
          :active="route.name === 'AdminNotifications'"
          class="sidebar-item"
        >
          <v-list-item-icon class="sidebar-icon">
            <v-icon>mdi-bell</v-icon>
          </v-list-item-icon>
          <v-list-item-title>通知管理</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

    <v-main class="admin-main">
      <router-view></router-view>
    </v-main>
  </v-app>

  <div v-else-if="isInitialized && !isAdmin" class="error-page">
    <div class="error-content">
      <v-icon size="128" color="primary" class="error-icon">mdi-shield-alert</v-icon>
      <h1 class="error-title">权限不足</h1>
      <p class="error-message">您没有访问管理后台的权限</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import api from '../api'

const router = useRouter()
const route = useRoute()

const isInitialized = ref(false)
const isAdmin = ref(false)
const isRefreshing = ref(false)
const deletionCount = ref(0)

const goToHome = () => {
  router.push('/')
}

const handleRefresh = () => {
  isRefreshing.value = true
  setTimeout(() => {
    window.location.reload()
  }, 1000)
}

const checkAdmin = async () => {
  try {
    const response = await api.get('/admin/check')
    isAdmin.value = response.data.is_admin
    if (!isAdmin.value) {
      console.log('不是管理员')
    }
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
    console.error('加载删除申请数量失败', error)
    deletionCount.value = 0
  }
}

onMounted(() => {
  checkAdmin()
  loadDeletionCount()
})

watch(() => route.path, () => {
  loadDeletionCount()
})
</script>

<style scoped>
.admin-header {
  background: linear-gradient(135deg, #6750A4 0%, #7C3AED 50%, #8B5CF6 100%);
  box-shadow: 0 4px 20px -4px rgba(103, 80, 164, 0.3);
  position: relative;
  overflow: hidden;
}

.admin-header::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -25%;
  width: 200px;
  height: 200px;
  background: rgba(255, 255, 255, 0.08);
  border-radius: 50%;
}

.admin-header::after {
  content: '';
  position: absolute;
  bottom: -30%;
  left: -10%;
  width: 150px;
  height: 150px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 50%;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  padding: 0 32px;
  position: relative;
  z-index: 1;
}

.header-brand {
  display: flex;
  align-items: center;
  gap: 14px;
}

.brand-icon {
  background: rgba(255, 255, 255, 0.18);
  backdrop-filter: blur(10px);
  padding: 10px;
  border-radius: 12px;
  transition: transform 0.3s ease;
}

.brand-icon:hover {
  transform: scale(1.05);
}

.brand-text {
  display: flex;
  flex-direction: column;
}

.brand-title {
  color: #fff;
  font-weight: 700;
  font-size: 20px;
  letter-spacing: -0.3px;
}

.brand-subtitle {
  color: rgba(255, 255, 255, 0.85);
  font-size: 12px;
  font-weight: 400;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.header-actions .v-btn {
  color: rgba(255, 255, 255, 0.95);
  background: rgba(255, 255, 255, 0.12);
  border-radius: 10px;
  padding: 8px 16px;
  transition: all 0.25s ease;
}

.header-actions .v-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: translateY(-2px);
}

.admin-sidebar {
  background: linear-gradient(180deg, #FFFFFF 0%, #F8F7FF 100%);
  border-right: 1px solid rgba(103, 80, 164, 0.1);
  box-shadow: 4px 0 20px -4px rgba(0, 0, 0, 0.05);
}

.sidebar-menu {
  padding: 20px 0;
}

.sidebar-item {
  margin: 6px 10px;
  border-radius: 10px;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
}

.sidebar-item::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 3px;
  background: linear-gradient(180deg, #6750A4 0%, #8B5CF6 100%);
  transform: scaleY(0);
  transition: transform 0.25s ease;
}

.sidebar-item:hover {
  background: rgba(103, 80, 164, 0.06);
  transform: translateX(4px);
}

.sidebar-item.active {
  background: rgba(103, 80, 164, 0.1);
  box-shadow: 0 4px 12px -4px rgba(103, 80, 164, 0.15);
}

.sidebar-item.active::before {
  transform: scaleY(1);
}

.sidebar-item.active .sidebar-icon,
.sidebar-item.active .v-list-item-title {
  color: #6750A4;
  font-weight: 500;
}

.sidebar-icon {
  color: #625B71;
  transition: color 0.25s ease;
}

.sidebar-item:hover .sidebar-icon {
  color: #6750A4;
}

.sidebar-badge {
  margin-left: auto;
  animation: badgePulse 2s ease-in-out infinite;
}

@keyframes badgePulse {
  0%, 100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.05);
  }
}

.sidebar-divider {
  margin: 16px 12px;
  background: linear-gradient(90deg, transparent 0%, rgba(103, 80, 164, 0.15) 50%, transparent 100%);
  height: 1px;
}

.admin-main {
  padding: 32px;
  min-height: calc(100vh - 72px);
  background: linear-gradient(135deg, #F8FAFC 0%, #F1F5F9 100%);
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.error-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #F8FAFC 0%, #EEF2FF 50%, #F1F5F9 100%);
}

.error-content {
  text-align: center;
  padding: 64px;
  background: #fff;
  border-radius: 24px;
  box-shadow: 0 20px 60px -20px rgba(103, 80, 164, 0.15);
  animation: slideUp 0.4s ease;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.error-icon {
  opacity: 0.9;
  background: linear-gradient(135deg, #FECACA 0%, #FCA5A5 100%);
  padding: 24px;
  border-radius: 50%;
  box-shadow: 0 8px 32px -8px rgba(239, 68, 68, 0.2);
}

.error-title {
  margin: 28px 0 14px;
  font-size: 32px;
  font-weight: 700;
  color: #1C1B1F;
  letter-spacing: -0.5px;
}

.error-message {
  color: #625B71;
  font-size: 15px;
  line-height: 1.6;
  max-width: 300px;
  margin: 0 auto;
}
</style>