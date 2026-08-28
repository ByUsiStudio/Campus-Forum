<template>
  <div class="auth-page-center">
    <div class="auth-card card-surface">
      <div class="auth-hero">
        <div class="auth-badge"><el-icon :size="26"><UserFilled /></el-icon></div>
        <h1 class="auth-title">{{ isInit ? '初始化系统' : '创建账号' }}</h1>
        <p class="auth-subtitle">
          {{ isInit ? '请创建管理员账号以开始使用' : '加入我们，开始分享' }}
        </p>
      </div>

      <div class="auth-body">
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
          @submit.prevent="handleRegister"
        >
          <el-row :gutter="16">
            <el-col :span="24" :md="12">
              <el-form-item label="用户名" prop="username">
                <el-input
                  v-model="form.username"
                  placeholder="请输入用户名"
                  :prefix-icon="User"
                />
              </el-form-item>
            </el-col>
            <el-col :span="24" :md="12">
              <el-form-item label="QQ号" prop="qq_number">
                <el-input
                  v-model="form.qq_number"
                  placeholder="请输入QQ号"
                  :prefix-icon="ChatDotRound"
                />
                <div class="form-tip text-secondary">头像将默认使用QQ头像</div>
              </el-form-item>
            </el-col>
          </el-row>

          <el-form-item label="显示名称" prop="display_name">
            <el-input
              v-model="form.display_name"
              placeholder="请输入显示名称"
              :prefix-icon="Postcard"
            />
          </el-form-item>

          <el-form-item label="密码" prop="password">
            <el-input
              v-model="form.password"
              type="password"
              placeholder="请输入密码（至少6个字符）"
              show-password
              :prefix-icon="Lock"
            />
          </el-form-item>

          <el-form-item v-if="!isInit" label="确认密码" prop="confirmPassword">
            <el-input
              v-model="confirmPassword"
              type="password"
              placeholder="请再次输入密码"
              show-password
              :prefix-icon="Lock"
              @keyup.enter="handleRegister"
            />
          </el-form-item>

          <el-button
            type="primary"
            size="large"
            class="w-full auth-submit"
            :loading="loading"
            @click="handleRegister"
          >
            {{ loading ? '处理中...' : (isInit ? '初始化系统' : '注册') }}
          </el-button>
        </el-form>

        <div v-if="!isInit" class="auth-footer">
          <div class="auth-register text-secondary">
            已有账号？
            <router-link to="/login" class="auth-link">立即登录</router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { User, UserFilled, Lock, ChatDotRound, Postcard } from '@element-plus/icons-vue'
import { authApi } from '@/api'
import { success as jcSuccess, error as showError } from '@/utils/message'

export default {
  name: 'Register',
  setup() {
    const router = useRouter()
    const formRef = ref(null)
    const form = ref({
      username: '',
      qq_number: '',
      display_name: '',
      password: ''
    })
    const confirmPassword = ref('')
    const error = ref('')
    const success = ref('')
    const loading = ref(false)
    const isInit = ref(false)

    const validateConfirmPassword = (rule, value, callback) => {
      if (value !== form.value.password) {
        callback(new Error('两次输入的密码不一致'))
      } else {
        callback()
      }
    }

    const rules = {
      username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
      qq_number: [{ required: true, message: '请输入QQ号', trigger: 'blur' }],
      display_name: [{ required: true, message: '请输入显示名称', trigger: 'blur' }],
      password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
        { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
      ],
      confirmPassword: [{ validator: validateConfirmPassword, trigger: 'blur' }]
    }

    const checkInit = async () => {
      try {
        const response = await authApi.checkInit()
        if (!response.data.initialized) {
          isInit.value = true
        }
      } catch (err) {
        console.error('检查初始化失败', err)
      }
    }

    const handleRegister = async () => {
      try {
        await formRef.value?.validate()
      } catch (e) {
        return
      }

      error.value = ''
      success.value = ''
      loading.value = true

      try {
        const payload = {
          username: form.value.username,
          qq_number: form.value.qq_number,
          display_name: form.value.display_name,
          password: form.value.password
        }

        if (isInit.value) {
          await authApi.initAdmin(payload)
          success.value = '系统初始化成功！即将跳转到登录页...'
          jcSuccess('系统初始化成功')
        } else {
          await authApi.register(payload)
          success.value = '注册成功！即将跳转到登录页...'
          jcSuccess('注册成功')
        }
        setTimeout(() => {
          router.push('/login')
        }, 2000)
      } catch (err) {
        error.value = err.response?.data?.error || '操作失败'
        showError(error.value)
      } finally {
        loading.value = false
      }
    }

    onMounted(() => {
      checkInit()
    })

    return {
      form,
      formRef,
      rules,
      confirmPassword,
      error,
      success,
      loading,
      isInit,
      User,
      UserFilled,
      Lock,
      ChatDotRound,
      Postcard,
      handleRegister
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
  max-width: 470px;
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

.auth-body {
  padding: 28px 32px 32px;
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

.form-tip {
  font-size: 12px;
  line-height: 1.4;
  margin-top: 4px;
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