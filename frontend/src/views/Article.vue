<template>
  <div v-if="articleError" class="error-page">
    <div class="error-card animate-fade-in">
      <div class="error-icon" :class="articleError.code === 404 ? 'error-icon-404' : 'error-icon-500'">
        <i class="fa-solid" :class="articleError.code === 404 ? 'fa-file-question' : 'fa-circle-exclamation'"></i>
      </div>
      <h2 class="error-code">{{ articleError.code }}</h2>
      <h3 class="error-title">{{ articleError.title }}</h3>
      <p class="error-message">{{ articleError.message }}</p>
      <p class="error-detail">{{ articleError.detail }}</p>
      <div class="error-actions">
        <button class="layui-btn layui-btn-primary" @click="router.push('/')">
          <i class="fa-solid fa-house"></i>返回首页
        </button>
        <button class="layui-btn" @click="articleError = null; loadArticle()">
          <i class="fa-solid fa-rotate-right"></i>刷新重试
        </button>
      </div>
    </div>
  </div>

  <div v-else-if="article" class="article-page">
    <div class="article-container">
      <div class="article-main">
        <div class="article-card animate-fade-in">
          <div class="article-header">
            <button class="back-btn" @click="router.back()">
              <i class="fa-solid fa-arrow-left"></i>返回
            </button>
            <h1 class="article-title">{{ article.title }}</h1>
            <div class="article-meta">
              <div 
                class="author-info" 
                @click="goToUserProfile(article.user.id)"
              >
                <img 
                  :src="article.user.avatar || `/api/users/${article.user.id}/avatar`" 
                  :alt="article.user.display_name"
                  class="author-avatar"
                >
                <span class="author-name">{{ article.user.display_name || article.user.username }}</span>
              </div>
              <span class="meta-divider">|</span>
              <span class="meta-item">
                <i class="fa-solid fa-clock"></i>{{ formatDate(article.created_at) }}
              </span>
              <span class="meta-divider">|</span>
              <span class="meta-item category-tag">
                <i class="fa-solid fa-folder"></i>{{ article.category.name }}
              </span>
              <span class="meta-divider">|</span>
              <span class="meta-item">
                <i class="fa-solid fa-eye"></i>{{ article.view_count }} 阅读
              </span>
            </div>
          </div>

          <div class="article-content" ref="contentRef" @click="handleContentClick">
            <MarkdownViewer :value="article.content" />
          </div>

          <div v-if="article.voice_url" class="voice-player">
            <div class="voice-player-header">
              <i class="fa-solid fa-volume-high"></i>
              <span>语音朗读</span>
            </div>
            <div class="voice-player-body">
              <button 
                class="play-btn" 
                @click="toggleVoicePlay"
                :class="{ playing: isPlaying }"
              >
                <i class="fa-solid" :class="isPlaying ? 'fa-pause' : 'fa-play'"></i>
              </button>
              <div class="progress-container">
                <div class="progress-bar">
                  <div class="progress-fill" :style="{ width: voiceProgress + '%' }"></div>
                </div>
                <span class="time-text">
                  {{ formatVoiceTime(currentVoiceTime) }} / {{ formatVoiceTime(voiceDuration) }}
                </span>
              </div>
            </div>
            <audio
              ref="audioRef"
              :src="article.voice_url"
              @timeupdate="onVoiceTimeUpdate"
              @loadedmetadata="onVoiceLoaded"
              @ended="onVoiceEnded"
            />
          </div>

          <div class="article-actions">
            <div class="action-group">
              <button 
                class="action-btn"
                :class="{ active: liked }"
                @click="toggleLike"
              >
                <i class="fa-solid fa-thumbs-up"></i>
                <span>{{ article.like_count }} 点赞</span>
              </button>
              <button 
                class="action-btn"
                :class="{ active: coined }"
                @click="coinArticle"
                :disabled="!token"
              >
                <i class="fa-solid fa-coins"></i>
                <span>{{ article.coin_count || 0 }} 投币</span>
              </button>
              <button 
                class="action-btn"
                :class="{ active: favorited }"
                @click="toggleFavorite"
              >
                <i class="fa-solid fa-heart"></i>
                <span>{{ article.favorite_count || 0 }} 收藏</span>
              </button>
              <button class="action-btn" @click="showShareDialog = true">
                <i class="fa-solid fa-share-nodes"></i>
                <span>分享</span>
              </button>
              <button 
                v-if="token && currentUser && currentUser.id !== article.user_id"
                class="action-btn action-btn-report"
                @click="showReportDialog = true"
              >
                <i class="fa-solid fa-flag"></i>
                <span>举报</span>
              </button>
            </div>
            <div v-if="canEdit" class="edit-group">
              <button class="edit-btn" :to="'/create?id=' + article.id">
                <i class="fa-solid fa-pen-to-square"></i>编辑
              </button>
              <button class="delete-btn" @click="deleteArticle">
                <i class="fa-solid fa-trash"></i>删除
              </button>
            </div>
          </div>
        </div>

        <div class="comment-section">
          <div class="comment-header">
            <i class="fa-solid fa-comments"></i>
            <span>评论 ({{ comments.length }})</span>
          </div>

          <div v-if="token" class="comment-input">
            <img 
              :src="currentUser?.avatar || `/api/users/${currentUser?.id}/avatar`" 
              :alt="currentUser?.display_name"
              class="comment-avatar"
            >
            <div class="comment-input-group">
              <textarea 
                v-model="commentContent" 
                placeholder="写下你的评论..."
                class="comment-textarea"
              ></textarea>
              <div class="comment-input-footer">
                <label class="checkbox-label">
                  <input type="checkbox" v-model="commentIsAnonymous">
                  <span>匿名评论</span>
                </label>
                <button 
                  class="submit-btn" 
                  @click="submitComment" 
                  :disabled="!commentContent.trim()"
                >发表</button>
              </div>
            </div>
          </div>

          <div v-else class="comment-login">
            <span class="login-tip">登录后参与评论</span>
            <button class="login-btn" @click="router.push('/login')">登录</button>
          </div>

          <div class="comment-list">
            <div 
              v-for="comment in comments" 
              :key="comment.id" 
              class="comment-item animate-fade-in-up"
            >
              <img 
                :src="comment.user.avatar || `/api/users/${comment.user.id}/avatar`" 
                :alt="comment.user.display_name"
                class="comment-item-avatar"
                @click="goToUserProfile(comment.user.id)"
              >
              <div class="comment-item-body">
                <div class="comment-item-header">
                  <span class="comment-item-author">
                    {{ comment.user.display_name || comment.user.username || '匿名用户' }}
                  </span>
                  <span class="comment-item-time">{{ formatDate(comment.created_at) }}</span>
                </div>
                <p class="comment-item-content">{{ comment.content }}</p>
                <div class="comment-item-actions">
                  <button 
                    class="comment-action"
                    :class="{ liked: commentLiked[comment.id] }"
                    @click="toggleCommentLike(comment)"
                  >
                    <i class="fa-solid fa-thumbs-up"></i>{{ comment.like_count }}
                  </button>
                  <button v-if="token" class="comment-action" @click="showReplyForm(comment.id)">
                    <i class="fa-solid fa-reply"></i>回复
                  </button>
                  <button 
                    v-if="canDeleteComment(comment)" 
                    class="comment-action comment-action-delete"
                    @click="deleteComment(comment.id, comment)"
                  >
                    <i class="fa-solid fa-trash"></i>
                  </button>
                </div>

                <div v-if="replyingTo === comment.id" class="reply-form">
                  <textarea 
                    v-model="replyContent"
                    :placeholder="'回复 ' + (comment.user.display_name || comment.user.username || '匿名用户')"
                    class="reply-textarea"
                  ></textarea>
                  <div class="reply-form-footer">
                    <label class="checkbox-label">
                      <input type="checkbox" v-model="replyIsAnonymous">
                      <span>匿名</span>
                    </label>
                    <button class="reply-btn" @click="submitReply(comment.id)">发送</button>
                    <button class="cancel-btn" @click="cancelReply">取消</button>
                  </div>
                </div>

                <CommentReply
                  v-if="comment.replies && comment.replies.length > 0"
                  :replies="comment.replies"
                  :commentLiked="commentLiked"
                  :token="token"
                  :currentUser="currentUser"
                  :replyingTo="replyingTo"
                  :localReplyContent="replyContent"
                  :localReplyIsAnonymous="replyIsAnonymous"
                  @toggleLike="toggleCommentLike"
                  @showReplyForm="showReplyForm"
                  @deleteComment="deleteComment"
                  @goToUserProfile="goToUserProfile"
                  @submitReply="handleNestedReply"
                  @cancelReply="cancelReply"
                />
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="article-sidebar">
        <div class="sidebar-card">
          <div class="sidebar-header">
            <img 
              :src="article.user.avatar || `/api/users/${article.user.id}/avatar`" 
              :alt="article.user.display_name"
              class="sidebar-avatar"
              @click="goToUserProfile(article.user.id)"
            >
            <div class="sidebar-user-info">
              <div class="sidebar-user-name">{{ article.user.display_name }}</div>
              <div class="sidebar-user-signature">{{ article.user.signature || '暂无签名' }}</div>
            </div>
          </div>
          <button 
            v-if="token && currentUser && currentUser.id !== article.user_id"
            class="follow-btn"
            :class="{ followed: followStatus.is_following }"
            @click="handleFollow"
          >
            <i class="fa-solid" :class="followStatus.is_following ? 'fa-check' : 'fa-plus'"></i>
            {{ followStatus.is_following ? '已关注' : followStatus.is_followed ? '回关' : '关注' }}
          </button>
        </div>
      </div>
    </div>
  </div>

  <div v-else class="loading-page">
  <div class="loading-spinner">
    <i class="fa-solid fa-spinner fa-spin"></i>
  </div>
