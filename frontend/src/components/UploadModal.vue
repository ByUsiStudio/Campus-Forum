<template>
  <el-dialog
    :model-value="show"
    :width="'500px'"
    :close-on-click-modal="false"
    append-to-body
    @update:model-value="val => emit('update:show', val)"
  >
    <template #header>
      <div class="modal-title">
        <el-icon class="title-icon" :size="24">
          <Picture v-if="uploadType === 'image'" />
          <VideoCamera v-else />
        </el-icon>
        <span>{{ uploadType === 'image' ? '上传图片' : '上传视频' }}</span>
      </div>
    </template>

    <div class="modal-body">
      <div class="upload-area" @click="triggerFileInput" @dragover.prevent @drop.prevent="handleDrop">
        <input
          ref="fileInput"
          type="file"
          :accept="uploadType === 'image' ? 'image/*' : 'video/*'"
          style="display: none"
          @change="handleFileSelect"
        />
        <el-icon :size="48" color="#6750a4" class="upload-icon">
          <PictureFilled v-if="uploadType === 'image'" />
          <VideoCameraFilled v-else />
        </el-icon>
        <div class="upload-text">点击或拖拽文件到此处</div>
        <div class="upload-hint">
          {{ uploadType === 'image' ? '支持 JPG、PNG、GIF 格式' : '支持 MP4、WebM 格式' }}
        </div>
      </div>

      <div v-if="uploading" class="upload-progress">
        <el-progress :percentage="progress" :stroke-width="8" />
        <div class="progress-text">{{ progress }}%</div>
      </div>

      <div v-if="uploadedFiles.length > 0" class="uploaded-list">
        <div v-for="(file, index) in uploadedFiles" :key="index" class="uploaded-item">
          <img v-if="uploadType === 'image'" :src="file.url" class="uploaded-preview" />
          <div v-else class="video-preview">
            <el-icon :size="32" color="#6750a4">
              <VideoPlay />
            </el-icon>
          </div>
          <div class="uploaded-info">
            <div class="uploaded-name">{{ file.name }}</div>
            <div class="uploaded-url">{{ file.url }}</div>
          </div>
          <el-button
            text
            type="danger"
            size="small"
            @click="removeFile(index)"
          >
            <el-icon><Delete /></el-icon>
          </el-button>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="modal-actions">
        <el-button @click="close">取消</el-button>
        <el-button
          type="primary"
          @click="confirm"
          :disabled="uploadedFiles.length === 0 || uploading"
        >
          确认插入
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref } from 'vue'
import {
  Picture,
  VideoCamera,
  PictureFilled,
  VideoCameraFilled,
  VideoPlay,
  Delete
} from '@element-plus/icons-vue'

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

const fileInput = ref(null)
const uploading = ref(false)
const progress = ref(0)
const uploadedFiles = ref([])

const triggerFileInput = () => {
  fileInput.value?.click()
}

const handleFileSelect = (event) => {
  const file = event.target.files[0]
  if (file) {
    uploadFile(file)
  }
  event.target.value = ''
}

const handleDrop = (event) => {
  const file = event.dataTransfer.files[0]
  if (file) {
    uploadFile(file)
  }
}

const uploadFile = (file) => {
  uploading.value = true
  progress.value = 0

  const formData = new FormData()
  formData.append(props.uploadType, file)

  return new Promise((resolve) => {
    const xhr = new XMLHttpRequest()
    const token = localStorage.getItem('token')
    xhr.open('POST', `/api/upload/${props.uploadType}`)

    if (token) {
      xhr.setRequestHeader('Authorization', `Bearer ${token}`)
    }

    xhr.upload.onprogress = (event) => {
      if (event.lengthComputable) {
        progress.value = Math.min(99, Math.round((event.loaded / event.total) * 100))
      }
    }

    xhr.onload = () => {
      try {
        const data = JSON.parse(xhr.responseText)
        if (xhr.status >= 200 && xhr.status < 300 && data.url) {
          progress.value = 100
          uploadedFiles.value.push({
            name: file.name,
            url: data.url
          })
        } else {
          progress.value = 0
          console.error('Upload failed:', data.error || data.message || `HTTP ${xhr.status}`)
        }
      } catch (error) {
        progress.value = 0
        console.error('Upload failed:', error)
      } finally {
        uploading.value = false
        resolve()
      }
    }

    xhr.onerror = () => {
      progress.value = 0
      uploading.value = false
      console.error('Upload failed: network error')
      resolve()
    }

    xhr.onabort = () => {
      progress.value = 0
      uploading.value = false
      resolve()
    }

    xhr.send(formData)
  })
}

const removeFile = (index) => {
  uploadedFiles.value.splice(index, 1)
}

const confirm = () => {
  if (uploadedFiles.value.length > 0) {
    emit('upload-success', {
      type: props.uploadType,
      files: uploadedFiles.value
    })
    close()
  }
}

const close = () => {
  uploadedFiles.value = []
  progress.value = 0
  uploading.value = false
  emit('update:show', false)
}
</script>

<style scoped>
.modal-title {
  display: flex;
  align-items: center;
  gap: 12px;
  font-weight: 600;
  font-size: 16px;
  color: #1a1a2e;
}

.title-icon {
  width: 36px;
  height: 36px;
  padding: 6px;
  border-radius: 10px;
  background: rgba(103, 80, 164, 0.1);
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.modal-body {
  padding: 4px 0;
}

.upload-area {
  border: 2px dashed rgba(103, 80, 164, 0.3);
  border-radius: 12px;
  padding: 32px;
  text-align: center;
  cursor: pointer;
  transition: border-color 0.2s, background-color 0.2s;
  background: #fafbfc;
}

.upload-area:hover {
  border-color: rgba(103, 80, 164, 0.6);
  background: #f5f6ff;
}

.upload-icon {
  margin-bottom: 12px;
}

.upload-text {
  font-size: 1rem;
  font-weight: 500;
  color: #1a1a2e;
  margin-bottom: 8px;
}

.upload-hint {
  font-size: 0.85rem;
  color: #6b7280;
}

.upload-progress {
  margin-top: 20px;
}

.progress-text {
  text-align: center;
  font-size: 0.85rem;
  color: #6b7280;
  margin-top: 8px;
}

.uploaded-list {
  margin-top: 20px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.uploaded-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #f8f9ff;
  border-radius: 10px;
}

.uploaded-preview {
  width: 60px;
  height: 60px;
  object-fit: cover;
  border-radius: 8px;
}

.video-preview {
  width: 60px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #e5e7eb;
  border-radius: 8px;
}

.uploaded-info {
  flex: 1;
  min-width: 0;
}

.uploaded-name {
  font-weight: 500;
  color: #1a1a2e;
  font-size: 0.9rem;
}

.uploaded-url {
  font-size: 0.75rem;
  color: #6b7280;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-top: 4px;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
