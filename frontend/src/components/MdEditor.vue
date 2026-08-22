<template>
  <div class="md-editor" :style="{ height: editorHeight + 'px' }">
    <!-- 工具栏 -->
    <div class="toolbar">
      <div class="toolbar-group">
        <el-tooltip content="粗体 (Ctrl+B)"><el-button text size="small" @click="insertFormat('**', '**')"><b>B</b></el-button></el-tooltip>
        <el-tooltip content="斜体 (Ctrl+I)"><el-button text size="small" @click="insertFormat('*', '*')"><i>I</i></el-button></el-tooltip>
        <el-tooltip content="删除线"><el-button text size="small" @click="insertFormat('~~', '~~')"><s>S</s></el-button></el-tooltip>
        <el-tooltip content="行内代码"><el-button text size="small" @click="insertFormat('`', '`')"><code>`</code></el-button></el-tooltip>
      </div>

      <div class="divider"></div>

      <div class="toolbar-group">
        <el-tooltip content="标题1"><el-button text size="small" @click="insertLine('# ')">H1</el-button></el-tooltip>
        <el-tooltip content="标题2"><el-button text size="small" @click="insertLine('## ')">H2</el-button></el-tooltip>
        <el-tooltip content="标题3"><el-button text size="small" @click="insertLine('### ')">H3</el-button></el-tooltip>
      </div>

      <div class="divider"></div>

      <div class="toolbar-group">
        <el-tooltip content="无序列表"><el-button text size="small" @click="insertLine('- ')"><el-icon><List /></el-icon></el-button></el-tooltip>
        <el-tooltip content="有序列表"><el-button text size="small" @click="insertLine('1. ')"><el-icon><List /></el-icon></el-button></el-tooltip>
        <el-tooltip content="任务列表"><el-button text size="small" @click="insertLine('- [ ] ')"><el-icon><Check /></el-icon></el-button></el-tooltip>
        <el-tooltip content="引用"><el-button text size="small" @click="insertLine('> ')"><el-icon><ChatSquare /></el-icon></el-button></el-tooltip>
      </div>

      <div class="divider"></div>

      <div class="toolbar-group">
        <el-tooltip content="链接 (Ctrl+K)"><el-button text size="small" @click="insertLink"><el-icon><Link /></el-icon></el-button></el-tooltip>
        <el-tooltip content="图片"><el-button text size="small" @click="insertImage"><el-icon><Picture /></el-icon></el-button></el-tooltip>
        <el-tooltip content="视频"><el-button text size="small" @click="insertVideo"><el-icon><VideoCamera /></el-icon></el-button></el-tooltip>
        <el-tooltip :content="isRecording ? '停止录音' : '语音输入'">
          <el-button text size="small" :type="isRecording ? 'danger' : 'default'" :loading="isRecording" @click="toggleRecording">
            <el-icon><Microphone /></el-icon>
          </el-button>
        </el-tooltip>
        <el-tooltip content="表格"><el-button text size="small" @click="insertTable"><el-icon><Grid /></el-icon></el-button></el-tooltip>
      </div>

      <div class="divider"></div>

      <div class="toolbar-group">
        <el-tooltip content="代码块"><el-button text size="small" @click="insertCodeBlock"><el-icon><Document /></el-icon></el-button></el-tooltip>
        <el-tooltip content="分割线"><el-button text size="small" @click="insertHr"><el-icon><Remove /></el-icon></el-button></el-tooltip>
        <el-tooltip content="数学公式"><el-button text size="small" @click="insertMath"><el-icon><Operation /></el-icon></el-button></el-tooltip>
        <el-dropdown trigger="click" @command="onContainerCommand">
          <el-button text size="small"><el-icon><Warning /></el-icon></el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="tip">💡 提示</el-dropdown-item>
              <el-dropdown-item command="warning">⚠️ 警告</el-dropdown-item>
              <el-dropdown-item command="danger">❌ 危险</el-dropdown-item>
              <el-dropdown-item command="info">ℹ️ 信息</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>

      <div class="toolbar-spacer"></div>

      <!-- 预览切换 -->
      <div class="toolbar-group">
        <el-button text size="small" :type="showPreview ? 'primary' : 'default'" @click="showPreview = !showPreview">
          <el-icon><View /></el-icon>&nbsp;预览
        </el-button>
      </div>
    </div>

    <!-- 编辑/预览主体 -->
    <div class="editor-body" :class="{ previewing: showPreview }">
      <textarea
        v-show="!showPreview"
        ref="textareaRef"
        class="editor-textarea"
        :value="content"
        :placeholder="placeholder"
        @input="onInput"
        @keydown="handleKeydown"
        @focus="triggerSuggestions"
        @scroll="syncScroll"
      ></textarea>

      <div v-show="showPreview" ref="previewRef" class="editor-preview" v-html="renderedContent"></div>

      <!-- 建议弹出层 -->
      <div
        v-if="suggestions.length > 0"
        class="suggestion-popover"
        :style="suggestionStyle"
      >
        <div
          v-for="(s, i) in suggestions"
          :key="s.text"
          class="suggestion-item"
          :class="{ active: i === selectedSuggestion }"
          @click="applySuggestion(s)"
          @mousedown.prevent
        >
          <span class="suggestion-text">{{ s.text }}</span>
          <span class="suggestion-detail text-secondary">{{ s.detail }}</span>
        </div>
      </div>
    </div>

    <!-- 底部统计 -->
    <div class="editor-footer">
      <span class="text-secondary">{{ stats.words }} 字 · {{ stats.paragraphs }} 段 · {{ stats.lines }} 行</span>
      <span v-if="showPreview" class="text-secondary">预览模式</span>
      <span class="ml-auto text-secondary">Markdown</span>
    </div>

  </div>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue'
import MarkdownIt from 'markdown-it'
import anchor from 'markdown-it-anchor'
import container from 'markdown-it-container'
import emoji from 'markdown-it-emoji'
import footnote from 'markdown-it-footnote'
import mark from 'markdown-it-mark'
import taskLists from 'markdown-it-task-lists'
import { katex as renderLatex } from 'katex'
import {
  List, Check, ChatSquare, Link, Picture, VideoCamera, Microphone,
  Grid, Remove, Operation, Warning, View, Document
} from '@element-plus/icons-vue'
import katex from 'katex'
import http from '@/api/http'
import { jcOpenHtml, jcFieldsConfig, jcCloseAll, jcToast } from '@/utils/jcu'

const props = defineProps({
  modelValue: { type: String, default: '' },
  placeholder: { type: String, default: '请输入 Markdown 内容…' },
  height: { type: Number, default: 400 }
})
const emit = defineEmits(['update:modelValue', 'voice-input'])

const content = ref(props.modelValue)
const textareaRef = ref(null)
const previewRef = ref(null)
const showPreview = ref(false)

const tableRows = ref(3)
const tableCols = ref(3)
const uploadType = ref('image')

const suggestions = ref([])
const selectedSuggestion = ref(0)
const suggestionPosition = ref({ top: 0, left: 0 })

// 语音识别
const isRecording = ref(false)
const recognition = ref(null)
const recognitionSupported = ref(false)

const editorHeight = computed(() => props.height)

const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
  breaks: true,
  highlight(str, lang) {
    if (lang && window.hljs) {
      try {
        const l = window.hljs.getLanguage(lang) ? lang : 'plaintext'
        return `<pre class="hljs-code" data-language="${l}"><code class="hljs language-${l}">${window.hljs.highlight(str, { language: l, ignoreIllegals: true }).value}</code></pre>`
      } catch (e) { /* ignore */ }
    }
    return `<pre class="hljs-code"><code class="hljs">${MarkdownIt().utils.escapeHtml(str)}</code></pre>`
  }
})
md.use(anchor, { permalink: anchor.permalink.ariaHidden({ placement: 'after', class: 'header-anchor', symbol: '#' }), level: [1, 2, 3, 4] })
;['tip', 'warning', 'danger', 'info'].forEach(t => {
  md.use(container, t, {
    render(tokens, idx) {
      if (tokens[idx].nesting === 1) {
        const label = { tip: '💡 提示', warning: '⚠️ 警告', danger: '❌ 危险', info: 'ℹ️ 信息' }[t]
        return `<div class="custom-container ${t}"><p class="custom-container-title">${label}</p>`
      }
      return '</div>'
    }
  })
})
md.use(emoji)
md.use(footnote)
md.use(mark)
md.use(taskLists, { enabled: true, label: true })

const katexBlock = (formula) => { try { return katex.renderToString(formula, { displayMode: true, throwOnError: false }) } catch (e) { return formula } }
const katexInline = (formula) => { try { return katex.renderToString(formula, { throwOnError: false }) } catch (e) { return formula } }

const renderMath = (text) => {
  let result = text
  result = result.replace(/\$\$([\s\S]*?)\$\$/g, (m, f) => katexBlock(f.trim()))
  result = result.replace(/\$(.+?)\$/g, (m, f) => katexInline(f.trim()))
  return result
}

const renderedContent = computed(() => {
  if (!content.value) return '<p class="text-secondary">预览区域</p>'
  return renderMath(md.render(content.value))
})

const stats = computed(() => {
  const text = content.value || ''
  const lines = text.split('\n').length
  const paragraphs = text.split(/\n\n+/).filter(p => p.trim()).length
  return { words: text.length, paragraphs, lines }
})

const suggestionStyle = computed(() => ({
  top: suggestionPosition.value.top + 'px',
  left: suggestionPosition.value.left + 'px'
}))

watch(() => props.modelValue, (v) => { if (v !== content.value) content.value = v })

const onInput = () => {
  content.value = textareaRef.value.value
  emit('update:modelValue', content.value)
  triggerSuggestions()
}

const handleInput = () => emit('update:modelValue', content.value)

const syncScroll = () => {
  // 简单保持预览滚动同步（可扩展）
}

const applySuggestion = (suggestion) => {
  const textarea = textareaRef.value
  if (!textarea) return
  const pos = textarea.selectionStart
  const before = content.value.substring(0, pos)
  const lastSpace = before.lastIndexOf(' ')
  const lastNewline = before.lastIndexOf('\n')
  const startIndex = Math.max(lastSpace, lastNewline) + 1
  const insert = suggestion.insert || ''
  const markerIndex = insert.indexOf('|')
  const clean = insert.replace(/\|/g, '')
  content.value = content.value.substring(0, startIndex) + clean + content.value.substring(pos)
  suggestions.value = []
  nextTick(() => {
    textarea.focus()
    const cursorPos = startIndex + (markerIndex >= 0 ? markerIndex : clean.length)
    textarea.setSelectionRange(cursorPos, cursorPos)
  })
  emit('update:modelValue', content.value)
}

const getSuggestions = (query) => {
  const list = [
    { text: '**粗体**', detail: '加粗', insert: '**|**' },
    { text: '*斜体*', detail: '倾斜', insert: '*|*' },
    { text: '~~删除线~~', detail: '删除', insert: '~~|~~' },
    { text: '`代码`', detail: '行内代码', insert: '`|`' },
    { text: '==高亮==', detail: '高亮', insert: '==|==' },
    { text: '# 标题', detail: '一级标题', insert: '# |' },
    { text: '## 标题', detail: '二级标题', insert: '## |' },
    { text: '### 标题', detail: '三级标题', insert: '### |' },
    { text: '- 列表项', detail: '无序列表', insert: '- |' },
    { text: '1. 列表项', detail: '有序列表', insert: '1. |' },
    { text: '- [ ] 任务', detail: '任务', insert: '- [ ] |' },
    { text: '> 引用', detail: '引用', insert: '> |' },
    { text: '[链接](url)', detail: '链接', insert: '[|]()' },
    { text: '![图片](url)', detail: '图片', insert: '![](|)' },
    { text: '```代码块', detail: '代码块', insert: '```\n|\n```' },
    { text: '---', detail: '分割线', insert: '\n---\n' },
    { text: ':::tip', detail: '提示框', insert: ':::tip\n|\n:::' },
    { text: ':::warning', detail: '警告框', insert: ':::warning\n|\n:::' },
    { text: ':::danger', detail: '危险框', insert: ':::danger\n|\n:::' },
    { text: ':::info', detail: '信息框', insert: ':::info\n|\n:::' },
    { text: '$$公式$$', detail: '公式', insert: '$$\n|\n$$' }
  ]
  return list.filter(s => s.text.toLowerCase().includes(query.toLowerCase()) || s.detail.toLowerCase().includes(query.toLowerCase()))
}

