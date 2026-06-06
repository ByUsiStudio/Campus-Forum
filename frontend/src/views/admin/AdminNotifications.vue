<template>
  <v-container fluid class="pa-4 pa-md-6">
    <v-row>
      <!-- 发送通知表单 -->
      <v-col cols="12" lg="5">
        <v-card variant="flat" rounded="lg">
          <v-card-title class="pb-2">
            <v-icon start size="20">mdi-bell</v-icon>
            发送通知
          </v-card-title>

          <v-card-text>
            <v-form ref="notificationFormRef" v-model="formValid">
              <v-select
                v-model="notificationForm.type"
                :items="notificationTypes"
                item-title="label"
                item-value="value"
                label="通知类型"
                variant="outlined"
                density="compact"
                prepend-inner-icon="mdi-alert"
                class="mb-3"
              />

              <v-select
                v-model="notificationForm.target"
                :items="notificationTargets"
                item-title="label"
                item-value="value"
                label="发送对象"
                variant="outlined"
                density="compact"
                prepend-inner-icon="mdi-account-group"
                class="mb-3"
              />

              <v-text-field
                v-model="notificationForm.title"
                label="通知标题"
                placeholder="请输入通知标题"
                variant="outlined"
                density="compact"
                :rules="[rules.required]"
                prepend-inner-icon="mdi-format-title"
                clearable
                class="mb-3"
              />

              <v-textarea
                v-model="notificationForm.content"
                label="通知内容"
                placeholder="请输入通知内容..."
                variant="outlined"
                density="compact"
                :rules="[rules.required]"
                prepend-inner-icon="mdi-text"
                rows="4"
                counter
                :maxlength="500"
              />
            </v-form>
          </v-card-text>

          <v-card-actions class="pa-4">
            <v-spacer />
            <v-btn
              color="primary"
              variant="flat"
              @click="handleSendNotification"
              :loading="sending"
              :disabled="!formValid"
            >
              <v-icon start>mdi-send</v-icon>
              发送通知
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>

      <!-- 通知历史列表 -->
      <v-col cols="12" lg="7">
        <v-card variant="flat" rounded="lg">
          <v-card-title class="pb-2">
            <v-icon start size="20">mdi-history</v-icon>
            通知历史
          </v-card-title>

          <v-list lines="two" v-if="notifications.length > 0">
            <v-list-item v-for="notification in notifications" :key="notification.id" class="py-3">
              <template v-slot:prepend>
                <v-avatar size="48" :color="getTypeColor(notification.type)" variant="tonal">
                  <v-icon>{{ getTypeIcon(notification.type) }}</v-icon>
                </v-avatar>
              </template>

              <v-list-item-title class="font-weight-medium mb-1">
                {{ notification.title }}
                <v-chip size="x-small" :color="getTypeColor(notification.type)" variant="tonal" class="ml-2">
                  {{ notification.type }}
                </v-chip>
              </v-list-item-title>

              <v-list-item-subtitle>
                <div class="mb-1">{{ notification.content }}</div>
                <div class="d-flex flex-wrap align-center ga-2">
                  <v-chip size="x-small" variant="outlined">
                    {{ notification.target }}
                  </v-chip>
                  <span class="text-caption text-medium-emphasis">
                    {{ formatTime(notification.created_at) }}
                  </span>
                </div>
              </v-list-item-subtitle>

              <template v-slot:append>
                <v-btn size="small" color="error" variant="text" @click="deleteNotification(notification.id)">
                  <v-icon>mdi-delete</v-icon>
                  <v-tooltip activator="parent">删除</v-tooltip>
                </v-btn>
              </template>
            </v-list-item>
          </v-list>

          <v-card-text v-else class="text-center py-8">
            <v-icon size="48" color="grey-lighten-1">mdi-inbox-outline</v-icon>
            <div class="text-body-1 text-medium-emphasis mt-2">
              暂无通知记录
            </div>
            <div class="text-caption text-medium-emphasis mt-1">
              发送的通知将显示在这里
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { adminNotificationApi } from '../../api/admin'
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