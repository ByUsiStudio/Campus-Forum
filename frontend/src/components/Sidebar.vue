<template>
  <div class="sidebar-component">
    <div class="user-info" v-if="user">
      <img :src="user.avatar" :alt="user.display_name" class="avatar">
      <div class="user-name">{{ user.display_name }}</div>
      <router-link to="/profile" class="profile-link">个人中心</router-link>
    </div>
    <div class="user-info" v-else>
      <div class="login-prompt">
        <router-link to="/login" class="btn btn-primary">登录</router-link>
        <router-link to="/register" class="btn btn-secondary">注册</router-link>
      </div>
    </div>
    
    <div class="sidebar-nav">
      <h3>导航</h3>
      <ul>
        <li v-for="item in sidebarItems" :key="item.link">
          <a :href="item.link">
            <span class="icon">{{ item.icon }}</span>
            {{ item.title }}
          </a>
        </li>
      </ul>
    </div>
    
    <div class="sidebar-stats">
      <h3>统计信息</h3>
      <div class="stats-item">总文章数：{{ stats.totalArticles }}</div>
      <div class="stats-item">总用户数：{{ stats.totalUsers }}</div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import api from '../api'

export default {
  name: 'Sidebar',
  setup() {
    const user = ref(null)
    const sidebarItems = ref([])
    const stats = ref({
      totalArticles: 0,
      totalUsers: 0
    })
    
    const loadUser = () => {
      const userStr = localStorage.getItem('user')
      if (userStr) {
        user.value = JSON.parse(userStr)
      }
    }
    
    const loadSidebarConfig = async () => {
      try {
        const response = await api.get('/sidebar-config')
        sidebarItems.value = response.data.items || []
      } catch (error) {
        console.error('加载侧边栏配置失败', error)
        // 默认配置
        sidebarItems.value = [
          { title: '首页', link: '/', icon: '🏠' }
        ]
      }
    }
    
    const loadStats = async () => {
      try {
        const response = await api.get('/articles', { params: { page: 1, page_size: 1 } })
        stats.value.totalArticles = response.data.total
        
        // 这里简单获取用户数，实际应该有个专门的API
        stats.value.totalUsers = 1
      } catch (error) {
        console.error('加载统计失败', error)
      }
    }
    
    onMounted(() => {
      loadUser()
      loadSidebarConfig()
      loadStats()
    })
    
    return {
      user,
      sidebarItems,
      stats
    }
  }
}
</script>

<style scoped>
.sidebar-component {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.user-info {
  text-align: center;
  padding: 20px;
  background: #f9fafb;
  border-radius: 8px;
}

.avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  margin-bottom: 10px;
  object-fit: cover;
}

.user-name {
  font-size: 16px;
  font-weight: 500;
  margin-bottom: 10px;
}

.login-prompt {
  display: flex;
  gap: 10px;
  justify-content: center;
}

.profile-link {
  display: inline-block;
  margin-top: 10px;
  color: #10b981;
  text-decoration: none;
  font-size: 14px;
}

.sidebar-nav h3,
.sidebar-stats h3 {
  margin-bottom: 15px;
  font-size: 16px;
  color: #1e293b;
}

.sidebar-nav ul {
  list-style: none;
  padding: 0;
}

.sidebar-nav li {
  margin-bottom: 10px;
}

.sidebar-nav a {
  text-decoration: none;
  color: #4b5563;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px;
  border-radius: 6px;
}

.icon {
  font-size: 18px;
}

.stats-item {
  padding: 8px 0;
  color: #6b7280;
  font-size: 14px;
  border-bottom: 1px solid #e5e7eb;
}
</style>