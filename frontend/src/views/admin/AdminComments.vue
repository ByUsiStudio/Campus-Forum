<template>
  <v-container fluid class="pa-6">
    <CommentsPanel
      :comments="allComments"
      :loading="loading"
      @delete="handleDeleteComment"
      @refresh="loadComments"
    />
  </v-container>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import CommentsPanel from './CommentsPanel.vue'
import api from '../api'
import { showConfirm, showSuccess, showError } from '../utils/modal'

const allComments = ref([])
const loading = ref(true)

const loadComments = async () => {
  loading.value = true
  try {
    const response = await api.get('/admin/comments')
    allComments.value = response.data
  } catch (error) {
    console.error('加载评论列表失败', error)
  } finally {
    loading.value = false
  }
}

const handleDeleteComment = async (commentId) => {
  const confirmed = await showConfirm('确定要删除此评论吗？')
  if (!confirmed) return
  try {
    await api.delete(`/admin/comments/${commentId}`)
    showSuccess('删除成功')
    loadComments()
  } catch (error) {
    console.error('删除评论失败', error)
    showError(error.response?.data?.error || '删除失败')
  }
}

onMounted(() => {
  loadComments()
})
</script>