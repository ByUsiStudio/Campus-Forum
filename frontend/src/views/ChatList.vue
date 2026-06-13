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

    <template v-else>
      <!-- 会话列表 -->
      <v-card v-if="conversations.length">
        <v-list lines="two" class="py-0">
          <template v-for="(conv, index) in conversations" :key="conv.id">
            <v-list-item
              :to="`/chat/${conv.id}`"
              class="conversation-item"
              @click="markAsRead(conv)"
            >
              <template v-slot:prepend>
                <div class="avatar-wrapper mr-3">
                  <v-avatar size="48" :color="getAvatarColor(conv.type)">
                    <v-img
                      v-if="conv.target_user?.avatar"
                      :src="conv.target_user.avatar"
                    ></v-img>
                    <v-icon v-else color="white">
                      {{ conv.type === 'group' ? 'mdi-account-group' : 'mdi-account' }}
                    </v-icon>
                  </v-avatar>
                  <v-badge
                    v-if="conv.unread_count > 0"
                    :content="conv.unread_count > 99 ? '99+' : conv.unread_count"
                    color="error"
                    class="unread-badge"
                  ></v-badge>
                </div>
              </template>

              <v-list-item-title class="font-weight-bold mb-1">
                {{ getConversationName(conv) }}
              </v-list-item-title>

              <v-list-item-subtitle class="text-truncate">
                {{ conv.last_message?.content || '暂无消息' }}
              </v-list-item-subtitle>

              <template v-slot:append>
                <div class="text-right">
                  <div class="text-caption text-grey">
                    {{ formatTime(conv.last_msg_time) }}
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
import { ref, reactive, onMounted } from 'vue'
import { chatApi } from '@/api'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()

const loading = ref(true)
const conversations = ref([])

const error = reactive({
  show: false,
  text: ''
})

// 获取会话列表
const fetchConversations = async () => {
  try {
    const res = await chatApi.getConversations()
    if (res.data?.conversations) {
      conversations.value = res.data.conversations
    }
  } catch (err) {
    console.error('获取会话列表失败:', err)
    error.text = '加载会话列表失败'
    error.show = true
  } finally {
    loading.value = false
  }
}

// 标记会话已读
const markAsRead = async (conv) => {
  if (conv.unread_count > 0) {
    try {
      await chatApi.markConversationRead(conv.id)
    } catch (err) {
      console.error('标记已读失败:', err)
    }
  }
}

// 获取会话名称
const getConversationName = (conv) => {
  if (conv.type === 'group') {
    return conv.group_name || '群聊'
  }
  return conv.target_user?.display_name || conv.target_user?.username || '未知用户'
}

// 获取头像颜色
const getAvatarColor = (type) => {
  return type === 'group' ? 'purple' : 'primary'
}

// 格式化时间
const formatTime = (timeStr) => {
  if (!timeStr) return ''

  const date = new Date(timeStr)
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

onMounted(() => {
  fetchConversations()
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
