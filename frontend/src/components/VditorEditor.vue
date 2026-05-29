<template>
  <div ref="editorRef" class="vditor-container"></div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import Vditor from 'vditor'
import 'vditor/dist/index.css'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  height: {
    type: Number,
    default: 400
  },
  mode: {
    type: String,
    default: 'sv'
  }
})

const emit = defineEmits(['update:modelValue'])

const editorRef = ref(null)
let vditor = null

const initVditor = () => {
  if (vditor) {
    vditor.destroy()
  }

  vditor = new Vditor(editorRef.value, {
    value: props.modelValue,
    height: props.height,
    mode: props.mode,
    placeholder: '使用 Markdown 语法编写文章...',
    toolbar: [
      'headings',
      'bold',
      'italic',
      'strike',
      '|',
      'line',
      'quote',
      'list',
      'ordered-list',
      'check',
      '|',
      'code',
      'inline-code',
      'link',
      'table',
      '|',
      'upload',
      'record',
      '|',
      'undo',
      'redo',
      '|',
      'preview',
      'fullscreen'
    ],
    upload: {
      accept: 'image/*',
      url: '/api/upload/image',
      fieldName: 'image',
      success (editor, msg) {
        const res = JSON.parse(msg)
        if (res.url) {
          vditor.insertValue(`![图片](${res.url})`)
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
      vditor.setValue(props.modelValue)
    }
  })
}

watch(() => props.modelValue, (newVal) => {
  if (vditor && vditor.getValue() !== newVal) {
    vditor.setValue(newVal)
  }
})

onMounted(() => {
  initVditor()
})

onBeforeUnmount(() => {
  if (vditor) {
    vditor.destroy()
    vditor = null
  }
})
</script>

<style scoped>
.vditor-container {
  border: 1px solid rgba(0, 0, 0, 0.23);
  border-radius: 4px;
}

:deep(.vditor) {
  border: none;
}

:deep(.vditor-toolbar) {
  border-bottom: 1px solid rgba(0, 0, 0, 0.23);
}
</style>