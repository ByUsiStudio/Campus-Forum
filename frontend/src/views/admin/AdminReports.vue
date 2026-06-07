<template>
  <div class="admin-reports">
    <!-- 页面标题 -->
    <div class="page-header mb-6">
      <div class="d-flex align-center justify-space-between flex-wrap ga-4">
        <div>
          <h1 class="text-h5 font-weight-bold mb-1">举报管理</h1>
          <p class="text-body-2 text-medium-emphasis">查看和处理用户提交的举报内容</p>
        </div>
        <div class="d-flex ga-3">
          <v-btn color="primary" @click="refreshData" :loading="loading" variant="tonal">
            <v-icon start>mdi-refresh</v-icon>
            刷新
          </v-btn>
        </div>
      </div>
    </div>

    <!-- 统计卡片 -->
    <v-row class="mb-6">
      <v-col cols="12" sm="4">
        <v-card elevation="0" class="stat-card stat-card-warning">
          <v-card-text>
            <div class="d-flex align-center justify-space-between">
              <div>
                <div class="text-h4 font-weight-bold">{{ stats.pending }}</div>
                <div class="text-body-2 text-medium-emphasis">待处理</div>
              </div>
              <v-avatar color="warning" size="48" rounded="lg">
                <v-icon>mdi-clock-alert</v-icon>
              </v-avatar>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" sm="4">
        <v-card elevation="0" class="stat-card stat-card-success">
          <v-card-text>
            <div class="d-flex align-center justify-space-between">
              <div>
                <div class="text-h4 font-weight-bold">{{ stats.resolved }}</div>
                <div class="text-body-2 text-medium-emphasis">已处理</div>
              </div>
              <v-avatar color="success" size="48" rounded="lg">
                <v-icon>mdi-check-circle</v-icon>
              </v-avatar>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" sm="4">
        <v-card elevation="0" class="stat-card stat-card-error">
          <v-card-text>
            <div class="d-flex align-center justify-space-between">
              <div>
                <div class="text-h4 font-weight-bold">{{ stats.rejected }}</div>
                <div class="text-body-2 text-medium-emphasis">已驳回</div>
              </div>
              <v-avatar color="error" size="48" rounded="lg">
                <v-icon>mdi-close-circle</v-icon>
              </v-avatar>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- 筛选和搜索 -->
    <v-card elevation="0" class="mb-4">
      <v-card-text>
        <v-row align="center" dense>
          <v-col cols="12" md="6">
            <v-text-field
              v-model="filters.search"
              placeholder="搜索举报内容..."
              variant="outlined"
              density="compact"
              prepend-inner-icon="mdi-magnify"
              hide-details
              clearable
              @update:model-value="debounceSearch"
            />
          </v-col>
          <v-col cols="6" md="3">
            <v-select
              v-model="filters.status"
              :items="statusOptions"
              label="状态"
              variant="outlined"
              density="compact"
              hide-details
              clearable
            />
          </v-col>
          <v-col cols="6" md="3">
            <v-select
              v-model="filters.target_type"
              :items="targetTypeOptions"
              label="举报类型"
              variant="outlined"
              density="compact"
              hide-details
              clearable
            />
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

    <!-- 举报列表 -->
    <v-card elevation="0">
      <v-card-title class="d-flex align-center justify-space-between">
        <span>举报列表 ({{ totalReports }})</span>
        <v-pagination
          v-model="page"
          :length="totalPages"
          :total-visible="5"
          density="compact"
          @update:model-value="loadReports"
        />
      </v-card-title>

      <v-data-table
        :headers="headers"
        :items="reports"
        :loading="loading"
        :items-per-page="pageSize"
        class="elevation-0"
      >
        <template v-slot:item.reporter="{ item }">
          <div class="d-flex align-center py-2">
            <v-avatar size="40" class="mr-3">
              <v-img v-if="item.reporter?.avatar" :src="item.reporter.avatar" />
              <span v-else class="text-primary font-weight-bold">
                {{ item.reporter?.display_name?.[0] || item.reporter?.username?.[0] || '?' }}
              </span>
            </v-avatar>
            <div>
              <div class="font-weight-medium">{{ item.reporter?.display_name || item.reporter?.username }}</div>
              <div class="text-caption text-medium-emphasis">{{ formatDate(item.created_at) }}</div>
            </div>
          </div>
        </template>

        <template v-slot:item.target="{ item }">
          <div>
            <v-chip size="small" :color="getTargetTypeColor(item.target_type)" variant="tonal">
              {{ getTargetTypeText(item.target_type) }}
            </v-chip>
            <div class="text-caption mt-1">ID: {{ item.target_id }}</div>
          </div>
        </template>

        <template v-slot:item.reason="{ item }">
          <div class="text-truncate" style="max-width: 200px;">
            {{ item.reason }}
          </div>
        </template>

        <template v-slot:item.status="{ item }">
          <v-chip
            :color="getStatusColor(item.status)"
            size="small"
            variant="tonal"
          >
            {{ getStatusText(item.status) }}
          </v-chip>
        </template>

        <template v-slot:item.actions="{ item }">
          <v-btn icon variant="text" size="small" @click="viewReport(item)">
            <v-icon>mdi-eye</v-icon>
          </v-btn>
          <v-menu v-if="item.status === 'pending'">
            <template v-slot:activator="{ props }">
              <v-btn icon variant="text" size="small" v-bind="props">
                <v-icon>mdi-dots-vertical</v-icon>
              </v-btn>
            </template>
            <v-list density="compact">
              <v-list-item @click="handleReport(item, 'resolved')" prepend-icon="mdi-check">
                <v-list-item-title>标记为已处理</v-list-item-title>
              </v-list-item>
              <v-list-item @click="handleReport(item, 'rejected')" prepend-icon="mdi-close">
                <v-list-item-title>驳回举报</v-list-item-title>
              </v-list-item>
            </v-list>
          </v-menu>
        </template>
      </v-data-table>
    </v-card>

    <!-- 查看举报详情对话框 -->
    <v-dialog v-model="viewDialog" max-width="700">
      <v-card v-if="selectedReport">
        <v-card-title class="bg-primary text-white pa-4">
          <v-icon class="mr-2">mdi-flag-variant</v-icon>
          举报详情
        </v-card-title>
        <v-card-text class="pa-4">
          <v-row dense>
            <v-col cols="6">
              <div class="text-caption text-medium-emphasis mb-1">举报人</div>
              <div class="d-flex align-center">
                <v-avatar size="32" class="mr-2">
                  <v-img v-if="selectedReport.reporter?.avatar" :src="selectedReport.reporter.avatar" />
                  <span v-else class="text-primary text-body-2">
                    {{ selectedReport.reporter?.display_name?.[0] || '?' }}
                  </span>
                </v-avatar>
                {{ selectedReport.reporter?.display_name || selectedReport.reporter?.username }}
              </div>
            </v-col>
            <v-col cols="6">
              <div class="text-caption text-medium-emphasis mb-1">举报类型</div>
              <v-chip size="small" :color="getTargetTypeColor(selectedReport.target_type)" variant="tonal">
                {{ getTargetTypeText(selectedReport.target_type) }}
              </v-chip>
            </v-col>
            <v-col cols="12">
              <div class="text-caption text-medium-emphasis mb-1">被举报内容ID</div>
              <div class="font-weight-medium">{{ selectedReport.target_id }}</div>
            </v-col>
            <v-col cols="12">
              <div class="text-caption text-medium-emphasis mb-1">举报原因</div>
              <div class="font-weight-medium">{{ selectedReport.reason }}</div>
            </v-col>
            <v-col cols="12" v-if="selectedReport.description">
              <div class="text-caption text-medium-emphasis mb-1">详细描述</div>
              <div class="text-body-2">{{ selectedReport.description }}</div>
            </v-col>
            <v-col cols="12">
              <div class="text-caption text-medium-emphasis mb-1">举报时间</div>
              <div>{{ formatDateTime(selectedReport.created_at) }}</div>
            </v-col>
            <v-col cols="12" v-if="selectedReport.status !== 'pending'">
              <v-divider class="mb-4" />
              <div class="text-caption text-medium-emphasis mb-1">处理信息</div>
              <div class="d-flex align-center mb-2">
                <v-chip size="small" :color="getStatusColor(selectedReport.status)" variant="tonal" class="mr-2">
                  {{ getStatusText(selectedReport.status) }}
                </v-chip>
                <span v-if="selectedReport.handler">处理人: {{ selectedReport.handler?.display_name || selectedReport.handler?.username }}</span>
              </div>
              <div v-if="selectedReport.handle_note" class="text-body-2">
                <div class="text-caption text-medium-emphasis">处理备注:</div>
                <div>{{ selectedReport.handle_note }}</div>
              </div>
              <div v-if="selectedReport.handled_at" class="text-caption text-medium-emphasis mt-2">
                处理时间: {{ formatDateTime(selectedReport.handled_at) }}
              </div>
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions class="pa-4" v-if="selectedReport.status === 'pending'">
          <v-spacer />
          <v-btn variant="text" @click="viewDialog = false">关闭</v-btn>
          <v-btn color="error" @click="handleReport(selectedReport, 'rejected')">驳回</v-btn>
          <v-btn color="success" @click="handleReport(selectedReport, 'resolved')">标记已处理</v-btn>
        </v-card-actions>
        <v-card-actions class="pa-4" v-else>
          <v-spacer />
          <v-btn variant="text" @click="viewDialog = false">关闭</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { reportApi } from '../../api'

