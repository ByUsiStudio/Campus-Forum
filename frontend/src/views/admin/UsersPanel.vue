<template>
  <div>
    <div class="layui-card mb-4">
      <div class="layui-card-body d-flex flex-wrap gap-3">
        <div class="search-box flex-1 min-w-[200px]">
          <div class="layui-input-group">
            <div class="layui-input-group-prepend">
              <span class="layui-input-group-text"><i class="fa-solid fa-magnifying-glass"></i></span>
            </div>
            <input type="text" v-model="searchQuery" placeholder="搜索用户..." class="layui-input" />
          </div>
        </div>
        <button class="layui-btn" @click="$emit('refresh')" :disabled="loading">
          <i class="fa-solid fa-rotate-right"></i>
          刷新
        </button>
      </div>
    </div>

    <div class="layui-card">
      <div v-if="filteredUsers.length > 0" class="users-list">
        <div v-for="user in filteredUsers" :key="user.id" class="user-item">
          <div class="user-avatar-wrap">
            <UserAvatar :user="user" :size="48" />
          </div>

          <div class="user-info">
            <div class="user-name-wrap">
              <span class="user-name font-weight-medium">{{ user.display_name }}</span>
              <span :class="['role-tag', getRoleClass(user.role)]">{{ getRoleText(user.role) }}</span>
              <span :class="['status-tag', user.status === 'banned' ? 'error' : 'success']">
                {{ user.status === 'banned' ? '已封禁' : '正常' }}
              </span>
            </div>
            <div class="user-meta">ID: {{ user.id }} | QQ: {{ user.qq_number || '-' }} | {{ formatDate(user.created_at) }}</div>
          </div>

          <div class="user-actions">
            <button v-if="canEditRole(user)" class="action-btn" @click="$emit('edit-role', user)" title="修改角色">
              <i class="fa-solid fa-user-pen"></i>
            </button>
            <button v-if="canBanUser(user)" :class="['action-btn', user.status === 'banned' ? 'success' : 'warning']" 
                    @click="$emit(user.status === 'banned' ? 'unban' : 'ban', user)" 
                    :title="user.status === 'banned' ? '解封' : '封禁'">
              <i :class="user.status === 'banned' ? 'fa-solid fa-lock-open' : 'fa-solid fa-lock'"></i>
            </button>
            <button v-if="canDeleteUser(user)" class="action-btn danger" @click="$emit('delete', user)" title="删除">
              <i class="fa-solid fa-trash"></i>
            </button>
          </div>
        </div>
      </div>

      <div v-else class="empty-state">
        <i class="fa-solid fa-user-search" style="font-size: 48px; color: #ccc;"></i>
        <div class="text-body-1 text-medium-emphasis mt-2">
          {{ searchQuery ? '未找到匹配的用户' : '暂无用户数据' }}
        </div>
      </div>
    </div>
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
    return { searchQuery: '' }
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
    getRoleClass(role) {
      return { admin: 'admin', system: 'system', user: 'user' }[role] || 'user'
    },
    getRoleText(role) {
      return { admin: '管理员', system: '系统管理员', user: '普通用户' }[role] || '用户'
    },
    formatDate(dateString) {
      if (!dateString) return '-'
      return new Date(dateString).toLocaleDateString('zh-CN')
    },
    canEditRole(user) {
      if (!this.currentUserId || user.id === this.currentUserId) return false
      if (this.currentUserRole === 'system') return true
      if (this.currentUserRole === 'admin' && user.role !== 'system') return true
      return false
    },
    canBanUser(user) {
      if (!this.currentUserId || user.id === this.currentUserId) return false
      if (this.currentUserRole === 'system') return true
      if (this.currentUserRole === 'admin' && user.role !== 'system') return true
      return false
    },
    canDeleteUser(user) {
      if (!this.currentUserId || user.id === this.currentUserId) return false
      if (this.currentUserRole === 'system') return true
      if (this.currentUserRole === 'admin' && user.role !== 'system') return true
      return false
    }
  }
}
</script>

<style scoped>
.d-flex {
  display: flex;
}

.flex-wrap {
  flex-wrap: wrap;
}

.gap-3 {
  gap: 12px;
}

.flex-1 {
  flex: 1;
}

.min-w-\[200px\] {
  min-width: 200px;
}

.mb-4 {
  margin-bottom: 16px;
}

.users-list {
  padding: 0;
}

.user-item {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid #f0f0f0;
  gap: 12px;
}

.user-item:last-child {
  border-bottom: none;
}

.user-avatar-wrap {
  flex-shrink: 0;
}

.user-info {
  flex: 1;
  min-width: 0;
}

.user-name-wrap {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.user-name {
  font-weight: 500;
}

.role-tag {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.role-tag.admin {
  background: #fff2f0;
  color: #ff4d4f;
}

.role-tag.system {
  background: #fffbe6;
  color: #faad14;
}

.role-tag.user {
  background: #e6f7ff;
  color: #1890ff;
}

.status-tag {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.status-tag.success {
  background: #f6ffed;
  color: #52c41a;
}

.status-tag.error {
  background: #fff2f0;
  color: #ff4d4f;
}

.user-meta {
  font-size: 13px;
  color: #999;
  margin-top: 4px;
}

.user-actions {
  display: flex;
  gap: 8px;
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

.action-btn.warning:hover {
  color: #FAAD14;
}

.action-btn.success:hover {
  color: #52C41A;
}

.action-btn.danger:hover {
  background: #fff2f0;
  color: #FF5722;
}

.empty-state {
  text-align: center;
  padding: 32px;
}

.text-body-1 {
  font-size: 16px;
}

.text-medium-emphasis {
  color: #999;
}

.font-weight-medium {
  font-weight: 500;
}

.mt-2 {
  margin-top: 8px;
}
</style>