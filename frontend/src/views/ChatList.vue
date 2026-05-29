<template>
  <v-card class="h-screen flex flex-col">
    <v-card-title class="text-h5">消息</v-card-title>

    <v-card-text class="flex-1 overflow-auto pa-0">
      <div v-if="sessions.length === 0" class="text-center pa-12 text-medium-emphasis">
        <v-icon size="64" color="grey-lighten-1">mdi-message-circle-off</v-icon>
        <div class="text-h6 mt-4">暂无消息</div>
        <div class="text-body-2">关注他人或等待他人关注你开始聊天</div>
      </div>

      <v-list v-else>
        <v-list-item
          v-for="session in sessions"
          :key="session.session_id"
          @click="goToChat(session.other_user.id)"
          class="cursor-pointer"
        >
          <template v-slot:prepend>
            <UserAvatar :user="session.other_user" :size="40" />
          </template>

          <v-list-item-title class="font-weight-bold">
            {{ session.other_user.display_name }}
          </v-list-item-title>
          <v-list-item-subtitle class="text-caption text-medium-emphasis">
            {{ session.last_message || '暂无消息' }}
          </v-list-item-subtitle>

          <template v-slot:append>
            <div class="text-right">
              <div class="text-caption text-medium-emphasis">
                {{ formatTime(session.last_message_at) }}
              </div>
              <v-badge
                v-if="session.unread_count > 0"
                :content="session.unread_count"
                color="error"
              ></v-badge>
            </div>
          </template>
        </v-list-item>
      </v-list>
    </v-card-text>
  </v-card>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '../api'
import UserAvatar from '../components/UserAvatar.vue'

export default {
  name: 'ChatList',
  components: {
    UserAvatar
  },
  setup() {
    const router = useRouter()
    const sessions = ref([])

    const loadSessions = async () => {
      try {
        const response = await api.get('/chat/sessions')
        sessions.value = response.data.sessions || []
      } catch (error) {
        console.error('加载会话失败', error)
      }
    }

    const goToChat = (userId) => {
      router.push(`/chat/${userId}`)
    }

    const formatTime = (date) => {
      if (!date) return ''
      const d = new Date(date)
      const now = new Date()
      const diff = now - d
      const hours = Math.floor(diff / 3600000)
      const days = Math.floor(diff / 86400000)

      if (hours < 1) return '刚刚'
      if (hours < 24) return `${hours}小时前`
      if (days < 7) return `${days}天前`
      return d.toLocaleDateString('zh-CN')
    }

    onMounted(() => {
      loadSessions()
    })

    return {
      sessions,
      goToChat,
      formatTime
    }
  }
}
</script>