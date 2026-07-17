<template>
  <div class="admin-site-config">
    <div class="layui-card">
      <div class="layui-card-header">
        <i class="fa-solid fa-globe"></i>
        <span>网站配置</span>
      </div>

      <div class="layui-card-body">
        <div class="section-title">基本配置</div>
        <div class="form-group">
          <label class="form-label">网站标题</label>
          <input
            v-model="siteConfigForm.siteTitle"
            type="text"
            placeholder="校园论坛 - 分享与交流"
            class="layui-input"
            maxlength="100"
          />
          <span class="char-count">{{ siteConfigForm.siteTitle.length }}/100</span>
        </div>

        <div class="divider"></div>

        <div class="section-title">备案信息</div>
        <div class="section-desc">以下字段可选，不填写则不在页面底部显示</div>
        <div class="form-row">
          <div class="form-group">
            <label class="form-label">ICP备案号</label>
            <input
              v-model="siteConfigForm.icpNumber"
              type="text"
              placeholder="京ICP备12345678号"
              class="layui-input"
              maxlength="50"
            />
          </div>
          <div class="form-group">
            <label class="form-label">公安联网备案号</label>
            <input
              v-model="siteConfigForm.publicSecurityNumber"
              type="text"
              placeholder="京公网安备 12345678901234567890号"
              class="layui-input"
              maxlength="50"
            />
          </div>
        </div>

        <div v-if="siteConfigForm.siteTitle" class="success-alert">
          <i class="fa-solid fa-check-circle"></i>
          <span>当前网站标题：<strong>{{ siteConfigForm.siteTitle }}</strong></span>
        </div>
      </div>

      <div class="layui-card-footer">
        <button class="layui-btn layui-btn-primary" @click="resetForm" :disabled="saving">
          <i class="fa-solid fa-rotate-right"></i>重置
        </button>
        <div class="spacer"></div>
        <button class="layui-btn" @click="saveSiteConfig" :disabled="saving || !siteConfigForm.siteTitle.trim()">
          <i class="fa-solid fa-floppy-disk"></i>保存配置
        </button>
      </div>
    </div>
  </div>
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

const saving = ref(false)

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

<style scoped>
.admin-site-config {
  padding: 20px;
}

.layui-card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
  font-weight: 600;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-bottom: 12px;
}

.section-desc {
  font-size: 13px;
  color: #999;
  margin-bottom: 16px;
}

.form-group {
  margin-bottom: 20px;
}

.form-label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin-bottom: 8px;
}

.char-count {
  display: block;
  text-align: right;
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.form-row {
  display: flex;
  gap: 24px;
}

.form-row .form-group {
  flex: 1;
}

.divider {
  height: 1px;
  background: #f0f0f0;
  margin: 24px 0;
}

.success-alert {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px;
  background: rgba(82, 196, 26, 0.1);
  color: #52C41A;
  border-radius: 6px;
  margin-top: 16px;
  font-size: 14px;
}

.layui-card-footer {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 20px;
}

.spacer {
  flex: 1;
}

@media (max-width: 768px) {
  .form-row {
    flex-direction: column;
  }
}
</style>