<template>
  <div class="admin-titles">
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

    <div v-if="grantDialog.show" class="modal-overlay" @click="grantDialog.show = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <div class="flex items-center gap-3">
            <div class="avatar" style="width: 40px; height: 40px; background: #1E9FFF; border-radius: 50%; display: flex; align-items: center; justify-content: center;">
              <i class="fa-solid fa-medal" style="color: white; font-size: 20px;"></i>
            </div>
            <span style="font-size: 18px; font-weight: 600;">授予头衔</span>
          </div>
          <button v-if="isMobile" class="layui-btn layui-btn-primary layui-btn-sm" @click="grantDialog.show = false">
            <i class="fa-solid fa-xmark"></i>
          </button>
        </div>
        <div class="modal-body">
          <div class="layui-form-item">
            <label class="layui-form-label">选择用户</label>
            <div class="layui-input-block">
              <select v-model="grantDialog.selectedUserId" class="layui-select">
                <option value="">请选择用户</option>
                <option v-for="user in usersForSelect" :key="user.id" :value="user.id">{{ user.display_name }}</option>
              </select>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="layui-btn layui-btn-primary" @click="grantDialog.show = false">取消</button>
          <button class="layui-btn" @click="handleGrant" :disabled="!grantDialog.selectedUserId">
            <i class="fa-solid fa-check mr-1"></i>
            确认授予
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import TitlesPanel from './TitlesPanel.vue'
import { adminTitleApi, adminUserApi } from '../../api/admin'
import { confirm, success, error } from '../../utils/modal'

const titles = ref([])
const users = ref([])
const loading = ref(true)
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
.admin-titles {
  padding: 24px;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}

.modal-content {
  background: #fff;
  border-radius: 8px;
  width: 90%;
  max-width: 500px;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px;
  background: linear-gradient(135deg, #f8f9ff 0%, #fff 100%);
  border-radius: 8px 8px 0 0;
}

.modal-body {
  padding: 20px;
}

.modal-footer {
  padding: 15px 20px;
  border-top: 1px solid #E8E8E8;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

@media (max-width: 599px) {
  .modal-content {
    width: 100%;
    height: 100%;
    border-radius: 0;
    display: flex;
    flex-direction: column;
  }
  .modal-body {
    flex: 1;
    overflow-y: auto;
  }
}
</style>