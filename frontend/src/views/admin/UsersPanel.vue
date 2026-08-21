<template>
  <div>
    <!-- 搜索栏 -->
    <v-card class="mb-4 pa-3" variant="flat" rounded="lg">
      <v-row dense align="center">
        <v-col cols="12" sm="8">
          <v-text-field
            v-model="searchQuery"
            placeholder="搜索用户..."
            prepend-inner-icon="mdi-magnify"
            variant="outlined"
            density="compact"
            hide-details
            clearable
          />
        </v-col>
        <v-col cols="12" sm="4" class="text-end">
          <v-btn variant="tonal" color="primary" @click="$emit('refresh')" :loading="loading" prepend-icon="mdi-refresh">
            刷新
          </v-btn>
        </v-col>
      </v-row>
    </v-card>

    <!-- 用户列表 -->
    <v-card variant="flat" rounded="lg">
      <v-list lines="two" v-if="filteredUsers.length > 0">
        <v-list-item v-for="user in filteredUsers" :key="user.id" class="py-3">
          <template v-slot:prepend>
            <UserAvatar :user="user" :size="48" />
          </template>

          <v-list-item-title class="font-weight-medium">
            {{ user.display_name }}
            <v-chip size="x-small" :color="getRoleColor(user.role)" variant="tonal" class="ml-2">
              {{ getRoleText(user.role) }}
            </v-chip>
            <v-chip size="x-small" :color="user.status === 'banned' ? 'error' : 'success'" variant="tonal" class="ml-1">
              {{ user.status === 'banned' ? '已封禁' : '正常' }}
            </v-chip>
          </v-list-item-title>

          <v-list-item-subtitle>
            ID: {{ user.id }} | QQ: {{ user.qq_number || '-' }} | {{ formatDate(user.created_at) }}
          </v-list-item-subtitle>

          <template v-slot:append>
            <v-btn-group variant="text" density="compact" divided>
              <v-btn size="small" color="primary" @click="$emit('edit-role', user)" v-if="canEditRole(user)">
                <v-icon>mdi-account-edit</v-icon>
                <v-tooltip activator="parent">修改角色</v-tooltip>
              </v-btn>
              <v-btn size="small" :color="user.status === 'banned' ? 'success' : 'warning'" 
                     @click="$emit(user.status === 'banned' ? 'unban' : 'ban', user)" v-if="canBanUser(user)">
                <v-icon>{{ user.status === 'banned' ? 'mdi-lock-open' : 'mdi-lock' }}</v-icon>
                <v-tooltip activator="parent">{{ user.status === 'banned' ? '解封' : '封禁' }}</v-tooltip>
              </v-btn>
              <v-btn size="small" color="error" @click="$emit('delete', user)" v-if="canDeleteUser(user)">
                <v-icon>mdi-delete</v-icon>
                <v-tooltip activator="parent">删除</v-tooltip>
              </v-btn>
            </v-btn-group>
          </template>
        </v-list-item>
      </v-list>

      <v-card-text v-else class="text-center py-8">
        <v-icon size="48" color="grey-lighten-1">mdi-account-search</v-icon>
        <div class="text-body-1 text-medium-emphasis mt-2">
          {{ searchQuery ? '未找到匹配的用户' : '暂无用户数据' }}
        </div>
      </v-card-text>
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
    getRoleColor(role) {
      return { admin: 'error', system: 'warning', user: 'primary' }[role] || 'grey'
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