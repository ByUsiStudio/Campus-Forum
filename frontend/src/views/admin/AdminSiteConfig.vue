<template>
  <v-container fluid class="pa-6">
    <v-card>
      <v-card-title>
        <v-icon class="mr-2">mdi-globe</v-icon>
        网站配置
      </v-card-title>
      <v-card-text>
        <v-text-field
          v-model="siteConfigForm.siteTitle"
          label="网站标题"
          variant="outlined"
          class="mb-4"
        ></v-text-field>
      </v-card-text>
      <v-card-actions>
        <v-btn color="primary" variant="flat" @click="saveSiteConfig">
          保存配置
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-container>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'
import { showSuccess, showError } from '../utils/modal'

const siteConfigForm = ref({
  siteTitle: ''
})

const loadSiteConfig = async () => {
  try {
    const response = await api.get('/site-config')
    siteConfigForm.value.siteTitle = response.data.site_title || ''
  } catch (error) {
    console.error('加载网站配置失败', error)
  }
}

const saveSiteConfig = async () => {
  try {
    await api.put('/site-config', { site_title: siteConfigForm.value.siteTitle })
    showSuccess('保存成功')
    if (siteConfigForm.value.siteTitle) {
      document.title = siteConfigForm.value.siteTitle
    }
  } catch (error) {
    console.error('保存网站配置失败', error)
    showError(error.response?.data?.error || '保存失败')
  }
}

onMounted(() => {
  loadSiteConfig()
})
</script>