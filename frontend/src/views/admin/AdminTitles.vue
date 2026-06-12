<script setup>
import { ref, onMounted } from 'vue'
import { adminTitleApi } from '../../api/admin'

const loading = ref(false)
const titles = ref([])
const editDialog = ref(false)
const createDialog = ref(false)
const selectedTitle = ref(null)
const editForm = ref({ name: '', description: '', color: '#6200EE' })
const createForm = ref({ name: '', description: '', color: '#6200EE' })
const saving = ref(false)

const loadTitles = async () => {
  loading.value = true
  try {
    const response = await adminTitleApi.getTitles()
    titles.value = response.data.titles || []
  } catch (error) {
    console.error('加载头衔失败:', error)
  } finally {
    loading.value = false
  }
}

const openCreateDialog = () => {
  createForm.value = { name: '', description: '', color: '#6200EE' }
  createDialog.value = true
}

const createTitle = async () => {
  saving.value = true
  try {
    await adminTitleApi.createTitle(createForm.value)
    createDialog.value = false
    loadTitles()
  } catch (error) {
    console.error('创建头衔失败:', error)
  } finally {
    saving.value = false
  }
}

const editTitle = (title) => {
  selectedTitle.value = title
  editForm.value = { name: title.name, description: title.description || '', color: title.color || '#6200EE' }
  editDialog.value = true
}

const saveTitle = async () => {
  saving.value = true
  try {
    await adminTitleApi.updateTitle(selectedTitle.value.id, editForm.value)
    editDialog.value = false
    loadTitles()
  } catch (error) {
    console.error('更新头衔失败:', error)
  } finally {
    saving.value = false
  }
}

const deleteTitle = async (title) => {
  if (!confirm(`确定要删除头衔 ${title.name} 吗？`)) return
  
  loading.value = true
  try {
    await adminTitleApi.deleteTitle(title.id)
    loadTitles()
  } catch (error) {
    console.error('删除头衔失败:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadTitles()
})
</script>

<template>
  <v-container fluid class="pa-0">
    <!-- 页面标题 -->
    <div class="mb-6">
      <h1 class="text-h5 font-weight-bold">头衔管理</h1>
      <p class="text-body-2 text-grey">管理用户头衔，创建、编辑、删除头衔</p>
    </div>

    <!-- 操作按钮 -->
    <v-card class="mb-4">
      <v-card-text>
        <v-btn color="primary" @click="openCreateDialog">
          <v-icon start>mdi-plus</v-icon>
          新建头衔
        </v-btn>
        <v-btn class="ml-2" @click="loadTitles" :loading="loading">
          <v-icon start>mdi-refresh</v-icon>
          刷新
        </v-btn>
      </v-card-text>
    </v-card>

    <!-- 头衔列表 -->
    <v-card>
      <v-card-title>头衔列表</v-card-title>
      <v-list>
        <v-list-item v-for="title in titles" :key="title.id">
          <v-list-item-title>
            <v-chip :color="title.color" size="small">{{ title.name }}</v-chip>
          </v-list-item-title>
          <v-list-item-subtitle>{{ title.description || '无描述' }}</v-list-item-subtitle>
          <v-list-item-actions>
            <v-btn icon variant="text" size="small" @click="editTitle(title)">
              <v-icon>mdi-pencil</v-icon>
            </v-btn>
            <v-btn icon variant="text" size="small" color="error" @click="deleteTitle(title)">
              <v-icon>mdi-delete</v-icon>
            </v-btn>
          </v-list-item-actions>
        </v-list-item>
      </v-list>
      <v-card-text v-if="titles.length === 0" class="text-center text-grey py-8">
        暂无头衔
      </v-card-text>
    </v-card>

    <!-- 创建头衔对话框 -->
    <v-dialog v-model="createDialog" max-width="400">
      <v-card>
        <v-card-title>新建头衔</v-card-title>
        <v-card-text>
          <v-text-field v-model="createForm.name" label="头衔名称" class="mb-3" />
          <v-textarea v-model="createForm.description" label="头衔描述" rows="2" class="mb-3" />
          <v-text-field v-model="createForm.color" label="颜色" type="color" />
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn text @click="createDialog = false">取消</v-btn>
          <v-btn color="primary" @click="createTitle" :loading="saving">创建</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- 编辑头衔对话框 -->
    <v-dialog v-model="editDialog" max-width="400">
      <v-card>
        <v-card-title>编辑头衔</v-card-title>
        <v-card-text>
          <v-text-field v-model="editForm.name" label="头衔名称" class="mb-3" />
          <v-textarea v-model="editForm.description" label="头衔描述" rows="2" class="mb-3" />
          <v-text-field v-model="editForm.color" label="颜色" type="color" />
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn text @click="editDialog = false">取消</v-btn>
          <v-btn color="primary" @click="saveTitle" :loading="saving">保存</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>