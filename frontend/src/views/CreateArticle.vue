<template>
  <div class="page-container">
    <div class="compose-page card-surface">
      <div class="form-header">
        <el-button circle text icon="ArrowLeft" class="back-btn" @click="goBack" />
        <div class="header-text">
          <span class="form-title">
            {{ isEdit ? '编辑文章' : '发布新文章' }}
          </span>
          <span class="form-subtitle text-secondary">
            {{ isEdit ? '修改稿件内容后保存更新' : '记录灵感，分享到校园' }}
          </span>
        </div>
      </div>

      <el-form label-position="top" class="compose-form" @submit.prevent="submitArticle">
        <div class="compose-grid">
          <!-- 侧栏：标题 + 分区 + 设置 + 提交 -->
          <div class="compose-side">
            <div class="side-card">
              <h3 class="side-title">文章信息</h3>

              <el-form-item label="文章标题">
                <el-input
                  v-model="form.title"
                  placeholder="请输入文章标题"
                  size="large"
                  clearable
                >
                  <template #prefix>
                    <el-icon><EditPen /></el-icon>
                  </template>
                </el-input>
              </el-form-item>

              <el-form-item label="选择分区">
                <el-select
                  v-model="form.category_id"
                  placeholder="请选择分区"
                  style="width: 100%"
                >
                  <el-option
                    v-for="cat in categories"
                    :key="cat.id"
                    :label="cat.name"
                    :value="cat.id"
                  />
                </el-select>
              </el-form-item>

              <el-form-item>
                <el-checkbox v-model="form.is_anonymous">匿名发布</el-checkbox>
              </el-form-item>
            </div>

            <div class="side-actions">
              <el-button class="action-btn" @click="goBack">取消</el-button>
              <el-button class="action-btn" :loading="isSavingDraft" @click="saveDraft">
                <el-icon v-if="!isSavingDraft" class="mr-1"><Document /></el-icon>
                {{ isSavingDraft ? '保存中...' : '保存草稿' }}
              </el-button>
              <el-button
                type="primary"
                class="action-btn publish-btn"
                native-type="submit"
                :loading="isSubmitting"
              >
                {{ isSubmitting ? '提交中...' : (isEdit ? '更新文章' : '发布文章') }}
              </el-button>
            </div>
          </div>

          <!-- 主栏：编辑器 + 语音 -->
          <div class="compose-main">
            <div class="editor-toolbar">
              <span class="editor-label">
                <el-icon><EditPen /></el-icon>
                文章内容（支持 Markdown）
              </span>
              <el-button text size="small" type="primary" @click="openVoiceDialog">
                <el-icon><Microphone /></el-icon>
                添加语音
              </el-button>
            </div>

            <MdEditor
              v-model="form.content"
              :height="520"
              @voice-input="handleVoiceInput"
            />

            <!-- 已添加的语音展示 -->
            <el-alert
              v-if="voiceUrl || isRecording"
              type="info"
              :closable="!isRecording"
              :show-icon="false"
              class="voice-alert"
              @close="removeVoice"
            >
              <div class="voice-status-row">
                <el-icon :size="18" :class="isRecording ? 'status-error' : 'status-primary'">
                  <Microphone />
                </el-icon>
                <div class="voice-status-main">
                  <div class="voice-status-label">
                    {{ isRecording ? '录音中...' : '语音已添加' }}
                  </div>
                  <div v-if="voiceUrl && !isRecording" class="text-secondary voice-filename">
                    {{ voiceUrl.split('/').pop() }}
                  </div>
                  <audio
                    v-if="voiceUrl && !isRecording"
                    :src="voiceUrl"
                    controls
                    class="voice-audio"
                  />
                </div>
              </div>
              <el-progress
                v-if="isRecording"
                :percentage="Math.min(recordingDuration / 60 * 100, 100)"
                :stroke-width="4"
                color="#f56c6c"
                class="mt-2"
              />
            </el-alert>
          </div>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import http from '@/api/http'