</div>

  <div v-if="showShareDialog" class="dialog-overlay" @click.self="showShareDialog = false">
    <div class="dialog-content animate-scale-in">
      <div class="dialog-header">
        <span>分享文章</span>
        <button class="dialog-close" @click="showShareDialog = false">
          <i class="fa-solid fa-xmark"></i>
        </button>
      </div>
      <div class="dialog-body">
        <p>复制以下链接分享给好友：</p>
        <div class="share-url-container">
          <input type="text" v-model="shareUrl" readonly class="share-url-input">
          <button class="copy-btn" @click="copyShareUrl">
            <i class="fa-solid fa-copy"></i>复制
          </button>
        </div>
        <span v-if="copySuccess" class="copy-success">已复制</span>
      </div>
      <div class="dialog-footer">
        <button class="layui-btn layui-btn-primary" @click="showShareDialog = false">关闭</button>
      </div>
    </div>
  </div>

  <div v-if="showReportDialog" class="dialog-overlay" @click.self="closeReportDialog">
    <div class="dialog-content animate-scale-in">
      <div class="dialog-header report-header">
        <i class="fa-solid fa-flag"></i>
        <span>举报文章</span>
        <button class="dialog-close" @click="closeReportDialog">
          <i class="fa-solid fa-xmark"></i>
        </button>
      </div>
      <div class="dialog-body">
        <div class="report-tip">
          <i class="fa-solid fa-info-circle"></i>
          <span>感谢您对平台环境的维护。我们会认真审核每一条举报，并在3个工作日内处理。</span>
        </div>
        <div class="report-reasons">
          <span class="reasons-title">请选择举报原因 <span class="required">*</span></span>
          <div class="reasons-list">
            <button 
              v-for="reason in reportReasons" 
              :key="reason.value"
              class="reason-chip"
              :class="{ selected: reportReason === reason.value }"
              @click="reportReason = reason.value"
            >
              <i :class="reason.icon"></i>
              {{ reason.title }}
            </button>
          </div>
        </div>
        <textarea 
          v-model="reportDescription"
          placeholder="请详细描述您举报的原因，包括具体内容和违规证据..."
          class="report-textarea"
        ></textarea>
      </div>
      <div class="dialog-footer">
        <button class="layui-btn layui-btn-primary" @click="closeReportDialog">取消</button>
        <button 
          class="layui-btn layui-btn-danger" 
          @click="submitReport" 
          :disabled="submittingReport"
        >
          <i v-if="submittingReport" class="layui-icon layui-icon-loading layui-anim layui-anim-spin"></i>
          提交举报
        </button>
      </div>
    </div>
  </div>

  <ImageViewer v-if="showImageViewer" :url="currentImageUrl" @close="closeImageViewer" />
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, computed, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '../api'
import ImageViewer from '../components/ImageViewer.vue'
import MarkdownViewer from '../components/MarkdownViewer.vue'
import CommentReply from '../components/CommentReply.vue'
import { confirm as showConfirm, prompt as showPrompt, success as showSuccess } from '../utils/modal'

