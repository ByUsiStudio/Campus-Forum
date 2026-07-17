<template>
  <div class="not-found-page">
    <div class="error-card">
      <div class="error-icon">
        <i class="fa-solid fa-circle-question"></i>
      </div>
      
      <h1 class="error-title">404 - 页面未找到</h1>
      
      <p class="error-message">您访问的页面不存在</p>
      <p class="error-detail">该页面可能已被删除或URL地址有误</p>
      
      <div class="error-actions">
        <button class="layui-btn layui-btn-normal" @click="goHome">
          <i class="fa-solid fa-house mr-2"></i>
          返回首页
        </button>
        <button class="layui-btn layui-btn-primary" @click="goBack">
          <i class="fa-solid fa-arrow-left mr-2"></i>
          返回上页
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { useRouter } from 'vue-router'

export default {
  name: 'NotFound',
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
.not-found-page {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: #f5f5f5;
  padding: 24px;
}

.error-card {
  background: white;
  border-radius: 12px;
  padding: 48px;
  text-align: center;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  max-width: 500px;
}

.error-icon {
  margin-bottom: 32px;
  
  i {
    font-size: 120px;
    color: var(--primary);
    animation: float 3s ease-in-out infinite;
  }
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px);
  }
  50% {
    transform: translateY(-10px);
  }
}

.error-title {
  font-size: 24px;
  font-weight: 600;
  color: #333;
  margin-bottom: 16px;
}

.error-message {
  font-size: 16px;
  color: #666;
  margin-bottom: 8px;
}

.error-detail {
  font-size: 14px;
  color: #999;
  margin-bottom: 32px;
}

.error-actions {
  display: flex;
  gap: 12px;
  justify-content: center;
}

.mr-2 {
  margin-right: 8px;
}
</style>
