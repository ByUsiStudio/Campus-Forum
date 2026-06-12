<script setup>
import { ref, inject, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { userApi } from '../api'

const router = useRouter()
const route = useRoute()
const user = inject('user')

const users = ref([])
const isFollowing = ref(false)

const loadUsers = async () => {
  try {
    const type = route.query.type || 'following'
    let response
    if (type === 'following') {
      response = await userApi.getFollowing()
    } else {
      response = await userApi.getFollowers()
    }
    users.value = response.data.users || []
  } catch (error) {
    console.error('加载失败:', error)
  }
}

const handleFollow = async (userId) => {
  try {
    await userApi.follow(userId)
    loadUsers()
  } catch (error) {
    console.error('关注失败:', error)
  }
}

const handleUnfollow = async (userId) => {
  try {
    await userApi.unfollow(userId)
    loadUsers()
  } catch (error) {
    console.error('取消关注失败:', error)
  }
}

onMounted(() => {
  if (!user.value) {
    router.push('/login')
    return
  }
  loadUsers()
})
</script>

<template>
  <v-app>
    <v-app-bar app>
      <v-btn icon @click="router.push('/profile')">
        <v-icon>mdi-arrow-left</v-icon>
      </v-btn>
      <v-toolbar-title>{{ route.query.type === 'followers' ? '粉丝列表' : '关注列表' }}</v-toolbar-title>
    </v-app-bar>
    
    <v-container class="py-6">
      <v-card v-if="users.length > 0">
        <v-list>
          <v-list-item v-for="u in users" :key="u.id">
            <v-list-item-avatar>
              <v-icon color="primary">mdi-account</v-icon>
            </v-list-item-avatar>
            <v-list-item-content>
              <v-list-item-title>{{ u.username }}</v-list-item-title>
            </v-list-item-content>
            <v-list-item-actions>
              <v-btn text color="primary" @click="handleFollow(u.id)">关注</v-btn>
            </v-list-item-actions>
          </v-list-item>
        </v-list>
      </v-card>
      
      <v-card v-else class="text-center py-12">
        <v-icon size="64" color="grey">mdi-users</v-icon>
        <p class="mt-4 text-grey">暂无{{ route.query.type === 'followers' ? '粉丝' : '关注' }}</p>
      </v-card>
    </v-container>
  </v-app>
</template>