import MdEditor from '@/components/MdEditor.vue'
import {
  ArrowLeft,
  EditPen,
  Microphone,
  Document
} from '@element-plus/icons-vue'
import { success as showSuccess, error as showError, warning as showWarning } from '@/utils/message'
import { jcOpenHtml, jcCloseAll } from '@/utils/jcu'

const route = useRoute()
const router = useRouter()

const categories = ref([])
const form = ref({
  title: '',
  category_id: null,
  content: '',
  is_anonymous: false
})
const isEdit = ref(false)
const articleId = ref(null)
const isSubmitting = ref(false)
const isSavingDraft = ref(false)
const voiceUrl = ref('')
const isRecording = ref(false)
const isPaused = ref(false)
const recordingDuration = ref(0)
const isUploadingVoice = ref(false)

const voiceFile = ref(null)
const recordedBlob = ref(null)
const recordedAudioUrl = ref('')

let mediaRecorder = null
let audioChunks = []
let recordingTimer = null
let voiceDialogEl = null
let activeVoiceTab = 'record'

const formatDuration = (seconds) => {
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
}

const goBack = () => {
  router.push('/profile')
}

const loadCategories = async () => {
  try {
    const response = await http.get('/categories')
    categories.value = response.data.categories
    if (categories.value.length > 0 && !form.value.category_id) {
      form.value.category_id = categories.value[0].id
    }
  } catch (error) {
    console.error('加载分区失败', error)
  }
}

const loadArticle = async () => {
  try {
    const response = await http.get(`/articles/${articleId.value}`)
    const article = response.data.article
    form.value.title = article.title
    form.value.category_id = article.category_id
    form.value.content = article.content
    form.value.is_anonymous = article.is_anonymous
    voiceUrl.value = article.voice_url || ''
  } catch (error) {
    console.error('加载文章失败', error)
  }
}

const saveDraft = async () => {
  if (!form.value.title.trim()) {
    await showWarning('请输入标题')
    return
  }

  isSavingDraft.value = true
  try {
    const submitData = {
      ...form.value,
      voice_url: voiceUrl.value,
      status: 'draft'
    }

    let response
    if (isEdit.value) {
      response = await http.put(`/articles/${articleId.value}`, submitData)
    } else {
      response = await http.post('/articles', submitData)
    }

    await showSuccess('草稿保存成功')
    await router.push('/profile')
  } catch (error) {
    console.error('保存草稿失败', error)
    await showError('保存失败: ' + (error.message || '未知错误'))
  } finally {
    isSavingDraft.value = false
  }
}

const submitArticle = async () => {
  if (!form.value.title.trim()) {
    await showWarning('请输入标题')
    return
  }
  if (!form.value.content.trim()) {
    await showWarning('请输入内容')
    return
  }

  isSubmitting.value = true

  try {
    let response
    const submitData = {
      title: form.value.title,
      content: form.value.content,
      category_id: form.value.category_id,
      voice_url: voiceUrl.value,
      is_anonymous: form.value.is_anonymous
    }

    if (isEdit.value) {
      response = await http.put(`/articles/${articleId.value}`, submitData)
    } else {
      response = await http.post('/articles', submitData)
    }

    const newArticleId = response.data.article?.id
    if (!newArticleId) {
      throw new Error('文章ID不存在')
    }

    await router.push(`/article/${newArticleId}`)
  } catch (error) {
    console.error('提交或跳转失败', error)
    await showError('提交失败: ' + (error.message || '未知错误'))
  } finally {
    isSubmitting.value = false
  }
}

const handleVoiceInput = (audioBlob) => {
  uploadVoice(audioBlob)
}

const uploadVoice = async (audioBlob) => {
  isUploadingVoice.value = true
  try {
    const formData = new FormData()
    const extension = audioBlob.type === 'audio/webm' ? 'webm' : 'mp3'
    formData.append('voice', audioBlob, `voice.${extension}`)
    const response = await http.post('/upload/voice', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })

    voiceUrl.value = response.data.url
    closeVoiceDialog()
    await showSuccess('语音上传成功')
  } catch (error) {
    console.error('语音上传失败', error)
    await showError('语音上传失败')
  } finally {
    isUploadingVoice.value = false
  }
}

