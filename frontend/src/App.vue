<template>
  <div class="app-container">
    <div v-if="announcementDialog" class="announcement-overlay" @click="closeAnnouncementDialog">
      <div class="announcement-modal animate-scale-in" @click.stop>
        <div class="announcement-header">
          <span class="announcement-title">
            <i class="layui-icon layui-icon-notice"></i>
            公告
          </span>
          <button class="close-btn" @click="closeAnnouncementDialog">
            <i class="layui-icon layui-icon-close"></i>
          </button>
        </div>
        <div class="announcement-content">
          <div class="markdown-body" v-html="announcementContent"></div>
        </div>
        <div class="announcement-footer">
          <label class="checkbox-label">
            <input type="checkbox" v-model="dontShowAgain">
            <span>不再显示</span>
          </label>
          <button class="layui-btn layui-btn-primary" @click="closeAnnouncementDialog">关闭</button>
        </div>
      </div>
    </div>

    <header v-if="!hideAppBar" class="navbar" :class="{ 'scrolled': isScrolled }">
      <div class="navbar-container">
        <div class="navbar-brand">
          <img src="/xylt.svg" alt="Logo" class="logo-img">
          <router-link :to="isAdminPage ? '/admin' : '/'" class="logo-text">
            {{ isAdminPage ? '管理后台' : siteTitle }}
          </router-link>
        </div>

        <div class="nav-search" v-if="!isAdminPage">
          <div class="layui-input-group">
            <input 
              type="text" 
              v-model="searchQuery" 
              placeholder="搜索..." 
              class="layui-input search-input"
              @keyup.enter="handleSearch"
            >
            <button class="layui-btn search-btn" @click="handleSearch">
              <i class="layui-icon layui-icon-search"></i>
            </button>
          </div>
        </div>

        <nav class="nav-links" v-if="!isMobile">
          <template v-if="!isAdminPage">
            <router-link to="/" class="nav-link">
              <i class="layui-icon layui-icon-home"></i>
              首页
            </router-link>
            <router-link to="/login" v-if="!store.isLoggedIn.value" class="nav-link">
              <i class="layui-icon layui-icon-username"></i>
              登录
            </router-link>
            <router-link to="/register" v-if="!store.isLoggedIn.value" class="nav-link">
              <i class="layui-icon layui-icon-user"></i>
              注册
            </router-link>
            <router-link to="/create" v-if="store.isLoggedIn.value" class="nav-link nav-link-primary">
              <i class="layui-icon layui-icon-edit"></i>
              写文章
            </router-link>
            <router-link to="/notifications" v-if="store.isLoggedIn.value" class="nav-link notification-link">
              <i class="layui-icon layui-icon-notice"></i>
              通知
              <span v-if="store.unreadCount.value > 0" class="badge">{{ store.unreadCount.value }}</span>
            </router-link>
            <router-link to="/profile" v-if="store.isLoggedIn.value" class="nav-link">
              <i class="layui-icon layui-icon-user"></i>
              我的
            </router-link>
            <router-link to="/admin" v-if="store.isAdmin.value" class="nav-link nav-link-admin">
              <i class="layui-icon layui-icon-set"></i>
              管理后台
            </router-link>
            <button class="nav-link logout-btn" v-if="store.isLoggedIn.value" @click="logout">
              <i class="layui-icon layui-icon-logout"></i>
              退出
            </button>
          </template>
          <template v-else>
            <router-link :to="{ name: 'AdminIndex' }" class="nav-link">
              <i class="layui-icon layui-icon-home"></i>
              控制台
            </router-link>
            <router-link :to="{ name: 'AdminUsers' }" class="nav-link">
              <i class="layui-icon layui-icon-user"></i>
              用户管理
            </router-link>
            <router-link :to="{ name: 'AdminArticles' }" class="nav-link">
              <i class="layui-icon layui-icon-file"></i>
              文章管理
            </router-link>
            <router-link :to="{ name: 'AdminSiteConfig' }" class="nav-link">
              <i class="layui-icon layui-icon-set"></i>
              系统设置
            </router-link>
            <router-link to="/" class="nav-link">
              <i class="layui-icon layui-icon-home"></i>
              返回首页
            </router-link>
          </template>
        </nav>

        <button class="mobile-menu-btn" @click="toggleMobileMenu">
          <i class="layui-icon" :class="mobileMenuOpen ? 'layui-icon-close' : 'layui-icon-menu'"></i>
        </button>
      </div>

      <div v-if="mobileMenuOpen" class="mobile-menu animate-slide-down">
        <div class="mobile-menu-header">
          <span>菜单</span>
          <button @click="mobileMenuOpen = false">
            <i class="layui-icon layui-icon-close"></i>
          </button>
        </div>
        <div class="mobile-menu-content">
          <router-link to="/" @click="mobileMenuOpen = false" class="mobile-menu-item">
            <i class="layui-icon layui-icon-home"></i>
            首页
          </router-link>
          <router-link to="/login" v-if="!store.isLoggedIn.value" @click="mobileMenuOpen = false" class="mobile-menu-item">
            <i class="layui-icon layui-icon-username"></i>
            登录
          </router-link>
          <router-link to="/register" v-if="!store.isLoggedIn.value" @click="mobileMenuOpen = false" class="mobile-menu-item">
            <i class="layui-icon layui-icon-user"></i>
            注册
          </router-link>
          <router-link to="/create" v-if="store.isLoggedIn.value" @click="mobileMenuOpen = false" class="mobile-menu-item">
            <i class="layui-icon layui-icon-edit"></i>
            写文章
          </router-link>
          <router-link to="/notifications" v-if="store.isLoggedIn.value" @click="mobileMenuOpen = false" class="mobile-menu-item">
            <i class="layui-icon layui-icon-notice"></i>
            通知
          </router-link>
          <router-link to="/signin" v-if="store.isLoggedIn.value" @click="mobileMenuOpen = false" class="mobile-menu-item">
            <i class="layui-icon layui-icon-calendar"></i>
            签到
          </router-link>
          <router-link to="/leaderboard" v-if="store.isLoggedIn.value" @click="mobileMenuOpen = false" class="mobile-menu-item">
            <i class="layui-icon layui-icon-trophy"></i>
            排行榜
          </router-link>
          <router-link to="/topics" v-if="store.isLoggedIn.value" @click="mobileMenuOpen = false" class="mobile-menu-item">
            <i class="layui-icon layui-icon-tags"></i>
            话题
          </router-link>
          <router-link to="/collections" v-if="store.isLoggedIn.value" @click="mobileMenuOpen = false" class="mobile-menu-item">
            <i class="layui-icon layui-icon-collection"></i>
            收藏夹
          </router-link>
          <router-link to="/admin" v-if="store.isAdmin.value" @click="mobileMenuOpen = false" class="mobile-menu-item admin-item">
            <i class="layui-icon layui-icon-set"></i>
            管理后台
          </router-link>
          <button v-if="store.isLoggedIn.value" @click="logout; mobileMenuOpen = false" class="mobile-menu-item logout-item">
            <i class="layui-icon layui-icon-logout"></i>
            退出登录
          </button>
        </div>
      </div>
    </header>

    <main class="main-content" :class="{ 'has-navbar': !hideAppBar }">
      <router-view />
    </main>

    <footer v-if="!hideAppBar" class="footer">
      <div class="footer-content">
        <div class="footer-links">
          <div class="footer-column">
            <h4>关于我们</h4>
            <ul>
              <li><a href="https://github.com/ByUsiStudio/Campus-Forum" target="_blank">GitHub</a></li>
              <li><a href="https://gitee.com/byusistudio/campus-forum" target="_blank">Gitee</a></li>
            </ul>
          </div>
          <div class="footer-column">
            <h4>联系方式</h4>
            <ul>
              <li>Email: contact@byusi.com</li>
            </ul>
          </div>
        </div>
        <div class="footer-bottom">
          <div class="version-info">
            <span>前端版本: {{ frontendVersion }}</span>
            <span> | </span>
            <span>后端版本: {{ backendVersion || 'unknown' }}</span>
          </div>
          <div v-if="icpNumber || publicSecurityNumber" class="cert-info">
            <span v-if="icpNumber">
              <a href="https://beian.miit.gov.cn" target="_blank">{{ icpNumber }}</a>
            </span>
            <span v-if="publicSecurityNumber" class="security-num">{{ publicSecurityNumber }}</span>
          </div>
        </div>
      </div>
    </footer>

    <div v-if="modalState.show" class="modal-overlay" @click="modalState.show = false">
      <div class="modal-content animate-scale-in" @click.stop>
        <div class="modal-header">
          <span class="modal-title">{{ modalState.title }}</span>
          <button @click="modalState.show = false">
            <i class="layui-icon layui-icon-close"></i>
          </button>
        </div>
        <div class="modal-body">
          <div v-if="modalState.icon" class="modal-icon" :class="modalState.iconColor">
            <i class="layui-icon" :class="modalState.icon"></i>
          </div>
          <p v-if="modalState.message">{{ modalState.message }}</p>
          <textarea 
            v-if="modalState.inputRows"
            v-model="modalInputValue"
            :rows="modalState.inputRows"
            :placeholder="modalState.inputPlaceholder"
            class="layui-textarea modal-textarea"
          ></textarea>
          <input 
            v-else-if="modalState.inputLabel"
            type="text"
            v-model="modalInputValue"
            :placeholder="modalState.inputPlaceholder"
            class="layui-input modal-input"
          >
        </div>
        <div class="modal-footer">
          <button v-if="modalState.cancelText" class="layui-btn layui-btn-primary" @click="handleCancel">
            {{ modalState.cancelText }}
          </button>
          <button 
            class="layui-btn" 
            :class="modalState.confirmColor === 'danger' ? 'layui-btn-danger' : ''"
            @click="() => handleConfirm(modalInputValue)"
          >
            {{ modalState.confirmText }}
          </button>
        </div>
      </div>
    </div>

    <button v-if="showBackToTop" class="back-to-top animate-bounce-in" @click="scrollToTop">
      <i class="layui-icon layui-icon-up"></i>
    </button>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useStore } from './stores'