const route = useRoute()
const router = useRouter()
const article = ref(null)
const articleError = ref(null)
const comments = ref([])
const liked = ref(false)
const favorited = ref(false)
const commentContent = ref('')
const commentIsAnonymous = ref(false)
const commentLiked = ref({})
const replyingTo = ref(null)
const replyContent = ref('')
const replyIsAnonymous = ref(false)
const showImageViewer = ref(false)
const currentImageUrl = ref('')
const showShareDialog = ref(false)
const shareUrl = ref('')
const copySuccess = ref(false)
const showReportDialog = ref(false)
const reportReason = ref('')
const reportDescription = ref('')
const submittingReport = ref(false)
const reportReasons = [
  { title: '垃圾广告', value: '垃圾广告', icon: 'fa-solid fa-trash' },
  { title: '色情低俗', value: '色情低俗', icon: 'fa-solid fa-triangle-exclamation' },
  { title: '暴力血腥', value: '暴力血腥', icon: 'fa-solid fa-swords' },
  { title: '政治敏感', value: '政治敏感', icon: 'fa-solid fa-user' },
  { title: '违法犯罪', value: '违法犯罪', icon: 'fa-solid fa-scale-unbalanced' },
  { title: '谣言虚假', value: '谣言虚假', icon: 'fa-solid fa-message' },
  { title: '侵犯隐私', value: '侵犯隐私', icon: 'fa-solid fa-eye-slash' },
  { title: '其他违规', value: '其他违规', icon: 'fa-solid fa-circle-question' }
]
const token = ref(localStorage.getItem('token'))
const currentUser = ref(null)
const contentRef = ref(null)
const followStatus = ref({
  is_following: false,
  is_followed: false,
  mutual: false
})
const siteTitle = ref('校园论坛')
const coined = ref(false)
const audioRef = ref(null)
const isPlaying = ref(false)
const voiceProgress = ref(0)
const voiceDuration = ref(0)
const currentVoiceTime = ref(0)
const voiceVolume = ref(0.8)

