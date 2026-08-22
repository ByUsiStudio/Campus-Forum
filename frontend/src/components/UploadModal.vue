<template>
  <!-- UploadModal 基于 JCuPupw：展示为自定义 HTML 弹窗 -->
  <div class="upload-modal-slot"></div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { jcOpenHtml, jcCloseAll } from '@/utils/jcu'

const props = defineProps({
  show: {
    type: Boolean,
    default: false
  },
  uploadType: {
    type: String,
    default: 'image',
    validator: (value) => ['image', 'video'].includes(value)
  }
})

const emit = defineEmits(['update:show', 'upload-success'])

const uploading = ref(false)
const progress = ref(0)
const uploadedFiles = ref([])
let contentEl = null

const resetState = () => {
  uploading.value = false
  progress.value = 0
  uploadedFiles.value = []
}

const renderProgress = () => {
  if (!contentEl) return
  const wrap = contentEl.querySelector('[data-uc-progress]')
  const bar = contentEl.querySelector('[data-uc-bar]')
  const pct = contentEl.querySelector('[data-uc-pct]')
  if (wrap) wrap.style.display = uploading.value ? 'block' : 'none'
  if (bar) bar.style.width = progress.value + '%'
  if (pct) pct.textContent = progress.value + '%'
}

const renderList = () => {
  if (!contentEl) return
  const list = contentEl.querySelector('[data-uc-list]')
  if (!list) return
  if (uploadedFiles.value.length) {
    const isImage = props.uploadType === 'image'
    list.style.display = 'flex'
    list.innerHTML = uploadedFiles.value.map((file, index) => `
      <div style="display:flex;align-items:center;gap:12px;padding:12px;background:var(--campus-primary-soft,rgba(79,110,247,.08));border-radius:10px;">
        ${isImage
          ? `<img src="${escAttr(file.url)}" style="width:56px;height:56px;object-fit:cover;border-radius:8px;flex-shrink:0;" />`
          : `<div style="width:56px;height:56px;display:flex;align-items:center;justify-content:center;background:var(--campus-border,#eef1f7);border-radius:8px;font-size:24px;flex-shrink:0;">🎬</div>`}
        <div style="flex:1;min-width:0;">
          <div style="font-weight:600;color:var(--campus-text,#333);font-size:14px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;">${esc(file.name)}</div>
          <div style="font-size:12px;color:var(--campus-text-secondary,#6b7280);overflow:hidden;text-overflow:ellipsis;white-space:nowrap;margin-top:4px;">${esc(file.url)}</div>
        </div>
        <button type="button" data-uc-remove="${index}" style="border:none;background:transparent;color:var(--campus-danger,#ef4444);cursor:pointer;font-size:13px;flex-shrink:0;">删除</button>
      </div>
    `).join('')
  } else {
    list.style.display = 'none'
    list.innerHTML = ''
  }
}

