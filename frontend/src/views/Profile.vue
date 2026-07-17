<template>
  <div v-if="user" class="profile-page">
    <div class="profile-container">
      <div class="profile-sidebar">
        <div class="user-card">
          <div class="user-avatar">
            <UserAvatar :user="user" :size="150" />
          </div>
          
          <div class="user-info">
            <div class="user-name">{{ user.display_name }}</div>
            <div class="user-username">@{{ user.username }}</div>
            <span v-if="user.role === 'admin'" class="admin-tag">
              <i class="fa-solid fa-shield-halved"></i>
              管理员
            </span>
          </div>

          <div class="user-actions">
            <button v-if="isOwnProfile" class="layui-btn layui-btn-normal layui-btn-sm w-full" @click="changeAvatar">
              <i class="fa-solid fa-camera mr-2"></i>
              更换头像
            </button>
            <button v-else class="layui-btn layui-btn-sm w-full" :class="followStatus.is_following ? 'layui-btn-primary' : 'layui-btn-normal'" @click="handleFollow">
              <i :class="followStatus.is_following ? 'fa-solid fa-check' : 'fa-solid fa-plus'" class="mr-2"></i>
              {{ followStatus.is_following ? '已关注' : followStatus.is_followed ? '回关' : '关注' }}
            </button>
          </div>

          <div class="stats-row">
            <div class="stat-item" @click="goToFollowing">
              <div class="stat-value">{{ followingCount }}</div>
              <div class="stat-label">关注</div>
            </div>
            <div class="stat-item" @click="goToFollowers">
              <div class="stat-value">{{ followersCount }}</div>
              <div class="stat-label">粉丝</div>
            </div>
            <div class="stat-item">
              <div class="stat-value">{{ articleCount }}</div>
              <div class="stat-label">文章</div>
            </div>
          </div>

          <div class="info-list">
            <div class="info-item">
              <i class="fa-solid fa-comment"></i>
              <span>QQ号：{{ user.qq_number || '未设置' }}</span>
            </div>
            <div class="info-item">
              <i class="fa-solid fa-pencil"></i>
              <span>签名：{{ user.signature || '暂无签名' }}</span>
            </div>
            <div class="info-item">
              <i class="fa-solid fa-calendar"></i>
              <span>注册时间：{{ formatDate(user.created_at) }}</span>
            </div>
          </div>
        </div>
      </div>

      <div class="profile-content">
        <div v-if="isOwnProfile" class="edit-card">
          <div class="card-header">
            <i class="fa-solid fa-user-pen mr-2"></i>
            编辑资料
          </div>
          
          <form @submit.prevent="updateProfile" class="edit-form">
            <div class="form-item">
              <label class="form-label">显示名称</label>
              <input 
                type="text"
                v-model="editForm.display_name"
                placeholder="请输入显示名称"
                class="layui-input"
              />
            </div>
            
            <div class="form-item">
              <label class="form-label">个性化签名</label>
              <textarea 
                v-model="editForm.signature"
                placeholder="请输入个性化签名"
                rows="3"
                class="layui-textarea"
              ></textarea>
            </div>
            
            <button type="submit" class="layui-btn layui-btn-normal">
              <i class="fa-solid fa-save mr-2"></i>
              保存修改
            </button>
          </form>
        </div>

        <div class="articles-card">
          <div class="card-header">
            <i class="fa-solid fa-file-lines mr-2"></i>
            {{ isOwnProfile ? '我的文章' : '他的文章' }}
          </div>

          <div v-if="isOwnProfile" class="tabs">
            <button 
              class="tab-btn" 
              :class="{ active: activeTab === 'published' }"
              @click="activeTab = 'published'"
            >
              已发布 ({{ myArticles.length }})
            </button>
            <button 
              class="tab-btn" 
              :class="{ active: activeTab === 'drafts' }"
              @click="activeTab = 'drafts'"
            >
              草稿 ({{ drafts.length }})
            </button>
          </div>

          <div v-if="activeTab === 'published'" class="articles-list">
            <div v-if="myArticles.length === 0" class="empty-state">
              <i class="fa-regular fa-file-lines"></i>
              <div>暂无文章</div>
              <button class="layui-btn layui-btn-normal" @click="goToCreate">
                <i class="fa-solid fa-pencil mr-2"></i>
                写文章
              </button>
            </div>

            <div v-else class="article-items">
              <div v-for="article in myArticles" :key="article.id" class="article-item">
                <div class="article-avatar">
                  <UserAvatar :user="article.user" :size="50" />
                </div>
                <div class="article-info">
                  <div class="article-title">
                    <router-link :to="'/article/' + article.id">{{ article.title }}</router-link>
                  </div>
                  <div class="article-meta">
                    <span class="meta-item">
                      <i class="fa-regular fa-clock"></i>
                      {{ formatDate(article.created_at) }}
                    </span>
                    <span class="meta-item">
                      <i class="fa-solid fa-heart text-error"></i>
                      {{ article.like_count }}
                    </span>
                    <span class="meta-item">
                      <i class="fa-regular fa-eye"></i>
                      {{ article.view_count }}
                    </span>
                    <span v-if="article.category" class="category-tag">{{ article.category.name }}</span>
                  </div>
                </div>
                <div class="article-actions" v-if="isOwnProfile">
                  <router-link :to="'/create?id=' + article.id" class="action-btn edit">
                    <i class="fa-solid fa-pencil"></i>
                  </router-link>
                  <button class="action-btn delete" @click="deleteArticle(article.id)">
                    <i class="fa-solid fa-trash"></i>
                  </button>
                </div>
              </div>
            </div>
          </div>

          <div v-if="activeTab === 'drafts'" class="articles-list">
            <div v-if="drafts.length === 0" class="empty-state">
              <i class="fa-regular fa-file-pen"></i>
              <div>暂无草稿</div>
              <button class="layui-btn layui-btn-normal" @click="goToCreate">
                <i class="fa-solid fa-pencil mr-2"></i>
                写文章
              </button>
            </div>

            <div v-else class="article-items">
              <div v-for="article in drafts" :key="article.id" class="article-item">
                <div class="article-avatar">
                  <i class="fa-regular fa-file-lines text-muted"></i>
                </div>
                <div class="article-info">
                  <div class="article-title">{{ article.title }}</div>
                  <div class="article-meta">
                    <span class="meta-item">
                      <i class="fa-regular fa-clock"></i>
                      最后修改：{{ formatDate(article.updated_at) }}
                    </span>
                    <span v-if="article.category" class="category-tag">{{ article.category.name }}</span>
                  </div>
                </div>
                <div class="article-actions">
                  <router-link :to="'/create?id=' + article.id" class="action-btn edit">
                    <i class="fa-solid fa-pencil"></i>
                  </router-link>
                  <button class="action-btn publish" @click="publishDraft(article.id)">
                    <i class="fa-solid fa-paper-plane"></i>
                  </button>
                  <button class="action-btn delete" @click="deleteDraft(article.id)">
                    <i class="fa-solid fa-trash"></i>
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <div v-else-if="!isLoggedIn" class="login-prompt">
    <div class="prompt-card">
      <i class="fa-solid fa-user-lock"></i>
      <div class="prompt-title">登录后可查看个人中心</div>
      <div class="prompt-desc">登录后您可以编辑个人资料、查看发布的文章、与其他人互动</div>
      <div class="prompt-actions">
        <router-link to="/login" class="layui-btn layui-btn-normal">
          <i class="fa-solid fa-right-to-bracket mr-2"></i>
          登录
        </router-link>
        <router-link to="/register" class="layui-btn layui-btn-primary">
          <i class="fa-solid fa-user-plus mr-2"></i>
          注册
        </router-link>
      </div>
    </div>
  </div>

  <div v-else class="loading-state">
    <div class="loading-spinner"></div>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api, { articleApi } from '../api'
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
    const drafts = ref([])
    const activeTab = ref('published')
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
        if (!targetUserId.value) {
          router.push('/')
        }
      }
    }

    const loadFollowStatus = async () => {
      if (!user.value || isOwnProfile.value || !currentUser.value) return

      try {
        const response = await api.get(`/friends/status/${user.value.id}`)
        followStatus.value = {
          is_following: response.data.is_friend,
          is_followed: response.data.is_friend,
          mutual: response.data.is_friend
        }
      } catch (error) {
        console.error('加载好友状态失败', error)
      }
    }

    const handleFollow = async () => {
      if (!currentUser.value) {
        router.push('/login')
        return
      }

      try {
        if (followStatus.value.is_following) {
          await api.delete(`/friends/${user.value.id}`)
          followStatus.value.is_following = false
          followStatus.value.mutual = false
        } else {
          await api.post('/friends/request', { user_id: user.value.id })
          alert('已发送好友请求')
        }
      } catch (error) {
        console.error('好友操作失败', error)
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
        const friendsRes = await api.get(`/friends/mutual/${targetId}`)
        followingCount.value = friendsRes.data.friends?.length || 0
        followersCount.value = friendsRes.data.friends?.length || 0
      } catch (error) {
        console.error('加载好友数据失败', error)
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

    const loadDrafts = async () => {
      if (!isOwnProfile.value) return
      try {
        const response = await articleApi.getDrafts()
        drafts.value = response.data.articles || []
      } catch (error) {
        console.error('加载草稿失败', error)
      }
    }

    const publishDraft = async (id) => {
      if (!confirm('确定要发布这篇草稿吗？')) return
      try {
        await articleApi.publishDraft(id)
        await loadDrafts()
        await loadMyArticles()
        alert('发布成功')
      } catch (error) {
        console.error('发布失败', error)
        alert('发布失败')
      }
    }

    const deleteDraft = async (id) => {
      if (!confirm('确定要删除这篇草稿吗？')) return
      try {
        await api.delete(`/articles/${id}`)
        await loadDrafts()
        alert('删除成功')
      } catch (error) {
        console.error('删除失败', error)
        alert('删除失败')
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

    const goToCreate = () => {
      router.push('/create')
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
        loadDrafts()
        loadFollowStatus()
        loadFollowCounts()
      })
    })

    return {
      user,
      myArticles,
      drafts,
      activeTab,
      editForm,
      updateProfile,
      changeAvatar,
      deleteArticle,
      loadDrafts,
      publishDraft,
      deleteDraft,
      formatDate,
      followStatus,
      handleFollow,
      goToFollowing,
      goToFollowers,
      followingCount,
      followersCount,
      articleCount,
      isOwnProfile,
      isLoggedIn,
      goToCreate
    }
  }
}
</script>

