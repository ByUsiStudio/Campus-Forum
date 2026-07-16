<template>
  <div class="md-editor">
    <div class="toolbar">
      <div class="toolbar-group">
        <v-btn variant="text" size="small" @click="insertFormat('**', '**')" title="粗体 (Ctrl+B)">
          <v-icon size="18">mdi-format-bold</v-icon>
        </v-btn>
        <v-btn variant="text" size="small" @click="insertFormat('*', '*')" title="斜体 (Ctrl+I)">
          <v-icon size="18">mdi-format-italic</v-icon>
        </v-btn>
        <v-btn variant="text" size="small" @click="insertFormat('~~', '~~')" title="删除线">
          <v-icon size="18">mdi-format-strikethrough</v-icon>
        </v-btn>
        <v-btn variant="text" size="small" @click="insertFormat('`', '`')" title="行内代码">
          <v-icon size="18">mdi-code-tags</v-icon>
        </v-btn>
      </div>

      <v-divider vertical class="mx-1" />

      <div class="toolbar-group">
        <v-btn variant="text" size="small" @click="insertLine('# ')" title="标题1">H1</v-btn>
        <v-btn variant="text" size="small" @click="insertLine('## ')" title="标题2">H2</v-btn>
        <v-btn variant="text" size="small" @click="insertLine('### ')" title="标题3">H3</v-btn>
      </div>

      <v-divider vertical class="mx-1" />

      <div class="toolbar-group">
        <v-btn variant="text" size="small" @click="insertLine('- ')" title="无序列表">
          <v-icon size="18">mdi-format-list-bulleted</v-icon>
        </v-btn>
        <v-btn variant="text" size="small" @click="insertLine('1. ')" title="有序列表">
          <v-icon size="18">mdi-format-list-numbered</v-icon>
        </v-btn>
        <v-btn variant="text" size="small" @click="insertLine('- [ ] ')" title="任务列表">
          <v-icon size="18">mdi-checkbox-marked-outline</v-icon>
        </v-btn>
        <v-btn variant="text" size="small" @click="insertLine('> ')" title="引用">
          <v-icon size="18">mdi-format-quote-close</v-icon>
        </v-btn>
      </div>

      <v-divider vertical class="mx-1" />

      <div class="toolbar-group">
        <v-btn variant="text" size="small" @click="insertLink" title="链接 (Ctrl+K)">
          <v-icon size="18">mdi-link</v-icon>
        </v-btn>
        <v-btn variant="text" size="small" @click="insertImage" title="图片">
          <v-icon size="18">mdi-image</v-icon>
        </v-btn>
        <v-btn variant="text" size="small" @click="insertVideo" title="视频">
          <v-icon size="18">mdi-video</v-icon>
        </v-btn>
        <v-btn variant="text" size="small" :color="isRecording ? 'red' : undefined" @click="toggleRecording" :loading="isRecording" :title="isRecording ? '停止录音' : '语音输入'">
          <v-icon size="18">mdi-microphone{{ isRecording ? '-off' : '' }}</v-icon>
        </v-btn>
        <v-btn variant="text" size="small" @click="insertTable" title="表格">
          <v-icon size="18">mdi-table</v-icon>
        </v-btn>
      </div>

      <v-divider vertical class="mx-1" />

      <div class="toolbar-group">
        <v-btn variant="text" size="small" @click="insertCodeBlock" title="代码块">
          <v-icon size="18">mdi-code-braces</v-icon>
        </v-btn>
        <v-btn variant="text" size="small" @click="insertHr" title="分割线">
          <v-icon size="18">mdi-minus</v-icon>
        </v-btn>
        <v-btn variant="text" size="small" @click="insertMath" title="数学公式">
          <v-icon size="18">mdi-function-variant</v-icon>
        </v-btn>
        <v-btn variant="text" size="small" @click="insertMark" title="高亮文本">
          <v-icon size="18">mdi-marker</v-icon>
        </v-btn>
        <v-btn variant="text" size="small" @click="insertFootnote" title="脚注">
          <v-icon size="18">mdi-text-box-outline</v-icon>
        </v-btn>
        <v-menu offset="8">
          <template #activator="{ props }">
            <v-btn variant="text" size="small" v-bind="props" title="提示框">
              <v-icon size="18">mdi-lightbulb-outline</v-icon>
              <v-icon size="14">mdi-chevron-down</v-icon>
            </v-btn>
          </template>
          <v-list density="compact" class="container-menu">
            <v-list-item @click="insertContainer('tip')">
              <template #prepend>
                <v-icon size="18" color="green">mdi-lightbulb</v-icon>
              </template>
              <v-list-item-title>💡 提示</v-list-item-title>
            </v-list-item>
            <v-list-item @click="insertContainer('info')">
              <template #prepend>
                <v-icon size="18" color="blue">mdi-information</v-icon>
              </template>
              <v-list-item-title>ℹ️ 信息</v-list-item-title>
            </v-list-item>
            <v-list-item @click="insertContainer('warning')">
              <template #prepend>
                <v-icon size="18" color="orange">mdi-alert</v-icon>
              </template>
              <v-list-item-title>⚠️ 警告</v-list-item-title>
            </v-list-item>
            <v-list-item @click="insertContainer('danger')">
              <template #prepend>
                <v-icon size="18" color="red">mdi-close-circle</v-icon>
              </template>
              <v-list-item-title>❌ 危险</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-menu>
      </div>

      <v-divider vertical class="mx-1" />

      <div class="toolbar-group">
        <v-btn variant="text" size="small" @click="formatDoc('undo')" title="撤销 (Ctrl+Z)">
          <v-icon size="18">mdi-undo</v-icon>
        </v-btn>
        <v-btn variant="text" size="small" @click="formatDoc('redo')" title="重做 (Ctrl+Y)">
          <v-icon size="18">mdi-redo</v-icon>
        </v-btn>
      </div>

      <v-spacer />

      <div class="toolbar-group">
        <v-btn variant="text" size="small" :color="viewMode === 'edit' ? 'primary' : undefined" @click="viewMode = 'edit'" title="编辑">
          <v-icon size="18">mdi-pencil</v-icon>
        </v-btn>
        <v-btn variant="text" size="small" :color="viewMode === 'split' ? 'primary' : undefined" @click="viewMode = 'split'" title="分屏">
          <v-icon size="18">mdi-view-split-vertical</v-icon>
        </v-btn>
        <v-btn variant="text" size="small" :color="viewMode === 'preview' ? 'primary' : undefined" @click="viewMode = 'preview'" title="预览">
          <v-icon size="18">mdi-eye</v-icon>
        </v-btn>
      </div>

      <v-divider vertical class="mx-1" />

      <div class="status-bar">
        <span class="status-item">
          <v-icon size="14">mdi-file-document-outline</v-icon>
          <span>{{ stats.words }} 字</span>
        </span>
        <span class="status-divider">|</span>
        <span class="status-item">
          <v-icon size="14">mdi-paragraph</v-icon>
          <span>{{ stats.paragraphs }} 段落</span>
        </span>
        <span class="status-divider">|</span>
        <span class="status-item">
          <v-icon size="14">mdi-format-list-numbered</v-icon>
          <span>{{ stats.lines }} 行</span>
        </span>
        <span class="status-divider">|</span>
        <span class="status-item">
          <v-icon size="14">mdi-tag</v-icon>
          <span>{{ stats.tokens }} 标记</span>
        </span>
      </div>
    </div>

    <div class="editor-body" :class="'view-' + viewMode">
      <div v-show="viewMode !== 'preview'" class="editor-pane">
        <textarea
          ref="textareaRef"
          v-model="content"
          class="editor-textarea"
          :placeholder="placeholder"
          @input="handleInput"
          @keydown="handleKeydown"
          @scroll="syncScroll"
        ></textarea>
        
        <div v-if="suggestions.length > 0" class="suggestion-box" :style="suggestionStyle">
          <div
            v-for="(suggestion, index) in suggestions"
            :key="suggestion.text"
            class="suggestion-item"
            :class="{ active: selectedSuggestion === index }"
            @click="applySuggestion(suggestion)"
            @mouseenter="selectedSuggestion = index"
          >
            <v-icon size="14" :color="suggestion.color">{{ suggestion.icon }}</v-icon>
            <span class="suggestion-text">{{ suggestion.text }}</span>
            <span class="suggestion-detail">{{ suggestion.detail }}</span>
          </div>
        </div>
      </div>
      
      <div v-show="viewMode !== 'edit'" class="preview-pane" ref="previewPaneRef" @scroll="handlePreviewScroll">
        <div class="preview-content markdown-body" v-html="renderedContent"></div>
      </div>
    </div>

    <div v-if="uploadModalVisible" class="upload-modal-overlay" @click.self="closeUploadModal">
      <div class="upload-modal-container">
        <div class="upload-modal-header">
          <div class="modal-title">
            <v-icon size="24" color="primary">{{ uploadType === 'image' ? 'mdi-image' : 'mdi-video' }}</v-icon>
            <h3>{{ uploadType === 'image' ? '上传图片' : '上传视频' }}</h3>
          </div>
          <v-btn icon variant="text" size="small" @click="closeUploadModal">
            <v-icon>mdi-close</v-icon>
          </v-btn>
        </div>
        
        <div class="upload-modal-body">
          <div class="upload-area" @click="triggerFileInput" @dragover.prevent="onDragOver" @dragleave="onDragLeave" @drop.prevent="onDrop">
            <div v-if="!hasFiles" class="upload-placeholder">
              <v-icon size="48" color="grey">mdi-cloud-upload</v-icon>
              <p class="upload-hint">点击或拖拽文件到此处</p>
              <p class="upload-format">支持格式：{{ uploadType === 'image' ? 'JPG, PNG, GIF, WebP' : 'MP4, WebM, MOV' }}</p>
            </div>
            
            <div v-else class="upload-preview">
              <div v-for="(file, index) in uploadFiles" :key="index" class="file-preview-item">
                <div class="preview-thumb">
                  <img v-if="uploadType === 'image'" :src="file.preview" alt="preview" />
                  <div v-else class="video-placeholder">
                    <v-icon size="32" color="grey">mdi-video</v-icon>
                  </div>
                </div>
                <div class="file-info">
                  <span class="file-name">{{ file.name }}</span>
                  <span class="file-size">{{ formatFileSize(file.size) }}</span>
                </div>
                <v-btn icon variant="text" size="small" color="error" @click="removeFile(index)">
                  <v-icon>mdi-delete</v-icon>
                </v-btn>
              </div>
            </div>
          </div>
          
          <input
            ref="fileInputRef"
            type="file"
            :accept="uploadType === 'image' ? 'image/*' : 'video/*'"
            multiple
            class="file-input"
            @change="handleFileChange"
          />
          
          <div v-if="uploadProgress > 0" class="upload-progress">
            <div class="progress-info">
              <span>正在上传...</span>
              <span>{{ uploadProgress }}%</span>
            </div>
            <v-progress-linear :model-value="uploadProgress" color="primary" height="4" rounded />
          </div>
          
          <div v-if="uploadSuccess" class="upload-success">
            <v-icon size="24" color="success">mdi-check-circle</v-icon>
            <span>上传成功！</span>
          </div>
          
          <div v-if="uploadError" class="upload-error">
            <v-icon size="24" color="error">mdi-alert-circle</v-icon>
            <span>{{ uploadError }}</span>
          </div>
        </div>
        
        <div class="upload-modal-footer">
          <v-btn variant="text" @click="closeUploadModal">取消</v-btn>
          <v-btn color="primary" :loading="isUploading" :disabled="!hasFiles || isUploading" @click="startUpload">
            {{ isUploading ? '上传中...' : '开始上传' }}
          </v-btn>
        </div>
      </div>
    </div>

    <div v-if="linkModalVisible" class="modal-overlay" @click.self="linkModalVisible = false">
      <div class="link-modal">
        <div class="modal-header">
          <h3>插入链接</h3>
          <v-btn icon variant="text" size="small" @click="linkModalVisible = false">
            <v-icon>mdi-close</v-icon>
          </v-btn>
        </div>
        <div class="modal-body">
          <v-text-field v-model="linkText" label="链接文字" variant="outlined" class="mb-3" />
          <v-text-field v-model="linkUrl" label="链接地址" variant="outlined" placeholder="https://" />
        </div>
        <div class="modal-footer">
          <v-btn variant="text" @click="linkModalVisible = false">取消</v-btn>
          <v-btn color="primary" @click="confirmLink">确定</v-btn>
        </div>
      </div>
    </div>

    <div v-if="tableModalVisible" class="modal-overlay" @click.self="tableModalVisible = false">
      <div class="table-modal">
        <div class="modal-header">
          <h3>插入表格</h3>
          <v-btn icon variant="text" size="small" @click="tableModalVisible = false">
            <v-icon>mdi-close</v-icon>
          </v-btn>
        </div>
        <div class="modal-body">
          <v-text-field v-model.number="tableRows" label="行数" variant="outlined" type="number" min="1" max="10" class="mb-3" />
          <v-text-field v-model.number="tableCols" label="列数" variant="outlined" type="number" min="1" max="10" />
        </div>
        <div class="modal-footer">
          <v-btn variant="text" @click="tableModalVisible = false">取消</v-btn>
          <v-btn color="primary" @click="confirmTable">确定</v-btn>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue'
