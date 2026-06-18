<template>
  <v-container v-if="articleError" class="min-h-screen">
    <ErrorPage
      :code="articleError.code"
      :title="articleError.title"
      :message="articleError.message"
      :detail="articleError.detail"
      :icon="articleError.code === 404 ? 'mdi-file-question' : 'mdi-alert-circle'"
      :actions="[
        { text: '返回首页', icon: 'mdi-home', color: 'primary', callback: () => router.push('/') },
        { text: '刷新重试', icon: 'mdi-refresh', variant: 'outlined', callback: () => { articleError.value = null; loadArticle() } }
      ]"
    />
  </v-container>
  
  <v-row v-else-if="article">
    <!-- 主内容区域 -->
    <v-col cols="12" md="9">
      <v-card class="article-card mb-4" variant="flat">
        <!-- 文章头部 -->
        <v-card-item>
          <template v-slot:prepend>
            <v-btn icon="mdi-arrow-left" variant="text" @click="router.back()"></v-btn>
          </template>
          <v-card-title class="text-h5 font-weight-bold">
            {{ article.title }}
          </v-card-title>
          <v-card-subtitle class="d-flex align-center flex-wrap gap-2 mt-2">
            <UserAvatar
              :user="article.user"
              :size="36"
              :showUsername="true"
              class="cursor-pointer"
              @click="goToUserProfile(article.user.id)"
            />
            <span class="text-caption text-medium-emphasis">
              {{ formatDate(article.created_at) }}
            </span>
            <v-chip size="small" color="primary" variant="flat" prepend-icon="mdi-folder">
              {{ article.category.name }}
            </v-chip>
            <span class="text-caption text-medium-emphasis">
              <v-icon size="small" class="mr-1">mdi-eye</v-icon>
              {{ article.view_count }} 阅读
            </span>
          </v-card-subtitle>
        </v-card-item>

        <v-divider></v-divider>

        <!-- 文章内容 -->
        <v-card-text class="pa-4">
          <div ref="contentRef" @click="handleContentClick">
            <MarkdownViewer :value="article.content" />
          </div>
        </v-card-text>

        <!-- 语音播放器 -->
        <v-card-actions v-if="article.voice_url" class="px-4 pb-4">
          <v-expand-transition>
            <v-sheet class="voice-player pa-3 rounded-lg" color="grey-lighten-4">
              <div class="d-flex align-center gap-3">
                <v-btn
                  :icon="isPlaying ? 'mdi-pause' : 'mdi-play'"
                  variant="flat"
                  color="primary"
                  size="large"
                  @click="toggleVoicePlay"
                />
                <div class="flex-grow-1">
                  <div class="d-flex align-center justify-space-between mb-1">
                    <span class="text-caption">
                      <v-icon size="small" class="mr-1">mdi-volume-high</v-icon>
                      语音朗读
                    </span>
                    <span class="text-caption text-medium-emphasis">
                      {{ formatVoiceTime(currentVoiceTime) }} / {{ formatVoiceTime(voiceDuration) }}
                    </span>
                  </div>
                  <v-progress-linear
                    v-model="voiceProgress"
                    color="primary"
                    height="4"
                    rounded
                  />
                </div>
              </div>
              <audio
                ref="audioRef"
                :src="article.voice_url"
                @timeupdate="onVoiceTimeUpdate"
                @loadedmetadata="onVoiceLoaded"
                @ended="onVoiceEnded"
              />
            </v-sheet>
          </v-expand-transition>
        </v-card-actions>

        <v-divider></v-divider>

        <!-- 互动栏 -->
        <v-card-actions class="pa-4">
          <div class="d-flex flex-wrap gap-2">
            <v-btn
              @click="toggleLike"
              :color="liked ? 'primary' : 'default'"
              :variant="liked ? 'flat' : 'outlined'"
              prepend-icon="mdi-thumb-up"
            >
              {{ article.like_count }} 点赞
            </v-btn>
            <v-btn
              @click="coinArticle"
              :color="coined ? 'amber' : 'default'"
              :variant="coined ? 'flat' : 'outlined'"
              prepend-icon="mdi-coins"
              :disabled="!token"
            >
              {{ article.coin_count || 0 }} 投币
            </v-btn>
            <v-btn
              @click="toggleFavorite"
              :color="favorited ? 'primary' : 'default'"
              :variant="favorited ? 'flat' : 'outlined'"
              prepend-icon="mdi-bookmark"
            >
              {{ article.favorite_count || 0 }} 收藏
            </v-btn>
            <v-btn
              @click="showShareDialog = true"
              variant="outlined"
              prepend-icon="mdi-share-variant"
            >
              分享
            </v-btn>
            <v-btn
              @click="showReportDialog = true"
              variant="outlined"
              color="error"
              prepend-icon="mdi-flag"
              v-if="token && currentUser && currentUser.id !== article.user_id"
            >
              举报
            </v-btn>
          </div>
          <v-spacer></v-spacer>
          <div v-if="canEdit" class="d-flex gap-2">
            <v-btn
              variant="outlined"
              color="primary"
              size="small"
              :to="'/create?id=' + article.id"
              prepend-icon="mdi-pencil"
            >
              编辑
            </v-btn>
            <v-btn
              variant="outlined"
              color="error"
              size="small"
              @click="deleteArticle"
              prepend-icon="mdi-delete"
            >
              删除
            </v-btn>
          </div>
        </v-card-actions>
      </v-card>

      <!-- 评论区域 -->
      <v-card variant="flat">
        <v-card-item>
          <template v-slot:prepend>
            <v-icon>mdi-comment-text</v-icon>
          </template>
          <v-card-title>评论 ({{ comments.length }})</v-card-title>
        </v-card-item>

        <v-card-text v-if="token" class="pa-4">
          <div class="d-flex gap-3">
            <UserAvatar :user="currentUser" :size="40" :showUsername="false" />
            <div class="flex-grow-1">
              <v-textarea
                v-model="commentContent"
                placeholder="写下你的评论..."
                variant="outlined"
                rows="3"
                hide-details
              />
              <div class="d-flex align-center justify-space-between mt-2">
                <v-checkbox
                  v-model="commentIsAnonymous"
                  label="匿名评论"
                  color="primary"
                  hide-details
                  density="compact"
                />
                <v-btn color="primary" @click="submitComment" :disabled="!commentContent.trim()">
                  发表
                </v-btn>
              </div>
            </div>
          </div>
        </v-card-text>

        <v-card-text v-else class="pa-4 text-center">
          <span class="text-medium-emphasis">登录后参与评论</span>
          <v-btn variant="text" color="primary" size="small" @click="router.push('/login')" class="ml-2">
            登录
          </v-btn>
        </v-card-text>

        <v-divider></v-divider>

        <!-- 评论列表 -->
        <v-list class="pa-2">
          <div v-for="comment in comments" :key="comment.id" class="comment-item mb-3">
            <v-list-item class="px-3">
              <template v-slot:prepend>
                <UserAvatar
                  :user="comment.user"
                  :size="40"
                  :showUsername="false"
                  class="cursor-pointer"
                  @click="goToUserProfile(comment.user.id)"
                />
              </template>

              <v-list-item-title class="d-flex align-center gap-2 mb-1">
                <span class="font-weight-medium">
                  {{ comment.user.display_name || comment.user.username || '匿名用户' }}
                </span>
                <span class="text-caption text-medium-emphasis">
                  {{ formatDate(comment.created_at) }}
                </span>
              </v-list-item-title>

              <v-list-item-subtitle class="comment-text">
                {{ comment.content }}
              </v-list-item-subtitle>

              <template v-slot:append>
                <div class="d-flex flex-column align-end gap-1">
                  <div class="d-flex gap-1">
                    <v-btn
                      variant="text"
                      size="x-small"
                      @click="toggleCommentLike(comment)"
                      :color="commentLiked[comment.id] ? 'primary' : 'default'"
                    >
                      <v-icon size="14">mdi-thumb-up</v-icon>
                      {{ comment.like_count }}
                    </v-btn>
                    <v-btn
                      variant="text"
                      size="x-small"
                      @click="showReplyForm(comment.id)"
                      v-if="token"
                    >
                      <v-icon size="14">mdi-reply</v-icon>
                      回复
                    </v-btn>
                    <v-btn
                      variant="text"
                      size="x-small"
                      color="error"
                      @click="deleteComment(comment.id, comment)"
                      v-if="canDeleteComment(comment)"
                    >
                      <v-icon size="14">mdi-delete</v-icon>
                    </v-btn>
                  </div>

                  <!-- 回复表单 -->
                  <v-expand-transition>
                    <div v-if="replyingTo === comment.id" class="reply-form mt-2">
                      <v-textarea
                        v-model="replyContent"
                        :placeholder="'回复 ' + (comment.user.display_name || comment.user.username || '匿名用户')"
                        variant="outlined"
                        rows="2"
                        density="compact"
                        hide-details
                      />
                      <div class="d-flex align-center gap-2 mt-2">
                        <v-checkbox
                          v-model="replyIsAnonymous"
                          label="匿名"
                          color="primary"
                          hide-details
                          density="compact"
                        />
                        <v-spacer></v-spacer>
                        <v-btn size="small" color="primary" @click="submitReply(comment.id)">发送</v-btn>
                        <v-btn size="small" variant="text" @click="cancelReply">取消</v-btn>
                      </div>
                    </div>
                  </v-expand-transition>
                </div>
              </template>
            </v-list-item>

            <!-- 回复列表 -->
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
        </v-list>
      </v-card>
    </v-col>

    <!-- 侧边栏 -->
    <v-col cols="12" md="3" class="d-none d-md-block">
      <v-card variant="flat" class="mb-4">
        <v-card-item>
          <v-card-title class="d-flex align-center gap-2">
            <UserAvatar
              :user="article.user"
              :size="48"
              :showUsername="false"
              class="cursor-pointer"
              @click="goToUserProfile(article.user.id)"
            />
            <div>
              <div class="text-subtitle-1 font-weight-medium">
                {{ article.user.display_name }}
              </div>
              <div class="text-caption text-medium-emphasis">
                {{ article.user.signature || '暂无签名' }}
              </div>
            </div>
          </v-card-title>
        </v-card-item>
        <v-card-actions v-if="token && currentUser && currentUser.id !== article.user_id">
          <v-btn
            variant="tonal"
            :color="followStatus.is_following ? 'default' : 'primary'"
            block
            @click="handleFollow"
          >
            <v-icon class="mr-1">{{ followStatus.is_following ? 'mdi-check' : 'mdi-plus' }}</v-icon>
            {{ followStatus.is_following ? '已关注' : followStatus.is_followed ? '回关' : '关注' }}
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-col>

    <!-- 分享对话框 -->
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

    <!-- 举报对话框 -->
    <v-dialog v-model="showReportDialog" max-width="500" persistent>
      <v-card>
        <v-card-title class="d-flex align-center bg-error text-white pa-4">
          <v-icon class="mr-2">mdi-flag-variant</v-icon>
          举报文章
        </v-card-title>
        <v-card-text class="pa-4">
          <v-alert
            type="info"
            variant="tonal"
            density="compact"
            class="mb-4"
          >
            感谢您对平台环境的维护。我们会认真审核每一条举报，并在3个工作日内处理。
          </v-alert>
          
          <v-form ref="reportForm">
            <div class="mb-4">
              <div class="text-subtitle-2 mb-2">请选择举报原因 <span class="text-error">*</span></div>
              <v-chip-group
                v-model="reportReason"
                column
                mandatory
              >
                <v-chip
                  v-for="reason in reportReasons"
                  :key="reason.value"
                  :value="reason.value"
                  filter
                  variant="outlined"
                  color="error"
                  size="small"
                >
                  <v-icon start size="16">{{ reason.icon }}</v-icon>
                  {{ reason.title }}
                </v-chip>
              </v-chip-group>
            </div>
            
            <v-textarea
              v-model="reportDescription"
              label="详细说明（必填）"
              variant="outlined"
              rows="4"
              :rules="[v => !!v || '请填写举报详细说明']"
              placeholder="请详细描述您举报的原因，包括具体内容和违规证据..."
              counter
              maxlength="500"
            />
          </v-form>
        </v-card-text>
        <v-card-actions class="pa-4 pt-0">
          <v-spacer />
          <v-btn variant="text" @click="closeReportDialog">取消</v-btn>
          <v-btn color="error" @click="submitReport" :loading="submittingReport" prepend-icon="mdi-send">
            提交举报
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- 图片查看器 -->
    <ImageViewer v-if="showImageViewer" :url="currentImageUrl" @close="closeImageViewer" />
  </v-row>

  <!-- 加载状态 -->
  <div v-else class="d-flex justify-center align-center" style="min-height: 60vh;">
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
import ErrorPage from '../components/ErrorPage.vue'
import CommentReply from '../components/CommentReply.vue'
import { confirm as showConfirm, prompt as showPrompt, success as showSuccess } from '../utils/modal'


