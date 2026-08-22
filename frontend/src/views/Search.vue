<template>
  <div class="page-container search-page">
    <el-card class="mb-4 search-card">
      <template #header>
        <div class="d-flex align-center gap-2 search-header">
          <el-icon :size="28" class="search-header-icon"><Search /></el-icon>
          <span class="search-title">搜索结果</span>
        </div>
      </template>

      <el-input
        v-model="searchQuery"
        placeholder="输入关键词搜索..."
        clearable
        size="large"
        @keyup.enter="search"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>

      <div class="mt-4" v-if="currentQuery">
        <div class="text-secondary search-summary">
          <el-icon class="summary-icon" :size="14"><Document /></el-icon>
          搜索 <strong>"{{ currentQuery }}"</strong> 找到 {{ total }} 篇文章
        </div>
      </div>
    </el-card>

    <!-- 加载状态 -->
    <div v-if="loading" class="text-center search-loading">
      <el-icon :size="48" class="is-loading search-loading-icon"><Loading /></el-icon>
    </div>

    <!-- 空状态 -->
    <div v-else-if="articles.length === 0 && currentQuery" class="text-center search-empty">
      <el-icon :size="80" class="search-empty-icon"><Search /></el-icon>
      <div class="text-secondary search-empty-title mt-4">没有找到相关文章</div>
      <div class="text-secondary search-empty-subtitle mt-2">换个关键词试试吧</div>
    </div>

    <!-- 结果列表 -->
    <ArticleList
      v-else
      :articles="articles"
      :show-categories="true"
      :show-pagination="true"
      :total="total"
      :page="page"
      :page-size="pageSize"
      @page-change="handlePageChange"
    />
  </div>
</template>

<script>
import { ref, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Search, Loading, Document } from '@element-plus/icons-vue'
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

    watch(() => route.query, () => applyUrlState())

    onMounted(() => applyUrlState(true))

    return {
      searchQuery,
      currentQuery,
      articles,
      loading,
      total,
      page,
      pageSize,
      search,
      handlePageChange
    }
  }
}
</script>

<style scoped>
.search-page {
  padding-top: 16px;
}

.search-card {
  border-radius: var(--campus-radius);
}

.search-header {
  font-size: 1.25rem;
  font-weight: 600;
}

.search-header-icon {
  color: var(--campus-primary);
}

.search-title {
  line-height: 1.3;
}

.search-summary {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.875rem;
}

.summary-icon {
  flex-shrink: 0;
}

.search-loading {
  padding: 40px 0;
}

.search-loading-icon {
  color: var(--campus-primary);
}

.search-empty {
  padding: 48px 0;
}

.search-empty-icon {
  color: var(--campus-border);
}

.search-empty-title {
  font-size: 1.25rem;
  font-weight: 500;
}

.search-empty-subtitle {
  font-size: 0.875rem;
  opacity: 0.75;
}
</style>
