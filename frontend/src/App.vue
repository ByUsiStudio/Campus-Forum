<template>
  <v-app>
    <!-- 主要内容区域 -->
    <v-main :class="{ 'pb-navigation': !hideAppBar && isMobile }" class="bg-grey-lighten-4">
      <!-- 移动端顶部搜索栏 -->
      <div v-if="!hideAppBar && isMobile && !isAdminPage" class="mobile-search-bar px-4 pt-3 pb-2">
        <v-text-field
          v-model="searchQuery"
          prepend-inner-icon="mdi-magnify"
          variant="solo"
          density="compact"
          placeholder="搜索..."
          hide-details
          rounded="pill"
          bg-color="white"
          @keyup.enter="handleSearch"
          class="mobile-search-field"
        ></v-text-field>
      </div>
      
      <v-container v-if="!hideAppBar" fluid class="pa-4">
        <router-view />
      </v-container>
      <v-container v-else>
        <router-view />
      </v-container>
    </v-main>

    <!-- 移动端底部导航栏 -->
    <v-bottom-navigation v-if="!hideAppBar && isMobile && !isAdminPage" grow color="primary" app>
      <v-btn to="/" value="home">
        <v-icon>mdi-home</v-icon>
        <span>首页</span>
      </v-btn>
      
      <v-btn to="/create" value="create" v-if="token">
        <v-icon>mdi-plus-circle</v-icon>
        <span>发布</span>
      </v-btn>
      
      <v-btn @click="appDrawer = true" value="drawer">
        <v-icon>mdi-grid</v-icon>
        <span>菜单</span>
      </v-btn>
      
      <v-btn :to="token ? '/profile' : '/login'" :value="token ? 'profile' : 'login'">
        <v-icon>mdi-account</v-icon>
        <span>{{ token ? '我的' : '登录' }}</span>
      </v-btn>
    </v-bottom-navigation>

    <!-- 移动端管理后台底部导航栏 -->
    <v-bottom-navigation v-if="!hideAppBar && isMobile && isAdminPage" grow color="primary" app>
      <v-btn :to="{ name: 'AdminDashboard' }" value="dashboard">
        <v-icon>mdi-view-dashboard</v-icon>
        <span>控制台</span>
      </v-btn>
      
      <v-btn :to="{ name: 'AdminUsers' }" value="users">
        <v-icon>mdi-account-multiple</v-icon>
        <span>用户管理</span>
      </v-btn>
      
      <v-btn :to="{ name: 'AdminArticles' }" value="articles">
        <v-icon>mdi-file-document</v-icon>
        <span>文章管理</span>
      </v-btn>
      
      <v-btn :to="{ name: 'AdminSettings' }" value="settings">
        <v-icon>mdi-settings</v-icon>
        <span>系统设置</span>
      </v-btn>
    </v-bottom-navigation>

    <!-- 移动端应用抽屉（从底部弹出） -->
    <Teleport to="body">
      <div v-if="appDrawer" class="app-drawer-overlay" @click="appDrawer = false">
        <div class="app-drawer-content" @click.stop>
          <div class="app-drawer-header">
            <div class="app-drawer-title">应用菜单</div>
            <v-btn icon @click="appDrawer = false" class="close-btn">
              <v-icon>mdi-close</v-icon>
            </v-btn>
          </div>
          
          <v-divider class="mb-4"></v-divider>
          
          <div class="app-drawer-grid">
            <v-btn 
              v-for="item in appDrawerItems" 
              :key="item.path"
              :to="item.path"
              class="drawer-item"
              @click="appDrawer = false"
            >
              <v-icon size="32" class="drawer-icon">{{ item.icon }}</v-icon>
              <span class="drawer-label">{{ item.label }}</span>
            </v-btn>
          </div>
          
          <v-divider class="mt-4 mb-3"></v-divider>
          
          <div class="app-drawer-footer">
            <v-btn @click="logout" class="footer-btn">
              <v-icon>mdi-logout</v-icon>
              <span>退出登录</span>
            </v-btn>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- 桌面端顶部导航栏 -->
    <v-app-bar v-if="!hideAppBar && !isMobile" elevation="2" color="primary" scroll-behavior="collapse">
      <v-app-bar-title class="d-flex align-center">
        <v-icon class="mr-2">{{ isAdminPage ? 'mdi-shield-crown' : 'mdi-forum' }}</v-icon>
        <router-link :to="isAdminPage ? '/admin' : '/'" class="logo-link">{{ isAdminPage ? '管理后台' : siteTitle }}</router-link>
      </v-app-bar-title>

      <v-spacer></v-spacer>

      <template v-slot:append>
        <div class="nav-links" v-if="!isAdminPage">
          <v-text-field
            v-model="searchQuery"
            prepend-inner-icon="mdi-magnify"
            variant="outlined"
            density="compact"
            color="white"
            placeholder="搜索..."
            @keyup.enter="handleSearch"
            @click:prepend-inner="handleSearch"
            class="search-field"
          ></v-text-field>
          <v-btn variant="text" to="/" color="white" prepend-icon="mdi-home">首页</v-btn>
          <v-btn variant="text" to="/login" v-if="!token" color="white" prepend-icon="mdi-login">登录</v-btn>
          <v-btn variant="text" to="/register" v-if="!token" color="white" prepend-icon="mdi-account-plus">注册</v-btn>
          <v-btn variant="text" to="/create" v-if="token" color="white" prepend-icon="mdi-pencil">写文章</v-btn>
          <NotificationBell v-if="token" />
          <v-btn variant="text" to="/profile" v-if="token" color="white" prepend-icon="mdi-account">我的</v-btn>
          <v-btn variant="text" to="/admin" v-if="isAdmin" color="error" prepend-icon="mdi-shield-crown">管理后台</v-btn>
          <v-btn variant="text" v-if="token" @click="logout" color="white" prepend-icon="mdi-logout">退出</v-btn>
        </div>
        <div class="nav-links" v-else>
          <v-btn variant="text" :to="{ name: 'AdminDashboard' }" color="white" prepend-icon="mdi-view-dashboard">控制台</v-btn>
          <v-btn variant="text" :to="{ name: 'AdminUsers' }" color="white" prepend-icon="mdi-account-multiple">用户管理</v-btn>
          <v-btn variant="text" :to="{ name: 'AdminArticles' }" color="white" prepend-icon="mdi-file-document">文章管理</v-btn>
          <v-btn variant="text" :to="{ name: 'AdminSettings' }" color="white" prepend-icon="mdi-settings">系统设置</v-btn>
          <v-btn variant="text" to="/" color="white" prepend-icon="mdi-home">返回首页</v-btn>
        </div>
      </template>
    </v-app-bar>

    <!-- 移动端抽屉导航（桌面端不使用） -->
    <v-navigation-drawer v-if="!hideAppBar && isMobile" v-model="drawer" temporary color="surface">
      <v-list nav density="comfortable">
        <template v-if="!isAdminPage">
          <v-list-item to="/" @click="drawer = false" prepend-icon="mdi-home" title="首页"></v-list-item>
          <template v-if="!token">
            <v-list-item to="/login" @click="drawer = false" prepend-icon="mdi-login" title="登录"></v-list-item>
            <v-list-item to="/register" @click="drawer = false" prepend-icon="mdi-account-plus" title="注册"></v-list-item>
          </template>
          <template v-if="token">
            <v-list-item to="/create" @click="drawer = false" prepend-icon="mdi-pencil" title="写文章"></v-list-item>
            <v-list-item to="/profile" @click="drawer = false" prepend-icon="mdi-account" title="我的"></v-list-item>
            <v-list-item to="/notifications" @click="drawer = false" prepend-icon="mdi-bell" title="通知"></v-list-item>
            <v-list-item v-if="isAdmin" to="/admin" @click="drawer = false" prepend-icon="mdi-shield-crown" title="管理后台" class="text-error"></v-list-item>
            <v-divider class="my-2"></v-divider>
            <v-list-item @click="logout" prepend-icon="mdi-logout" title="退出" class="text-secondary"></v-list-item>
          </template>
        </template>
        <template v-else>
          <v-list-item :to="{ name: 'AdminDashboard' }" @click="drawer = false" prepend-icon="mdi-view-dashboard" title="控制台"></v-list-item>
          <v-list-item :to="{ name: 'AdminUsers' }" @click="drawer = false" prepend-icon="mdi-account-multiple" title="用户管理"></v-list-item>
          <v-list-item :to="{ name: 'AdminArticles' }" @click="drawer = false" prepend-icon="mdi-file-document" title="文章管理"></v-list-item>
          <v-list-item :to="{ name: 'AdminSettings' }" @click="drawer = false" prepend-icon="mdi-settings" title="系统设置"></v-list-item>
          <v-divider class="my-2"></v-divider>
          <v-list-item to="/" @click="drawer = false" prepend-icon="mdi-home" title="返回首页"></v-list-item>
          <v-list-item @click="logout" prepend-icon="mdi-logout" title="退出" class="text-secondary"></v-list-item>
        </template>
      </v-list>
    </v-navigation-drawer>

    <!-- 页脚 -->
    <v-footer v-if="!hideAppBar" class="open-source-footer bg-transparent">
      <v-container fluid>
        <v-row justify="center" align="center">
          <v-col cols="12" sm="auto">
            <div class="d-flex align-center flex-wrap justify-center gap-2">
              <span class="text-caption text-medium-emphasis">本论坛基于</span>
              <a href="https://github.com/ByUsiStudio/Campus-Forum" target="_blank" class="open-source-link">
                <v-icon size="small" class="mr-1">mdi-github</v-icon>GitHub
              </a>
              <span class="text-caption text-medium-emphasis">与</span>
              <a href="https://gitee.com/byusistudio/campus-forum" target="_blank" class="open-source-link">
                <v-icon size="small" class="mr-1">mdi-git</v-icon>Gitee
              </a>
              <span class="text-caption text-medium-emphasis">开源</span>
            </div>
            <div class="text-caption text-medium-emphasis text-center mt-1">
              <a href="https://github.com/ByUsiStudio/Campus-Forum" target="_blank" class="open-source-link">
                https://github.com/ByUsiStudio/Campus-Forum
              </a>
            </div>
            <div class="text-caption text-medium-emphasis text-center mt-1">
              <v-icon size="x-small" class="mr-1">mdi-tag</v-icon>
              前端版本: {{ frontendVersion }} | 后端版本: {{ backendVersion || 'unknown' }}
            </div>
            <div v-if="icpNumber || publicSecurityNumber" class="d-flex align-center justify-center flex-wrap gap-4 mt-2">
              <div v-if="icpNumber" class="icp-info">
                <v-icon size="x-small" class="mr-1">mdi-shield-check</v-icon>
                <a href="https://beian.miit.gov.cn" target="_blank" class="open-source-link">{{ icpNumber }}</a>
              </div>
              <div v-if="publicSecurityNumber" class="security-info">
                <v-icon size="x-small" class="mr-1">mdi-police-badge</v-icon>
                <span>{{ publicSecurityNumber }}</span>
              </div>
            </div>
          </v-col>
        </v-row>
      </v-container>
    </v-footer>

    <!-- 全局弹窗 -->
    <AppModal
      :show="modalState.show"
      :type="modalState.type"
      :title="modalState.title"
      :message="modalState.message"
      :icon="modalState.icon"
      :icon-color="modalState.iconColor"
      :confirm-text="modalState.confirmText"
      :cancel-text="modalState.cancelText"
      :confirm-color="modalState.confirmColor"
      :input-value="modalState.inputValue"
      :input-label="modalState.inputLabel"
      :input-type="modalState.inputType"
      :input-placeholder="modalState.inputPlaceholder"
      :input-rows="modalState.inputRows"
      @update:show="modalState.show = $event"
      @confirm="(value) => handleConfirm(value)"
      @cancel="handleCancel"
    />
  </v-app>
