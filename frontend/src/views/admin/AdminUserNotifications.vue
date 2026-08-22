<template>
  <div class="page-container">
    <div class="admin-notifications">
      <!-- 页面标题 -->
      <div class="page-header mb-6">
        <div class="page-header-row">
          <div>
            <h1 class="page-title">用户通知管理</h1>
            <p class="page-subtitle">向指定用户发送单独通知，管理权限组</p>
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
        <el-col :xs="24" :sm="8">
          <el-card shadow="never" class="stat-card stat-card-primary">
            <div class="stat-row">
              <div>
                <div class="stat-value">{{ stats.totalPermissions }}</div>
                <div class="stat-label">权限组数量</div>
              </div>
              <el-avatar :size="48" class="stat-avatar avatar-primary">
                <el-icon><Lock /></el-icon>
              </el-avatar>
            </div>
          </el-card>
        </el-col>
        <el-col :xs="24" :sm="8">
          <el-card shadow="never" class="stat-card stat-card-success">
            <div class="stat-row">
              <div>
                <div class="stat-value">{{ stats.activeGroups }}</div>
                <div class="stat-label">活跃权限组</div>
              </div>
              <el-avatar :size="48" class="stat-avatar avatar-success">
                <el-icon><CircleCheck /></el-icon>
              </el-avatar>
            </div>
          </el-card>
        </el-col>
        <el-col :xs="24" :sm="8">
          <el-card shadow="never" class="stat-card stat-card-info">
            <div class="stat-row">
              <div>
                <div class="stat-value">{{ stats.defaultGroup }}</div>
                <div class="stat-label">默认权限组</div>
              </div>
              <el-avatar :size="48" class="stat-avatar avatar-info">
                <el-icon><Star /></el-icon>
              </el-avatar>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 标签页 -->
      <el-tabs v-model="activeTab" class="mb-4">
        <el-tab-pane label="发送通知" name="send-notification">
          <el-icon class="mr-1" style="vertical-align: -2px"><BellFilled /></el-icon>
          发送通知
          <el-card shadow="never" class="mt-3">
            <template #header>
              <span>发送单独通知</span>
            </template>
            <el-form ref="notificationFormRef" :model="notificationForm" label-position="top">
              <el-row :gutter="16">
                <el-col :xs="24" :md="12">
                  <el-form-item label="选择用户">
                    <el-select
                      v-model="notificationForm.user_id"
                      filterable
                      clearable
                      :loading="loadingUsers"
                      placeholder="搜索并选择用户"
                      class="w-full mb-3"
                    >
                      <el-option
                        v-for="u in userOptions"
                        :key="u.id"
                        :value="u.id"
                        :label="u.display_name || u.username"
                      >
                        <div class="d-flex align-center">
                          <el-avatar :size="20" class="mr-2">
                            <img
                              v-if="u.avatar"
                              :src="u.avatar"
                              style="width:100%;height:100%;object-fit:cover"
                            />
                            <template v-else>{{ (u.display_name || u.username)[0] }}</template>
                          </el-avatar>
                          <span>{{ u.display_name || u.username }}</span>
                          <span class="user-username ml-2">@{{ u.username }}</span>
                        </div>
                      </el-option>
                    </el-select>
                  </el-form-item>
                </el-col>
                <el-col :xs="24" :md="12">
                  <el-form-item label="通知类型">
                    <el-select v-model="notificationForm.type" class="w-full mb-3">
                      <el-option
                        v-for="t in notificationTypes"
                        :key="t.value"
                        :value="t.value"
                        :label="t.title"
                      />
                    </el-select>
                  </el-form-item>
                </el-col>
                <el-col :xs="24">
                  <el-form-item label="通知标题">
                    <el-input v-model="notificationForm.title" placeholder="请输入通知标题" class="mb-3" />
                  </el-form-item>
                </el-col>
                <el-col :xs="24">
                  <el-form-item label="通知内容">
                    <el-input
                      v-model="notificationForm.content"
                      type="textarea"
                      :rows="4"
                      placeholder="请输入通知内容"
                      class="mb-3"
                    />
                  </el-form-item>
                </el-col>
                <el-col :xs="24" :md="12">
                  <el-form-item label="优先级">
                    <el-select v-model="notificationForm.priority" class="w-full">
                      <el-option
                        v-for="p in priorityOptions"
                        :key="p.value"
                        :value="p.value"
                        :label="p.title"
                      />
                    </el-select>
                  </el-form-item>
                </el-col>
                <el-col :xs="24" :md="12">
                  <el-form-item label="跳转链接（可选）">
                    <el-input v-model="notificationForm.link" placeholder="/article/123" />
                  </el-form-item>
                </el-col>
              </el-row>
            </el-form>
            <div class="form-actions">
              <span class="flex-1"></span>
              <el-button @click="resetNotificationForm">重置</el-button>
              <el-button type="primary" :loading="sending" @click="sendNotification">
                <el-icon class="mr-2"><Promotion /></el-icon>
                发送通知
              </el-button>
            </div>
          </el-card>
        </el-tab-pane>

        <el-tab-pane label="权限组管理" name="permission-groups">
          <el-card shadow="never">
            <template #header>
              <div class="group-header">
                <span>权限组列表</span>
                <div>
                  <el-button
                    size="small"
                    type="primary"
                    plain
                    :loading="initializing"
                    @click="initDefaultGroups"
                  >
                    初始化默认组
                  </el-button>
                  <el-button size="small" type="primary" @click="openCreateGroup">
                    <el-icon class="mr-1"><Plus /></el-icon>
                    创建权限组
                  </el-button>
                </div>
              </div>
            </template>

            <div class="table-responsive">
              <el-table :data="permissionGroups" v-loading="loading" stripe>
              <el-table-column prop="name" label="名称" min-width="120" />
              <el-table-column prop="description" label="描述" min-width="180" show-overflow-tooltip />
              <el-table-column label="级别" width="100">
                <template #default="{ row }">
                  <el-tag size="small" type="primary" effect="plain">Level {{ row.level }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column label="默认" width="100">
                <template #default="{ row }">
                  <el-tag v-if="row.is_default" size="small" type="success" effect="plain">默认</el-tag>
                  <span v-else class="text-empty">-</span>
                </template>
              </el-table-column>
              <el-table-column label="权限" min-width="220">
                <template #default="{ row }">
                  <el-tag
                    v-for="perm in parsePermissions(row.permissions).slice(0, 3)"
                    :key="perm"
                    size="small"
                    effect="plain"
                    class="mr-1 mb-1"
                  >
                    {{ perm }}
                  </el-tag>
                  <span v-if="parsePermissions(row.permissions).length > 3" class="text-secondary">
                    +{{ parsePermissions(row.permissions).length - 3 }} more
                  </span>
                </template>
              </el-table-column>
              <el-table-column label="操作" width="120">
                <template #default="{ row }">
                  <el-button link type="primary" :icon="Edit" @click="editGroup(row)" />
                  <el-button link type="danger" :icon="Delete" @click="deleteGroup(row)" />
                </template>
              </el-table-column>
              </el-table>
            </div>
          </el-card>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { adminUserNotificationApi, permissionGroupApi, adminUserApi } from '../../api'
import { confirm, success, error, warning } from '../../utils/message'
import { Refresh, Lock, CircleCheck, Star, BellFilled, Promotion, Plus, Edit, Delete } from '@element-plus/icons-vue'
import { jcCloseAll, jcOpenHtml } from '@/utils/jcu'

const activeTab = ref('send-notification')
const loading = ref(false)
const sending = ref(false)
const initializing = ref(false)
const saving = ref(false)
const loadingUsers = ref(false)
const permissionGroups = ref([])
const userOptions = ref([])

const stats = ref({
  totalPermissions: 0,
  activeGroups: 0,
  defaultGroup: 0
})

const notificationForm = ref({
  user_id: null,
  type: 'system',
  title: '',
  content: '',
  priority: 'normal',
  link: ''
})

const editingGroup = ref(null)

const notificationTypes = [
  { title: '系统通知', value: 'system' },
  { title: '警告', value: 'warning' },
  { title: '晋升通知', value: 'promotion' },
  { title: '提醒', value: 'reminder' },
  { title: '活动通知', value: 'activity' }
]

const priorityOptions = [
  { title: '高优先级', value: 'high' },
  { title: '普通', value: 'normal' },
  { title: '低优先级', value: 'low' }
]

const availablePermissions = [
  'article:read', 'article:create', 'article:edit:own', 'article:edit:any',
  'article:delete:own', 'article:delete:any', 'article:pin', 'article:featured',
  'comment:create', 'comment:edit:own', 'comment:edit:any',
  'comment:delete:own', 'comment:delete:any',
  'user:profile:view', 'user:profile:edit', 'user:ban',
  'category:manage', 'report:view', 'report:handle',
  '*'
]

const loadPermissionGroups = async () => {
  loading.value = true
  try {
    const response = await permissionGroupApi.getGroups()
    permissionGroups.value = response.data.groups || []

    stats.value.totalPermissions = permissionGroups.value.length
    stats.value.activeGroups = permissionGroups.value.filter(g => g.is_active).length
    stats.value.defaultGroup = permissionGroups.value.filter(g => g.is_default).length
  } catch (err) {
    console.error('加载权限组失败:', err)
    error('加载权限组失败')
  } finally {
    loading.value = false
  }
}

const loadUsers = async () => {
  loadingUsers.value = true
  try {
    const response = await adminUserApi.getUsers()
    userOptions.value = response.data.users || []
  } catch (err) {
    console.error('加载用户列表失败:', err)
    error('加载用户列表失败')
  } finally {
    loadingUsers.value = false
  }
}

const refreshData = () => {
  loadPermissionGroups()
}

const sendNotification = async () => {
  if (!notificationForm.value.user_id) {
    warning('请选择用户')
    return
  }
  if (!notificationForm.value.title || !notificationForm.value.content) {
    warning('请填写标题和内容')
    return
  }

  sending.value = true
  try {
    await adminUserNotificationApi.sendNotification(notificationForm.value)
    success('通知发送成功')
    resetNotificationForm()
  } catch (err) {
    console.error('发送通知失败:', err)
    error('发送失败')
  } finally {
    sending.value = false
  }
}

const resetNotificationForm = () => {
  notificationForm.value = {
    user_id: null,
    type: 'system',
    title: '',
    content: '',
    priority: 'normal',
    link: ''
  }
}

const parsePermissions = (permissions) => {
  try {
    return JSON.parse(permissions)
  } catch {
    return []
  }
}

const initDefaultGroups = async () => {
  const ok = await confirm('确定要初始化默认权限组吗？这将创建4个默认权限组。')
  if (!ok) return

  initializing.value = true
  try {
    await permissionGroupApi.initDefaults()
    success('默认权限组初始化成功')
    loadPermissionGroups()
  } catch (err) {
    console.error('初始化失败:', err)
    error(err.response?.data?.error || '初始化失败')
  } finally {
    initializing.value = false
  }
}

const openCreateGroup = () => {
  editingGroup.value = null
  openGroupDialog({
    name: '',
    description: '',
    level: 1,
    is_default: false,
    permissions: []
  })
}

const editGroup = (group) => {
  editingGroup.value = group
  openGroupDialog({
    name: group.name,
    description: group.description,
    level: group.level,
    is_default: group.is_default,
    permissions: parsePermissions(group.permissions)
  })
}

const openGroupDialog = (form) => {
  const isEdit = !!editingGroup.value
  const esc = (s) => String(s ?? '').replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;')
  const field = (label, inner, hint) =>
    `<div style="margin-bottom:14px;"><div style="margin:0 0 4px;font-weight:600;font-size:13px;color:var(--jc-text,#333);">${label}</div>${inner}${hint ? `<div style="margin-top:4px;font-size:12px;color:var(--jc-text-2,#64748b);">${hint}</div>` : ''}</div>`

  const permChips = availablePermissions
    .map((p) => {
      const checked = form.permissions.includes(p) ? ' checked' : ''
      return `<label style="display:inline-flex;align-items:center;gap:4px;margin:0 8px 8px 0;padding:6px 10px;border:1px solid var(--jc-border,#e6e9f2);border-radius:8px;cursor:pointer;font-size:13px;color:var(--jc-text,#0f172a);background:var(--jc-surface,#fff);"><input type="checkbox" class="jf-perm" value="${esc(p)}"${checked} /> ${esc(p)}</label>`
    })
    .join('')

  const content =
    field('权限组名称', `<input data-group="name" class="jc-modal__input" type="text" value="${esc(form.name)}" placeholder="请输入权限组名称" required />`) +
    field('描述', `<textarea data-group="description" class="jc-modal__input" rows="2" style="resize:vertical;">${esc(form.description)}</textarea>`) +
    field('权限级别', `<input data-group="level" class="jc-modal__input" type="number" min="1" value="${form.level}" /><span style="margin-left:8px;font-size:12px;color:var(--jc-text-2,#64748b);">数字越大权限越高</span>`) +
    field('设为默认权限组', `<label style="display:inline-flex;align-items:center;gap:6px;cursor:pointer;font-size:14px;color:var(--jc-text,#0f172a);"><input type="checkbox" data-group="is_default" ${form.is_default ? 'checked' : ''} /> 设为默认</label>`) +
    field('权限列表', `<div>${permChips}</div>`, '选择该权限组包含的权限')

  jcOpenHtml({
    title: isEdit ? '编辑权限组' : '创建权限组',
    content,
    width: 700,
    size: 'md',
    buttons: [
      { text: '取消', type: 'default', action: () => jcCloseAll() },
      {
        text: isEdit ? '保存' : '创建',
        type: 'primary',
        action: (inst) => {
          const root = inst.modalContent
          const name = (root.querySelector('[data-group="name"]')?.value || '').trim()
          if (!name) {
            warning('请输入权限组名称')
            return
          }
          const level = Number(root.querySelector('[data-group="level"]')?.value || 1)
          const perms = Array.from(root.querySelectorAll('.jf-perm:checked')).map((el) => el.value)
          saveGroup({
            name,
            description: root.querySelector('[data-group="description"]')?.value || '',
            level,
            is_default: root.querySelector('[data-group="is_default"]')?.checked || false,
            permissions: perms
          })
        }
      }
    ]
  })
}

const saveGroup = async (form) => {
  saving.value = true
  try {
    if (editingGroup.value) {
      await permissionGroupApi.updateGroup(editingGroup.value.id, form)
    } else {
      await permissionGroupApi.createGroup(form)
    }
    success('保存成功')
    jcCloseAll()
    editingGroup.value = null
    loadPermissionGroups()
  } catch (err) {
    console.error('保存失败:', err)
    error('保存失败')
  } finally {
    saving.value = false
  }
}

const deleteGroup = async (group) => {
  const ok = await confirm(`确定要删除权限组 "${group.name}" 吗？`)
  if (!ok) return

  try {
    await permissionGroupApi.deleteGroup(group.id)
    success('删除成功')
    loadPermissionGroups()
  } catch (err) {
    console.error('删除失败:', err)
    error(err.response?.data?.error || '删除失败')
  }
}

onMounted(() => {
  loadPermissionGroups()
  loadUsers()
})
</script>

<style scoped>
.admin-notifications {
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
  border-left-color: var(--el-color-primary);
  background: linear-gradient(135deg, rgba(64, 158, 255, 0.08), rgba(64, 158, 255, 0.02));
}

.stat-card-success {
  border-left-color: var(--el-color-success);
  background: linear-gradient(135deg, rgba(103, 194, 58, 0.08), rgba(103, 194, 58, 0.02));
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
  font-size: 24px;
  font-weight: 700;
}

.stat-label {
  font-size: 13px;
  color: var(--el-text-color-secondary);
}

.stat-avatar {
  --el-avatar-bg-color: transparent;
  border-radius: 8px;
}

.avatar-primary {
  color: var(--el-color-primary);
}

.avatar-success {
  color: var(--el-color-success);
}

.avatar-info {
  color: var(--el-color-info);
}

.group-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.form-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
  padding-top: 8px;
}

.flex-1 {
  flex: 1;
}

.user-username {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.text-empty {
  color: var(--el-text-color-placeholder);
}

.text-secondary {
  color: var(--el-text-color-secondary);
  font-size: 12px;
}

.d-flex {
  display: flex;
}

.align-center {
  align-items: center;
}

.page-title {
  font-size: 20px;
  font-weight: 700;
  margin: 0 0 4px;
}

.page-subtitle {
  margin: 0;
  color: var(--el-text-color-secondary);
  font-size: 13px;
}

.table-responsive {
  overflow-x: auto;
}
</style>
