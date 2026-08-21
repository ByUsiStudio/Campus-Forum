<template>
  <v-dialog :model-value="show" @update:model-value="val => emit('update:show', val)" max-width="500px">
    <v-card class="upload-modal">
      <v-card-title class="modal-title">
        <v-icon class="title-icon">{{ uploadType === 'image' ? 'mdi-image' : 'mdi-video' }}</v-icon>
        {{ uploadType === 'image' ? '上传图片' : '上传视频' }}
      </v-card-title>
      
      <v-card-text class="modal-body">
        <div class="upload-area" @click="triggerFileInput" @dragover.prevent @drop.prevent="handleDrop">
          <input
            ref="fileInput"
            type="file"
            :accept="uploadType === 'image' ? 'image/*' : 'video/*'"
            style="display: none"
            @change="handleFileSelect"
          />
          <v-icon size="48" color="primary" class="upload-icon">
            {{ uploadType === 'image' ? 'mdi-image-plus' : 'mdi-video-plus' }}
          </v-icon>
          <div class="upload-text">点击或拖拽文件到此处</div>
          <div class="upload-hint">
            {{ uploadType === 'image' ? '支持 JPG、PNG、GIF 格式' : '支持 MP4、WebM 格式' }}
          </div>
        </div>

        <div v-if="uploading" class="upload-progress">
          <v-progress-linear
            :value="progress"
            color="primary"
            height="8"
            rounded
          ></v-progress-linear>
          <div class="progress-text">{{ progress }}%</div>
        </div>

        <div v-if="uploadedFiles.length > 0" class="uploaded-list">
          <div v-for="(file, index) in uploadedFiles" :key="index" class="uploaded-item">
            <img v-if="uploadType === 'image'" :src="file.url" class="uploaded-preview" />
            <div v-else class="video-preview">
              <v-icon size="32" color="primary">mdi-play-circle</v-icon>
            </div>
            <div class="uploaded-info">
              <div class="uploaded-name">{{ file.name }}</div>
              <div class="uploaded-url">{{ file.url }}</div>
            </div>
            <v-btn
              variant="text"
              color="error"
              size="small"
              @click="removeFile(index)"
            >
              <v-icon>mdi-delete</v-icon>
            </v-btn>
          </div>
        </div>
      </v-card-text>

      <v-card-actions class="modal-actions">
        <v-btn variant="text" @click="close">取消</v-btn>
        <v-btn
          color="primary"
          variant="flat"
          @click="confirm"
          :disabled="uploadedFiles.length === 0 || uploading"
        >
          确认插入
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
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
.upload-modal {
  border-radius: 16px;
  overflow: hidden;
}

.modal-title {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 20px 24px;
  background: linear-gradient(135deg, #f8f9ff 0%, #fff 100%);
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  font-weight: 600;
}

.title-icon {
  width: 36px;
  height: 36px;
  padding: 6px;
  border-radius: 10px;
  background: rgba(103, 80, 164, 0.1);
}

.modal-body {
  padding: 24px;
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
  padding: 16px 24px;
  gap: 12px;
  border-top: 1px solid rgba(0, 0, 0, 0.05);
}
</style>
