<template>
  <div class="admin-notifications">
    <div class="page-header mb-6">
      <div class="flex items-center justify-between flex-wrap gap-4">
        <div>
          <h1 style="font-size: 18px; font-weight: 600; margin-bottom: 5px;">用户通知管理</h1>
          <p style="color: #999; font-size: 14px;">向指定用户发送单独通知，管理权限组</p>
        </div>
        <div class="flex gap-3">
          <button class="layui-btn" @click="refreshData" :disabled="loading">
            <i class="fa-solid fa-refresh mr-2"></i>
            刷新
          </button>
        </div>
      </div>
    </div>

    <div class="layui-row mb-6">
      <div class="layui-col-xs12 layui-col-sm4">
        <div class="stat-card stat-card-primary">
          <div class="stat-content">
            <div style="font-size: 28px; font-weight: 700;">{{ stats.totalPermissions }}</div>
            <div style="color: #666; font-size: 14px;">权限组数量</div>
          </div>
          <div class="stat-icon">
            <i class="fa-solid fa-shield-halved"></i>
          </div>
        </div>
      </div>
      <div class="layui-col-xs12 layui-col-sm4">
        <div class="stat-card stat-card-success">
          <div class="stat-content">
            <div style="font-size: 28px; font-weight: 700;">{{ stats.activeGroups }}</div>
            <div style="color: #666; font-size: 14px;">活跃权限组</div>
          </div>
          <div class="stat-icon" style="background: rgba(82, 196, 26, 0.1); color: #52C41A;">
            <i class="fa-solid fa-check-circle"></i>
          </div>
        </div>
      </div>
      <div class="layui-col-xs12 layui-col-sm4">
        <div class="stat-card stat-card-info">
          <div class="stat-content">
            <div style="font-size: 28px; font-weight: 700;">{{ stats.defaultGroup }}</div>
            <div style="color: #666; font-size: 14px;">默认权限组</div>
          </div>
          <div class="stat-icon" style="background: rgba(30, 159, 255, 0.1); color: #1E9FFF;">
            <i class="fa-solid fa-star"></i>
          </div>
        </div>
      </div>
    </div>

    <div class="layui-tab layui-tab-brief mb-4">
      <ul class="layui-tab-title">
        <li :class="{ 'layui-this': activeTab === 'send-notification' }" @click="activeTab = 'send-notification'">
          <i class="fa-solid fa-bell mr-2"></i>发送通知
        </li>
        <li :class="{ 'layui-this': activeTab === 'permission-groups' }" @click="activeTab = 'permission-groups'">
          <i class="fa-solid fa-shield-halved mr-2"></i>权限组管理
        </li>
      </ul>
    </div>

    <div v-show="activeTab === 'send-notification'" class="layui-card mb-4">
      <div class="layui-card-header">发送单独通知</div>
      <div class="layui-card-body">
        <div class="layui-row">
          <div class="layui-col-xs12 layui-col-md6">
            <div class="layui-form-item">
              <label class="layui-form-label">选择用户</label>
              <div class="layui-input-block">
                <select v-model="notificationForm.user_id" class="layui-select">
                  <option value="">请选择用户</option>
                  <option v-for="user in userOptions" :key="user.id" :value="user.id">{{ user.display_name || user.username }}</option>
                </select>
              </div>
            </div>
          </div>
          <div class="layui-col-xs12 layui-col-md6">
            <div class="layui-form-item">
              <label class="layui-form-label">通知类型</label>
              <div class="layui-input-block">
                <select v-model="notificationForm.type" class="layui-select">
                  <option v-for="type in notificationTypes" :key="type.value" :value="type.value">{{ type.title }}</option>
                </select>
              </div>
            </div>
          </div>
          <div class="layui-col-xs12">
            <div class="layui-form-item">
              <label class="layui-form-label">通知标题</label>
              <div class="layui-input-block">
                <input type="text" v-model="notificationForm.title" class="layui-input" />
              </div>
            </div>
          </div>
          <div class="layui-col-xs12">
            <div class="layui-form-item">
              <label class="layui-form-label">通知内容</label>
              <div class="layui-input-block">
                <textarea v-model="notificationForm.content" rows="4" class="layui-textarea"></textarea>
              </div>
            </div>
          </div>
          <div class="layui-col-xs12 layui-col-md6">
            <div class="layui-form-item">
              <label class="layui-form-label">优先级</label>
              <div class="layui-input-block">
                <select v-model="notificationForm.priority" class="layui-select">
                  <option v-for="priority in priorityOptions" :key="priority.value" :value="priority.value">{{ priority.title }}</option>
                </select>
              </div>
            </div>
          </div>
          <div class="layui-col-xs12 layui-col-md6">
            <div class="layui-form-item">
              <label class="layui-form-label">跳转链接（可选）</label>
              <div class="layui-input-block">
                <input type="text" v-model="notificationForm.link" placeholder="/article/123" class="layui-input" />
              </div>
            </div>
          </div>
        </div>
      </div>
      <div class="layui-card-body border-t flex justify-end gap-2">
        <button class="layui-btn layui-btn-primary" @click="resetNotificationForm">重置</button>
        <button class="layui-btn" @click="sendNotification" :disabled="sending">
          <i class="fa-solid fa-paper-plane mr-2"></i>
          发送通知
        </button>
      </div>
    </div>

    <div v-show="activeTab === 'permission-groups'" class="layui-card">
      <div class="layui-card-header flex items-center justify-between">
        <span>权限组列表</span>
        <div class="flex gap-2">
          <button class="layui-btn layui-btn-primary layui-btn-sm" @click="initDefaultGroups" :disabled="initializing">
            初始化默认组
          </button>
          <button class="layui-btn layui-btn-sm" @click="showCreateGroupDialog = true">
            <i class="fa-solid fa-plus mr-1"></i>
            创建权限组
          </button>
        </div>
      </div>
      <div class="layui-card-body">
        <table class="layui-table" v-if="permissionGroups.length > 0">
          <thead>
            <tr>
              <th>名称</th>
              <th>描述</th>
              <th>级别</th>
              <th>默认</th>
              <th>权限</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="group in permissionGroups" :key="group.id">
              <td>{{ group.name }}</td>
              <td>{{ group.description }}</td>
              <td><span class="layui-badge">Level {{ group.level }}</span></td>
              <td><span v-if="group.is_default" class="layui-badge layui-bg-green">默认</span><span v-else>-</span></td>
              <td>
                <span v-for="perm in parsePermissions(group.permissions).slice(0, 3)" :key="perm" class="layui-badge layui-bg-gray mr-1 mb-1">{{ perm }}</span>
                <span v-if="parsePermissions(group.permissions).length > 3" class="text-muted text-sm">+{{ parsePermissions(group.permissions).length - 3 }} more</span>
              </td>
              <td>
                <button class="layui-btn layui-btn-sm layui-btn-primary" @click="editGroup(group)">
                  <i class="fa-solid fa-pencil"></i>
                </button>
                <button class="layui-btn layui-btn-sm layui-btn-danger" @click="deleteGroup(group)">
                  <i class="fa-solid fa-trash"></i>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
        <div v-else class="text-center py-8 text-muted">暂无权限组数据</div>
      </div>
    </div>

    <div v-if="showCreateGroupDialog" class="modal-overlay" @click="closeGroupDialog">
      <div class="modal-content" @click.stop>
        <div class="modal-header" style="background: #1E9FFF; color: white; padding: 15px; border-radius: 4px 4px 0 0;">
          <i class="fa-solid" :class="editingGroup ? 'fa-pencil' : 'fa-plus'" style="margin-right: 8px;"></i>
          {{ editingGroup ? '编辑权限组' : '创建权限组' }}
        </div>
        <div class="modal-body" style="padding: 20px;">
          <div class="layui-form-item">
            <label class="layui-form-label">权限组名称</label>
            <div class="layui-input-block">
              <input type="text" v-model="groupForm.name" class="layui-input" />
            </div>
          </div>
          <div class="layui-form-item">
            <label class="layui-form-label">描述</label>
            <div class="layui-input-block">
              <textarea v-model="groupForm.description" rows="2" class="layui-textarea"></textarea>
            </div>
          </div>
          <div class="layui-form-item">
            <label class="layui-form-label">权限级别</label>
            <div class="layui-input-block">
              <input type="number" v-model.number="groupForm.level" placeholder="数字越大权限越高" class="layui-input" />
            </div>
          </div>
          <div class="layui-form-item">
            <label class="layui-form-label">设为默认权限组</label>
            <div class="layui-input-block">
              <input type="checkbox" v-model="groupForm.is_default" lay-skin="switch" />
            </div>
          </div>
          <div class="layui-form-item">
            <label class="layui-form-label">权限列表</label>
            <div class="layui-input-block">
              <select v-model="groupForm.permissions" multiple class="layui-select" style="height: 150px;">
                <option v-for="perm in availablePermissions" :key="perm" :value="perm">{{ perm }}</option>
              </select>
            </div>
          </div>
        </div>
        <div class="modal-footer" style="padding: 15px; border-top: 1px solid #E8E8E8; display: flex; justify-content: flex-end; gap: 10px;">
          <button class="layui-btn layui-btn-primary" @click="closeGroupDialog">取消</button>
          <button class="layui-btn" @click="saveGroup" :disabled="saving">
            {{ editingGroup ? '保存' : '创建' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { userNotificationApi, permissionGroupApi, adminUserApi } from '../../api'

export default {
  name: 'AdminUserNotifications',
  setup() {
    const activeTab = ref('send-notification')
    const loading = ref(false)
    const sending = ref(false)
    const initializing = ref(false)
    const saving = ref(false)
    const loadingUsers = ref(false)
    const permissionGroups = ref([])
    const userOptions = ref([])
    const userSearch = ref('')

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

    const groupForm = ref({
      name: '',
      description: '',
      level: 1,
      is_default: false,
      permissions: []
    })

    const editingGroup = ref(null)
    const showCreateGroupDialog = ref(false)

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
      } catch (error) {
        console.error('加载权限组失败:', error)
      } finally {
        loading.value = false
      }
    }

    const loadUsers = async () => {
      loadingUsers.value = true
      try {
        const response = await adminUserApi.getUsers()
        userOptions.value = response.data.users || []
      } catch (error) {
        console.error('加载用户列表失败:', error)
      } finally {
        loadingUsers.value = false
      }
    }

    const refreshData = () => {
      loadPermissionGroups()
    }

    const sendNotification = async () => {
      if (!notificationForm.value.user_id) {
        alert('请选择用户')
        return
      }
      if (!notificationForm.value.title || !notificationForm.value.content) {
        alert('请填写标题和内容')
        return
      }

      sending.value = true
      try {
        await userNotificationApi.sendNotification(notificationForm.value)
        alert('通知发送成功')
        resetNotificationForm()
      } catch (error) {
        console.error('发送通知失败:', error)
        alert('发送失败')
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
      if (!confirm('确定要初始化默认权限组吗？这将创建4个默认权限组。')) return

      initializing.value = true
      try {
        await permissionGroupApi.initDefaults()
        alert('默认权限组初始化成功')
        loadPermissionGroups()
      } catch (error) {
        console.error('初始化失败:', error)
        alert('初始化失败: ' + (error.response?.data?.error || '未知错误'))
      } finally {
        initializing.value = false
      }
    }

    const editGroup = (group) => {
      editingGroup.value = group
      groupForm.value = {
        name: group.name,
        description: group.description,
        level: group.level,
        is_default: group.is_default,
        permissions: parsePermissions(group.permissions)
      }
      showCreateGroupDialog.value = true
    }

    const saveGroup = async () => {
      if (!groupForm.value.name) {
        alert('请输入权限组名称')
        return
      }

      saving.value = true
      try {
        if (editingGroup.value) {
          await permissionGroupApi.updateGroup(editingGroup.value.id, groupForm.value)
        } else {
          await permissionGroupApi.createGroup(groupForm.value)
        }
        closeGroupDialog()
        loadPermissionGroups()
      } catch (error) {
        console.error('保存失败:', error)
        alert('保存失败')
      } finally {
        saving.value = false
      }
    }

    const deleteGroup = async (group) => {
      if (!confirm(`确定要删除权限组 "${group.name}" 吗？`)) return

      try {
        await permissionGroupApi.deleteGroup(group.id)
        loadPermissionGroups()
      } catch (error) {
        console.error('删除失败:', error)
        alert('删除失败: ' + (error.response?.data?.error || '未知错误'))
      }
    }

    const closeGroupDialog = () => {
      showCreateGroupDialog.value = false
      editingGroup.value = null
      groupForm.value = {
        name: '',
        description: '',
        level: 1,
        is_default: false,
        permissions: []
      }
    }

    onMounted(() => {
      loadPermissionGroups()
      loadUsers()
    })

    return {
      activeTab,
      loading,
      sending,
      initializing,
      saving,
      loadingUsers,
      permissionGroups,
      userOptions,
      userSearch,
      stats,
      notificationForm,
      groupForm,
      editingGroup,
      showCreateGroupDialog,
      notificationTypes,
      priorityOptions,
      availablePermissions,
      loadPermissionGroups,
      loadUsers,
      refreshData,
      sendNotification,
      resetNotificationForm,
      parsePermissions,
      initDefaultGroups,
      editGroup,
      saveGroup,
      deleteGroup,
      closeGroupDialog
    }
  }
}
</script>

<style scoped>
.admin-notifications {
  max-width: 1600px;
  margin: 0 auto;
}

.stat-card {
  border-left: 4px solid;
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: 0 2px 8px rgba(0,0,0,0.06);
  transition: transform 0.2s;
}

.stat-card:hover {
  transform: translateY(-4px);
}

.stat-card-primary {
  border-left-color: #1E9FFF;
}

.stat-card-primary .stat-icon {
  background: rgba(30, 159, 255, 0.1);
  color: #1E9FFF;
}

.stat-icon {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}

.modal-content {
  background: #fff;
  border-radius: 4px;
  width: 90%;
  max-width: 600px;
  max-height: 90vh;
  overflow-y: auto;
}

.border-t {
  border-top: 1px solid #E8E8E8;
}
</style>