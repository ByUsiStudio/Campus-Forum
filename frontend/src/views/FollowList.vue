<template>
  <div class="follow-list-page">
    <el-card shadow="never" class="card-surface">
      <div class="page-header">
        <el-button circle text @click="router.back()">
          <el-icon><ArrowLeft /></el-icon>
        </el-button>
        <span class="page-title">{{ pageTitle }}</span>
      </div>

      <el-tabs v-model="tab" class="mb-4">
        <el-tab-pane :label="`关注 (${followingCount})`" name="following" />
        <el-tab-pane :label="`粉丝 (${followersCount})`" name="followers" />
      </el-tabs>

      <div v-if="loading" class="loading-wrap">
        <div class="loading"></div>
      </div>

      <el-empty
        v-else-if="currentList.length === 0"
        :description="`暂无${tab === 'following' ? '关注' : '粉丝'}`"
      />

      <el-list v-else>
        <el-list-item v-for="userItem in currentList" :key="userItem.id">
          <div class="user-row">
            <UserAvatar :user="userItem" :size="50" :show-username="false" />

            <div class="user-main">
              <router-link :to="`/profile?id=${userItem.id}`" class="user-name-link">
                {{ userItem.display_name || userItem.username }}
              </router-link>
              <div class="user-sub">
                <span class="text-secondary">@{{ userItem.username }}</span>
                <el-tag v-if="userItem.role === 'admin'" type="danger" size="small" class="ml-1">
                  管理员
                </el-tag>
              </div>
              <div v-if="userItem.signature" class="user-signature text-truncate">
                {{ userItem.signature }}
              </div>
            </div>

            <el-button
              v-if="!isOwnProfile && currentUser && currentUser.id !== userItem.id"
              :type="isFollowing(userItem.id) ? 'default' : 'primary'"
              plain
              size="small"
              @click="handleFollowToggle(userItem)"
            >
              {{ isFollowing(userItem.id) ? '已关注' : '关注' }}
            </el-button>
          </div>
        </el-list-item>
      </el-list>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ArrowLeft } from '@element-plus/icons-vue'
import http from '@/api/http'
import { useUserStore } from '@/stores/user'
import UserAvatar from '@/components/UserAvatar.vue'
import { success as showSuccess } from '@/utils/message'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

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
    const res = await http.get('/friends')
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
    const res = await http.get(`/friends/mutual/${userId.value}`)
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
      await http.delete(`/friends/${user.id}`)
      following.value = following.value.filter(u => u.id !== user.id)
      followingCount.value = following.value.length
    } else {
      // 后端好友请求接口期望 friend_id 字段
      await http.post('/friends/request', { friend_id: user.id })
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
  if (userStore.isLoggedIn) {
    if (!userStore.user) {
      await userStore.initUser()
    }
    currentUser.value = userStore.user
  } else {
    currentUser.value = null
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

.page-header {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
}

.page-title {
  font-size: 18px;
  font-weight: 700;
}

.mb-4 {
  margin-bottom: 16px;
}

.loading-wrap {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 32px;
}

.loading {
  width: 32px;
  height: 32px;
  border: 3px solid var(--campus-primary-light, rgba(30, 136, 229, 0.3));
  border-top-color: var(--campus-primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.user-row {
  display: flex;
  align-items: center;
  width: 100%;
  gap: 12px;
}

.user-main {
  flex: 1;
  min-width: 0;
}

.user-name-link {
  font-weight: 700;
  text-decoration: none;
  color: var(--campus-text, inherit);
}

.user-name-link:hover {
  color: var(--campus-primary);
}

.user-sub {
  display: flex;
  align-items: center;
  font-size: 13px;
}

.ml-1 {
  margin-left: 4px;
}

.user-signature {
  font-size: 13px;
  color: var(--campus-text-secondary, #909399);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
