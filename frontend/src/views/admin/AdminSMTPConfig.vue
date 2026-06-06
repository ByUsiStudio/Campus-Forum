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
                :rules="[rules.required]"
                :type="showPassword ? 'text' : 'password'"
                prepend-inner-icon="mdi-lock"
                :append-inner-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
                @click:append-inner="showPassword = !showPassword"
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
            测试配置前请确保所有参数填写正确，测试邮件将发送到您的发件人邮箱。
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

const rules = {
  required: v => !!v || '此字段为必填项',
  email: v => /.+@.+\..+/.test(v) || '请输入有效的邮箱地址',
  number: v => !isNaN(parseFloat(v)) && isFinite(v) || '请输入有效的数字'
}

const loadSmtpConfig = async () => {
  try {
    const response = await adminSiteConfigApi.getConfig()
    smtpConfigForm.value = {
      host: response.data.smtp_host || '',
      port: response.data.smtp_port || 587,
      username: response.data.smtp_username || '',
      password: response.data.smtp_password || '',
      from: response.data.smtp_from || '',
      fromName: response.data.smtp_from_name || ''
    }
  } catch (err) {
    console.error('加载SMTP配置失败', err)
    error('加载SMTP配置失败')
  }
}

const saveSmtpConfig = async () => {
  saving.value = true
  try {
    await adminSiteConfigApi.updateConfig({
      smtp_host: smtpConfigForm.value.host,
      smtp_port: smtpConfigForm.value.port,
      smtp_username: smtpConfigForm.value.username,
      smtp_password: smtpConfigForm.value.password,
      smtp_from: smtpConfigForm.value.from,
      smtp_from_name: smtpConfigForm.value.fromName
    })
    success('SMTP配置保存成功')
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