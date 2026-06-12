<script setup>
import { provide, ref, onMounted } from 'vue'
import axios from 'axios'
import AppLayout from './components/AppLayout.vue'
import AppModal from './components/AppModal.vue'
import { alert } from './utils/modal'

const user = ref(null)
const siteConfig = ref({
  site_title: '校园论坛',
  site_description: ''
})

const isLoading = ref(true)

const setUser = (userData) => {
  user.value = userData
}

const clearUser = () => {
  user.value = null
  localStorage.removeItem('token')
  localStorage.removeItem('user')
}

provide('user', user)
provide('setUser', setUser)
provide('clearUser', clearUser)
provide('siteConfig', siteConfig)

axios.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

axios.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      clearUser()
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

const showAnnouncement = async () => {
  try {
    const response = await axios.get('/api/announcement')
    const content = response.data.content
    if (content && content.trim()) {
      await alert(content, {
        title: '站点公告',
        icon: 'mdi-bullhorn',
        iconColor: 'primary',
        markdown: true
      })
    }
  } catch (error) {
    console.log('加载公告失败:', error)
  }
}

onMounted(async () => {
  try {
    const token = localStorage.getItem('token')
    if (token) {
      const response = await axios.get('/api/profile')
      user.value = response.data
    }
    
    const configResponse = await axios.get('/api/site-config')
    siteConfig.value = configResponse.data
  } catch (error) {
    console.log('初始化失败:', error)
  } finally {
    isLoading.value = false
    // 显示公告弹窗
    showAnnouncement()
  }
})
</script>

<template>
  <v-app>
    <AppLayout v-if="!isLoading" />
    <v-progress-circular 
      v-else 
      indeterminate 
      color="primary" 
      class="position-fixed top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2"
    />
    <AppModal />
  </v-app>
</template>

<style>
body {
  margin: 0;
  padding: 0;
}
</style>
