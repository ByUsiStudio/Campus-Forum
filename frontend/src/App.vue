<template>
  <div>
    <div class="header">
      <div class="header-content">
        <router-link to="/" class="logo">校园论坛</router-link>
        <div class="nav-links">
          <router-link to="/">首页</router-link>
          <router-link v-if="!token" to="/login">登录</router-link>
          <router-link v-if="!token" to="/register">注册</router-link>
          <router-link v-if="token" to="/create">写文章</router-link>
          <router-link v-if="token" to="/profile">个人中心</router-link>
          <router-link v-if="isAdmin" to="/admin">管理后台</router-link>
          <a v-if="token" href="#" @click.prevent="logout">退出</a>
        </div>
      </div>
    </div>
    <div class="container">
      <router-view />
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from './api'

export default {
  name: 'App',
  setup() {
    const router = useRouter()
    const token = ref(localStorage.getItem('token'))
    const user = ref(null)
    
    const isAdmin = computed(() => {
      return user.value && user.value.role === 'admin'
    })
    
    const logout = () => {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      token.value = null
      user.value = null
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
      const storedUser = localStorage.getItem('user')
      if (storedUser) {
        user.value = JSON.parse(storedUser)
      }
      loadUser()
    })
    
    return {
      token,
      user,
      isAdmin,
      logout
    }
  }
}
</script>