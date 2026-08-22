<template>
  <div class="app-shell">
    <!-- ===================== 顶部导航 ===================== -->
    <header class="app-header" :class="{ 'is-admin': isAdminPage }">
      <div class="app-header-inner">
        <!-- 品牌 -->
        <div class="app-brand">
          <span class="brand-mark">
            <el-icon :size="22"><component :is="isAdminPage ? 'Setting' : 'Forum'" /></el-icon>
          </span>
          <router-link :to="isAdminPage ? '/admin' : '/'" class="brand-title">
            {{ isAdminPage ? '管理后台' : appStore.siteTitle }}
          </router-link>
        </div>

        <!-- 桌面导航 -->
        <nav class="app-nav" v-if="!isAdminPage">
          <router-link
            v-for="item in navItems"
            :key="item.path"
            :to="item.path"
            class="nav-link"
            :class="{ 'is-active': route.path === item.path || (item.exact && route.path.startsWith(item.path)) }"
          >
            <el-icon v-if="item.icon" :size="15"><component :is="item.icon" /></el-icon>
            <span>{{ item.label }}</span>
          </router-link>
        </nav>

        <div class="app-actions">
          <!-- 桌面搜索 -->
          <div class="search-box">
            <el-input
              v-model="searchQuery"
              placeholder="搜索文章…"
              clearable
              class="search-input-desktop"
              @keyup.enter="handleSearch"
            >
              <template #prefix><el-icon><Search /></el-icon></template>
            </el-input>
          </div>
          <el-button class="search-mobile-btn" circle text aria-label="搜索" @click="go('/search')">
            <el-icon :size="20"><Search /></el-icon>
          </el-button>

          <!-- 写文章 -->
          <el-button
            v-if="userStore.isLoggedIn && !isAdminPage"
            class="write-btn"
            type="primary"
            @click="go('/create')"
          >
            <el-icon :size="15"><EditPen /></el-icon>
            <span class="write-txt">写文章</span>
          </el-button>

          <!-- 通知铃铛 -->
          <el-badge
            v-if="userStore.isLoggedIn"
            :value="notificationStore.unreadCount"
            :hidden="notificationStore.unreadCount === 0"
            :offset="[-4, 4]"
          >
            <el-button class="icon-btn" circle text @click="go('/notifications')">
              <el-icon :size="20"><Bell /></el-icon>
            </el-button>
          </el-badge>

          <!-- 用户 -->
          <template v-if="userStore.user">
            <el-dropdown trigger="click" @command="onUserCommand">
              <span class="user-trigger cursor-pointer">
                <el-avatar :size="32" :src="userStore.user.avatar" class="user-avatar">
                  {{ (userStore.user.display_name || userStore.user.username || 'U')[0] }}
                </el-avatar>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="profile">
                    <el-icon><User /></el-icon>个人中心
                  </el-dropdown-item>
                  <el-dropdown-item v-if="userStore.isAdmin" command="admin" divided>
                    <el-icon><Setting /></el-icon>管理后台
                  </el-dropdown-item>
                  <el-dropdown-item command="logout" divided>
                    <el-icon><SwitchButton /></el-icon>退出登录
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
          <template v-else>
            <el-button class="desktop-login" plain round @click="go('/login')">登录</el-button>
          </template>

          <!-- 移动端菜单按钮 -->
          <el-button
            v-if="!isAdminPage"
            class="mobile-toggle"
            circle
            text
            @click="mobileMenuOpen = !mobileMenuOpen"
            :aria-label="mobileMenuOpen ? '收起菜单' : '打开菜单'"
          >
            <el-icon :size="22"><Menu v-if="!mobileMenuOpen" /><Close v-else /></el-icon>
          </el-button>
        </div>
      </div>
    </header>

    <!-- 移动端抽屉 -->
    <transition name="slide-in">
      <div class="mobile-drawer-mask" v-if="mobileMenuOpen && !isAdminPage" @click="mobileMenuOpen = false"></div>
    </transition>
    <transition name="drawer">
      <aside class="mobile-drawer" v-if="mobileMenuOpen && !isAdminPage">
        <div class="drawer-head">
          <div class="drawer-profile" v-if="userStore.user">
            <el-avatar :size="46" :src="userStore.user.avatar">
              {{ (userStore.user.display_name || userStore.user.username || 'U')[0] }}
            </el-avatar>
            <div>
              <div class="drawer-name">{{ userStore.user.display_name }}</div>
              <div class="drawer-sub">@{{ userStore.user.username }}</div>
            </div>
          </div>
          <div v-else class="drawer-guest">登录后解锁更多功能</div>
          <el-button class="drawer-close" circle text @click="mobileMenuOpen = false">
            <el-icon><Close /></el-icon>
          </el-button>
        </div>

        <nav class="drawer-nav">
          <router-link
            v-for="item in drawerItems"
            :key="item.path"
            :to="item.path"
            class="drawer-item"
            @click="mobileMenuOpen = false"
          >
            <el-icon :size="19"><component :is="item.icon" /></el-icon>
            <span>{{ item.label }}</span>
          </router-link>
        </nav>

        <div class="drawer-auth" v-if="!userStore.isLoggedIn">
          <el-button type="primary" round @click="go('/login')">登录</el-button>
          <el-button round @click="go('/register')">注册</el-button>
        </div>
        <div class="drawer-auth" v-else>
          <el-button type="danger" plain round @click="onLogout">退出登录</el-button>
        </div>
      </aside>
    </transition>

    <!-- ===================== 主内容 ===================== -->
    <main class="app-main">
      <router-view />
    </main>

    <!-- ===================== 页脚 ===================== -->
    <footer class="app-footer">
      <div class="app-footer-inner">
        <span class="footer-brand">{{ appStore.siteTitle }}</span>
        <span class="text-muted footer-version">
          前端 {{ appStore.frontendVersion || '—' }} / 后端 {{ appStore.backendVersion || '—' }}
        </span>
        <div v-if="appStore.icpNumber || appStore.publicSecurityNumber" class="footer-icp">
          <span v-if="appStore.icpNumber">{{ appStore.icpNumber }}</span>
          <span v-if="appStore.publicSecurityNumber">{{ appStore.publicSecurityNumber }}</span>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted, h } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import MarkdownIt from 'markdown-it'
