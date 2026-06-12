<script setup>
import { ref, inject, onMounted, onUnmounted, nextTick, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getMessages, sendMessage, markConversationRead } from '../api/chat'
import { getFriendList } from '../api/follow'

const router = useRouter()
const route = useRoute()
const user = inject('user')
const isMobile = computed(() => window.innerWidth < 1024)

const messages = ref([])
const newMessage = ref('')
const conversationId = ref(route.params.id)
const targetUser = ref(null)
const messagesContainer = ref(null)

const loadMessages = async () => {
  try {
    const response = await getMessages(conversationId.value)
    messages.value = response.data.messages || []
    await nextTick(() => {
      scrollToBottom()
    })
  } catch (error) {
    console.error('加载消息失败:', error)
  }
}

const loadFriendInfo = async () => {
  try {
    const response = await getFriendList()
    const friends = response.data.friends || []
    targetUser.value = friends.find(f => f.id === parseInt(conversationId.value))
  } catch (error) {
    console.error('加载好友信息失败:', error)
  }
}

const handleSend = async () => {
  if (!newMessage.value.trim()) return
  try {
    await sendMessage(conversationId.value, newMessage.value)
    newMessage.value = ''
    loadMessages()
  } catch (error) {
    console.error('发送消息失败:', error)
  }
}

const scrollToBottom = () => {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

const handleBack = () => {
  markConversationRead(conversationId.value)
  router.push('/chat')
}

let refreshInterval = null

onMounted(() => {
  if (!user.value) {
    router.push('/login')
    return
  }
  loadMessages()
  loadFriendInfo()
  markConversationRead(conversationId.value)
  
  refreshInterval = setInterval(() => {
    loadMessages()
  }, 2000)
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})
</script>

<template>
  <v-container class="py-6">
    <!-- 移动端返回按钮和标题 -->
    <div v-if="isMobile" class="mb-4 flex items-center">
      <v-btn text @click="handleBack">
        <v-icon>mdi-arrow-left</v-icon>
      </v-btn>
      <h2 class="text-h5 ml-2">{{ targetUser?.display_name || targetUser?.friend?.username || '聊天' }}</h2>
    </div>

    <!-- PC端标题 -->
    <v-card-title v-else class="text-h5 mb-4">
      {{ targetUser?.display_name || targetUser?.friend?.username || '聊天' }}
    </v-card-title>

    <!-- 聊天区域 -->
    <v-card class="h-full" style="overflow: hidden;">
      <div 
        ref="messagesContainer" 
        class="h-full overflow-y-auto p-4"
        style="max-height: calc(100vh - 220px);"
      >
        <div v-if="messages.length === 0" class="text-center py-12">
          <v-icon size="64" color="grey">mdi-message-circle</v-icon>
          <p class="mt-4 text-grey">暂无消息，开始聊天吧</p>
        </div>

        <div v-for="message in messages" :key="message.id" class="mb-4">
          <div 
            class="flex"
            :class="message.sender_id === user?.id ? 'justify-end' : 'justify-start'"
          >
            <div 
              class="max-w-xs px-4 py-2 rounded-lg"
              :class="message.sender_id === user?.id 
                ? 'bg-primary text-white rounded-br-sm' 
                : 'bg-grey-light text-black rounded-bl-sm'"
            >
              <p>{{ message.content }}</p>
              <p class="text-xs mt-1 opacity-60">{{ message.created_at }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- 输入框 -->
      <v-card-actions class="p-4">
        <v-text-field
          v-model="newMessage"
          placeholder="输入消息..."
          outlined
          class="flex-1"
          @keyup.enter="handleSend"
        ></v-text-field>
        <v-btn color="primary" @click="handleSend">
          <v-icon>mdi-send</v-icon>
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-container>
</template>

<style scoped>
.bg-grey-light {
  background-color: #f5f5f5;
}
</style>
