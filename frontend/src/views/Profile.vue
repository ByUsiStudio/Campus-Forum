<template>
  <div class="page-container">
    <!-- 加载状态 -->
    <div v-if="loading" class="profile-loading">
      <div class="loading"></div>
    </div>

    <!-- 未登录提示 -->
    <div v-else-if="!isLoggedIn" class="auth-page-center">
      <el-card class="auth-card card-surface text-center">
        <el-icon class="not-login-icon"><Lock /></el-icon>
        <div class="not-login-title">登录后可查看个人中心</div>
        <div class="not-login-sub text-secondary mb-4">
          登录后您可以编辑个人资料、查看发布的文章、与其他人互动
        </div>
        <div class="d-flex gap-2 justify-center">
          <router-link to="/login">
            <el-button type="primary"><el-icon style="margin-right:4px"><SwitchButton /></el-icon>登录</el-button>
          </router-link>
          <router-link to="/register">
            <el-button plain><el-icon style="margin-right:4px"><User /></el-icon>注册</el-button>
          </router-link>
        </div>
      </el-card>
    </div>

    <!-- 用户信息 -->
    <template v-else-if="user">
      <el-row :gutter="16">
        <!-- 用户信息卡片 -->
        <el-col :xs="24" :md="8" class="mb-4">
          <el-card class="card-surface profile-card">
            <div class="text-center">
              <div class="avatar-wrap mb-4">
                <UserAvatar :user="user" :size="150" />
              </div>

              <div class="mb-4">
                <div class="profile-name">{{ user.display_name }}</div>
                <div class="text-secondary profile-username">@{{ user.username }}</div>
                <el-tag v-if="user.role === 'admin'" type="danger" size="small" class="mt-2">
                  <el-icon style="margin-right:4px"><Lock /></el-icon>
                  管理员
                </el-tag>
              </div>

              <el-button v-if="isOwnProfile" plain type="primary" class="w-full" @click="changeAvatar">
                <el-icon style="margin-right:4px"><Camera /></el-icon>
                更换头像
              </el-button>
              <el-button
                v-else
                :plain="!followStatus.is_following"
                :type="followStatus.is_following ? 'default' : 'primary'"
                class="w-full"
                @click="handleFollow"
              >
                <el-icon style="margin-right:4px">
                  <Check v-if="followStatus.is_following" /><Plus v-else />
                </el-icon>
                {{ followStatus.is_following ? '已关注' : followStatus.is_followed ? '回关' : '关注' }}
              </el-button>
            </div>

            <el-divider class="my-4" />

            <!-- 统计数据 -->
            <el-row>
              <el-col :span="8" class="stat-col cursor-pointer" @click="goToFollowing">
                <div class="stat-num">{{ followingCount }}</div>
                <div class="text-secondary">关注</div>
              </el-col>
              <el-col :span="8" class="stat-col cursor-pointer" @click="goToFollowers">
                <div class="stat-num">{{ followersCount }}</div>
                <div class="text-secondary">粉丝</div>
              </el-col>
              <el-col :span="8" class="stat-col">
                <div class="stat-num">{{ articleCount }}</div>
                <div class="text-secondary">文章</div>
              </el-col>
            </el-row>

            <el-divider class="my-4" />

            <!-- 用户信息列表 -->
            <div class="user-info-list">
              <div class="info-item">
                <el-icon class="info-icon"><Message /></el-icon>
                <span>QQ号：{{ user.qq_number || '未设置' }}</span>
              </div>
              <div class="info-item">
                <el-icon class="info-icon"><EditPen /></el-icon>
                <span class="whitespace-pre-line">签名：{{ user.signature || '暂无签名' }}</span>
              </div>
              <div class="info-item">
                <el-icon class="info-icon"><Calendar /></el-icon>
                <span>注册时间：{{ formatDate(user.created_at) }}</span>
              </div>
            </div>
          </el-card>
        </el-col>

        <!-- 文章和编辑资料 -->
        <el-col :xs="24" :md="16">
          <!-- 编辑资料卡片 -->
          <el-card v-if="isOwnProfile" class="card-surface mb-4">
            <template #header>
              <div class="card-header-title">
                <el-icon style="margin-right:8px"><EditPen /></el-icon>
                编辑资料
              </div>
            </template>

            <el-form :model="editForm" label-position="top" @submit.prevent="updateProfile">
              <el-form-item label="显示名称">
                <el-input
                  v-model="editForm.display_name"
                  placeholder="请输入显示名称"
                  :prefix-icon="Postcard"
                />
              </el-form-item>

              <el-form-item label="个性化签名">
                <el-input
                  v-model="editForm.signature"
                  type="textarea"
                  :rows="3"
                  maxlength="200"
                  show-word-limit
                  placeholder="最多200个字符"
                />
              </el-form-item>

              <el-button type="primary" native-type="submit" :loading="saving">
                <el-icon style="margin-right:4px"><Files /></el-icon>
                保存修改
              </el-button>
            </el-form>
          </el-card>

          <!-- 文章列表 -->
          <el-card class="card-surface">
            <template #header>
              <div class="card-header-title">
                <el-icon style="margin-right:8px"><Document /></el-icon>
                {{ isOwnProfile ? '我的文章' : '他的文章' }}
              </div>
            </template>

            <el-tabs v-if="isOwnProfile" v-model="activeTab" class="mb-4">
              <el-tab-pane :label="`已发布 (${myArticles.length})`" name="published">
                <div v-if="myArticles.length === 0" class="empty-state">
                  <el-icon class="empty-icon"><DocumentRemove /></el-icon>
                  <div class="text-secondary">暂无文章</div>
                  <router-link to="/create">
                    <el-button type="primary" class="mt-4">
                      <el-icon style="margin-right:4px"><EditPen /></el-icon>
                      写文章
                    </el-button>
                  </router-link>
                </div>

                <div v-else>
                  <div v-for="article in myArticles" :key="article.id" class="article-row">
                    <el-avatar :size="50" :src="article.user?.avatar || ''" class="article-avatar">
                      {{ (article.user?.display_name || article.user?.username || 'U')[0]?.toUpperCase() }}
                    </el-avatar>
                    <div class="article-body">
                      <div class="article-title">
                        <router-link :to="'/article/' + article.id" class="text-primary">
                          {{ article.title }}
                        </router-link>
                      </div>
                      <div class="article-meta">
                        <span><el-icon><Clock /></el-icon> {{ formatDate(article.created_at) }}</span>
                        <span class="text-error"><el-icon><Star /></el-icon> {{ article.like_count }}</span>
                        <span><el-icon><View /></el-icon> {{ article.view_count }}</span>
                        <el-tag v-if="article.category" size="small" type="primary" effect="plain">
                          {{ article.category.name }}
                        </el-tag>
                      </div>
                    </div>
                    <div class="article-actions">
                      <router-link :to="'/create?id=' + article.id">
                        <el-button text size="small" type="primary"><el-icon><EditPen /></el-icon></el-button>
                      </router-link>
                      <el-button text size="small" type="danger" @click="deleteArticle(article.id)">
                        <el-icon><Delete /></el-icon>
                      </el-button>
                    </div>
                  </div>
                </div>
              </el-tab-pane>

              <el-tab-pane :label="`草稿 (${drafts.length})`" name="drafts">
                <div v-if="drafts.length === 0" class="empty-state">
                  <el-icon class="empty-icon"><DocumentRemove /></el-icon>
                  <div class="text-secondary">暂无草稿</div>
                  <router-link to="/create">
                    <el-button type="primary" class="mt-4">
                      <el-icon style="margin-right:4px"><EditPen /></el-icon>
                      写文章
                    </el-button>
                  </router-link>
                </div>

                <div v-else>
                  <div v-for="article in drafts" :key="article.id" class="article-row">
                    <el-icon class="article-avatar draft-icon"><DocumentRemove /></el-icon>
                    <div class="article-body">
                      <div class="article-title">{{ article.title }}</div>
                      <div class="article-meta">
                        <span><el-icon><Clock /></el-icon> 最后修改：{{ formatDate(article.updated_at) }}</span>
                        <el-tag v-if="article.category" size="small" type="info" effect="plain">
                          {{ article.category.name }}
                        </el-tag>
                      </div>
                    </div>
                    <div class="article-actions">
                      <router-link :to="'/create?id=' + article.id">
                        <el-button text size="small" type="primary"><el-icon><EditPen /></el-icon></el-button>
                      </router-link>
                      <el-button text size="small" type="success" @click="publishDraft(article.id)">
                        <el-icon><Promotion /></el-icon>
                      </el-button>
                      <el-button text size="small" type="danger" @click="deleteDraft(article.id)">
                        <el-icon><Delete /></el-icon>
                      </el-button>
                    </div>
                  </div>
                </div>
              </el-tab-pane>
            </el-tabs>

            <!-- 非本人：文章列表 -->
            <template v-if="!isOwnProfile">
              <div v-if="myArticles.length === 0" class="empty-state">
                <el-icon class="empty-icon"><DocumentRemove /></el-icon>
                <div class="text-secondary">暂无文章</div>
              </div>
              <div v-else>
                <div v-for="article in myArticles" :key="article.id" class="article-row">
                  <el-avatar :size="50" :src="article.user?.avatar || ''" class="article-avatar">
                    {{ (article.user?.display_name || article.user?.username || 'U')[0]?.toUpperCase() }}
                  </el-avatar>
                  <div class="article-body">
                    <div class="article-title">
                      <router-link :to="'/article/' + article.id" class="text-primary">
                        {{ article.title }}
                      </router-link>
                    </div>
                    <div class="article-meta">
                      <span><el-icon><Clock /></el-icon> {{ formatDate(article.created_at) }}</span>
                      <span class="text-error"><el-icon><Star /></el-icon> {{ article.like_count }}</span>
                      <span><el-icon><View /></el-icon> {{ article.view_count }}</span>
                      <el-tag v-if="article.category" size="small" type="primary" effect="plain">
                        {{ article.category.name }}
                      </el-tag>
                    </div>
                  </div>
                </div>
              </div>
            </template>
          </el-card>
        </el-col>
      </el-row>
    </template>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  Lock, SwitchButton, User, Camera, Check, Plus, Message, EditPen,
  Calendar, Postcard, Files, Document, DocumentRemove, Clock, Star,
  View, Delete, Promotion
} from '@element-plus/icons-vue'
import api, { articleApi, userApi, friendApi, uploadApi } from '../api'
import UserAvatar from '../components/UserAvatar.vue'
import { useUserStore } from '../stores/user'
import { success, error, confirm, prompt } from '../utils/message'

