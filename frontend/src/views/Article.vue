<script setup>
import { ref, inject, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { articleApi, commentApi, favoriteApi } from '../api'

const router = useRouter()
const route = useRoute()
const user = inject('user')

const article = ref(null)
const comments = ref([])
const isLoading = ref(false)
const isLiked = ref(false)
const isFavorited = ref(false)
const newComment = ref('')
const replyTarget = ref(null)
const replyContent = ref('')

const loadArticle = async () => {
  isLoading.value = true
  try {
    const response = await articleApi.getArticle(route.params.id)
    article.value = response.data
    
    if (user.value) {
      try {
        const favoriteResponse = await favoriteApi.checkFavorite(route.params.id)
        isFavorited.value = favoriteResponse.data.is_favorited
      } catch (error) {
        console.error('检查收藏状态失败:', error)
      }
    }
  } catch (error) {
    console.error('加载文章失败:', error)
    router.push('/')
  } finally {
    isLoading.value = false
  }
}

const loadComments = async () => {
  try {
    const response = await articleApi.getArticle(route.params.id)
    comments.value = response.data.comments || []
  } catch (error) {
    console.error('加载评论失败:', error)
  }
}

const handleLike = async () => {
  if (!user.value) {
    router.push('/login')
    return
  }
  try {
    if (isLiked.value) {
      await articleApi.unlikeArticle(route.params.id)
    } else {
      await articleApi.likeArticle(route.params.id)
    }
    isLiked.value = !isLiked.value
    article.value.like_count += isLiked.value ? 1 : -1
  } catch (error) {
    console.error('点赞失败:', error)
  }
}

const handleFavorite = async () => {
  if (!user.value) {
    router.push('/login')
    return
  }
  try {
    if (isFavorited.value) {
      await favoriteApi.removeFavorite(route.params.id)
    } else {
      await favoriteApi.addFavorite(route.params.id)
    }
    isFavorited.value = !isFavorited.value
  } catch (error) {
    console.error('收藏失败:', error)
  }
}

const handleComment = async () => {
  if (!user.value || !newComment.value.trim()) return
  try {
    await commentApi.createComment(route.params.id, {
      content: newComment.value
    })
    newComment.value = ''
    loadComments()
  } catch (error) {
    console.error('评论失败:', error)
  }
}

const handleReply = async (parentId) => {
  if (!user.value || !replyContent.value.trim()) return
  try {
    await commentApi.createComment(route.params.id, {
      content: replyContent.value,
      parent_id: parentId
    })
    replyContent.value = ''
    replyTarget.value = null
    loadComments()
  } catch (error) {
    console.error('回复失败:', error)
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
  loadArticle()
  loadComments()
})
</script>

<template>
  <v-container class="max-w-4xl mx-auto px-4 py-8" v-if="article">
    <!-- 返回按钮 -->
    <v-btn 
      text 
      color="gray-600" 
      class="mb-6 hover:text-primary transition-colors"
      @click="router.push('/')"
    >
      <v-icon class="mr-2" size="20">mdi-arrow-left</v-icon>
      返回首页
    </v-btn>
    
    <!-- 文章卡片 -->
    <v-card rounded="2xl" elevation="4" class="mb-8 overflow-hidden">
      <v-card-title class="text-center py-6 px-8">
        <div class="flex items-center justify-center gap-2 mb-4">
          <v-chip 
            v-if="article.category" 
            size="small" 
            class="tag-purple"
          >
            {{ article.category.name }}
          </v-chip>
        </div>
        <h1 class="text-2xl md:text-3xl font-bold text-gray-800">{{ article.title }}</h1>
      </v-card-title>
      
      <v-card-subtitle class="flex items-center justify-center pb-4 px-8">
        <v-avatar size="44" color="primary" class="avatar-hover">
          <v-icon size="20" color="white">mdi-account</v-icon>
        </v-avatar>
        <div class="ml-3 text-center">
          <span class="font-medium text-gray-800">{{ article.user?.display_name || article.user?.username }}</span>
          <p class="text-xs text-gray-400 mt-1">{{ formatTime(article.created_at) }}</p>
        </div>
      </v-card-subtitle>
      
      <v-divider class="mx-8 mb-6"></v-divider>
      
      <v-card-text class="px-8 pb-8">
        <div class="prose prose-lg max-w-none">
          <p class="text-gray-700 leading-loose whitespace-pre-wrap text-base md:text-lg">
            {{ article.content }}
          </p>
        </div>
      </v-card-text>
      
      <v-card-actions class="justify-center py-6 bg-gray-50">
        <v-btn 
          :class="isLiked ? 'text-primary' : 'text-gray-500'"
          icon
          @click="handleLike"
          class="mx-6 transition-all"
        >
          <v-icon :size="28" :class="isLiked ? 'text-primary' : ''">{{ isLiked ? 'mdi-heart' : 'mdi-heart-outline' }}</v-icon>
        </v-btn>
        <span :class="isLiked ? 'text-primary font-medium' : 'text-gray-500'">{{ article.like_count }}</span>
        
        <v-btn 
          :class="isFavorited ? 'text-primary' : 'text-gray-500'"
          icon
          @click="handleFavorite"
          class="mx-6 transition-all"
        >
          <v-icon :size="28">{{ isFavorited ? 'mdi-bookmark' : 'mdi-bookmark-outline' }}</v-icon>
        </v-btn>
        <span :class="isFavorited ? 'text-primary font-medium' : 'text-gray-500'">{{ article.favorite_count }}</span>
        
        <v-btn 
          color="gray-500"
          icon
          class="mx-6"
        >
          <v-icon size="28">mdi-eye</v-icon>
        </v-btn>
        <span class="text-gray-500">{{ article.view_count }}</span>
        
        <v-btn 
          color="gray-500"
          icon
          class="mx-6"
        >
          <v-icon size="28">mdi-message</v-icon>
        </v-btn>
        <span class="text-gray-500">{{ article.comment_count }}</span>
      </v-card-actions>
    </v-card>
    
    <!-- 评论区 -->
    <v-card rounded="2xl" elevation="4" class="overflow-hidden">
      <v-card-title class="gradient-purple text-white py-4 px-6">
        <v-icon class="mr-3" size="20">mdi-comment</v-icon>
        <span class="font-bold">评论 ({{ comments.length }})</span>
      </v-card-title>
      
      <v-card-text class="px-6 py-4">
        <!-- 发表评论 -->
        <div v-if="user" class="mb-8">
          <v-textarea
            v-model="newComment"
            label="发表评论"
            placeholder="写下你的评论..."
            rows="3"
            rounded="xl"
            color="primary"
            hide-details="auto"
          ></v-textarea>
          <v-btn 
            class="btn-gradient mt-3"
            @click="handleComment"
            :disabled="!newComment.trim()"
            :loading="false"
          >
            <v-icon class="mr-2" size="18">mdi-send</v-icon>
            发表评论
          </v-btn>
        </div>
        
        <div v-else class="text-center py-8 text-gray-400">
          <v-btn class="btn-gradient" @click="router.push('/login')">登录后发表评论</v-btn>
        </div>
        
        <!-- 评论列表 -->
        <div v-if="comments.length > 0">
          <div 
            v-for="comment in comments" 
            :key="comment.id" 
            class="border-b border-gray-100 pb-6 mb-6 last:border-0 last:mb-0"
          >
            <div class="flex items-start">
              <v-avatar size="44" color="primary" class="flex-shrink-0 avatar-hover">
                <v-icon size="20" color="white">mdi-account</v-icon>
              </v-avatar>
              <div class="ml-4 flex-1">
                <div class="flex items-center justify-between">
                  <span class="font-medium text-gray-800">{{ comment.user?.display_name || comment.user?.username }}</span>
                  <span class="text-sm text-gray-400">{{ formatTime(comment.created_at) }}</span>
                </div>
                <p class="mt-2 text-gray-700 leading-relaxed">{{ comment.content }}</p>
                
                <div class="flex items-center mt-3 gap-4">
                  <v-btn text color="gray-500" size="small" hover-color="primary">
                    <v-icon size="16" class="mr-1">mdi-heart</v-icon>
                    {{ comment.like_count }}
                  </v-btn>
                  <v-btn 
                    text 
                    color="gray-500" 
                    size="small" 
                    hover-color="primary"
                    @click="replyTarget = comment"
                  >
                    <v-icon size="16" class="mr-1">mdi-reply</v-icon>
                    回复
                  </v-btn>
                </div>
                
                <!-- 回复输入框 -->
                <div v-if="replyTarget?.id === comment.id" class="mt-4 p-4 bg-gray-50 rounded-xl">
                  <v-textarea
                    v-model="replyContent"
                    :placeholder="`回复 ${comment.user?.display_name || comment.user?.username}`"
                    rows="2"
                    rounded="lg"
                    color="primary"
                    hide-details="auto"
                    class="mb-3"
                  ></v-textarea>
                  <div class="flex gap-2">
                    <v-btn 
                      class="btn-gradient"
                      size="small"
                      @click="handleReply(comment.id)"
                      :disabled="!replyContent.trim()"
                    >
                      发送回复
                    </v-btn>
                    <v-btn 
                      text 
                      color="gray-500" 
                      size="small" 
                      @click="replyTarget = null"
                    >
                      取消
                    </v-btn>
                  </div>
                </div>
                
                <!-- 子回复 -->
                <div v-if="comment.replies && comment.replies.length > 0" class="mt-6 ml-4 pl-4 border-l-2 border-primary/30">
                  <div 
                    v-for="reply in comment.replies" 
                    :key="reply.id" 
                    class="pb-4 mb-4 last:mb-0"
                  >
                    <div class="flex items-start">
                      <v-avatar size="36" color="primary" class="flex-shrink-0">
                        <v-icon size="16" color="white">mdi-account</v-icon>
                      </v-avatar>
                      <div class="ml-3 flex-1">
                        <div class="flex items-center justify-between">
                          <span class="font-medium text-gray-700 text-sm">{{ reply.user?.display_name || reply.user?.username }}</span>
                          <span class="text-xs text-gray-400">{{ formatTime(reply.created_at) }}</span>
                        </div>
                        <p class="mt-1 text-gray-600 text-sm">{{ reply.content }}</p>
                        <v-btn 
                          text 
                          color="gray-400" 
                          size="small" 
                          hover-color="primary"
                          @click="replyTarget = reply"
                          class="mt-2"
                        >
                          <v-icon size="14" class="mr-1">mdi-reply</v-icon>
                          回复
                        </v-btn>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <div v-else class="empty-state">
          <v-icon size="96" color="gray-200" class="empty-state-icon">mdi-comment-outline</v-icon>
          <p class="text-gray-400">暂无评论</p>
          <p class="text-gray-400 text-sm mt-1">快来发表第一条评论吧</p>
        </div>
      </v-card-text>
    </v-card>
  </v-container>
  
  <div v-else-if="isLoading" class="loading-center">
    <v-progress-circular indeterminate color="primary" :size="48" />
  </div>
</template>