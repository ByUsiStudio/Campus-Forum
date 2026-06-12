<script setup>
import { ref, inject } from 'vue'
import { useRouter } from 'vue-router'
import { authApi } from '../api'

const router = useRouter()
const setUser = inject('setUser')

const form = ref({
  username: '',
  password: ''
})

const isLoading = ref(false)
const error = ref('')

const handleLogin = async () => {
  if (!form.value.username || !form.value.password) {
    error.value = '请填写用户名和密码'
    return
  }
  isLoading.value = true
  try {
    const response = await authApi.login(form.value)
    localStorage.setItem('token', response.data.token)
    setUser(response.data.user)
    router.push('/')
  } catch (err) {
    error.value = err.response?.data?.error || '登录失败'
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <v-container fluid class="min-h-screen d-flex align-center justify-center">
    <v-card width="400" max-width="90%" elevation="8">
      <v-card-title class="text-center">
        <v-icon size="48" color="primary">mdi-forum</v-icon>
        <h2 class="text-h5 font-weight-bold mt-2">校园论坛</h2>
      </v-card-title>
      
      <v-card-text>
        <v-alert v-if="error" type="error" dense>
          {{ error }}
        </v-alert>
        
        <v-form @submit.prevent="handleLogin">
          <v-text-field
            v-model="form.username"
            label="用户名"
            prepend-icon="mdi-account"
            class="mb-4"
            required
          />
          
          <v-text-field
            v-model="form.password"
            label="密码"
            type="password"
            prepend-icon="mdi-lock"
            append-icon="mdi-eye-off"
            class="mb-6"
            required
          />
          
          <v-btn
            type="submit"
            color="primary"
            block
            :loading="isLoading"
          >
            登录
          </v-btn>
        </v-form>
        
        <div class="d-flex justify-between mt-4">
          <v-btn text color="primary" @click="router.push('/forgot-password')">
            忘记密码
          </v-btn>
          <v-btn text color="primary" @click="router.push('/register')">
            注册账号
          </v-btn>
        </div>
      </v-card-text>
    </v-card>
  </v-container>
</template>
