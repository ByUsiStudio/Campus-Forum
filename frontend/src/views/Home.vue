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
const searchQuery = ref('')
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
  <v-container class="max-w-7xl mx-auto px-4 py-6">
    <!-- 头部导航 -->
    <v-app-bar 
      class="mb-8 rounded-2xl shadow-purple-md"
      background-color="surface"
      elevation="4"
    >
      <v-container class="flex items-center justify-between px-6">
        <v-toolbar-title class="text-xl font-bold">
          <v-icon size="32" class="mr-2 text-primary">mdi-forum</v-icon>
          <span class="text-gradient">校园论坛</span>
        </v-toolbar-title>
        
        <v-spacer></v-spacer>
        
        <v-text-field
          v-model="searchQuery"
          label="搜索文章..."
          prepend-icon="mdi-magnify"
          rounded="xl"
          class="hidden md:flex w-80"
          color="primary"
          hide-details="auto"
          @keyup.enter="router.push(`/search?keyword=${searchQuery}`)"
        />
        
        <v-spacer></v-spacer>
        
        <div class="flex items-center gap-2">
          <v-btn 
            v-if="user" 
            icon 
            class="hidden sm:flex"
            @click="router.push('/notifications')"
          >
            <v-icon size="24" color="gray-600">mdi-bell</v-icon>
          </v-btn>
          
          <v-btn 
            v-if="user" 
            icon 
            class="btn-gradient"
            @click="router.push('/create')"
          >
            <v-icon size="20" color="white">mdi-plus</v-icon>
          </v-btn>
          
          <template v-if="user">
            <v-menu>
              <template v-slot:activator="{ props }">
                <v-btn 
                  icon 
                  v-bind="props"
                  class="ml-2"
                >
                  <v-avatar size="44" color="primary" class="avatar-hover">
                    <v-icon size="20" color="white">mdi-account</v-icon>
                  </v-avatar>
                </v-btn>
              </template>
              <v-list rounded="xl" elevation="8">
                <v-list-item @click="router.push('/profile')" class="list-item-hover">
                  <v-list-item-icon>
                    <v-icon color="primary">mdi-user</v-icon>
                  </v-list-item-icon>
                  <v-list-item-title>个人中心</v-list-item-title>
                </v-list-item>
                <v-list-item v-if="user.role === 'admin'" @click="router.push('/admin')" class="list-item-hover">
                  <v-list-item-icon>
                    <v-icon color="primary">mdi-settings</v-icon>
                  </v-list-item-icon>
                  <v-list-item-title>管理后台</v-list-item-title>
                </v-list-item>
                <v-divider></v-divider>
                <v-list-item @click="handleLogout" class="list-item-hover text-error">
                  <v-list-item-icon>
                    <v-icon color="error">mdi-logout</v-icon>
                  </v-list-item-icon>
                  <v-list-item-title>退出登录</v-list-item-title>
                </v-list-item>
              </v-list>
            </v-menu>
          </template>
          
          <template v-else>
            <v-btn text color="gray-600" class="hidden sm:flex" @click="router.push('/login')">登录</v-btn>
            <v-btn class="btn-gradient" @click="router.push('/register')">注册</v-btn>
          </template>
        </div>
      </v-container>
    </v-app-bar>
    
    <v-row>
      <!-- 左侧分类栏 -->
      <v-col md="3" class="mb-6">
        <v-card rounded="2xl" elevation="4" class="sticky top-6 overflow-hidden">
          <v-card-title class="gradient-purple text-white py-4 px-6">
            <v-icon class="mr-3" size="20">mdi-folder-open</v-icon>
            <span class="font-bold">文章分类</span>
          </v-card-title>
          <v-card-text class="px-0">
            <v-list rounded="none">
              <v-list-item 
                :class="currentCategory === null ? 'selected-item' : ''"
                @click="handleCategoryClick(null)"
                class="list-item-hover px-6"
              >
                <v-list-item-icon>
                  <v-icon :color="currentCategory === null ? 'primary' : 'gray-400'">mdi-home</v-icon>
                </v-list-item-icon>
                <v-list-item-title :class="currentCategory === null ? 'text-primary font-medium' : ''">全部文章</v-list-item-title>
              </v-list-item>
              <v-list-item 
                v-for="category in categories" 
                :key="category.id"
                :class="currentCategory === category.id ? 'selected-item' : ''"
                @click="handleCategoryClick(category.id)"
                class="list-item-hover px-6"
              >
                <v-list-item-icon>
                  <v-icon :color="currentCategory === category.id ? 'primary' : 'gray-400'">mdi-folder</v-icon>
                </v-list-item-icon>
                <v-list-item-title :class="currentCategory === category.id ? 'text-primary font-medium' : ''">{{ category.name }}</v-list-item-title>
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>
        
        <!-- 签到卡片 -->
        <v-card 
          v-if="user"
          rounded="2xl" 
          elevation="4" 
          class="mt-4 cursor-pointer card-hover overflow-hidden"
          @click="handleSignin"
        >
          <v-card-title class="gradient-purple-light py-4 px-6">
            <v-icon class="mr-3" size="24" color="primary">mdi-calendar-check</v-icon>
            <span class="font-bold text-gray-800">每日签到</span>
          </v-card-title>
          <v-card-text class="text-center py-4">
            <v-btn 
              v-if="!signinStatus.hasSignedIn"
              class="btn-gradient mb-3"
              size="large"
            >
              <v-icon class="mr-2" size="18">mdi-calendar-plus</v-icon>
              立即签到
            </v-btn>
            <div v-else class="flex items-center justify-center text-success mb-3">
              <v-icon class="mr-2" size="24">mdi-check-circle</v-icon>
              <span class="font-medium">今日已签到</span>
            </div>
            <div class="flex items-center justify-center gap-4 text-sm text-gray-500">
              <span>
                <v-icon class="mr-1" size="14">mdi-flame</v-icon>
                连续 {{ signinStatus.signInDays }} 天
              </span>
              <span class="text-gray-300">|</span>
              <span>
                <v-icon class="mr-1" size="14">mdi-star</v-icon>
                累计 {{ signinStatus.totalSignIns }} 次
              </span>
            </div>
          </v-card-text>
        </v-card>
        
        <!-- 移动端搜索 -->
        <v-text-field
          v-model="searchQuery"
          label="搜索文章..."
          prepend-icon="mdi-magnify"
          rounded="xl"
          class="md:hidden mt-4"
          color="primary"
          hide-details="auto"
          @keyup.enter="router.push(`/search?keyword=${searchQuery}`)"
        />
      </v-col>
      
      <!-- 右侧文章列表 -->
      <v-col md="9">
        <v-card 
          v-for="article in articles" 
          :key="article.id" 
          rounded="2xl" 
          elevation="2" 
          class="mb-4 card-hover cursor-pointer overflow-hidden"
          @click="router.push(`/article/${article.id}`)"
        >
          <v-card-title class="px-6 py-5">
            <div class="flex items-start justify-between w-full">
              <div class="flex-1">
                <div class="flex items-center gap-2 mb-2">
                  <v-chip 
                    v-if="article.category" 
                    size="small" 
                    class="tag-purple"
                  >
                    {{ article.category.name }}
                  </v-chip>
                </div>
                <h3 class="text-xl font-bold text-gray-800 mb-2">{{ article.title }}</h3>
                <p class="text-gray-600 card-text-limit">{{ article.content }}</p>
              </div>
            </div>
          </v-card-title>
          <v-card-subtitle class="px-6 py-3 flex items-center justify-between bg-gray-50">
            <div class="flex items-center gap-3">
              <v-avatar size="36" color="primary" class="avatar-hover">
                <v-icon size="16" color="white">mdi-account</v-icon>
              </v-avatar>
              <div>
                <span class="font-medium text-gray-800">{{ article.user?.display_name || article.user?.username }}</span>
                <p class="text-xs text-gray-400">{{ formatTime(article.created_at) }}</p>
              </div>
            </div>
            <div class="flex items-center gap-4 text-sm text-gray-500">
              <span class="flex items-center gap-1">
                <v-icon size="16">mdi-eye</v-icon>
                {{ article.view_count }}
              </span>
              <span class="flex items-center gap-1">
                <v-icon size="16">mdi-heart</v-icon>
                {{ article.like_count }}
              </span>
              <span class="flex items-center gap-1">
                <v-icon size="16">mdi-comment</v-icon>
                {{ article.comment_count }}
              </span>
            </div>
          </v-card-subtitle>
        </v-card>
        
        <div v-if="isLoading" class="loading-center">
          <v-progress-circular indeterminate color="primary" :size="48" />
        </div>
        
        <v-btn 
          v-if="page < totalPages && !isLoading"
          class="btn-gradient mx-auto d-block mt-4"
          size="large"
          @click="loadArticles(currentCategory, page + 1)"
        >
          <v-icon class="mr-2" size="18">mdi-chevron-down</v-icon>
          加载更多
        </v-btn>
        
        <div v-if="articles.length === 0 && !isLoading" class="empty-state">
          <v-icon size="96" color="gray-300" class="empty-state-icon">mdi-file-question</v-icon>
          <p class="text-gray-400 text-lg">暂无文章</p>
          <p class="text-gray-400 text-sm mt-2">快来发布第一篇文章吧！</p>
        </div>
      </v-col>
    </v-row>
  </v-container>
</template>