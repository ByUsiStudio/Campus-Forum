<script setup>
import { ref, inject, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { adminApi } from '../api'

const router = useRouter()
const user = inject('user')

const activeTab = ref('users')
const users = ref([])
const articles = ref([])
const comments = ref([])
const reports = ref([])
const isLoading = ref(false)

const loadUsers = async () => {
  isLoading.value = true
  try {
    const response = await adminApi.getUsers()
    users.value = response.data.users
  } catch (error) {
    console.error('加载用户失败:', error)
  } finally {
    isLoading.value = false
  }
}

const loadArticles = async () => {
  isLoading.value = true
  try {
    const response = await adminApi.getArticles()
    articles.value = response.data.articles
  } catch (error) {
    console.error('加载文章失败:', error)
  } finally {
    isLoading.value = false
  }
}

const loadComments = async () => {
  isLoading.value = true
  try {
    const response = await adminApi.getComments()
    comments.value = response.data.comments
  } catch (error) {
    console.error('加载评论失败:', error)
  } finally {
    isLoading.value = false
  }
}

const loadReports = async () => {
  isLoading.value = true
  try {
    const response = await adminApi.getReports()
    reports.value = response.data.reports
  } catch (error) {
    console.error('加载举报失败:', error)
  } finally {
    isLoading.value = false
  }
}

const handleBanUser = async (userId) => {
  try {
    await adminApi.banUser(userId, { reason: '违规操作' })
    users.value = users.value.map(u => u.id === userId ? { ...u, is_banned: true } : u)
  } catch (error) {
    console.error('封禁用户失败:', error)
  }
}

const handleUnbanUser = async (userId) => {
  try {
    await adminApi.unbanUser(userId)
    users.value = users.value.map(u => u.id === userId ? { ...u, is_banned: false } : u)
  } catch (error) {
    console.error('解封用户失败:', error)
  }
}

const handleDeleteArticle = async (articleId) => {
  try {
    await adminApi.deleteArticle(articleId)
    articles.value = articles.value.filter(a => a.id !== articleId)
  } catch (error) {
    console.error('删除文章失败:', error)
  }
}

const handleDeleteComment = async (commentId) => {
  try {
    await adminApi.deleteComment(commentId)
    comments.value = comments.value.filter(c => c.id !== commentId)
  } catch (error) {
    console.error('删除评论失败:', error)
  }
}

const handleReport = async (reportId) => {
  try {
    await adminApi.handleReport(reportId)
    reports.value = reports.value.map(r => r.id === reportId ? { ...r, status: 'handled' } : r)
  } catch (error) {
    console.error('处理举报失败:', error)
  }
}

const formatTime = (timeStr) => {
  const date = new Date(timeStr)
  const now = new Date()
  const diff = now - date
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  
  if (days === 0) {
    const hours = Math.floor(diff / (1000 * 60 * 60))
    if (hours === 0) {
      const minutes = Math.floor(diff / (1000 * 60))
      return minutes <= 0 ? '刚刚' : `${minutes}分钟前`
    }
    return `${hours}小时前`
  } else if (days < 7) {
    return `${days}天前`
  } else {
    return date.toLocaleDateString('zh-CN')
  }
}

onMounted(() => {
  if (!user.value || user.value.role !== 'admin') {
    router.push('/')
    return
  }
})
</script>

<template>
  <v-container class="max-w-7xl mx-auto px-4 py-8">
    <!-- 返回按钮 -->
    <v-btn 
      text 
      color="gray-600" 
      class="mb-6 hover:text-primary transition-colors"
      @click="router.push('/')"
    >
      <v-icon class="mr-2" size="20">mdi-arrow-left</v-icon>
      返回首页
    </v-btn>
    
    <v-card rounded="2xl" elevation="4" class="overflow-hidden">
      <v-card-title class="gradient-purple text-white py-6 px-8">
        <v-icon class="mr-3" size="24">mdi-settings</v-icon>
        <span class="font-bold text-xl">管理后台</span>
      </v-card-title>
      
      <v-tabs 
        v-model="activeTab" 
        background-color="surface"
        class="border-b border-gray-100"
      >
        <v-tab 
          value="users" 
          @click="loadUsers"
          class="text-gray-600 hover:text-primary transition-colors"
        >
          <v-icon class="mr-2" size="18">mdi-users</v-icon>
          用户管理
        </v-tab>
        <v-tab 
          value="articles" 
          @click="loadArticles"
          class="text-gray-600 hover:text-primary transition-colors"
        >
          <v-icon class="mr-2" size="18">mdi-file</v-icon>
          文章管理
        </v-tab>
        <v-tab 
          value="comments" 
          @click="loadComments"
          class="text-gray-600 hover:text-primary transition-colors"
        >
          <v-icon class="mr-2" size="18">mdi-comment</v-icon>
          评论管理
        </v-tab>
        <v-tab 
          value="reports" 
          @click="loadReports"
          class="text-gray-600 hover:text-primary transition-colors"
        >
          <v-icon class="mr-2" size="18">mdi-flag</v-icon>
          举报管理
        </v-tab>
      </v-tabs>
      
      <v-tabs-items v-model="activeTab" class="p-6">
        <!-- 用户管理 -->
        <v-tab-item value="users">
          <div v-if="isLoading" class="loading-center">
            <v-progress-circular indeterminate color="primary" :size="48" />
          </div>
          <v-data-table
            v-else
            :headers="[
              { text: 'ID', value: 'id', width: '80px' }, 
              { text: '用户名', value: 'username' }, 
              { text: '昵称', value: 'display_name' }, 
              { text: 'QQ号', value: 'qq_number', width: '100px' }, 
              { text: '角色', value: 'role', width: '100px' }, 
              { text: '状态', value: 'is_banned', width: '100px' }, 
              { text: '注册时间', value: 'created_at' }, 
              { text: '操作', value: 'actions', width: '120px' }
            ]"
            :items="users"
            class="elevation-0 rounded-xl"
            hide-default-footer
          >
            <template v-slot:item.role="{ item }">
              <v-chip 
                size="small" 
                :color="item.role === 'admin' ? 'primary' : 'gray'"
                class="tag-purple"
              >
                {{ item.role === 'admin' ? '管理员' : '普通用户' }}
              </v-chip>
            </template>
            <template v-slot:item.is_banned="{ item }">
              <span :class="item.is_banned ? 'text-error font-medium' : 'text-success font-medium'">
                {{ item.is_banned ? '已封禁' : '正常' }}
              </span>
            </template>
            <template v-slot:item.created_at="{ item }">
              <span class="text-sm text-gray-500">{{ formatTime(item.created_at) }}</span>
            </template>
            <template v-slot:item.actions="{ item }">
              <v-btn 
                text 
                :color="item.is_banned ? 'success' : 'error'"
                size="small"
                @click="item.is_banned ? handleUnbanUser(item.id) : handleBanUser(item.id)"
              >
                {{ item.is_banned ? '解封' : '封禁' }}
              </v-btn>
            </template>
          </v-data-table>
          <div v-if="users.length === 0 && !isLoading" class="empty-state">
            <v-icon size="96" color="gray-200" class="empty-state-icon">mdi-users</v-icon>
            <p class="text-gray-400">暂无用户</p>
          </div>
        </v-tab-item>
        
        <!-- 文章管理 -->
        <v-tab-item value="articles">
          <div v-if="isLoading" class="loading-center">
            <v-progress-circular indeterminate color="primary" :size="48" />
          </div>
          <v-data-table
            v-else
            :headers="[
              { text: 'ID', value: 'id', width: '80px' }, 
              { text: '标题', value: 'title' }, 
              { text: '作者', value: 'user_name' }, 
              { text: '分类', value: 'category_name', width: '100px' }, 
              { text: '浏览量', value: 'view_count', width: '100px' }, 
              { text: '创建时间', value: 'created_at' }, 
              { text: '操作', value: 'actions', width: '100px' }
            ]"
            :items="articles.map(a => ({ ...a, user_name: a.user?.display_name || a.user?.username, category_name: a.category?.name }))"
            class="elevation-0 rounded-xl"
            hide-default-footer
          >
            <template v-slot:item.created_at="{ item }">
              <span class="text-sm text-gray-500">{{ formatTime(item.created_at) }}</span>
            </template>
            <template v-slot:item.actions="{ item }">
              <v-btn 
                text 
                color="error"
                size="small"
                @click="handleDeleteArticle(item.id)"
              >
                <v-icon class="mr-1" size="16">mdi-delete</v-icon>
                删除
              </v-btn>
            </template>
          </v-data-table>
          <div v-if="articles.length === 0 && !isLoading" class="empty-state">
            <v-icon size="96" color="gray-200" class="empty-state-icon">mdi-file</v-icon>
            <p class="text-gray-400">暂无文章</p>
          </div>
        </v-tab-item>
        
        <!-- 评论管理 -->
        <v-tab-item value="comments">
          <div v-if="isLoading" class="loading-center">
            <v-progress-circular indeterminate color="primary" :size="48" />
          </div>
          <v-data-table
            v-else
            :headers="[
              { text: 'ID', value: 'id', width: '80px' }, 
              { text: '内容', value: 'content' }, 
              { text: '作者', value: 'user_name' }, 
              { text: '文章', value: 'article_title' }, 
              { text: '创建时间', value: 'created_at' }, 
              { text: '操作', value: 'actions', width: '100px' }
            ]"
            :items="comments.map(c => ({ ...c, user_name: c.user?.display_name || c.user?.username, article_title: c.article?.title }))"
            class="elevation-0 rounded-xl"
            hide-default-footer
          >
            <template v-slot:item.content="{ item }">
              <span class="text-sm text-gray-600">{{ item.content.slice(0, 50) }}{{ item.content.length > 50 ? '...' : '' }}</span>
            </template>
            <template v-slot:item.created_at="{ item }">
              <span class="text-sm text-gray-500">{{ formatTime(item.created_at) }}</span>
            </template>
            <template v-slot:item.actions="{ item }">
              <v-btn 
                text 
                color="error"
                size="small"
                @click="handleDeleteComment(item.id)"
              >
                <v-icon class="mr-1" size="16">mdi-delete</v-icon>
                删除
              </v-btn>
            </template>
          </v-data-table>
          <div v-if="comments.length === 0 && !isLoading" class="empty-state">
            <v-icon size="96" color="gray-200" class="empty-state-icon">mdi-comment</v-icon>
            <p class="text-gray-400">暂无评论</p>
          </div>
        </v-tab-item>
        
        <!-- 举报管理 -->
        <v-tab-item value="reports">
          <div v-if="isLoading" class="loading-center">
            <v-progress-circular indeterminate color="primary" :size="48" />
          </div>
          <v-data-table
            v-else
            :headers="[
              { text: 'ID', value: 'id', width: '80px' }, 
              { text: '举报类型', value: 'type', width: '120px' }, 
              { text: '举报内容', value: 'content' }, 
              { text: '举报者', value: 'reporter_name' }, 
              { text: '状态', value: 'status', width: '100px' }, 
              { text: '创建时间', value: 'created_at' }, 
              { text: '操作', value: 'actions', width: '100px' }
            ]"
            :items="reports.map(r => ({ ...r, reporter_name: r.reporter?.display_name || r.reporter?.username }))"
            class="elevation-0 rounded-xl"
            hide-default-footer
          >
            <template v-slot:item.type="{ item }">
              <v-chip 
                size="small" 
                color="primary"
                class="tag-purple"
              >
                {{ item.type === 'article' ? '文章' : '评论' }}
              </v-chip>
            </template>
            <template v-slot:item.content="{ item }">
              <span class="text-sm text-gray-600">{{ item.content.slice(0, 50) }}{{ item.content.length > 50 ? '...' : '' }}</span>
            </template>
            <template v-slot:item.status="{ item }">
              <v-chip 
                size="small" 
                :color="item.status === 'pending' ? 'warning' : 'success'"
              >
                {{ item.status === 'pending' ? '待处理' : '已处理' }}
              </v-chip>
            </template>
            <template v-slot:item.created_at="{ item }">
              <span class="text-sm text-gray-500">{{ formatTime(item.created_at) }}</span>
            </template>
            <template v-slot:item.actions="{ item }">
              <v-btn 
                v-if="item.status === 'pending'"
                text 
                class="btn-gradient"
                size="small"
                @click="handleReport(item.id)"
              >
                <v-icon class="mr-1" size="16">mdi-check</v-icon>
                处理
              </v-btn>
              <span v-else class="text-gray-400 text-sm">已处理</span>
            </template>
          </v-data-table>
          <div v-if="reports.length === 0 && !isLoading" class="empty-state">
            <v-icon size="96" color="gray-200" class="empty-state-icon">mdi-flag</v-icon>
            <p class="text-gray-400">暂无举报</p>
          </div>
        </v-tab-item>
      </v-tabs-items>
    </v-card>
  </v-container>
</template>