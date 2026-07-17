<template>
  <div class="create-article-page">
    <div class="article-container">
      <div class="article-card">
        <div class="article-header">
          <button class="back-btn" @click="router.back()">
            <i class="fa-solid fa-arrow-left"></i>
          </button>
          <h1 class="article-title">{{ isEdit ? '编辑文章' : '发布新文章' }}</h1>
        </div>

        <div class="divider"></div>

        <form @submit.prevent="submitArticle" class="article-form">
          <div class="form-item">
            <label class="form-label">文章标题</label>
            <div class="input-group">
              <i class="fa-solid fa-font"></i>
              <input 
                type="text"
                v-model="form.title"
                placeholder="请输入标题"
                class="layui-input"
                required
              />
            </div>
          </div>

          <div class="form-item">
            <label class="form-label">选择分区</label>
            <div class="input-group">
              <i class="fa-solid fa-folder"></i>
              <select v-model="form.category_id" class="layui-input">
                <option :value="null">请选择分区</option>
                <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
              </select>
            </div>
          </div>

          <div class="content-card">
            <div class="content-header">
              <div class="content-label">
                <i class="fa-solid fa-pencil mr-2"></i>
                文章内容 (支持 Markdown)
              </div>
              <button type="button" class="layui-btn layui-btn-normal layui-btn-sm" @click="showVoiceDialog = true">
                <i class="fa-solid fa-microphone mr-2"></i>
                添加语音
              </button>
            </div>
            <MdEditor
              v-model="form.content"
              :height="400"
              @voice-input="handleVoiceInput"
            />
          </div>

          <div v-if="voiceUrl || isRecording" class="voice-alert" :class="{ recording: isRecording }">
            <div class="alert-content">
              <i :class="isRecording ? 'fa-solid fa-microphone text-error' : 'fa-solid fa-volume-high'"></i>
              <div class="alert-info">
                <div class="alert-title">{{ isRecording ? '录音中...' : '语音已添加' }}</div>
                <div v-if="voiceUrl && !isRecording" class="alert-subtitle">{{ voiceUrl.split('/').pop() }}</div>
              </div>
              <button type="button" class="close-btn" @click="removeVoice">
                <i class="fa-solid fa-xmark"></i>
              </button>
            </div>
            <div v-if="isRecording" class="progress-bar">
              <div class="progress-fill" :style="{ width: (recordingDuration / 60) * 100 + '%' }"></div>
            </div>
          </div>

          <div class="form-actions">
            <label class="checkbox-item">
              <input type="checkbox" v-model="form.is_anonymous" />
              <span>匿名发布</span>
            </label>

            <div class="action-buttons">
              <button type="button" class="layui-btn layui-btn-primary" @click="router.back()">
                取消
              </button>
              <button 
                type="button" 
                class="layui-btn layui-btn-normal"
                :disabled="savingDraft"
                @click="saveDraft"
              >
                <i class="fa-solid fa-save mr-2"></i>
                {{ savingDraft ? '保存中...' : '保存草稿' }}
              </button>
              <button 
                type="submit" 
                class="layui-btn layui-btn-primary"
                :disabled="submitting"
              >
                <i class="fa-solid fa-paper-plane mr-2"></i>
                {{ submitting ? '提交中...' : (isEdit ? '更新文章' : '发布文章') }}
              </button>
            </div>
          </div>
        </form>
      </div>
    </div>

    <div v-if="showVoiceDialog" class="voice-dialog-overlay" @click.self="showVoiceDialog = false">
      <div class="voice-dialog">
        <div class="dialog-header">
          <i class="fa-solid fa-microphone mr-2"></i>
          添加语音
        </div>
        
        <div class="dialog-body">
          <div class="voice-tabs">
            <button 
              class="tab-btn" 
              :class="{ active: voiceTab === 'record' }"
              @click="voiceTab = 'record'"
            >
              <i class="fa-solid fa-microphone mr-2"></i>
              录音
            </button>
            <button 
              class="tab-btn" 
              :class="{ active: voiceTab === 'upload' }"
              @click="voiceTab = 'upload'"
            >
              <i class="fa-solid fa-upload mr-2"></i>
              上传
            </button>
          </div>

          <div v-if="voiceTab === 'record'" class="record-section">
            <div class="record-time">
              <i :class="isRecording ? 'fa-solid fa-microphone text-error' : 'fa-solid fa-microphone text-muted'" style="font-size: 64px;"></i>
              <div class="time-display">{{ formatDuration(recordingDuration) }}</div>
              <div class="record-status">{{ isPaused ? '已暂停' : (isRecording ? '录音中...' : '点击开始录音') }}</div>
            </div>

            <div class="record-controls">
              <button 
                v-if="!isRecording && !recordedBlob"
                class="layui-btn layui-btn-danger"
                @click="startRecording"
              >
                <i class="fa-solid fa-microphone mr-2"></i>
                开始录音
              </button>

              <template v-else>
                <button 
                  v-if="isRecording"
                  class="layui-btn layui-btn-warm"
                  @click="pauseRecording"
                >
                  <i class="fa-solid fa-pause mr-2"></i>
                  暂停
                </button>
                <button 
                  v-if="isPaused"
                  class="layui-btn layui-btn-normal"
                  @click="resumeRecording"
                >
                  <i class="fa-solid fa-play mr-2"></i>
                  继续
                </button>
                <button 
                  class="layui-btn layui-btn-danger"
                  @click="stopRecording"
                >
                  <i class="fa-solid fa-square mr-2"></i>
                  结束
                </button>
              </template>
            </div>

            <div v-if="recordedBlob" class="recorded-preview">
              <audio :src="recordedAudioUrl" controls class="audio-preview" />
              <div class="preview-actions">
                <button class="layui-btn layui-btn-normal" @click="confirmVoice" :disabled="isUploadingVoice">
                  使用此录音
                </button>
                <button class="layui-btn layui-btn-primary" @click="discardRecording">
                  丢弃
                </button>
              </div>
            </div>
          </div>

          <div v-if="voiceTab === 'upload'" class="upload-section">
            <div class="upload-input-wrapper">
              <i class="fa-solid fa-file-audio"></i>
              <input 
                type="file" 
                accept="audio/*"
                class="upload-input"
                @change="handleFileChange"
              />
              <span v-if="voiceFile">已选择: {{ voiceFile.name }}</span>
              <span v-else>点击选择音频文件</span>
            </div>
            <button 
              class="layui-btn layui-btn-normal w-full"
              :disabled="!voiceFile"
              :class="{ loading: isUploadingVoice }"
              @click="uploadSelectedFile"
            >
              <i class="fa-solid fa-upload mr-2"></i>
              {{ isUploadingVoice ? '上传中...' : '上传音频' }}
            </button>
            <div class="upload-hint">支持 MP3、WAV、AAC、OGG 等常见音频格式</div>
          </div>
        </div>

        <div class="dialog-footer">
          <button class="layui-btn layui-btn-primary" @click="showVoiceDialog = false">关闭</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '../api'