import MarkdownIt from 'markdown-it'
import anchor from 'markdown-it-anchor'
import container from 'markdown-it-container'
import * as emoji from 'markdown-it-emoji'
import footnote from 'markdown-it-footnote'
import mark from 'markdown-it-mark'
import taskLists from 'markdown-it-task-lists'
import katex from 'katex'
import 'katex/dist/katex.min.css'
import api from '../api'

const katexInline = (str) => {
  try {
    return katex.renderToString(str, {
      displayMode: false,
      throwOnError: false
    })
  } catch (e) {
    return str
  }
}

const katexBlock = (str) => {
  try {
    return '<div class="katex-block">' + katex.renderToString(str, {
      displayMode: true,
      throwOnError: false
    }) + '</div>'
  } catch (e) {
    return '<pre>' + str + '</pre>'
  }
}

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  height: {
    type: Number,
    default: 500
  },
  placeholder: {
    type: String,
    default: '输入内容...（支持 Markdown 语法）'
  }
})

const emit = defineEmits(['update:modelValue'])

const textareaRef = ref(null)
const fileInputRef = ref(null)
const previewPaneRef = ref(null)
const content = ref(props.modelValue)
const viewMode = ref('split')
const linkModalVisible = ref(false)
const tableModalVisible = ref(false)
const linkText = ref('')
const linkUrl = ref('')
const tableRows = ref(3)
const tableCols = ref(3)
const isScrolling = ref(false)

