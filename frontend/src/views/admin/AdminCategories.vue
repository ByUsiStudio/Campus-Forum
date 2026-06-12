<script setup>
import { ref, onMounted } from 'vue'
import { adminApi, categoryApi } from '../../api'

const loading = ref(false)
const categories = ref([])
const editDialog = ref(false)
const createDialog = ref(false)
const selectedCategory = ref(null)
const editForm = ref({ name: '', description: '' })
const createForm = ref({ name: '', description: '' })
const saving = ref(false)

const loadCategories = async () => {
  loading.value = true
  try {
    const response = await categoryApi.getCategories()
    categories.value = response.data.categories || []
  } catch (error) {
    console.error('加载分类失败:', error)
  } finally {
    loading.value = false
  }
}

const openCreateDialog = () => {
  createForm.value = { name: '', description: '' }
  createDialog.value = true
}

const createCategory = async () => {
  saving.value = true
  try {
    await categoryApi.createCategory(createForm.value)
    createDialog.value = false
    loadCategories()
  } catch (error) {
    console.error('创建分类失败:', error)
  } finally {
    saving.value = false
  }
}

const editCategory = (category) => {
  selectedCategory.value = category
  editForm.value = { name: category.name, description: category.description || '' }
  editDialog.value = true
}

const saveCategory = async () => {
  saving.value = true
  try {
    await categoryApi.updateCategory(selectedCategory.value.id, editForm.value)
    editDialog.value = false
    loadCategories()
  } catch (error) {
    console.error('更新分类失败:', error)
  } finally {
    saving.value = false
  }
}

const deleteCategory = async (category) => {
  if (!confirm(`确定要删除分类 ${category.name} 吗？`)) return
  
  loading.value = true
  try {
    await categoryApi.deleteCategory(category.id)
    loadCategories()
  } catch (error) {
    console.error('删除分类失败:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadCategories()
})
</script>

<template>
  <v-container fluid class="pa-0">
    <!-- 页面标题 -->
    <div class="mb-6">
      <h1 class="text-h5 font-weight-bold">分区管理</h1>
      <p class="text-body-2 text-grey">管理文章分区，创建、编辑、删除分区</p>
    </div>

    <!-- 操作按钮 -->
    <v-card class="mb-4">
      <v-card-text>
        <v-btn color="primary" @click="openCreateDialog">
          <v-icon start>mdi-plus</v-icon>
          新建分区
        </v-btn>
        <v-btn class="ml-2" @click="loadCategories" :loading="loading">
          <v-icon start>mdi-refresh</v-icon>
          刷新
        </v-btn>
      </v-card-text>
    </v-card>

    <!-- 分类列表 -->
    <v-card>
      <v-card-title>分区列表</v-card-title>
      <v-list>
        <v-list-item v-for="category in categories" :key="category.id">
          <v-list-item-title>{{ category.name }}</v-list-item-title>
          <v-list-item-subtitle>{{ category.description || '无描述' }}</v-list-item-subtitle>
          <v-list-item-actions>
            <v-btn icon variant="text" size="small" @click="editCategory(category)">
              <v-icon>mdi-pencil</v-icon>
            </v-btn>
            <v-btn icon variant="text" size="small" color="error" @click="deleteCategory(category)">
              <v-icon>mdi-delete</v-icon>
            </v-btn>
          </v-list-item-actions>
        </v-list-item>
      </v-list>
      <v-card-text v-if="categories.length === 0" class="text-center text-grey py-8">
        暂无分区
      </v-card-text>
    </v-card>

    <!-- 创建分类对话框 -->
    <v-dialog v-model="createDialog" max-width="400">
      <v-card>
        <v-card-title>新建分区</v-card-title>
        <v-card-text>
          <v-text-field v-model="createForm.name" label="分区名称" class="mb-3" />
          <v-textarea v-model="createForm.description" label="分区描述" rows="2" />
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn text @click="createDialog = false">取消</v-btn>
          <v-btn color="primary" @click="createCategory" :loading="saving">创建</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- 编辑分类对话框 -->
    <v-dialog v-model="editDialog" max-width="400">
      <v-card>
        <v-card-title>编辑分区</v-card-title>
        <v-card-text>
          <v-text-field v-model="editForm.name" label="分区名称" class="mb-3" />
          <v-textarea v-model="editForm.description" label="分区描述" rows="2" />
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn text @click="editDialog = false">取消</v-btn>
          <v-btn color="primary" @click="saveCategory" :loading="saving">保存</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>