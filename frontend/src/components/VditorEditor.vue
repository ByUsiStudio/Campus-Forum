<template>
  <div ref="editorRef" class="vditor-container">
    <input
      ref="videoInputRef"
      type="file"
      accept="video/*"
      style="display: none"
      @change="handleVideoUpload"
    />
  </div>
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
const videoInputRef = ref(null)

const handleVideoUpload = async (event) => {
  const file = event.target.files[0]
  if (!file || !vditor) return

  const formData = new FormData()
  formData.append('video', file)

  try {
    const token = localStorage.getItem('token')
    const response = await fetch('/api/upload/video', {
      method: 'POST',
      headers: {
        'Authorization': token ? `Bearer ${token}` : ''
      },
      body: formData
    })

    const data = await response.json()
    if (data.url) {
      vditor.insertValue(`<video src="${data.url}" controls></video>`)
    } else {
      console.error('Upload failed:', data.error)
    }
  } catch (error) {
    console.error('Video upload failed:', error)
  }

  event.target.value = ''
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
        dataTooltip: '标题'
      },
      {
        name: 'bold',
        tip: '加粗',
        dataTooltip: '加粗'
      },
      {
        name: 'italic',
        tip: '斜体',
        dataTooltip: '斜体'
      },
      {
        name: 'strike',
        tip: '删除线',
        dataTooltip: '删除线'
      },
      '|',
      {
        name: 'line',
        tip: '分割线',
        dataTooltip: '分割线'
      },
      {
        name: 'quote',
        tip: '引用',
        dataTooltip: '引用'
      },
      {
        name: 'list',
        tip: '无序列表',
        dataTooltip: '无序列表'
      },
      {
        name: 'ordered-list',
        tip: '有序列表',
        dataTooltip: '有序列表'
      },
      {
        name: 'check',
        tip: '任务列表',
        dataTooltip: '任务列表'
      },
      '|',
      {
        name: 'code',
        tip: '代码块',
        dataTooltip: '代码块'
      },
      {
        name: 'inline-code',
        tip: '行内代码',
        dataTooltip: '行内代码'
      },
      {
        name: 'link',
        tip: '链接',
        dataTooltip: '链接'
      },
      {
        name: 'table',
        tip: '表格',
        dataTooltip: '表格'
      },
      '|',
      {
        name: 'upload',
        tip: '上传图片',
        dataTooltip: '上传图片'
      },
      {
        tip: '上传视频',
        dataTooltip: '上传视频',
        icon: '<svg viewBox="0 0 24 24" width="18" height="18"><path d="M17 10.5V7c0-.55-.45-1-1-1H4c-.55 0-1 .45-1 1v10c0 .55.45 1 1 1h12c.55 0 1-.45 1-1v-3.5l4 4v-11l-4 4z" fill="currentColor"></path></svg>',
        click () {
          videoInputRef.value?.click()
        }
      },
      '|',
      {
        name: 'undo',
        tip: '撤销',
        dataTooltip: '撤销'
      },
      {
        name: 'redo',
        tip: '重做',
        dataTooltip: '重做'
      },
      '|',
      {
        name: 'preview',
        tip: '预览',
        dataTooltip: '预览'
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
            vditor.insertValue(`![${res.url}](${res.url})`)
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

:deep(.vditor-toolbar__item[data-tip]) {
  position: relative;
}

:deep(.vditor-toolbar__item[data-tip]:hover::after) {
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
}

:deep(.vditor-preview__action) {
  display: none !important;
}
</style>