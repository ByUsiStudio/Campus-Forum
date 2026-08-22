<template>
  <div class="announcement-page">
    <!-- 页面标题 -->
    <div class="page-head">
      <div>
        <h1 class="page-title">公告管理</h1>
        <p class="page-desc">编辑将在网站首页顶部展示的公告内容</p>
      </div>
    </div>

    <!-- 公告编辑卡片 -->
    <div class="editor-card card-surface">
      <div class="card-head">
        <div class="card-head-icon">
          <el-icon :size="20"><Bell /></el-icon>
        </div>
        <div class="card-head-text">
          <div class="card-head-title">编辑公告</div>
          <div class="card-head-desc">支持 Markdown 语法，保存后立即生效</div>
        </div>
      </div>

      <el-form ref="announcementForm" :model="formModel" :rules="rules" class="announcement-form">
        <el-form-item prop="content">
          <el-input
            v-model="announcementContent"
            type="textarea"
            :rows="6"
            maxlength="500"
            show-word-limit
            placeholder="请输入公告内容，例如：# 欢迎加入校园论坛"
            class="announcement-input"
          />
        </el-form-item>

        <el-alert
          v-if="announcementContent"
          type="info"
          :closable="false"
          show-icon
          class="announcement-tip"
        >
          公告将显示在网站首页顶部，内容过长可能会影响显示效果。
        </el-alert>
      </el-form>

      <div class="card-footer">
        <div class="footer-hint">
          <el-icon :size="15"><InfoFilled /></el-icon>
          <span>支持 Markdown：标题、加粗、链接、代码等</span>
        </div>
        <el-button
          type="primary"
          round
          :loading="saving"
          :disabled="!formValid"
          @click="saveAnnouncement"
        >
          <el-icon class="mr-1"><Promotion /></el-icon>
          保存公告
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { Bell, Promotion, InfoFilled } from '@element-plus/icons-vue'
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
.announcement-page {
  display: flex;
  flex-direction: column;
}

.page-head {
  margin-bottom: 20px;
}

.page-title {
  margin: 0 0 4px;
  font-size: 22px;
  font-weight: 700;
  letter-spacing: -0.01em;
  color: var(--campus-text);
}

.page-desc {
  margin: 0;
  font-size: 14px;
  color: var(--campus-text-secondary);
}

/* ---------------- 编辑器卡片 ---------------- */
.editor-card {
  padding: 24px;
  max-width: 860px;
}

.card-head {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid var(--campus-border);
}

.card-head-icon {
  width: 42px;
  height: 42px;
  flex-shrink: 0;
  display: grid;
  place-items: center;
  border-radius: 12px;
  background: var(--campus-primary-soft);
  color: var(--campus-primary);
}

.card-head-title {
  font-size: 16px;
  font-weight: 700;
  color: var(--campus-text);
}

.card-head-desc {
  font-size: 13px;
  color: var(--campus-text-muted);
  margin-top: 2px;
}

.announcement-input {
  font-size: 14px;
}

.announcement-input :deep(.el-textarea__inner) {
  padding: 14px;
  line-height: 1.7;
  border-radius: 12px;
  font-family: inherit;
}

.announcement-tip {
  border-radius: 12px;
  align-items: center;
}

.card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-top: 20px;
  padding-top: 18px;
  border-top: 1px solid var(--campus-border);
  flex-wrap: wrap;
}

.footer-hint {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: var(--campus-text-muted);
}

.mr-1 {
  margin-right: 4px;
}

@media (max-width: 575px) {
  .editor-card {
    padding: 16px;
  }
  .card-footer {
    flex-direction: column;
    align-items: stretch;
  }
  .card-footer .el-button {
    width: 100%;
  }
}
</style>
