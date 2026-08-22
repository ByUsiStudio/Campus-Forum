<template>
  <div class="app-shell">
    <!-- 顶部导航 -->
    <header class="app-header">
      <div class="app-header-inner">
        <div class="app-brand">
          <el-icon class="brand-icon"><Forum /></el-icon>
          <router-link :to="isAdminPage ? '/admin' : '/'" class="brand-title">
            {{ isAdminPage ? '管理后台' : appStore.siteTitle }}
          </router-link>
        </div>

        <!-- 桌面端导航 -->
        <nav class="app-nav" v-if="!isAdminPage">
          <el-button v-if="!userStore.isLoggedIn" text @click="go('/login')">登录</el-button>
          <el-button v-if="!userStore.isLoggedIn" text @click="go('/register')">注册</el-button>
          <el-button v-if="userStore.isLoggedIn" text @click="go('/create')">
            <el-icon><EditPen /></el-icon>&nbsp;写文章
          </el-button>
          <el-button v-if="userStore.isLoggedIn" text @click="go('/collections')">收藏</el-button>
          <el-button v-if="userStore.isLoggedIn" text @click="go('/signin')">签到</el-button>
          <el-button v-if="userStore.isLoggedIn" text @click="go('/topics')">话题</el-button>
          <el-button v-if="userStore.isAdmin" text @click="go('/admin')">
            <el-icon><Setting /></el-icon>&nbsp;管理后台
          </el-button>
        </nav>

        <!-- 移动端导航触发 -->
        <el-button v-if="!isAdminPage" class="mobile-menu-toggle d-lg-none" text circle @click="mobileMenuOpen = !mobileMenuOpen" :aria-label="mobileMenuOpen ? '收起菜单' : '打开菜单'">
          <el-icon :size="20"><Menu v-if="!mobileMenuOpen" /><Close v-else /></el-icon>
        </el-button>

        <div class="app-actions">
          <el-input
            v-model="searchQuery"
            placeholder="搜索…"
            clearable
            class="search-input"
            @keyup.enter="handleSearch"
          >
            <template #prefix><el-icon><Search /></el-icon></template>
          </el-input>

          <el-badge v-if="userStore.isLoggedIn" :value="notificationStore.unreadCount" :hidden="notificationStore.unreadCount === 0" :offset="[-4,4]">
            <el-button circle text @click="go('/notifications')">
              <el-icon :size="20"><Bell /></el-icon>
            </el-button>
          </el-badge>

          <template v-if="userStore.user">
            <el-dropdown trigger="click" @command="onUserCommand">
              <span class="user-trigger cursor-pointer">
                <el-avatar :size="32" :src="userStore.user.avatar">
                  {{ (userStore.user.display_name || userStore.user.username || 'U')[0] }}
                </el-avatar>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="profile">
                    <el-icon><User /></el-icon>个人中心
                  </el-dropdown-item>
                  <el-dropdown-item command="logout" divided>
                    <el-icon><SwitchButton /></el-icon>退出登录
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
          <template v-else>
            <el-button type="primary" @click="go('/login')">登录</el-button>
          </template>
        </div>
      </div>
    </header>

    <!-- 移动端下拉导航 -->
    <transition name="slide-down">
      <div class="mobile-menu-panel" v-if="mobileMenuOpen && !isAdminPage" @click="mobileMenuOpen = false">
        <div class="mobile-menu-links">
          <template v-if="!userStore.isLoggedIn">
            <a class="mobile-menu-item" @click="go('/login')"><el-icon><User /></el-icon>登录</a>
            <a class="mobile-menu-item" @click="go('/register')"><el-icon><Postcard /></el-icon>注册</a>
          </template>
          <template v-else>
            <a class="mobile-menu-item" @click="go('/create')"><el-icon><EditPen /></el-icon>写文章</a>
            <a class="mobile-menu-item" @click="go('/collections')"><el-icon><Folder /></el-icon>我的收藏</a>
            <a class="mobile-menu-item" @click="go('/signin')"><el-icon><Calendar /></el-icon>每日签到</a>
            <a class="mobile-menu-item" @click="go('/topics')"><el-icon><Collection /></el-icon>话题</a>
            <a class="mobile-menu-item" @click="go('/profile')"><el-icon><User /></el-icon>个人中心</a>
          </template>
          <a class="mobile-menu-item" @click="go('/leaderboard')"><el-icon><Trophy /></el-icon>排行榜</a>
          <a class="mobile-menu-item" @click="go('/notifications')"><el-icon><Bell /></el-icon>我的通知</a>
          <a v-if="userStore.isAdmin" class="mobile-menu-item" @click="go('/admin')"><el-icon><Setting /></el-icon>管理后台</a>
        </div>
      </div>
    </transition>

    <!-- 内容区 -->
    <main class="app-main">
      <router-view />
    </main>

    <!-- 页脚 -->
    <footer class="app-footer">
      <div class="app-footer-inner">
        <span class="text-secondary">{{ appStore.siteTitle }}</span>
        <span class="text-secondary">
          前端 {{ appStore.frontendVersion || '—' }} / 后端 {{ appStore.backendVersion || '—' }}
        </span>
        <div v-if="appStore.icpNumber || appStore.publicSecurityNumber" class="footer-icp">
          <span v-if="appStore.icpNumber">{{ appStore.icpNumber }}</span>
          <span v-if="appStore.publicSecurityNumber">{{ appStore.publicSecurityNumber }}</span>
        </div>
      </div>
    </footer>

    <!-- 公告弹窗 -->
    <el-dialog
      v-model="announcementVisible"
      title="公告"
      width="560px"
      append-to-body
    >
      <div class="markdown-body" v-if="announcementContent" v-html="announcementContent"></div>
      <template #footer>
        <el-checkbox v-model="dontShowAgain">不再显示</el-checkbox>
        <el-button type="primary" @click="closeAnnouncement">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import MarkdownIt from 'markdown-it'
