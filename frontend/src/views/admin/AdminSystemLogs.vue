<template>
  <div class="admin-system-logs">
    <!-- 页面标题 -->
    <div class="page-header mb-6">
      <div class="d-flex align-center justify-space-between flex-wrap ga-4">
        <div>
          <h1 class="text-h5 font-weight-bold mb-1">系统操作日志</h1>
          <p class="text-body-2 text-medium-emphasis">查看用户和管理员的所有操作记录</p>
        </div>
        <div class="d-flex ga-3">
          <v-btn
            color="error"
            variant="tonal"
            @click="deleteOldLogs"
            :loading="deleting"
            :disabled="deleting"
          >
            <v-icon start>mdi-delete-sweep</v-icon>
            删除90天前日志
          </v-btn>
          <v-btn
            color="primary"
            @click="refreshData"
            :loading="loading"
          >
            <v-icon start>mdi-refresh</v-icon>
            刷新
          </v-btn>
        </div>
      </div>
    </div>

    <!-- 筛选器 -->
    <v-card class="mb-6">
      <v-card-text>
        <v-row>
          <v-col cols="12" sm="6" md="4" lg="3">
            <v-select
              v-model="filter.module"
              :items="modules"
              label="筛选模块"
              variant="outlined"
              clearable
              hint="按模块筛选"
            >
              <template v-slot:prepend-item>
                <v-list-item
                  @click="filter.module = ''"
                  :title="'全部模块'"
                ></v-list-item>
              </template>
            </v-select>
          </v-col>
          <v-col cols="12" sm="6" md="4" lg="3">
            <v-text-field
              v-model="filter.user_id"
              label="用户ID筛选"
              type="number"
              variant="outlined"
              clearable
            ></v-text-field>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

    <!-- 日志列表 -->
    <v-card>
      <v-progress-linear
        v-if="loading"
        indeterminate
        color="primary"
      ></v-progress-linear>

      <v-data-table
        :headers="headers"
        :items="logs"
        :loading="loading"
        :items-per-page="20"
        :options.sync="pagination"
        :server-items-length="total"
        class="logs-table"
      >
        <template v-slot:item.user="{ item }">
          <div v-if="item.user">
            <div class="font-weight-bold text-primary">{{ item.user.username || item.user.nickname || '用户' }}</div>
            <div class="text-caption text-medium-emphasis">ID: {{ item.user_id }}</div>
          </div>
          <div v-else>
            <div class="text-caption text-medium-emphasis">用户ID: {{ item.user_id }}</div>
          </div>
        </template>

        <template v-slot:item.action="{ item }">
          <v-chip
            :color="getActionColor(item.action)"
            size="small"
          >
            {{ getActionText(item.action) }}
          </v-chip>
        </template>

        <template v-slot:item.module="{ item }">
          <v-chip
            color="blue-lighten-5"
            text-color="blue-darken-2"
            size="small"
          >
            {{ getModuleText(item.module) }}
          </v-chip>
        </template>

        <template v-slot:item.created_at="{ item }">
          {{ formatDate(item.created_at) }}
        </template>
      </v-data-table>
    </v-card>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import api from '../../api'

export default {
  name: 'AdminSystemLogs',
  setup() {
    const loading = ref(false)
    const deleting = ref(false)
    const logs = ref([])
    const modules = ref([])
    const total = ref(0)

    const pagination = ref({
      page: 1,
      itemsPerPage: 20,
      sortBy: ['created_at'],
      sortDesc: [true]
    })

    const filter = ref({
      module: '',
      user_id: ''
    })

    const headers = ref([
      { title: '时间', key: 'created_at', sortable: true },
      { title: '用户', key: 'user', sortable: false, width: 180 },
      { title: '操作', key: 'action', sortable: false, width: 120 },
      { title: '模块', key: 'module', sortable: false, width: 120 },
      { title: '详情', key: 'details', sortable: false },
      { title: 'IP地址', key: 'ip', sortable: false, width: 140 }
    ])

    const moduleMap = {
      'user': '用户',
      'article': '文章',
      'comment': '评论',
      'notification': '通知',
      'permission': '权限',
      'report': '举报',
      'announcement': '公告',
      'category': '分类',
      'title': '头衔',
      'site_config': '网站配置',
      'deletion': '删除申请'
    }

    const actionMap = {
      'create': '创建',
      'update': '更新',
      'delete': '删除',
      'login': '登录',
      'logout': '登出',
      'send': '发送',
      'grant': '授予',
      'revoke': '撤销',
      'ban': '封禁',
      'unban': '解封',
      'pin': '置顶',
      'unpin': '取消置顶',
      'approve': '批准',
      'reject': '拒绝'
    }

    const actionColorMap = {
      'create': 'green',
      'update': 'blue',
      'delete': 'red',
      'login': 'purple',
      'logout': 'grey',
      'send': 'cyan',
      'grant': 'green-darken-2',
      'revoke': 'orange',
      'ban': 'red-darken-2',
      'unban': 'green-lighten-2',
      'pin': 'amber',
      'unpin': 'grey',
      'approve': 'green-lighten-2',
      'reject': 'red-lighten-2'
    }

    const formatDate = (dateString) => {
      if (!dateString) return '-'
      const date = new Date(dateString)
      return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
      })
    }

    const getModuleText = (module) => {
      return moduleMap[module] || module
    }

    const getActionText = (action) => {
      return actionMap[action] || action
    }

    const getActionColor = (action) => {
      return actionColorMap[action] || 'grey'
    }

    const loadLogs = async () => {
      loading.value = true
      try {
        const params = {
          page: pagination.value.page,
          page_size: pagination.value.itemsPerPage
        }
        if (filter.value.module) {
          params.module = filter.value.module
        }
        if (filter.value.user_id) {
          params.user_id = filter.value.user_id
        }
        const response = await api.get('/api/system-logs', { params })
        logs.value = response.data.logs
        total.value = response.data.total
      } catch (error) {
        console.error('加载日志失败:', error)
      } finally {
        loading.value = false
      }
    }

    const loadModules = async () => {
      try {
        const response = await api.get('/api/system-logs/modules')
        modules.value = response.data.modules
      } catch (error) {
        console.error('加载模块列表失败:', error)
      }
    }

    const deleteOldLogs = async () => {
      if (!confirm('确定要删除90天前的所有日志吗？此操作不可恢复。')) {
        return
      }
      deleting.value = true
      try {
        await api.delete('/api/system-logs/old', { params: { days: 90 } })
        alert('旧日志删除成功')
        await loadLogs()
      } catch (error) {
        console.error('删除旧日志失败:', error)
        alert('删除旧日志失败')
      } finally {
        deleting.value = false
      }
    }

    const refreshData = async () => {
      await Promise.all([
        loadLogs(),
        loadModules()
      ])
    }

    onMounted(() => {
      refreshData()
    })

    return {
      loading,
      deleting,
      logs,
      modules,
      total,
      pagination,
      filter,
      headers,
      formatDate,
      getModuleText,
      getActionText,
      getActionColor,
      deleteOldLogs,
      refreshData
    }
  }
}
</script>

<style scoped>
.admin-system-logs {
  padding: 24px;
}

.page-header {
  margin-bottom: 24px;
}

.logs-table {
  width: 100%;
}

.text-ellipsis {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
