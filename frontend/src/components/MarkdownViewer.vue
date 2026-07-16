<template>
  <div class="markdown-viewer" v-html="renderedContent"></div>
</template>

<script setup>
import { computed, onMounted, nextTick, ref } from 'vue'
import MarkdownIt from 'markdown-it'
import anchor from 'markdown-it-anchor'
import container from 'markdown-it-container'
import * as emoji from 'markdown-it-emoji'
import footnote from 'markdown-it-footnote'
import mark from 'markdown-it-mark'
import taskLists from 'markdown-it-task-lists'
import katex from 'katex'
import 'katex/dist/katex.min.css'

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
  value: {
    type: String,
    default: ''
  }
})

const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
  breaks: true,
  highlight: function (str, lang) {
    if (lang && window.hljs) {
      try {
        const language = window.hljs.getLanguage(lang) ? lang : 'plaintext'
        return '<pre class="code-block" data-language="' + language + '"><code class="hljs language-' + language + '">' +
          window.hljs.highlight(str, { language: language, ignoreIllegals: true }).value +
          '</code></pre>'
      } catch (e) {}
    }
    return '<pre class="code-block"><code class="hljs">' + md.utils.escapeHtml(str) + '</code></pre>'
  }
})

// 自定义链接渲染规则：外链自动在新标签打开
const defaultLinkRenderer = md.renderer.rules.link_open || function(tokens, idx, options, env, self) {
  return self.renderToken(tokens, idx, options)
}

md.renderer.rules.link_open = function(tokens, idx, options, env, self) {
  const token = tokens[idx]
  const hrefIndex = token.attrIndex('href')
  
  if (hrefIndex >= 0) {
    const href = token.attrs[hrefIndex][1]
    
    // 判断是否为外链
    const isExternal = href && (
      href.startsWith('http://') ||
      href.startsWith('https://') ||
      href.startsWith('ftp://')
    )
    
    // 判断是否为本站链接
    const isInternal = href && (
      href.startsWith('/') ||
      href.startsWith('#') ||
      href.startsWith('?') ||
      href.startsWith('./') ||
      href.startsWith('../')
    )
    
    // 如果是外链且没有 target 属性，添加 target="_blank"
    if (isExternal && !isInternal) {
      const targetIndex = token.attrIndex('target')
      if (targetIndex < 0) {
        token.attrPush(['target', '_blank'])
      }
      
      // 添加 rel 属性以提高安全性
      const relIndex = token.attrIndex('rel')
      if (relIndex < 0) {
        token.attrPush(['rel', 'noopener noreferrer'])
      }
    }
  }
  
  return defaultLinkRenderer(tokens, idx, options, env, self)
}

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

const renderMath = (content) => {
  let result = content
  
  result = result.replace(/\$\$([\s\S]*?)\$\$/g, (match, formula) => {
    return katexBlock(formula.trim())
  })
  
  result = result.replace(/\$(.+?)\$/g, (match, formula) => {
    return katexInline(formula.trim())
  })
  
  return result
}

const renderedContent = computed(() => {
  if (!props.value) return '<p class="empty-hint">暂无内容</p>'
  const rendered = md.render(props.value)
  return renderMath(rendered)
})
</script>

<style scoped>
.markdown-viewer {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  font-size: 16px;
  line-height: 1.9;
  color: #2d2d2d;
  word-wrap: break-word;
  white-space: pre-wrap;
}

.markdown-viewer :deep(h1),
.markdown-viewer :deep(h2),
.markdown-viewer :deep(h3),
.markdown-viewer :deep(h4),
.markdown-viewer :deep(h5),
.markdown-viewer :deep(h6) {
  margin: 28px 0 16px 0;
  font-weight: 600;
  line-height: 1.4;
  color: #1a1a1a;
  position: relative;
}

.markdown-viewer :deep(h1) {
  font-size: 2rem;
  padding-bottom: 12px;
  border-bottom: 2px solid #e5e5e5;
}

.markdown-viewer :deep(h2) {
  font-size: 1.6rem;
  padding-bottom: 8px;
  border-bottom: 1px solid #e5e5e5;
}

.markdown-viewer :deep(h3) {
  font-size: 1.3rem;
}

.markdown-viewer :deep(h4) {
  font-size: 1.1rem;
}

.markdown-viewer :deep(.header-anchor) {
  position: absolute;
  left: -24px;
  color: #999;
  text-decoration: none;
  opacity: 0;
  transition: opacity 0.2s;
}

.markdown-viewer :deep(h1:hover .header-anchor),
.markdown-viewer :deep(h2:hover .header-anchor),
.markdown-viewer :deep(h3:hover .header-anchor),
.markdown-viewer :deep(h4:hover .header-anchor) {
  opacity: 1;
}

