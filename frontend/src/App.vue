<script setup>
import { provide, ref, onMounted } from 'vue'
import axios from 'axios'

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
}

provide('user', user)
provide('setUser', setUser)
provide('clearUser', clearUser)
provide('siteConfig', siteConfig)

// 使用完整API路径，不在baseURL中添加前缀

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
  }
})
</script>

<template>
  <v-app>
    <v-main v-if="!isLoading">
      <router-view />
    </v-main>
    <v-progress-circular 
      v-else 
      indeterminate 
      color="primary" 
      class="position-fixed top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2"
    />
  </v-app>
</template>

<style scoped>
</style>