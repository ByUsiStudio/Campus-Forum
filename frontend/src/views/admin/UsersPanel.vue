<template>
  <div class="users-panel">
    <!-- 面板头部 -->
    <div class="panel-header">
      <div class="header-left">
        <h2 class="panel-title">用户管理</h2>
        <v-chip size="small" variant="tonal" color="primary" class="ml-3">
          {{ users.length }} 位用户
        </v-chip>
      </div>
      <div class="header-right">
        <v-text-field
          v-model="searchQuery"
          placeholder="搜索用户..."
          prepend-inner-icon="mdi-magnify"
          variant="outlined"
          density="compact"
          hide-details
          class="search-field"
          clearable
        />
        <v-btn variant="text" color="primary" @click="$emit('refresh')" :loading="loading">
          <v-icon start size="18">mdi-refresh</v-icon>
          刷新
        </v-btn>
      </div>
    </div>

    <!-- 用户表格 -->
    <v-card class="table-card" elevation="1">
      <v-table class="users-table" hover>
        <thead>
          <tr>
            <th class="text-left user-col">用户信息</th>
            <th class="text-left contact-col">联系方式</th>
            <th class="text-center role-col">角色</th>
            <th class="text-center status-col">状态</th>
            <th class="text-left time-col">注册时间</th>
            <th class="text-center action-col">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="filteredUsers.length === 0">
            <td colspan="6" class="empty-state">
              <div class="empty-content">
                <v-icon size="48" color="grey-lighten-1">mdi-account-search-outline</v-icon>
                <div class="mt-3 text-body-1 text-medium-emphasis">
                  {{ searchQuery ? '未找到匹配的用户' : '暂无用户数据' }}
                </div>
              </div>
            </td>
          </tr>
          <tr v-for="user in filteredUsers" :key="user.id" class="user-row">
            <td class="user-cell">
              <div class="user-info-wrapper">
                <UserAvatar :user="user" :size="40" />
                <div class="user-details">
                  <div class="user-name">{{ user.display_name }}</div>
                  <div class="user-id">ID: {{ user.id }}</div>
                </div>
              </div>
            </td>
            <td class="contact-cell">
              <div class="contact-info">
                <div v-if="user.qq_number" class="contact-item">
                  <v-icon size="16" class="mr-1">mdi-qqchat</v-icon>
                  <span>{{ user.qq_number }}</span>
                </div>
                <div v-else class="text-caption text-medium-emphasis">-</div>
              </div>
            </td>
            <td class="role-cell">
              <v-chip 
                size="small" 
                :color="getRoleColor(user.role)" 
                variant="tonal"
                class="role-chip"
              >
                {{ getRoleText(user.role) }}
              </v-chip>
            </td>
            <td class="status-cell">
              <v-chip
                size="small"
                :color="user.status === 'banned' ? 'error' : 'success'"
                variant="tonal"
                class="status-chip"
              >
                <v-icon start size="14">
                  {{ user.status === 'banned' ? 'mdi-block-helper' : 'mdi-check-circle' }}
                </v-icon>
                {{ user.status === 'banned' ? '已封禁' : '正常' }}
              </v-chip>
            </td>
            <td class="time-cell">
              <div class="time-info">
                <div class="time-date">{{ formatDate(user.created_at) }}</div>
                <div class="time-relative text-caption text-medium-emphasis">
                  {{ formatRelativeTime(user.created_at) }}
                </div>
              </div>
            </td>
            <td class="action-cell">
              <div class="action-buttons">
                <v-tooltip text="编辑角色" location="top">
                  <template #activator="{ props }">
                    <v-btn
                      v-bind="props"
                      variant="text"
                      size="small"
                      color="primary"
                      @click="$emit('edit-role', user)"
                      v-if="canEditRole(user)"
                      icon
                    >
                      <v-icon size="18">mdi-account-edit</v-icon>
                    </v-btn>
                  </template>
                </v-tooltip>
                
                <v-tooltip :text="user.status === 'banned' ? '解封用户' : '封禁用户'" location="top">
                  <template #activator="{ props }">
                    <v-btn
                      v-bind="props"
                      variant="text"
                      size="small"
                      :color="user.status === 'banned' ? 'success' : 'warning'"
                      @click="$emit(user.status === 'banned' ? 'unban' : 'ban', user)"
                      v-if="canBanUser(user)"
                      icon
                    >
                      <v-icon size="18">
                        {{ user.status === 'banned' ? 'mdi-unlock' : 'mdi-lock' }}
                      </v-icon>
                    </v-btn>
                  </template>
                </v-tooltip>
                
                <v-tooltip text="删除用户" location="top">
                  <template #activator="{ props }">
                    <v-btn
                      v-bind="props"
                      variant="text"
                      size="small"
                      color="error"
                      @click="$emit('delete', user)"
                      v-if="canDeleteUser(user)"
                      icon
                    >
                      <v-icon size="18">mdi-delete</v-icon>
                    </v-btn>
                  </template>
                </v-tooltip>
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
  data() {
    return {
      searchQuery: ''
    }
  },
  computed: {
    filteredUsers() {
      if (!this.searchQuery) return this.users
      const query = this.searchQuery.toLowerCase()
      return this.users.filter(user => 
        user.display_name.toLowerCase().includes(query) ||
        user.qq_number?.toLowerCase().includes(query) ||
        user.id.toString().includes(query)
      )
    }
  },
  methods: {
    getRoleColor(role) {
      const colors = { 
        admin: 'error', 
        system: 'warning', 
        user: 'primary' 
      }
      return colors[role] || 'grey'
    },
    getRoleText(role) {
      const texts = { 
        admin: '管理员', 
        system: '系统管理员', 
        user: '普通用户' 
      }
      return texts[role] || '用户'
    },
    formatDate(dateString) {
      if (!dateString) return '-'
      const date = new Date(dateString)
      return date.toLocaleDateString('zh-CN', { 
        year: 'numeric', 
        month: '2-digit', 
        day: '2-digit' 
      })
    },
    formatRelativeTime(dateString) {
      if (!dateString) return ''
      const date = new Date(dateString)
      const now = new Date()
      const diff = now - date
      const days = Math.floor(diff / (1000 * 60 * 60 * 24))
      
      if (days === 0) return '今天'
      if (days === 1) return '昨天'
      if (days < 7) return `${days}天前`
      if (days < 30) return `${Math.floor(days / 7)}周前`
      if (days < 365) return `${Math.floor(days / 30)}个月前`
      return `${Math.floor(days / 365)}年前`
    },
    canEditRole(user) {
      if (!this.currentUserId || user.id === this.currentUserId) return false
      if (this.currentUserRole === 'system') return true
      if (this.currentUserRole === 'admin' && user.role !== 'system' && user.role !== 'System') return true
      return false
    },
    canBanUser(user) {
      if (!this.currentUserId || user.id === this.currentUserId) return false
      if (this.currentUserRole === 'system') return true
      if (this.currentUserRole === 'admin' && user.role !== 'system' && user.role !== 'System') return true
      return false
    },
    canDeleteUser(user) {
      if (!this.currentUserId || user.id === this.currentUserId) return false
      if (this.currentUserRole === 'system') return true
      if (this.currentUserRole === 'admin' && user.role !== 'system' && user.role !== 'System') return true
      return false
    }
  }
}
</script>

