<template>
  <!--
    AppModal 已基于 JCuPupw 实现。
    保持与原 Element Plus 版完全一致的 props/emits 契约：
    props: show/type/title/message/icon/iconColor/confirmText/cancelText/confirmColor
           inputValue/inputLabel/inputType/inputPlaceholder/inputRows
    emits: update:show/confirm/cancel
  -->
  <div class="app-modal-slot"></div>
</template>

<script>
import { ref, computed, watch } from 'vue'
import { jcOpenHtml, jcCloseAll } from '@/utils/jcu'

function esc(s) {
  return String(s == null ? '' : s)
    .replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
}

const ICON_GLYPH = {
  'mdi-alert-circle': '⚠️', 'warning': '⚠️', 'alert': '⚠️',
  'mdi-information': 'ℹ️', 'info': 'ℹ️',
  'mdi-check-circle': '✅', 'success': '✅',
  'mdi-edit': '✏️', 'edit': '✏️'
}

export default {
  name: 'AppModal',
  props: {
    show: { type: Boolean, default: false },
    type: { type: String, default: 'alert', validator: (v) => ['alert', 'confirm', 'prompt'].includes(v) },
    title: { type: String, default: '' },
    message: { type: String, default: '' },
    icon: { type: String, default: 'mdi-alert-circle' },
    iconColor: { type: String, default: 'warning' },
    confirmText: { type: String, default: '确定' },
    cancelText: { type: String, default: '取消' },
    confirmColor: { type: String, default: 'primary' },
    inputValue: { type: String, default: '' },
    inputLabel: { type: String, default: '' },
    inputType: { type: String, default: 'text' },
    inputPlaceholder: { type: String, default: '' },
    inputRows: { type: Number, default: 1 }
  },
  emits: ['update:show', 'confirm', 'cancel'],
  setup(props, { emit }) {
    const internalValue = ref(props.inputValue)
    let contentEl = null

    const showCancel = computed(() => props.type === 'confirm' || props.type === 'prompt')

    const confirmBtnType = computed(() => {
      const c = props.confirmColor
      if (c === 'error' || c === 'danger') return 'danger'
      if (c === 'success') return 'primary'
      if (c === 'warning') return 'danger'
      if (c === 'secondary') return 'default'
      return 'primary'
    })

    watch(() => props.show, (val) => { if (val) open() })
    watch(() => props.inputValue, (val) => { internalValue.value = val })

    const resetInternal = () => { internalValue.value = props.inputValue }

    const open = () => {
      contentEl = null
      const glyph = ICON_GLYPH[props.icon] || (props.iconColor === 'success' ? '✅' : '⚠️')
      let content = ''

      if (props.type === 'prompt') {
        const labelHtml = props.inputLabel
          ? `<div style="margin-bottom:6px;font-weight:600;font-size:13px;color:var(--jc-text,#333);">${esc(props.inputLabel)}</div>`
          : ''
        const isArea = props.inputType === 'textarea'
        content = `
          <div style="margin-bottom:14px;">${glyph} <span>${esc(props.message)}</span></div>
          ${labelHtml}
          ${isArea
            ? `<textarea data-app-input class="jc-modal__input" rows="${props.inputRows || 1}" style="resize:vertical;" placeholder="${esc(props.inputPlaceholder || '')}">${esc(props.inputValue)}</textarea>`
            : `<input data-app-input class="jc-modal__input" type="text" value="${esc(props.inputValue)}" placeholder="${esc(props.inputPlaceholder || '')}" />`}
        `
      } else {
        content = `<div style="display:flex;align-items:flex-start;gap:10px;"><span style="font-size:22px;line-height:1.4;">${glyph}</span><span>${esc(props.message)}</span></div>`
      }

      const buttons = []
      if (showCancel.value) {
        buttons.push({
          text: props.cancelText,
          type: 'default',
          action: () => {
            emit('cancel')
            emit('update:show', false)
            jcCloseAll()
          }
        })
      }
      buttons.push({
        text: props.confirmText,
        type: confirmBtnType.value,
        action: () => {
          let val = props.inputValue
          if (props.type === 'prompt' && contentEl) {
            const input = contentEl.querySelector('[data-app-input]')
            if (input) val = input.value
          }
          internalValue.value = val
          emit('confirm', val)
          emit('update:show', false)
          jcCloseAll()
        }
      })

      jcOpenHtml({
        title: props.title || '提示',
        content,
        width: 500,
        size: 'sm',
        closeOnOverlay: false,
        onMount: (el) => { contentEl = el },
        onClose: () => { resetInternal() },
        buttons
      })
    }

    return { open, internalValue, showCancel, confirmBtnType, resetInternal }
  }
}
</script>

<style scoped>
.app-modal-slot { display: none; }
</style>
