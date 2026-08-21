<template>
  <div class="follow-list-page">
    <v-card class="pa-6">
      <div class="d-flex align-center mb-4">
        <v-btn icon variant="text" @click="router.back()">
          <v-icon>mdi-arrow-left</v-icon>
        </v-btn>
        <v-card-title class="pa-0">{{ pageTitle }}</v-card-title>
      </div>

      <v-tabs v-model="tab" align-tabs="start" class="mb-4">
        <v-tab value="following">关注 ({{ followingCount }})</v-tab>
        <v-tab value="followers">粉丝 ({{ followersCount }})</v-tab>
      </v-tabs>

      <div v-if="loading" class="text-center pa-8">
        <v-progress-circular indeterminate color="primary"></v-progress-circular>
      </div>

      <div v-else-if="currentList.length === 0" class="text-center pa-8 text-medium-emphasis">
        暂无{{ tab === 'following' ? '关注' : '粉丝' }}
      </div>

      <v-list v-else>
        <v-list-item
          v-for="userItem in currentList"
          :key="userItem.id"
          class="px-0"
        >
          <template v-slot:prepend>
            <UserAvatar :user="userItem" :size="50" />
          </template>

          <v-list-item-title class="font-weight-bold">
            <router-link :to="'/profile?id=' + userItem.id" class="text-decoration-none">
              {{ userItem.display_name || userItem.username }}
            </router-link>
          </v-list-item-title>

          <v-list-item-subtitle class="text-caption">
            @{{ userItem.username }}
            <span v-if="userItem.role === 'admin'" class="text-primary ml-1">管理员</span>
          </v-list-item-subtitle>

          <v-list-item-subtitle v-if="userItem.signature" class="text-truncate">
            {{ userItem.signature }}
          </v-list-item-subtitle>

          <template v-slot:append>
            <v-btn
              v-if="!isOwnProfile && currentUser && currentUser.id !== userItem.id"
              variant="outlined"
              size="small"
              :color="isFollowing(userItem.id) ? 'default' : 'primary'"
              @click="handleFollowToggle(userItem)"
            >
              {{ isFollowing(userItem.id) ? '已关注' : '关注' }}
            </v-btn>
          </template>
        </v-list-item>
      </v-list>
    </v-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import api from '../api'
import UserAvatar from '../components/UserAvatar.vue'
import { success as showSuccess } from '../utils/modal'

const router = useRouter()
const route = useRoute()

const tab = ref(route.query.tab || 'following')
const loading = ref(false)
const following = ref([])
const followers = ref([])
const currentUser = ref(null)
const userId = ref(route.query.userId || route.query.id || null)
const followingCount = ref(0)
const followersCount = ref(0)

const isOwnProfile = computed(
  () => !userId.value || currentUser.value?.id === parseInt(userId.value, 10)
)

const pageTitle = computed(() => {
  if (isOwnProfile.value) {
    return tab.value === 'following' ? '我的关注' : '我的粉丝'
  }
  return tab.value === 'following' ? '他的关注' : '他的粉丝'
})

const currentList = computed(() => {
  return tab.value === 'following' ? following.value : followers.value
})

const followingIds = computed(() => {
  return new Set(following.value.map(u => u.id))
})

const isFollowing = (targetUserId) => {
  return followingIds.value.has(targetUserId)
}

const applyLists = (list) => {
  const items = list || []
  following.value = items
  followers.value = items
  followingCount.value = items.length
  followersCount.value = items.length
}

const loadOwnFollowData = async () => {
  if (!currentUser.value?.id) return
  loading.value = true
  try {
    const res = await api.get('/friends')
    applyLists(res.data.friends)
  } catch (error) {
    console.error('获取关注/粉丝列表失败', error)
  } finally {
    loading.value = false
  }
}

const loadOtherUserFollowData = async () => {
  if (!userId.value || !currentUser.value?.id) return
  loading.value = true
  try {
    const res = await api.get(`/friends/mutual/${userId.value}`)
    applyLists(res.data.mutual_friends || res.data.friends)
  } catch (error) {
    console.error('获取对方的关注/粉丝列表失败', error)
  } finally {
    loading.value = false
  }
}

const handleFollowToggle = async (user) => {
  try {
    if (isFollowing(user.id)) {
      await api.delete(`/friends/${user.id}`)
      following.value = following.value.filter(u => u.id !== user.id)
      followingCount.value = following.value.length
    } else {
      // 后端好友请求接口期望 friend_id 字段
      await api.post('/friends/request', { friend_id: user.id })
      await showSuccess('已发送好友请求')
    }
  } catch (error) {
    console.error('操作失败', error)
  }
}

const loadData = () => {
  if (userId.value) {
    loadOtherUserFollowData()
  } else {
    loadOwnFollowData()
  }
}

onMounted(async () => {
  try {
    const token = localStorage.getItem('token')
    if (token) {
      const profileRes = await api.get('/profile')
      currentUser.value = profileRes.data
    }
  } catch (error) {
    console.error('获取当前用户信息失败', error)
  }
  loadData()
})

watch(
  () => [route.query.userId, route.query.id, route.query.tab],
  ([newUserId, newId, newTab]) => {
    userId.value = newUserId || newId || null
    tab.value = newTab || 'following'
    loadData()
  }
)
</script>

<style scoped>
.follow-list-page {
  max-width: 800px;
  margin: 0 auto;
  padding: 16px;
}
</style>
