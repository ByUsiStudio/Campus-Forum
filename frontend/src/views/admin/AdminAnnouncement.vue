<template>
  <v-container fluid class="pa-6">
    <v-card class="pa-6">
      <v-card-title class="d-flex align-center mb-6">
        <v-avatar color="warning" size="44" class="mr-3">
          <v-icon color="white">mdi-bullhorn</v-icon>
        </v-avatar>
        <div>
          <div class="text-h5 font-weight-bold">公告管理</div>
          <div class="text-caption text-medium-emphasis">发布网站公告信息</div>
        </div>
      </v-card-title>

      <v-divider class="mb-6"></v-divider>

      <v-form ref="announcementForm" v-model="formValid">
        <v-textarea
          v-model="announcementContent"
          label="公告内容"
          placeholder="请输入公告内容..."
          variant="outlined"
          density="comfortable"
          :rules="[rules.required]"
          prepend-inner-icon="mdi-text"
          rows="5"
          counter
          :maxlength="500"
          class="mb-4"
        >
          <template #label>
            <span class="text-body-2">公告内容</span>
          </template>
        </v-textarea>

        <v-alert
          v-if="announcementContent"
          type="info"
          variant="tonal"
          density="compact"
          class="mb-4"
          icon="mdi-information"
        >
          <span class="text-body-2">
            公告将显示在网站首页顶部，内容过长可能会影响显示效果。
          </span>
        </v-alert>
      </v-form>

      <v-divider class="my-6"></v-divider>

      <v-card-actions class="px-0">
        <v-spacer></v-spacer>
        <v-btn
          color="primary"
          variant="flat"
          size="large"
          @click="saveAnnouncement"
          :loading="saving"
          :disabled="!formValid"
        >
          <v-icon class="mr-2">mdi-content-save</v-icon>
          保存公告
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-container>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../../api'
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
    const response = await api.get('/announcement')
    announcementContent.value = response.data.content || ''
  } catch (err) {
    console.error('加载公告失败', err)
  }
}

const saveAnnouncement = async () => {
  saving.value = true
  try {
    await api.put('/announcement', { content: announcementContent.value })
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
:deep(.v-field) {
  border-radius: 12px;
}

:deep(.v-field--outlined .v-field__outline) {
  border-color: rgba(148, 163, 184, 0.3);
}

:deep(.v-field--focused .v-field__outline) {
  border-color: rgb(var(--v-theme-primary));
}

:deep(.v-textarea .v-field__input) {
  min-height: 120px;
}
</style>