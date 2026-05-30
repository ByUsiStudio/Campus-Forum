<template>
  <v-app v-if="isInitialized && isAdmin">
    <v-app-bar flat class="admin-header" height="72">
      <div class="header-content">
        <div class="header-brand">
          <v-icon size="28" color="primary" class="brand-icon">mdi-shield-crown</v-icon>
          <div class="brand-text">
            <div class="brand-title">管理后台</div>
            <div class="brand-subtitle">校园论坛管理中心</div>
          </div>
        </div>

        <v-tabs v-model="activeTab" class="header-tabs" slider-color="primary">
          <v-tab value="overview" class="tab-item">
            <v-icon size="20" class="tab-icon">mdi-view-dashboard</v-icon>
            <span>概览</span>
          </v-tab>
          <v-tab value="users" class="tab-item">
            <v-icon size="20" class="tab-icon">mdi-account-group</v-icon>
            <span>用户</span>
          </v-tab>
          <v-tab value="articles" class="tab-item">
            <v-icon size="20" class="tab-icon">mdi-file-document-edit</v-icon>
            <span>文章</span>
          </v-tab>
          <v-tab value="comments" class="tab-item">
            <v-icon size="20" class="tab-icon">mdi-comment-text-multiple</v-icon>
            <span>评论</span>
          </v-tab>
          <v-tab value="categories" class="tab-item">
            <v-icon size="20" class="tab-icon">mdi-shape</v-icon>
            <span>分区</span>
          </v-tab>
          <v-tab value="titles" class="tab-item">
            <v-icon size="20" class="tab-icon">mdi-medal</v-icon>
            <span>头衔</span>
          </v-tab>
          <v-tab value="content" class="tab-item">
            <v-icon size="20" class="tab-icon">mdi-web</v-icon>
            <span>内容</span>
            <v-chip v-if="deletionRequests.length" size="x-small" color="error" class="tab-badge">
              {{ deletionRequests.length }}
            </v-chip>
          </v-tab>
          <v-tab value="system" class="tab-item">
            <v-icon size="20" class="tab-icon">mdi-cog</v-icon>
            <span>系统</span>
          </v-tab>
        </v-tabs>

        <div class="header-actions">
          <v-btn icon variant="text" size="small" @click="goToHome">
            <v-icon>mdi-home</v-icon>
          </v-btn>
          <v-btn icon variant="text" size="small" @click="handleRefresh">
            <v-icon :class="{ 'spin': isRefreshing }">mdi-refresh</v-icon>
          </v-btn>
        </div>
      </div>
    </v-app-bar>

    <v-main class="admin-main">
      <v-container fluid class="pa-6">
        <v-window v-model="activeTab">
          <v-window-item value="overview">
            <OverviewPanel :statistics="statistics" :loading="loading" @refresh="loadStatistics" />
          </v-window-item>

          <v-window-item value="users">
            <UsersPanel
              :users="users"
              :loading="loading"
              :current-user-id="currentUserId"
              :current-user-role="currentUserRole"
              @edit-role="showEditRoleDialog"
              @edit-user="showEditUserDialog"
              @ban="showBanDialog"
              @unban="handleUnban"
              @delete="handleDeleteUser"
              @refresh="loadUsers"
            />
          </v-window-item>

          <v-window-item value="articles">
            <ArticlesPanel
              :articles="articles"
              :loading="loading"
              :page="articlePage"
              :total-pages="articleTotalPages"
              :filter="articleFilter"
              :status-options="articleStatusOptions"
              :current-user-role="currentUserRole"
              @change-status="showStatusDialog"
              @delete="handleDeleteArticle"
              @refresh="loadArticles"
              @update:page="articlePage = $event; loadArticles()"
              @update:filter="articleFilter = $event; articlePage = 1; loadArticles()"
            />
          </v-window-item>

          <v-window-item value="comments">
            <CommentsPanel
              :comments="allComments"
              :loading="loading"
              @delete="handleDeleteComment"
              @refresh="loadComments"
            />
          </v-window-item>

          <v-window-item value="categories">
            <CategoriesPanel
              :categories="categories"
              :loading="loading"
              @add="addCategory"
              @edit="showEditCategoryDialog"
              @delete="handleDeleteCategory"
              @refresh="loadCategories"
            />
          </v-window-item>

          <v-window-item value="titles">
            <TitlesPanel
              :titles="titles"
              :users="users"
              :loading="loading"
              @add-title="addTitle"
              @grant="grantTitle"
              @revoke="revokeTitle"
              @delete-title="handleDeleteTitle"
              @refresh="loadTitles"
            />
          </v-window-item>

          <v-window-item value="content">
            <ContentPanel
              :sidebar-items="sidebarItems"
              :deletion-requests="deletionRequests"
              :announcement="announcementContent"
              @add-sidebar-item="addSidebarItem"
              @remove-sidebar-item="removeSidebarItem"
              @save-sidebar="saveSidebarConfig"
              @approve-deletion="approveDeletion"
              @reject-deletion="rejectDeletion"
              @save-announcement="saveAnnouncement"
              @refresh-sidebar="loadSidebarConfig"
              @refresh-deletions="loadDeletionRequests"
              @refresh-announcement="loadAnnouncement"
            />
          </v-window-item>

          <v-window-item value="system">
            <SystemPanel
              :site-config="siteConfigForm"
              :smtp-config="smtpConfigForm"
              :notifications="notifications"
              :notification-form="notificationForm"
              :notification-types="notificationTypes"
              @save-site="saveSiteConfig"
              @save-smtp="saveSmtpConfig"
              @test-smtp="testSmtpConfig"
              @send-notification="handleSendNotification"
              @delete-notification="deleteNotification"
              @refresh-site="loadSiteConfig"
              @refresh-smtp="loadSmtpConfig"
              @refresh-notifications="loadNotifications"
            />
          </v-window-item>
        </v-window>
      </v-container>
    </v-main>

    <v-dialog v-model="editRoleDialog.show" max-width="480">
      <v-card class="dialog-card">
        <v-card-title class="dialog-title">
          <v-icon class="title-icon">mdi-account-edit</v-icon>
          修改用户角色
        </v-card-title>
        <v-card-text class="dialog-body">
          <div class="user-info-card">
            <UserAvatar :user="editRoleDialog.user || {}" :size="48" />
            <div class="user-info-text">
              <div class="user-name">{{ editRoleDialog.user?.display_name }}</div>
              <div class="user-meta">ID: {{ editRoleDialog.user?.id }}</div>
            </div>
          </div>
          <v-radio-group v-model="editRoleDialog.role" class="mt-4">
            <v-radio label="普通用户" value="user" color="primary"></v-radio>
            <v-radio label="管理员" value="admin" color="error"></v-radio>
            <v-radio label="系统管理员" value="system" color="warning"></v-radio>
          </v-radio-group>
        </v-card-text>
        <v-card-actions class="dialog-actions">
          <v-btn variant="text" @click="editRoleDialog.show = false">取消</v-btn>
          <v-btn color="primary" variant="flat" @click="handleEditRole">确认修改</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="banDialog.show" max-width="480">
      <v-card class="dialog-card">
        <v-card-title class="dialog-title">
          <v-icon class="title-icon text-error">mdi-gavel</v-icon>
          封禁用户
        </v-card-title>
        <v-card-text class="dialog-body">
          <div class="user-info-card">
            <UserAvatar :user="banDialog.user || {}" :size="48" />
            <div class="user-info-text">
              <div class="user-name">{{ banDialog.user?.display_name }}</div>
              <div class="user-meta">即将被封禁</div>
            </div>
          </div>
          <v-textarea
            v-model="banDialog.reason"
            label="封禁原因"
            placeholder="请输入封禁原因..."
            variant="outlined"
            rows="3"
            class="mt-4"
          ></v-textarea>
        </v-card-text>
        <v-card-actions class="dialog-actions">
          <v-btn variant="text" @click="banDialog.show = false">取消</v-btn>
          <v-btn color="error" variant="flat" @click="handleBan">确认封禁</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="editUserDialog.show" max-width="480">
      <v-card class="dialog-card">
        <v-card-title class="dialog-title">
          <v-icon class="title-icon">mdi-account-edit</v-icon>
          编辑用户
        </v-card-title>
        <v-card-text class="dialog-body">
          <div class="user-info-card">
            <UserAvatar :user="editUserDialog.user || {}" :size="48" />
            <div class="user-info-text">
              <div class="user-name">{{ editUserDialog.user?.display_name }}</div>
              <div class="user-meta">ID: {{ editUserDialog.user?.id }}</div>
            </div>
          </div>
          <v-text-field
            v-model="editUserDialog.displayName"
            label="显示名称"
            variant="outlined"
            class="mt-4 mb-4"
          ></v-text-field>
          <v-text-field
            v-model="editUserDialog.email"
            label="邮箱"
            variant="outlined"
          ></v-text-field>
        </v-card-text>
        <v-card-actions class="dialog-actions">
          <v-btn variant="text" @click="editUserDialog.show = false">取消</v-btn>
          <v-btn color="primary" variant="flat" @click="handleEditUser">保存</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="statusDialog.show" max-width="480">
      <v-card class="dialog-card">
        <v-card-title class="dialog-title">
          <v-icon class="title-icon">mdi-file-edit</v-icon>
          修改文章状态
        </v-card-title>
        <v-card-text class="dialog-body">
          <div class="article-preview">
            <div class="preview-label">文章</div>
            <div class="preview-title">{{ statusDialog.article?.title }}</div>
          </div>
          <v-select
            v-model="statusDialog.status"
            :items="[
              { title: '待审核', value: 'pending' },
              { title: '已发布', value: 'published' },
              { title: '已拒绝', value: 'rejected' }
            ]"
            label="选择状态"
            variant="outlined"
            class="mt-4"
          ></v-select>
        </v-card-text>
        <v-card-actions class="dialog-actions">
          <v-btn variant="text" @click="statusDialog.show = false">取消</v-btn>
          <v-btn color="primary" variant="flat" @click="handleEditStatus">确认</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="editCategoryDialog.show" max-width="480">
      <v-card class="dialog-card">
        <v-card-title class="dialog-title">
          <v-icon class="title-icon">mdi-folder-edit</v-icon>
          编辑分区
        </v-card-title>
        <v-card-text class="dialog-body">
          <v-text-field
            v-model="editCategoryDialog.name"
            label="分区名称"
            variant="outlined"
            class="mb-4"
          ></v-text-field>
          <v-text-field
            v-model="editCategoryDialog.description"
            label="描述"
            variant="outlined"
          ></v-text-field>
        </v-card-text>
        <v-card-actions class="dialog-actions">
          <v-btn variant="text" @click="editCategoryDialog.show = false">取消</v-btn>
          <v-btn color="primary" variant="flat" @click="handleEditCategory">保存</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-app>

  <v-app v-else-if="isInitialized && !isAdmin">
    <v-main class="error-page">
      <v-container class="fill-height">
        <v-row align="center" justify="center">
          <v-col cols="12" sm="8" md="5" lg="4">
            <div class="error-content">
              <v-icon size="80" color="error" class="error-icon">mdi-lock</v-icon>
              <h2 class="error-title">访问受限</h2>
              <p class="error-message">您没有权限访问管理后台</p>
              <v-btn color="primary" size="large" @click="goToHome" class="mt-6">
                返回首页
              </v-btn>
            </div>
          </v-col>
        </v-row>
      </v-container>
    </v-main>
  </v-app>
