<template>
  <v-container class="min-h-screen">
    <div class="error-page">
      <div class="error-icon">
        <v-icon color="warning" size="150">mdi-alert-lock</v-icon>
      </div>
      
      <v-card-title class="text-center text-h3 font-weight-bold mb-4">
        403 - 禁止访问
      </v-card-title>
      
      <v-card-text class="text-center">
        <p class="text-body-1 mb-2">您没有权限访问此页面</p>
        <p class="text-body-2 text-secondary">请检查您的权限或联系管理员获取帮助</p>
      </v-card-text>
      
      <v-card-actions class="justify-center gap-3 mt-6">
        <v-btn color="primary" @click="goHome">
          <v-icon start>mdi-home</v-icon>
          返回首页
        </v-btn>
        <v-btn variant="outlined" @click="goBack">
          <v-icon start>mdi-arrow-left</v-icon>
          返回上页
        </v-btn>
      </v-card-actions>
    </div>
  </v-container>
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
}

.error-icon {
  margin-bottom: 32px;
  animation: pulse 2s ease-in-out infinite;
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
