<template>
  <v-container fluid class="pa-4 pa-md-6">
    <!-- SMTP配置表单 -->
    <v-card variant="flat" rounded="lg">
      <v-card-title class="pb-2">
        <v-icon start size="20">mdi-email-settings</v-icon>
        SMTP 配置
      </v-card-title>

      <v-card-text>
        <v-form ref="smtpForm" v-model="formValid">
          <v-row dense>
            <v-col cols="12" sm="6">
              <v-text-field
                v-model="smtpConfigForm.host"
                label="SMTP 主机"
                placeholder="smtp.example.com"
                variant="outlined"
                density="compact"
                :rules="[rules.required]"
                prepend-inner-icon="mdi-server"
                clearable
              />
            </v-col>
            <v-col cols="12" sm="6">
              <v-text-field
                v-model.number="smtpConfigForm.port"
                label="SMTP 端口"
                placeholder="587"
                variant="outlined"
                density="compact"
                :rules="[rules.required, rules.number]"
                prepend-inner-icon="mdi-numeric"
                type="number"
                clearable
              />
            </v-col>
            <v-col cols="12" sm="6">
              <v-text-field
                v-model="smtpConfigForm.username"
                label="用户名 / 邮箱"
                placeholder="your-email@example.com"
                variant="outlined"
                density="compact"
                :rules="[rules.required, rules.email]"
                prepend-inner-icon="mdi-account"
                clearable
              />
            </v-col>
            <v-col cols="12" sm="6">
              <v-text-field
                v-model="smtpConfigForm.password"
                label="密码 / 授权码"
                variant="outlined"
                density="compact"
                :rules="passwordSet && !passwordChanged ? [] : [rules.required]"
                :type="showPassword ? 'text' : 'password'"
                prepend-inner-icon="mdi-lock"
                :append-inner-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
                @click:append-inner="showPassword = !showPassword"
                @update:model-value="onPasswordChange"
                :hint="passwordSet && !passwordChanged ? '密码已设置（不可见），输入新密码以更新' : ''"
                persistent-hint
              />
            </v-col>
            <v-col cols="12" sm="6">
              <v-text-field
                v-model="smtpConfigForm.from"
                label="发件人邮箱"
                placeholder="noreply@example.com"
                variant="outlined"
                density="compact"
                :rules="[rules.required, rules.email]"
                prepend-inner-icon="mdi-email"
                clearable
              />
            </v-col>
            <v-col cols="12" sm="6">
              <v-text-field
                v-model="smtpConfigForm.fromName"
                label="发件人名称"
                placeholder="校园论坛"
                variant="outlined"
                density="compact"
                :rules="[rules.required]"
                prepend-inner-icon="mdi-account-circle"
                clearable
              />
            </v-col>
          </v-row>

          <v-alert
            v-if="smtpConfigForm.host && smtpConfigForm.port"
            type="info"
            variant="tonal"
            density="compact"
            class="mt-4"
            icon="mdi-information"
          >
            <template v-if="passwordSet && !passwordChanged">
              密码已设置但不可见。如需测试或更新密码，请先输入新密码。
            </template>
            <template v-else>
              测试配置前请确保所有参数填写正确，测试邮件将发送到您的发件人邮箱。
            </template>
          </v-alert>
        </v-form>
      </v-card-text>

      <v-card-actions class="pa-4">
        <v-spacer />
        <v-btn
          color="secondary"
          variant="tonal"
          @click="testSmtpConfig"
          :loading="testing"
          class="mr-2"
        >
          <v-icon start>mdi-test-tube</v-icon>
          测试配置
        </v-btn>
        <v-btn
          color="primary"
          variant="flat"
          @click="saveSmtpConfig"
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

// 密码占位符，表示密码已设置但不可见
const PASSWORD_PLACEHOLDER = '••••••••••••••••'

const smtpConfigForm = ref({
  host: '',
  port: 587,
  username: '',
  password: '',
  from: '',
  fromName: ''
})

const smtpForm = ref(null)
const formValid = ref(false)
const showPassword = ref(false)
const saving = ref(false)
const testing = ref(false)
const passwordSet = ref(false) // 标记密码是否已设置
const passwordChanged = ref(false) // 标记密码是否被用户修改

const rules = {
  required: v => !!v || '此字段为必填项',
  email: v => /.+@.+\..+/.test(v) || '请输入有效的邮箱地址',
  number: v => !isNaN(parseFloat(v)) && isFinite(v) || '请输入有效的数字'
}

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
  // 当用户修改密码字段时，标记密码已更改
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
    
    // 只有当密码被修改时才发送新密码
    // 如果密码是占位符或未修改，发送空字符串表示不更新密码
    if (passwordChanged.value && smtpConfigForm.value.password !== PASSWORD_PLACEHOLDER) {
      updateData.smtp_password = smtpConfigForm.value.password
    } else {
      updateData.smtp_password = '' // 空字符串表示不修改密码
    }
    
    await adminSiteConfigApi.updateConfig(updateData)
    success('SMTP配置保存成功')
    // 重新加载配置以更新状态
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
    // 测试时需要使用实际密码，如果密码未修改则无法测试
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