const uploadModalVisible = ref(false)
const uploadType = ref('image')
const uploadFiles = ref([])
const isUploading = ref(false)
const uploadProgress = ref(0)
const uploadSuccess = ref(false)
const uploadError = ref('')

const suggestions = ref([])
const selectedSuggestion = ref(0)
const suggestionPosition = ref({ top: 0, left: 0 })

// 语音识别相关
const isRecording = ref(false)
const recognition = ref(null)
const recognitionSupported = ref(false)

// 初始化语音识别
const initRecognition = () => {
  if ('webkitSpeechRecognition' in window || 'SpeechRecognition' in window) {
    const SpeechRecognition = window.SpeechRecognition || window.webkitSpeechRecognition
    recognition.value = new SpeechRecognition()
    recognition.value.lang = 'zh-CN' // 设置为中文
    recognition.value.interimResults = true // 显示中间结果
    recognition.value.continuous = false // 不连续识别，说完一句话就停止
    
    recognition.value.onstart = () => {
      isRecording.value = true
    }
    
    recognition.value.onresult = (event) => {
      let finalTranscript = ''
      let interimTranscript = ''
      
      for (let i = event.resultIndex; i < event.results.length; i++) {
        const transcript = event.results[i][0].transcript
        if (event.results[i].isFinal) {
          finalTranscript += transcript
        } else {
          interimTranscript += transcript
        }
      }
      
      // 插入识别结果到编辑器
      if (finalTranscript) {
        insertAtCursor(finalTranscript)
      }
    }
    
    recognition.value.onerror = (event) => {
      console.error('语音识别错误:', event.error)
      isRecording.value = false
      if (event.error === 'no-speech') {
        // 没有检测到语音，不显示错误
      } else if (event.error === 'not-allowed') {
        alert('请允许麦克风访问权限')
      }
    }
    
    recognition.value.onend = () => {
      isRecording.value = false
    }
    
    recognitionSupported.value = true
  }
}

// 初始化时调用
initRecognition()

const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
  breaks: true,
  highlight: function (str, lang) {
    if (lang && window.hljs) {
      try {
        const language = window.hljs.getLanguage(lang) ? lang : 'plaintext'
        return '<pre class="hljs-code code-block" data-language="' + language + '"><code class="hljs language-' + language + '">' +
          window.hljs.highlight(str, { language: language, ignoreIllegals: true }).value +
          '</code></pre>'
      } catch (e) {}
    }
    return '<pre class="hljs-code code-block"><code class="hljs">' + MarkdownIt().utils.escapeHtml(str) + '</code></pre>'
  }
})

md.use(anchor, {
  permalink: anchor.permalink.ariaHidden({
    placement: 'after',
    class: 'header-anchor',
    symbol: '#'
  }),
  level: [1, 2, 3, 4]
})

md.use(container, 'tip', {
  render: function (tokens, idx) {
    if (tokens[idx].nesting === 1) {
      return '<div class="custom-container tip"><p class="custom-container-title">💡 提示</p>'
    } else {
      return '</div>'
    }
  }
})

md.use(container, 'warning', {
  render: function (tokens, idx) {
    if (tokens[idx].nesting === 1) {
      return '<div class="custom-container warning"><p class="custom-container-title">⚠️ 警告</p>'
    } else {
      return '</div>'
    }
  }
})

md.use(container, 'danger', {
  render: function (tokens, idx) {
    if (tokens[idx].nesting === 1) {
      return '<div class="custom-container danger"><p class="custom-container-title">❌ 危险</p>'
    } else {
      return '</div>'
    }
  }
})

md.use(container, 'info', {
  render: function (tokens, idx) {
    if (tokens[idx].nesting === 1) {
      return '<div class="custom-container info"><p class="custom-container-title">ℹ️ 信息</p>'
    } else {
      return '</div>'
    }
  }
})

md.use(emoji)
md.use(footnote)
md.use(mark)
md.use(taskLists, { enabled: true, label: true })

const renderMath = (text) => {
  let result = text

  result = result.replace(/\$\$([\s\S]*?)\$\$/g, (match, formula) => {
    return katexBlock(formula.trim())
  })

  result = result.replace(/\$(.+?)\$/g, (match, formula) => {
    return katexInline(formula.trim())
  })

  return result
}

const renderedContent = computed(() => {
  if (!content.value) return '<p class="placeholder-text">预览区域</p>'
  const rendered = md.render(content.value)
  return renderMath(rendered)
})

