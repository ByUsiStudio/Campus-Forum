<template>
  <div class="admin-users">
    <div class="page-header mb-6">
      <div class="d-flex align-center justify-between flex-wrap gap-4">
        <div>
          <h1 class="text-h5 font-weight-bold mb-1">用户管理</h1>
          <p class="text-body-2 text-medium-emphasis">管理系统用户，查看在线状态，管理用户权限</p>
        </div>
        <button class="layui-btn" @click="refreshData" :disabled="loading">
          <i class="fa-solid fa-rotate-right"></i>
          刷新
        </button>
      </div>
    </div>

    <div class="stats-grid">
      <div class="stat-card stat-card-primary">
        <div class="stat-content">
          <div class="stat-value">{{ stats.total }}</div>
          <div class="stat-label">总用户数</div>
        </div>
        <div class="stat-icon">
          <i class="fa-solid fa-users"></i>
        </div>
      </div>
      <div class="stat-card stat-card-success">
        <div class="stat-content">
          <div class="stat-value">{{ stats.online }}</div>
          <div class="stat-label">在线用户</div>
        </div>
        <div class="stat-icon">
          <i class="fa-solid fa-user-check"></i>
        </div>
      </div>
      <div class="stat-card stat-card-warning">
        <div class="stat-content">
          <div class="stat-value">{{ stats.banned }}</div>
          <div class="stat-label">被封禁用户</div>
        </div>
        <div class="stat-icon">
          <i class="fa-solid fa-user-xmark"></i>
        </div>
      </div>
      <div class="stat-card stat-card-info">
        <div class="stat-content">
          <div class="stat-value">{{ stats.admins }}</div>
          <div class="stat-label">管理员</div>
        </div>
        <div class="stat-icon">
          <i class="fa-solid fa-shield-halved"></i>
        </div>
      </div>
    </div>

    <div class="layui-card mb-6 online-users-card">
      <div class="layui-card-header d-flex align-center justify-between">
        <div class="d-flex align-center gap-2">
          <i class="fa-solid fa-circle" style="color: #52C41A;"></i>
          <span>在线用户 ({{ onlineUsers.length }})</span>
        </div>
        <button class="layui-btn layui-btn-sm" @click="loadOnlineUsers" :disabled="loadingOnline">
          <i class="fa-solid fa-rotate-right"></i>
          刷新
        </button>
      </div>
      <div class="layui-card-body">
        <div v-if="onlineUsers.length > 0" class="online-users-wrap">
          <router-link
            v-for="user in onlineUsers.slice(0, 20)"
            :key="user.id"
            :to="'/profile?id=' + user.id"
            class="online-user-chip"
          >
            <div class="chip-avatar">
              <img v-if="user.avatar" :src="user.avatar" />
              <span v-else>{{ user.display_name?.[0] || user.username?.[0] }}</span>
            </div>
            {{ user.display_name || user.username }}
          </router-link>
          <div v-if="onlineUsers.length > 20" class="text-caption text-medium-emphasis mt-2">
            还有 {{ onlineUsers.length - 20 }} 位用户在线...
          </div>
        </div>
        <div v-else class="empty-online">
          <i class="fa-solid fa-user-xmark" style="font-size: 48px; color: #ccc;"></i>
          <div class="mt-2">暂无在线用户</div>
        </div>
      </div>
    </div>

    <div class="layui-card mb-4">
      <div class="layui-card-body">
        <div class="filter-form">
          <div class="filter-item">
            <div class="layui-input-group">
              <div class="layui-input-group-prepend">
                <span class="layui-input-group-text"><i class="fa-solid fa-magnifying-glass"></i></span>
              </div>
              <input type="text" v-model="filters.search" placeholder="搜索用户名、显示名称..." class="layui-input" @input="debounceSearch" />
            </div>
          </div>
          <div class="filter-item">
            <select v-model="filters.role" class="layui-select">
              <option :value="null">全部角色</option>
              <option value="admin">管理员</option>
              <option value="user">普通用户</option>
            </select>
          </div>
          <div class="filter-item">
            <select v-model="filters.status" class="layui-select">
              <option :value="null">全部状态</option>
              <option value="normal">正常</option>
              <option value="banned">已封禁</option>
            </select>
          </div>
          <div class="filter-item">
            <select v-model="filters.online" class="layui-select">
              <option :value="null">全部在线状态</option>
              <option value="online">在线</option>
              <option value="offline">离线</option>
            </select>
          </div>
          <div class="filter-item">
            <button class="layui-btn" @click="loadUsers" :disabled="loading">应用筛选</button>
          </div>
        </div>
      </div>
    </div>

    <div class="layui-card">
      <div class="layui-card-header d-flex align-center justify-between">
        <span>用户列表 ({{ totalUsers }})</span>
        <div class="pagination" v-if="totalPages > 1">
          <button class="layui-laypage-prev" :disabled="page <= 1" @click="page--; loadUsers()">上一页</button>
          <span v-for="p in pageRange" :key="p" :class="['layui-laypage-curr', { active: p === page }]" @click="page = p; loadUsers()">{{ p }}</span>
          <button class="layui-laypage-next" :disabled="page >= totalPages" @click="page++; loadUsers()">下一页</button>
        </div>
      </div>

      <div class="users-table">
        <table class="layui-table">
          <thead>
            <tr>
              <th>用户</th>
              <th>角色</th>
              <th>状态</th>
              <th>在线状态</th>
              <th>注册时间</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in users" :key="item.id">
              <td>
                <div class="user-cell">
                  <div class="user-avatar">
                    <img v-if="item.avatar" :src="item.avatar" />
                    <span v-else>{{ item.display_name?.[0] || item.username?.[0] }}</span>
                  </div>
                  <div>
                    <div class="font-weight-medium">{{ item.display_name || item.username }}</div>
                    <div class="text-caption text-medium-emphasis">@{{ item.username }}</div>
                  </div>
                </div>
              </td>
              <td>
                <span :class="['role-chip', getRoleClass(item.role)]">{{ getRoleText(item.role) }}</span>
              </td>
              <td>
                <span :class="['status-chip', item.status === 'normal' ? 'success' : 'error']">{{ item.status === 'normal' ? '正常' : '已封禁' }}</span>
              </td>
              <td>
                <div class="online-status">
                  <i :class="['fa-solid', item.online_status === 'online' ? 'fa-circle' : 'fa-circle']" :style="{ color: item.online_status === 'online' ? '#52C41A' : '#999' }"></i>
                  <span>{{ item.online_status === 'online' ? '在线' : '离线' }}</span>
                  <span v-if="item.last_active_at" class="text-medium-emphasis">{{ formatLastActive(item.last_active_at) }}</span>
                </div>
              </td>
              <td>{{ formatDate(item.created_at) }}</td>
              <td>
                <div class="actions-cell">
                  <button class="action-btn" @click="editUser(item)" title="查看资料">
                    <i class="fa-solid fa-user"></i>
                  </button>
                  <button class="action-btn" @click="editUser(item)" title="编辑用户">
                    <i class="fa-solid fa-pencil"></i>
                  </button>
                  <button v-if="item.status === 'normal'" class="action-btn" @click="banUser(item)" title="封禁用户">
                    <i class="fa-solid fa-user-xmark"></i>
                  </button>
                  <button v-else class="action-btn" @click="unbanUser(item)" title="解封用户">
                    <i class="fa-solid fa-user-check"></i>
                  </button>
                  <button class="action-btn danger" @click="deleteUser(item)" title="删除用户">
                    <i class="fa-solid fa-trash"></i>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div v-if="editDialog" class="modal-overlay" @click.self="editDialog = false">
      <div class="modal-content">
        <div class="layui-card">
          <div class="layui-card-header bg-primary text-white">
            <i class="fa-solid fa-user-pen"></i>
            编辑用户
          </div>
          <div class="layui-card-body">
            <div class="layui-form-item">
              <label class="layui-form-label">显示名称</label>
              <div class="layui-input-block">
                <input type="text" v-model="editForm.display_name" class="layui-input" />
              </div>
            </div>
            <div class="layui-form-item">
              <label class="layui-form-label">签名</label>
              <div class="layui-input-block">
                <textarea v-model="editForm.signature" class="layui-textarea" rows="2"></textarea>
              </div>
            </div>
            <div class="layui-form-item">
              <label class="layui-form-label">角色</label>
              <div class="layui-input-block">
                <select v-model="editForm.role" class="layui-select">
                  <option value="user">普通用户</option>
                  <option value="admin">管理员</option>
                </select>
              </div>
            </div>
            <div class="layui-form-item">
              <label class="layui-form-label">用户名</label>
              <div class="layui-input-block">
                <input type="text" :value="editForm.username" class="layui-input" disabled />
              </div>
            </div>
          </div>
          <div class="layui-card-footer d-flex justify-end gap-3">
            <button class="layui-btn layui-btn-primary" @click="editDialog = false">取消</button>
            <button class="layui-btn" @click="saveUser" :disabled="saving">保存</button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="banDialog" class="modal-overlay" @click.self="banDialog = false">
      <div class="modal-content">
        <div class="layui-card">
          <div class="layui-card-header bg-warning text-white">
            <i class="fa-solid fa-user-xmark"></i>
            封禁用户
          </div>
          <div class="layui-card-body">
            <p class="mb-4">确定要封禁用户 <strong>{{ selectedUser?.display_name || selectedUser?.username }}</strong> 吗？</p>
            <div class="layui-form-item">
              <label class="layui-form-label">封禁原因</label>
              <div class="layui-input-block">
                <textarea v-model="banReason" class="layui-textarea" rows="3" placeholder="请输入封禁原因..."></textarea>
              </div>
            </div>
          </div>
          <div class="layui-card-footer d-flex justify-end gap-3">
            <button class="layui-btn layui-btn-primary" @click="banDialog = false">取消</button>
            <button class="layui-btn layui-btn-warning" @click="confirmBan" :disabled="actionLoading">确认封禁</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import api from '../../api'
