<template>
  <div class="smtp-config-page">
    <!-- 页面标题 -->
    <div class="page-head">
      <div>
        <h1 class="page-title">邮件配置</h1>
        <p class="page-desc">配置 SMTP 服务，用于发送系统邮件通知</p>
      </div>
    </div>

    <!-- SMTP配置表单 -->
    <div class="card-surface config-card">
      <div class="card-head">
        <div class="card-head-icon">
          <el-icon :size="20"><Message /></el-icon>
        </div>
        <div class="card-head-text">
          <div class="card-head-title">SMTP 配置</div>
          <div class="card-head-desc">测试配置前请确保所有参数填写正确</div>
        </div>
      </div>

      <el-form ref="smtpForm" :model="smtpConfigForm" :rules="rules" label-position="top" class="config-form">
        <el-row :gutter="16">
          <el-col :xs="24" :sm="12">
            <el-form-item label="SMTP 主机" prop="host">
              <el-input
                v-model="smtpConfigForm.host"
                placeholder="smtp.example.com"
                clearable
                :prefix-icon="Cpu"
              />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12">
            <el-form-item label="SMTP 端口" prop="port">
              <el-input-number
                v-model="smtpConfigForm.port"
                :min="1"
                :max="65535"
                :prefix-icon="Operation"
                controls-position="right"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12">
            <el-form-item label="用户名 / 邮箱" prop="username">
              <el-input
                v-model="smtpConfigForm.username"
                placeholder="your-email@example.com"
                clearable
                :prefix-icon="User"
              />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12">
            <el-form-item label="密码 / 授权码" prop="password">
              <el-input
                v-model="smtpConfigForm.password"
                :type="showPassword ? 'text' : 'password'"
                clearable
                :prefix-icon="Lock"
                @input="onPasswordChange"
              >
                <template #suffix>
                  <el-icon
                    class="cursor-pointer"
                    @click="showPassword = !showPassword"
                  >
                    <View v-if="!showPassword" />
                    <Hide v-else />
                  </el-icon>
                </template>
              </el-input>
              <div v-if="passwordSet && !passwordChanged" class="form-item-hint">
                密码已设置（不可见），输入新密码以更新
              </div>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="16">
          <el-col :xs="24" :sm="12">
            <el-form-item label="发件人邮箱" prop="from">
              <el-input
                v-model="smtpConfigForm.from"
                placeholder="noreply@example.com"
                clearable
                :prefix-icon="Message"
              />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12">
            <el-form-item label="发件人名称" prop="fromName">
              <el-input
                v-model="smtpConfigForm.fromName"
                placeholder="校园论坛"
                clearable
                :prefix-icon="UserFilled"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-alert
          v-if="smtpConfigForm.host && smtpConfigForm.port"
          type="info"
          :closable="false"
          class="mt-4"
          show-icon
        >
          <template v-if="passwordSet && !passwordChanged">
            密码已设置但不可见。如需测试或更新密码，请先输入新密码。
          </template>
          <template v-else>
            测试配置前请确保所有参数填写正确，测试邮件将发送到您的发件人邮箱。
          </template>
        </el-alert>
      </el-form>

      <div class="form-actions">
        <el-button
          type="warning"
          plain
          round
          @click="testSmtpConfig"
          :loading="testing"
        >
          <el-icon class="mr-1"><Memo /></el-icon>
          测试配置
        </el-button>
        <el-button
          type="primary"
          round
          @click="saveSmtpConfig"
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
import { ref, computed, onMounted } from 'vue'
import { Cpu, Operation, User, Lock, UserFilled, Message, View, Hide, Memo, Check } from '@element-plus/icons-vue'
import { adminSiteConfigApi } from '../../api/admin'
import { success, error } from '@/utils/message'

// 密码占位符，表示密码已设置但不可见
const PASSWORD_PLACEHOLDER = '••••••••••••••••'

/** SMTP 配置表单数据 */
const smtpConfigForm = ref({
  host: '',
  port: 587,
  username: '',
  password: '',
  from: '',
  fromName: ''
})

