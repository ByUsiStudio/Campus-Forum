<template>
  <div class="grid-layout">
    <div class="sidebar">
      <Sidebar />
    </div>
    
    <div class="main-content">
      <v-card class="pa-6 mb-4">
        <v-card-title class="text-h5">{{ categoryName }}</v-card-title>
      </v-card>
      
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
