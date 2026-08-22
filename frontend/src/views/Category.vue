<template>
  <div class="grid-layout">
    <div class="sidebar">
      <Sidebar />
    </div>
    
    <div class="main-content">
      <el-card class="category-card mb-4">
        <span class="category-title">{{ categoryName }}</span>
      </el-card>
      
      <ArticleList :articles="articles" />
      
      <div class="pagination-bar" v-if="totalPages > 1">
        <el-button
          @click="prevPage"
          :disabled="page === 1"
          plain
        >
          上一页
        </el-button>
        <span class="page-info">第 {{ page }} / {{ totalPages }} 页</span>
        <el-button
          @click="nextPage"
          :disabled="page === totalPages"
          plain
        >
          下一页
        </el-button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { articleApi } from '../api'
import Sidebar from '../components/Sidebar.vue'
import ArticleList from '../components/ArticleList.vue'

export default {
  name: 'Category',
  components: {
    Sidebar,
    ArticleList
  },
  setup() {
    const route = useRoute()
    const articles = ref([])
    const categoryName = ref('')
    const page = ref(1)
    const totalPages = ref(1)

    const loadArticles = async () => {
      try {
        const response = await articleApi.getArticles({
          category_id: route.params.id,
          page: page.value,
          page_size: 20
        })
        articles.value = response.data.articles
        totalPages.value = response.data.total_pages

        if (articles.value.length > 0 && articles.value[0].category) {
          categoryName.value = articles.value[0].category.name
        }
      } catch (error) {
        console.error('加载文章失败', error)
      }
    }

    const prevPage = () => {
      if (page.value > 1) {
        page.value--
        loadArticles()
      }
    }

    const nextPage = () => {
      if (page.value < totalPages.value) {
        page.value++
        loadArticles()
      }
    }

    watch(() => route.params.id, () => {
      page.value = 1
      categoryName.value = ''
      loadArticles()
    })

    onMounted(() => {
      loadArticles()
    })

    return {
      articles,
      categoryName,
      page,
      totalPages,
      loadArticles,
      prevPage,
      nextPage
    }
  }
}
</script>

<style scoped>
.main-content {
  min-width: 0;
}

.category-card {
  width: 100%;
  margin-bottom: 1rem;
  --el-card-padding: 24px;
}

.category-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--campus-primary);
}

.pagination-bar {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 16px;
  margin-top: 1rem;
}

.page-info {
  font-size: 0.875rem;
  color: var(--campus-text-secondary);
}
</style>
