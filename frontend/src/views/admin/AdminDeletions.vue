<template>
  <v-container fluid class="pa-4 pa-md-6">
    <!-- 删除申请列表 -->
    <v-card variant="flat" rounded="lg">
      <v-card-title class="pb-2">
        <v-icon start size="20">mdi-delete-forever</v-icon>
        删除申请
        <v-chip v-if="deletionRequests.length" size="small" color="error" class="ml-2">
          {{ deletionRequests.length }}
        </v-chip>
      </v-card-title>

      <v-list lines="two" v-if="deletionRequests.length > 0">
        <v-list-item v-for="request in deletionRequests" :key="request.id" class="py-3">
          <template v-slot:prepend>
            <v-avatar size="48" color="error" variant="tonal">
              <v-icon>mdi-file-document</v-icon>
            </v-avatar>
          </template>

          <v-list-item-title class="font-weight-medium mb-1">
            {{ request.Article?.title || '未知文章' }}
          </v-list-item-title>

          <v-list-item-subtitle>
            <div class="d-flex flex-wrap align-center ga-2">
              <span class="d-flex align-center text-caption">
                <v-icon size="14" color="primary" class="mr-1">mdi-account</v-icon>
                用户 ID: {{ request.user_id }}
              </span>
              <span class="d-flex align-center text-caption">
                <v-icon size="14" color="grey" class="mr-1">mdi-clock-outline</v-icon>
                {{ formatTime(request.created_at) }}
              </span>
            </div>
          </v-list-item-subtitle>

          <template v-slot:append>
            <v-btn-group variant="text" density="compact" divided>
              <v-btn size="small" color="success" @click="approveDeletion(request.id)">
                <v-icon>mdi-check</v-icon>
                <v-tooltip activator="parent">批准</v-tooltip>
              </v-btn>
              <v-btn size="small" color="error" @click="rejectDeletion(request.id)">
                <v-icon>mdi-close</v-icon>
                <v-tooltip activator="parent">拒绝</v-tooltip>
              </v-btn>
            </v-btn-group>
          </template>
        </v-list-item>
      </v-list>

      <v-card-text v-else class="text-center py-8">
        <v-icon size="48" color="grey-lighten-1">mdi-inbox</v-icon>
        <div class="text-body-1 text-medium-emphasis mt-2">
          暂无删除申请
        </div>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { adminDeletionApi } from '../../api/admin'
import { confirm, success, error } from '../../utils/modal'

const deletionRequests = ref([])

const loadDeletionRequests = async () => {
  try {
    const response = await adminDeletionApi.getRequests()
    deletionRequests.value = response.data.requests || []
  } catch (err) {
    console.error('加载删除申请失败', err)
    error('加载删除申请失败')
  }
}

const formatTime = (timeString) => {
  if (!timeString) return '-'
  const date = new Date(timeString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const approveDeletion = async (id) => {
  const confirmed = await confirm('确定要批准此删除申请吗？')
  if (!confirmed) return
  
  try {
    await adminDeletionApi.approveRequest(id)
    success('已批准删除')
    loadDeletionRequests()
  } catch (err) {
    console.error('批准删除失败', err)
    error(err.response?.data?.error || '操作失败')
  }
}

const rejectDeletion = async (id) => {
  const confirmed = await confirm('确定要拒绝此删除申请吗？')
  if (!confirmed) return
  
  try {
    await adminDeletionApi.rejectRequest(id)
    success('已拒绝删除')
    loadDeletionRequests()
  } catch (err) {
    console.error('拒绝删除失败', err)
    error(err.response?.data?.error || '操作失败')
  }
}

onMounted(() => {
  loadDeletionRequests()
})
</script>