<template>
  <v-container fluid class="pa-6">
    <ArticlesPanel
      :articles="articles"
      :loading="loading"
      :page="articlePage"
      :total-pages="articleTotalPages"
      :filter="articleFilter"
      :status-options="articleStatusOptions"
      :current-user-role="currentUserRole"
      @change-status="showStatusDialog"
      @delete="handleDeleteArticle"
      @refresh="loadArticles"
      @update:page="articlePage = $event; loadArticles()"
      @update:filter="articleFilter = $event; articlePage = 1; loadArticles()"
    />
  </v-container>

  <v-dialog v-model="statusDialog.show" max-width="480">
    <v-card class="dialog-card">
      <v-card-title class="dialog-title">
        <v-icon class="title-icon">mdi-file-edit</v-icon>
        修改文章状态
      </v-card-title>
      <v-card-text class="dialog-body">
        <div class="article-preview">
          <div class="preview-label">文章预览</div>
          <div class="preview-title">{{ statusDialog.article?.title }}</div>
        </div>
        <v-radio-group v-model="statusDialog.status" class="mt-4">
          <v-radio label="待审核" value="pending" color="warning"></v-radio>
          <v-radio label="已发布" value="published" color="success"></v-radio>
          <v-radio label="已拒绝" value="rejected" color="error"></v-radio>
        </v-radio-group>
      </v-card-text>
      <v-card-actions class="dialog-actions">
        <v-btn variant="text" @click="statusDialog.show = false">取消</v-btn>
        <v-btn color="primary" variant="flat" @click="handleEditStatus">确认修改</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import ArticlesPanel from './ArticlesPanel.vue'
import api from '../api'
import { showConfirm, showSuccess, showError } from '../utils/modal'

const router = useRouter()
const articles = ref([])
const loading = ref(true)
const articlePage = ref(1)
const articleTotalPages = ref(1)
const articleFilter = ref('all')
const currentUserRole = ref(null)

const articleStatusOptions = [
  { value: 'all', label: '全部' },
  { value: 'pending', label: '待审核' },
  { value: 'published', label: '已发布' },
  { value: 'rejected', label: '已拒绝' }
]

const statusDialog = ref({
  show: false,
  article: null,
  status: 'published'
})

const loadArticles = async () => {
  loading.value = true
  try {
    const response = await api.get(`/admin/articles?page=${articlePage.value}&status=${articleFilter.value}`)
    articles.value = response.data.articles
    articleTotalPages.value = response.data.total_pages
  } catch (error) {
    console.error('加载文章列表失败', error)
  } finally {
    loading.value = false
  }
}

const loadCurrentUser = async () => {
  try {
    const response = await api.get('/user')
    currentUserRole.value = response.data.role
  } catch (error) {
    console.error('加载当前用户失败', error)
    router.push('/login')
  }
}

const showStatusDialog = (article) => {
  statusDialog.value = {
    show: true,
    article,
    status: article.status
  }
}

const handleEditStatus = async () => {
  try {
    await api.put(`/admin/articles/${statusDialog.value.article.id}/status`, {
      status: statusDialog.value.status
    })
    showSuccess('修改成功')
    statusDialog.value.show = false
    loadArticles()
  } catch (error) {
    console.error('修改状态失败', error)
    showError(error.response?.data?.error || '修改失败')
  }
}

const handleDeleteArticle = async (article) => {
  const confirmed = await showConfirm(`确定要删除文章 "${article.title}" 吗？`)
  if (!confirmed) return
  try {
    await api.delete(`/admin/articles/${article.id}`)
    showSuccess('删除成功')
    loadArticles()
  } catch (error) {
    console.error('删除文章失败', error)
    showError(error.response?.data?.error || '删除失败')
  }
}

onMounted(() => {
  loadCurrentUser()
  loadArticles()
})
</script>

<style scoped>
.dialog-card {
  border-radius: 20px !important;
  overflow: hidden;
}

.dialog-title {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 24px 24px 16px;
  font-size: 1.2rem;
  font-weight: 700;
  background: linear-gradient(135deg, #f8f9ff 0%, #fff 100%);
}

.title-icon {
  width: 40px;
  height: 40px;
  padding: 8px;
  border-radius: 10px;
  background: rgba(103, 80, 164, 0.1);
}

.dialog-body {
  padding: 24px !important;
}

.dialog-actions {
  padding: 16px 24px 24px;
  gap: 12px;
}

.article-preview {
  padding: 16px;
  background: #f8f9ff;
  border-radius: 12px;
}

.preview-label {
  font-size: 0.75rem;
  color: #6b7280;
  margin-bottom: 8px;
}

.preview-title {
  font-size: 1rem;
  font-weight: 600;
  color: #1a1a2e;
}
</style>