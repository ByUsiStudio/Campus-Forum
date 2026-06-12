<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { authApi } from '../api'

const router = useRouter()

const form = ref({
  username: '',
  email: '',
  password: '',
  confirmPassword: ''
})

const isLoading = ref(false)
const error = ref('')

const handleRegister = async () => {
  if (!form.value.username || !form.value.email || !form.value.password) {
    error.value = '请填写所有必填项'
    return
  }
  if (form.value.password !== form.value.confirmPassword) {
    error.value = '两次输入的密码不一致'
    return
  }
  isLoading.value = true
  try {
    await authApi.register(form.value)
    router.push('/login')
  } catch (err) {
    error.value = err.response?.data?.error || '注册失败'
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <v-container fluid class="min-h-screen d-flex align-center justify-center">
    <v-card width="400" max-width="90%" elevation="8">
      <v-card-title class="text-center">
        <v-icon size="48" color="primary">mdi-user-plus</v-icon>
        <h2 class="text-h5 font-weight-bold mt-2">注册账号</h2>
      </v-card-title>
      
      <v-card-text>
        <v-alert v-if="error" type="error" dense>
          {{ error }}
        </v-alert>
        
        <v-form @submit.prevent="handleRegister">
          <v-text-field
            v-model="form.username"
            label="用户名"
            prepend-icon="mdi-account"
            class="mb-4"
            required
          />
          
          <v-text-field
            v-model="form.email"
            label="邮箱"
            type="email"
            prepend-icon="mdi-email"
            class="mb-4"
            required
          />
          
          <v-text-field
            v-model="form.password"
            label="密码"
            type="password"
            prepend-icon="mdi-lock"
            class="mb-4"
            required
          />
          
          <v-text-field
            v-model="form.confirmPassword"
            label="确认密码"
            type="password"
            prepend-icon="mdi-lock-check"
            class="mb-6"
            required
          />
          
          <v-btn
            type="submit"
            color="primary"
            block
            :loading="isLoading"
          >
            注册
          </v-btn>
        </v-form>
        
        <v-btn text color="primary" block class="mt-4" @click="router.push('/login')">
          已有账号？立即登录
        </v-btn>
      </v-card-text>
    </v-card>
  </v-container>
</template>
