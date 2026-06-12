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
          info: '#3B82F6',         // 信息蓝色
          success: '#22C55E',      // 成功绿色
          warning: '#F59E0B',      // 警告黄色
          'primary-light': '#E6D5F0',
          'primary-dark': '#7B68EE',
        }
      },
      dark: {
        colors: {
          primary: '#A855F7',      // 深紫色主色
          secondary: '#C084FC',    // 中等紫色
          accent: '#D8B4FE',       // 浅紫色
          background: '#1A1A2E',   // 深色背景
          surface: '#252542',      // 表面深色
          error: '#F87171',        // 错误红色
          info: '#60A5FA',         // 信息蓝色
          success: '#4ADE80',      // 成功绿色
          warning: '#FBBF24',      // 警告黄色
        }
      }
    }
  }
})