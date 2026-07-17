<template>
  <div class="admin-deletions">
    <div class="layui-card">
      <div class="layui-card-header">
        <i class="fa-solid fa-trash"></i>
        <span>删除申请</span>
        <span v-if="deletionRequests.length" class="pending-badge">{{ deletionRequests.length }}</span>
      </div>

      <div v-if="deletionRequests.length > 0" class="request-list">
        <div v-for="request in deletionRequests" :key="request.id" class="request-item">
          <div class="request-avatar">
            <i class="fa-solid fa-file-lines"></i>
          </div>
          <div class="request-content">
            <div class="request-title">{{ request.Article?.title || '未知文章' }}</div>
            <div class="request-meta">
              <span class="meta-item">
                <i class="fa-solid fa-user"></i>用户 ID: {{ request.user_id }}
              </span>
              <span class="meta-item">
                <i class="fa-solid fa-clock"></i>{{ formatTime(request.created_at) }}
              </span>
            </div>
          </div>
          <div class="request-actions">
            <button class="action-btn approve" @click="approveDeletion(request.id)">
              <i class="fa-solid fa-check"></i>批准
            </button>
            <button class="action-btn reject" @click="rejectDeletion(request.id)">
              <i class="fa-solid fa-xmark"></i>拒绝
            </button>
          </div>
        </div>
      </div>

      <div v-else class="empty-state">
        <i class="fa-solid fa-inbox"></i>
        <div class="empty-text">暂无删除申请</div>
      </div>
    </div>
  </div>
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

<style scoped>
.admin-deletions {
  padding: 20px;
}

.layui-card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
  font-weight: 600;
}

.pending-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 20px;
  height: 20px;
  padding: 0 6px;
  background: #FF5722;
  color: white;
  font-size: 12px;
  font-weight: 600;
  border-radius: 10px;
}

.request-list {
  padding: 8px 0;
}

.request-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;

  &:last-child {
    border-bottom: none;
  }
}

.request-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: rgba(255, 87, 34, 0.1);
  color: #FF5722;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  flex-shrink: 0;
}

.request-content {
  flex: 1;
  min-width: 0;
}

.request-title {
  font-size: 15px;
  font-weight: 500;
  color: #333;
  margin-bottom: 8px;
}

.request-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #999;
}

.request-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  border: none;
  border-radius: 4px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.3s;

  &.approve {
    background: rgba(82, 196, 26, 0.1);
    color: #52C41A;

    &:hover {
      background: rgba(82, 196, 26, 0.2);
    }
  }

  &.reject {
    background: rgba(255, 87, 34, 0.1);
    color: #FF5722;

    &:hover {
      background: rgba(255, 87, 34, 0.2);
    }
  }
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
}

.empty-state i {
  font-size: 48px;
  color: #e0e0e0;
  margin-bottom: 16px;
}

.empty-text {
  font-size: 15px;
  color: #999;
}
</style>