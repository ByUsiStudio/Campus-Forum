<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { articleApi, categoryApi } from '../api'

const router = useRouter()
const route = useRoute()

const articles = ref([])
const category = ref(null)
const isLoading = ref(false)

const loadArticles = async () => {
  isLoading.value = true
  try {
    const response = await articleApi.getArticles({ category_id: route.params.id })
    articles.value = response.data.articles
  } catch (error) {
    console.error('加载文章失败:', error)
  } finally {
    isLoading.value = false
  }
}

const loadCategory = async () => {
  try {
    const response = await categoryApi.getCategories()
    category.value = response.data.categories.find(c => c.id === parseInt(route.params.id))
  } catch (error) {
    console.error('加载分类失败:', error)
  }
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
  loadCategory()
})
</script>

<template>
  <v-container class="max-w-4xl mx-auto px-4 py-8">
    <!-- 返回按钮 -->
    <v-btn 
      text 
      color="gray-600" 
      class="mb-6 hover:text-primary transition-colors"
      @click="router.push('/')"
    >
      <v-icon class="mr-2" size="20">mdi-arrow-left</v-icon>
      返回首页
    </v-btn>
    
    <!-- 分类信息卡片 -->
    <v-card rounded="2xl" elevation="4" class="mb-6 overflow-hidden">
      <div class="gradient-purple p-6">
        <div class="flex items-center">
          <v-icon class="mr-3" size="28" color="white">mdi-folder</v-icon>
          <div>
            <h2 class="text-white font-bold text-xl">{{ category?.name || '文章分类' }}</h2>
            <p class="text-white/70 text-sm mt-1">{{ category?.description || '暂无描述' }}</p>
          </div>
        </div>
      </div>
      <v-card-subtitle class="px-6 py-3 bg-gray-50">
        <span class="text-gray-500 text-sm">共 {{ articles.length }} 篇文章</span>
      </v-card-subtitle>
    </v-card>
    
    <!-- 文章列表 -->
    <div v-if="isLoading" class="loading-center">
      <v-progress-circular indeterminate color="primary" :size="48" />
    </div>
    
    <div v-else-if="articles.length > 0" class="space-y-4">
      <v-card 
        v-for="article in articles" 
        :key="article.id" 
        rounded="2xl" 
        elevation="2" 
        class="card-hover cursor-pointer overflow-hidden"
        @click="router.push(`/article/${article.id}`)"
      >
        <v-card-title class="px-6 py-4">
          <h3 class="text-lg font-bold text-gray-800">{{ article.title }}</h3>
        </v-card-title>
        <v-card-text class="px-6 pb-2">
          <p class="text-gray-600 text-sm card-text-limit">{{ article.content }}</p>
        </v-card-text>
        <v-card-subtitle class="px-6 py-3 bg-gray-50 flex items-center justify-between">
          <div class="flex items-center">
            <v-avatar size="32" color="primary" class="avatar-hover">
              <v-icon size="16" color="white">mdi-account</v-icon>
            </v-avatar>
            <span class="ml-2 text-gray-600 text-sm">{{ article.user?.display_name || article.user?.username }}</span>
          </div>
          <div class="flex items-center text-gray-400 text-sm">
            <span>{{ formatTime(article.created_at) }}</span>
            <span class="mx-2">·</span>
            <v-icon class="mr-1" size="16">mdi-eye</v-icon>
            <span>{{ article.view_count }}</span>
            <span class="mx-2">·</span>
            <v-icon class="mr-1" size="16">mdi-heart</v-icon>
            <span>{{ article.like_count }}</span>
          </div>
        </v-card-subtitle>
      </v-card>
    </div>
    
    <div v-else class="empty-state">
      <v-icon size="96" color="gray-200" class="empty-state-icon">mdi-file-question</v-icon>
      <p class="text-gray-400">该分类暂无文章</p>
      <p class="text-gray-400 text-sm mt-1">成为第一个发布文章的人吧</p>
    </div>
  </v-container>
</template>