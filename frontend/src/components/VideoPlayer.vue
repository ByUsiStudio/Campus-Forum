<template>
  <div class="video-player-container">
    <div class="video-wrapper">
      <!-- 视频元素 -->
      <video
        ref="videoRef"
        :src="src"
        :poster="poster"
        :controls="true"
        :autoplay="false"
        :loop="false"
        :muted="false"
        playsinline
        preload="metadata"
        @loadstart="handleLoadStart"
        @progress="handleProgress"
        @loadedmetadata="handleLoadedMetadata"
        @error="handleError"
      ></video>
      
      <!-- 加载覆盖层 -->
      <div class="video-overlay" v-if="isLoading">
        <div class="loading-content">
          <svg class="loading-spinner" viewBox="0 0 50 50">
            <circle class="path" cx="25" cy="25" r="20" fill="none" stroke-width="4"></circle>
          </svg>
          <span class="loading-text">加载中...</span>
        </div>
        <!-- 加载进度条 -->
        <div class="loading-progress-bar">
          <div class="loading-progress" :style="{ width: bufferedPercent + '%' }"></div>
        </div>
      </div>
      
      <!-- 错误覆盖层 -->
      <div class="video-overlay video-error" v-if="hasError">
        <span class="error-icon">⚠️</span>
        <span class="error-text">{{ errorMessage }}</span>
        <button class="retry-btn" @click="retryLoad">重试</button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'

export default {
  name: 'VideoPlayer',
  props: {
    src: {
      type: String,
      required: true
    },
    poster: {
      type: String,
      default: ''
    }
  },
  setup(props) {
    const videoRef = ref(null)
    const isLoading = ref(true)
    const hasError = ref(false)
    const errorMessage = ref('')
    const bufferedPercent = ref(0)
    
    const handleLoadStart = () => {
      isLoading.value = true
      hasError.value = false
      bufferedPercent.value = 0
    }
    
    const handleProgress = () => {
      const video = videoRef.value
      if (video && video.buffered.length > 0) {
        const bufferedEnd = video.buffered.end(video.buffered.length - 1)
        const duration = video.duration
        if (duration > 0) {
          bufferedPercent.value = Math.min((bufferedEnd / duration) * 100, 99)
        }
      }
    }
    
    const handleLoadedMetadata = () => {
      isLoading.value = false
      bufferedPercent.value = 100
    }
    
    const handleError = (event) => {
      const video = event.target
      isLoading.value = false
      
      switch (video.error?.code) {
        case MediaError.MEDIA_ERR_ABORTED:
          errorMessage.value = '视频加载被中断'
          break
        case MediaError.MEDIA_ERR_NETWORK:
          errorMessage.value = '网络错误，无法加载视频'
          break
        case MediaError.MEDIA_ERR_DECODE:
          errorMessage.value = '视频解码失败'
          break
        case MediaError.MEDIA_ERR_SRC_NOT_SUPPORTED:
          errorMessage.value = '不支持的视频格式'
          break
        default:
          errorMessage.value = '视频加载失败，请重试'
      }
      
      hasError.value = true
      console.error('视频加载错误:', props.src, video.error)
    }
    
    const retryLoad = () => {
      hasError.value = false
      isLoading.value = true
      bufferedPercent.value = 0
      
      if (videoRef.value) {
        videoRef.value.load()
      }
    }
    
    return {
      videoRef,
      isLoading,
      hasError,
      errorMessage,
      bufferedPercent,
      handleLoadStart,
      handleProgress,
      handleLoadedMetadata,
      handleError,
      retryLoad
    }
  }
}
</script>

<style scoped>
.video-player-container {
  width: 100%;
  max-width: 800px;
  margin: 1em auto;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  background: #000;
}

.video-wrapper {
  position: relative;
  width: 100%;
  aspect-ratio: 16 / 9;
  background: #000;
}

.video-wrapper video {
  width: 100%;
  height: 100%;
  display: block;
  object-fit: contain;
}

/* 视频覆盖层 */
.video-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 10;
}

/* 加载状态 */
.loading-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.loading-spinner {
  animation: rotate 2s linear infinite;
  width: 50px;
  height: 50px;
}

.loading-spinner .path {
  stroke: #1976d2;
  stroke-linecap: round;
  animation: dash 1.5s ease-in-out infinite;
}

@keyframes rotate {
  100% {
    transform: rotate(360deg);
  }
}

@keyframes dash {
  0% {
    stroke-dasharray: 1, 150;
    stroke-dashoffset: 0;
  }
  50% {
    stroke-dasharray: 90, 150;
    stroke-dashoffset: -35;
  }
  100% {
    stroke-dasharray: 90, 150;
    stroke-dashoffset: -124;
  }
}

.loading-text {
  color: #fff;
  font-size: 14px;
}

/* 加载进度条 */
.loading-progress-bar {
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 3px;
  background: rgba(255, 255, 255, 0.2);
}

.loading-progress {
  height: 100%;
  background: #1976d2;
  transition: width 0.3s ease;
}

/* 错误状态 */
.video-error {
  background: rgba(0, 0, 0, 0.9);
}

.error-icon {
  font-size: 48px;
  margin-bottom: 12px;
}

.error-text {
  color: #fff;
  font-size: 16px;
  margin-bottom: 16px;
}

.retry-btn {
  padding: 10px 24px;
  background: #1976d2;
  color: #fff;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: background 0.3s ease;
}

.retry-btn:hover {
  background: #1565c0;
}
</style>