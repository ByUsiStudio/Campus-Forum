<template>
  <div ref="viewerRef" class="vditor-reset article-content"></div>
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

const renderMarkdown = async () => {
  if (!viewerRef.value || !props.value) return

  await Vditor.preview(viewerRef.value, props.value, {
    cdn: 'https://unpkg.com/vditor@3.10.9',
    mode: 'light',
    lazyLoadImage: false,
    after: () => {
      processMediaElements()
    }
  })
}

const processMediaElements = () => {
  if (!viewerRef.value) return

  const images = viewerRef.value.querySelectorAll('img')
  images.forEach(img => {
    img.style.maxWidth = '100%'
    img.style.height = 'auto'
    img.style.borderRadius = '8px'
    img.style.margin = '16px 0'
    img.style.display = 'block'
    img.style.boxShadow = '0 2px 8px rgba(0, 0, 0, 0.1)'
    img.loading = 'lazy'
  })

  const videos = viewerRef.value.querySelectorAll('video')
  videos.forEach(video => {
    video.style.maxWidth = '100%'
    video.style.maxHeight = '480px'
    video.style.height = 'auto'
    video.style.borderRadius = '8px'
    video.style.margin = '16px 0'
    video.style.display = 'block'
    video.style.boxShadow = '0 2px 8px rgba(0, 0, 0, 0.1)'
    video.controls = true
    video.playsInline = true
  })
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
  padding: 16px;
  overflow-wrap: break-word;
}

.article-content {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  line-height: 1.8;
  color: #333;
}

:deep(img) {
  max-width: 100%;
  max-height: 600px;
  height: auto;
  border-radius: 8px;
  margin: 16px 0;
  display: block;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: transform 0.2s;
}

:deep(img:hover) {
  transform: scale(1.02);
}

:deep(video) {
  max-width: 100%;
  max-height: 480px;
  height: auto;
  border-radius: 8px;
  margin: 16px 0;
  display: block;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  background: #000;
}

:deep(pre) {
  background-color: #f6f8fa;
  padding: 16px;
  border-radius: 8px;
  overflow-x: auto;
  margin: 16px 0;
  border: 1px solid #e1e4e8;
}

:deep(code) {
  background-color: #f6f8fa;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'SF Mono', Consolas, 'Liberation Mono', Menlo, monospace;
  font-size: 14px;
}

:deep(pre code) {
  background-color: transparent;
  padding: 0;
  font-size: 14px;
}

:deep(.vditor-reset) {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  line-height: 1.8;
}

:deep(h1), :deep(h2), :deep(h3), :deep(h4), :deep(h5), :deep(h6) {
  margin-top: 24px;
  margin-bottom: 16px;
  font-weight: 600;
  line-height: 1.25;
  color: #24292e;
}

:deep(h1) {
  font-size: 2em;
  border-bottom: 1px solid #e1e4e8;
  padding-bottom: 8px;
}

:deep(h2) {
  font-size: 1.5em;
  border-bottom: 1px solid #e1e4e8;
  padding-bottom: 8px;
}

:deep(h3) {
  font-size: 1.25em;
}

:deep(p) {
  margin-bottom: 16px;
}

:deep(ul), :deep(ol) {
  padding-left: 2em;
  margin-bottom: 16px;
}

:deep(li) {
  margin-bottom: 4px;
}

:deep(blockquote) {
  padding: 16px 20px;
  color: #6a737d;
  border-left: 4px solid #dfe2e5;
  margin: 16px 0;
  background: #f6f8fa;
  border-radius: 4px;
}

:deep(blockquote p) {
  margin-bottom: 0;
}

:deep(table) {
  border-collapse: collapse;
  width: 100%;
  margin: 16px 0;
  display: block;
  overflow-x: auto;
}

:deep(table th), :deep(table td) {
  padding: 12px 16px;
  border: 1px solid #dfe2e5;
}

:deep(table th) {
  font-weight: 600;
  background-color: #f6f8fa;
}

:deep(table tr:nth-child(2n)) {
  background-color: #f6f8fa;
}

:deep(a) {
  color: rgb(var(--v-theme-primary));
  text-decoration: none;
}

:deep(a:hover) {
  text-decoration: underline;
}

:deep(hr) {
  border: none;
  height: 2px;
  background: #e1e4e8;
  margin: 24px 0;
}

:deep(.task-list-item) {
  list-style: none;
}

:deep(.task-list-item input) {
  margin-right: 8px;
}

:deep(.vditor-preview__action) {
  display: none !important;
}

:deep(.vditor-preview__action--current) {
  display: none !important;
}
</style>