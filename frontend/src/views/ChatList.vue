<template>
  <v-container fluid class="fill-height pa-0">
    <v-row no-gutters class="fill-height">
      <v-col cols="12" class="fill-height">
        <v-card class="fill-height d-flex flex-column" elevation="0">
          <v-card-item class="pa-4">
            <div class="d-flex align-center">
              <v-btn icon="mdi-arrow-left" variant="text" @click="router.back()" class="mr-2"></v-btn>
              <v-card-title class="text-h5">消息中心</v-card-title>
            </div>
          </v-card-item>
          
          <v-divider></v-divider>

          <v-card-text class="flex-grow-1 overflow-auto pa-0">
            <div v-if="sessions.length === 0" class="text-center pa-12">
              <v-icon size="80" color="grey-lighten-2">mdi-message-circle-outline</v-icon>
              <div class="text-h6 text-medium-emphasis mt-4">暂无消息</div>
              <div class="text-body-2 text-medium-emphasis mt-2">
                关注他人或等待他人关注你开始聊天
              </div>
            </div>

            <v-list v-else lines="two" class="pa-2">
              <v-list-item
                v-for="session in sessions"
                :key="session.session_id"
                @click="goToChat(session.other_user.id)"
                class="mb-2 rounded-lg"
                color="primary"
              >
                <template v-slot:prepend>
                  <v-badge
                    v-if="session.unread_count > 0"
                    :content="session.unread_count"
                    color="error"
                    offset-x="-5"
                    offset-y="-5"
                  >
                    <UserAvatar :user="session.other_user" :size="50" />
                  </v-badge>
                  <UserAvatar v-else :user="session.other_user" :size="50" />
                </template>

                <v-list-item-title class="font-weight-bold mb-1">
                  {{ session.other_user.display_name }}
                </v-list-item-title>
                
                <v-list-item-subtitle class="text-body-2 text-medium-emphasis d-flex align-center gap-2">
                  <span class="text-truncate">{{ session.last_message || '暂无消息' }}</span>
                </v-list-item-subtitle>

                <template v-slot:append>
                  <div class="d-flex flex-column align-end gap-1">
                    <span class="text-caption text-medium-emphasis">
                      {{ formatTime(session.last_message_at) }}
                    </span>
                  </div>
                </template>
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
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