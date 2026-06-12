<template>
  <v-container class="max-w-2xl mx-auto px-4 py-8">
    <!-- 返回按钮 -->
    <v-btn 
      text 
      color="gray-600" 
      class="mb-6 hover:text-primary transition-colors"
      @click="router.back()"
    >
      <v-icon class="mr-2" size="20">mdi-arrow-left</v-icon>
      返回
    </v-btn>
    
    <v-card rounded="2xl" elevation="4" class="overflow-hidden">
      <!-- 标题 -->
      <v-card-title class="gradient-purple text-white py-6 px-8">
        <v-icon class="mr-3" size="24">{{ tab === 'following' ? 'mdi-user-follow' : 'mdi-users' }}</v-icon>
        <span class="font-bold text-xl">{{ pageTitle }}</span>
      </v-card-title>
      
      <!-- Tab 切换 -->
      <v-tabs 
        v-model="tab" 
        background-color="surface"
        class="border-b border-gray-100"
      >
        <v-tab 
          value="following" 
          class="text-gray-600 hover:text-primary transition-colors"
        >
          <v-icon class="mr-2" size="18">mdi-user-follow</v-icon>
          关注 ({{ followingCount }})
        </v-tab>
        <v-tab 
          value="followers" 
          class="text-gray-600 hover:text-primary transition-colors"
        >
          <v-icon class="mr-2" size="18">mdi-users</v-icon>
          粉丝 ({{ followersCount }})
        </v-tab>
      </v-tabs>
      
      <!-- 内容区域 -->
      <div class="p-6">
        <div v-if="loading" class="loading-center">
          <v-progress-circular indeterminate color="primary" :size="48" />
        </div>
        
        <div v-else-if="currentList.length === 0" class="empty-state">
          <v-icon size="96" color="gray-200" class="empty-state-icon">
            {{ tab === 'following' ? 'mdi-user-follow' : 'mdi-users' }}
          </v-icon>
          <p class="text-gray-400">暂无{{ tab === 'following' ? '关注' : '粉丝' }}</p>
        </div>
        
        <v-list v-else class="space-y-3">
          <v-card 
            v-for="userItem in currentList"
            :key="userItem.id"
            rounded="xl"
            class="card-hover"
          >
            <v-list-item class="px-4 py-3">
              <template v-slot:prepend>
                <v-avatar size="56" color="primary" class="avatar-hover">
                  <v-icon size="28" color="white">mdi-account</v-icon>
                </v-avatar>
              </template>
              
              <v-list-item-content>
                <div class="flex items-center">
                  <router-link 
                    :to="'/profile?id=' + userItem.id" 
                    class="text-gray-800 font-bold hover:text-primary transition-colors"
                  >
                    {{ userItem.display_name || userItem.username }}
                  </router-link>
                  <v-chip 
                    v-if="userItem.role === 'admin'" 
                    size="small" 
                    class="ml-2 tag-purple"
                  >
                    管理员
                  </v-chip>
                </div>
                <p class="text-gray-500 text-sm">@{{ userItem.username }}</p>
                <p v-if="userItem.signature" class="text-gray-400 text-sm mt-1">{{ userItem.signature }}</p>
              </v-list-item-content>
              
              <template v-slot:append>
                <v-btn
                  v-if="!isOwnProfile && currentUser && currentUser.id !== userItem.id"
                  :class="[
                    'btn-outline-purple',
                    isFollowing(userItem.id) ? 'bg-gray-100 text-gray-600' : ''
                  ]"
                  size="small"
                  @click="handleFollowToggle(userItem)"
                >
                  {{ isFollowing(userItem.id) ? '已关注' : '+ 关注' }}
                </v-btn>
              </template>
            </v-list-item>
          </v-card>
        </v-list>
      </div>
    </v-card>
  </v-container>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { authApi, followApi } from '../api'

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
    const [followingRes, followersRes] = await Promise.all([
      followApi.getFollowing(),
      followApi.getFollowers()
    ])

    following.value = followingRes.data.following || []
    followers.value = followersRes.data.followers || []
    followingCount.value = following.value.length
    followersCount.value = followers.value.length
  } catch (error) {
    console.error('获取关注列表失败', error)
  } finally {
    loading.value = false
  }
}

const loadOtherUserFollowData = async () => {
  if (!userId.value) return

  loading.value = true

  try {
    const [followingRes, followersRes] = await Promise.all([
      followApi.getUserFollowing(userId.value),
      followApi.getUserFollowers(userId.value)
    ])

    following.value = followingRes.data.following || []
    followers.value = followersRes.data.followers || []
    followingCount.value = following.value.length
    followersCount.value = followers.value.length
  } catch (error) {
    console.error('获取用户关注列表失败', error)
  } finally {
    loading.value = false
  }
}

const handleFollowToggle = async (user) => {
  try {
    if (isFollowing(user.id)) {
      await followApi.unfollow(user.id)
      following.value = following.value.filter(u => u.id !== user.id)
      followingCount.value--
    } else {
      await followApi.follow(user.id)
      following.value.push(user)
      followingCount.value++
    }
  } catch (error) {
    console.error('操作失败', error)
  }
}

onMounted(async () => {
  try {
    const profileRes = await authApi.getProfile()
    currentUser.value = profileRes.data
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