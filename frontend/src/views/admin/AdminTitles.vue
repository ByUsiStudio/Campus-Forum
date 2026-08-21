<template>
  <v-container fluid class="pa-4 pa-md-6">
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

  <v-dialog v-model="grantDialog.show" :max-width="isMobile ? '100%' : '500'" :fullscreen="isMobile" transition="dialog-bottom-transition">
    <v-card class="dialog-card">
      <v-card-title class="dialog-title">
        <v-avatar color="primary" size="40" class="mr-3">
          <v-icon color="white" size="20">mdi-medal</v-icon>
        </v-avatar>
        <span>授予头衔</span>
        <v-spacer />
        <v-btn v-if="isMobile" icon variant="text" @click="grantDialog.show = false">
          <v-icon>mdi-close</v-icon>
        </v-btn>
      </v-card-title>
      <v-card-text class="dialog-body">
        <v-form ref="grantForm" v-model="formValid">
          <v-select
            v-model="grantDialog.selectedUserId"
            :items="usersForSelect"
            item-title="display_name"
            item-value="id"
            label="选择用户"
            variant="outlined"
            density="comfortable"
            :rules="[rules.required]"
            prepend-inner-icon="mdi-account"
            class="mb-2"
          >
            <template #label>
              <span class="text-body-2">选择用户</span>
            </template>
            <template #item="{ props, item }">
              <v-list-item v-bind="props" :title="item.raw.display_name">
                <template #prepend>
                  <v-avatar size="32" class="mr-3">
                    <v-img :src="item.raw.avatar || '/default-avatar.png'"></v-img>
                  </v-avatar>
                </template>
                <template #subtitle>
                  <span class="text-caption">@{{ item.raw.username }}</span>
                </template>
              </v-list-item>
            </template>
          </v-select>
        </v-form>
      </v-card-text>
      <v-card-actions class="dialog-actions" :class="{ 'flex-column': isMobile }">
        <v-btn variant="text" @click="grantDialog.show = false" :block="isMobile" class="mb-2 mb-md-0">
          取消
        </v-btn>
        <v-btn
          color="primary"
          variant="flat"
          @click="handleGrant"
          :disabled="!formValid"
          :block="isMobile"
        >
          <v-icon class="mr-1">mdi-check</v-icon>
          确认授予
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import TitlesPanel from './TitlesPanel.vue'
import { adminTitleApi, adminUserApi } from '../../api/admin'
import { confirm, success, error } from '../../utils/modal'

const titles = ref([])
const users = ref([])
const loading = ref(true)
const grantForm = ref(null)
const formValid = ref(false)
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

const rules = {
  required: v => !!v || '此字段为必填项'
}

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

const addTitle = async () => {
  const titleName = prompt('请输入头衔名称：')
  if (!titleName) return

  try {
    await adminTitleApi.createTitle({ name: titleName })
    success('添加成功')
    loadTitles()
  } catch (err) {
    console.error('添加头衔失败', err)
    error(err.response?.data?.error || '添加失败')
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

const revokeTitle = async (titleId, userId) => {
  const confirmed = await confirm('确定要撤销此头衔吗？')
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
  const confirmed = await confirm(`确定要删除头衔 "${title.name}" 吗？此操作不可恢复。`)
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
.dialog-card {
  border-radius: 20px !important;
  overflow: hidden;
}

@media (max-width: 599px) {
  .dialog-card {
    border-radius: 0 !important;
  }
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

:deep(.v-select .v-field) {
  border-radius: 12px;
}
</style>