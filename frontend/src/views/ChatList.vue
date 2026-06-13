<template>
  <v-container class="pa-4">
    <!-- 页面标题 -->
    <div class="d-flex align-center mb-4">
      <v-btn icon variant="text" @click="$router.back()" class="mr-2">
        <v-icon>mdi-arrow-left</v-icon>
      </v-btn>
      <h2 class="text-h5 font-weight-bold">我的消息</h2>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="d-flex justify-center align-center py-10">
      <v-progress-circular indeterminate color="primary"></v-progress-circular>
    </div>

    <!-- IM未连接提示 -->
    <v-card v-else-if="!isConnected" class="pa-6 text-center mb-4">
      <v-icon size="48" color="warning" class="mb-4">mdi-cloud-off-outline</v-icon>
      <div class="text-body-1 mb-4">IM服务未连接</div>
      <v-btn color="primary" @click="initIM" :loading="initing">
        重新连接
      </v-btn>
    </v-card>

    <template v-else>
      <!-- 会话列表 -->
      <v-card v-if="conversations.length">
        <v-list lines="two" class="py-0">
          <template v-for="(conv, index) in conversations" :key="conv.conversation.targetId">
            <v-list-item
              @click="openConversation(conv)"
              class="conversation-item"
            >
              <template v-slot:prepend>
                <div class="avatar-wrapper mr-3">
                  <v-avatar size="48" color="primary">
                    <v-img
                      v-if="conv.conversation.type === 'single' && getUserAvatar(conv.conversation.targetId)"
                      :src="getUserAvatar(conv.conversation.targetId)"
                    ></v-img>
                    <v-icon v-else color="white">
                      {{ conv.conversation.type === 'group' ? 'mdi-account-group' : 'mdi-account' }}
                    </v-icon>
                  </v-avatar>
                  <v-badge
                    v-if="conv.unreadCount > 0"
                    :content="conv.unreadCount > 99 ? '99+' : conv.unreadCount"
                    color="error"
                    class="unread-badge"
                  ></v-badge>
                </div>
              </template>

              <v-list-item-title class="font-weight-bold mb-1">
                {{ getConversationName(conv) }}
              </v-list-item-title>

              <v-list-item-subtitle class="text-truncate">
                {{ getLastMessage(conv) }}
              </v-list-item-subtitle>

              <template v-slot:append>
                <div class="text-right">
                  <div class="text-caption text-grey">
                    {{ formatTime(conv.lastMessage?.timestamp) }}
                  </div>
                </div>
              </template>
            </v-list-item>

            <v-divider v-if="index < conversations.length - 1"></v-divider>
          </template>
        </v-list>
      </v-card>

      <!-- 空状态 -->
      <v-card v-else class="pa-6 text-center">
        <v-icon size="64" color="grey-lighten-1" class="mb-4">mdi-message-text-outline</v-icon>
        <div class="text-body-1 text-grey mb-4">暂无会话消息</div>
        <v-btn color="primary" to="/" variant="tonal">
          去首页看看
        </v-btn>
      </v-card>
    </template>

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
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useGoIM } from '@/composables/useGoIM'
import axios from 'axios'

const router = useRouter()

const {
  isConnected,
  connectWebSocket,
  onMessage,
  getOnlineUsers,
} = useGoIM()

const loading = ref(true)
const initing = ref(false)
const conversations = ref([])
const currentUserId = ref(parseInt(localStorage.getItem('userId') || '0'))

// 获取API基础URL
const getApiBase = () => import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

const error = reactive({
  show: false,
  text: ''
})

// 初始化IM
const initIM = async () => {
  initing.value = true
  try {
    if (currentUserId.value) {
      connectWebSocket(currentUserId.value)
    }
    await loadConversations()
  } catch (err) {
    error.show = true
    error.text = err.message || '初始化失败'
  } finally {
    initing.value = false
  }
}

// 加载会话列表
const loadConversations = async () => {
  try {
    const apiBase = getApiBase()
    // 获取当前用户的所有聊天会话
    const response = await axios.get(`${apiBase}/api/conversations`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    conversations.value = response.data.data || []
  } catch (err) {
    console.error('加载会话列表失败:', err)
    // 即使失败也继续，显示空列表
    conversations.value = []
  }
}

// 获取会话名称
const getConversationName = (conv) => {
  if (conv.conversation?.type === 'group') {
    return conv.conversation.targetId || '群聊'
  }
  // 单聊显示对方用户名
  return conv.conversation?.targetId || '未知用户'
}

// 获取最后一条消息
const getLastMessage = (conv) => {
  if (!conv.lastMessage) return '暂无消息'

  const { MessageContentType } = conv.lastMessage

  // 文本消息
  if (conv.lastMessage.content?.content) {
    return conv.lastMessage.content.content
  }

  // 其他类型消息
  const typeMap = {
    1: '[图片]',
    2: '[语音]',
    3: '[视频]',
    4: '[文件]',
    5: '[位置]',
    100: '[自定义消息]',
  }
  return typeMap[MessageContentType] || '[消息]'
}

// 格式化时间
const formatTime = (timestamp) => {
  if (!timestamp) return ''

  const date = new Date(timestamp)
  const now = new Date()
  const diff = now - date

  // 超过7天显示日期
  if (diff > 7 * 24 * 60 * 60 * 1000) {
    return `${date.getMonth() + 1}/${date.getDate()}`
  }

  // 超过24小时显示星期
  if (diff > 24 * 60 * 60 * 1000) {
    const weekdays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
    return weekdays[date.getDay()]
  }

  // 一天内显示时间
  return `${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`
}

// 获取用户头像（这里可以关联用户系统获取真实头像）
const getUserAvatar = (userId) => {
  // TODO: 可以通过用户API获取真实头像
  return null
}

// 打开会话
const openConversation = (conv) => {
  const type = conv.conversation.type === 'group' ? 'group' : 'private'
  router.push(`/chat/${type}/${conv.conversation.targetId}`)
}

let removeMessageListener = null

onMounted(async () => {
  loading.value = true
  await initIM()

  // 监听新消息
  removeMessageListener = onMessage((message) => {
    // 更新会话列表或显示新消息通知
    loadConversations()
  })

  loading.value = false
})

onUnmounted(() => {
  if (removeMessageListener) {
    removeMessageListener()
  }
})
</script>

<style scoped>
.conversation-item {
  cursor: pointer;
  transition: background-color 0.2s;
}

.conversation-item:hover {
  background-color: rgba(0, 0, 0, 0.04);
}

.avatar-wrapper {
  position: relative;
}

.unread-badge {
  position: absolute;
  top: -4px;
  right: -4px;
}
</style>
