// 全局弹窗服务
import { ref } from 'vue'

// 弹窗状态
const modalState = ref({
  show: false,
  type: 'alert',
  title: '',
  message: '',
  icon: 'mdi-alert-circle',
  iconColor: 'warning',
  confirmText: '确定',
  cancelText: '取消',
  confirmColor: 'primary',
  inputValue: '',
  inputLabel: '',
  inputType: 'text',
  inputPlaceholder: '',
  inputRows: 1
})

// 回调函数
let confirmCallback = null
let cancelCallback = null

// Alert 弹窗
export const alert = (message, options = {}) => {
  return new Promise((resolve) => {
    modalState.value = {
      show: true,
      type: 'alert',
      title: options.title || '提示',
      message: message,
      icon: options.icon || 'mdi-information',
      iconColor: options.iconColor || 'primary',
      confirmText: options.confirmText || '确定',
      confirmColor: options.confirmColor || 'primary'
    }
    confirmCallback = resolve
    cancelCallback = null
  })
}

// Confirm 弹窗
export const confirm = (message, options = {}) => {
  return new Promise((resolve, reject) => {
    modalState.value = {
      show: true,
      type: 'confirm',
      title: options.title || '确认',
      message: message,
      icon: options.icon || 'mdi-alert-circle',
      iconColor: options.iconColor || 'warning',
      confirmText: options.confirmText || '确定',
      cancelText: options.cancelText || '取消',
      confirmColor: options.confirmColor || 'error'
    }
    confirmCallback = resolve
    cancelCallback = reject
  })
}

// Prompt 弹窗
export const prompt = (message, options = {}) => {
  return new Promise((resolve, reject) => {
    modalState.value = {
      show: true,
      type: 'prompt',
      title: options.title || '输入',
      message: message,
      icon: options.icon || 'mdi-edit',
      iconColor: options.iconColor || 'primary',
      confirmText: options.confirmText || '确定',
      cancelText: options.cancelText || '取消',
      confirmColor: options.confirmColor || 'primary',
      inputValue: options.defaultValue || '',
      inputLabel: options.inputLabel || '',
      inputType: options.inputType || 'text',
      inputPlaceholder: options.placeholder || '',
      inputRows: options.rows || 1
    }
    confirmCallback = resolve
    cancelCallback = reject
  })
}

// 成功提示
export const success = (message, options = {}) => {
  return alert(message, {
    title: options.title || '成功',
    icon: 'mdi-check-circle',
    iconColor: 'success',
    ...options
  })
}

// 错误提示
export const error = (message, options = {}) => {
  return alert(message, {
    title: options.title || '错误',
    icon: 'mdi-alert-circle',
    iconColor: 'error',
    ...options
  })
}

// 警告提示
export const warning = (message, options = {}) => {
  return alert(message, {
    title: options.title || '警告',
    icon: 'mdi-alert-triangle',
    iconColor: 'warning',
    ...options
  })
}

// 处理确认
export const handleConfirm = (value) => {
  if (confirmCallback) {
    confirmCallback(value)
    confirmCallback = null
  }
}

// 处理取消
export const handleCancel = () => {
  if (cancelCallback) {
    cancelCallback()
    cancelCallback = null
  }
}

export { modalState }
