<template>
  <div class="create-article">
    <h2>{{ isEdit ? '编辑文章' : '发布新文章' }}</h2>
    
    <form @submit.prevent="submitArticle">
      <div class="form-group">
        <label>标题</label>
        <input type="text" v-model="form.title" required>
      </div>
      
      <div class="form-group">
        <label>分区</label>
        <select v-model="form.category_id" required>
          <option v-for="cat in categories" :key="cat.id" :value="cat.id">
            {{ cat.name }}
          </option>
        </select>
      </div>
      
      <div class="form-group">
        <label>内容 (支持Markdown)</label>
        <div class="editor-toolbar">
          <button type="button" @click="insertImage" class="toolbar-btn">插入图片</button>
          <button type="button" @click="insertVideo" class="toolbar-btn">插入视频</button>
        </div>
        <textarea v-model="form.content" rows="20" required placeholder="使用Markdown语法编写文章..."></textarea>
      </div>
      
      <div class="form-group" v-if="uploadProgress > 0">
        <div class="progress-bar">
          <div class="progress-fill" :style="{ width: uploadProgress + '%' }">
            {{ uploadProgress }}%
          </div>
        </div>
      </div>
      
      <div class="form-actions">
        <button type="submit" class="btn btn-primary" :disabled="submitting">
          {{ submitting ? '提交中...' : (isEdit ? '更新文章' : '发布文章') }}
        </button>
        <button type="button" @click="preview" class="btn btn-secondary">预览</button>
      </div>
    </form>
    
    <!-- 预览模态框 -->
    <div v-if="showPreview" class="preview-modal" @click="closePreview">
      <div class="preview-content" @click.stop>
        <h3>{{ form.title }}</h3>
        <div class="markdown-body" v-html="previewHtml"></div>
        <button @click="closePreview" class="btn btn-secondary">关闭</button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '../api'
import MarkdownIt from 'markdown-it'

const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true
})

export default {
  name: 'CreateArticle',
  setup() {
    const route = useRoute()
    const router = useRouter()
    const categories = ref([])
    const form = ref({
      title: '',
      category_id: null,
      content: ''
    })
    const isEdit = ref(false)
    const articleId = ref(null)
    const submitting = ref(false)
    const showPreview = ref(false)
    const uploadProgress = ref(0)
    
    const previewHtml = computed(() => {
      return md.render(form.value.content)
    })
    
    const loadCategories = async () => {
      try {
        const response = await api.get('/categories')
        categories.value = response.data.categories
        if (categories.value.length > 0 && !form.value.category_id) {
          form.value.category_id = categories.value[0].id
        }
      } catch (error) {
        console.error('加载分区失败', error)
      }
    }
    
    const loadArticle = async () => {
      try {
        const response = await api.get(`/articles/${articleId.value}`)
        const article = response.data.article
        form.value.title = article.title
        form.value.category_id = article.category_id
        form.value.content = article.content
      } catch (error) {
        console.error('加载文章失败', error)
      }
    }
    
    const insertImage = async () => {
      const input = document.createElement('input')
      input.type = 'file'
      input.accept = 'image/*'
      input.onchange = async (e) => {
        const file = e.target.files[0]
        if (!file) return
        
        const formData = new FormData()
        formData.append('image', file)
        
        uploadProgress.value = 0
        try {
          const response = await api.post('/upload/image', formData, {
            onUploadProgress: (progressEvent) => {
              const percent = Math.round((progressEvent.loaded * 100) / progressEvent.total)
              uploadProgress.value = percent
            }
          })
          const imageUrl = response.data.url
          const markdownImage = `![图片](${imageUrl})`
          insertAtCursor(markdownImage)
        } catch (error) {
          console.error('上传图片失败', error)
          alert('上传失败')
        } finally {
          setTimeout(() => {
            uploadProgress.value = 0
          }, 1000)
        }
      }
      input.click()
    }
    
    const insertVideo = async () => {
      const input = document.createElement('input')
      input.type = 'file'
      input.accept = 'video/*'
      input.onchange = async (e) => {
        const file = e.target.files[0]
        if (!file) return
        
        const formData = new FormData()
        formData.append('video', file)
        
        uploadProgress.value = 0
        try {
          const response = await api.post('/upload/video', formData, {
            onUploadProgress: (progressEvent) => {
              const percent = Math.round((progressEvent.loaded * 100) / progressEvent.total)
              uploadProgress.value = percent
            }
          })
          const videoUrl = response.data.url
          const videoHtml = `\n<video controls src="${videoUrl}"></video>\n`
          insertAtCursor(videoHtml)
        } catch (error) {
          console.error('上传视频失败', error)
          alert('上传失败')
        } finally {
          setTimeout(() => {
            uploadProgress.value = 0
          }, 1000)
        }
      }
      input.click()
    }
    
    const insertAtCursor = (text) => {
      const textarea = document.querySelector('textarea')
      const start = textarea.selectionStart
      const end = textarea.selectionEnd
      const content = form.value.content
      form.value.content = content.substring(0, start) + text + content.substring(end)
      setTimeout(() => {
        textarea.selectionStart = textarea.selectionEnd = start + text.length
        textarea.focus()
      }, 0)
    }
    
    const submitArticle = async () => {
      if (!form.value.title.trim()) {
        alert('请输入标题')
        return
      }
      if (!form.value.content.trim()) {
        alert('请输入内容')
        return
      }
      
      submitting.value = true
      
      try {
        let response
        if (isEdit.value) {
          response = await api.put(`/articles/${articleId.value}`, form.value)
        } else {
          response = await api.post('/articles', form.value)
        }
        router.push(`/article/${response.data.article.id}`)
      } catch (error) {
        console.error('提交失败', error)
        alert('提交失败')
      } finally {
        submitting.value = false
      }
    }
    
    const preview = () => {
      showPreview.value = true
    }
    
    const closePreview = () => {
      showPreview.value = false
    }
    
    onMounted(() => {
      const id = route.query.id
      if (id) {
        isEdit.value = true
        articleId.value = parseInt(id)
        loadArticle()
      }
      loadCategories()
    })
    
    return {
      categories,
      form,
      isEdit,
      submitting,
      showPreview,
      previewHtml,
      uploadProgress,
      submitArticle,
      insertImage,
      insertVideo,
      preview,
      closePreview
    }
  }
}
</script>

<style scoped>
.create-article {
  background: white;
  border-radius: 12px;
  padding: 30px;
  max-width: 900px;
  margin: 0 auto;
}

.create-article h2 {
  margin-bottom: 30px;
  color: #059669;
}

.editor-toolbar {
  margin-bottom: 10px;
  display: flex;
  gap: 10px;
}

.toolbar-btn {
  padding: 6px 12px;
  background: #f3f4f6;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  cursor: pointer;
}

.form-actions {
  display: flex;
  gap: 15px;
  margin-top: 20px;
}

.preview-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.preview-content {
  background: white;
  border-radius: 12px;
  padding: 30px;
  max-width: 800px;
  max-height: 80vh;
  overflow-y: auto;
  width: 90%;
}

.preview-content h3 {
  margin-bottom: 20px;
}
</style>