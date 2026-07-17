<template>
  <div class="admin-announcement">
    <div class="layui-card">
      <div class="layui-card-header">
        <i class="fa-solid fa-bullhorn"></i>
        <span>公告管理</span>
      </div>
      <div class="layui-card-body">
        <textarea
          v-model="announcementContent"
          placeholder="请输入公告内容..."
          class="layui-textarea"
          rows="5"
          maxlength="500"
        ></textarea>
        <span class="char-count">{{ announcementContent.length }}/500</span>
        
        <div v-if="announcementContent" class="info-alert">
          <i class="fa-solid fa-circle-info"></i>
          <span>公告将显示在网站首页顶部，内容过长可能会影响显示效果。</span>
        </div>
      </div>
      <div class="layui-card-footer">
        <button
          class="layui-btn"
          @click="saveAnnouncement"
          :disabled="saving || !announcementContent.trim()"
        >
          <i class="fa-solid fa-floppy-disk"></i>保存公告
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { adminAnnouncementApi } from '../../api/admin'
import { success, error } from '../../utils/modal'

const announcementContent = ref('')
const saving = ref(false)

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
.admin-announcement {
  padding: 20px;
}

.layui-card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
  font-weight: 600;
}

.info-alert {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px;
  background: rgba(30, 159, 255, 0.1);
  color: #1E9FFF;
  border-radius: 6px;
  margin-top: 16px;
  font-size: 14px;
}

.char-count {
  display: block;
  text-align: right;
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.layui-card-footer {
  display: flex;
  justify-content: flex-end;
  padding: 16px 20px;
}
</style>