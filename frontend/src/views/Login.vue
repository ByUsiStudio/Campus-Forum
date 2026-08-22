<template>
  <div class="auth-page-center">
    <div class="auth-card card-surface">
      <div class="auth-hero">
        <div class="auth-badge"><el-icon :size="26"><User /></el-icon></div>
        <h1 class="auth-title">欢迎回来</h1>
        <p class="auth-subtitle text-secondary">登录校园论坛，继续你的分享</p>
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

        <el-form
          :model="form"
          :rules="rules"
          ref="formRef"
          label-position="top"
          class="auth-form"
          @submit.prevent="handleLogin"
        >
          <el-form-item label="用户名" prop="username">
            <el-input
              v-model="form.username"
              placeholder="请输入用户名"
              size="large"
              :prefix-icon="User"
            />
          </el-form-item>

          <el-form-item label="密码" prop="password">
            <el-input
              v-model="form.password"
              type="password"
              placeholder="请输入密码"
              size="large"
              show-password
              :prefix-icon="Lock"
              @keyup.enter="handleLogin"
            />
          </el-form-item>

          <el-button
            type="primary"
            size="large"
            class="w-full auth-submit"
            :loading="loading"
            @click="handleLogin"
          >
            {{ loading ? '登录中...' : '登录' }}
          </el-button>
        </el-form>

        <div class="auth-footer">
          <router-link to="/forgot-password" class="forgot-link">忘记密码？</router-link>
          <div class="auth-footer-sep"></div>
          <div class="auth-register text-secondary">
            还没有账号？
            <router-link to="/register" class="auth-link">立即注册</router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { User, Lock } from '@element-plus/icons-vue'
import { authApi } from '@/api'
import { useUserStore } from '@/stores/user'
import { success, error as showError } from '@/utils/message'

export default {
  name: 'Login',
  setup() {
    const router = useRouter()
    const route = useRoute()
    const userStore = useUserStore()

    const formRef = ref(null)
    const form = ref({
      username: '',
      password: ''
    })
    const error = ref('')
    const loading = ref(false)

    const rules = {
      username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
      password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
    }

    const checkInit = async () => {
      try {
        const response = await authApi.checkInit()
        if (!response.data.initialized) {
          router.push('/register')
        }
      } catch (err) {
        console.error('检查初始化失败', err)
      }
    }

    const handleLogin = async () => {
      try {
        await formRef.value?.validate()
      } catch (e) {
        return
      }

      error.value = ''
      loading.value = true

      try {
        const response = await authApi.login(form.value)
        const loginRes = response.data

        userStore.persist(loginRes)

        // 拉取完整资料字段（qq_number/signature/created_at 等），避免字段缺失
        userStore.refreshProfile().catch(() => {})

        success('登录成功')

        // 登录成功后跳转：优先回到被拦截的目标页，否则回首页
        const redirect = route.query.redirect
        if (redirect && redirect.startsWith('/') && !redirect.startsWith('/login')) {
          router.push(redirect)
        } else {
          router.push('/')
        }
      } catch (err) {
        error.value = err.response?.data?.error || '登录失败'
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
      rules,
      formRef,
      error,
      loading,
      User,
      Lock,
      handleLogin
    }
  }
}
</script>

<style scoped>
.auth-page-center {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  padding: 24px 16px;
  box-sizing: border-box;
  background:
    radial-gradient(60% 60% at 20% 0%, rgba(79, 110, 247, 0.08), transparent 70%),
    radial-gradient(50% 50% at 90% 100%, rgba(109, 139, 251, 0.1), transparent 70%);
}

.auth-card {
  width: 100%;
  max-width: 420px;
  border-radius: var(--campus-radius-lg, 20px);
  box-shadow: var(--campus-shadow-lg);
  overflow: hidden;
}

.auth-hero {
  position: relative;
  padding: 36px 32px 28px;
  text-align: center;
  background: linear-gradient(135deg, #4f6ef7 0%, #6d8bfb 60%, #8aa2fd 100%);
  color: #fff;
}

.auth-badge {
  width: 56px;
  height: 56px;
  margin: 0 auto 16px;
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
  font-size: 26px;
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

.forgot-link {
  font-size: 13px;
}

.auth-footer-sep {
  width: 40px;
  height: 1px;
  background: var(--campus-border, #e6e9f2);
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

/* 移动端：缩小内边距、卡片贴边 */
@media (max-width: 768px) {
  .auth-page-center {
    min-height: calc(100vh - 40px);
    padding: 12px;
  }
  .auth-hero {
    padding: 28px 24px 22px;
  }
  .auth-body {
    padding: 24px 20px 28px;
  }
  .auth-title {
    font-size: 23px;
  }
}
</style>
