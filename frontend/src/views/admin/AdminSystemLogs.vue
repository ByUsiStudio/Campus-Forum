<template>
  <div class="admin-system-logs page-container">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="header-row">
        <div>
          <h1 class="page-title">系统操作日志</h1>
          <p class="page-desc">查看用户和管理员的所有操作记录</p>
        </div>
        <div class="header-actions">
          <el-button
            type="danger"
            plain
            :loading="deleting"
            :disabled="deleting"
            @click="deleteOldLogs"
          >
            <el-icon><Delete /></el-icon>
            删除90天前日志
          </el-button>
          <el-button
            type="primary"
            :loading="loading"
            @click="refreshData"
          >
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </div>
    </div>

    <!-- 筛选器 -->
    <el-card class="mb-6" shadow="never">
      <div class="filter-row">
        <el-select
          v-model="filter.module"
          clearable
          placeholder="筛选模块"
          class="filter-module"
        >
          <el-option
            v-for="m in allModules"
            :key="m"
            :label="m"
            :value="m"
          />
        </el-select>
        <el-input
          v-model="filter.user_id"
          type="number"
          clearable
          placeholder="用户ID筛选"
          class="filter-user"
        />
      </div>
    </el-card>

    <!-- 日志列表 -->
    <el-card shadow="never">
      <el-table
        v-loading="loading"
        :data="logs"
        class="logs-table"
      >
        <el-table-column prop="created_at" label="时间" sortable>
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="用户" width="180">
          <template #default="{ row }">
            <div v-if="row.user">
              <div class="user-name">{{ row.user.username || row.user.nickname || '用户' }}</div>
              <div class="user-id">ID: {{ row.user_id }}</div>
            </div>
            <div v-else class="user-id">用户ID: {{ row.user_id }}</div>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120">
          <template #default="{ row }">
            <el-tag :type="getActionType(row.action)" size="small">
              {{ getActionText(row.action) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="模块" width="120">
          <template #default="{ row }">
            <el-tag type="info" size="small">
              {{ getModuleText(row.module) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="details" label="详情" />
        <el-table-column prop="ip" label="IP地址" width="140" />
      </el-table>

      <div class="pagination-row">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.itemsPerPage"
          :page-sizes="[10, 20, 50]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="loadLogs"
          @size-change="handleSizeChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { Delete, Refresh } from '@element-plus/icons-vue'
import api from '@/api'
import { success, error, confirm } from '@/utils/message'

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

    const allModules = computed(() => {
      if (modules.value.length > 0) {
        return modules.value.map((m) => getModuleText(m))
      }
      return Object.values(moduleMap)
    })

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

    const actionTypeMap = {
      'create': 'success',
      'update': 'primary',
      'delete': 'danger',
      'login': 'warning',
      'logout': 'info',
      'send': 'primary',
      'grant': 'success',
      'revoke': 'warning',
      'ban': 'danger',
      'unban': 'success',
      'pin': 'warning',
      'unpin': 'info',
      'approve': 'success',
      'reject': 'danger'
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

    const getActionType = (action) => {
      return actionTypeMap[action] || 'info'
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
        const response = await api.get('/system-logs', { params })
        logs.value = response.data.logs
        total.value = response.data.total
      } catch (err) {
        console.error('加载日志失败:', err)
      } finally {
        loading.value = false
      }
    }

    const handleSizeChange = () => {
      pagination.value.page = 1
      loadLogs()
    }

    const loadModules = async () => {
      try {
        const response = await api.get('/system-logs/modules')
        modules.value = response.data.modules
      } catch (err) {
        console.error('加载模块列表失败:', err)
      }
    }

    const deleteOldLogs = async () => {
      let confirmed
      try {
        confirmed = await confirm('确定要删除90天前的所有日志吗？此操作不可恢复。')
      } catch (err) {
        confirmed = null
      }
      if (!confirmed) return

      deleting.value = true
      try {
        await api.delete('/system-logs/old', { params: { days: 90 } })
        success('旧日志删除成功')
        await loadLogs()
      } catch (err) {
        console.error('删除旧日志失败:', err)
        error('删除旧日志失败')
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
      allModules,
      formatDate,
      getModuleText,
      getActionText,
      getActionType,
      deleteOldLogs,
      refreshData,
      handleSizeChange
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

.header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 16px;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.page-title {
  margin: 0;
  font-size: 20px;
  font-weight: 700;
  color: #303133;
}

.page-desc {
  margin: 4px 0 0;
  font-size: 14px;
  color: #909399;
}

.filter-row {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
}

.filter-module {
  width: 200px;
}

.filter-user {
  width: 200px;
}

.logs-table {
  width: 100%;
}

.user-name {
  font-weight: 600;
  color: #409eff;
}

.user-id {
  font-size: 12px;
  color: #909399;
}

.pagination-row {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.mb-6 {
  margin-bottom: 24px;
}
</style>
