<template>
  <div class="follow-page">
    <div class="page-head">
      <button class="back-btn" @click="router.back()" aria-label="返回">
        <el-icon :size="20"><ArrowLeft /></el-icon>
      </button>
      <h2 class="page-title">{{ pageTitle }}</h2>
    </div>

    <!-- Tab 切换 -->
    <div class="seg-tabs">
      <button
        class="seg-item"
        :class="{ active: tab === 'following' }"
        @click="switchTab('following')"
      >
        关注 <span class="seg-count">{{ followingCount }}</span>
      </button>
      <button
        class="seg-item"
        :class="{ active: tab === 'followers' }"
        @click="switchTab('followers')"
      >
        粉丝 <span class="seg-count">{{ followersCount }}</span>
      </button>
    </div>

    <div v-if="loading" class="loading-wrap">
      <div class="loading"></div>
    </div>

    <div v-else-if="currentList.length === 0" class="empty-state">
      <el-icon :size="56" class="empty-icon"><UserFilled /></el-icon>
      <div class="empty-text">暂无{{ tab === 'following' ? '关注' : '粉丝' }}</div>
    </div>

    <div v-else class="user-grid">
      <div v-for="userItem in currentList" :key="userItem.id" class="card-surface user-card">
        <div class="card-top">
          <UserAvatar :user="userItem" :size="54" :show-username="false" />
          <div class="user-main">
            <router-link :to="`/profile?id=${userItem.id}`" class="user-name-link">
              {{ userItem.display_name || userItem.username }}
            </router-link>
            <div class="user-sub">
              <span class="text-secondary">@{{ userItem.username }}</span>
              <span v-if="userItem.role === 'admin'" class="admin-chip">管理员</span>
            </div>
          </div>
        </div>

        <p v-if="userItem.signature" class="user-signature">{{ userItem.signature }}</p>

        <button
          v-if="!isOwnProfile && currentUser && currentUser.id !== userItem.id"
          class="follow-btn"
          :class="isFollowing(userItem.id) ? 'followed' : ''"
          @click="handleFollowToggle(userItem)"
        >
          {{ isFollowing(userItem.id) ? '已关注' : '关注' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ArrowLeft, UserFilled } from '@element-plus/icons-vue'
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

const switchTab = (name) => {
  tab.value = name
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
.follow-page {
  max-width: 900px;
  margin: 0 auto;
  padding: 16px;
}

.page-head {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.back-btn {
  width: 38px;
  height: 38px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  border-radius: 12px;
  background: var(--campus-surface);
  color: var(--campus-text);
  cursor: pointer;
  box-shadow: var(--campus-shadow-sm);
  transition: var(--campus-transition);
}

.back-btn:hover {
  background: var(--campus-primary-soft);
  color: var(--campus-primary);
}

.page-title {
  margin: 0;
  font-size: 20px;
  font-weight: 800;
  color: var(--campus-text);
}

/* ---------------- 分段切换 ---------------- */
.seg-tabs {
  display: inline-flex;
  background: var(--campus-surface-2);
  border: 1px solid var(--campus-border);
  border-radius: 14px;
  padding: 4px;
  gap: 4px;
  margin-bottom: 20px;
}

.seg-item {
  border: none;
  background: transparent;
  padding: 8px 20px;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 600;
  color: var(--campus-text-secondary);
  cursor: pointer;
  transition: var(--campus-transition);
}

.seg-item.active {
  background: var(--campus-surface);
  color: var(--campus-primary);
  box-shadow: var(--campus-shadow-sm);
}

.seg-count {
  margin-left: 4px;
  font-size: 12px;
  opacity: 0.7;
}

/* ---------------- loading / empty ---------------- */
.loading-wrap {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 40px 0;
}

.loading {
  width: 32px;
  height: 32px;
  border: 3px solid var(--campus-primary-light);
  border-top-color: var(--campus-primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.empty-state {
  text-align: center;
  padding: 56px 0;
}

.empty-icon {
  color: var(--campus-border);
  margin-bottom: 12px;
}

.empty-text {
  font-size: 15px;
  color: var(--campus-text-secondary);
}

/* ---------------- 用户卡片网格 ---------------- */
.user-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}

.user-card {
  padding: 20px;
  border-radius: var(--campus-radius);
  box-shadow: var(--campus-shadow-sm);
  transition: var(--campus-transition);
}

.user-card:hover {
  box-shadow: var(--campus-shadow);
  transform: translateY(-2px);
}

.card-top {
  display: flex;
  align-items: center;
  gap: 14px;
}

.card-top :deep(.user-avatar .avatar) {
  background: var(--campus-primary);
}

.user-main {
  flex: 1;
  min-width: 0;
}

.user-name-link {
  font-weight: 700;
  font-size: 16px;
  color: var(--campus-text);
  text-decoration: none;
  word-break: break-word;
  transition: var(--campus-transition);
}

.user-name-link:hover {
  color: var(--campus-primary);
}

.user-sub {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 6px;
  font-size: 13px;
  margin-top: 2px;
}

.admin-chip {
  font-size: 11px;
  font-weight: 600;
  padding: 1px 8px;
  border-radius: 999px;
  color: #fff;
  background: linear-gradient(120deg, #ef4444, #f59e0b);
}

.user-signature {
  margin: 14px 0 0;
  font-size: 13px;
  color: var(--campus-text-secondary);
  line-height: 1.6;
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  word-break: break-word;
}

.follow-btn {
  width: 100%;
  margin-top: 16px;
  padding: 9px 0;
  border: none;
  border-radius: 12px;
  background: var(--campus-primary);
  color: #fff;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: var(--campus-transition);
}

.follow-btn:hover {
  background: var(--campus-primary-dark);
}

.follow-btn.followed {
  background: var(--campus-surface-2);
  color: var(--campus-text-secondary);
  border: 1px solid var(--campus-border);
}

.follow-btn.followed:hover {
  border-color: var(--campus-danger);
  color: var(--campus-danger);
}

/* ---------------- 响应式 ---------------- */
@media (max-width: 640px) {
  .user-grid {
    grid-template-columns: 1fr;
  }

  .seg-item {
    flex: 1;
    text-align: center;
  }

  .seg-tabs {
    display: flex;
    width: 100%;
  }
}
</style>
