<template>
  <div v-if="article">
    <v-card class="pa-6 mb-4">
      <v-card-title class="text-h4 mb-2">{{ article.title }}</v-card-title>
      
      <!-- 作者信息 -->
      <div class="d-flex align-center gap-4 mb-4">
        <v-avatar
          size="48"
          color="primary"
          class="cursor-pointer"
          @click="goToUserProfile(article.user.id)"
        >
          <v-img :src="article.user.avatar"></v-img>
        </v-avatar>
        <div class="flex-grow-1">
          <div
            class="font-weight-bold text-body-1 cursor-pointer"
            @click="goToUserProfile(article.user.id)"
          >
            {{ article.user.display_name }}
          </div>
          <div v-if="article.user.signature" class="text-caption text-medium-emphasis">
            {{ article.user.signature }}
          </div>
        </div>
        <v-btn
          v-if="token && currentUser && currentUser.id !== article.user_id"
          variant="outlined"
          size="small"
          :color="followStatus.is_following ? 'default' : 'primary'"
          @click="handleFollow"
        >
          {{ followStatus.is_following ? '已关注' : followStatus.is_followed ? '回关' : '关注' }}
        </v-btn>
      </div>
      
      <v-card-subtitle class="d-flex flex-wrap gap-3 mb-4">
        <v-chip size="small" color="secondary" variant="tonal">
          <v-icon start size="small">mdi-folder</v-icon>
          {{ article.category.name }}
        </v-chip>
        <v-chip size="small" variant="text">
          <v-icon start size="small">mdi-clock</v-icon>
          {{ formatDate(article.created_at) }}
        </v-chip>
        <v-chip size="small" variant="text">
          <v-icon start size="small">mdi-eye</v-icon>
          {{ article.view_count }}
        </v-chip>
      </v-card-subtitle>
      
      <v-card-actions v-if="canEdit">
        <v-btn variant="outlined" color="primary" :to="'/create?id=' + article.id">
          <v-icon start>mdi-pencil</v-icon>
          编辑
        </v-btn>
        <v-btn variant="outlined" color="error" @click="deleteArticle">
          <v-icon start>mdi-delete</v-icon>
          删除
        </v-btn>
      </v-card-actions>
    </v-card>
    
    <v-card class="pa-6 mb-4">
      <div ref="contentRef" class="article-content markdown-body" @click="handleContentClick" v-html="renderedHtml"></div>
    </v-card>
    
    <v-card class="pa-4 mb-4">
      <v-btn
        @click="toggleLike"
        :color="liked ? 'primary' : 'default'"
        :variant="liked ? 'flat' : 'outlined'"
      >
        <v-icon start>mdi-thumb-up</v-icon>
        {{ article.like_count }} 点赞
      </v-btn>
    </v-card>
    
    <v-card class="pa-6">
      <v-card-title class="text-h6 mb-4">
        评论 ({{ comments.length }})
      </v-card-title>
      
      <div v-if="token" class="mb-6">
        <v-textarea
          v-model="commentContent"
          label="写下你的评论..."
          variant="outlined"
          rows="3"
          hide-details
          class="mb-3"
        ></v-textarea>
        <v-btn color="primary" @click="submitComment">发表评论</v-btn>
      </div>
      
      <div v-for="comment in comments" :key="comment.id" class="comment-item mb-4">
        <div class="d-flex gap-3">
          <v-avatar color="primary" size="40" class="cursor-pointer" @click="goToUserProfile(comment.user.id)">
            <v-img :src="comment.user.avatar"></v-img>
          </v-avatar>
          <div class="flex-grow-1">
            <div class="d-flex align-center gap-2">
              <span class="font-weight-bold cursor-pointer" @click="goToUserProfile(comment.user.id)">{{ comment.user.display_name }}</span>
              <span class="text-caption text-medium-emphasis">{{ formatDate(comment.created_at) }}</span>
            </div>
            <div class="mt-1">{{ comment.content }}</div>
            <div class="d-flex gap-4 mt-2">
              <v-btn variant="text" size="small" @click="toggleCommentLike(comment)" :color="commentLiked[comment.id] ? 'primary' : 'default'">
                <v-icon start size="small">mdi-thumb-up</v-icon>
                {{ comment.like_count }} 点赞
              </v-btn>
              <v-btn variant="text" size="small" @click="showReplyForm(comment.id)" v-if="token">
                <v-icon start size="small">mdi-reply</v-icon>
                回复 ({{ comment.reply_count }})
              </v-btn>
              <v-btn variant="text" size="small" color="error" @click="deleteComment(comment.id)" v-if="canDeleteComment(comment)">
                <v-icon start size="small">mdi-delete</v-icon>
                删除
              </v-btn>
            </div>
            
            <!-- 回复表单 -->
            <div v-if="replyingTo === comment.id" class="mt-3">
              <v-textarea
                v-model="replyContent"
                :placeholder="'回复 ' + comment.user.display_name + '...'"
                variant="outlined"
                rows="2"
                hide-details
                class="mb-2"
              ></v-textarea>
              <div class="d-flex gap-2">
                <v-btn size="small" color="primary" @click="submitReply(comment.id)">发送</v-btn>
                <v-btn size="small" variant="text" @click="cancelReply">取消</v-btn>
              </div>
            </div>
            
            <!-- 回复列表 -->
            <div v-if="comment.replies && comment.replies.length > 0" class="replies-list mt-3 ml-4 pl-3 border-left">
              <div v-for="reply in comment.replies" :key="reply.id" class="reply-item mb-3">
                <div class="d-flex gap-2">
                  <v-avatar color="primary" size="32" class="cursor-pointer" @click="goToUserProfile(reply.user.id)">
                    <v-img :src="reply.user.avatar"></v-img>
                  </v-avatar>
                  <div class="flex-grow-1">
                    <div class="d-flex align-center gap-2">
                      <span class="font-weight-bold text-body-2 cursor-pointer" @click="goToUserProfile(reply.user.id)">{{ reply.user.display_name }}</span>
                      <span class="text-caption text-medium-emphasis">{{ formatDate(reply.created_at) }}</span>
                    </div>
                    <div class="mt-1 text-body-2">{{ reply.content }}</div>
                    <div class="d-flex gap-3 mt-1">
                      <v-btn variant="text" size="small" @click="toggleCommentLike(reply)" :color="commentLiked[reply.id] ? 'primary' : 'default'" density="compact">
                        <v-icon start size="x-small">mdi-thumb-up</v-icon>
                        {{ reply.like_count }}
                      </v-btn>
                      <v-btn variant="text" size="small" color="error" @click="deleteComment(reply.id)" v-if="canDeleteComment(reply)" density="compact">
                        <v-icon start size="x-small">mdi-delete</v-icon>
                        删除
                      </v-btn>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </v-card>
    
    <ImageViewer v-if="showImageViewer" :url="currentImageUrl" @close="closeImageViewer" />
  </div>
  
  <div v-else class="d-flex justify-center align-center" style="min-height: 50vh;">
    <v-progress-circular indeterminate color="primary" size="64"></v-progress-circular>
  </div>
