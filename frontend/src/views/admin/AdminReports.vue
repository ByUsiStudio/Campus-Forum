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
          <el-button type="primary" :loading="loading" @click="refreshData">
            <template #icon>
              <el-icon><Refresh /></el-icon>
            </template>
            刷新
          </el-button>
        </div>
      </div>
    </div>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="mb-6">
      <el-col :xs="24" :sm="8">
        <el-card shadow="never" class="stat-card stat-card-warning">
          <div class="d-flex align-center justify-space-between">
            <div>
              <div class="text-h4 font-weight-bold">{{ stats.pending }}</div>
              <div class="text-body-2 text-medium-emphasis">待处理</div>
            </div>
            <el-avatar :size="48" class="stat-avatar-warning">
              <el-icon :size="26"><AlarmClock /></el-icon>
            </el-avatar>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="8">
        <el-card shadow="never" class="stat-card stat-card-success">
          <div class="d-flex align-center justify-space-between">
            <div>
              <div class="text-h4 font-weight-bold">{{ stats.resolved }}</div>
              <div class="text-body-2 text-medium-emphasis">已处理</div>
            </div>
            <el-avatar :size="48" class="stat-avatar-success">
              <el-icon :size="26"><CircleCheck /></el-icon>
            </el-avatar>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="8">
        <el-card shadow="never" class="stat-card stat-card-error">
          <div class="d-flex align-center justify-space-between">
            <div>
              <div class="text-h4 font-weight-bold">{{ stats.rejected }}</div>
              <div class="text-body-2 text-medium-emphasis">已驳回</div>
            </div>
            <el-avatar :size="48" class="stat-avatar-error">
              <el-icon :size="26"><CircleClose /></el-icon>
            </el-avatar>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 筛选和搜索 -->
    <el-card shadow="never" class="mb-4">
      <el-row :gutter="16" align="middle">
        <el-col :xs="24" :md="10">
          <el-input
            v-model="filters.search"
            placeholder="搜索举报内容..."
            clearable
            @input="debounceSearch"
            @clear="loadReports"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-col>
        <el-col :xs="12" :md="7">
          <el-select
            v-model="filters.status"
            placeholder="状态"
            clearable
            style="width: 100%"
            @change="loadReports"
          >
            <el-option
              v-for="opt in statusOptions"
              :key="opt.value"
              :label="opt.title"
              :value="opt.value"
            />
          </el-select>
        </el-col>
        <el-col :xs="12" :md="7">
          <el-select
            v-model="filters.target_type"
            placeholder="举报类型"
            clearable
            style="width: 100%"
            @change="loadReports"
          >
            <el-option
              v-for="opt in targetTypeOptions"
              :key="opt.value"
              :label="opt.title"
              :value="opt.value"
            />
          </el-select>
        </el-col>
      </el-row>
    </el-card>

    <!-- 举报列表 -->
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>举报列表 ({{ totalReports }})</span>
          <el-pagination
            v-model:current-page="page"
            :page-size="pageSize"
            layout="prev, pager, next"
            :total="totalReports"
            :pager-count="5"
            @current-change="loadReports"
          />
        </div>
      </template>

      <div class="table-responsive">
      <el-table
        :data="reports"
        v-loading="loading"
        style="width: 100%"
      >
        <el-table-column label="举报人" min-width="200">
          <template #default="{ row }">
            <div class="d-flex align-center py-1">
              <el-avatar :size="40" class="mr-3">
                <img v-if="row.reporter?.avatar" :src="row.reporter.avatar" alt="" />
                <span v-else class="avatar-fallback">
                  {{ row.reporter?.display_name?.[0] || row.reporter?.username?.[0] || '?' }}
                </span>
              </el-avatar>
              <div>
                <div class="font-weight-medium">{{ row.reporter?.display_name || row.reporter?.username }}</div>
                <div class="el-table-caption">{{ formatDate(row.created_at) }}</div>
              </div>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="举报类型" width="150">
          <template #default="{ row }">
            <div>
              <el-tag
                size="small"
                :type="getTargetTypeTagType(row.target_type)"
                effect="light"
              >
                {{ getTargetTypeText(row.target_type) }}
              </el-tag>
              <div class="el-table-caption mt-1">ID: {{ row.target_id }}</div>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="举报原因" min-width="200">
          <template #default="{ row }">
            <div class="text-truncate" style="max-width: 200px;">
              {{ row.reason }}
            </div>
          </template>
        </el-table-column>

        <el-table-column label="状态" width="120" align="center">
          <template #default="{ row }">
            <el-tag
              :type="getStatusTagType(row.status)"
              size="small"
              effect="light"
            >
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>

        <el-table-column label="操作" width="120" align="center">
          <template #default="{ row }">
            <el-button
              type="primary"
              link
              :icon="View"
              @click="viewReport(row)"
            />
            <el-dropdown
              v-if="row.status === 'pending'"
              trigger="click"
              @command="(cmd) => handleDropdown(row, cmd)"
            >
              <el-button type="primary" link :icon="MoreFilled" />
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="resolved">
                    <el-icon><Check /></el-icon>
                    标记为已处理
                  </el-dropdown-item>
                  <el-dropdown-item command="rejected" divided>
                    <el-icon><Close /></el-icon>
                    驳回举报
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>
      </div>
    </el-card>
  </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import {
  Refresh,
  Search,
  AlarmClock,
  CircleCheck,
  CircleClose,
  View,
  MoreFilled,
  Check,
  Close
} from '@element-plus/icons-vue'
import { reportApi } from '../../api'
import { prompt, error, success } from '@/utils/message'
import { jcCloseAll, jcOpenHtml } from '@/utils/jcu'

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
    const selectedReport = ref(null)
    const narrow = ref(window.innerWidth < 576)

    const updateNarrow = () => { narrow.value = window.innerWidth < 576 }
    window.addEventListener('resize', updateNarrow)

    const filters = ref({
      search: '',
      status: null,
      target_type: null
    })

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
        page.value = 1
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
        if (filters.value.search) params.search = filters.value.search
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

    const getTargetTypeTagType = (type) => {
      const types = {
        article: 'primary',
        comment: 'warning',
        user: 'success'
      }
      return types[type] || 'info'
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

    const getStatusTagType = (status) => {
      const types = {
        pending: 'warning',
        resolved: 'success',
        rejected: 'danger'
      }
      return types[status] || 'info'
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
      const r = report
      const pending = r.status === 'pending'

      const esc = (s) => String(s ?? '').replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;')
      const item = (label, value) =>
        `<div style="margin-bottom:12px;"><div style="font-size:12px;color:var(--jc-text-2,#64748b);font-weight:600;margin-bottom:4px;">${label}</div><div style="font-size:14px;color:var(--jc-text,#0f172a);">${value}</div></div>`

      let html = ''
      html += item('举报人', esc(r.reporter?.display_name || r.reporter?.username || '未知'))
      html += `<div style="display:flex;gap:16px;flex-wrap:wrap;">
        ${item('举报类型', esc(getTargetTypeText(r.target_type)))}
        ${item('被举报内容 ID', esc(r.target_id))}
      </div>`
      html += item('举报原因', esc(r.reason))
      if (r.description) html += item('详细描述', esc(r.description))
      html += item('举报时间', esc(formatDateTime(r.created_at)))

      if (r.status !== 'pending') {
        html += `<div style="border-top:1px solid var(--jc-border,#e6e9f2);padding-top:12px;margin-top:12px;">`
        html += `<div style="font-size:12px;color:var(--jc-text-2,#64748b);font-weight:600;margin-bottom:8px;">处理信息</div>`
        html += item('状态', esc(getStatusText(r.status)))
        if (r.handler) html += item('处理人', esc(r.handler?.display_name || r.handler?.username || ''))
        if (r.handle_note) html += item('处理备注', esc(r.handle_note))
        if (r.handled_at) html += item('处理时间', esc(formatDateTime(r.handled_at)))
        html += `</div>`
      }

      const buttons = []
      if (pending) {
        buttons.push({ text: '驳回', type: 'danger', action: () => handleReport(r, 'rejected') })
        buttons.push({ text: '标记已处理', type: 'primary', action: () => handleReport(r, 'resolved') })
      } else {
        buttons.push({ text: '关闭', type: 'default', action: () => jcCloseAll() })
      }

      jcOpenHtml({
        title: '举报详情',
        content: html,
        width: 700,
        size: 'md',
        buttons
      })
    }

    const handleDropdown = (report, command) => {
      handleReport(report, command)
    }

    const handleReport = async (report, status) => {
      const note = await prompt(
        `请输入处理备注 (${status === 'resolved' ? '已处理' : '驳回'}):`,
        { placeholder: '处理备注' }
      ).catch(() => null)
      if (note === null) return

      try {
        await reportApi.handleReport(report.id, {
          status,
          handle_note: note
        })
        success('举报已处理')
        jcCloseAll()
        loadReports()
      } catch (error) {
        console.error('处理举报失败:', error)
        error('处理失败')
      }
    }

    onMounted(() => {
      loadReports()
    })

    onUnmounted(() => {
      window.removeEventListener('resize', updateNarrow)
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
      statusOptions,
      targetTypeOptions,
      selectedReport,
      narrow,
      debounceSearch,
      loadReports,
      refreshData,
      getTargetTypeColor,
      getTargetTypeTagType,
      getTargetTypeText,
      getStatusColor,
      getStatusTagType,
      getStatusText,
      formatDate,
      formatDateTime,
      viewReport,
      handleDropdown,
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

.mb-6 {
  margin-bottom: 24px;
}

.mb-4 {
  margin-bottom: 16px;
}

.mb-1 {
  margin-bottom: 4px;
}

.mb-2 {
  margin-bottom: 8px;
}

.mt-1 {
  margin-top: 4px;
}

.mt-2 {
  margin-top: 8px;
}

.mr-2 {
  margin-right: 8px;
}

.mr-3 {
  margin-right: 12px;
}

.d-flex {
  display: flex;
}

.align-center {
  align-items: center;
}

.justify-space-between {
  justify-content: space-between;
}

.flex-wrap {
  flex-wrap: wrap;
}

.ga-4 {
  gap: 16px;
}

.ga-3 {
  gap: 12px;
}

.py-1 {
  padding-top: 4px;
  padding-bottom: 4px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.table-responsive {
  overflow-x: auto;
}

.text-h4 {
  font-size: 24px;
  font-weight: 700;
  line-height: 1.2;
}

.text-h5 {
  font-size: 20px;
  font-weight: 700;
}

.font-weight-bold {
  font-weight: 700;
}

.font-weight-medium {
  font-weight: 500;
}

.text-body-2 {
  font-size: 14px;
  line-height: 1.5;
}

.text-caption {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.text-medium-emphasis {
  color: var(--el-text-color-secondary);
}

.text-primary {
  color: var(--el-color-primary);
}

.text-truncate {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.avatar-fallback {
  font-size: 14px;
  font-weight: 700;
  color: var(--el-color-primary);
}

.el-table-caption {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.stat-card {
  border-left: 4px solid;
  transition: transform 0.2s;
}

.stat-card:hover {
  transform: translateY(-4px);
}

.stat-card-warning {
  --stat-border: var(--el-color-warning);
  border-left-color: var(--stat-border);
}

.stat-card-success {
  --stat-border: var(--el-color-success);
  border-left-color: var(--stat-border);
}

.stat-card-error {
  --stat-border: var(--el-color-danger);
  border-left-color: var(--stat-border);
}

.stat-avatar-warning {
  background: var(--el-color-warning-light-9);
  color: var(--el-color-warning);
}

.stat-avatar-success {
  background: var(--el-color-success-light-9);
  color: var(--el-color-success);
}

.stat-avatar-error {
  background: var(--el-color-danger-light-9);
  color: var(--el-color-danger);
}

.page-header h1 {
  color: var(--el-text-color-primary);
}
</style>
