<template>
  <v-row v-if="user">
    <!-- 用户信息卡片 -->
    <v-col cols="12" md="4">
      <v-card class="pa-6 text-center" elevation="2">
        <v-avatar size="150" class="mb-4">
          <UserAvatar :user="user" :size="150" />
        </v-avatar>
        
        <div class="mb-4">
          <div class="text-h5 font-weight-bold mb-1">{{ user.display_name }}</div>
          <div class="text-body-2 text-medium-emphasis">@{{ user.username }}</div>
          <v-chip v-if="user.role === 'admin'" color="error" size="small" class="mt-2">
            <v-icon start size="small">mdi-shield-crown</v-icon>
            管理员
          </v-chip>
        </div>

        <v-btn v-if="isOwnProfile" variant="outlined" color="primary" @click="changeAvatar" block>
          <v-icon start>mdi-camera</v-icon>
          更换头像
        </v-btn>
        <v-btn v-else variant="tonal" :color="followStatus.is_following ? 'default' : 'primary'" @click="handleFollow" block>
          <v-icon start>{{ followStatus.is_following ? 'mdi-check' : 'mdi-plus' }}</v-icon>
          {{ followStatus.is_following ? '已关注' : followStatus.is_followed ? '回关' : '关注' }}
        </v-btn>

        <v-divider class="my-4"></v-divider>

        <!-- 统计数据 -->
        <v-row dense>
          <v-col cols="4" class="text-center cursor-pointer" @click="goToFollowing">
            <div class="text-h6 font-weight-bold">{{ followingCount }}</div>
            <div class="text-caption text-medium-emphasis">关注</div>
          </v-col>
          <v-col cols="4" class="text-center cursor-pointer" @click="goToFollowers">
            <div class="text-h6 font-weight-bold">{{ followersCount }}</div>
            <div class="text-caption text-medium-emphasis">粉丝</div>
          </v-col>
          <v-col cols="4" class="text-center">
            <div class="text-h6 font-weight-bold">{{ articleCount }}</div>
            <div class="text-caption text-medium-emphasis">文章</div>
          </v-col>
        </v-row>

        <v-divider class="my-4"></v-divider>

        <!-- 用户信息列表 -->
        <v-list density="compact" class="text-left">
          <v-list-item prepend-icon="mdi-qqchat">
            <v-list-item-title>QQ号：{{ user.qq_number || '未设置' }}</v-list-item-title>
          </v-list-item>
          <v-list-item prepend-icon="mdi-pencil">
            <v-list-item-title class="whitespace-pre-line">签名：{{ user.signature || '暂无签名' }}</v-list-item-title>
          </v-list-item>
          <v-list-item prepend-icon="mdi-calendar">
            <v-list-item-title>注册时间：{{ formatDate(user.created_at) }}</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-card>
    </v-col>

    <!-- 文章和编辑资料 -->
    <v-col cols="12" md="8">
      <!-- 编辑资料卡片 -->
      <v-card v-if="isOwnProfile" class="pa-6 mb-4" elevation="2">
        <v-card-title class="d-flex align-center mb-4">
          <v-icon class="mr-2">mdi-account-edit</v-icon>
          编辑资料
        </v-card-title>
        
        <v-form @submit.prevent="updateProfile">
          <v-text-field
            v-model="editForm.display_name"
            label="显示名称"
            variant="outlined"
            prepend-inner-icon="mdi-card-account-details"
            class="mb-4"
          ></v-text-field>
          
          <v-textarea
            v-model="editForm.signature"
            label="个性化签名"
            variant="outlined"
            rows="3"
            class="mb-4"
            hint="最多200个字符"
            persistent-hint
            counter="200"
          ></v-textarea>
          
          <v-btn type="submit" color="primary" prepend-icon="mdi-content-save">
            保存修改
          </v-btn>
        </v-form>
      </v-card>

      <!-- 文章列表 -->
      <v-card class="pa-6" elevation="2">
        <v-card-title class="d-flex align-center mb-4">
          <v-icon class="mr-2">mdi-file-document</v-icon>
          {{ isOwnProfile ? '我的文章' : '他的文章' }}
          <v-spacer></v-spacer>
          <span class="text-body-2 text-medium-emphasis">{{ myArticles.length }} 篇</span>
        </v-card-title>

        <div v-if="myArticles.length === 0" class="text-center pa-8">
          <v-icon size="64" color="grey-lighten-1" class="mb-4">mdi-file-document-outline</v-icon>
          <div class="text-body-1 text-medium-emphasis">暂无文章</div>
          <v-btn v-if="isOwnProfile" color="primary" class="mt-4" to="/create">
            <v-icon start>mdi-pencil</v-icon>
            写文章
          </v-btn>
        </div>

        <v-list v-else lines="two" class="pa-0">
          <v-list-item
            v-for="article in myArticles"
            :key="article.id"
            class="px-0"
          >
            <template v-slot:prepend>
              <UserAvatar :user="article.user" :size="50" />
            </template>

            <v-list-item-title class="font-weight-bold mb-1">
              <router-link :to="'/article/' + article.id" class="text-decoration-none text-primary">
                {{ article.title }}
              </router-link>
            </v-list-item-title>
            
            <v-list-item-subtitle class="d-flex align-center flex-wrap gap-2">
              <span>
                <v-icon size="x-small">mdi-clock</v-icon>
                {{ formatDate(article.created_at) }}
              </span>
              <span>
                <v-icon size="x-small" class="text-error">mdi-heart</v-icon>
                {{ article.like_count }}
              </span>
              <span>
                <v-icon size="x-small">mdi-eye</v-icon>
                {{ article.view_count }}
              </span>
              <v-chip v-if="article.category" size="x-small" color="primary" variant="flat">
                {{ article.category.name }}
              </v-chip>
            </v-list-item-subtitle>

            <template v-slot:append v-if="isOwnProfile">
              <v-btn
                variant="text"
                size="small"
                :to="'/create?id=' + article.id"
                color="primary"
                icon="mdi-pencil"
              ></v-btn>
              <v-btn
                variant="text"
                size="small"
                color="error"
                @click="deleteArticle(article.id)"
                icon="mdi-delete"
              ></v-btn>
            </template>
          </v-list-item>
        </v-list>
      </v-card>
    </v-col>
  </v-row>

  <!-- 未登录提示 -->
  <v-container v-else-if="!isLoggedIn" fluid class="pa-4 fill-height">
    <v-row justify="center" align="center" class="fill-height">
      <v-col cols="12" sm="8" md="6" lg="4">
        <v-card class="text-center pa-8" elevation="2">
          <v-icon size="80" color="grey-lighten-1" class="mb-4">mdi-account-lock</v-icon>
          <div class="text-h6 mb-4">登录后可查看个人中心</div>
          <div class="text-body-2 text-medium-emphasis mb-6">
            登录后您可以编辑个人资料、查看发布的文章、与其他人互动
          </div>
          <div class="d-flex gap-2 justify-center">
            <v-btn color="primary" to="/login" prepend-icon="mdi-login">
              登录
            </v-btn>
            <v-btn variant="outlined" to="/register" prepend-icon="mdi-account-plus">
              注册
            </v-btn>
          </div>
        </v-card>
      </v-col>
    </v-row>
  </v-container>

  <!-- 加载状态 -->
  <div v-else class="d-flex justify-center align-center" style="min-height: 60vh;">
    <v-progress-circular indeterminate color="primary" size="64"></v-progress-circular>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '../api'
