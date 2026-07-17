<template>
  <div>
    <div class="layui-card mb-4">
      <div class="layui-card-body">
        <div class="search-bar">
          <div class="search-input-wrap">
            <i class="fa-solid fa-magnifying-glass"></i>
            <input type="text" v-model="searchQuery" placeholder="搜索文章标题..." class="layui-input" />
          </div>
          <select v-model="localFilter" @change="$emit('update:filter', localFilter)" class="layui-select">
            <option v-for="option in statusOptions" :key="option.value" :value="option.value">
              {{ option.label }}
            </option>
          </select>
          <button class="layui-btn" @click="$emit('refresh')" :disabled="loading">
            <i class="fa-solid fa-rotate-right"></i>
            刷新
          </button>
        </div>
      </div>
    </div>

    <div class="layui-card">
      <div v-if="filteredArticles.length > 0" class="articles-list">
        <div v-for="article in filteredArticles" :key="article.id" class="article-item">
          <div class="article-icon">
            <i class="fa-solid fa-file-lines"></i>
          </div>

          <div class="article-info">
            <div class="article-title-wrap">
              <span class="article-title font-weight-medium">{{ article.title }}</span>
              <span :class="['status-tag', getStatusClass(article.status)]">{{ getStatusText(article.status) }}</span>
            </div>
            <div class="article-meta">
              <span class="meta-item"><i class="fa-solid fa-hashtag"></i> ID: {{ article.id }}</span>
              <span class="meta-item">
                <UserAvatar :user="article.User || {}" :size="20" />
                {{ article.User?.display_name || '-' }}
              </span>
              <span v-if="article.Category" class="meta-item category-tag">{{ article.Category.name }}</span>
            </div>
            <div class="article-stats">
              <span class="stat-item"><i class="fa-solid fa-heart" style="color: #FF5722;"></i> {{ article.like_count || 0 }}</span>
              <span class="stat-item"><i class="fa-solid fa-eye" style="color: #1E9FFF;"></i> {{ article.view_count || 0 }}</span>
              <span class="stat-item"><i class="fa-solid fa-clock"></i> {{ formatDate(article.created_at) }}</span>
            </div>
          </div>

          <div class="article-actions">
            <button v-if="article.status !== 'deleted'" class="action-btn" @click="$emit('change-status', article)" title="修改状态">
              <i class="fa-solid fa-refresh-cw"></i>
            </button>
            <button v-if="article.status === 'deleted'" class="action-btn success" @click="$emit('restore', article)" title="恢复文章">
              <i class="fa-solid fa-rotate-left"></i>
            </button>
            <router-link :to="`/article/${article.id}`" target="_blank" class="action-btn" title="查看文章">
              <i class="fa-solid fa-eye"></i>
            </router-link>
            <button v-if="canDeleteArticle() && article.status !== 'deleted'" class="action-btn danger" @click="$emit('delete', article)" title="删除">
              <i class="fa-solid fa-trash"></i>
            </button>
          </div>
        </div>
      </div>

      <div v-else class="empty-state">
        <i class="fa-solid fa-file-lines" style="font-size: 48px; color: #ccc;"></i>
        <div class="text-body-1 text-medium-emphasis mt-2">
          {{ searchQuery ? '未找到匹配的文章' : '暂无文章数据' }}
        </div>
      </div>
    </div>

    <div v-if="totalPages > 1" class="pagination-wrap">
      <div class="layui-laypage">
        <button class="layui-laypage-prev" :disabled="page <= 1" @click="handlePrevPage">上一页</button>
        <span v-for="p in pageRange" :key="p" :class="['layui-laypage-curr', { active: p === page }]" @click="$emit('update:page', p)">
          {{ p }}
        </span>
        <button class="layui-laypage-next" :disabled="page >= totalPages" @click="handleNextPage">下一页</button>
      </div>
      <div class="page-info">第 {{ page }} / {{ totalPages }} 页</div>
    </div>
  </div>
</template>

<script>
import { ref, watch, computed } from 'vue'
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

    const pageRange = computed(() => {
      const range = []
      const total = props.totalPages
      const current = props.page
      const visible = 5
      
      let start = Math.max(1, current - Math.floor(visible / 2))
      let end = Math.min(total, start + visible - 1)
      
      if (end - start + 1 < visible) {
        start = Math.max(1, end - visible + 1)
      }
      
      for (let i = start; i <= end; i++) {
        range.push(i)
      }
      return range
    })

    const getStatusClass = (status) => {
      const classes = {
        pending: 'pending',
        published: 'published',
        rejected: 'rejected'
      }
      return classes[status] || ''
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
      return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit'
      })
    }

    const canDeleteArticle = () => {
      return props.currentUserRole === 'system' || props.currentUserRole === 'admin'
    }

    const handlePrevPage = () => {
      if (props.page > 1) {
        emit('update:page', props.page - 1)
      }
    }

    const handleNextPage = () => {
      if (props.page < props.totalPages) {
        emit('update:page', props.page + 1)
      }
    }

    return {
      localPage,
      localFilter,
      searchQuery,
      filteredArticles,
      pageRange,
      getStatusClass,
      getStatusText,
      formatDate,
      canDeleteArticle,
      handlePrevPage,
      handleNextPage
    }
  }
}
</script>

<style scoped>
.search-bar {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.search-input-wrap {
  flex: 1;
  min-width: 200px;
  position: relative;
}

.search-input-wrap i {
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: #999;
}

.search-input-wrap input {
  padding-left: 36px;
}

.layui-select {
  min-width: 150px;
}

.mb-4 {
  margin-bottom: 16px;
}

.articles-list {
  padding: 0;
}

.article-item {
  display: flex;
  align-items: flex-start;
  padding: 16px;
  border-bottom: 1px solid #f0f0f0;
  gap: 12px;
}

.article-item:last-child {
  border-bottom: none;
}

.article-icon {
  width: 48px;
  height: 48px;
  background: #e6f7ff;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #1E9FFF;
  flex-shrink: 0;
}

.article-info {
  flex: 1;
  min-width: 0;
}

.article-title-wrap {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.article-title {
  font-weight: 500;
  font-size: 16px;
}

.status-tag {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.status-tag.pending {
  background: #fffbe6;
  color: #faad14;
}

.status-tag.published {
  background: #f6ffed;
  color: #52c41a;
}

.status-tag.rejected {
  background: #fff2f0;
  color: #ff4d4f;
}

.category-tag {
  background: #e6f7ff;
  color: #1890ff;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.article-meta {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
  margin-top: 8px;
  font-size: 13px;
  color: #666;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.meta-item i {
  font-size: 12px;
}

.article-stats {
  display: flex;
  gap: 16px;
  margin-top: 8px;
  font-size: 13px;
  color: #999;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.stat-item i {
  font-size: 12px;
}

.article-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  width: 32px;
  height: 32px;
  border: none;
  background: transparent;
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #999;
  transition: all 0.3s;
}

.action-btn:hover {
  background: #f0f0f0;
  color: #1E9FFF;
}

.action-btn.success:hover {
  color: #52C41A;
}

.action-btn.danger:hover {
  background: #fff2f0;
  color: #FF5722;
}

.empty-state {
  text-align: center;
  padding: 32px;
}

.text-body-1 {
  font-size: 16px;
}

.text-medium-emphasis {
  color: #999;
}

.font-weight-medium {
  font-weight: 500;
}

.mt-2 {
  margin-top: 8px;
}

.pagination-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  margin-top: 16px;
}

.page-info {
  font-size: 13px;
  color: #999;
}
</style>