const startRecording = async () => {
  try {
    const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
    mediaRecorder = new MediaRecorder(stream)
    audioChunks = []

    mediaRecorder.ondataavailable = (event) => {
      if (event.data.size > 0) {
        audioChunks.push(event.data)
      }
    }

    mediaRecorder.onstop = () => {
      stream.getTracks().forEach(track => track.stop())
      if (audioChunks.length > 0) {
        recordedBlob.value = new Blob(audioChunks, { type: 'audio/webm' })
        recordedAudioUrl.value = URL.createObjectURL(recordedBlob.value)
      }
    }

    mediaRecorder.start(100)
    isRecording.value = true
    isPaused.value = false
    recordingDuration.value = 0

    recordingTimer = setInterval(() => {
      if (!isPaused.value) {
        recordingDuration.value += 1
      }
      if (voiceDialogEl) renderVoiceContent(voiceDialogEl)
    }, 1000)
  } catch (error) {
    console.error('无法访问麦克风', error)
    await showError('无法访问麦克风，请检查权限设置')
  }
}

const pauseRecording = () => {
  if (mediaRecorder && mediaRecorder.state === 'recording') {
    mediaRecorder.pause()
    isPaused.value = true
  }
}

const resumeRecording = () => {
  if (mediaRecorder && mediaRecorder.state === 'paused') {
    mediaRecorder.resume()
    isPaused.value = false
  }
}

const stopRecording = () => {
  if (mediaRecorder && mediaRecorder.state !== 'inactive') {
    mediaRecorder.stop()
  }
  isRecording.value = false
  isPaused.value = false
  if (recordingTimer) {
    clearInterval(recordingTimer)
    recordingTimer = null
  }
}

const confirmVoice = async () => {
  if (recordedBlob.value) {
    await uploadVoice(recordedBlob.value)
    discardRecording()
  }
}

const discardRecording = () => {
  recordedBlob.value = null
  if (recordedAudioUrl.value) {
    URL.revokeObjectURL(recordedAudioUrl.value)
    recordedAudioUrl.value = ''
  }
  recordingDuration.value = 0
}

const escHtml = (s) => String(s == null ? '' : s)
  .replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;')

