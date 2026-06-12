<script setup>
import { ref, inject, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { articleApi, categoryApi, signinApi } from '../api'

const router = useRouter()
const user = inject('user')
const clearUser = inject('clearUser')
const isMobile = computed(() => window.innerWidth < 1024)

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

const loadMore = () => {
  if (page.value < totalPages.value && !isLoading.value) {
    loadArticles(currentCategory.value, page.value + 1)
  }
}

const goToArticle = (articleId) => {
  router.push(`/article/${articleId}`)
}

onMounted(() => {
  loadArticles()
  loadCategories()
  loadSigninStatus()
})
</script>

<template>
  <v-container class="py-6" fluid>
    <v-row>
      <!-- PC端左侧边栏 -->
      <v-col v-if="!isMobile" cols="12" md="3">
        <div class="sidebar-sticky">
          <!-- 分类卡片 -->
          <v-card class="mb-4">
            <v-card-title class="text-subtitle-1 font-weight-bold pb-2">
              <v-icon size="20" class="mr-2">mdi-folder-outline</v-icon>
              文章分类
            </v-card-title>
            <v-divider></v-divider>
            <v-list density="compact" class="py-2">
              <v-list-item
                :active="currentCategory === null"
                :class="{ 'category-active': currentCategory === null }"
                @click="handleCategoryClick(null)"
                class="category-item"
              >
                <template #prepend>
                  <v-icon size="18" class="mr-2">mdi-view-grid-outline</v-icon>
                </template>
                <v-list-item-title class="text-body-2">全部文章</v-list-item-title>
              </v-list-item>
              <v-list-item
                v-for="category in categories"
                :key="category.id"
                :active="currentCategory === category.id"
                :class="{ 'category-active': currentCategory === category.id }"
                @click="handleCategoryClick(category.id)"
                class="category-item"
              >
                <template #prepend>
                  <v-icon size="18" class="mr-2">mdi-folder-outline</v-icon>
                </template>
                <v-list-item-title class="text-body-2">{{ category.name }}</v-list-item-title>
              </v-list-item>
            </v-list>
          </v-card>

          <!-- 签到卡片 -->
          <v-card v-if="user" class="signin-card" @click="handleSignin">
            <v-card-text class="text-center pa-4">
              <div class="d-flex align-center justify-center mb-2">
                <v-icon color="primary" size="24" class="mr-2">mdi-calendar-check</v-icon>
                <span class="text-subtitle-2 font-weight-bold">每日签到</span>
              </div>
              <div class="signin-days">
                <span class="text-h3 font-weight-bold text-primary">{{ signinStatus.signInDays }}</span>
              </div>
              <div class="text-caption text-grey mb-3">连续签到天数</div>
              <div class="text-caption text-grey mb-3">
                累计签到 {{ signinStatus.totalSignIns }} 次
              </div>
              <v-btn
                :color="signinStatus.hasSignedIn ? 'grey' : 'primary'"
                :variant="signinStatus.hasSignedIn ? 'outlined' : 'flat'"
                size="small"
                block
                :disabled="signinStatus.hasSignedIn"
              >
                {{ signinStatus.hasSignedIn ? '今日已签到' : '立即签到' }}
              </v-btn>
            </v-card-text>
          </v-card>
        </div>
      </v-col>

      <!-- 右侧文章列表 -->
      <v-col cols="12" :md="isMobile ? 12 : 9">
        <!-- 移动端分类横向滚动 -->
        <v-scroll-x v-if="isMobile" class="mb-4 category-scroll">
          <v-chip-group
            v-model="currentCategory"
            selected-class="bg-primary text-white"
            class="d-inline-flex"
          >
            <v-chip
              :value="null"
              variant="outlined"
              class="mx-1"
              @click="handleCategoryClick(null)"
            >
              全部
            </v-chip>
            <v-chip
              v-for="category in categories"
              :key="category.id"
              :value="category.id"
              variant="outlined"
              class="mx-1"
              @click="handleCategoryClick(category.id)"
            >
              {{ category.name }}
            </v-chip>
          </v-chip-group>
        </v-scroll-x>

        <!-- 文章列表 -->
        <div v-if="articles.length > 0" class="article-list">
          <v-card
            v-for="article in articles"
            :key="article.id"
            class="article-card mb-3"
            @click="goToArticle(article.id)"
          >
            <v-card-text class="pa-4">
              <!-- 标题 -->
              <h3 class="text-title-1 font-weight-bold text-grey-darken-3 mb-2 article-title">
                {{ article.title }}
              </h3>

              <!-- 摘要 -->
              <p class="text-body-2 text-grey mb-3 article-summary">
                {{ article.content.replace(/[#*`]/g, '').substring(0, 120) }}...
              </p>

              <!-- 元信息 -->
              <div class="d-flex align-center flex-wrap text-caption text-grey">
                <v-avatar size="20" class="mr-2">
                <v-img v-if="article.user?.avatar" :src="article.user.avatar"></v-img>
                <v-icon v-else size="16">mdi-account</v-icon>
              </v-avatar>
              <span class="mr-3">{{ article.user?.display_name || article.user?.username || '未知用户' }}</span>

                <v-chip
                  size="x-small"
                  variant="tonal"
                  color="primary"
                  class="mr-3"
                >
                  {{ article.category?.name || '未分类' }}
                </v-chip>

                <v-icon size="14" class="mr-1">mdi-clock-outline</v-icon>
                <span class="mr-3">{{ formatTime(article.created_at) }}</span>

                <v-spacer></v-spacer>

                <div class="d-flex align-center">
                  <v-icon size="14" class="mr-1">mdi-eye-outline</v-icon>
                  <span class="mr-3">{{ article.view_count || 0 }}</span>
                  <v-icon size="14" class="mr-1">mdi-heart-outline</v-icon>
                  <span>{{ article.like_count || 0 }}</span>
                </div>
              </div>
            </v-card-text>
          </v-card>
        </div>

        <!-- 空状态 -->
        <v-card v-else class="text-center py-16 empty-state">
          <v-icon size="80" color="primary-lighten-2" class="mb-4">mdi-file-document-outline</v-icon>
          <p class="text-h6 text-grey mb-2">暂无文章</p>
          <p class="text-body-2 text-grey">快去发布第一篇文章吧</p>
          <v-btn
            color="primary"
            class="mt-4"
            to="/create"
          >
            发布文章
          </v-btn>
        </v-card>

        <!-- 加载更多 -->
        <div v-if="page < totalPages" class="text-center mt-6">
          <v-btn
            variant="outlined"
            color="primary"
            :loading="isLoading"
            @click="loadMore"
            class="load-more-btn"
          >
            加载更多
          </v-btn>
        </div>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>
.sidebar-sticky {
  position: sticky;
  top: 80px;
}

.category-item {
  border-radius: 8px;
  margin: 2px 8px;
}

.category-active {
  background-color: rgba(149, 117, 205, 0.12) !important;
}

.category-active .v-list-item__overlay {
  opacity: 0;
}

.signin-card {
  border: 1px solid rgba(149, 117, 205, 0.2);
}

.signin-days {
  line-height: 1.2;
}

.article-card {
  cursor: pointer;
  transition: all 0.3s ease !important;
  border: 1px solid transparent;
}

.article-card:hover {
  border-color: rgba(149, 117, 205, 0.3);
  transform: translateY(-2px);
}

.article-title {
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.article-summary {
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.category-scroll {
  overflow-x: auto;
  padding-bottom: 8px;
}

.category-scroll::-webkit-scrollbar {
  height: 4px;
}

.load-more-btn {
  border-radius: 20px;
  padding: 0 32px;
}

.empty-state {
  background: linear-gradient(135deg, #FAFAFA 0%, #F3E5F5 100%);
}
</style>
