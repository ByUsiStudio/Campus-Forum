<template>
  <div v-if="article">
    <v-card class="pa-6 mb-4" v-if="article">
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
      <div class="article-content markdown-body" v-html="article.content_html" @click="handleImageClick"></div>
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
    
    <!-- 图片查看器 -->
    <ImageViewer v-if="showImageViewer" :url="currentImageUrl" @close="closeImageViewer" />
  </div>
  
  <div v-else class="d-flex justify-center align-center" style="min-height: 50vh;">
    <v-progress-circular indeterminate color="primary" size="64"></v-progress-circular>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
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
      } catch (error) {
        console.error('加载文章失败', error)
        router.push('/')
      }
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
        // 用户取消
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
        // 用户取消
      }
    }
    
    const canDeleteComment = (comment) => {
      if (!currentUser.value) return false
      return currentUser.value.id === comment.user_id || currentUser.value.role === 'admin'
    }
    
    const handleImageClick = (event) => {
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
      toggleLike,
      submitComment,
      deleteArticle,
      deleteComment,
      canDeleteComment,
      handleImageClick,
      closeImageViewer,
      formatDate
    }
  }
}
</script>

<style scoped>
.article-content {
  min-height: 300px;
}
</style>
