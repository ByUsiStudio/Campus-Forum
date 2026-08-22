<template>
  <div class="page-container home-page">
    <!-- 顶部欢迎横幅 -->
    <header class="hero-banner">
      <div class="hero-content">
        <div class="hero-badge">
          <el-icon class="hero-badge-icon"><Promotion /></el-icon>
          <span class="hero-badge-text">校园社区</span>
        </div>
        <h1 class="hero-title">探索新鲜校园动态</h1>
        <p class="hero-subtitle">发现有趣话题，分享校园生活的每一个精彩瞬间</p>
        <div class="hero-actions">
          <el-button type="primary" round class="hero-btn hero-btn-write" @click="$router.push('/create')">
            <el-icon class="btn-icon"><EditPen /></el-icon>
            写文章
          </el-button>
          <el-button
            round
            class="hero-btn hero-btn-browse"
            plain
            @click="scrollToFeed"
          >
            浏览文章
            <el-icon class="btn-icon"><ArrowRight /></el-icon>
          </el-button>
        </div>
      </div>
      <div class="hero-art" aria-hidden="true">
        <div class="hero-art-circle hero-art-circle-1"></div>
        <div class="hero-art-circle hero-art-circle-2"></div>
        <div class="hero-art-circle hero-art-circle-3"></div>
      </div>
    </header>

    <div class="home-layout">
      <!-- 主内容（左宽列） -->
      <main class="main-column" ref="feedEl">
        <!-- 公告 -->
        <section v-if="announcement.content" class="announcement-card">
          <div class="announcement-head">
            <el-icon class="announcement-icon"><BellFilled /></el-icon>
            <span class="announcement-title">公告</span>
          </div>
          <div class="markdown-body announcement-body" v-html="announcement.content_html"></div>
        </section>

        <!-- 分区和操作栏 -->
        <section class="card-surface filter-bar">
          <div class="filter-row">
            <div class="filter-left">
              <span class="filter-label">
                <el-icon class="filter-label-icon"><FolderOpened /></el-icon>
                分区浏览
              </span>
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
            </div>

            <span class="text-secondary total-pages">
              共 {{ totalPages }} 页
            </span>
          </div>
        </section>

        <!-- 文章列表 -->
        <ArticleList :articles="articles" :loading="loading" />

        <!-- 分页 -->
        <div v-if="totalPages > 1" class="card-surface mt-4 pagination-bar">
          <div class="pagination-inner">
            <el-button
              @click="prevPage"
              :disabled="page === 1"
              type="primary"
              plain
              size="default"
              round
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
              size="default"
              round
            >
              下一页
              <el-icon class="ml-1"><ArrowRight /></el-icon>
            </el-button>
          </div>
        </div>
      </main>

      <!-- 侧边栏（右宽列，桌面端 sticky） -->
      <aside class="sidebar-column">
        <!-- 移动端快捷操作 -->
        <div class="mobile-quick">
          <div class="card-surface mobile-quick-card">
            <el-button type="primary" class="mobile-quick-btn" round @click="$router.push('/create')">
              <el-icon class="mr-1"><EditPen /></el-icon>
              写文章
            </el-button>
          </div>
        </div>

        <div class="sidebar-sticky">
          <Sidebar />
        </div>
      </aside>
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
    const feedEl = ref(null)

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

    const scrollToFeed = () => {
      if (feedEl.value) {
        feedEl.value.scrollIntoView({ behavior: 'smooth', block: 'start' })
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
      feedEl,
      loadArticles,
      prevPage,
      nextPage,
      scrollToFeed
    }
  }
}
</script>

<style scoped>
.home-page {
  padding-top: 16px;
}

/* ===== 顶部欢迎横幅 ===== */
.hero-banner {
  position: relative;
  overflow: hidden;
  border-radius: calc(var(--campus-radius) + 6px);
  background: linear-gradient(135deg, #4f6ef7 0%, #7c5cf5 55%, #a45cf0 100%);
  color: #fff;
  padding: 32px 28px;
  margin-bottom: 20px;
  box-shadow: 0 10px 30px rgba(79, 110, 247, 0.28);
}

.hero-content {
  position: relative;
  z-index: 2;
  max-width: 640px;
}

.hero-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: rgba(255, 255, 255, 0.18);
  border: 1px solid rgba(255, 255, 255, 0.35);
  backdrop-filter: blur(4px);
  padding: 4px 12px;
  border-radius: 999px;
  font-size: 13px;
  font-weight: 500;
  margin-bottom: 14px;
}

