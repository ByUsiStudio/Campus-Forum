<template>
  <div class="register-container">
    <div class="register-box">
      <h2>{{ isInit ? '初始化系统' : '注册新账号' }}</h2>
      <div v-if="error" class="error-message">{{ error }}</div>
      <div v-if="success" class="success-message">{{ success }}</div>
      
      <form @submit.prevent="handleRegister">
        <div class="form-group">
          <label>用户名</label>
          <input type="text" v-model="form.username" required>
        </div>
        <div class="form-group">
          <label>QQ号</label>
          <input type="text" v-model="form.qq_number" required>
          <small>头像将默认使用QQ头像</small>
        </div>
        <div class="form-group">
          <label>显示名称</label>
          <input type="text" v-model="form.display_name" required>
        </div>
        <div class="form-group">
          <label>密码</label>
          <input type="password" v-model="form.password" required minlength="6">
        </div>
        <div class="form-group" v-if="!isInit">
          <label>确认密码</label>
          <input type="password" v-model="confirmPassword" required>
        </div>
        <button type="submit" class="btn btn-primary" :disabled="loading">
          {{ loading ? '处理中...' : (isInit ? '初始化系统' : '注册') }}
        </button>
      </form>
      
      <div class="login-link" v-if="!isInit">
        已有账号？ <router-link to="/login">立即登录</router-link>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '../api'

export default {
  name: 'Register',
  setup() {
    const router = useRouter()
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

    const checkInit = async () => {
      try {
        const response = await api.get('/auth/check-init')
        if (!response.data.initialized) {
          isInit.value = true
        }
      } catch (err) {
        console.error('检查初始化失败', err)
      }
    }

    const handleRegister = async () => {
      error.value = ''
      success.value = ''
      
      if (!isInit.value && form.value.password !== confirmPassword.value) {
        error.value = '两次输入的密码不一致'
        return
      }
      
      if (form.value.password.length < 6) {
        error.value = '密码长度不能少于6位'
        return
      }
      
      loading.value = true
      
      try {
        let response
        if (isInit.value) {
          response = await api.post('/auth/init-admin', form.value)
          success.value = '系统初始化成功！即将跳转到登录页...'
          setTimeout(() => {
            router.push('/login')
          }, 2000)
        } else {
          response = await api.post('/auth/register', form.value)
          success.value = '注册成功！即将跳转到登录页...'
          setTimeout(() => {
            router.push('/login')
          }, 2000)
        }
      } catch (err) {
        error.value = err.response?.data?.error || '操作失败'
      } finally {
        loading.value = false
      }
    }

    onMounted(() => {
      checkInit()
    })

    return {
      form,
      confirmPassword,
      error,
      success,
      loading,
      isInit,
      handleRegister
    }
  }
}
</script>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 80vh;
}

.register-box {
  background: white;
  padding: 40px;
  border-radius: 12px;
  width: 100%;
  max-width: 450px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.register-box h2 {
  text-align: center;
  margin-bottom: 30px;
  color: #059669;
}

.register-box small {
  color: #6b7280;
  font-size: 12px;
}

.login-link {
  text-align: center;
  margin-top: 20px;
  font-size: 14px;
}

.login-link a {
  color: #10b981;
  text-decoration: none;
}
</style>