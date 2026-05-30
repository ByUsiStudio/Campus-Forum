import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'
import './styles/style.css'

// Highlight.js 语法高亮
import hljs from 'highlight.js'
import 'highlight.js/styles/github-dark.css'
window.hljs = hljs

// Vuetify
import 'vuetify/styles'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import '@mdi/font/css/materialdesignicons.css'

const vuetify = createVuetify({
  components,
  directives,
  theme: {
    defaultTheme: 'light',
    themes: {
      light: {
        colors: {
          primary: '#6750A4',
          secondary: '#625B71',
          error: '#B3261E',
          background: '#FEF7FF',
          surface: '#FFFFFF',
          surfaceVariant: '#E7E0EC',
          onPrimary: '#FFFFFF',
          onSecondary: '#FFFFFF',
          onBackground: '#1C1B1F',
          onSurface: '#1C1B1F',
        }
      }
    }
  },
  defaults: {
    VBtn: {
      variant: 'flat',
    },
    VCard: {
      elevation: 1,
    }
  }
})

// 导入组件
import Home from './views/Home.vue'
import Login from './views/Login.vue'
import Register from './views/Register.vue'
import Article from './views/Article.vue'
import CreateArticle from './views/CreateArticle.vue'
import Profile from './views/Profile.vue'
import Admin from './views/Admin.vue'
import Category from './views/Category.vue'
import VideoPlayerPage from './views/VideoPlayerPage.vue'
import Notifications from './views/Notifications.vue'
import ChatList from './views/ChatList.vue'
import StandaloneChat from './views/StandaloneChat.vue'
import ForgotPassword from './views/ForgotPassword.vue'

// 导入后台管理子组件
import AdminIndex from './views/admin/AdminIndex.vue'
import AdminUsers from './views/admin/AdminUsers.vue'
import AdminArticles from './views/admin/AdminArticles.vue'
import AdminComments from './views/admin/AdminComments.vue'
import AdminCategories from './views/admin/AdminCategories.vue'
import AdminTitles from './views/admin/AdminTitles.vue'
import AdminSidebar from './views/admin/AdminSidebar.vue'
import AdminDeletions from './views/admin/AdminDeletions.vue'
import AdminAnnouncement from './views/admin/AdminAnnouncement.vue'
import AdminSiteConfig from './views/admin/AdminSiteConfig.vue'
import AdminSMTPConfig from './views/admin/AdminSMTPConfig.vue'
import AdminNotifications from './views/admin/AdminNotifications.vue'

const routes = [
    { path: '/', component: Home, name: 'Home' },
    { path: '/login', component: Login, name: 'Login' },
    { path: '/register', component: Register, name: 'Register' },
    { path: '/forgot-password', component: ForgotPassword, name: 'ForgotPassword' },
    { path: '/article/:id', component: Article, name: 'Article' },
    { path: '/create', component: CreateArticle, name: 'CreateArticle' },
    { path: '/profile', component: Profile, name: 'Profile' },
    { 
      path: '/admin', 
      component: Admin, 
      name: 'Admin',
      children: [
        { path: '', redirect: { name: 'AdminIndex' } },
        { path: 'index', name: 'AdminIndex', component: AdminIndex },
        { path: 'users', name: 'AdminUsers', component: AdminUsers },
        { path: 'articles', name: 'AdminArticles', component: AdminArticles },
        { path: 'comments', name: 'AdminComments', component: AdminComments },
        { path: 'categories', name: 'AdminCategories', component: AdminCategories },
        { path: 'titles', name: 'AdminTitles', component: AdminTitles },
        { path: 'sidebar', name: 'AdminSidebar', component: AdminSidebar },
        { path: 'deletions', name: 'AdminDeletions', component: AdminDeletions },
        { path: 'announcement', name: 'AdminAnnouncement', component: AdminAnnouncement },
        { path: 'siteconfig', name: 'AdminSiteConfig', component: AdminSiteConfig },
        { path: 'smtpconfig', name: 'AdminSMTPConfig', component: AdminSMTPConfig },
        { path: 'notifications', name: 'AdminNotifications', component: AdminNotifications }
      ]
    },
    { path: '/category/:id', component: Category, name: 'Category' },
    { path: '/video', component: VideoPlayerPage, name: 'VideoPlayer' },
    { path: '/notifications', component: Notifications, name: 'Notifications' },
    { path: '/chat', component: ChatList, name: 'ChatList' },
    { path: '/chat/:id', component: StandaloneChat, name: 'StandaloneChat', meta: { hideAppBar: true } }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

// 路由守卫
const publicPaths = ['/login', '/register', '/forgot-password', '/']

router.beforeEach((to, from, next) => {
    const token = localStorage.getItem('token')

    if (publicPaths.includes(to.path)) {
        if (token && (to.path === '/login' || to.path === '/register')) {
            next('/')
        } else {
            next()
        }
    } else {
        if (!token) {
            next('/login')
        } else {
            next()
        }
    }
})

const app = createApp(App)
app.use(vuetify)
app.use(router)
app.mount('#app')

// 注册 Service Worker（PWA 支持）
if ('serviceWorker' in navigator) {
  window.addEventListener('load', () => {
    navigator.serviceWorker.register('/sw.js')
      .then((registration) => {
        console.log('ServiceWorker registration successful with scope: ', registration.scope)
      })
      .catch((err) => {
        console.log('ServiceWorker registration failed: ', err)
      })
  })
}
