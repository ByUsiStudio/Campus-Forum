<template>
  <div class="users-panel">
    <div class="panel-header">
      <h2 class="panel-title">用户管理</h2>
      <v-btn variant="text" color="primary" @click="$emit('refresh')" :loading="loading">
        <v-icon start size="18">mdi-refresh</v-icon>
        刷新
      </v-btn>
    </div>

    <v-card class="table-card">
      <v-table class="users-table">
        <thead>
          <tr>
            <th class="text-left">用户</th>
            <th class="text-left">QQ号码</th>
            <th class="text-left">角色</th>
            <th class="text-left">状态</th>
            <th class="text-left">注册时间</th>
            <th class="text-center">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="users.length === 0">
            <td colspan="6" class="text-center pa-8">
              <v-icon size="40" color="grey">mdi-account-search</v-icon>
              <div class="mt-2 text-grey text-body-2">暂无用户数据</div>
            </td>
          </tr>
          <tr v-for="user in users" :key="user.id" class="user-row">
            <td>
              <div class="user-cell">
                <UserAvatar :user="user" :size="32" />
                <div class="user-info">
                  <div class="user-name">{{ user.display_name }}</div>
                  <div class="user-id">ID: {{ user.id }}</div>
                </div>
              </div>
            </td>
            <td class="text-body-2">{{ user.qq_number || '-' }}</td>
            <td>
              <v-chip size="small" :color="getRoleColor(user.role)" variant="tonal">
                {{ getRoleText(user.role) }}
              </v-chip>
            </td>
            <td>
              <v-chip
                size="small"
                :color="user.status === 'banned' ? 'error' : 'success'"
                variant="tonal"
              >
                {{ user.status === 'banned' ? '已封禁' : '正常' }}
              </v-chip>
            </td>
            <td class="text-body-2">{{ formatDate(user.created_at) }}</td>
            <td>
              <div class="action-cell">
                <v-btn variant="text" size="small" color="primary" @click="$emit('edit-role', user)" v-if="canEditRole(user)">
                  改角色
                </v-btn>
                <v-btn variant="text" size="small" color="warning" @click="$emit('ban', user)" v-if="canBanUser(user)">
                  {{ user.status === 'banned' ? '解封' : '封禁' }}
                </v-btn>
                <v-btn variant="text" size="small" color="error" @click="$emit('delete', user)" v-if="canDeleteUser(user)">
                  删除
                </v-btn>
              </div>
            </td>
          </tr>
        </tbody>
      </v-table>
    </v-card>
  </div>
</template>

<script>
import UserAvatar from '../../components/UserAvatar.vue'

export default {
  name: 'UsersPanel',
  components: { UserAvatar },
  props: {
    users: { type: Array, default: () => [] },
    loading: { type: Boolean, default: false },
    currentUserId: { type: [Number, String], default: null },
    currentUserRole: { type: String, default: null }
  },
  emits: ['edit-role', 'edit-user', 'ban', 'unban', 'delete', 'refresh'],
  setup(props) {
    const getRoleColor = (role) => {
      const colors = { admin: 'error', system: 'warning', user: 'default' }
      return colors[role] || 'default'
    }

    const getRoleText = (role) => {
      const texts = { admin: '管理员', system: '系统管理员', user: '用户' }
      return texts[role] || '用户'
    }

    const formatDate = (dateString) => {
      if (!dateString) return '-'
      const date = new Date(dateString)
      return date.toLocaleString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' })
    }

    const canEditRole = (user) => {
      if (!props.currentUserId || user.id === props.currentUserId) return false
      if (props.currentUserRole === 'system') return true
      if (props.currentUserRole === 'admin' && user.role !== 'system' && user.role !== 'System') return true
      return false
    }

    const canBanUser = (user) => {
      if (!props.currentUserId || user.id === props.currentUserId) return false
      if (props.currentUserRole === 'system') return true
      if (props.currentUserRole === 'admin' && user.role !== 'system' && user.role !== 'System') return true
      return false
    }

    const canDeleteUser = (user) => {
      if (!props.currentUserId || user.id === props.currentUserId) return false
      if (props.currentUserRole === 'system') return true
      if (props.currentUserRole === 'admin' && user.role !== 'system' && user.role !== 'System') return true
      return false
    }

    return { getRoleColor, getRoleText, formatDate, canEditRole, canBanUser, canDeleteUser }
  }
}
</script>

<style scoped>
.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.panel-title {
  font-size: 20px;
  font-weight: 600;
  color: #1A1A1A;
  margin: 0;
}

.table-card {
  border-radius: 12px;
  border: 1px solid #F0F0F0;
  overflow: hidden;
}

.users-table th {
  font-weight: 600;
  font-size: 12px;
  color: #666;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  padding: 14px 16px !important;
  background: #FAFAFA;
  border-bottom: 1px solid #F0F0F0;
}

.users-table td {
  padding: 12px 16px !important;
  font-size: 14px;
  color: #333;
  border-bottom: 1px solid #F5F5F5;
  vertical-align: middle;
}

.user-row:hover td {
  background: #FAFAFA;
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-info {
  display: flex;
  flex-direction: column;
}

.user-name {
  font-weight: 500;
  font-size: 14px;
  color: #1A1A1A;
}

.user-id {
  font-size: 12px;
  color: #999;
}

.action-cell {
  display: flex;
  gap: 4px;
  justify-content: center;
}
</style>