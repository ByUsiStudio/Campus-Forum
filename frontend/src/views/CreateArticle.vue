<template>
  <div>
    <v-card class="pa-6" max-width="900" style="margin: 0 auto;">
      <v-card-title class="text-h5 mb-4" style="color: rgb(var(--v-theme-primary));">
        {{ isEdit ? '编辑文章' : '发布新文章' }}
      </v-card-title>

      <v-form @submit.prevent="submitArticle">
        <v-text-field
          v-model="form.title"
          label="标题"
          variant="outlined"
          required
          class="mb-4"
        ></v-text-field>

        <v-select
          v-model="form.category_id"
          :items="categories"
          item-title="name"
          item-value="id"
          label="分区"
          variant="outlined"
          class="mb-4"
        ></v-select>

        <div class="mb-4">
          <MdEditor
            v-model="form.content"
            :height="500"
          />
        </div>

        <div class="d-flex gap-4">
          <v-btn
            type="submit"
            color="primary"
            size="large"
            :loading="submitting"
          >
            {{ submitting ? '提交中...' : (isEdit ? '更新文章' : '发布文章') }}
          </v-btn>
        </div>
      </v-form>
    </v-card>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '../api'
import MdEditor from '../components/MdEditor.vue'
import { alert as showAlert, error as showError } from '../utils/modal'

export default {
  name: 'CreateArticle',
  components: {
    MdEditor
  },
  setup() {
    const route = useRoute()
    const router = useRouter()
    const categories = ref([])
    const form = ref({
      title: '',
      category_id: null,
      content: ''
    })
    const isEdit = ref(false)
    const articleId = ref(null)
    const submitting = ref(false)

    const loadCategories = async () => {
      try {
        const response = await api.get('/categories')
        categories.value = response.data.categories
        if (categories.value.length > 0 && !form.value.category_id) {
          form.value.category_id = categories.value[0].id
        }
      } catch (error) {
        console.error('加载分区失败', error)
      }
    }

    const loadArticle = async () => {
      try {
        const response = await api.get(`/articles/${articleId.value}`)
        const article = response.data.article
        form.value.title = article.title
        form.value.category_id = article.category_id
        form.value.content = article.content
      } catch (error) {
        console.error('加载文章失败', error)
      }
    }

    const submitArticle = async () => {
      if (!form.value.title.trim()) {
        await showAlert('请输入标题')
        return
      }
      if (!form.value.content.trim()) {
        await showAlert('请输入内容')
        return
      }

      submitting.value = true

      try {
        let response
        if (isEdit.value) {
          response = await api.put(`/articles/${articleId.value}`, form.value)
        } else {
          response = await api.post('/articles', form.value)
        }
        router.push(`/article/${response.data.article.id}`)
      } catch (error) {
        console.error('提交失败', error)
        await showError('提交失败')
      } finally {
        submitting.value = false
      }
    }

    onMounted(() => {
      const id = route.query.id
      if (id) {
        isEdit.value = true
        articleId.value = parseInt(id)
        loadArticle()
      }
      loadCategories()
    })

    return {
      categories,
      form,
      isEdit,
      submitting,
      submitArticle
    }
  }
}
</script>