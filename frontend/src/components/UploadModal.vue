<template>
  <div v-if="show" class="upload-modal-overlay" @click="close">
    <div class="upload-modal" @click.stop>
      <div class="modal-header">
        <i :class="uploadType === 'image' ? 'fa-solid fa-image' : 'fa-solid fa-video'" class="modal-icon"></i>
        <span class="modal-title">{{ uploadType === 'image' ? '上传图片' : '上传视频' }}</span>
        <button class="close-btn" @click="close">
          <i class="fa-solid fa-xmark"></i>
        </button>
      </div>
      
      <div class="modal-body">
        <div class="upload-area" @click="triggerFileInput" @dragover.prevent @drop.prevent="handleDrop">
          <input
            ref="fileInput"
            type="file"
            :accept="uploadType === 'image' ? 'image/*' : 'video/*'"
            style="display: none"
            @change="handleFileSelect"
          />
          <i :class="uploadType === 'image' ? 'fa-solid fa-image-plus' : 'fa-solid fa-video-plus'" class="upload-icon"></i>
          <div class="upload-text">点击或拖拽文件到此处</div>
          <div class="upload-hint">
            {{ uploadType === 'image' ? '支持 JPG、PNG、GIF 格式' : '支持 MP4、WebM 格式' }}
          </div>
        </div>

        <div v-if="uploading" class="upload-progress">
          <div class="progress-bar">
            <div class="progress-fill" :style="{ width: progress + '%' }"></div>
          </div>
          <div class="progress-text">{{ progress }}%</div>
        </div>

        <div v-if="uploadedFiles.length > 0" class="uploaded-list">
          <div v-for="(file, index) in uploadedFiles" :key="index" class="uploaded-item">
            <img v-if="uploadType === 'image'" :src="file.url" class="uploaded-preview" />
            <div v-else class="video-preview">
              <i class="fa-solid fa-play-circle"></i>
            </div>
            <div class="uploaded-info">
              <div class="uploaded-name">{{ file.name }}</div>
              <div class="uploaded-url">{{ file.url }}</div>
            </div>
            <button class="delete-btn" @click="removeFile(index)">
              <i class="fa-solid fa-trash"></i>
            </button>
          </div>
        </div>
      </div>

      <div class="modal-footer">
        <button class="cancel-btn" @click="close">取消</button>
        <button
          class="confirm-btn"
          @click="confirm"
          :disabled="uploadedFiles.length === 0 || uploading"
        >
          确认插入
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, defineProps, defineEmits } from 'vue'

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

const uploadFile = async (file) => {
  uploading.value = true
  progress.value = 0

  const formData = new FormData()
  formData.append(props.uploadType, file)

  try {
    const token = localStorage.getItem('token')
    const response = await fetch(`/api/upload/${props.uploadType}`, {
      method: 'POST',
      headers: {
        'Authorization': token ? `Bearer ${token}` : ''
      },
      body: formData
    })

    const data = await response.json()
    if (data.url) {
      uploadedFiles.value.push({
        name: file.name,
        url: data.url
      })
    } else {
      console.error('Upload failed:', data.error)
    }
  } catch (error) {
    console.error('Upload failed:', error)
  } finally {
    uploading.value = false
    progress.value = 100
  }
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
  z-index: 9999;
}

.upload-modal {
  background: #fff;
  border-radius: 16px;
  width: 90%;
  max-width: 500px;
  overflow: hidden;
}

.modal-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 20px 24px;
  background: linear-gradient(135deg, #f8f9ff 0%, #fff 100%);
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.modal-icon {
  width: 36px;
  height: 36px;
  padding: 6px;
  border-radius: 10px;
  background: rgba(30, 159, 255, 0.1);
  color: #1E9FFF;
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal-title {
  flex: 1;
  font-size: 16px;
  font-weight: 600;
}

.close-btn {
  background: none;
  border: none;
  font-size: 20px;
  color: #999;
  cursor: pointer;
  padding: 4px;

  &:hover {
    color: #666;
  }
}

.modal-body {
  padding: 24px;
}

.upload-area {
  border: 2px dashed rgba(30, 159, 255, 0.3);
  border-radius: 12px;
  padding: 32px;
  text-align: center;
  cursor: pointer;
  transition: border-color 0.2s, background-color 0.2s;
  background: #fafbfc;

  &:hover {
    border-color: rgba(30, 159, 255, 0.6);
    background: #f5f6ff;
  }
}

.upload-icon {
  font-size: 48px;
  color: #1E9FFF;
  margin-bottom: 12px;
}

.upload-text {
  font-size: 1rem;
  font-weight: 500;
  color: #333;
  margin-bottom: 8px;
}

.upload-hint {
  font-size: 0.85rem;
  color: #666;
}

.upload-progress {
  margin-top: 20px;
}

.progress-bar {
  height: 8px;
  background: #e8e8e8;
  border-radius: 4px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: #1E9FFF;
  border-radius: 4px;
  transition: width 0.3s ease;
}

.progress-text {
  text-align: center;
  font-size: 0.85rem;
  color: #666;
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
  font-size: 24px;
  color: #1E9FFF;
}

.uploaded-info {
  flex: 1;
  min-width: 0;
}

.uploaded-name {
  font-weight: 500;
  color: #333;
  font-size: 0.9rem;
}

.uploaded-url {
  font-size: 0.75rem;
  color: #666;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-top: 4px;
}

.delete-btn {
  background: none;
  border: none;
  color: #FF5722;
  cursor: pointer;
  padding: 8px;
  border-radius: 6px;

  &:hover {
    background: rgba(255, 87, 34, 0.1);
  }
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px;
  border-top: 1px solid #f0f0f0;
}

.cancel-btn {
  padding: 8px 20px;
  background: #f5f5f5;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  color: #666;

  &:hover {
    background: #e8e8e8;
  }
}

.confirm-btn {
  padding: 8px 20px;
  background: #1E9FFF;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  color: #fff;

  &:hover:not(:disabled) {
    background: #0086E6;
  }

  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
}
</style>