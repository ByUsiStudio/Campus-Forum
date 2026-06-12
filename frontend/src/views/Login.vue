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
  <div class="min-h-screen flex items-center justify-center p-4 relative overflow-hidden">
    <div class="absolute inset-0 overflow-hidden">
      <div class="absolute -top-40 -right-40 w-80 h-80 bg-primary/10 rounded-full blur-3xl"></div>
      <div class="absolute -bottom-40 -left-40 w-80 h-80 bg-accent/10 rounded-full blur-3xl"></div>
      <div class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-96 h-96 bg-secondary/5 rounded-full blur-3xl"></div>
    </div>
    
    <v-card 
      class="w-full max-w-md relative z-10" 
      elevation="12"
      rounded="2xl"
      background-color="surface"
    >
      <v-card-title class="text-center py-8 px-6">
        <div class="mb-6">
          <v-icon 
            size="56" 
            color="primary"
            class="mx-auto"
            style="filter: drop-shadow(0 4px 12px rgba(147, 112, 219, 0.3));"
          >
            mdi-forum
          </v-icon>
        </div>
        <h1 class="text-2xl font-bold text-gradient">校园论坛</h1>
        <p class="text-gray-500 mt-2 text-sm">欢迎回来，开始你的校园之旅</p>
      </v-card-title>
      
      <v-card-text class="px-8 pb-8">
        <v-form @submit.prevent="handleLogin">
          <v-text-field
            v-model="form.username"
            label="用户名"
            placeholder="请输入用户名"
            prepend-icon="mdi-account"
            class="mb-5"
            color="primary"
            hide-details="auto"
            density="comfortable"
          />
          
          <v-text-field
            v-model="form.password"
            label="密码"
            type="password"
            placeholder="请输入密码"
            prepend-icon="mdi-lock"
            append-icon="mdi-eye-off"
            class="mb-6"
            color="primary"
            hide-details="auto"
            density="comfortable"
          />
          
          <div v-if="error" class="mb-6">
            <v-alert 
              type="error" 
              border="bottom" 
              class="text-sm rounded-lg"
              dense
            >
              {{ error }}
            </v-alert>
          </div>
          
          <v-btn
            type="submit"
            class="w-full mb-6 btn-gradient"
            :loading="isLoading"
            size="large"
            block
          >
            <span v-if="!isLoading">登录</span>
            <span v-else>登录中...</span>
          </v-btn>
        </v-form>
        
        <div class="flex items-center justify-between">
          <v-btn 
            text 
            color="gray-500" 
            size="small"
            @click="router.push('/forgot-password')"
          >
            忘记密码？
          </v-btn>
          <div class="flex items-center">
            <span class="text-gray-500 text-sm">还没有账号？</span>
            <v-btn 
              text 
              color="primary" 
              size="small"
              @click="router.push('/register')"
            >
              立即注册
            </v-btn>
          </div>
        </div>
      </v-card-text>
    </v-card>
  </div>
</template>