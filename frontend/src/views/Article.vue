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
      <div ref="contentRef" class="article-content markdown-body" @click="handleContentClick" v-html="article.content_html"></div>
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
import { ref, onMounted, computed, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '../api'
import ImageViewer from '../components/ImageViewer.vue'
import { confirm as showConfirm, prompt as showPrompt, success as showSuccess } from '../utils/modal'

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
        
        await nextTick()
        processVideoTags()
      } catch (error) {
        console.error('加载文章失败', error)
        router.push('/')
      }
    }
    
    const extractVideosFromContent = () => {
      if (!article.value || !article.value.content) return []
      
      const videoRegex = /<video[^>]*src="([^"]+)"[^>]*\/?\s*>/gi
      const videos = []
      let match
      
      while ((match = videoRegex.exec(article.value.content)) !== null) {
        videos.push({
          src: match[1],
          poster: ''
        })
      }
      
      console.log('🔍 从原始 content 中提取到的视频:', videos)
      return videos
    }
    
    const processVideoTags = () => {
      console.log('=== processVideoTags 开始 ===')
      
      if (!contentRef.value) {
        console.log('❌ contentRef.value 为空')
        return
      }
      
      if (!article.value) {
        console.log('❌ article.value 为空')
        return
      }
      
      const videos = extractVideosFromContent()
      
      if (videos.length === 0) {
        console.log('⚠️ 没有找到任何视频')
        return
      }
      
      console.log('🔍 找到的视频数量:', videos.length)
      
      const pElements = contentRef.value.querySelectorAll('p')
      let foundPlaceholder = false
      
      pElements.forEach((p) => {
        const content = p.innerHTML
        if (content.includes('raw HTML omitted') || content.trim() === '') {
          foundPlaceholder = true
          
          const container = document.createElement('div')
          container.className = 'video-buttons-container'
          
          videos.forEach((video, index) => {
            const button = document.createElement('button')
            button.className = 'video-play-button'
            button.innerHTML = `
              <span class="play-icon-wrapper">
                <span class="play-icon">▶</span>
              </span>
              <span class="button-content">
                <span class="button-title">视频 ${index + 1}</span>
                <span class="button-hint">点击播放</span>
              </span>
            `
            
            button.addEventListener('click', () => {
              window.location.href = `/video?src=${encodeURIComponent(video.src)}&poster=${encodeURIComponent(video.poster)}&articleId=${article.value.id}`
            })
            
            container.appendChild(button)
          })
          
          p.parentNode.replaceChild(container, p)
          console.log(`✅ 已替换占位符，创建 ${videos.length} 个视频按钮`)
        }
      })
      
      if (!foundPlaceholder && videos.length > 0) {
        const container = document.createElement('div')
        container.className = 'video-buttons-container'
        
        videos.forEach((video, index) => {
          const button = document.createElement('button')
          button.className = 'video-play-button'
          button.innerHTML = `
            <span class="play-icon-wrapper">
              <span class="play-icon">▶</span>
            </span>
            <span class="button-content">
              <span class="button-title">视频 ${index + 1}</span>
              <span class="button-hint">点击播放</span>
            </span>
          `
          
          button.addEventListener('click', () => {
            window.location.href = `/video?src=${encodeURIComponent(video.src)}&poster=${encodeURIComponent(video.poster)}&articleId=${article.value.id}`
          })
          
          container.appendChild(button)
        })
        
        const contentDiv = contentRef.value.querySelector('div') || contentRef.value
        contentDiv.appendChild(container)
        console.log(`✅ 未找到占位符，在内容末尾添加 ${videos.length} 个视频按钮`)
      }
      
      console.log('=== processVideoTags 结束 ===')
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

.video-buttons-container {
  margin: 24px 0;
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  justify-content: center;
}

.video-play-button {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 18px 28px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  border: none;
  border-radius: 16px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 
    0 4px 20px rgba(102, 126, 234, 0.35),
    0 1px 3px rgba(0, 0, 0, 0.1);
  min-width: 220px;
  overflow: hidden;
  position: relative;
}

.video-play-button::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255,255,255,0.2), transparent);
  transition: left 0.5s ease;
}

.video-play-button:hover::before {
  left: 100%;
}

.video-play-button:hover {
  transform: translateY(-3px) scale(1.02);
  box-shadow: 
    0 10px 30px rgba(102, 126, 234, 0.45),
    0 2px 8px rgba(0, 0, 0, 0.15);
}

.video-play-button:active {
  transform: translateY(-1px) scale(0.98);
}

.play-icon-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 50%;
  transition: all 0.3s ease;
}

.video-play-button:hover .play-icon-wrapper {
  background: rgba(255, 255, 255, 0.3);
  transform: scale(1.1);
}

.play-icon {
  font-size: 18px;
  font-weight: bold;
  margin-left: 2px;
}

.button-content {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 2px;
}

.button-title {
  font-size: 16px;
  font-weight: 600;
}

.button-hint {
  font-size: 12px;
  opacity: 0.85;
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
