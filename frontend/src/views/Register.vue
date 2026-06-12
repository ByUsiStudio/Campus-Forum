<script setup>
import { ref, inject } from 'vue'
import { useRouter } from 'vue-router'
import { authApi } from '../api'

const router = useRouter()
const setUser = inject('setUser')

const form = ref({
  username: '',
  qq_number: '',
  display_name: '',
  password: '',
  confirm_password: ''
})

const isLoading = ref(false)
const error = ref('')

const handleRegister = async () => {
  if (!form.value.username || !form.value.password || !form.value.qq_number) {
    error.value = '请填写必填项'
    return
  }
  
  if (form.value.password !== form.value.confirm_password) {
    error.value = '两次输入的密码不一致'
    return
  }

  isLoading.value = true
  try {
    await authApi.register(form.value)
    const loginResponse = await authApi.login({
      username: form.value.username,
      password: form.value.password
    })
    localStorage.setItem('token', loginResponse.data.token)
    setUser(loginResponse.data.user)
    router.push('/')
  } catch (err) {
    error.value = err.response?.data?.error || '注册失败'
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
            mdi-user-plus
          </v-icon>
        </div>
        <h1 class="text-2xl font-bold text-gradient">注册账号</h1>
        <p class="text-gray-500 mt-2 text-sm">加入我们的校园社区</p>
      </v-card-title>
      
      <v-card-text class="px-8 pb-8">
        <v-form @submit.prevent="handleRegister">
          <v-text-field
            v-model="form.username"
            label="用户名"
            placeholder="请输入用户名"
            prepend-icon="mdi-account"
            class="mb-4"
            color="primary"
            hide-details="auto"
            density="comfortable"
          />
          
          <v-text-field
            v-model="form.qq_number"
            label="QQ号"
            placeholder="请输入QQ号"
            prepend-icon="mdi-message-circle"
            class="mb-4"
            color="primary"
            hide-details="auto"
            density="comfortable"
          />
          
          <v-text-field
            v-model="form.display_name"
            label="昵称（可选）"
            placeholder="请输入昵称"
            prepend-icon="mdi-face"
            class="mb-4"
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
            class="mb-4"
            color="primary"
            hide-details="auto"
            density="comfortable"
          />
          
          <v-text-field
            v-model="form.confirm_password"
            label="确认密码"
            type="password"
            placeholder="请再次输入密码"
            prepend-icon="mdi-lock-check"
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
            <span v-if="!isLoading">注册</span>
            <span v-else>注册中...</span>
          </v-btn>
        </v-form>
        
        <div class="text-center">
          <span class="text-gray-500 text-sm">已有账号？</span>
          <v-btn 
            text 
            color="primary" 
            size="small"
            @click="router.push('/login')"
          >
            立即登录
          </v-btn>
        </div>
      </v-card-text>
    </v-card>
  </div>
</template>