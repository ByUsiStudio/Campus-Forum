<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { articleApi } from '../api'

const router = useRouter()
const route = useRoute()

const keyword = ref('')
const articles = ref([])
const isLoading = ref(false)

const search = async () => {
  if (!keyword.value.trim()) return
  isLoading.value = true
  try {
    const response = await articleApi.searchArticles({ keyword: keyword.value })
    articles.value = response.data.articles || []
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
    return hours === 0 ? '刚刚' : `${hours}小时前`
  }
  return `${days}天前`
}

onMounted(() => {
  const queryKeyword = route.query.keyword
  if (queryKeyword) {
    keyword.value = queryKeyword
    search()
  }
})
</script>

<template>
  <v-app>
    <v-app-bar app>
      <v-btn icon @click="router.push('/')">
        <v-icon>mdi-arrow-left</v-icon>
      </v-btn>
      <v-text-field
        v-model="keyword"
        placeholder="搜索文章..."
        prepend-icon="mdi-magnify"
        class="flex-1 mx-4"
        @keyup.enter="search"
      />
      <v-btn color="primary" @click="search">搜索</v-btn>
    </v-app-bar>
    
    <v-container class="py-6">
      <v-card v-if="articles.length > 0">
        <v-list>
          <v-list-item
            v-for="article in articles"
            :key="article.id"
            @click="router.push(`/article/${article.id}`)"
          >
            <v-list-item-content>
              <v-list-item-title>{{ article.title }}</v-list-item-title>
              <v-list-item-subtitle>
                {{ article.content.substring(0, 50) }}...
                <span class="ml-2 text-grey">{{ formatTime(article.created_at) }}</span>
              </v-list-item-subtitle>
            </v-list-item-content>
          </v-list-item>
        </v-list>
      </v-card>
      
      <v-card v-else class="text-center py-12">
        <v-icon size="64" color="grey">mdi-search</v-icon>
        <p class="mt-4 text-grey">{{ keyword ? '未找到相关文章' : '请输入关键词进行搜索' }}</p>
      </v-card>
    </v-container>
  </v-app>
</template>