const triggerSuggestions = () => {
  const textarea = textareaRef.value
  if (!textarea) return
  const pos = textarea.selectionStart
  const before = content.value.substring(0, pos)
  const lastSpace = before.lastIndexOf(' ')
  const lastNewline = before.lastIndexOf('\n')
  const startIndex = Math.max(lastSpace, lastNewline) + 1
  const word = before.substring(startIndex)
  if (word.length > 0) {
    const matched = getSuggestions(word)
    if (matched.length) {
      suggestions.value = matched
      selectedSuggestion.value = 0
      updateSuggestionPosition()
      return
    }
  }
  suggestions.value = []
}

const updateSuggestionPosition = () => {
  const textarea = textareaRef.value
  if (!textarea) return
  const rect = textarea.getBoundingClientRect()
  const editorRect = textarea.closest('.editor-body')?.getBoundingClientRect() || rect
  suggestionPosition.value = {
    top: rect.top - editorRect.top + textarea.clientHeight - 40,
    left: 8
  }
}

const insertFormat = (prefix, suffix) => {
  const textarea = textareaRef.value
  if (!textarea) return
  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const selected = content.value.substring(start, end)
  content.value = content.value.substring(0, start) + prefix + selected + suffix + content.value.substring(end)
  nextTick(() => {
    textarea.focus()
    textarea.setSelectionRange(start + prefix.length, end + prefix.length)
  })
  emit('update:modelValue', content.value)
}

