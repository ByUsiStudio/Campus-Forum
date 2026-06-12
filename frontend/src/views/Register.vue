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
  <div class="min-h-screen gradient-purple-light d-flex align-items-center justify-center p-4">
    <v-card 
      class="w-full max-w-md" 
      elevation="8"
      rounded="xl"
    >
      <v-card-title class="text-center py-6">
        <v-icon 
          size="48" 
          color="primary"
          class="mx-auto mb-3"
        >
          mdi-user-plus
        </v-icon>
        <h2 class="text-2xl font-bold text-primary">注册账号</h2>
        <p class="text-gray-500 mt-2">加入我们的校园社区</p>
      </v-card-title>
      
      <v-card-text>
        <v-form @submit.prevent="handleRegister">
          <v-text-field
            v-model="form.username"
            label="用户名"
            placeholder="请输入用户名"
            prepend-icon="mdi-account"
            class="mb-3"
          />
          
          <v-text-field
            v-model="form.qq_number"
            label="QQ号"
            placeholder="请输入QQ号"
            prepend-icon="mdi-message-circle"
            class="mb-3"
          />
          
          <v-text-field
            v-model="form.display_name"
            label="昵称（可选）"
            placeholder="请输入昵称"
            prepend-icon="mdi-face"
            class="mb-3"
          />
          
          <v-text-field
            v-model="form.password"
            label="密码"
            type="password"
            placeholder="请输入密码"
            prepend-icon="mdi-lock"
            class="mb-3"
          />
          
          <v-text-field
            v-model="form.confirm_password"
            label="确认密码"
            type="password"
            placeholder="请再次输入密码"
            prepend-icon="mdi-lock-check"
            class="mb-6"
          />
          
          <v-btn
            type="submit"
            color="primary"
            class="w-full mb-4"
            :loading="isLoading"
            block
            size="large"
          >
            <span v-if="!isLoading">注册</span>
            <span v-else>注册中...</span>
          </v-btn>
        </v-form>
        
        <div v-if="error" class="text-center mb-4">
          <v-alert 
            type="error" 
            border="bottom" 
            class="text-sm"
          >
            {{ error }}
          </v-alert>
        </div>
        
        <div class="text-center">
          <span class="text-gray-500">已有账号？</span>
          <v-btn 
            text 
            color="primary" 
            @click="router.push('/login')"
          >
            立即登录
          </v-btn>
        </div>
      </v-card-text>
    </v-card>
  </div>
</template>