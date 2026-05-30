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
  animation: panelFadeIn 0.4s ease;
}

@keyframes panelFadeIn {
  from {
    opacity: 0;
    transform: translateY(15px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.table-card {
  border-radius: 20px;
  overflow: hidden;
  background: #fff;
  box-shadow: 0 4px 20px -4px rgba(0, 0, 0, 0.06);
}

.users-table {
  width: 100%;
  border-collapse: collapse;
}

.users-table thead {
  background: linear-gradient(135deg, #F8F7FF 0%, #F1F5F9 100%);
}

.users-table th {
  font-weight: 600;
  color: #49454F;
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 0.8px;
  padding: 18px 20px !important;
  text-align: left;
  border-bottom: 2px solid #E7E0EC;
}

.users-table td {
  padding: 16px 20px !important;
  font-size: 14px;
  color: #49454F;
  border-bottom: 1px solid #F2F0F4;
  transition: background 0.2s ease;
}

.user-row {
  transition: all 0.25s ease;
}

.user-row:hover {
  background: #F8F7FF;
  transform: scale(1.002);
}

.user-row:hover td {
  background: #F8F7FF;
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 14px;
}

.user-info {
  min-width: 0;
  display: flex;
  flex-direction: column;
}

.user-name {
  font-weight: 600;
  color: #1C1B1F;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  font-size: 14px;
}

.user-id {
  font-size: 12px;
  color: #938F99;
  margin-top: 2px;
}

.qq-cell {
  color: #625B71;
  font-family: 'Monaco', 'Consolas', monospace;
  font-size: 13px;
}

.date-cell {
  color: #938F99;
  font-size: 13px;
}

.action-cell {
  display: flex;
  gap: 8px;
  justify-content: center;
  flex-wrap: wrap;
}

.action-cell .v-btn {
  border-radius: 8px;
  font-size: 12px;
  padding: 6px 12px;
  transition: all 0.2s ease;
}

.action-cell .v-btn:hover {
  transform: translateY(-1px);
}

@media (max-width: 960px) {
  .users-table {
    font-size: 13px;
  }

  .users-table th,
  .users-table td {
    padding: 14px 12px !important;
  }

  .qq-cell,
  .date-cell {
    display: none;
  }
}

@media (max-width: 600px) {
  .action-cell {
    flex-direction: column;
    gap: 4px;
  }

  .action-cell .v-btn {
    width: 100%;
  }
}
</style>
