<template>
  <div class="error-page">
    <v-card class="error-card">
      <div class="error-icon">
        <v-icon :color="iconColor" :size="120">
          {{ icon }}
        </v-icon>
      </div>
      
      <v-card-title class="text-center text-h4 mb-2">
        {{ title }}
      </v-card-title>
      
      <v-card-text class="text-center">
        <p class="text-body-1 mb-2">{{ message }}</p>
        <p v-if="detail" class="text-body-2 text-secondary">{{ detail }}</p>
      </v-card-text>
      
      <v-card-actions class="justify-center gap-3 mt-4">
        <v-btn
          v-for="action in actions"
          :key="action.text"
          :color="action.color || 'primary'"
          :variant="action.variant || 'flat'"
          @click="handleAction(action)"
        >
          <v-icon v-if="action.icon" start>{{ action.icon }}</v-icon>
          {{ action.text }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </div>
</template>

<script>
import { computed } from 'vue'

export default {
  name: 'ErrorPage',
  props: {
    code: {
      type: [String, Number],
      default: '500'
    },
    title: {
      type: String,
      default: '出错了'
    },
    message: {
      type: String,
      default: '服务器内部错误'
    },
    detail: {
      type: String,
      default: ''
    },
    icon: {
      type: String,
      default: 'mdi-alert-circle'
    },
    iconColor: {
      type: String,
      default: 'error'
    },
    actions: {
      type: Array,
      default: () => []
    }
  },
  setup(props) {
    const handleAction = (action) => {
      if (action.callback) {
        action.callback()
      }
    }
    
    return {
      handleAction
    }
  }
}
</script>

<style scoped>
.error-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.error-card {
  text-align: center;
  max-width: 500px;
  width: 100%;
  padding: 40px;
}

.error-icon {
  margin-bottom: 24px;
}
</style>
