<template>
  <div class="pa-4 pa-md-6">
    <!-- 网站配置表单 -->
    <el-card shadow="never" class="page-container">
      <template #header>
        <div class="card-header">
          <el-icon :size="20" class="card-header-icon"><OfficeBuilding /></el-icon>
          <span>网站配置</span>
        </div>
      </template>

      <el-form ref="siteForm" :model="siteConfigForm" :rules="rules" label-position="top">
        <!-- 网站基本配置 -->
        <div class="form-section-title">基本配置</div>
        <el-form-item label="网站标题" prop="siteTitle" class="mb-4">
          <el-input
            v-model="siteConfigForm.siteTitle"
            placeholder="校园论坛 - 分享与交流"
            :prefix-icon="Document"
            clearable
            show-word-limit
            maxlength="100"
          />
        </el-form-item>

        <el-divider class="my-4" />

        <!-- 备案信息配置 -->
        <div class="form-section-title">备案信息</div>
        <div class="form-section-desc mb-3">以下字段可选，不填写则不在页面底部显示</div>
        <el-row :gutter="16">
          <el-col :xs="24" :sm="12">
            <el-form-item label="ICP备案号" prop="icpNumber">
              <el-input
                v-model="siteConfigForm.icpNumber"
                placeholder="京ICP备12345678号"
                :prefix-icon="CircleCheck"
                clearable
                show-word-limit
                maxlength="50"
              />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12">
            <el-form-item label="公安联网备案号" prop="publicSecurityNumber">
              <el-input
                v-model="siteConfigForm.publicSecurityNumber"
                placeholder="京公网安备 12345678901234567890号"
                :prefix-icon="UserFilled"
                clearable
                show-word-limit
                maxlength="50"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-alert
          v-if="siteConfigForm.siteTitle"
          type="success"
          :closable="false"
          class="mt-4"
          :icon="CircleCheck"
        >
          当前网站标题：<strong>{{ siteConfigForm.siteTitle }}</strong>
        </el-alert>
      </el-form>

      <template #footer>
        <div class="form-actions">
          <el-button
            type="warning"
            plain
            @click="resetForm"
            :disabled="saving"
          >
            <el-icon class="mr-1"><Refresh /></el-icon>
            重置
          </el-button>
          <el-button
            type="primary"
            @click="saveSiteConfig"
            :loading="saving"
            :disabled="!formValid"
          >
            <el-icon class="mr-1"><Check /></el-icon>
            保存配置
          </el-button>
        </div>
      </template>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { OfficeBuilding, Document, CircleCheck, UserFilled, Refresh, Check } from '@element-plus/icons-vue'
import { adminSiteConfigApi } from '../../api/admin'
import { success, error } from '@/utils/message'

const siteConfigForm = reactive({
  siteTitle: '',
  icpNumber: '',
  publicSecurityNumber: ''
})

const siteForm = ref(null)
const saving = ref(false)

const rules = {
  siteTitle: [
    { required: true, message: '此字段为必填项', trigger: 'blur' }
  ]
}

const originalConfig = ref({
  siteTitle: '',
  icpNumber: '',
  publicSecurityNumber: ''
})

// 表单整体校验结果，用于控制保存按钮
const formValid = computed(() => !!(siteConfigForm.siteTitle && siteConfigForm.siteTitle.trim()))

const loadSiteConfig = async () => {
  try {
    const response = await adminSiteConfigApi.getConfig()
    siteConfigForm.siteTitle = response.data.site_title || ''
    siteConfigForm.icpNumber = response.data.icp_number || ''
    siteConfigForm.publicSecurityNumber = response.data.public_security_number || ''
    // 保存原始配置用于重置
    originalConfig.value = { ...siteConfigForm }
  } catch (err) {
    console.error('加载网站配置失败', err)
    error('加载网站配置失败')
  }
}

const resetForm = () => {
  siteConfigForm.siteTitle = originalConfig.value.siteTitle
  siteConfigForm.icpNumber = originalConfig.value.icpNumber
  siteConfigForm.publicSecurityNumber = originalConfig.value.publicSecurityNumber
}

const saveSiteConfig = async () => {
  if (siteForm.value) {
    const valid = await siteForm.value.validate().catch(() => false)
    if (!valid) return
  }
  saving.value = true
  try {
    await adminSiteConfigApi.updateConfig({
      site_title: siteConfigForm.siteTitle,
      icp_number: siteConfigForm.icpNumber,
      public_security_number: siteConfigForm.publicSecurityNumber
    })
    success('网站配置保存成功')
    if (siteConfigForm.siteTitle) {
      document.title = siteConfigForm.siteTitle
      window.dispatchEvent(new CustomEvent('site-title-updated', {
        detail: siteConfigForm.siteTitle
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
.card-header {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 1.1rem;
  font-weight: 700;
}

.card-header-icon {
  color: var(--el-color-primary);
}

.form-section-title {
  font-size: 0.9rem;
  font-weight: 700;
  margin-bottom: 12px;
}

.form-section-desc {
  font-size: 0.8rem;
  color: var(--el-text-color-secondary);
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 12px;
}

.mr-1 {
  margin-right: 4px;
}

.mb-3 {
  margin-bottom: 12px;
}

.mb-4 {
  margin-bottom: 16px;
}

.mt-4 {
  margin-top: 16px;
}
</style>