import { commonApi } from './api'
import { marked } from 'marked'

const router = useRouter()
const route = useRoute()
const store = useStore()

const isMobile = ref(false)
const hideAppBar = ref(false)
const backendVersion = ref(null)
const siteTitle = ref('校园论坛')
const frontendVersion = ref(typeof __FRONTEND_VERSION__ !== 'undefined' ? __FRONTEND_VERSION__ : 'unknown')
const icpNumber = ref(null)
const publicSecurityNumber = ref(null)
const searchQuery = ref('')
const announcementDialog = ref(false)
const announcementContent = ref('')
const dontShowAgain = ref(false)
const mobileMenuOpen = ref(false)
const isScrolled = ref(false)
const showBackToTop = ref(false)
const modalInputValue = ref('')

const modalState = ref({
  show: false,
  type: 'info',
  title: '',
  message: '',
  icon: '',
  iconColor: '',
  confirmText: '确定',
  cancelText: '取消',
  confirmColor: '',
  inputValue: '',
  inputLabel: '',
  inputType: '',
  inputPlaceholder: '',
  inputRows: 0
})

let modalCallback = null

const handleConfirm = (value) => {
  if (modalCallback) {
    modalCallback(value)
  }
  modalState.value.show = false
  modalInputValue.value = ''
}

