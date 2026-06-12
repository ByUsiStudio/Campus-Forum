<script setup>
import { ref, inject, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { profileApi, articleApi, friendApi, favoriteApi } from '../api'
import UserAvatar from '../components/UserAvatar.vue'

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
    const response = await profileApi.getProfile()
    profile.value = response.data
  } catch (error) {
    console.error('加载个人信息失败:', error)
  }
}

const loadArticles = async () => {
  try {
    const response = await articleApi.getMyArticles({ page: 1, page_size: 10 })
    articles.value = response.data.articles || []
  } catch (error) {
    console.error('加载文章失败:', error)
    articles.value = []
  }
}

const loadFavorites = async () => {
  try {
    const response = await favoriteApi.getFavorites({ page: 1, page_size: 10 })
    favorites.value = response.data.articles || []
  } catch (error) {
    console.error('加载收藏失败:', error)
    favorites.value = []
  }
}

const loadFollowing = async () => {
  try {
    const response = await friendApi.getFriends()
    following.value = response.data?.friends || []
  } catch (error) {
    console.error('加载关注列表失败:', error)
    following.value = []
  }
}

const loadFollowers = async () => {
  try {
    const response = await friendApi.getRequests()
    followers.value = response.data?.requests || []
  } catch (error) {
    console.error('加载粉丝列表失败:', error)
    followers.value = []
  }
}

const handleFollow = async (userId) => {
  try {
    await friendApi.sendRequest(userId)
    loadFollowers()
  } catch (error) {
    console.error('关注失败:', error)
  }
}

const handleUnfollow = async (userId) => {
  try {
    await friendApi.deleteFriend(userId)
    loadFollowing()
  } catch (error) {
    console.error('取消关注失败:', error)
  }
}

const handleTabChange = async (tab) => {
  activeTab.value = tab
  if (tab === 'favorites' && favorites.value.length === 0) {
    await loadFavorites()
  } else if (tab === 'following' && following.value.length === 0) {
    await loadFollowing()
  } else if (tab === 'followers' && followers.value.length === 0) {
    await loadFollowers()
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
  loadArticles()
  loadFavorites()
  loadFollowing()
  loadFollowers()
})
</script>

