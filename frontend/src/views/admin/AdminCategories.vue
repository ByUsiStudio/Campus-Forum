<template>
  <v-container fluid class="pa-6">
    <CategoriesPanel
      :categories="categories"
      :loading="loading"
      @add="addCategory"
      @edit="showEditCategoryDialog"
      @delete="handleDeleteCategory"
      @refresh="loadCategories"
    />
  </v-container>

  <v-dialog v-model="editCategoryDialog.show" max-width="480">
    <v-card class="dialog-card">
      <v-card-title class="dialog-title">
        <v-icon class="title-icon">{{ editCategoryDialog.category ? 'mdi-pencil' : 'mdi-plus' }}</v-icon>
        {{ editCategoryDialog.category ? '编辑分区' : '添加分区' }}
      </v-card-title>
      <v-card-text class="dialog-body">
        <v-text-field
          v-model="editCategoryDialog.name"
          label="分区名称"
          variant="outlined"
          class="mt-4"
        ></v-text-field>
        <v-text-field
          v-model="editCategoryDialog.description"
          label="描述"
          variant="outlined"
          class="mt-4"
        ></v-text-field>
      </v-card-text>
      <v-card-actions class="dialog-actions">
        <v-btn variant="text" @click="closeCategoryDialog">取消</v-btn>
        <v-btn color="primary" variant="flat" @click="handleEditCategory">保存</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import CategoriesPanel from './CategoriesPanel.vue'
import api from '../api'
import { showConfirm, showSuccess, showError } from '../utils/modal'

const categories = ref([])
const loading = ref(true)

const editCategoryDialog = ref({
  show: false,
  category: null,
  name: '',
  description: ''
})

const loadCategories = async () => {
  loading.value = true
  try {
    const response = await api.get('/categories')
    categories.value = response.data
  } catch (error) {
    console.error('加载分区列表失败', error)
  } finally {
    loading.value = false
  }
}

const addCategory = () => {
  editCategoryDialog.value = {
    show: true,
    category: null,
    name: '',
    description: ''
  }
}

const showEditCategoryDialog = (category) => {
  editCategoryDialog.value = {
    show: true,
    category,
    name: category.name,
    description: category.description || ''
  }
}

const closeCategoryDialog = () => {
  editCategoryDialog.value.show = false
}

const handleEditCategory = async () => {
  if (!editCategoryDialog.value.name) {
    showError('请输入分区名称')
    return
  }
  
  const data = {
    name: editCategoryDialog.value.name,
    description: editCategoryDialog.value.description
  }

  try {
    if (editCategoryDialog.value.category) {
      await api.put(`/categories/${editCategoryDialog.value.category.id}`, data)
    } else {
      await api.post('/categories', data)
    }
    showSuccess('保存成功')
    closeCategoryDialog()
    loadCategories()
  } catch (error) {
    console.error('保存分区失败', error)
    showError(error.response?.data?.error || '保存失败')
  }
}

const handleDeleteCategory = async (category) => {
  const confirmed = await showConfirm(`确定要删除分区 "${category.name}" 吗？`)
  if (!confirmed) return
  try {
    await api.delete(`/categories/${category.id}`)
    showSuccess('删除成功')
    loadCategories()
  } catch (error) {
    console.error('删除分区失败', error)
    showError(error.response?.data?.error || '删除失败')
  }
}

onMounted(() => {
  loadCategories()
})
</script>

<style scoped>
.dialog-card {
  border-radius: 20px !important;
  overflow: hidden;
}

.dialog-title {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 24px 24px 16px;
  font-size: 1.2rem;
  font-weight: 700;
  background: linear-gradient(135deg, #f8f9ff 0%, #fff 100%);
}

.title-icon {
  width: 40px;
  height: 40px;
  padding: 8px;
  border-radius: 10px;
  background: rgba(103, 80, 164, 0.1);
}

.dialog-body {
  padding: 24px !important;
}

.dialog-actions {
  padding: 16px 24px 24px;
  gap: 12px;
}
</style>