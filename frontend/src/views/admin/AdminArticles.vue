<script setup>
import { ref, computed, onMounted } from 'vue'
import { adminArticleApi } from '../../api/admin'

const loading = ref(false)
const articles = ref([])
const page = ref(1)
const pageSize = ref(20)
const totalArticles = ref(0)
const totalPages = computed(() => Math.ceil(totalArticles.value / pageSize.value))
const statusFilter = ref('all')
const statusDialog = ref(false)
const selectedArticle = ref(null)
const newStatus = ref('published')

const headers = [
  { title: '文章', key: 'title', sortable: false },
  { title: '作者', key: 'author', width: '120px' },
  { title: '分类', key: 'category', width: '100px' },
  { title: '状态', key: 'status', width: '100px' },
  { title: '浏览', key: 'view_count', width: '80px' },
  { title: '点赞', key: 'like_count', width: '80px' },
  { title: '创建时间', key: 'created_at', width: '150px' },
  { title: '操作', key: 'actions', width: '80px', sortable: false }
]

const statusOptions = [
  { title: '全部', value: 'all' },
  { title: '已发布', value: 'published' },
  { title: '待审核', value: 'pending' },
  { title: '已拒绝', value: 'rejected' },
  { title: '已删除', value: 'deleted' }
]

const loadArticles = async () => {
  loading.value = true
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value
    }
    if (statusFilter.value !== 'all') {
      params.status = statusFilter.value
    }
    const response = await adminArticleApi.getArticles(params)
    articles.value = response.data.articles || []
    totalArticles.value = response.data.total || 0
  } catch (error) {
    console.error('加载文章失败:', error)
  } finally {
    loading.value = false
  }
}

const getStatusColor = (status) => {
  const colors = {
    published: 'success',
    pending: 'warning',
    rejected: 'error',
    deleted: 'grey'
  }
  return colors[status] || 'grey'
}

const getStatusText = (status) => {
  const texts = {
    published: '已发布',
    pending: '待审核',
    rejected: '已拒绝',
    deleted: '已删除'
  }
  return texts[status] || '未知'
}

const changeStatus = (article) => {
  selectedArticle.value = article
  newStatus.value = article.status
  statusDialog.value = true
}

const saveStatus = async () => {
  loading.value = true
  try {
    await adminArticleApi.updateArticleStatus(selectedArticle.value.id, newStatus.value)
    statusDialog.value = false
    loadArticles()
  } catch (error) {
    console.error('更新状态失败:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadArticles()
})
</script>

<template>
  <v-container fluid class="pa-0">
    <!-- 页面标题 -->
    <div class="mb-6">
      <h1 class="text-h5 font-weight-bold">文章管理</h1>
      <p class="text-body-2 text-grey">管理所有文章，审核、发布、删除文章</p>
    </div>

    <!-- 筛选 -->
    <v-card class="mb-4">
      <v-card-text>
        <v-row align="center" dense>
          <v-col cols="12" md="4">
            <v-select
              v-model="statusFilter"
              :items="statusOptions"
              label="文章状态"
              hide-details
              @update:model-value="loadArticles"
            />
          </v-col>
          <v-col cols="12" md="8" class="text-right">
            <v-btn color="primary" @click="loadArticles" :loading="loading">
              <v-icon start>mdi-refresh</v-icon>
              刷新
            </v-btn>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

    <!-- 文章列表 -->
    <v-card>
      <v-card-title class="d-flex align-center justify-space-between">
        <span>文章列表 ({{ totalArticles }})</span>
        <v-pagination
          v-model="page"
          :length="totalPages"
          :total-visible="5"
          density="compact"
          @update:model-value="loadArticles"
        />
      </v-card-title>

      <v-data-table
        :headers="headers"
        :items="articles"
        :loading="loading"
        :items-per-page="pageSize"
      >
        <template v-slot:item.title="{ item }">
          <div class="py-2">
            <div class="font-weight-medium">{{ item.title }}</div>
            <div class="text-caption text-grey">{{ item.summary || '无摘要' }}</div>
          </div>
        </template>

        <template v-slot:item.author="{ item }">
          <span>{{ item.user?.username || '未知' }}</span>
        </template>

        <template v-slot:item.category="{ item }">
          <span>{{ item.category?.name || '无分类' }}</span>
        </template>

        <template v-slot:item.status="{ item }">
          <v-chip :color="getStatusColor(item.status)" size="small">
            {{ getStatusText(item.status) }}
          </v-chip>
        </template>

        <template v-slot:item.view_count="{ item }">
          <span>{{ item.view_count || 0 }}</span>
        </template>

        <template v-slot:item.like_count="{ item }">
          <span>{{ item.like_count || 0 }}</span>
        </template>

        <template v-slot:item.created_at="{ item }">
          <span class="text-caption">{{ new Date(item.created_at).toLocaleDateString('zh-CN') }}</span>
        </template>

        <template v-slot:item.actions="{ item }">
          <v-menu>
            <template v-slot:activator="{ props }">
              <v-btn icon variant="text" size="small" v-bind="props">
                <v-icon>mdi-dots-vertical</v-icon>
              </v-btn>
            </template>
            <v-list density="compact">
              <v-list-item prepend-icon="mdi-eye">
                <v-list-item-title>查看文章</v-list-item-title>
              </v-list-item>
              <v-list-item @click="changeStatus(item)" prepend-icon="mdi-pencil">
                <v-list-item-title>修改状态</v-list-item-title>
              </v-list-item>
            </v-list>
          </v-menu>
        </template>
      </v-data-table>
    </v-card>

    <!-- 修改状态对话框 -->
    <v-dialog v-model="statusDialog" max-width="400">
      <v-card>
        <v-card-title>修改文章状态</v-card-title>
        <v-card-text>
          <p class="mb-4">文章: <strong>{{ selectedArticle?.title }}</strong></p>
          <v-radio-group v-model="newStatus">
            <v-radio label="已发布" value="published" color="success"></v-radio>
            <v-radio label="待审核" value="pending" color="warning"></v-radio>
            <v-radio label="已拒绝" value="rejected" color="error"></v-radio>
            <v-radio label="已删除" value="deleted"></v-radio>
          </v-radio-group>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn text @click="statusDialog = false">取消</v-btn>
          <v-btn color="primary" @click="saveStatus" :loading="loading">保存</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>