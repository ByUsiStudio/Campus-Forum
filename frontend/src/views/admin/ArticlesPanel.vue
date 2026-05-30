<template>
  <div class="articles-panel">
    <div class="panel-header">
      <div class="header-left">
        <h2 class="panel-title">文章管理</h2>
        <p class="panel-subtitle">审核与管理用户发布的文章</p>
      </div>
      <div class="header-actions">
        <v-select
          v-model="localFilter"
          :items="statusOptions"
          label="状态筛选"
          variant="outlined"
          density="compact"
          hide-details
          class="filter-select"
          @update:model-value="$emit('update:filter', $event)"
        ></v-select>
        <v-btn variant="outlined" color="primary" @click="$emit('refresh')" :loading="loading">
          <v-icon start>mdi-refresh</v-icon>
          刷新
        </v-btn>
      </div>
    </div>

    <v-card class="table-card">
      <v-table class="articles-table">
        <thead>
          <tr>
            <th class="text-left">ID</th>
            <th class="text-left">标题</th>
            <th class="text-left">作者</th>
            <th class="text-left">分区</th>
            <th class="text-center">点赞</th>
            <th class="text-center">浏览</th>
            <th class="text-left">状态</th>
            <th class="text-left">发布时间</th>
            <th class="text-center">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="articles.length === 0">
            <td colspan="9" class="text-center pa-8">
              <v-icon size="48" color="grey">mdi-file-document-search</v-icon>
              <div class="mt-2 text-grey">暂无文章数据</div>
            </td>
          </tr>
          <tr v-for="article in articles" :key="article.id" class="article-row">
            <td class="id-cell">{{ article.id }}</td>
            <td class="title-cell">
              <div class="article-title">{{ article.title }}</div>
            </td>
            <td>
              <div class="author-cell">
                <UserAvatar :user="article.User || {}" :size="28" />
                <span class="author-name">{{ article.User?.display_name || '-' }}</span>
              </div>
            </td>
            <td>
              <v-chip size="small" variant="tonal" color="primary">
                {{ article.Category?.name || '-' }}
              </v-chip>
            </td>
            <td class="text-center count-cell">
              <v-icon size="16" color="pink">mdi-heart</v-icon>
              {{ article.like_count || 0 }}
            </td>
            <td class="text-center count-cell">
              <v-icon size="16" color="blue">mdi-eye</v-icon>
              {{ article.view_count || 0 }}
            </td>
            <td>
              <v-chip size="small" :color="getStatusColor(article.status)" variant="tonal">
                <v-icon start size="x-small">{{ getStatusIcon(article.status) }}</v-icon>
                {{ getStatusText(article.status) }}
              </v-chip>
            </td>
            <td class="date-cell">{{ formatDate(article.created_at) }}</td>
            <td>
              <div class="action-cell">
                <v-btn variant="text" size="small" color="primary" @click="$emit('change-status', article)">
                  改状态
                </v-btn>
                <v-btn
                  variant="text"
                  size="small"
                  color="error"
                  @click="$emit('delete', article)"
                  v-if="canDeleteArticle()"
                >
                  删除
                </v-btn>
              </div>
            </td>
          </tr>
        </tbody>
      </v-table>
    </v-card>

    <div class="pagination-wrapper" v-if="totalPages > 1">
      <v-pagination
        v-model="localPage"
        :length="totalPages"
        :total-visible="5"
        rounded="lg"
        @update:model-value="$emit('update:page', $event)"
      ></v-pagination>
      <div class="page-info">第 {{ page }} / {{ totalPages }} 页</div>
    </div>
  </div>
</template>

<script>
import { ref, watch } from 'vue'
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
  emits: ['change-status', 'delete', 'refresh', 'update:page', 'update:filter'],
  setup(props, { emit }) {
    const localPage = ref(props.page)
    const localFilter = ref(props.filter)

    watch(() => props.page, (val) => {
      localPage.value = val
    })

    watch(() => props.filter, (val) => {
      localFilter.value = val
    })

    const getStatusColor = (status) => {
      const colors = {
        pending: 'warning',
        published: 'success',
        rejected: 'error'
      }
      return colors[status] || 'default'
    }

    const getStatusIcon = (status) => {
      const icons = {
        pending: 'mdi-clock-outline',
        published: 'mdi-check-circle',
        rejected: 'mdi-close-circle'
      }
      return icons[status] || 'mdi-help-circle'
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
      getStatusColor,
      getStatusIcon,
      getStatusText,
      formatDate,
      canDeleteArticle
    }
  }
}
</script>

<style scoped>
.articles-panel {
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 16px;
}

.panel-title {
  font-size: 1.5rem;
  font-weight: 700;
  color: #1a1a2e;
  margin: 0 0 4px 0;
}

.panel-subtitle {
  font-size: 0.9rem;
  color: #6b7280;
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 12px;
  align-items: center;
  flex-wrap: wrap;
}

.filter-select {
  width: 140px;
}

.table-card {
  border-radius: 16px;
  overflow: hidden;
}

.articles-table {
  width: 100%;
}

.articles-table th {
  font-weight: 600;
  color: #6b7280;
  font-size: 0.85rem;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  background: #f8f9ff !important;
  padding: 16px 12px !important;
}

.articles-table td {
  padding: 14px 12px !important;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.article-row:hover {
  background: rgba(103, 80, 164, 0.03);
}

.id-cell {
  font-family: monospace;
  color: #9ca3af;
  font-size: 0.85rem;
}

.title-cell {
  max-width: 240px;
}

.article-title {
  font-weight: 500;
  color: #1a1a2e;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.author-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.author-name {
  font-size: 0.9rem;
  color: #4b5563;
}

.count-cell {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  font-size: 0.9rem;
  color: #6b7280;
}

.date-cell {
  color: #6b7280;
  font-size: 0.85rem;
  white-space: nowrap;
}

.action-cell {
  display: flex;
  justify-content: center;
  gap: 4px;
}

.pagination-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  margin-top: 24px;
}

.page-info {
  font-size: 0.9rem;
  color: #6b7280;
}

@media (max-width: 960px) {
  .articles-table {
    font-size: 0.8rem;
  }

  .articles-table th,
  .articles-table td {
    padding: 10px 6px !important;
  }

  .author-cell,
  .date-cell {
    display: none;
  }

  .title-cell {
    max-width: 160px;
  }
}
</style>
