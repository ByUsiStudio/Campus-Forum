<template>
  <div class="page-container">
    <div class="admin-users">
      <!-- 页面标题 -->
      <div class="page-header mb-6">
        <div class="page-header-row">
          <div>
            <h1 class="page-title">用户管理</h1>
            <p class="page-subtitle">管理系统用户，查看在线状态，管理用户权限</p>
          </div>
          <div>
            <el-button type="primary" :loading="loading" @click="refreshData">
              <el-icon class="mr-2"><Refresh /></el-icon>
              刷新
            </el-button>
          </div>
        </div>
      </div>

      <!-- 统计卡片 -->
      <el-row :gutter="16" class="mb-6">
        <el-col :xs="24" :sm="12" :md="6">
          <el-card shadow="never" class="stat-card stat-card-primary">
            <div class="stat-row">
              <div>
                <div class="stat-value">{{ stats.total }}</div>
                <div class="stat-label">总用户数</div>
              </div>
              <el-avatar :size="48" class="stat-avatar avatar-primary">
                <el-icon><User /></el-icon>
              </el-avatar>
            </div>
          </el-card>
        </el-col>
        <el-col :xs="24" :sm="12" :md="6">
          <el-card shadow="never" class="stat-card stat-card-success">
            <div class="stat-row">
              <div>
                <div class="stat-value">{{ stats.online }}</div>
                <div class="stat-label">在线用户</div>
              </div>
              <el-avatar :size="48" class="stat-avatar avatar-success">
                <el-icon><User /></el-icon>
              </el-avatar>
            </div>
          </el-card>
        </el-col>
        <el-col :xs="24" :sm="12" :md="6">
          <el-card shadow="never" class="stat-card stat-card-warning">
            <div class="stat-row">
              <div>
                <div class="stat-value">{{ stats.banned }}</div>
                <div class="stat-label">被封禁用户</div>
              </div>
              <el-avatar :size="48" class="stat-avatar avatar-warning">
                <el-icon><UserFilled /></el-icon>
              </el-avatar>
            </div>
          </el-card>
        </el-col>
        <el-col :xs="24" :sm="12" :md="6">
          <el-card shadow="never" class="stat-card stat-card-info">
            <div class="stat-row">
              <div>
                <div class="stat-value">{{ stats.admins }}</div>
                <div class="stat-label">管理员</div>
              </div>
              <el-avatar :size="48" class="stat-avatar avatar-info">
                <el-icon><Key /></el-icon>
              </el-avatar>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 在线用户快捷查看 -->
      <el-card shadow="never" class="mb-6 online-users-card">
        <template #header>
          <div class="online-users-header">
            <div class="online-users-title">
              <el-icon color="var(--el-color-success)" class="mr-2"><CircleCheckFilled /></el-icon>
              在线用户 ({{ onlineUsers.length }})
            </div>
            <div class="flex-1"></div>
            <el-button size="small" text type="primary" :loading="loadingOnline" @click="loadOnlineUsers">
              <el-icon class="mr-2"><Refresh /></el-icon>
              刷新
            </el-button>
          </div>
        </template>
        <div v-if="onlineUsers.length > 0">
          <el-tag
            v-for="user in onlineUsers.slice(0, 20)"
            :key="user.id"
            type="success"
            effect="plain"
            round
            class="online-chip"
            @click="goProfile(user.id)"
          >
            <el-avatar :size="20" class="chip-avatar">
              <img v-if="user.avatar" :src="user.avatar" class="chip-avatar-img" />
              <template v-else>
                {{ user.display_name?.[0] || user.username?.[0] }}
              </template>
            </el-avatar>
            <span class="chip-name">{{ user.display_name || user.username }}</span>
          </el-tag>
          <div v-if="onlineUsers.length > 20" class="mt-2 online-more">
            还有 {{ onlineUsers.length - 20 }} 位用户在线...
          </div>
        </div>
        <div v-else class="online-empty">
          <el-icon :size="48" color="var(--el-color-info-light-5)"><User /></el-icon>
          <div class="mt-2">暂无在线用户</div>
        </div>
      </el-card>

      <!-- 筛选和搜索 -->
      <el-card shadow="never" class="mb-4">
        <el-row :gutter="12" align="middle">
          <el-col :xs="24" :md="8">
            <el-input
              v-model="filters.search"
              placeholder="搜索用户名、显示名称..."
              clearable
              :prefix-icon="Search"
              @update:model-value="debounceSearch"
            />
          </el-col>
          <el-col :xs="12" :md="4">
            <el-select v-model="filters.role" placeholder="角色" clearable class="w-full">
              <el-option
                v-for="opt in roleOptions"
                :key="opt.value"
                :label="opt.title"
                :value="opt.value"
              />
            </el-select>
          </el-col>
          <el-col :xs="12" :md="4">
            <el-select v-model="filters.status" placeholder="状态" clearable class="w-full">
              <el-option
                v-for="opt in statusOptions"
                :key="opt.value"
                :label="opt.title"
                :value="opt.value"
              />
            </el-select>
          </el-col>
          <el-col :xs="12" :md="4">
            <el-select v-model="filters.online" placeholder="在线状态" clearable class="w-full">
              <el-option
                v-for="opt in onlineOptions"
                :key="opt.value"
                :label="opt.title"
                :value="opt.value"
              />
            </el-select>
          </el-col>
          <el-col :xs="12" :md="4">
            <el-button type="primary" :loading="loading" class="w-full" @click="loadUsers">
              应用筛选
            </el-button>
          </el-col>
        </el-row>
      </el-card>

      <!-- 用户列表 -->
      <el-card shadow="never">
        <template #header>
          <div class="table-header">
            <span>用户列表 ({{ totalUsers }})</span>
            <el-pagination
              v-model:current-page="page"
              :page-size="pageSize"
              :total="totalUsers"
              :pager-count="5"
              layout="prev, pager, next"
              @current-change="loadUsers"
            />
          </div>
        </template>

        <div class="table-responsive">
        <el-table :data="users" v-loading="loading">
          <el-table-column label="用户" min-width="220">
            <template #default="{ row }">
              <div class="user-cell">
                <el-avatar :size="40" class="mr-3 user-avatar">
                  <img v-if="row.avatar" :src="row.avatar" class="user-avatar-img" />
                  <template v-else>
                    <span class="user-avatar-char">{{ row.display_name?.[0] || row.username?.[0] }}</span>
                  </template>
                </el-avatar>
                <div>
                  <div class="user-name">{{ row.display_name || row.username }}</div>
                  <div class="user-handle">@{{ row.username }}</div>
                </div>
              </div>
            </template>
          </el-table-column>

          <el-table-column label="角色" width="100">
            <template #default="{ row }">
              <el-tag
                :color="getRoleColor(row.role)"
                effect="light"
                size="small"
              >
                {{ getRoleText(row.role) }}
              </el-tag>
            </template>
          </el-table-column>

          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag
                :type="row.status === 'normal' ? 'success' : 'danger'"
                effect="light"
                size="small"
              >
                {{ row.status === 'normal' ? '正常' : '已封禁' }}
              </el-tag>
            </template>
          </el-table-column>

          <el-table-column label="在线状态" width="200">
            <template #default="{ row }">
              <div class="online-cell">
                <el-icon
                  :color="row.online_status === 'online' ? 'var(--el-color-success)' : 'var(--el-color-info)'"
                  size="12"
                  class="mr-1"
                >
                  <CircleCheckFilled />
                </el-icon>
                <span class="online-text">
                  {{ row.online_status === 'online' ? '在线' : '离线' }}
                </span>
                <span v-if="row.last_active_at" class="online-last ml-2">
                  {{ formatLastActive(row.last_active_at) }}
                </span>
              </div>
            </template>
          </el-table-column>

          <el-table-column prop="created_at" label="注册时间" width="180" />

          <el-table-column label="操作" width="90">
            <template #default="{ row }">
              <el-dropdown trigger="click" @command="(cmd) => handleAction(cmd, row)">
                <el-button text size="small" :icon="MoreFilled" />
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="view" :icon="User">
                      查看资料
                    </el-dropdown-item>
                    <el-dropdown-item command="edit" :icon="Edit">
                      编辑用户
                    </el-dropdown-item>
                    <el-dropdown-item
                      v-if="row.status === 'normal'"
                      command="ban"
                      :icon="CircleCloseFilled"
                    >
                      封禁用户
                    </el-dropdown-item>
                    <el-dropdown-item
                      v-else
                      command="unban"
                      :icon="CircleCheckFilled"
                    >
                      解封用户
                    </el-dropdown-item>
                    <el-dropdown-item command="delete" :icon="Delete" divided class="danger-item">
                      删除用户
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </template>
          </el-table-column>
        </el-table>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '../../api'