import { useAppStore } from '@/stores/app'
import { useUserStore } from '@/stores/user'
import { useNotificationStore } from '@/stores/notification'
import http from '@/api/http'

// 公告使用 markdown-it 渲染
const renderMd = new MarkdownIt({ html: true, linkify: true, breaks: true })

const route = useRoute()
const router = useRouter()

const appStore = useAppStore()
const userStore = useUserStore()
const notificationStore = useNotificationStore()

const searchQuery = ref('')
const announcementVisible = ref(false)
const announcementContent = ref('')
const dontShowAgain = ref(false)
const mobileMenuOpen = ref(false)

const isAdminPage = computed(() => route.path.startsWith('/admin'))

const go = (path) => router.push(path)

const handleSearch = () => {
  const q = searchQuery.value.trim()
  if (q) router.push({ path: '/search', query: { q } })
}

// 路由变化时关闭移动端菜单
watch(() => route.path, () => {
  mobileMenuOpen.value = false
})

const onUserCommand = (cmd) => {
  if (cmd === 'profile') router.push('/profile')
  if (cmd === 'logout') {
    userStore.logout()
    router.push('/login')
  }
}

const closeAnnouncement = () => {
  if (dontShowAgain.value) localStorage.setItem('hideAnnouncement', 'true')
  announcementVisible.value = false
}

const loadAnnouncement = async () => {
  try {
    const { data } = await http.get('/announcement')
    if (data.content) {
      announcementContent.value = renderMd.render(data.content)
      if (!localStorage.getItem('hideAnnouncement')) {
        announcementVisible.value = true
      }
    }
  } catch (e) { /* ignore */ }
}

onMounted(() => {
  appStore.initApp()
  userStore.initUser()
  if (userStore.isLoggedIn) notificationStore.startPolling()
  loadAnnouncement()
})

onUnmounted(() => {
  notificationStore.stopPolling()
})
</script>

<style scoped>
.app-shell {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}
.app-header {
  position: sticky;
  top: 0;
  z-index: 100;
  background: #fff;
  border-bottom: 1px solid var(--campus-border);
  box-shadow: 0 1px 6px rgba(0, 0, 0, 0.04);
}
.app-header-inner {
  max-width: 1200px;
  margin: 0 auto;
  padding: 10px 16px;
  display: flex;
  align-items: center;
  gap: 16px;
}
.app-brand {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}
.brand-icon {
  color: var(--campus-primary);
  font-size: 22px;
}
.brand-title {
  font-size: 18px;
  font-weight: 700;
  color: var(--campus-text);
  white-space: nowrap;
}
.app-nav {
  display: flex;
  align-items: center;
  gap: 4px;
  flex: 1;
  overflow-x: auto;
}
.app-actions {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-shrink: 0;
}
.search-input {
  width: 220px;
}
.user-trigger {
  display: flex;
  align-items: center;
  line-height: 1;
}
.app-main {
  flex: 1;
}
.app-footer {
  border-top: 1px solid var(--campus-border);
  background: #fff;
}
.app-footer-inner {
  max-width: 1200px;
  margin: 0 auto;
  padding: 16px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  font-size: 13px;
}
.footer-icp {
  display: flex;
  gap: 16px;
}
.mobile-menu-toggle {
  display: none;
  color: var(--campus-text);
}
.mobile-menu-panel {
  position: sticky;
  top: 57px;
  z-index: 99;
  background: #fff;
  border-bottom: 1px solid var(--campus-border);
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.06);
}
.mobile-menu-links {
  display: flex;
  flex-direction: column;
  max-width: 1200px;
  margin: 0 auto;
  padding: 8px;
}
.mobile-menu-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 14px;
  border-radius: var(--campus-radius);
  color: var(--campus-text);
  font-size: 15px;
  cursor: pointer;
}
.mobile-menu-item:active {
  background: var(--campus-bg);
}
.mobile-menu-item .el-icon {
  color: var(--campus-primary);
}

/* 下拉动画 */
.slide-down-enter-active,
.slide-down-leave-active {
  transition: all 0.2s ease;
}
.slide-down-enter-from,
.slide-down-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

@media (max-width: 768px) {
  .app-nav { display: none !important; }
  .mobile-menu-toggle { display: flex !important; }
  .search-input { width: 140px; }
  .app-header-inner { gap: 8px; }
}
</style>
