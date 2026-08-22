<template>
  <div class="auth-page-center">
    <el-card class="auth-card card-surface">
      <template #header>
        <div class="auth-header">
          <el-icon class="auth-logo"><Key /></el-icon>
          <span class="auth-title">找回密码</span>
        </div>
      </template>

      <el-steps :active="step" align-center class="mb-6">
        <el-step title="发送验证码" />
        <el-step title="重置密码" />
      </el-steps>

      <el-alert
        v-if="error"
        :title="error"
        type="error"
        show-icon
        :closable="false"
        class="mb-4"
      />
      <el-alert
        v-if="success"
        :title="success"
        type="success"
        show-icon
        :closable="false"
        class="mb-4"
      />

      <el-form
        :model="form"
        :rules="rules"
        ref="formRef"
        label-position="top"
        @submit.prevent="handleSubmit"
      >
        <el-form-item label="QQ号码" prop="qq_number">
          <el-input
            v-model="form.qq_number"
            placeholder="请输入QQ号码"
            :prefix-icon="ChatDotRound"
            :disabled="step === 2"
          />
        </el-form-item>

        <div v-if="step === 2">
          <el-form-item label="验证码" prop="code">
            <el-input
              v-model="form.code"
              placeholder="请输入验证码"
              :prefix-icon="CircleCheck"
            />
          </el-form-item>

          <el-form-item label="新密码" prop="password">
            <el-input
              v-model="form.password"
              type="password"
              placeholder="请输入新密码"
              show-password
              :prefix-icon="Lock"
            />
          </el-form-item>

          <el-form-item label="确认密码" prop="confirm_password">
            <el-input
              v-model="form.confirm_password"
              type="password"
              placeholder="请再次输入新密码"
              show-password
              :prefix-icon="Lock"
            />
          </el-form-item>
        </div>

        <el-button
          type="primary"
          size="large"
          class="w-full"
          :loading="loading"
          native-type="submit"
        >
          {{ loading ? '处理中...' : buttonText }}
        </el-button>
      </el-form>

      <div class="auth-footer text-secondary">
        想起密码了？
        <router-link to="/login" class="auth-link">返回登录</router-link>
      </div>
    </el-card>
  </div>
</template>

<script>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { Key, ChatDotRound, CircleCheck, Lock } from '@element-plus/icons-vue'
import { authApi } from '@/api'
import { success, error } from '@/utils/message'

export default {
  name: 'ForgotPassword',
  setup() {
    const router = useRouter()
    const formRef = ref(null)
    const step = ref(1)
    const loading = ref(false)
    const errMsg = ref('')
    const successMsg = ref('')
    const resetIdentifier = ref('')

    const form = ref({
      qq_number: '',
      code: '',
      password: '',
      confirm_password: ''
    })

    const rules = {
      qq_number: [{ required: true, message: '请输入QQ号码', trigger: 'blur' }],
      code: [{ required: true, message: '请输入验证码', trigger: 'blur' }],
      password: [
        { required: true, message: '请输入新密码', trigger: 'blur' },
        { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
      ],
      confirm_password: [
        { required: true, message: '请再次输入新密码', trigger: 'blur' },
        {
          validator: (rule, value, callback) => {
            if (value !== form.value.password) {
              callback(new Error('两次输入的密码不一致'))
            } else {
              callback()
            }
          },
          trigger: 'blur'
        }
      ]
    }

    const buttonText = computed(() => {
      return step.value === 1 ? '发送验证码' : '重置密码'
    })

    const handleSubmit = async () => {
      errMsg.value = ''
      successMsg.value = ''

      try {
        await formRef.value?.validate()
      } catch (e) {
        return
      }

      if (step.value === 1) {
        loading.value = true
        try {
          const response = await authApi.sendResetCode({
            qq_number: form.value.qq_number
          })
          const msg = response.data.message || '验证码已发送到您的QQ邮箱'
          successMsg.value = msg
          success(msg)
          resetIdentifier.value = response.data.identifier || ''
          step.value = 2
        } catch (err) {
          errMsg.value = err.response?.data?.error || err.message || '发送验证码失败'
          error(errMsg.value)
        } finally {
          loading.value = false
        }
      } else {
        loading.value = true
        try {
          await authApi.resetPassword({
            qq_number: form.value.qq_number,
            code: form.value.code,
            identifier: resetIdentifier.value,
            password: form.value.password
          })
          successMsg.value = '密码重置成功！'
          success('密码重置成功！')
          setTimeout(() => {
            router.push('/login')
          }, 1500)
        } catch (err) {
          errMsg.value = err.response?.data?.error || err.message || '重置密码失败'
          error(errMsg.value)
        } finally {
          loading.value = false
        }
      }
    }

    return {
      form,
      formRef,
      rules,
      step,
      loading,
      error: errMsg,
      success: successMsg,
      buttonText,
      Key,
      ChatDotRound,
      CircleCheck,
      Lock,
      handleSubmit
    }
  }
}
</script>

<style scoped>
.auth-page-center {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 80vh;
}

.auth-card {
  width: 100%;
  max-width: 400px;
  padding: 8px 8px 16px;
}

.auth-header {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 12px 0 4px;
}

.auth-logo {
  font-size: 26px;
  color: var(--campus-primary);
}

.auth-title {
  font-size: 20px;
  font-weight: 600;
  color: var(--campus-primary);
}

.auth-footer {
  text-align: center;
  margin-top: 16px;
  font-size: 13px;
}

.auth-link {
  color: var(--campus-primary);
}

.mb-4 {
  margin-bottom: 16px;
}

.mb-6 {
  margin-bottom: 24px;
}

.w-full {
  width: 100%;
}
</style>