import { adminUserApi } from '../../api/admin'
import { confirm, success, error } from '../../utils/message'
import { jcCloseAll, jcOpenHtml } from '@/utils/jcu'
import {
  Refresh,
  Search,
  MoreFilled,
  User,
  UserFilled,
  Key,
  CircleCheckFilled,
  CircleCloseFilled,
  Edit,
  Delete
} from '@element-plus/icons-vue'

const router = useRouter()

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
const selectedUser = ref(null)
const editForm = ref({
  id: null,
  username: '',
  display_name: '',
  signature: '',
  role: ''
})

const filters = ref({
  search: '',
  role: null,
  status: null,
  online: null
})

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

    // adminUserApi.getUsers() does not accept params, so use the default api for
    // the paginated/filtered listing (endpoint must stay without the /api prefix).
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
    error('加载用户失败')
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
    error('加载在线用户失败')
  } finally {
    loadingOnline.value = false
  }
}

const refreshData = () => {
  loadUsers()
  loadOnlineUsers()
}

const goProfile = (id) => {
  router.push('/profile?id=' + id)
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
  openEditDialog()
}

const openEditDialog = () => {
  const roleOptions = [
    { title: '管理员', value: 'admin' },
    { title: '普通用户', value: 'user' }
  ]
  const roleHint =
    `<select data-user-role class="jc-modal__input" style="width:100%;margin-top:4px;">` +
    roleOptions.map((r) => `<option value="${r.value}" ${r.value === editForm.value.role ? 'selected' : ''}>${r.title}</option>`).join('') +
    `</select>`

  const esc = (s) => String(s ?? '').replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;')
  const field = (label, inner) =>
    `<div style="margin-bottom:12px;"><div style="margin:0 0 4px;font-weight:600;font-size:13px;color:var(--jc-text,#333);">${label}</div>${inner}</div>`

  const content =
    field('显示名称', `<input data-user-input="display_name" class="jc-modal__input" type="text" value="${esc(editForm.value.display_name)}" placeholder="显示名称" />`) +
    field('签名', `<textarea data-user-input="signature" class="jc-modal__input" rows="2" style="resize:vertical;">${esc(editForm.value.signature)}</textarea>`) +
    field('角色', roleHint) +
    field('用户名', `<div style="padding:10px 12px;background:var(--jc-bg,#f3f5fb);border-radius:8px;font-size:13px;color:var(--jc-text-2,#64748b);">${esc(editForm.value.username || '')}</div>`)

  jcOpenHtml({
    title: '编辑用户',
    content,
    width: 600,
    size: 'md',
    buttons: [
      { text: '取消', type: 'default', action: () => jcCloseAll() },
      {
        text: '保存',
        type: 'primary',
        action: (inst) => {
          const root = inst.modalContent
          const display_name = (root.querySelector('[data-user-input="display_name"]')?.value || '').trim()
          if (!display_name) {
            error('请输入显示名称')
            return
          }
          saveUser({
            display_name,
            signature: root.querySelector('[data-user-input="signature"]')?.value || '',
            role: root.querySelector('[data-user-role]')?.value || ''
          })
        }
      }
    ]
  })
}

