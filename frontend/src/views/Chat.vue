<template>
  <div class="chat-wrapper">
    <v-card class="chat-container">
      <v-card-title class="chat-header d-flex align-center gap-2">
        <v-btn icon variant="text" size="small" @click="goBack">
          <v-icon>mdi-arrow-left</v-icon>
        </v-btn>
        <div style="flex: 1; min-width: 0;">
          <UserAvatar :user="currentUser" :size="32" />
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

      <v-card-text class="chat-messages">
        <div
          v-for="message in messages"
          :key="message.id"
          class="message-wrapper"
          :class="{ 'justify-end': isOwnMessage(message) }"
        >
          <template v-if="!isOwnMessage(message)">
            <UserAvatar :user="currentUser" :size="32" class="message-avatar" />
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
            <UserAvatar :user="myUser" :size="32" class="message-avatar" />
          </template>
        </div>

        <div v-if="!canSendMore && !followStatus.mutual" class="text-center text-warning mt-4">
          <v-icon size="small">mdi-alert-circle</v-icon>
          <span class="text-caption">未互相关注，仅能发送2条消息</span>
        </div>
      </v-card-text>

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
import UserAvatar from '../components/UserAvatar.vue'

export default {
  name: 'Chat',
  components: {
    UserAvatar
  },
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

    const myUser = ref(JSON.parse(localStorage.getItem('user') || '{}'))

    const otherUserId = parseInt(route.params.id)
    const token = ref(localStorage.getItem('token'))

    const connectWebSocket = () => {
      const user = JSON.parse(localStorage.getItem('user') || '{}')
      if (!user.id || !token.value) return

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

            if (data.error) {
              console.error('WebSocket server error:', data.error)
              return
            }

            const isFromOther = data.sender_id === otherUserId && data.receiver_id === user.id
            const isToOther = data.sender_id === user.id && data.receiver_id === otherUserId

            if (isFromOther || isToOther) {
              const exists = messages.value.some(m => m.id === data.id)
              if (!exists) {
                console.log('Adding new message:', data)
                messages.value.push(data)
                scrollToBottom()

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
          setTimeout(connectWebSocket, 3000)
        }

        websocket.value.onerror = (error) => {
          console.error('WebSocket connection error:', error)
        }
      } catch (error) {
        console.error('Failed to create WebSocket:', error)
        setTimeout(connectWebSocket, 3000)
      }
    }

    const loadUser = async () => {
      try {
        const response = await api.get(`/follow/status/${otherUserId}`)

        followStatus.value = {
          is_following: response.data.is_following,
          is_followed: response.data.is_followed,
          mutual: response.data.mutual
        }

        console.log('关注状态:', followStatus.value)

        if (response.data.following_user) {
          currentUser.value = response.data.following_user
        } else {
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
        const response = await api.post(`/chat/messages/${otherUserId}`)
        messages.value = response.data.messages || []
        scrollToBottom()
      } catch (error) {
        console.error('加载消息失败', error)
      }
    }

    const checkMessageLimit = async () => {
      try {
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
      myUser
    }
  }
}
</script>

<style scoped>
.chat-wrapper {
  height: calc(100vh - 64px);
  display: flex;
  flex-direction: column;
  margin: 0;
  padding: 0;
  overflow: hidden;
}

.chat-container {
  height: 100%;
  display: flex;
  flex-direction: column;
  border-radius: 0;
  overflow: hidden;
}

.chat-header {
  flex-shrink: 0;
  padding: 8px 12px;
  border-bottom: 1px solid #e0e0e0;
  display: flex;
  align-items: center;
  min-height: 56px;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  background-color: #f5f5f5;
  padding: 12px;
  display: flex;
  flex-direction: column;
}

.message-wrapper {
  display: flex;
  align-items: flex-start;
  margin-bottom: 10px;
}

.message-wrapper.justify-end {
  justify-content: flex-end;
}

.message-avatar {
  flex-shrink: 0;
  margin-right: 8px;
}

.message-wrapper.justify-end .message-avatar {
  margin-right: 0;
  margin-left: 8px;
}

.message-content {
  max-width: 70%;
  min-width: 60px;
}

.message-card {
  border-radius: 12px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
  word-break: break-word;
}

.message-wrapper.justify-end .message-card {
  border-bottom-right-radius: 4px;
}

.message-wrapper:not(.justify-end) .message-card {
  border-bottom-left-radius: 4px;
}

.chat-input {
  flex-shrink: 0;
  padding: 10px 12px;
  border-top: 1px solid #e0e0e0;
  display: flex;
  align-items: center;
}

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

@media (min-width: 601px) and (max-width: 768px) {
  .chat-wrapper {
    height: calc(100vh - 64px);
    margin: 0;
  }
}

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