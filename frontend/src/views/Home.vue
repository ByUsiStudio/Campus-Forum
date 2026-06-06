<template>
  <v-row>
    <!-- 侧边栏 - 桌面端显示 -->
    <v-col cols="12" md="3" class="d-none d-md-block">
      <Sidebar />
    </v-col>
    
    <!-- 主内容 -->
    <v-col cols="12" md="9">
      <!-- 公告 -->
      <v-alert
        v-if="announcement.content"
        type="info"
        variant="tonal"
        class="mb-4"
        prominent
        icon="mdi-bullhorn"
      >
        <div class="announcement-title font-weight-bold mb-2 d-flex align-center">
          <v-icon class="mr-2" size="small">mdi-bullhorn</v-icon>
          公告
        </div>
        <div class="markdown-body" v-html="announcement.content_html"></div>
      </v-alert>
      
      <!-- 分区和操作栏 -->
      <v-card class="mb-4 pa-3" variant="flat">
        <div class="d-flex align-center flex-wrap gap-3">
          <v-select
            v-model="selectedCategory"
            :items="categoryOptions"
            label="选择分区"
            variant="outlined"
            density="compact"
            hide-details
            clearable
            @update:model-value="loadArticles"
            style="max-width: 200px;"
            prepend-inner-icon="mdi-folder"
          ></v-select>
          
          <v-spacer></v-spacer>
          
          <span class="text-caption text-medium-emphasis">
            共 {{ totalPages }} 页
          </span>
        </div>
      </v-card>
      
      <!-- 文章列表 -->
      <ArticleList :articles="articles" :loading="loading" />
      
      <!-- 分页 -->
      <v-card v-if="totalPages > 1" class="mt-4 pa-3" variant="flat">
        <div class="d-flex justify-center align-center flex-wrap gap-2">
          <v-btn 
            @click="prevPage" 
            :disabled="page === 1" 
            variant="outlined"
            color="primary"
            size="small"
            prepend-icon="mdi-chevron-left"
          >
            上一页
          </v-btn>
          
          <v-pagination
            v-model="page"
            :length="totalPages"
            :total-visible="5"
            @update:model-value="loadArticles"
            rounded="circle"
            density="compact"
          ></v-pagination>
          
          <v-btn 
            @click="nextPage" 
            :disabled="page === totalPages" 
            variant="outlined"
            color="primary"
            size="small"
            append-icon="mdi-chevron-right"
          >
            下一页
          </v-btn>
        </div>
      </v-card>
    </v-col>
    
    <!-- 移动端快捷操作 -->
    <v-col cols="12" class="d-md-none mt-4">
      <v-card variant="flat" class="pa-3">
        <div class="d-flex justify-center gap-2">
          <v-btn
            color="primary"
            to="/create"
            prepend-icon="mdi-pencil"
            size="small"
          >
            写文章
          </v-btn>
        </div>
      </v-card>
    </v-col>
  </v-row>
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
