<template>
  <div ref="editorRef" class="vditor-container">
    <UploadModal
      v-model:show="showImageModal"
      upload-type="image"
      @upload-success="handleUploadSuccess"
    />
    <UploadModal
      v-model:show="showVideoModal"
      upload-type="video"
      @upload-success="handleUploadSuccess"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, watch, nextTick } from 'vue'
import Vditor from 'vditor'
import 'vditor/dist/index.css'
import UploadModal from './UploadModal.vue'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  height: {
    type: Number,
    default: 500
  },
  mode: {
    type: String,
    default: 'sv'
  }
})

const emit = defineEmits(['update:modelValue'])

const editorRef = ref(null)
let vditor = null
let isInitialized = false
const showImageModal = ref(false)
const showVideoModal = ref(false)

const handleUploadSuccess = (data) => {
  if (!vditor) return
  
  data.files.forEach(file => {
    if (data.type === 'image') {
      vditor.insertValue(`![${file.name}](${file.url})`)
    } else if (data.type === 'video') {
      vditor.insertValue(`<video src="${file.url}" controls width="100%"></video>`)
    }
  })
}

const initVditor = async () => {
  if (!editorRef.value) return

  await nextTick()

  if (vditor) {
    vditor.destroy()
    vditor = null
  }

  vditor = new Vditor(editorRef.value, {
    value: (typeof props.modelValue === 'string' ? props.modelValue : '') || '',
    height: props.height,
    mode: props.mode,
    placeholder: '使用 Markdown 语法编写文章...',
    cache: {
      enable: false
    },
    typewriterMode: true,
    toolbar: [
      {
        name: 'headings',
        tip: '标题',
        icon: '<svg viewBox="0 0 24 24" width="18" height="18"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8l-6-6zm2 18H6V4h6v6h6v12zM10 9H8v2h2V9zm0 4H8v2h2v-2zm4-4h-2v2h2V9zm0 4h-2v2h2v-2z" fill="currentColor"/></svg>'
      },
      {
        name: 'bold',
        tip: '粗体',
        icon: '<svg viewBox="0 0 24 24" width="18" height="18"><path d="M15.6 10.79c.97-.67 1.65-1.77 1.65-2.99 0-2.26-1.85-4.1-4.14-4.1-2.68 0-4.64 2.07-4.64 4.65 0 2.11 1.2 3.59 2.98 4.28.33.12.45.51.32.76-.1.18-.29.29-.49.29-.27 0-.52-.17-.61-.41-.65-1.56-.99-3.28-.99-5.21 0-3.22 2.13-5.66 5-5.66 3.03 0 5.26 2.56 5.26 5.93 0 3.45-2.19 5.87-5.36 6.32v.5H13v-.5c-.75-.13-1.47-.44-2.12-.88l-.01-.01c-.27-.18-.4-.48-.4-.82 0-.66.56-1.2 1.23-1.2.45 0 .84.27 1.05.68.28.58.42 1.2.42 1.85 0 2.32-1.48 4.41-3.8 4.41-2.52 0-4.5-2.01-4.5-4.53 0-2.39 1.79-4.33 4.18-4.33 1.95 0 3.53 1.48 3.53 3.39 0 1.05-.6 2.01-1.5 2.51-.63.35-1.3.52-2.01.52-.86 0-1.82-.48-2.39-1.27-.16-.24-.23-.51-.18-.79.05-.27.27-.48.54-.48.28 0 .53.19.65.45.13.28.38.5.68.67.88.58 1.97.92 3.16.92 3.37 0 6-2.91 6-6.48 0-3.37-2.75-6.11-6.13-6.11-3.86 0-6.98 3.2-6.98 7.15 0 4.19 2.98 7.69 7.03 7.69 1.76 0 3.37-.63 4.61-1.77l.01-.01c.21-.17.38-.38.5-.62.13-.25.19-.52.19-.8z" fill="currentColor"/></svg>'
      },
      {
        name: 'italic',
        tip: '斜体',
        icon: '<svg viewBox="0 0 24 24" width="18" height="18"><path d="M10 4v3h2.21l-3.42 8H6v3h8v-3h-2.21l3.42-8H18V4z" fill="currentColor"/></svg>'
      },
      {
        name: 'strike',
        tip: '删除线',
        icon: '<svg viewBox="0 0 24 24" width="18" height="18"><path d="M19 4H5a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6a2 2 0 0 0-2-2zm0 16H5V6h14v14zM7 10h5v2H7v-2zm5 4h5v2h-5v-2z" fill="currentColor"/></svg>'
      },
      '|',
      {
        name: 'line',
        tip: '分隔线',
        icon: '<svg viewBox="0 0 24 24" width="18" height="18"><path d="M3 12h18M3 8h18M3 16h18" stroke="currentColor" stroke-width="2" fill="none"/></svg>'
      },
      {
        name: 'quote',
        tip: '引用',
        icon: '<svg viewBox="0 0 24 24" width="18" height="18"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8l-6-6zm2 16H6V4h6v6h6v12zM10 9H8v2h2V9zm0 4H8v2h2v-2zm4-4h-2v2h2V9zm0 4h-2v2h2v-2z" fill="currentColor"/></svg>'
      },
      {
        name: 'list',
        tip: '无序列表',
        icon: '<svg viewBox="0 0 24 24" width="18" height="18"><path d="M4 4h16v2H4zm0 4h16v2H4zm0 4h16v2H4zm0 4h16v2H4z" fill="currentColor"/></svg>'
      },
      {
        name: 'ordered-list',
        tip: '有序列表',
        icon: '<svg viewBox="0 0 24 24" width="18" height="18"><path d="M2 17h2v.5H3v1h1v.5H2v1h3v-4H2v1zm1-9h1V4H2v1h1v3zm-1 3h1.8L2 13.1v.9h3v-1H3.2L5 10.9V10H2v1zm5-6v2h14V5H7zm0 4v2h14V9H7zm0 4v2h14v-2H7z" fill="currentColor"/></svg>'
      },
      {
        name: 'check',
        tip: '任务列表',
        icon: '<svg viewBox="0 0 24 24" width="18" height="18"><path d="M19 3H5a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V5a2 2 0 0 0-2-2zm0 16H5V5h14v14z" fill="currentColor"/><path d="M7 12l2 2 4-4M17 17l1.41-1.41L12 11.17l-3.41 3.42L7 14l5-5 7 7z" fill="currentColor"/></svg>'
      },
      '|',
      {
        name: 'code',
        tip: '代码块',
        icon: '<svg viewBox="0 0 24 24" width="18" height="18"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8l-6-6zm2 16H6V4h6v6h6v12zM13 9H8v2h5V9zm0 4H8v2h5v-2z" fill="currentColor"/></svg>'
      },
      {
        name: 'inline-code',
        tip: '行内代码',
        icon: '<svg viewBox="0 0 24 24" width="18" height="18"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z" fill="currentColor"/></svg>'
      },
      {
        name: 'link',
        tip: '链接',
        icon: '<svg viewBox="0 0 24 24" width="18" height="18"><path d="M3.9 12c0-1.71 1.39-3.1 3.1-3.1h4V7H7c-2.76 0-5 2.24-5 5s2.24 5 5 5h4v-1.9H7c-1.71 0-3.1-1.39-3.1-3.1zM8 13h8v-2H8v2zm9-6h-4v1.9h4c1.71 0 3.1 1.39 3.1 3.1s-1.39 3.1-3.1 3.1h-4V17h4c2.76 0 5-2.24 5-5s-2.24-5-5-5z" fill="currentColor"/></svg>'
      },
      {
        name: 'table',
        tip: '表格',
        icon: '<svg viewBox="0 0 24 24" width="18" height="18"><path d="M19 3H5a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V5a2 2 0 0 0-2-2zm0 16H5V5h14v14zM7 7h2v2H7V7zm0 4h2v2H7v-2zm0 4h2v2H7v-2zm4-8h2v2h-2V7zm0 4h2v2h-2v-2zm0 4h2v2h-2v-2zm4-8h2v2h-2V7zm0 4h2v2h-2v-2z" fill="currentColor"/></svg>'
      },
      '|',
      {
        name: 'image-upload',
        tip: '上传图片',
        icon: '<svg viewBox="0 0 24 24" width="18" height="18"><path d="M19 9h-4V3H9v6H5l7 7 7-7zM5 18v2h14v-2H5z" fill="currentColor"/></svg>',
        click: () => {
          showImageModal.value = true
        }
      },
      {
        name: 'video-upload',
        tip: '上传视频',
        icon: '<svg viewBox="0 0 24 24" width="18" height="18"><path d="M8 5v14l11-7z" fill="currentColor"/></svg>',
        click: () => {
          showVideoModal.value = true
        }
      },
      '|',
      {
        name: 'undo',
        tip: '撤销',
        icon: '<svg viewBox="0 0 24 24" width="18" height="18"><path d="M12.5 8c-2.65 0-5.05.99-6.9 2.6L2 7v9h9l-3.62-3.62c1.39-1.16 3.16-1.88 5.12-1.88 3.42 0 6.17 2.46 7.06 5.65.25.84-.51 1.61-1.35 1.61H4.83c-.52 0-.95-.39-.95-.89V17c0-.28.11-.53.29-.71L9.5 13H4c-.55 0-1-.45-1-1s.45-1 1-1h5.5z" fill="currentColor"/></svg>'
      },
      {
        name: 'redo',
        tip: '重做',
        icon: '<svg viewBox="0 0 24 24" width="18" height="18"><path d="M18.4 10.6C16.55 8.99 14.15 8 11.5 8c-4.67 0-8.5 3.83-8.5 8.5s3.83 8.5 8.5 8.5c3.34 0 6.14-2.11 7.35-5.14.15-.41-.17-.82-.58-.82h-4.8c-.27 0-.52-.17-.61-.42l-.81-2.31c-.14-.42.12-.85.54-.85h1.71c.43 0 .75-.34.75-.78 0-.63-.62-1.02-1.19-.72l-.66.25c-.98.36-2.06.57-3.2.57-2.65 0-5.05-.99-6.9-2.6L2 17V8h9l3.62 3.62c-.13.11-.26.23-.37.38zm-1.5-1.9c.94-.98 1.5-2.23 1.5-3.63 0-3.31-2.69-6-6-6S5 4.69 5 8s2.69 6 6 6c1.44 0 2.77-.58 3.75-1.53l2.25 2.25V8l-2.35 2.35z" fill="currentColor"/></svg>'
      },
      '|',
      {
        name: 'preview',
        tip: '预览',
        icon: '<svg viewBox="0 0 24 24" width="18" height="18"><path d="M19 3H5a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V5a2 2 0 0 0-2-2zm-4 8l-4 4-4-4 1.41-1.41L11 10.17l6.59-6.59L19 5l-8 8z" fill="currentColor"/></svg>'
      }
    ],
    toolbarConfig: {
      hide: false,
      pin: true
    },
    input (value) {
      emit('update:modelValue', value)
    },
    after () {
      isInitialized = true
    }
  })
}

