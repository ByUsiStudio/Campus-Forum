/**
 * 统一交互反馈服务（基于 JCuPupw）
 *
 * 所有弹窗 / 轻提示统一走 JCuPupw 库（window.JCuPupw）：
 *   - toast  : JCuPupw.toast({ content, type: 'success'|'error'|'warning'|'info', duration })
 *   - alert  : JCuPupw.alert({ title, content, type, buttonText })
 *   - confirm: JCuPupw.confirm({ title, content, confirmText, cancelText }) -> Promise<boolean>
 *   - prompt : JCuPupw.prompt({ title, content, placeholder, defaultValue, ... }) -> Promise<string|null>
 *
 * 对外保持与原 Element Plus 版一致的签名，方便整个前端平滑切换。
 */

function jc() {
  if (typeof window !== 'undefined' && window.JCuPupw) return window.JCuPupw
  throw new Error('JCuPupw 组件库未加载，请检查 index.html 是否引入 libs/js/jcupupw.umd.js')
}

export function success(message, title = '成功') {
  try {
    jc().toast({ content: String(message), type: 'success' })
  } catch (e) { /* ignore */ }
  return Promise.resolve(true)
}

export function error(message, title = '错误') {
  try {
    jc().toast({ content: String(message), type: 'error' })
  } catch (e) { /* ignore */ }
  return Promise.resolve(false)
}

export function warning(message, title = '警告') {
  try {
    jc().toast({ content: String(message), type: 'warning' })
  } catch (e) { /* ignore */ }
  return Promise.resolve(true)
}

export function info(message, title = '提示') {
  try {
    jc().toast({ content: String(message), type: 'info' })
  } catch (e) { /* ignore */ }
  return Promise.resolve(true)
}

/**
 * 轻提示（别名，面向需要细粒度控制的场景）
 * options: { type: 'success'|'error'|'warning'|'info', duration }
 */
export function toast(message, options = {}) {
  const { type = 'info', duration } = options
  try {
    return jc().toast({ content: String(message), type, duration })
  } catch (e) { /* ignore */ }
  return null
}

export function alert(message, options = {}) {
  try {
    return jc().alert({
      title: options.title || '提示',
      content: String(message),
      type: 'primary',
      buttonText: options.confirmText || '确定'
    })
  } catch (e) {
    return Promise.resolve()
  }
}

export function confirm(message, options = {}) {
  try {
    return jc().confirm({
      title: options.title || '确认',
      content: String(message),
      confirmText: options.confirmText || '确定',
      cancelText: options.cancelText || '取消'
    })
  } catch (e) {
    return Promise.resolve(false)
  }
}

export function prompt(message, options = {}) {
  try {
    return jc().prompt({
      title: options.title || '输入',
      content: String(message),
      placeholder: options.placeholder || '',
      defaultValue: options.defaultValue || '',
      confirmText: options.confirmText || '确定',
      cancelText: options.cancelText || '取消'
    })
  } catch (e) {
    return Promise.resolve(null)
  }
}

// JCuPupw 全局引用便捷导出（库本体挂载于 window.JCuPupw）
export { jc as JCuPupw }

/** 提供 window.JCuPupw 便捷引用 */
export function getJCuPupw() {
  return (typeof window !== 'undefined' && window.JCuPupw) ? window.JCuPupw : null
}