import { useAppStore } from '@/stores/app'
import { useUserStore } from '@/stores/user'
import { useNotificationStore } from '@/stores/notification'
import http from '@/api/http'

const renderMd = new MarkdownIt({ html: true, linkify: true, breaks: true })

const route = useRoute()
const router = useRouter()

const appStore = useAppStore()
const userStore = useUserStore()
const notificationStore = useNotificationStore()

const searchQuery = ref('')
const mobileMenuOpen = ref(false)

const isAdminPage = computed(() => route.path.startsWith('/admin'))

const navItems = computed(() => {
  const items = [
    { path: '/', label: '首页', icon: 'Home', exact: false },
    { path: '/topics', label: '话题', icon: 'Collection' },
    { path: '/leaderboard', label: '排行榜', icon: 'Trophy' }
  ]
  return items
})

const drawerItems = computed(() => {
  const items = [
    { path: '/', label: '首页', icon: 'Home' },
    { path: '/topics', label: '话题', icon: 'Collection' },
    { path: '/leaderboard', label: '排行榜', icon: 'Trophy' }
  ]
  if (userStore.isLoggedIn) {
    items.push(
      { path: '/create', label: '写文章', icon: 'EditPen' },
      { path: '/collections', label: '我的收藏', icon: 'Folder' },
      { path: '/signin', label: '每日签到', icon: 'Calendar' },
      { path: '/profile', label: '个人中心', icon: 'User' },
      { path: '/notifications', label: '我的通知', icon: 'Bell' }
    )
  } else {
    items.push({ path: '/login', label: '登录', icon: 'User' })
  }
  if (userStore.isAdmin) items.push({ path: '/admin', label: '管理后台', icon: 'Setting' })
  return items
})

const go = (path) => router.push(path)

const handleSearch = () => {
  const q = searchQuery.value.trim()
  if (q) router.push({ path: '/search', query: { q } })
  else router.push('/search')
}

watch(() => route.path, () => {
  mobileMenuOpen.value = false
})

const onLogout = () => {
  mobileMenuOpen.value = false
  userStore.logout()
  router.push('/login')
}

const onUserCommand = (cmd) => {
  if (cmd === 'profile') router.push('/profile')
  if (cmd === 'admin') router.push('/admin')
  if (cmd === 'logout') onLogout()
}

// ---------------- 公告（JCuPupw）----------------
async function loadAnnouncement() {
  try {
    const { data } = await http.get('/announcement')
    if (data && data.content && !localStorage.getItem('hideAnnouncement')) {
      showAnnouncement(renderMd.render(data.content))
    }
  } catch (e) { /* ignore */ }
}

function showAnnouncement(html) {
  const jc = window.JCuPupw
  if (!jc) return
  const notShowBtn = {
    text: '不再显示',
    type: 'default',
    action: () => {
      localStorage.setItem('hideAnnouncement', 'true')
      jc.closeAll()
    }
  }
  const closeBtn = {
    text: '知道了',
    type: 'primary',
    action: () => jc.closeAll()
  }
  jc.instance().open({
    title: '公告',
    content: `<div class="markdown-body" style="max-height:60vh;overflow:auto;">${html}</div>`,
    width: 560,
    size: 'md',
    buttons: [notShowBtn, closeBtn]
  })
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
  background: var(--campus-bg);
}

/* ---------- Header ---------- */
.app-header {
  position: sticky;
  top: 0;
  z-index: 100;
  background: rgba(255, 255, 255, 0.92);
  backdrop-filter: saturate(180%) blur(12px);
  -webkit-backdrop-filter: saturate(180%) blur(12px);
  border-bottom: 1px solid var(--campus-border);
  height: var(--campus-header-h);
  display: flex;
  align-items: center;
}
.app-header-inner {
  max-width: var(--campus-max-width);
  margin: 0 auto;
  padding: 0 var(--campus-gutter);
  width: 100%;
  display: flex;
  align-items: center;
  gap: 20px;
  box-sizing: border-box;
}
.app-brand {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-shrink: 0;
}
.brand-mark {
  width: 34px;
  height: 34px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  background: linear-gradient(135deg, var(--campus-primary), var(--campus-primary-light));
  box-shadow: 0 4px 10px rgba(79, 110, 247, 0.3);
}
.brand-title {
  font-size: 18px;
  font-weight: 800;
  color: var(--campus-text);
  white-space: nowrap;
  letter-spacing: -0.01em;
}

