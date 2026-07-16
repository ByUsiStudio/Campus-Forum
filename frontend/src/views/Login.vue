<template>
  <div class="login-page">
    <div class="login-bg"></div>
    <div class="login-container">
      <div class="login-card animate-scale-in">
        <div class="login-header">
          <div class="login-icon">
            <i class="layui-icon layui-icon-username"></i>
          </div>
          <h2 class="login-title">登录论坛</h2>
        </div>

        <div class="login-body">
          <div v-if="error" class="error-message">
            <i class="layui-icon layui-icon-close-circle"></i>
            {{ error }}
          </div>

          <form @submit.prevent="handleLogin" class="login-form">
            <div class="form-item">
              <label class="form-label">用户名</label>
              <div class="layui-input-group">
                <span class="layui-input-group-prepend">
                  <i class="layui-icon layui-icon-user"></i>
                </span>
                <input 
                  type="text" 
                  v-model="form.username" 
                  placeholder="请输入用户名"
                  class="layui-input"
                  required
                >
              </div>
            </div>

            <div class="form-item">
              <label class="form-label">密码</label>
              <div class="layui-input-group">
                <span class="layui-input-group-prepend">
                  <i class="layui-icon layui-icon-password"></i>
                </span>
                <input 
                  type="password" 
                  v-model="form.password" 
                  placeholder="请输入密码"
                  class="layui-input"
                  required
                >
              </div>
            </div>

            <button 
              type="submit" 
              class="login-btn"
              :disabled="loading"
            >
              <i v-if="loading" class="layui-icon layui-icon-loading layui-anim layui-anim-spin"></i>
              {{ loading ? '登录中...' : '登录' }}
            </button>
          </form>

          <div class="login-links">
            <router-link to="/forgot-password" class="link-item">忘记密码？</router-link>
            <span class="link-divider">|</span>
            <router-link to="/register" class="link-item register-link">立即注册</router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useStore } from '../stores'
import { commonApi } from '../api'

const router = useRouter()
const store = useStore()

const form = ref({
  username: '',
  password: ''
})
const error = ref('')
const loading = ref(false)

const checkInit = async () => {
  try {
    const response = await commonApi.checkInit()
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
    const response = await commonApi.login(form.value)
    const { token, refresh_token, expires_in, user } = response.data

    store.login(user, token, refresh_token, expires_in)
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
</script>

<style lang="less" scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;

  .login-bg {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(135deg, rgba(30, 159, 255, 0.1) 0%, rgba(255, 87, 34, 0.05) 100%);

    &::before {
      content: '';
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      width: 600px;
      height: 600px;
      background: radial-gradient(circle, rgba(30, 159, 255, 0.08) 0%, transparent 60%);
      border-radius: 50%;
      animation: glowPulse 4s ease-in-out infinite;
    }
  }

  .login-container {
    position: relative;
    z-index: 1;
    width: 90%;
    max-width: 420px;
  }

  .login-card {
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(12px);
    border-radius: 16px;
    padding: 40px;
    box-shadow: 0 12px 40px rgba(0, 0, 0, 0.1);
    border: 1px solid rgba(255, 255, 255, 0.5);

    .login-header {
      text-align: center;
      margin-bottom: 32px;

      .login-icon {
        width: 64px;
        height: 64px;
        margin: 0 auto 16px;
        background: linear-gradient(135deg, #1E9FFF, #0086E6);
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 32px;
        color: #fff;
        box-shadow: 0 8px 24px rgba(30, 159, 255, 0.4);
        animation: float 3s ease-in-out infinite;
      }

      .login-title {
        font-size: 24px;
        font-weight: 700;
        color: #333;
        margin: 0;
      }
    }

    .login-body {
      .error-message {
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 12px 16px;
        background: rgba(255, 87, 34, 0.1);
        color: #FF5722;
        border-radius: 6px;
        margin-bottom: 20px;
        font-size: 14px;
      }

      .login-form {
        .form-item {
          margin-bottom: 20px;

          .form-label {
            display: block;
            font-size: 14px;
            font-weight: 500;
            color: #333;
            margin-bottom: 8px;
          }

          .layui-input-group {
            .layui-input-group-prepend {
              background: #f5f5f5;
              border-radius: 6px 0 0 6px;
            }

            .layui-input {
              border-radius: 0 6px 6px 0;
            }
          }
        }

        .login-btn {
          width: 100%;
          padding: 14px;
          background: linear-gradient(135deg, #1E9FFF, #0086E6);
          color: #fff;
          border: none;
          border-radius: 8px;
          font-size: 16px;
          font-weight: 600;
          cursor: pointer;
          display: flex;
          align-items: center;
          justify-content: center;
          gap: 8px;
          transition: all 0.3s @ease-out-back;
          box-shadow: 0 4px 16px rgba(30, 159, 255, 0.4);

          &:hover:not(:disabled) {
            transform: translateY(-3px);
            box-shadow: 0 8px 24px rgba(30, 159, 255, 0.5);
          }

          &:disabled {
            opacity: 0.7;
            cursor: not-allowed;
          }
        }
      }

      .login-links {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 12px;
        margin-top: 24px;

        .link-item {
          color: #999;
          text-decoration: none;
          font-size: 14px;
          transition: color 0.3s ease;

          &:hover {
            color: #1E9FFF;
          }

          &.register-link {
            color: #1E9FFF;
            font-weight: 500;
          }
        }

        .link-divider {
          color: #e6e6e6;
        }
      }
    }
  }
}

@keyframes glowPulse {
  0%, 100% {
    transform: translate(-50%, -50%) scale(1);
    opacity: 0.5;
  }
  50% {
    transform: translate(-50%, -50%) scale(1.2);
    opacity: 0.8;
  }
}

@keyframes float {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-8px);
  }
}
</style>
