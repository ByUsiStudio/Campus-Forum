<template>
  <v-container fluid class="pa-4 pa-md-6">
    <!-- 网站配置表单 -->
    <v-card variant="flat" rounded="lg">
      <v-card-title class="pb-2">
        <v-icon start size="20">mdi-globe</v-icon>
        网站配置
      </v-card-title>

      <v-card-text>
        <v-form ref="siteForm" v-model="formValid">
          <!-- 网站基本配置 -->
          <div class="text-subtitle-2 font-weight-bold mb-3">基本配置</div>
          <v-text-field
            v-model="siteConfigForm.siteTitle"
            label="网站标题"
            placeholder="校园论坛 - 分享与交流"
            variant="outlined"
            density="compact"
            :rules="[rules.required]"
            prepend-inner-icon="mdi-web"
            clearable
            counter="100"
            maxlength="100"
            class="mb-4"
          />

          <v-divider class="my-4" />

          <!-- 备案信息配置 -->
          <div class="text-subtitle-2 font-weight-bold mb-2">备案信息</div>
          <div class="text-caption text-medium-emphasis mb-3">
            以下字段可选，不填写则不在页面底部显示
          </div>
          <v-row dense>
            <v-col cols="12" sm="6">
              <v-text-field
                v-model="siteConfigForm.icpNumber"
                label="ICP备案号"
                placeholder="京ICP备12345678号"
                variant="outlined"
                density="compact"
                prepend-inner-icon="mdi-shield-check"
                clearable
                counter="50"
                maxlength="50"
              />
            </v-col>
            <v-col cols="12" sm="6">
              <v-text-field
                v-model="siteConfigForm.publicSecurityNumber"
                label="公安联网备案号"
                placeholder="京公网安备 12345678901234567890号"
                variant="outlined"
                density="compact"
                prepend-inner-icon="mdi-police-badge"
                clearable
                counter="50"
                maxlength="50"
              />
            </v-col>
          </v-row>

          <v-alert
            v-if="siteConfigForm.siteTitle"
            type="success"
            variant="tonal"
            density="compact"
            class="mt-4"
            icon="mdi-check-circle"
          >
            当前网站标题：<strong>{{ siteConfigForm.siteTitle }}</strong>
          </v-alert>
        </v-form>
      </v-card-text>

      <v-card-actions class="pa-4">
        <v-btn
          color="warning"
          variant="outlined"
          @click="resetForm"
          :disabled="saving"
        >
          <v-icon start>mdi-refresh</v-icon>
          重置
        </v-btn>
        <v-spacer />
        <v-btn
          color="primary"
          variant="flat"
          @click="saveSiteConfig"
          :loading="saving"
          :disabled="!formValid"
        >
          <v-icon start>mdi-content-save</v-icon>
          保存配置
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-container>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { adminSiteConfigApi } from '../../api/admin'
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

const originalConfig = ref({
  siteTitle: '',
  icpNumber: '',
  publicSecurityNumber: ''
})

const loadSiteConfig = async () => {
  try {
    const response = await adminSiteConfigApi.getConfig()
    siteConfigForm.value.siteTitle = response.data.site_title || ''
    siteConfigForm.value.icpNumber = response.data.icp_number || ''
    siteConfigForm.value.publicSecurityNumber = response.data.public_security_number || ''
    // 保存原始配置用于重置
    originalConfig.value = { ...siteConfigForm.value }
  } catch (err) {
    console.error('加载网站配置失败', err)
    error('加载网站配置失败')
  }
}

const resetForm = () => {
  siteConfigForm.value = { ...originalConfig.value }
}

const saveSiteConfig = async () => {
  saving.value = true
  try {
    await adminSiteConfigApi.updateConfig({
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