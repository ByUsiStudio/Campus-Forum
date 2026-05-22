<template>
  <div>
    <v-card class="pa-6" max-width="900" style="margin: 0 auto;">
      <v-card-title class="text-h5 mb-4" style="color: rgb(var(--v-theme-primary));">
        {{ isEdit ? '编辑文章' : '发布新文章' }}
      </v-card-title>
      
      <v-form @submit.prevent="submitArticle">
        <v-text-field
          v-model="form.title"
          label="标题"
          variant="outlined"
          required
          class="mb-4"
        ></v-text-field>
        
        <v-select
          v-model="form.category_id"
          :items="categories"
          item-title="name"
          item-value="id"
          label="分区"
          variant="outlined"
          class="mb-4"
        ></v-select>
        
        <div class="mb-2">
          <v-btn-toggle class="mb-2">
            <v-btn size="small" variant="outlined" @click="insertImage">
              <v-icon start>mdi-image</v-icon>
              插入图片
            </v-btn>
            <v-btn size="small" variant="outlined" @click="insertVideo">
              <v-icon start>mdi-video</v-icon>
              插入视频
            </v-btn>
          </v-btn-toggle>
        </div>
        
        <v-textarea
          v-model="form.content"
          label="内容 (支持Markdown)"
          variant="outlined"
          rows="20"
          required
          placeholder="使用Markdown语法编写文章..."
          class="mb-4"
        ></v-textarea>
        
        <v-progress-linear
          v-if="uploadProgress > 0"
          :model-value="uploadProgress"
          color="primary"
          height="20"
          class="mb-4"
        >
          {{ uploadProgress }}%
        </v-progress-linear>
        
        <div class="d-flex gap-4">
          <v-btn
            type="submit"
            color="primary"
            size="large"
            :loading="submitting"
          >
            {{ submitting ? '提交中...' : (isEdit ? '更新文章' : '发布文章') }}
          </v-btn>
          <v-btn
            variant="outlined"
            size="large"
            @click="showPreview = true"
          >
            预览
          </v-btn>
        </div>
      </v-form>
    </v-card>
    
    <!-- 预览对话框 -->
    <v-dialog v-model="showPreview" max-width="800" scrollable>
      <v-card>
        <v-card-title class="text-h6">{{ form.title }}</v-card-title>
        <v-divider></v-divider>
        <v-card-text class="pa-4">
          <div class="markdown-body" v-html="previewHtml"></div>
        </v-card-text>
        <v-divider></v-divider>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn variant="text" @click="showPreview = false">关闭</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
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
      if (!textarea) return
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
      insertVideo
    }
  }
}
</script>
