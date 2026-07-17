<template>
  <div class="admin-articles">
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

    <div v-if="statusDialog.show" class="modal-overlay" @click.self="statusDialog.show = false">
      <div class="modal-content">
        <div class="layui-card">
          <div class="layui-card-header d-flex align-center gap-3">
            <i class="fa-solid fa-file-pen"></i>
            <span>修改文章状态</span>
          </div>
          <div class="layui-card-body">
            <div class="article-preview">
              <div class="preview-label">文章预览</div>
              <div class="preview-title">{{ statusDialog.article?.title }}</div>
            </div>
            <div class="status-radio-group mt-4">
              <label :class="['radio-label', { active: statusDialog.status === 'pending' }]">
                <input type="radio" v-model="statusDialog.status" value="pending" />
                <i class="fa-solid fa-circle-dot"></i>
                <span>待审核</span>
              </label>
              <label :class="['radio-label', { active: statusDialog.status === 'published' }]">
                <input type="radio" v-model="statusDialog.status" value="published" />
                <i class="fa-solid fa-circle-dot"></i>
                <span>已发布</span>
              </label>
              <label :class="['radio-label', { active: statusDialog.status === 'rejected' }]">
                <input type="radio" v-model="statusDialog.status" value="rejected" />
                <i class="fa-solid fa-circle-dot"></i>
                <span>已拒绝</span>
              </label>
            </div>
          </div>
          <div class="layui-card-footer d-flex justify-end gap-3">
            <button class="layui-btn layui-btn-primary" @click="statusDialog.show = false">取消</button>
            <button class="layui-btn" @click="handleEditStatus">确认修改</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import ArticlesPanel from './ArticlesPanel.vue'
import api from '../../api'
import { adminArticleApi } from '../../api/admin'
import { confirm, success, error } from '../../utils/modal'

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
  } catch (error) {
    console.error('加载文章列表失败', error)
  } finally {
    loading.value = false
  }
}

const loadCurrentUser = async () => {
  try {
    const response = await api.get('/profile')
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
    await adminArticleApi.updateArticleStatus(statusDialog.value.article.id, statusDialog.value.status)
    success('修改成功')
    statusDialog.value.show = false
    loadArticles()
  } catch (error) {
    console.error('修改状态失败', error)
    error(error.response?.data?.error || '修改失败')
  }
}

const handleDeleteArticle = async (article) => {
  const confirmed = await confirm(`确定要删除文章 "${article.title}" 吗？`)
  if (!confirmed) return
  try {
    await api.delete(`/articles/${article.id}`)
    success('删除成功')
    loadArticles()
  } catch (error) {
    console.error('删除文章失败', error)
    error(error.response?.data?.error || '删除失败')
  }
}

const handleRestoreArticle = async (article) => {
  const confirmed = await confirm(`确定要恢复文章 "${article.title}" 吗？`)
  if (!confirmed) return
  try {
    await api.post(`/articles/${article.id}/restore`)
    success('恢复成功')
    loadArticles()
  } catch (error) {
    console.error('恢复文章失败', error)
    error(error.response?.data?.error || '恢复失败')
  }
}

onMounted(() => {
  loadCurrentUser()
  loadArticles()
})
</script>

<style scoped>
.admin-articles {
  max-width: 1600px;
  margin: 0 auto;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  width: 90%;
  max-width: 480px;
  background: white;
  border-radius: 12px;
  overflow: hidden;
}

.d-flex {
  display: flex;
}

.align-center {
  align-items: center;
}

.justify-end {
  justify-content: flex-end;
}

.gap-3 {
  gap: 12px;
}

.mt-4 {
  margin-top: 16px;
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

.status-radio-group {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.radio-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 12px;
  border-radius: 8px;
  transition: background 0.3s;
}

.radio-label:hover {
  background: #f5f5f5;
}

.radio-label.active {
  background: #e6f7ff;
}

.radio-label input {
  display: none;
}

.radio-label i {
  font-size: 18px;
  color: #ccc;
}

.radio-label.active i {
  color: #1E9FFF;
}
</style>