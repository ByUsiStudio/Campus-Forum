<template>
  <div>
    <!-- 搜索栏 -->
    <el-card class="page-container mb-4" shadow="never">
      <el-row :gutter="12" align="middle">
        <el-col :xs="24" :sm="18">
          <el-input
            v-model="searchQuery"
            placeholder="搜索用户..."
            :prefix-icon="Search"
            clearable
          />
        </el-col>
        <el-col :xs="24" :sm="6" class="search-actions">
          <el-button
            type="primary"
            plain
            :loading="loading"
            :icon="Refresh"
            @click="$emit('refresh')"
          >
            刷新
          </el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- 用户列表 -->
    <el-card class="page-container" shadow="never">
      <el-table :data="filteredUsers" v-loading="loading" v-if="filteredUsers.length > 0">
        <el-table-column label="用户" min-width="220">
          <template #default="{ row }">
            <div class="user-cell">
              <UserAvatar :user="row" :size="48" />
            </div>
          </template>
        </el-table-column>

        <el-table-column label="角色" width="130">
          <template #default="{ row }">
            <el-tag :type="getRoleType(row.role)" effect="light" size="small">
              {{ getRoleText(row.role) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag
              :type="row.status === 'banned' ? 'danger' : 'success'"
              effect="light"
              size="small"
            >
              {{ row.status === 'banned' ? '已封禁' : '正常' }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="ID / 联系 / 注册时间" min-width="240">
          <template #default="{ row }">
            <div class="user-meta">
              <span>ID: {{ row.id }}</span>
              <span>QQ: {{ row.qq_number || '-' }}</span>
              <span>{{ formatDate(row.created_at) }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <div class="action-btns">
              <el-tooltip content="修改角色" placement="top">
                <el-button
                  size="small"
                  text
                  type="primary"
                  :icon="Edit"
                  v-if="canEditRole(row)"
                  @click="$emit('edit-role', row)"
                />
              </el-tooltip>

              <el-tooltip :content="row.status === 'banned' ? '解封' : '封禁'" placement="top">
                <el-button
                  size="small"
                  text
                  :type="row.status === 'banned' ? 'success' : 'warning'"
                  :icon="row.status === 'banned' ? Unlock : Lock"
                  v-if="canBanUser(row)"
                  @click="$emit(row.status === 'banned' ? 'unban' : 'ban', row)"
                />
              </el-tooltip>

              <el-tooltip content="删除" placement="top">
                <el-button
                  size="small"
                  text
                  type="danger"
                  :icon="Delete"
                  v-if="canDeleteUser(row)"
                  @click="$emit('delete', row)"
                />
              </el-tooltip>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <el-card
        v-else
        shadow="never"
        class="empty-card"
        :body-style="{ textAlign: 'center', padding: '32px' }"
      >
        <el-icon :size="48" color="var(--el-color-info-light-5)"><User /></el-icon>
        <div class="empty-text mt-2">
          {{ searchQuery ? '未找到匹配的用户' : '暂无用户数据' }}
        </div>
      </el-card>
    </el-card>
  </div>
</template>

<script>
import UserAvatar from '../../components/UserAvatar.vue'
import {
  Search,
  Refresh,
  Edit,
  Lock,
  Unlock,
  Delete,
  User
} from '@element-plus/icons-vue'

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
    getRoleType(role) {
      return { admin: 'danger', system: 'warning', user: 'primary' }[role] || 'info'
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
.search-actions {
  display: flex;
  justify-content: flex-end;
}

.user-cell {
  padding: 4px 0;
}

.user-meta {
  display: flex;
  flex-direction: column;
  gap: 2px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.action-btns {
  display: flex;
  align-items: center;
  gap: 2px;
}

.empty-text {
  color: var(--el-text-color-secondary);
}

.mt-2 {
  margin-top: 8px;
}

.mb-4 {
  margin-bottom: 16px;
}
</style>