/* ---------- Nav ---------- */
.app-nav {
  display: flex;
  gap: 4px;
  flex: 1;
  align-items: center;
  overflow-x: auto;
}
.nav-link {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 7px 14px;
  border-radius: 10px;
  color: var(--campus-text-secondary);
  font-size: 14px;
  font-weight: 500;
  white-space: nowrap;
  transition: var(--campus-transition);
}
.nav-link:hover {
  color: var(--campus-primary);
  background: var(--campus-primary-soft);
}
.nav-link.is-active {
  color: var(--campus-primary);
  background: var(--campus-primary-soft);
  font-weight: 600;
}

.app-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}
.search-box {
  display: flex;
  align-items: center;
}
.search-input-desktop {
  width: 200px;
  --el-input-border-radius: 20px;
}
.search-mobile-btn {
  display: none;
  color: var(--campus-text);
}
.icon-btn {
  color: var(--campus-text);
}
.user-trigger {
  display: flex;
  align-items: center;
  line-height: 1;
}
.user-avatar {
  border: 2px solid #fff;
  box-shadow: 0 0 0 1px var(--campus-border);
}
.desktop-login {
  color: var(--campus-primary);
}
.write-btn {
  border-radius: 10px;
}

/* ---------- Mobile Drawer ---------- */
.mobile-toggle {
  display: none;
  color: var(--campus-text);
}
.mobile-drawer-mask {
  position: fixed;
  inset: 0;
  z-index: 199;
  background: rgba(15, 23, 42, 0.4);
  backdrop-filter: blur(2px);
}
.mobile-drawer {
  position: fixed;
  top: 0;
  right: 0;
  bottom: 0;
  width: 280px;
  max-width: 85vw;
  z-index: 200;
  background: var(--campus-surface);
  box-shadow: -8px 0 30px rgba(15, 23, 42, 0.15);
  display: flex;
  flex-direction: column;
  padding: 20px 16px;
  box-sizing: border-box;
}
.drawer-head {
  display: flex;
  align-items: center;
  gap: 12px;
  padding-bottom: 16px;
  border-bottom: 1px solid var(--campus-border);
}
.drawer-profile {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
}
.drawer-name {
  font-weight: 700;
}
.drawer-sub {
  font-size: 12px;
  color: var(--campus-text-secondary);
}
.drawer-guest {
  flex: 1;
  color: var(--campus-text-secondary);
}
.drawer-close {
  color: var(--campus-text);
}
.drawer-nav {
  display: flex;
  flex-direction: column;
  gap: 4px;
  margin-top: 16px;
  flex: 1;
  overflow-y: auto;
}
.drawer-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 12px;
  border-radius: 10px;
  color: var(--campus-text);
  font-size: 15px;
  font-weight: 500;
  transition: var(--campus-transition);
}
.drawer-item:hover {
  color: var(--campus-primary);
  background: var(--campus-primary-soft);
}
.drawer-item .el-icon {
  color: var(--campus-primary);
}
.drawer-auth {
  display: flex;
  gap: 10px;
  padding-top: 16px;
  border-top: 1px solid var(--campus-border);
}
.drawer-auth .el-button {
  flex: 1;
}

/* ---------- Main / Footer ---------- */
.app-main {
  flex: 1;
  display: flex;
  flex-direction: column;
}
.app-footer {
  border-top: 1px solid var(--campus-border);
  background: var(--campus-surface);
  margin-top: 32px;
}
.app-footer-inner {
  max-width: var(--campus-max-width);
  margin: 0 auto;
  padding: 20px var(--campus-gutter);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  font-size: 13px;
}
.footer-brand {
  font-weight: 700;
  color: var(--campus-text);
}
.footer-version {
  font-size: 12px;
}
.footer-icp {
  display: flex;
  gap: 16px;
  color: var(--campus-text-muted);
  font-size: 12px;
}

/* ---------- Animations ---------- */
.slide-in-enter-active,
.slide-in-leave-active {
  transition: opacity 0.25s ease;
}
.slide-in-enter-from,
.slide-in-leave-to {
  opacity: 0;
}
.drawer-enter-active,
.drawer-leave-active {
  transition: transform 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}
.drawer-enter-from,
.drawer-leave-to {
  transform: translateX(100%);
}

/* ---------- Responsive ---------- */
@media (max-width: 768px) {
  .app-nav { display: none !important; }
  .mobile-toggle { display: inline-flex !important; }
  .search-input-desktop { display: none !important; }
  .search-mobile-btn { display: inline-flex !important; }
  .write-txt { display: none; }
  .desktop-login { display: none !important; }
  .brand-title { font-size: 16px; }
  .app-header-inner { gap: 10px; }
}
</style>
