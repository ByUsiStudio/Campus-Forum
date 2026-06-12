<script setup>
import { ref, onMounted } from 'vue'
import { adminApi } from '../../api'

const loading = ref(false)
const sending = ref(false)
const form = ref({
  user_id: null,
  title: '',
  content: ''
})
const batchForm = ref({
  title: '',
  content: ''
})
const batchDialog = ref(false)

const sendNotification = async () => {
  if (!form.value.user_id || !form.value.title || !form.value.content) return
  
  sending.value = true
  try {
    await adminApi.sendUserNotification(form.value)
    form.value = { user_id: null, title: '', content: '' }
  } catch (error) {
    console.error('发送通知失败:', error)
  } finally {
    sending.value = false
  }
}

const openBatchDialog = () => {
  batchForm.value = { title: '', content: '' }
  batchDialog.value = true
}

const sendBatchNotifications = async () => {
  if (!batchForm.value.title || !batchForm.value.content) return
  
  sending.value = true
  try {
    await userNotificationApi.sendBatchNotifications(batchForm.value)
    batchDialog.value = false
  } catch (error) {
    console.error('批量发送通知失败:', error)
  } finally {
    sending.value = false
  }
}

onMounted(() => {
  // 初始化
})
</script>

<template>
  <v-container fluid class="pa-0">
    <!-- 页面标题 -->
    <div class="mb-6">
      <h1 class="text-h5 font-weight-bold">用户通知</h1>
      <p class="text-body-2 text-grey">发送通知给用户</p>
    </div>

    <!-- 发送单个通知 -->
    <v-card class="mb-4">
      <v-card-title>发送单个通知</v-card-title>
      <v-card-text>
        <v-text-field v-model="form.user_id" label="用户ID" type="number" class="mb-3" />
        <v-text-field v-model="form.title" label="通知标题" class="mb-3" />
        <v-textarea v-model="form.content" label="通知内容" rows="3" class="mb-3" />
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="primary" @click="sendNotification" :loading="sending">
          <v-icon start>mdi-send</v-icon>
          发送通知
        </v-btn>
      </v-card-actions>
    </v-card>

    <!-- 批量发送 -->
    <v-card>
      <v-card-title>批量发送通知</v-card-title>
      <v-card-text>
        <p class="text-grey mb-4">向所有用户发送通知</p>
        <v-btn color="warning" @click="openBatchDialog">
          <v-icon start>mdi-send-all</v-icon>
          批量发送通知
        </v-btn>
      </v-card-text>
    </v-card>

    <!-- 批量发送对话框 -->
    <v-dialog v-model="batchDialog" max-width="500">
      <v-card>
        <v-card-title>批量发送通知</v-card-title>
        <v-card-text>
          <v-text-field v-model="batchForm.title" label="通知标题" class="mb-3" />
          <v-textarea v-model="batchForm.content" label="通知内容" rows="5" class="mb-3" />
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn text @click="batchDialog = false">取消</v-btn>
          <v-btn color="warning" @click="sendBatchNotifications" :loading="sending">
            批量发送
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>