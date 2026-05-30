<template>
  <div class="users-panel">
    <div class="panel-header">
      <div class="header-left">
        <h2 class="panel-title">用户管理</h2>
        <p class="panel-subtitle">管理系统用户、角色与权限</p>
      </div>
      <v-btn variant="outlined" color="primary" @click="$emit('refresh')" :loading="loading">
        <v-icon start>mdi-refresh</v-icon>
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
              <v-icon size="48" color="grey">mdi-account-search</v-icon>
              <div class="mt-2 text-grey">暂无用户数据</div>
            </td>
          </tr>
          <tr v-for="user in users" :key="user.id" class="user-row">
            <td>
              <div class="user-cell">
                <UserAvatar :user="user" :size="36" />
                <div class="user-info">
                  <div class="user-name">{{ user.display_name }}</div>
                  <div class="user-id">ID: {{ user.id }}</div>
                </div>
              </div>
            </td>
            <td class="qq-cell">{{ user.qq_number || '-' }}</td>
            <td>
              <v-chip
                size="small"
                :color="getRoleColor(user.role)"
                variant="tonal"
              >
                <v-icon start size="x-small">{{ getRoleIcon(user.role) }}</v-icon>
                {{ getRoleText(user.role) }}
              </v-chip>
            </td>
            <td>
              <v-chip
                size="small"
                :color="user.status === 'banned' ? 'error' : 'success'"
                variant="tonal"
              >
                <v-icon start size="x-small">{{ user.status === 'banned' ? 'mdi-close-circle' : 'mdi-check-circle' }}</v-icon>
                {{ user.status === 'banned' ? '已封禁' : '正常' }}
              </v-chip>
            </td>
            <td class="date-cell">{{ formatDate(user.created_at) }}</td>
            <td>
              <div class="action-cell">
                <v-btn
                  variant="text"
                  size="small"
                  color="primary"
                  @click="$emit('edit-role', user)"
                  v-if="canEditRole(user)"
                >
                  改角色
                </v-btn>
                <v-btn
                  variant="text"
                  size="small"
                  color="info"
                  @click="$emit('edit-user', user)"
                  v-if="canEditUser(user)"
                >
                  编辑
                </v-btn>
                <v-btn
                  variant="text"
                  size="small"
                  color="warning"
                  @click="$emit('ban', user)"
                  v-if="canBanUser(user)"
                >
                  封禁
                </v-btn>
                <v-btn
                  variant="text"
                  size="small"
                  color="success"
                  @click="$emit('unban', user)"
                  v-if="canUnbanUser(user)"
                >
                  解封
                </v-btn>
                <v-btn
                  variant="text"
                  size="small"
                  color="error"
                  @click="$emit('delete', user)"
                  v-if="canDeleteUser(user)"
                >
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
  components: {
    UserAvatar
  },
  props: {
    users: {
      type: Array,
      default: () => []
    },
    loading: {
      type: Boolean,
      default: false
    },
    currentUserId: {
      type: [Number, String],
      default: null
    },
    currentUserRole: {
      type: String,
      default: null
    }
  },
  emits: ['edit-role', 'edit-user', 'ban', 'unban', 'delete', 'refresh'],
  setup(props) {
    const getRoleColor = (role) => {
      const colors = {
        admin: 'error',
        system: 'warning',
        user: 'default'
      }
      return colors[role] || 'default'
    }

    const getRoleIcon = (role) => {
      const icons = {
        admin: 'mdi-shield-account',
        system: 'mdi-crown',
        user: 'mdi-account'
      }
      return icons[role] || 'mdi-account'
    }

    const getRoleText = (role) => {
      const texts = {
        admin: '管理员',
        system: '系统管理员',
        user: '用户'
      }
      return texts[role] || '用户'
    }

    const formatDate = (dateString) => {
      if (!dateString) return '-'
      const date = new Date(dateString)
      return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit'
      })
    }

    const canEditRole = (user) => {
      if (!props.currentUserId || user.id === props.currentUserId) return false
      if (props.currentUserRole === 'system') return true
      if (props.currentUserRole === 'admin' && user.role !== 'system' && user.role !== 'System') return true
      return false
    }

    const canEditUser = (user) => {
      if (!props.currentUserId || user.id === props.currentUserId) return false
      if (props.currentUserRole === 'system') return true
      if (props.currentUserRole === 'admin' && user.role !== 'system' && user.role !== 'System') return true
      return false
    }

    const canBanUser = (user) => {
      if (!props.currentUserId || user.id === props.currentUserId) return false
      if (user.status === 'banned') return false
      if (props.currentUserRole === 'system') return true
      if (props.currentUserRole === 'admin' && user.role !== 'system' && user.role !== 'System') return true
      return false
    }

    const canUnbanUser = (user) => {
      if (!props.currentUserId || user.id === props.currentUserId) return false
      if (user.status !== 'banned') return false
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

    return {
      getRoleColor,
      getRoleIcon,
      getRoleText,
      formatDate,
      canEditRole,
      canEditUser,
      canBanUser,
      canUnbanUser,
      canDeleteUser
    }
  }
}
</script>

<style scoped>
.users-panel {
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 16px;
}

.panel-title {
  font-size: 1.5rem;
  font-weight: 700;
  color: #1a1a2e;
  margin: 0 0 4px 0;
}

.panel-subtitle {
  font-size: 0.9rem;
  color: #6b7280;
  margin: 0;
}

.table-card {
  border-radius: 16px;
  overflow: hidden;
}

.users-table {
  width: 100%;
}

.users-table th {
  font-weight: 600;
  color: #6b7280;
  font-size: 0.85rem;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  background: #f8f9ff !important;
  padding: 16px 12px !important;
}

.users-table td {
  padding: 14px 12px !important;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.user-row:hover {
  background: rgba(103, 80, 164, 0.03);
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-info {
  min-width: 0;
}

.user-name {
  font-weight: 600;
  color: #1a1a2e;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-id {
  font-size: 0.75rem;
  color: #9ca3af;
}

.qq-cell {
  color: #6b7280;
  font-family: monospace;
}

.date-cell {
  color: #6b7280;
  font-size: 0.9rem;
}

.action-cell {
  display: flex;
  gap: 4px;
  justify-content: center;
  flex-wrap: wrap;
}

@media (max-width: 960px) {
  .users-table {
    font-size: 0.85rem;
  }

  .users-table th,
  .users-table td {
    padding: 10px 8px !important;
  }

  .qq-cell,
  .date-cell {
    display: none;
  }
}

@media (max-width: 600px) {
  .action-cell {
    flex-direction: column;
  }
}
</style>
