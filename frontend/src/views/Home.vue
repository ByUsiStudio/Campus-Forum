<template>
  <div class="home-page">
    <div class="layui-container">
      <div class="layui-row">
        <div class="layui-col-md-3 d-none d-md-block">
          <div class="sidebar-wrapper animate-fade-in-left">
            <Sidebar />
          </div>
        </div>

        <div class="layui-col-md-9">
          <div v-if="announcement.content" class="announcement-card animate-fade-in-up">
            <div class="layui-card">
              <div class="layui-card-header">
                <i class="layui-icon layui-icon-notice"></i>
                <span>公告</span>
              </div>
              <div class="layui-card-body">
                <div class="markdown-body" v-html="announcement.content_html"></div>
              </div>
            </div>
          </div>

          <div class="filter-card animate-fade-in-up delay-200">
            <div class="layui-card">
              <div class="layui-card-body">
                <div class="filter-row">
                  <div class="filter-item">
                    <select 
                      v-model="selectedCategory" 
                      @change="loadArticles"
                      class="layui-form-select"
                    >
                      <option value="">全部分区</option>
                      <option v-for="cat in categories" :key="cat.id" :value="cat.id">
                        {{ cat.name }}
                      </option>
                    </select>
                  </div>
                  <div class="filter-info">
                    共 {{ totalPages }} 页
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="article-list-wrapper animate-fade-in-up delay-300">
            <ArticleList :articles="articles" :loading="loading" />
          </div>

          <div v-if="totalPages > 1" class="pagination-card animate-fade-in-up delay-400">
            <div class="layui-card">
              <div class="layui-card-body">
                <div class="pagination-container">
                  <button 
                    class="layui-btn layui-btn-primary" 
                    :disabled="page === 1"
                    @click="prevPage"
                  >
                    <i class="layui-icon layui-icon-left"></i>
                    上一页
                  </button>
                  <div class="pagination-numbers">
                    <button 
                      v-for="p in visiblePages" 
                      :key="p"
                      class="page-btn"
                      :class="{ active: p === page }"
                      @click="goToPage(p)"
                    >
                      {{ p }}
                    </button>
                  </div>
                  <button 
                    class="layui-btn layui-btn-primary" 
                    :disabled="page === totalPages"
                    @click="nextPage"
                  >
                    下一页
                    <i class="layui-icon layui-icon-right"></i>
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="layui-col-md-12 d-md-none mt-4">
          <div class="mobile-action-card animate-fade-in-up">
            <div class="layui-card">
              <div class="layui-card-body">
                <div class="mobile-action-row">
                  <button class="layui-btn layui-btn-fluid" @click="goToCreate">
                    <i class="layui-icon layui-icon-edit"></i>
                    写文章
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '../api'
import Sidebar from '../components/Sidebar.vue'
import ArticleList from '../components/ArticleList.vue'

const router = useRouter()

const articles = ref([])
const categories = ref([])
const announcement = ref({ content: '', content_html: '' })
const page = ref(1)
const totalPages = ref(1)
const selectedCategory = ref('')
const loading = ref(false)

const visiblePages = computed(() => {
  const pages = []
  const total = totalPages.value
  const current = page.value
  
  if (total <= 5) {
    for (let i = 1; i <= total; i++) {
      pages.push(i)
    }
  } else {
    if (current <= 3) {
      pages.push(1, 2, 3, 4, 5)
    } else if (current >= total - 2) {
      pages.push(total - 4, total - 3, total - 2, total - 1, total)
    } else {
      pages.push(current - 2, current - 1, current, current + 1, current + 2)
    }
  }
  return pages
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
    const response = await api.get('/articles', { params })
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
    const response = await api.get('/categories')
    categories.value = response.data.categories
  } catch (error) {
    console.error('加载分区失败', error)
  }
}

const loadAnnouncement = async () => {
  try {
    const response = await api.get('/announcement')
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

const goToPage = (p) => {
  if (p !== page.value) {
    page.value = p
    loadArticles()
  }
}

const goToCreate = () => {
  router.push('/create')
}

onMounted(() => {
  loadArticles()
  loadCategories()
  loadAnnouncement()
})
</script>

<style lang="less" scoped>
.home-page {
  padding: 20px 0;
}

.sidebar-wrapper {
  margin-right: 15px;
}

.announcement-card,
.filter-card,
.article-list-wrapper,
.pagination-card,
.mobile-action-card {
  margin-bottom: 20px;
}

.filter-row {
  display: flex;
  align-items: center;
  justify-content: space-between;

  .filter-item {
    flex: 1;
    max-width: 200px;
  }

  .filter-info {
    font-size: 14px;
    color: #999;
  }
}

.pagination-container {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;

  .pagination-numbers {
    display: flex;
    gap: 4px;
  }

  .page-btn {
    width: 36px;
    height: 36px;
    border: 1px solid #e6e6e6;
    background: #fff;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
    transition: all 0.3s @ease-out-back;

    &:hover {
      border-color: #1E9FFF;
      color: #1E9FFF;
      transform: translateY(-2px);
    }

    &.active {
      background: #1E9FFF;
      color: #fff;
      border-color: #1E9FFF;
    }

    &:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }
  }
}

.mobile-action-row {
  display: flex;
  gap: 10px;

  button {
    flex: 1;
  }
}

.markdown-body {
  font-size: 0.95rem;
  line-height: 1.6;
}

.markdown-body :deep(h1),
.markdown-body :deep(h2),
.markdown-body :deep(h3) {
  margin-top: 0.5rem;
  margin-bottom: 0.5rem;
}

.markdown-body :deep(p) {
  margin: 0.5rem 0;
}
</style>