watch(() => props.modelValue, (newVal) => {
  if (isInitialized && vditor) {
    const strVal = typeof newVal === 'string' ? newVal : ''
    if (vditor.getValue() !== strVal) {
      vditor.setValue(strVal)
    }
  }
})

onMounted(() => {
  initVditor()
})

onBeforeUnmount(() => {
  if (vditor) {
    vditor.destroy()
    vditor = null
    isInitialized = false
  }
})

defineExpose({
  getInstance: () => vditor,
  focus: () => vditor?.focus()
})
</script>

<style scoped>
.vditor-container {
  border: 1px solid rgba(0, 0, 0, 0.12);
  border-radius: 8px;
  min-height: 500px;
  overflow: visible;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.vditor-container:focus-within {
  border-color: rgb(var(--v-theme-primary));
  box-shadow: 0 0 0 2px rgba(var(--v-theme-primary), 0.1);
}

:deep(.vditor) {
  border: none;
  border-radius: 8px;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
}

:deep(.vditor-toolbar) {
  border-bottom: 1px solid rgba(0, 0, 0, 0.08);
  padding: 8px 12px;
  background: #fafbfc;
  position: sticky;
  top: 0;
  z-index: 10;
  overflow: visible;
}

:deep(.vditor-toolbar__item) {
  margin: 0 4px;
  padding: 6px 8px;
  border-radius: 4px;
  transition: background-color 0.2s;
  position: relative;
  cursor: pointer;
}

:deep(.vditor-toolbar__item:hover) {
  background-color: rgba(0, 0, 0, 0.06);
}

:deep(.vditor-toolbar__item:hover::after) {
  content: attr(data-tip);
  position: absolute;
  bottom: 100%;
  left: 50%;
  transform: translateX(-50%);
  padding: 6px 10px;
  background-color: rgba(0, 0, 0, 0.85);
  color: #fff;
  font-size: 12px;
  white-space: nowrap;
  border-radius: 4px;
  opacity: 1;
  visibility: visible;
  transition: opacity 0.2s, visibility 0.2s;
  z-index: 9999;
  pointer-events: none;
  margin-bottom: 8px;
  display: block;
}

:deep(.vditor-toolbar__icon) {
  width: 20px;
  height: 20px;
}

:deep(.vditor-panel) {
  min-width: 200px;
  max-width: 400px;
}

:deep(.vditor-sv) {
  font-size: 15px;
}

:deep(.vditor-sv .vditor-reset) {
  padding: 16px 20px;
}

:deep(.vditor-sv .vditor-edit) {
  padding: 16px 20px;
}

:deep(.vditor-preview) {
  padding: 16px 20px;
  background: #fff;
}

:deep(.vditor-content) {
  min-height: 400px;
}

:deep(.vditor-edit textarea) {
  font-family: 'SF Mono', Consolas, 'Liberation Mono', Menlo, monospace;
  line-height: 1.6;
}

:deep(.vditor-preview__action) {
  display: none !important;
}

@media (max-width: 600px) {
  :deep(.vditor-toolbar) {
    overflow-x: auto;
    overflow-y: hidden;
    white-space: nowrap;
  }
  
  :deep(.vditor-toolbar__item) {
    flex-shrink: 0;
  }
}
</style>
