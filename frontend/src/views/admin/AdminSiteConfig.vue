<template>
  <v-container fluid class="pa-6">
    <v-card class="pa-6">
      <v-card-title class="d-flex align-center mb-6">
        <v-avatar color="primary" size="44" class="mr-3">
          <v-icon color="white">mdi-globe</v-icon>
        </v-avatar>
        <div>
          <div class="text-h5 font-weight-bold">网站配置</div>
          <div class="text-caption text-medium-emphasis">设置网站基本信息</div>
        </div>
      </v-card-title>

      <v-divider class="mb-6"></v-divider>

      <v-form ref="siteForm" v-model="formValid">
        <v-row>
          <v-col cols="12" md="8">
            <v-text-field
              v-model="siteConfigForm.siteTitle"
              label="网站标题"
              placeholder="校园论坛 - 分享与交流"
              variant="outlined"
              density="comfortable"
              :rules="[rules.required]"
              prepend-inner-icon="mdi-web"
              clearable
              class="mb-4"
            >
              <template #label>
                <span class="text-body-2">网站标题</span>
              </template>
              <template #hint>
                <span class="text-caption">显示在浏览器标签页和网站顶部</span>
              </template>
            </v-text-field>
          </v-col>
        </v-row>

        <v-divider class="my-4"></v-divider>

        <div class="text-body-2 font-weight-bold mb-4">备案信息</div>
        <v-row>
          <v-col cols="12" md="6">
            <v-text-field
              v-model="siteConfigForm.icpNumber"
              label="ICP备案号"
              placeholder="例如：京ICP备XXXXXXXXXXXX号"
              variant="outlined"
              density="comfortable"
              prepend-inner-icon="mdi-shield-check"
              clearable
              class="mb-4"
            >
              <template #label>
                <span class="text-body-2">ICP备案号</span>
              </template>
              <template #hint>
                <span class="text-caption">可选，留空则不显示</span>
              </template>
            </v-text-field>
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field
              v-model="siteConfigForm.publicSecurityNumber"
              label="公安联网备案号"
              placeholder="例如：12345678901234567890"
              variant="outlined"
              density="comfortable"
              prepend-inner-icon="mdi-police-badge"
              clearable
              class="mb-4"
            >
              <template #label>
                <span class="text-body-2">公安联网备案号</span>
              </template>
              <template #hint>
                <span class="text-caption">可选，留空则不显示</span>
              </template>
            </v-text-field>
          </v-col>
        </v-row>

        <v-alert
          v-if="siteConfigForm.siteTitle"
          type="success"
          variant="tonal"
          density="compact"
          class="mb-4"
          icon="mdi-check-circle"
        >
          <span class="text-body-2">
            当前网站标题：<strong>{{ siteConfigForm.siteTitle }}</strong>
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
          @click="saveSiteConfig"
          :loading="saving"
          :disabled="!formValid"
        >
          <v-icon class="mr-2">mdi-content-save</v-icon>
          保存配置
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-container>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../../api'
import { success, error } from '../../utils/modal'

const siteConfigForm = ref({
  siteTitle: '',
  icpNumber: '',
  publicSecurityNumber: ''
})

const siteForm = ref(null)
const formValid = ref(false)
const saving = ref(false)

const rules = {
  required: v => !!v || '此字段为必填项'
}

const loadSiteConfig = async () => {
  try {
    const response = await api.get('/site-config')
    siteConfigForm.value.siteTitle = response.data.site_title || ''
    siteConfigForm.value.icpNumber = response.data.icp_number || ''
    siteConfigForm.value.publicSecurityNumber = response.data.public_security_number || ''
  } catch (err) {
    console.error('加载网站配置失败', err)
  }
}

const saveSiteConfig = async () => {
  saving.value = true
  try {
    await api.put('/site-config', {
      site_title: siteConfigForm.value.siteTitle,
      icp_number: siteConfigForm.value.icpNumber,
      public_security_number: siteConfigForm.value.publicSecurityNumber
    })
    success('网站配置保存成功')
    if (siteConfigForm.value.siteTitle) {
      document.title = siteConfigForm.value.siteTitle
      window.dispatchEvent(new CustomEvent('site-title-updated', {
        detail: siteConfigForm.value.siteTitle
      }))
    }
  } catch (err) {
    console.error('保存网站配置失败', err)
    error(err.response?.data?.error || '保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadSiteConfig()
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

:deep(.v-switch .v-switch__track) {
  border-radius: 12px;
}
</style>