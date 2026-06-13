<template>
  <v-container fluid class="pa-0 fill-height d-flex flex-column">
    <!-- 顶部导航 -->
    <div class="chat-header pa-3 d-flex align-center">
      <v-btn icon variant="text" @click="$router.back()" class="mr-2">
        <v-icon>mdi-arrow-left</v-icon>
      </v-btn>

      <v-avatar size="40" class="mr-3">
        <v-img v-if="targetAvatar" :src="targetAvatar"></v-img>
        <v-icon v-else color="primary">mdi-account</v-icon>
      </v-avatar>

      <div>
        <div class="text-body-1 font-weight-bold">{{ targetName }}</div>
        <div class="text-caption text-grey">
          <v-icon size="12" class="mr-1">mdi-circle</v-icon>
          {{ isConnected ? '在线' : '离线' }}
        </div>
      </div>

      <v-spacer></v-spacer>
    </div>

    <!-- 消息区域 -->
    <div ref="messagesContainer" class="messages-container flex-grow-1 overflow-y-auto pa-3">
      <!-- 加载状态 -->
      <div v-if="loading" class="d-flex justify-center align-center py-10">
        <v-progress-circular indeterminate color="primary"></v-progress-circular>
      </div>

      <!-- 消息列表 -->
      <div v-else class="messages-list">
        <template v-for="(msg, index) in messages" :key="msg.messageId">
          <!-- 日期分割线 -->
          <div v-if="showDateDivider(index)" class="date-divider text-center my-4">
            <span class="date-text">{{ formatDate(msg.timestamp) }}</span>
          </div>

          <!-- 消息气泡 -->
          <div
            class="message-wrapper d-flex mb-3"
            :class="isMyMessage(msg) ? 'justify-end' : 'justify-start'"
          >
            <!-- 发送者头像 -->
            <v-avatar
              v-if="!isMyMessage(msg)"
              size="36"
              class="mr-2 flex-shrink-0"
            >
              <v-img v-if="getSenderAvatar(msg.senderId)" :src="getSenderAvatar(msg.senderId)"></v-img>
              <v-icon v-else color="grey">mdi-account</v-icon>
            </v-avatar>

            <!-- 消息内容 -->
            <div
              class="message-bubble"
              :class="isMyMessage(msg) ? 'my-message' : 'other-message'"
            >
              <div class="message-content">{{ getMessageContent(msg) }}</div>
              <div class="message-time text-caption">
                {{ formatTime(msg.timestamp) }}
              </div>
            </div>

            <!-- 我的头像 -->
            <v-avatar
              v-if="isMyMessage(msg)"
              size="36"
              class="ml-2 flex-shrink-0"
            >
              <v-img v-if="currentUser?.avatar" :src="currentUser.avatar"></v-img>
              <v-icon v-else color="grey">mdi-account</v-icon>
            </v-avatar>
          </div>
        </template>
      </div>

      <!-- 空状态 -->
      <div v-if="!loading && !messages.length" class="text-center py-10">
        <v-icon size="64" color="grey-lighten-1" class="mb-4">mdi-message-outline</v-icon>
        <div class="text-body-1 text-grey">开始和对方聊天吧</div>
      </div>
    </div>

    <!-- 输入区域 -->
    <div class="chat-input-area pa-3">
      <div class="d-flex align-end gap-2">
        <v-text-field
          v-model="newMessage"
          variant="outlined"
          density="compact"
          placeholder="输入消息..."
          hide-details
          rounded="pill"
          @keyup.enter="sendTextMessage"
          class="flex-grow-1"
        ></v-text-field>

        <v-btn
          icon
          color="primary"
          :disabled="!newMessage.trim()"
          :loading="sending"
          @click="sendTextMessage"
        >
          <v-icon>mdi-send</v-icon>
        </v-btn>
      </div>
    </div>

    <!-- 错误提示 -->
    <v-snackbar v-model="error.show" color="error" location="top">
      {{ error.text }}
      <template v-slot:actions>
        <v-btn variant="text" @click="error.show = false">关闭</v-btn>
      </template>
    </v-snackbar>
  </v-container>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { useGoIM } from '@/composables/useGoIM'
import axios from 'axios'

const route = useRoute()

const {
  isConnected,
  connectWebSocket,
  sendPrivateMessage,
  sendGroupMessage,
  onMessage,
} = useGoIM()

// 从路由获取会话信息
const conversationType = computed(() => route.params.type === 'group' ? 'group' : 'private')
const targetId = computed(() => route.params.targetId)
const targetName = ref('聊天')
const targetAvatar = ref(null)

const loading = ref(true)
const sending = ref(false)
const messages = ref([])
const newMessage = ref('')
const currentUserId = ref(parseInt(localStorage.getItem('userId') || '0'))

// 获取API基础URL
const getApiBase = () => import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

const currentUser = computed(() => {
  const userStr = localStorage.getItem('user')
  return userStr ? JSON.parse(userStr) : null
})

const error = reactive({
  show: false,
  text: ''
})

