<template>
  <div class="page-container profile-page">
    <!-- 加载状态 -->
    <div v-if="loading" class="profile-loading">
      <div class="loading"></div>
    </div>

    <!-- 未登录提示 -->
    <div v-else-if="!isLoggedIn" class="auth-page-center">
      <div class="card-surface auth-card">
        <div class="not-login-icon">
          <el-icon :size="56"><Lock /></el-icon>
        </div>
        <div class="not-login-title">登录后可查看个人中心</div>
        <div class="not-login-sub">
          登录后您可以编辑个人资料、查看发布的文章、与其他人互动
        </div>
        <div class="auth-actions">
          <router-link to="/login">
            <button class="btn btn-primary">
              <el-icon style="margin-right:4px"><SwitchButton /></el-icon>登录
            </button>
          </router-link>
          <router-link to="/register">
            <button class="btn btn-ghost">
              <el-icon style="margin-right:4px"><User /></el-icon>注册
            </button>
          </router-link>
        </div>
      </div>
    </div>

    <!-- 用户主页 -->
    <template v-else-if="user">
      <!-- 顶部封面 + 头像 -->
      <div class="cover-banner">
        <div class="cover-pattern"></div>
      </div>

      <div class="profile-body">
        <!-- 左侧资料卡片 -->
        <aside class="profile-side">
          <div class="card-surface profile-card">
            <div class="avatar-box">
              <UserAvatar :user="user" :size="110" />
              <button
                v-if="isOwnProfile"
                class="avatar-edit-btn"
                title="更换头像"
                @click="changeAvatar"
              >
                <el-icon :size="16"><Camera /></el-icon>
              </button>
            </div>

            <div class="ident">
              <div class="profile-name">{{ user.display_name }}</div>
              <div class="profile-username">@{{ user.username }}</div>
              <div v-if="user.role === 'admin'" class="admin-badge">
                <el-icon style="margin-right:4px"><Lock /></el-icon>
                管理员
              </div>
              <span v-else class="level-badge">
                <el-icon style="margin-right:4px"><Medal /></el-icon>
                {{ levelLabel }}
              </span>
            </div>

            <p v-if="user.signature" class="bio">{{ user.signature }}</p>

            <div class="action-row">
              <button
                v-if="isOwnProfile"
                class="btn btn-primary w-full"
                @click="openEditProfile"
              >
                <el-icon style="margin-right:6px"><EditPen /></el-icon>
                编辑资料
              </button>
              <button
                v-else
                class="btn w-full"
                :class="followStatus.is_following ? 'btn-ghost' : 'btn-primary'"
                @click="handleFollow"
              >
                <el-icon style="margin-right:6px">
                  <Check v-if="followStatus.is_following" /><Plus v-else />
                </el-icon>
                {{ followStatus.is_following ? '已关注' : followStatus.is_followed ? '回关' : '关注' }}
              </button>
            </div>

            <!-- 统计数据 -->
            <div class="stats-row">
              <div class="stat-cell cursor-pointer" @click="goToFollowing">
                <div class="stat-num">{{ followingCount }}</div>
                <div class="stat-label">关注</div>
              </div>
              <div class="stat-cell cursor-pointer" @click="goToFollowers">
                <div class="stat-num">{{ followersCount }}</div>
                <div class="stat-label">粉丝</div>
              </div>
              <div class="stat-cell">
                <div class="stat-num">{{ articleCount }}</div>
                <div class="stat-label">文章</div>
              </div>
            </div>

            <!-- 用户信息列表 -->
            <div class="info-list">
              <div class="info-item">
                <el-icon class="info-icon"><Message /></el-icon>
                <span>QQ号：{{ user.qq_number || '未设置' }}</span>
              </div>
              <div class="info-item">
                <el-icon class="info-icon"><Calendar /></el-icon>
                <span>注册时间：{{ formatDate(user.created_at) }}</span>
              </div>
            </div>

            <!-- 快捷入口 -->
            <div class="quick-nav">
              <router-link to="/collections" class="quick-link">
                <el-icon><Folder /></el-icon> 我的收藏夹
              </router-link>
              <router-link to="/follow-list" class="quick-link">
                <el-icon><UserFilled /></el-icon> 好友列表
              </router-link>
            </div>
          </div>
        </aside>

        <!-- 右侧内容 -->
        <main class="profile-main">
          <div class="card-surface main-card">
            <el-tabs v-if="isOwnProfile" v-model="activeTab" class="main-tabs">
              <el-tab-pane :label="`已发布 (${myArticles.length})`" name="published">
                <div v-if="myArticles.length === 0" class="empty-state">
                  <el-icon class="empty-icon" :size="64"><DocumentRemove /></el-icon>
                  <div class="empty-title">暂无文章</div>
                  <router-link to="/create">
                    <button class="btn btn-primary">
                      <el-icon style="margin-right:4px"><EditPen /></el-icon>
                      写文章
                    </button>
                  </router-link>
                </div>
                <div v-else>
                  <div v-for="article in myArticles" :key="article.id" class="article-row">
                    <el-avatar :size="44" :src="article.user?.avatar || ''" class="article-avatar">
                      {{ (article.user?.display_name || article.user?.username || 'U')[0]?.toUpperCase() }}
                    </el-avatar>
                    <div class="article-body">
                      <div class="article-title">
                        <router-link :to="'/article/' + article.id" class="article-link">
                          {{ article.title }}
                        </router-link>
                      </div>
                      <div class="article-meta">
                        <span><el-icon><Clock /></el-icon> {{ formatDate(article.created_at) }}</span>
                        <span class="like-count"><el-icon><Star /></el-icon> {{ article.like_count }}</span>
                        <span><el-icon><View /></el-icon> {{ article.view_count }}</span>
                        <el-tag v-if="article.category" size="small" effect="plain" class="cat-tag">
                          {{ article.category.name }}
                        </el-tag>
                      </div>
                    </div>
                    <div class="article-actions">
                      <router-link :to="'/create?id=' + article.id">
                        <button class="icon-btn" title="编辑"><el-icon><EditPen /></el-icon></button>
                      </router-link>
                      <button class="icon-btn danger" title="删除" @click="deleteArticle(article.id)">
                        <el-icon><Delete /></el-icon>
                      </button>
                    </div>
                  </div>
                </div>
              </el-tab-pane>

              <el-tab-pane :label="`草稿 (${drafts.length})`" name="drafts">
                <div v-if="drafts.length === 0" class="empty-state">
                  <el-icon class="empty-icon" :size="64"><DocumentRemove /></el-icon>
                  <div class="empty-title">暂无草稿</div>
                  <router-link to="/create">
                    <button class="btn btn-primary">
                      <el-icon style="margin-right:4px"><EditPen /></el-icon>
                      写文章
                    </button>
                  </router-link>
                </div>
                <div v-else>
                  <div v-for="article in drafts" :key="article.id" class="article-row">
                    <div class="draft-avatar"><el-icon :size="22"><DocumentRemove /></el-icon></div>
                    <div class="article-body">
                      <div class="article-title">{{ article.title }}</div>
                      <div class="article-meta">
                        <span><el-icon><Clock /></el-icon> 最后修改：{{ formatDate(article.updated_at) }}</span>
                        <el-tag v-if="article.category" size="small" effect="plain">
                          {{ article.category.name }}
                        </el-tag>
                      </div>
                    </div>
                    <div class="article-actions">
                      <router-link :to="'/create?id=' + article.id">
                        <button class="icon-btn" title="编辑"><el-icon><EditPen /></el-icon></button>
                      </router-link>
                      <button class="icon-btn success" title="发布" @click="publishDraft(article.id)">
                        <el-icon><Promotion /></el-icon>
                      </button>
                      <button class="icon-btn danger" title="删除" @click="deleteDraft(article.id)">
                        <el-icon><Delete /></el-icon>
                      </button>
                    </div>
                  </div>
                </div>
              </el-tab-pane>
            </el-tabs>

            <!-- 非本人：文章列表 -->
            <template v-if="!isOwnProfile">
              <div class="block-head">
                <el-icon><Document /></el-icon>
                <span>{{ isOwnProfile ? '我的文章' : '他的文章' }}</span>
              </div>
              <div v-if="myArticles.length === 0" class="empty-state">
                <el-icon class="empty-icon" :size="64"><DocumentRemove /></el-icon>
                <div class="empty-title">暂无文章</div>
              </div>
              <div v-else>
                <div v-for="article in myArticles" :key="article.id" class="article-row">
                  <el-avatar :size="44" :src="article.user?.avatar || ''" class="article-avatar">
                    {{ (article.user?.display_name || article.user?.username || 'U')[0]?.toUpperCase() }}
                  </el-avatar>
                  <div class="article-body">
                    <div class="article-title">
                      <router-link :to="'/article/' + article.id" class="article-link">
                        {{ article.title }}
                      </router-link>
                    </div>
                    <div class="article-meta">
                      <span><el-icon><Clock /></el-icon> {{ formatDate(article.created_at) }}</span>
                      <span class="like-count"><el-icon><Star /></el-icon> {{ article.like_count }}</span>
                      <span><el-icon><View /></el-icon> {{ article.view_count }}</span>
                      <el-tag v-if="article.category" size="small" effect="plain" class="cat-tag">
                        {{ article.category.name }}
                      </el-tag>
                    </div>
                  </div>
                </div>
              </div>
            </template>
          </div>
        </main>
      </div>
    </template>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  Lock, SwitchButton, User, Camera, Check, Plus, Message, EditPen,
  Calendar, Document, DocumentRemove, Clock, Star, View, Delete,
  Promotion, Medal, Folder, UserFilled
} from '@element-plus/icons-vue'
import api, { articleApi, userApi, friendApi, uploadApi } from '../api'
import UserAvatar from '../components/UserAvatar.vue'
import { useUserStore } from '../stores/user'
import { success, error, confirm, prompt } from '../utils/message'
import { jcOpenHtml, jcFieldsConfig, jcCloseAll } from '../utils/jcu'

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
    const levelLabel = computed(() => {
      if (!user.value) return '注册用户'
      return 'Lv.1 注册用户'
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
        followingCount.value = friendsRes.data.mutual_friends?.length || friendsRes.data.friends?.length || 0
        followersCount.value = friendsRes.data.mutual_friends?.length || friendsRes.data.friends?.length || 0
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

    // 打开"编辑资料" JCuPupw 表单弹窗
    const openEditProfile = () => {
      const cfg = jcFieldsConfig([
        {
          name: 'display_name',
          label: '显示名称',
          placeholder: '请输入显示名称',
          value: editForm.value.display_name,
          required: true
        },
        {
          name: 'signature',
          label: '个性化签名',
          type: 'textarea',
          rows: 4,
          placeholder: '最多200个字符',
          value: editForm.value.signature
        }
      ])
      jcOpenHtml({
        title: '编辑资料',
        content: cfg.html,
        width: 480,
        size: 'sm',
        buttons: [
          { text: '取消', type: 'default', action: () => jcCloseAll() },
          {
            text: '保存修改',
            type: 'primary',
            action: () => {
              if (!cfg.validate(document)) return
              const v = cfg.collect(document)
              updateProfile(v.display_name, v.signature)
              jcCloseAll()
            }
          }
        ]
      })
    }

    const updateProfile = async (displayName, signature) => {
      saving.value = true
      editForm.value.display_name = displayName
      editForm.value.signature = signature
      try {
        await api.put('/profile', {
          display_name: displayName ?? '',
          signature: signature ?? ''
        })
        user.value.display_name = displayName ?? ''
        user.value.signature = signature ?? ''
        if (userStore.user) {
          const next = { ...userStore.user, display_name: displayName ?? '', signature: signature ?? '' }
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
      openEditProfile,
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
      levelLabel,
      loading,
      Lock, SwitchButton, User, Camera, Check, Plus, Message, EditPen,
      Calendar, Document, DocumentRemove, Clock, Star, View, Delete,
      Promotion, Medal, Folder, UserFilled
    }
  }
}
</script>

<style scoped>
.profile-page {
  padding-bottom: 40px;
}

/* ---------------- loading / auth ---------------- */
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
  padding: 40px 32px;
  text-align: center;
  background: linear-gradient(160deg, var(--campus-surface) 0%, var(--campus-surface-2) 100%);
  border-radius: var(--campus-radius-lg);
}

.not-login-icon {
  display: flex;
  justify-content: center;
  color: var(--campus-primary);
  margin-bottom: 16px;
}

.not-login-title {
  font-size: 20px;
  font-weight: 700;
  margin-bottom: 8px;
  color: var(--campus-text);
}

.not-login-sub {
  font-size: 14px;
  color: var(--campus-text-secondary);
  margin-bottom: 24px;
}

.auth-actions {
  display: flex;
  gap: 12px;
  justify-content: center;
}

/* ---------------- 覆盖横幅 ---------------- */
.cover-banner {
  position: relative;
  height: 200px;
  border-radius: var(--campus-radius-lg);
  background: linear-gradient(120deg, #4f6ef7 0%, #6d8bfb 45%, #9db4ff 100%);
  overflow: hidden;
  box-shadow: var(--campus-shadow);
}

.cover-pattern {
  position: absolute;
  inset: 0;
  background:
    radial-gradient(circle at 75% 20%, rgba(255, 255, 255, 0.25) 0, transparent 35%),
    radial-gradient(circle at 15% 90%, rgba(255, 255, 255, 0.18) 0, transparent 40%);
}

/* ---------------- 主体布局 ---------------- */
.profile-body {
  display: grid;
  grid-template-columns: 300px 1fr;
  gap: 20px;
  margin-top: -56px;
  padding: 0 16px;
}

.profile-side {
  min-width: 0;
}

.profile-main {
  min-width: 0;
}

/* ---------------- 资料卡片 ---------------- */
.profile-card {
  padding: 24px;
  border-radius: var(--campus-radius-lg);
  box-shadow: var(--campus-shadow);
}

.avatar-box {
  position: relative;
  display: flex;
  justify-content: flex-start;
  margin-top: -86px;
}

.avatar-box :deep(.user-avatar .avatar) {
  border: 4px solid var(--campus-surface);
  box-shadow: var(--campus-shadow);
  background: var(--campus-primary);
}

.avatar-edit-btn {
  position: absolute;
  bottom: 4px;
  right: 8px;
  width: 30px;
  height: 30px;
  border-radius: 50%;
  border: none;
  background: var(--campus-primary);
  color: #fff;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: var(--campus-shadow-sm);
  transition: var(--campus-transition);
}

.avatar-edit-btn:hover {
  background: var(--campus-primary-dark);
  transform: scale(1.05);
}

.ident {
  margin-top: 14px;
  text-align: left;
}

.profile-name {
  font-size: 24px;
  font-weight: 800;
  color: var(--campus-text);
  letter-spacing: -0.01em;
}

.profile-username {
  font-size: 14px;
  color: var(--campus-text-secondary);
  margin-top: 2px;
}

.admin-badge {
  display: inline-flex;
  align-items: center;
  margin-top: 8px;
  padding: 2px 10px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 600;
  color: #fff;
  background: linear-gradient(120deg, #ef4444, #f59e0b);
}

.level-badge {
  display: inline-flex;
  align-items: center;
  margin-top: 8px;
  padding: 2px 10px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 600;
  color: var(--campus-primary);
  background: var(--campus-primary-soft);
  border: 1px solid rgba(79, 110, 247, 0.2);
}

.bio {
  margin: 14px 0 0;
  padding: 12px;
  background: var(--campus-surface-2);
  border-radius: var(--campus-radius-sm);
  font-size: 13px;
  line-height: 1.6;
  color: var(--campus-text);
  white-space: pre-line;
  word-break: break-word;
}

.action-row {
  margin-top: 16px;
}

/* ---------------- 统计 ---------------- */
.stats-row {
  display: flex;
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px solid var(--campus-border);
}

.stat-cell {
  flex: 1;
  text-align: center;
}

.stat-num {
  font-size: 22px;
  font-weight: 800;
  color: var(--campus-text);
}

.stat-label {
  font-size: 12px;
  color: var(--campus-text-secondary);
  margin-top: 2px;
}

/* ---------------- 信息列表 ---------------- */
.info-list {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid var(--campus-border);
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: var(--campus-text);
}

.info-icon {
  color: var(--campus-primary);
  flex-shrink: 0;
}

/* ---------------- 快捷导航 ---------------- */
.quick-nav {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid var(--campus-border);
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.quick-link {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  border-radius: var(--campus-radius-sm);
  background: var(--campus-surface-2);
  font-size: 13px;
  font-weight: 600;
  color: var(--campus-text);
  transition: var(--campus-transition);
}

.quick-link:hover {
  background: var(--campus-primary-soft);
  color: var(--campus-primary);
}

/* ---------------- 按钮 ---------------- */
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  padding: 9px 18px;
  border-radius: 12px;
  border: none;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: var(--campus-transition);
}

.btn-primary {
  background: var(--campus-primary);
  color: #fff;
}

.btn-primary:hover {
  background: var(--campus-primary-dark);
  transform: translateY(-1px);
  box-shadow: var(--campus-shadow);
}

.btn-ghost {
  background: var(--campus-surface);
  color: var(--campus-primary);
  border: 1px solid var(--campus-border);
}

.btn-ghost:hover {
  border-color: var(--campus-primary);
  background: var(--campus-primary-soft);
}

.w-full {
  width: 100%;
}

/* ---------------- 主卡片 ---------------- */
.main-card {
  border-radius: var(--campus-radius-lg);
  padding: 24px;
  box-shadow: var(--campus-shadow);
}

.block-head {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 700;
  color: var(--campus-text);
  margin-bottom: 16px;
}

.main-tabs :deep(.el-tabs__item) {
  font-weight: 600;
  font-size: 15px;
}

.main-tabs :deep(.el-tabs__active-bar) {
  height: 3px;
  border-radius: 3px;
}

/* ---------------- 空状态 ---------------- */
.empty-state {
  text-align: center;
  padding: 48px 16px;
}

.empty-icon {
  color: var(--campus-border);
  margin-bottom: 12px;
}

.empty-title {
  font-size: 15px;
  color: var(--campus-text-secondary);
  margin-bottom: 16px;
}

/* ---------------- 文章行 ---------------- */
.article-row {
  display: flex;
  align-items: flex-start;
  gap: 14px;
  padding: 16px 0;
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

.draft-avatar {
  width: 44px;
  height: 44px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 12px;
  background: var(--campus-primary-soft);
  color: var(--campus-primary);
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

.article-link {
  color: var(--campus-text);
  transition: var(--campus-transition);
}

.article-link:hover {
  color: var(--campus-primary);
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

.like-count {
  color: var(--campus-danger);
}

.cat-tag {
  margin-left: 2px;
}

.article-actions {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  gap: 4px;
}

.icon-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: transparent;
  color: var(--campus-text-secondary);
  border-radius: 8px;
  cursor: pointer;
  transition: var(--campus-transition);
}

.icon-btn:hover {
  background: var(--campus-primary-soft);
  color: var(--campus-primary);
}

.icon-btn.danger:hover {
  background: rgba(239, 68, 68, 0.1);
  color: var(--campus-danger);
}

.icon-btn.success:hover {
  background: rgba(34, 197, 94, 0.12);
  color: var(--campus-success);
}

/* ---------------- 响应式 ---------------- */
@media (max-width: 900px) {
  .profile-body {
    grid-template-columns: 1fr;
    padding: 0 4px;
  }

  .profile-side {
    max-width: 520px;
    margin: 0 auto;
    width: 100%;
  }

  .cover-banner {
    height: 150px;
  }

  .main-card {
    padding: 16px;
  }
}

@media (max-width: 480px) {
  .profile-body {
    margin-top: -40px;
  }

  .avatar-box {
    margin-top: -68px;
  }

  .cover-banner {
    height: 120px;
  }

  .profile-card {
    padding: 18px;
  }

  .main-card {
    padding: 12px;
  }
}
</style>
