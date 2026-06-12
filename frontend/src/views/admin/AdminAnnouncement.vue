<script setup>
import { ref, onMounted } from 'vue'
import { adminAnnouncementApi } from '../../api/admin'

const loading = ref(false)
const saving = ref(false)
const content = ref('')

const loadAnnouncement = async () => {
  loading.value = true
  try {
    const response = await adminAnnouncementApi.getAnnouncement()
    content.value = response.data.content || ''
  } catch (error) {
    console.error('加载公告失败:', error)
  } finally {
    loading.value = false
  }
}

const saveAnnouncement = async () => {
  saving.value = true
  try {
    await adminAnnouncementApi.updateAnnouncement(content.value)
  } catch (error) {
    console.error('保存公告失败:', error)
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadAnnouncement()
})
</script>

<template>
  <v-container fluid class="pa-0">
    <!-- 页面标题 -->
    <div class="mb-6">
      <h1 class="text-h5 font-weight-bold">公告管理</h1>
      <p class="text-body-2 text-grey">编辑网站公告内容</p>
    </div>

    <!-- 公告编辑 -->
    <v-card>
      <v-card-title>公告内容</v-card-title>
      <v-card-text>
        <v-textarea
          v-model="content"
          label="公告内容"
          rows="10"
          :loading="loading"
          placeholder="请输入公告内容..."
        />
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="primary" @click="saveAnnouncement" :loading="saving">
          <v-icon start>mdi-check</v-icon>
          保存公告
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-container>
</template>