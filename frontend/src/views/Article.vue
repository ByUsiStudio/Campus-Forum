<template>
  <div class="page-container article-page">
    <!-- 加载错误 -->
    <div v-if="articleError">
      <ErrorPage
        :code="articleError.code"
        :title="articleError.title"
        :message="articleError.message"
        :detail="articleError.detail"
        :actions="[
          { text: '返回首页', color: 'primary', callback: () => router.push('/') },
          { text: '刷新重试', callback: () => { articleError = null; loadArticle() } }
        ]"
      />
    </div>

    <!-- 内容 -->
    <div v-else-if="article" class="article-layout">
      <!-- 主内容 -->
      <div class="article-main">
        <el-card shadow="never" class="article-card">
          <template #header>
            <div class="article-header">
              <el-button circle text @click="router.back()">
                <el-icon><ArrowLeft /></el-icon>
              </el-button>
              <div class="header-info">
                <h1 class="article-title">{{ article.title }}</h1>
                <div class="article-meta">
                  <span class="cursor-pointer" @click="goToUserProfile(article.user?.id)">
                    <UserAvatar :user="article.user" :size="36" :show-username="true" />
                  </span>
                  <span class="text-secondary">{{ formatDate(article.created_at) }}</span>
                  <el-tag v-if="article.category" size="small">{{ article.category.name }}</el-tag>
                  <span class="text-secondary">
                    <el-icon><View /></el-icon> {{ article.view_count }} 阅读
                  </span>
                </div>
              </div>
            </div>
          </template>

          <div ref="contentRef" @click="handleContentClick">
            <MarkdownViewer :value="article.content" />
          </div>

          <!-- 语音播放器 -->
          <div v-if="article.voice_url" class="voice-player mt-4">
            <div class="voice-row">
              <el-button circle :type="isPlaying ? 'primary' : 'default'" @click="toggleVoicePlay">
                <el-icon><component :is="isPlaying ? 'VideoPause' : 'VideoPlay'" /></el-icon>
              </el-button>
              <div class="voice-body">
                <div class="voice-label text-secondary">
                  <el-icon><Headset /></el-icon> 语音朗读
                  <span class="ml-2">{{ formatVoiceTime(currentVoiceTime) }} / {{ formatVoiceTime(voiceDuration) }}</span>
                </div>
                <el-progress :percentage="voiceProgress" :stroke-width="6" :show-text="false" />
                <audio
                  ref="audioRef"
                  :src="article.voice_url"
                  @timeupdate="onVoiceTimeUpdate"
                  @loadedmetadata="onVoiceLoaded"
                  @ended="onVoiceEnded"
                  style="display:none"
                />
              </div>
            </div>
          </div>
        </el-card>

        <!-- 互动栏 -->
        <el-card shadow="never" class="interact-card">
          <div class="interact-row">
            <el-button :type="liked ? 'primary' : 'default'" @click="toggleLike">
              <el-icon><component :is="liked ? 'GoodsFilled' : 'Goods'"/></el-icon>&nbsp;{{ article.like_count }} 点赞
            </el-button>
            <el-button :type="coined ? 'warning' : 'default'" :disabled="!token" @click="coinArticle">
              <el-icon><Coin /></el-icon>&nbsp;{{ article.coin_count || 0 }} 投币
            </el-button>
            <el-button :type="favorited ? 'primary' : 'default'" @click="toggleFavorite">
              <el-icon><Star /></el-icon>&nbsp;{{ article.favorite_count || 0 }} 收藏
            </el-button>
            <el-button @click="showShareDialog = true">
              <el-icon><Share /></el-icon>&nbsp;分享
            </el-button>
            <el-button v-if="token && currentUser && currentUser.id !== article.user_id" @click="showReportDialog = true">
              <el-icon><Warning /></el-icon>&nbsp;举报
            </el-button>
            <div class="spacer"></div>
            <template v-if="canEdit">
              <el-button type="primary" plain size="small" :to="'/create?id=' + article.id">
                <el-icon><EditPen /></el-icon>&nbsp;编辑
              </el-button>
              <el-button type="danger" plain size="small" @click="deleteArticle">
                <el-icon><Delete /></el-icon>&nbsp;删除
              </el-button>
            </template>
          </div>
        </el-card>

        <!-- 评论 -->
        <el-card shadow="never" class="comment-card">
          <template #header>
            <span>
              <el-icon><ChatDotRound /></el-icon>&nbsp;评论（{{ comments.length }}）
            </span>
          </template>

          <!-- 发表评论 -->
          <div v-if="token" class="comment-input">
            <UserAvatar :user="currentUser" :size="40" :show-username="false" />
            <div class="comment-input-body">
              <el-input
                v-model="commentContent"
                type="textarea"
                :rows="3"
                placeholder="写下你的评论..."
              />
              <div class="comment-tools">
                <el-checkbox v-model="commentIsAnonymous">匿名评论</el-checkbox>
                <el-button type="primary" :disabled="!commentContent.trim()" @click="submitComment">
                  发表
                </el-button>
              </div>
            </div>
          </div>
          <div v-else class="comment-login text-center">
            <span class="text-secondary">登录后参与评论</span>
            <el-button type="primary" link @click="router.push('/login')">登录</el-button>
          </div>

          <!-- 评论列表 -->
          <div class="comment-list">
            <div v-for="comment in comments" :key="comment.id" class="comment-item">
              <UserAvatar
                :user="comment.user"
                :size="40"
                :show-username="false"
                class="cursor-pointer"
                @click="goToUserProfile(comment.user?.id)"
              />
              <div class="comment-body">
                <div class="comment-meta">
                  <span class="comment-name">{{ comment.user?.display_name || comment.user?.username || '匿名用户' }}</span>
                  <span class="text-secondary">{{ formatDate(comment.created_at) }}</span>
                </div>
                <div class="comment-content">{{ comment.content }}</div>
                <div class="comment-actions">
                  <el-button text size="small" :type="commentLiked[comment.id] ? 'primary' : 'default'" @click="toggleCommentLike(comment)">
                    <el-icon><Goods /></el-icon>&nbsp;{{ comment.like_count }}
                  </el-button>
                  <el-button v-if="token" text size="small" @click="showReplyForm(comment.id)">
                    <el-icon><ChatLineRound /></el-icon>&nbsp;回复
                  </el-button>
                  <el-button v-if="canDeleteComment(comment)" text size="small" type="danger" @click="deleteComment(comment.id, comment)">
                    <el-icon><Delete /></el-icon>
                  </el-button>
                </div>

                <!-- 回复表单 -->
                <div v-if="replyingTo === comment.id" class="reply-form">
                  <el-input
                    v-model="replyContent"
                    :placeholder="'回复 ' + (comment.user?.display_name || comment.user?.username || '匿名用户')"
                    type="textarea"
                    :rows="2"
                  />
                  <div class="reply-tools">
                    <el-checkbox v-model="replyIsAnonymous">匿名</el-checkbox>
                    <div class="spacer"></div>
                    <el-button size="small" type="primary" @click="submitReply(comment.id)">发送</el-button>
                    <el-button size="small" text @click="cancelReply">取消</el-button>
                  </div>
                </div>

                <!-- 嵌套回复 -->
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
        </el-card>
      </div>

      <!-- 侧边作者卡片 -->
      <el-card v-if="article.user" shadow="never" class="author-card" style="display:none">
      </el-card>
    </div>

    <!-- 加载中 -->
    <div v-else class="loading" style="display:flex;justify-content:center;padding:80px 0">
      <el-icon class="is-loading" :size="32"><Loading /></el-icon>
    </div>

    <!-- 分享对话框 -->
    <el-dialog v-model="showShareDialog" title="分享文章" width="420px" append-to-body>
      <el-input v-model="shareUrl" readonly>
        <template #append>
          <el-button @click="copyShareUrl">复制</el-button>
        </template>
      </el-input>
      <el-tag v-if="copySuccess" type="success" class="mt-2">已复制</el-tag>
    </el-dialog>

    <!-- 举报对话框 -->
    <el-dialog v-model="showReportDialog" title="举报文章" width="520px" append-to-body>
      <el-alert type="info" :closable="false" class="mb-3">
        感谢您对平台环境的维护。我们会认真审核每一条举报，并在3个工作日内处理。
      </el-alert>
      <div class="report-field">
        <div class="text-secondary mb-1">请选择举报原因 <span style="color:#f56c6c">*</span></div>
        <el-radio-group v-model="reportReason">
          <el-radio v-for="r in reportReasons" :key="r.value" :value="r.value">{{ r.title }}</el-radio>
        </el-radio-group>
      </div>
      <el-input
        v-model="reportDescription"
        type="textarea"
        :rows="4"
        maxlength="500"
        placeholder="请详细描述您举报的原因，包括具体内容和违规证据..."
      />
      <template #footer>
        <el-button @click="closeReportDialog">取消</el-button>
        <el-button type="danger" :loading="submittingReport" @click="submitReport">提交举报</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, computed, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, View, Headset, Coin, Star, Share, Warning, Goods, GoodsFilled, Delete, EditPen, ChatDotRound, ChatLineRound, Loading, VideoPause, VideoPlay } from '@element-plus/icons-vue'
