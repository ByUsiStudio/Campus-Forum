<template>
  <div class="page-container">
    <el-row :gutter="16">
      <!-- 发送通知表单 -->
      <el-col :xs="24" :lg="9">
        <el-card shadow="never" class="mb-4">
          <template #header>
            <div class="card-header">
              <el-icon class="mr-1"><Bell /></el-icon>
              <span>发送通知</span>
            </div>
          </template>

          <el-form
            ref="notificationFormRef"
            :model="notificationForm"
            :rules="formRules"
            label-position="top"
          >
            <el-form-item label="通知类型" prop="type">
              <el-select v-model="notificationForm.type" class="w-full">
                <el-option
                  v-for="item in notificationTypes"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                />
              </el-select>
            </el-form-item>

            <el-form-item label="发送对象" prop="target">
              <el-select v-model="notificationForm.target" class="w-full">
                <el-option
                  v-for="item in notificationTargets"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                />
              </el-select>
            </el-form-item>

            <el-form-item label="通知标题" prop="title">
              <el-input
                v-model="notificationForm.title"
                placeholder="请输入通知标题"
                clearable
              />
            </el-form-item>

            <el-form-item label="通知内容" prop="content">
              <el-input
                v-model="notificationForm.content"
                type="textarea"
                placeholder="请输入通知内容..."
                :rows="4"
                :maxlength="500"
                show-word-limit
              />
            </el-form-item>

            <el-form-item>
              <div class="w-full flex justify-end">
                <el-button
                  type="primary"
                  :loading="sending"
                  @click="handleSendNotification"
                >
                  <el-icon class="mr-1"><Promotion /></el-icon>
                  发送通知
                </el-button>
              </div>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>

      <!-- 通知历史列表 -->
      <el-col :xs="24" :lg="15">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <el-icon class="mr-1"><Clock /></el-icon>
              <span>通知历史</span>
            </div>
          </template>

          <template v-if="notifications.length > 0">
            <div
              v-for="notification in notifications"
              :key="notification.id"
              class="notification-item"
            >
              <div class="flex items-start gap-3">
                <div class="flex-shrink-0">
                  <el-avatar :size="48" :style="{ backgroundColor: getTypeColor(notification.type), color: '#fff' }">
                    <el-icon><component :is="getTypeIcon(notification.type)" /></el-icon>
                  </el-avatar>
                </div>

                <div class="min-w-0 flex-1">
                  <div class="flex items-center gap-2 mb-1">
                    <span class="font-medium">{{ notification.title }}</span>
                    <el-tag
                      type="info"
                      size="small"
                      :style="{ color: getTypeColor(notification.type), borderColor: getTypeColor(notification.type), backgroundColor: 'transparent' }"
                    >
                      {{ notification.type }}
                    </el-tag>
                  </div>
                  <div class="mb-1 text-body">{{ notification.content }}</div>
                  <div class="flex items-center gap-2">
                    <el-tag size="small" effect="plain">{{ notification.target }}</el-tag>
                    <span class="text-caption">{{ formatTime(notification.created_at) }}</span>
                  </div>
                </div>

                <el-tooltip content="删除" placement="top">
                  <el-button
                    type="danger"
                    text
                    size="small"
                    @click="deleteNotification(notification.id)"
                  >
                    <el-icon><Delete /></el-icon>
                  </el-button>
                </el-tooltip>
              </div>
            </div>
          </template>

          <el-empty
            v-else
            description="暂无通知记录"
            :image-size="80"
          >
            <p class="empty-hint">发送的通知将显示在这里</p>
          </el-empty>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import {
  Bell,
  Promotion,
  Clock,
  InfoFilled,
  Calendar,
  Warning,
  Delete
} from '@element-plus/icons-vue'
import { adminNotificationApi } from '@/api'
import { confirm, success, error } from '@/utils/message'

// ---- 通知历史 ----
const notifications = ref([])
const notificationFormRef = ref(null)
const sending = ref(false)

// ---- 发送表单 ----
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

const formRules = {
  title: [{ required: true, message: '请输入通知标题', trigger: 'blur' }],
  content: [{ required: true, message: '请输入通知内容', trigger: 'blur' }]
}

const typeStyles = {
  primary: '#409eff',
  success: '#67c23a',
  warning: '#e6a23c',
  grey: '#909399'
}

const getTypeColor = (type) => {
  const colors = {
    system: 'primary',
    activity: 'success',
    warning: 'warning'
  }
  const name = colors[type] || 'grey'
  return typeStyles[name]
}

const getTypeIcon = (type) => {
  const icons = {
    system: InfoFilled,
    activity: Calendar,
    warning: Warning
  }
  return icons[type] || Bell
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
  let valid = false
  if (notificationFormRef.value) {
    valid = await notificationFormRef.value.validate().then(() => true).catch(() => false)
  } else {
    valid = !!(notificationForm.value.title && notificationForm.value.content)
  }
  if (!valid) {
    error('请填写通知标题和内容')
    return
  }

  sending.value = true
  try {
    await adminNotificationApi.createNotification(notificationForm.value)
    success('发送成功')
    notificationForm.value = { type: 'system', target: 'all', title: '', content: '' }
    if (notificationFormRef.value) notificationFormRef.value.clearValidate()
    await loadNotifications()
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
    await loadNotifications()
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
.card-header {
  display: flex;
  align-items: center;
  font-weight: 600;
}
.notification-item {
  display: flex;
  padding: 16px 4px;
  border-bottom: 1px solid var(--el-border-color-lighter);
}
.notification-item:last-child {
  border-bottom: none;
}
.font-medium {
  font-weight: 500;
}
.text-body {
  font-size: 14px;
  color: var(--el-text-color-primary);
}
.text-caption {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}
.empty-hint {
  margin: 4px 0 0;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}
.w-full {
  width: 100%;
}
.mb-1 {
  margin-bottom: 4px;
}
.mb-4 {
  margin-bottom: 16px;
}
</style>