const escAttr = (s) => String(s ?? '').replace(/"/g, '&quot;').replace(/</g, '&lt;')
const esc = (s) => String(s ?? '').replace(/</g, '&lt;').replace(/>/g, '&gt;')

const open = () => {
  resetState()
  const isImage = props.uploadType === 'image'
  const html = `
    <div data-uc-area style="border:2px dashed var(--campus-primary-soft,rgba(79,110,247,.35));border-radius:12px;padding:30px;text-align:center;cursor:pointer;background:var(--campus-surface-2,#f8f9fc);transition:all .15s;">
      <input type="file" data-uc-file accept="${isImage ? 'image/*' : 'video/*'}" style="display:none;" />
      <div style="font-size:40px;margin-bottom:10px;color:var(--campus-primary,#4f6ef7);">${isImage ? '🖼️' : '🎬'}</div>
      <div style="font-weight:600;font-size:15px;color:var(--campus-text,#333);">点击或拖拽文件到此处</div>
      <div style="font-size:13px;color:var(--campus-text-secondary,#6b7280);margin-top:8px;">${isImage ? '支持 JPG、PNG、GIF 格式' : '支持 MP4、WebM 格式'}</div>
    </div>
    <div data-uc-progress style="display:none;margin-top:16px;">
      <div style="height:8px;background:var(--campus-border,#e6e9f2);border-radius:99px;overflow:hidden;">
        <div data-uc-bar style="height:100%;width:0%;background:var(--campus-primary,#4f6ef7);border-radius:99px;transition:width .2s;"></div>
      </div>
      <div data-uc-pct style="text-align:center;font-size:13px;color:var(--campus-text-secondary,#6b7280);margin-top:8px;">0%</div>
    </div>
    <div data-uc-list style="margin-top:16px;display:none;flex-direction:column;gap:10px;"></div>
  `

  jcOpenHtml({
    title: isImage ? '上传图片' : '上传视频',
    content: html,
    width: 500,
    size: 'sm',
    closeOnOverlay: false,
    onMount: (root) => {
      contentEl = root
      const area = root.querySelector('[data-uc-area]')
      const file = root.querySelector('[data-uc-file]')
      const list = root.querySelector('[data-uc-list]')

      area.addEventListener('click', () => file.click())
      area.addEventListener('dragover', (e) => {
        e.preventDefault()
        area.style.borderColor = 'var(--campus-primary,#4f6ef7)'
      })
      area.addEventListener('dragleave', () => { area.style.borderColor = '' })
      area.addEventListener('drop', (e) => {
        e.preventDefault()
        area.style.borderColor = ''
        const f = e.dataTransfer.files[0]
        if (f) uploadFile(f)
      })
      file.addEventListener('change', (e) => {
        const f = e.target.files[0]
        e.target.value = ''
        if (f) uploadFile(f)
      })
      list.addEventListener('click', (e) => {
        const btn = e.target.closest('[data-uc-remove]')
        if (btn) removeFile(Number(btn.getAttribute('data-uc-remove')))
      })
    },
    onClose: () => {
      resetState()
      emit('update:show', false)
      contentEl = null
    },
    buttons: [
      { text: '取消', type: 'default', action: () => jcCloseAll() },
      { text: '确认插入', type: 'primary', close: false, action: () => confirmInsert() }
    ]
  })
}

watch(() => props.show, (val) => {
  if (val) open()
})

const confirmInsert = () => {
  if (uploadedFiles.value.length === 0 || uploading.value) return
  emit('upload-success', { type: props.uploadType, files: [...uploadedFiles.value] })
  jcCloseAll()
}

const removeFile = (index) => {
  uploadedFiles.value.splice(index, 1)
  renderList()
}

const uploadFile = (file) => {
  uploading.value = true
  progress.value = 0
  renderProgress()
  renderList()

  const formData = new FormData()
  formData.append(props.uploadType, file)

  const xhr = new XMLHttpRequest()
  const token = localStorage.getItem('token')
  xhr.open('POST', `/api/upload/${props.uploadType}`)
  if (token) xhr.setRequestHeader('Authorization', `Bearer ${token}`)

  xhr.upload.onprogress = (event) => {
    if (event.lengthComputable) {
      progress.value = Math.min(99, Math.round((event.loaded / event.total) * 100))
      renderProgress()
    }
  }

  const done = () => {
    uploading.value = false
    renderProgress()
    renderList()
  }

  xhr.onload = () => {
    try {
      const data = JSON.parse(xhr.responseText)
      if (xhr.status >= 200 && xhr.status < 300 && data.url) {
        progress.value = 100
        uploadedFiles.value.push({ name: file.name, url: data.url })
      } else {
        console.error('Upload failed:', data.error || data.message || `HTTP ${xhr.status}`)
      }
    } catch (error) {
      console.error('Upload failed:', error)
    }
    done()
  }

  xhr.onerror = () => {
    console.error('Upload failed: network error')
    done()
  }

  xhr.onabort = () => {
    done()
  }

  xhr.send(formData)
}
</script>

<style scoped>
.upload-modal-slot { display: none; }
</style>