.hero-badge-icon {
  font-size: 15px;
}

.hero-title {
  font-size: 1.9rem;
  font-weight: 800;
  line-height: 1.25;
  margin: 0 0 8px;
}

.hero-subtitle {
  font-size: 0.975rem;
  opacity: 0.92;
  margin: 0 0 20px;
  line-height: 1.6;
}

.hero-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.hero-btn {
  font-weight: 600;
}

.hero-btn-write {
  --el-button-bg-color: #fff;
  --el-button-text-color: #4f6ef7;
  --el-button-border-color: #fff;
  --el-button-hover-bg-color: #f1f3ff;
  --el-button-hover-text-color: #3d5bf0;
  --el-button-hover-border-color: #f1f3ff;
}

.hero-btn-browse {
  --el-button-bg-color: transparent;
  --el-button-text-color: #fff;
  --el-button-border-color: rgba(255, 255, 255, 0.7);
  --el-button-hover-bg-color: rgba(255, 255, 255, 0.15);
  --el-button-hover-text-color: #fff;
  --el-button-hover-border-color: #fff;
}

.btn-icon {
  margin-right: 4px;
}

/* 装饰圆 */
.hero-art {
  position: absolute;
  inset: 0;
  pointer-events: none;
}

.hero-art-circle {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.12);
}

.hero-art-circle-1 {
  width: 220px;
  height: 220px;
  right: -40px;
  top: -70px;
}

.hero-art-circle-2 {
  width: 140px;
  height: 140px;
  right: 120px;
  bottom: -60px;
}

.hero-art-circle-3 {
  width: 60px;
  height: 60px;
  right: 40px;
  bottom: 30px;
  background: rgba(255, 255, 255, 0.2);
}

/* ===== 双列布局 ===== */
.home-layout {
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

/* ===== 公告 ===== */
.announcement-card {
  background: linear-gradient(135deg, #eef2ff 0%, #f5f3ff 100%);
  border: 1px solid #dfe4ff;
  border-radius: var(--campus-radius);
  padding: 14px 18px;
  margin-bottom: 16px;
}

.announcement-head {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 6px;
}

.announcement-icon {
  color: var(--campus-primary);
  font-size: 16px;
}

.announcement-title {
  font-weight: 700;
  color: var(--campus-primary);
  font-size: 14px;
}

.announcement-body :deep(.markdown-body) {
  font-size: 0.875rem;
  color: var(--campus-text);
  line-height: 1.6;
}

/* ===== 分区过滤 ===== */
.filter-bar {
  padding: 12px 14px;
}

.filter-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  flex-wrap: wrap;
}

.filter-left {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.filter-label {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  font-weight: 600;
  color: var(--campus-text);
  white-space: nowrap;
}

.filter-label-icon {
  color: var(--campus-primary);
}

.category-select {
  width: 200px;
  max-width: 200px;
}

.total-pages {
  font-size: 13px;
  white-space: nowrap;
}

/* ===== 分页 ===== */
.pagination-bar {
  padding: 14px;
}

.pagination-inner {
  display: flex;
  justify-content: center;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px 12px;
}

.pagination-inner .el-pagination {
  white-space: nowrap;
}

/* ===== 移动端快捷操作 ===== */
.mobile-quick {
  display: none;
}

.mobile-quick-card {
  padding: 12px;
  margin-bottom: 16px;
}

.mobile-quick-btn {
  width: 100%;
}

.mr-1 {
  margin-right: 4px;
}

.ml-1 {
  margin-left: 4px;
}

/* ===== 响应式 ===== */
@media (max-width: 991px) {
  .home-layout {
    grid-template-columns: 1fr;
  }

  .sidebar-sticky {
    position: static;
  }

  .hero-title {
    font-size: 1.6rem;
  }
}

@media (max-width: 767px) {
  .mobile-quick {
    display: block;
  }

  .hero-banner {
    padding: 24px 20px;
  }

  .hero-title {
    font-size: 1.4rem;
  }

  .hero-actions {
    flex-direction: column;
  }

  .hero-btn {
    width: 100%;
  }

  .hero-art-circle-1 {
    width: 160px;
    height: 160px;
    right: -50px;
    top: -60px;
  }
}

@media (max-width: 575px) {
  .category-select {
    width: 100%;
    max-width: 100%;
  }

  .filter-row {
    flex-direction: column;
    align-items: stretch;
  }

  .total-pages {
    align-self: flex-end;
  }

  .pagination-inner {
    flex-direction: column;
    gap: 8px;
  }
}
</style>
