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
  <v-container class="max-w-4xl mx-auto py-8">
    <v-card rounded="xl" elevation="4" class="mb-6">
      <v-card-text>
        <div class="flex">
          <v-text-field
            v-model="keyword"
            label="搜索文章"
            prepend-icon="mdi-magnify"
            class="flex-1"
            @keyup.enter="searchArticles"
          />
          <v-btn 
            color="primary" 
            class="ml-4"
            @click="searchArticles"
          >
            搜索
          </v-btn>
        </div>
      </v-card-text>
    </v-card>
    
    <!-- 搜索结果 -->
    <div v-if="isLoading" class="text-center py-8">
      <v-progress-circular indeterminate color="primary" />
    </div>
    
    <v-card 
      v-for="article in articles" 
      :key="article.id" 
      rounded="xl" 
      elevation="4" 
      class="mb-4 card-hover cursor-pointer"
      @click="router.push(`/article/${article.id}`)"
    >
      <v-card-title>
        <h3 class="text-lg font-bold text-gray-800">{{ article.title }}</h3>
      </v-card-title>
      <v-card-text>
        <p class="text-gray-600">{{ article.content.slice(0, 150) }}...</p>
      </v-card-text>
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
          <span>{{ article.like_count }}</span>
        </div>
      </v-card-subtitle>
    </v-card>
    
    <div v-if="!isLoading && articles.length === 0 && keyword" class="text-center py-12">
      <v-icon size="64" color="gray" class="mx-auto mb-4">mdi-search-off</v-icon>
      <p class="text-gray-500">未找到相关文章</p>
    </div>
  </v-container>
</template>