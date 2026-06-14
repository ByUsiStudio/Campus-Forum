<template>
  <v-container class="min-h-screen">
    <div class="error-page">
      <div class="error-icon">
        <v-icon color="primary" size="150">mdi-file-question</v-icon>
      </div>
      
      <v-card-title class="text-center text-h3 font-weight-bold mb-4">
        404 - 页面未找到
      </v-card-title>
      
      <v-card-text class="text-center">
        <p class="text-body-1 mb-2">您访问的页面不存在</p>
        <p class="text-body-2 text-secondary">该页面可能已被删除或URL地址有误</p>
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
.error-page {
  text-align: center;
  padding: 40px;
  max-width: 500px;
  margin: 0 auto;
}

.error-icon {
  margin-bottom: 32px;
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px);
  }
  50% {
    transform: translateY(-10px);
  }
}
</style>
