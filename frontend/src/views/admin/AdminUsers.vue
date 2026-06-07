<template>
  <div class="admin-users">
    <!-- 页面标题 -->
    <div class="page-header mb-6">
      <div class="d-flex align-center justify-space-between flex-wrap ga-4">
        <div>
          <h1 class="text-h5 font-weight-bold mb-1">用户管理</h1>
          <p class="text-body-2 text-medium-emphasis">管理系统用户，查看在线状态，管理用户权限</p>
        </div>
        <div class="d-flex ga-3">
          <v-btn color="primary" @click="refreshData" :loading="loading" variant="tonal">
            <v-icon start>mdi-refresh</v-icon>
            刷新
          </v-btn>
        </div>
      </div>
    </div>

    <!-- 统计卡片 -->
    <v-row class="mb-6">
      <v-col cols="12" sm="6" md="3">
        <v-card elevation="0" class="stat-card stat-card-primary">
          <v-card-text>
            <div class="d-flex align-center justify-space-between">
              <div>
                <div class="text-h4 font-weight-bold">{{ stats.total }}</div>
                <div class="text-body-2 text-medium-emphasis">总用户数</div>
              </div>
              <v-avatar color="primary" size="48" rounded="lg">
                <v-icon>mdi-account-group</v-icon>
              </v-avatar>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" sm="6" md="3">
        <v-card elevation="0" class="stat-card stat-card-success">
          <v-card-text>
            <div class="d-flex align-center justify-space-between">
              <div>
                <div class="text-h4 font-weight-bold">{{ stats.online }}</div>
                <div class="text-body-2 text-medium-emphasis">在线用户</div>
              </div>
              <v-avatar color="success" size="48" rounded="lg">
                <v-icon>mdi-account-check</v-icon>
              </v-avatar>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" sm="6" md="3">
        <v-card elevation="0" class="stat-card stat-card-warning">
          <v-card-text>
            <div class="d-flex align-center justify-space-between">
              <div>
                <div class="text-h4 font-weight-bold">{{ stats.banned }}</div>
                <div class="text-body-2 text-medium-emphasis">被封禁用户</div>
              </div>
              <v-avatar color="warning" size="48" rounded="lg">
                <v-icon>mdi-account-off</v-icon>
              </v-avatar>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" sm="6" md="3">
        <v-card elevation="0" class="stat-card stat-card-info">
          <v-card-text>
            <div class="d-flex align-center justify-space-between">
              <div>
                <div class="text-h4 font-weight-bold">{{ stats.admins }}</div>
                <div class="text-body-2 text-medium-emphasis">管理员</div>
              </div>
              <v-avatar color="info" size="48" rounded="lg">
                <v-icon>mdi-shield-account</v-icon>
              </v-avatar>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- 在线用户快捷查看 -->
    <v-card elevation="0" class="mb-6 online-users-card">
      <v-card-title class="d-flex align-center">
        <v-icon class="mr-2" color="success">mdi-circle</v-icon>
        在线用户 ({{ onlineUsers.length }})
        <v-spacer />
        <v-btn
          size="small"
          variant="text"
          color="primary"
          @click="loadOnlineUsers"
          :loading="loadingOnline"
        >
          <v-icon start>mdi-refresh</v-icon>
          刷新
        </v-btn>
      </v-card-title>
      <v-card-text v-if="onlineUsers.length > 0">
        <v-chip-group>
          <v-chip
            v-for="user in onlineUsers.slice(0, 20)"
            :key="user.id"
            :to="'/profile?id=' + user.id"
            size="small"
            color="success"
            variant="tonal"
          >
            <v-avatar start size="20">
              <v-img v-if="user.avatar" :src="user.avatar" />
              <span v-else class="text-caption">
                {{ user.display_name?.[0] || user.username?.[0] }}
              </span>
            </v-avatar>
            {{ user.display_name || user.username }}
          </v-chip>
        </v-chip-group>
        <div v-if="onlineUsers.length > 20" class="mt-2 text-caption text-medium-emphasis">
          还有 {{ onlineUsers.length - 20 }} 位用户在线...
        </div>
      </v-card-text>
      <v-card-text v-else class="text-center text-medium-emphasis">
        <v-icon size="48" color="grey-lighten-2">mdi-account-off-outline</v-icon>
        <div class="mt-2">暂无在线用户</div>
      </v-card-text>
    </v-card>

    <!-- 筛选和搜索 -->
    <v-card elevation="0" class="mb-4">
      <v-card-text>
        <v-row align="center" dense>
          <v-col cols="12" md="4">
            <v-text-field
              v-model="filters.search"
              placeholder="搜索用户名、显示名称..."
              variant="outlined"
              density="compact"
              prepend-inner-icon="mdi-magnify"
              hide-details
              clearable
              @update:model-value="debounceSearch"
            />
          </v-col>
          <v-col cols="6" md="2">
            <v-select
              v-model="filters.role"
              :items="roleOptions"
              label="角色"
              variant="outlined"
              density="compact"
              hide-details
              clearable
            />
          </v-col>
          <v-col cols="6" md="2">
            <v-select
              v-model="filters.status"
              :items="statusOptions"
              label="状态"
              variant="outlined"
              density="compact"
              hide-details
              clearable
            />
          </v-col>
          <v-col cols="6" md="2">
            <v-select
              v-model="filters.online"
              :items="onlineOptions"
              label="在线状态"
              variant="outlined"
              density="compact"
              hide-details
              clearable
            />
          </v-col>
          <v-col cols="6" md="2">
            <v-btn
              color="primary"
              block
              @click="loadUsers"
              :loading="loading"
            >
              应用筛选
            </v-btn>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

    <!-- 用户列表 -->
    <v-card elevation="0">
      <v-card-title class="d-flex align-center justify-space-between">
        <span>用户列表 ({{ totalUsers }})</span>
        <v-pagination
          v-model="page"
          :length="totalPages"
          :total-visible="5"
          density="compact"
          @update:model-value="loadUsers"
        />
      </v-card-title>

      <v-data-table
        :headers="headers"
        :items="users"
        :loading="loading"
        :items-per-page="pageSize"
        class="elevation-0"
      >
        <template v-slot:item.user="{ item }">
          <div class="d-flex align-center py-2">
            <v-avatar size="40" class="mr-3">
              <v-img v-if="item.avatar" :src="item.avatar" />
              <span v-else class="text-primary font-weight-bold">
                {{ item.display_name?.[0] || item.username?.[0] }}
              </span>
            </v-avatar>
            <div>
              <div class="font-weight-medium">{{ item.display_name || item.username }}</div>
              <div class="text-caption text-medium-emphasis">@{{ item.username }}</div>
            </div>
          </div>
        </template>

        <template v-slot:item.role="{ item }">
          <v-chip
            :color="getRoleColor(item.role)"
            size="small"
            variant="tonal"
          >
            {{ getRoleText(item.role) }}
          </v-chip>
        </template>

        <template v-slot:item.status="{ item }">
          <v-chip
            :color="item.status === 'normal' ? 'success' : 'error'"
            size="small"
            variant="tonal"
          >
            {{ item.status === 'normal' ? '正常' : '已封禁' }}
          </v-chip>
        </template>

        <template v-slot:item.online_status="{ item }">
          <div class="d-flex align-center">
            <v-icon
              :color="item.online_status === 'online' ? 'success' : 'grey'"
              size="small"
              class="mr-1"
            >
              mdi-circle
            </v-icon>
            <span class="text-caption">
              {{ item.online_status === 'online' ? '在线' : '离线' }}
            </span>
            <span v-if="item.last_active_at" class="text-caption text-medium-emphasis ml-2">
              {{ formatLastActive(item.last_active_at) }}
            </span>
          </div>
        </template>

        <template v-slot:item.actions="{ item }">
          <v-menu>
            <template v-slot:activator="{ props }">
              <v-btn icon variant="text" size="small" v-bind="props">
                <v-icon>mdi-dots-vertical</v-icon>
              </v-btn>
            </template>
            <v-list density="compact">
              <v-list-item :to="'/profile?id=' + item.id" prepend-icon="mdi-account">
                <v-list-item-title>查看资料</v-list-item-title>
              </v-list-item>
              <v-list-item @click="editUser(item)" prepend-icon="mdi-pencil">
                <v-list-item-title>编辑用户</v-list-item-title>
              </v-list-item>
              <v-list-item
                v-if="item.status === 'normal'"
                @click="banUser(item)"
                prepend-icon="mdi-account-off"
              >
                <v-list-item-title>封禁用户</v-list-item-title>
              </v-list-item>
              <v-list-item
                v-else
                @click="unbanUser(item)"
                prepend-icon="mdi-account-check"
              >
                <v-list-item-title>解封用户</v-list-item-title>
              </v-list-item>
              <v-divider />
              <v-list-item
                @click="deleteUser(item)"
                prepend-icon="mdi-delete"
                class="text-error"
              >
                <v-list-item-title>删除用户</v-list-item-title>
              </v-list-item>
            </v-list>
          </v-menu>
        </template>
      </v-data-table>
    </v-card>

    <!-- 编辑用户对话框 -->
    <v-dialog v-model="editDialog" max-width="600">
      <v-card>
        <v-card-title class="bg-primary text-white pa-4">
          <v-icon class="mr-2">mdi-account-edit</v-icon>
          编辑用户
        </v-card-title>
        <v-card-text class="pa-4">
          <v-form ref="editForm">
            <v-text-field
              v-model="editForm.display_name"
              label="显示名称"
              variant="outlined"
              class="mb-3"
            />
            <v-textarea
              v-model="editForm.signature"
              label="签名"
              variant="outlined"
              rows="2"
              class="mb-3"
            />
            <v-select
              v-model="editForm.role"
              :items="roleOptions.slice(1)"
              label="角色"
              variant="outlined"
              class="mb-3"
            />
            <v-text-field
              :model-value="editForm.username"
              label="用户名"
              variant="outlined"
              disabled
            />
          </v-form>
        </v-card-text>
        <v-card-actions class="pa-4">
          <v-spacer />
          <v-btn variant="text" @click="editDialog = false">取消</v-btn>
          <v-btn color="primary" @click="saveUser" :loading="saving">保存</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- 封禁用户对话框 -->
    <v-dialog v-model="banDialog" max-width="500">
      <v-card>
        <v-card-title class="bg-warning text-white pa-4">
          <v-icon class="mr-2">mdi-account-off</v-icon>
          封禁用户
        </v-card-title>
        <v-card-text class="pa-4">
          <p class="mb-4">确定要封禁用户 <strong>{{ selectedUser?.display_name || selectedUser?.username }}</strong> 吗？</p>
          <v-textarea
            v-model="banReason"
            label="封禁原因"
            variant="outlined"
            rows="3"
            placeholder="请输入封禁原因..."
          />
        </v-card-text>
        <v-card-actions class="pa-4">
          <v-spacer />
          <v-btn variant="text" @click="banDialog = false">取消</v-btn>
          <v-btn color="warning" @click="confirmBan" :loading="actionLoading">确认封禁</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import api from '../../api'

