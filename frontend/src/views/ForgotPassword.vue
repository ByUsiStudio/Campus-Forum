<template>
  <div class="d-flex justify-center align-center" style="min-height: 80vh;">
    <v-card width="100%" max-width="400" class="pa-6">
      <v-card-title class="text-h5 text-center pb-4" style="color: rgb(var(--v-theme-primary));">
        找回密码
      </v-card-title>

      <v-card-text>
        <v-alert v-if="error" type="error" variant="tonal" class="mb-4">{{ error }}</v-alert>
        <v-alert v-if="success" type="success" variant="tonal" class="mb-4">{{ success }}</v-alert>

        <v-form @submit.prevent="handleSubmit">
          <v-text-field
            v-model="form.qq_number"
            label="QQ号码"
            variant="outlined"
            required
            prepend-inner-icon="mdi-qqchat"
            class="mb-4"
          ></v-text-field>

          <div v-if="step === 2">
            <v-text-field
              v-model="form.code"
              label="验证码"
              variant="outlined"
              required
              prepend-inner-icon="mdi-shield-check"
              class="mb-4"
            ></v-text-field>

            <v-text-field
              v-model="form.password"
              label="新密码"
              variant="outlined"
              type="password"
              required
              prepend-inner-icon="mdi-lock"
              class="mb-4"
            ></v-text-field>

            <v-text-field
              v-model="form.confirm_password"
              label="确认密码"
              variant="outlined"
              type="password"
              required
              prepend-inner-icon="mdi-lock-check"
              class="mb-4"
            ></v-text-field>
          </div>

          <v-btn
            type="submit"
            color="primary"
            block
            size="large"
            :loading="loading"
          >
            {{ loading ? '处理中...' : buttonText }}
          </v-btn>
        </v-form>

        <div class="text-center mt-4 text-body-2">
          想起密码了？ <router-link to="/login" class="text-primary">返回登录</router-link>
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
          step.value = 2
        } catch (err) {
          error.value = err.message || '发送验证码失败'
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
            password: form.value.password
          })
          success.value = '密码重置成功！'
          setTimeout(() => {
            router.push('/login')
          }, 1500)
        } catch (err) {
          error.value = err.message || '重置密码失败'
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
