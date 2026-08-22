<template>
  <div class="page-container search-page">
    <!-- 搜索头部 -->
    <section class="search-hero">
      <div class="search-hero-inner">
        <div class="search-hero-icon">
          <el-icon :size="26"><Search /></el-icon>
        </div>
        <div class="search-hero-text">
          <h1 class="search-title">搜索结果</h1>
          <p class="search-subtitle">在全站中检索感兴趣的内容</p>
        </div>
      </div>

      <div class="search-box">
        <el-input
          v-model="searchQuery"
          placeholder="输入关键词搜索..."
          clearable
          size="large"
          class="search-input"
          @keyup.enter="search"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
          <template #append>
            <el-button class="search-submit-btn" type="primary" @click="search">
              搜索
            </el-button>
          </template>
        </el-input>
      </div>

      <div v-if="currentQuery" class="search-summary">
        <el-icon class="summary-icon" :size="15"><Document /></el-icon>
        <span>
          搜索
          <strong class="summary-keyword">"{{ currentQuery }}"</strong>
          找到 {{ total }} 篇文章
        </span>
      </div>
    </section>

    <!-- 结果区 -->
    <div class="search-results">
      <!-- 加载状态 -->
      <div v-if="loading" class="state-wrap">
        <el-icon :size="40" class="is-loading state-spin"><Loading /></el-icon>
        <div class="state-hint">正在搜索...</div>
      </div>

      <!-- 空状态 -->
      <div v-else-if="articles.length === 0 && currentQuery" class="state-wrap empty-state">
        <div class="empty-icon">
          <el-icon :size="64"><Search /></el-icon>
        </div>
        <div class="empty-title">没有找到相关文章</div>
        <div class="text-secondary empty-subtitle">换个关键词试试吧</div>
        <el-button round type="primary" plain class="empty-btn" @click="clearSearch">
          <el-icon class="mr-1"><RefreshLeft /></el-icon>
          重新搜索
        </el-button>
      </div>

      <!-- 结果列表 -->
      <template v-else-if="articles.length > 0">
        <ArticleList
          :articles="articles"
          :highlight-query="currentQuery"
          :show-categories="true"
        />

        <!-- 分页 -->
        <div v-if="totalPages > 1" class="pagination-bar">
          <div class="pagination-inner">
            <el-button
              @click="handlePageChange(page - 1)"
              :disabled="page === 1"
              round
              plain
              type="primary"
            >
              <el-icon class="mr-1"><ArrowLeft /></el-icon>
              上一页
            </el-button>

            <el-pagination
              :current-page="page"
              :page-count="totalPages"
              :pager-count="5"
              layout="pager"
              :hide-on-single-page="false"
              background
              @current-change="handlePageChange"
            />

            <el-button
              @click="handlePageChange(page + 1)"
              :disabled="page === totalPages"
              round
              plain
              type="primary"
            >
              下一页
              <el-icon class="ml-1"><ArrowRight /></el-icon>
            </el-button>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script>
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Search, Loading, Document, RefreshLeft, ArrowLeft, ArrowRight } from '@element-plus/icons-vue'
import ArticleList from '../components/ArticleList.vue'
import { articleApi } from '../api'

export default {
  name: 'Search',
  components: {
    ArticleList
  },
  setup() {
    const route = useRoute()
    const router = useRouter()

    const searchQuery = ref('')
    const currentQuery = ref('')
    const articles = ref([])
    const loading = ref(false)
    const total = ref(0)
    const page = ref(1)
    const pageSize = ref(20)

    const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize.value)))

    const loadArticles = async () => {
      const keyword = currentQuery.value.trim()
      if (!keyword) {
        articles.value = []
        total.value = 0
        return
      }

      loading.value = true
      try {
        const response = await articleApi.searchArticles({
          keyword,
          page: page.value,
          page_size: pageSize.value
        })
        const data = response.data || {}
        articles.value = data.articles || []
        total.value = data.total || 0
      } catch (error) {
        console.error('搜索失败', error)
      } finally {
        loading.value = false
      }
    }

    const applyUrlState = (forceReload = false) => {
      const q = typeof route.query.q === 'string' ? route.query.q : ''
      const p = parseInt(route.query.page, 10) || 1
      const qChanged = q !== currentQuery.value
      const pChanged = p !== page.value

      searchQuery.value = q
      currentQuery.value = q
      page.value = p

      if (forceReload || qChanged || pChanged) {
        loadArticles()
      }
    }

    const search = () => {
      const q = searchQuery.value.trim()
      if (!q) return

      if (q === currentQuery.value && page.value === 1) {
        loadArticles()
        return
      }

      router.replace({ path: '/search', query: { q, page: 1 } })
    }

    const handlePageChange = (newPage) => {
      if (newPage === page.value) return
      router.replace({ path: '/search', query: { q: currentQuery.value, page: newPage } })
    }

    const clearSearch = () => {
      searchQuery.value = ''
      currentQuery.value = ''
      router.replace({ path: '/search' })
    }

    watch(() => route.query, () => applyUrlState())

    onMounted(() => applyUrlState(true))

    return {
      searchQuery,
      currentQuery,
      articles,
      loading,
      total,
      totalPages,
      page,
      pageSize,
      search,
      handlePageChange,
      clearSearch
    }
  }
}
</script>

