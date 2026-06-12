<script setup>
import { ref, inject, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { authApi, articleApi, userApi, favoriteApi } from '../api'

const router = useRouter()
const user = inject('user')

const profile = ref(null)
const articles = ref([])
const favorites = ref([])
const following = ref([])
const followers = ref([])
const activeTab = ref('articles')
const isLoading = ref(false)

const loadProfile = async () => {
  isLoading.value = true
  try {
    const response = await authApi.getProfile()
    profile.value = response.data
    
    const [articlesRes, favoritesRes, followingRes, followersRes] = await Promise.all([
      articleApi.getMyArticles({ page: 1, page_size: 10 }),
      favoriteApi.getFavorites({ page: 1, page_size: 10 }),
      userApi.getFollowing(),
      userApi.getFollowers()
    ])
    
    articles.value = articlesRes.data.articles || []
    favorites.value = favoritesRes.data.articles || []
    following.value = followingRes.data.users || []
    followers.value = followersRes.data.users || []
  } catch (error) {
    console.error('加载个人信息失败:', error)
  } finally {
    isLoading.value = false
  }
}

const handleFollow = async (userId) => {
  try {
    await userApi.follow(userId)
    loadProfile()
  } catch (error) {
    console.error('关注失败:', error)
  }
}

const handleUnfollow = async (userId) => {
  try {
    await userApi.unfollow(userId)
    loadProfile()
  } catch (error) {
    console.error('取消关注失败:', error)
  }
}

const formatTime = (timeStr) => {
  const date = new Date(timeStr)
  return date.toLocaleDateString('zh-CN')
}

onMounted(() => {
  if (!user.value) {
    router.push('/login')
    return
  }
  loadProfile()
})
</script>

<template>
  <v-app>
    <v-app-bar app>
      <v-btn icon @click="router.push('/')">
        <v-icon>mdi-arrow-left</v-icon>
      </v-btn>
      <v-toolbar-title>个人中心</v-toolbar-title>
    </v-app-bar>
    
    <v-container class="py-6">
      <v-card v-if="profile" class="mb-6">
        <v-card-text class="text-center">
          <v-avatar size="120" color="primary" class="mx-auto mb-4">
            <v-icon size="60" color="white">mdi-account</v-icon>
          </v-avatar>
          <h2 class="text-h4 font-bold">{{ profile.username }}</h2>
          <p class="text-grey">{{ profile.email }}</p>
          <div class="flex justify-center gap-8 mt-4">
            <div>
              <div class="text-2xl font-bold">{{ articles.length }}</div>
              <div class="text-sm text-grey">文章</div>
            </div>
            <div>
              <div class="text-2xl font-bold">{{ following.length }}</div>
              <div class="text-sm text-grey">关注</div>
            </div>
            <div>
              <div class="text-2xl font-bold">{{ followers.length }}</div>
              <div class="text-sm text-grey">粉丝</div>
            </div>
          </div>
        </v-card-text>
      </v-card>
      
      <!-- Tab切换 -->
      <v-tabs v-model="activeTab" background-color="transparent" centered>
        <v-tab value="articles">我的文章</v-tab>
        <v-tab value="favorites">我的收藏</v-tab>
        <v-tab value="following">我关注的</v-tab>
        <v-tab value="followers">我的粉丝</v-tab>
      </v-tabs>
      
      <!-- 我的文章 -->
      <v-card v-if="activeTab === 'articles'" class="mt-4">
        <v-list v-if="articles.length > 0">
          <v-list-item
            v-for="article in articles"
            :key="article.id"
            @click="router.push(`/article/${article.id}`)"
          >
            <v-list-item-content>
              <v-list-item-title>{{ article.title }}</v-list-item-title>
              <v-list-item-subtitle>{{ formatTime(article.created_at) }}</v-list-item-subtitle>
            </v-list-item-content>
          </v-list-item>
        </v-list>
        <v-card-text v-else class="text-center text-grey py-8">
          暂无文章
        </v-card-text>
      </v-card>
      
      <!-- 我的收藏 -->
      <v-card v-if="activeTab === 'favorites'" class="mt-4">
        <v-list v-if="favorites.length > 0">
          <v-list-item
            v-for="article in favorites"
            :key="article.id"
            @click="router.push(`/article/${article.id}`)"
          >
            <v-list-item-content>
              <v-list-item-title>{{ article.title }}</v-list-item-title>
              <v-list-item-subtitle>{{ article.author?.username }}</v-list-item-subtitle>
            </v-list-item-content>
          </v-list-item>
        </v-list>
        <v-card-text v-else class="text-center text-grey py-8">
          暂无收藏
        </v-card-text>
      </v-card>
      
      <!-- 我关注的 -->
      <v-card v-if="activeTab === 'following'" class="mt-4">
        <v-list v-if="following.length > 0">
          <v-list-item v-for="u in following" :key="u.id">
            <v-list-item-avatar>
              <v-icon color="primary">mdi-account</v-icon>
            </v-list-item-avatar>
            <v-list-item-content>
              <v-list-item-title>{{ u.username }}</v-list-item-title>
            </v-list-item-content>
            <v-list-item-actions>
              <v-btn text color="error" @click="handleUnfollow(u.id)">取消关注</v-btn>
            </v-list-item-actions>
          </v-list-item>
        </v-list>
        <v-card-text v-else class="text-center text-grey py-8">
          暂无关注
        </v-card-text>
      </v-card>
      
      <!-- 我的粉丝 -->
      <v-card v-if="activeTab === 'followers'" class="mt-4">
        <v-list v-if="followers.length > 0">
          <v-list-item v-for="u in followers" :key="u.id">
            <v-list-item-avatar>
              <v-icon color="primary">mdi-account</v-icon>
            </v-list-item-avatar>
            <v-list-item-content>
              <v-list-item-title>{{ u.username }}</v-list-item-title>
            </v-list-item-content>
            <v-list-item-actions>
              <v-btn text color="primary" @click="handleFollow(u.id)">关注</v-btn>
            </v-list-item-actions>
          </v-list-item>
        </v-list>
        <v-card-text v-else class="text-center text-grey py-8">
          暂无粉丝
        </v-card-text>
      </v-card>
    </v-container>
  </v-app>
</template>
