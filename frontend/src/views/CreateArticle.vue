<template>
  <div class="page-container">
    <el-card class="card-surface article-form-card">
      <div class="form-header">
        <el-button circle text icon="ArrowLeft" @click="goBack" />
        <span class="form-title text-primary">
          {{ isEdit ? '编辑文章' : '发布新文章' }}
        </span>
      </div>

      <el-divider class="mb-4" />

      <el-form label-position="top" @submit.prevent="submitArticle">
        <el-form-item label="文章标题">
          <el-input
            v-model="form.title"
            placeholder="请输入文章标题"
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

        <el-form-item label="文章内容（支持 Markdown）">
          <div class="content-editor-wrap">
            <div class="editor-toolbar">
              <el-button text size="small" @click="showVoiceDialog = true">
                <el-icon><Microphone /></el-icon>
                添加语音
              </el-button>
            </div>
            <MdEditor
              v-model="form.content"
              :height="400"
              @voice-input="handleVoiceInput"
            />
          </div>
        </el-form-item>

        <!-- 已添加的语音展示 -->
        <el-alert
          v-if="voiceUrl || isRecording"
          type="info"
          :closable="!isRecording"
          :show-icon="false"
          class="mb-4"
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

        <el-card shadow="never" class="options-card mb-4">
          <el-checkbox v-model="form.is_anonymous">匿名发布</el-checkbox>

          <div class="form-actions">
            <el-button @click="goBack">取消</el-button>
            <el-button :loading="isSavingDraft" @click="saveDraft">
              {{ isSavingDraft ? '保存中...' : '保存草稿' }}
            </el-button>
            <el-button type="primary" native-type="submit" :loading="isSubmitting">
              {{ isSubmitting ? '提交中...' : (isEdit ? '更新文章' : '发布文章') }}
            </el-button>
          </div>
        </el-card>
      </el-form>
    </el-card>

    <!-- 语音对话框 -->
    <el-dialog
      v-model="showVoiceDialog"
      title="添加语音"
      width="520px"
      :close-on-click-modal="false"
    >
      <el-tabs v-model="voiceTab">
        <el-tab-pane name="record">
          <template #label>
            <el-icon><Microphone /></el-icon>
            录音
          </template>

          <div class="record-section">
            <div class="record-time">
              <el-icon :size="48" :color="isRecording ? '#f56c6c' : '#909399'" class="mb-3">
                <Microphone />
              </el-icon>
              <div class="time-display">{{ formatDuration(recordingDuration) }}</div>
              <div class="record-status">
                {{ isPaused ? '已暂停' : (isRecording ? '录音中...' : '点击开始录音') }}
              </div>
            </div>

            <div class="record-controls">
              <el-button
                v-if="!isRecording && !recordedBlob"
                type="danger"
                size="large"
                @click="startRecording"
              >
                <el-icon><Microphone /></el-icon>
                开始录音
              </el-button>

              <template v-else>
                <el-button
                  v-if="isRecording && !isPaused"
                  type="warning"
                  size="large"
                  @click="pauseRecording"
                >
                  <el-icon><VideoPause /></el-icon>
                  暂停
                </el-button>
                <el-button
                  v-if="isPaused"
                  type="success"
                  size="large"
                  @click="resumeRecording"
                >
                  <el-icon><VideoPlay /></el-icon>
                  继续
                </el-button>
                <el-button type="danger" size="large" @click="stopRecording">
                  <el-icon><VideoCameraFilled /></el-icon>
                  结束
                </el-button>
              </template>
            </div>

            <div v-if="recordedBlob" class="recorded-preview">
              <audio :src="recordedAudioUrl" controls class="voice-audio" />
              <div class="recorded-actions">
                <el-button type="primary" :loading="isUploadingVoice" @click="confirmVoice">
                  使用此录音
                </el-button>
                <el-button @click="discardRecording">丢弃</el-button>
              </div>
            </div>
          </div>
        </el-tab-pane>

        <el-tab-pane name="upload">
          <template #label>
            <el-icon><Upload /></el-icon>
            上传
          </template>

          <div class="upload-section">
            <el-upload
              ref="voiceUploadRef"
              :auto-upload="false"
              accept="audio/*"
              :limit="1"
              :on-change="handleUploadChange"
              :on-remove="handleUploadRemove"
            >
              <el-button :icon="Upload">选择音频文件</el-button>
              <template #tip>
                <div class="el-upload__tip upload-hint">
                  支持 MP3、WAV、AAC、OGG 等常见音频格式
                </div>
              </template>
            </el-upload>
            <el-button
              type="primary"
              :disabled="!voiceFile"
              :loading="isUploadingVoice"
              class="mt-2"
              @click="uploadSelectedFile"
            >
              上传音频
            </el-button>
          </div>
        </el-tab-pane>
      </el-tabs>

      <template #footer>
        <el-button @click="showVoiceDialog = false">关闭</el-button>
      </template>
    </el-dialog>
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
  Upload,
  VideoPause,
  VideoPlay,
  VideoCameraFilled
} from '@element-plus/icons-vue'
import { success as showSuccess, error as showError, warning as showWarning } from '@/utils/message'

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

const showVoiceDialog = ref(false)
const voiceTab = ref('record')
const voiceUploadRef = ref(null)
const voiceFile = ref(null)
const recordedBlob = ref(null)
const recordedAudioUrl = ref('')

let mediaRecorder = null
let audioChunks = []
let recordingTimer = null

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
    showVoiceDialog.value = false
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

const handleUploadChange = (file) => {
  voiceFile.value = file.raw || null
}

const handleUploadRemove = () => {
  voiceFile.value = null
}

const uploadSelectedFile = async () => {
  if (!voiceFile.value) return

  isUploadingVoice.value = true
  try {
    const formData = new FormData()
    formData.append('voice', voiceFile.value)

    const response = await http.post('/upload/voice', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })

    voiceUrl.value = response.data.url
    showVoiceDialog.value = false
    voiceFile.value = null
    if (voiceUploadRef.value) {
      voiceUploadRef.value.clearFiles()
    }
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
.article-form-card {
  max-width: 960px;
  margin: 0 auto;
  border-radius: var(--campus-radius);
}

.form-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.form-title {
  font-size: 1.25rem;
  font-weight: 700;
}

.content-editor-wrap {
  width: 100%;
}

.editor-toolbar {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 8px;
}

.options-card {
  border: 1px solid var(--campus-border);
  border-radius: var(--campus-radius);
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
}

.form-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
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
</style>
