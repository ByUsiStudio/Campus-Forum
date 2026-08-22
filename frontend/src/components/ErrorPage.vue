<template>
  <div class="error-page">
    <el-result :icon="icon" :title="title" :sub-title="message">
      <template #extra>
        <p v-if="detail" class="text-secondary">{{ detail }}</p>
        <div class="error-actions">
          <el-button
            v-for="action in actions"
            :key="action.text"
            :type="action.color || 'primary'"
            @click="handleAction(action)"
          >
            <el-icon v-if="action.icon"><component :is="action.icon" /></el-icon>
            <span v-if="action.icon">&nbsp;</span>{{ action.text }}
          </el-button>
        </div>
      </template>
    </el-result>
  </div>
</template>

<script setup>
defineProps({
  code: { type: [String, Number], default: '500' },
  title: { type: String, default: '出错了' },
  message: { type: String, default: '服务器内部错误' },
  detail: { type: String, default: '' },
  icon: { type: String, default: '500' },
  iconColor: { type: String, default: 'error' },
  actions: { type: Array, default: () => [] }
})

const handleAction = (action) => {
  if (action.callback) action.callback()
}
</script>

<style scoped>
.error-page {
  min-height: 60vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}
.error-actions {
  margin-top: 12px;
  display: flex;
  gap: 8px;
  justify-content: center;
}
</style>
