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
            mdi-key
          </v-icon>
        </div>
        <h1 class="text-2xl font-bold text-gradient">找回密码</h1>
        <p class="text-gray-500 mt-2 text-sm">通过QQ邮箱验证重置密码</p>
      </v-card-title>
      
      <v-card-text class="px-8 pb-8">
        <v-alert v-if="error" type="error" variant="tonal" class="mb-4 text-sm rounded-lg"></v-alert>
        <v-alert v-if="success" type="success" variant="tonal" class="mb-4 text-sm rounded-lg"></v-alert>

        <v-form @submit.prevent="handleSubmit">
          <v-text-field
            v-model="form.qq_number"
            label="QQ号码"
            variant="outlined"
            required
            prepend-inner-icon="mdi-qqchat"
            class="mb-4"
            color="primary"
            hide-details="auto"
            density="comfortable"
          ></v-text-field>

          <div v-if="step === 2">
            <v-text-field
              v-model="form.code"
              label="验证码"
              variant="outlined"
              required
              prepend-inner-icon="mdi-shield-check"
              class="mb-4"
              color="primary"
              hide-details="auto"
              density="comfortable"
            ></v-text-field>

            <v-text-field
              v-model="form.password"
              label="新密码"
              variant="outlined"
              type="password"
              required
              prepend-inner-icon="mdi-lock"
              append-icon="mdi-eye-off"
              class="mb-4"
              color="primary"
              hide-details="auto"
              density="comfortable"
            ></v-text-field>

            <v-text-field
              v-model="form.confirm_password"
              label="确认密码"
              variant="outlined"
              type="password"
              required
              prepend-inner-icon="mdi-lock-check"
              append-icon="mdi-eye-off"
              class="mb-6"
              color="primary"
              hide-details="auto"
              density="comfortable"
            ></v-text-field>
          </div>

          <v-btn
            type="submit"
            class="w-full mb-6 btn-gradient"
            size="large"
            block
            :loading="loading"
          >
            {{ loading ? '处理中...' : buttonText }}
          </v-btn>
        </v-form>

        <div class="text-center">
          <span class="text-gray-500 text-sm">想起密码了？</span>
          <router-link to="/login" class="text-primary ml-1 text-sm">返回登录</router-link>
        </div>
      </v-card-text>
    </v-card>
  </div>
</template>

<script>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import api from '../api'

export default {
  name: 'ForgotPassword',
  setup() {
    const router = useRouter()
    const step = ref(1)
    const loading = ref(false)
    const error = ref('')
    const success = ref('')
    const resetIdentifier = ref('')

    const form = ref({
      qq_number: '',
      code: '',
      password: '',
      confirm_password: ''
    })

    const buttonText = computed(() => {
      return step.value === 1 ? '发送验证码' : '重置密码'
    })

    const handleSubmit = async () => {
      error.value = ''
      success.value = ''

      if (step.value === 1) {
        if (!form.value.qq_number.trim()) {
          error.value = '请输入QQ号码'
          return
        }

        loading.value = true
        try {
          const response = await api.post('/password/reset-code', {
            qq_number: form.value.qq_number
          })
          success.value = response.data.message || '验证码已发送到您的QQ邮箱'
          resetIdentifier.value = response.data.identifier || ''
          step.value = 2
        } catch (err) {
          error.value = err.response?.data?.error || err.message || '发送验证码失败'
        } finally {
          loading.value = false
        }
      } else {
        if (!form.value.code.trim()) {
          error.value = '请输入验证码'
          return
        }

        if (form.value.password.length < 6) {
          error.value = '密码长度不能少于6位'
          return
        }

        if (form.value.password !== form.value.confirm_password) {
          error.value = '两次输入的密码不一致'
          return
        }

        loading.value = true
        try {
          await api.post('/password/reset', {
            qq_number: form.value.qq_number,
            code: form.value.code,
            identifier: resetIdentifier.value,
            password: form.value.password
          })
          success.value = '密码重置成功！'
          setTimeout(() => {
            router.push('/login')
          }, 1500)
        } catch (err) {
          error.value = err.response?.data?.error || err.message || '重置密码失败'
        } finally {
          loading.value = false
        }
      }
    }

    return {
      form,
      step,
      loading,
      error,
      success,
      buttonText,
      handleSubmit
    }
  }
}
</script>