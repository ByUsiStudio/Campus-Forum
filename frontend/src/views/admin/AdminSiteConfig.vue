<script setup>
import { ref, onMounted } from 'vue'
import { adminApi, siteApi } from '../../api'

const loading = ref(false)
const saving = ref(false)
const config = ref({
  site_title: '校园论坛',
  site_description: '',
  site_keywords: ''
})

const loadConfig = async () => {
  loading.value = true
  try {
    const response = await adminSiteConfigApi.getConfig()
    config.value = response.data || config.value
  } catch (error) {
    console.error('加载配置失败:', error)
  } finally {
    loading.value = false
  }
}

const saveConfig = async () => {
  saving.value = true
  try {
    await adminApi.updateSiteConfig(config.value)
  } catch (error) {
    console.error('保存配置失败:', error)
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadConfig()
})
</script>

<template>
  <v-container fluid class="pa-0">
    <!-- 页面标题 -->
    <div class="mb-6">
      <h1 class="text-h5 font-weight-bold">网站配置</h1>
      <p class="text-body-2 text-grey">配置网站基本信息</p>
    </div>

    <!-- 配置表单 -->
    <v-card>
      <v-card-title>基本配置</v-card-title>
      <v-card-text>
        <v-text-field v-model="config.site_title" label="网站标题" class="mb-3" />
        <v-textarea v-model="config.site_description" label="网站描述" rows="3" class="mb-3" />
        <v-text-field v-model="config.site_keywords" label="网站关键词" class="mb-3" />
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="primary" @click="saveConfig" :loading="saving">
          <v-icon start>mdi-check</v-icon>
          保存配置
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-container>
</template>