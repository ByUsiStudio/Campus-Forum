<script setup>
import { ref, inject, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { articleApi, categoryApi } from '../api'

const router = useRouter()
const user = inject('user')

const form = ref({
  title: '',
  content: '',
  category_id: ''
})

const categories = ref([])
const isLoading = ref(false)
const error = ref('')

const loadCategories = async () => {
  try {
    const response = await categoryApi.getCategories()
    categories.value = response.data.categories
  } catch (error) {
    console.error('加载分类失败:', error)
  }
}

const handleSubmit = async () => {
  if (!form.value.title.trim()) {
    error.value = '请输入文章标题'
    return
  }
  if (!form.value.content.trim()) {
    error.value = '请输入文章内容'
    return
  }
  if (!form.value.category_id) {
    error.value = '请选择分类'
    return
  }
  isLoading.value = true
  try {
    await articleApi.create(form.value)
    router.push('/')
  } catch (err) {
    error.value = err.response?.data?.error || '发布失败'
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  if (!user.value) {
    router.push('/login')
    return
  }
  loadCategories()
})
</script>

<template>
  <v-app>
    <v-app-bar app>
      <v-btn icon @click="router.push('/')">
        <v-icon>mdi-arrow-left</v-icon>
      </v-btn>
      <v-toolbar-title>发布文章</v-toolbar-title>
    </v-app-bar>
    
    <v-container class="py-6">
      <v-card>
        <v-card-text>
          <v-alert v-if="error" type="error" dense>
            {{ error }}
          </v-alert>
          
          <v-text-field
            v-model="form.title"
            label="文章标题"
            class="mb-4"
          />
          
          <v-select
            v-model="form.category_id"
            :items="categories"
            item-text="name"
            item-value="id"
            label="选择分类"
            class="mb-4"
          />
          
          <v-textarea
            v-model="form.content"
            label="文章内容"
            rows="15"
            class="mb-6"
          />
          
          <v-btn
            color="primary"
            :loading="isLoading"
            @click="handleSubmit"
          >
            发布文章
          </v-btn>
          
          <v-btn text @click="router.push('/')" class="ml-4">
            取消
          </v-btn>
        </v-card-text>
      </v-card>
    </v-container>
  </v-app>
</template>