export default {
  name: 'Profile',
  components: {
    UserAvatar
  },
  setup() {
    const route = useRoute()
    const router = useRouter()
    const userStore = useUserStore()

    const loading = ref(true)
    const user = ref(null)
    const myArticles = ref([])
    const drafts = ref([])
    const activeTab = ref('published')
    const saving = ref(false)
    const editForm = ref({
      display_name: '',
      signature: ''
    })
    const followStatus = ref({
      is_following: false,
      is_followed: false,
      mutual: false
    })
    const followingCount = ref(0)
    const followersCount = ref(0)
    const articleCount = ref(0)

    const targetUserId = computed(() => route.query.id || route.query.userId)
    const isLoggedIn = computed(() => userStore.isLoggedIn)
    const isOwnProfile = computed(() => {
      return !targetUserId.value || (userStore.user && userStore.user.id === user.value?.id)
    })

    const loadProfile = async () => {
      if (!userStore.token && !targetUserId.value) {
        return false
      }

      try {
        if (targetUserId.value) {
          const response = await userApi.getUserByID(targetUserId.value)
          user.value = response.data
        } else {
          const response = await api.get('/profile')
          user.value = response.data
          editForm.value.display_name = user.value.display_name || ''
          editForm.value.signature = user.value.signature || ''
        }
        return true
      } catch (err) {
        console.error('加载用户信息失败', err)
        if (!targetUserId.value) {
          router.push('/')
        }
        return false
      }
    }

    const loadFollowStatus = async () => {
      if (!user.value || isOwnProfile.value || !userStore.user) return

      try {
        const response = await friendApi.checkFriendStatus(user.value.id)
        const isFriend = !!response.data.is_friend
        followStatus.value = {
          is_following: isFriend,
          is_followed: isFriend,
          mutual: isFriend
        }
      } catch (err) {
        console.error('加载好友状态失败', err)
      }
    }

    const handleFollow = async () => {
      if (!userStore.user) {
        router.push('/login')
        return
      }

      try {
        if (followStatus.value.is_following) {
          await friendApi.deleteFriend(user.value.id)
          followStatus.value.is_following = false
          followStatus.value.is_followed = false
          followStatus.value.mutual = false
        } else {
          await friendApi.sendFriendRequest({ user_id: user.value.id })
          success('已发送好友请求')
        }
      } catch (err) {
        console.error('好友操作失败', err)
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
        const friendsRes = await friendApi.getMutualFriends(targetId)
        followingCount.value = friendsRes.data.friends?.length || 0
        followersCount.value = friendsRes.data.friends?.length || 0
      } catch (err) {
        console.error('加载好友数据失败', err)
      }
    }

    const loadMyArticles = async () => {
      try {
        if (targetUserId.value) {
          const response = await userApi.getUserArticles(targetUserId.value)
          myArticles.value = response.data.articles || []
        } else {
          const response = await articleApi.getMyArticles()
          myArticles.value = response.data.articles || response.data || []
        }
        articleCount.value = myArticles.value.length
      } catch (err) {
        console.error('加载文章失败', err)
      }
    }

    const loadDrafts = async () => {
      if (!isOwnProfile.value) return
      try {
        const response = await articleApi.getDrafts()
        drafts.value = response.data.articles || []
      } catch (err) {
        console.error('加载草稿失败', err)
      }
    }

    const publishDraft = async (id) => {
      try {
        await confirm('确定要发布这篇草稿吗？')
      } catch (e) {
        return
      }
      try {
        await articleApi.publishDraft(id)
        await loadDrafts()
        await loadMyArticles()
        success('发布成功')
      } catch (err) {
        console.error('发布失败', err)
        error('发布失败')
      }
    }

    const deleteDraft = async (id) => {
      try {
        await confirm('确定要删除这篇草稿吗？')
      } catch (e) {
        return
      }
      try {
        await articleApi.deleteArticle(id)
        await loadDrafts()
        success('删除成功')
      } catch (err) {
        console.error('删除失败', err)
        error('删除失败')
      }
    }

    const updateProfile = async () => {
      saving.value = true
      try {
        await api.put('/profile', {
          display_name: editForm.value.display_name,
          signature: editForm.value.signature
        })
        user.value.display_name = editForm.value.display_name
        user.value.signature = editForm.value.signature
        if (userStore.user) {
          const next = { ...userStore.user, display_name: editForm.value.display_name, signature: editForm.value.signature }
          userStore.setUser(next)
        }
        success('更新成功')
      } catch (err) {
        console.error('更新失败', err)
        error('更新失败')
      } finally {
        saving.value = false
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
          const response = await uploadApi.uploadAvatar(formData)
          user.value.avatar = response.data.url
          if (userStore.user) {
            userStore.setUser({ ...userStore.user, avatar: response.data.url })
          }
          success('头像更新成功')
        } catch (err) {
          console.error('上传头像失败', err)
          error('上传失败')
        }
      }
      input.click()
    }

    const deleteArticle = async (id) => {
      try {
        await confirm('确定要删除这篇文章吗？')
      } catch (e) {
        return
      }

      const userRole = user.value?.role
      if (userRole !== 'admin') {
        let reason
        try {
          reason = await prompt('请输入删除原因（管理员将审核）：', { placeholder: '删除原因' })
        } catch (e) {
          return
        }
        if (!reason) return

        try {
          await articleApi.deleteArticle(id, { data: { reason } })
          success('删除申请已提交，等待管理员审核')
          myArticles.value = myArticles.value.filter(a => a.id !== id)
        } catch (err) {
          console.error('提交删除申请失败', err)
          error('提交失败')
        }
      } else {
        try {
          await articleApi.deleteArticle(id)
          myArticles.value = myArticles.value.filter(a => a.id !== id)
          success('删除成功')
        } catch (err) {
          console.error('删除失败', err)
          error('删除失败')
        }
      }
    }

    const formatDate = (date) => {
      if (!date) return ''
      return new Date(date).toLocaleString('zh-CN')
    }

    onMounted(async () => {
      loading.value = true
      const ok = await loadProfile()
      if (ok) {
        try {
          await Promise.all([
            loadMyArticles(),
            loadDrafts(),
            loadFollowStatus(),
            loadFollowCounts()
          ])
        } finally {
          loading.value = false
        }
      } else {
        loading.value = false
      }
    })

    return {
      user,
      myArticles,
      drafts,
      activeTab,
      editForm,
      saving,
      updateProfile,
      changeAvatar,
      deleteArticle,
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
      loading,
      Lock, SwitchButton, User, Camera, Check, Plus, Message, EditPen,
      Calendar, Postcard, Files, Document, DocumentRemove, Clock, Star,
      View, Delete, Promotion
    }
  }
}
</script>