const stats = computed(() => {
  const text = content.value
  
  const words = text.length
  const lines = text.split('\n').length
  
  const paragraphs = text.split(/\n\n+/).filter(p => p.trim()).length
  
  const tokenPattern = /(#+|```|`|\*\*|\*|~~|>|\[|\]|\(|)|\d+\./g
  const tokens = (text.match(tokenPattern) || []).length
  
  return { words, paragraphs, lines, tokens }
})

const suggestionStyle = computed(() => ({
  top: suggestionPosition.value.top + 'px',
  left: suggestionPosition.value.left + 'px'
}))

const hasFiles = computed(() => uploadFiles.value.length > 0)

watch(() => props.modelValue, (newVal) => {
  if (newVal !== content.value) {
    content.value = newVal
  }
})

const handleInput = () => {
  emit('update:modelValue', content.value)
  triggerSuggestions()
  autoCompleteList()
}

const autoCompleteList = () => {
  const textarea = textareaRef.value
  if (!textarea) return

  const start = textarea.selectionStart
  const beforeCursor = content.value.substring(0, start)
  
  const bulletMatch = beforeCursor.match(/(\n|^)-\s*$/)
  const numberMatch = beforeCursor.match(/(\n|^)(\d+)\.\s*$/)
  
  if (bulletMatch && !bulletMatch[0].endsWith(' ')) {
    const insertPos = start
    content.value = content.value.substring(0, insertPos) + ' ' + content.value.substring(insertPos)
    nextTick(() => {
      textarea.setSelectionRange(insertPos + 1, insertPos + 1)
    })
  } else if (numberMatch && !numberMatch[0].endsWith(' ')) {
    const insertPos = start
    content.value = content.value.substring(0, insertPos) + ' ' + content.value.substring(insertPos)
    nextTick(() => {
      textarea.setSelectionRange(insertPos + 1, insertPos + 1)
    })
  }
}

const handleKeydown = (e) => {
  if (e.ctrlKey || e.metaKey) {
    switch (e.key.toLowerCase()) {
      case 'b':
        e.preventDefault()
        insertFormat('**', '**')
        break
      case 'i':
        e.preventDefault()
        insertFormat('*', '*')
        break
      case 'k':
        e.preventDefault()
        insertLink()
        break
      case 'z':
        e.preventDefault()
        formatDoc('undo')
        break
      case 'y':
        e.preventDefault()
        formatDoc('redo')
        break
      case 's':
        e.preventDefault()
        break
    }
  }

  if (e.key === 'Tab') {
    e.preventDefault()
    if (suggestions.value.length > 0) {
      applySuggestion(suggestions.value[selectedSuggestion.value])
    } else {
      insertAtCursor('  ')
    }
  }

  if (e.key === 'Enter') {
    if (suggestions.value.length > 0) {
      e.preventDefault()
      applySuggestion(suggestions.value[selectedSuggestion.value])
      return
    }
    
    const textarea = textareaRef.value
    if (!textarea) return

    const start = textarea.selectionStart
    const line = getCurrentLine()

    if (line.match(/^(\d+)\. $/)) {
      const match = line.match(/^(\d+)\. $/)
      const num = parseInt(match[1]) + 1
      insertAtCursor(`\n${num}. `)
      e.preventDefault()
    } else if (line.match(/^- $/)) {
      insertAtCursor('\n- ')
      e.preventDefault()
    } else if (line.match(/^- \[ \] $/)) {
      insertAtCursor('\n- [ ] ')
      e.preventDefault()
    } else if (line.match(/^> $/)) {
      insertAtCursor('\n> ')
      e.preventDefault()
    } else if (line.match(/^#+$/)) {
      insertAtCursor('\n')
      e.preventDefault()
    }
  }

  if (e.key === 'ArrowDown') {
    if (suggestions.value.length > 0) {
      e.preventDefault()
      selectedSuggestion.value = (selectedSuggestion.value + 1) % suggestions.value.length
    }
  }

  if (e.key === 'ArrowUp') {
    if (suggestions.value.length > 0) {
      e.preventDefault()
      selectedSuggestion.value = (selectedSuggestion.value - 1 + suggestions.value.length) % suggestions.value.length
    }
  }

  if (e.key === 'Escape') {
    suggestions.value = []
  }
}

const getCurrentLine = () => {
  const textarea = textareaRef.value
  if (!textarea) return ''

  const start = textarea.selectionStart
  const lineStart = content.value.lastIndexOf('\n', start - 1) + 1
  const lineEnd = content.value.indexOf('\n', start)
  return content.value.substring(lineStart, lineEnd === -1 ? content.value.length : lineEnd)
}

const formatDoc = (cmd) => {
  const textarea = textareaRef.value
  if (!textarea) return
  textarea.focus()
  document.execCommand(cmd, false, null)
}

const insertFormat = (prefix, suffix) => {
  const textarea = textareaRef.value
  if (!textarea) return

  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const selectedText = content.value.substring(start, end)

  const newText = content.value.substring(0, start) + prefix + selectedText + suffix + content.value.substring(end)
  content.value = newText

  nextTick(() => {
    textarea.focus()
    if (selectedText) {
      textarea.setSelectionRange(start + prefix.length, end + prefix.length)
    } else {
      textarea.setSelectionRange(start + prefix.length, start + prefix.length)
    }
  })

  handleInput()
}

const insertLine = (prefix) => {
  const textarea = textareaRef.value
  if (!textarea) return

  const start = textarea.selectionStart
  const lineStart = content.value.lastIndexOf('\n', start - 1) + 1

  const newText = content.value.substring(0, lineStart) + prefix + content.value.substring(lineStart)
  content.value = newText

  nextTick(() => {
    textarea.focus()
    const newPos = lineStart + prefix.length
    textarea.setSelectionRange(newPos, newPos)
  })

  handleInput()
}

const insertCodeBlock = () => {
  const textarea = textareaRef.value
  if (!textarea) return

  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const selectedText = content.value.substring(start, end) || ''

  const codeBlock = `\n\`\`\`\n${selectedText}\n\`\`\`\n`
  const newText = content.value.substring(0, start) + codeBlock + content.value.substring(end)
  content.value = newText

  nextTick(() => {
    textarea.focus()
    const cursorPos = start + 4
    textarea.setSelectionRange(cursorPos, cursorPos + selectedText.length)
  })

  handleInput()
}

const insertLink = () => {
  const textarea = textareaRef.value
  if (!textarea) return

  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const selectedText = content.value.substring(start, end) || '链接文字'

  linkText.value = selectedText
  linkUrl.value = ''
  linkModalVisible.value = true

  nextTick(() => {
    textarea.focus()
  })
}

const confirmLink = () => {
  if (!linkUrl.value) return

  const textarea = textareaRef.value
  if (!textarea) return

  const start = textarea.selectionStart
  const linkMarkdown = `[${linkText.value}](${linkUrl.value})`
  const newText = content.value.substring(0, start) + linkMarkdown + content.value.substring(start)
  content.value = newText

  linkModalVisible.value = false

  nextTick(() => {
    textarea.focus()
    const newPos = start + linkMarkdown.length
    textarea.setSelectionRange(newPos, newPos)
  })

  handleInput()
}

const insertImage = () => {
  uploadType.value = 'image'
  uploadFiles.value = []
  uploadSuccess.value = false
  uploadError.value = ''
  uploadModalVisible.value = true
}

const insertVideo = () => {
  uploadType.value = 'video'
  uploadFiles.value = []
  uploadSuccess.value = false
  uploadError.value = ''
  uploadModalVisible.value = true
}

const insertTable = () => {
  tableRows.value = 3
  tableCols.value = 3
  tableModalVisible.value = true
}

// 语音识别控制函数
const toggleRecording = () => {
  if (!recognitionSupported.value) {
    alert('您的浏览器不支持语音识别功能，请使用Chrome浏览器')
    return
  }
  
  if (isRecording.value) {
    // 停止录音
    recognition.value?.stop()
  } else {
    // 开始录音
    try {
      recognition.value?.start()
    } catch (error) {
      console.error('启动语音识别失败:', error)
      isRecording.value = false
    }
  }
}

const confirmTable = () => {
  const table = generateTable(tableRows.value, tableCols.value)
  insertAtCursor('\n' + table + '\n')
  tableModalVisible.value = false
}

const generateTable = (rows, cols) => {
  const header = '| ' + Array(cols).fill('列').map((h, i) => `${h}${i + 1}`).join(' | ') + ' |'
  const separator = '| ' + Array(cols).fill('---').join(' | ') + ' |'
  const body = Array(rows - 1).fill(null).map(() => '| ' + Array(cols).fill('').join(' | ') + ' |').join('\n')
  return header + '\n' + separator + '\n' + body
}

const insertHr = () => {
  insertAtCursor('\n---\n\n')
}

const insertMath = () => {
  const textarea = textareaRef.value
  if (!textarea) return

  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const selectedText = content.value.substring(start, end) || 'x^2 + y^2 = z^2'

  const mathBlock = `\n\$\$\n${selectedText}\n\$\$\n`
  const newText = content.value.substring(0, start) + mathBlock + content.value.substring(end)
  content.value = newText

  nextTick(() => {
    textarea.focus()
    const cursorPos = start + 3
    textarea.setSelectionRange(cursorPos, cursorPos + selectedText.length)
  })

  handleInput()
}

const insertMark = () => {
  insertFormat('==', '==')
}

const insertFootnote = () => {
  const textarea = textareaRef.value
  if (!textarea) return

  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const selectedText = content.value.substring(start, end) || '脚注文本'

  const footnoteId = 'fn-' + Date.now()
  const footnoteRef = `[^${footnoteId}]`
  const footnoteDef = `[^${footnoteId}]: ${selectedText}`

  const newText = content.value.substring(0, start) + footnoteRef + content.value.substring(end)
  content.value = newText

  nextTick(() => {
    textarea.focus()
    content.value = content.value + '\n\n' + footnoteDef
    const newPos = start + footnoteRef.length
    textarea.setSelectionRange(newPos, newPos)
  })

  handleInput()
}

const insertContainer = (type) => {
  const textarea = textareaRef.value
  if (!textarea) return

  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const selectedText = content.value.substring(start, end) || '这里放置内容'

  const containerBlock = `\n:::${type}\n${selectedText}\n:::\n`
  const newText = content.value.substring(0, start) + containerBlock + content.value.substring(end)
  content.value = newText

  nextTick(() => {
    textarea.focus()
    const cursorPos = start + type.length + 4
    textarea.setSelectionRange(cursorPos, cursorPos + selectedText.length)
  })

  handleInput()
}

const insertAtCursor = (text) => {
  const textarea = textareaRef.value
  if (!textarea) return

  const start = textarea.selectionStart
  const end = textarea.selectionEnd

  const newText = content.value.substring(0, start) + text + content.value.substring(end)
  content.value = newText

  nextTick(() => {
    textarea.focus()
    const newPos = start + text.length
    textarea.setSelectionRange(newPos, newPos)
  })

  handleInput()
}

const triggerSuggestions = () => {
  const textarea = textareaRef.value
  if (!textarea) return

  const pos = textarea.selectionStart
  const beforeCursor = content.value.substring(0, pos)
  const lastSpace = beforeCursor.lastIndexOf(' ')
  const lastNewline = beforeCursor.lastIndexOf('\n')
  const startIndex = Math.max(lastSpace, lastNewline) + 1
  const word = beforeCursor.substring(startIndex)

  if (word.length > 0) {
    const matchedSuggestions = getSuggestions(word)
    if (matchedSuggestions.length > 0) {
      suggestions.value = matchedSuggestions
      selectedSuggestion.value = 0
      updateSuggestionPosition()
      return
    }
  }

  suggestions.value = []
}

const getSuggestions = (query) => {
  const markdownSuggestions = [
    { text: '**粗体**', detail: '加粗文本', icon: 'mdi-format-bold', color: '#6750A4', insert: '**|**' },
    { text: '*斜体*', detail: '倾斜文本', icon: 'mdi-format-italic', color: '#6750A4', insert: '*|*' },
    { text: '~~删除线~~', detail: '删除文本', icon: 'mdi-format-strikethrough', color: '#6750A4', insert: '~~|~~' },
    { text: '`代码`', detail: '行内代码', icon: 'mdi-code-tags', color: '#EC4899', insert: '`|`' },
    { text: '==高亮==', detail: '高亮文本', icon: 'mdi-marker', color: '#FFC107', insert: '==|==' },
    { text: '# 标题', detail: '一级标题', icon: 'mdi-format-header-1', color: '#1976D2', insert: '# |' },
    { text: '## 标题', detail: '二级标题', icon: 'mdi-format-header-2', color: '#1976D2', insert: '## |' },
    { text: '### 标题', detail: '三级标题', icon: 'mdi-format-header-3', color: '#1976D2', insert: '### |' },
    { text: '- 列表项', detail: '无序列表', icon: 'mdi-format-list-bulleted', color: '#43A047', insert: '- |' },
    { text: '1. 列表项', detail: '有序列表', icon: 'mdi-format-list-numbered', color: '#43A047', insert: '1. |' },
    { text: '- [ ] 任务', detail: '待办任务', icon: 'mdi-checkbox-marked-outline', color: '#FB8C00', insert: '- [ ] |' },
    { text: '> 引用', detail: '引用文本', icon: 'mdi-format-quote-close', color: '#7E57C2', insert: '> |' },
    { text: '[链接](url)', detail: '插入链接', icon: 'mdi-link', color: '#0288D1', insert: '[|]()' },
    { text: '![图片](url)', detail: '插入图片', icon: 'mdi-image', color: '#00ACC1', insert: '![](|)' },
    { text: '```代码块', detail: '代码块', icon: 'mdi-code-braces', color: '#EC407A', insert: '```\n|\n```' },
    { text: '---', detail: '分割线', icon: 'mdi-minus', color: '#9E9E9E', insert: '\n---\n' },
    { text: ':::tip', detail: '提示框', icon: 'mdi-lightbulb', color: '#4CAF50', insert: ':::tip\n|\n:::' },
    { text: ':::warning', detail: '警告框', icon: 'mdi-alert', color: '#FF9800', insert: ':::warning\n|\n:::' },
    { text: ':::danger', detail: '危险框', icon: 'mdi-close-circle', color: '#f44336', insert: ':::danger\n|\n:::' },
    { text: ':::info', detail: '信息框', icon: 'mdi-information', color: '#2196F3', insert: ':::info\n|\n:::' },
    { text: '$$公式$$', detail: '数学公式', icon: 'mdi-function-variant', color: '#9C27B0', insert: '$$\n|\n$$' },
  ]

  return markdownSuggestions.filter(s =>
    s.text.toLowerCase().includes(query.toLowerCase()) ||
    s.detail.toLowerCase().includes(query.toLowerCase())
  )
}

const updateSuggestionPosition = () => {
  const textarea = textareaRef.value
  if (!textarea) return

  const pos = textarea.selectionStart
  const textBeforeCursor = content.value.substring(0, pos)
  const lines = textBeforeCursor.split('\n')
  const currentLineIndex = lines.length - 1
  const currentLineText = lines[currentLineIndex]

  const lineHeight = 24
  const charWidth = 8.4

  const top = (currentLineIndex + 1) * lineHeight + 16
  const left = currentLineText.length * charWidth + 16

  suggestionPosition.value = { top, left }
}

const applySuggestion = (suggestion) => {
  const textarea = textareaRef.value
  if (!textarea) return

  const pos = textarea.selectionStart
  const beforeCursor = content.value.substring(0, pos)
  const lastSpace = beforeCursor.lastIndexOf(' ')
  const lastNewline = beforeCursor.lastIndexOf('\n')
  const startIndex = Math.max(lastSpace, lastNewline) + 1

  const newText = content.value.substring(0, startIndex) + suggestion.insert + content.value.substring(pos)
  content.value = newText

  suggestions.value = []

  nextTick(() => {
    textarea.focus()
    const cursorPos = startIndex + suggestion.insert.indexOf('|')
    textarea.setSelectionRange(cursorPos, cursorPos)
  })

  handleInput()
}

const closeUploadModal = () => {
  uploadModalVisible.value = false
  uploadFiles.value = []
  uploadProgress.value = 0
  uploadSuccess.value = false
  uploadError.value = ''
}

const triggerFileInput = () => {
  fileInputRef.value?.click()
}

const handleFileChange = (e) => {
  const files = Array.from(e.target.files)
  addFiles(files)
}

const onDragOver = (e) => {
  e.currentTarget.classList.add('drag-over')
}

const onDragLeave = (e) => {
  e.currentTarget.classList.remove('drag-over')
}

const onDrop = (e) => {
  e.currentTarget.classList.remove('drag-over')
  const files = Array.from(e.dataTransfer.files)
  addFiles(files)
}

const addFiles = (files) => {
  files.forEach(file => {
    const reader = new FileReader()
    reader.onload = (e) => {
      uploadFiles.value.push({
        file,
        name: file.name,
        size: file.size,
        preview: e.target.result
      })
    }
    reader.readAsDataURL(file)
  })
}

const removeFile = (index) => {
  uploadFiles.value.splice(index, 1)
}

const formatFileSize = (bytes) => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}

const startUpload = async () => {
  if (uploadFiles.value.length === 0) return

  isUploading.value = true
  uploadProgress.value = 0
  uploadSuccess.value = false
  uploadError.value = ''

  const isImage = uploadType.value === 'image'
  const uploadTypeStr = isImage ? 'image' : 'video'

  try {
    for (let i = 0; i < uploadFiles.value.length; i++) {
      const item = uploadFiles.value[i]
      const formData = new FormData()
      formData.append(isImage ? 'image' : 'video', item.file)

      const response = await api.post(`/upload/${uploadTypeStr}`, formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
        onUploadProgress: (progressEvent) => {
          const percent = Math.round((progressEvent.loaded * 100) / progressEvent.total)
          uploadProgress.value = (percent / uploadFiles.value.length) + (i / uploadFiles.value.length) * 100
        }
      })

      const fileData = {
        name: item.name,
        url: response.data.url
      }

      if (isImage) {
        insertAtCursor(`\n![${item.name}](${fileData.url})\n`)
      } else {
        insertAtCursor(`\n<video src="${fileData.url}" controls></video>\n`)
      }
    }

    uploadSuccess.value = true
    setTimeout(() => {
      closeUploadModal()
    }, 1500)
  } catch (error) {
    uploadError.value = error.response?.data?.message || '上传失败，请重试'
  } finally {
    isUploading.value = false
  }
}

