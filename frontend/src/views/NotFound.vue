<template>
  <div class="error-page">
    <div class="error-art" aria-hidden="true">
      <div class="map-shape">
        <el-icon :size="56" color="#fff"><Compass /></el-icon>
      </div>
      <div class="deco deco-1"></div>
      <div class="deco deco-2"></div>
    </div>

    <div class="error-code">
      <span class="four">4</span>
      <span class="zero">0</span>
      <span class="four">4</span>
    </div>
    <h1 class="error-title">页面不存在</h1>
    <p class="error-desc">您访问的页面不存在或已被移除</p>

    <div class="error-actions">
      <el-button type="primary" size="large" @click="goHome">
        <el-icon class="mr-1"><HomeFilled /></el-icon>
        返回首页
      </el-button>
      <el-button size="large" @click="goBack">
        <el-icon class="mr-1"><ArrowLeft /></el-icon>
        返回上页
      </el-button>
    </div>

    <p class="error-hint">该页面可能已被删除或 URL 地址有误</p>
  </div>
</template>

<script>
import { useRouter } from 'vue-router'
import { HomeFilled, ArrowLeft, Compass } from '@element-plus/icons-vue'

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
      ArrowLeft,
      Compass
    }
  }
}
</script>

<style scoped>
.error-page {
  min-height: calc(100vh - 120px);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 40px 16px;
}

.error-art {
  position: relative;
  margin-bottom: 24px;
}

.map-shape {
  width: 120px;
  height: 120px;
  border-radius: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, var(--campus-primary-light), var(--campus-primary-dark));
  box-shadow: 0 16px 32px rgba(79, 110, 247, 0.35);
  position: relative;
  z-index: 2;
}

.deco {
  position: absolute;
  border-radius: 50%;
  background: var(--campus-primary-soft);
}

.deco-1 {
  width: 46px;
  height: 46px;
  top: -14px;
  left: -20px;
}

.deco-2 {
  width: 28px;
  height: 28px;
  bottom: -10px;
  right: -16px;
  background: rgba(245, 158, 11, 0.2);
}

.error-code {
  font-size: 80px;
  font-weight: 800;
  line-height: 1;
  letter-spacing: 0.04em;
}

.four {
  background: linear-gradient(120deg, var(--campus-primary), var(--campus-primary-dark));
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.zero {
  background: linear-gradient(120deg, var(--campus-warning), #f97316);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.error-title {
  font-size: 22px;
  font-weight: 700;
  color: var(--campus-text);
  margin: 12px 0 8px;
}

.error-desc {
  font-size: 15px;
  color: var(--campus-text-secondary);
  margin: 0 0 24px;
}

.error-actions {
  display: flex;
  gap: 12px;
  justify-content: center;
  flex-wrap: wrap;
}

.error-hint {
  font-size: 13px;
  color: var(--campus-text-muted);
  margin: 20px 0 0;
}

.mr-1 {
  margin-right: 4px;
}

@media (max-width: 480px) {
  .error-code {
    font-size: 60px;
  }
}
</style>
