<template>
  <div class="auth-page-center">
    <el-card class="auth-card card-surface">
      <template #header>
        <div class="auth-header">
          <el-icon class="auth-logo"><User /></el-icon>
          <span class="auth-title">登录论坛</span>
        </div>
      </template>

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
        @submit.prevent="handleLogin"
      >
        <el-form-item label="用户名" prop="username">
          <el-input
            v-model="form.username"
            placeholder="请输入用户名"
            :prefix-icon="User"
          />
        </el-form-item>

        <el-form-item label="密码" prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="请输入密码"
            show-password
            :prefix-icon="Lock"
            @keyup.enter="handleLogin"
          />
        </el-form-item>

        <el-button
          type="primary"
          size="large"
          class="w-full"
          :loading="loading"
          @click="handleLogin"
        >
          {{ loading ? '登录中...' : '登录' }}
        </el-button>
      </el-form>

      <div class="auth-footer text-secondary">
        <router-link to="/forgot-password" class="mr-4">忘记密码？</router-link>
        还没有账号？
        <router-link to="/register" class="auth-link">立即注册</router-link>
      </div>
    </el-card>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { User, Lock } from '@element-plus/icons-vue'
import { authApi } from '@/api'
import { useUserStore } from '@/stores/user'
import { success, error } from '@/utils/message'

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
        error(error.value)
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
.mr-4 {
  margin-right: 16px;
}
.w-full {
  width: 100%;
}
</style>
