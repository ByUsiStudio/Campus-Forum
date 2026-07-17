<template>
  <div class="error-page">
    <div class="error-card layui-card">
      <div class="error-icon">
        <i :class="getIconClass(icon)" :style="{ color: getIconColor(iconColor) }" style="font-size: 120px;"></i>
      </div>
      
      <h2 class="error-title">{{ title }}</h2>
      
      <div class="error-content">
        <p class="text-body-1 mb-2">{{ message }}</p>
        <p v-if="detail" class="text-body-2 text-secondary">{{ detail }}</p>
      </div>
      
      <div class="error-actions d-flex justify-center gap-3 mt-4">
        <button
          v-for="action in actions"
          :key="action.text"
          :class="getActionClass(action)"
          @click="handleAction(action)"
        >
          <i v-if="action.icon" :class="getIconClass(action.icon)" class="mr-1"></i>
          {{ action.text }}
        </button>
      </div>
    </div>
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
    const mdiToFa = {
      'mdi-alert-circle': 'fa-solid fa-circle-exclamation',
      'mdi-check-circle': 'fa-solid fa-circle-check',
      'mdi-alert-triangle': 'fa-solid fa-triangle-exclamation',
      'mdi-information': 'fa-solid fa-circle-info',
      'mdi-home': 'fa-solid fa-house',
      'mdi-arrow-left': 'fa-solid fa-arrow-left',
      'mdi-refresh': 'fa-solid fa-rotate-right',
      'mdi-close': 'fa-solid fa-xmark',
      'mdi-chevron-left': 'fa-solid fa-chevron-left'
    }

    const getIconClass = (icon) => {
      if (!icon) return 'fa-solid fa-circle-exclamation'
      if (icon.startsWith('fa-')) return icon
      return mdiToFa[icon] || 'fa-solid fa-circle-exclamation'
    }

    const getIconColor = (color) => {
      const colors = {
        error: '#FF5722',
        warning: '#FAAD14',
        success: '#52C41A',
        primary: '#1E9FFF',
        info: '#1E9FFF'
      }
      return colors[color] || '#FF5722'
    }

    const getActionClass = (action) => {
      const color = action.color || 'primary'
      const variant = action.variant || 'flat'
      
      if (variant === 'outlined') {
        return `layui-btn layui-btn-outline`
      }
      if (color === 'error') {
        return 'layui-btn layui-btn-danger'
      }
      if (color === 'success') {
        return 'layui-btn layui-btn-normal'
      }
      return 'layui-btn'
    }

    const handleAction = (action) => {
      if (action.callback) {
        action.callback()
      }
    }
    
    return {
      handleAction,
      getIconClass,
      getIconColor,
      getActionClass
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

.error-title {
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 16px;
  color: #333;
}

.error-content {
  margin-bottom: 24px;
}

.text-body-1 {
  font-size: 16px;
  color: #333;
}

.text-body-2 {
  font-size: 14px;
}

.text-secondary {
  color: #999;
}

.mb-2 {
  margin-bottom: 8px;
}

.mt-4 {
  margin-top: 16px;
}

.d-flex {
  display: flex;
}

.justify-center {
  justify-content: center;
}

.gap-3 {
  gap: 12px;
}

.mr-1 {
  margin-right: 4px;
}
</style>