import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'
import './style.css'

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

const routes = [
    { path: '/', component: Home, name: 'Home' },
    { path: '/login', component: Login, name: 'Login' },
    { path: '/register', component: Register, name: 'Register' },
    { path: '/forgot-password', component: ForgotPassword, name: 'ForgotPassword' },
    { path: '/article/:id', component: Article, name: 'Article' },
    { path: '/create', component: CreateArticle, name: 'CreateArticle' },
    { path: '/profile', component: Profile, name: 'Profile' },
    { path: '/admin', component: Admin, name: 'Admin' },
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
router.beforeEach((to, from, next) => {
    const token = localStorage.getItem('token')
    const isInit = localStorage.getItem('isInit')
    
    if (to.path === '/login' || to.path === '/register') {
        if (token) {
            next('/')
        } else {
            next()
        }
    } else {
        if (!token && to.path !== '/') {
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