const insertLine = (prefix) => {
  const textarea = textareaRef.value
  if (!textarea) return
  const start = textarea.selectionStart
  const lineStart = content.value.lastIndexOf('\n', start - 1) + 1
  content.value = content.value.substring(0, lineStart) + prefix + content.value.substring(lineStart)
  nextTick(() => {
    textarea.focus()
    const p = lineStart + prefix.length
    textarea.setSelectionRange(p, p)
  })
  emit('update:modelValue', content.value)
}

const insertAtCursor = (text) => {
  const textarea = textareaRef.value
  if (!textarea) { content.value += text; emit('update:modelValue', content.value); return }
  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  content.value = content.value.substring(0, start) + text + content.value.substring(end)
  nextTick(() => {
    textarea.focus()
    const p = start + text.length
    textarea.setSelectionRange(p, p)
  })
  emit('update:modelValue', content.value)
}

const insertCodeBlock = () => {
  const textarea = textareaRef.value
  const start = textarea?.selectionStart ?? content.value.length
  const end = textarea?.selectionEnd ?? content.value.length
  const selected = content.value.substring(start, end) || ''
  insertAtCursor(`\n\`\`\`\n${selected}\n\`\`\`\n`)
}

const insertLink = () => {
  const textarea = textareaRef.value
  const start = textarea?.selectionStart ?? 0
  const end = textarea?.selectionEnd ?? 0
  const cfg = jcFieldsConfig([
    { name: 'text', label: '链接文字', value: content.value.substring(start, end) || '链接文字' },
    { name: 'url', label: '链接地址', placeholder: 'https://', required: true }
  ])
  jcOpenHtml({
    title: '插入链接',
    content: cfg.html,
    width: 460,
    size: 'sm',
    buttons: [
      { text: '取消', type: 'default', action: () => jcCloseAll() },
      {
        text: '确认',
        type: 'primary',
        action: () => {
          if (!cfg.validate(document)) return
          const v = cfg.collect(document)
          if (!v.url) return
          const markdown = `[${v.text}](${v.url})`
          insertAtCursor(markdown)
          jcCloseAll()
        }
      }
    ]
  })
}