export default {
  name: 'AdminUsers',
  setup() {
    const loading = ref(false)
    const loadingOnline = ref(false)
    const saving = ref(false)
    const actionLoading = ref(false)
    const users = ref([])
    const onlineUsers = ref([])
    const stats = ref({
      total: 0,
      online: 0,
      banned: 0,
      admins: 0
    })
    const page = ref(1)
    const pageSize = ref(20)
    const totalUsers = ref(0)
    const totalPages = computed(() => Math.ceil(totalUsers.value / pageSize.value))
    const editDialog = ref(false)
    const banDialog = ref(false)
    const selectedUser = ref(null)
    const editForm = ref({
      id: null,
      username: '',
      display_name: '',
      signature: '',
      role: ''
    })
    const banReason = ref('')

    const filters = ref({
      search: '',
      role: null,
      status: null,
      online: null
    })

    const headers = [
      { title: '用户', key: 'user', sortable: false },
      { title: '角色', key: 'role', width: '100px' },
      { title: '状态', key: 'status', width: '100px' },
      { title: '在线状态', key: 'online_status', width: '200px' },
      { title: '注册时间', key: 'created_at', width: '180px' },
      { title: '操作', key: 'actions', width: '80px', sortable: false }
    ]

    const roleOptions = [
      { title: '全部', value: null },
      { title: '管理员', value: 'admin' },
      { title: '普通用户', value: 'user' }
    ]

    const statusOptions = [
      { title: '全部', value: null },
      { title: '正常', value: 'normal' },
      { title: '已封禁', value: 'banned' }
    ]

    const onlineOptions = [
      { title: '全部', value: null },
      { title: '在线', value: 'online' },
      { title: '离线', value: 'offline' }
    ]

    let searchTimeout = null

    const debounceSearch = () => {
      if (searchTimeout) clearTimeout(searchTimeout)
      searchTimeout = setTimeout(() => {
        loadUsers()
      }, 500)
    }

    const loadUsers = async () => {
      loading.value = true
      try {
        const params = {
          page: page.value,
          page_size: pageSize.value
        }
        if (filters.value.search) params.search = filters.value.search
        if (filters.value.role) params.role = filters.value.role
        if (filters.value.status) params.status = filters.value.status
        if (filters.value.online) params.online = filters.value.online

        const response = await api.get('/admin/users', { params })
        users.value = response.data.users || []
        totalUsers.value = response.data.total || 0

        // 更新统计数据
        stats.value.total = response.data.total || 0
        stats.value.banned = response.data.banned_count || 0
        stats.value.admins = response.data.admin_count || 0
        stats.value.online = response.data.online_count || 0
      } catch (error) {
        console.error('加载用户失败:', error)
      } finally {
        loading.value = false
      }
    }

    const loadOnlineUsers = async () => {
      loadingOnline.value = true
      try {
        const response = await api.get('/users/status')
        onlineUsers.value = response.data.users || []
        stats.value.online = response.data.online_count || 0
      } catch (error) {
        console.error('加载在线用户失败:', error)
      } finally {
        loadingOnline.value = false
      }
    }

    const refreshData = () => {
      loadUsers()
      loadOnlineUsers()
    }

    const getRoleColor = (role) => {
      const colors = {
        admin: 'purple',
        user: 'blue'
      }
      return colors[role] || 'grey'
    }

    const getRoleText = (role) => {
      const texts = {
        admin: '管理员',
        user: '用户'
      }
      return texts[role] || '未知'
    }

    const formatLastActive = (timestamp) => {
      if (!timestamp) return ''
      const date = new Date(timestamp)
      const now = new Date()
      const diff = now - date
      const minutes = Math.floor(diff / 60000)
      
      if (minutes < 1) return '刚刚'
      if (minutes < 60) return `${minutes}分钟前`
      const hours = Math.floor(minutes / 60)
      if (hours < 24) return `${hours}小时前`
      return date.toLocaleDateString('zh-CN')
    }

    const editUser = (user) => {
      selectedUser.value = user
      editForm.value = {
        id: user.id,
        username: user.username,
        display_name: user.display_name || '',
        signature: user.signature || '',
        role: user.role
      }
      editDialog.value = true
    }

    const saveUser = async () => {
      saving.value = true
      try {
        await api.put(`/admin/users/${editForm.value.id}`, {
          display_name: editForm.value.display_name,
          signature: editForm.value.signature,
          role: editForm.value.role
        })
        editDialog.value = false
        loadUsers()
      } catch (error) {
        console.error('保存用户失败:', error)
      } finally {
        saving.value = false
      }
    }

    const banUser = (user) => {
      selectedUser.value = user
      banReason.value = ''
      banDialog.value = true
    }

    const confirmBan = async () => {
      actionLoading.value = true
      try {
        await api.post(`/admin/users/${selectedUser.value.id}/ban`, {
          reason: banReason.value
        })
        banDialog.value = false
        loadUsers()
      } catch (error) {
        console.error('封禁用户失败:', error)
      } finally {
        actionLoading.value = false
      }
    }

    const unbanUser = async (user) => {
      if (!confirm(`确定要解封用户 ${user.display_name || user.username} 吗？`)) return
      
      actionLoading.value = true
      try {
        await api.post(`/admin/users/${user.id}/unban`)
        loadUsers()
      } catch (error) {
        console.error('解封用户失败:', error)
      } finally {
        actionLoading.value = false
      }
    }

    const deleteUser = async (user) => {
      if (!confirm(`确定要删除用户 ${user.display_name || user.username} 吗？此操作不可恢复！`)) return
      
      actionLoading.value = true
      try {
        await api.delete(`/admin/users/${user.id}`)
        loadUsers()
      } catch (error) {
        console.error('删除用户失败:', error)
      } finally {
        actionLoading.value = false
      }
    }

    onMounted(() => {
      loadUsers()
      loadOnlineUsers()
    })

    return {
      loading,
      loadingOnline,
      saving,
      actionLoading,
      users,
      onlineUsers,
      stats,
      page,
      pageSize,
      totalUsers,
      totalPages,
      filters,
      headers,
      roleOptions,
      statusOptions,
      onlineOptions,
      editDialog,
      banDialog,
      selectedUser,
      editForm,
      banReason,
      debounceSearch,
      loadUsers,
      loadOnlineUsers,
      refreshData,
      getRoleColor,
      getRoleText,
      formatLastActive,
      editUser,
      saveUser,
      banUser,
      confirmBan,
      unbanUser,
      deleteUser
    }
  }
}
</script>

