<template>
  <div class="admin-notifications">
    <div class="notifications-container">
      <div class="left-panel">
        <div class="layui-card">
          <div class="layui-card-header">
            <i class="fa-solid fa-bell"></i>
            <span>发送通知</span>
          </div>
          <div class="layui-card-body">
            <div class="form-group">
              <label class="form-label">通知类型</label>
              <select v-model="notificationForm.type" class="layui-input">
                <option v-for="type in notificationTypes" :key="type.value" :value="type.value">
                  {{ type.label }}
                </option>
              </select>
            </div>

            <div class="form-group">
              <label class="form-label">发送对象</label>
              <select v-model="notificationForm.target" class="layui-input">
                <option v-for="target in notificationTargets" :key="target.value" :value="target.value">
                  {{ target.label }}
                </option>
              </select>
            </div>

            <div class="form-group">
              <label class="form-label">通知标题</label>
              <input
                v-model="notificationForm.title"
                type="text"
                placeholder="请输入通知标题"
                class="layui-input"
              />
            </div>

            <div class="form-group">
              <label class="form-label">通知内容</label>
              <textarea
                v-model="notificationForm.content"
                placeholder="请输入通知内容..."
                class="layui-textarea"
                rows="4"
                maxlength="500"
              ></textarea>
              <span class="char-count">{{ notificationForm.content.length }}/500</span>
            </div>
          </div>
          <div class="layui-card-footer">
            <button
              class="layui-btn"
              @click="handleSendNotification"
              :disabled="sending || !notificationForm.title.trim() || !notificationForm.content.trim()"
            >
              <i class="fa-solid fa-paper-plane"></i>发送通知
            </button>
          </div>
        </div>
      </div>

      <div class="right-panel">
        <div class="layui-card">
          <div class="layui-card-header">
            <i class="fa-solid fa-history"></i>
            <span>通知历史</span>
          </div>

          <div v-if="notifications.length > 0" class="notification-list">
            <div v-for="notification in notifications" :key="notification.id" class="notification-item">
              <div class="notification-avatar" :class="getTypeColor(notification.type)">
                <i class="fa-solid" :class="getTypeIcon(notification.type)"></i>
              </div>
              <div class="notification-content">
                <div class="notification-title">
                  {{ notification.title }}
                  <span class="type-tag" :class="getTypeColor(notification.type)">{{ notification.type }}</span>
                </div>
                <div class="notification-body">{{ notification.content }}</div>
                <div class="notification-meta">
                  <span class="target-tag">{{ notification.target }}</span>
                  <span class="time-text">{{ formatTime(notification.created_at) }}</span>
                </div>
              </div>
              <button class="delete-btn" @click="deleteNotification(notification.id)">
                <i class="fa-solid fa-trash"></i>
              </button>
            </div>
          </div>

          <div v-else class="empty-state">
            <i class="fa-solid fa-inbox"></i>
            <div class="empty-text">暂无通知记录</div>
            <div class="empty-hint">发送的通知将显示在这里</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { adminNotificationApi } from '../../api/admin'
import { confirm, success, error } from '../../utils/modal'

const notifications = ref([])
const sending = ref(false)

const notificationForm = ref({
  type: 'system',
  target: 'all',
  title: '',
  content: ''
})

const notificationTypes = [
  { value: 'system', label: '系统通知' },
  { value: 'activity', label: '活动通知' },
  { value: 'warning', label: '警告通知' }
]

const notificationTargets = [
  { value: 'all', label: '全部用户' },
  { value: 'admin', label: '管理员' }
]

const getTypeColor = (type) => {
  const colors = {
    system: 'primary',
    activity: 'success',
    warning: 'warning'
  }
  return colors[type] || 'grey'
}

const getTypeIcon = (type) => {
  const icons = {
    system: 'fa-circle-info',
    activity: 'fa-calendar',
    warning: 'fa-triangle-exclamation'
  }
  return icons[type] || 'fa-bell'
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

const loadNotifications = async () => {
  try {
    const response = await adminNotificationApi.getNotifications()
    notifications.value = response.data.notifications || []
  } catch (err) {
    console.error('加载通知列表失败', err)
    error('加载通知列表失败')
  }
}

const handleSendNotification = async () => {
  if (!notificationForm.value.title || !notificationForm.value.content) {
    error('请填写通知标题和内容')
    return
  }

  sending.value = true
  try {
    await adminNotificationApi.createNotification(notificationForm.value)
    success('发送成功')
    notificationForm.value = { type: 'system', target: 'all', title: '', content: '' }
    loadNotifications()
  } catch (err) {
    console.error('发送通知失败', err)
    error(err.response?.data?.error || '发送失败')
  } finally {
    sending.value = false
  }
}

const deleteNotification = async (id) => {
  const confirmed = await confirm('确定要删除此通知吗？')
  if (!confirmed) return

  try {
    await adminNotificationApi.deleteNotification(id)
    success('删除成功')
    loadNotifications()
  } catch (err) {
    console.error('删除通知失败', err)
    error(err.response?.data?.error || '删除失败')
  }
}

onMounted(() => {
  loadNotifications()
})
</script>

<style scoped>
.admin-notifications {
  padding: 20px;
}

.notifications-container {
  display: flex;
  gap: 24px;
  max-width: 1200px;
}

.left-panel {
  width: 45%;
}

.right-panel {
  flex: 1;
}

.layui-card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
  font-weight: 600;
}

.form-group {
  margin-bottom: 20px;
}

.form-label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin-bottom: 8px;
}

.char-count {
  display: block;
  text-align: right;
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.layui-card-footer {
  display: flex;
  justify-content: flex-end;
  padding: 16px 20px;
}

.notification-list {
  padding: 8px 0;
}

.notification-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;

  &:last-child {
    border-bottom: none;
  }
}

.notification-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  flex-shrink: 0;

  &.primary {
    background: rgba(30, 159, 255, 0.1);
    color: #1E9FFF;
  }

  &.success {
    background: rgba(82, 196, 26, 0.1);
    color: #52C41A;
  }

  &.warning {
    background: rgba(250, 173, 20, 0.1);
    color: #FAAD14;
  }
}

.notification-content {
  flex: 1;
  min-width: 0;
}

.notification-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 15px;
  font-weight: 500;
  color: #333;
  margin-bottom: 6px;
}

.type-tag {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;

  &.primary {
    background: rgba(30, 159, 255, 0.1);
    color: #1E9FFF;
  }

  &.success {
    background: rgba(82, 196, 26, 0.1);
    color: #52C41A;
  }

  &.warning {
    background: rgba(250, 173, 20, 0.1);
    color: #FAAD14;
  }
}

.notification-body {
  font-size: 14px;
  color: #666;
  line-height: 1.5;
  margin-bottom: 8px;
}

.notification-meta {
  display: flex;
  align-items: center;
  gap: 12px;
}

.target-tag {
  padding: 2px 8px;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  font-size: 12px;
  color: #666;
}

.time-text {
  font-size: 12px;
  color: #999;
}

.delete-btn {
  padding: 4px;
  background: transparent;
  border: none;
  color: #999;
  font-size: 16px;
  cursor: pointer;
  flex-shrink: 0;

  &:hover {
    color: #FF5722;
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
  margin-bottom: 8px;
}

.empty-hint {
  font-size: 13px;
  color: #bbb;
}

@media (max-width: 768px) {
  .notifications-container {
    flex-direction: column;
  }

  .left-panel {
    width: 100%;
  }
}
</style>