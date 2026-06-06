<template>
  <v-app>
    <!-- 根据路由meta决定是否显示顶部栏 -->
    <v-app-bar v-if="!hideAppBar" elevation="2" color="primary" scroll-behavior="collapse">
      <v-app-bar-nav-icon v-if="isMobile" @click="drawer = !drawer"></v-app-bar-nav-icon>
      
      <v-app-bar-title class="d-flex align-center">
        <v-icon class="mr-2">mdi-forum</v-icon>
        <router-link to="/" class="logo-link">{{ siteTitle }}</router-link>
      </v-app-bar-title>

      <v-spacer></v-spacer>

      <!-- 桌面端导航 -->
      <template v-if="!isMobile" v-slot:append>
        <div class="nav-links">
          <v-btn variant="text" to="/" color="white" prepend-icon="mdi-home">首页</v-btn>
          <v-btn variant="text" to="/login" v-if="!token" color="white" prepend-icon="mdi-login">登录</v-btn>
          <v-btn variant="text" to="/register" v-if="!token" color="white" prepend-icon="mdi-account-plus">注册</v-btn>
          <v-btn variant="text" to="/create" v-if="token" color="white" prepend-icon="mdi-pencil">写文章</v-btn>
          <v-btn variant="text" to="/profile" v-if="token" color="white" prepend-icon="mdi-account">个人中心</v-btn>
          <NotificationBell v-if="token" />
          <v-btn variant="text" to="/admin" v-if="isAdmin" color="error" prepend-icon="mdi-shield-crown">管理后台</v-btn>
          <v-btn variant="text" v-if="token" @click="logout" color="white" prepend-icon="mdi-logout">退出</v-btn>
        </div>
      </template>

      <!-- 移动端导航 -->
      <template v-else v-slot:append>
        <NotificationBell v-if="token" />
      </template>
    </v-app-bar>

    <!-- 移动端抽屉导航 -->
    <v-navigation-drawer v-if="!hideAppBar" v-model="drawer" temporary color="surface">
      <v-list nav density="comfortable">
        <v-list-item to="/" @click="drawer = false" prepend-icon="mdi-home" title="首页"></v-list-item>
        <template v-if="!token">
          <v-list-item to="/login" @click="drawer = false" prepend-icon="mdi-login" title="登录"></v-list-item>
          <v-list-item to="/register" @click="drawer = false" prepend-icon="mdi-account-plus" title="注册"></v-list-item>
        </template>
        <template v-if="token">
          <v-list-item to="/create" @click="drawer = false" prepend-icon="mdi-pencil" title="写文章"></v-list-item>
          <v-list-item to="/profile" @click="drawer = false" prepend-icon="mdi-account" title="个人中心"></v-list-item>
          <v-list-item v-if="isAdmin" to="/admin" @click="drawer = false" prepend-icon="mdi-shield-crown" title="管理后台" class="text-error"></v-list-item>
          <v-divider class="my-2"></v-divider>
          <v-list-item @click="logout" prepend-icon="mdi-logout" title="退出" class="text-secondary"></v-list-item>
        </template>
      </v-list>
    </v-navigation-drawer>

    <!-- 根据路由meta决定是否使用v-main和v-container包装 -->
    <v-main v-if="!hideAppBar" class="bg-grey-lighten-4">
      <v-container fluid class="pa-4">
        <router-view />
      </v-container>
    </v-main>

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

    <!-- 独立页面直接渲染router-view -->
    <router-view v-else />

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
    const chatUnreadCount = ref(0)
    const drawer = ref(false)
    const isMobile = ref(false)
    const hideAppBar = ref(false)
    const backendVersion = ref(null)
    const siteTitle = ref('校园论坛')
    const frontendVersion = ref(typeof __FRONTEND_VERSION__ !== 'undefined' ? __FRONTEND_VERSION__ : 'unknown')
    const icpNumber = ref(null)
    const publicSecurityNumber = ref(null)
    let pollInterval = null
    
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
    
    const loadChatUnreadCount = async () => {
      if (!token.value) return
      try {
        const response = await api.get('/chat/unread-count')
        chatUnreadCount.value = response.data.unread_count || 0
      } catch (error) {
        console.error('加载未读消息数量失败', error)
      }
    }
    
    const goToChat = () => {
      router.push('/chat')
    }
    
    const logout = () => {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.dispatchEvent(new Event('user-logout'))
      drawer.value = false
      router.push('/login')
    }
    
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
      loadChatUnreadCount()
      loadVersion()
      loadSiteTitle()

      // 每分钟刷新未读消息数量
      pollInterval = setInterval(loadChatUnreadCount, 60000)

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
      logout,
      modalState,
      handleConfirm,
      handleCancel,
      chatUnreadCount,
      goToChat,
      isMobile,
      drawer,
      backendVersion,
      siteTitle,
      updateSiteTitle,
      frontendVersion,
      icpNumber,
      publicSecurityNumber
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
</style>
