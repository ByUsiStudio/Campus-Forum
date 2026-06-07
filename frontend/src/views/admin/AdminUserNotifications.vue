<template>
  <div class="admin-notifications">
    <!-- 页面标题 -->
    <div class="page-header mb-6">
      <div class="d-flex align-center justify-space-between flex-wrap ga-4">
        <div>
          <h1 class="text-h5 font-weight-bold mb-1">用户通知管理</h1>
          <p class="text-body-2 text-medium-emphasis">向指定用户发送单独通知，管理权限组</p>
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
      <v-col cols="12" sm="4">
        <v-card elevation="0" class="stat-card stat-card-primary">
          <v-card-text>
            <div class="d-flex align-center justify-space-between">
              <div>
                <div class="text-h4 font-weight-bold">{{ stats.totalPermissions }}</div>
                <div class="text-body-2 text-medium-emphasis">权限组数量</div>
              </div>
              <v-avatar color="primary" size="48" rounded="lg">
                <v-icon>mdi-shield-account</v-icon>
              </v-avatar>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" sm="4">
        <v-card elevation="0" class="stat-card stat-card-success">
          <v-card-text>
            <div class="d-flex align-center justify-space-between">
              <div>
                <div class="text-h4 font-weight-bold">{{ stats.activeGroups }}</div>
                <div class="text-body-2 text-medium-emphasis">活跃权限组</div>
              </div>
              <v-avatar color="success" size="48" rounded="lg">
                <v-icon>mdi-check-circle</v-icon>
              </v-avatar>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" sm="4">
        <v-card elevation="0" class="stat-card stat-card-info">
          <v-card-text>
            <div class="d-flex align-center justify-space-between">
              <div>
                <div class="text-h4 font-weight-bold">{{ stats.defaultGroup }}</div>
                <div class="text-body-2 text-medium-emphasis">默认权限组</div>
              </div>
              <v-avatar color="info" size="48" rounded="lg">
                <v-icon>mdi-star</v-icon>
              </v-avatar>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- 标签页 -->
    <v-tabs v-model="activeTab" color="primary" class="mb-4">
      <v-tab value="send-notification">
        <v-icon start>mdi-bell-ring</v-icon>
        发送通知
      </v-tab>
      <v-tab value="permission-groups">
        <v-icon start>mdi-shield-account</v-icon>
        权限组管理
      </v-tab>
    </v-tabs>

    <v-window v-model="activeTab" class="mb-4">
      <!-- 发送通知标签页 -->
      <v-window-item value="send-notification">
      <v-card elevation="0">
        <v-card-title>发送单独通知</v-card-title>
        <v-card-text>
          <v-form ref="notificationForm">
            <v-row>
              <v-col cols="12" md="6">
                <v-autocomplete
                  v-model="notificationForm.user_id"
                  :items="userOptions"
                  item-title="display_name"
                  item-value="id"
                  label="选择用户"
                  variant="outlined"
                  :loading="loadingUsers"
                  :search-input.sync="userSearch"
                  cache-items
                  hide-no-data
                  class="mb-3"
                >
                  <template v-slot:item="{ props, item }">
                    <v-list-item v-bind="props" :title="item.raw.display_name || item.raw.username">
                      <template v-slot:prepend>
                        <v-avatar size="32">
                          <v-img v-if="item.raw.avatar" :src="item.raw.avatar" />
                          <span v-else class="text-primary">{{ (item.raw.display_name || item.raw.username)[0] }}</span>
                        </v-avatar>
                      </template>
                      <template v-slot:subtitle>
                        @{{ item.raw.username }}
                      </template>
                    </v-list-item>
                  </template>
                </v-autocomplete>
              </v-col>
              <v-col cols="12" md="6">
                <v-select
                  v-model="notificationForm.type"
                  :items="notificationTypes"
                  label="通知类型"
                  variant="outlined"
                  class="mb-3"
                />
              </v-col>
              <v-col cols="12">
                <v-text-field
                  v-model="notificationForm.title"
                  label="通知标题"
                  variant="outlined"
                  class="mb-3"
                />
              </v-col>
              <v-col cols="12">
                <v-textarea
                  v-model="notificationForm.content"
                  label="通知内容"
                  variant="outlined"
                  rows="4"
                  class="mb-3"
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-select
                  v-model="notificationForm.priority"
                  :items="priorityOptions"
                  label="优先级"
                  variant="outlined"
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="notificationForm.link"
                  label="跳转链接（可选）"
                  variant="outlined"
                  placeholder="/article/123"
                />
              </v-col>
            </v-row>
          </v-form>
        </v-card-text>
        <v-card-actions class="pa-4">
          <v-spacer />
          <v-btn variant="outlined" @click="resetNotificationForm">重置</v-btn>
          <v-btn color="primary" @click="sendNotification" :loading="sending" prepend-icon="mdi-send">
            发送通知
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-window-item>

    <!-- 权限组管理标签页 -->
    <v-window-item value="permission-groups">
      <v-card elevation="0">
        <v-card-title class="d-flex align-center justify-space-between">
          <span>权限组列表</span>
          <div class="d-flex ga-2">
            <v-btn size="small" variant="tonal" color="primary" @click="initDefaultGroups" :loading="initializing">
              初始化默认组
            </v-btn>
            <v-btn size="small" color="primary" @click="showCreateGroupDialog = true" prepend-icon="mdi-plus">
              创建权限组
            </v-btn>
          </div>
        </v-card-title>

        <v-data-table
          :headers="groupHeaders"
          :items="permissionGroups"
          :loading="loading"
          class="elevation-0"
        >
          <template v-slot:item.level="{ item }">
            <v-chip size="small" color="primary" variant="tonal">
              Level {{ item.level }}
            </v-chip>
          </template>

          <template v-slot:item.is_default="{ item }">
            <v-chip v-if="item.is_default" size="small" color="success" variant="tonal">
              默认
            </v-chip>
            <span v-else class="text-medium-emphasis">-</span>
          </template>

          <template v-slot:item.permissions="{ item }">
            <v-chip size="small" variant="tonal" class="mr-1 mb-1" v-for="perm in parsePermissions(item.permissions).slice(0, 3)" :key="perm">
              {{ perm }}
            </v-chip>
            <span v-if="parsePermissions(item.permissions).length > 3" class="text-caption text-medium-emphasis">
              +{{ parsePermissions(item.permissions).length - 3 }} more
            </span>
          </template>

          <template v-slot:item.actions="{ item }">
            <v-btn icon variant="text" size="small" @click="editGroup(item)">
              <v-icon>mdi-pencil</v-icon>
            </v-btn>
            <v-btn icon variant="text" size="small" color="error" @click="deleteGroup(item)">
              <v-icon>mdi-delete</v-icon>
            </v-btn>
          </template>
        </v-data-table>
      </v-card>
    </v-window-item>
    </v-window>

    <!-- 创建/编辑权限组对话框 -->
    <v-dialog v-model="showCreateGroupDialog" max-width="700">
      <v-card>
        <v-card-title class="bg-primary text-white pa-4">
          <v-icon class="mr-2">{{ editingGroup ? 'mdi-pencil' : 'mdi-plus' }}</v-icon>
          {{ editingGroup ? '编辑权限组' : '创建权限组' }}
        </v-card-title>
        <v-card-text class="pa-4">
          <v-form ref="groupForm">
            <v-text-field
              v-model="groupForm.name"
              label="权限组名称"
              variant="outlined"
              class="mb-3"
            />
            <v-textarea
              v-model="groupForm.description"
              label="描述"
              variant="outlined"
              rows="2"
              class="mb-3"
            />
            <v-text-field
              v-model.number="groupForm.level"
              label="权限级别"
              type="number"
              variant="outlined"
              hint="数字越大权限越高"
              persistent-hint
              class="mb-3"
            />
            <v-switch
              v-model="groupForm.is_default"
              label="设为默认权限组"
              color="success"
              class="mb-3"
            />
            <v-select
              v-model="groupForm.permissions"
              :items="availablePermissions"
              label="权限列表"
              variant="outlined"
              multiple
              chips
              closable-chips
              hint="选择该权限组包含的权限"
              persistent-hint
            />
          </v-form>
        </v-card-text>
        <v-card-actions class="pa-4">
          <v-spacer />
          <v-btn variant="text" @click="closeGroupDialog">取消</v-btn>
          <v-btn color="primary" @click="saveGroup" :loading="saving">
            {{ editingGroup ? '保存' : '创建' }}
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
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

    const groupHeaders = [
      { title: '名称', key: 'name' },
      { title: '描述', key: 'description' },
      { title: '级别', key: 'level', width: '100px' },
      { title: '默认', key: 'is_default', width: '100px' },
      { title: '权限', key: 'permissions' },
      { title: '操作', key: 'actions', width: '120px', sortable: false }
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
      groupHeaders,
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

.stat-card-info {
  border-left-color: rgb(var(--v-theme-info));
  background: linear-gradient(135deg, rgba(var(--v-theme-info), 0.08), rgba(var(--v-theme-info), 0.02));
}

.page-header h1 {
  color: rgb(var(--v-theme-on-surface));
}
</style>
