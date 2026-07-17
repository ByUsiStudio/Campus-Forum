<template>
  <div class="admin-system-logs">
    <div class="page-header mb-6">
      <div class="flex items-center justify-between flex-wrap gap-4">
        <div>
          <h1 style="font-size: 18px; font-weight: 600; margin-bottom: 5px;">系统操作日志</h1>
          <p style="color: #999; font-size: 14px;">查看用户和管理员的所有操作记录</p>
        </div>
        <div class="flex gap-3">
          <button class="layui-btn layui-btn-danger" @click="deleteOldLogs" :disabled="deleting">
            <i class="fa-solid fa-trash-can mr-2"></i>
            删除90天前日志
          </button>
          <button class="layui-btn" @click="refreshData" :disabled="loading">
            <i class="fa-solid fa-refresh mr-2"></i>
            刷新
          </button>
        </div>
      </div>
    </div>

    <div class="layui-card mb-6">
      <div class="layui-card-body">
        <div class="layui-row">
          <div class="layui-col-xs12 layui-col-sm6 layui-col-md4 layui-col-lg3">
            <div class="layui-form-item">
              <label class="layui-form-label">筛选模块</label>
              <div class="layui-input-block">
                <select v-model="filter.module" class="layui-select">
                  <option value="">全部模块</option>
                  <option v-for="module in modules" :key="module" :value="module">{{ getModuleText(module) }}</option>
                </select>
              </div>
            </div>
          </div>
          <div class="layui-col-xs12 layui-col-sm6 layui-col-md4 layui-col-lg3">
            <div class="layui-form-item">
              <label class="layui-form-label">用户ID筛选</label>
              <div class="layui-input-block">
                <input type="number" v-model="filter.user_id" class="layui-input" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="layui-card">
      <div v-if="loading" style="height: 4px; background: #E8E8E8;">
        <div style="height: 100%; background: #1E9FFF; animation: loading 1.5s infinite;"></div>
      </div>

      <div class="layui-card-body">
        <table class="layui-table" v-if="logs.length > 0">
          <thead>
            <tr>
              <th>时间</th>
              <th style="width: 180px;">用户</th>
              <th style="width: 120px;">操作</th>
              <th style="width: 120px;">模块</th>
              <th>详情</th>
              <th style="width: 140px;">IP地址</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="log in logs" :key="log.id">
              <td>{{ formatDate(log.created_at) }}</td>
              <td>
                <div v-if="log.user">
                  <div style="font-weight: 600; color: #1E9FFF;">{{ log.user.username || log.user.nickname || '用户' }}</div>
                  <div style="color: #999; font-size: 12px;">ID: {{ log.user_id }}</div>
                </div>
                <div v-else style="color: #999; font-size: 12px;">用户ID: {{ log.user_id }}</div>
              </td>
              <td><span class="layui-badge" :style="{ backgroundColor: getActionColor(log.action) }">{{ getActionText(log.action) }}</span></td>
              <td><span class="layui-badge layui-bg-blue">{{ getModuleText(log.module) }}</span></td>
              <td class="text-ellipsis" :title="log.details">{{ log.details }}</td>
              <td>{{ log.ip }}</td>
            </tr>
          </tbody>
        </table>
        <div v-else class="text-center py-8 text-muted">暂无日志数据</div>
      </div>

      <div v-if="total > 20" class="layui-card-body border-t flex items-center justify-center gap-4">
        <div class="layui-laypage">
          <button 
            v-for="page in visiblePages" 
            :key="page"
            class="layui-laypage-btn"
            :class="{ 'layui-laypage-curr': page === pagination.page }"
            @click="pagination.page = page; loadLogs()"
          >
            {{ page }}
          </button>
        </div>
        <div style="color: #999; font-size: 14px;">
          第 {{ pagination.page }} 页 (共 {{ total }} 条)
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
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
      'create': '#52C41A',
      'update': '#1E9FFF',
      'delete': '#FF4D4F',
      'login': '#722ED1',
      'logout': '#999',
      'send': '#13C2C2',
      'grant': '#389E0D',
      'revoke': '#FA8C16',
      'ban': '#CF1322',
      'unban': '#95DE64',
      'pin': '#FAAD14',
      'unpin': '#999',
      'approve': '#95DE64',
      'reject': '#FF7875'
    }

    const totalPages = computed(() => Math.ceil(total.value / pagination.value.itemsPerPage))

    const visiblePages = computed(() => {
      const total = totalPages.value
      const current = pagination.value.page
      const pages = []
      const start = Math.max(1, current - 2)
      const end = Math.min(total, current + 2)
      for (let i = start; i <= end; i++) {
        pages.push(i)
      }
      return pages
    })

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
      return actionColorMap[action] || '#999'
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
      visiblePages,
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

.text-ellipsis {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 300px;
}

.border-t {
  border-top: 1px solid #E8E8E8;
}

@keyframes loading {
  0% { width: 0%; }
  50% { width: 50%; }
  100% { width: 100%; }
}
</style>