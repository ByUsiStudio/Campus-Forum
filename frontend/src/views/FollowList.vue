<template>
  <div class="follow-list-page">
    <div class="follow-card">
      <div class="card-header">
        <button class="back-btn" @click="router.back()">
          <i class="fa-solid fa-arrow-left"></i>
        </button>
        <h2 class="page-title">{{ pageTitle }}</h2>
      </div>

      <div class="tabs-container">
        <button 
          class="tab-btn" 
          :class="{ active: tab === 'following' }"
          @click="tab = 'following'"
        >
          关注 ({{ followingCount }})
        </button>
        <button 
          class="tab-btn" 
          :class="{ active: tab === 'followers' }"
          @click="tab = 'followers'"
        >
          粉丝 ({{ followersCount }})
        </button>
      </div>

      <div v-if="loading" class="loading-state">
        <div class="loading-spinner"></div>
      </div>

      <div v-else-if="currentList.length === 0" class="empty-state">
        暂无{{ tab === 'following' ? '关注' : '粉丝' }}
      </div>

      <div v-else class="follow-list">
        <div 
          v-for="userItem in currentList" 
          :key="userItem.id"
          class="follow-item"
        >
          <router-link :to="'/profile?id=' + userItem.id" class="avatar-link">
            <div class="user-avatar">
              <i class="fa-solid fa-user"></i>
            </div>
          </router-link>

          <div class="user-info">
            <div class="user-name">
              <router-link :to="'/profile?id=' + userItem.id">
                {{ userItem.display_name || userItem.username }}
              </router-link>
              <span v-if="userItem.role === 'admin'" class="admin-tag">管理员</span>
            </div>
            <div class="user-username">@{{ userItem.username }}</div>
            <div v-if="userItem.signature" class="user-signature">{{ userItem.signature }}</div>
          </div>

          <button
            v-if="!isOwnProfile && currentUser && currentUser.id !== userItem.id"
            class="follow-btn"
            :class="{ followed: isFollowing(userItem.id) }"
            @click="handleFollowToggle(userItem)"
          >
            {{ isFollowing(userItem.id) ? '已关注' : '关注' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import api from '../api'

const router = useRouter()
const route = useRoute()

const tab = ref(route.query.tab || 'following')
const loading = ref(false)
const following = ref([])
const followers = ref([])
const currentUser = ref(null)
const userId = ref(route.query.userId || route.query.id || null)
const isOwnProfile = computed(() => !userId.value || currentUser.value?.id === parseInt(userId.value))
const followingCount = ref(0)
const followersCount = ref(0)

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

const loadFollowData = async () => {
  if (!userId.value && !currentUser.value?.id) return

  const targetId = userId.value || currentUser.value.id
  loading.value = true

  try {
    const friendsRes = await api.get('/friends')

    following.value = friendsRes.data.friends || []
    followers.value = friendsRes.data.friends || []
    followingCount.value = following.value.length
    followersCount.value = followers.value.length
  } catch (error) {
    console.error('获取好友列表失败', error)
  } finally {
    loading.value = false
  }
}

const loadOtherUserFollowData = async () => {
  if (!userId.value) return

  loading.value = true

  try {
    const friendsRes = await api.get(`/friends/mutual/${userId.value}`)

    following.value = friendsRes.data.friends || []
    followers.value = friendsRes.data.friends || []
    followingCount.value = following.value.length
    followersCount.value = followers.value.length
  } catch (error) {
    console.error('获取好友列表失败', error)
  } finally {
    loading.value = false
  }
}

const handleFollowToggle = async (user) => {
  try {
    if (isFollowing(user.id)) {
      await api.delete(`/friends/${user.id}`)
      following.value = following.value.filter(u => u.id !== user.id)
      followingCount.value--
    } else {
      await api.post('/friends/request', { user_id: user.id })
    }
  } catch (error) {
    console.error('操作失败', error)
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
    console.error('获取当前用户失败', error)
  }

  if (userId.value) {
    loadOtherUserFollowData()
  } else {
    loadFollowData()
  }
})

watch(() => [route.query.userId, route.query.id, route.query.tab], ([newUserId, newId, newTab]) => {
  userId.value = newUserId || newId || null
  tab.value = newTab || 'following'
  if (userId.value) {
    loadOtherUserFollowData()
  } else {
    loadFollowData()
  }
})
</script>

<style scoped>
.follow-list-page {
  max-width: 800px;
  margin: 0 auto;
  padding: 16px;
}

.follow-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
  overflow: hidden;
}

.card-header {
  display: flex;
  align-items: center;
  padding: 16px 24px;
  border-bottom: 1px solid #f0f0f0;
}

.back-btn {
  background: none;
  border: none;
  font-size: 20px;
  color: #333;
  cursor: pointer;
  margin-right: 16px;
  padding: 8px;
  border-radius: 6px;
  
  &:hover {
    background: #f5f5f5;
  }
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin: 0;
}

.tabs-container {
  display: flex;
  border-bottom: 1px solid #f0f0f0;
}

.tab-btn {
  flex: 1;
  padding: 16px;
  background: none;
  border: none;
  font-size: 14px;
  font-weight: 500;
  color: #666;
  cursor: pointer;
  position: relative;
  transition: color 0.3s ease;
  
  &.active {
    color: var(--primary);
    
    &::after {
      content: '';
      position: absolute;
      bottom: 0;
      left: 50%;
      transform: translateX(-50%);
      width: 40px;
      height: 2px;
      background: var(--primary);
      border-radius: 2px;
    }
  }
  
  &:hover {
    color: var(--primary);
  }
}

.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48px;
}

.loading-spinner {
  width: 48px;
  height: 48px;
  border: 4px solid #f0f0f0;
  border-top-color: var(--primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.empty-state {
  text-align: center;
  padding: 48px;
  color: #999;
}

.follow-list {
  display: flex;
  flex-direction: column;
}

.follow-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 24px;
  border-bottom: 1px solid #f0f0f0;
  
  &:last-child {
    border-bottom: none;
  }
}

.avatar-link {
  flex-shrink: 0;
}

.user-avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  background: var(--primary);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  
  i {
    font-size: 24px;
  }
}

.user-info {
  flex: 1;
  min-width: 0;
}

.user-name {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 15px;
  font-weight: 600;
  color: #333;
  margin-bottom: 4px;
  
  a {
    color: #333;
    text-decoration: none;
    
    &:hover {
      color: var(--primary);
    }
  }
}

.admin-tag {
  font-size: 11px;
  color: white;
  background: var(--primary);
  padding: 2px 8px;
  border-radius: 4px;
}

.user-username {
  font-size: 13px;
  color: #999;
  margin-bottom: 4px;
}

.user-signature {
  font-size: 13px;
  color: #666;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.follow-btn {
  padding: 6px 16px;
  border: 1px solid var(--primary);
  border-radius: 20px;
  font-size: 13px;
  font-weight: 500;
  color: var(--primary);
  background: white;
  cursor: pointer;
  transition: all 0.3s ease;
  
  &:hover {
    background: rgba(30, 159, 255, 0.1);
  }
  
  &.followed {
    background: #f5f5f5;
    border-color: #e8e8e8;
    color: #666;
  }
}
</style>