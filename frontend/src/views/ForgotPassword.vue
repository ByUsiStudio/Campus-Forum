<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const form = ref({
  email: '',
  code: '',
  password: ''
})

const step = ref(1)
const isLoading = ref(false)
const error = ref('')

const sendCode = async () => {
  if (!form.value.email) {
    error.value = '请输入邮箱'
    return
  }
  isLoading.value = true
  try {
    await new Promise(resolve => setTimeout(resolve, 1000))
    step.value = 2
  } catch (err) {
    error.value = '发送失败'
  } finally {
    isLoading.value = false
  }
}

const resetPassword = async () => {
  if (!form.value.code || !form.value.password) {
    error.value = '请填写完整信息'
    return
  }
  isLoading.value = true
  try {
    await new Promise(resolve => setTimeout(resolve, 1000))
    step.value = 3
  } catch (err) {
    error.value = '重置失败'
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <v-container fluid class="min-h-screen d-flex align-center justify-center">
    <v-card width="400" max-width="90%" elevation="8">
      <!-- 步骤1：输入邮箱 -->
      <template v-if="step === 1">
        <v-card-title class="text-center">
          <v-icon size="48" color="primary">mdi-lock-reset</v-icon>
          <h2 class="text-h5 font-weight-bold mt-2">找回密码</h2>
        </v-card-title>
        
        <v-card-text>
          <v-alert v-if="error" type="error" dense>
            {{ error }}
          </v-alert>
          
          <v-text-field
            v-model="form.email"
            label="邮箱"
            type="email"
            prepend-icon="mdi-email"
            class="mb-6"
          />
          
          <v-btn
            color="primary"
            block
            :loading="isLoading"
            @click="sendCode"
          >
            发送验证码
          </v-btn>
          
          <v-btn text color="primary" block class="mt-4" @click="router.push('/login')">
            返回登录
          </v-btn>
        </v-card-text>
      </template>
      
      <!-- 步骤2：输入验证码和新密码 -->
      <template v-if="step === 2">
        <v-card-title class="text-center">
          <v-icon size="48" color="primary">mdi-lock</v-icon>
          <h2 class="text-h5 font-weight-bold mt-2">设置新密码</h2>
        </v-card-title>
        
        <v-card-text>
          <v-alert v-if="error" type="error" dense>
            {{ error }}
          </v-alert>
          
          <v-text-field
            v-model="form.code"
            label="验证码"
            prepend-icon="mdi-code"
            class="mb-4"
          />
          
          <v-text-field
            v-model="form.password"
            label="新密码"
            type="password"
            prepend-icon="mdi-lock"
            class="mb-6"
          />
          
          <v-btn
            color="primary"
            block
            :loading="isLoading"
            @click="resetPassword"
          >
            确认重置
          </v-btn>
          
          <v-btn text color="primary" block class="mt-4" @click="step = 1">
            返回上一步
          </v-btn>
        </v-card-text>
      </template>
      
      <!-- 步骤3：重置成功 -->
      <template v-if="step === 3">
        <v-card-title class="text-center">
          <v-icon size="48" color="success">mdi-check-circle</v-icon>
          <h2 class="text-h5 font-weight-bold mt-2">重置成功</h2>
        </v-card-title>
        
        <v-card-text class="text-center">
          <p>您的密码已重置成功</p>
          <v-btn color="primary" class="mt-4" @click="router.push('/login')">
            返回登录
          </v-btn>
        </v-card-text>
      </template>
    </v-card>
  </v-container>
</template>
