<template>
  <div class="articles-panel">
    <!-- 搜索和筛选栏 -->
    <el-card class="mb-4" shadow="never" body-class="articles-filter">
      <el-row :gutter="12" align="middle">
        <el-col :xs="24" :sm="12" :md="8">
          <el-input
            v-model="searchQuery"
            placeholder="搜索文章标题..."
            clearable
            class="articles-search"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-col>
        <el-col :xs="24" :sm="12" :md="8">
          <el-select
            v-model="localFilter"
            style="width: 100%"
            placeholder="状态筛选"
            @update:model-value="$emit('update:filter', $event)"
          >
            <el-option
              v-for="opt in statusOptions"
              :key="opt.value"
              :label="opt.label"
              :value="opt.value"
            />
          </el-select>
        </el-col>
        <el-col :xs="24" :sm="24" :md="8" class="articles-actions">
          <el-button type="primary" @click="$emit('refresh')" :loading="loading">
            <el-icon class="btn-icon"><Refresh /></el-icon>
            刷新
          </el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- 文章列表 -->
    <el-card shadow="never" body-class="articles-list-wrap">
      <el-empty
        v-if="filteredArticles.length === 0"
        :description="searchQuery ? '未找到匹配的文章' : '暂无文章数据'"
      >
        <template #image>
          <el-icon :size="56" class="empty-icon"><Document /></el-icon>
        </template>
      </el-empty>

      <template v-else>
        <div v-for="article in filteredArticles" :key="article.id" class="article-item">
          <div class="article-item-main">
            <div class="article-avatar">
              <el-icon :size="22"><Document /></el-icon>
            </div>
            <div class="article-content">
              <div class="article-title-row">
                <span class="article-title">{{ article.title }}</span>
                <el-tag
                  size="small"
                  :type="getStatusType(article.status)"
                  effect="light"
                  class="ml-2"
                >
                  {{ getStatusText(article.status) }}
                </el-tag>
              </div>
              <div class="article-subtitle">
                <span class="meta-item">
                  <el-icon :size="14" class="meta-icon"><Key /></el-icon>
                  ID: {{ article.id }}
                </span>
                <span class="meta-item">
                  <UserAvatar :user="article.User || {}" :size="20" class="meta-avatar" />
                  {{ article.User?.display_name || '-' }}
                </span>
                <el-tag size="small" type="primary" effect="light" v-if="article.Category">
                  {{ article.Category.name }}
                </el-tag>
              </div>
              <div class="article-metrics">
                <span class="metric-item">
                  <el-icon :size="14" class="metric-icon metric-pink"><Star /></el-icon>
                  {{ article.like_count || 0 }}
                </span>
                <span class="metric-item">
                  <el-icon :size="14" class="metric-icon metric-blue"><View /></el-icon>
                  {{ article.view_count || 0 }}
                </span>
                <span class="metric-item">
                  <el-icon :size="14" class="metric-icon metric-grey"><Clock /></el-icon>
                  {{ formatDate(article.created_at) }}
                </span>
              </div>
            </div>
          </div>
          <div class="article-actions">
            <el-tooltip content="修改状态" placement="top">
              <el-button
                size="small"
                type="primary"
                link
                @click="$emit('change-status', article)"
                v-if="article.status !== 'deleted'"
              >
                <el-icon><Refresh /></el-icon>
                <span>修改状态</span>
              </el-button>
            </el-tooltip>
            <el-tooltip content="恢复文章" placement="top">
              <el-button
                size="small"
                type="success"
                link
                @click="$emit('restore', article)"
                v-if="article.status === 'deleted'"
              >
                <el-icon><RefreshLeft /></el-icon>
                <span>恢复</span>
              </el-button>
            </el-tooltip>
            <el-tooltip content="查看文章" placement="top">
              <el-button
                size="small"
                type="primary"
                link
                :href="`/article/${article.id}`"
                target="_blank"
                tag="a"
              >
                <el-icon><View /></el-icon>
                <span>查看</span>
              </el-button>
            </el-tooltip>
            <el-tooltip content="删除" placement="top">
              <el-button
                size="small"
                type="danger"
                link
                @click="$emit('delete', article)"
                v-if="canDeleteArticle() && article.status !== 'deleted'"
              >
                <el-icon><Delete /></el-icon>
                <span>删除</span>
              </el-button>
            </el-tooltip>
          </div>
        </div>
      </template>
    </el-card>

    <!-- 分页 -->
    <div class="articles-pagination" v-if="totalPages > 1">
      <el-pagination
        v-model:current-page="localPage"
        :page-size="1"
        :total="totalPages"
        :pager-count="5"
        layout="pager, total"
        @current-change="$emit('update:page', $event)"
      />
      <div class="pagination-text">第 {{ page }} / {{ totalPages }} 页</div>
    </div>
  </div>
