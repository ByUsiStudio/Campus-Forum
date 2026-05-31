<template>
  <div>
    <v-card class="pa-6" max-width="900" style="margin: 0 auto;">
      <v-card-title class="text-h5 mb-4" style="color: rgb(var(--v-theme-primary));">
        {{ isEdit ? '编辑文章' : '发布新文章' }}
      </v-card-title>

      <v-form @submit.prevent="submitArticle">
        <v-text-field
          v-model="form.title"
          label="标题"
          variant="outlined"
          required
          class="mb-4"
        ></v-text-field>

        <v-select
          v-model="form.category_id"
          :items="categories"
          item-title="name"
          item-value="id"
          label="分区"
          variant="outlined"
          class="mb-4"
        ></v-select>

        <div class="mb-4">
          <MdEditor
            v-model="form.content"
            :height="500"
            @voice-input="handleVoiceInput"
          />
        </div>

        <v-expand-transition>
          <div v-if="voiceUrl || isRecording" class="mb-4">
            <v-card variant="outlined" class="pa-3">
              <div class="d-flex align-center gap-3">
                <v-icon :color="isRecording ? 'error' : 'primary'">
                  {{ isRecording ? 'mdi-microphone' : 'mdi-volume-high' }}
                </v-icon>
                <div class="flex-grow-1">
                  <div class="text-subtitle-2">
                    {{ isRecording ? '录音中...' : '语音已添加' }}
                  </div>
                  <div v-if="voiceUrl && !isRecording" class="text-caption text-medium-emphasis">
                    {{ voiceUrl.split('/').pop() }}
                  </div>
                </div>
                <v-btn
                  v-if="!isRecording"
                  variant="text"
                  color="error"
                  size="small"
                  @click="removeVoice"
                >
                  删除
                </v-btn>
              </div>
              <v-progress-linear
                v-if="isRecording"
                :model-value="recordingDuration"
                color="error"
                height="4"
                class="mt-2"
              />
            </v-card>
          </div>
        </v-expand-transition>

        <div class="d-flex align-center gap-4 mb-4">
          <v-btn
            variant="outlined"
            color="primary"
            size="small"
            @click="showVoiceDialog = true"
          >
            <v-icon left>mdi-microphone</v-icon>
            添加语音
          </v-btn>

          <v-checkbox
            v-model="form.is_anonymous"
            label="匿名发布"
            color="primary"
            hide-details
            class="mt-0"
          ></v-checkbox>
        </div>

        <div class="d-flex gap-4">
          <v-btn
            type="submit"
            color="primary"
            size="large"
            :loading="submitting"
          >
            {{ submitting ? '提交中...' : (isEdit ? '更新文章' : '发布文章') }}
          </v-btn>
        </div>
      </v-form>
    </v-card>

    <v-dialog v-model="showVoiceDialog" max-width="500">
      <v-card class="voice-dialog-card">
        <v-card-title class="voice-dialog-title">
          <v-icon class="mr-2">mdi-microphone</v-icon>
          添加语音
        </v-card-title>
        <v-card-text class="voice-dialog-body">
          <v-tabs v-model="voiceTab" color="primary" grow>
            <v-tab value="record">
              <v-icon start>mdi-microphone</v-icon>
              录音
            </v-tab>
            <v-tab value="upload">
              <v-icon start>mdi-upload</v-icon>
              上传
            </v-tab>
          </v-tabs>

          <v-window v-model="voiceTab" class="mt-4">
            <v-window-item value="record">
              <div class="record-section">
                <div class="record-time">
                  <v-icon :color="isRecording ? 'error' : 'grey'" size="48" class="mb-3">
                    mdi-microphone
                  </v-icon>
                  <div class="time-display">
                    {{ formatDuration(recordingDuration) }}
                  </div>
                  <div class="record-status">
                    {{ isPaused ? '已暂停' : (isRecording ? '录音中...' : '点击开始录音') }}
                  </div>
                </div>

                <div class="record-controls">
                  <v-btn
                    v-if="!isRecording && !recordedBlob"
                    color="error"
                    size="large"
                    @click="startRecording"
                  >
                    <v-icon start>mdi-microphone</v-icon>
                    开始录音
                  </v-btn>

                  <template v-else>
                    <v-btn
                      v-if="isRecording"
                      color="warning"
                      size="large"
                      @click="pauseRecording"
                    >
                      <v-icon start>mdi-pause</v-icon>
                      暂停
                    </v-btn>
                    <v-btn
                      v-if="isPaused"
                      color="success"
                      size="large"
                      @click="resumeRecording"
                    >
                      <v-icon start>mdi-play</v-icon>
                      继续
                    </v-btn>
                    <v-btn
                      color="error"
                      size="large"
                      @click="stopRecording"
                    >
                      <v-icon start>mdi-stop</v-icon>
                      结束
                    </v-btn>
                  </template>
                </div>

                <div v-if="recordedBlob" class="recorded-preview">
                  <audio ref="recordedAudioRef" :src="recordedAudioUrl" controls class="audio-preview" />
                  <div class="d-flex gap-2 mt-3">
                    <v-btn color="primary" @click="confirmVoice" :loading="isUploadingVoice">
                      使用此录音
                    </v-btn>
                    <v-btn variant="outlined" @click="discardRecording">
                      丢弃
                    </v-btn>
                  </div>
                </div>
              </div>
            </v-window-item>

            <v-window-item value="upload">
              <div class="upload-section">
                <v-file-input
                  ref="voiceFileInput"
                  v-model="voiceFile"
                  label="选择音频文件"
                  accept="audio/*"
                  variant="outlined"
                  prepend-icon="mdi-file-music"
                  show-size
                  class="mb-4"
                />
                <v-btn
                  color="primary"
                  block
                  :disabled="!voiceFile"
                  :loading="isUploadingVoice"
                  @click="uploadSelectedFile"
                >
                  <v-icon start>mdi-upload</v-icon>
                  上传音频
                </v-btn>
                <div class="text-caption text-medium-emphasis mt-3">
                  支持 MP3、WAV、AAC、OGG 等常见音频格式
                </div>
              </div>
            </v-window-item>
          </v-window>
        </v-card-text>
        <v-card-actions class="voice-dialog-actions">
          <v-btn variant="text" @click="showVoiceDialog = false">关闭</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted, computed } from 'vue'
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
    const recordedAudioRef = ref(null)

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

    const escapeRegExp = (value) => {
      return value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
    }

    const insertVoiceTagIntoContent = (url) => {
      if (!url) return
      const trimmedContent = form.value.content.trim()
      const audioMarker = `<audio controls src="${url}"></audio>`
      if (trimmedContent.includes(audioMarker) || trimmedContent.includes(`src="${url}"`)) {
        return
      }
      form.value.content = trimmedContent
        ? `${trimmedContent}\n\n${audioMarker}\n`
        : `${audioMarker}\n`
    }

    const removeVoiceTagFromContent = (url) => {
      if (!url) return
      const escapedUrl = escapeRegExp(url)
      const audioRegex = new RegExp(`\\n?<audio[^>]*src=\"${escapedUrl}\"[^>]*>(?:<\\/audio>)?\\n?`, 'g')
      form.value.content = form.value.content.replace(audioRegex, '').trim()
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

      if (voiceUrl.value) {
        insertVoiceTagIntoContent(voiceUrl.value)
      }

      submitting.value = true

      try {
        let response
        const submitData = {
          ...form.value,
          voice_url: voiceUrl.value
        }

        if (isEdit.value) {
          response = await api.put(`/articles/${articleId.value}`, submitData)
        } else {
          response = await api.post('/articles', submitData)
        }

        const newArticleId = response.data.article?.id
        if (!newArticleId) {
          throw new Error('文章ID不存在')
        }

        console.log('准备跳转到文章页面:', `/article/${newArticleId}`)
        await router.push(`/article/${newArticleId}`)
        console.log('跳转成功')
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
        insertVoiceTagIntoContent(voiceUrl.value)
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
      removeVoiceTagFromContent(voiceUrl.value)
      voiceUrl.value = ''
    }

    const resetVoiceDialog = () => {
      voiceTab.value = 'record'
      voiceFile.value = null
      discardRecording()
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
      voiceUrl,
      isRecording,
      isPaused,
      recordingDuration,
      isUploadingVoice,
      submitArticle,
      handleVoiceInput,
      showVoiceDialog,
      voiceTab,
      voiceFile,
      recordedBlob,
      recordedAudioUrl,
      recordedAudioRef,
      formatDuration,
      startRecording,
      pauseRecording,
      resumeRecording,
      stopRecording,
      confirmVoice,
      discardRecording,
      uploadSelectedFile,
      removeVoice,
      resetVoiceDialog
    }
  }
}
</script>

<style scoped>
.voice-dialog-card {
  border-radius: 16px;
}

.voice-dialog-title {
  background: linear-gradient(135deg, rgb(var(--v-theme-primary)), rgb(var(--v-theme-secondary)));
  color: white;
  padding: 16px 24px;
  font-size: 1.1rem;
}

.voice-dialog-body {
  padding: 16px 24px;
}

.voice-dialog-actions {
  padding: 8px 24px 16px;
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
  font-family: 'Roboto Mono', monospace;
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

.audio-preview {
  width: 100%;
  max-width: 350px;
  height: 40px;
}

.upload-section {
  padding: 16px 0;
}
</style>