<style scoped>
.search-page {
  padding-top: 16px;
}

/* ===== 搜索头部 ===== */
.search-hero {
  position: relative;
  overflow: hidden;
  background: linear-gradient(135deg, #eef2ff 0%, #f5f0ff 55%, #fdf2ff 100%);
  border: 1px solid #e3e4ff;
  border-radius: calc(var(--campus-radius) + 4px);
  padding: 26px 28px;
  margin-bottom: 20px;
  box-shadow: var(--campus-shadow);
}

.search-hero-inner {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 18px;
}

.search-hero-icon {
  width: 54px;
  height: 54px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 16px;
  background: linear-gradient(135deg, var(--campus-primary), #8b5cf6);
  color: #fff;
  box-shadow: 0 6px 16px rgba(79, 110, 247, 0.28);
}

.search-title {
  font-size: 1.5rem;
  font-weight: 800;
  color: var(--campus-text);
  margin: 0 0 2px;
  line-height: 1.2;
}

.search-subtitle {
  font-size: 0.875rem;
  color: var(--campus-text-secondary);
  margin: 0;
}

/* 搜索框 */
.search-box {
  max-width: 640px;
}

.search-input {
  border-radius: 999px;
  overflow: hidden;
}

.search-input :deep(.el-input__wrapper) {
  border-radius: 999px;
  box-shadow: 0 2px 10px rgba(15, 23, 42, 0.06);
}

.search-input :deep(.el-input-group__append) {
  border-radius: 0 999px 999px 0;
  overflow: hidden;
  padding: 0;
  box-shadow: none;
}

.search-submit-btn {
  height: 100%;
  margin: 0;
}

/* 摘要 */
.search-summary {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.9rem;
  color: var(--campus-text-secondary);
  margin-top: 14px;
  flex-wrap: wrap;
}

.summary-icon {
  flex-shrink: 0;
  color: var(--campus-primary);
}

.summary-keyword {
  color: var(--campus-primary);
}

/* ===== 分页 ===== */
.pagination-bar {
  margin-top: 20px;
  padding: 12px;
  background: var(--campus-surface);
  border: 1px solid var(--campus-border);
  border-radius: var(--campus-radius);
  box-shadow: var(--campus-shadow);
}

.pagination-inner {
  display: flex;
  justify-content: center;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}

.pagination-inner .el-pagination {
  white-space: nowrap;
}

@media (max-width: 575px) {
  .pagination-inner {
    flex-direction: column;
    gap: 8px;
  }

  .pagination-inner .el-button {
    width: 100%;
  }
}

/* ===== 结果区 ===== */
.search-results {
  min-height: 200px;
}

.state-wrap {
  padding: 48px 0;
  text-align: center;
  background: var(--campus-surface);
  border: 1px solid var(--campus-border);
  border-radius: var(--campus-radius);
}

.state-spin {
  color: var(--campus-primary);
}

.state-hint {
  margin-top: 12px;
  font-size: 0.875rem;
  color: var(--campus-text-secondary);
}

.empty-icon {
  margin-bottom: 14px;
  color: var(--campus-border);
}

.empty-title {
  font-size: 1.15rem;
  font-weight: 700;
  color: var(--campus-text);
  margin-bottom: 6px;
}

.empty-subtitle {
  font-size: 0.875rem;
  margin-bottom: 18px;
}

.empty-btn {
  margin-top: 4px;
}

.mr-1 {
  margin-right: 4px;
}

@media (max-width: 575px) {
  .search-hero {
    padding: 20px 16px;
  }

  .search-hero-icon {
    width: 46px;
    height: 46px;
  }

  .search-title {
    font-size: 1.25rem;
  }

  .search-submit-btn {
    font-size: 0.875rem;
    padding: 0 16px;
  }
}
</style>
