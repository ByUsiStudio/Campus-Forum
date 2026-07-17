<template>
  <div class="admin-smtp-config">
    <div class="layui-card">
      <div class="layui-card-header">
        <i class="fa-solid fa-envelope"></i>
        <span>SMTP 配置</span>
      </div>

      <div class="layui-card-body">
        <div class="form-row">
          <div class="form-group">
            <label class="form-label">SMTP 主机</label>
            <input
              v-model="smtpConfigForm.host"
              type="text"
              placeholder="smtp.example.com"
              class="layui-input"
            />
          </div>
          <div class="form-group">
            <label class="form-label">SMTP 端口</label>
            <input
              v-model.number="smtpConfigForm.port"
              type="number"
              placeholder="587"
              class="layui-input"
            />
          </div>
        </div>

        <div class="form-row">
          <div class="form-group">
            <label class="form-label">用户名 / 邮箱</label>
            <input
              v-model="smtpConfigForm.username"
              type="text"
              placeholder="your-email@example.com"
              class="layui-input"
            />
          </div>
          <div class="form-group">
            <label class="form-label">密码 / 授权码</label>
            <div class="password-input-wrapper">
              <input
                v-model="smtpConfigForm.password"
                :type="showPassword ? 'text' : 'password'"
                placeholder="请输入密码"
                class="layui-input"
                @input="onPasswordChange"
              />
              <button class="password-toggle" @click="showPassword = !showPassword">
                <i class="fa-solid" :class="showPassword ? 'fa-eye-slash' : 'fa-eye'"></i>
              </button>
            </div>
            <span v-if="passwordSet && !passwordChanged" class="password-hint">密码已设置（不可见），输入新密码以更新</span>
          </div>
        </div>

        <div class="form-row">
          <div class="form-group">
            <label class="form-label">发件人邮箱</label>
            <input
              v-model="smtpConfigForm.from"
              type="text"
              placeholder="noreply@example.com"
              class="layui-input"
            />
          </div>
          <div class="form-group">
            <label class="form-label">发件人名称</label>
            <input
              v-model="smtpConfigForm.fromName"
              type="text"
              placeholder="校园论坛"
              class="layui-input"
            />
          </div>
        </div>

        <div v-if="smtpConfigForm.host && smtpConfigForm.port" class="info-alert">
          <i class="fa-solid fa-circle-info"></i>
          <template v-if="passwordSet && !passwordChanged">
            密码已设置但不可见。如需测试或更新密码，请先输入新密码。
          </template>
          <template v-else>
            测试配置前请确保所有参数填写正确，测试邮件将发送到您的发件人邮箱。
          </template>
        </div>
      </div>

      <div class="layui-card-footer">
        <div class="spacer"></div>
        <button class="layui-btn layui-btn-secondary" @click="testSmtpConfig" :disabled="testing">
          <i class="fa-solid fa-flask"></i>测试配置
        </button>
        <button class="layui-btn" @click="saveSmtpConfig" :disabled="saving || !isFormValid">
          <i class="fa-solid fa-floppy-disk"></i>保存配置
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { adminSiteConfigApi } from '../../api/admin'
import { success, error } from '../../utils/modal'

const PASSWORD_PLACEHOLDER = '••••••••••••••••'

const smtpConfigForm = ref({
  host: '',
  port: 587,
  username: '',
  password: '',
  from: '',
  fromName: ''
})

const showPassword = ref(false)
const saving = ref(false)
const testing = ref(false)
const passwordSet = ref(false)
const passwordChanged = ref(false)

const isFormValid = computed(() => {
  return smtpConfigForm.value.host && 
         smtpConfigForm.value.port && 
         smtpConfigForm.value.username && 
         smtpConfigForm.value.from && 
         smtpConfigForm.value.fromName
})

const loadSmtpConfig = async () => {
  try {
    const response = await adminSiteConfigApi.getConfig()
    passwordSet.value = response.data.smtp_password_set || false
    smtpConfigForm.value = {
      host: response.data.smtp_host || '',
      port: response.data.smtp_port || 587,
      username: response.data.smtp_username || '',
      password: passwordSet.value ? PASSWORD_PLACEHOLDER : '',
      from: response.data.smtp_from || '',
      fromName: response.data.smtp_from_name || ''
    }
    passwordChanged.value = false
  } catch (err) {
    console.error('加载SMTP配置失败', err)
    error('加载SMTP配置失败')
  }
}

const onPasswordChange = () => {
  passwordChanged.value = true
}

const saveSmtpConfig = async () => {
  saving.value = true
  try {
    const updateData = {
      smtp_host: smtpConfigForm.value.host,
      smtp_port: smtpConfigForm.value.port,
      smtp_username: smtpConfigForm.value.username,
      smtp_from: smtpConfigForm.value.from,
      smtp_from_name: smtpConfigForm.value.fromName
    }
    
    if (passwordChanged.value && smtpConfigForm.value.password !== PASSWORD_PLACEHOLDER) {
      updateData.smtp_password = smtpConfigForm.value.password
    } else {
      updateData.smtp_password = ''
    }
    
    await adminSiteConfigApi.updateConfig(updateData)
    success('SMTP配置保存成功')
    loadSmtpConfig()
  } catch (err) {
    console.error('保存SMTP配置失败', err)
    error(err.response?.data?.error || '保存失败')
  } finally {
    saving.value = false
  }
}

const testSmtpConfig = async () => {
  testing.value = true
  try {
    if (!passwordChanged.value && smtpConfigForm.value.password === PASSWORD_PLACEHOLDER) {
      error('密码已设置但不可见，无法测试。请先输入新密码后再测试。')
      testing.value = false
      return
    }
    
    await adminSiteConfigApi.testSmtp({
      smtp_host: smtpConfigForm.value.host,
      smtp_port: smtpConfigForm.value.port,
      smtp_username: smtpConfigForm.value.username,
      smtp_password: smtpConfigForm.value.password,
      smtp_from: smtpConfigForm.value.from,
      smtp_to: smtpConfigForm.value.from
    })
    success('测试邮件发送成功！请检查您的收件箱。')
  } catch (err) {
    console.error('测试SMTP配置失败', err)
    error(err.response?.data?.error || '测试失败，请检查配置是否正确')
  } finally {
    testing.value = false
  }
}

onMounted(() => {
  loadSmtpConfig()
})
</script>

<style scoped>
.admin-smtp-config {
  padding: 20px;
}

.layui-card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
  font-weight: 600;
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

.form-row {
  display: flex;
  gap: 24px;
}

.form-row .form-group {
  flex: 1;
}

.password-input-wrapper {
  position: relative;
}

.password-toggle {
  position: absolute;
  right: 10px;
  top: 50%;
  transform: translateY(-50%);
  background: transparent;
  border: none;
  color: #999;
  font-size: 16px;
  cursor: pointer;
  z-index: 1;

  &:hover {
    color: #333;
  }
}

.password-hint {
  display: block;
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.info-alert {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  padding: 12px;
  background: rgba(30, 159, 255, 0.1);
  color: #1E9FFF;
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