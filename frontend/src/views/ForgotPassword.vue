<template>
  <div class="auth-page-center">
    <div class="auth-card card-surface">
      <div class="auth-hero">
        <div class="auth-badge"><el-icon :size="26"><Key /></el-icon></div>
        <h1 class="auth-title">找回密码</h1>
        <p class="auth-subtitle text-center second-line">
          通过QQ邮箱验证重置你的密码
        </p>
      </div>

      <div class="auth-body">
        <div class="step-indicator">
          <div
            class="step-item"
            :class="{ 'is-active': step === 1, 'is-done': step > 1 }"
          >
            <span class="step-dot">{{ step > 1 ? '✓' : '1' }}</span>
            <span class="step-label">发送验证码</span>
          </div>
          <div class="step-line" :class="{ 'is-done': step > 1 }"></div>
          <div class="step-item" :class="{ 'is-active': step === 2 }">
            <span class="step-dot">{{ step > 2 ? '✓' : '2' }}</span>
            <span class="step-label">重置密码</span>
          </div>
        </div>

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
          size="large"
          class="auth-form"
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

          <template v-if="step === 2">
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
          </template>

          <el-button
            type="primary"
            size="large"
            class="w-full auth-submit"
            :loading="loading"
            native-type="submit"
          >
            {{ loading ? '处理中...' : buttonText }}
          </el-button>
        </el-form>

        <div class="auth-footer">
          <div class="auth-register text-secondary">
            想起密码了？
            <router-link to="/login" class="auth-link">返回登录</router-link>
          </div>
        </div>
      </div>
    </div>
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
  align-items: flex-start;
  min-height: 100vh;
  padding: 32px 16px;
  box-sizing: border-box;
  background:
    radial-gradient(60% 60% at 20% 0%, rgba(79, 110, 247, 0.08), transparent 70%),
    radial-gradient(50% 50% at 90% 100%, rgba(109, 139, 251, 0.1), transparent 70%);
}

.auth-card {
  width: 100%;
  max-width: 440px;
  border-radius: var(--campus-radius-lg, 20px);
  box-shadow: var(--campus-shadow-lg);
  overflow: hidden;
}

.auth-hero {
  position: relative;
  padding: 34px 32px 26px;
  text-align: center;
  background: linear-gradient(135deg, #4f6ef7 0%, #6d8bfb 60%, #8aa2fd 100%);
  color: #fff;
}

.auth-badge {
  width: 54px;
  height: 54px;
  margin: 0 auto 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.22);
  color: #fff;
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.18);
}

.auth-title {
  margin: 0;
  font-size: 25px;
  font-weight: 700;
  letter-spacing: -0.01em;
  color: #fff;
}

.auth-subtitle {
  margin: 8px 0 0;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.85);
}

.text-center {
  text-align: center;
}

.auth-body {
  padding: 28px 32px 32px;
}

/* 步骤指示器 */
.step-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 24px;
}

.step-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  min-width: 72px;
}

.step-dot {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
  background: var(--campus-surface-2, #f1f3f8);
  color: var(--campus-text-secondary, #64748b);
  border: 1px solid var(--campus-border, #e6e9f2);
  transition: all 0.25s;
}

.step-item.is-active .step-dot {
  background: var(--campus-primary, #4f6ef7);
  border-color: var(--campus-primary, #4f6ef7);
  color: #fff;
  box-shadow: 0 0 0 4px rgba(79, 110, 247, 0.12);
}

.step-item.is-done .step-dot {
  background: var(--campus-success, #22c55e);
  border-color: var(--campus-success, #22c55e);
  color: #fff;
}

.step-label {
  font-size: 13px;
  color: var(--campus-text-secondary, #64748b);
}

.step-item.is-active .step-label {
  color: var(--campus-primary, #4f6ef7);
  font-weight: 600;
}

.step-line {
  flex: 1;
  max-width: 56px;
  height: 2px;
  background: var(--campus-border, #e6e9f2);
  margin: 0 10px 22px;
  border-radius: 2px;
  transition: background 0.25s;
}

.step-line.is-done {
  background: var(--campus-success, #22c55e);
}

.auth-form :deep(.el-form-item__label) {
  font-weight: 600;
  font-size: 14px;
  padding-bottom: 8px;
}

.auth-form :deep(.el-input__wrapper) {
  border-radius: 12px;
  padding: 2px 14px;
  box-shadow: 0 0 0 1px var(--el-border-color, #e6e9f2) inset;
  transition: box-shadow 0.2s;
}

.auth-form :deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px var(--campus-primary-light, #6d8bfb) inset;
}

.auth-form :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 2px var(--campus-primary-soft, rgba(79, 110, 247, 0.35)) inset;
}

.auth-submit {
  height: 46px;
  margin-top: 4px;
  border-radius: 12px;
  font-size: 15px;
  letter-spacing: 0.02em;
}

.auth-footer {
  margin-top: 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.auth-register {
  font-size: 14px;
}

.auth-link {
  font-weight: 600;
}

.mb-4 {
  margin-bottom: 16px;
}
.w-full {
  width: 100%;
}

/* 移动端 */
@media (max-width: 768px) {
  .auth-page-center {
    min-height: calc(100vh - 40px);
    padding: 12px;
  }
  .auth-hero {
    padding: 26px 24px 20px;
  }
  .auth-body {
    padding: 24px 20px 28px;
  }
  .auth-title {
    font-size: 22px;
  }
}
</style>