</template>

<script>
import { ref, watch, computed } from 'vue'
import {
  Search,
  Refresh,
  Document,
  Key,
  View,
  Star,
  Clock,
  RefreshLeft,
  Delete
} from '@element-plus/icons-vue'
import UserAvatar from '../../components/UserAvatar.vue'

export default {
  name: 'ArticlesPanel',
  components: {
    UserAvatar
  },
  props: {
    articles: {
      type: Array,
      default: () => []
    },
    loading: {
      type: Boolean,
      default: false
    },
    page: {
      type: Number,
      default: 1
    },
    totalPages: {
      type: Number,
      default: 1
    },
    filter: {
      type: String,
      default: ''
    },
    statusOptions: {
      type: Array,
      default: () => []
    },
    currentUserRole: {
      type: String,
      default: null
    }
  },
  emits: ['change-status', 'delete', 'restore', 'refresh', 'update:page', 'update:filter'],
  setup(props, { emit }) {
    const localPage = ref(props.page)
    const localFilter = ref(props.filter)
    const searchQuery = ref('')

    watch(() => props.page, (val) => {
      localPage.value = val
    })

    watch(() => props.filter, (val) => {
      localFilter.value = val
    })

    const filteredArticles = computed(() => {
      if (!searchQuery.value) return props.articles
      const query = searchQuery.value.toLowerCase()
      return props.articles.filter(article =>
        article.title?.toLowerCase().includes(query) ||
        article.User?.display_name?.toLowerCase().includes(query)
      )
    })

    const getStatusType = (status) => {
      const types = {
        pending: 'warning',
        published: 'success',
        rejected: 'danger'
      }
      return types[status] || 'info'
    }

    const getStatusText = (status) => {
      const texts = {
        pending: '待审核',
        published: '已发布',
        rejected: '已拒绝'
      }
      return texts[status] || status
    }

    const formatDate = (dateString) => {
      if (!dateString) return '-'
      const date = new Date(dateString)
      return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit'
      })
    }

    const canDeleteArticle = () => {
      return props.currentUserRole === 'system' || props.currentUserRole === 'admin'
    }

    return {
      localPage,
      localFilter,
      searchQuery,
      filteredArticles,
      getStatusType,
      getStatusText,
      formatDate,
      canDeleteArticle
    }
  }
}
</script>

<style scoped>
.articles-filter {
  display: flex;
  align-items: center;
}
.articles-search {
  width: 100%;
}
.articles-actions {
  display: flex;
  justify-content: flex-end;
}
.btn-icon {
  margin-right: 4px;
}
.articles-list-wrap {
  padding-top: 4px;
}
.empty-icon {
  color: #c0c4cc;
}
.article-item {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  padding: 12px 4px;
  border-bottom: 1px solid #f0f0f0;
}
.article-item:last-child {
  border-bottom: none;
}
.article-item-main {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  min-width: 0;
}
.article-avatar {
  flex-shrink: 0;
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(103, 80, 164, 0.1);
  color: #6750a4;
}
.article-content {
  min-width: 0;
}
.article-title-row {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  font-weight: 500;
  margin-bottom: 4px;
}
.article-title {
  color: #1a1a2e;
}
.ml-2 {
  margin-left: 8px;
}
.article-subtitle {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 12px;
  color: #6b7280;
  font-size: 13px;
}
.meta-item {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}
.meta-avatar {
  margin-right: 2px;
}
.meta-icon {
  color: #6b7280;
}
.article-metrics {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 16px;
  margin-top: 6px;
  font-size: 12px;
  color: #909399;
}
.metric-item {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}
.metric-icon {
  margin-right: 2px;
}
.metric-pink {
  color: #f56c6c;
}
.metric-blue {
  color: #409eff;
}
.metric-grey {
  color: #909399;
}
.article-actions {
  display: flex;
  align-items: center;
  flex-shrink: 0;
  gap: 4px;
}
.article-actions .el-button + .el-button {
  margin-left: 0;
}
.article-actions .el-button > span {
  display: inline-flex;
  align-items: center;
  gap: 2px;
}
.articles-pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  margin-top: 16px;
}
.pagination-text {
  font-size: 13px;
  color: #909399;
}
</style>
