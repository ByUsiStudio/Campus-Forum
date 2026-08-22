<template>
  <div class="auth-page-center">
    <el-card class="auth-card card-surface">
      <template #header>
        <div class="auth-header">
          <el-icon class="auth-logo"><UserFilled /></el-icon>
          <span class="auth-title">{{ isInit ? '初始化系统' : '注册新账号' }}</span>
        </div>
        <div class="auth-subtitle text-secondary">
          {{ isInit ? '请创建管理员账号以开始使用' : '加入我们，开始分享' }}
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
        @submit.prevent="handleRegister"
      >
        <el-form-item label="用户名" prop="username">
          <el-input
            v-model="form.username"
            placeholder="请输入用户名"
            :prefix-icon="User"
          />
        </el-form-item>

        <el-form-item label="QQ号" prop="qq_number">
          <el-input
            v-model="form.qq_number"
            placeholder="请输入QQ号"
            :prefix-icon="ChatDotRound"
          />
          <div class="el-form-item__tip text-secondary">头像将默认使用QQ头像</div>
        </el-form-item>

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
          class="w-full"
          :loading="loading"
          @click="handleRegister"
        >
          {{ loading ? '处理中...' : (isInit ? '初始化系统' : '注册') }}
        </el-button>
      </el-form>

      <div v-if="!isInit" class="auth-footer text-secondary">
        已有账号？
        <router-link to="/login" class="auth-link">立即登录</router-link>
      </div>
    </el-card>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { User, UserFilled, Lock, ChatDotRound, Postcard } from '@element-plus/icons-vue'
import { authApi } from '@/api'
import { success, error } from '@/utils/message'

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
      if (!isInit.value && value !== form.value.password) {
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
        // 构造提交数据，剔除确认密码字段
        const payload = {
          username: form.value.username,
          qq_number: form.value.qq_number,
          display_name: form.value.display_name,
          password: form.value.password
        }

        if (isInit.value) {
          await authApi.initAdmin(payload)
          success.value = '系统初始化成功！即将跳转到登录页...'
          success('系统初始化成功')
        } else {
          await authApi.register(payload)
          success.value = '注册成功！即将跳转到登录页...'
          success('注册成功')
        }
        setTimeout(() => {
          router.push('/login')
        }, 2000)
      } catch (err) {
        error.value = err.response?.data?.error || '操作失败'
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

.auth-subtitle {
  text-align: center;
  font-size: 13px;
  margin-bottom: 4px;
}

.auth-footer {
  text-align: center;
  margin-top: 16px;
  font-size: 13px;
}

.auth-link {
  color: var(--campus-primary);
}

.el-form-item__tip {
  font-size: 12px;
  line-height: 1.4;
  margin-top: 4px;
}

.mb-4 {
  margin-bottom: 16px;
}
.w-full {
  width: 100%;
}
</style>
