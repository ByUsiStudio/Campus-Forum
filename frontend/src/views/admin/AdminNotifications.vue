<template>
  <v-container fluid class="pa-6">
    <v-card>
      <v-card-title>
        <v-icon class="mr-2">mdi-bell</v-icon>
        发送通知
      </v-card-title>
      <v-card-text>
        <v-select
          v-model="notificationForm.type"
          :items="notificationTypes"
          item-title="label"
          item-value="value"
          label="通知类型"
          variant="outlined"
          class="mb-4"
        ></v-select>
        
        <v-select
          v-model="notificationForm.target"
          :items="notificationTargets"
          item-title="label"
          item-value="value"
          label="发送对象"
          variant="outlined"
          class="mb-4"
        ></v-select>
        
        <v-text-field
          v-model="notificationForm.title"
          label="通知标题"
          variant="outlined"
          class="mb-4"
        ></v-text-field>
        
        <v-text-field
          v-model="notificationForm.content"
          label="通知内容"
          variant="outlined"
          multiline
          rows="4"
          class="mb-4"
        ></v-text-field>
      </v-card-text>
      <v-card-actions>
        <v-btn color="primary" variant="flat" @click="handleSendNotification">
          发送通知
        </v-btn>
      </v-card-actions>
    </v-card>

    <v-card class="mt-6">
      <v-card-title>
        <v-icon class="mr-2">mdi-history</v-icon>
        通知历史
      </v-card-title>
      <v-card-text>
        <v-list v-if="notifications.length > 0">
          <v-list-item
            v-for="notification in notifications"
            :key="notification.id"
            class="align-center"
          >
            <v-list-item-content>
              <v-list-item-title class="font-weight-bold mb-1">
                {{ notification.title }}
              </v-list-item-title>
              <v-list-item-subtitle class="mb-2">
                {{ notification.content }}
              </v-list-item-subtitle>
              <v-list-item-subtitle class="text-caption text-medium-emphasis">
                {{ notification.type }} | {{ notification.target }} | {{ formatTime(notification.created_at) }}
              </v-list-item-subtitle>
            </v-list-item-content>
            <v-list-item-action>
              <v-btn
                variant="text"
                color="error"
                size="small"
                @click="deleteNotification(notification.id)"
              >
                <v-icon>mdi-delete</v-icon>
              </v-btn>
            </v-list-item-action>
          </v-list-item>
        </v-list>
        
        <div v-else class="text-center text-gray-400 py-8">
          <v-icon size="48" color="grey lighten-3">mdi-inbox</v-icon>
          <p class="mt-2">暂无通知记录</p>
        </div>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../../api'
import { confirm, success, error } from '../../utils/modal'

const notifications = ref([])
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
  { value: 'online', label: '在线用户' },
  { value: 'admin', label: '管理员' }
]

const loadNotifications = async () => {
  try {
    const response = await api.get('/notifications/admin')
    notifications.value = response.data.notifications || []
  } catch (error) {
    console.error('加载通知列表失败', error)
  }
}

const formatTime = (timeString) => {
  if (!timeString) return ''
  const date = new Date(timeString)
  return date.toLocaleString('zh-CN')
}

const handleSendNotification = async () => {
  if (!notificationForm.value.title || !notificationForm.value.content) {
    error('请填写通知标题和内容')
    return
  }
  
  try {
    await api.post('/notifications', notificationForm.value)
    success('发送成功')
    notificationForm.value = { type: 'system', target: 'all', title: '', content: '' }
    loadNotifications()
  } catch (error) {
    console.error('发送通知失败', error)
    error(error.response?.data?.error || '发送失败')
  }
}

const deleteNotification = async (id) => {
  const confirmed = await confirm('确定要删除此通知吗？')
  if (!confirmed) return
  
  try {
    await api.delete(`/notifications/${id}`)
    success('删除成功')
    loadNotifications()
  } catch (error) {
    console.error('删除通知失败', error)
    error(error.response?.data?.error || '删除失败')
  }
}

onMounted(() => {
  loadNotifications()
})
</script>