</template>

<script>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import api from '../api'
import UserAvatar from '../components/UserAvatar.vue'
import { confirm as showConfirm, success as showSuccess, error as showError } from '../utils/modal'
import OverviewPanel from './admin/OverviewPanel.vue'
import UsersPanel from './admin/UsersPanel.vue'
import ArticlesPanel from './admin/ArticlesPanel.vue'
import CommentsPanel from './admin/CommentsPanel.vue'
import CategoriesPanel from './admin/CategoriesPanel.vue'
import TitlesPanel from './admin/TitlesPanel.vue'
import ContentPanel from './admin/ContentPanel.vue'
import SystemPanel from './admin/SystemPanel.vue'

export default {
  name: 'Admin',
  components: {
    UserAvatar,
    OverviewPanel,
    UsersPanel,
    ArticlesPanel,
    CommentsPanel,
    CategoriesPanel,
    TitlesPanel,
    ContentPanel,
    SystemPanel
  },
  setup() {
    const router = useRouter()
    const activeTab = ref('overview')
    const isAdmin = ref(false)
    const currentUserId = ref(null)
    const currentUserRole = ref(null)
    const isInitialized = ref(false)
    const loading = ref(false)
    const isRefreshing = ref(false)
    const deletionRequests = ref([])

    const statistics = ref({
      user_count: 0,
      article_count: 0,
      comment_count: 0,
      view_count: 0
    })

    const users = ref([])
    const articles = ref([])
    const articlePage = ref(1)
    const articleTotalPages = ref(1)
    const articleFilter = ref('')
    const articleStatusOptions = [
      { title: '全部', value: '' },
      { title: '待审核', value: 'pending' },
      { title: '已发布', value: 'published' },
      { title: '已拒绝', value: 'rejected' }
    ]

    const comments = ref([])
    const allComments = computed(() => comments.value)

    const categories = ref([])
    const titles = ref([])
    const sidebarItems = ref([])
    const announcementContent = ref('')
    const siteConfigForm = ref({ siteTitle: '' })
    const smtpConfigForm = ref({
      host: '',
      port: 465,
      username: '',
      password: '',
      from: '',
      fromName: '',
      ssl: true
    })
    const notifications = ref([])
    const notificationForm = ref({
      type: 'system',
      target: 'all',
      title: '',
      content: ''
    })
    const notificationTypes = [
      { title: '系统通知', value: 'system' },
      { title: '活动通知', value: 'activity' },
      { title: '更新通知', value: 'update' },
      { title: '警告通知', value: 'warning' },
      { title: '全部用户', value: 'all' },
      { title: '仅管理员', value: 'admin' }
    ]

    const categoryForm = ref({
      name: '',
      description: '',
      color: '#6750A4'
    })
    const titleForm = ref({
      name: '',
      description: '',
      color: '#6750A4',
      icon: ''
    })
    const grantForm = ref({
      user_id: null,
      title_id: null,
      reason: ''
    })

    const editRoleDialog = ref({
      show: false,
      user: null,
      role: 'user'
    })
    const banDialog = ref({
      show: false,
      user: null,
      reason: ''
    })
    const editUserDialog = ref({
      show: false,
      user: null,
      displayName: '',
      email: ''
    })
    const statusDialog = ref({
      show: false,
      article: null,
      status: 'published'
    })
    const editCategoryDialog = ref({
      show: false,
      category: null,
      name: '',
      description: ''
    })

    const checkAdmin = async () => {
      try {
        const token = localStorage.getItem('token')
        if (!token) {
          router.push('/login')
          return
        }
        const response = await api.get('/profile')
        const userData = response.data
        const userRole = userData.role || userData.user?.role

        if (userRole !== 'admin' && userRole !== 'system' && userRole !== 'Admin' && userRole !== 'System') {
          router.push('/')
          return
        }
        currentUserId.value = userData.id || userData.user?.id
        currentUserRole.value = userRole
        isAdmin.value = true
        isInitialized.value = true
      } catch (error) {
        console.error('检查管理员权限失败', error)
        localStorage.removeItem('token')
        router.push('/login')
      }
    }

    const goToHome = () => {
      router.push('/')
    }

    const handleRefresh = async () => {
      isRefreshing.value = true
      await loadDataForCurrentTab()
      setTimeout(() => {
        isRefreshing.value = false
      }, 500)
    }

    const loadDataForCurrentTab = async () => {
      loading.value = true
      try {
        switch (activeTab.value) {
          case 'overview': await loadStatistics(); break
          case 'users': await loadUsers(); break
          case 'articles': await loadArticles(); break
          case 'comments': await loadComments(); break
          case 'categories': await loadCategories(); break
          case 'titles': await loadTitles(); break
          case 'content':
            await Promise.all([loadSidebarConfig(), loadDeletionRequests(), loadAnnouncement()])
            break
          case 'system':
            await Promise.all([loadSiteConfig(), loadSmtpConfig(), loadNotifications()])
            break
        }
      } finally {
        loading.value = false
      }
    }

    const loadStatistics = async () => {
      try {
        const response = await api.get('/admin/statistics')
        statistics.value = response.data
      } catch (error) {
        console.error('加载统计失败', error)
      }
    }

    const loadUsers = async () => {
      try {
        const response = await api.get('/admin/users')
        users.value = response.data.users
      } catch (error) {
        console.error('加载用户失败', error)
      }
    }

    const loadArticles = async () => {
      try {
        const params = { page: articlePage.value, page_size: 20 }
        if (articleFilter.value) {
          params.status = articleFilter.value
        }
        const response = await api.get('/admin/articles', { params })
        articles.value = response.data.articles
        articleTotalPages.value = response.data.total_pages
      } catch (error) {
        console.error('加载文章失败', error)
      }
    }

    const loadComments = async () => {
      try {
        const response = await api.get('/admin/comments', { params: { page: 1, page_size: 50 } })
        comments.value = response.data.comments
      } catch (error) {
        console.error('加载评论失败', error)
      }
    }

    const loadCategories = async () => {
      try {
        const response = await api.get('/categories')
        categories.value = response.data.categories || []
      } catch (error) {
        console.error('加载分区失败', error)
      }
    }

    const loadTitles = async () => {
      try {
        const response = await api.get('/titles')
        titles.value = response.data.titles || []
      } catch (error) {
        console.error('加载头衔失败', error)
      }
    }

    const loadSidebarConfig = async () => {
      try {
        const response = await api.get('/sidebar-config')
        sidebarItems.value = response.data.items || []
      } catch (error) {
        console.error('加载侧边栏配置失败', error)
      }
    }

    const loadAnnouncement = async () => {
      try {
        const response = await api.get('/announcement')
        announcementContent.value = response.data.content || ''
      } catch (error) {
        console.error('加载公告失败', error)
      }
    }

    const loadSiteConfig = async () => {
      try {
        const response = await api.get('/site-config')
        siteConfigForm.value.siteTitle = response.data.site_title || response.data.SiteTitle || ''
      } catch (error) {
        console.error('加载网站配置失败', error)
      }
    }

    const loadSmtpConfig = async () => {
      try {
        const response = await api.get('/site-config')
        smtpConfigForm.value = {
          host: response.data.smtp_host || response.data.SMTPHost || '',
          port: response.data.smtp_port || response.data.SMTPPort || 587,
          username: response.data.smtp_username || response.data.SMTPUsername || '',
          password: response.data.smtp_password || response.data.SMTPPassword || '',
          from: response.data.smtp_from || response.data.SMTPFrom || '',
          fromName: response.data.smtp_from_name || response.data.SMTPFromName || '',
          ssl: response.data.smtp_port == 465
        }
      } catch (error) {
        console.error('加载SMTP配置失败', error)
      }
    }

    const loadNotifications = async () => {
      try {
        const response = await api.get('/notifications/admin')
        notifications.value = response.data.notifications || []
      } catch (error) {
        console.error('加载通知失败', error)
      }
    }

    const loadDeletionRequests = async () => {
      try {
        const response = await api.get('/deletion-requests')
        deletionRequests.value = response.data.requests || []
      } catch (error) {
        console.error('加载删除请求失败', error)
      }
    }

    const showEditRoleDialog = (user) => {
      editRoleDialog.value = {
        show: true,
        user,
        role: user.role
      }
    }

    const showBanDialog = (user) => {
      banDialog.value = {
        show: true,
        user,
        reason: ''
      }
    }

    const showEditUserDialog = (user) => {
      editUserDialog.value = {
        show: true,
        user,
        displayName: user.display_name || '',
        email: user.email || ''
      }
    }

    const showStatusDialog = (article) => {
      statusDialog.value = {
        show: true,
        article,
        status: article.status
      }
    }

    const showEditCategoryDialog = (category) => {
      editCategoryDialog.value = {
        show: true,
        category,
        name: category.name,
        description: category.description || ''
      }
    }

    const handleEditRole = async () => {
      try {
        await api.put(`/admin/users/${editRoleDialog.value.user.id}/role`, { role: editRoleDialog.value.role })
        showSuccess('修改成功')
        editRoleDialog.value.show = false
        loadUsers()
      } catch (error) {
        console.error('修改角色失败', error)
        showError(error.response?.data?.error || '修改失败')
      }
    }

    const handleBan = async () => {
      if (!banDialog.value.reason) {
        showError('请输入封禁原因')
        return
      }
      try {
        await api.post(`/admin/users/${banDialog.value.user.id}/ban`, { reason: banDialog.value.reason })
        showSuccess('封禁成功')
        banDialog.value.show = false
        loadUsers()
      } catch (error) {
        console.error('封禁用户失败', error)
        showError(error.response?.data?.error || '封禁失败')
      }
    }

    const handleEditUser = async () => {
      if (!editUserDialog.value.displayName) {
        showError('请输入显示名称')
        return
      }
      try {
        await api.put(`/admin/users/${editUserDialog.value.user.id}`, {
          display_name: editUserDialog.value.displayName,
          email: editUserDialog.value.email
        })
        showSuccess('保存成功')
        editUserDialog.value.show = false
        loadUsers()
      } catch (error) {
        console.error('编辑用户失败', error)
        showError(error.response?.data?.error || '保存失败')
      }
    }

    const handleUnban = async (user) => {
      const confirmed = await showConfirm(`确定要解封用户 "${user.display_name}" 吗？`)
      if (!confirmed) return
      try {
        await api.post(`/admin/users/${user.id}/unban`)
        showSuccess('解封成功')
        loadUsers()
      } catch (error) {
        console.error('解封用户失败', error)
        showError(error.response?.data?.error || '解封失败')
      }
    }

    const handleDeleteUser = async (user) => {
      const confirmed = await showConfirm(`确定要删除用户 "${user.display_name}" 吗？`)
      if (!confirmed) return
      try {
        await api.delete(`/admin/users/${user.id}`)
        showSuccess('删除成功')
        loadUsers()
      } catch (error) {
        console.error('删除用户失败', error)
        showError(error.response?.data?.error || '删除失败')
      }
    }

    const handleEditStatus = async () => {
      try {
        await api.put(`/admin/articles/${statusDialog.value.article.id}/status`, { status: statusDialog.value.status })
        showSuccess('修改成功')
        statusDialog.value.show = false
        loadArticles()
      } catch (error) {
        console.error('修改状态失败', error)
        showError(error.response?.data?.error || '修改失败')
      }
    }

    const handleDeleteArticle = async (article) => {
      const confirmed = await showConfirm(`确定要删除文章 "${article.title}" 吗？`)
      if (!confirmed) return
      try {
        await api.delete(`/admin/articles/${article.id}`)
        showSuccess('删除成功')
        loadArticles()
      } catch (error) {
        console.error('删除文章失败', error)
        showError(error.response?.data?.error || '删除失败')
      }
    }

    const handleDeleteComment = async (commentId) => {
      const confirmed = await showConfirm('确定要删除这条评论吗？')
      if (!confirmed) return
      try {
        await api.delete(`/admin/comments/${commentId}`)
        showSuccess('删除成功')
        loadComments()
      } catch (error) {
        console.error('删除评论失败', error)
        showError(error.response?.data?.error || '删除失败')
      }
    }

    const addCategory = async (formData) => {
      const catData = formData || categoryForm.value
      if (!catData.name) {
        showError('请输入分区名称')
        return
      }
      try {
        await api.post('/categories', catData)
        showSuccess('添加成功')
        categoryForm.value = { name: '', description: '', color: '#6750A4' }
        loadCategories()
      } catch (error) {
        console.error('添加分区失败', error)
        showError(error.response?.data?.error || '添加失败')
      }
    }

    const handleEditCategory = async () => {
      try {
        await api.put(`/categories/${editCategoryDialog.value.category.id}`, {
          name: editCategoryDialog.value.name,
          description: editCategoryDialog.value.description
        })
        showSuccess('保存成功')
        editCategoryDialog.value.show = false
        loadCategories()
      } catch (error) {
        console.error('保存分区失败', error)
        showError(error.response?.data?.error || '保存失败')
      }
    }

    const handleDeleteCategory = async (id) => {
      const confirmed = await showConfirm('确定要删除这个分区吗？')
      if (!confirmed) return
      try {
        await api.delete(`/categories/${id}`)
        showSuccess('删除成功')
        loadCategories()
      } catch (error) {
        console.error('删除分区失败', error)
        showError(error.response?.data?.error || '删除失败')
      }
    }

    const addTitle = async (formData) => {
      const titleData = formData || titleForm.value
      if (!titleData.name) {
        showError('请输入头衔名称')
        return
      }
      try {
        await api.post('/titles', titleData)
        showSuccess('添加成功')
        titleForm.value = { name: '', description: '', color: '#6750A4', icon: '' }
        loadTitles()
      } catch (error) {
        console.error('添加头衔失败', error)
        showError(error.response?.data?.error || '添加失败')
      }
    }

    const grantTitle = async (formData) => {
      const grantData = formData || grantForm.value
      if (!grantData.user_id || !grantData.title_id) {
        showError('请选择用户和头衔')
        return
      }
      try {
        await api.post('/titles/grant', grantData)
        showSuccess('授予成功')
        grantForm.value = { user_id: null, title_id: null, reason: '' }
        loadUsers()
      } catch (error) {
        console.error('授予头衔失败', error)
        showError(error.response?.data?.error || '授予失败')
      }
    }

    const revokeTitle = async (userId, titleId) => {
      const confirmed = await showConfirm('确定要撤销该用户的头衔吗？')
      if (!confirmed) return
      try {
        await api.post('/titles/revoke', { user_id: userId, title_id: titleId })
        showSuccess('撤销成功')
        loadUsers()
      } catch (error) {
        console.error('撤销头衔失败', error)
        showError(error.response?.data?.error || '撤销失败')
      }
    }

    const handleDeleteTitle = async (id) => {
      const confirmed = await showConfirm('确定要删除这个头衔吗？')
      if (!confirmed) return
      try {
        await api.delete(`/titles/${id}`)
        showSuccess('删除成功')
        loadTitles()
      } catch (error) {
        console.error('删除头衔失败', error)
        showError(error.response?.data?.error || '删除失败')
      }
    }

    const addSidebarItem = () => {
      sidebarItems.value.push({ title: '', link: '', icon: 'mdi-link' })
    }

    const removeSidebarItem = (index) => {
      sidebarItems.value.splice(index, 1)
    }

    const saveSidebarConfig = async () => {
      try {
        await api.put('/sidebar-config', { items: sidebarItems.value })
        showSuccess('保存成功')
      } catch (error) {
        console.error('保存侧边栏配置失败', error)
        showError(error.response?.data?.error || '保存失败')
      }
    }

    const saveAnnouncement = async () => {
      try {
        await api.put('/announcement', { content: announcementContent.value })
        showSuccess('保存成功')
      } catch (error) {
        console.error('保存公告失败', error)
        showError(error.response?.data?.error || '保存失败')
      }
    }

    const saveSiteConfig = async (formData) => {
      const config = formData || siteConfigForm.value
      try {
        await api.put('/site-config', { site_title: config.siteTitle })
        showSuccess('保存成功')
        if (config.siteTitle) {
          document.title = config.siteTitle
        }
      } catch (error) {
        console.error('保存网站配置失败', error)
        showError(error.response?.data?.error || '保存失败')
      }
    }

    const saveSmtpConfig = async (formData) => {
      const config = formData || smtpConfigForm.value
      try {
        await api.put('/site-config', {
          smtp_host: config.host,
          smtp_port: config.port,
          smtp_username: config.username,
          smtp_password: config.password,
          smtp_from: config.from,
          smtp_from_name: config.fromName
        })
        showSuccess('保存成功')
      } catch (error) {
        console.error('保存SMTP配置失败', error)
        showError(error.response?.data?.error || '保存失败')
      }
    }

    const testSmtpConfig = async (formData) => {
      const config = formData || smtpConfigForm.value
      try {
        await api.post('/site-config/test-smtp', {
          smtp_host: config.host,
          smtp_port: config.port,
          smtp_username: config.username,
          smtp_password: config.password,
          smtp_from: config.from,
          smtp_to: config.from
        })
        showSuccess('测试邮件发送成功')
      } catch (error) {
        console.error('测试SMTP配置失败', error)
        showError(error.response?.data?.error || '测试失败')
      }
    }

    const handleSendNotification = async () => {
      if (!notificationForm.value.title || !notificationForm.value.content) {
        showError('请填写通知标题和内容')
        return
      }
      try {
        await api.post('/notifications', notificationForm.value)
        showSuccess('发送成功')
        notificationForm.value = { type: 'system', target: 'all', title: '', content: '' }
        loadNotifications()
      } catch (error) {
        console.error('发送通知失败', error)
        showError(error.response?.data?.error || '发送失败')
      }
    }

    const deleteNotification = async (id) => {
      const confirmed = await showConfirm('确定要删除此通知吗？')
      if (!confirmed) return
      try {
        await api.delete(`/notifications/${id}`)
        showSuccess('删除成功')
        loadNotifications()
      } catch (error) {
        console.error('删除通知失败', error)
        showError(error.response?.data?.error || '删除失败')
      }
    }

    const approveDeletion = async (id) => {
      const confirmed = await showConfirm('确定要批准此删除申请吗？')
      if (!confirmed) return
      try {
        await api.post(`/admin/deletion-requests/${id}/approve`)
        showSuccess('已批准删除')
        loadDeletionRequests()
      } catch (error) {
        console.error('批准删除请求失败', error)
        showError(error.response?.data?.error || '操作失败')
      }
    }

    const rejectDeletion = async (id) => {
      const confirmed = await showConfirm('确定要拒绝此删除申请吗？')
      if (!confirmed) return
      try {
        await api.post(`/admin/deletion-requests/${id}/reject`)
        showSuccess('已拒绝')
        loadDeletionRequests()
      } catch (error) {
        console.error('拒绝删除请求失败', error)
        showError(error.response?.data?.error || '操作失败')
      }
    }

    watch(activeTab, (newTab) => {
      if (!isInitialized.value) return
      loadDataForCurrentTab()
    })

    onMounted(async () => {
      await checkAdmin()
      if (isAdmin.value) {
        loadStatistics()
        loadDeletionRequests()
      }
    })

    return {
      activeTab,
      isAdmin,
      isInitialized,
      currentUserId,
      loading,
      isRefreshing,
      statistics,
      users,
      articles,
      articlePage,
      articleTotalPages,
      articleFilter,
      articleStatusOptions,
      allComments,
      categories,
      titles,
      sidebarItems,
      announcementContent,
      siteConfigForm,
      smtpConfigForm,
      notifications,
      notificationForm,
      notificationTypes,
      deletionRequests,
      categoryForm,
      titleForm,
      grantForm,
      editRoleDialog,
      banDialog,
      editUserDialog,
      statusDialog,
      editCategoryDialog,
      showEditRoleDialog,
      showBanDialog,
      showEditUserDialog,
      showStatusDialog,
      showEditCategoryDialog,
      handleEditRole,
      handleBan,
      handleEditUser,
      handleUnban,
      handleDeleteUser,
      handleEditStatus,
      handleDeleteArticle,
      handleDeleteComment,
      addCategory,
      handleEditCategory,
      handleDeleteCategory,
      addTitle,
      grantTitle,
      revokeTitle,
      handleDeleteTitle,
      addSidebarItem,
      removeSidebarItem,
      saveSidebarConfig,
      saveAnnouncement,
      saveSiteConfig,
      saveSmtpConfig,
      testSmtpConfig,
      handleSendNotification,
      deleteNotification,
      approveDeletion,
      rejectDeletion,
      goToHome,
      handleRefresh,
      loadStatistics,
      loadUsers,
      loadArticles,
      loadComments,
      loadCategories,
      loadTitles,
      loadSidebarConfig,
      loadDeletionRequests,
      loadAnnouncement,
      loadSiteConfig,
      loadSmtpConfig,
      loadNotifications
    }
  }
}
</script>