const syncScroll = () => {
  if (isScrolling.value || viewMode.value !== 'split') return

  const textarea = textareaRef.value
  const preview = previewPaneRef.value
  if (!textarea || !preview) return

  isScrolling.value = true
  const scrollRatio = textarea.scrollTop / (textarea.scrollHeight - textarea.clientHeight)
  preview.scrollTop = scrollRatio * (preview.scrollHeight - preview.clientHeight)

  setTimeout(() => { isScrolling.value = false }, 50)
}

const handlePreviewScroll = () => {
  if (isScrolling.value || viewMode.value !== 'split') return

  const textarea = textareaRef.value
  const preview = previewPaneRef.value
  if (!textarea || !preview) return

  isScrolling.value = true
  const scrollRatio = preview.scrollTop / (preview.scrollHeight - preview.clientHeight)
  textarea.scrollTop = scrollRatio * (textarea.scrollHeight - textarea.clientHeight)

  setTimeout(() => { isScrolling.value = false }, 50)
}
</script>

<style scoped>
.md-editor {
  width: 100%;
  max-width: 100%;
  max-height: 80vh;
  border: 1px solid #E5E5E5;
  border-radius: 12px;
  overflow: hidden;
  background: #fff;
  display: flex;
  flex-direction: column;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.toolbar {
  display: flex;
  align-items: center;
  padding: 8px 16px;
  background: #FAFAFA;
  border-bottom: 1px solid #E5E5E5;
  flex-wrap: wrap;
  gap: 4px;
}

.toolbar-group {
  display: flex;
  align-items: center;
  gap: 2px;
}

.status-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 12px;
  background: #F5F5F5;
  border-radius: 6px;
  font-size: 12px;
  color: #666;
}

