<template>
  <v-container fluid class="fill-height bg-grey-lighten-4">
    <v-row justify="center" align="center">
      <v-col cols="12" sm="8" md="5" lg="4">
        <v-card class="pa-6" elevation="2">
          <div class="text-center mb-6">
            <v-icon size="64" color="primary" class="mb-4">mdi-forum</v-icon>
            <v-card-title class="text-h5 font-weight-bold" style="color: rgb(var(--v-theme-primary));">
              登录论坛
            </v-card-title>
            <v-card-subtitle class="text-body-2 text-medium-emphasis">
              欢迎回来！请登录您的账号
            </v-card-subtitle>
          </div>
          
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
                :rules="[v => !!v || '请输入用户名']"
              ></v-text-field>
              
              <v-text-field
                v-model="form.password"
                label="密码"
                variant="outlined"
                type="password"
                required
                prepend-inner-icon="mdi-lock"
                class="mb-4"
                :rules="[v => !!v || '请输入密码']"
              ></v-text-field>
              
              <v-btn
                type="submit"
                color="primary"
                block
                size="large"
                :loading="loading"
                class="mb-4"
              >
                <v-icon start>mdi-login</v-icon>
                {{ loading ? '登录中...' : '登录' }}
              </v-btn>
            </v-form>
            
            <div class="d-flex justify-space-between align-center text-body-2">
              <router-link to="/forgot-password" class="text-secondary">
                <v-icon size="small" class="mr-1">mdi-key</v-icon>
                忘记密码？
              </router-link>
              <span class="text-medium-emphasis">
                还没有账号？
                <router-link to="/register" class="text-primary font-weight-medium">立即注册</router-link>
              </span>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
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
        
        // 通知 App.vue 更新状态
        window.dispatchEvent(new Event('user-updated'))
        
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