</template>

<script>
import { ref, onMounted, onBeforeUnmount, computed, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '../api'
import ImageViewer from '../components/ImageViewer.vue'
import { confirm as showConfirm, prompt as showPrompt, success as showSuccess } from '../utils/modal'
import MarkdownIt from 'markdown-it'
import videojs from 'video.js'
import 'video.js/dist/video-js.css'

export default {
  name: 'Article',
  components: {
    ImageViewer
  },
  setup() {
    const route = useRoute()
    const router = useRouter()
    const article = ref(null)
    const comments = ref([])
    const liked = ref(false)
    const commentContent = ref('')
    const commentLiked = ref({})
    const replyingTo = ref(null)
    const replyContent = ref('')
    const showImageViewer = ref(false)
    const currentImageUrl = ref('')
    const token = ref(localStorage.getItem('token'))
    const currentUser = ref(null)
    const contentRef = ref(null)
    const renderedHtml = ref('')
    const followStatus = ref({
      is_following: false,
      is_followed: false,
      mutual: false
    })
    
    // 存储 Video.js 实例
    const videoPlayers = ref([])
    
    // 初始化 markdown-it
    const md = new MarkdownIt({
      html: true,
      linkify: true,
      typographer: true
    })
    
    const canEdit = computed(() => {
      if (!currentUser.value || !article.value) return false
      return currentUser.value.id === article.value.user_id || currentUser.value.role === 'admin'
    })
    
    const loadArticle = async () => {
      try {
        const response = await api.get(`/articles/${route.params.id}`)
        article.value = response.data.article
        comments.value = response.data.comments
        liked.value = response.data.liked || false
        commentLiked.value = response.data.comment_liked || {}
        
        // 只使用 article.content 中的数据，不使用 content_html
        renderedHtml.value = md.render(article.value.content)
        
        await nextTick()
        
        // 初始化 Video.js 播放器
        initVideoPlayers()
        
        // 加载关注状态
        loadFollowStatus()
      } catch (error) {
        console.error('加载文章失败', error)
        router.push('/')
      }
    }
    
    // 初始化所有视频为 Video.js 播放器
    const initVideoPlayers = () => {
      if (!contentRef.value) return
      
      const videoElements = contentRef.value.querySelectorAll('video')
      console.log('找到视频元素:', videoElements.length)
      
      videoElements.forEach((videoEl, index) => {
        // 获取视频属性
        const src = videoEl.src || videoEl.getAttribute('src')
        const poster = videoEl.poster || videoEl.getAttribute('poster') || ''
        
        if (!src) return
        
        // 创建包装器
        const wrapper = document.createElement('div')
        wrapper.className = 'video-js-container'
        wrapper.style.cssText = 'width: 100%; max-width: 800px; margin: 16px auto; border-radius: 8px; overflow: hidden;'
        
        // 将原始 video 元素移入包装器
        videoEl.parentNode.insertBefore(wrapper, videoEl)
        wrapper.appendChild(videoEl)
        
        // 移除 controls 属性，让 Video.js 自己控制
        videoEl.removeAttribute('controls')
        videoEl.className = 'video-js vjs-big-play-centered'
        videoEl.style.cssText = 'width: 100%; height: auto; max-height: none;'
        
        // 获取视频原始尺寸
        const originalWidth = parseInt(videoEl.getAttribute('width')) || 1280
        const originalHeight = parseInt(videoEl.getAttribute('height')) || 720
        const aspectRatio = originalHeight / originalWidth
        
        // 初始化 Video.js
        const player = videojs(videoEl, {
          controls: true,
          autoplay: false,
          preload: 'auto',
          fluid: true,
          responsive: true,
          maintainAspectRatio: true,
          playbackRates: [0.5, 0.75, 1, 1.25, 1.5, 2],
          poster: poster,
          sources: [{
            src: src,
            type: 'video/mp4'
          }],
          controlBar: {
            children: [
              'playToggle',
              'volumePanel',
              'currentTimeDisplay',
              'timeDivider',
              'durationDisplay',
              'progressControl',
              'playbackRateMenuButton',
              'fullscreenToggle'
            ]
          }
        })
        
        // 设置视频容器的宽高比
        const playerContainer = player.el().parentElement
        playerContainer.style.aspectRatio = `${originalWidth} / ${originalHeight}`
        
        videoPlayers.value.push(player)
        console.log(`初始化视频播放器 ${index + 1}, 宽高比: ${originalWidth}x${originalHeight}`)
      })
    }
    
    const toggleLike = async () => {
      if (!token.value) {
        router.push('/login')
        return
      }
      
      try {
        if (liked.value) {
          await api.delete(`/articles/${article.value.id}/like`)
          article.value.like_count--
          liked.value = false
        } else {
          await api.post(`/articles/${article.value.id}/like`)
          article.value.like_count++
          liked.value = true
        }
      } catch (error) {
        console.error('点赞失败', error)
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
      
      try {
        if (commentLiked.value[comment.id]) {
          await api.delete(`/comments/${comment.id}/like`)
          comment.like_count--
          commentLiked.value[comment.id] = false
        } else {
          await api.post(`/comments/${comment.id}/like`)
          comment.like_count++
          commentLiked.value[comment.id] = true
        }
      } catch (error) {
        console.error('评论点赞失败', error)
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
    
    const deleteComment = async (commentId) => {
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
      return currentUser.value.id === comment.user_id || currentUser.value.role === 'admin'
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
    
    // 组件卸载时销毁所有 Video.js 实例
    onBeforeUnmount(() => {
      videoPlayers.value.forEach(player => {
        if (player && !player.isDisposed()) {
          player.dispose()
        }
      })
      videoPlayers.value = []
    })
    
    return {
      article,
      comments,
      liked,
      commentContent,
      commentLiked,
      replyingTo,
      replyContent,
      token,
      currentUser,
      canEdit,
      showImageViewer,
      currentImageUrl,
      contentRef,
      renderedHtml,
      followStatus,
      toggleLike,
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
.article-content {
  min-height: 300px;
  line-height: 1.8;
  font-size: 16px;
  color: #333;
  word-wrap: break-word;
  overflow-wrap: break-word;
}

/* 标题样式 */
.article-content :deep(h1) {
  font-size: 2em;
  font-weight: bold;
  margin: 1.5em 0 1em;
  padding-bottom: 0.3em;
  border-bottom: 2px solid #e0e0e0;
  color: #1976d2;
}

.article-content :deep(h2) {
  font-size: 1.6em;
  font-weight: bold;
  margin: 1.3em 0 0.8em;
  padding-bottom: 0.2em;
  border-bottom: 1px solid #e0e0e0;
  color: #1976d2;
}

.article-content :deep(h3) {
  font-size: 1.4em;
  font-weight: bold;
  margin: 1.2em 0 0.6em;
  color: #1976d2;
}

.article-content :deep(h4) {
  font-size: 1.2em;
  font-weight: bold;
  margin: 1em 0 0.5em;
  color: #1976d2;
}

.article-content :deep(h5) {
  font-size: 1.1em;
  font-weight: bold;
  margin: 1em 0 0.5em;
  color: #1976d2;
}

.article-content :deep(h6) {
  font-size: 1em;
  font-weight: bold;
  margin: 1em 0 0.5em;
  color: #1976d2;
}

/* 段落样式 */
.article-content :deep(p) {
  margin: 1em 0;
  line-height: 1.8;
}

/* 列表样式 */
.article-content :deep(ul),
.article-content :deep(ol) {
  margin: 1em 0;
  padding-left: 2em;
}

.article-content :deep(ul) {
  list-style-type: disc;
}

.article-content :deep(ol) {
  list-style-type: decimal;
}

.article-content :deep(li) {
  margin: 0.5em 0;
  line-height: 1.6;
}

.article-content :deep(ul ul),
.article-content :deep(ol ol),
.article-content :deep(ul ol),
.article-content :deep(ol ul) {
  margin: 0.5em 0;
}

/* 引用样式 */
.article-content :deep(blockquote) {
  margin: 1.5em 0;
  padding: 1em 1.5em;
  border-left: 4px solid #1976d2;
  background-color: #f5f5f5;
  border-radius: 4px;
  color: #666;
}

.article-content :deep(blockquote p) {
  margin: 0.5em 0;
}

/* 代码块样式 */
.article-content :deep(pre) {
  margin: 1.5em 0;
  padding: 1em;
  background-color: #2d2d2d;
  border-radius: 8px;
  overflow-x: auto;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.article-content :deep(pre code) {
  color: #f8f8f2;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.6;
  background: none;
  padding: 0;
}

/* 行内代码样式 */
.article-content :deep(code) {
  background-color: #f5f5f5;
  color: #e53935;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 0.9em;
}

/* 链接样式 */
.article-content :deep(a) {
  color: #1976d2;
  text-decoration: none;
  border-bottom: 1px solid transparent;
  transition: border-color 0.2s;
}

.article-content :deep(a:hover) {
  border-bottom-color: #1976d2;
}

/* 图片样式 */
.article-content :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 8px;
  margin: 1.5em auto;
  display: block;
  cursor: pointer;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.article-content :deep(img:hover) {
  transform: scale(1.02);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* 表格样式 */
.article-content :deep(table) {
  width: 100%;
  margin: 1.5em 0;
  border-collapse: collapse;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.article-content :deep(thead) {
  background-color: #1976d2;
  color: white;
}

.article-content :deep(th) {
  padding: 12px 16px;
  text-align: left;
  font-weight: bold;
}

.article-content :deep(td) {
  padding: 12px 16px;
  border-bottom: 1px solid #e0e0e0;
}

.article-content :deep(tbody tr:hover) {
  background-color: #f5f5f5;
}

.article-content :deep(tbody tr:last-child td) {
  border-bottom: none;
}

/* 分隔线样式 */
.article-content :deep(hr) {
  margin: 2em 0;
  border: none;
  height: 2px;
  background: linear-gradient(to right, transparent, #e0e0e0, transparent);
}

/* 删除线样式 */
.article-content :deep(del) {
  color: #999;
  text-decoration: line-through;
}

/* 强调样式 */
.article-content :deep(strong) {
  font-weight: bold;
  color: #1976d2;
}

.article-content :deep(em) {
  font-style: italic;
  color: #666;
}

/* 视频容器样式 */
.video-js-container {
  width: 100%;
  max-width: 800px;
  margin: 16px auto;
  border-radius: 8px;
  overflow: hidden;
}

/* 评论样式 */
.comment-item {
  padding-bottom: 16px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.08);
}

.comment-item:last-child {
  border-bottom: none;
}

.replies-list {
  border-left: 2px solid rgba(0, 0, 0, 0.12);
}

.reply-item:last-child {
  margin-bottom: 0 !important;
}

.cursor-pointer {
  cursor: pointer;
}

/* 移动端适配 */
@media (max-width: 600px) {
  .article-content {
    font-size: 15px;
    line-height: 1.7;
  }
  
  .article-content :deep(h1) {
    font-size: 1.6em;
  }
  
  .article-content :deep(h2) {
    font-size: 1.4em;
  }
  
  .article-content :deep(h3) {
    font-size: 1.2em;
  }
  
  .article-content :deep(pre) {
    padding: 0.8em;
    font-size: 13px;
  }
  
  .article-content :deep(table) {
    font-size: 14px;
  }
  
  .article-content :deep(th),
  .article-content :deep(td) {
    padding: 8px 12px;
  }
}
</style>
