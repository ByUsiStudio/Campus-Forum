<template>
  <div>
    <el-card class="page-container" shadow="never">
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
        @restore="handleRestoreArticle"
        @refresh="loadArticles"
        @update:page="articlePage = $event; loadArticles()"
        @update:filter="articleFilter = $event; articlePage = 1; loadArticles()"
      />
    </el-card>

    <el-dialog v-model="statusDialog.show" width="480px" class="status-dialog">
      <template #header>
        <div class="dialog-title">
          <el-icon class="title-icon"><Edit /></el-icon>
          <span>修改文章状态</span>
        </div>
      </template>

      <div class="dialog-body">
        <div class="article-preview">
          <div class="preview-label">文章预览</div>
          <div class="preview-title">{{ statusDialog.article?.title }}</div>
        </div>
        <el-radio-group v-model="statusDialog.status" class="mt-4">
          <el-radio value="pending">待审核</el-radio>
          <el-radio value="published">已发布</el-radio>
          <el-radio value="rejected">已拒绝</el-radio>
        </el-radio-group>
      </div>

      <template #footer>
        <div class="dialog-actions">
          <el-button @click="statusDialog.show = false">取消</el-button>
          <el-button type="primary" @click="handleEditStatus">确认修改</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Edit } from '@element-plus/icons-vue'
import ArticlesPanel from './ArticlesPanel.vue'
import api from '@/api'
import { adminArticleApi } from '@/api/admin'
import { confirm, success, error } from '@/utils/message'

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
  { value: 'rejected', label: '已拒绝' },
  { value: 'deleted', label: '已删除' }
]

const statusDialog = ref({
  show: false,
  article: null,
  status: 'published'
})

const loadArticles = async () => {
  loading.value = true
  try {
    const response = await adminArticleApi.getArticles({
      page: articlePage.value,
      status: articleFilter.value
    })
    articles.value = response.data.articles
    articleTotalPages.value = response.data.total_pages
  } catch (e) {
    console.error('加载文章列表失败', e)
  } finally {
    loading.value = false
  }
}

const loadCurrentUser = async () => {
  try {
    const response = await api.get('/profile')
    currentUserRole.value = response.data.role
  } catch (e) {
    console.error('加载当前用户失败', e)
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
    await adminArticleApi.updateArticleStatus(statusDialog.value.article.id, statusDialog.value.status)
    success('修改成功')
    statusDialog.value.show = false
    loadArticles()
  } catch (e) {
    console.error('修改状态失败', e)
    error(e.response?.data?.error || '修改失败')
  }
}

const handleDeleteArticle = async (article) => {
  const confirmed = await confirm(`确定要删除文章 "${article.title}" 吗？`).catch(() => null)
  if (!confirmed) return
  try {
    await api.delete(`/articles/${article.id}`)
    success('删除成功')
    loadArticles()
  } catch (e) {
    console.error('删除文章失败', e)
    error(e.response?.data?.error || '删除失败')
  }
}

const handleRestoreArticle = async (article) => {
  const confirmed = await confirm(`确定要恢复文章 "${article.title}" 吗？`).catch(() => null)
  if (!confirmed) return
  try {
    await api.post(`/articles/${article.id}/restore`)
    success('恢复成功')
    loadArticles()
  } catch (e) {
    console.error('恢复文章失败', e)
    error(e.response?.data?.error || '恢复失败')
  }
}

onMounted(() => {
  loadCurrentUser()
  loadArticles()
})
</script>

<style scoped>
.dialog-title {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 1.2rem;
  font-weight: 700;
}

.title-icon {
  width: 40px;
  height: 40px;
  font-size: 20px;
  padding: 8px;
  border-radius: 10px;
  background: rgba(103, 80, 164, 0.1);
  color: #6750a4;
}

.dialog-body {
  padding: 8px 4px;
}

.dialog-actions {
  display: flex;
  justify-content: flex-end;
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
