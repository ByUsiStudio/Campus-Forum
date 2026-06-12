<script setup>
import { ref, inject, onMounted, computed } from 'vue'
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
  return date.toLocaleString('zh-CN')
}

onMounted(() => {
  loadArticle()
  loadComments()
})
</script>

<template>
  <v-container class="max-w-4xl mx-auto py-8" v-if="article">
    <!-- 返回按钮 -->
    <v-btn 
      text 
      color="primary" 
      class="mb-6"
      @click="router.push('/')"
    >
      <v-icon class="mr-1">mdi-arrow-left</v-icon>
      返回首页
    </v-btn>
    
    <!-- 文章卡片 -->
    <v-card rounded="xl" elevation="4" class="mb-6">
      <v-card-title class="text-center pb-4">
        <h1 class="text-2xl font-bold text-gray-800">{{ article.title }}</h1>
      </v-card-title>
      
      <v-card-subtitle class="flex items-center justify-center mb-4">
        <v-avatar size="40" color="secondary">
          <v-icon>mdi-account</v-icon>
        </v-avatar>
        <span class="ml-2 font-medium">{{ article.user?.display_name || article.user?.username }}</span>
        <span class="mx-2 text-gray-400">·</span>
        <span>{{ article.category?.name }}</span>
        <span class="mx-2 text-gray-400">·</span>
        <span>{{ formatTime(article.created_at) }}</span>
      </v-card-subtitle>
      
      <v-divider class="mx-4 mb-4"></v-divider>
      
      <v-card-text class="px-8">
        <p class="text-gray-700 leading-relaxed whitespace-pre-wrap text-lg">
          {{ article.content }}
        </p>
      </v-card-text>
      
      <v-card-actions class="justify-center py-4">
        <v-btn 
          :color="isLiked ? 'primary' : 'gray'"
          icon
          @click="handleLike"
          class="mx-4"
        >
          <v-icon :size="24">{{ isLiked ? 'mdi-heart' : 'mdi-heart-outline' }}</v-icon>
        </v-btn>
        <span>{{ article.like_count }}</span>
        
        <v-btn 
          :color="isFavorited ? 'primary' : 'gray'"
          icon
          @click="handleFavorite"
          class="mx-4"
        >
          <v-icon :size="24">{{ isFavorited ? 'mdi-bookmark' : 'mdi-bookmark-outline' }}</v-icon>
        </v-btn>
        <span>{{ article.favorite_count }}</span>
        
        <v-btn 
          color="gray"
          icon
          class="mx-4"
        >
          <v-icon size="24">mdi-eye</v-icon>
        </v-btn>
        <span>{{ article.view_count }}</span>
        
        <v-btn 
          color="gray"
          icon
          class="mx-4"
        >
          <v-icon size="24">mdi-message</v-icon>
        </v-btn>
        <span>{{ article.comment_count }}</span>
      </v-card-actions>
    </v-card>
    
    <!-- 评论区 -->
    <v-card rounded="xl" elevation="4">
      <v-card-title class="gradient-purple text-white">
        <v-icon class="mr-2">mdi-comment</v-icon>
        <span class="font-bold">评论 ({{ comments.length }})</span>
      </v-card-title>
      
      <v-card-text>
        <!-- 发表评论 -->
        <div v-if="user" class="mb-6">
          <v-textarea
            v-model="newComment"
            label="发表评论"
            placeholder="写下你的评论..."
            rows="3"
            rounded="lg"
          ></v-textarea>
          <v-btn 
            color="primary" 
            class="mt-2"
            @click="handleComment"
            :disabled="!newComment.trim()"
          >
            发表评论
          </v-btn>
        </div>
        
        <div v-else class="text-center py-4 text-gray-500">
          <v-btn text color="primary" @click="router.push('/login')">登录后发表评论</v-btn>
        </div>
        
        <!-- 评论列表 -->
        <div v-if="comments.length > 0">
          <div 
            v-for="comment in comments" 
            :key="comment.id" 
            class="border-b border-gray-200 pb-4 mb-4 last:border-0"
          >
            <div class="flex items-start">
              <v-avatar size="40" color="secondary" class="flex-shrink-0">
                <v-icon>mdi-account</v-icon>
              </v-avatar>
              <div class="ml-3 flex-1">
                <div class="flex items-center justify-between">
                  <span class="font-medium">{{ comment.user?.display_name || comment.user?.username }}</span>
                  <span class="text-sm text-gray-400">{{ formatTime(comment.created_at) }}</span>
                </div>
                <p class="mt-2 text-gray-700">{{ comment.content }}</p>
                
                <div class="flex items-center mt-2">
                  <v-btn text color="primary" size="small" @click="handleLikeComment(comment.id)">
                    <v-icon size="16" class="mr-1">mdi-heart</v-icon>
                    {{ comment.like_count }}
                  </v-btn>
                  <v-btn 
                    text 
                    color="primary" 
                    size="small" 
                    @click="replyTarget = comment"
                    class="ml-4"
                  >
                    <v-icon size="16" class="mr-1">mdi-reply</v-icon>
                    回复
                  </v-btn>
                </div>
                
                <!-- 回复输入框 -->
                <div v-if="replyTarget?.id === comment.id" class="mt-3">
                  <v-textarea
                    v-model="replyContent"
                    placeholder="回复 {{ comment.user?.display_name }}"
                    rows="2"
                    rounded="lg"
                    class="mb-2"
                  ></v-textarea>
                  <v-btn 
                    color="primary" 
                    size="small"
                    @click="handleReply(comment.id)"
                    :disabled="!replyContent.trim()"
                  >
                    回复
                  </v-btn>
                  <v-btn 
                    text 
                    color="gray" 
                    size="small" 
                    class="ml-2"
                    @click="replyTarget = null"
                  >
                    取消
                  </v-btn>
                </div>
                
                <!-- 子回复 -->
                <div v-if="comment.replies && comment.replies.length > 0" class="mt-4 ml-4 border-l-2 border-primary pl-4">
                  <div 
                    v-for="reply in comment.replies" 
                    :key="reply.id" 
                    class="border-b border-gray-100 pb-3 mb-3 last:border-0"
                  >
                    <div class="flex items-center">
                      <v-avatar size="32" color="secondary">
                        <v-icon>mdi-account</v-icon>
                      </v-avatar>
                      <div class="ml-2">
                        <span class="font-medium">{{ reply.user?.display_name || reply.user?.username }}</span>
                        <p class="text-sm text-gray-700">{{ reply.content }}</p>
                        <span class="text-xs text-gray-400">{{ formatTime(reply.created_at) }}</span>
                      </div>
                    </div>
                    <v-btn 
                      text 
                      color="primary" 
                      size="small" 
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
        
        <div v-else class="text-center py-8 text-gray-500">
          <v-icon size="48" class="mx-auto mb-2">mdi-comment-outline</v-icon>
          <p>暂无评论，快来发表第一条评论吧</p>
        </div>
      </v-card-text>
    </v-card>
  </v-container>
  
  <div v-else-if="isLoading" class="text-center py-12">
    <v-progress-circular indeterminate color="primary" />
  </div>
</template>