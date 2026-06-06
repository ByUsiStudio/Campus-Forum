<template>
  <v-container fluid class="fill-height pa-0">
    <v-row no-gutters class="fill-height">
      <v-col cols="12" class="fill-height">
        <v-card class="fill-height d-flex flex-column" elevation="0">
          <v-card-item class="chat-header pa-2">
            <div class="d-flex align-center gap-2">
              <v-btn icon variant="text" size="small" @click="goBack">
                <v-icon>mdi-arrow-left</v-icon>
              </v-btn>
              <UserAvatar :user="currentUser" :size="36" />
              <div class="text-body-1 font-weight-bold flex-grow-1">{{ currentUser?.display_name }}</div>
              <v-btn
                v-if="!followStatus.is_following"
                variant="tonal"
                color="primary"
                size="small"
                @click="handleFollow"
              >
                <v-icon start size="small">mdi-plus</v-icon>
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
            </div>
          </v-card-item>

          <v-divider></v-divider>

          <v-card-text class="chat-messages flex-grow-1" ref="messagesContainer">
            <div
              v-for="message in sortedMessages"
              :key="message.id"
              class="message-wrapper"
              :class="{ 'justify-end': isOwnMessage(message) }"
            >
              <template v-if="!isOwnMessage(message)">
                <UserAvatar :user="currentUser" :size="36" class="message-avatar" />
                <div class="message-content">
                  <div class="text-caption text-medium-emphasis mb-1">{{ message.sender_name }}</div>
                  <v-card class="bg-white message-card" variant="elevated">
                    <v-card-text class="pa-2 pa-sm-3 text-body-2">
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
                  <v-card class="bg-primary" variant="elevated">
                    <v-card-text class="pa-2 pa-sm-3 text-body-2 text-white">
                      {{ message.content }}
                      <div class="text-caption mt-1 text-white opacity-70">
                        {{ formatTime(message.created_at) }}
                      </div>
                    </v-card-text>
                  </v-card>
                </div>
                <UserAvatar :user="myUser" :size="36" class="message-avatar" />
              </template>
            </div>

            <div v-if="!canSendMore && !followStatus.mutual" class="text-center mt-4">
              <v-alert type="warning" variant="tonal" density="compact">
                <v-icon size="small" class="mr-1">mdi-alert-circle</v-icon>
                未互相关注，仅能发送2条消息
              </v-alert>
            </div>
          </v-card-text>

          <v-divider></v-divider>

          <v-card-actions class="chat-input pa-2">
            <v-text-field
              v-model="messageInput"
              placeholder="输入消息..."
              variant="outlined"
              density="compact"
              hide-details
              class="flex-grow-1"
              @keydown.enter="sendMessage"
              :disabled="!canSendMore"
              maxlength="500"
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
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '../api'
import { success, error } from '../utils/modal'
import UserAvatar from '../components/UserAvatar.vue'

export default {
  name: 'StandaloneChat',
  components: {
    UserAvatar
  },
  setup() {
    const route = useRoute()
    const router = useRouter()

    const otherUserId = parseInt(route.params.id)
    const token = ref(localStorage.getItem('token'))

    const messages = ref([])
    const messageInput = ref('')
    const currentUser = ref(null)
    const followStatus = ref({
      is_following: false,
      is_followed: false,
      mutual: false
    })
    const canSendMore = ref(true)
    const websocket = ref(null)
    const messagesContainer = ref(null)

    const myUser = ref(JSON.parse(localStorage.getItem('user') || '{}'))

    const messageTotalCount = ref(0)
    const lastMessageId = ref(0)

    const sortedMessages = computed(() => {
      return [...messages.value].sort((a, b) => a.id - b.id)
    })

    const isOwnMessage = (message) => {
      const user = JSON.parse(localStorage.getItem('user') || '{}')
      return message.sender_id === user.id
    }

    const formatTime = (time) => {
      const date = new Date(time)
      return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
    }

    const scrollToBottom = () => {
      setTimeout(() => {
        if (messagesContainer.value) {
          messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
        }
      }, 100)
    }

    const connectWebSocket = () => {
      const user = JSON.parse(localStorage.getItem('user') || '{}')
      if (!user.id || !token.value) return

      const wsUrl = `/ws/chat?user_id=${user.id}`

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
              if (data.total_count) {
                messageTotalCount.value = data.total_count
              }
              if (data.last_message_id) {
                lastMessageId.value = data.last_message_id
              }

              const exists = messages.value.some(m => m.id === data.id)
              if (!exists) {
                messages.value.push(data)
                scrollToBottom()

                if (data.sender_id === otherUserId) {
                  checkMessageLimit()
                }
              }

              validateMessageIntegrity()
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

    const validateMessageIntegrity = async () => {
      if (messages.value.length !== messageTotalCount.value) {
        console.log(`消息数量不匹配: 本地${messages.value.length}, 服务器${messageTotalCount.value}`)
        await loadMessages()
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

        messageTotalCount.value = messages.value.length
        if (messages.value.length > 0) {
          lastMessageId.value = messages.value[messages.value.length - 1].id
        }

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
      } catch (err) {
        console.error('发送消息失败', err)
        messages.value = messages.value.filter(m => m.id !== localMessage.id)
        error('发送失败，请重试')
      }
    }

    const handleFollow = async () => {
      try {
        await api.post(`/follow/${otherUserId}`)
        success('关注成功')
        await loadUser()
        await checkMessageLimit()
      } catch (err) {
        console.error('关注失败', err)
        error('关注失败，请重试')
      }
    }

    const handleUnfollow = async () => {
      try {
        await api.delete(`/follow/${otherUserId}`)
        success('取消关注成功')
        await loadUser()
      } catch (err) {
        console.error('取消关注失败', err)
        error('取消关注失败，请重试')
      }
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

    return {
      currentUser,
      messages,
      sortedMessages,
      messageInput,
      followStatus,
      canSendMore,
      sendMessage,
      handleFollow,
      handleUnfollow,
      isOwnMessage,
      formatTime,
      goBack,
      myUser,
      messagesContainer
    }
  }
}
</script>

<style scoped>
.standalone-chat {
  height: 100vh;
  width: 100%;
  display: flex;
  flex-direction: column;
  margin: 0;
  padding: 0;
  overflow: hidden;
  background-color: #f5f5f5;
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
  .standalone-chat {
    max-width: 800px;
    margin: 0 auto;
  }

  .chat-container {
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }
}

@media (max-width: 600px) {
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