</template>

<script>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import api from './api'
import AppModal from './components/AppModal.vue'
import NotificationBell from './components/NotificationBell.vue'
import { modalState, handleConfirm, handleCancel } from './utils/modal'

export default {
  name: 'App',
  components: {
    AppModal,
    NotificationBell
  },
  setup() {
    const router = useRouter()
    const route = useRoute()
    const token = ref(localStorage.getItem('token'))
    const user = ref(null)
    const drawer = ref(false)
    const appDrawer = ref(false)
    const isMobile = ref(false)
    const hideAppBar = ref(false)
    const backendVersion = ref(null)
    const siteTitle = ref('校园论坛')
    const frontendVersion = ref(typeof __FRONTEND_VERSION__ !== 'undefined' ? __FRONTEND_VERSION__ : 'unknown')
    const icpNumber = ref(null)
    const publicSecurityNumber = ref(null)
    const searchQuery = ref('')
    
    const handleSearch = () => {
      if (searchQuery.value.trim()) {
        router.push({ path: '/search', query: { q: searchQuery.value.trim() } })
      }
    }
    
    // 检测屏幕宽度
    const checkMobile = () => {
      isMobile.value = window.innerWidth < 768
    }
    
    // 监听路由变化更新hideAppBar
    watch(() => route.meta, (meta) => {
      hideAppBar.value = meta?.hideAppBar || false
    }, { immediate: true })
    
    const isAdmin = computed(() => {
      return user.value && (user.value.role === 'admin' || user.value.role === 'system' || user.value.role === 'Admin' || user.value.role === 'System')
    })
    
    const isAdminPage = computed(() => {
      return route.path.startsWith('/admin')
    })
    
    const logout = () => {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.dispatchEvent(new Event('user-logout'))
      drawer.value = false
      appDrawer.value = false
      router.push('/login')
    }
    
    const appDrawerItems = computed(() => {
      const items = [
        { path: '/', icon: 'mdi-home', label: '首页' },
        { path: '/create', icon: 'mdi-pencil', label: '写文章' },
        { path: '/notifications', icon: 'mdi-bell', label: '通知' },
        { path: '/profile', icon: 'mdi-account', label: '我的' }
      ]
      if (isAdmin.value) {
        items.push({ path: '/admin', icon: 'mdi-shield-crown', label: '管理后台' })
      }
      return items.filter(item => {
        if (item.path === '/create' || item.path === '/notifications' || item.path === '/profile') {
          return token.value
        }
        return true
      })
    })
    
    const loadUser = async () => {
      if (token.value) {
        try {
          const response = await api.get('/profile')
          user.value = response.data
          localStorage.setItem('user', JSON.stringify(response.data))
        } catch (error) {
          console.error('加载用户信息失败', error)
          localStorage.removeItem('token')
          token.value = null
        }
      }
    }

    const loadVersion = async () => {
      try {
        const response = await api.get('/version')
        backendVersion.value = response.data.backend?.version || response.data.backend_version || response.data.version
        frontendVersion.value = response.data.frontend?.version || frontendVersion.value
      } catch (error) {
        console.error('加载版本信息失败', error)
      }
    }

    const loadSiteTitle = async () => {
      try {
        const response = await api.get('/site-config')
        siteTitle.value = response.data.site_title || '校园论坛'
        icpNumber.value = response.data.icp_number || null
        publicSecurityNumber.value = response.data.public_security_number || null
        document.title = siteTitle.value
      } catch (error) {
        console.error('加载网站标题失败', error)
      }
    }

    const updateSiteTitle = (newTitle) => {
      if (newTitle) {
        siteTitle.value = newTitle
        document.title = newTitle
      }
    }

    watch(siteTitle, (newTitle) => {
      if (newTitle) {
        document.title = newTitle
      }
    })
    
    onMounted(() => {
      checkMobile()
      window.addEventListener('resize', checkMobile)

      const storedUser = localStorage.getItem('user')
      if (storedUser) {
        user.value = JSON.parse(storedUser)
      }
      loadUser()
      loadVersion()
      loadSiteTitle()

      // 监听登录/登出事件
      window.addEventListener('user-updated', () => {
        const storedUser = localStorage.getItem('user')
        if (storedUser) {
          user.value = JSON.parse(storedUser)
        }
        token.value = localStorage.getItem('token')
      })

      window.addEventListener('site-title-updated', (event) => {
        if (event.detail) {
          siteTitle.value = event.detail
          document.title = event.detail
        }
      })

      window.addEventListener('user-logout', () => {
        user.value = null
        token.value = null
      })
    })
    
    return {
      token,
      user,
      isAdmin,
      isAdminPage,
      logout,
      modalState,
      handleConfirm,
      handleCancel,
      isMobile,
      drawer,
      appDrawer,
      appDrawerItems,
      backendVersion,
      siteTitle,
      updateSiteTitle,
      frontendVersion,
      icpNumber,
      publicSecurityNumber,
      searchQuery,
      handleSearch
    }
  }
}
</script>

