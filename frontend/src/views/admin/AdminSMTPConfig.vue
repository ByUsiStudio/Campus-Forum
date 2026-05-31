<template>
  <v-container fluid class="pa-6">
    <v-card class="pa-6">
      <v-card-title class="d-flex align-center mb-6">
        <v-avatar color="primary" size="44" class="mr-3">
          <v-icon color="white">mdi-email-settings</v-icon>
        </v-avatar>
        <div>
          <div class="text-h5 font-weight-bold">SMTP 配置</div>
          <div class="text-caption text-medium-emphasis">配置邮件发送服务参数</div>
        </div>
      </v-card-title>

      <v-divider class="mb-6"></v-divider>

      <v-form ref="smtpForm" v-model="formValid">
        <v-row>
          <v-col cols="12" md="6">
            <v-text-field
              v-model="smtpConfigForm.host"
              label="SMTP 主机"
              placeholder="smtp.example.com"
              variant="outlined"
              density="comfortable"
              :rules="[rules.required]"
              prepend-inner-icon="mdi-server"
              clearable
              class="mb-2"
            >
              <template #label>
                <span class="text-body-2">SMTP 主机</span>
              </template>
            </v-text-field>
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field
              v-model.number="smtpConfigForm.port"
              label="SMTP 端口"
              placeholder="587"
              variant="outlined"
              density="comfortable"
              :rules="[rules.required, rules.number]"
              prepend-inner-icon="mdi-numeric"
              type="number"
              clearable
              class="mb-2"
            >
              <template #label>
                <span class="text-body-2">SMTP 端口</span>
              </template>
            </v-text-field>
          </v-col>
        </v-row>

        <v-row>
          <v-col cols="12" md="6">
            <v-text-field
              v-model="smtpConfigForm.username"
              label="用户名"
              placeholder="your-email@example.com"
              variant="outlined"
              density="comfortable"
              :rules="[rules.required, rules.email]"
              prepend-inner-icon="mdi-account"
              clearable
              class="mb-2"
            >
              <template #label>
                <span class="text-body-2">用户名 / 邮箱</span>
              </template>
            </v-text-field>
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field
              v-model="smtpConfigForm.password"
              label="密码"
              variant="outlined"
              density="comfortable"
              :rules="[rules.required]"
              :type="showPassword ? 'text' : 'password'"
              prepend-inner-icon="mdi-lock"
              :append-inner-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
              @click:append-inner="showPassword = !showPassword"
              class="mb-2"
            >
              <template #label>
                <span class="text-body-2">密码 / 授权码</span>
              </template>
            </v-text-field>
          </v-col>
        </v-row>

        <v-row>
          <v-col cols="12" md="6">
            <v-text-field
              v-model="smtpConfigForm.from"
              label="发件人邮箱"
              placeholder="noreply@example.com"
              variant="outlined"
              density="comfortable"
              :rules="[rules.required, rules.email]"
              prepend-inner-icon="mdi-email"
              clearable
              class="mb-2"
            >
              <template #label>
                <span class="text-body-2">发件人邮箱</span>
              </template>
            </v-text-field>
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field
              v-model="smtpConfigForm.fromName"
              label="发件人名称"
              placeholder="校园论坛"
              variant="outlined"
              density="comfortable"
              :rules="[rules.required]"
              prepend-inner-icon="mdi-account-circle"
              clearable
              class="mb-2"
            >
              <template #label>
                <span class="text-body-2">发件人名称</span>
              </template>
            </v-text-field>
          </v-col>
        </v-row>

        <v-alert
          v-if="smtpConfigForm.host && smtpConfigForm.port"
          type="info"
          variant="tonal"
          density="compact"
          class="mb-4 mt-2"
          icon="mdi-information"
        >
          <span class="text-body-2">
            测试配置前请确保所有参数填写正确，测试邮件将发送到您的发件人邮箱。
          </span>
        </v-alert>
      </v-form>

      <v-divider class="my-6"></v-divider>

      <v-card-actions class="px-0">
        <v-spacer></v-spacer>
        <v-btn
          color="secondary"
          variant="tonal"
          size="large"
          @click="testSmtpConfig"
          :loading="testing"
          class="mr-4"
        >
          <v-icon class="mr-2">mdi-test-tube</v-icon>
          测试配置
        </v-btn>
        <v-btn
          color="primary"
          variant="flat"
          size="large"
          @click="saveSmtpConfig"
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
    const response = await api.get('/site-config')
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
  }
}

const saveSmtpConfig = async () => {
  saving.value = true
  try {
    await api.put('/site-config', {
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
    await api.post('/site-config/test-smtp', {
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
:deep(.v-field) {
  border-radius: 12px;
}

:deep(.v-field--outlined .v-field__outline) {
  border-color: rgba(148, 163, 184, 0.3);
}

:deep(.v-field--focused .v-field__outline) {
  border-color: rgb(var(--v-theme-primary));
}
</style>