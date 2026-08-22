import { ElMessage, ElMessageBox } from 'element-plus'

/**
 * 统一交互反馈服务（基于 Element Plus）
 * confirm / prompt / alert / success / error / warning
 */

export function success(message, title = '成功') {
  ElMessage({ type: 'success', message })
  return Promise.resolve(true)
}

export function error(message, title = '错误') {
  ElMessage({ type: 'error', message })
  return Promise.resolve(false)
}

export function warning(message, title = '警告') {
  ElMessage({ type: 'warning', message })
  return Promise.resolve(true)
}

export function alert(message, options = {}) {
  return ElMessageBox.alert(message, options.title || '提示', {
    confirmButtonText: options.confirmText || '确定',
    type: options.type || 'info'
  })
}

export function confirm(message, options = {}) {
  return ElMessageBox.confirm(message, options.title || '确认', {
    confirmButtonText: options.confirmText || '确定',
    cancelButtonText: options.cancelText || '取消',
    type: options.type || 'warning',
    confirmButtonClass: options.confirmButtonClass || ''
  })
}

export function prompt(message, options = {}) {
  return ElMessageBox.prompt(message, options.title || '输入', {
    confirmButtonText: options.confirmText || '确定',
    cancelButtonText: options.cancelText || '取消',
    inputValue: options.defaultValue || '',
    inputPlaceholder: options.placeholder || '',
    inputType: options.inputType || 'text',
    inputValidator: () => true
  }).then(({ value }) => value)
}

export { ElMessage }
