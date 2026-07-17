<template>
  <div class="grid-layout">
    <div class="sidebar">
      <Sidebar />
    </div>
    
    <div class="main-content">
      <div class="layui-card mb-4" style="padding: 24px;">
        <h2 style="font-size: 24px; font-weight: 600; margin: 0;">{{ categoryName }}</h2>
      </div>
      
      <ArticleList :articles="articles" />
      
      <div class="flex justify-center items-center gap-4 mt-4" v-if="totalPages > 1">
        <button 
          @click="prevPage" 
          :disabled="page === 1" 
          class="layui-btn layui-btn-primary"
        >
          上一页
        </button>
        <span style="color: #666;">第 {{ page }} / {{ totalPages }} 页</span>
        <button 
          @click="nextPage" 
          :disabled="page === totalPages" 
          class="layui-btn layui-btn-primary"
        >
          下一页
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import api from '../api'
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
        const response = await api.get('/articles', {
          params: {
            category_id: route.params.id,
            page: page.value,
            page_size: 20
          }
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
</style>