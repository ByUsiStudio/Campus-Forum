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
import { adminCommentApi } from '../../api/admin'
import { confirm, success, error } from '../../utils/modal'

const allComments = ref([])
const loading = ref(true)

const loadComments = async () => {
  loading.value = true
  try {
    const response = await adminCommentApi.getComments()
    allComments.value = response.data.comments || []
  } catch (error) {
    console.error('加载评论列表失败', error)
  } finally {
    loading.value = false
  }
}

const handleDeleteComment = async (commentId) => {
  const confirmed = await confirm('确定要删除此评论吗？')
  if (!confirmed) return
  try {
    await adminCommentApi.deleteComment(commentId)
    success('删除成功')
    loadComments()
  } catch (error) {
    console.error('删除评论失败', error)
    error(error.response?.data?.error || '删除失败')
  }
}

onMounted(() => {
  loadComments()
})
</script>