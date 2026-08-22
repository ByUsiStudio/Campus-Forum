<template>
  <div class="page-container error-container">
    <div class="card-surface error-page">
      <el-result icon="404" title="404" sub-title="您访问的页面不存在">
        <template #extra>
          <p class="text-secondary mb-4">该页面可能已被删除或URL地址有误</p>
          <div class="action-row">
            <el-button type="primary" @click="goHome">
              <el-icon class="mr-1"><HomeFilled /></el-icon>
              返回首页
            </el-button>
            <el-button @click="goBack">
              <el-icon class="mr-1"><ArrowLeft /></el-icon>
              返回上页
            </el-button>
          </div>
        </template>
      </el-result>
    </div>
  </div>
</template>

<script>
import { useRouter } from 'vue-router'
import { HomeFilled, ArrowLeft } from '@element-plus/icons-vue'

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
      goBack,
      HomeFilled,
      ArrowLeft
    }
  }
}
</script>

<style scoped>
.error-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 16px;
  box-sizing: border-box;
}

.error-page {
  width: 100%;
  max-width: 500px;
  padding: 40px;
  text-align: center;
}

.action-row {
  display: flex;
  justify-content: center;
  gap: 12px;
}
</style>
