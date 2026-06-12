<script setup>
import { ref, inject, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { authApi, articleApi, favoriteApi, signinApi } from '../api'

const router = useRouter()
const user = inject('user')

const profile = ref(null)
const articles = ref([])
const favorites = ref([])
const activeTab = ref('articles')
const isLoading = ref(false)
const signinStatus = ref({
  hasSignedIn: false,
  signInDays: 0,
  totalSignIns: 0
})
const editMode = ref(false)
const editedProfile = ref({
  display_name: '',
  bio: ''
})

const loadProfile = async () => {
  isLoading.value = true
  try {
    const response = await authApi.getProfile()
    profile.value = response.data
    editedProfile.value = {
      display_name: profile.value.display_name,
      bio: profile.value.bio || ''
    }
  } catch (error) {
    console.error('加载个人资料失败:', error)
  } finally {
    isLoading.value = false
  }
}

const loadArticles = async () => {
  try {
    const response = await articleApi.getArticles({ user_id: profile.value.id })
    articles.value = response.data.articles
  } catch (error) {
    console.error('加载文章失败:', error)
  }
}

const loadFavorites = async () => {
  try {
    const response = await favoriteApi.getFavorites()
    favorites.value = response.data.favorites
  } catch (error) {
    console.error('加载收藏失败:', error)
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

const handleUpdateProfile = async () => {
  try {
    await authApi.updateProfile(editedProfile.value)
    profile.value.display_name = editedProfile.value.display_name
    profile.value.bio = editedProfile.value.bio
    editMode.value = false
  } catch (error) {
    console.error('更新资料失败:', error)
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

const handleTabChange = (tab) => {
  activeTab.value = tab
  if (tab === 'articles') {
    loadArticles()
  } else if (tab === 'favorites') {
    loadFavorites()
  }
}

onMounted(() => {
  if (!user.value) {
    router.push('/login')
    return
  }
  loadProfile()
  loadSigninStatus()
})
</script>

<template>
  <v-container class="max-w-4xl mx-auto px-4 py-8">
    <!-- 用户信息卡片 -->
    <v-card rounded="2xl" elevation="4" class="mb-8 overflow-hidden" v-if="profile">
      <div class="gradient-purple p-8">
        <div class="flex flex-col sm:flex-row items-center sm:items-start">
          <v-avatar size="120" color="white" class="mb-4 sm:mb-0 avatar-hover">
            <v-icon size="60" color="primary">mdi-account</v-icon>
          </v-avatar>
          <div class="ml-0 sm:ml-6 text-center sm:text-left">
            <h2 class="text-2xl font-bold text-white">{{ profile.display_name || profile.username }}</h2>
            <p class="text-white/80 mt-1">@{{ profile.username }}</p>
            <p class="text-white/60 text-sm mt-2 max-w-md">{{ profile.bio || '暂无简介' }}</p>
          </div>
        </div>
        
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 mt-8">
          <div class="text-center p-3 bg-white/10 rounded-xl">
            <span class="text-2xl font-bold text-white">{{ profile.article_count }}</span>
            <p class="text-white/60 text-sm mt-1">文章</p>
          </div>
          <div 
            class="text-center p-3 bg-white/10 rounded-xl cursor-pointer hover:bg-white/20 transition-colors"
            @click="router.push('/follow-list?type=follower')"
          >
            <span class="text-2xl font-bold text-white">{{ profile.follower_count }}</span>
            <p class="text-white/60 text-sm mt-1">粉丝</p>
          </div>
          <div 
            class="text-center p-3 bg-white/10 rounded-xl cursor-pointer hover:bg-white/20 transition-colors"
            @click="router.push('/follow-list?type=following')"
          >
            <span class="text-2xl font-bold text-white">{{ profile.following_count }}</span>
            <p class="text-white/60 text-sm mt-1">关注</p>
          </div>
          <div class="text-center p-3 bg-white/10 rounded-xl">
            <span class="text-2xl font-bold text-white">{{ signinStatus.totalSignIns }}</span>
            <p class="text-white/60 text-sm mt-1">签到</p>
          </div>
        </div>
      </div>
      
      <v-card-actions class="justify-center py-4 bg-gray-50">
        <v-btn 
          v-if="!signinStatus.hasSignedIn"
          class="btn-gradient"
          @click="handleSignin"
          size="large"
        >
          <v-icon class="mr-2" size="18">mdi-calendar-check</v-icon>
          立即签到
        </v-btn>
        <div v-else class="flex items-center text-success">
          <v-icon class="mr-2" size="20">mdi-check-circle</v-icon>
          <span class="font-medium">今日已签到</span>
          <span class="text-gray-400 mx-2">·</span>
          <span>连续 {{ signinStatus.signInDays }} 天</span>
        </div>
      </v-card-actions>
    </v-card>
    
    <!-- Tab 切换 -->
    <v-card rounded="2xl" elevation="4" class="overflow-hidden">
      <v-tabs 
        v-model="activeTab" 
        background-color="surface"
        class="border-b border-gray-100"
        @change="handleTabChange"
      >
        <v-tab 
          value="articles" 
          class="text-gray-600 hover:text-primary transition-colors"
        >
          <v-icon class="mr-2" size="18">mdi-file</v-icon>
          我的文章
        </v-tab>
        <v-tab 
          value="favorites" 
          class="text-gray-600 hover:text-primary transition-colors"
        >
          <v-icon class="mr-2" size="18">mdi-bookmark</v-icon>
          我的收藏
        </v-tab>
        <v-tab 
          value="settings" 
          class="text-gray-600 hover:text-primary transition-colors"
        >
          <v-icon class="mr-2" size="18">mdi-settings</v-icon>
          个人设置
        </v-tab>
      </v-tabs>
      
      <v-tabs-items v-model="activeTab" class="p-6">
        <!-- 我的文章 -->
        <v-tab-item value="articles">
          <v-list v-if="articles.length > 0" class="space-y-3">
            <v-card 
              v-for="article in articles" 
              :key="article.id"
              rounded="xl"
              class="card-hover cursor-pointer"
              @click="router.push(`/article/${article.id}`)"
            >
              <v-card-title class="px-6 py-4">
                <h3 class="font-bold text-gray-800">{{ article.title }}</h3>
              </v-card-title>
              <v-card-text class="px-6 pb-2">
                <p class="text-gray-600 text-sm card-text-limit">{{ article.content }}</p>
              </v-card-text>
              <v-card-subtitle class="px-6 py-3 bg-gray-50 flex items-center justify-between">
                <span class="text-gray-500 text-sm">{{ article.category?.name }}</span>
                <span class="text-gray-400 text-xs">{{ formatTime(article.created_at) }}</span>
              </v-card-subtitle>
            </v-card>
          </v-list>
          <div v-else class="empty-state">
            <v-icon size="96" color="gray-200" class="empty-state-icon">mdi-file-question</v-icon>
            <p class="text-gray-400">暂无文章</p>
            <v-btn class="btn-gradient mt-4" @click="router.push('/create')">
              <v-icon class="mr-2" size="18">mdi-plus</v-icon>
              发布第一篇文章
            </v-btn>
          </div>
        </v-tab-item>
        
        <!-- 我的收藏 -->
        <v-tab-item value="favorites">
          <v-list v-if="favorites.length > 0" class="space-y-3">
            <v-card 
              v-for="item in favorites" 
              :key="item.article.id"
              rounded="xl"
              class="card-hover cursor-pointer"
              @click="router.push(`/article/${item.article.id}`)"
            >
              <v-card-title class="px-6 py-4">
                <h3 class="font-bold text-gray-800">{{ item.article.title }}</h3>
              </v-card-title>
              <v-card-text class="px-6 pb-2">
                <p class="text-gray-500 text-sm">作者：{{ item.article.user?.display_name || item.article.user?.username }}</p>
              </v-card-text>
              <v-card-subtitle class="px-6 py-3 bg-gray-50 flex items-center justify-between">
                <span class="text-gray-500 text-sm">{{ item.article.category?.name }}</span>
                <span class="text-gray-400 text-xs">{{ formatTime(item.created_at) }}</span>
              </v-card-subtitle>
            </v-card>
          </v-list>
          <div v-else class="empty-state">
            <v-icon size="96" color="gray-200" class="empty-state-icon">mdi-bookmark-outline</v-icon>
            <p class="text-gray-400">暂无收藏</p>
            <p class="text-gray-400 text-sm mt-1">收藏喜欢的文章，方便随时查看</p>
          </div>
        </v-tab-item>
        
        <!-- 个人设置 -->
        <v-tab-item value="settings">
          <div v-if="!editMode">
            <v-card-text>
              <v-list>
                <v-list-item class="list-item-hover">
                  <v-list-item-icon>
                    <v-icon color="primary">mdi-face</v-icon>
                  </v-list-item-icon>
                  <v-list-item-content>
                    <v-list-item-title>昵称</v-list-item-title>
                    <v-list-item-subtitle>{{ profile.display_name || '未设置' }}</v-list-item-subtitle>
                  </v-list-item-content>
                </v-list-item>
                <v-list-item class="list-item-hover">
                  <v-list-item-icon>
                    <v-icon color="primary">mdi-file-text</v-icon>
                  </v-list-item-icon>
                  <v-list-item-content>
                    <v-list-item-title>个人简介</v-list-item-title>
                    <v-list-item-subtitle>{{ profile.bio || '未设置' }}</v-list-item-subtitle>
                  </v-list-item-content>
                </v-list-item>
                <v-list-item class="list-item-hover">
                  <v-list-item-icon>
                    <v-icon color="primary">mdi-email</v-icon>
                  </v-list-item-icon>
                  <v-list-item-content>
                    <v-list-item-title>QQ邮箱</v-list-item-title>
                    <v-list-item-subtitle>{{ profile.qq_number }}@qq.com</v-list-item-subtitle>
                  </v-list-item-content>
                </v-list-item>
              </v-list>
            </v-card-text>
            <v-card-actions class="justify-center py-4">
              <v-btn class="btn-gradient" @click="editMode = true">
                <v-icon class="mr-2" size="18">mdi-edit</v-icon>
                编辑资料
              </v-btn>
            </v-card-actions>
          </div>
          
          <div v-else>
            <v-card-text>
              <v-text-field
                v-model="editedProfile.display_name"
                label="昵称"
                rounded="xl"
                color="primary"
                hide-details="auto"
                class="mb-4"
              />
              <v-textarea
                v-model="editedProfile.bio"
                label="个人简介"
                rows="3"
                rounded="xl"
                color="primary"
                hide-details="auto"
                class="mb-6"
              />
              <div class="flex gap-3">
                <v-btn class="btn-gradient" @click="handleUpdateProfile">
                  <v-icon class="mr-2" size="18">mdi-check</v-icon>
                  保存设置
                </v-btn>
                <v-btn text color="gray-500" @click="editMode = false">
                  <v-icon class="mr-2" size="18">mdi-close</v-icon>
                  取消
                </v-btn>
              </div>
            </v-card-text>
          </div>
        </v-tab-item>
      </v-tabs-items>
    </v-card>
  </v-container>
</template>