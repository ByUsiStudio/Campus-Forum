<template>
  <div class="admin-reports">
    <div class="page-header">
      <div class="header-left">
        <h1>举报管理</h1>
        <p>查看和处理用户提交的举报内容</p>
      </div>
      <div class="header-right">
        <button class="layui-btn" @click="refreshData" :disabled="loading">
          <i class="fa-solid fa-rotate-right"></i>刷新
        </button>
      </div>
    </div>

    <div class="stats-row">
      <div class="stat-card stat-card-warning">
        <div class="stat-content">
          <div class="stat-value">{{ stats.pending }}</div>
          <div class="stat-label">待处理</div>
        </div>
        <div class="stat-icon">
          <i class="fa-solid fa-clock"></i>
        </div>
      </div>
      <div class="stat-card stat-card-success">
        <div class="stat-content">
          <div class="stat-value">{{ stats.resolved }}</div>
          <div class="stat-label">已处理</div>
        </div>
        <div class="stat-icon">
          <i class="fa-solid fa-check-circle"></i>
        </div>
      </div>
      <div class="stat-card stat-card-error">
        <div class="stat-content">
          <div class="stat-value">{{ stats.rejected }}</div>
          <div class="stat-label">已驳回</div>
        </div>
        <div class="stat-icon">
          <i class="fa-solid fa-xmark-circle"></i>
        </div>
      </div>
    </div>

    <div class="layui-card filter-card">
      <div class="layui-card-body">
        <div class="filter-row">
          <div class="filter-item">
            <input
              v-model="filters.search"
              type="text"
              placeholder="搜索举报内容..."
              class="layui-input"
              @input="debounceSearch"
            />
          </div>
          <div class="filter-item">
            <select v-model="filters.status" class="layui-input" @change="loadReports">
              <option :value="null">全部状态</option>
              <option value="pending">待处理</option>
              <option value="resolved">已处理</option>
              <option value="rejected">已驳回</option>
            </select>
          </div>
          <div class="filter-item">
            <select v-model="filters.target_type" class="layui-input" @change="loadReports">
              <option :value="null">全部类型</option>
              <option value="article">文章</option>
              <option value="comment">评论</option>
              <option value="user">用户</option>
            </select>
          </div>
        </div>
      </div>
    </div>

    <div class="layui-card">
      <div class="layui-card-header">
        <span>举报列表 ({{ totalReports }})</span>
        <div class="pagination-wrapper">
          <button class="page-btn" :disabled="page <= 1" @click="page--; loadReports()">
            <i class="fa-solid fa-chevron-left"></i>
          </button>
          <span class="page-info">{{ page }} / {{ totalPages }}</span>
          <button class="page-btn" :disabled="page >= totalPages" @click="page++; loadReports()">
            <i class="fa-solid fa-chevron-right"></i>
          </button>
        </div>
      </div>

      <div v-if="loading" class="loading-state">
        <i class="fa-solid fa-spinner fa-spin"></i>
      </div>

      <div v-else class="report-table">
        <table>
          <thead>
            <tr>
              <th>举报人</th>
              <th>举报类型</th>
              <th>举报原因</th>
              <th>状态</th>
              <th>时间</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="report in reports" :key="report.id">
              <td>
                <div class="reporter-cell">
                  <div class="reporter-avatar">
                    <img v-if="report.reporter?.avatar" :src="report.reporter.avatar" />
                    <span v-else>{{ report.reporter?.display_name?.[0] || report.reporter?.username?.[0] || '?' }}</span>
                  </div>
                  <div>
                    <div>{{ report.reporter?.display_name || report.reporter?.username }}</div>
                    <div class="time-caption">{{ formatDate(report.created_at) }}</div>
                  </div>
                </div>
              </td>
              <td>
                <span class="type-chip" :class="getTargetTypeColor(report.target_type)">
                  {{ getTargetTypeText(report.target_type) }}
                </span>
                <div class="id-caption">ID: {{ report.target_id }}</div>
              </td>
              <td class="reason-cell">{{ report.reason }}</td>
              <td>
                <span class="status-chip" :class="getStatusColor(report.status)">
                  {{ getStatusText(report.status) }}
                </span>
              </td>
              <td>{{ formatDate(report.created_at) }}</td>
              <td>
                <div class="action-buttons">
                  <button class="action-btn" @click="viewReport(report)">
                    <i class="fa-solid fa-eye"></i>
                  </button>
                  <div v-if="report.status === 'pending'" class="action-menu">
                    <button class="action-btn" @click="toggleMenu(report.id)">
                      <i class="fa-solid fa-ellipsis-vertical"></i>
                    </button>
                    <div v-if="activeMenu === report.id" class="menu-dropdown">
                      <button @click="handleReport(report, 'resolved'); toggleMenu(null)">
                        <i class="fa-solid fa-check"></i>标记为已处理
                      </button>
                      <button @click="handleReport(report, 'rejected'); toggleMenu(null)">
                        <i class="fa-solid fa-xmark"></i>驳回举报
                      </button>
                    </div>
                  </div>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div v-if="viewDialog" class="dialog-overlay" @click.self="viewDialog = false">
      <div class="dialog-content">
        <div class="dialog-header bg-primary">
          <i class="fa-solid fa-flag"></i>
          <span>举报详情</span>
          <button class="dialog-close" @click="viewDialog = false">
            <i class="fa-solid fa-xmark"></i>
          </button>
        </div>
        <div v-if="selectedReport" class="dialog-body">
          <div class="detail-row">
            <div class="detail-item">
              <div class="detail-label">举报人</div>
              <div class="detail-value">
                <div class="reporter-avatar small">
                  <img v-if="selectedReport.reporter?.avatar" :src="selectedReport.reporter.avatar" />
                  <span v-else>{{ selectedReport.reporter?.display_name?.[0] || '?' }}</span>
                </div>
                {{ selectedReport.reporter?.display_name || selectedReport.reporter?.username }}
              </div>
            </div>
            <div class="detail-item">
              <div class="detail-label">举报类型</div>
              <span class="type-chip" :class="getTargetTypeColor(selectedReport.target_type)">
                {{ getTargetTypeText(selectedReport.target_type) }}
              </span>
            </div>
          </div>
          <div class="detail-row">
            <div class="detail-item full">
              <div class="detail-label">被举报内容ID</div>
              <div class="detail-value">{{ selectedReport.target_id }}</div>
            </div>
          </div>
          <div class="detail-row">
            <div class="detail-item full">
              <div class="detail-label">举报原因</div>
              <div class="detail-value">{{ selectedReport.reason }}</div>
            </div>
          </div>
          <div v-if="selectedReport.description" class="detail-row">
            <div class="detail-item full">
              <div class="detail-label">详细描述</div>
              <div class="detail-value">{{ selectedReport.description }}</div>
            </div>
          </div>
          <div class="detail-row">
            <div class="detail-item full">
              <div class="detail-label">举报时间</div>
              <div class="detail-value">{{ formatDateTime(selectedReport.created_at) }}</div>
            </div>
          </div>
          <div v-if="selectedReport.status !== 'pending'" class="detail-row">
            <div class="detail-item full">
              <div class="detail-label">处理信息</div>
              <div class="detail-value">
                <div class="handle-info">
                  <span class="status-chip" :class="getStatusColor(selectedReport.status)">
                    {{ getStatusText(selectedReport.status) }}
                  </span>
                  <span v-if="selectedReport.handler">处理人: {{ selectedReport.handler?.display_name || selectedReport.handler?.username }}</span>
                </div>
                <div v-if="selectedReport.handle_note" class="handle-note">
                  <div class="note-label">处理备注:</div>
                  <div>{{ selectedReport.handle_note }}</div>
                </div>
                <div v-if="selectedReport.handled_at" class="handle-time">
                  处理时间: {{ formatDateTime(selectedReport.handled_at) }}
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="dialog-footer" v-if="selectedReport">
          <div v-if="selectedReport.status === 'pending'">
            <button class="layui-btn layui-btn-primary" @click="viewDialog = false">关闭</button>
            <button class="layui-btn layui-btn-danger" @click="handleReport(selectedReport, 'rejected')">驳回</button>
            <button class="layui-btn layui-btn-success" @click="handleReport(selectedReport, 'resolved')">标记已处理</button>
          </div>
          <div v-else>
            <button class="layui-btn" @click="viewDialog = false">关闭</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { reportApi } from '../../api'

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
const activeMenu = ref(null)