import MdEditor from '../components/MdEditor.vue'
import { alert as showAlert, error as showError } from '../utils/modal'

export default {
  name: 'CreateArticle',
  components: {
    MdEditor
  },
  setup() {
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
    const submitting = ref(false)
    const savingDraft = ref(false)
    const voiceUrl = ref('')
    const isRecording = ref(false)
    const isPaused = ref(false)
    const recordingDuration = ref(0)
    const isUploadingVoice = ref(false)

    const showVoiceDialog = ref(false)
    const voiceTab = ref('record')
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

    const loadCategories = async () => {
      try {
        const response = await api.get('/categories')
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
        const response = await api.get(`/articles/${articleId.value}`)
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
        await showAlert('请输入标题')
        return
      }

      savingDraft.value = true
      try {
        const submitData = {
          ...form.value,
          voice_url: voiceUrl.value,
          status: 'draft'
        }

        if (isEdit.value) {
          await api.put(`/articles/${articleId.value}`, submitData)
        } else {
          await api.post('/articles', submitData)
        }

        await showAlert('草稿保存成功')
        await router.push('/profile')
      } catch (error) {
        console.error('保存草稿失败', error)
        await showError('保存失败: ' + (error.message || '未知错误'))
      } finally {
        savingDraft.value = false
      }
    }

    const submitArticle = async () => {
      if (!form.value.title.trim()) {
        await showAlert('请输入标题')
        return
      }
      if (!form.value.content.trim()) {
        await showAlert('请输入内容')
        return
      }

      submitting.value = true

      try {
        const submitData = {
          ...form.value,
          voice_url: voiceUrl.value
        }

        let response
        if (isEdit.value) {
          response = await api.put(`/articles/${articleId.value}`, submitData)
        } else {
          response = await api.post('/articles', submitData)
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
        submitting.value = false
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

        const response = await api.post('/upload/voice', formData, {
          headers: { 'Content-Type': 'multipart/form-data' }
        })

        voiceUrl.value = response.data.url
        showVoiceDialog.value = false
        await showAlert('语音上传成功')
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

    const handleFileChange = (event) => {
      const file = event.target.files[0]
      if (file) {
        voiceFile.value = file
      }
    }

    const uploadSelectedFile = async () => {
      if (!voiceFile.value) return

      isUploadingVoice.value = true
      try {
        const formData = new FormData()
        formData.append('voice', voiceFile.value)

        const response = await api.post('/upload/voice', formData, {
          headers: { 'Content-Type': 'multipart/form-data' }
        })

        voiceUrl.value = response.data.url
        showVoiceDialog.value = false
        voiceFile.value = null
        await showAlert('音频上传成功')
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

    return {
      categories,
      form,
      isEdit,
      submitting,
      savingDraft,
      voiceUrl,
      isRecording,
      isPaused,
      recordingDuration,
      isUploadingVoice,
      submitArticle,
      saveDraft,
      handleVoiceInput,
      showVoiceDialog,
      voiceTab,
      voiceFile,
      recordedBlob,
      recordedAudioUrl,
      formatDuration,
      startRecording,
      pauseRecording,
      resumeRecording,
      stopRecording,
      confirmVoice,
      discardRecording,
      uploadSelectedFile,
      removeVoice,
      handleFileChange
    }
  }
}
</script>

<style scoped>
.create-article-page {
  padding: 24px 0;
  min-height: 100vh;
  background: #f5f5f5;
}

.article-container {
  max-width: 900px;
  margin: 0 auto;
  padding: 0 15px;
}

.article-card {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
}

.article-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;
}

.back-btn {
  background: none;
  border: none;
  font-size: 20px;
  color: #666;
  cursor: pointer;
  padding: 8px;
  border-radius: 6px;
  
  &:hover {
    background: #f5f5f5;
    color: var(--primary);
  }
}

.article-title {
  font-size: 20px;
  font-weight: 600;
  color: var(--primary);
  margin: 0;
}

.divider {
  height: 1px;
  background: #f0f0f0;
  margin-bottom: 24px;
}

.article-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-label {
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

.input-group {
  position: relative;
  
  i {
    position: absolute;
    left: 12px;
    top: 50%;
    transform: translateY(-50%);
    font-size: 14px;
    color: #999;
  }
  
  .layui-input {
    padding-left: 40px;
  }
}

.content-card {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 16px;
}

.content-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.content-label {
  font-size: 14px;
  color: #666;
}

.voice-alert {
  background: rgba(30, 159, 255, 0.08);
  border-left: 4px solid var(--primary);
  border-radius: 8px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  
  &.recording {
    background: rgba(255, 77, 79, 0.08);
    border-left-color: #FF4D4F;
  }
}

.alert-content {
  display: flex;
  align-items: center;
  gap: 12px;
}

.alert-content i {
  font-size: 20px;
}

.text-error {
  color: #FF4D4F;
}

.alert-info {
  flex: 1;
}

.alert-title {
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

.alert-subtitle {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.close-btn {
  background: none;
  border: none;
  font-size: 16px;
  color: #999;
  cursor: pointer;
  padding: 4px;
  
  &:hover {
    color: #666;
  }
}

.progress-bar {
  height: 4px;
  background: #e8e8e8;
  border-radius: 2px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: #FF4D4F;
  border-radius: 2px;
  transition: width 1s linear;
}

.form-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 16px;
}

.checkbox-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #666;
  cursor: pointer;
  
  input {
    width: 18px;
    height: 18px;
  }
}

.action-buttons {
  display: flex;
  gap: 8px;
}

.voice-dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 24px;
}

.voice-dialog {
  background: white;
  border-radius: 12px;
  width: 100%;
  max-width: 500px;
  overflow: hidden;
}

.dialog-header {
  background: var(--primary);
  color: white;
  padding: 16px 24px;
  font-size: 16px;
  font-weight: 600;
}

.dialog-body {
  padding: 24px;
}

.voice-tabs {
  display: flex;
  gap: 0;
  border-bottom: 1px solid #f0f0f0;
  margin-bottom: 20px;
}

.tab-btn {
  flex: 1;
  background: none;
  border: none;
  padding: 12px;
  font-size: 14px;
  color: #666;
  cursor: pointer;
  border-bottom: 2px solid transparent;
  transition: all 0.2s;
  
  &:hover {
    color: var(--primary);
  }
  
  &.active {
    color: var(--primary);
    border-bottom-color: var(--primary);
  }
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

.text-muted {
  color: #999;
}

.time-display {
  font-size: 2.5rem;
  font-weight: 600;
  font-family: 'Roboto Mono', monospace;
  color: #333;
  margin-top: 16px;
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

.audio-preview {
  width: 100%;
  max-width: 350px;
  height: 40px;
}

.preview-actions {
  display: flex;
  gap: 8px;
  justify-content: center;
  margin-top: 16px;
}

.upload-section {
  padding: 16px 0;
}

.upload-input-wrapper {
  border: 2px dashed #e8e8e8;
  border-radius: 8px;
  padding: 24px;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s;
  
  &:hover {
    border-color: var(--primary);
    background: rgba(30, 159, 255, 0.05);
  }
  
  i {
    font-size: 32px;
    color: #999;
    margin-bottom: 12px;
    display: block;
  }
  
  span {
    font-size: 14px;
    color: #666;
  }
}

.upload-input {
  display: none;
}

.upload-hint {
  font-size: 12px;
  color: #999;
  text-align: center;
  margin-top: 12px;
}

.dialog-footer {
  padding: 16px 24px;
  border-top: 1px solid #f0f0f0;
  text-align: right;
}

.w-full {
  width: 100%;
}

.mr-2 {
  margin-right: 8px;
}
</style>