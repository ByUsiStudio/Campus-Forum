<template>
  <div class="page-container">
    <el-card shadow="never" class="deletions-card">
      <template #header>
        <div class="card-title">
          <el-icon :size="20" class="mr-1"><DeleteFilled /></el-icon>
          <span>删除申请</span>
          <el-tag
            v-if="deletionRequests.length"
            type="danger"
            size="small"
            class="ml-2"
          >
            {{ deletionRequests.length }}
          </el-tag>
        </div>
      </template>

      <el-empty
        v-if="deletionRequests.length === 0"
        description="暂无删除申请"
        :image-size="80"
      />

      <el-list
        v-else
        :data="deletionRequests"
        class="deletion-list"
      >
        <el-list-item
          v-for="request in deletionRequests"
          :key="request.id"
          class="deletion-item"
        >
          <div class="deletion-item-inner">
            <el-avatar :size="48" class="doc-avatar">
              <el-icon><Document /></el-icon>
            </el-avatar>

            <div class="deletion-info">
              <div class="deletion-title">
                {{ request.Article?.title || '未知文章' }}
              </div>
              <div class="deletion-meta">
                <span class="meta-item">
                  <el-icon :size="14"><User /></el-icon>
                  用户 ID: {{ request.user_id }}
                </span>
                <span class="meta-item">
                  <el-icon :size="14"><Clock /></el-icon>
                  {{ formatTime(request.created_at) }}
                </span>
              </div>
            </div>

            <div class="deletion-actions">
              <el-tooltip content="批准">
                <el-button
                  type="success"
                  size="small"
                  circle
                  @click="approveDeletion(request.id)"
                >
                  <el-icon><Check /></el-icon>
                </el-button>
              </el-tooltip>
              <el-tooltip content="拒绝">
                <el-button
                  type="danger"
                  size="small"
                  circle
                  @click="rejectDeletion(request.id)"
                >
                  <el-icon><Close /></el-icon>
                </el-button>
              </el-tooltip>
            </div>
          </div>
        </el-list-item>
      </el-list>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import {
  Check,
  Close,
  Clock,
  DeleteFilled,
  Document,
  User
} from '@element-plus/icons-vue'
import { adminDeletionApi } from '../../api/admin'
import { confirm, success, error } from '@/utils/message'

const deletionRequests = ref([])

const loadDeletionRequests = async () => {
  try {
    const response = await adminDeletionApi.getRequests()
    deletionRequests.value = response.data?.requests || response.data || []
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
  const confirmed = await confirm('确定要批准此删除申请吗？').catch(() => null)
  if (!confirmed) return

  try {
    await adminDeletionApi.approveRequest(id)
    success('已批准删除')
    await loadDeletionRequests()
  } catch (err) {
    console.error('批准删除失败', err)
    error(err.response?.data?.error || '操作失败')
  }
}

const rejectDeletion = async (id) => {
  const confirmed = await confirm('确定要拒绝此删除申请吗？').catch(() => null)
  if (!confirmed) return

  try {
    await adminDeletionApi.rejectRequest(id)
    success('已拒绝删除')
    await loadDeletionRequests()
  } catch (err) {
    console.error('拒绝删除失败', err)
    error(err.response?.data?.error || '操作失败')
  }
}

onMounted(loadDeletionRequests)
</script>

<style scoped>
.page-container {
  padding: 16px;
}
.card-title {
  display: inline-flex;
  align-items: center;
  font-weight: 600;
}
.mr-1 {
  margin-right: 4px;
}
.ml-2 {
  margin-left: 8px;
}
.deletion-list {
  --el-list-item-padding: 0;
}
.deletion-item {
  padding: 12px 0;
  border-bottom: 1px solid var(--el-border-color-lighter);
}
.deletion-item:last-child {
  border-bottom: none;
}
.deletion-item-inner {
  display: flex;
  align-items: center;
  gap: 16px;
  width: 100%;
}
.doc-avatar {
  background: var(--el-color-danger-light-9);
  color: var(--el-color-danger);
  flex-shrink: 0;
}
.deletion-info {
  flex: 1;
  min-width: 0;
}
.deletion-title {
  font-weight: 500;
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.deletion-meta {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 16px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}
.meta-item {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}
.deletion-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}
</style>
