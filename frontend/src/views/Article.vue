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
const replyContent = ref('')
const replyTarget = ref(null)

const loadArticle = async () => {
  isLoading.value = true
  try {
    const response = await articleApi.getArticle(route.params.id)
    article.value = response.data
    comments.value = response.data.comments || []
    
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
    loadArticle()
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
    loadArticle()
  } catch (error) {
    console.error('回复失败:', error)
  }
}

const deleteComment = async (commentId) => {
  if (!user.value) return
  try {
    await commentApi.deleteComment(commentId)
    loadArticle()
  } catch (error) {
    console.error('删除评论失败:', error)
  }
}

const formatTime = (timeStr) => {
  const date = new Date(timeStr)
  return date.toLocaleString('zh-CN')
}

onMounted(() => {
  loadArticle()
})
</script>

<template>
  <v-app>
    <v-app-bar app>
      <v-btn icon @click="router.push('/')">
        <v-icon>mdi-arrow-left</v-icon>
      </v-btn>
      <v-toolbar-title>文章详情</v-toolbar-title>
    </v-app-bar>
    
    <v-container class="py-6">
      <v-card v-if="article">
        <v-card-title>
          <h1 class="text-h4">{{ article.title }}</h1>
        </v-card-title>
        
        <v-card-subtitle class="px-6">
          <v-chip size="small" color="primary" text-color="white">
            {{ article.category?.name || '未分类' }}
          </v-chip>
          <span class="ml-3">{{ article.author?.username }}</span>
          <span class="ml-2 text-grey">{{ formatTime(article.created_at) }}</span>
        </v-card-subtitle>
        
        <v-card-text class="px-6">
          <p>{{ article.content }}</p>
        </v-card-text>
        
        <v-card-actions class="px-6">
          <v-btn icon :color="isLiked ? 'red' : 'grey'" @click="handleLike">
            <v-icon>{{ isLiked ? 'mdi-heart' : 'mdi-heart-outline' }}</v-icon>
            <span class="ml-1">{{ article.like_count }}</span>
          </v-btn>
          
          <v-btn icon :color="isFavorited ? 'primary' : 'grey'" @click="handleFavorite">
            <v-icon>{{ isFavorited ? 'mdi-bookmark' : 'mdi-bookmark-outline' }}</v-icon>
            <span class="ml-1">收藏</span>
          </v-btn>
          
          <v-btn icon color="grey">
            <v-icon>mdi-eye</v-icon>
            <span class="ml-1">{{ article.view_count }}</span>
          </v-btn>
        </v-card-actions>
      </v-card>
      
      <!-- 评论区 -->
      <v-card class="mt-6">
        <v-card-title>
          <v-icon>mdi-comment</v-icon>
          <span class="ml-2">评论 ({{ comments.length }})</span>
        </v-card-title>
        
        <v-card-text>
          <v-textarea
            v-model="newComment"
            label="写下你的评论..."
            rows="3"
            class="mb-4"
          />
          <v-btn color="primary" @click="handleComment">发表评论</v-btn>
        </v-card-text>
        
        <!-- 评论列表 -->
        <v-list v-if="comments.length > 0">
          <template v-for="comment in comments" :key="comment.id">
            <v-list-item>
              <v-list-item-avatar>
                <v-icon color="primary">mdi-account</v-icon>
              </v-list-item-avatar>
              <v-list-item-content>
                <v-list-item-title>{{ comment.user?.username }}</v-list-item-title>
                <v-list-item-subtitle>
                  {{ comment.content }}
                  <span class="ml-2 text-sm text-grey">{{ formatTime(comment.created_at) }}</span>
                </v-list-item-subtitle>
              </v-list-item-content>
              <v-list-item-actions>
                <v-btn text small @click="replyTarget = comment.id">回复</v-btn>
                <v-btn v-if="user?.id === comment.user_id || user?.role === 'admin'" text small color="error" @click="deleteComment(comment.id)">删除</v-btn>
              </v-list-item-actions>
            </v-list-item>
            
            <!-- 回复输入框 -->
            <v-list-item v-if="replyTarget === comment.id">
              <v-textarea
                v-model="replyContent"
                placeholder="输入回复内容..."
                rows="2"
                class="flex-1"
              />
              <v-btn color="primary" @click="handleReply(comment.id)">回复</v-btn>
              <v-btn text @click="replyTarget = null">取消</v-btn>
            </v-list-item>
            
            <!-- 嵌套回复 -->
            <v-list v-if="comment.replies && comment.replies.length > 0" class="ml-8">
              <v-list-item v-for="reply in comment.replies" :key="reply.id">
                <v-list-item-avatar>
                  <v-icon color="primary">mdi-account</v-icon>
                </v-list-item-avatar>
                <v-list-item-content>
                  <v-list-item-title>{{ reply.user?.username }}</v-list-item-title>
                  <v-list-item-subtitle>
                    {{ reply.content }}
                    <span class="ml-2 text-sm text-grey">{{ formatTime(reply.created_at) }}</span>
                  </v-list-item-subtitle>
                </v-list-item-content>
                <v-list-item-actions>
                  <v-btn text small @click="replyTarget = reply.id">回复</v-btn>
                  <v-btn v-if="user?.id === reply.user_id || user?.role === 'admin'" text small color="error" @click="deleteComment(reply.id)">删除</v-btn>
                </v-list-item-actions>
              </v-list-item>
            </v-list>
          </template>
        </v-list>
        
        <v-card-text v-else class="text-center text-grey py-8">
          暂无评论，快来发表第一条评论吧！
        </v-card-text>
      </v-card>
    </v-container>
  </v-app>
</template>
