<template>
  <div class="bytemd-editor">
    <div ref="editorRef" class="bytemd-container"></div>
    <UploadModal
      v-model:show="showImageModal"
      upload-type="image"
      @upload-success="handleUploadSuccess"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, watch, nextTick } from 'vue'
import { Editor, Viewer } from 'bytemd'
import gfm from '@bytemd/plugin-gfm'
import highlight from '@bytemd/plugin-highlight'
import math from '@bytemd/plugin-math'
import mediumZoom from '@bytemd/plugin-medium-zoom'
import 'bytemd/dist/index.css'
import 'katex/dist/katex.min.css'
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
    default: 'split'
  },
  preview: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['update:modelValue'])

const editorRef = ref(null)
let editor = null
let isInitialized = false
const showImageModal = ref(false)
const pendingInsert = ref('')

const plugins = [
  gfm(),
  highlight(),
  math(),
  mediumZoom()
]

const handleUploadSuccess = (data) => {
  showImageModal.value = false
  data.files.forEach(file => {
    if (data.type === 'image') {
      const markdown = `![${file.name}](${file.url})`
      if (editor) {
        const pos = editor.value().selection
        editor.value().replaceSelection(markdown)
      } else {
        emit('update:modelValue', props.modelValue + markdown)
      }
    }
  })
}

const initEditor = () => {
  if (!editorRef.value || isInitialized) return

  editor = new Editor({
    target: editorRef.value,
    props: {
      value: props.modelValue,
      plugins
    }
  })

  editor.on('change', (e) => {
    emit('update:modelValue', e.detail.value)
  })

  isInitialized = true
}

onMounted(() => {
  nextTick(() => {
    initEditor()
  })
})

onBeforeUnmount(() => {
  if (editor) {
    editor.$destroy()
    editor = null
    isInitialized = false
  }
})

watch(() => props.modelValue, (newVal) => {
  if (editor && newVal !== editor.value().innerText) {
    editor.$set({ value: newVal })
  }
})

watch(() => props.height, (newVal) => {
  if (editorRef.value) {
    editorRef.value.style.minHeight = `${newVal}px`
  }
}, { immediate: true })
</script>

<style scoped>
.bytemd-editor {
  width: 100%;
  border: 1px solid #E7E0EC;
  border-radius: 12px;
  overflow: hidden;
  background: #fff;
}

.bytemd-container {
  width: 100%;
  min-height: v-bind('props.height + "px"');
}

.bytemd-container :deep(.bytemd) {
  height: v-bind('props.height + "px"');
}

.bytemd-container :deep(.bytemd-editor) {
  background: #fff;
}

.bytemd-container :deep(.bytemd-preview) {
  background: #FAFAFA;
  padding: 16px 20px;
}

.bytemd-container :deep(.bytemd-toolbar) {
  background: linear-gradient(135deg, #F8F7FF 0%, #F1F5F9 100%);
  border-bottom: 1px solid #E7E0EC;
  padding: 8px 12px;
}

.bytemd-container :deep(.bytemd-toolbar button) {
  border-radius: 6px;
  transition: all 0.2s ease;
}

.bytemd-container :deep(.bytemd-toolbar button:hover) {
  background: rgba(103, 80, 164, 0.1);
}

.bytemd-container :deep(.bytemd-statusbar) {
  background: #F8F7FF;
  border-top: 1px solid #E7E0EC;
  padding: 6px 12px;
  font-size: 12px;
  color: #625B71;
}

.bytemd-container :deep(.bytemd-preview) {
  border-left: 1px solid #E7E0EC;
}
</style>