<template>
  <v-container fluid class="pa-6">
    <v-card>
      <v-card-title>
        <v-icon class="mr-2">mdi-bullhorn</v-icon>
        公告管理
      </v-card-title>
      <v-card-text>
        <v-text-field
          v-model="announcementContent"
          label="公告内容"
          variant="outlined"
          multiline
          rows="4"
          class="mb-4"
        ></v-text-field>
      </v-card-text>
      <v-card-actions>
        <v-btn color="primary" variant="flat" @click="saveAnnouncement">
          保存公告
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-container>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'
import { showSuccess, showError } from '../utils/modal'

const announcementContent = ref('')

const loadAnnouncement = async () => {
  try {
    const response = await api.get('/announcement')
    announcementContent.value = response.data.content || ''
  } catch (error) {
    console.error('加载公告失败', error)
  }
}

const saveAnnouncement = async () => {
  try {
    await api.put('/announcement', { content: announcementContent.value })
    showSuccess('保存成功')
  } catch (error) {
    console.error('保存公告失败', error)
    showError(error.response?.data?.error || '保存失败')
  }
}

onMounted(() => {
  loadAnnouncement()
})
</script>