const smtpForm = ref(null)
const showPassword = ref(false) // 是否明文显示密码
const saving = ref(false) // 是否正在保存
const testing = ref(false) // 是否正在测试
const passwordSet = ref(false) // 后端是否已设置密码
const passwordChanged = ref(false) // 用户是否修改了密码字段

/** 表单校验规则 */
const rules = {
  host: [{ required: true, message: '此字段为必填项', trigger: 'blur' }],
  port: [
    { required: true, message: '此字段为必填项', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value === null || value === undefined || value === '') {
          callback('此字段为必填项')
        } else if (isNaN(parseFloat(value)) || !isFinite(value)) {
          callback('请输入有效的数字')
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  username: [
    { required: true, message: '此字段为必填项', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value === null || value === undefined || value === '') {
          callback('此字段为必填项')
        } else if (!/.+@.+\..+/.test(value)) {
          callback('请输入有效的邮箱地址')
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  password: [
    {
      validator: (rule, value, callback) => {
        if (passwordSet.value && !passwordChanged.value) {
          callback()
        } else if (!value || value === '') {
          callback('此字段为必填项')
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  from: [
    { required: true, message: '此字段为必填项', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value === null || value === undefined || value === '') {
          callback('此字段为必填项')
        } else if (!/.+@.+\..+/.test(value)) {
          callback('请输入有效的邮箱地址')
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  fromName: [{ required: true, message: '此字段为必填项', trigger: 'blur' }]
}

/** 表单整体校验结果，用于控制保存按钮 */
const formValid = computed(() => {
  const f = smtpConfigForm.value
  const hostOk = !!f.host && f.host.trim() !== ''
  const portOk = f.port !== null && f.port !== undefined && f.port !== ''
  const usernameOk = !!f.username && f.username.trim() !== '' && /.+@.+\..+/.test(f.username)
  const passwordOk = (passwordSet.value && !passwordChanged.value) || (!!f.password && f.password !== '')
  const fromOk = !!f.from && f.from.trim() !== '' && /.+@.+\..+/.test(f.from)
  const fromNameOk = !!f.fromName && f.fromName.trim() !== ''
  return hostOk && portOk && usernameOk && passwordOk && fromOk && fromNameOk
})

/** 加载当前 SMTP 配置 */
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

/** 用户修改密码字段时标记密码已更改 */
const onPasswordChange = () => {
  passwordChanged.value = true
}

/** 保存 SMTP 配置 */
const saveSmtpConfig = async () => {
  if (smtpForm.value) {
    const valid = await smtpForm.value.validate().catch(() => false)
    if (!valid) return
  }
  saving.value = true
  try {
    const updateData = {
      smtp_host: smtpConfigForm.value.host,
      smtp_port: smtpConfigForm.value.port,
      smtp_username: smtpConfigForm.value.username,
      smtp_from: smtpConfigForm.value.from,
      smtp_from_name: smtpConfigForm.value.fromName,
      // 密码字段为占位符或未修改时传空字符串，表示不更新密码
      smtp_password: passwordChanged.value && smtpConfigForm.value.password !== PASSWORD_PLACEHOLDER
        ? smtpConfigForm.value.password
        : ''
    }

    await adminSiteConfigApi.updateConfig(updateData)
    success('SMTP配置保存成功')
    // 重新加载配置以刷新密码设置状态
    loadSmtpConfig()
  } catch (err) {
    console.error('保存SMTP配置失败', err)
    error(err.response?.data?.error || '保存失败')
  } finally {
    saving.value = false
  }
}

/** 测试 SMTP 配置 */
const testSmtpConfig = async () => {
  testing.value = true
  try {
    // 密码未修改且为占位符时，后端不持有明文，无法测试
    if (!passwordChanged.value && smtpConfigForm.value.password === PASSWORD_PLACEHOLDER) {
      error('密码已设置但不可见，无法测试。请先输入新密码后再测试。')
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
.smtp-config-page {
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

.form-item-hint {
  font-size: 12px;
  line-height: 1.5;
  color: var(--campus-text-muted);
  margin-top: 4px;
}

.mt-4 {
  margin-top: 16px;
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

.cursor-pointer {
  cursor: pointer;
}

.mr-1 {
  margin-right: 4px;
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