const saveUser = async (data) => {
  saving.value = true
  try {
    await adminUserApi.updateUser(editForm.value.id, data)
    jcCloseAll()
    success('保存成功')
    loadUsers()
  } catch (error) {
    console.error('保存用户失败:', error)
    error(error.response?.data?.error || '保存用户失败')
  } finally {
    saving.value = false
  }
}

const banUser = (user) => {
  selectedUser.value = user
  openBanDialog()
}

const openBanDialog = () => {
  const esc = (s) => String(s ?? '').replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;')
  const name = selectedUser.value?.display_name || selectedUser.value?.username || ''
  const content = `
    <div style="margin-bottom:12px;font-size:14px;color:var(--jc-text,#0f172a);">
      确定要封禁用户 <strong>${esc(name)}</strong> 吗？
    </div>
    <div><div style="margin:0 0 4px;font-weight:600;font-size:13px;color:var(--jc-text,#333);">封禁原因</div>
      <textarea data-ban-reason class="jc-modal__input" rows="3" style="resize:vertical;" placeholder="请输入封禁原因..."></textarea>
    </div>
  `
  jcOpenHtml({
    title: '封禁用户',
    content,
    width: 500,
    size: 'md',
    buttons: [
      { text: '取消', type: 'default', action: () => jcCloseAll() },
      {
        text: '确认封禁',
        type: 'danger',
        action: (inst) => {
          const reason = inst.modalContent.querySelector('[data-ban-reason]')?.value || ''
          confirmBan(reason)
        }
      }
    ]
  })
}

const confirmBan = async (reason) => {
  actionLoading.value = true
  try {
    await adminUserApi.banUser(selectedUser.value.id, reason)
    jcCloseAll()
    success('封禁成功')
    loadUsers()
  } catch (error) {
    console.error('封禁用户失败:', error)
    error(error.response?.data?.error || '封禁用户失败')
  } finally {
    actionLoading.value = false
  }
}

const unbanUser = async (user) => {
  let ok
  try {
    ok = await confirm(`确定要解封用户 ${user.display_name || user.username} 吗？`)
  } catch (e) {
    return
  }
  if (!ok) return

  actionLoading.value = true
  try {
    await adminUserApi.unbanUser(user.id)
    success('解封成功')
    loadUsers()
  } catch (error) {
    console.error('解封用户失败:', error)
    error(error.response?.data?.error || '解封用户失败')
  } finally {
    actionLoading.value = false
  }
}

