<template>
  <div class="profile-container" v-if="user">
    <div class="profile-header">
      <div class="avatar-section">
        <img :src="user.avatar" :alt="user.display_name" class="avatar-large">
        <button @click="changeAvatar" class="btn btn-secondary">更换头像</button>
      </div>
      <div class="info-section">
        <h2>{{ user.display_name }}</h2>
        <p>用户名：{{ user.username }}</p>
        <p>QQ号：{{ user.qq_number }}</p>
        <p>角色：{{ user.role === 'admin' ? '管理员' : '普通用户' }}</p>
        <p>注册时间：{{ formatDate(user.created_at) }}</p>
      </div>
    </div>
    
    <div class="profile-edit">
      <h3>编辑资料</h3>
      <form @submit.prevent="updateProfile">
        <div class="form-group">
          <label>显示名称</label>
          <input type="text" v-model="editForm.display_name">
        </div>
        <button type="submit" class="btn btn-primary">保存修改</button>
      </form>
    </div>
    
    <div class="my-articles">
      <h3>我的文章</h3>
      <div v-if="myArticles.length === 0" class="empty">暂无文章</div>
      <div v-for="article in myArticles" :key="article.id" class="article-item">
        <div class="article-info">
          <router-link :to="'/article/' + article.id" class="article-title">{{ article.title }}</router-link>
          <div class="article-meta">
            <span>{{ formatDate(article.created_at) }}</span>
            <span>❤️ {{ article.like_count }}</span>
            <span>👁️ {{ article.view_count }}</span>
          </div>
        </div>
        <div class="article-actions">
          <router-link :to="'/create?id=' + article.id" class="btn btn-secondary">编辑</router-link>
          <button @click="deleteArticle(article.id)" class="btn btn-danger">删除</button>
        </div>
      </div>
    </div>
  </div>
  
  <div v-else class="loading">加载中...</div>
</template>

<script>
import { ref, onMounted } from 'vue'
import api from '../api'

export default {
  name: 'Profile',
  setup() {
    const user = ref(null)
    const myArticles = ref([])
    const editForm = ref({
      display_name: ''
    })
    
    const loadProfile = async () => {
      try {
        const response = await api.get('/profile')
        user.value = response.data
        editForm.value.display_name = user.value.display_name
      } catch (error) {
        console.error('加载用户信息失败', error)
      }
    }
    
    const loadMyArticles = async () => {
      try {
        const response = await api.get('/articles', {
          params: { page: 1, page_size: 50 }
        })
        // 过滤出当前用户的文章
        if (user.value) {
          myArticles.value = response.data.articles.filter(
            a => a.user_id === user.value.id
          )
        }
      } catch (error) {
        console.error('加载文章失败', error)
      }
    }
    
    const updateProfile = async () => {
      try {
        await api.put('/profile', {
          display_name: editForm.value.display_name
        })
        user.value.display_name = editForm.value.display_name
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
          // 更新本地存储
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
      loadProfile().then(() => {
        loadMyArticles()
      })
    })
    
    return {
      user,
      myArticles,
      editForm,
      updateProfile,
      changeAvatar,
      deleteArticle,
      formatDate
    }
  }
}
</script>

<style scoped>
.profile-container {
  background: white;
  border-radius: 12px;
  padding: 30px;
}

.profile-header {
  display: grid;
  grid-template-columns: 200px 1fr;
  gap: 30px;
  padding-bottom: 30px;
  border-bottom: 1px solid #e5e7eb;
  margin-bottom: 30px;
}

.avatar-section {
  text-align: center;
}

.avatar-large {
  width: 150px;
  height: 150px;
  border-radius: 50%;
  object-fit: cover;
  margin-bottom: 15px;
  border: 3px solid #10b981;
}

.info-section h2 {
  margin-bottom: 15px;
  color: #1e293b;
}

.info-section p {
  margin-bottom: 10px;
  color: #4b5563;
}

.profile-edit {
  padding-bottom: 30px;
  border-bottom: 1px solid #e5e7eb;
  margin-bottom: 30px;
}

.profile-edit h3 {
  margin-bottom: 20px;
}

.my-articles h3 {
  margin-bottom: 20px;
}

.article-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
  border-bottom: 1px solid #e5e7eb;
}

.article-info {
  flex: 1;
}

.article-title {
  font-size: 16px;
  font-weight: 500;
  text-decoration: none;
  color: #1e293b;
}

.article-meta {
  margin-top: 8px;
  font-size: 12px;
  color: #6b7280;
  display: flex;
  gap: 15px;
}

.article-actions {
  display: flex;
  gap: 10px;
}

.empty {
  text-align: center;
  padding: 40px;
  color: #6b7280;
}
</style>