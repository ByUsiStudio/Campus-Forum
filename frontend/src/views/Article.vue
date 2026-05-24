<template>
  <div v-if="article">
    <v-card class="pa-6 mb-4">
      <v-card-title class="text-h4 mb-2">{{ article.title }}</v-card-title>
      
      <v-card-subtitle class="d-flex flex-wrap gap-3 mb-4">
        <v-chip size="small" color="primary" variant="tonal">
          <v-icon start size="small">mdi-account</v-icon>
          {{ article.user.display_name }}
        </v-chip>
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
      
      <v-list lines="two">
        <v-list-item v-for="comment in comments" :key="comment.id" class="px-0">
          <template v-slot:prepend>
            <v-avatar color="primary" size="40">
              <v-img :src="comment.user.avatar"></v-img>
            </v-avatar>
          </template>
          
          <v-list-item-title>{{ comment.user.display_name }}</v-list-item-title>
          <v-list-item-subtitle>{{ formatDate(comment.created_at) }}</v-list-item-subtitle>
          
          <v-list-item-content class="mt-2">
            {{ comment.content }}
          </v-list-item-content>
          
          <template v-slot:append v-if="canDeleteComment(comment)">
            <v-btn
              variant="text"
              color="error"
              size="small"
              @click="deleteComment(comment.id)"
            >
              删除
            </v-btn>
          </template>
        </v-list-item>
      </v-list>
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
    const showImageViewer = ref(false)
    const currentImageUrl = ref('')
    const token = ref(localStorage.getItem('token'))
    const currentUser = ref(null)
    const contentRef = ref(null)
    const renderedHtml = ref('')
    
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
        
        // 只使用 article.content 中的数据，不使用 content_html
        renderedHtml.value = md.render(article.value.content)
        
        await nextTick()
        
        // 初始化 Video.js 播放器
        initVideoPlayers()
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
        videoEl.style.cssText = 'width: 100%; height: auto;'
        
        // 初始化 Video.js
        const player = videojs(videoEl, {
          controls: true,
          autoplay: false,
          preload: 'auto',
          fluid: true,
          playbackRates: [0.5, 0.75, 1, 1.25, 1.5, 2],
          poster: poster,
          sources: [{
            src: src,
            type: 'video/mp4'
          }]
        })
        
        videoPlayers.value.push(player)
        console.log(`初始化视频播放器 ${index + 1}`)
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
    
    const submitComment = async () => {
      if (!commentContent.value.trim()) return
      
      try {
        const response = await api.post(`/articles/${article.value.id}/comments`, {
          content: commentContent.value
        })
        comments.value.unshift(response.data.comment)
        commentContent.value = ''
      } catch (error) {
        console.error('评论失败', error)
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
          comments.value = comments.value.filter(c => c.id !== commentId)
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
      token,
      currentUser,
      canEdit,
      showImageViewer,
      currentImageUrl,
      contentRef,
      renderedHtml,
      toggleLike,
      submitComment,
      deleteArticle,
      deleteComment,
      canDeleteComment,
      handleContentClick,
      closeImageViewer,
      formatDate
    }
  }
}
</script>

<style scoped>
.article-content {
  min-height: 300px;
  line-height: 1.8;
}

.video-js-container {
  width: 100%;
  max-width: 800px;
  margin: 16px auto;
  border-radius: 8px;
  overflow: hidden;
}

.article-content :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 8px;
  margin: 16px 0;
  cursor: pointer;
  transition: transform 0.2s ease;
}

.article-content :deep(img:hover) {
  transform: scale(1.02);
}

.article-content :deep(h1),
.article-content :deep(h2),
.article-content :deep(h3),
.article-content :deep(h4),
.article-content :deep(h5),
.article-content :deep(h6) {
  margin-top: 24px;
  margin-bottom: 16px;
  font-weight: 600;
}

.article-content :deep(p) {
  margin-bottom: 16px;
}

.article-content :deep(a) {
  color: rgb(var(--v-theme-primary));
  text-decoration: none;
}

.article-content :deep(a:hover) {
  text-decoration: underline;
}

.article-content :deep(ul),
.article-content :deep(ol) {
  margin: 16px 0;
  padding-left: 32px;
}

.article-content :deep(blockquote) {
  margin: 16px 0;
  padding: 12px 20px;
  border-left: 4px solid rgb(var(--v-theme-primary));
  background: rgba(var(--v-theme-primary), 0.05);
}

.article-content :deep(code) {
  background: rgba(var(--v-theme-on-surface), 0.1);
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  font-size: 0.9em;
}

.article-content :deep(pre) {
  background: rgba(var(--v-theme-on-surface), 0.95);
  color: #fff;
  padding: 16px;
  border-radius: 8px;
  overflow-x: auto;
  margin: 16px 0;
}

.article-content :deep(pre code) {
  background: none;
  padding: 0;
}
</style>
