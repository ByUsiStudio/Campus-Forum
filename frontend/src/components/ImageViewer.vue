<template>
  <div class="lightbox-overlay" @click="close">
    <button class="lightbox-close" @click="close">
      <v-icon>mdi-close</v-icon>
    </button>
    <img :src="url" @click.stop>
  </div>
</template>

<script>
import { onMounted, onBeforeUnmount } from 'vue'

export default {
  name: 'ImageViewer',
  props: {
    url: {
      type: String,
      required: true
    }
  },
  emits: ['close'],
  setup(props, { emit }) {
    const close = () => {
      emit('close')
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
      close
    }
  }
}
</script>