const filters = ref({
  search: '',
  status: null,
  target_type: null
})

let searchTimeout = null

const debounceSearch = () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    page.value = 1
    loadReports()
  }, 500)
}

const toggleMenu = (id) => {
  activeMenu.value = activeMenu.value === id ? null : id
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
  page.value = 1
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
</script>

<style scoped>
.admin-reports {
  padding: 20px;
  max-width: 1600px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 16px;
}

.header-left h1 {
  font-size: 24px;
  font-weight: 700;
  color: #333;
  margin: 0 0 8px;
}

.header-left p {
  font-size: 14px;
  color: #999;
  margin: 0;
}

.stats-row {
  display: flex;
  gap: 20px;
  margin-bottom: 24px;
}

.stat-card {
  flex: 1;
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  border-left: 4px solid;
  display: flex;
  justify-content: space-between;
  align-items: center;
  transition: transform 0.2s;

  &:hover {
    transform: translateY(-4px);
  }
}

.stat-card-warning {
  border-left-color: #FAAD14;
  background: linear-gradient(135deg, rgba(250, 173, 20, 0.08), rgba(250, 173, 20, 0.02));
}

.stat-card-success {
  border-left-color: #52C41A;
  background: linear-gradient(135deg, rgba(82, 196, 26, 0.08), rgba(82, 196, 26, 0.02));
}

.stat-card-error {
  border-left-color: #FF5722;
  background: linear-gradient(135deg, rgba(255, 87, 34, 0.08), rgba(255, 87, 34, 0.02));
}

.stat-value {
  font-size: 36px;
  font-weight: 700;
  color: #333;
}

.stat-label {
  font-size: 14px;
  color: #999;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
}

.stat-card-warning .stat-icon {
  background: rgba(250, 173, 20, 0.1);
  color: #FAAD14;
}

.stat-card-success .stat-icon {
  background: rgba(82, 196, 26, 0.1);
  color: #52C41A;
}

.stat-card-error .stat-icon {
  background: rgba(255, 87, 34, 0.1);
  color: #FF5722;
}

.filter-card {
  margin-bottom: 20px;
}

.filter-row {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.filter-item {
  flex: 1;
  min-width: 200px;
}

.layui-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 16px;
  font-weight: 600;
}

.pagination-wrapper {
  display: flex;
  align-items: center;
  gap: 12px;
}

.page-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  background: #fff;
  cursor: pointer;
  color: #666;

  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  &:hover:not(:disabled) {
    background: #f5f5f5;
    border-color: #1E9FFF;
    color: #1E9FFF;
  }
}

