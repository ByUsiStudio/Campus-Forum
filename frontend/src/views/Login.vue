<template>
  <div class="d-flex justify-center align-center" style="min-height: 80vh;">
    <v-card width="100%" max-width="400" class="pa-6">
      <v-card-title class="text-h5 text-center pb-4" style="color: rgb(var(--v-theme-primary));">
        登录论坛
      </v-card-title>
      
      <v-card-text>
        <v-alert v-if="error" type="error" variant="tonal" class="mb-4">{{ error }}</v-alert>
        
        <v-form @submit.prevent="handleLogin">
          <v-text-field
            v-model="form.username"
            label="用户名"
            variant="outlined"
            required
            prepend-inner-icon="mdi-account"
            class="mb-4"
          ></v-text-field>
          
          <v-text-field
            v-model="form.password"
            label="密码"
            variant="outlined"
            type="password"
            required
            prepend-inner-icon="mdi-lock"
            class="mb-4"
          ></v-text-field>
          
          <v-btn
            type="submit"
            color="primary"
            block
            size="large"
            :loading="loading"
          >
            {{ loading ? '登录中...' : '登录' }}
          </v-btn>
        </v-form>
        
        <div class="text-center mt-4 text-body-2">
          还没有账号？ <router-link to="/register" class="text-primary">立即注册</router-link>
        </div>
      </v-card-text>
    </v-card>
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
