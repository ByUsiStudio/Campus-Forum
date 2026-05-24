<template>
  <div class="chat-wrapper">
    <v-card class="chat-container">
      <!-- 头部 -->
      <v-card-title class="chat-header d-flex align-center gap-2">
        <v-btn icon variant="text" size="small" @click="goBack">
          <v-icon>mdi-arrow-left</v-icon>
        </v-btn>
        <div style="flex: 1; min-width: 0;">
          <div class="font-weight-bold text-truncate">{{ currentUser?.display_name }}</div>
          <div class="text-caption text-medium-emphasis">
            {{ followStatus.mutual ? '互相关注' : followStatus.is_following ? '已关注' : '未关注' }}
          </div>
        </div>
        <v-btn
          v-if="!followStatus.is_following"
          variant="outlined"
          color="primary"
          size="small"
          @click="handleFollow"
        >
          关注
        </v-btn>
        <v-btn
          v-else
          variant="text"
          color="secondary"
          size="small"
          @click="handleUnfollow"
        >
          取消关注
        </v-btn>
      </v-card-title>

      <!-- 消息列表 -->
      <v-card-text class="chat-messages">
        <div
          v-for="message in messages"
          :key="message.id"
          class="message-wrapper"
          :class="{ 'justify-end': isOwnMessage(message) }"
        >
          <!-- 对方消息 -->
          <template v-if="!isOwnMessage(message)">
            <v-avatar color="primary" size="32" class="message-avatar">
              <v-img v-if="currentUser?.avatar" :src="currentUser.avatar"></v-img>
              <span v-else>{{ currentUser?.display_name?.[0] || 'U' }}</span>
            </v-avatar>
            <div class="message-content">
              <div class="text-caption text-medium-emphasis mb-1">{{ message.sender_name }}</div>
              <v-card class="bg-white message-card">
                <v-card-text class="pa-2 pa-sm-3">
                  {{ message.content }}
                  <div class="text-caption mt-1 text-medium-emphasis">
                    {{ formatTime(message.created_at) }}
                  </div>
                </v-card-text>
              </v-card>
            </div>
          </template>
          
          <!-- 自己的消息 -->
          <template v-else>
            <div class="message-content">
              <v-card class="bg-primary text-white message-card">
                <v-card-text class="pa-2 pa-sm-3">
                  {{ message.content }}
                  <div class="text-caption mt-1 text-white opacity-70">
                    {{ formatTime(message.created_at) }}
                  </div>
                </v-card-text>
              </v-card>
            </div>
            <v-avatar color="primary" size="32" class="message-avatar">
              <v-img v-if="myAvatar" :src="myAvatar"></v-img>
              <span v-else>{{ myDisplayName?.[0] || 'U' }}</span>
            </v-avatar>
          </template>
        </div>

        <div v-if="!canSendMore && !followStatus.mutual" class="text-center text-warning mt-4">
          <v-icon size="small">mdi-alert-circle</v-icon>
          <span class="text-caption">未互相关注，仅能发送2条消息</span>
        </div>
      </v-card-text>

      <!-- 输入框 -->
      <v-card-actions class="chat-input">
        <v-text-field
          v-model="messageInput"
          label="输入消息..."
          variant="outlined"
          density="compact"
          hide-details
          class="flex-grow-1"
          @keydown.enter="sendMessage"
          :disabled="!canSendMore"
        ></v-text-field>
        <v-btn
          color="primary"
          icon
          @click="sendMessage"
          :disabled="!messageInput.trim() || !canSendMore"
          class="ml-2"
        >
          <v-icon>mdi-send</v-icon>
        </v-btn>
      </v-card-actions>
    </v-card>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted, watch, nextTick, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '../api'

