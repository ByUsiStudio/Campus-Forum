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

  <v-dialog v-model="editCategoryDialog.show" max-width="500">
    <v-card class="dialog-card">
      <v-card-title class="dialog-title">
        <v-avatar :color="editCategoryDialog.category ? 'primary' : 'success'" size="40" class="mr-3">
          <v-icon color="white" size="20">{{ editCategoryDialog.category ? 'mdi-pencil' : 'mdi-plus' }}</v-icon>
        </v-avatar>
        {{ editCategoryDialog.category ? '编辑分区' : '添加分区' }}
      </v-card-title>
      <v-card-text class="dialog-body">
        <v-form ref="categoryForm" v-model="formValid">
          <v-text-field
            v-model="editCategoryDialog.name"
            label="分区名称"
            placeholder="例如：表白墙"
            variant="outlined"
            density="comfortable"
            :rules="[rules.required]"
            prepend-inner-icon="mdi-tag"
            clearable
            class="mb-4"
          >
            <template #label>
              <span class="text-body-2">分区名称</span>
            </template>
          </v-text-field>

          <v-textarea
            v-model="editCategoryDialog.description"
            label="分区描述"
            placeholder="描述分区的内容和用途..."
            variant="outlined"
            density="comfortable"
            prepend-inner-icon="mdi-text"
            rows="3"
            counter
            :maxlength="200"
          >
            <template #label>
              <span class="text-body-2">分区描述</span>
            </template>
          </v-textarea>
        </v-form>
      </v-card-text>
      <v-card-actions class="dialog-actions">
        <v-btn variant="text" @click="closeCategoryDialog" class="mr-2">
          取消
        </v-btn>
        <v-btn
          color="primary"
          variant="flat"
          @click="handleEditCategory"
          :disabled="!formValid"
        >
          <v-icon class="mr-1">mdi-check</v-icon>
          保存
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import CategoriesPanel from './CategoriesPanel.vue'
import api from '../../api'
import { confirm, success, error } from '../../utils/modal'

const categories = ref([])
const loading = ref(true)
const categoryForm = ref(null)
const formValid = ref(false)

const editCategoryDialog = ref({
  show: false,
  category: null,
  name: '',
  description: ''
})

const rules = {
  required: v => !!v || '此字段为必填项'
}

const loadCategories = async () => {
  loading.value = true
  try {
    const response = await api.get('/categories')
    categories.value = response.data.categories || []
  } catch (err) {
    console.error('加载分区列表失败', err)
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
  formValid.value = false
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
    error('请输入分区名称')
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
    success('保存成功')
    closeCategoryDialog()
    loadCategories()
  } catch (err) {
    console.error('保存分区失败', err)
    error(err.response?.data?.error || '保存失败')
  }
}

const handleDeleteCategory = async (category) => {
  const confirmed = await confirm(`确定要删除分区 "${category.name}" 吗？此操作不可恢复。`)
  if (!confirmed) return
  try {
    await api.delete(`/categories/${category.id}`)
    success('删除成功')
    loadCategories()
  } catch (err) {
    console.error('删除分区失败', err)
    error(err.response?.data?.error || '删除失败')
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

.dialog-body {
  padding: 24px !important;
}

.dialog-actions {
  padding: 16px 24px 24px;
  gap: 12px;
}

:deep(.v-field) {
  border-radius: 12px;
}

:deep(.v-field--outlined .v-field__outline) {
  border-color: rgba(148, 163, 184, 0.3);
}

:deep(.v-field--focused .v-field__outline) {
  border-color: rgb(var(--v-theme-primary));
}
</style>