<style scoped>
.admin-users {
  max-width: 1600px;
  margin: 0 auto;
}

.stat-card {
  border-left: 4px solid;
  transition: transform 0.2s;
}

.stat-card:hover {
  transform: translateY(-4px);
}

.stat-card-primary {
  border-left-color: rgb(var(--v-theme-primary));
  background: linear-gradient(135deg, rgba(var(--v-theme-primary), 0.08), rgba(var(--v-theme-primary), 0.02));
}

.stat-card-success {
  border-left-color: rgb(var(--v-theme-success));
  background: linear-gradient(135deg, rgba(var(--v-theme-success), 0.08), rgba(var(--v-theme-success), 0.02));
}

.stat-card-warning {
  border-left-color: rgb(var(--v-theme-warning));
  background: linear-gradient(135deg, rgba(var(--v-theme-warning), 0.08), rgba(var(--v-theme-warning), 0.02));
}

.stat-card-info {
  border-left-color: rgb(var(--v-theme-info));
  background: linear-gradient(135deg, rgba(var(--v-theme-info), 0.08), rgba(var(--v-theme-info), 0.02));
}

.online-users-card {
  border: 2px dashed rgba(var(--v-theme-success), 0.3);
  background: rgba(var(--v-theme-success), 0.02);
}

.page-header h1 {
  color: rgb(var(--v-theme-on-surface));
}

@media (max-width: 600px) {
  .stat-card {
    margin-bottom: 12px;
  }
}
</style>
