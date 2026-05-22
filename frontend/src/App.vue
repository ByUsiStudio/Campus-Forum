<template>
  <v-app>
    <v-app-bar elevation="1" color="surface">
      <v-app-bar-title>
        <router-link to="/" class="logo-link">校园论坛</router-link>
      </v-app-bar-title>

      <template v-slot:append>
        <div class="nav-links">
          <v-btn variant="text" to="/" color="primary">首页</v-btn>
          <v-btn variant="text" to="/login" v-if="!token" color="primary">登录</v-btn>
          <v-btn variant="text" to="/register" v-if="!token" color="primary">注册</v-btn>
          <v-btn variant="text" to="/create" v-if="token" color="primary">写文章</v-btn>
          <v-btn variant="text" to="/profile" v-if="token" color="primary">个人中心</v-btn>
          <v-btn variant="text" to="/admin" v-if="isAdmin" color="error">管理后台</v-btn>
          <v-btn variant="text" v-if="token" @click="logout" color="secondary">退出</v-btn>
        </div>
      </template>
    </v-app-bar>

    <v-main>
      <v-container fluid>
        <router-view />
      </v-container>
    </v-main>
  </v-app>
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
      window.dispatchEvent(new Event('user-logout'))
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
      logout
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
</style>