import { confirm, success, error } from '../../utils/modal'

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

    const pageRange = computed(() => {
      const range = []
      const total = totalPages.value
      const current = page.value
      const visible = 5
      
      let start = Math.max(1, current - Math.floor(visible / 2))
      let end = Math.min(total, start + visible - 1)
      
      if (end - start + 1 < visible) {
        start = Math.max(1, end - visible + 1)
      }
      
      for (let i = start; i <= end; i++) {
        range.push(i)
      }
      return range
    })

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

    const getRoleClass = (role) => {
      const classes = {
        admin: 'admin',
        user: 'user'
      }
      return classes[role] || 'user'
    }

    const getRoleText = (role) => {
      const texts = {
        admin: '管理员',
        user: '用户'
      }
      return texts[role] || '未知'
    }

    const formatDate = (dateString) => {
      if (!dateString) return '-'
      const date = new Date(dateString)
      return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
      })
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
        success('保存成功')
        loadUsers()
      } catch (error) {
        console.error('保存用户失败:', error)
        error('保存失败')
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
        success('封禁成功')
        loadUsers()
      } catch (error) {
        console.error('封禁用户失败:', error)
        error('封禁失败')
      } finally {
        actionLoading.value = false
      }
    }

    const unbanUser = async (user) => {
      const confirmed = await confirm(`确定要解封用户 ${user.display_name || user.username} 吗？`)
      if (!confirmed) return
      
      actionLoading.value = true
      try {
        await api.post(`/admin/users/${user.id}/unban`)
        success('解封成功')
        loadUsers()
      } catch (error) {
        console.error('解封用户失败:', error)
        error('解封失败')
      } finally {
        actionLoading.value = false
      }
    }

    const deleteUser = async (user) => {
      const confirmed = await confirm(`确定要删除用户 ${user.display_name || user.username} 吗？此操作不可恢复！`)
      if (!confirmed) return
      
      actionLoading.value = true
      try {
        await api.delete(`/admin/users/${user.id}`)
        success('删除成功')
        loadUsers()
      } catch (error) {
        console.error('删除用户失败:', error)
        error('删除失败')
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
      pageRange,
      filters,
      editDialog,
      banDialog,
      selectedUser,
      editForm,
      banReason,
      debounceSearch,
      loadUsers,
      loadOnlineUsers,
      refreshData,
      getRoleClass,
      getRoleText,
      formatDate,
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

.page-header {
  margin-bottom: 24px;
}

.text-h5 {
  font-size: 24px;
}

.font-weight-bold {
  font-weight: 700;
}

.text-body-2 {
  font-size: 14px;
}

.text-medium-emphasis {
  color: #999;
}

.d-flex {
  display: flex;
}

.align-center {
  align-items: center;
}

.justify-between {
  justify-content: space-between;
}

.flex-wrap {
  flex-wrap: wrap;
}

.gap-4 {
  gap: 16px;
}

.mb-1 {
  margin-bottom: 4px;
}

.mb-4 {
  margin-bottom: 16px;
}

.mb-6 {
  margin-bottom: 24px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  padding: 20px;
  border-radius: 8px;
  border-left: 4px solid;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stat-card-primary {
  border-left-color: #1E9FFF;
  background: linear-gradient(135deg, rgba(30, 159, 255, 0.08), rgba(30, 159, 255, 0.02));
}

.stat-card-success {
  border-left-color: #52C41A;
  background: linear-gradient(135deg, rgba(82, 196, 26, 0.08), rgba(82, 196, 26, 0.02));
}

.stat-card-warning {
  border-left-color: #FAAD14;
  background: linear-gradient(135deg, rgba(250, 173, 20, 0.08), rgba(250, 173, 20, 0.02));
}

.stat-card-info {
  border-left-color: #9B59B6;
  background: linear-gradient(135deg, rgba(155, 89, 182, 0.08), rgba(155, 89, 182, 0.02));
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
}

.stat-label {
  font-size: 14px;
  color: #666;
}

.stat-icon {
  font-size: 32px;
  opacity: 0.5;
}

.online-users-card {
  border: 2px dashed rgba(82, 196, 26, 0.3);
  background: rgba(82, 196, 26, 0.02);
}

.online-users-wrap {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.online-user-chip {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: #f6ffed;
  color: #52c41a;
  border-radius: 20px;
  font-size: 13px;
  text-decoration: none;
}

.chip-avatar {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: #52c41a;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 10px;
  overflow: hidden;
}

.chip-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.empty-online {
  text-align: center;
  padding: 20px;
}

.text-caption {
  font-size: 12px;
}

.filter-form {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.filter-item {
  flex: 1;
  min-width: 150px;
}

.filter-item:last-child {
  flex: 0;
}

.users-table {
  overflow-x: auto;
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #1E9FFF;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  overflow: hidden;
}

.user-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.role-chip {
  padding: 4px 12px;
  border-radius: 4px;
  font-size: 12px;
}

.role-chip.admin {
  background: #f5e6ff;
  color: #9B59B6;
}

.role-chip.user {
  background: #e6f7ff;
  color: #1E9FFF;
}

.status-chip {
  padding: 4px 12px;
  border-radius: 4px;
  font-size: 12px;
}

.status-chip.success {
  background: #f6ffed;
  color: #52c41a;
}

.status-chip.error {
  background: #fff2f0;
  color: #ff4d4f;
}

.online-status {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
}

.actions-cell {
  display: flex;
  gap: 6px;
}

.action-btn {
  width: 32px;
  height: 32px;
  border: none;
  background: transparent;
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #999;
  transition: all 0.3s;
}

.action-btn:hover {
  background: #f0f0f0;
  color: #1E9FFF;
}

.action-btn.danger:hover {
  background: #fff2f0;
  color: #FF5722;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  width: 90%;
  max-width: 600px;
  background: white;
  border-radius: 12px;
  overflow: hidden;
}

.bg-primary {
  background: #1E9FFF;
}

.bg-warning {
  background: #FAAD14;
}

.text-white {
  color: white;
}

.justify-end {
  justify-content: flex-end;
}

.gap-3 {
  gap: 12px;
}
</style>