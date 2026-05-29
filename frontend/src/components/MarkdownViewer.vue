<template>
  <div ref="viewerRef" class="vditor-reset markdown-body"></div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import Vditor from 'vditor'
import 'vditor/dist/index.css'

const props = defineProps({
  value: {
    type: String,
    default: ''
  }
})

const viewerRef = ref(null)
let vditor = null

const renderMarkdown = () => {
  if (!viewerRef.value) return

  if (!vditor) {
    vditor = new Vditor(viewerRef.value, {
      value: props.value,
      mode: 'preview',
      height: 'auto',
      language: 'zh_CN',
      globalReset: false,
      input: () => {},
      after: () => {}
    })
  } else {
    vditor.setValue(props.value)
  }
}

watch(() => props.value, () => {
  renderMarkdown()
})

onMounted(() => {
  renderMarkdown()
})
</script>

<style scoped>
.vditor-reset {
  max-width: 100%;
}

:deep(img) {
  max-width: 100%;
  height: auto;
}

:deep(video) {
  max-width: 100%;
  height: auto;
}

:deep(pre) {
  background-color: #f6f8fa;
  padding: 16px;
  border-radius: 6px;
  overflow-x: auto;
}

:deep(code) {
  background-color: #f6f8fa;
  padding: 2px 6px;
  border-radius: 3px;
  font-family: 'SF Mono', Consolas, monospace;
}

:deep(pre code) {
  background-color: transparent;
  padding: 0;
}
</style>