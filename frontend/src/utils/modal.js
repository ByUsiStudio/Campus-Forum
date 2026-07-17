import { ref } from 'vue'

const modalState = ref({
  show: false,
  type: 'alert',
  title: '',
  message: '',
  icon: 'fa-solid fa-circle-exclamation',
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

let confirmCallback = null
let cancelCallback = null

export const alert = (message, options = {}) => {
  return new Promise((resolve) => {
    modalState.value = {
      show: true,
      type: 'alert',
      title: options.title || '提示',
      message: message,
      icon: options.icon || 'fa-solid fa-circle-info',
      iconColor: options.iconColor || 'primary',
      confirmText: options.confirmText || '确定',
      confirmColor: options.confirmColor || 'primary'
    }
    confirmCallback = resolve
    cancelCallback = null
  })
}

export const confirm = (message, options = {}) => {
  return new Promise((resolve) => {
    modalState.value = {
      show: true,
      type: 'confirm',
      title: options.title || '确认',
      message: message,
      icon: options.icon || 'fa-solid fa-circle-exclamation',
      iconColor: options.iconColor || 'warning',
      confirmText: options.confirmText || '确定',
      cancelText: options.cancelText || '取消',
      confirmColor: options.confirmColor || 'error'
    }
    confirmCallback = () => resolve(true)
    cancelCallback = () => resolve(false)
  })
}

export const prompt = (message, options = {}) => {
  return new Promise((resolve) => {
    modalState.value = {
      show: true,
      type: 'prompt',
      title: options.title || '输入',
      message: message,
      icon: options.icon || 'fa-solid fa-pencil',
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
    cancelCallback = () => resolve(null)
  })
}

export const success = (message, options = {}) => {
  return alert(message, {
    title: options.title || '成功',
    icon: 'fa-solid fa-circle-check',
    iconColor: 'success',
    ...options
  })
}

export const error = (message, options = {}) => {
  return alert(message, {
    title: options.title || '错误',
    icon: 'fa-solid fa-circle-exclamation',
    iconColor: 'error',
    ...options
  })
}

export const warning = (message, options = {}) => {
  return alert(message, {
    title: options.title || '警告',
    icon: 'fa-solid fa-triangle-exclamation',
    iconColor: 'warning',
    ...options
  })
}

export const handleConfirm = (value) => {
  if (confirmCallback) {
    confirmCallback(value)
    confirmCallback = null
  }
}

export const handleCancel = () => {
  if (cancelCallback) {
    cancelCallback()
    cancelCallback = null
  }
}

export { modalState }