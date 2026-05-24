<template>
  <v-app>
    <!-- 根据路由meta决定是否显示顶部栏 -->
    <v-app-bar v-if="!hideAppBar" elevation="1" color="surface">
      <v-app-bar-nav-icon v-if="isMobile" @click="drawer = !drawer"></v-app-bar-nav-icon>
      <v-app-bar-title>
        <router-link to="/" class="logo-link">校园论坛</router-link>
      </v-app-bar-title>

      <!-- 桌面端导航 -->
      <template v-if="!isMobile" v-slot:append>
        <div class="nav-links">
          <v-btn variant="text" to="/" color="primary">首页</v-btn>
          <v-btn variant="text" to="/login" v-if="!token" color="primary">登录</v-btn>
          <v-btn variant="text" to="/register" v-if="!token" color="primary">注册</v-btn>
          <v-btn variant="text" to="/create" v-if="token" color="primary">写文章</v-btn>
          <v-btn variant="text" to="/profile" v-if="token" color="primary">个人中心</v-btn>
          <!-- 临时去除聊天界面的入口 -->
          <!-- <v-btn variant="text" v-if="token" @click="goToChat">
            消息
          </v-btn> -->
          <NotificationBell v-if="token" />
          <v-btn variant="text" to="/admin" v-if="isAdmin" color="error">管理后台</v-btn>
          <v-btn variant="text" v-if="token" @click="logout" color="secondary">退出</v-btn>
        </div>
      </template>

      <!-- 移动端导航 -->
      <template v-else v-slot:append>
        <!-- 临时去除聊天界面的入口 -->
        <!-- <v-btn icon @click="goToChat" class="mr-2">
          <v-icon>mdi-message</v-icon>
        </v-btn> -->
        <NotificationBell v-if="token" />
      </template>
    </v-app-bar>

    <!-- 移动端抽屉导航 -->
    <v-navigation-drawer v-if="!hideAppBar" v-model="drawer" temporary>
      <v-list>
        <v-list-item to="/" @click="drawer = false">
          <v-list-item-title>首页</v-list-item-title>
        </v-list-item>
        <template v-if="!token">
          <v-list-item to="/login" @click="drawer = false">
            <v-list-item-title>登录</v-list-item-title>
          </v-list-item>
          <v-list-item to="/register" @click="drawer = false">
            <v-list-item-title>注册</v-list-item-title>
          </v-list-item>
        </template>
        <template v-if="token">
          <v-list-item to="/create" @click="drawer = false">
            <v-list-item-title>写文章</v-list-item-title>
          </v-list-item>
          <v-list-item to="/profile" @click="drawer = false">
            <v-list-item-title>个人中心</v-list-item-title>
          </v-list-item>
          <v-list-item v-if="isAdmin" to="/admin" @click="drawer = false">
            <v-list-item-title class="text-error">管理后台</v-list-item-title>
          </v-list-item>
          <v-list-item @click="logout">
            <v-list-item-title class="text-secondary">退出</v-list-item-title>
          </v-list-item>
        </template>
      </v-list>
    </v-navigation-drawer>

    <!-- 根据路由meta决定是否使用v-main和v-container包装 -->
    <v-main v-if="!hideAppBar">
      <v-container fluid>
        <router-view />
      </v-container>
    </v-main>
    
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
      return user.value && user.value.role === 'admin'
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
    
    onMounted(() => {
      checkMobile()
      window.addEventListener('resize', checkMobile)
      
      const storedUser = localStorage.getItem('user')
      if (storedUser) {
        user.value = JSON.parse(storedUser)
      }
      loadUser()
      loadChatUnreadCount()

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
      drawer
    }
  }
}
</script>

<style>
.logo-link {
  color: rgb(var(--v-theme-primary));
  text-decoration: none;
  font-weight: bold;
  font-size: 1.25rem;
}

.nav-links {
  display: flex;
  gap: 4px;
  align-items: center;
}

.chat-badge {
  position: absolute;
  top: -4px;
  right: -8px;
}
</style>