const handleCancel = () => {
  modalState.value.show = false
  modalInputValue.value = ''
}

const showModal = (options) => {
  return new Promise((resolve) => {
    modalCallback = resolve
    modalState.value = {
      show: true,
      type: options.type || 'info',
      title: options.title || '',
      message: options.message || '',
      icon: options.icon || '',
      iconColor: options.iconColor || '',
      confirmText: options.confirmText || '确定',
      cancelText: options.cancelText || '取消',
      confirmColor: options.confirmColor || '',
      inputValue: options.inputValue || '',
      inputLabel: options.inputLabel || '',
      inputType: options.inputType || '',
      inputPlaceholder: options.inputPlaceholder || '',
      inputRows: options.inputRows || 0
    }
  })
}

window.showModal = showModal

const loadAnnouncement = async () => {
  try {
    const response = await commonApi.getAnnouncement()
    if (response.data.content) {
      announcementContent.value = marked(response.data.content)
      const hideAnnouncement = localStorage.getItem('hideAnnouncement')
      if (!hideAnnouncement) {
        announcementDialog.value = true
      }
    }
  } catch (error) {
    console.error('加载公告失败', error)
  }
}

const closeAnnouncementDialog = () => {
  if (dontShowAgain.value) {
    localStorage.setItem('hideAnnouncement', 'true')
  }
  announcementDialog.value = false
}

const handleSearch = () => {
  if (searchQuery.value.trim()) {
    router.push({ path: '/search', query: { q: searchQuery.value.trim() } })
  }
}

const checkMobile = () => {
  isMobile.value = window.innerWidth < 768
}

const toggleMobileMenu = () => {
  mobileMenuOpen.value = !mobileMenuOpen.value
}

const handleScroll = () => {
  isScrolled.value = window.scrollY > 50
  showBackToTop.value = window.scrollY > 300
}