const insertImage = () => { uploadType.value = 'image'; openUploadDialog() }
const insertVideo = () => { uploadType.value = 'video'; openUploadDialog() }

const openUploadDialog = () => {
  const isImage = uploadType.value === 'image'
  const id = 'jc-md-upload'
  const html = `
    <div data-jc-upload-area style="border:2px dashed var(--campus-border,#e6e9f2);border-radius:12px;padding:28px;text-align:center;cursor:pointer;transition:all .15s;">
      <input type="file" data-jc-file accept="${isImage ? 'image/*' : 'video/*'}" style="display:none;" />
      <div style="font-size:38px;color:var(--campus-primary,#4f6ef7);margin-bottom:8px;">${isImage ? '🖼️' : '🎬'}</div>
      <div style="font-weight:600;color:var(--campus-text,#0f172a);">点击选择${isImage ? '图片' : '视频'}</div>
      <div style="font-size:12px;color:var(--campus-text-secondary,#6b7280);margin-top:6px;">${isImage ? '支持 JPG、PNG、GIF' : '支持 MP4、WebM'}</div>
    </div>
    <div data-jc-upload-progress style="display:none;margin-top:16px;">
      <div style="height:8px;background:var(--campus-border,#e6e9f2);border-radius:99px;overflow:hidden;">
        <div data-jc-upload-bar style="height:100%;width:0%;background:var(--campus-primary,#4f6ef7);border-radius:99px;transition:width .2s;"></div>
      </div>
      <div data-jc-upload-pct style="text-align:center;font-size:12px;color:var(--campus-text-secondary,#6b7280);margin-top:6px;">0%</div>
    </div>
    <div data-jc-upload-error style="display:none;color:var(--campus-danger,#ef4444);font-size:13px;margin-top:12px;">上传失败，请重试</div>
    <div data-jc-upload-status style="display:none;color:var(--campus-success,#22c55e);font-size:13px;margin-top:12px;">上传成功，即将插入…</div>
  `
  jcOpenHtml({
    id,
    title: isImage ? '插入图片' : '插入视频',
    content: html,
    width: 480,
    size: 'sm',
    closeOnOverlay: false,
    onMount: (root) => {
      const area = root.querySelector('[data-jc-upload-area]')
      const file = root.querySelector('[data-jc-file]')
      const progressWrap = root.querySelector('[data-jc-upload-progress]')
      const bar = root.querySelector('[data-jc-upload-bar]')
      const pct = root.querySelector('[data-jc-upload-pct]')
      const err = root.querySelector('[data-jc-upload-error]')
      const ok = root.querySelector('[data-jc-upload-status]')

      const setProgress = () => {
        progressWrap.style.display = 'block'
        bar.style.width = window.__jcUploadProgress + '%'
        pct.textContent = window.__jcUploadProgress + '%'
      }

      area.addEventListener('click', () => file.click())
      area.addEventListener('dragover', (e) => { e.preventDefault(); area.style.borderColor = 'var(--campus-primary)' })
      area.addEventListener('dragleave', () => { area.style.borderColor = '' })
      area.addEventListener('drop', (e) => {
        e.preventDefault()
        const f = e.dataTransfer.files[0]
        if (f) doUpload(f, { progressWrap, bar, pct, err, ok })
      })
      file.addEventListener('change', (e) => {
        const f = e.target.files[0]
        e.target.value = ''
        if (f) doUpload(f, { progressWrap, bar, pct, err, ok })
      })
    },
    buttons: [{ text: '关闭', type: 'primary', action: () => jcCloseAll() }]
  })
}