.status-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.status-divider {
  color: #DDD;
}

.editor-body {
  display: flex;
  flex: 1;
  min-height: v-bind('props.height + "px"');
}

.editor-body.view-edit .editor-pane {
  flex: 1;
}

.editor-body.view-preview .preview-pane {
  flex: 1;
}

.editor-body.view-split .editor-pane,
.editor-body.view-split .preview-pane {
  flex: 1;
}

.editor-pane {
  flex: 1;
  border-right: 1px solid #E5E5E5;
  display: flex;
  position: relative;
}

.editor-textarea {
  width: 100%;
  min-height: v-bind('props.height + "px"');
  padding: 16px;
  border: none;
  outline: none;
  resize: none;
  font-family: 'JetBrains Mono', 'Monaco', 'Consolas', monospace;
  font-size: 14px;
  line-height: 1.7;
  background: #fff;
  tab-size: 2;
  color: #333;
}

.editor-textarea::placeholder {
  color: #999;
}

.suggestion-box {
  position: absolute;
  background: #fff;
  border: 1px solid #E5E5E5;
  border-radius: 8px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
  min-width: 240px;
  max-height: 300px;
  overflow-y: auto;
  z-index: 1000;
}

.suggestion-item {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  cursor: pointer;
  gap: 8px;
  transition: background 0.15s;
}

