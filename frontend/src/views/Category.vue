<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { articleApi, categoryApi } from '../api'

const router = useRouter()
const route = useRoute()

const articles = ref([])
const category = ref(null)
const isLoading = ref(false)

const loadCategory = async () => {
  isLoading.value = true
  try {
    const [articlesRes, categoriesRes] = await Promise.all([
      articleApi.getArticles({ category_id: route.params.id }),
      categoryApi.getCategories()
    ])
    articles.value = articlesRes.data.articles || []
    category.value = categoriesRes.data.categories.find(c => c.id === parseInt(route.params.id))
  } catch (error) {
    console.error('加载分类失败:', error)
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
  loadCategory()
})
</script>

<template>
  <v-app>
    <v-app-bar app>
      <v-btn icon @click="router.push('/')">
        <v-icon>mdi-arrow-left</v-icon>
      </v-btn>
      <v-toolbar-title>{{ category?.name || '分类' }}</v-toolbar-title>
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
                {{ article.author?.username }}
                <span class="ml-2 text-grey">{{ formatTime(article.created_at) }}</span>
              </v-list-item-subtitle>
            </v-list-item-content>
          </v-list-item>
        </v-list>
      </v-card>
      
      <v-card v-else class="text-center py-12">
        <v-icon size="64" color="grey">mdi-folder-open</v-icon>
        <p class="mt-4 text-grey">该分类暂无文章</p>
      </v-card>
    </v-container>
  </v-app>
</template>
