<template>
  <v-container fluid class="pa-6">
    <v-card>
      <v-card-title>
        <v-icon class="mr-2">mdi-email-settings</v-icon>
        SMTP 配置
      </v-card-title>
      <v-card-text>
        <v-text-field
          v-model="smtpConfigForm.host"
          label="SMTP 主机"
          variant="outlined"
          class="mb-4"
        ></v-text-field>
        <v-text-field
          v-model.number="smtpConfigForm.port"
          label="SMTP 端口"
          variant="outlined"
          class="mb-4"
        ></v-text-field>
        <v-text-field
          v-model="smtpConfigForm.username"
          label="用户名"
          variant="outlined"
          class="mb-4"
        ></v-text-field>
        <v-text-field
          v-model="smtpConfigForm.password"
          label="密码"
          type="password"
          variant="outlined"
          class="mb-4"
        ></v-text-field>
        <v-text-field
          v-model="smtpConfigForm.from"
          label="发件人邮箱"
          variant="outlined"
          class="mb-4"
        ></v-text-field>
        <v-text-field
          v-model="smtpConfigForm.fromName"
          label="发件人名称"
          variant="outlined"
          class="mb-4"
        ></v-text-field>
      </v-card-text>
      <v-card-actions>
        <v-btn color="primary" variant="flat" @click="testSmtpConfig">
          <v-icon class="mr-1">mdi-test-tube</v-icon>
          测试配置
        </v-btn>
        <v-btn color="primary" variant="flat" @click="saveSmtpConfig">
          保存配置
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-container>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'
import { showSuccess, showError } from '../utils/modal'

const smtpConfigForm = ref({
  host: '',
  port: 587,
  username: '',
  password: '',
  from: '',
  fromName: ''
})

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
  } catch (error) {
    console.error('加载SMTP配置失败', error)
  }
}

const saveSmtpConfig = async () => {
  try {
    await api.put('/site-config', {
      smtp_host: smtpConfigForm.value.host,
      smtp_port: smtpConfigForm.value.port,
      smtp_username: smtpConfigForm.value.username,
      smtp_password: smtpConfigForm.value.password,
      smtp_from: smtpConfigForm.value.from,
      smtp_from_name: smtpConfigForm.value.fromName
    })
    showSuccess('保存成功')
  } catch (error) {
    console.error('保存SMTP配置失败', error)
    showError(error.response?.data?.error || '保存失败')
  }
}

const testSmtpConfig = async () => {
  try {
    await api.post('/site-config/test-smtp', {
      smtp_host: smtpConfigForm.value.host,
      smtp_port: smtpConfigForm.value.port,
      smtp_username: smtpConfigForm.value.username,
      smtp_password: smtpConfigForm.value.password,
      smtp_from: smtpConfigForm.value.from,
      smtp_to: smtpConfigForm.value.from
    })
    showSuccess('测试邮件发送成功')
  } catch (error) {
    console.error('测试SMTP配置失败', error)
    showError(error.response?.data?.error || '测试失败')
  }
}

onMounted(() => {
  loadSmtpConfig()
})
</script>