.page-info {
  font-size: 14px;
  color: #666;
  min-width: 60px;
  text-align: center;
}

.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 60px;
  font-size: 32px;
  color: #1E9FFF;
}

.report-table {
  overflow-x: auto;
}

.report-table table {
  width: 100%;
  border-collapse: collapse;
}

.report-table th,
.report-table td {
  padding: 12px 16px;
  text-align: left;
  border-bottom: 1px solid #f0f0f0;
}

.report-table th {
  background: #fafafa;
  font-weight: 600;
  color: #666;
  font-size: 14px;
}

.report-table tr:hover {
  background: #f9f9f9;
}

.reporter-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.reporter-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #f0f0f0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  color: #1E9FFF;
  overflow: hidden;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  &.small {
    width: 32px;
    height: 32px;
    font-size: 14px;
  }
}

.time-caption {
  font-size: 12px;
  color: #999;
}

.id-caption {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.reason-cell {
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.type-chip,
.status-chip {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  display: inline-block;
}

.type-chip.blue {
  background: rgba(30, 159, 255, 0.1);
  color: #1E9FFF;
}

.type-chip.purple {
  background: rgba(156, 39, 176, 0.1);
  color: #9C27B0;
}

.type-chip.orange {
  background: rgba(255, 152, 0, 0.1);
  color: #FF9800;
}

.status-chip.warning {
  background: rgba(250, 173, 20, 0.1);
  color: #FAAD14;
}

.status-chip.success {
  background: rgba(82, 196, 26, 0.1);
  color: #52C41A;
}

.status-chip.error {
  background: rgba(255, 87, 34, 0.1);
  color: #FF5722;
}

.action-buttons {
  display: flex;
  align-items: center;
  gap: 8px;
}

.action-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: transparent;
  border-radius: 4px;
  color: #999;
  cursor: pointer;

  &:hover {
    background: #f0f0f0;
    color: #1E9FFF;
  }
}

.action-menu {
  position: relative;
}

.menu-dropdown {
  position: absolute;
  top: 100%;
  right: 0;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  z-index: 100;
  margin-top: 4px;

  button {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 16px;
    width: 100%;
    border: none;
    background: transparent;
    text-align: left;
    font-size: 14px;
    cursor: pointer;
    color: #333;

    &:hover {
      background: #f5f5f5;
    }
  }
}

.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.dialog-content {
  background: #fff;
  border-radius: 12px;
  width: 90%;
  max-width: 700px;
  overflow: hidden;
}

.dialog-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  font-size: 18px;
  font-weight: 600;

  &.bg-primary {
    background: #1E9FFF;
    color: #fff;
  }
}

.dialog-close {
  padding: 4px;
  background: transparent;
  border: none;
  color: inherit;
  font-size: 20px;
  cursor: pointer;

  &:hover {
    opacity: 0.8;
  }
}

.dialog-body {
  padding: 24px;
}

.detail-row {
  display: flex;
  gap: 24px;
  margin-bottom: 16px;
}

.detail-item {
  flex: 1;

  &.full {
    flex: 100%;
  }
}

.detail-label {
  font-size: 12px;
  color: #999;
  margin-bottom: 8px;
}

.detail-value {
  font-size: 14px;
  color: #333;
  display: flex;
  align-items: center;
  gap: 8px;
}

.handle-info {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.handle-note {
  background: #f8f9fa;
  padding: 12px;
  border-radius: 6px;
  margin-bottom: 8px;

  .note-label {
    font-size: 12px;
    color: #999;
    margin-bottom: 4px;
  }
}

.handle-time {
  font-size: 12px;
  color: #999;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px;
  border-top: 1px solid #f0f0f0;
}

@media (max-width: 768px) {
  .stats-row {
    flex-direction: column;
  }

  .filter-row {
    flex-direction: column;
  }

  .filter-item {
    min-width: 100%;
  }

  .detail-row {
    flex-direction: column;
  }
}
</style>