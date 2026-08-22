<template>
  <div class="page-container error-page">
    <el-result icon="warning" title="403 - 禁止访问" sub-title="您没有权限访问此页面">
      <template #extra>
        <div class="error-actions">
          <el-button type="primary" :icon="Home" @click="goHome">
            返回首页
          </el-button>
          <el-button :icon="ArrowLeft" @click="goBack">
            返回上页
          </el-button>
        </div>
        <p class="text-secondary error-hint">请检查您的权限或联系管理员获取帮助</p>
      </template>
    </el-result>
  </div>
</template>

<script>
import { useRouter } from 'vue-router'
import { Home, ArrowLeft } from '@element-plus/icons-vue'

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
      goBack,
      Home,
      ArrowLeft
    }
  }
}
</script>

<style scoped>
.error-page {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: calc(100vh - 120px);
  padding: 40px;
}

.error-actions {
  display: flex;
  gap: 16px;
  justify-content: center;
  margin-bottom: 16px;
}

.error-hint {
  font-size: 14px;
}
</style>