const escAttr = (s) => String(s == null ? '' : s).replace(/"/g, '&quot;').replace(/</g, '&lt;')

// ---------------- 语音 JCuPupw 弹窗 ----------------
const openVoiceDialog = () => {
  voiceFile.value = null
  activeVoiceTab = 'record'
  jcOpenHtml({
    title: '添加语音',
    width: 520,
    size: 'md',
    closeOnOverlay: false,
    onMount: (el) => {
      voiceDialogEl = el
      renderVoiceContent(el)
    },
    onClose: () => {
      stopRecording()
      if (voiceDialogEl) { voiceDialogEl = null }
    },
    buttons: [{ text: '关闭', type: 'default', action: () => jcCloseAll() }]
  })
}

const closeVoiceDialog = () => {
  stopRecording()
  jcCloseAll()
}

const renderVoiceContent = (el) => {
  el.innerHTML = `
    <div data-jc-vtab style="display:flex;gap:6px;margin-bottom:16px;border-bottom:1px solid var(--campus-border,#e6e9f2);">
      <button type="button" data-jc-tab="record" style="${tabBtnStyle('record')}">🎙️ 录音</button>
      <button type="button" data-jc-tab="upload" style="${tabBtnStyle('upload')}">⬆️ 上传</button>
    </div>
    <div data-jc-vtab-body>${renderVoiceTabBody()}</div>
  `

  el.querySelectorAll('[data-jc-tab]').forEach((btn) => {
    btn.addEventListener('click', () => {
      activeVoiceTab = btn.getAttribute('data-jc-tab')
      renderVoiceContent(el)
    })
  })
  bindVoiceBody(el)
}

const tabBtnStyle = (tab) => {
  const active = activeVoiceTab === tab
  return `cursor:pointer;padding:8px 16px;border:none;border-radius:8px 8px 0 0;font-size:14px;font-weight:600;background:${active ? 'var(--campus-primary-soft,rgba(79,110,247,.1))' : 'transparent'};color:${active ? 'var(--campus-primary,#4f6ef7)' : 'var(--campus-text-secondary,#64748b)'};border-bottom:${active ? '2px solid var(--campus-primary,#4f6ef7)' : '2px solid transparent'};`
}

const renderVoiceTabBody = () => {
  if (activeVoiceTab === 'upload') {
    return `
      <div data-jc-upload-area style="border:2px dashed var(--campus-border,#e6e9f2);border-radius:12px;padding:22px;text-align:center;cursor:pointer;">
        <input type="file" data-jc-voice-file accept="audio/*" style="display:none;" />
        <div style="font-size:34px;">🎵</div>
        <div style="font-weight:600;margin:6px 0;color:var(--campus-text,#0f172a);">${voiceFile.value ? escHtml(voiceFile.value.name) : '选择音频文件'}</div>
        <div style="font-size:12px;color:var(--campus-text-secondary,#6b7280);">支持 MP3、WAV、AAC、OGG 等格式</div>
      </div>
      <button data-jc-upload-btn type="button" ${voiceFile.value ? '' : 'disabled'} style="margin-top:14px;width:100%;padding:10px;border:none;border-radius:10px;background:var(--campus-primary,#4f6ef7);color:#fff;font-size:14px;font-weight:600;cursor:${voiceFile.value ? 'pointer' : 'not-allowed'};opacity:${voiceFile.value ? 1 : .5};">${isUploadingVoice.value ? '上传中…' : '上传音频'}</button>
    `
  }
  // record tab
  const status = isPaused.value ? '已暂停' : (isRecording.value ? '录音中...' : '点击开始录音')
  const time = formatDuration(recordingDuration.value)
  let controls = ''
  if (!isRecording.value && !recordedBlob.value) {
    controls = `<button data-jc-rec-start type="button" style="padding:10px 22px;border:none;border-radius:99px;background:var(--campus-danger,#ef4444);color:#fff;font-size:14px;font-weight:600;cursor:pointer;">🎙️ 开始录音</button>`
  } else if (!recordedBlob.value) {
    controls = `<div style="display:flex;gap:10px;justify-content:center;flex-wrap:wrap;">
      ${isPaused.value
        ? `<button data-jc-rec-resume type="button" style="padding:10px 20px;border:none;border-radius:99px;background:var(--campus-success,#22c55e);color:#fff;cursor:pointer;font-weight:600;">继续</button>`
        : `<button data-jc-rec-pause type="button" style="padding:10px 20px;border:none;border-radius:99px;background:var(--campus-warning,#f59e0b);color:#fff;cursor:pointer;font-weight:600;">暂停</button>`}
      <button data-jc-rec-stop type="button" style="padding:10px 20px;border:none;border-radius:99px;background:var(--campus-danger,#ef4444);color:#fff;cursor:pointer;font-weight:600;">结束</button>
    </div>`
  } else {
    controls = `<div style="display:flex;gap:10px;justify-content:center;flex-wrap:wrap;">
      <button data-jc-rec-use type="button" style="padding:10px 20px;border:none;border-radius:99px;background:var(--campus-primary,#4f6ef7);color:#fff;cursor:pointer;font-weight:600;">${isUploadingVoice.value ? '上传中…' : '使用此录音'}</button>
      <button data-jc-rec-discard type="button" style="padding:10px 20px;border:none;border-radius:99px;background:transparent;border:1px solid var(--campus-border);color:var(--campus-text);cursor:pointer;">丢弃</button>
    </div>`
  }
  const preview = recordedBlob.value && recordedAudioUrl.value
    ? `<audio src="${escAttr(recordedAudioUrl.value)}" controls style="width:100%;margin-top:14px;"></audio>`
    : ''
  return `
    <div style="text-align:center;">
      <div style="font-size:44px;margin-bottom:6px;color:${isRecording.value ? 'var(--campus-danger,#ef4444)' : 'var(--campus-text-muted,#94a3b8)'};animation:${isRecording.value ? 'pulse 1.2s infinite' : 'none'};">🎙️</div>
      <div style="font-size:28px;font-weight:700;font-variant-numeric:tabular-nums;color:var(--campus-text,#0f172a);">${time}</div>
      <div style="font-size:13px;color:var(--campus-text-secondary,#64748b);margin:4px 0 16px;">${status}</div>
      ${controls}
      ${preview}
    </div>
  `
}

const bindVoiceBody = (el) => {
  const uploadArea = el.querySelector('[data-jc-upload-area]')
  const fileInput = el.querySelector('[data-jc-voice-file]')
  if (uploadArea) uploadArea.addEventListener('click', () => fileInput && fileInput.click())
  if (fileInput) {
    fileInput.addEventListener('change', (e) => {
      const f = e.target.files[0]
      e.target.value = ''
      if (f) { voiceFile.value = f; renderVoiceContent(el) }
    })
  }

  const startBtn = el.querySelector('[data-jc-rec-start]')
  if (startBtn) startBtn.addEventListener('click', () => { startRecording(); renderVoiceContent(el) })
  const recPause = el.querySelector('[data-jc-rec-pause]')
  if (recPause) recPause.addEventListener('click', () => { pauseRecording(); renderVoiceContent(el) })
  const recResume = el.querySelector('[data-jc-rec-resume]')
  if (recResume) recResume.addEventListener('click', () => { resumeRecording(); renderVoiceContent(el) })
  const recStop = el.querySelector('[data-jc-rec-stop]')
  if (recStop) recStop.addEventListener('click', () => { stopRecording(); renderVoiceContent(el) })
  const recUse = el.querySelector('[data-jc-rec-use]')
  if (recUse) recUse.addEventListener('click', () => confirmVoice())
  const recDiscard = el.querySelector('[data-jc-rec-discard]')
  if (recDiscard) recDiscard.addEventListener('click', () => { discardRecording(); renderVoiceContent(el) })

  const uploadBtn = el.querySelector('[data-jc-upload-btn]')
  if (uploadBtn) uploadBtn.addEventListener('click', () => uploadSelectedFile())
}

const uploadSelectedFile = async () => {
  if (!voiceFile.value) return

  isUploadingVoice.value = true
  renderVoiceContent(voiceDialogEl)
  try {
    const formData = new FormData()
    formData.append('voice', voiceFile.value)

    const response = await http.post('/upload/voice', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })

    voiceUrl.value = response.data.url
    voiceFile.value = null
    closeVoiceDialog()
    await showSuccess('音频上传成功')
  } catch (error) {
    console.error('音频上传失败', error)
    await showError('音频上传失败: ' + (error.message || '未知错误'))
  } finally {
    isUploadingVoice.value = false
  }
}

