<template>
  <div v-if="user">
    <v-row>
      <v-col cols="12" md="4">
        <v-card class="pa-6 text-center">
          <v-avatar size="150" class="mb-4">
            <v-img :src="user.avatar" :alt="user.display_name"></v-img>
          </v-avatar>
          <v-btn v-if="isOwnProfile" variant="outlined" color="primary" @click="changeAvatar" block>
            <v-icon start>mdi-camera</v-icon>
            更换头像
          </v-btn>
          <div v-else class="d-flex flex-column gap-2">
            <v-btn variant="outlined" :color="followStatus.is_following ? 'default' : 'primary'" @click="handleFollow" block>
              {{ followStatus.is_following ? '已关注' : followStatus.is_followed ? '回关' : '关注' }}
            </v-btn>
            <!-- 临时去除聊天界面的入口 -->
            <!-- <v-btn v-if="followStatus.is_following || followStatus.mutual" variant="outlined" color="primary" :to="'/chat/' + user.id" block>
              <v-icon start>mdi-message</v-icon>
              发消息
            </v-btn> -->
          </div>
        </v-card>
      </v-col>
      
      <v-col cols="12" md="8">
        <v-card class="pa-6 mb-4">
          <v-card-title class="text-h5 mb-2">{{ user.display_name }}</v-card-title>
          <div v-if="user.titles && user.titles.length > 0" class="mb-4">
            <v-chip
              v-for="title in user.titles"
              :key="title.id"
              size="small"
              :color="title.color"
              class="mr-2"
            >
              <v-icon v-if="title.icon" start size="x-small">{{ title.icon }}</v-icon>
              {{ title.name }}
            </v-chip>
          </div>
          <v-list density="compact">
            <v-list-item>
              <template v-slot:prepend>
                <v-icon>mdi-account</v-icon>
              </template>
              <v-list-item-title>用户名：{{ user.username }}</v-list-item-title>
            </v-list-item>
            <v-list-item>
              <template v-slot:prepend>
                <v-icon>mdi-qqchat</v-icon>
              </template>
              <v-list-item-title>QQ号：{{ user.qq_number }}</v-list-item-title>
            </v-list-item>
            <v-list-item>
              <template v-slot:prepend>
                <v-icon>mdi-pencil</v-icon>
              </template>
              <v-list-item-title class="whitespace-pre-line">签名：{{ user.signature || '暂无签名' }}</v-list-item-title>
            </v-list-item>
            <v-list-item>
              <template v-slot:prepend>
                <v-icon>mdi-shield-account</v-icon>
              </template>
              <v-list-item-title>
                角色：{{ user.role === 'admin' ? '管理员' : '普通用户' }}
              </v-list-item-title>
            </v-list-item>
            <v-list-item>
              <template v-slot:prepend>
                <v-icon>mdi-calendar</v-icon>
              </template>
              <v-list-item-title>注册时间：{{ formatDate(user.created_at) }}</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-card>
        
        <v-card v-if="isOwnProfile" class="pa-6 mb-4">
          <v-card-subtitle class="text-h6 pa-0 mb-4">编辑资料</v-card-subtitle>
          <v-form @submit.prevent="updateProfile">
            <v-text-field
              v-model="editForm.display_name"
              label="显示名称"
              variant="outlined"
              class="mb-4"
            ></v-text-field>
            <v-textarea
              v-model="editForm.signature"
              label="个性化签名"
              variant="outlined"
              rows="3"
              class="mb-4"
              hint="最多200个字符"
              :maxlength="200"
            ></v-textarea>
            <v-btn type="submit" color="primary">保存修改</v-btn>
          </v-form>
        </v-card>
      </v-col>
    </v-row>
    
    <v-card class="pa-6 mt-4">
      <v-card-title class="text-h6 mb-4">{{ isOwnProfile ? '我的文章' : '他的文章' }}</v-card-title>
      
      <div v-if="myArticles.length === 0" class="text-center pa-8 text-medium-emphasis">
        暂无文章
      </div>
      
      <v-list v-else lines="two">
        <v-list-item
          v-for="article in myArticles"
          :key="article.id"
          class="px-0"
        >
          <template v-slot:prepend>
            <v-avatar color="primary" size="50" class="mr-4">
              <v-img v-if="article.user.avatar" :src="article.user.avatar"></v-img>
              <span v-else class="text-h6">{{ article.user.display_name?.[0] || 'U' }}</span>
            </v-avatar>
          </template>
          
          <v-list-item-title class="font-weight-bold">
            <router-link :to="'/article/' + article.id" class="text-decoration-none">
              {{ article.title }}
            </router-link>
          </v-list-item-title>
          <v-list-item-subtitle>
            <v-icon size="small">mdi-clock</v-icon>
            {{ formatDate(article.created_at) }}
            <v-icon size="small" class="ml-2">mdi-heart</v-icon>
            {{ article.like_count }}
            <v-icon size="small" class="ml-2">mdi-eye</v-icon>
            {{ article.view_count }}
          </v-list-item-subtitle>
          
          <template v-slot:append>
            <v-btn
              variant="text"
              size="small"
              :to="'/create?id=' + article.id"
              color="primary"
            >
              编辑
            </v-btn>
            <v-btn
              variant="text"
              size="small"
              color="error"
              @click="deleteArticle(article.id)"
            >
              删除
            </v-btn>
          </template>
        </v-list-item>
      </v-list>
    </v-card>
  </div>
  
  <div v-else class="d-flex justify-center align-center" style="min-height: 50vh;">
    <v-progress-circular indeterminate color="primary" size="64"></v-progress-circular>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '../api'

export default {
  name: 'Profile',
  setup() {
    const route = useRoute()
    const router = useRouter()
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
    const currentUser = ref(null)
    const isOwnProfile = computed(() => {
      return !route.query.id || (currentUser.value && currentUser.value.id === user.value?.id)
    })
    
    const loadProfile = async () => {
      try {
        if (route.query.id) {
          const response = await api.get(`/users/${route.query.id}`)
          user.value = response.data
        } else {
          const response = await api.get('/profile')
          user.value = response.data
          editForm.value.display_name = user.value.display_name
          editForm.value.signature = user.value.signature || ''
        }
      } catch (error) {
        console.error('加载用户信息失败', error)
        router.push('/')
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
    
    const loadMyArticles = async () => {
      try {
        if (route.query.id) {
          const response = await api.get(`/users/${route.query.id}/articles`)
          myArticles.value = response.data.articles || []
        } else {
          const response = await api.get('/articles', {
            params: { page: 1, page_size: 50 }
          })
          if (user.value) {
            myArticles.value = response.data.articles.filter(
              a => a.user_id === user.value.id
            )
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
      isOwnProfile
    }
  }
}
</script>