// 判断是否是我的消息
const isMyMessage = (msg) => {
  return String(msg.sender_id) === String(currentUserId.value)
}

// 获取消息内容
const getMessageContent = (msg) => {
  if (!msg.content) return msg.content || ''
  return msg.content
}

// 发送文本消息
const sendTextMessage = async () => {
  const content = newMessage.value.trim()
  if (!content || sending.value || !isConnected.value) return

  sending.value = true
  newMessage.value = ''

  try {
    if (conversationType.value === 'group') {
      await sendGroupMessage(targetId.value, content)
    } else {
      await sendPrivateMessage(targetId.value, content)
    }
    
    // 添加发送的消息到列表
    messages.value.push({
      message_id: `temp-${Date.now()}`,
      content: content,
      sender_id: currentUserId.value,
      target_id: targetId.value,
      timestamp: Date.now(),
      status: 'sending'
    })
    scrollToBottom()
  } catch (err) {
    console.error('发送消息失败:', err)
    error.text = '发送消息失败'
    error.show = true
    newMessage.value = content
  } finally {
    sending.value = false
  }
}

// 加载消息历史
const loadMessages = async () => {
  loading.value = true
  try {
    const apiBase = getApiBase()
    const response = await axios.get(`${apiBase}/api/messages/history`, {
      params: {
        conversation_id: conversationType.value === 'group' ? `group_${targetId.value}` : `private_${currentUserId.value}_${targetId.value}`,
        limit: 50
      },
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    messages.value = response.data.data || []
  } catch (err) {
    console.error('加载消息历史失败:', err)
    messages.value = []
  } finally {
    loading.value = false
    nextTick(() => scrollToBottom())
  }
}

// 滚动到底部
const scrollToBottom = () => {
  nextTick(() => {
    const container = messagesContainer.value
    if (container) {
      container.scrollTop = container.scrollHeight
    }
  })
}

// 显示日期分割线
const showDateDivider = (index) => {
  if (index === 0) return true
  const current = new Date(messages.value[index].timestamp)
  const previous = new Date(messages.value[index - 1].timestamp)
  return current.toDateString() !== previous.toDateString()
}

// 格式化日期
const formatDate = (timestamp) => {
  const date = new Date(timestamp)
  const today = new Date()
  const yesterday = new Date(today)
  yesterday.setDate(yesterday.getDate() - 1)

  if (date.toDateString() === today.toDateString()) {
    return '今天'
  }
  if (date.toDateString() === yesterday.toDateString()) {
    return '昨天'
  }
  return `${date.getMonth() + 1}月${date.getDate()}日`
}

// 格式化时间
const formatTime = (timestamp) => {
  const date = new Date(timestamp)
  return `${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`
}

// 获取发送者头像
const getSenderAvatar = (senderId) => {
  if (String(senderId) === String(currentUserId.value)) {
    return currentUser.value?.avatar
  }
  return targetAvatar.value
}

const messagesContainer = ref(null)
let removeMessageListener = null

onMounted(async () => {
  // 设置目标名称
  targetName.value = targetId.value

  // 连接WebSocket
  if (currentUserId.value) {
    connectWebSocket(currentUserId.value)
  }

  // 加载消息
  await loadMessages()

  // 监听新消息
  removeMessageListener = onMessage((message) => {
    // 只处理当前会话的消息
    if (
      message.target_id === targetId.value ||
      (conversationType.value === 'group' && message.type === 'group')
    ) {
      const exists = messages.value.find(m => m.message_id === message.message_id)
      if (!exists) {
        messages.value.push(message)
        scrollToBottom()
      }
    }
  })
})

onUnmounted(() => {
  if (removeMessageListener) {
    removeMessageListener()
  }
})
</script>

<style scoped>
.chat-header {
  background: rgb(var(--v-theme-surface));
  border-bottom: 1px solid rgba(0, 0, 0, 0.08);
}

.messages-container {
  background: #f5f5f5;
}

.messages-list {
  display: flex;
  flex-direction: column;
}

.message-wrapper {
  max-width: 80%;
}

.message-bubble {
  max-width: 100%;
  padding: 10px 14px;
  border-radius: 16px;
  position: relative;
}

.my-message {
  background: rgb(var(--v-theme-primary));
  color: white;
  border-bottom-right-radius: 4px;
}

.other-message {
  background: white;
  color: rgb(var(--v-theme-on-surface));
  border-bottom-left-radius: 4px;
}

.message-content {
  word-break: break-word;
  line-height: 1.4;
}

.message-time {
  display: flex;
  align-items: center;
  gap: 2px;
  margin-top: 4px;
  opacity: 0.7;
}

.my-message .message-time {
  justify-content: flex-end;
}

.date-divider {
  position: relative;
}

.date-text {
  background: #e0e0e0;
  padding: 2px 12px;
  border-radius: 10px;
  font-size: 12px;
  color: #666;
}

.chat-input-area {
  background: rgb(var(--v-theme-surface));
  border-top: 1px solid rgba(0, 0, 0, 0.08);
}
</style>