.suggestion-item:hover,
.suggestion-item.active {
  background: #F5F5F5;
}

.suggestion-text {
  flex: 1;
  font-family: 'JetBrains Mono', monospace;
  font-size: 13px;
  color: #333;
}

.suggestion-detail {
  font-size: 12px;
  color: #999;
}

.preview-pane {
  flex: 1;
  overflow-y: auto;
  background: #FAFAFA;
}

.preview-content {
  padding: 16px;
  min-height: v-bind('props.height + "px"');
}

.markdown-body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  line-height: 1.8;
  color: #333;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.markdown-body :deep(h1),
.markdown-body :deep(h2),
.markdown-body :deep(h3),
.markdown-body :deep(h4) {
  margin: 20px 0 12px 0;
  font-weight: 600;
  line-height: 1.4;
  color: #1A1A1A;
}

.markdown-body :deep(h1) { font-size: 1.75rem; border-bottom: 1px solid #E5E5E5; padding-bottom: 8px; }
.markdown-body :deep(h2) { font-size: 1.5rem; }
.markdown-body :deep(h3) { font-size: 1.25rem; }

.markdown-body :deep(p) {
  margin: 0 0 12px 0;
}

.markdown-body :deep(code) {
  background: #f0f0f0;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'JetBrains Mono', 'Monaco', monospace;
  font-size: 13px;
}

.markdown-body :deep(pre) {
  background: #1a1a1a;
  border-radius: 8px;
  padding: 16px;
  overflow-x: auto;
  margin: 12px 0;
}

.markdown-body :deep(pre code) {
  background: none;
  padding: 0;
  color: #e0e0e0;
}

.markdown-body :deep(blockquote) {
  border-left: 4px solid #6750A4;
  margin: 12px 0;
  padding: 8px 16px;
  background: #f8f8f8;
  border-radius: 0 8px 8px 0;
  color: #666;
}

.markdown-body :deep(ul),
.markdown-body :deep(ol) {
  padding-left: 24px;
  margin: 12px 0;
}

.markdown-body :deep(li) {
  white-space: pre-wrap;
}

.markdown-body :deep(img) {
  max-width: 100%;
  border-radius: 8px;
  margin: 12px 0;
  display: block;
}

.markdown-body :deep(video) {
  max-width: 100%;
  border-radius: 8px;
  margin: 12px 0;
  display: block;
  background: #000;
}

.markdown-body :deep(audio) {
  width: 100%;
  max-width: 100%;
  margin: 16px 0;
  display: block;
}

.markdown-body :deep(table) {
  border-collapse: collapse;
  width: 100%;
  margin: 12px 0;
}

.markdown-body :deep(th),
.markdown-body :deep(td) {
  border: 1px solid #e0e0e0;
  padding: 8px 12px;
  text-align: left;
}

.markdown-body :deep(th) {
  background: #f5f5f5;
  font-weight: 600;
}

.markdown-body :deep(hr) {
  border: none;
  border-top: 2px solid #e0e0e0;
  margin: 24px 0;
}

.markdown-body :deep(a) {
  color: #6750A4;
  text-decoration: none;
}

.markdown-body :deep(a:hover) {
  text-decoration: underline;
}

.markdown-body :deep(.placeholder-text) {
  color: #999;
  text-align: center;
  padding-top: 40px;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
}

.link-modal,
.table-modal {
  background: #fff;
  border-radius: 12px;
  width: 420px;
  max-width: 90vw;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid #E5E5E5;
}

.modal-header h3 {
  margin: 0;
  font-size: 1.1rem;
  font-weight: 600;
}

.modal-body {
  padding: 20px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 16px 20px;
  border-top: 1px solid #E5E5E5;
}

.upload-modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
  backdrop-filter: blur(4px);
}

.upload-modal-container {
  background: #fff;
  border-radius: 16px;
  width: 520px;
  max-width: 90vw;
  max-height: 85vh;
  overflow: hidden;
  box-shadow: 0 12px 48px rgba(0, 0, 0, 0.2);
}

.upload-modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid #E5E5E5;
  background: #FAFAFA;
}

.modal-title {
  display: flex;
  align-items: center;
  gap: 10px;
}

.modal-title h3 {
  margin: 0;
  font-size: 1.2rem;
  font-weight: 600;
  color: #1A1A1A;
}

.upload-modal-body {
  padding: 24px;
  overflow-y: auto;
  max-height: 50vh;
}

.upload-area {
  border: 2px dashed #DDD;
  border-radius: 12px;
  padding: 32px;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s;
  background: #FAFAFA;
}

.upload-area:hover {
  border-color: #6750A4;
  background: #F5F3FF;
}

.upload-area.drag-over {
  border-color: #6750A4;
  background: #F5F3FF;
  transform: scale(1.02);
}

.upload-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.upload-hint {
  margin: 0;
  font-size: 14px;
  color: #333;
  font-weight: 500;
}

.upload-format {
  margin: 0;
  font-size: 12px;
  color: #999;
}

.upload-preview {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.file-preview-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #fff;
  border-radius: 8px;
  border: 1px solid #E5E5E5;
}

.preview-thumb {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  overflow: hidden;
  background: #F5F5F5;
  flex-shrink: 0;
}

