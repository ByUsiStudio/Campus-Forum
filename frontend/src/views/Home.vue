<template>
  <div class="grid-layout">
    <!-- 侧边栏 -->
    <div class="sidebar">
      <Sidebar />
    </div>
    
    <!-- 主内容 -->
    <div class="main-content">
      <v-alert v-if="announcement.content" type="info" variant="tonal" class="mb-4">
        <div class="announcement-title font-weight-bold mb-2">公告</div>
        <div class="markdown-body" v-html="announcement.content_html"></div>
      </v-alert>
      
      <div class="d-flex align-center gap-4 mb-4">
        <v-select
          v-model="selectedCategory"
          :items="categoryOptions"
          label="选择分区"
          variant="outlined"
          density="compact"
          hide-details
          @update:model-value="loadArticles"
          style="max-width: 200px;"
        ></v-select>
      </div>
      
      <ArticleList :articles="articles" />
      
      <div class="d-flex justify-center align-center gap-4 mt-4" v-if="totalPages > 1">
        <v-btn 
          @click="prevPage" 
          :disabled="page === 1" 
          variant="outlined"
          color="primary"
        >
          上一页
        </v-btn>
        <span class="text-body-2">第 {{ page }} / {{ totalPages }} 页</span>
        <v-btn 
          @click="nextPage" 
          :disabled="page === totalPages" 
          variant="outlined"
          color="primary"
        >
          下一页
        </v-btn>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
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
    const selectedCategory = ref(null)
    
    const categoryOptions = computed(() => {
      return [
        { title: '全部分区', value: null },
        ...categories.value.map(cat => ({ title: cat.name, value: cat.id }))
      ]
    })
    
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
      categoryOptions,
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
</style>