const doUpload = (file, els) => {
  const isImage = uploadType.value === 'image'
  const formData = new FormData()
  formData.append(uploadType.value, file)
  window.__jcUploadProgress = 0
  els.err.style.display = 'none'
  els.progressWrap.style.display = 'block'
  els.bar.style.width = '0%'
  els.pct.textContent = '0%'

  http.post(`/upload/${uploadType.value}`, formData, {
    onUploadProgress: (e) => {
      if (e.total) {
        window.__jcUploadProgress = Math.round((e.loaded / e.total) * 100)
        els.bar.style.width = window.__jcUploadProgress + '%'
        els.pct.textContent = window.__jcUploadProgress + '%'
      }
    }
  }).then(({ data }) => {
    const url = data.url
    els.progressWrap.style.display = 'none'
    els.ok.style.display = 'block'
    if (isImage) insertAtCursor(`\n![图片](${url})\n`)
    else insertAtCursor(`\n<video src="${url}" controls></video>\n`)
    setTimeout(() => jcCloseAll(), 500)
  }).catch(() => {
    els.progressWrap.style.display = 'none'
    els.err.style.display = 'block'
  })
}

const insertTable = () => {
  const cfg = jcFieldsConfig([
    { name: 'rows', label: '行数', type: 'number', value: tableRows.value },
    { name: 'cols', label: '列数', type: 'number', value: tableCols.value }
  ])
  jcOpenHtml({
    title: '插入表格',
    content: cfg.html,
    width: 380,
    size: 'sm',
    buttons: [
      { text: '取消', type: 'default', action: () => jcCloseAll() },
      {
        text: '插入',
        type: 'primary',
        action: () => {
          const v = cfg.collect(document)
          const rows = Math.min(20, Math.max(2, Number(v.rows) || 3))
          const cols = Math.min(10, Math.max(1, Number(v.cols) || 3))
          tableRows.value = rows
          tableCols.value = cols
          confirmTable()
          jcCloseAll()
        }
      }
    ]
  })
}

const confirmTable = () => {
  const header = '| ' + Array(tableCols.value).fill('列').map((h, i) => `${h}${i + 1}`).join(' | ') + ' |'
  const sep = '| ' + Array(tableCols.value).fill('---').join(' | ') + ' |'
  const body = Array(Math.max(1, tableRows.value - 1)).fill(null)
    .map(() => '| ' + Array(tableCols.value).fill('').join(' | ') + ' |').join('\n')
  insertAtCursor('\n' + header + '\n' + sep + '\n' + body + '\n')
}

const insertHr = () => insertAtCursor('\n---\n\n')
const insertMark = () => insertFormat('==', '==')

