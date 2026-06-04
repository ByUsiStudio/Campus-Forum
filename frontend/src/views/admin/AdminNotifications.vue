<template>
  <v-container fluid class="pa-6">
    <v-row>
      <v-col cols="12" lg="5">
        <v-card class="pa-6">
          <v-card-title class="d-flex align-center mb-6">
            <v-avatar color="primary" size="44" class="mr-3">
              <v-icon color="white">mdi-bell</v-icon>
            </v-avatar>
            <div>
              <div class="text-h5 font-weight-bold">发送通知</div>
              <div class="text-caption text-medium-emphasis">向用户推送系统通知</div>
            </div>
          </v-card-title>

          <v-divider class="mb-6"></v-divider>

          <v-form ref="notificationFormRef" v-model="formValid">
            <v-select
              v-model="notificationForm.type"
              :items="notificationTypes"
              item-title="label"
              item-value="value"
              label="通知类型"
              variant="outlined"
              density="comfortable"
              prepend-inner-icon="mdi-alert"
              class="mb-4"
            >
              <template #label>
                <span class="text-body-2">通知类型</span>
              </template>
            </v-select>

            <v-select
              v-model="notificationForm.target"
              :items="notificationTargets"
              item-title="label"
              item-value="value"
              label="发送对象"
              variant="outlined"
              density="comfortable"
              prepend-inner-icon="mdi-account-group"
              class="mb-4"
            >
              <template #label>
                <span class="text-body-2">发送对象</span>
              </template>
            </v-select>

            <v-text-field
              v-model="notificationForm.title"
              label="通知标题"
              placeholder="请输入通知标题"
              variant="outlined"
              density="comfortable"
              :rules="[rules.required]"
              prepend-inner-icon="mdi-format-title"
              clearable
              class="mb-4"
            >
              <template #label>
                <span class="text-body-2">通知标题</span>
              </template>
            </v-text-field>

            <v-textarea
              v-model="notificationForm.content"
              label="通知内容"
              placeholder="请输入通知内容..."
              variant="outlined"
              density="comfortable"
              :rules="[rules.required]"
              prepend-inner-icon="mdi-text"
              rows="4"
              counter
              :maxlength="500"
              class="mb-4"
            >
              <template #label>
                <span class="text-body-2">通知内容</span>
              </template>
            </v-textarea>
          </v-form>

          <v-divider class="my-6"></v-divider>

          <v-card-actions class="px-0">
            <v-spacer></v-spacer>
            <v-btn
              color="primary"
              variant="flat"
              size="large"
              @click="handleSendNotification"
              :loading="sending"
              :disabled="!formValid"
            >
              <v-icon class="mr-2">mdi-send</v-icon>
              发送通知
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>

      <v-col cols="12" lg="7">
        <v-card class="pa-6">
          <v-card-title class="d-flex align-center mb-4">
            <v-avatar color="secondary" size="44" class="mr-3">
              <v-icon color="white">mdi-history</v-icon>
            </v-avatar>
            <div>
              <div class="text-h5 font-weight-bold">通知历史</div>
              <div class="text-caption text-medium-emphasis">查看已发送的通知记录</div>
            </div>
          </v-card-title>

          <v-divider class="mb-6"></v-divider>

          <div v-if="notifications.length > 0">
            <v-list lines="two" class="pa-0">
              <v-list-item
                v-for="notification in notifications"
                :key="notification.id"
                class="notification-item pa-4 mb-3"
              >
                <template #prepend>
                  <v-avatar :color="getTypeColor(notification.type)" size="40" class="mr-3">
                    <v-icon color="white" size="20">{{ getTypeIcon(notification.type) }}</v-icon>
                  </v-avatar>
                </template>

                <v-list-item-title class="font-weight-bold mb-1">
                  {{ notification.title }}
                </v-list-item-title>
                <v-list-item-subtitle class="text-body-2 mb-2" style="white-space: pre-wrap;">
                  {{ notification.content }}
                </v-list-item-subtitle>
                <v-list-item-subtitle class="text-caption text-medium-emphasis">
                  <v-chip
                    size="x-small"
                    :color="getTypeColor(notification.type)"
                    variant="tonal"
                    class="mr-2"
                  >
                    {{ notification.type }}
                  </v-chip>
                  <v-chip size="x-small" variant="outlined" class="mr-2">
                    {{ notification.target }}
                  </v-chip>
                  <span class="ml-1">{{ formatTime(notification.created_at) }}</span>
                </v-list-item-subtitle>

                <template #append>
                  <v-btn
                    variant="text"
                    color="error"
                    size="small"
                    icon="mdi-delete"
                    @click="deleteNotification(notification.id)"
                  ></v-btn>
                </template>
              </v-list-item>
            </v-list>
          </div>

          <div v-else class="empty-state">
            <v-icon size="64" color="grey lighten-3">mdi-inbox-outline</v-icon>
            <p class="mt-4 text-body-1 text-medium-emphasis">暂无通知记录</p>
            <p class="text-caption text-grey">发送的通知将显示在这里</p>
          </div>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../../api'
import { confirm, success, error } from '../../utils/modal'

const notifications = ref([])
const notificationFormRef = ref(null)
const formValid = ref(false)
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

const rules = {
  required: v => !!v || '此字段为必填项'
}

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
    system: 'mdi-information',
    activity: 'mdi-calendar',
    warning: 'mdi-alert'
  }
  return icons[type] || 'mdi-bell'
}

const formatTime = (timeString) => {
  if (!timeString) return ''
  const date = new Date(timeString)
  return date.toLocaleString('zh-CN')
}

const loadNotifications = async () => {
  try {
    const response = await api.get('/notifications/admin')
    notifications.value = response.data.notifications || []
  } catch (err) {
    console.error('加载通知列表失败', err)
  }
}

const handleSendNotification = async () => {
  if (!notificationForm.value.title || !notificationForm.value.content) {
    error('请填写通知标题和内容')
    return
  }

  sending.value = true
  try {
    await api.post('/notifications', notificationForm.value)
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
    await api.delete(`/notifications/${id}`)
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
.notification-item {
  border: 1px solid rgba(148, 163, 184, 0.2);
  border-radius: 12px;
  transition: all 0.2s ease;
}

.notification-item:hover {
  border-color: rgba(148, 163, 184, 0.4);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 300px;
  padding: 40px;
}

:deep(.v-field) {
  border-radius: 12px;
}

:deep(.v-field--outlined .v-field__outline) {
  border-color: rgba(148, 163, 184, 0.3);
}

:deep(.v-field--focused .v-field__outline) {
  border-color: rgb(var(--v-theme-primary));
}

:deep(.v-textarea .v-field__input) {
  min-height: 100px;
}
</style>