<template>
  <div class="grid-layout">
    <!-- 侧边栏 -->
    <div class="sidebar">
      <Sidebar />
    </div>
    
    <!-- 主内容 -->
    <div class="main-content">
      <div class="announcement" v-if="announcement.content">
        <div class="announcement-title">公告</div>
        <div class="markdown-body" v-html="announcement.content_html"></div>
      </div>
      
      <div class="filters">
        <select v-model="selectedCategory" @change="loadArticles">
          <option value="">全部分区</option>
          <option v-for="cat in categories" :key="cat.id" :value="cat.id">
            {{ cat.name }}
          </option>
        </select>
      </div>
      
      <ArticleList :articles="articles" />
      
      <div class="pagination" v-if="totalPages > 1">
        <button @click="prevPage" :disabled="page === 1">上一页</button>
        <span>第 {{ page }} / {{ totalPages }} 页</span>
        <button @click="nextPage" :disabled="page === totalPages">下一页</button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import api from '../api'
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
    const selectedCategory = ref('')
    
    const loadArticles = async () => {
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
      loadArticles,
      prevPage,
      nextPage
    }
  }
}
</script>

<style scoped>
.filters {
  margin-bottom: 20px;
}

.filters select {
  padding: 8px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
}

.announcement {
  background: #fef3c7;
  border-left: 4px solid #f59e0b;
  padding: 15px;
  margin-bottom: 20px;
  border-radius: 8px;
}

.announcement-title {
  font-weight: bold;
  margin-bottom: 10px;
  color: #92400e;
}

.pagination {
  display: flex;
  justify-content: center;
  gap: 15px;
  margin-top: 30px;
  align-items: center;
}

.pagination button {
  padding: 8px 16px;
  border: 1px solid #d1d5db;
  background: white;
  border-radius: 6px;
  cursor: pointer;
}

.pagination button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>