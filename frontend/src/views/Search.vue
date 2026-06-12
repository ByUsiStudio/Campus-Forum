<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { articleApi } from '../api'

const router = useRouter()
const route = useRoute()

const keyword = ref('')
const articles = ref([])
const isLoading = ref(false)

const searchArticles = async () => {
  if (!keyword.value.trim()) return
  
  isLoading.value = true
  try {
    const response = await articleApi.searchArticles({ keyword: keyword.value })
    articles.value = response.data.articles
  } catch (error) {
    console.error('搜索失败:', error)
  } finally {
    isLoading.value = false
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
  const queryKeyword = route.query.keyword
  if (queryKeyword) {
    keyword.value = queryKeyword
    searchArticles()
  }
})
</script>

<template>
  <v-container class="max-w-4xl mx-auto px-4 py-8">
    <!-- 返回按钮 -->
    <v-btn 
      text 
      color="gray-600" 
      class="mb-6 hover:text-primary transition-colors"
      @click="router.back()"
    >
      <v-icon class="mr-2" size="20">mdi-arrow-left</v-icon>
      返回
    </v-btn>
    
    <!-- 搜索框 -->
    <v-card rounded="2xl" elevation="4" class="mb-6 overflow-hidden">
      <div class="gradient-purple p-6">
        <h2 class="text-white font-bold text-lg mb-4">搜索文章</h2>
        <div class="flex gap-4">
          <v-text-field
            v-model="keyword"
            label="输入关键词"
            prepend-icon="mdi-magnify"
            background-color="white"
            rounded="xl"
            class="flex-1"
            @keyup.enter="searchArticles"
          />
          <v-btn 
            class="btn-gradient"
            @click="searchArticles"
          >
            <v-icon class="mr-2" size="18">mdi-magnify</v-icon>
            搜索
          </v-btn>
        </div>
      </div>
    </v-card>
    
    <!-- 搜索结果 -->
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
            <v-chip v-if="article.category" size="small" class="ml-3 tag-purple">
              {{ article.category.name }}
            </v-chip>
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
    
    <div v-else-if="!isLoading && keyword" class="empty-state">
      <v-icon size="96" color="gray-200" class="empty-state-icon">mdi-search-off</v-icon>
      <p class="text-gray-400">未找到相关文章</p>
      <p class="text-gray-400 text-sm mt-1">试试其他关键词吧</p>
    </div>
    
    <div v-else class="empty-state">
      <v-icon size="96" color="gray-200" class="empty-state-icon">mdi-magnify</v-icon>
      <p class="text-gray-400">输入关键词开始搜索</p>
      <p class="text-gray-400 text-sm mt-1">查找你感兴趣的文章</p>
    </div>
  </v-container>
</template>