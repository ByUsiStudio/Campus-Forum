/**
 * JCuPupw 组件库薄封装
 *
 * 库本体：frontend/public/libs/js/jcupupw.umd.js
 * 全局对象：window.JCuPupw（加载于 index.html）
 *
 * 兼容 Vue SFC 的便捷封装，尤其用于“复杂表单/自定义弹窗”。
 * 所有弹窗统一走本模块，保证全站一致。
 */
function jc() {
  if (typeof window !== 'undefined' && window.JCuPupw) return window.JCuPupw
  throw new Error('JCuPupw 组件库未加载，请检查 index.html 是否正确引入 libs/js/jcupupw.umd.js')
}

/** 轻提示（不阻塞） */
export function jcToast(content, type = 'info', duration = 2500) {
  try {
    return jc().toast({ content: String(content), type, duration })
  } catch (e) {
    console.error('[JCuPupw] toast 失败：', e)
    return null
  }
}

/** 确认框 -> Promise<boolean> */
export function jcConfirm({ title = '确认', content = '', confirmText = '确定', cancelText = '取消', ...rest } = {}) {
  try {
    return jc().confirm({ title, content, confirmText, cancelText, ...rest })
  } catch (e) {
    console.error('[JCuPupw] confirm 失败：', e)
    return Promise.resolve(false)
  }
}

/** 输入框 -> Promise<string | null> */
export function jcPrompt({ title = '输入', content = '', placeholder = '', defaultValue = '', confirmText = '确定', cancelText = '取消', ...rest } = {}) {
  try {
    return jc().prompt({ title, content, placeholder, defaultValue, confirmText, cancelText, ...rest })
  } catch (e) {
    console.error('[JCuPupw] prompt 失败：', e)
    return Promise.resolve(null)
  }
}

/**
 * 自定义 HTML 弹窗（复杂表单场景）
 * @param {object} opts
 *   title   - 标题
 *   content - HTML 字符串（放入 .jc-modal__content）
 *   buttons - [{ text, type: 'default'|'primary'|'danger', action }]，action 内的 this 绑定实例
 *   width   - 数字 px 或 CSS 宽度
 *   size    - 'sm' | 'md' | 'lg' | 'auto'
 *   onOpen  - (instance) => void，打开后回调（用于挂载 DOM 事件 / 读取输入）
 *   onClose - (instance) => void
 *   onMount - (contentEl) => void，DOM 就绪后回调，可在此绑定监听器
 * @returns JCuPupw instance（open 返回的 Promise）
 */
export function jcOpenHtml({ title = '', content = '', buttons = [], width, size, onOpen, onClose, onMount, ...rest } = {}) {
  const instance = jc().instance()
  const fallback = [{ text: '关闭', type: 'primary', action: () => instance.close() }]
  return instance.open({
    ...rest,
    title,
    content,
    width,
    size,
    buttons: buttons.length ? buttons : fallback,
    onOpen: (inst) => {
      if (typeof onMount === 'function') {
        const el = inst.modalContent
        if (el) onMount(el)
      }
      if (typeof onOpen === 'function') onOpen(inst)
    },
    onClose: (inst) => {
      if (typeof onClose === 'function') onClose(inst)
    }
  })
}

/** 关闭所有 JCuPupw 弹窗 */
export function jcCloseAll() {
  try {
    jc().closeAll()
  } catch (e) { /* ignore */ }
}

/** 打开复杂弹窗并返回 (instance, contentEl)，供调用方即时取用 */
export function jcOpenRaw(opts = {}) {
  const instance = jc().instance()
  let contentEl = null
  const onOpen = (inst) => {
    contentEl = inst.modalContent
    if (opts.onMount) opts.onMount(inst.modalContent, inst)
  }
  instance.open({ ...opts, onOpen }).catch(() => {})
  return {
    instance,
    get content() {
      return contentEl
    }
  }
}

/**
 * 表单弹窗（把一组字段渲染成 .jc-modal__input 风格的 HTML，并生成配置对象）
 *
 * field 配置：
 *   { name, label, type: 'text'|'password'|'textarea'|'number', placeholder, value, required, rows }
 *
 * 返回一个对象：
 *   { html, collect: (root) => valuesObject, validate: (root) => boolean }
 * collect 从已渲染的 DOM 读取各字段值；validate 做必填校验并 toast 提示。
 */
export function jcFieldsConfig(fields = []) {
  const esc = (s) => String(s ?? '').replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;')

  const html = fields.map((f) => {
    const label = f.label ? `<div style="margin:0 0 4px;font-weight:600;font-size:13px;color:var(--jc-text,#333);">${esc(f.label)}</div>` : ''
    const placeholder = f.placeholder ? ` placeholder="${esc(f.placeholder)}"` : ''
    const required = f.required ? ' required' : ''
    const value = f.value != null ? esc(f.value) : ''
    if (f.type === 'textarea') {
      return `<div style="margin-bottom:12px;">${label}<textarea data-jf="${esc(f.name)}" class="jc-modal__input" rows="${f.rows || 3}" style="resize:vertical;"${placeholder}${required}>${value}</textarea></div>`
    }
    const inputType = f.type === 'number' ? 'number' : (f.type === 'password' ? 'password' : 'text')
    return `<div style="margin-bottom:12px;">${label}<input data-jf="${esc(f.name)}" class="jc-modal__input" type="${inputType}" value="${value}"${placeholder}${required} /></div>`
  }).join('')

  const collect = (root) => {
    const values = {}
    const nodes = (root || document).querySelectorAll('[data-jf]')
    nodes.forEach((el) => {
      const name = el.getAttribute('data-jf')
      let v = el.value || ''
      if (el.getAttribute('type') === 'number') v = v === '' ? '' : Number(v)
      values[name] = v
    })
    return values
  }

  const validate = (root) => {
    const nodes = (root || document).querySelectorAll('[data-jf][required]')
    for (const el of nodes) {
      if (!(el.value || '').trim()) {
        const label = el.closest('div')
        const lbl = label ? label.querySelector('div')?.textContent : ''
        jcToast(`请填写${lbl || '必填项'}`, 'warning')
        return false
      }
    }
    return true
  }

  return { html, collect, validate }
}

export { jc }
