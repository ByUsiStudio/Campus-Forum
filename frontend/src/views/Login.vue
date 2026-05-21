<template>
  <div class="login-container">
    <div class="login-box">
      <h2>登录论坛</h2>
      <div v-if="error" class="error-message">{{ error }}</div>
      <form @submit.prevent="handleLogin">
        <div class="form-group">
          <label>用户名</label>
          <input type="text" v-model="form.username" required>
        </div>
        <div class="form-group">
          <label>密码</label>
          <input type="password" v-model="form.password" required>
        </div>
        <button type="submit" class="btn btn-primary" :disabled="loading">
          {{ loading ? '登录中...' : '登录' }}
        </button>
      </form>
      <div class="register-link">
        还没有账号？ <router-link to="/register">立即注册</router-link>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '../api'

export default {
  name: 'Login',
  setup() {
    const router = useRouter()
    const form = ref({
      username: '',
      password: ''
    })
    const error = ref('')
    const loading = ref(false)

    const checkInit = async () => {
      try {
        const response = await api.get('/auth/check-init')
        if (!response.data.initialized) {
          router.push('/register')
        }
      } catch (err) {
        console.error('检查初始化失败', err)
      }
    }

    const handleLogin = async () => {
      error.value = ''
      loading.value = true
      
      try {
        const response = await api.post('/auth/login', form.value)
        const { token, user } = response.data
        
        localStorage.setItem('token', token)
        localStorage.setItem('user', JSON.stringify(user))
        
        router.push('/')
      } catch (err) {
        error.value = err.response?.data?.error || '登录失败'
      } finally {
        loading.value = false
      }
    }

    onMounted(() => {
      checkInit()
    })

    return {
      form,
      error,
      loading,
      handleLogin
    }
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 80vh;
}

.login-box {
  background: white;
  padding: 40px;
  border-radius: 12px;
  width: 100%;
  max-width: 400px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.login-box h2 {
  text-align: center;
  margin-bottom: 30px;
  color: #059669;
}

.register-link {
  text-align: center;
  margin-top: 20px;
  font-size: 14px;
}

.register-link a {
  color: #10b981;
  text-decoration: none;
}
</style>