export default {
  name: 'Chat',
  setup() {
    const route = useRoute()
    const router = useRouter()
    const currentUser = ref(null)
    const messages = ref([])
    const messageInput = ref('')
    const followStatus = ref({
      is_following: false,
      is_followed: false,
      mutual: false
    })
    const canSendMore = ref(true)
    const websocket = ref(null)
    
    // 当前用户信息（用于显示自己的头像）
    const myUser = ref(JSON.parse(localStorage.getItem('user') || '{}'))
    const myAvatar = computed(() => myUser.value.avatar)
    const myDisplayName = computed(() => myUser.value.display_name || myUser.value.username)

    const otherUserId = parseInt(route.params.id)
    const token = ref(localStorage.getItem('token'))

    // WebSocket连接
    const connectWebSocket = () => {
      const user = JSON.parse(localStorage.getItem('user') || '{}')
      if (!user.id || !token.value) return
      
      // 直接使用相对路径，nginx已代理 /ws 到后端
      const wsUrl = `/ws/chat?user_id=${user.id}`
      
      console.log('Connecting to WebSocket:', wsUrl)
      
      try {
        websocket.value = new WebSocket(wsUrl)
        
        websocket.value.onopen = () => {
          console.log('WebSocket connected successfully')
        }
        
        websocket.value.onmessage = (event) => {
          try {
            const data = JSON.parse(event.data)
            
            // 如果是错误消息
            if (data.error) {
              console.error('WebSocket server error:', data.error)
              return
            }
            
            // 检查消息是否是当前聊天的
            // data.sender_id 和 data.receiver_id 来自 ChatMessageResponse
            const isFromOther = data.sender_id === otherUserId && data.receiver_id === user.id
            const isToOther = data.sender_id === user.id && data.receiver_id === otherUserId
            
            if (isFromOther || isToOther) {
              // 检查消息是否已经存在（避免重复添加）
              const exists = messages.value.some(m => m.id === data.id)
              if (!exists) {
                console.log('Adding new message:', data)
                messages.value.push(data)
                scrollToBottom()
                
                // 如果是对方发送的消息，更新消息限制
                if (data.sender_id === otherUserId) {
                  checkMessageLimit()
                }
              }
            }
          } catch (error) {
            console.error('解析WebSocket消息失败:', error)
          }
        }
        
        websocket.value.onclose = () => {
          console.log('WebSocket disconnected, attempting to reconnect in 3 seconds...')
          // 3秒后重连
          setTimeout(connectWebSocket, 3000)
        }
        
        websocket.value.onerror = (error) => {
          console.error('WebSocket connection error:', error)
        }
      } catch (error) {
        console.error('Failed to create WebSocket:', error)
        // 3秒后重试
        setTimeout(connectWebSocket, 3000)
      }
    }

    const loadUser = async () => {
      try {
        const response = await api.get(`/follow/status/${otherUserId}`)
        
        // 正确设置关注状态
        followStatus.value = {
          is_following: response.data.is_following,
          is_followed: response.data.is_followed,
          mutual: response.data.mutual
        }
        
        console.log('关注状态:', followStatus.value)
        
        // 尝试从follow状态响应中获取用户信息
        if (response.data.following_user) {
          currentUser.value = response.data.following_user
        } else {
          // 如果API没有返回用户信息，创建一个基本的user对象
          currentUser.value = {
            id: otherUserId,
            display_name: `用户${otherUserId}`,
            username: `user${otherUserId}`,
            avatar: '',
            role: 'user',
            signature: ''
          }
        }
      } catch (error) {
        console.error('加载用户信息失败', error)
        // 出错时也创建基本的user对象
        currentUser.value = {
          id: otherUserId,
          display_name: `用户${otherUserId}`,
          username: `user${otherUserId}`,
          avatar: '',
          role: 'user',
          signature: ''
        }
      }
    }

    const loadMessages = async () => {
      try {
        // 使用POST拉取服务器聊天记录
        const response = await api.post(`/chat/messages/${otherUserId}`)
        messages.value = response.data.messages || []
        scrollToBottom()
      } catch (error) {
        console.error('加载消息失败', error)
      }
    }

    const checkMessageLimit = async () => {
      try {
        // 统计已发送消息数量
        const response = await api.post(`/chat/messages/${otherUserId}`)
        const myMessages = response.data.messages.filter(m => m.sender_id !== otherUserId)
        canSendMore.value = myMessages.length < 2 || followStatus.value.mutual
      } catch (error) {
        console.error('检查消息限制失败', error)
      }
    }

    const sendMessage = async () => {
      if (!messageInput.value.trim() || !canSendMore.value) return

      const currentUserInfo = JSON.parse(localStorage.getItem('user') || '{}')
      
      // 本地回显消息
      const localMessage = {
        id: Date.now(),
        sender_id: currentUserInfo.id,
        receiver_id: otherUserId,
        content: messageInput.value,
        created_at: new Date().toISOString(),
        sender_name: currentUserInfo.display_name || currentUserInfo.username
      }
      messages.value.push(localMessage)
      scrollToBottom()

      try {
        await api.post('/chat/send', {
          receiver_id: otherUserId,
          content: messageInput.value
        })

        messageInput.value = ''
        await loadMessages()
        await checkMessageLimit()
      } catch (error) {
        console.error('发送消息失败', error)
        // 移除本地回显的消息
        messages.value.pop()
        alert(error.response?.data?.error || '发送失败')
      }
    }

    const handleFollow = async () => {
      try {
        await api.post(`/follow/${otherUserId}`)
        followStatus.value.is_following = true
        followStatus.value.mutual = followStatus.value.is_followed
        canSendMore.value = followStatus.value.mutual
      } catch (error) {
        console.error('关注失败', error)
      }
    }

    const handleUnfollow = async () => {
      try {
        await api.delete(`/follow/${otherUserId}`)
        followStatus.value.is_following = false
        followStatus.value.mutual = false
        await checkMessageLimit()
      } catch (error) {
        console.error('取消关注失败', error)
      }
    }

    const isOwnMessage = (message) => {
      const user = JSON.parse(localStorage.getItem('user') || '{}')
      return message.sender_id === user.id
    }

    const formatTime = (date) => {
      const d = new Date(date)
      return d.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
    }

    const scrollToBottom = () => {
      nextTick(() => {
        const container = document.querySelector('.chat-messages')
        if (container) {
          container.scrollTop = container.scrollHeight
        }
      })
    }

    const goBack = () => {
      router.push('/chat')
    }

    onMounted(() => {
      loadUser()
      loadMessages()
      checkMessageLimit()
      connectWebSocket()
    })

    onUnmounted(() => {
      if (websocket.value) {
        websocket.value.close()
      }
    })

    watch(followStatus, () => {
      checkMessageLimit()
    })

    return {
      currentUser,
      messages,
      messageInput,
      followStatus,
      canSendMore,
      sendMessage,
      handleFollow,
      handleUnfollow,
      isOwnMessage,
      formatTime,
      goBack,
      myAvatar,
      myDisplayName
    }
  }
}
</script>