const canEdit = computed(() => {
  if (!currentUser.value || !article.value) return false
  return currentUser.value.id === article.value.user_id || currentUser.value.role === 'admin'
})

const toggleVoicePlay = () => {
  if (!audioRef.value) return
  if (isPlaying.value) {
    audioRef.value.pause()
  } else {
    audioRef.value.play()
  }
  isPlaying.value = !isPlaying.value
}

const onVoiceTimeUpdate = () => {
  if (!audioRef.value) return
  currentVoiceTime.value = audioRef.value.currentTime
  if (audioRef.value.duration) {
    voiceProgress.value = (audioRef.value.currentTime / audioRef.value.duration) * 100
  }
}

const onVoiceLoaded = () => {
  if (!audioRef.value) return
  voiceDuration.value = audioRef.value.duration
  audioRef.value.volume = voiceVolume.value
}

const onVoiceEnded = () => {
  isPlaying.value = false
  voiceProgress.value = 0
  currentVoiceTime.value = 0
}

const formatVoiceTime = (seconds) => {
  if (!seconds || isNaN(seconds)) return '00:00'
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
}

const loadArticle = async () => {
  try {
    if (!route.params.id) {
      throw new Error('文章ID为空')
    }
    
    const articleRes = await api.get(`/articles/${route.params.id}`)
    const articleData = articleRes.data.article || articleRes.data
    
    if (!articleData) {
      throw new Error('文章数据为空')
    }
    
    article.value = articleData
    comments.value = articleRes.data.comments || []
    liked.value = articleRes.data.liked || false
    commentLiked.value = articleRes.data.comment_liked || {}
    
    try {
      const siteConfigRes = await api.get('/site-config')
      siteTitle.value = siteConfigRes.data.site_title || '校园论坛'
    } catch (configError) {
      siteTitle.value = '校园论坛'
    }
    
    if (token.value && article.value.id) {
      try {
        const favoriteRes = await api.get(`/articles/${article.value.id}/favorite/check`)
        favorited.value = favoriteRes.data.favorited || false
      } catch (error) {
        favorited.value = false
      }
    }
    
    shareUrl.value = `${window.location.origin}/article/${article.value.id}`
    document.title = `${article.value.title} - ${siteTitle.value}`

    await nextTick()
    if (audioRef.value) {
      audioRef.value.load()
    }
    
    initVideoPlayers()
    loadFollowStatus()
  } catch (error) {
    const status = error.response?.status
    if (status === 404) {
      articleError.value = {
        code: 404,
        title: '资源未找到',
        message: '请求的文章不存在',
        detail: error.response?.data?.error || '该文章可能已被删除或ID错误'
      }
    } else {
      articleError.value = {
        code: status || 500,
        title: '加载失败',
        message: '加载文章时发生错误',
        detail: error.response?.data?.error || error.message
      }
    }
  }
}

const initVideoPlayers = () => {
  if (!contentRef.value) return
  
  const videoElements = contentRef.value.querySelectorAll('video')
  videoElements.forEach((videoEl) => {
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
  }
}

const loadFollowStatus = async () => {
  if (!article.value || !token.value) return

  try {
    const response = await api.get(`/friends/status/${article.value.user_id}`)
    followStatus.value = {
      is_following: response.data.is_friend,
      is_followed: response.data.is_friend,
      mutual: response.data.is_friend
    }
  } catch (error) {
    console.error('加载好友状态失败', error)
  }
}

const handleFollow = async () => {
  if (!token.value) {
    router.push('/login')
    return
  }

  try {
    if (followStatus.value.is_following) {
      await api.delete(`/friends/${article.value.user_id}`)
      followStatus.value.is_following = false
      followStatus.value.mutual = false
    } else {
      await api.post('/friends/request', { user_id: article.value.user_id })
      await showSuccess('已发送好友请求')
    }
  } catch (error) {
    console.error('好友操作失败', error)
  }
}

const goToUserProfile = (userId) => {
  router.push(`/profile?id=${userId}`)
}

const coinArticle = async () => {
  if (!token.value) {
    router.push('/login')
    return
  }

  try {
    await api.post(`/articles/${article.value.id}/coin`)
    article.value.coin_count = (article.value.coin_count || 0) + 1
    coined.value = true
    await showSuccess('投币成功')
  } catch (error) {
    console.error('投币失败', error)
  }
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
  }
}

const showReplyForm = (commentId) => {
  replyingTo.value = commentId
  replyContent.value = ''
}

const cancelReply = () => {
  replyingTo.value = null
  replyContent.value = ''
  replyIsAnonymous.value = false
}

const submitComment = async () => {
  if (!commentContent.value.trim()) return

  try {
    await api.post(`/articles/${article.value.id}/comments`, {
      content: commentContent.value,
      is_anonymous: commentIsAnonymous.value
    })
    await loadArticle()
    commentContent.value = ''
    commentIsAnonymous.value = false
  } catch (error) {
    console.error('评论失败', error)
  }
}