export default {
  name: 'Article',
  components: {
    ImageViewer,
    UserAvatar,
    MarkdownViewer,
    CommentReply
  },
  setup() {
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
      { title: '垃圾广告', value: '垃圾广告', icon: 'mdi-trash-can' },
      { title: '色情低俗', value: '色情低俗', icon: 'mdi-alert' },
      { title: '暴力血腥', value: '暴力血腥', icon: 'mdi-knife' },
      { title: '政治敏感', value: '政治敏感', icon: 'mdi-account-alert' },
      { title: '违法犯罪', value: '违法犯罪', icon: 'mdi-gavel' },
      { title: '谣言虚假', value: '谣言虚假', icon: 'mdi-chat-alert' },
      { title: '侵犯隐私', value: '侵犯隐私', icon: 'mdi-eye-off' },
      { title: '其他违规', value: '其他违规', icon: 'mdi-help-circle' }
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
        console.log('尝试加载文章，ID:', route.params.id)
        console.log('当前路由:', route.fullPath)
        console.log('路由参数:', JSON.stringify(route.params))
        
        if (!route.params.id) {
          throw new Error('文章ID为空')
        }
        
        const articleRes = await api.get(`/articles/${route.params.id}`)
        console.log('API状态码:', articleRes.status)
        console.log('API返回数据:', JSON.stringify(articleRes.data, null, 2))
        
        // 尝试多种数据结构
        const articleData = articleRes.data.article || articleRes.data
        
        if (!articleData) {
          throw new Error('文章数据为空')
        }
        
        article.value = articleData
        comments.value = articleRes.data.comments || []
        liked.value = articleRes.data.liked || false
        commentLiked.value = articleRes.data.comment_liked || {}
        
        // 获取站点配置
        try {
          const siteConfigRes = await api.get('/site-config')
          siteTitle.value = siteConfigRes.data.site_title || '校园论坛'
        } catch (configError) {
          console.warn('加载站点配置失败', configError)
          siteTitle.value = '校园论坛'
        }
        
        // 检查收藏状态
        if (token.value && article.value.id) {
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
        if (audioRef.value) {
          audioRef.value.load()
        }
        
        initVideoPlayers()
        loadFollowStatus()
        
        console.log('文章加载成功')
      } catch (error) {
        console.error('加载文章失败', error)
        console.error('错误详情:', error.response?.data || error.message)
        
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
          // 删除好友
          await api.delete(`/friends/${article.value.user_id}`)
          followStatus.value.is_following = false
          followStatus.value.mutual = false
        } else {
          // 发送好友请求
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
      replyIsAnonymous.value = false
    }
    
    const submitComment = async () => {
      if (!commentContent.value.trim()) return

      try {
        const response = await api.post(`/articles/${article.value.id}/comments`, {
          content: commentContent.value,
          is_anonymous: commentIsAnonymous.value
        })
        
        // 重新加载文章和评论，确保数据一致性
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
        const response = await api.post(`/articles/${article.value.id}/comments`, {
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
          // 从列表中移除评论
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
      showReportDialog,
      reportReason,
      reportDescription,
      reportReasons,
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
      submitReport,
      closeReportDialog,
      formatDate,
      loadFollowStatus,
      handleFollow,
      goToUserProfile
    }
  }
}
</script>

<style scoped>
.article-card {
  border-radius: 12px;
}

.cursor-pointer {
  cursor: pointer;
}

.comment-text {
  white-space: pre-wrap;
  word-break: break-word;
}
</style>
