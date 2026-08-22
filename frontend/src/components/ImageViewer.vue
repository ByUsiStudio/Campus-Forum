<template>
  <div class="image-viewer overlay" @click="close">
    <el-button
      class="image-viewer-close"
      circle
      :icon="Close"
      aria-label="关闭"
      @click.stop="close"
    />
    <el-image
      class="image-viewer-image"
      :src="url"
      :preview-src-list="[url]"
      :initial-index="0"
      hide-on-click-modal
      fit="contain"
      @load="onImageLoad"
      @error="onImageError"
    >
      <template #error>
        <div class="image-viewer-error">
          <el-icon :size="48"><PictureFilled /></el-icon>
          <p class="text-secondary">图片加载失败</p>
        </div>
      </template>
    </el-image>
    <!-- 点按图片本身不关闭（交由 el-image 预览） -->
  </div>
</template>

<script>
import { onMounted, onBeforeUnmount } from 'vue'
import { Close, PictureFilled } from '@element-plus/icons-vue'

export default {
  name: 'ImageViewer',
  props: {
    url: {
      type: String,
      required: true
    }
  },
  emits: ['close'],
  components: {
    Close,
    PictureFilled
  },
  setup(props, { emit }) {
    const close = () => {
      emit('close')
    }

    const onImageLoad = () => {
      // 图片加载完成，无额外处理
    }

    const onImageError = () => {
      // 图片加载失败，展示占位提示即可
    }

    const handleKeydown = (e) => {
      if (e.key === 'Escape') {
        close()
      }
    }

    onMounted(() => {
      window.addEventListener('keydown', handleKeydown)
    })

    onBeforeUnmount(() => {
      window.removeEventListener('keydown', handleKeydown)
    })

    return {
      Close,
      close,
      onImageLoad,
      onImageError
    }
  }
}
</script>

<style scoped>
.overlay {
  position: fixed;
  inset: 0;
  z-index: 3000;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.85);
  cursor: pointer;
}

.image-viewer-close {
  position: absolute;
  top: 20px;
  right: 24px;
  z-index: 1;
}

.image-viewer-image {
  max-width: 90vw;
  max-height: 90vh;
  background: transparent;
  cursor: default;
}

.image-viewer-error {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--campus-text-secondary);
  min-width: 120px;
  min-height: 120px;
}

.image-viewer-error p {
  margin: 8px 0 0;
}
</style>