<style scoped>
.loading {
  display: inline-block;
  width: 48px;
  height: 48px;
  box-sizing: border-box;
  border: 4px solid var(--campus-border);
  border-top-color: var(--campus-primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.profile-loading {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 60vh;
}

.auth-page-center {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 60vh;
}

.auth-card {
  width: 100%;
  max-width: 460px;
  padding: 16px;
}

.not-login-icon {
  font-size: 64px;
  color: var(--campus-text-secondary);
  margin-bottom: 16px;
}

.not-login-title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 8px;
}

.not-login-sub {
  font-size: 13px;
}

.profile-card {
  height: 100%;
}

.avatar-wrap {
  display: flex;
  justify-content: center;
}

.profile-name {
  font-size: 22px;
  font-weight: 700;
  margin-bottom: 4px;
}

.profile-username {
  font-size: 14px;
}

.stat-col {
  text-align: center;
}

.stat-num {
  font-size: 18px;
  font-weight: 700;
}

.user-info-list {
  text-align: left;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: var(--campus-text);
}

.info-icon {
  color: var(--campus-text-secondary);
  flex-shrink: 0;
}

.card-header-title {
  display: flex;
  align-items: center;
  font-weight: 600;
}

.empty-state {
  text-align: center;
  padding: 32px;
}

.empty-icon {
  font-size: 64px;
  color: var(--campus-border);
  margin-bottom: 16px;
}

.article-row {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px 0;
  border-bottom: 1px solid var(--campus-border);
}

.article-row:last-child {
  border-bottom: none;
}

.article-avatar {
  flex-shrink: 0;
  background-color: var(--campus-primary);
  color: #fff;
  font-weight: 600;
}

.article-body {
  flex: 1;
  min-width: 0;
}

.article-title {
  font-weight: 700;
  margin-bottom: 6px;
  word-break: break-word;
}

.article-meta {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
  font-size: 13px;
  color: var(--campus-text-secondary);
}

.article-meta .el-icon {
  vertical-align: -2px;
  margin-right: 3px;
}

.text-error {
  color: var(--el-color-danger);
}

.article-actions {
  flex-shrink: 0;
  display: flex;
  align-items: center;
}

.whitespace-pre-line {
  white-space: pre-line;
}

.draft-icon {
  width: 50px;
  height: 50px;
  flex-shrink: 0;
  font-size: 30px;
  color: var(--campus-text-secondary);
  margin-top: 10px;
}

:deep(.w-full) {
  width: 100%;
}

:deep(.mt-2) { margin-top: 8px; }
:deep(.mt-4) { margin-top: 16px; }
:deep(.mb-4) { margin-bottom: 16px; }
:deep(.my-4) { margin-top: 16px; margin-bottom: 16px; }

.d-flex { display: flex; }
.justify-center { justify-content: center; }
.gap-2 { gap: 8px; }
.text-center { text-align: center; }
.mb-4 { margin-bottom: 16px; }
.mt-4 { margin-top: 16px; }
.mt-2 { margin-top: 8px; }
.my-4 { margin-top: 16px; margin-bottom: 16px; }
.text-secondary { color: var(--campus-text-secondary); }
.text-primary { color: var(--campus-primary); }
.cursor-pointer { cursor: pointer; }
</style>
