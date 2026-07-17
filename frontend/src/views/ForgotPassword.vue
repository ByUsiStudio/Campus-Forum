<template>
  <div class="d-flex justify-center align-center" style="min-height: 80vh;">
    <div class="layui-card" style="width: 100%; max-width: 400px; padding: 24px;">
      <div class="text-center mb-4">
        <h2 style="color: var(--primary); font-size: 20px; font-weight: 600;">找回密码</h2>
      </div>

      <div class="layui-card-body">
        <div v-if="error" class="layui-alert layui-alert-danger mb-4">{{ error }}</div>
        <div v-if="success" class="layui-alert layui-alert-success mb-4">{{ success }}</div>

        <form @submit.prevent="handleSubmit" class="layui-form">
          <div class="layui-form-item">
            <label class="layui-form-label">QQ号码</label>
            <div class="layui-input-group">
              <div class="layui-input-group-prepend">
                <span class="layui-input-group-text"><i class="fa-solid fa-message"></i></span>
              </div>
              <input type="text" v-model="form.qq_number" class="layui-input" placeholder="请输入QQ号码" required />
            </div>
          </div>

          <div v-if="step === 2">
            <div class="layui-form-item">
              <label class="layui-form-label">验证码</label>
              <div class="layui-input-group">
                <div class="layui-input-group-prepend">
                  <span class="layui-input-group-text"><i class="fa-solid fa-shield-check"></i></span>
                </div>
                <input type="text" v-model="form.code" class="layui-input" placeholder="请输入验证码" required />
              </div>
            </div>

            <div class="layui-form-item">
              <label class="layui-form-label">新密码</label>
              <div class="layui-input-group">
                <div class="layui-input-group-prepend">
                  <span class="layui-input-group-text"><i class="fa-solid fa-lock"></i></span>
                </div>
                <input type="password" v-model="form.password" class="layui-input" placeholder="请输入新密码" required />
              </div>
            </div>

            <div class="layui-form-item">
              <label class="layui-form-label">确认密码</label>
              <div class="layui-input-group">
                <div class="layui-input-group-prepend">
                  <span class="layui-input-group-text"><i class="fa-solid fa-lock-check"></i></span>
                </div>
                <input type="password" v-model="form.confirm_password" class="layui-input" placeholder="请再次输入密码" required />
              </div>
            </div>
          </div>

          <button
            type="submit"
            class="layui-btn layui-btn-lg layui-btn-fluid"
            :disabled="loading"
          >
            {{ loading ? '处理中...' : buttonText }}
          </button>
        </form>

        <div class="text-center mt-4 text-secondary" style="font-size: 14px;">
          想起密码了？ <router-link to="/login" style="color: var(--primary);">返回登录</router-link>
        </div>
      </div>
    </div>
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