<template>
  <div>
    <!-- 搜索栏 -->
    <div class="search-bar mb-4">
      <el-input
        v-model="searchQuery"
        placeholder="搜索评论内容或用户名..."
        clearable
        class="search-input"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
      <el-button type="primary" :loading="loading" @click="$emit('refresh')">
        <el-icon style="margin-right: 4px"><Refresh /></el-icon>
        刷新
      </el-button>
    </div>

    <!-- 评论列表 -->
    <el-table
      :data="filteredComments"
      v-loading="loading"
      empty-text="暂无评论数据"
    >
      <el-table-column label="评论用户" min-width="160">
        <template #default="{ row }">
          <div class="user-cell">
            <UserAvatar :user="getUserInfo(row)" :size="36" :show-username="false" />
            <span class="user-name">{{ getUserName(row) }}</span>
            <el-tag
              v-if="row.is_anonymous"
              size="small"
              type="info"
              effect="plain"
              class="anon-tag"
            >
              匿名
            </el-tag>
          </div>
        </template>
      </el-table-column>

      <el-table-column label="评论内容" min-width="240" show-overflow-tooltip>
        <template #default="{ row }">
          <span>{{ row.content }}</span>
        </template>
      </el-table-column>

      <el-table-column label="所属文章" min-width="180" show-overflow-tooltip>
        <template #default="{ row }">
          <el-link
            v-if="getArticleInfo(row)"
            type="primary"
            :href="`/article/${row.article_id}`"
            target="_blank"
          >
            <el-icon style="margin-right: 4px"><Document /></el-icon>
            {{ getArticleInfo(row).title }}
          </el-link>
          <span v-else>-</span>
        </template>
      </el-table-column>

      <el-table-column label="互动" width="130">
        <template #default="{ row }">
          <div class="stats-cell">
            <el-icon color="#f56c6c"><ChatDotRound /></el-icon>
            <span class="stat-num">{{ row.reply_count || 0 }}</span>
            <el-icon color="#409eff"><Pointer /></el-icon>
            <span class="stat-num">{{ row.like_count || 0 }}</span>
          </div>
        </template>
      </el-table-column>

      <el-table-column label="评论时间" width="170">
        <template #default="{ row }">
          <el-icon style="margin-right: 4px; vertical-align: middle"><Clock /></el-icon>
          <span>{{ formatDate(row.created_at) }}</span>
        </template>
      </el-table-column>

      <el-table-column label="操作" width="130" fixed="right">
        <template #default="{ row }">
          <el-button
            v-if="row.article_id"
            size="small"
            type="info"
            plain
            @click="openArticle(row)"
          >
            查看文章
          </el-button>
          <el-button size="small" type="danger" plain @click="$emit('delete', row.id)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-empty
      v-if="filteredComments.length === 0"
      :description="searchQuery ? '未找到匹配的评论' : '暂无评论数据'"
      :image-size="80"
    />

    <!-- 分页 -->
    <div class="pagination-row" v-if="pagination.totalPages > 1">
      <el-pagination
        :current-page="currentPage"
        :page-size="pagination.pageSize"
        :total="pagination.total"
        layout="prev, pager, next, total"
        background
        @current-change="handlePageChange"
      />
      <span class="page-text">
        第 {{ pagination.page }} / {{ pagination.totalPages }} 页 (共 {{ pagination.total }} 条)
      </span>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import {
  Search,
  Refresh,
  Document,
  ChatDotRound,
  Pointer,
  Clock
} from '@element-plus/icons-vue'
import UserAvatar from '../../components/UserAvatar.vue'

const props = defineProps({
  comments: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  },
  pagination: {
    type: Object,
    default: () => ({
      page: 1,
      pageSize: 20,
      total: 0,
      totalPages: 0
    })
  }
})

const emit = defineEmits(['delete', 'refresh', 'page-change'])

const currentPage = ref(props.pagination.page || 1)
const searchQuery = ref('')

watch(
  () => props.pagination.page,
  (newVal) => {
    currentPage.value = newVal
  }
)

const filteredComments = computed(() => {
  if (!searchQuery.value) return props.comments
  const query = searchQuery.value.toLowerCase()
  return props.comments.filter(
    (comment) =>
      comment.content?.toLowerCase().includes(query) ||
      getUserName(comment).toLowerCase().includes(query)
  )
})

const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getUserInfo = (comment) => comment.user || comment.User || {}

const getUserName = (comment) => {
  const user = getUserInfo(comment)
  return (
    user.display_name ||
    user.DisplayName ||
    user.username ||
    user.Username ||
    '未知用户'
  )
}

const getArticleInfo = (comment) => comment.article || comment.Article || null

const handlePageChange = (page) => {
  emit('page-change', page)
}

const openArticle = (row) => {
  if (row.article_id) {
    window.open(`/article/${row.article_id}`, '_blank')
  }
}
</script>

<style scoped>
.search-bar {
  display: flex;
  align-items: center;
  gap: 12px;
}

.search-input {
  max-width: 360px;
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.user-name {
  font-weight: 500;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.anon-tag {
  flex-shrink: 0;
}

.stats-cell {
  display: flex;
  align-items: center;
  gap: 6px;
}

.stat-num {
  margin-right: 8px;
  color: #606266;
}

.pagination-row {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  margin-top: 16px;
}

.page-text {
  font-size: 12px;
  color: #909399;
}
</style>
