<template>
  <div class="article-container" v-if="article">
    <div class="article-header">
      <h1>{{ article.title }}</h1>
      <div class="article-meta">
        <span>作者：{{ article.user.display_name }}</span>
        <span>分区：{{ article.category.name }}</span>
        <span>发布时间：{{ formatDate(article.created_at) }}</span>
        <span>阅读：{{ article.view_count }}</span>
      </div>
      <div class="article-actions" v-if="canEdit">
        <router-link :to="'/create?id=' + article.id" class="btn btn-secondary">编辑</router-link>
        <button @click="deleteArticle" class="btn btn-danger">删除</button>
      </div>
    </div>
    
    <div class="article-content markdown-body" v-html="article.content_html" @click="handleImageClick"></div>
    
    <div class="article-footer">
      <button @click="toggleLike" class="btn-like" :class="{ liked: liked }">
        {{ article.like_count }} 点赞
      </button>
    </div>
    
    <div class="comments-section">
      <h3>评论 ({{ comments.length }})</h3>
      <div class="comment-form" v-if="token">
        <textarea v-model="commentContent" placeholder="写下你的评论..." rows="3"></textarea>
        <button @click="submitComment" class="btn btn-primary">发表评论</button>
      </div>
      
      <div class="comments-list">
        <div v-for="comment in comments" :key="comment.id" class="comment-item">
          <div class="comment-header">
            <span>{{ comment.user.display_name }}</span>
            <span>{{ formatDate(comment.created_at) }}</span>
          </div>
          <div class="comment-content">{{ comment.content }}</div>
          <button v-if="canDeleteComment(comment)" @click="deleteComment(comment.id)" class="btn-delete">删除</button>
        </div>
      </div>
    </div>
    
    <!-- 图片查看器 -->
    <ImageViewer v-if="showImageViewer" :url="currentImageUrl" @close="closeImageViewer" />
  </div>
  
  <div v-else class="loading">加载中...</div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '../api'
import ImageViewer from '../components/ImageViewer.vue'

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
      if (!confirm('确定要删除这篇文章吗？')) return
      
      if (currentUser.value?.role !== 'admin') {
        const reason = prompt('请输入删除原因（管理员将审核）：')
        if (!reason) return
        
        try {
          await api.delete(`/articles/${article.value.id}`, { data: { reason } })
          alert('删除申请已提交，等待管理员审核')
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
    }
    
    const deleteComment = async (commentId) => {
      if (!confirm('确定要删除这条评论吗？')) return
      
      try {
        await api.delete(`/comments/${commentId}`)
        comments.value = comments.value.filter(c => c.id !== commentId)
      } catch (error) {
        console.error('删除评论失败', error)
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
.article-container {
  background: white;
  border-radius: 12px;
  padding: 30px;
}

.article-header {
  border-bottom: 1px solid #e5e7eb;
  padding-bottom: 20px;
  margin-bottom: 20px;
}

.article-header h1 {
  margin-bottom: 15px;
  color: #1e293b;
}

.article-meta {
  display: flex;
  gap: 20px;
  font-size: 14px;
  color: #6b7280;
}

.article-actions {
  margin-top: 15px;
  display: flex;
  gap: 10px;
}

.article-content {
  min-height: 300px;
  margin-bottom: 30px;
}

.article-footer {
  border-top: 1px solid #e5e7eb;
  padding-top: 20px;
}

.btn-like {
  padding: 8px 20px;
  border: 1px solid #d1d5db;
  background: white;
  border-radius: 20px;
  cursor: pointer;
  font-size: 14px;
}

.btn-like.liked {
  background: #10b981;
  color: white;
  border-color: #10b981;
}

.comments-section {
  margin-top: 40px;
}

.comments-section h3 {
  margin-bottom: 20px;
}

.comment-form {
  margin-bottom: 30px;
}

.comment-form textarea {
  width: 100%;
  padding: 10px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  margin-bottom: 10px;
  resize: vertical;
}

.btn-delete {
  background: none;
  border: none;
  color: #ef4444;
  cursor: pointer;
  font-size: 12px;
  margin-top: 5px;
}
</style>