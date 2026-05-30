<template>
  <v-container fluid class="pa-6">
    <v-card>
      <v-card-title>
        <v-icon class="mr-2">mdi-delete-forever</v-icon>
        删除申请
        <v-chip v-if="deletionRequests.length" size="x-small" color="error" class="ml-2">
          {{ deletionRequests.length }}
        </v-chip>
      </v-card-title>
      <v-card-text>
        <v-list v-if="deletionRequests.length > 0">
          <v-list-item
            v-for="request in deletionRequests"
            :key="request.id"
            class="align-center"
          >
            <v-list-item-content>
              <v-list-item-title>{{ request.Article?.title }}</v-list-item-title>
              <v-list-item-subtitle>
                用户 ID: {{ request.user_id }} | 申请时间: {{ formatTime(request.created_at) }}
              </v-list-item-subtitle>
            </v-list-item-content>
            <v-list-item-action>
              <v-btn
                variant="text"
                color="success"
                size="small"
                @click="approveDeletion(request.id)"
              >
                <v-icon>mdi-check</v-icon>
                批准
              </v-btn>
              <v-btn
                variant="text"
                color="error"
                size="small"
                @click="rejectDeletion(request.id)"
              >
                <v-icon>mdi-x</v-icon>
                拒绝
              </v-btn>
            </v-list-item-action>
          </v-list-item>
        </v-list>
        
        <div v-else class="text-center text-gray-400 py-8">
          <v-icon size="48" color="grey lighten-3">mdi-inbox</v-icon>
          <p class="mt-2">暂无删除申请</p>
        </div>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'
import { showConfirm, showSuccess, showError } from '../utils/modal'

const deletionRequests = ref([])

const loadDeletionRequests = async () => {
  try {
    const response = await api.get('/deletion-requests')
    deletionRequests.value = response.data
  } catch (error) {
    console.error('加载删除申请失败', error)
  }
}

const formatTime = (timeString) => {
  if (!timeString) return ''
  const date = new Date(timeString)
  return date.toLocaleString('zh-CN')
}

const approveDeletion = async (id) => {
  const confirmed = await showConfirm('确定要批准此删除申请吗？')
  if (!confirmed) return
  
  try {
    await api.post(`/deletion-requests/${id}/approve`)
    showSuccess('已批准删除')
    loadDeletionRequests()
  } catch (error) {
    console.error('批准删除失败', error)
    showError(error.response?.data?.error || '操作失败')
  }
}

const rejectDeletion = async (id) => {
  const confirmed = await showConfirm('确定要拒绝此删除申请吗？')
  if (!confirmed) return
  
  try {
    await api.post(`/deletion-requests/${id}/reject`)
    showSuccess('已拒绝删除')
    loadDeletionRequests()
  } catch (error) {
    console.error('拒绝删除失败', error)
    showError(error.response?.data?.error || '操作失败')
  }
}

onMounted(() => {
  loadDeletionRequests()
})
</script>