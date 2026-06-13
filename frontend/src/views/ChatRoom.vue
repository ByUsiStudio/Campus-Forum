<template>
  <v-container fluid class="pa-0 fill-height d-flex flex-column">
    <!-- 顶部导航 -->
    <div class="chat-header pa-3 d-flex align-center">
      <v-btn icon variant="text" @click="$router.back()" class="mr-2">
        <v-icon>mdi-arrow-left</v-icon>
      </v-btn>

      <v-avatar v-if="conversation?.target_user" size="40" class="mr-3">
        <v-img
          v-if="conversation.target_user.avatar"
          :src="conversation.target_user.avatar"
        ></v-img>
        <v-icon v-else color="white">mdi-account</v-icon>
      </v-avatar>

      <div>
        <div class="text-body-1 font-weight-bold">
          {{ conversation?.target_user?.display_name || conversation?.target_user?.username || '聊天' }}
        </div>
        <div class="text-caption text-grey">
          <v-icon size="12" class="mr-1">mdi-circle</v-icon>
          在线
        </div>
      </div>

      <v-spacer></v-spacer>

      <v-btn icon variant="text" @click="showConversationInfo = true">
        <v-icon>mdi-information-outline</v-icon>
      </v-btn>
    </div>

    <!-- 消息区域 -->
    <div ref="messagesContainer" class="messages-container flex-grow-1 overflow-y-auto pa-3">
      <!-- 加载更多 -->
      <div v-if="loadingMore" class="text-center py-2">
        <v-progress-circular indeterminate size="20" color="primary"></v-progress-circular>
      </div>

      <!-- 加载状态 -->
      <div v-if="loading && !messages.length" class="d-flex justify-center align-center py-10">
        <v-progress-circular indeterminate color="primary"></v-progress-circular>
      </div>

      <!-- 消息列表 -->
      <div class="messages-list">
        <template v-for="(msg, index) in messages" :key="msg.id">
          <!-- 日期分割线 -->
          <div v-if="showDateDivider(index)" class="date-divider text-center my-4">
            <span class="date-text">{{ formatDate(msg.created_at) }}</span>
          </div>

          <!-- 消息气泡 -->
          <div
            class="message-wrapper d-flex mb-3"
            :class="msg.sender_id === currentUserId ? 'justify-end' : 'justify-start'"
          >
            <!-- 发送者头像 -->
            <v-avatar
              v-if="msg.sender_id !== currentUserId"
              size="36"
              class="mr-2 flex-shrink-0"
            >
              <v-img
                v-if="getSenderAvatar(msg.sender_id)"
                :src="getSenderAvatar(msg.sender_id)"
              ></v-img>
              <v-icon v-else color="grey">mdi-account</v-icon>
            </v-avatar>

            <!-- 消息内容 -->
            <div
              class="message-bubble"
              :class="msg.sender_id === currentUserId ? 'my-message' : 'other-message'"
            >
              <div class="message-content">{{ msg.content }}</div>
              <div class="message-time text-caption">
                {{ formatTime(msg.created_at) }}
                <v-icon
                  v-if="msg.sender_id === currentUserId"
                  size="12"
                  :color="msg.status >= 1 ? 'primary' : 'grey'"
                >
                  {{ msg.status >= 2 ? 'mdi-check-all' : 'mdi-check' }}
                </v-icon>
              </div>
            </div>

            <!-- 我的头像 -->
            <v-avatar
              v-if="msg.sender_id === currentUserId"
              size="36"
              class="ml-2 flex-shrink-0"
            >
              <v-img
                v-if="currentUser?.avatar"
                :src="currentUser.avatar"
              ></v-img>
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
          @keyup.enter="sendMessage"
          class="flex-grow-1"
        ></v-text-field>

        <v-btn
          icon
          color="primary"
          :disabled="!newMessage.trim()"
          :loading="sending"
          @click="sendMessage"
        >
          <v-icon>mdi-send</v-icon>
        </v-btn>
      </div>
    </div>

    <!-- 会话信息侧滑 -->
    <v-navigation-drawer
      v-model="showConversationInfo"
      location="right"
      temporary
      width="280"
    >
      <v-list>
        <v-list-subheader>会话信息</v-list-subheader>

        <v-list-item>
          <v-list-item-title class="text-body-2 text-grey">会话类型</v-list-item-title>
          <template v-slot:append>
            <v-list-item-title>{{ conversation?.type === 'group' ? '群聊' : '私聊' }}</v-list-item-title>
          </template>
        </v-list-item>

        <v-list-item>
          <v-list-item-title class="text-body-2 text-grey">对方</v-list-item-title>
          <template v-slot:append>
            <v-list-item-title>{{ conversation?.target_user?.display_name }}</v-list-item-title>
          </template>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

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
import { ref, reactive, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useRoute } from 'vue-router'
import { chatApi } from '@/api'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const userStore = useUserStore()

const conversationId = computed(() => parseInt(route.params.id))

const loading = ref(true)
const loadingMore = ref(false)
const sending = ref(false)
const messages = ref([])
const newMessage = ref('')
const showConversationInfo = ref(false)

const conversation = ref({
  id: null,
  type: 'private',
  target_user: null
})

const currentUser = computed(() => userStore.user)
const currentUserId = computed(() => userStore.user?.id)

const error = reactive({
  show: false,
  text: ''
})