import http from '@/api/http'
import ImageViewer from '../components/ImageViewer.vue'
import UserAvatar from '../components/UserAvatar.vue'
import MarkdownViewer from '../components/MarkdownViewer.vue'
import ErrorPage from '../components/ErrorPage.vue'
import CommentReply from '../components/CommentReply.vue'
import { confirm as showConfirm, prompt as showPrompt, success as showSuccess } from '@/utils/message'

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
const showShareDialog = ref(false)
const shareUrl = ref('')
const copySuccess = ref(false)
const showReportDialog = ref(false)
const reportReason = ref('')
const reportDescription = ref('')
const submittingReport = ref(false)
const reportReasons = [
  '垃圾广告', '色情低俗', '暴力血腥', '政治敏感',
  '违法犯罪', '谣言虚假', '侵犯隐私', '其他违规'
]
const token = ref(localStorage.getItem('token'))
const currentUser = ref(null)
const contentRef = ref(null)
const followStatus = ref({ is_following: false, is_followed: false, mutual: false })
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
  if (isPlaying.value) audioRef.value.pause()
  else audioRef.value.play()
  isPlaying.value = !isPlaying.value
}
const onVoiceTimeUpdate = () => {
  if (!audioRef.value) return
  currentVoiceTime.value = audioRef.value.currentTime
  if (audioRef.value.duration) voiceProgress.value = (audioRef.value.currentTime / audioRef.value.duration) * 100
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
    if (!route.params.id) throw new Error('文章ID为空')
    const articleRes = await http.get(`/articles/${route.params.id}`)
    const articleData = articleRes.data.article || articleRes.data
    if (!articleData) throw new Error('文章数据为空')
    article.value = articleData
    comments.value = articleRes.data.comments || []
    liked.value = articleRes.data.liked || false
    commentLiked.value = articleRes.data.comment_liked || {}
    try {
      const siteConfigRes = await http.get('/site-config')
      siteTitle.value = siteConfigRes.data.site_title || '校园论坛'
    } catch (configError) { siteTitle.value = '校园论坛' }
    if (token.value && article.value.id) {
      try {
        const favoriteRes = await http.get(`/articles/${article.value.id}/favorite/check`)
        favorited.value = favoriteRes.data.favorited || false
      } catch (error) { favorited.value = false }
    }
    shareUrl.value = `${window.location.origin}/article/${article.value.id}`
    document.title = `${article.value.title} - ${siteTitle.value}`
    await nextTick()
    if (audioRef.value) audioRef.value.load()
    initVideoPlayers()
    loadFollowStatus()
  } catch (error) {
    const status = error.response?.status
    if (status === 404) {
      articleError.value = { code: 404, title: '资源未找到', message: '请求的文章不存在', detail: error.response?.data?.error || '该文章可能已被删除或ID错误' }
    } else {
      articleError.value = { code: status || 500, title: '加载失败', message: '加载文章时发生错误', detail: error.response?.data?.error || error.message }
    }
  }
}

