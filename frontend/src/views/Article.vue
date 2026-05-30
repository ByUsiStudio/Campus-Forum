<template>
  <div v-if="article" class="article-page">
    <div class="article-container">
      <header class="article-header">
        <h1 class="article-title">{{ article.title }}</h1>

        <div class="article-meta">
          <div class="author-info">
            <UserAvatar
              :user="article.user"
              :size="44"
              class="cursor-pointer"
              @click="goToUserProfile(article.user.id)"
            />
            <div class="author-details">
              <span class="author-name">{{ article.user.display_name }}</span>
              <span class="meta-separator">·</span>
              <span class="publish-date">{{ formatDate(article.created_at) }}</span>
            </div>
            <v-btn
              v-if="token && currentUser && currentUser.id !== article.user_id"
              variant="tonal"
              size="small"
              :color="followStatus.is_following ? 'default' : 'primary'"
              class="follow-btn"
              @click="handleFollow"
            >
              {{ followStatus.is_following ? '已关注' : followStatus.is_followed ? '回关' : '关注' }}
            </v-btn>
          </div>

          <div class="article-tags">
            <v-chip size="small" color="primary" variant="flat" class="category-chip">
              <v-icon start size="14">mdi-folder</v-icon>
              {{ article.category.name }}
            </v-chip>
            <span class="view-count">
              <v-icon size="14">mdi-eye</v-icon>
              {{ article.view_count }} 阅读
            </span>
          </div>
        </div>

        <div v-if="canEdit" class="article-actions">
          <v-btn variant="outlined" color="primary" size="small" :to="'/create?id=' + article.id">
            <v-icon start size="16">mdi-pencil</v-icon>
            编辑
          </v-btn>
          <v-btn variant="outlined" color="error" size="small" @click="deleteArticle">
            <v-icon start size="16">mdi-delete</v-icon>
            删除
          </v-btn>
        </div>
      </header>

      <div ref="contentRef" class="article-body" @click="handleContentClick">
        <MarkdownViewer :value="article.content" />
      </div>

      <footer class="article-footer">
        <div class="interaction-bar">
          <v-btn
            @click="toggleLike"
            :color="liked ? 'primary' : 'default'"
            :variant="liked ? 'flat' : 'outlined'"
            class="action-btn"
          >
            <v-icon start>mdi-thumb-up</v-icon>
            {{ article.like_count }}
          </v-btn>
          <v-btn
            @click="toggleFavorite"
            :color="favorited ? 'primary' : 'default'"
            :variant="favorited ? 'flat' : 'outlined'"
            class="action-btn"
          >
            <v-icon start>mdi-bookmark</v-icon>
            {{ article.favorite_count || 0 }}
          </v-btn>
          <v-btn
            @click="showShareDialog = true"
            variant="outlined"
            class="action-btn"
          >
            <v-icon start>mdi-share-variant</v-icon>
            分享
          </v-btn>
        </div>
      </footer>

      <section class="comments-section">
        <h3 class="section-title">
          <v-icon class="mr-2">mdi-comment-text</v-icon>
          评论 ({{ comments.length }})
        </h3>

        <div v-if="token" class="comment-form">
          <UserAvatar :user="currentUser" :size="40" />
          <div class="comment-input-wrapper">
            <v-textarea
              v-model="commentContent"
              placeholder="写下你的评论..."
              variant="outlined"
              rows="3"
              hide-details
              class="comment-textarea"
            />
            <v-btn color="primary" size="small" @click="submitComment" class="submit-comment-btn">
              发表
            </v-btn>
          </div>
        </div>

        <div v-else class="login-hint">
          <span>登录后参与评论</span>
          <v-btn variant="text" color="primary" size="small" @click="router.push('/login')">登录</v-btn>
        </div>

        <div class="comments-list">
          <div v-for="comment in comments" :key="comment.id" class="comment-item">
            <UserAvatar
              :user="comment.user"
              :size="40"
              class="cursor-pointer"
              @click="goToUserProfile(comment.user.id)"
            />
            <div class="comment-content-wrapper">
              <div class="comment-header">
                <span class="comment-author">{{ comment.user.display_name }}</span>
                <span class="comment-time">{{ formatDate(comment.created_at) }}</span>
              </div>
              <p class="comment-text">{{ comment.content }}</p>
              <div class="comment-actions">
                <v-btn
                  variant="text"
                  size="x-small"
                  @click="toggleCommentLike(comment)"
                  :color="commentLiked[comment.id] ? 'primary' : 'default'"
                >
                  <v-icon start size="14">mdi-thumb-up</v-icon>
                  {{ comment.like_count }}
                </v-btn>
                <v-btn
                  variant="text"
                  size="x-small"
                  @click="showReplyForm(comment.id)"
                  v-if="token"
                >
                  <v-icon start size="14">mdi-reply</v-icon>
                  回复 ({{ comment.reply_count }})
                </v-btn>
                <v-btn
                  variant="text"
                  size="x-small"
                  color="error"
                  @click="deleteComment(comment.id, comment)"
                  v-if="canDeleteComment(comment)"
                >
                  <v-icon start size="14">mdi-delete</v-icon>
                  删除
                </v-btn>
              </div>

              <div v-if="replyingTo === comment.id" class="reply-form">
                <v-textarea
                  v-model="replyContent"
                  :placeholder="'回复 ' + comment.user.display_name + '...'"
                  variant="outlined"
                  rows="2"
                  hide-details
                  density="compact"
                />
                <div class="reply-form-actions">
                  <v-btn size="small" color="primary" @click="submitReply(comment.id)">发送</v-btn>
                  <v-btn size="small" variant="text" @click="cancelReply">取消</v-btn>
                </div>
              </div>

              <div v-if="comment.replies && comment.replies.length > 0" class="replies-list">
                <div v-for="reply in comment.replies" :key="reply.id" class="reply-item">
                  <UserAvatar
                    :user="reply.user"
                    :size="32"
                    class="cursor-pointer"
                    @click="goToUserProfile(reply.user.id)"
                  />
                  <div class="reply-content-wrapper">
                    <div class="comment-header">
                      <span class="comment-author">{{ reply.user.display_name }}</span>
                      <span class="comment-time">{{ formatDate(reply.created_at) }}</span>
                    </div>
                    <p class="comment-text">{{ reply.content }}</p>
                    <div class="comment-actions">
                      <v-btn
                        variant="text"
                        size="x-small"
                        @click="toggleCommentLike(reply)"
                        :color="commentLiked[reply.id] ? 'primary' : 'default'"
                        density="compact"
                      >
                        <v-icon start size="12">mdi-thumb-up</v-icon>
                        {{ reply.like_count }}
                      </v-btn>
                      <v-btn
                        variant="text"
                        size="x-small"
                        color="error"
                        @click="deleteComment(reply.id, reply)"
                        v-if="canDeleteComment(reply)"
                        density="compact"
                      >
                        <v-icon start size="12">mdi-delete</v-icon>
                        删除
                      </v-btn>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>
    </div>

    <v-dialog v-model="showShareDialog" max-width="400">
      <v-card>
        <v-card-title class="text-h6">分享文章</v-card-title>
        <v-card-text>
          <p class="mb-3 text-body-2">复制以下链接分享给好友：</p>
          <v-text-field
            v-model="shareUrl"
            readonly
            variant="outlined"
            density="compact"
            hide-details
            append-inner-icon="mdi-content-copy"
            @click:append-inner="copyShareUrl"
          />
          <v-chip v-if="copySuccess" color="success" class="mt-3" size="small">
            已复制
          </v-chip>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showShareDialog = false">关闭</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <ImageViewer v-if="showImageViewer" :url="currentImageUrl" @close="closeImageViewer" />
  </div>

  <div v-else class="loading-container">
    <v-progress-circular indeterminate color="primary" size="48" />
  </div>
