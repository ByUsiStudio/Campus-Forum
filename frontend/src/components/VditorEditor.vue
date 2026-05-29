<template>
  <div ref="editorRef" class="vditor-container"></div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, watch, nextTick } from 'vue'
import Vditor from 'vditor'
import 'vditor/dist/index.css'

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

const initVditor = async () => {
  if (!editorRef.value) return

  await nextTick()

  if (vditor) {
    vditor.destroy()
    vditor = null
  }

  vditor = new Vditor(editorRef.value, {
    value: props.modelValue || '',
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
        tip: '标题'
      },
      {
        name: 'bold',
        tip: '加粗'
      },
      {
        name: 'italic',
        tip: '斜体'
      },
      {
        name: 'strike',
        tip: '删除线'
      },
      '|',
      {
        name: 'line',
        tip: '分割线'
      },
      {
        name: 'quote',
        tip: '引用'
      },
      {
        name: 'list',
        tip: '无序列表'
      },
      {
        name: 'ordered-list',
        tip: '有序列表'
      },
      {
        name: 'check',
        tip: '任务列表'
      },
      '|',
      {
        name: 'code',
        tip: '代码块'
      },
      {
        name: 'inline-code',
        tip: '行内代码'
      },
      {
        name: 'link',
        tip: '链接'
      },
      {
        name: 'table',
        tip: '表格'
      },
      '|',
      {
        name: 'upload',
        tip: '上传图片'
      },
      '|',
      {
        name: 'undo',
        tip: '撤销'
      },
      {
        name: 'redo',
        tip: '重做'
      },
      '|',
      {
        name: 'preview',
        tip: '预览'
      }
    ],
    toolbarConfig: {
      hide: false,
      pin: true
    },
    upload: {
      accept: 'image/*',
      url: '/api/upload/image',
      fieldName: 'image',
      headers: {
        'Authorization': localStorage.getItem('token') || ''
      },
      success (editor, msg) {
        try {
          const res = JSON.parse(msg)
          if (res.url) {
            vditor.insertValue(`![图片](${res.url})`)
          }
        } catch (e) {
          console.error('Parse upload response failed:', e)
        }
      },
      error (msg) {
        console.error('Upload failed:', msg)
      }
    },
    input (value) {
      emit('update:modelValue', value)
    },
    after () {
      isInitialized = true
      if (props.modelValue) {
        vditor.setValue(props.modelValue)
      }
    }
  })
}

watch(() => props.modelValue, (newVal) => {
  if (isInitialized && vditor && vditor.getValue() !== newVal) {
    vditor.setValue(newVal || '')
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
  overflow: hidden;
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
}

:deep(.vditor-toolbar__item) {
  margin: 0 4px;
  padding: 6px 8px;
  border-radius: 4px;
  transition: background-color 0.2s;
  position: relative;
}

:deep(.vditor-toolbar__item:hover) {
  background-color: rgba(0, 0, 0, 0.06);
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

:deep(.vditor-tooltipped) {
  position: relative;
}

:deep(.vditor-tooltipped::after) {
  content: attr(data-tooltip);
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
  opacity: 0;
  visibility: hidden;
  transition: opacity 0.2s, visibility 0.2s;
  z-index: 1000;
  pointer-events: none;
  margin-bottom: 6px;
}

:deep(.vditor-tooltipped:hover::after) {
  opacity: 1;
  visibility: visible;
}

:deep(.vditor-preview__action) {
  display: none !important;
}
</style>