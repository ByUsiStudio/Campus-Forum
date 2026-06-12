<script setup>
import { ref, inject, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { getConversations, getUnreadCount } from '../api/chat'
import { getFriendList } from '../api/follow'

const router = useRouter()
const user = inject('user')
const conversations = ref([])
const friends = ref([])
const unreadCount = ref(0)
const searchQuery = ref('')

const loadConversations = async () => {
  try {
    const response = await getConversations()
    conversations.value = response.data.conversations || []
  } catch (error) {
    console.error('加载会话列表失败:', error)
  }
}

const loadUnreadCount = async () => {
  try {
    const response = await getUnreadCount()
    unreadCount.value = response.data.unread_count || 0
  } catch (error) {
    console.error('加载未读消息数失败:', error)
  }
}

const loadFriends = async () => {
  try {
    const response = await getFriendList()
    friends.value = response.data.friends || []
  } catch (error) {
    console.error('加载好友列表失败:', error)
  }
}

const goToChatRoom = (conversationId) => {
  router.push(`/chat/${conversationId}`)
}

const filteredConversations = () => {
  if (!searchQuery.value) return conversations.value
  const query = searchQuery.value.toLowerCase()
  return conversations.value.filter(c => 
    (c.target_user && c.target_user.username && c.target_user.username.toLowerCase().includes(query)) ||
    (c.target_user && c.target_user.display_name && c.target_user.display_name.toLowerCase().includes(query))
  )
}

let refreshInterval = null

onMounted(() => {
  if (!user.value) {
    router.push('/login')
    return
  }
  loadConversations()
  loadUnreadCount()
  loadFriends()
  
  refreshInterval = setInterval(() => {
    loadUnreadCount()
  }, 3000)
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})
</script>

<template>
  <v-app>
    <v-app-bar app>
      <v-btn icon @click="router.push('/profile')">
        <v-icon>mdi-arrow-left</v-icon>
      </v-btn>
      <v-toolbar-title>聊天</v-toolbar-title>
      <v-btn icon @click="loadConversations">
        <v-icon>mdi-refresh</v-icon>
      </v-btn>
    </v-app-bar>

    <v-container class="py-6">
      <!-- 搜索框 -->
      <v-card class="mb-4">
        <v-text-field
          v-model="searchQuery"
          placeholder="搜索好友或会话..."
          prepend-inner-icon="mdi-magnify"
          outlined
        ></v-text-field>
      </v-card>

      <!-- 会话列表 -->
      <v-card v-if="filteredConversations().length > 0">
        <v-list>
          <v-list-item 
            v-for="conv in filteredConversations()" 
            :key="conv.id"
            @click="goToChatRoom(conv.id)"
            class="cursor-pointer"
          >
            <v-list-item-avatar>
              <v-icon color="primary">mdi-account</v-icon>
            </v-list-item-avatar>
            <v-list-item-content>
              <v-list-item-title>
                {{ conv.target_user?.display_name || conv.target_user?.username }}
                <span v-if="conv.unread_count > 0" class="ml-2">
                  <span class="badge badge-pill badge-danger">{{ conv.unread_count }}</span>
                </span>
              </v-list-item-title>
              <v-list-item-subtitle>
                {{ conv.last_message?.content || '暂无消息' }}
              </v-list-item-subtitle>
            </v-list-item-content>
            <v-list-item-actions>
              <span class="text-grey text-sm">{{ conv.last_msg_time }}</span>
            </v-list-item-actions>
          </v-list-item>
        </v-list>
      </v-card>

      <!-- 空状态 -->
      <v-card v-else class="text-center py-12">
        <v-icon size="64" color="grey">mdi-message-circle</v-icon>
        <p class="mt-4 text-grey">暂无会话</p>
        <p class="text-grey text-sm">点击下方按钮开始聊天</p>
      </v-card>

      <!-- 好友列表（快速发起聊天） -->
      <v-card class="mt-6">
        <v-card-title>好友</v-card-title>
        <v-list>
          <v-list-item 
            v-for="friend in friends" 
            :key="friend.id"
            @click="goToChatRoom(friend.id)"
            class="cursor-pointer"
          >
            <v-list-item-avatar>
              <v-icon color="primary">mdi-account</v-icon>
            </v-list-item-avatar>
            <v-list-item-content>
              <v-list-item-title>{{ friend.display_name || friend.friend?.username }}</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
        </v-list>
      </v-card>
    </v-container>
  </v-app>
</template>
