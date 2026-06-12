<script setup>
import { ref, computed, onMounted } from 'vue'
import { adminCommentApi } from '../../api/admin'

const loading = ref(false)
const comments = ref([])
const page = ref(1)
const pageSize = ref(20)
const totalComments = ref(0)
const totalPages = computed(() => Math.ceil(totalComments.value / pageSize.value))

const headers = [
  { title: '评论内容', key: 'content', sortable: false },
  { title: '用户', key: 'user', width: '120px' },
  { title: '文章', key: 'article', width: '150px' },
  { title: '创建时间', key: 'created_at', width: '150px' },
  { title: '操作', key: 'actions', width: '80px', sortable: false }
]

const loadComments = async () => {
  loading.value = true
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value
    }
    const response = await adminCommentApi.getComments(params)
    comments.value = response.data.comments || []
    totalComments.value = response.data.total || 0
  } catch (error) {
    console.error('加载评论失败:', error)
  } finally {
    loading.value = false
  }
}

const deleteComment = async (comment) => {
  if (!confirm(`确定要删除这条评论吗？此操作不可恢复！`)) return
  
  loading.value = true
  try {
    await adminCommentApi.deleteComment(comment.id)
    loadComments()
  } catch (error) {
    console.error('删除评论失败:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadComments()
})
</script>

<template>
  <v-container fluid class="pa-0">
    <!-- 页面标题 -->
    <div class="mb-6">
      <h1 class="text-h5 font-weight-bold">评论管理</h1>
      <p class="text-body-2 text-grey">管理所有评论，删除违规评论</p>
    </div>

    <!-- 刷新按钮 -->
    <v-card class="mb-4">
      <v-card-text>
        <v-btn color="primary" @click="loadComments" :loading="loading">
          <v-icon start>mdi-refresh</v-icon>
          刷新
        </v-btn>
      </v-card-text>
    </v-card>

    <!-- 评论列表 -->
    <v-card>
      <v-card-title class="d-flex align-center justify-space-between">
        <span>评论列表 ({{ totalComments }})</span>
        <v-pagination
          v-model="page"
          :length="totalPages"
          :total-visible="5"
          density="compact"
          @update:model-value="loadComments"
        />
      </v-card-title>

      <v-data-table
        :headers="headers"
        :items="comments"
        :loading="loading"
        :items-per-page="pageSize"
      >
        <template v-slot:item.content="{ item }">
          <div class="py-2">
            <div class="text-body-2">{{ item.content }}</div>
          </div>
        </template>

        <template v-slot:item.user="{ item }">
          <span>{{ item.user?.username || '未知' }}</span>
        </template>

        <template v-slot:item.article="{ item }">
          <span>{{ item.article?.title || '未知文章' }}</span>
        </template>

        <template v-slot:item.created_at="{ item }">
          <span class="text-caption">{{ new Date(item.created_at).toLocaleDateString('zh-CN') }}</span>
        </template>

        <template v-slot:item.actions="{ item }">
          <v-btn icon variant="text" size="small" color="error" @click="deleteComment(item)">
            <v-icon>mdi-delete</v-icon>
          </v-btn>
        </template>
      </v-data-table>
    </v-card>
  </v-container>
</template>