.markdown-viewer :deep(p) {
  margin: 16px 0;
}

.markdown-viewer :deep(a) {
  color: #6750A4;
  text-decoration: none;
  border-bottom: 1px solid transparent;
  transition: border-color 0.2s;
}

.markdown-viewer :deep(a:hover) {
  border-bottom-color: #6750A4;
}

.markdown-viewer :deep(audio) {
  width: 100%;
  max-width: 100%;
  margin: 16px 0;
  display: block;
}

.markdown-viewer :deep(code) {
  background: #f5f5f5;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'JetBrains Mono', 'Monaco', 'Consolas', monospace;
  font-size: 0.9em;
  color: #e53935;
}

.markdown-viewer :deep(pre) {
  background: #1e1e1e;
  border-radius: 8px;
  padding: 16px;
  overflow-x: auto;
  margin: 20px 0;
}

.markdown-viewer :deep(pre code) {
  background: none;
  padding: 0;
  color: #d4d4d4;
  font-size: 14px;
  line-height: 1.6;
  white-space: pre;
}

.markdown-viewer :deep(blockquote) {
  margin: 20px 0;
  padding: 12px 20px;
  border-left: 4px solid #6750A4;
  background: linear-gradient(135deg, #f8f7ff 0%, #f0effe 100%);
  border-radius: 0 8px 8px 0;
  color: #555;
}

.markdown-viewer :deep(blockquote p) {
  margin: 0;
}

.markdown-viewer :deep(blockquote blockquote) {
  margin: 12px 0;
  border-left-color: #9575CD;
  background: linear-gradient(135deg, #f5f3ff 0%, #ede7f6 100%);
}

.markdown-viewer :deep(ul),
.markdown-viewer :deep(ol) {
  margin: 16px 0;
  padding-left: 28px;
}

.markdown-viewer :deep(ul) {
  list-style-type: disc;
}

.markdown-viewer :deep(ol) {
  list-style-type: decimal;
}

.markdown-viewer :deep(li) {
  margin: 8px 0;
  line-height: 1.8;
  white-space: pre-wrap;
}

.markdown-viewer :deep(li > ul),
.markdown-viewer :deep(li > ol) {
  margin: 4px 0;
}

.markdown-viewer :deep(table) {
  width: 100%;
  margin: 20px 0;
  border-collapse: collapse;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.markdown-viewer :deep(thead) {
  background: linear-gradient(135deg, #6750A4 0%, #7E67C8 100%);
  color: #fff;
}

.markdown-viewer :deep(th) {
  padding: 14px 16px;
  text-align: left;
  font-weight: 600;
  font-size: 14px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.markdown-viewer :deep(td) {
  padding: 12px 16px;
  border-bottom: 1px solid #eee;
}

.markdown-viewer :deep(tbody tr:nth-child(even)) {
  background: #fafafa;
}

.markdown-viewer :deep(tbody tr:hover) {
  background: #f0f0f0;
}

.markdown-viewer :deep(img) {
  max-width: 100%;
  max-height: 600px;
  width: auto;
  height: auto;
  border-radius: 8px;
  margin: 20px auto;
  display: block;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s, box-shadow 0.2s;
  object-fit: contain;
}

.markdown-viewer :deep(img:hover) {
  transform: scale(1.01);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
}

.markdown-viewer :deep(img[data-width="large"]),
.markdown-viewer :deep(img[src*=".gif"]) {
  max-height: none;
}

.markdown-viewer :deep(hr) {
  margin: 32px 0;
  border: none;
  height: 2px;
  background: linear-gradient(to right, transparent, #e0e0e0, transparent);
}

.markdown-viewer :deep(del) {
  color: #999;
  text-decoration: line-through;
}

.markdown-viewer :deep(mark) {
  background: linear-gradient(120deg, #ffeaa7 0%, #fdcb6e 100%);
  padding: 2px 4px;
  border-radius: 3px;
}

.markdown-viewer :deep(.footnotes) {
  margin-top: 40px;
  padding-top: 20px;
  border-top: 2px solid #e5e5e5;
  font-size: 14px;
  color: #666;
}

.markdown-viewer :deep(.footnotes-list) {
  padding-left: 20px;
}

.markdown-viewer :deep(.footnote-item) {
  margin: 12px 0;
}

.markdown-viewer :deep(.footnote-ref) {
  color: #6750A4;
  font-weight: bold;
  text-decoration: none;
}

.markdown-viewer :deep(.custom-container) {
  margin: 20px 0;
  padding: 16px 20px;
  border-radius: 8px;
  border-left: 4px solid;
}

.markdown-viewer :deep(.custom-container p) {
  margin: 0;
  font-weight: 600;
}

.markdown-viewer :deep(.custom-container.tip) {
  background: linear-gradient(135deg, #e8f5e9 0%, #c8e6c9 100%);
  border-color: #4CAF50;
}

.markdown-viewer :deep(.custom-container.warning) {
  background: linear-gradient(135deg, #fff3e0 0%, #ffe0b2 100%);
  border-color: #FF9800;
}

.markdown-viewer :deep(.custom-container.danger) {
  background: linear-gradient(135deg, #ffebee 0%, #ffcdd2 100%);
  border-color: #f44336;
}

.markdown-viewer :deep(.custom-container.info) {
  background: linear-gradient(135deg, #e3f2fd 0%, #bbdefb 100%);
  border-color: #2196F3;
}

.markdown-viewer :deep(.task-list-item) {
  list-style: none;
  margin-left: -24px;
}

.markdown-viewer :deep(.task-list-item input[type="checkbox"]) {
  margin-right: 8px;
  width: 16px;
  height: 16px;
  accent-color: #6750A4;
}

.markdown-viewer :deep(.emoji) {
  display: inline;
  vertical-align: middle;
}

.markdown-viewer :deep(video) {
  max-width: 100%;
  max-height: 500px;
  width: auto;
  height: auto;
  border-radius: 8px;
  margin: 20px 0;
  background: #000;
  object-fit: contain;
}

.markdown-viewer :deep(.empty-hint) {
  color: #999;
  text-align: center;
  padding: 40px;
}

.markdown-viewer :deep(kbd) {
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

.markdown-viewer :deep(.katex-block) {
  margin: 20px 0;
  padding: 16px;
  background: #fafafa;
  border-radius: 8px;
  overflow-x: auto;
}

.markdown-viewer :deep(.katex) {
  font-size: 1.1em;
}

.markdown-viewer :deep(.code-block) {
  position: relative;
  border-radius: 8px;
  overflow: hidden;
}

.markdown-viewer :deep(.code-block::before) {
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

.markdown-viewer :deep(.code-block pre) {
  margin: 0;
  border-radius: 0;
}

.markdown-viewer :deep(.code-block code) {
  display: block;
  padding: 16px;
  overflow-x: auto;
  font-family: 'JetBrains Mono', 'Fira Code', 'Consolas', monospace;
  font-size: 14px;
  line-height: 1.6;
}

.markdown-viewer :deep(.hljs) {
  background: #0d1117 !important;
  color: #c9d1d9;
}

.markdown-viewer :deep(.hljs-keyword),
.markdown-viewer :deep(.hljs-selector-tag),
.markdown-viewer :deep(.hljs-built_in),
.markdown-viewer :deep(.hljs-name),
.markdown-viewer :deep(.hljs-tag) {
  color: #ff7b72;
}

.markdown-viewer :deep(.hljs-string),
.markdown-viewer :deep(.hljs-title),
.markdown-viewer :deep(.hljs-section),
.markdown-viewer :deep(.hljs-attribute),
.markdown-viewer :deep(.hljs-literal),
.markdown-viewer :deep(.hljs-template-tag),
.markdown-viewer :deep(.hljs-template-variable),
.markdown-viewer :deep(.hljs-type),
.markdown-viewer :deep(.hljs-addition) {
  color: #a5d6ff;
}

.markdown-viewer :deep(.hljs-comment),
.markdown-viewer :deep(.hljs-quote),
.markdown-viewer :deep(.hljs-deletion),
.markdown-viewer :deep(.hljs-meta) {
  color: #8b949e;
}

.markdown-viewer :deep(.hljs-number),
.markdown-viewer :deep(.hljs-regexp),
.markdown-viewer :deep(.hljs-literal),
.markdown-viewer :deep(.hljs-bullet),
.markdown-viewer :deep(.hljs-link) {
  color: #ffa657;
}

.markdown-viewer :deep(.hljs-function),
.markdown-viewer :deep(.hljs-title.function_) {
  color: #d2a8ff;
}

.markdown-viewer :deep(.hljs-variable),
.markdown-viewer :deep(.hljs-params) {
  color: #c9d1d9;
}

.markdown-viewer :deep(.hljs-property) {
  color: #79c0ff;
}

.markdown-viewer :deep(.hljs-operator) {
  color: #ff7b72;
}

.markdown-viewer :deep(.hljs-class .hljs-title) {
  color: #ffa657;
}

.markdown-viewer :deep(.hljs-attr) {
  color: #a5d6ff;
}
</style>