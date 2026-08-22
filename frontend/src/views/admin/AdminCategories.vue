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
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import CategoriesPanel from './CategoriesPanel.vue'
import { adminCategoryApi } from '../../api/admin'
import { confirm, success, error } from '@/utils/message'
import { jcCloseAll, jcFieldsConfig, jcOpenHtml } from '@/utils/jcu'

const categories = ref([])
const loading = ref(true)

const emptyDialog = () => ({
  category: null,
  name: '',
  description: ''
})

const editCategoryDialog = ref(emptyDialog())

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
  openCategoryDialog()
}

const showEditCategoryDialog = (category) => {
  editCategoryDialog.value = {
    category,
    name: category.name,
    description: category.description || ''
  }
  openCategoryDialog()
}

const openCategoryDialog = () => {
  const isEdit = !!editCategoryDialog.value.category
  const cfg = jcFieldsConfig([
    { name: 'name', label: '分区名称', type: 'text', placeholder: '例如：表白墙', value: editCategoryDialog.value.name, required: true },
    { name: 'description', label: '分区描述', type: 'textarea', placeholder: '描述分区的内容和用途...', value: editCategoryDialog.value.description, rows: 3 }
  ])
  jcOpenHtml({
    title: isEdit ? '编辑分区' : '添加分区',
    content: cfg.html,
    width: 520,
    buttons: [
      { text: '取消', type: 'default', action: () => jcCloseAll() },
      {
        text: '保存',
        type: 'primary',
        action: (inst) => {
          if (!cfg.validate(inst.modalContent)) return
          handleEditCategory(cfg.collect(inst.modalContent))
        }
      }
    ]
  })
}

const handleEditCategory = async (data) => {
  if (!data.name) {
    error('请输入分区名称')
    return
  }

  const payload = {
    name: data.name,
    description: data.description
  }

  try {
    if (editCategoryDialog.value.category) {
      await adminCategoryApi.updateCategory(editCategoryDialog.value.category.id, payload)
    } else {
      await adminCategoryApi.createCategory(payload)
    }
    success('保存成功')
    jcCloseAll()
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
