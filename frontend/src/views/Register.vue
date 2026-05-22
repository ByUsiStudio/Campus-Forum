<template>
  <div class="d-flex justify-center align-center" style="min-height: 80vh;">
    <v-card width="100%" max-width="450" class="pa-6">
      <v-card-title class="text-h5 text-center pb-4" style="color: rgb(var(--v-theme-primary));">
        {{ isInit ? '初始化系统' : '注册新账号' }}
      </v-card-title>
      
      <v-card-text>
        <v-alert v-if="error" type="error" variant="tonal" class="mb-4">{{ error }}</v-alert>
        <v-alert v-if="success" type="success" variant="tonal" class="mb-4">{{ success }}</v-alert>
        
        <v-form @submit.prevent="handleRegister">
          <v-text-field
            v-model="form.username"
            label="用户名"
            variant="outlined"
            required
            prepend-inner-icon="mdi-account"
            class="mb-4"
          ></v-text-field>
          
          <v-text-field
            v-model="form.qq_number"
            label="QQ号"
            variant="outlined"
            required
            prepend-inner-icon="mdi-qqchat"
            class="mb-4"
            hint="头像将默认使用QQ头像"
          ></v-text-field>
          
          <v-text-field
            v-model="form.display_name"
            label="显示名称"
            variant="outlined"
            required
            prepend-inner-icon="mdi-card-account-details"
            class="mb-4"
          ></v-text-field>
          
          <v-text-field
            v-model="form.password"
            label="密码"
            variant="outlined"
            type="password"
            required
            minlength="6"
            prepend-inner-icon="mdi-lock"
            class="mb-4"
          ></v-text-field>
          
          <v-text-field
            v-if="!isInit"
            v-model="confirmPassword"
            label="确认密码"
            variant="outlined"
            type="password"
            required
            prepend-inner-icon="mdi-lock-check"
            class="mb-4"
          ></v-text-field>
          
          <v-btn
            type="submit"
            color="primary"
            block
            size="large"
            :loading="loading"
          >
            {{ loading ? '处理中...' : (isInit ? '初始化系统' : '注册') }}
          </v-btn>
        </v-form>
        
        <div class="text-center mt-4 text-body-2" v-if="!isInit">
          已有账号？ <router-link to="/login" class="text-primary">立即登录</router-link>
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