</template>

<script>
import { ref, onMounted, onBeforeUnmount, computed, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '../api'
import ImageViewer from '../components/ImageViewer.vue'
import UserAvatar from '../components/UserAvatar.vue'
import MarkdownViewer from '../components/MarkdownViewer.vue'
import { confirm as showConfirm, prompt as showPrompt, success as showSuccess } from '../utils/modal'


export default {
  name: 'Article',
  components: {
    ImageViewer,
    UserAvatar,
    MarkdownViewer
  },
  setup() {
    const route = useRoute()
    const router = useRouter()
    const article = ref(null)
    const comments = ref([])
    const liked = ref(false)
    const favorited = ref(false)
    const commentContent = ref('')
    const commentLiked = ref({})
    const replyingTo = ref(null)
    const replyContent = ref('')
    const showImageViewer = ref(false)
    const currentImageUrl = ref('')
    const showShareDialog = ref(false)
    const shareUrl = ref('')
    const copySuccess = ref(false)
    const token = ref(localStorage.getItem('token'))
    const currentUser = ref(null)
    const contentRef = ref(null)
    const followStatus = ref({
      is_following: false,
      is_followed: false,
      mutual: false
    })
    const siteTitle = ref('校园论坛')


    const canEdit = computed(() => {
      if (!currentUser.value || !article.value) return false
      return currentUser.value.id === article.value.user_id || currentUser.value.role === 'admin'
    })
    
    const loadArticle = async () => {
      try {
        const [articleRes, siteConfigRes] = await Promise.all([
          api.get(`/articles/${route.params.id}`),
          api.get('/site-config')
        ])
        
        article.value = articleRes.data.article
        comments.value = articleRes.data.comments
        liked.value = articleRes.data.liked || false
        commentLiked.value = articleRes.data.comment_liked || {}
        siteTitle.value = siteConfigRes.data.site_title || '校园论坛'
        
        // 检查收藏状态
        if (token.value) {
          try {
            const favoriteRes = await api.get(`/articles/${article.value.id}/favorite/check`)
            favorited.value = favoriteRes.data.favorited || false
          } catch (error) {
            favorited.value = false
          }
        }
        
        // 设置分享链接
        shareUrl.value = `${window.location.origin}/article/${article.value.id}`

        document.title = `${article.value.title} - ${siteTitle.value}`

        await nextTick()
        
        initVideoPlayers()
        loadFollowStatus()
      } catch (error) {
        console.error('加载文章失败', error)
        router.push('/')
      }
    }
    
    // 使用浏览器默认视频播放器
    const initVideoPlayers = () => {
      if (!contentRef.value) return
      
      const videoElements = contentRef.value.querySelectorAll('video')
      videoElements.forEach((videoEl) => {
        // 使用浏览器默认播放器
        videoEl.controls = true
        videoEl.playsInline = true
        videoEl.style.maxWidth = '100%'
        videoEl.style.maxHeight = '500px'
        videoEl.style.height = 'auto'
        videoEl.style.borderRadius = '8px'
        videoEl.style.margin = '16px 0'
        videoEl.style.display = 'block'
        videoEl.style.background = '#000'
      })
    }
    
    const toggleLike = async () => {
      if (!token.value) {
        router.push('/login')
        return
      }
      
      const isLiked = liked.value
      
      try {
        if (isLiked) {
          await api.delete(`/articles/${article.value.id}/like`)
          article.value.like_count--
          liked.value = false
        } else {
          await api.post(`/articles/${article.value.id}/like`)
          article.value.like_count++
          liked.value = true
        }
      } catch (error) {
        liked.value = isLiked
        console.error('点赞操作失败', error)
      }
    }

    const toggleFavorite = async () => {
      if (!token.value) {
        router.push('/login')
        return
      }

      const isFavorited = favorited.value

      try {
        if (isFavorited) {
          await api.delete(`/articles/${article.value.id}/favorite`)
          article.value.favorite_count--
          favorited.value = false
        } else {
          await api.post(`/articles/${article.value.id}/favorite`)
          article.value.favorite_count++
          favorited.value = true
        }
      } catch (error) {
        favorited.value = isFavorited
        console.error('收藏操作失败', error)
      }
    }

    const loadFollowStatus = async () => {
      if (!article.value || !token.value) return
      
      try {
        const response = await api.get(`/follow/status/${article.value.user_id}`)
        followStatus.value = response.data
      } catch (error) {
        console.error('加载关注状态失败', error)
      }
    }
    
    const handleFollow = async () => {
      if (!token.value) {
        router.push('/login')
        return
      }
      
      try {
        if (followStatus.value.is_following) {
          await api.delete(`/follow/${article.value.user_id}`)
          followStatus.value.is_following = false
          followStatus.value.mutual = false
        } else {
          await api.post(`/follow/${article.value.user_id}`)
          followStatus.value.is_following = true
          followStatus.value.mutual = followStatus.value.is_followed
        }
      } catch (error) {
        console.error('关注失败', error)
      }
    }
    
    const goToUserProfile = (userId) => {
      router.push(`/profile?id=${userId}`)
    }
    
    const toggleCommentLike = async (comment) => {
      if (!token.value) {
        router.push('/login')
        return
      }
      
      const isCommentLiked = commentLiked.value[comment.id]
      
      try {
        if (isCommentLiked) {
          await api.delete(`/comments/${comment.id}/like`)
          comment.like_count--
          commentLiked.value[comment.id] = false
        } else {
          await api.post(`/comments/${comment.id}/like`)
          comment.like_count++
          commentLiked.value[comment.id] = true
        }
      } catch (error) {
        commentLiked.value[comment.id] = isCommentLiked
        console.error('评论点赞操作失败', error)
      }
    }
    
    const showReplyForm = (commentId) => {
      replyingTo.value = commentId
      replyContent.value = ''
    }
    
    const cancelReply = () => {
      replyingTo.value = null
      replyContent.value = ''
    }
    
    const submitComment = async () => {
      if (!commentContent.value.trim()) return
      
      try {
        const response = await api.post(`/articles/${article.value.id}/comments`, {
          content: commentContent.value
        })
        // 新评论添加到列表
        comments.value.unshift({
          ...response.data.comment,
          replies: []
        })
        commentContent.value = ''
      } catch (error) {
        console.error('评论失败', error)
      }
    }
    
    const submitReply = async (parentId) => {
      if (!replyContent.value.trim()) return
      
      try {
        const response = await api.post(`/articles/${article.value.id}/comments`, {
          content: replyContent.value,
          parent_id: parentId
        })
        
        // 找到父评论并添加回复
        const parentComment = comments.value.find(c => c.id === parentId)
        if (parentComment) {
          if (!parentComment.replies) {
            parentComment.replies = []
          }
          parentComment.replies.push(response.data.comment)
          parentComment.reply_count++
        }
        
        cancelReply()
      } catch (error) {
        console.error('回复失败', error)
      }
    }
    
    const deleteArticle = async () => {
      try {
        const confirmed = await showConfirm('确定要删除这篇文章吗？', {
          title: '确认删除',
          icon: 'mdi-alert-circle',
          iconColor: 'error'
        })
        
        if (!confirmed) return
        
        if (currentUser.value?.role !== 'admin') {
          const reason = await showPrompt('请输入删除原因（管理员将审核）：', {
            title: '删除原因',
            inputLabel: '删除原因',
            placeholder: '请说明删除原因...',
            rows: 3
          })
          
          if (!reason) return
          
          try {
            await api.delete(`/articles/${article.value.id}`, { data: { reason } })
            await showSuccess('删除申请已提交，等待管理员审核')
            router.push('/')
          } catch (error) {
            console.error('提交删除申请失败', error)
          }
        } else {
          try {
            await api.delete(`/articles/${article.value.id}`)
            router.push('/')
          } catch (error) {
            console.error('删除失败', error)
          }
        }
      } catch (error) {
      }
    }
    
    const deleteComment = async (commentId, comment) => {
      // 双重检查权限
      if (!canDeleteComment(comment)) {
        console.error('无权限删除该评论')
        return
      }
      try {
        const confirmed = await showConfirm('确定要删除这条评论吗？', {
          title: '确认删除',
          icon: 'mdi-alert-circle',
          iconColor: 'error'
        })
        
        if (!confirmed) return
        
        try {
          await api.delete(`/comments/${commentId}`)
          // 从列表中移除
          for (let i = 0; i < comments.value.length; i++) {
            if (comments.value[i].id === commentId) {
              comments.value.splice(i, 1)
              break
            }
            // 检查是否是回复
            if (comments.value[i].replies) {
              for (let j = 0; j < comments.value[i].replies.length; j++) {
                if (comments.value[i].replies[j].id === commentId) {
                  comments.value[i].replies.splice(j, 1)
                  comments.value[i].reply_count--
                  break
                }
              }
            }
          }
        } catch (error) {
          console.error('删除评论失败', error)
        }
      } catch (error) {
      }
    }
    
    const canDeleteComment = (comment) => {
      if (!currentUser.value) return false
      const currentUserId = Number(currentUser.value.id)
      const commentUserId = Number(comment.user_id)
      return currentUserId === commentUserId || currentUser.value.role === 'admin'
    }
    
    const handleContentClick = (event) => {
      const target = event.target
      if (target.tagName === 'IMG') {
        currentImageUrl.value = target.src
        showImageViewer.value = true
      }
    }
    
    const closeImageViewer = () => {
      showImageViewer.value = false
    }
    
    const copyShareUrl = async () => {
      try {
        await navigator.clipboard.writeText(shareUrl.value)
        copySuccess.value = true
        setTimeout(() => {
          copySuccess.value = false
        }, 2000)
      } catch (error) {
        console.error('复制失败', error)
      }
    }
    
    const formatDate = (date) => {
      return new Date(date).toLocaleString('zh-CN')
    }
    
    onMounted(() => {
      const user = localStorage.getItem('user')
      if (user) {
        currentUser.value = JSON.parse(user)
      }
      loadArticle()
    })
    
    onBeforeUnmount(() => {
      document.title = siteTitle.value
    })
    
    return {
      article,
      comments,
      liked,
      favorited,
      commentContent,
      commentLiked,
      replyingTo,
      replyContent,
      token,
      currentUser,
      canEdit,
      showImageViewer,
      currentImageUrl,
      showShareDialog,
      shareUrl,
      copySuccess,
      contentRef,
      followStatus,
      toggleLike,
      toggleFavorite,
      copyShareUrl,
      toggleCommentLike,
      showReplyForm,
      cancelReply,
      submitComment,
      submitReply,
      deleteArticle,
      deleteComment,
      canDeleteComment,
      handleContentClick,
      closeImageViewer,
      formatDate,
      loadFollowStatus,
      handleFollow,
      goToUserProfile
    }
  }
}
</script>

<style scoped>
.article-page {
  max-width: 900px;
  margin: 0 auto;
  padding: 24px 16px;
}

.article-container {
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  overflow: hidden;
}

.article-header {
  padding: 32px 40px 24px;
  border-bottom: 1px solid #f0f0f0;
}

.article-title {
  font-size: 2rem;
  font-weight: 700;
  color: #1a1a1a;
  line-height: 1.4;
  margin: 0 0 20px 0;
}

.article-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 16px;
}

