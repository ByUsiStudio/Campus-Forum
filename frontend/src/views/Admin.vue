<script setup>
import { ref, inject, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { adminApi } from '../api'

const router = useRouter()
const user = inject('user')

const activeTab = ref('users')
const users = ref([])
const articles = ref([])
const isLoading = ref(false)

const loadUsers = async () => {
  isLoading.value = true
  try {
    const response = await adminApi.getUsers()
    users.value = response.data.users || []
  } catch (error) {
    console.error('加载用户失败:', error)
  } finally {
    isLoading.value = false
  }
}

const loadArticles = async () => {
  isLoading.value = true
  try {
    const response = await adminApi.getArticles()
    articles.value = response.data.articles || []
  } catch (error) {
    console.error('加载文章失败:', error)
  } finally {
    isLoading.value = false
  }
}

const handleBanUser = async (userId) => {
  try {
    await adminApi.banUser(userId, { reason: '违规' })
    loadUsers()
  } catch (error) {
    console.error('封禁失败:', error)
  }
}

const handleDeleteArticle = async (articleId) => {
  try {
    await adminApi.deleteArticle(articleId)
    loadArticles()
  } catch (error) {
    console.error('删除文章失败:', error)
  }
}

onMounted(() => {
  if (!user.value || user.value.role !== 'admin') {
    router.push('/')
    return
  }
  loadUsers()
})
</script>

<template>
  <v-app>
    <v-app-bar app>
      <v-btn icon @click="router.push('/')">
        <v-icon>mdi-arrow-left</v-icon>
      </v-btn>
      <v-toolbar-title>管理后台</v-toolbar-title>
    </v-app-bar>
    
    <v-container class="py-6">
      <v-tabs v-model="activeTab" background-color="transparent" centered>
        <v-tab value="users">用户管理</v-tab>
        <v-tab value="articles">文章管理</v-tab>
      </v-tabs>
      
      <!-- 用户管理 -->
      <v-card v-if="activeTab === 'users'" class="mt-4">
        <v-list v-if="users.length > 0">
          <v-list-item v-for="u in users" :key="u.id">
            <v-list-item-avatar>
              <v-icon color="primary">mdi-account</v-icon>
            </v-list-item-avatar>
            <v-list-item-content>
              <v-list-item-title>{{ u.username }}</v-list-item-title>
              <v-list-item-subtitle>{{ u.email }} - {{ u.role }}</v-list-item-subtitle>
            </v-list-item-content>
            <v-list-item-actions>
              <v-btn text color="error" @click="handleBanUser(u.id)">封禁</v-btn>
            </v-list-item-actions>
          </v-list-item>
        </v-list>
        <v-card-text v-else class="text-center text-grey py-8">
          暂无用户
        </v-card-text>
      </v-card>
      
      <!-- 文章管理 -->
      <v-card v-if="activeTab === 'articles'" class="mt-4">
        <v-list v-if="articles.length > 0">
          <v-list-item v-for="article in articles" :key="article.id">
            <v-list-item-content>
              <v-list-item-title>{{ article.title }}</v-list-item-title>
              <v-list-item-subtitle>{{ article.author?.username }}</v-list-item-subtitle>
            </v-list-item-content>
            <v-list-item-actions>
              <v-btn text @click="router.push(`/article/${article.id}`)">查看</v-btn>
              <v-btn text color="error" @click="handleDeleteArticle(article.id)">删除</v-btn>
            </v-list-item-actions>
          </v-list-item>
        </v-list>
        <v-card-text v-else class="text-center text-grey py-8">
          暂无文章
        </v-card-text>
      </v-card>
    </v-container>
  </v-app>
</template>