.preview-thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.video-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.file-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.file-name {
  font-size: 13px;
  color: #333;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-size {
  font-size: 12px;
  color: #999;
}

.file-input {
  display: none;
}

.upload-progress {
  margin-top: 16px;
}

.progress-info {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
  font-size: 13px;
  color: #666;
}

.upload-success {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-top: 16px;
  padding: 12px;
  background: #E8F5E9;
  border-radius: 8px;
  color: #2E7D32;
}

.upload-error {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-top: 16px;
  padding: 12px;
  background: #FFEBEE;
  border-radius: 8px;
  color: #C62828;
}

.upload-modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 20px 24px;
  border-top: 1px solid #E5E5E5;
  background: #FAFAFA;
}

.container-menu {
  min-width: 140px;
  border-radius: 8px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
}

.markdown-body :deep(.header-anchor) {
  position: absolute;
  left: -24px;
  color: #999;
  text-decoration: none;
  opacity: 0;
  transition: opacity 0.2s;
}

.markdown-body :deep(h1:hover .header-anchor),
.markdown-body :deep(h2:hover .header-anchor),
.markdown-body :deep(h3:hover .header-anchor),
.markdown-body :deep(h4:hover .header-anchor) {
  opacity: 1;
}

.markdown-body :deep(mark) {
  background: linear-gradient(120deg, #ffeaa7 0%, #fdcb6e 100%);
  padding: 2px 4px;
  border-radius: 3px;
}

.markdown-body :deep(.footnotes) {
  margin-top: 40px;
  padding-top: 20px;
  border-top: 2px solid #e5e5e5;
  font-size: 14px;
  color: #666;
}

.markdown-body :deep(.footnotes-list) {
  padding-left: 20px;
}

.markdown-body :deep(.footnote-item) {
  margin: 12px 0;
}

.markdown-body :deep(.footnote-ref) {
  color: #6750A4;
  font-weight: bold;
  text-decoration: none;
}

.markdown-body :deep(.custom-container) {
  margin: 16px 0;
  padding: 14px 18px;
  border-radius: 8px;
  border-left: 4px solid;
}

.markdown-body :deep(.custom-container p) {
  margin: 0;
  font-weight: 600;
}

.markdown-body :deep(.custom-container.tip) {
  background: linear-gradient(135deg, #e8f5e9 0%, #c8e6c9 100%);
  border-color: #4CAF50;
}

.markdown-body :deep(.custom-container.warning) {
  background: linear-gradient(135deg, #fff3e0 0%, #ffe0b2 100%);
  border-color: #FF9800;
}

.markdown-body :deep(.custom-container.danger) {
  background: linear-gradient(135deg, #ffebee 0%, #ffcdd2 100%);
  border-color: #f44336;
}

.markdown-body :deep(.custom-container.info) {
  background: linear-gradient(135deg, #e3f2fd 0%, #bbdefb 100%);
  border-color: #2196F3;
}

.markdown-body :deep(.task-list-item) {
  list-style: none;
  margin-left: -24px;
}

.markdown-body :deep(.task-list-item input[type="checkbox"]) {
  margin-right: 8px;
  width: 16px;
  height: 16px;
  accent-color: #6750A4;
}

.markdown-body :deep(.emoji) {
  display: inline;
  vertical-align: middle;
}

.markdown-body :deep(kbd) {
  display: inline-block;
  padding: 2px 8px;
  font-size: 0.9em;
  font-family: 'JetBrains Mono', monospace;
  color: #333;
  background: #f5f5f5;
  border: 1px solid #ddd;
  border-radius: 4px;
  box-shadow: 0 2px 0 #ccc;
}

.markdown-body :deep(.code-block) {
  position: relative;
  border-radius: 8px;
  overflow: hidden;
  margin: 16px 0;
}

.markdown-body :deep(.code-block::before) {
  content: attr(data-language);
  position: absolute;
  top: 0;
  right: 0;
  padding: 4px 12px;
  font-size: 11px;
  color: #888;
  text-transform: uppercase;
  background: rgba(0, 0, 0, 0.2);
  font-weight: 500;
}

.markdown-body :deep(.code-block pre) {
  margin: 0;
  border-radius: 0;
  background: #0d1117;
}

.markdown-body :deep(.code-block code) {
  display: block;
  padding: 16px;
  overflow-x: auto;
  font-family: 'JetBrains Mono', 'Fira Code', 'Consolas', monospace;
  font-size: 13px;
  line-height: 1.6;
  white-space: pre;
}

.markdown-body :deep(.hljs) {
  background: #0d1117 !important;
  color: #c9d1d9;
}

.markdown-body :deep(.hljs-keyword),
.markdown-body :deep(.hljs-selector-tag),
.markdown-body :deep(.hljs-built_in),
.markdown-body :deep(.hljs-name),
.markdown-body :deep(.hljs-tag) {
  color: #ff7b72;
}

.markdown-body :deep(.hljs-string),
.markdown-body :deep(.hljs-title),
.markdown-body :deep(.hljs-section),
.markdown-body :deep(.hljs-attribute),
.markdown-body :deep(.hljs-literal),
.markdown-body :deep(.hljs-template-tag),
.markdown-body :deep(.hljs-template-variable),
.markdown-body :deep(.hljs-type),
.markdown-body :deep(.hljs-addition) {
  color: #a5d6ff;
}

.markdown-body :deep(.hljs-comment),
.markdown-body :deep(.hljs-quote),
.markdown-body :deep(.hljs-deletion),
.markdown-body :deep(.hljs-meta) {
  color: #8b949e;
}

.markdown-body :deep(.hljs-number),
.markdown-body :deep(.hljs-regexp),
.markdown-body :deep(.hljs-literal),
.markdown-body :deep(.hljs-bullet),
.markdown-body :deep(.hljs-link) {
  color: #ffa657;
}

.markdown-body :deep(.hljs-function),
.markdown-body :deep(.hljs-title.function_) {
  color: #d2a8ff;
}

.markdown-body :deep(.hljs-variable),
.markdown-body :deep(.hljs-params) {
  color: #c9d1d9;
}

.markdown-body :deep(.hljs-property) {
  color: #79c0ff;
}

.markdown-body :deep(.hljs-operator) {
  color: #ff7b72;
}

.markdown-body :deep(.hljs-class .hljs-title) {
  color: #ffa657;
}

.markdown-body :deep(.hljs-attr) {
  color: #a5d6ff;
}

.markdown-body :deep(.katex-block) {
  margin: 16px 0;
  padding: 16px;
  background: #fafafa;
  border-radius: 8px;
  overflow-x: auto;
}

.markdown-body :deep(.katex) {
  font-size: 1.1em;
}
</style>