<style scoped>
.admin-header {
  background: linear-gradient(135deg, #f8f9ff 0%, #fff 100%) !important;
  border-bottom: 1px solid rgba(103, 80, 164, 0.1);
}

.header-content {
  display: flex;
  align-items: center;
  width: 100%;
  padding: 0 24px;
  gap: 32px;
}

.header-brand {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-shrink: 0;
}

.brand-icon {
  background: linear-gradient(135deg, #6750A4 0%, #9c8cd8 100%);
  border-radius: 12px;
  padding: 8px;
  color: white !important;
}

.brand-title {
  font-size: 1.1rem;
  font-weight: 700;
  color: #1a1a2e;
  line-height: 1.3;
}

.brand-subtitle {
  font-size: 0.75rem;
  color: #6b7280;
}

.header-tabs {
  flex: 1;
  max-width: 800px;
}

.tab-item {
  min-width: auto !important;
  padding: 0 16px !important;
  font-size: 0.875rem;
  font-weight: 500;
  gap: 8px;
}

.tab-icon {
  margin-right: 4px;
}

.tab-badge {
  margin-left: 4px;
}

.header-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.admin-main {
  background: #f8f9ff;
  min-height: calc(100vh - 72px);
}

.dialog-card {
  border-radius: 20px !important;
  overflow: hidden;
}

.dialog-title {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 24px 24px 16px;
  font-size: 1.2rem;
  font-weight: 700;
  background: linear-gradient(135deg, #f8f9ff 0%, #fff 100%);
}

.title-icon {
  width: 40px;
  height: 40px;
  padding: 8px;
  border-radius: 10px;
  background: rgba(103, 80, 164, 0.1);
}

.dialog-body {
  padding: 24px !important;
}

.dialog-actions {
  padding: 16px 24px 24px;
  gap: 12px;
}

.user-info-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  background: #f8f9ff;
  border-radius: 12px;
}

.user-info-text {
  flex: 1;
}

.user-name {
  font-size: 1rem;
  font-weight: 600;
  color: #1a1a2e;
}

.user-meta {
  font-size: 0.85rem;
  color: #6b7280;
  margin-top: 4px;
}

.article-preview {
  padding: 16px;
  background: #f8f9ff;
  border-radius: 12px;
}

.preview-label {
  font-size: 0.75rem;
  color: #6b7280;
  margin-bottom: 8px;
}

.preview-title {
  font-size: 1rem;
  font-weight: 600;
  color: #1a1a2e;
}

.error-page {
  background: linear-gradient(135deg, #f8f9ff 0%, #fff 100%);
}

.error-content {
  text-align: center;
  padding: 48px 24px;
}

.error-icon {
  margin-bottom: 24px;
}

.error-title {
  font-size: 1.75rem;
  font-weight: 700;
  color: #1a1a2e;
  margin-bottom: 12px;
}

.error-message {
  font-size: 1rem;
  color: #6b7280;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.spin {
  animation: spin 1s linear infinite;
}

@media (max-width: 960px) {
  .header-content {
    padding: 0 16px;
    gap: 16px;
  }

  .brand-subtitle {
    display: none;
  }

  .header-tabs {
    max-width: none;
  }

  .tab-item span {
    display: none;
  }
}

@media (max-width: 640px) {
  .header-tabs {
    overflow-x: auto;
  }

  .tab-item {
    padding: 0 12px !important;
  }

  .dialog-card {
    margin: 16px !important;
    max-width: calc(100vw - 32px) !important;
    width: 100%;
  }

  .dialog-title {
    padding: 16px 20px 12px !important;
    font-size: 1.1rem !important;
  }

  .title-icon {
    width: 32px !important;
    height: 32px !important;
    padding: 6px !important;
  }

  .dialog-body {
    padding: 16px 20px !important;
  }

  .dialog-actions {
    padding: 12px 20px 16px !important;
    flex-wrap: wrap;
    gap: 8px;
  }

  .user-info-card {
    padding: 12px !important;
    gap: 12px !important;
  }
}
</style>