const scrollToTop = () => {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

const isAdminPage = computed(() => {
  return route.path.startsWith('/admin')
})

const logout = () => {
  store.logout()
  mobileMenuOpen.value = false
  router.push('/login')
}

const loadUser = async () => {
  if (store.isLoggedIn.value) {
    try {
      const response = await commonApi.getProfile()
      store.updateUser(response.data)
    } catch (error) {
      console.error('加载用户信息失败', error)
      store.logout()
    }
  }
}

const loadUnreadCount = async () => {
  if (store.isLoggedIn.value) {
    try {
      const response = await commonApi.getUnreadCount()
      store.setUnreadCount(response.data.unread_count || response.data.count || 0)
    } catch (error) {
      console.error('加载未读通知数失败', error)
    }
  }
}

const loadVersion = async () => {
  try {
    const response = await commonApi.getVersion()
    backendVersion.value = response.data.backend?.version || response.data.backend_version || response.data.version
    frontendVersion.value = response.data.frontend?.version || frontendVersion.value
  } catch (error) {
    console.error('加载版本信息失败', error)
  }
}

const loadSiteTitle = async () => {
  try {
    const response = await commonApi.getSiteConfig()
    siteTitle.value = response.data.site_title || response.data.name || '校园论坛'
    icpNumber.value = response.data.icp_number || response.data.icp || null
    publicSecurityNumber.value = response.data.public_security_number || null
    document.title = siteTitle.value
    store.setSiteConfig(response.data)
  } catch (error) {
    console.error('加载网站标题失败', error)
  }
}

watch(() => route.meta, (meta) => {
  hideAppBar.value = meta?.hideAppBar || false
}, { immediate: true })

watch(siteTitle, (newTitle) => {
  if (newTitle) {
    document.title = newTitle
  }
})

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  window.addEventListener('scroll', handleScroll)

  loadUser()
  loadVersion()
  loadSiteTitle()
  loadAnnouncement()
  loadUnreadCount()

  window.addEventListener('site-title-updated', (event) => {
    if (event.detail) {
      siteTitle.value = event.detail
      document.title = event.detail
    }
  })
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
  window.removeEventListener('scroll', handleScroll)
})
</script>

<style lang="less" scoped>
.app-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.announcement-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}

.announcement-modal {
  background: #fff;
  border-radius: 12px;
  width: 90%;
  max-width: 600px;
  overflow: hidden;
}

.announcement-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #e6e6e6;

  .announcement-title {
    font-size: 18px;
    font-weight: 600;
    color: #333;
    display: flex;
    align-items: center;
    gap: 8px;

    i {
      color: #1E9FFF;
    }
  }

  .close-btn {
    background: none;
    border: none;
    font-size: 20px;
    cursor: pointer;
    color: #999;
    padding: 4px 8px;

    &:hover {
      color: #666;
    }
  }
}

.announcement-content {
  padding: 20px;
  max-height: 400px;
  overflow-y: auto;
}

.announcement-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-top: 1px solid #e6e6e6;

  .checkbox-label {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 14px;
    color: #666;
  }
}

.navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1000;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(12px);
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  transition: all 0.4s @ease-out-cubic;

  &.scrolled {
    background: rgba(255, 255, 255, 0.98);
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  }

  .navbar-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 15px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 60px;
  }

  .navbar-brand {
    display: flex;
    align-items: center;
    gap: 10px;

    .logo-img {
      height: 36px;
    }

    .logo-text {
      font-size: 20px;
      font-weight: 700;
      color: #333;
      text-decoration: none;
    }
  }

  .nav-search {
    flex: 1;
    max-width: 300px;
    margin: 0 20px;

    .search-input {
      border-radius: 20px;
    }

    .search-btn {
      border-radius: 0 20px 20px 0;
    }
  }

  .nav-links {
    display: flex;
    align-items: center;
    gap: 4px;

    .nav-link {
      display: flex;
      align-items: center;
      gap: 6px;
      padding: 8px 16px;
      color: #333;
      text-decoration: none;
      font-size: 14px;
      border-radius: 6px;
      transition: all 0.3s @ease-out-back;

      &:hover {
        background: rgba(30, 159, 255, 0.1);
        color: #1E9FFF;
        transform: translateY(-2px);
      }

      &.nav-link-primary {
        background: #1E9FFF;
        color: #fff;

        &:hover {
          background: #0086E6;
        }
      }

      &.nav-link-admin {
        color: #FF5722;

        &:hover {
          background: rgba(255, 87, 34, 0.1);
        }
      }

      &.notification-link {
        position: relative;

        .badge {
          position: absolute;
          top: 4px;
          right: 4px;
          background: #FF5722;
          color: #fff;
          font-size: 10px;
          padding: 1px 5px;
          border-radius: 10px;
          min-width: 16px;
          text-align: center;
        }
      }

      &.logout-btn {
        background: none;
        border: none;
        cursor: pointer;

        &:hover {
          background: rgba(255, 87, 34, 0.1);
          color: #FF5722;
        }
      }
    }
  }

  .mobile-menu-btn {
    display: none;
    background: none;
    border: none;
    font-size: 24px;
    cursor: pointer;
    color: #333;
    padding: 8px;

    @media (max-width: 768px) {
      display: block;
    }
  }
}