const removeVoice = () => {
  voiceUrl.value = ''
}

onUnmounted(() => {
  if (recordingTimer) {
    clearInterval(recordingTimer)
  }
  if (mediaRecorder && mediaRecorder.state !== 'inactive') {
    mediaRecorder.stop()
  }
  if (recordedAudioUrl.value) {
    URL.revokeObjectURL(recordedAudioUrl.value)
  }
})

onMounted(() => {
  const id = route.query.id
  if (id) {
    isEdit.value = true
    articleId.value = parseInt(id)
    loadArticle()
  }
  loadCategories()
})
</script>

<style scoped>
.compose-page {
  max-width: 1080px;
  margin: 0 auto;
  border-radius: var(--campus-radius-lg, 20px);
  box-shadow: var(--campus-shadow);
  padding: 24px;
}

.form-header {
  display: flex;
  align-items: center;
  gap: 14px;
  padding-bottom: 20px;
  border-bottom: 1px solid var(--campus-border, #e6e9f2);
  margin-bottom: 22px;
}

.back-btn {
  flex-shrink: 0;
}

.header-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.form-title {
  font-size: 1.3rem;
  font-weight: 700;
  letter-spacing: -0.01em;
  color: var(--campus-text, #0f172a);
}

.form-subtitle {
  font-size: 13px;
}

/* 双栏布局：侧栏 + 主编辑区 */
.compose-grid {
  display: grid;
  grid-template-columns: 300px 1fr;
  gap: 22px;
  align-items: start;
}

.compose-side {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.side-card {
  border: 1px solid var(--campus-border, #e6e9f2);
  border-radius: var(--campus-radius, 14px);
  background: var(--campus-surface-2, #f8f9fc);
  padding: 20px 18px;
}

.side-title {
  margin: 0 0 16px;
  font-size: 15px;
  font-weight: 700;
  color: var(--campus-text, #0f172a);
}

.side-card :deep(.el-form-item__label) {
  font-weight: 600;
  font-size: 13px;
  padding-bottom: 8px;
}

.side-card :deep(.el-input__wrapper),
.side-card :deep(.el-select__wrapper) {
  border-radius: 10px;
  background: #fff;
}

.side-actions {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.action-btn {
  width: 100%;
  margin-left: 0;
  height: 42px;
  border-radius: 10px;
  font-size: 14px;
}

.publish-btn {
  height: 46px;
  font-size: 15px;
}

/* 主编辑区 */
.compose-main {
  min-width: 0;
}

.editor-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 12px;
  flex-wrap: wrap;
}

.editor-label {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-weight: 600;
  font-size: 14px;
  color: var(--campus-text, #0f172a);
}

.voice-alert {
  margin-top: 14px;
  border-radius: 12px;
}

.voice-status-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.status-error {
  color: #f56c6c;
}

.status-primary {
  color: var(--campus-primary);
}

.voice-status-main {
  flex: 1;
}

.voice-status-label {
  font-size: 0.9rem;
  font-weight: 500;
}

.voice-filename {
  font-size: 0.8rem;
}

.voice-audio {
  width: 100%;
  max-width: 350px;
  height: 40px;
  margin-top: 4px;
}

.record-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px 0;
}

.record-time {
  text-align: center;
  margin-bottom: 24px;
}

.time-display {
  font-size: 2.5rem;
  font-weight: 600;
  font-family: 'Roboto Mono', 'JetBrains Mono', monospace;
  color: #333;
}

.record-status {
  font-size: 0.9rem;
  color: #666;
  margin-top: 8px;
}

.record-controls {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  justify-content: center;
}

.recorded-preview {
  margin-top: 24px;
  width: 100%;
  text-align: center;
}

.recorded-actions {
  display: flex;
  gap: 8px;
  justify-content: center;
  margin-top: 12px;
}

.upload-section {
  padding: 16px 0;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.upload-hint {
  color: #999;
  font-size: 12px;
  line-height: 1.5;
  margin-top: 8px;
}

.mr-1 {
  margin-right: 4px;
}
.mt-2 {
  margin-top: 8px;
}

/* 移动端：单列堆叠 */
@media (max-width: 768px) {
  .compose-page {
    padding: 16px;
    border-radius: var(--campus-radius, 14px);
  }

  .compose-grid {
    grid-template-columns: 1fr;
    gap: 18px;
  }

  .compose-side {
    order: 2;
  }

  .compose-main {
    order: 1;
  }

  .side-actions {
    flex-direction: row;
    flex-wrap: wrap;
  }

  .action-btn {
    flex: 1;
    min-width: 120px;
  }

  .form-title {
    font-size: 1.15rem;
  }
}
</style>