import UserAvatar from '../components/UserAvatar.vue'

export default {
  name: 'Profile',
  components: {
    UserAvatar
  },
  setup() {
    const route = useRoute()
    const router = useRouter()
    const token = ref(localStorage.getItem('token'))
    const user = ref(null)
    const myArticles = ref([])
    const editForm = ref({
      display_name: '',
      signature: ''
    })
    const followStatus = ref({
      is_following: false,
      is_followed: false,
      mutual: false
    })
    const targetUserId = computed(() => route.query.id || route.query.userId)
    const currentUser = ref(null)
    const followingCount = ref(0)
    const followersCount = ref(0)
    const articleCount = ref(0)
    const isOwnProfile = computed(() => {
      return !targetUserId.value || (currentUser.value && currentUser.value.id === user.value?.id)
    })
    const isLoggedIn = computed(() => !!token.value)

    const loadProfile = async () => {
      if (!token.value && !targetUserId.value) return
      
      try {
        if (targetUserId.value) {
          const response = await api.get(`/users/${targetUserId.value}`)
          user.value = response.data
        } else {
          const response = await api.get('/profile')
          user.value = response.data
          editForm.value.display_name = user.value.display_name
          editForm.value.signature = user.value.signature || ''
        }
      } catch (error) {
        console.error('加载用户信息失败', error)
        if (targetUserId.value) {
          // 查看他人资料失败，保持加载状态
        } else {
          router.push('/')
        }
      }
    }

    const loadFollowStatus = async () => {
      if (!user.value || isOwnProfile.value || !currentUser.value) return

      try {
        const response = await api.get(`/follow/status/${user.value.id}`)
        followStatus.value = response.data
      } catch (error) {
        console.error('加载关注状态失败', error)
      }
    }

    const handleFollow = async () => {
      if (!currentUser.value) {
        router.push('/login')
        return
      }

      try {
        if (followStatus.value.is_following) {
          await api.delete(`/follow/${user.value.id}`)
          followStatus.value.is_following = false
          followStatus.value.mutual = false
        } else {
          await api.post(`/follow/${user.value.id}`)
          followStatus.value.is_following = true
          followStatus.value.mutual = followStatus.value.is_followed
        }
      } catch (error) {
        console.error('关注失败', error)
      }
    }

    const goToFollowing = () => {
      if (user.value) {
        router.push({ path: '/follow-list', query: { id: user.value.id } })
      }
    }

    const goToFollowers = () => {
      if (user.value) {
        router.push({ path: '/follow-list', query: { id: user.value.id, tab: 'followers' } })
      }
    }

    const loadFollowCounts = async () => {
      if (!user.value) return
      try {
        const targetId = targetUserId.value || user.value.id
        const [followingRes, followersRes] = await Promise.all([
          api.get(`/users/${targetId}/following`),
          api.get(`/users/${targetId}/followers`)
        ])
        followingCount.value = followingRes.data.following?.length || 0
        followersCount.value = followersRes.data.followers?.length || 0
      } catch (error) {
        console.error('加载关注数据失败', error)
      }
    }

    const loadMyArticles = async () => {
      try {
        if (targetUserId.value) {
          const response = await api.get(`/users/${targetUserId.value}/articles`)
          myArticles.value = response.data.articles || []
          articleCount.value = myArticles.value.length
        } else {
          const response = await api.get('/articles', {
            params: { page: 1, page_size: 50 }
          })
          if (user.value) {
            myArticles.value = response.data.articles.filter(
              a => a.user_id === user.value.id
            )
            articleCount.value = myArticles.value.length
          }
        }
      } catch (error) {
        console.error('加载文章失败', error)
      }
    }

    const updateProfile = async () => {
      try {
        await api.put('/profile', {
          display_name: editForm.value.display_name,
          signature: editForm.value.signature
        })
        user.value.display_name = editForm.value.display_name
        user.value.signature = editForm.value.signature
        const storedUser = JSON.parse(localStorage.getItem('user') || '{}')
        storedUser.display_name = editForm.value.display_name
        storedUser.signature = editForm.value.signature
        localStorage.setItem('user', JSON.stringify(storedUser))
        alert('更新成功')
      } catch (error) {
        console.error('更新失败', error)
        alert('更新失败')
      }
    }

    const changeAvatar = () => {
      const input = document.createElement('input')
      input.type = 'file'
      input.accept = 'image/*'
      input.onchange = async (e) => {
        const file = e.target.files[0]
        if (!file) return

        const formData = new FormData()
        formData.append('avatar', file)

        try {
          const response = await api.post('/upload/avatar', formData)
          user.value.avatar = response.data.url
          const storedUser = JSON.parse(localStorage.getItem('user') || '{}')
          storedUser.avatar = response.data.url
          localStorage.setItem('user', JSON.stringify(storedUser))
          alert('头像更新成功')
        } catch (error) {
          console.error('上传头像失败', error)
          alert('上传失败')
        }
      }
      input.click()
    }

    const deleteArticle = async (id) => {
      if (!confirm('确定要删除这篇文章吗？')) return

      const userRole = user.value?.role
      if (userRole !== 'admin') {
        const reason = prompt('请输入删除原因（管理员将审核）：')
        if (!reason) return

        try {
          await api.delete(`/articles/${id}`, { data: { reason } })
          alert('删除申请已提交，等待管理员审核')
          myArticles.value = myArticles.value.filter(a => a.id !== id)
        } catch (error) {
          console.error('提交删除申请失败', error)
          alert('提交失败')
        }
      } else {
        try {
          await api.delete(`/articles/${id}`)
          myArticles.value = myArticles.value.filter(a => a.id !== id)
          alert('删除成功')
        } catch (error) {
          console.error('删除失败', error)
          alert('删除失败')
        }
      }
    }

    const formatDate = (date) => {
      return new Date(date).toLocaleString('zh-CN')
    }

    onMounted(() => {
      const storedUser = localStorage.getItem('user')
      if (storedUser) {
        currentUser.value = JSON.parse(storedUser)
      }
      loadProfile().then(() => {
        loadMyArticles()
        loadFollowStatus()
        loadFollowCounts()
      })
    })

    return {
      user,
      myArticles,
      editForm,
      updateProfile,
      changeAvatar,
      deleteArticle,
      formatDate,
      followStatus,
      handleFollow,
      goToFollowing,
      goToFollowers,
      followingCount,
      followersCount,
      articleCount,
      isOwnProfile,
      isLoggedIn
    }
  }
}
</script>