<style scoped>
.profile-page {
  padding: 24px 0;
}

.profile-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 15px;
  display: flex;
  gap: 24px;
}

.profile-sidebar {
  width: 360px;
  flex-shrink: 0;
}

.profile-content {
  flex: 1;
}

.user-card {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
}

.user-avatar {
  text-align: center;
  margin-bottom: 16px;
}

.user-info {
  text-align: center;
  margin-bottom: 16px;
}

.user-name {
  font-size: 20px;
  font-weight: 600;
  color: #333;
}

.user-username {
  font-size: 14px;
  color: #999;
  margin-top: 4px;
}

.admin-tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #FF5722;
  background: rgba(255, 87, 34, 0.1);
  padding: 2px 8px;
  border-radius: 4px;
  margin-top: 8px;
}

.user-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 20px;
}

.stats-row {
  display: flex;
  justify-content: space-around;
  padding: 16px 0;
  border-top: 1px solid #f0f0f0;
  border-bottom: 1px solid #f0f0f0;
  margin-bottom: 16px;
}

.stat-item {
  text-align: center;
  cursor: pointer;
  transition: color 0.2s;

  &:hover {
    color: #1E9FFF;
  }
}

.stat-value {
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.stat-label {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.info-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #666;

  i {
    font-size: 14px;
    color: #999;
  }
}

.edit-card, .articles-card {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
  margin-bottom: 20px;
}

.card-header {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-bottom: 20px;
}

.edit-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-label {
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

.tabs {
  display: flex;
  gap: 0;
  margin-bottom: 20px;
  border-bottom: 1px solid #f0f0f0;
}

.tab-btn {
  background: none;
  border: none;
  padding: 12px 20px;
  font-size: 14px;
  color: #666;
  cursor: pointer;
  border-bottom: 2px solid transparent;
  transition: all 0.2s;

  &:hover {
    color: #1E9FFF;
  }

  &.active {
    color: #1E9FFF;
    border-bottom-color: #1E9FFF;
  }
}

.articles-list {
  min-height: 200px;
}

.empty-state {
  text-align: center;
  padding: 48px 24px;
  color: #999;

  i {
    font-size: 64px;
    margin-bottom: 16px;
    color: #e8e8e8;
  }

  div {
    font-size: 16px;
    margin-bottom: 16px;
    color: #666;
  }
}

.article-items {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.article-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  background: #f8f9fa;
  border-radius: 8px;
  transition: background 0.2s;

  &:hover {
    background: #f0f0f0;
  }
}

.article-avatar {
  flex-shrink: 0;

  i {
    font-size: 32px;
    color: #999;
  }
}

.article-info {
  flex: 1;
  min-width: 0;
}

.article-title {
  font-size: 15px;
  font-weight: 500;
  color: #333;
  margin-bottom: 8px;

  a {
    color: inherit;
    text-decoration: none;

    &:hover {
      color: #1E9FFF;
    }
  }
}

.article-meta {
  display: flex;
  align-items: center;
  gap: 16px;
  font-size: 13px;
  color: #999;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.category-tag {
  font-size: 12px;
  color: #1E9FFF;
  background: rgba(30, 159, 255, 0.1);
  padding: 2px 8px;
  border-radius: 4px;
}

.article-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  width: 36px;
  height: 36px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  transition: all 0.2s;

  &.edit {
    background: rgba(30, 159, 255, 0.1);
    color: #1E9FFF;

    &:hover {
      background: rgba(30, 159, 255, 0.2);
    }
  }

  &.publish {
    background: rgba(82, 196, 26, 0.1);
    color: #52C41A;

    &:hover {
      background: rgba(82, 196, 26, 0.2);
    }
  }

  &.delete {
    background: rgba(255, 87, 34, 0.1);
    color: #FF5722;

    &:hover {
      background: rgba(255, 87, 34, 0.2);
    }
  }
}

.text-error {
  color: #FF5722;
}

.text-muted {
  color: #999;
}

.login-prompt {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 60vh;
  padding: 24px;
}

.prompt-card {
  background: #fff;
  border-radius: 12px;
  padding: 48px;
  text-align: center;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  max-width: 400px;
}

.prompt-card i {
  font-size: 64px;
  color: #e8e8e8;
  margin-bottom: 16px;
}

.prompt-title {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin-bottom: 8px;
}

.prompt-desc {
  font-size: 14px;
  color: #999;
  margin-bottom: 24px;
}

.prompt-actions {
  display: flex;
  gap: 12px;
  justify-content: center;
}

.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 60vh;
}

.loading-spinner {
  width: 48px;
  height: 48px;
  border: 4px solid #f0f0f0;
  border-top-color: #1E9FFF;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.w-full {
  width: 100%;
}

.mr-2 {
  margin-right: 8px;
}

@media (max-width: 768px) {
  .profile-container {
    flex-direction: column;
  }

  .profile-sidebar {
    width: 100%;
  }
}
</style>