<template>
  <div>
    <el-card class="page-container" shadow="never">
      <CategoriesPanel
        :categories="categories"
        :loading="loading"
        @add="addCategory"
        @edit="showEditCategoryDialog"
        @delete="handleDeleteCategory"
        @refresh="loadCategories"
      />
    </el-card>

    <el-dialog
      v-model="editCategoryDialog.show"
      width="500px"
      class="dialog-card"
      :title="null"
    >
      <template #header>
        <div class="dialog-header">
          <el-icon :size="20" class="dialog-header-icon">
            <component :is="editCategoryDialog.category ? Edit : Plus" />
          </el-icon>
          <span>{{ editCategoryDialog.category ? '编辑分区' : '添加分区' }}</span>
        </div>
      </template>

      <el-form ref="categoryForm" :model="editCategoryDialog" :rules="rules" label-position="top">
        <el-form-item label="分区名称" prop="name" class="mb-4">
          <el-input
            v-model="editCategoryDialog.name"
            placeholder="例如：表白墙"
            clearable
            :prefix-icon="CollectionTag"
          />
        </el-form-item>

        <el-form-item label="分区描述" prop="description">
          <el-input
            v-model="editCategoryDialog.description"
            type="textarea"
            placeholder="描述分区的内容和用途..."
            :rows="3"
            maxlength="200"
            show-word-limit
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="closeCategoryDialog" class="mr-2">取消</el-button>
        <el-button type="primary" @click="handleEditCategory">
          <el-icon class="mr-1"><Check /></el-icon>
          保存
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Edit, Plus, CollectionTag, Check } from '@element-plus/icons-vue'
import CategoriesPanel from './CategoriesPanel.vue'
import { adminCategoryApi } from '../../api/admin'
import { confirm, success, error } from '@/utils/message'

const categories = ref([])
const loading = ref(true)
const categoryForm = ref(null)

const emptyDialog = () => ({
  show: false,
  category: null,
  name: '',
  description: ''
})

const editCategoryDialog = ref(emptyDialog())

const rules = {
  name: [
    { required: true, message: '此字段为必填项', trigger: 'blur' }
  ]
}

const loadCategories = async () => {
  loading.value = true
  try {
    const response = await adminCategoryApi.getCategories()
    categories.value = response.data.categories || []
  } catch (err) {
    console.error('加载分区列表失败', err)
    error('加载分区列表失败')
  } finally {
    loading.value = false
  }
}

const addCategory = () => {
  editCategoryDialog.value = emptyDialog()
  editCategoryDialog.value.show = true
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
      await adminCategoryApi.updateCategory(editCategoryDialog.value.category.id, data)
    } else {
      await adminCategoryApi.createCategory(data)
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
  const confirmed = await confirm(`确定要删除分区 "${category.name}" 吗？此操作不可恢复。`).catch(() => null)
  if (!confirmed) return
  try {
    await adminCategoryApi.deleteCategory(category.id)
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
  border-radius: 12px;
  overflow: hidden;
}

.dialog-header {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 1.2rem;
  font-weight: 700;
  background: linear-gradient(135deg, #f8f9ff 0%, #fff 100%);
  padding: 8px 4px;
}

.dialog-header-icon {
  color: var(--el-color-primary);
}
</style>
