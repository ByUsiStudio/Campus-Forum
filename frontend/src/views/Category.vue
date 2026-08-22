<template>
  <div class="page-container category-page">
    <div class="category-layout">
      <!-- 主内容 -->
      <main class="main-column">
        <!-- 分类头部 -->
        <section class="category-header">
          <div class="category-header-icon">
            <el-icon :size="28"><FolderOpened /></el-icon>
          </div>
          <div class="category-header-text">
            <h1 class="category-title">{{ categoryName || '分类文章' }}</h1>
            <p class="category-subtitle">共 {{ totalPages }} 页 · 精选该分区的优质内容</p>
          </div>
        </section>

        <ArticleList :articles="articles" />

        <div class="pagination-bar" v-if="totalPages > 1">
          <el-button
            @click="prevPage"
            :disabled="page === 1"
            round
            plain
            type="primary"
            class="page-btn"
          >
            <el-icon class="mr-1"><ArrowLeft /></el-icon>
            上一页
          </el-button>
          <span class="page-info">第 {{ page }} / {{ totalPages }} 页</span>
          <el-button
            @click="nextPage"
            :disabled="page === totalPages"
            round
            plain
            type="primary"
            class="page-btn"
          >
            下一页
            <el-icon class="ml-1"><ArrowRight /></el-icon>
          </el-button>
        </div>
      </main>

      <!-- 侧边栏 - 桌面端显示 -->
      <aside class="sidebar-column">
        <div class="sidebar-sticky">
          <Sidebar />
        </div>
      </aside>
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
.category-page {
  padding-top: 16px;
}

.category-layout {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 300px;
  gap: 20px;
  align-items: start;
}

.main-column {
  min-width: 0;
}

.sidebar-column {
  min-width: 0;
}

.sidebar-sticky {
  position: sticky;
  top: 84px;
}

.category-header {
  position: relative;
  overflow: hidden;
  display: flex;
  align-items: center;
  gap: 16px;
  background: linear-gradient(135deg, #eef2ff 0%, #f5f0ff 100%);
  border: 1px solid #e3e4ff;
  border-radius: var(--campus-radius);
  padding: 24px 28px;
  margin-bottom: 18px;
}

.category-header::after {
  content: '';
  position: absolute;
  width: 110px;
  height: 110px;
  border-radius: 50%;
  background: rgba(79, 110, 247, 0.08);
  right: -30px;
  top: -40px;
}

.category-header-icon {
  flex-shrink: 0;
  width: 58px;
  height: 58px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 16px;
  background: linear-gradient(135deg, var(--campus-primary), #8b5cf6);
  color: #fff;
  box-shadow: 0 6px 16px rgba(79, 110, 247, 0.3);
}

.category-header-text {
  min-width: 0;
  position: relative;
  z-index: 1;
}

.category-title {
  font-size: 1.5rem;
  font-weight: 800;
  color: var(--campus-text);
  margin: 0 0 4px;
  line-height: 1.3;
}

.category-subtitle {
  font-size: 0.875rem;
  color: var(--campus-text-secondary);
  margin: 0;
}

.pagination-bar {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 16px;
  margin-top: 20px;
  padding: 12px;
  background: var(--campus-surface);
  border: 1px solid var(--campus-border);
  border-radius: var(--campus-radius);
  box-shadow: var(--campus-shadow);
}

.page-btn {
  min-width: 100px;
}

.page-info {
  font-size: 0.875rem;
  color: var(--campus-text-secondary);
  white-space: nowrap;
}

.mr-1 {
  margin-right: 4px;
}

.ml-1 {
  margin-left: 4px;
}

@media (max-width: 991px) {
  .category-layout {
    grid-template-columns: 1fr;
  }

  .sidebar-sticky {
    position: static;
  }
}

@media (max-width: 575px) {
  .category-header {
    padding: 18px;
  }

  .category-header-icon {
    width: 48px;
    height: 48px;
  }

  .category-title {
    font-size: 1.25rem;
  }

  .pagination-bar {
    flex-direction: column;
    gap: 10px;
  }

  .page-btn {
    width: 100%;
  }
}
</style>
