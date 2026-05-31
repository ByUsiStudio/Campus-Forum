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
            :loading="isRecording || isUploadingVoice"
            @click="toggleRecording"
          >
            <v-icon left>{{ isRecording ? 'mdi-stop' : 'mdi-microphone' }}</v-icon>
            {{ isRecording ? '停止录音' : '添加语音' }}
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
    const voiceUrl = ref('')
    const isRecording = ref(false)
    const recordingDuration = ref(0)
    const isUploadingVoice = ref(false)

    let mediaRecorder = null
    let audioChunks = []
    let recordingTimer = null

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
        formData.append('voice', audioBlob, 'voice.webm')

        const response = await api.post('/upload/voice', formData, {
          headers: { 'Content-Type': 'multipart/form-data' }
        })

        voiceUrl.value = response.data.url
        await showAlert('语音上传成功')
      } catch (error) {
        console.error('语音上传失败', error)
        await showError('语音上传失败')
      } finally {
        isUploadingVoice.value = false
      }
    }

    const toggleRecording = async () => {
      if (isRecording.value) {
        stopRecording()
      } else {
        await startRecording()
      }
    }

    const startRecording = async () => {
      try {
        const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
        mediaRecorder = new MediaRecorder(stream)
        audioChunks = []

        mediaRecorder.ondataavailable = (event) => {
          audioChunks.push(event.data)
        }

        mediaRecorder.onstop = async () => {
          const audioBlob = new Blob(audioChunks, { type: 'audio/webm' })
          stream.getTracks().forEach(track => track.stop())
          await uploadVoice(audioBlob)
        }

        mediaRecorder.start()
        isRecording.value = true
        recordingDuration.value = 0

        recordingTimer = setInterval(() => {
          recordingDuration.value += 1
        }, 1000)

      } catch (error) {
        console.error('无法访问麦克风', error)
        await showError('无法访问麦克风，请检查权限设置')
      }
    }

    const stopRecording = () => {
      if (mediaRecorder && mediaRecorder.state !== 'inactive') {
        mediaRecorder.stop()
      }
      isRecording.value = false
      if (recordingTimer) {
        clearInterval(recordingTimer)
        recordingTimer = null
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
      recordingDuration,
      isUploadingVoice,
      submitArticle,
      handleVoiceInput,
      toggleRecording,
      removeVoice
    }
  }
}
</script>