<style scoped>
.users-panel {
  width: 100%;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  gap: 16px;
  flex-wrap: wrap;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.panel-title {
  font-size: 20px;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.search-field {
  width: 240px;
}

.table-card {
  border-radius: 16px;
  border: 1px solid #e8eaed;
  overflow: hidden;
  background: linear-gradient(135deg, #ffffff 0%, #fafbfc 100%);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  transition: all 0.3s ease;
}

.table-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.users-table {
  width: 100%;
}

.users-table th {
  font-weight: 600;
  font-size: 12px;
  color: #666;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  padding: 16px !important;
  background: linear-gradient(135deg, #f8f9fa 0%, #fafbfc 100%);
  border-bottom: 2px solid #e8eaed;
  white-space: nowrap;
}

.users-table td {
  padding: 16px !important;
  font-size: 14px;
  color: #333;
  border-bottom: 1px solid #f0f0f0;
  vertical-align: middle;
}

.user-row:hover td {
  background: linear-gradient(135deg, rgba(103, 80, 164, 0.04) 0%, rgba(156, 39, 176, 0.04) 100%);
  transition: background 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.user-col { width: 25%; }
.contact-col { width: 15%; }
.role-col { width: 12%; }
.status-col { width: 12%; }
.time-col { width: 16%; }
.action-col { width: 20%; }

.user-cell {
  padding: 12px 16px !important;
}

.user-info-wrapper {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-details {
  display: flex;
  flex-direction: column;
}

.user-name {
  font-weight: 600;
  font-size: 14px;
  color: #1a1a1a;
  margin-bottom: 2px;
}

.user-id {
  font-size: 12px;
  color: #999;
}

.contact-cell {
  padding: 12px 16px !important;
}

.contact-info {
  display: flex;
  align-items: center;
}

.contact-item {
  display: flex;
  align-items: center;
  font-size: 13px;
  color: #666;
}

.role-cell, .status-cell {
  text-align: center;
  padding: 12px 16px !important;
}

.role-chip, .status-chip {
  font-weight: 500;
}

.time-cell {
  padding: 12px 16px !important;
}

.time-info {
  display: flex;
  flex-direction: column;
}

.time-date {
  font-size: 13px;
  color: #333;
  margin-bottom: 2px;
}

.time-relative {
  font-size: 11px;
}

.action-cell {
  text-align: center;
  padding: 12px 16px !important;
}

.action-buttons {
  display: flex;
  justify-content: center;
  gap: 4px;
}

.empty-state {
  padding: 48px 16px !important;
  text-align: center;
}

.empty-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

@media (max-width: 768px) {
  .panel-header {
    flex-direction: column;
    align-items: stretch;
  }
  
  .header-left, .header-right {
    justify-content: space-between;
  }
  
  .search-field {
    width: 100%;
  }
  
  .users-table {
    font-size: 13px;
  }
  
  .users-table th,
  .users-table td {
    padding: 12px 8px !important;
  }
  
  .user-col { width: 30%; }
  .contact-col { display: none; }
  .role-col { width: 15%; }
  .status-col { width: 15%; }
  .time-col { display: none; }
  .action-col { width: 25%; }
}

@media (max-width: 480px) {
  .action-buttons {
    flex-direction: column;
    gap: 2px;
  }
  
  .user-info-wrapper {
    gap: 8px;
  }
  
  .user-name {
    font-size: 13px;
  }
  
  .user-id {
    font-size: 11px;
  }
}
</style>