<template>
  <div class="page-container">
    <el-row :gutter="16">
      <!-- 侧边栏 - 桌面端显示 -->
      <el-col :xs="0" :sm="0" :md="6" :lg="6" class="sidebar-col">
        <Sidebar />
      </el-col>

      <!-- 主内容 -->
      <el-col :xs="24" :sm="24" :md="18" :lg="18">
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
      </el-col>
    </el-row>
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
.sidebar-col {
  display: block;
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

@media (max-width: 991px) {
  .sidebar-col {
    display: none;
  }
}
</style>
