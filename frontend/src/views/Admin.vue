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

        <v-tabs class="header-tabs" slider-color="primary">
          <v-tab :to="{ name: 'AdminIndex' }" class="tab-item">
            <v-icon size="20" class="tab-icon">mdi-view-dashboard</v-icon>
            <span>概览</span>
          </v-tab>
          <v-tab :to="{ name: 'AdminUsers' }" class="tab-item">
            <v-icon size="20" class="tab-icon">mdi-account-group</v-icon>
            <span>用户</span>
          </v-tab>
          <v-tab :to="{ name: 'AdminArticles' }" class="tab-item">
            <v-icon size="20" class="tab-icon">mdi-file-document-edit</v-icon>
            <span>文章</span>
          </v-tab>
          <v-tab :to="{ name: 'AdminComments' }" class="tab-item">
            <v-icon size="20" class="tab-icon">mdi-comment-text-multiple</v-icon>
            <span>评论</span>
          </v-tab>
          <v-tab :to="{ name: 'AdminCategories' }" class="tab-item">
            <v-icon size="20" class="tab-icon">mdi-shape</v-icon>
            <span>分区</span>
          </v-tab>
          <v-tab :to="{ name: 'AdminTitles' }" class="tab-item">
            <v-icon size="20" class="tab-icon">mdi-medal</v-icon>
            <span>头衔</span>
          </v-tab>
          <v-tab :to="{ name: 'AdminSidebar' }" class="tab-item">
            <v-icon size="20" class="tab-icon">mdi-web</v-icon>
            <span>侧边栏</span>
          </v-tab>
          <v-tab :to="{ name: 'AdminDeletions' }" class="tab-item">
            <v-icon size="20" class="tab-icon">mdi-delete-forever</v-icon>
            <span>删除申请</span>
            <v-chip v-if="deletionCount > 0" size="x-small" color="error" class="tab-badge">
              {{ deletionCount }}
            </v-chip>
          </v-tab>
          <v-tab :to="{ name: 'AdminAnnouncement' }" class="tab-item">
            <v-icon size="20" class="tab-icon">mdi-bullhorn</v-icon>
            <span>公告</span>
          </v-tab>
          <v-tab :to="{ name: 'AdminSiteConfig' }" class="tab-item">
            <v-icon size="20" class="tab-icon">mdi-globe</v-icon>
            <span>网站配置</span>
          </v-tab>
          <v-tab :to="{ name: 'AdminSMTPConfig' }" class="tab-item">
            <v-icon size="20" class="tab-icon">mdi-email-settings</v-icon>
            <span>邮件配置</span>
          </v-tab>
          <v-tab :to="{ name: 'AdminNotifications' }" class="tab-item">
            <v-icon size="20" class="tab-icon">mdi-bell</v-icon>
            <span>通知管理</span>
          </v-tab>
        </v-tabs>

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
    const response = await api.get('/admin/deletions')
    deletionCount.value = response.data.length
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
  background: linear-gradient(135deg, #f8f9ff 0%, #fff 100%) !important;
  border-bottom: 1px solid rgba(103, 80, 164, 0.1);
}

.header-content {
  display: flex;
  align-items: center;
  width: 100%;
  padding: 0 24px;
  gap: 32px;
}

.header-brand {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-shrink: 0;
}

.brand-icon {
  background: linear-gradient(135deg, #6750A4 0%, #9c8cd8 100%);
  border-radius: 12px;
  padding: 8px;
  color: white !important;
}

.brand-title {
  font-size: 1.1rem;
  font-weight: 700;
  color: #1a1a2e;
  line-height: 1.3;
}

.brand-subtitle {
  font-size: 0.75rem;
  color: #6b7280;
}

.header-tabs {
  flex: 1;
  max-width: 1200px;
}

.tab-item {
  min-width: auto !important;
  padding: 0 16px !important;
  font-size: 0.875rem;
  font-weight: 500;
  gap: 8px;
}

.tab-icon {
  margin-right: 4px;
}

.tab-badge {
  margin-left: 4px;
}

.header-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.admin-main {
  background: #f8f9ff;
  min-height: calc(100vh - 72px);
}

.error-page {
  background: linear-gradient(135deg, #f8f9ff 0%, #fff 100%);
}

.error-content {
  text-align: center;
  padding: 48px 24px;
}

.error-icon {
  margin-bottom: 24px;
}

.error-title {
  font-size: 1.75rem;
  font-weight: 700;
  color: #1a1a2e;
  margin-bottom: 12px;
}

.error-message {
  font-size: 1rem;
  color: #6b7280;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.spin {
  animation: spin 1s linear infinite;
}

@media (max-width: 960px) {
  .header-content {
    padding: 0 16px;
    gap: 16px;
  }

  .brand-subtitle {
    display: none;
  }

  .header-tabs {
    max-width: none;
  }

  .tab-item span {
    display: none;
  }
}

@media (max-width: 640px) {
  .header-tabs {
    overflow-x: auto;
  }

  .tab-item {
    padding: 0 12px !important;
  }
}
</style>
