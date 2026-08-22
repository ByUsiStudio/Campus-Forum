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
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import ArticlesPanel from './ArticlesPanel.vue'
import api from '@/api'
import { adminArticleApi } from '@/api/admin'
import { confirm, success, error } from '@/utils/message'
import { jcCloseAll, jcFieldsConfig, jcOpenHtml } from '@/utils/jcu'

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
    article,
    status: article.status
  }

  const statusLabels = {
    pending: '待审核',
    published: '已发布',
    rejected: '已拒绝'
  }
  const current = statusDialog.value.article.status
  const cfg = jcFieldsConfig([
    {
      name: 'status',
      label: '设置状态',
      type: 'text',
      placeholder: '待审核 / 已发布 / 已拒绝',
      value: statusLabels[current] || ''
    }
  ])
  jcOpenHtml({
    title: `修改文章状态：${statusDialog.value.article.title}`,
    content: `
      <div style="margin-bottom:12px;font-size:13px;color:var(--jc-text-2,#64748b);">
        为文章设置新状态（${Object.entries(statusLabels).map(([k, v]) => `${k}: ${v}`).join('　')}）
      </div>
      ${cfg.html}
    `,
    width: 480,
    buttons: [
      { text: '取消', type: 'default', action: () => jcCloseAll() },
      {
        text: '确认修改',
        type: 'primary',
        action: (inst) => {
          if (!cfg.validate(inst.modalContent)) return
          const values = cfg.collect(inst.modalContent)
          const target = Object.keys(statusLabels).find(
            (k) => statusLabels[k] === values.status
          ) || values.status
          handleEditStatus(target)
        }
      }
    ]
  })
}

const handleEditStatus = async (newStatus) => {
  try {
    await adminArticleApi.updateArticleStatus(statusDialog.value.article.id, newStatus)
    success('修改成功')
    jcCloseAll()
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
