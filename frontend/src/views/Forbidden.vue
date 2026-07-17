<template>
  <div class="container error-page">
    <div class="error-icon">
      <i class="fa-solid fa-shield-halved" style="font-size: 150px; color: #FAAD14;"></i>
    </div>
    
    <h1 class="error-title">403 - 禁止访问</h1>
    
    <div class="error-text">
      <p class="text-lg mb-2">您没有权限访问此页面</p>
      <p class="text-secondary">请检查您的权限或联系管理员获取帮助</p>
    </div>
    
    <div class="error-actions mt-6 d-flex justify-center gap-3">
      <button class="layui-btn" @click="goHome">
        <i class="fa-solid fa-house mr-1"></i>
        返回首页
      </button>
      <button class="layui-btn layui-btn-primary" @click="goBack">
        <i class="fa-solid fa-arrow-left mr-1"></i>
        返回上页
      </button>
    </div>
  </div>
</template>

<script>
import { useRouter } from 'vue-router'

export default {
  name: 'Forbidden',
  setup() {
    const router = useRouter()
    
    const goHome = () => {
      router.push('/')
    }
    
    const goBack = () => {
      if (window.history.length > 1) {
        router.back()
      } else {
        router.push('/')
      }
    }
    
    return {
      goHome,
      goBack
    }
  }
}
</script>

<style scoped>
.error-page {
  text-align: center;
  padding: 40px;
  max-width: 500px;
  margin: 0 auto;
  min-height: 80vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.error-icon {
  margin-bottom: 32px;
  animation: pulse 2s ease-in-out infinite;
}

.error-title {
  font-size: 32px;
  font-weight: 700;
  margin-bottom: 20px;
  color: #333;
}

.error-text {
  margin-bottom: 30px;
}

.text-lg {
  font-size: 18px;
}

.text-secondary {
  color: #999;
}

.d-flex {
  display: flex;
}

.justify-center {
  justify-content: center;
}

.gap-3 {
  gap: 15px;
}

.mr-1 {
  margin-right: 5px;
}

.mb-2 {
  margin-bottom: 10px;
}

.mt-6 {
  margin-top: 30px;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.7;
  }
}
</style>