let ws = null
let reconnectTimer = null
let heartbeatTimer = null

// WebSocket连接
const connectWebSocket = () => {
  const token = localStorage.getItem('token')
  if (!token) return

  try {
    ws = new WebSocket(`ws://localhost:3620/ws?token=${token}`)

    ws.onopen = () => {
      console.log('WebSocket连接成功')
      startHeartbeat()
    }

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        if (data.type === 'chat' && data.conversation_id === conversationId.value) {
          // 添加新消息
          const exists = messages.value.find(m => m.id === data.message_id)
          if (!exists) {
            messages.value.push({
              id: data.message_id,
              conversation_id: data.conversation_id,
              sender_id: data.sender_id,
              content: data.content,
              type: data.type || 'text',
              status: 1,
              created_at: data.created_at || new Date().toISOString()
            })
            scrollToBottom()
          }
        }
      } catch (err) {
        console.error('解析WebSocket消息失败:', err)
      }
    }

    ws.onerror = (err) => {
      console.error('WebSocket错误:', err)
    }

    ws.onclose = () => {
      console.log('WebSocket连接关闭')
      stopHeartbeat()
      // 尝试重连
      if (!reconnectTimer) {
        reconnectTimer = setTimeout(() => {
          reconnectTimer = null
          connectWebSocket()
        }, 3000)
      }
    }
  } catch (err) {
    console.error('建立WebSocket连接失败:', err)
  }
}

// 开始心跳
const startHeartbeat = () => {
  stopHeartbeat()
  heartbeatTimer = setInterval(() => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({ type: 'ping' }))
    }
  }, 30000)
}

// 停止心跳
const stopHeartbeat = () => {
  if (heartbeatTimer) {
    clearInterval(heartbeatTimer)
    heartbeatTimer = null
  }
}

// 获取会话信息
const fetchConversation = async () => {
  try {
    const res = await chatApi.getConversations()
    const conv = res.data?.conversations?.find(c => c.id === conversationId.value)
    if (conv) {
      conversation.value = conv
    }
  } catch (err) {
    console.error('获取会话信息失败:', err)
  }
}

// 获取消息历史
const fetchMessages = async (loadMore = false) => {
  if (loadMore) {
    loadingMore.value = true
  } else {
    loading.value = true
  }

  try {
    const res = await chatApi.getMessages({
      conversation_id: conversationId.value,
      limit: 20,
      offset: loadMore ? messages.value.length : 0
    })

    if (res.data?.messages) {
      if (loadMore) {
        messages.value = [...res.data.messages.reverse(), ...messages.value]
      } else {
        messages.value = res.data.messages.reverse()
      }
    }
  } catch (err) {
    console.error('获取消息历史失败:', err)
    error.text = '加载消息失败'
    error.show = true
  } finally {
    loading.value = false
    loadingMore.value = false
    if (!loadMore) {
      nextTick(() => scrollToBottom())
    }
  }
}

// 发送消息
const sendMessage = async () => {
  const content = newMessage.value.trim()
  if (!content || sending.value) return

  sending.value = true
  newMessage.value = ''

  // 构建消息对象
  const tempId = Date.now()
  const tempMsg = {
    id: tempId,
    conversation_id: conversationId.value,
    sender_id: currentUserId.value,
    content: content,
    type: 'text',
    status: 0,
    created_at: new Date().toISOString()
  }

  try {
    // 先添加到本地
    messages.value.push(tempMsg)
    scrollToBottom()

    // 尝试WebSocket发送
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({
        type: 'chat',
        data: {
          conversation_id: conversationId.value,
          content: content,
          type: 'text'
        }
      }))
      tempMsg.status = 1
    } else {
      // 降级到HTTP
      const res = await chatApi.sendMessage({
        conversation_id: conversationId.value,
        content: content,
        type: 'text'
      })
      if (res.data?.message) {
        // 更新消息ID
        const idx = messages.value.findIndex(m => m.id === tempId)
        if (idx !== -1) {
          messages.value[idx] = res.data.message
        }
      }
    }
  } catch (err) {
    console.error('发送消息失败:', err)
    error.text = '发送消息失败'
    error.show = true
    // 移除临时消息
    messages.value = messages.value.filter(m => m.id !== tempId)
    newMessage.value = content
  } finally {
    sending.value = false
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
  const current = new Date(messages.value[index].created_at)
  const previous = new Date(messages.value[index - 1].created_at)
  return current.toDateString() !== previous.toDateString()
}

// 格式化日期
const formatDate = (timeStr) => {
  const date = new Date(timeStr)
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
const formatTime = (timeStr) => {
  const date = new Date(timeStr)
  return `${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`
}

// 获取发送者头像
const getSenderAvatar = (senderId) => {
  if (senderId === currentUserId.value) {
    return currentUser.value?.avatar
  }
  return conversation.value?.target_user?.avatar
}

const messagesContainer = ref(null)

// 监听消息变化，自动滚动
watch(messages, () => {
  scrollToBottom()
}, { deep: true })

onMounted(async () => {
  await fetchConversation()
  await fetchMessages()
  connectWebSocket()
})

onUnmounted(() => {
  if (ws) {
    ws.close()
    ws = null
  }
  stopHeartbeat()
  if (reconnectTimer) {
    clearTimeout(reconnectTimer)
    reconnectTimer = null
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