export default {
  name: 'AdminReports',
  setup() {
    const loading = ref(false)
    const reports = ref([])
    const stats = ref({
      pending: 0,
      resolved: 0,
      rejected: 0
    })
    const page = ref(1)
    const pageSize = ref(20)
    const totalReports = ref(0)
    const totalPages = computed(() => Math.ceil(totalReports.value / pageSize.value))
    const viewDialog = ref(false)
    const selectedReport = ref(null)
    const handleNote = ref('')

    const filters = ref({
      search: '',
      status: null,
      target_type: null
    })

    const headers = [
      { title: '举报人', key: 'reporter', sortable: false },
      { title: '举报类型', key: 'target', width: '150px' },
      { title: '举报原因', key: 'reason' },
      { title: '状态', key: 'status', width: '120px' },
      { title: '时间', key: 'created_at', width: '180px' },
      { title: '操作', key: 'actions', width: '100px', sortable: false }
    ]

    const statusOptions = [
      { title: '全部', value: null },
      { title: '待处理', value: 'pending' },
      { title: '已处理', value: 'resolved' },
      { title: '已驳回', value: 'rejected' }
    ]

    const targetTypeOptions = [
      { title: '全部', value: null },
      { title: '文章', value: 'article' },
      { title: '评论', value: 'comment' },
      { title: '用户', value: 'user' }
    ]

    let searchTimeout = null

    const debounceSearch = () => {
      if (searchTimeout) clearTimeout(searchTimeout)
      searchTimeout = setTimeout(() => {
        loadReports()
      }, 500)
    }

    const loadReports = async () => {
      loading.value = true
      try {
        const params = {
          page: page.value,
          page_size: pageSize.value
        }
        if (filters.value.status) params.status = filters.value.status
        if (filters.value.target_type) params.target_type = filters.value.target_type

        const response = await reportApi.getReports(params)
        reports.value = response.data.reports || []
        totalReports.value = response.data.total || 0

        // 更新统计数据
        stats.value.pending = response.data.pending_count || 0
        stats.value.resolved = response.data.resolved_count || 0
        stats.value.rejected = response.data.rejected_count || 0
      } catch (error) {
        console.error('加载举报列表失败:', error)
      } finally {
        loading.value = false
      }
    }

    const refreshData = () => {
      loadReports()
    }

    const getTargetTypeColor = (type) => {
      const colors = {
        article: 'blue',
        comment: 'purple',
        user: 'orange'
      }
      return colors[type] || 'grey'
    }

    const getTargetTypeText = (type) => {
      const texts = {
        article: '文章',
        comment: '评论',
        user: '用户'
      }
      return texts[type] || '未知'
    }

    const getStatusColor = (status) => {
      const colors = {
        pending: 'warning',
        resolved: 'success',
        rejected: 'error'
      }
      return colors[status] || 'grey'
    }

    const getStatusText = (status) => {
      const texts = {
        pending: '待处理',
        resolved: '已处理',
        rejected: '已驳回'
      }
      return texts[status] || '未知'
    }

    const formatDate = (date) => {
      if (!date) return ''
      return new Date(date).toLocaleDateString('zh-CN', {
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      })
    }

    const formatDateTime = (date) => {
      if (!date) return ''
      return new Date(date).toLocaleString('zh-CN')
    }

    const viewReport = (report) => {
      selectedReport.value = report
      handleNote.value = ''
      viewDialog.value = true
    }

    const handleReport = async (report, status) => {
      const note = prompt(`请输入处理备注 (${status === 'resolved' ? '已处理' : '驳回'}):`)
      if (note === null) return

      try {
        await reportApi.handleReport(report.id, {
          status,
          handle_note: note
        })
        viewDialog.value = false
        loadReports()
      } catch (error) {
        console.error('处理举报失败:', error)
        alert('处理失败')
      }
    }

    onMounted(() => {
      loadReports()
    })

    return {
      loading,
      reports,
      stats,
      page,
      pageSize,
      totalReports,
      totalPages,
      filters,
      headers,
      statusOptions,
      targetTypeOptions,
      viewDialog,
      selectedReport,
      handleNote,
      debounceSearch,
      loadReports,
      refreshData,
      getTargetTypeColor,
      getTargetTypeText,
      getStatusColor,
      getStatusText,
      formatDate,
      formatDateTime,
      viewReport,
      handleReport
    }
  }
}
</script>

<style scoped>
.admin-reports {
  max-width: 1600px;
  margin: 0 auto;
}

.stat-card {
  border-left: 4px solid;
  transition: transform 0.2s;
}

.stat-card:hover {
  transform: translateY(-4px);
}

.stat-card-warning {
  border-left-color: rgb(var(--v-theme-warning));
  background: linear-gradient(135deg, rgba(var(--v-theme-warning), 0.08), rgba(var(--v-theme-warning), 0.02));
}

.stat-card-success {
  border-left-color: rgb(var(--v-theme-success));
  background: linear-gradient(135deg, rgba(var(--v-theme-success), 0.08), rgba(var(--v-theme-success), 0.02));
}

.stat-card-error {
  border-left-color: rgb(var(--v-theme-error));
  background: linear-gradient(135deg, rgba(var(--v-theme-error), 0.08), rgba(var(--v-theme-error), 0.02));
}

.page-header h1 {
  color: rgb(var(--v-theme-on-surface));
}
</style>