.mobile-menu {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: #fff;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  z-index: 999;

  .mobile-menu-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 16px;
    border-bottom: 1px solid #e6e6e6;
    font-size: 16px;
    font-weight: 600;

    button {
      background: none;
      border: none;
      font-size: 20px;
      cursor: pointer;
      color: #999;
    }
  }

  .mobile-menu-content {
    padding: 8px 0;

    .mobile-menu-item {
      display: flex;
      align-items: center;
      gap: 12px;
      padding: 12px 16px;
      color: #333;
      text-decoration: none;
      font-size: 15px;
      transition: background 0.2s;

      &:hover {
        background: #f5f5f5;
      }

      &.admin-item {
        color: #FF5722;
      }

      &.logout-item {
        background: none;
        border: none;
        cursor: pointer;
        color: #FF5722;
        width: 100%;
      }
    }
  }
}

.main-content {
  flex: 1;
  padding-top: 60px;

  &.has-navbar {
    padding-top: 60px;
  }
}

.footer {
  background: #23262E;
  color: #fff;
  padding: 40px 0 20px;

  .footer-content {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 15px;
  }

  .footer-links {
    display: flex;
    justify-content: center;
    gap: 60px;
    margin-bottom: 30px;

    .footer-column {
      h4 {
        margin-bottom: 16px;
        font-size: 14px;
        color: rgba(255, 255, 255, 0.8);
      }

      ul {
        list-style: none;
        padding: 0;

        li {
          margin-bottom: 8px;
        }

        a {
          color: rgba(255, 255, 255, 0.6);
          text-decoration: none;
          font-size: 13px;

          &:hover {
            color: #1E9FFF;
          }
        }
      }
    }
  }

  .footer-bottom {
    border-top: 1px solid rgba(255, 255, 255, 0.1);
    padding-top: 20px;
    text-align: center;

    .version-info {
      font-size: 12px;
      color: rgba(255, 255, 255, 0.5);
      margin-bottom: 10px;
    }

    .cert-info {
      font-size: 12px;
      color: rgba(255, 255, 255, 0.5);

      a {
        color: rgba(255, 255, 255, 0.6);
        text-decoration: none;

        &:hover {
          color: #1E9FFF;
        }
      }

      .security-num {
        margin-left: 20px;
      }
    }
  }
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}

.modal-content {
  background: #fff;
  border-radius: 12px;
  width: 90%;
  max-width: 450px;
  overflow: hidden;

  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px 20px;
    border-bottom: 1px solid #e6e6e6;

    .modal-title {
      font-size: 16px;
      font-weight: 600;
      color: #333;
    }

    button {
      background: none;
      border: none;
      font-size: 20px;
      cursor: pointer;
      color: #999;
    }
  }

  .modal-body {
    padding: 20px;
    text-align: center;

    .modal-icon {
      font-size: 48px;
      margin-bottom: 16px;

      &.layui-icon-face-smile {
        color: #1E9FFF;
      }

      &.layui-icon-face-cry {
        color: #FF5722;
      }

      &.layui-icon-help-circle {
        color: #FFB800;
      }
    }

    p {
      color: #666;
      line-height: 1.6;
    }

    .modal-input, .modal-textarea {
      margin-top: 16px;
      width: 100%;
    }
  }

  .modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    padding: 16px 20px;
    border-top: 1px solid #e6e6e6;
  }
}

.back-to-top {
  position: fixed;
  bottom: 30px;
  right: 30px;
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: #1E9FFF;
  color: #fff;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  box-shadow: 0 4px 12px rgba(30, 159, 255, 0.4);
  z-index: 999;
  transition: all 0.3s @ease-out-back;

  &:hover {
    transform: translateY(-6px) scale(1.1);
    box-shadow: 0 8px 20px rgba(30, 159, 255, 0.5);
  }
}

@media (max-width: 768px) {
  .nav-search {
    display: none;
  }

  .nav-links {
    display: none;
  }

  .footer-links {
    flex-direction: column;
    align-items: center;
    gap: 20px;
  }
}
</style>