<style scoped>
/* 聊天包装器 - 全屏显示，考虑顶部栏高度 */
.chat-wrapper {
  height: calc(100vh - 64px);
  display: flex;
  flex-direction: column;
  margin: 0;
  padding: 0;
  overflow: hidden;
}

/* 聊天容器 - 响应式布局 */
.chat-container {
  height: 100%;
  display: flex;
  flex-direction: column;
  border-radius: 0;
  overflow: hidden;
}

/* 聊天头部 - 响应式 */
.chat-header {
  flex-shrink: 0;
  padding: 8px 12px;
  border-bottom: 1px solid #e0e0e0;
  display: flex;
  align-items: center;
  min-height: 56px;
}

/* 聊天消息区域 - 响应式 */
.chat-messages {
  flex: 1;
  overflow-y: auto;
  background-color: #f5f5f5;
  padding: 12px;
  display: flex;
  flex-direction: column;
}

/* 单个消息行 - 响应式 */
.message-wrapper {
  display: flex;
  align-items: flex-start;
  margin-bottom: 10px;
}

/* 自己的消息 - 右对齐 */
.message-wrapper.justify-end {
  justify-content: flex-end;
}

/* 消息头像 */
.message-avatar {
  flex-shrink: 0;
  margin-right: 8px;
}

.message-wrapper.justify-end .message-avatar {
  margin-right: 0;
  margin-left: 8px;
}

/* 消息内容区域 - 响应式宽度 */
.message-content {
  max-width: 70%;
  min-width: 60px;
}

/* 消息气泡卡片 */
.message-card {
  border-radius: 12px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
  word-break: break-word;
}

/* 自己消息的气泡 - 右下角直角 */
.message-wrapper.justify-end .message-card {
  border-bottom-right-radius: 4px;
}

/* 对方消息的气泡 - 左下角直角 */
.message-wrapper:not(.justify-end) .message-card {
  border-bottom-left-radius: 4px;
}

/* 聊天输入区域 - 响应式 */
.chat-input {
  flex-shrink: 0;
  padding: 10px 12px;
  border-top: 1px solid #e0e0e0;
  display: flex;
  align-items: center;
}

/* 桌面端 - 居中显示 */
@media (min-width: 769px) {
  .chat-wrapper {
    height: calc(100vh - 64px);
    max-width: 800px;
    margin: 0 auto;
  }
  
  .chat-container {
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }
}

/* 平板适配 */
@media (min-width: 601px) and (max-width: 768px) {
  .chat-wrapper {
    height: calc(100vh - 64px);
    margin: 0;
  }
}

/* 移动端适配 */
@media (max-width: 600px) {
  .chat-wrapper {
    height: calc(100vh - 56px);
    margin: 0;
  }
  
  .chat-header {
    padding: 6px 10px;
    min-height: 52px;
  }
  
  .chat-messages {
    padding: 10px;
  }
  
  .message-content {
    max-width: 80%;
  }
  
  .message-avatar {
    width: 28px;
    height: 28px;
    margin-right: 6px;
  }
  
  .message-wrapper.justify-end .message-avatar {
    margin-left: 6px;
  }
  
  .message-wrapper {
    margin-bottom: 8px;
  }
  
  .chat-input {
    padding: 8px 10px;
  }
}

/* 超小屏幕 */
@media (max-width: 400px) {
  .message-content {
    max-width: 85%;
  }
  
  .chat-header {
    padding: 6px 8px;
  }
  
  .chat-messages {
    padding: 8px;
  }
}
</style>
