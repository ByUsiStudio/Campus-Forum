<template>
  <div class="page-container home-page">
    <el-row :gutter="16">
      <!-- 侧边栏 - 桌面端显示 -->
      <el-col :xs="0" :sm="0" :md="6" :lg="6" class="sidebar-col">
        <Sidebar />
      </el-col>

      <!-- 主内容 -->
      <el-col :xs="24" :sm="24" :md="18" :lg="18">
        <!-- 公告 -->
        <el-alert
          v-if="announcement.content"
          type="info"
          :closable="false"
          show-icon
          class="mb-4 announcement-alert"
        >
          <template #title>
            <div class="announcement-title d-flex align-center">
              <el-icon class="mr-1"><Bell /></el-icon>
              <span>公告</span>
            </div>
          </template>
          <div class="markdown-body" v-html="announcement.content_html"></div>
        </el-alert>

        <!-- 分区和操作栏 -->
        <div class="card-surface mb-4 filter-bar">
          <div class="d-flex align-center flex-wrap gap-3">
            <el-select
              v-model="selectedCategory"
              placeholder="选择分区"
              clearable
              class="category-select"
              @change="loadArticles"
            >
              <el-option
                v-for="option in categoryOptions"
                :key="String(option.value)"
                :label="option.title"
                :value="option.value"
              />
            </el-select>

            <div class="filter-spacer"></div>

            <span class="text-secondary total-pages">
              共 {{ totalPages }} 页
            </span>
          </div>
        </div>

        <!-- 文章列表 -->
        <ArticleList :articles="articles" :loading="loading" />

        <!-- 分页 -->
        <div v-if="totalPages > 1" class="card-surface mt-4 pagination-bar">
          <div class="d-flex justify-center align-center flex-wrap gap-2">
            <el-button
              @click="prevPage"
              :disabled="page === 1"
              type="primary"
              plain
              size="small"
            >
              <el-icon class="mr-1"><ArrowLeft /></el-icon>
              上一页
            </el-button>

            <el-pagination
              v-model:current-page="page"
              :page-count="totalPages"
              :pager-count="5"
              layout="pager"
              :hide-on-single-page="false"
              @current-change="loadArticles"
              background
            />

            <el-button
              @click="nextPage"
              :disabled="page === totalPages"
              type="primary"
              plain
              size="small"
            >
              下一页
              <el-icon class="ml-1"><ArrowRight /></el-icon>
            </el-button>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- 移动端快捷操作 -->
    <div class="mobile-quick mt-4">
      <div class="card-surface pa-3">
        <div class="d-flex justify-center gap-2">
          <el-button type="primary" size="small" @click="$router.push('/create')">
            <el-icon class="mr-1"><EditPen /></el-icon>
            写文章
          </el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { articleApi } from '../api'
import { commonApi } from '../api'
import Sidebar from '../components/Sidebar.vue'
import ArticleList from '../components/ArticleList.vue'

export default {
  name: 'Home',
  components: {
    Sidebar,
    ArticleList
  },
  setup() {
    const articles = ref([])
    const categories = ref([])
    const announcement = ref({ content: '', content_html: '' })
    const page = ref(1)
    const totalPages = ref(1)
    const selectedCategory = ref(null)
    const loading = ref(false)

    const categoryOptions = computed(() => {
      return [
        { title: '全部分区', value: null },
        ...categories.value.map(cat => ({ title: cat.name, value: cat.id }))
      ]
    })

    const loadArticles = async () => {
      loading.value = true
      try {
        const params = {
          page: page.value,
          page_size: 20
        }
        if (selectedCategory.value) {
          params.category_id = selectedCategory.value
        }
        const response = await articleApi.getArticles(params)
        articles.value = response.data.articles
        totalPages.value = response.data.total_pages
      } catch (error) {
        console.error('加载文章失败', error)
      } finally {
        loading.value = false
      }
    }

    const loadCategories = async () => {
      try {
        const response = await commonApi.getCategories()
        categories.value = response.data.categories
      } catch (error) {
        console.error('加载分区失败', error)
      }
    }

    const loadAnnouncement = async () => {
      try {
        const response = await commonApi.getAnnouncement()
        announcement.value = response.data
      } catch (error) {
        console.error('加载公告失败', error)
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

    onMounted(() => {
      loadArticles()
      loadCategories()
      loadAnnouncement()
    })

    return {
      articles,
      categories,
      announcement,
      page,
      totalPages,
      selectedCategory,
      categoryOptions,
      loading,
      loadArticles,
      prevPage,
      nextPage
    }
  }
}
</script>

<style scoped>
.home-page {
  padding-top: 16px;
}

.sidebar-col {
  display: block;
}

.mobile-quick {
  display: none;
}

.filter-bar {
  padding: 12px;
}

.category-select {
  width: 200px;
  max-width: 200px;
}

.filter-spacer {
  flex: 1;
}

.total-pages {
  font-size: 13px;
  white-space: nowrap;
}

.pagination-bar {
  padding: 12px;
}

.pagination-bar .el-pagination {
  white-space: nowrap;
}

.announcement-title {
  font-weight: 600;
}

.announcement-title .mr-1 {
  margin-right: 4px;
}

.mr-1 {
  margin-right: 4px;
}

.ml-1 {
  margin-left: 4px;
}

/* Element Plus el-alert body uses padding for description, keep markdown margin clean */
.announcement-alert :deep(.markdown-body) {
  margin-top: 4px;
}

@media (max-width: 991px) {
  .sidebar-col {
    display: none;
  }
}

@media (max-width: 767px) {
  .mobile-quick {
    display: block;
  }
}
</style>
