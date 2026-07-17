<template>
  <div class="search-page">
    <div class="search-card">
      <div class="card-header">
        <i class="fa-solid fa-magnifying-glass mr-2"></i>
        <span>搜索结果</span>
      </div>
      <div class="card-body">
        <div class="search-input-wrapper">
          <i class="fa-solid fa-magnifying-glass"></i>
          <input 
            v-model="searchQuery"
            type="text"
            placeholder="输入关键词搜索..."
            class="search-input"
            @keyup.enter="search"
          />
          <button v-if="searchQuery" class="clear-btn" @click="searchQuery = ''">
            <i class="fa-solid fa-xmark"></i>
          </button>
        </div>
        <div v-if="currentQuery" class="search-info">
          <i class="fa-solid fa-file-lines mr-1"></i>
          搜索 <strong>"{{ currentQuery }}"</strong> 找到 {{ total }} 篇文章
        </div>
      </div>
    </div>

    <div v-if="loading" class="loading-state">
      <div class="loading-spinner"></div>
    </div>

    <div v-else-if="articles.length === 0 && currentQuery" class="empty-state">
      <i class="fa-solid fa-search-slash"></i>
      <div class="empty-title">没有找到相关文章</div>
      <div class="empty-desc">换个关键词试试吧</div>
    </div>

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
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
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

    const search = async () => {
      if (!searchQuery.value.trim()) {
        return
      }

      currentQuery.value = searchQuery.value.trim()
      page.value = 1
      router.replace({ path: '/search', query: { q: currentQuery.value, page: 1 } })
      await loadArticles()
    }

    const loadArticles = async () => {
      loading.value = true
      try {
        const response = await articleApi.searchArticles({
          keyword: currentQuery.value,
          page: page.value,
          page_size: pageSize.value
        })
        articles.value = response.data.articles || []
        total.value = response.data.total || 0
      } catch (error) {
        console.error('搜索失败', error)
      } finally {
        loading.value = false
      }
    }

    const handlePageChange = (newPage) => {
      page.value = newPage
      router.replace({ path: '/search', query: { q: currentQuery.value, page: newPage } })
      loadArticles()
    }

    onMounted(() => {
      if (route.query.q) {
        searchQuery.value = route.query.q
        currentQuery.value = route.query.q
        if (route.query.page) {
          page.value = parseInt(route.query.page) || 1
        }
        loadArticles()
      }
    })

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
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;
}

.search-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
  margin-bottom: 24px;
}

.card-header {
  display: flex;
  align-items: center;
  padding: 16px 24px;
  border-bottom: 1px solid #f0f0f0;
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.card-body {
  padding: 24px;
}

.search-input-wrapper {
  position: relative;
}

.search-input-wrapper i {
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 16px;
  color: #999;
}

.search-input {
  width: 100%;
  padding: 12px 40px;
  border: 2px solid #e8e8e8;
  border-radius: 8px;
  font-size: 14px;
  outline: none;
  
  &:focus {
    border-color: var(--primary);
  }
}

.clear-btn {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  font-size: 16px;
  color: #999;
  cursor: pointer;
  
  &:hover {
    color: #666;
  }
}

.search-info {
  margin-top: 12px;
  font-size: 14px;
  color: #666;
}

.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48px;
}

.loading-spinner {
  width: 48px;
  height: 48px;
  border: 4px solid #f0f0f0;
  border-top-color: var(--primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.empty-state {
  text-align: center;
  padding: 48px 24px;
  
  i {
    font-size: 64px;
    color: #e8e8e8;
    margin-bottom: 16px;
  }
}

.empty-title {
  font-size: 16px;
  color: #666;
  margin-bottom: 8px;
}

.empty-desc {
  font-size: 14px;
  color: #999;
}

.mr-1 {
  margin-right: 4px;
}

.mr-2 {
  margin-right: 8px;
}
</style>
