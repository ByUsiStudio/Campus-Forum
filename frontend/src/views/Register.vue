<template>
  <div class="register-page">
    <div class="register-bg"></div>
    <div class="register-container">
      <div class="register-card animate-scale-in">
        <div class="register-header">
          <div class="register-icon">
            <i class="layui-icon" :class="isInit ? 'layui-icon-set' : 'layui-icon-user-add'"></i>
          </div>
          <h2 class="register-title">{{ isInit ? '初始化系统' : '注册新账号' }}</h2>
          <p class="register-subtitle">{{ isInit ? '请创建管理员账号以开始使用' : '加入我们，开始分享' }}</p>
        </div>

        <div class="register-body">
          <div v-if="error" class="error-message">
            <i class="layui-icon layui-icon-close-circle"></i>
            {{ error }}
          </div>
          <div v-if="success" class="success-message">
            <i class="layui-icon layui-icon-check-circle"></i>
            {{ success }}
          </div>

          <form @submit.prevent="handleRegister" class="register-form">
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
              <label class="form-label">QQ号</label>
              <div class="layui-input-group">
                <span class="layui-input-group-prepend">
                  <i class="layui-icon layui-icon-qq"></i>
                </span>
                <input 
                  type="text" 
                  v-model="form.qq_number" 
                  placeholder="请输入QQ号"
                  class="layui-input"
                  required
                >
              </div>
              <span class="form-hint">头像将默认使用QQ头像</span>
            </div>

            <div class="form-item">
              <label class="form-label">显示名称</label>
              <div class="layui-input-group">
                <span class="layui-input-group-prepend">
                  <i class="layui-icon layui-icon-card"></i>
                </span>
                <input 
                  type="text" 
                  v-model="form.display_name" 
                  placeholder="请输入显示名称"
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
              <span class="form-hint">至少6个字符</span>
            </div>

            <div class="form-item" v-if="!isInit">
              <label class="form-label">确认密码</label>
              <div class="layui-input-group">
                <span class="layui-input-group-prepend">
                  <i class="layui-icon layui-icon-verification"></i>
                </span>
                <input 
                  type="password" 
                  v-model="confirmPassword" 
                  placeholder="请再次输入密码"
                  class="layui-input"
                  required
                >
              </div>
            </div>

            <button 
              type="submit" 
              class="register-btn"
              :disabled="loading"
            >
              <i v-if="loading" class="layui-icon layui-icon-loading layui-anim layui-anim-spin"></i>
              <i v-else class="layui-icon" :class="isInit ? 'layui-icon-set' : 'layui-icon-user-add'"></i>
              {{ loading ? '处理中...' : (isInit ? '初始化系统' : '注册') }}
            </button>
          </form>

          <div v-if="!isInit" class="register-links">
            <span class="link-text">已有账号？</span>
            <router-link to="/login" class="link-item">立即登录</router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '../api'

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
    if (isInit.value) {
      await api.post('/auth/init-admin', form.value)
      success.value = '系统初始化成功！即将跳转到登录页...'
      setTimeout(() => {
        router.push('/login')
      }, 2000)
    } else {
      await api.post('/auth/register', form.value)
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
</script>

<style lang="less" scoped>
.register-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;

  .register-bg {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(135deg, rgba(255, 87, 34, 0.1) 0%, rgba(30, 159, 255, 0.05) 100%);

    &::before {
      content: '';
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      width: 600px;
      height: 600px;
      background: radial-gradient(circle, rgba(255, 87, 34, 0.08) 0%, transparent 60%);
      border-radius: 50%;
      animation: glowPulse 4s ease-in-out infinite;
    }
  }

  .register-container {
    position: relative;
    z-index: 1;
    width: 90%;
    max-width: 460px;
  }

  .register-card {
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(12px);
    border-radius: 16px;
    padding: 40px;
    box-shadow: 0 12px 40px rgba(0, 0, 0, 0.1);
    border: 1px solid rgba(255, 255, 255, 0.5);

    .register-header {
      text-align: center;
      margin-bottom: 32px;

      .register-icon {
        width: 64px;
        height: 64px;
        margin: 0 auto 16px;
        background: linear-gradient(135deg, #FF5722, #E64A19);
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 32px;
        color: #fff;
        box-shadow: 0 8px 24px rgba(255, 87, 34, 0.4);
        animation: float 3s ease-in-out infinite;
      }

      .register-title {
        font-size: 24px;
        font-weight: 700;
        color: #333;
        margin: 0 0 8px;
      }

      .register-subtitle {
        font-size: 14px;
        color: #999;
        margin: 0;
      }
    }

    .register-body {
      .error-message {
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 12px 16px;
        background: rgba(255, 87, 34, 0.1);
        color: #FF5722;
        border-radius: 6px;
        margin-bottom: 16px;
        font-size: 14px;
      }

      .success-message {
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 12px 16px;
        background: rgba(82, 196, 26, 0.1);
        color: #52C41A;
        border-radius: 6px;
        margin-bottom: 16px;
        font-size: 14px;
      }

      .register-form {
        .form-item {
          margin-bottom: 16px;

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

          .form-hint {
            display: block;
            font-size: 12px;
            color: #999;
            margin-top: 4px;
          }
        }

        .register-btn {
          width: 100%;
          padding: 14px;
          background: linear-gradient(135deg, #FF5722, #E64A19);
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
          box-shadow: 0 4px 16px rgba(255, 87, 34, 0.4);

          &:hover:not(:disabled) {
            transform: translateY(-3px);
            box-shadow: 0 8px 24px rgba(255, 87, 34, 0.5);
          }

          &:disabled {
            opacity: 0.7;
            cursor: not-allowed;
          }
        }
      }

      .register-links {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 8px;
        margin-top: 24px;

        .link-text {
          color: #999;
          font-size: 14px;
        }

        .link-item {
          color: #FF5722;
          text-decoration: none;
          font-size: 14px;
          font-weight: 500;
          transition: color 0.3s ease;

          &:hover {
            color: #E64A19;
          }
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
