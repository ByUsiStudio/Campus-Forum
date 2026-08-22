<template>
  <div>
    <el-card class="page-container" shadow="never">
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
    </el-card>

    <el-dialog
      v-model="grantDialog.show"
      :width="isMobile ? '100%' : '500px'"
      :fullscreen="isMobile"
      class="dialog-card"
      :title="null"
    >
      <template #header>
        <div class="dialog-header">
          <el-icon :size="20" class="dialog-header-icon">
            <Medal />
          </el-icon>
          <span>授予头衔</span>
        </div>
      </template>

      <el-form ref="grantForm" :model="grantDialog" label-position="top">
        <el-form-item label="选择用户" prop="selectedUserId">
          <el-select
            v-model="grantDialog.selectedUserId"
            :placeholder="'请选择用户'"
            clearable
            filterable
            style="width: 100%"
            class="mb-2"
          >
            <el-option
              v-for="u in usersForSelect"
              :key="u.id"
              :label="u.display_name"
              :value="u.id"
            >
              <div class="d-flex align-center">
                <el-avatar :size="28" :src="u.avatar || '/default-avatar.png'" class="mr-2" />
                <div>
                  <div>{{ u.display_name }}</div>
                  <div class="text-body-2 text-medium-emphasis">@{{ u.username }}</div>
                </div>
              </div>
            </el-option>
          </el-select>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="grantDialog.show = false" class="mr-2">取消</el-button>
        <el-button type="primary" :disabled="!grantDialog.selectedUserId" @click="handleGrant">
          <el-icon class="mr-1"><Check /></el-icon>
          确认授予
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { Medal, Check } from '@element-plus/icons-vue'
import TitlesPanel from './TitlesPanel.vue'
import { adminTitleApi, adminUserApi } from '../../api/admin'
import { confirm, success, error } from '@/utils/message'

const titles = ref([])
const users = ref([])
const loading = ref(true)
const grantForm = ref(null)
const isMobile = ref(false)

const grantDialog = ref({
  show: false,
  titleId: null,
  selectedUserId: null
})

const usersForSelect = computed(() => {
  return users.value.map(u => ({
    id: u.id,
    display_name: u.display_name,
    username: u.username,
    avatar: u.avatar
  }))
})

const checkMobile = () => {
  isMobile.value = window.innerWidth < 600
}

const loadTitles = async () => {
  loading.value = true
  try {
    const response = await adminTitleApi.getTitles()
    titles.value = response.data.titles || []
  } catch (err) {
    console.error('加载头衔列表失败', err)
  } finally {
    loading.value = false
  }
}

const loadUsers = async () => {
  try {
    const response = await adminUserApi.getUsers()
    users.value = response.data.users || []
  } catch (err) {
    console.error('加载用户列表失败', err)
  }
}

// 子组件通过 @add-title 传入完整表单 { name, description, color, icon }
const addTitle = async (payload) => {
  const titleData = payload && typeof payload === 'object'
    ? { name: payload.name, description: payload.description, color: payload.color, icon: payload.icon }
    : { name: payload }

  const name = (titleData.name || '').trim()
  if (!name) {
    error('请输入头衔名称')
    return
  }

  try {
    await adminTitleApi.createTitle(titleData)
    success('添加成功')
    loadTitles()
  } catch (err) {
    console.error('添加头衔失败', err)
    error(err.response?.data?.error || '添加失败')
  }
}

// 子组件通过 @grant 传入 grantForm 或触发授予对话框时传入 titleId
const grantTitle = (payload) => {
  const isObject = payload && typeof payload === 'object'
  const titleId = isObject ? payload.title_id : payload
  grantDialog.value = {
    show: true,
    titleId,
    // 面板传入了 user_id 时预选用户
    selectedUserId: isObject && payload.user_id ? payload.user_id : null
  }
}

const handleGrant = async () => {
  if (!grantDialog.value.selectedUserId) {
    error('请选择用户')
    return
  }

  try {
    await adminTitleApi.grantTitle({
      title_id: grantDialog.value.titleId,
      user_id: grantDialog.value.selectedUserId
    })
    success('授予成功')
    grantDialog.value.show = false
    loadTitles()
  } catch (err) {
    console.error('授予头衔失败', err)
    error(err.response?.data?.error || '授予失败')
  }
}

// 子组件通过 @revoke 传入 (userId, titleId)
const revokeTitle = async (userId, titleId) => {
  const confirmed = await confirm('确定要撤销此头衔吗？').catch(() => null)
  if (!confirmed) return

  try {
    await adminTitleApi.revokeTitle({ title_id: titleId, user_id: userId })
    success('撤销成功')
    loadTitles()
  } catch (err) {
    console.error('撤销头衔失败', err)
    error(err.response?.data?.error || '撤销失败')
  }
}

const handleDeleteTitle = async (title) => {
  const confirmed = await confirm(`确定要删除头衔 "${title.name}" 吗？此操作不可恢复。`).catch(() => null)
  if (!confirmed) return

  try {
    await adminTitleApi.deleteTitle(title.id)
    success('删除成功')
    loadTitles()
  } catch (err) {
    console.error('删除头衔失败', err)
    error(err.response?.data?.error || '删除失败')
  }
}

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  loadTitles()
  loadUsers()
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})
</script>

<style scoped>
.page-container {
  border-radius: 12px;
}

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