const deleteUser = async (user) => {
  let ok
  try {
    ok = await confirm(`确定要删除用户 ${user.display_name || user.username} 吗？此操作不可恢复！`)
  } catch (e) {
    return
  }
  if (!ok) return

  actionLoading.value = true
  try {
    await adminUserApi.deleteUser(user.id)
    success('删除成功')
    loadUsers()
  } catch (error) {
    console.error('删除用户失败:', error)
    error(error.response?.data?.error || '删除用户失败')
  } finally {
    actionLoading.value = false
  }
}

const handleAction = (command, row) => {
  switch (command) {
    case 'view':
      goProfile(row.id)
      break
    case 'edit':
      editUser(row)
      break
    case 'ban':
      banUser(row)
      break
    case 'unban':
      unbanUser(row)
      break
    case 'delete':
      deleteUser(row)
      break
  }
}

onMounted(() => {
  loadUsers()
  loadOnlineUsers()
})
</script>

<style scoped>
.admin-users {
  max-width: 1600px;
  margin: 0 auto;
}

.page-header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 16px;
}

.page-title {
  margin: 0 0 4px;
  font-size: 20px;
  font-weight: 700;
  color: var(--el-text-color-primary);
}

.page-subtitle {
  margin: 0;
  font-size: 13px;
  color: var(--el-text-color-secondary);
}

.mb-6 {
  margin-bottom: 24px;
}

.mb-4 {
  margin-bottom: 16px;
}

.mb-3 {
  margin-bottom: 12px;
}

.mr-2 {
  margin-right: 8px;
}

.mr-3 {
  margin-right: 12px;
}

.ml-2 {
  margin-left: 8px;
}

.mt-2 {
  margin-top: 8px;
}

.flex-1 {
  flex: 1;
}

.w-full {
  width: 100%;
}

.stat-card {
  border-left: 4px solid;
  transition: transform 0.2s;
  margin-bottom: 12px;
}

.stat-card:hover {
  transform: translateY(-4px);
}

.stat-card-primary {
  border-left-color: var(--el-color-primary);
  background: linear-gradient(135deg, rgba(64, 158, 255, 0.08), rgba(64, 158, 255, 0.02));
}

.stat-card-success {
  border-left-color: var(--el-color-success);
  background: linear-gradient(135deg, rgba(103, 194, 58, 0.08), rgba(103, 194, 58, 0.02));
}

.stat-card-warning {
  border-left-color: var(--el-color-warning);
  background: linear-gradient(135deg, rgba(230, 162, 60, 0.08), rgba(230, 162, 60, 0.02));
}

.stat-card-info {
  border-left-color: var(--el-color-info);
  background: linear-gradient(135deg, rgba(144, 147, 153, 0.08), rgba(144, 147, 153, 0.02));
}

.stat-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.stat-value {
  font-size: 26px;
  font-weight: 700;
  color: var(--el-text-color-primary);
}

.stat-label {
  font-size: 13px;
  color: var(--el-text-color-secondary);
}

.stat-avatar {
  color: #fff;
}

.avatar-primary {
  background: var(--el-color-primary);
}

.avatar-success {
  background: var(--el-color-success);
}

.avatar-warning {
  background: var(--el-color-warning);
}

.avatar-info {
  background: var(--el-color-info);
}

.online-users-card {
  border: 2px dashed rgba(103, 194, 58, 0.3);
  background: rgba(103, 194, 58, 0.02);
}

.online-users-header {
  display: flex;
  align-items: center;
}

.online-users-title {
  display: flex;
  align-items: center;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.online-chip {
  margin: 4px 8px 4px 0;
  cursor: pointer;
}

.online-chip:hover {
  opacity: 0.85;
}

.chip-avatar {
  margin-right: 6px;
  background: var(--el-color-success);
  color: #fff;
  font-size: 12px;
  vertical-align: middle;
}

.chip-avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 50%;
}

.chip-name {
  vertical-align: middle;
}

.online-more {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.online-empty {
  text-align: center;
  color: var(--el-text-color-secondary);
}

.table-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.user-cell {
  display: flex;
  align-items: center;
  padding: 8px 0;
}

.user-avatar {
  background: var(--el-color-primary);
}

.user-avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.user-avatar-char {
  color: #fff;
  font-weight: 700;
}

.user-name {
  font-weight: 500;
  color: var(--el-text-color-primary);
}

.user-handle {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.online-cell {
  display: flex;
  align-items: center;
}

.online-text {
  font-size: 12px;
  color: var(--el-text-color-regular);
}

.online-last {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.danger-item {
  color: var(--el-color-danger);
}

.danger-item:hover {
  background: var(--el-color-danger-light-9);
}

.table-responsive {
  overflow-x: auto;
}

@media (max-width: 600px) {
  .stat-card {
    margin-bottom: 12px;
  }
}
</style>