<template>
  <v-container class="py-6" fluid>
    <!-- 用户信息卡片 -->
    <v-card class="profile-card mb-6">
      <v-card-text class="text-center pa-6">
        <UserAvatar :user="profile" :size="100" class="mb-4" />
        <h2 class="text-h5 font-weight-bold text-grey-darken-3 mb-1">
          {{ profile?.username || user?.username || '用户' }}
        </h2>
        <p class="text-body-2 text-grey mb-4">{{ profile?.email || user?.email || '' }}</p>

        <v-row justify="center" class="stats-row">
          <v-col cols="4" class="text-center">
            <div class="text-h6 font-weight-bold text-primary">{{ articles.length }}</div>
            <div class="text-caption text-grey">文章</div>
          </v-col>
          <v-col cols="4" class="text-center">
            <div class="text-h6 font-weight-bold text-primary">{{ following.length }}</div>
            <div class="text-caption text-grey">关注</div>
          </v-col>
          <v-col cols="4" class="text-center">
            <div class="text-h6 font-weight-bold text-primary">{{ followers.length }}</div>
            <div class="text-caption text-grey">粉丝</div>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

    <!-- Tab切换 -->
    <v-tabs v-model="activeTab" background-color="transparent" color="primary" centered class="mb-4">
      <v-tab value="articles" @click="handleTabChange('articles')">我的文章</v-tab>
      <v-tab value="favorites" @click="handleTabChange('favorites')">我的收藏</v-tab>
      <v-tab value="following" @click="handleTabChange('following')">我关注的</v-tab>
      <v-tab value="followers" @click="handleTabChange('followers')">我的粉丝</v-tab>
    </v-tabs>

    <!-- 我的文章 -->
    <v-card v-if="activeTab === 'articles'" class="article-list-card">
      <div v-if="articles.length > 0">
        <v-list lines="two">
          <v-list-item
            v-for="article in articles"
            :key="article.id"
            @click="router.push(`/article/${article.id}`)"
            class="article-item"
          >
            <template #prepend>
              <v-avatar size="48" color="primary-lighten-2" class="mr-3">
                <v-img v-if="article.cover" :src="article.cover"></v-img>
                <v-icon v-else color="white">mdi-file-document-outline</v-icon>
              </v-avatar>
            </template>
            <v-list-item-title class="font-weight-medium mb-1">{{ article.title }}</v-list-item-title>
            <v-list-item-subtitle class="text-grey">
              {{ formatTime(article.created_at) }} · {{ article.view_count || 0 }} 阅读
            </v-list-item-subtitle>
            <template #append>
              <v-icon color="grey-lighten-1">mdi-chevron-right</v-icon>
            </template>
          </v-list-item>
        </v-list>
      </div>
      <v-card-text v-else class="text-center py-12 empty-state">
        <v-icon size="64" color="primary-lighten-2" class="mb-4">mdi-file-document-outline</v-icon>
        <p class="text-body-1 text-grey mb-4">暂无文章</p>
        <v-btn color="primary" to="/create">发布文章</v-btn>
      </v-card-text>
    </v-card>

    <!-- 我的收藏 -->
    <v-card v-if="activeTab === 'favorites'" class="article-list-card">
      <div v-if="favorites.length > 0">
        <v-list lines="two">
          <v-list-item
            v-for="article in favorites"
            :key="article.id"
            @click="router.push(`/article/${article.id}`)"
            class="article-item"
          >
            <template #prepend>
              <v-avatar size="48" color="error-lighten-2" class="mr-3">
                <v-icon color="white">mdi-heart</v-icon>
              </v-avatar>
            </template>
            <v-list-item-title class="font-weight-medium mb-1">{{ article.title }}</v-list-item-title>
            <v-list-item-subtitle class="text-grey">
              {{ article.author?.username || '未知作者' }} · {{ formatTime(article.created_at) }}
            </v-list-item-subtitle>
            <template #append>
              <v-icon color="grey-lighten-1">mdi-chevron-right</v-icon>
            </template>
          </v-list-item>
        </v-list>
      </div>
      <v-card-text v-else class="text-center py-12 empty-state">
        <v-icon size="64" color="primary-lighten-2" class="mb-4">mdi-heart-outline</v-icon>
        <p class="text-body-1 text-grey">暂无收藏</p>
      </v-card-text>
    </v-card>

    <!-- 我关注的 -->
    <v-card v-if="activeTab === 'following'" class="article-list-card">
      <div v-if="following.length > 0">
        <v-list lines="two">
          <v-list-item
            v-for="u in following"
            :key="u.id"
            class="article-item"
          >
            <template #prepend>
              <UserAvatar :user="u" :size="48" class="mr-3" />
            </template>
            <v-list-item-title class="font-weight-medium">{{ u.username }}</v-list-item-title>
            <v-list-item-subtitle class="text-grey">{{ u.bio || '暂无简介' }}</v-list-item-subtitle>
            <template #append>
              <v-btn text color="error" size="small" @click.stop="handleUnfollow(u.id)">取消关注</v-btn>
            </template>
          </v-list-item>
        </v-list>
      </div>
      <v-card-text v-else class="text-center py-12 empty-state">
        <v-icon size="64" color="primary-lighten-2" class="mb-4">mdi-account-multiple-outline</v-icon>
        <p class="text-body-1 text-grey">暂无关注</p>
      </v-card-text>
    </v-card>

    <!-- 我的粉丝 -->
    <v-card v-if="activeTab === 'followers'" class="article-list-card">
      <div v-if="followers.length > 0">
        <v-list lines="two">
          <v-list-item
            v-for="u in followers"
            :key="u.id"
            class="article-item"
          >
            <template #prepend>
              <UserAvatar :user="u" :size="48" class="mr-3" />
            </template>
            <v-list-item-title class="font-weight-medium">{{ u.username }}</v-list-item-title>
            <v-list-item-subtitle class="text-grey">{{ u.bio || '暂无简介' }}</v-list-item-subtitle>
            <template #append>
              <v-btn text color="primary" size="small" @click.stop="handleFollow(u.id)">关注</v-btn>
            </template>
          </v-list-item>
        </v-list>
      </div>
      <v-card-text v-else class="text-center py-12 empty-state">
        <v-icon size="64" color="primary-lighten-2" class="mb-4">mdi-account-group-outline</v-icon>
        <p class="text-body-1 text-grey">暂无粉丝</p>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<style scoped>
.profile-card {
  border-radius: 16px !important;
  border: 1px solid rgba(149, 117, 205, 0.12) !important;
}

.article-list-card {
  border-radius: 16px !important;
  border: 1px solid rgba(149, 117, 205, 0.08) !important;
}

.article-item {
  transition: background-color 0.2s ease;
}

.article-item:hover {
  background-color: rgba(149, 117, 205, 0.04);
}

.empty-state {
  background: linear-gradient(135deg, #FAFAFA 0%, #F3E5F5 100%);
}

.stats-row {
  margin-top: 16px;
}
</style>
