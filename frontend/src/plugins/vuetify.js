import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import '@mdi/font/css/materialdesignicons.css'

export default createVuetify({
  components,
  directives,
  theme: {
    defaultTheme: 'light',
    themes: {
      light: {
        colors: {
          primary: '#9370DB',      // 淡紫色主色
          secondary: '#DDA0DD',    // 浅紫色
          accent: '#BA55D3',       // 中等紫色
          background: '#FDF5FF',   // 淡紫色背景
          surface: '#FFFFFF',      // 表面白色
          error: '#EF4444',        // 错误红色
          info: '#6366F1',         // 信息靛蓝色
          success: '#10B981',      // 成功绿色
          warning: '#F59E0B',      // 警告黄色
          'primary-light': '#E6D5F0',
          'primary-dark': '#7B68EE',
          'primary-soft': '#F5E6FF',
          'purple-50': '#FAF5FF',
          'purple-100': '#F3E8FF',
          'purple-200': '#E9D5FF',
          'purple-300': '#D8B4FE',
          'purple-400': '#C084FC',
          'purple-500': '#A855F7',
          'purple-600': '#9333EA',
          'purple-700': '#7E22CE',
          'gray-50': '#F9FAFB',
          'gray-100': '#F3F4F6',
          'gray-200': '#E5E7EB',
          'gray-300': '#D1D5DB',
          'gray-400': '#9CA3AF',
          'gray-500': '#6B7280',
          'gray-600': '#4B5563',
          'gray-700': '#374151',
          'gray-800': '#1F2937',
          'gray-900': '#111827',
        },
        variables: {
          'border-radius': '12px',
          'border-radius-sm': '8px',
          'border-radius-lg': '16px',
          'border-radius-xl': '24px',
          'shadow-key-umbra-opacity': '0.05',
          'shadow-key-penumbra-opacity': '0.03',
          'shadow-key-ambient-opacity': '0',
        }
      },
      dark: {
        colors: {
          primary: '#A855F7',      // 深紫色主色
          secondary: '#C084FC',    // 中等紫色
          accent: '#D8B4FE',       // 浅紫色
          background: '#0F0F23',   // 深色背景
          surface: '#1A1A2E',      // 表面深色
          error: '#F87171',        // 错误红色
          info: '#818CF8',         // 信息靛蓝色
          success: '#4ADE80',      // 成功绿色
          warning: '#FBBF24',      // 警告黄色
          'primary-light': '#C084FC',
          'primary-dark': '#7E22CE',
        },
        variables: {
          'border-radius': '12px',
          'border-radius-sm': '8px',
          'border-radius-lg': '16px',
          'border-radius-xl': '24px',
        }
      }
    }
  },
  defaults: {
    VBtn: {
      rounded: 'lg',
      class: 'btn-hover',
    },
    VCard: {
      rounded: 'xl',
      elevation: 2,
      class: 'card-hover',
    },
    VTextField: {
      rounded: 'lg',
      variant: 'outlined',
    },
    VTextarea: {
      rounded: 'lg',
      variant: 'outlined',
    },
    VSelect: {
      rounded: 'lg',
      variant: 'outlined',
    },
    VList: {
      rounded: 'lg',
    },
    VListTile: {
      class: 'list-item-hover',
    },
    VChip: {
      rounded: 'lg',
    },
    VProgressCircular: {
      color: 'primary',
    },
  }
})