<template>
  <v-container fluid class="pa-6">
    <TitlesPanel
      :titles="titles"
      :users="users"
      :loading="loading"
      @add-title="addTitle"
      @grant="grantTitle"
      @revoke="revokeTitle"
      @delete-title="handleDeleteTitle"
      @refresh="loadTitles"
    />
  </v-container>

  <v-dialog v-model="grantDialog.show" max-width="480">
    <v-card class="dialog-card">
      <v-card-title class="dialog-title">
        <v-icon class="title-icon">mdi-medal</v-icon>
        授予头衔
      </v-card-title>
      <v-card-text class="dialog-body">
        <v-select
          v-model="grantDialog.selectedUserId"
          :items="usersForSelect"
          item-title="display_name"
          item-value="id"
          label="选择用户"
          variant="outlined"
          class="mt-4"
        ></v-select>
      </v-card-text>
      <v-card-actions class="dialog-actions">
        <v-btn variant="text" @click="grantDialog.show = false">取消</v-btn>
        <v-btn color="primary" variant="flat" @click="handleGrant">确认授予</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import TitlesPanel from './TitlesPanel.vue'
import api from '../../api'
import { confirm, success, error } from '../../utils/modal'

const titles = ref([])
const users = ref([])
const loading = ref(true)

const grantDialog = ref({
  show: false,
  titleId: null,
  selectedUserId: null
})

const usersForSelect = computed(() => {
  return users.value.map(u => ({
    id: u.id,
    display_name: u.display_name
  }))
})

const loadTitles = async () => {
  loading.value = true
  try {
    const response = await api.get('/titles')
    titles.value = response.data.titles || []
  } catch (error) {
    console.error('加载头衔列表失败', error)
  } finally {
    loading.value = false
  }
}

const loadUsers = async () => {
  try {
    const response = await api.get('/admin/users')
    users.value = response.data.users || []
  } catch (error) {
    console.error('加载用户列表失败', error)
  }
}

const addTitle = async () => {
  const titleName = prompt('请输入头衔名称：')
  if (!titleName) return
  
  try {
    await api.post('/titles', { name: titleName })
    success('添加成功')
    loadTitles()
  } catch (error) {
    console.error('添加头衔失败', error)
    error(error.response?.data?.error || '添加失败')
  }
}

const grantTitle = (titleId) => {
  grantDialog.value = {
    show: true,
    titleId,
    selectedUserId: null
  }
}

const handleGrant = async () => {
  if (!grantDialog.value.selectedUserId) {
    error('请选择用户')
    return
  }
  
  try {
    await api.post('/titles/grant', {
      title_id: grantDialog.value.titleId,
      user_id: grantDialog.value.selectedUserId
    })
    success('授予成功')
    grantDialog.value.show = false
    loadTitles()
  } catch (error) {
    console.error('授予头衔失败', error)
    error(error.response?.data?.error || '授予失败')
  }
}

const revokeTitle = async (titleId, userId) => {
  const confirmed = await confirm('确定要撤销此头衔吗？')
  if (!confirmed) return
  
  try {
    await api.post('/titles/revoke', { title_id: titleId, user_id: userId })
    success('撤销成功')
    loadTitles()
  } catch (error) {
    console.error('撤销头衔失败', error)
    error(error.response?.data?.error || '撤销失败')
  }
}

const handleDeleteTitle = async (title) => {
  const confirmed = await confirm(`确定要删除头衔 "${title.name}" 吗？`)
  if (!confirmed) return
  
  try {
    await api.delete(`/titles/${title.id}`)
    success('删除成功')
    loadTitles()
  } catch (error) {
    console.error('删除头衔失败', error)
    error(error.response?.data?.error || '删除失败')
  }
}

onMounted(() => {
  loadTitles()
  loadUsers()
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