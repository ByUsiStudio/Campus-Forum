import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import '@mdi/font/css/materialdesignicons.css'

export default createVuetify({
  components,
  directives,
  theme: {
    defaultTheme: 'lightPurple',
    themes: {
      lightPurple: {
        dark: false,
        colors: {
          primary: '#9575CD',
          'primary-lighten-1': '#B39DDB',
          'primary-lighten-2': '#D1C4E9',
          'primary-darken-1': '#7E57C2',
          secondary: '#CE93D8',
          accent: '#F3E5F5',
          background: '#FAFAFA',
          surface: '#FFFFFF',
          error: '#EF5350',
          info: '#42A5F5',
          success: '#66BB6A',
          warning: '#FFA726',
          // Text colors
          'on-primary': '#FFFFFF',
          'on-secondary': '#FFFFFF',
          'on-surface': '#424242',
          'on-background': '#424242',
        }
      },
      dark: {
        colors: {
          primary: '#B39DDB',
          secondary: '#CE93D8',
          accent: '#1E1E1E',
          error: '#EF5350',
          info: '#42A5F5',
          success: '#66BB6A',
          warning: '#FFA726',
        }
      }
    }
  }
})