.author-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.author-details {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.author-name {
  font-weight: 500;
  color: #333;
  font-size: 14px;
}

.meta-separator {
  color: #ccc;
}

.publish-date {
  font-size: 13px;
  color: #888;
}

.follow-btn {
  margin-left: 8px;
}

.article-tags {
  display: flex;
  align-items: center;
  gap: 16px;
}

.category-chip {
  font-size: 12px;
}

.view-count {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #888;
}

.article-actions {
  display: flex;
  gap: 12px;
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px dashed #e5e5e5;
}

.article-body {
  padding: 32px 40px;
  min-height: 400px;
}

.article-footer {
  padding: 20px 40px;
  border-top: 1px solid #f0f0f0;
  background: #fafafa;
}

.interaction-bar {
  display: flex;
  gap: 12px;
}

.action-btn {
  min-width: 80px;
}

.comments-section {
  padding: 32px 40px;
  border-top: 1px solid #f0f0f0;
}

.section-title {
  display: flex;
  align-items: center;
  font-size: 1.1rem;
  font-weight: 600;
  color: #333;
  margin: 0 0 24px 0;
}

.comment-form {
  display: flex;
  gap: 16px;
  margin-bottom: 32px;
}

.comment-input-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.comment-textarea {
  flex: 1;
}

.submit-comment-btn {
  align-self: flex-end;
}

.login-hint {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 16px 20px;
  background: #f5f5f5;
  border-radius: 8px;
  color: #666;
  margin-bottom: 24px;
}

.comments-list {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.comment-item {
  display: flex;
  gap: 16px;
}

.comment-content-wrapper {
  flex: 1;
}

.comment-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
}

.comment-author {
  font-weight: 500;
  color: #333;
  font-size: 14px;
}

.comment-time {
  font-size: 12px;
  color: #999;
}

.comment-text {
  margin: 0;
  font-size: 14px;
  line-height: 1.7;
  color: #444;
}

.comment-actions {
  display: flex;
  gap: 16px;
  margin-top: 8px;
}

.reply-form {
  margin-top: 12px;
  padding: 12px;
  background: #f5f5f5;
  border-radius: 8px;
}

.reply-form-actions {
  display: flex;
  gap: 8px;
  margin-top: 8px;
}

.replies-list {
  margin-top: 16px;
  padding-left: 20px;
  border-left: 2px solid #e8e8e8;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.reply-item {
  display: flex;
  gap: 12px;
}

.reply-content-wrapper {
  flex: 1;
}

.loading-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 50vh;
}

@media (max-width: 768px) {
  .article-page {
    padding: 12px;
  }

  .article-header {
    padding: 20px;
  }

  .article-title {
    font-size: 1.5rem;
  }

  .article-meta {
    flex-direction: column;
    align-items: flex-start;
  }

  .article-body {
    padding: 20px;
  }

  .article-footer {
    padding: 16px 20px;
  }

  .comments-section {
    padding: 20px;
  }

  .comment-form {
    flex-direction: column;
  }

  .comment-item {
    flex-direction: column;
    gap: 8px;
  }

  .replies-list {
    padding-left: 12px;
  }
}
</style>
