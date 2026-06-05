<template>
  <v-container fluid class="pa-6">
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
  </v-container>

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
          class="mt-4"
        ></v-text-field>
      </v-card-text>
      <v-card-actions class="dialog-actions">
        <v-btn variant="text" @click="editUserDialog.show = false">取消</v-btn>
        <v-btn color="primary" variant="flat" @click="handleEditUser">保存</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-dialog v-model="banDialog.show" max-width="480">
    <v-card class="dialog-card">
      <v-card-title class="dialog-title">
        <v-icon class="title-icon">mdi-ban</v-icon>
        封禁用户
      </v-card-title>
      <v-card-text class="dialog-body">
        <div class="user-info-card">
          <UserAvatar :user="banDialog.user || {}" :size="48" />
          <div class="user-info-text">
            <div class="user-name">{{ banDialog.user?.display_name }}</div>
            <div class="user-meta">ID: {{ banDialog.user?.id }}</div>
          </div>
        </div>
        <v-text-field
          v-model="banDialog.reason"
          label="封禁原因"
          variant="outlined"
          class="mt-4"
        ></v-text-field>
      </v-card-text>
      <v-card-actions class="dialog-actions">
        <v-btn variant="text" @click="banDialog.show = false">取消</v-btn>
        <v-btn color="error" variant="flat" @click="handleBan">确认封禁</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import UsersPanel from './UsersPanel.vue'
import UserAvatar from '../../components/UserAvatar.vue'
import { adminUserApi } from '../../api/admin'
import api from '../../api'
import { confirm, success, error } from '../../utils/modal'

const router = useRouter()
const users = ref([])
const loading = ref(true)
const currentUserId = ref(null)
const currentUserRole = ref(null)

const editRoleDialog = ref({
  show: false,
  user: null,
  role: ''
})

const editUserDialog = ref({
  show: false,
  user: null,
  displayName: ''
})

const banDialog = ref({
  show: false,
  user: null,
  reason: ''
})

const loadUsers = async () => {
  loading.value = true
  try {
    const response = await adminUserApi.getUsers()
    users.value = response.data.users || []
  } catch (error) {
    console.error('加载用户列表失败', error)
  } finally {
    loading.value = false
  }
}

const loadCurrentUser = async () => {
  try {
    const response = await api.get('/profile')
    currentUserId.value = response.data.id
    currentUserRole.value = response.data.role
  } catch (error) {
    console.error('加载当前用户失败', error)
    router.push('/login')
  }
}

const showEditRoleDialog = (user) => {
  editRoleDialog.value = {
    show: true,
    user,
    role: user.role
  }
}

const showEditUserDialog = (user) => {
  editUserDialog.value = {
    show: true,
    user,
    displayName: user.display_name
  }
}

const showBanDialog = (user) => {
  banDialog.value = {
    show: true,
    user,
    reason: ''
  }
}

const handleEditRole = async () => {
  try {
    await adminUserApi.updateUserRole(editRoleDialog.value.user.id, editRoleDialog.value.role)
    success('修改成功')
    editRoleDialog.value.show = false
    loadUsers()
  } catch (error) {
    console.error('修改角色失败', error)
    error(error.response?.data?.error || '修改失败')
  }
}

const handleEditUser = async () => {
  if (!editUserDialog.value.displayName) {
    error('请输入显示名称')
    return
  }
  try {
    await adminUserApi.updateUser(editUserDialog.value.user.id, {
      display_name: editUserDialog.value.displayName
    })
    success('保存成功')
    editUserDialog.value.show = false
    loadUsers()
  } catch (error) {
    console.error('编辑用户失败', error)
    error(error.response?.data?.error || '保存失败')
  }
}

const handleBan = async () => {
  try {
    await adminUserApi.banUser(banDialog.value.user.id, banDialog.value.reason)
    success('封禁成功')
    banDialog.value.show = false
    loadUsers()
  } catch (error) {
    console.error('封禁用户失败', error)
    error(error.response?.data?.error || '封禁失败')
  }
}

const handleUnban = async (userId) => {
  const confirmed = await confirm('确定要解封此用户吗？')
  if (!confirmed) return
  try {
    await adminUserApi.unbanUser(userId)
    success('解封成功')
    loadUsers()
  } catch (error) {
    console.error('解封用户失败', error)
    error(error.response?.data?.error || '解封失败')
  }
}

const handleDeleteUser = async (user) => {
  const confirmed = await confirm(`确定要删除用户 "${user.display_name}" 吗？`)
  if (!confirmed) return
  try {
    await adminUserApi.deleteUser(user.id)
    success('删除成功')
    loadUsers()
  } catch (error) {
    console.error('删除用户失败', error)
    error(error.response?.data?.error || '删除失败')
  }
}

onMounted(() => {
  loadCurrentUser()
  loadUsers()
})
</script>

<style scoped>
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
</style>