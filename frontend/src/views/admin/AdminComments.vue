<template>
  <v-container fluid class="pa-6">
    <CommentsPanel
      :comments="allComments"
      :loading="loading"
      :pagination="pagination"
      @delete="handleDeleteComment"
      @refresh="loadComments"
      @page-change="handlePageChange"
    />
  </v-container>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import CommentsPanel from './CommentsPanel.vue'
import { adminCommentApi } from '../../api/admin'
import { confirm, success, error } from '../../utils/modal'

const allComments = ref([])
const loading = ref(true)
const pagination = ref({
  page: 1,
  pageSize: 20,
  total: 0,
  totalPages: 0
})

const loadComments = async (page = 1) => {
  loading.value = true
  try {
    console.log('开始加载评论列表，页码:', page)
    const response = await adminCommentApi.getComments({
      page: page,
      page_size: pagination.value.pageSize
    })
    console.log('评论列表响应:', response.data)
    allComments.value = response.data.comments || []
    pagination.value = {
      page: response.data.page || 1,
      pageSize: response.data.page_size || 20,
      total: response.data.total || 0,
      totalPages: response.data.total_pages || 0
    }
    console.log('加载的评论数量:', allComments.value.length)
    console.log('分页信息:', pagination.value)
  } catch (err) {
    console.error('加载评论列表失败', err)
    error('加载评论列表失败')
  } finally {
    loading.value = false
  }
}

const handlePageChange = (newPage) => {
  loadComments(newPage)
}

const handleDeleteComment = async (commentId) => {
  const confirmed = await confirm('确定要删除此评论吗？')
  if (!confirmed) return
  try {
    await adminCommentApi.deleteComment(commentId)
    success('删除成功')
    loadComments(pagination.value.page)
  } catch (err) {
    console.error('删除评论失败', err)
    error(err.response?.data?.error || '删除失败')
  }
}

onMounted(() => {
  loadComments()
})
</script>