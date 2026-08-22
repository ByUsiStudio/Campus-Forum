<template>
  <div class="error-page">
    <div class="error-card">
      <div class="error-code">{{ code }}</div>
      <div class="error-icon-wrap">
        <el-icon :size="42" :class="'error-icon ' + (iconColorIconClass)"><component :is="iconComponent" /></el-icon>
      </div>
      <div class="error-title">{{ title }}</div>
      <div class="error-message">{{ message }}</div>
      <p v-if="detail" class="text-secondary error-detail">{{ detail }}</p>
      <div class="error-actions" v-if="actions.length">
        <el-button
          v-for="action in actions"
          :key="action.text"
          :type="action.color || 'primary'"
          round
          @click="handleAction(action)"
        >
          <el-icon v-if="action.icon"><component :is="action.icon" /></el-icon>
          <span v-if="action.icon">&nbsp;</span>{{ action.text }}
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { Warning, CircleCheck } from '@element-plus/icons-vue'

const props = defineProps({
  code: { type: [String, Number], default: '500' },
  title: { type: String, default: '出错了' },
  message: { type: String, default: '服务器内部错误' },
  detail: { type: String, default: '' },
  icon: { type: String, default: '500' },
  iconColor: { type: String, default: 'error' },
  actions: { type: Array, default: () => [] }
})

const iconComponent = computed(() => {
  const i = String(props.icon)
  if (i === 'success' || i === 'CircleCheck' || i === 'circle-check') return CircleCheck
  return Warning
})

const iconColorIconClass = computed(() => {
  const c = props.iconColor
  if (c === 'success') return 'is-success'
  if (c === 'warning') return 'is-warning'
  return 'is-error'
})

const handleAction = (action) => {
  if (action.callback) action.callback()
}
</script>

<style scoped>
.error-page {
  min-height: 65vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}
.error-card {
  background: var(--campus-surface);
  border: 1px solid var(--campus-border);
  border-radius: var(--campus-radius-lg);
  box-shadow: var(--campus-shadow);
  padding: 44px 40px;
  text-align: center;
  max-width: 460px;
  width: 100%;
}
.error-code {
  font-size: 64px;
  font-weight: 800;
  letter-spacing: -0.03em;
  line-height: 1;
  background: linear-gradient(135deg, var(--campus-primary), var(--campus-primary-light));
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  margin-bottom: 12px;
}
.error-icon-wrap {
  display: flex;
  justify-content: center;
  margin-bottom: 8px;
}
.error-icon.is-error { color: var(--campus-danger); }
.error-icon.is-warning { color: var(--campus-warning); }
.error-icon.is-success { color: var(--campus-success); }
.error-title {
  font-size: 20px;
  font-weight: 700;
  color: var(--campus-text);
  margin-bottom: 6px;
}
.error-message {
  color: var(--campus-text-secondary);
  font-size: 14px;
}
.error-detail {
  color: var(--campus-text-muted);
  font-size: 13px;
  margin-top: 4px;
}
.error-actions {
  margin-top: 20px;
  display: flex;
  gap: 10px;
  justify-content: center;
  flex-wrap: wrap;
}
</style>

