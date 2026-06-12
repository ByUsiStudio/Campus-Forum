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
  <v-container class="max-w-4xl mx-auto py-8">
    <v-btn 
      text 
      color="primary" 
      class="mb-6"
      @click="router.push('/')"
    >
      <v-icon class="mr-1">mdi-arrow-left</v-icon>
      返回首页
    </v-btn>
    
    <v-card rounded="xl" elevation="4" class="mb-6">
      <v-card-title class="gradient-purple text-white">
        <v-icon class="mr-2">mdi-folder</v-icon>
        <span class="font-bold text-xl">{{ category?.name || '文章分类' }}</span>
      </v-card-title>
      <v-card-text class="text-gray-600">
        {{ category?.description || '暂无描述' }}
      </v-card-text>
    </v-card>
    
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
    
    <div v-if="!isLoading && articles.length === 0" class="text-center py-12">
      <v-icon size="64" color="gray" class="mx-auto mb-4">mdi-file-question</v-icon>
      <p class="text-gray-500">该分类暂无文章</p>
    </div>
  </v-container>
</template>