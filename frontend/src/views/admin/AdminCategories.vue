<template>
  <div class="admin-categories">
    <CategoriesPanel
      :categories="categories"
      :loading="loading"
      @add="addCategory"
      @edit="showEditCategoryDialog"
      @delete="handleDeleteCategory"
      @refresh="loadCategories"
    />

    <div v-if="editCategoryDialog.show" class="dialog-overlay" @click.self="closeCategoryDialog">
      <div class="dialog-content">
        <div class="dialog-header">
          <div class="dialog-title">
            <div class="title-avatar" :class="editCategoryDialog.category ? 'primary' : 'success'">
              <i class="fa-solid" :class="editCategoryDialog.category ? 'fa-pencil' : 'fa-plus'"></i>
            </div>
            <span>{{ editCategoryDialog.category ? '编辑分区' : '添加分区' }}</span>
          </div>
          <button class="dialog-close" @click="closeCategoryDialog">
            <i class="fa-solid fa-xmark"></i>
          </button>
        </div>
        <div class="dialog-body">
          <div class="form-group">
            <label class="form-label">分区名称</label>
            <input
              v-model="editCategoryDialog.name"
              type="text"
              placeholder="例如：表白墙"
              class="layui-input"
            />
          </div>
          <div class="form-group">
            <label class="form-label">分区描述</label>
            <textarea
              v-model="editCategoryDialog.description"
              placeholder="描述分区的内容和用途..."
              class="layui-textarea"
              rows="3"
              maxlength="200"
            ></textarea>
            <span class="char-count">{{ editCategoryDialog.description.length }}/200</span>
          </div>
        </div>
        <div class="dialog-footer">
          <button class="layui-btn layui-btn-primary" @click="closeCategoryDialog">取消</button>
          <button
            class="layui-btn"
            @click="handleEditCategory"
            :disabled="!editCategoryDialog.name.trim()"
          >
            <i class="fa-solid fa-check"></i>保存
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import CategoriesPanel from './CategoriesPanel.vue'
import { adminCategoryApi } from '../../api/admin'
import { confirm, success, error } from '../../utils/modal'

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
    const response = await adminCategoryApi.getCategories()
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
  const confirmed = await confirm(`确定要删除分区 "${category.name}" 吗？此操作不可恢复。`)
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
.admin-categories {
  padding: 20px;
}

.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.dialog-content {
  background: #fff;
  border-radius: 12px;
  width: 90%;
  max-width: 500px;
  overflow: hidden;
}

.dialog-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24px;
  background: linear-gradient(135deg, #f8f9ff 0%, #fff 100%);
  border-bottom: 1px solid #f0f0f0;
}

.dialog-title {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 18px;
  font-weight: 600;
}

.title-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 20px;

  &.primary {
    background: #1E9FFF;
  }

  &.success {
    background: #52C41A;
  }
}

.dialog-close {
  padding: 4px;
  background: transparent;
  border: none;
  color: #999;
  font-size: 20px;
  cursor: pointer;

  &:hover {
    color: #333;
  }
}

.dialog-body {
  padding: 24px;
}

.form-group {
  margin-bottom: 20px;
}

.form-label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin-bottom: 8px;
}

.char-count {
  display: block;
  text-align: right;
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px;
  border-top: 1px solid #f0f0f0;
}
</style>