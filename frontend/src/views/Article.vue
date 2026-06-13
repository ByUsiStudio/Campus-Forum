<script setup>
import { ref, inject, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { articleApi, commentApi, favoriteApi, reportApi } from '../api'
import MarkdownViewer from '../components/MarkdownViewer.vue'
import UserAvatar from '../components/UserAvatar.vue'
import { alert, confirm } from '../utils/modal'

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
const isMobile = computed(() => window.innerWidth < 1024)

const loadArticle = async () => {
  isLoading.value = true
  try {
    const response = await articleApi.getArticle(route.params.id)
    article.value = response.data.article
    comments.value = response.data.comments || []

    if (user.value) {
      try {
        const favoriteResponse = await favoriteApi.check(route.params.id)
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
      await articleApi.unlike(route.params.id)
    } else {
      await articleApi.like(route.params.id)
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
      await favoriteApi.remove(route.params.id)
    } else {
      await favoriteApi.add(route.params.id)
    }
    isFavorited.value = !isFavorited.value
  } catch (error) {
    console.error('收藏失败:', error)
  }
}

const handleComment = async () => {
  if (!user.value || !newComment.value.trim()) return
  try {
    await commentApi.create(route.params.id, {
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
    await commentApi.create(route.params.id, {
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
    await commentApi.delete(commentId)
    loadArticle()
  } catch (error) {
    console.error('删除评论失败:', error)
  }
}

const handleReport = async () => {
  if (!user.value) {
    router.push('/login')
    return
  }
  try {
    const confirmed = await confirm('确定要举报这篇文章吗？', {
      title: '举报文章',
      icon: 'mdi-alert-circle',
      iconColor: 'warning'
    })
    if (confirmed) {
      await reportApi.create({ article_id: route.params.id, type: 'article' })
      alert('举报已提交，感谢您的反馈', {
        title: '举报成功',
        icon: 'mdi-check-circle',
        iconColor: 'success'
      })
    }
  } catch (error) {
    console.error('举报失败:', error)
    alert('举报失败，请稍后重试', {
      title: '举报失败',
      icon: 'mdi-alert-circle',
      iconColor: 'error'
    })
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

const goBack = () => {
  router.back()
}

onMounted(() => {
  loadArticle()
})
</script>

<template>
  <v-container class="py-6" fluid>
    <v-row justify="center">
      <v-col cols="12" :md="isMobile ? 12 : 9">
        <!-- 返回按钮 -->
        <v-btn
          variant="text"
          color="primary"
          class="mb-4"
          @click="goBack"
        >
          <v-icon class="mr-1">mdi-arrow-left</v-icon>
          返回
        </v-btn>

        <!-- 加载状态 -->
        <v-card v-if="isLoading" class="pa-8 text-center">
          <v-progress-circular indeterminate color="primary"></v-progress-circular>
        </v-card>

        <!-- 文章内容 -->
        <v-card v-else-if="article" class="article-card">
          <!-- 文章头部 -->
          <v-card-text class="pa-6 pb-0">
            <!-- 分类标签 -->
            <v-chip
              size="small"
              variant="tonal"
              color="primary"
              class="mb-3"
            >
              {{ article.category?.name || '未分类' }}
            </v-chip>

            <!-- 标题 -->
            <h1 class="text-h4 text-md-h3 font-weight-bold text-grey-darken-3 mb-4 article-title">
              {{ article.title }}
            </h1>

            <!-- 作者信息 -->
            <div class="d-flex align-center flex-wrap mb-4">
              <UserAvatar :user="article.user" :size="40" class="mr-3" />
              <div>
                <div class="text-subtitle-1 font-weight-medium">
                  {{ article.user?.display_name || article.user?.username || '未知用户' }}
                </div>
                <div class="text-caption text-grey">
                  发布于 {{ formatTime(article.created_at) }}
                  <span class="mx-2">·</span>
                  阅读 {{ article.view_count || 0 }}
                </div>
              </div>
            </div>
          </v-card-text>

          <v-divider class="mx-6"></v-divider>

          <!-- 文章正文 -->
          <v-card-text class="pa-6">
            <MarkdownViewer :value="article.content" />
          </v-card-text>

          <v-divider class="mx-6"></v-divider>

          <!-- 操作栏 -->
          <v-card-text class="pa-4">
            <div class="d-flex align-center justify-space-between">
              <div class="d-flex align-center">
                <!-- 点赞 -->
                <v-btn
                  variant="tonal"
                  :color="isLiked ? 'error' : 'grey'"
                  :class="{ 'like-btn': true, 'liked': isLiked }"
                  @click="handleLike"
                  class="mr-4"
                >
                  <v-icon :class="{ 'mr-1': true, 'text-error': isLiked }">
                    {{ isLiked ? 'mdi-heart' : 'mdi-heart-outline' }}
                  </v-icon>
                  <span>{{ article.like_count || 0 }}</span>
                </v-btn>

                <!-- 收藏 -->
                <v-btn
                  variant="tonal"
                  :color="isFavorited ? 'primary' : 'grey'"
                  @click="handleFavorite"
                >
                  <v-icon :class="{ 'mr-1': true, 'text-primary': isFavorited }">
                    {{ isFavorited ? 'mdi-bookmark' : 'mdi-bookmark-outline' }}
                  </v-icon>
                  <span>{{ isFavorited ? '已收藏' : '收藏' }}</span>
                </v-btn>
              </div>

              <!-- 举报按钮 -->
              <v-btn variant="text" color="grey" @click="handleReport">
                <v-icon>mdi-flag-outline</v-icon>
              </v-btn>
            </div>
          </v-card-text>
        </v-card>

        <!-- 评论区 -->
        <v-card class="mt-6 comment-card">
          <v-card-title class="d-flex align-center pa-4">
            <v-icon color="primary" class="mr-2">mdi-comment-outline</v-icon>
            <span class="text-subtitle-1 font-weight-bold">评论</span>
            <v-chip size="small" variant="tonal" color="primary" class="ml-2">
              {{ comments.length }}
            </v-chip>
          </v-card-title>

          <v-divider></v-divider>

          <!-- 评论输入框 -->
          <v-card-text class="pa-4" v-if="user">
            <div class="d-flex align-start">
              <UserAvatar :user="user" :size="36" class="mr-3" />
              <div class="flex-grow-1">
                <v-textarea
                  v-model="newComment"
                  placeholder="写下你的想法..."
                  variant="outlined"
                  rows="3"
                  hide-details
                  class="mb-3"
                ></v-textarea>
                <div class="d-flex justify-end">
                  <v-btn
                    color="primary"
                    @click="handleComment"
                    :disabled="!newComment.trim()"
                  >
                    发布评论
                  </v-btn>
                </div>
              </div>
            </div>
          </v-card-text>

          <v-card-text v-else class="text-center py-6">
            <p class="text-grey mb-4">登录后才能评论</p>
            <v-btn color="primary" to="/login">登录</v-btn>
          </v-card-text>

          <v-divider></v-divider>

          <!-- 评论列表 -->
          <div v-if="comments.length > 0" class="comment-list">
            <div
              v-for="comment in comments"
              :key="comment.id"
              class="comment-item"
            >
              <div class="d-flex pa-4">
                <UserAvatar :user="comment.user" :size="36" class="mr-3" />
                <div class="flex-grow-1">
                  <div class="d-flex align-center mb-1">
                    <span class="text-subtitle-2 font-weight-medium mr-2">
                                          {{ comment.user?.display_name || comment.user?.username || '匿名用户' }}
                    </span>
                    <span class="text-caption text-grey">
                      {{ formatTime(comment.created_at) }}
                    </span>
                  </div>
                  <p class="text-body-2 mb-2 comment-content">{{ comment.content }}</p>
                  <div class="d-flex align-center">
                    <v-btn
                      variant="text"
                      size="small"
                      color="grey"
                      @click="replyTarget = comment.id"
                    >
                      <v-icon size="14" class="mr-1">mdi-reply</v-icon>
                      回复
                    </v-btn>
                    <v-btn
                      v-if="user?.id === comment.user_id || user?.role === 'admin'"
                      variant="text"
                      size="small"
                      color="error"
                      @click="deleteComment(comment.id)"
                    >
                      删除
                    </v-btn>
                  </div>

                  <!-- 回复输入框 -->
                  <div v-if="replyTarget === comment.id" class="mt-3">
                    <v-textarea
                      v-model="replyContent"
                      placeholder="输入回复内容..."
                      variant="outlined"
                      rows="2"
                      hide-details
                      class="mb-2"
                    ></v-textarea>
                    <div class="d-flex justify-end">
                      <v-btn size="small" variant="text" @click="replyTarget = null">
                        取消
                      </v-btn>
                      <v-btn
                        size="small"
                        color="primary"
                        :disabled="!replyContent.trim()"
                        @click="handleReply(comment.id)"
                        class="ml-2"
                      >
                        回复
                      </v-btn>
                    </div>
                  </div>

                  <!-- 嵌套回复 -->
                  <div
                    v-if="comment.replies && comment.replies.length > 0"
                    class="reply-list mt-3"
                  >
                    <div
                      v-for="reply in comment.replies"
                      :key="reply.id"
                      class="reply-item"
                    >
                      <div class="d-flex pa-3">
                        <UserAvatar :user="reply.user" :size="28" class="mr-2" />
                        <div class="flex-grow-1">
                          <div class="d-flex align-center mb-1">
                            <span class="text-body-2 font-weight-medium mr-2">
                              {{ reply.user?.display_name || reply.user?.username || '匿名用户' }}
                            </span>
                            <span class="text-caption text-grey">
                              {{ formatTime(reply.created_at) }}
                            </span>
                          </div>
                          <p class="text-body-2 mb-1 comment-content">{{ reply.content }}</p>
                          <div class="d-flex align-center">
                            <v-btn
                              variant="text"
                              size="x-small"
                              color="grey"
                              @click="replyTarget = reply.id"
                            >
                              回复
                            </v-btn>
                            <v-btn
                              v-if="user?.id === reply.user_id || user?.role === 'admin'"
                              variant="text"
                              size="x-small"
                              color="error"
                              @click="deleteComment(reply.id)"
                            >
                              删除
                            </v-btn>
                          </div>

                          <!-- 回复回复输入框 -->
                          <div v-if="replyTarget === reply.id" class="mt-2">
                            <v-textarea
                              v-model="replyContent"
                              placeholder="输入回复内容..."
                              variant="outlined"
                              rows="2"
                              hide-details
                              class="mb-2"
                            ></v-textarea>
                            <div class="d-flex justify-end">
                              <v-btn size="small" variant="text" @click="replyTarget = null">
                                取消
                              </v-btn>
                              <v-btn
                                size="small"
                                color="primary"
                                :disabled="!replyContent.trim()"
                                @click="handleReply(reply.id)"
                                class="ml-2"
                              >
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
              <v-divider></v-divider>
            </div>
          </div>

          <v-card-text v-else class="text-center py-12">
            <v-icon size="64" color="primary-lighten-2" class="mb-4">mdi-comment-outline</v-icon>
            <p class="text-body-1 text-grey">暂无评论，快来发表第一条评论吧</p>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>
.article-card {
  border-radius: 16px !important;
}

.article-title {
  line-height: 1.3;
}

.comment-card {
  border-radius: 16px !important;
}

.comment-item {
  transition: background-color 0.2s ease;
}

.comment-item:hover {
  background-color: rgba(149, 117, 205, 0.04);
}

.comment-content {
  line-height: 1.6;
  white-space: pre-wrap;
}

.reply-list {
  background-color: rgba(149, 117, 205, 0.04);
  border-radius: 12px;
}

.reply-item {
  transition: background-color 0.2s ease;
}

.reply-item:hover {
  background-color: rgba(149, 117, 205, 0.08);
}

.like-btn.liked {
  background-color: rgba(239, 83, 80, 0.1) !important;
}
</style>