const initVideoPlayers = () => {
  if (!contentRef.value) return
  contentRef.value.querySelectorAll('video').forEach((videoEl) => {
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
  if (!token.value) { router.push('/login'); return }
  const isLiked = liked.value
  try {
    if (isLiked) { await http.delete(`/articles/${article.value.id}/like`); article.value.like_count--; liked.value = false }
    else { await http.post(`/articles/${article.value.id}/like`); article.value.like_count++; liked.value = true }
  } catch (error) { liked.value = isLiked }
}

const toggleFavorite = async () => {
  if (!token.value) { router.push('/login'); return }
  const isFavorited = favorited.value
  try {
    if (isFavorited) { await http.delete(`/articles/${article.value.id}/favorite`); article.value.favorite_count--; favorited.value = false }
    else { await http.post(`/articles/${article.value.id}/favorite`); article.value.favorite_count++; favorited.value = true }
  } catch (error) { favorited.value = isFavorited }
}

const loadFollowStatus = async () => {
  if (!article.value || !token.value) return
  try {
    const response = await http.get(`/friends/status/${article.value.user_id}`)
    followStatus.value = { is_following: response.data.is_friend, is_followed: response.data.is_friend, mutual: response.data.is_friend }
  } catch (error) { /* ignore */ }
}

const handleFollow = async () => {
  if (!token.value) { router.push('/login'); return }
  try {
    if (followStatus.value.is_following) {
      await http.delete(`/friends/${article.value.user_id}`)
      followStatus.value.is_following = false
      followStatus.value.mutual = false
    } else {
      await http.post('/friends/request', { friend_id: article.value.user_id })
      await showSuccess('已发送好友请求')
    }
  } catch (error) { /* ignore */ }
}

const goToUserProfile = (userId) => { if (userId) router.push(`/profile?id=${userId}`) }

const coinArticle = async () => {
  if (!token.value) { router.push('/login'); return }
  try {
    await http.post(`/articles/${article.value.id}/coin`)
    article.value.coin_count = (article.value.coin_count || 0) + 1
    coined.value = true
    await showSuccess('投币成功')
  } catch (error) { /* ignore */ }
}

const toggleCommentLike = async (comment) => {
  if (!token.value) { router.push('/login'); return }
  const isCommentLiked = commentLiked.value[comment.id]
  try {
    if (isCommentLiked) { await http.delete(`/comments/${comment.id}/like`); comment.like_count--; commentLiked.value[comment.id] = false }
    else { await http.post(`/comments/${comment.id}/like`); comment.like_count++; commentLiked.value[comment.id] = true }
  } catch (error) { commentLiked.value[comment.id] = isCommentLiked }
}

const showReplyForm = (commentId) => { replyingTo.value = commentId; replyContent.value = '' }
const cancelReply = () => { replyingTo.value = null; replyContent.value = ''; replyIsAnonymous.value = false }

const submitComment = async () => {
  if (!commentContent.value.trim()) return
  try {
    await http.post(`/articles/${article.value.id}/comments`, { content: commentContent.value, is_anonymous: commentIsAnonymous.value })
    await loadArticle()
    commentContent.value = ''
    commentIsAnonymous.value = false
  } catch (error) { /* ignore */ }
}

const submitReply = async (parentId) => {
  if (!replyContent.value.trim()) return
  try {
    await http.post(`/articles/${article.value.id}/comments`, { content: replyContent.value, parent_id: parentId, is_anonymous: replyIsAnonymous.value })
    await loadArticle()
    cancelReply()
  } catch (error) { /* ignore */ }
}

const handleNestedReply = async ({ parentId, content, isAnonymous }) => {
  if (!content.trim()) return
  try {
    await http.post(`/articles/${article.value.id}/comments`, { content, parent_id: parentId, is_anonymous: isAnonymous })
    await loadArticle()
    cancelReply()
  } catch (error) { /* ignore */ }
}

const deleteArticle = async () => {
  try {
    const confirmed = await showConfirm('确定要删除这篇文章吗？', { title: '确认删除' }).catch(() => false)
    if (!confirmed) return
    if (currentUser.value?.role !== 'admin') {
      const reason = await showPrompt('请输入删除原因（管理员将审核）：', { title: '删除原因' }).catch(() => null)
      if (!reason) return
      try {
        await http.delete(`/articles/${article.value.id}`, { data: { reason } })
        await showSuccess('删除申请已提交，等待管理员审核')
        router.push('/')
      } catch (error) { /* ignore */ }
    } else {
      try {
        await http.delete(`/articles/${article.value.id}`)
        router.push('/')
      } catch (error) { /* ignore */ }
    }
  } catch (error) { /* ignore */ }
}

const deleteComment = async (commentId, comment) => {
  if (!canDeleteComment(comment)) return
  try {
    const confirmed = await showConfirm('确定要删除这条评论吗？', { title: '确认删除' }).catch(() => false)
    if (!confirmed) return
    try {
      await http.delete(`/comments/${commentId}`)
      for (let i = 0; i < comments.value.length; i++) {
        if (comments.value[i].id === commentId) { comments.value.splice(i, 1); break }
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
    } catch (error) { /* ignore */ }
  } catch (error) { /* ignore */ }
}

const canDeleteComment = (comment) => {
  if (!currentUser.value) return false
  const currentUserId = Number(currentUser.value.id)
  const commentUserId = Number(comment.user_id)
  return currentUserId === commentUserId || currentUser.value.role === 'admin'
}

const handleContentClick = (event) => {
  const target = event.target
  if (target.tagName === 'IMG') { currentImageUrl.value = target.src; showImageViewer.value = true }
}

const copyShareUrl = async () => {
  try {
    await navigator.clipboard.writeText(shareUrl.value)
    copySuccess.value = true
    setTimeout(() => { copySuccess.value = false }, 2000)
  } catch (error) { /* ignore */ }
}

const submitReport = async () => {
  if (!reportReason.value) { showSuccess('请选择举报原因'); return }
  if (!reportDescription.value.trim()) { showSuccess('请填写详细说明'); return }
  submittingReport.value = true
  try {
    await http.post('/reports', { target_type: 'article', target_id: article.value.id, reason: reportReason.value, description: reportDescription.value.trim() })
    showSuccess('举报已提交，感谢您的反馈')
    closeReportDialog()
  } catch (error) { /* ignore */ } finally { submittingReport.value = false }
}

const closeReportDialog = () => {
  showReportDialog.value = false
  reportReason.value = ''
  reportDescription.value = ''
}

const formatDate = (date) => new Date(date).toLocaleString('zh-CN')

onMounted(() => {
  const user = localStorage.getItem('user')
  if (user) currentUser.value = JSON.parse(user)
  loadArticle()
})
onBeforeUnmount(() => { document.title = siteTitle.value })
</script>

<style scoped>
.article-layout {
  display: grid;
  grid-template-columns: 1fr;
  gap: 16px;
}
.article-header {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}
.header-info { flex: 1; min-width: 0; }
.article-title { font-size: 22px; font-weight: 700; margin: 4px 0 8px; }
.article-meta { display: flex; align-items: center; gap: 12px; flex-wrap: wrap; font-size: 13px; }
.interact-row { display: flex; flex-wrap: wrap; gap: 8px; }
.spacer { flex: 1; }
.voice-row { display: flex; gap: 12px; align-items: center; }
.voice-body { flex: 1; }
.voice-label { display: flex; align-items: center; gap: 6px; font-size: 13px; margin-bottom: 6px; }
.comment-card { margin-top: 16px; }
.comment-input { display: flex; gap: 12px; margin-bottom: 16px; }
.comment-input-body { flex: 1; }
.comment-tools { display: flex; justify-content: space-between; align-items: center; margin-top: 8px; }
.comment-login { padding: 24px 0; }
.comment-item { display: flex; gap: 12px; padding: 12px 0; border-bottom: 1px solid var(--campus-border); }
.comment-body { flex: 1; min-width: 0; }
.comment-meta { display: flex; gap: 8px; align-items: center; margin-bottom: 4px; }
.comment-name { font-weight: 600; }
.comment-content { white-space: pre-wrap; word-break: break-word; margin-bottom: 6px; }
.comment-actions { display: flex; gap: 4px; }
.reply-form { margin-top: 10px; }
.reply-tools { display: flex; align-items: center; gap: 8px; margin-top: 8px; }
.text-center { text-align: center; }
@media (max-width: 768px) {
  .article-title { font-size: 18px; }
}
</style>
