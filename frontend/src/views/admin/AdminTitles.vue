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
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import TitlesPanel from './TitlesPanel.vue'
import { adminTitleApi, adminUserApi } from '../../api/admin'
import { confirm, success, error } from '@/utils/message'
import { jcCloseAll, jcOpenHtml } from '@/utils/jcu'

const titles = ref([])
const users = ref([])
const loading = ref(true)
const isMobile = ref(false)

const grantDialog = ref({
  titleId: null,
  selectedUserId: null
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
    titleId,
    selectedUserId: isObject && payload.user_id ? payload.user_id : null
  }
  openGrantDialog()
}

const openGrantDialog = () => {
  const esc = (s) => String(s ?? '').replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;')
  const options = users.value
    .map((u) => {
      const selected = grantDialog.value.selectedUserId && String(u.id) === String(grantDialog.value.selectedUserId) ? ' selected' : ''
      return `<option value="${u.id}"${selected}>${esc(u.display_name)}（@${esc(u.username)}）</option>`
    })
    .join('')

  const content = `
    <div style="margin-bottom:12px;"><div style="margin:0 0 4px;font-weight:600;font-size:13px;color:var(--jc-text,#333);">选择用户</div>
      ${options ? `<select data-grant-user class="jc-modal__input" style="width:100%;">${options}</select>` : '<div style="font-size:13px;color:#64748b;">暂无可用用户</div>'}
    </div>
  `

  jcOpenHtml({
    title: '授予头衔',
    content,
    width: isMobile.value ? '100%' : 500,
    size: isMobile.value ? 'sm' : 'md',
    buttons: [
      { text: '取消', type: 'default', action: () => jcCloseAll() },
      {
        text: '确认授予',
        type: 'primary',
        action: (inst) => {
          const sel = inst.modalContent.querySelector('[data-grant-user]')
          if (!sel || !sel.value) {
            error('请选择用户')
            return
          }
          handleGrant(Number(sel.value))
        }
      }
    ]
  })
}

const handleGrant = async (userId) => {
  if (!userId) {
    error('请选择用户')
    return
  }

  try {
    await adminTitleApi.grantTitle({
      title_id: grantDialog.value.titleId,
      user_id: userId
    })
    success('授予成功')
    jcCloseAll()
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
</style>