const insertMath = () => {
  const textarea = textareaRef.value
  const start = textarea?.selectionStart ?? 0
  const end = textarea?.selectionEnd ?? 0
  const selected = content.value.substring(start, end) || 'x^2 + y^2 = z^2'
  insertAtCursor(`\n$$\n${selected}\n$$\n`)
}

const onContainerCommand = (cmd) => {
  const textarea = textareaRef.value
  const start = textarea?.selectionStart ?? 0
  const end = textarea?.selectionEnd ?? 0
  const selected = content.value.substring(start, end) || '这里放置内容'
  insertAtCursor(`\n:::${cmd}\n${selected}\n:::\n`)
}

const handleKeydown = (e) => {
  if (e.ctrlKey || e.metaKey) {
    const k = e.key.toLowerCase()
    if (k === 'b') { e.preventDefault(); insertFormat('**', '**') }
    else if (k === 'i') { e.preventDefault(); insertFormat('*', '*') }
    else if (k === 'k') { e.preventDefault(); insertLink() }
    else if (k === 's') e.preventDefault()
  }
  if (e.key === 'Tab') {
    e.preventDefault()
    if (suggestions.value.length) applySuggestion(suggestions.value[selectedSuggestion.value])
    else insertAtCursor('  ')
  }
  if (e.key === 'Enter' && suggestions.value.length) { e.preventDefault(); applySuggestion(suggestions.value[selectedSuggestion.value]) }
}

// 语音识别
const initRecognition = () => {
  if ('webkitSpeechRecognition' in window || 'SpeechRecognition' in window) {
    const SR = window.SpeechRecognition || window.webkitSpeechRecognition
    recognition.value = new SR()
    recognition.value.lang = 'zh-CN'
    recognition.value.interimResults = true
    recognition.value.continuous = false
    recognition.value.onstart = () => { isRecording.value = true }
    recognition.value.onresult = (event) => {
      let finalTranscript = ''
      for (let i = event.resultIndex; i < event.results.length; i++) {
        if (event.results[i].isFinal) finalTranscript += event.results[i][0].transcript
      }
      if (finalTranscript) { insertAtCursor(finalTranscript); emit('voice-input', finalTranscript) }
    }
    recognition.value.onerror = () => { isRecording.value = false }
    recognition.value.onend = () => { isRecording.value = false }
    recognitionSupported.value = true
  }
}
const toggleRecording = () => {
  if (!recognitionSupported.value) { jcToast('您的浏览器不支持语音识别，请使用 Chrome', 'warning'); return }
  if (isRecording.value) recognition.value?.stop()
  else try { recognition.value?.start() } catch (e) { isRecording.value = false }
}

initRecognition()
</script>

<style scoped>
.md-editor {
  display: flex;
  flex-direction: column;
  border: 1px solid var(--campus-border);
  border-radius: 8px;
  overflow: hidden;
  background: #fff;
}
.toolbar {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 2px;
  padding: 4px 6px;
  border-bottom: 1px solid var(--campus-border);
  background: #fafafa;
}
.toolbar-group { display: flex; align-items: center; }
.toolbar-spacer { flex: 1; }
.divider { width: 1px; height: 18px; background: var(--campus-border); margin: 0 6px; }
.editor-body { position: relative; flex: 1; min-height: 0; }
.editor-textarea {
  width: 100%;
  height: 100%;
  border: none;
  outline: none;
  resize: none;
  padding: 12px;
  font-family: 'Fira Code', Consolas, monospace;
  font-size: 14px;
  line-height: 1.7;
  box-sizing: border-box;
}
.editor-preview {
  height: 100%;
  overflow: auto;
  padding: 12px 16px;
  box-sizing: border-box;
}
.editor-footer {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 12px;
  padding: 4px 8px;
  border-top: 1px solid var(--campus-border);
  color: var(--campus-text-secondary);
}
.ml-auto { margin-left: auto; }
.suggestion-popover {
  position: absolute;
  z-index: 20;
  width: 260px;
  max-height: 240px;
  overflow: auto;
  background: #fff;
  border: 1px solid var(--campus-border);
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.12);
}
.suggestion-item { padding: 6px 10px; cursor: pointer; display: flex; justify-content: space-between; gap: 8px; font-size: 13px; }
.suggestion-item:hover, .suggestion-item.active { background: #f0f2ff; }
.text-secondary { color: var(--campus-text-secondary); }
.mt-2 { margin-top: 8px; }
</style>
