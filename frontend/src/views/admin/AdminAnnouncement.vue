<template>
  <v-container fluid class="pa-4 pa-md-6">
    <!-- 公告编辑表单 -->
    <v-card variant="flat" rounded="lg">
      <v-card-title class="pb-2">
        <v-icon start size="20">mdi-bullhorn</v-icon>
        公告管理
      </v-card-title>

      <v-card-text>
        <v-form ref="announcementForm" v-model="formValid">
          <v-textarea
            v-model="announcementContent"
            label="公告内容"
            placeholder="请输入公告内容..."
            variant="outlined"
            density="compact"
            :rules="[rules.required]"
            prepend-inner-icon="mdi-text"
            rows="5"
            counter
            :maxlength="500"
            class="mb-4"
          />

          <v-alert
            v-if="announcementContent"
            type="info"
            variant="tonal"
            density="compact"
            class="mb-4"
            icon="mdi-information"
          >
            公告将显示在网站首页顶部，内容过长可能会影响显示效果。
          </v-alert>
        </v-form>
      </v-card-text>

      <v-card-actions class="pa-4">
        <v-btn
          color="primary"
          variant="flat"
          @click="saveAnnouncement"
          :loading="saving"
          :disabled="!formValid"
        >
          <v-icon start>mdi-content-save</v-icon>
          保存公告
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-container>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { adminAnnouncementApi } from '../../api/admin'
import { success, error } from '../../utils/modal'

const announcementContent = ref('')
const announcementForm = ref(null)
const formValid = ref(false)
const saving = ref(false)

const rules = {
  required: v => !!v || '此字段为必填项'
}

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