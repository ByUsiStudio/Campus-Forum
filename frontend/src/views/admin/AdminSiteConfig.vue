<template>
  <div class="site-config-page">
    <!-- 页面标题 -->
    <div class="page-head">
      <div>
        <h1 class="page-title">网站配置</h1>
        <p class="page-desc">管理站点基本信息与备案信息</p>
      </div>
    </div>

    <!-- 网站配置表单 -->
    <div class="card-surface config-card">
      <div class="card-head">
        <div class="card-head-icon">
          <el-icon :size="20"><OfficeBuilding /></el-icon>
        </div>
        <div class="card-head-text">
          <div class="card-head-title">基本配置</div>
          <div class="card-head-desc">以下设置将影响前台展示</div>
        </div>
      </div>

      <el-form ref="siteForm" :model="siteConfigForm" :rules="rules" label-position="top" class="config-form">
        <!-- 网站基本配置 -->
        <el-form-item label="网站标题" prop="siteTitle">
          <el-input
            v-model="siteConfigForm.siteTitle"
            placeholder="校园论坛 - 分享与交流"
            :prefix-icon="Document"
            clearable
            show-word-limit
            maxlength="100"
          />
        </el-form-item>

        <el-divider class="section-divider" />

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
          show-icon
        >
          当前网站标题：<strong>{{ siteConfigForm.siteTitle }}</strong>
        </el-alert>
      </el-form>

      <div class="form-actions">
        <el-button
          type="warning"
          plain
          round
          @click="resetForm"
          :disabled="saving"
        >
          <el-icon class="mr-1"><Refresh /></el-icon>
          重置
        </el-button>
        <el-button
          type="primary"
          round
          @click="saveSiteConfig"
          :loading="saving"
          :disabled="!formValid"
        >
          <el-icon class="mr-1"><Check /></el-icon>
          保存配置
        </el-button>
      </div>
    </div>
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
.site-config-page {
  display: flex;
  flex-direction: column;
}

.page-head {
  margin-bottom: 20px;
}

.page-title {
  margin: 0 0 4px;
  font-size: 22px;
  font-weight: 700;
  letter-spacing: -0.01em;
  color: var(--campus-text);
}

.page-desc {
  margin: 0;
  font-size: 14px;
  color: var(--campus-text-secondary);
}

/* ---------------- 配置卡片 ---------------- */
.config-card {
  padding: 24px;
  max-width: 860px;
}

.card-head {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid var(--campus-border);
}

.card-head-icon {
  width: 42px;
  height: 42px;
  flex-shrink: 0;
  display: grid;
  place-items: center;
  border-radius: 12px;
  background: var(--campus-primary-soft);
  color: var(--campus-primary);
}

.card-head-title {
  font-size: 16px;
  font-weight: 700;
  color: var(--campus-text);
}

.card-head-desc {
  font-size: 13px;
  color: var(--campus-text-muted);
  margin-top: 2px;
}

.config-form {
  max-width: 720px;
}

.section-divider {
  margin: 20px 0;
}

.form-section-title {
  font-size: 15px;
  font-weight: 700;
  color: var(--campus-text);
  margin-bottom: 8px;
}

.form-section-desc {
  font-size: 13px;
  color: var(--campus-text-muted);
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 12px;
  margin-top: 20px;
  padding-top: 18px;
  border-top: 1px solid var(--campus-border);
}

.mt-4 {
  margin-top: 16px;
}

.mr-1 {
  margin-right: 4px;
}

.mb-3 {
  margin-bottom: 12px;
}

@media (max-width: 575px) {
  .config-card {
    padding: 16px;
  }
  .form-actions {
    flex-direction: column-reverse;
  }
  .form-actions .el-button {
    width: 100%;
  }
}
</style>
