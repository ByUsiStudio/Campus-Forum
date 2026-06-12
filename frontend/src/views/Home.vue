<script setup>
import { ref, inject, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { articleApi, categoryApi, signinApi } from '../api'

const router = useRouter()
const user = inject('user')
const clearUser = inject('clearUser')

const articles = ref([])
const categories = ref([])
const currentCategory = ref(null)
const page = ref(1)
const totalPages = ref(1)
const isLoading = ref(false)
const signinStatus = ref({
  hasSignedIn: false,
  signInDays: 0,
  totalSignIns: 0
})

const loadArticles = async (categoryId = null, pageNum = 1) => {
  isLoading.value = true
  try {
    const params = {
      page: pageNum,
      page_size: 10
    }
    if (categoryId) {
      params.category_id = categoryId
    }
    const response = await articleApi.getArticles(params)
    if (pageNum === 1) {
      articles.value = response.data.articles
    } else {
      articles.value = [...articles.value, ...response.data.articles]
    }
    totalPages.value = response.data.total_pages
    page.value = pageNum
  } catch (error) {
    console.error('加载文章失败:', error)
  } finally {
    isLoading.value = false
  }
}

const loadCategories = async () => {
  try {
    const response = await categoryApi.getCategories()
    categories.value = response.data.categories
  } catch (error) {
    console.error('加载分类失败:', error)
  }
}

const loadSigninStatus = async () => {
  if (!user.value) return
  try {
    const response = await signinApi.getStatus()
    signinStatus.value = response.data
  } catch (error) {
    console.error('加载签到状态失败:', error)
  }
}

const handleSignin = async () => {
  if (!user.value) {
    router.push('/login')
    return
  }
  try {
    const response = await signinApi.signin()
    signinStatus.value = {
      hasSignedIn: true,
      signInDays: response.data.sign_in_days,
      totalSignIns: response.data.total_sign_ins
    }
  } catch (error) {
    console.error('签到失败:', error)
  }
}

const handleCategoryClick = (categoryId) => {
  currentCategory.value = categoryId
  loadArticles(categoryId, 1)
}

const handleLogout = () => {
  clearUser()
  router.push('/login')
}

const formatTime = (timeStr) => {
  const date = new Date(timeStr)
  const now = new Date()
  const diff = now - date
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  
  if (days === 0) {
    const hours = Math.floor(diff / (1000 * 60 * 60))
    if (hours === 0) {
      const minutes = Math.floor(diff / (1000 * 60))
      return minutes <= 0 ? '刚刚' : `${minutes}分钟前`
    }
    return `${hours}小时前`
  } else if (days < 7) {
    return `${days}天前`
  } else {
    return date.toLocaleDateString('zh-CN')
  }
}

onMounted(() => {
  loadArticles()
  loadCategories()
  loadSigninStatus()
})
</script>

<template>
  <v-container class="max-w-6xl mx-auto py-8">
    <!-- 头部导航 -->
    <v-app-bar 
      color="primary" 
      dark 
      rounded="lg" 
      class="mb-6"
    >
      <v-container class="flex items-center justify-between">
        <v-toolbar-title class="text-xl font-bold">
          <v-icon size="28" class="mr-2">mdi-forum</v-icon>
          校园论坛
        </v-toolbar-title>
        
        <v-spacer></v-spacer>
        
        <v-text-field
          v-model="searchQuery"
          label="搜索文章..."
          prepend-icon="mdi-magnify"
          rounded="lg"
          class="hidden md:flex w-64"
          @keyup.enter="router.push(`/search?keyword=${searchQuery}`)"
        />
        
        <v-spacer></v-spacer>
        
        <v-btn 
          v-if="user" 
          icon 
          class="mr-2"
          @click="router.push('/notifications')"
        >
          <v-icon size="24">mdi-bell</v-icon>
        </v-btn>
        
        <v-btn 
          v-if="user" 
          icon 
          class="mr-2"
          @click="router.push('/create')"
        >
          <v-icon size="24">mdi-plus</v-icon>
        </v-btn>
        
        <template v-if="user">
          <v-menu>
            <template v-slot:activator="{ props }">
              <v-btn 
                icon 
                v-bind="props"
              >
                <v-avatar size="40" color="secondary">
                  <v-icon>mdi-account</v-icon>
                </v-avatar>
              </v-btn>
            </template>
            <v-list>
              <v-list-item @click="router.push('/profile')">
                <v-list-item-icon>
                  <v-icon>mdi-user</v-icon>
                </v-list-item-icon>
                <v-list-item-title>个人中心</v-list-item-title>
              </v-list-item>
              <v-list-item v-if="user.role === 'admin'" @click="router.push('/admin')">
                <v-list-item-icon>
                  <v-icon>mdi-settings</v-icon>
                </v-list-item-icon>
                <v-list-item-title>管理后台</v-list-item-title>
              </v-list-item>
              <v-list-item @click="handleLogout">
                <v-list-item-icon>
                  <v-icon>mdi-logout</v-icon>
                </v-list-item-icon>
                <v-list-item-title>退出登录</v-list-item-title>
              </v-list-item>
            </v-list>
          </v-menu>
        </template>
        
        <template v-else>
          <v-btn text color="white" @click="router.push('/login')">登录</v-btn>
          <v-btn text color="white" @click="router.push('/register')">注册</v-btn>
        </template>
      </v-container>
    </v-app-bar>
    
    <v-row>
      <!-- 左侧分类栏 -->
      <v-col md="3" class="mb-6">
        <v-card rounded="xl" elevation="4" class="sticky top-4">
          <v-card-title class="gradient-purple text-white">
            <v-icon class="mr-2">mdi-folder-open</v-icon>
            <span class="font-bold">文章分类</span>
          </v-card-title>
          <v-card-text>
            <v-list>
              <v-list-item 
                :class="currentCategory === null ? 'active' : ''"
                @click="handleCategoryClick(null)"
              >
                <v-list-item-icon>
                  <v-icon :color="currentCategory === null ? 'primary' : 'gray'">mdi-home</v-icon>
                </v-list-item-icon>
                <v-list-item-title>全部文章</v-list-item-title>
              </v-list-item>
              <v-list-item 
                v-for="category in categories" 
                :key="category.id"
                :class="currentCategory === category.id ? 'active' : ''"
                @click="handleCategoryClick(category.id)"
              >
                <v-list-item-icon>
                  <v-icon :color="currentCategory === category.id ? 'primary' : 'gray'">mdi-folder</v-icon>
                </v-list-item-icon>
                <v-list-item-title>{{ category.name }}</v-list-item-title>
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>
        
        <!-- 签到卡片 -->
        <v-card 
          v-if="user"
          rounded="xl" 
          elevation="4" 
          class="mt-4 cursor-pointer card-hover"
          @click="handleSignin"
        >
          <v-card-title class="gradient-purple-light">
            <v-icon class="mr-2" size="24">mdi-calendar-check</v-icon>
            <span class="font-bold">每日签到</span>
          </v-card-title>
          <v-card-text class="text-center">
            <v-btn 
              v-if="!signinStatus.hasSignedIn"
              color="primary" 
              rounded="lg"
              class="mb-2"
            >
              立即签到
            </v-btn>
            <div v-else class="text-success">
              <v-icon class="mr-1">mdi-check-circle</v-icon>
              今日已签到
            </div>
            <p class="text-sm text-gray-500 mt-2">
              连续签到 {{ signinStatus.signInDays }} 天 · 累计 {{ signinStatus.totalSignIns }} 次
            </p>
          </v-card-text>
        </v-card>
      </v-col>
      
      <!-- 右侧文章列表 -->
      <v-col md="9">
        <v-card 
          v-for="article in articles" 
          :key="article.id" 
          rounded="xl" 
          elevation="4" 
          class="mb-4 card-hover cursor-pointer"
          @click="router.push(`/article/${article.id}`)"
        >
          <v-card-title>
            <div class="flex items-start justify-between w-full">
              <div class="flex-1">
                <h3 class="text-lg font-bold text-gray-800 mb-1">{{ article.title }}</h3>
                <p class="text-sm text-gray-500">{{ article.content.slice(0, 100) }}...</p>
              </div>
            </div>
          </v-card-title>
          <v-card-subtitle class="flex items-center justify-between text-sm">
            <div class="flex items-center">
              <v-avatar size="32" color="secondary">
                <v-icon>mdi-account</v-icon>
              </v-avatar>
              <span class="ml-2">{{ article.user?.display_name || article.user?.username }}</span>
              <span class="mx-2 text-gray-400">·</span>
              <span>{{ article.category?.name }}</span>
              <span class="mx-2 text-gray-400">·</span>
              <span>{{ formatTime(article.created_at) }}</span>
            </div>
            <div class="flex items-center">
              <v-icon class="mr-1" size="16">mdi-eye</v-icon>
              <span class="mr-3">{{ article.view_count }}</span>
              <v-icon class="mr-1" size="16">mdi-heart</v-icon>
              <span class="mr-3">{{ article.like_count }}</span>
              <v-icon class="mr-1" size="16">mdi-comment</v-icon>
              <span>{{ article.comment_count }}</span>
            </div>
          </v-card-subtitle>
        </v-card>
        
        <div v-if="isLoading" class="text-center py-8">
          <v-progress-circular indeterminate color="primary" />
        </div>
        
        <v-btn 
          v-if="page < totalPages && !isLoading"
          color="primary" 
          class="mx-auto d-block"
          @click="loadArticles(currentCategory, page + 1)"
        >
          加载更多
        </v-btn>
        
        <div v-if="articles.length === 0 && !isLoading" class="text-center py-12">
          <v-icon size="64" color="gray" class="mx-auto mb-4">mdi-file-question</v-icon>
          <p class="text-gray-500">暂无文章</p>
        </div>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>
.active {
  background-color: rgba(147, 112, 219, 0.1);
  color: #9370DB;
}

.active .v-icon {
  color: #9370DB !important;
}
</style>