<style>
.logo-link {
  color: white !important;
  text-decoration: none;
  font-weight: bold;
  font-size: 1.25rem;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.nav-links {
  display: flex;
  gap: 4px;
  align-items: center;
}

.nav-links .v-btn {
  font-weight: 500;
}

.search-field {
  width: 240px;
}

.search-field :deep(.v-field__outline) {
  color: rgba(255, 255, 255, 0.7);
}

.search-field :deep(.v-label),
.search-field :deep(.v-field__input) {
  color: white !important;
}

.search-field :deep(.v-icon) {
  color: white !important;
}

.chat-badge {
  position: absolute;
  top: -4px;
  right: -8px;
}

.open-source-footer {
  margin-top: auto;
  padding: 16px 0;
}

.open-source-link {
  color: rgb(var(--v-theme-primary));
  text-decoration: none;
  font-weight: 500;
}

.open-source-link:hover {
  text-decoration: underline;
}

/* 移动端底部导航栏 padding */
.pb-navigation {
  padding-bottom: 56px !important;
}

/* 移动端搜索栏样式 */
.mobile-search-bar {
  background: rgb(var(--v-theme-background));
  position: sticky;
  top: 0;
  z-index: 1;
}

.mobile-search-field {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.mobile-search-field :deep(.v-field) {
  border-radius: 24px;
}

/* 移动端应用抽屉样式 */
.app-drawer-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 9999;
  display: flex;
  align-items: flex-end;
  animation: fadeIn 0.2s ease-out;
}

.app-drawer-content {
  width: 100%;
  max-height: 80vh;
  background-color: #fff;
  border-radius: 24px 24px 0 0;
  padding: 16px;
  animation: slideUp 0.3s ease-out;
}

.app-drawer-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 16px;
}

.app-drawer-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: #1a1a2e;
}

.close-btn {
  color: #6b7280;
}

.app-drawer-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  padding: 0 8px;
}

.drawer-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16px 8px;
  border-radius: 12px;
  transition: background-color 0.2s;
}

.drawer-item:hover {
  background-color: #f3f4f6;
}

.drawer-icon {
  color: #4b5563;
  margin-bottom: 8px;
}

.drawer-label {
  font-size: 0.75rem;
  color: #374151;
  font-weight: 500;
}

.app-drawer-footer {
  padding: 8px 16px;
}

.footer-btn {
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 8px;
  color: #ef4444;
  font-weight: 500;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes slideUp {
  from {
    transform: translateY(100%);
  }
  to {
    transform: translateY(0);
  }
}
</style>
