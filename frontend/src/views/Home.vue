<script setup>
import { ref, inject, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { articleApi, categoryApi, signinApi } from '../api'

const router = useRouter()
const route = useRoute()
const user = inject('user')
const clearUser = inject('clearUser')
const isMobile = computed(() => window.innerWidth < 1024)

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

const loadMore = () => {
  if (page.value < totalPages.value && !isLoading.value) {
    loadArticles(currentCategory.value, page.value + 1)
  }
}

onMounted(() => {
  loadArticles()
  loadCategories()
  loadSigninStatus()
})
</script>

<template>
  <v-container class="py-6">
    <!-- 移动端布局：分类横向滚动 + 文章列表 -->
    <template v-if="isMobile">
      <!-- 分类横向滚动 -->
      <v-scroll-x class="mb-4">
        <v-btn
          v-for="category in [{ id: null, name: '全部' }, ...categories]"
          :key="category.id || 'all'"
          :class="currentCategory === category.id ? 'primary' : ''"
          @click="handleCategoryClick(category.id)"
        >
          {{ category.name }}
        </v-btn>
      </v-scroll-x>

      <!-- 签到卡片 -->
      <v-card v-if="user" class="mb-4" @click="handleSignin">
        <v-card-title class="d-flex align-center justify-between">
          <span>每日签到</span>
          <v-icon color="primary">mdi-calendar-check</v-icon>
        </v-card-title>
        <v-card-text>
          <div class="text-center">
            <div class="text-4xl font-bold text-primary">{{ signinStatus.signInDays }}</div>
            <div class="text-sm text-grey">连续签到天数</div>
            <div class="text-sm mt-2">
              累计签到 {{ signinStatus.totalSignIns }} 次
            </div>
          </div>
          <v-btn v-if="!signinStatus.hasSignedIn" block color="primary" class="mt-4">
            立即签到
          </v-btn>
          <v-btn v-else block disabled class="mt-4">
            今日已签到
          </v-btn>
        </v-card-text>
      </v-card>

      <!-- 文章列表 -->
      <v-card
        v-for="article in articles"
        :key="article.id"
        class="mb-4 cursor-pointer"
        @click="router.push(`/article/${article.id}`)"
      >
        <v-card-title>
          <h3 class="text-h6">{{ article.title }}</h3>
        </v-card-title>
        <v-card-text>
          <p>{{ article.content.substring(0, 100) }}...</p>
        </v-card-text>
        <v-card-actions>
          <v-chip size="small" color="primary" text-color="white">
            {{ article.category?.name || '未分类' }}
          </v-chip>
          <span class="ml-2 text-sm text-grey">{{ article.author?.username }}</span>
          <span class="ml-auto text-sm text-grey">{{ formatTime(article.created_at) }}</span>
        </v-card-actions>
      </v-card>

      <!-- 加载更多 -->
      <v-card v-if="page < totalPages" class="text-center">
        <v-card-text>
          <v-btn color="primary" :loading="isLoading" @click="loadMore">
            加载更多
          </v-btn>
        </v-card-text>
      </v-card>

      <v-card v-if="articles.length === 0" class="text-center py-12">
        <v-icon size="64" color="grey">mdi-file-question</v-icon>
        <p class="mt-4 text-grey">暂无文章</p>
      </v-card>
    </template>

    <!-- PC端布局：左侧分类栏 + 右侧文章列表 -->
    <template v-else>
      <v-row>
        <!-- 左侧分类栏 -->
        <v-col md="3" class="mb-6">
          <v-card>
            <v-card-title>文章分类</v-card-title>
            <v-list>
              <v-list-item
                :class="currentCategory === null ? 'active' : ''"
                @click="handleCategoryClick(null)"
              >
                <v-list-item-icon>
                  <v-icon :color="currentCategory === null ? 'primary' : 'grey'">mdi-home</v-icon>
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
                  <v-icon :color="currentCategory === category.id ? 'primary' : 'grey'">mdi-folder</v-icon>
                </v-list-item-icon>
                <v-list-item-title>{{ category.name }}</v-list-item-title>
              </v-list-item>
            </v-list>
          </v-card>
          
          <!-- 签到卡片 -->
          <v-card v-if="user" class="mt-4 cursor-pointer" @click="handleSignin">
            <v-card-title class="d-flex align-center justify-between">
              <span>每日签到</span>
              <v-icon color="primary">mdi-calendar-check</v-icon>
            </v-card-title>
            <v-card-text>
              <div class="text-center">
                <div class="text-4xl font-bold text-primary">{{ signinStatus.signInDays }}</div>
                <div class="text-sm text-grey">连续签到天数</div>
                <div class="text-sm mt-2">
                  累计签到 {{ signinStatus.totalSignIns }} 次
                </div>
              </div>
              <v-btn v-if="!signinStatus.hasSignedIn" block color="primary" class="mt-4">
                立即签到
              </v-btn>
              <v-btn v-else block disabled class="mt-4">
                今日已签到
              </v-btn>
            </v-card-text>
          </v-card>
        </v-col>
        
        <!-- 右侧文章列表 -->
        <v-col md="9">
          <v-card
            v-for="article in articles"
            :key="article.id"
            class="mb-4 cursor-pointer"
            @click="router.push(`/article/${article.id}`)"
          >
            <v-card-title>
              <h3 class="text-h6">{{ article.title }}</h3>
            </v-card-title>
            <v-card-text>
              <p>{{ article.content.substring(0, 100) }}...</p>
            </v-card-text>
            <v-card-actions>
              <v-chip size="small" color="primary" text-color="white">
                {{ article.category?.name || '未分类' }}
              </v-chip>
              <span class="ml-2 text-sm text-grey">{{ article.author?.username }}</span>
              <span class="ml-auto text-sm text-grey">{{ formatTime(article.created_at) }}</span>
              <v-icon class="ml-2 text-grey">mdi-eye</v-icon>
              <span class="text-sm text-grey">{{ article.view_count }}</span>
              <v-icon class="ml-2 text-grey">mdi-heart</v-icon>
              <span class="text-sm text-grey">{{ article.like_count }}</span>
            </v-card-actions>
          </v-card>
          
          <!-- 加载更多 -->
          <v-card v-if="page < totalPages" class="text-center">
            <v-card-text>
              <v-btn color="primary" :loading="isLoading" @click="loadMore">
                加载更多
              </v-btn>
            </v-card-text>
          </v-card>
          
          <v-card v-if="articles.length === 0" class="text-center py-12">
            <v-icon size="64" color="grey">mdi-file-question</v-icon>
            <p class="mt-4 text-grey">暂无文章</p>
          </v-card>
        </v-col>
      </v-row>
    </template>
  </v-container>
</template>

<style scoped>
.active {
  background-color: rgba(98, 0, 238, 0.1);
}
</style>
