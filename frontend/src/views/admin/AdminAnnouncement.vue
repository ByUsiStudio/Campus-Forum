<template>
  <div class="page-container">
    <!-- 公告编辑表单 -->
    <el-card shadow="never" class="mb-4">
      <template #header>
        <div class="card-header">
          <el-icon><Bell /></el-icon>
          <span>公告管理</span>
        </div>
      </template>

      <el-form ref="announcementForm" :model="formModel" :rules="rules">
        <el-form-item prop="content">
          <el-input
            v-model="announcementContent"
            type="textarea"
            :rows="5"
            maxlength="500"
            show-word-limit
            placeholder="请输入公告内容..."
          />
        </el-form-item>

        <el-alert
          v-if="announcementContent"
          type="info"
          :closable="false"
          show-icon
        >
          公告将显示在网站首页顶部，内容过长可能会影响显示效果。
        </el-alert>
      </el-form>

      <div class="card-footer">
        <el-button
          type="primary"
          :loading="saving"
          :disabled="!formValid"
          @click="saveAnnouncement"
        >
          <el-icon><Promotion /></el-icon>
          保存公告
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { adminAnnouncementApi } from '@/api/admin'
import { success, error } from '@/utils/message'

const announcementContent = ref('')
const announcementForm = ref(null)
const saving = ref(false)

const activeRule = (rule, value, callback) => {
  if (!value || !value.trim()) {
    callback(new Error('此字段为必填项'))
  } else {
    callback()
  }
}

const rules = {
  content: [{ validator: activeRule, trigger: 'blur' }]
}

const formModel = computed(() => ({
  content: announcementContent.value
}))

const formValid = computed(() => {
  return !!(announcementContent.value && announcementContent.value.trim())
})

const loadAnnouncement = async () => {
  try {
    const response = await adminAnnouncementApi.getAnnouncement()
    announcementContent.value = response.data.content || ''
  } catch (err) {
    console.error('加载公告失败', err)
    error('加载公告失败')
  }
}

const saveAnnouncement = async () => {
  saving.value = true
  try {
    await adminAnnouncementApi.updateAnnouncement({ content: announcementContent.value })
    success('公告保存成功')
  } catch (err) {
    console.error('保存公告失败', err)
    error(err.response?.data?.error || '保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadAnnouncement()
})
</script>

<style scoped>
.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.card-footer {
  margin-top: 16px;
  display: flex;
  justify-content: flex-start;
}
</style>
