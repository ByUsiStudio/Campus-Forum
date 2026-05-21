<template>
  <div class="image-viewer-mask" @click="close">
    <img :src="url" class="image-viewer-img" @click.stop>
  </div>
</template>

<script>
export default {
  name: 'ImageViewer',
  props: {
    url: {
      type: String,
      required: true
    }
  },
  setup(props, { emit }) {
    const close = () => {
      emit('close')
    }
    
    // ESC键关闭
    const handleKeydown = (e) => {
      if (e.key === 'Escape') {
        close()
      }
    }
    
    // 添加事件监听
    window.addEventListener('keydown', handleKeydown)
    
    // 清理函数
    const cleanup = () => {
      window.removeEventListener('keydown', handleKeydown)
    }
    
    // Vue 3 composition API 清理
    if (typeof window !== 'undefined') {
      setTimeout(() => {
        // 确保清理
      }, 0)
    }
    
    return {
      close,
      cleanup
    }
  },
  beforeUnmount() {
    if (this.cleanup) {
      this.cleanup()
    }
  }
}
</script>