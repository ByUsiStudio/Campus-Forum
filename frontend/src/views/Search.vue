<template>
  <v-container fluid class="py-6">
    <v-card class="mb-6">
      <v-card-title class="d-flex align-center gap-2">
        <v-icon size="28">mdi-magnify</v-icon>
        <span class="text-h5">搜索结果</span>
      </v-card-title>
      <v-card-text>
        <v-text-field
          v-model="searchQuery"
          prepend-inner-icon="mdi-magnify"
          placeholder="输入关键词搜索..."
          variant="outlined"
          @keyup.enter="search"
          clearable
        ></v-text-field>
        <div class="mt-4" v-if="currentQuery">
          <div class="text-subtitle-2 text-medium-emphasis">
            <v-icon class="mr-1" size="small">mdi-file-document</v-icon>
            搜索 <strong>"{{ currentQuery }}"</strong> 找到 {{ total }} 篇文章
          </div>
        </div>
      </v-card-text>
    </v-card>

    <div v-if="loading" class="text-center py-8">
      <v-progress-circular indeterminate size="48" color="primary"></v-progress-circular>
    </div>

    <div v-else-if="articles.length === 0 && currentQuery" class="text-center py-12">
      <v-icon size="80" class="text-grey-lighten-2">mdi-file-search</v-icon>
      <div class="text-h5 text-medium-emphasis mt-4">没有找到相关文章</div>
      <div class="text-subtitle-2 text-grey-lighten-1 mt-2">换个关键词试试吧</div>
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
  </v-container>
</template>

<script>
import { ref, watch, onMounted } from 'vue'
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