const submitReply = async (parentId) => {
  if (!replyContent.value.trim()) return
  
  try {
    await api.post(`/articles/${article.value.id}/comments`, {
      content: replyContent.value,
      parent_id: parentId,
      is_anonymous: replyIsAnonymous.value
    })
    await loadArticle()
    cancelReply()
  } catch (error) {
    console.error('回复失败', error)
  }
}

const handleNestedReply = async ({ parentId, content, isAnonymous }) => {
  if (!content.trim()) return
  
  try {
    await api.post(`/articles/${article.value.id}/comments`, {
      content: content,
      parent_id: parentId,
      is_anonymous: isAnonymous
    })
    await loadArticle()
    cancelReply()
  } catch (error) {
    console.error('回复失败', error)
  }
}

const deleteArticle = async () => {
  try {
    const confirmed = await showConfirm('确定要删除这篇文章吗？', {
      title: '确认删除',
      icon: 'fa-solid fa-circle-exclamation',
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
  if (!canDeleteComment(comment)) {
    return
  }
  try {
    const confirmed = await showConfirm('确定要删除这条评论吗？', {
      title: '确认删除',
      icon: 'fa-solid fa-circle-exclamation',
      iconColor: 'error'
    })
    
    if (!confirmed) return
    
    try {
      await api.delete(`/comments/${commentId}`)
      for (let i = 0; i < comments.value.length; i++) {
        if (comments.value[i].id === commentId) {
          comments.value.splice(i, 1)
          break
        }
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

const submitReport = async () => {
  if (!reportReason.value) {
    alert('请选择举报原因')
    return
  }

  if (!reportDescription.value.trim()) {
    alert('请填写详细说明')
    return
  }

  submittingReport.value = true
  try {
    await api.post('/reports', {
      target_type: 'article',
      target_id: article.value.id,
      reason: reportReason.value,
      description: reportDescription.value.trim()
    })
    showSuccess('举报已提交，感谢您的反馈')
    closeReportDialog()
  } catch (error) {
    console.error('提交举报失败', error)
    alert('提交失败，请重试')
  } finally {
    submittingReport.value = false
  }
}

const closeReportDialog = () => {
  showReportDialog.value = false
  reportReason.value = ''
  reportDescription.value = ''
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
</script>

<style lang="less" scoped>
.error-page {
  min-height: 60vh;
  display: flex;
  align-items: center;
  justify-content: center;

  .error-card {
    text-align: center;
    padding: 60px 40px;
    background: #fff;
    border-radius: 16px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);

    .error-icon {
      width: 80px;
      height: 80px;
      margin: 0 auto 20px;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 40px;

      &-404 {
        background: rgba(30, 159, 255, 0.1);
        color: #1E9FFF;
      }

      &-500 {
        background: rgba(255, 87, 34, 0.1);
        color: #FF5722;
      }
    }

    .error-code {
      font-size: 48px;
      font-weight: 700;
      color: #333;
      margin: 0 0 8px;
    }

    .error-title {
      font-size: 20px;
      font-weight: 600;
      color: #333;
      margin: 0 0 8px;
    }

    .error-message {
      font-size: 14px;
      color: #666;
      margin: 0 0 4px;
    }

    .error-detail {
      font-size: 13px;
      color: #999;
      margin: 0 0 24px;
    }

    .error-actions {
      display: flex;
      gap: 12px;
      justify-content: center;
    }
  }
}

.article-page {
  padding: 30px 0;

  .article-container {
    display: flex;
    gap: 30px;
    max-width: 1200px;
    margin: 0 auto;

    .article-main {
      flex: 1;
      min-width: 0;
    }

    .article-sidebar {
      width: 280px;
      flex-shrink: 0;
    }
  }
}

.article-card {
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
  margin-bottom: 24px;
  overflow: hidden;

  .article-header {
    padding: 24px;
    border-bottom: 1px solid #f0f0f0;

    .back-btn {
      display: inline-flex;
      align-items: center;
      gap: 4px;
      padding: 6px 12px;
      background: #f5f5f5;
      border: none;
      border-radius: 6px;
      color: #666;
      font-size: 14px;
      cursor: pointer;
      transition: all 0.3s @ease-out-back;
      margin-bottom: 16px;

      &:hover {
        background: #e8e8e8;
        color: #333;
      }
    }

    .article-title {
      font-size: 24px;
      font-weight: 700;
      color: #333;
      margin: 0 0 16px;
      line-height: 1.4;
    }

    .article-meta {
      display: flex;
      align-items: center;
      flex-wrap: wrap;
      gap: 8px;

      .author-info {
        display: flex;
        align-items: center;
        gap: 8px;
        cursor: pointer;

        .author-avatar {
          width: 36px;
          height: 36px;
          border-radius: 50%;
          object-fit: cover;
          border: 2px solid #f0f0f0;
        }

        .author-name {
          font-size: 14px;
          font-weight: 500;
          color: #333;
        }
      }

      .meta-divider {
        color: #e0e0e0;
      }

      .meta-item {
        display: flex;
        align-items: center;
        gap: 4px;
        font-size: 13px;
        color: #999;

        &.category-tag {
          background: rgba(30, 159, 255, 0.1);
          color: #1E9FFF;
          padding: 2px 8px;
          border-radius: 4px;
        }
      }
    }
  }

  .article-content {
    padding: 24px;
    font-size: 16px;
    line-height: 1.8;
    color: #333;

    :deep(h1), :deep(h2), :deep(h3), :deep(h4), :deep(h5), :deep(h6) {
      font-weight: 600;
      margin: 24px 0 16px;
      color: #333;
    }

    :deep(p) {
      margin-bottom: 16px;
    }

    :deep(pre) {
      background: #f8f9fa;
      border-radius: 8px;
      padding: 16px;
      overflow-x: auto;
      margin: 16px 0;
    }

    :deep(code) {
      background: #f0f0f0;
      padding: 2px 6px;
      border-radius: 4px;
      font-size: 14px;
    }

    :deep(pre code) {
      background: none;
      padding: 0;
    }

    :deep(blockquote) {
      border-left: 4px solid #1E9FFF;
      padding-left: 16px;
      margin: 16px 0;
      color: #666;
      background: rgba(30, 159, 255, 0.05);
      padding: 12px 16px;
      border-radius: 0 8px 8px 0;
    }

    :deep(img) {
      max-width: 100%;
      border-radius: 8px;
      cursor: pointer;
      margin: 8px 0;
    }

    :deep(ul), :deep(ol) {
      padding-left: 24px;
      margin: 16px 0;
    }

    :deep(li) {
      margin-bottom: 8px;
    }
  }

  .voice-player {
    margin: 0 24px 24px;
    padding: 16px;
    background: #f8f9fa;
    border-radius: 8px;

    .voice-player-header {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 14px;
      color: #333;
      margin-bottom: 12px;
    }

    .voice-player-body {
      display: flex;
      align-items: center;
      gap: 16px;

      .play-btn {
        width: 44px;
        height: 44px;
        border-radius: 50%;
        background: linear-gradient(135deg, #1E9FFF, #0086E6);
        border: none;
        color: #fff;
        font-size: 20px;
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
        transition: all 0.3s @ease-out-back;
        box-shadow: 0 4px 12px rgba(30, 159, 255, 0.4);

        &:hover {
          transform: scale(1.1);
        }

        &.playing {
          background: linear-gradient(135deg, #FF5722, #E64A19);
          box-shadow: 0 4px 12px rgba(255, 87, 34, 0.4);
        }
      }

      .progress-container {
        flex: 1;

        .progress-bar {
          height: 6px;
          background: #e0e0e0;
          border-radius: 3px;
          overflow: hidden;

          .progress-fill {
            height: 100%;
            background: linear-gradient(90deg, #1E9FFF, #0086E6);
            border-radius: 3px;
            transition: width 0.1s linear;
          }
        }

        .time-text {
          display: block;
          font-size: 12px;
          color: #999;
          margin-top: 4px;
          text-align: right;
        }
      }
    }
  }

  .article-actions {
    padding: 20px 24px;
    border-top: 1px solid #f0f0f0;
    display: flex;
    justify-content: space-between;
    align-items: center;

    .action-group {
      display: flex;
      gap: 16px;

      .action-btn {
        display: flex;
        align-items: center;
        gap: 6px;
        padding: 8px 16px;
        background: #f5f5f5;
        border: none;
        border-radius: 6px;
        color: #666;
        font-size: 14px;
        cursor: pointer;
        transition: all 0.3s @ease-out-back;

        &:hover {
          background: #e8e8e8;
          color: #333;
        }

        &.active {
          background: rgba(30, 159, 255, 0.1);
          color: #1E9FFF;
        }

        &:disabled {
          opacity: 0.5;
          cursor: not-allowed;
        }

        &.action-btn-report {
          background: rgba(255, 87, 34, 0.1);
          color: #FF5722;

          &:hover {
            background: rgba(255, 87, 34, 0.2);
          }
        }
      }
    }

    .edit-group {
      display: flex;
      gap: 8px;

      .edit-btn {
        display: flex;
        align-items: center;
        gap: 4px;
        padding: 8px 16px;
        background: #f5f5f5;
        border: none;
        border-radius: 6px;
        color: #666;
        font-size: 14px;
        cursor: pointer;
        transition: all 0.3s @ease-out-back;

        &:hover {
          background: #e8e8e8;
          color: #333;
        }
      }

      .delete-btn {
        display: flex;
        align-items: center;
        gap: 4px;
        padding: 8px 16px;
        background: rgba(255, 87, 34, 0.1);
        border: none;
        border-radius: 6px;
        color: #FF5722;
        font-size: 14px;
        cursor: pointer;
        transition: all 0.3s @ease-out-back;

        &:hover {
          background: rgba(255, 87, 34, 0.2);
        }
      }
    }
  }
}

.comment-section {
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
  overflow: hidden;

  .comment-header {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 20px 24px;
    border-bottom: 1px solid #f0f0f0;
    font-size: 16px;
    font-weight: 600;
    color: #333;
  }

  .comment-input {
    display: flex;
    gap: 16px;
    padding: 20px 24px;
    border-bottom: 1px solid #f0f0f0;

    .comment-avatar {
      width: 48px;
      height: 48px;
      border-radius: 50%;
      object-fit: cover;
      flex-shrink: 0;
    }

    .comment-input-group {
      flex: 1;

      .comment-textarea {
        width: 100%;
        padding: 12px;
        border: 1px solid #e0e0e0;
        border-radius: 8px;
        font-size: 14px;
        resize: none;
        min-height: 80px;
        transition: border-color 0.3s ease;

        &:focus {
          outline: none;
          border-color: #1E9FFF;
        }
      }

      .comment-input-footer {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-top: 12px;

        .checkbox-label {
          display: flex;
          align-items: center;
          gap: 6px;
          font-size: 14px;
          color: #666;
          cursor: pointer;
        }

        .submit-btn {
          padding: 8px 24px;
          background: linear-gradient(135deg, #1E9FFF, #0086E6);
          color: #fff;
          border: none;
          border-radius: 6px;
          font-size: 14px;
          font-weight: 500;
          cursor: pointer;
          transition: all 0.3s @ease-out-back;

          &:hover:not(:disabled) {
            transform: translateY(-2px);
            box-shadow: 0 4px 12px rgba(30, 159, 255, 0.4);
          }

          &:disabled {
            opacity: 0.5;
            cursor: not-allowed;
          }
        }
      }
    }
  }

  .comment-login {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 12px;
    padding: 24px;
    border-bottom: 1px solid #f0f0f0;

    .login-tip {
      font-size: 14px;
      color: #999;
    }

    .login-btn {
      padding: 6px 16px;
      background: #1E9FFF;
      color: #fff;
      border: none;
      border-radius: 4px;
      font-size: 14px;
      cursor: pointer;

      &:hover {
        background: #0086E6;
      }
    }
  }

  .comment-list {
    padding: 16px 24px;

    .comment-item {
      display: flex;
      gap: 16px;
      padding: 16px 0;
      border-bottom: 1px solid #f5f5f5;

      &:last-child {
        border-bottom: none;
      }

      .comment-item-avatar {
        width: 44px;
        height: 44px;
        border-radius: 50%;
        object-fit: cover;
        flex-shrink: 0;
        cursor: pointer;
      }

      .comment-item-body {
        flex: 1;
        min-width: 0;

        .comment-item-header {
          display: flex;
          align-items: center;
          gap: 12px;
          margin-bottom: 8px;

          .comment-item-author {
            font-size: 14px;
            font-weight: 500;
            color: #333;
          }

          .comment-item-time {
            font-size: 12px;
            color: #999;
          }
        }

        .comment-item-content {
          font-size: 15px;
          color: #333;
          line-height: 1.6;
          margin: 0 0 12px;
          word-break: break-word;
        }

        .comment-item-actions {
          display: flex;
          gap: 16px;

          .comment-action {
            display: flex;
            align-items: center;
            gap: 4px;
            padding: 4px 8px;
            background: transparent;
            border: none;
            color: #999;
            font-size: 13px;
            cursor: pointer;
            transition: color 0.3s ease;

            &:hover {
              color: #1E9FFF;
            }

            &.liked {
              color: #1E9FFF;
            }

            &.comment-action-delete {
              &:hover {
                color: #FF5722;
              }
            }
          }
        }

        .reply-form {
          margin-top: 12px;
          padding: 12px;
          background: #f8f9fa;
          border-radius: 8px;

          .reply-textarea {
            width: 100%;
            padding: 10px;
            border: 1px solid #e0e0e0;
            border-radius: 6px;
            font-size: 14px;
            resize: none;
            min-height: 60px;

            &:focus {
              outline: none;
              border-color: #1E9FFF;
            }
          }

          .reply-form-footer {
            display: flex;
            align-items: center;
            gap: 12px;
            margin-top: 10px;

            .checkbox-label {
              display: flex;
              align-items: center;
              gap: 6px;
              font-size: 13px;
              color: #666;
              cursor: pointer;
            }

            .reply-btn {
              padding: 6px 16px;
              background: #1E9FFF;
              color: #fff;
              border: none;
              border-radius: 4px;
              font-size: 13px;
              cursor: pointer;

              &:hover {
                background: #0086E6;
              }
            }

            .cancel-btn {
              padding: 6px 16px;
              background: #f5f5f5;
              color: #666;
              border: none;
              border-radius: 4px;
              font-size: 13px;
              cursor: pointer;

              &:hover {
                background: #e8e8e8;
              }
            }
          }
        }
      }
    }
  }
}

.sidebar-card {
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
  padding: 20px;

  .sidebar-header {
    display: flex;
    gap: 12px;
    margin-bottom: 16px;

    .sidebar-avatar {
      width: 56px;
      height: 56px;
      border-radius: 50%;
      object-fit: cover;
      cursor: pointer;
    }

    .sidebar-user-info {
      flex: 1;

      .sidebar-user-name {
        font-size: 16px;
        font-weight: 600;
        color: #333;
        margin-bottom: 4px;
      }

      .sidebar-user-signature {
        font-size: 13px;
        color: #999;
      }
    }
  }

  .follow-btn {
    width: 100%;
    padding: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 6px;
    background: linear-gradient(135deg, #1E9FFF, #0086E6);
    color: #fff;
    border: none;
    border-radius: 6px;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.3s @ease-out-back;

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 4px 12px rgba(30, 159, 255, 0.4);
    }

    &.followed {
      background: #f5f5f5;
      color: #666;

      &:hover {
        transform: none;
        box-shadow: none;
        background: #e8e8e8;
      }
    }
  }
}

.loading-page {
  min-height: 60vh;
  display: flex;
  align-items: center;
  justify-content: center;

  .loading-spinner {
    font-size: 48px;
    color: #1E9FFF;
  }
}

.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.dialog-content {
  background: #fff;
  border-radius: 12px;
  width: 90%;
  max-width: 480px;
  overflow: hidden;

  .dialog-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px 20px;
    border-bottom: 1px solid #f0f0f0;
    font-size: 16px;
    font-weight: 600;
    color: #333;

    &.report-header {
      background: #FF5722;
      color: #fff;

      .dialog-close {
        color: rgba(255, 255, 255, 0.8);
      }
    }

    .dialog-close {
      padding: 4px;
      background: transparent;
      border: none;
      color: #999;
      font-size: 20px;
      cursor: pointer;

      &:hover {
        color: #333;
      }
    }
  }

  .dialog-body {
    padding: 20px;

    p {
      font-size: 14px;
      color: #666;
      margin: 0 0 12px;
    }

    .share-url-container {
      display: flex;
      gap: 8px;

      .share-url-input {
        flex: 1;
        padding: 10px;
        border: 1px solid #e0e0e0;
        border-radius: 6px;
        font-size: 14px;
        background: #f8f9fa;
      }

      .copy-btn {
        padding: 10px 16px;
        background: #1E9FFF;
        color: #fff;
        border: none;
        border-radius: 6px;
        font-size: 14px;
        cursor: pointer;

        &:hover {
          background: #0086E6;
        }
      }
    }

    .copy-success {
      display: inline-block;
      padding: 6px 12px;
      background: rgba(82, 196, 26, 0.1);
      color: #52C41A;
      border-radius: 4px;
      font-size: 13px;
      margin-top: 12px;
    }

    .report-tip {
      display: flex;
      align-items: flex-start;
      gap: 8px;
      padding: 12px;
      background: rgba(30, 159, 255, 0.1);
      color: #1E9FFF;
      border-radius: 6px;
      font-size: 13px;
      margin-bottom: 16px;
    }

    .report-reasons {
      margin-bottom: 16px;

      .reasons-title {
        display: block;
        font-size: 14px;
        font-weight: 500;
        color: #333;
        margin-bottom: 10px;

        .required {
          color: #FF5722;
        }
      }

      .reasons-list {
        display: flex;
        flex-wrap: wrap;
        gap: 8px;

        .reason-chip {
          display: flex;
          align-items: center;
          gap: 6px;
          padding: 6px 12px;
          background: #f5f5f5;
          border: 1px solid #e0e0e0;
          border-radius: 20px;
          font-size: 13px;
          color: #666;
          cursor: pointer;
          transition: all 0.3s @ease-out-back;

          &:hover {
            border-color: #1E9FFF;
            color: #1E9FFF;
          }

          &.selected {
            background: rgba(30, 159, 255, 0.1);
            border-color: #1E9FFF;
            color: #1E9FFF;
          }
        }
      }
    }

    .report-textarea {
      width: 100%;
      padding: 12px;
      border: 1px solid #e0e0e0;
      border-radius: 8px;
      font-size: 14px;
      resize: none;
      min-height: 100px;

      &:focus {
        outline: none;
        border-color: #1E9FFF;
      }
    }
  }

  .dialog-footer {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    padding: 16px 20px;
    border-top: 1px solid #f0f0f0;
  }
}

@media (max-width: 768px) {
  .article-container {
    flex-direction: column;
    gap: 20px;
  }

  .article-sidebar {
    width: 100%;
  }

  .article-header {
    padding: 16px !important;

    .article-title {
      font-size: 20px !important;
    }

    .article-meta {
      flex-wrap: wrap;
      gap: 6px;
    }
  }

  .article-content {
    padding: 16px !important;
  }

  .article-actions {
    padding: 16px !important;
    flex-direction: column;
    gap: 12px;
    align-items: flex-start !important;

    .action-group {
      flex-wrap: wrap;
      gap: 8px;
    }
  }

  .comment-input {
    padding: 16px !important;
  }
}
</style>
