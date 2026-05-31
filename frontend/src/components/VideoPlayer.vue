<template>
  <div class="video-player-wrapper" ref="wrapperRef">
    <video
      ref="videoRef"
      class="plyr-video"
      :poster="poster"
      playsinline
      controls
    />
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, watch, nextTick } from 'vue'
import Plyr from 'plyr'
import 'plyr/dist/plyr.css'

const props = defineProps({
  src: {
    type: String,
    required: true
  },
  poster: {
    type: String,
    default: ''
  },
  autoplay: {
    type: Boolean,
    default: false
  },
  loop: {
    type: Boolean,
    default: false
  },
  muted: {
    type: Boolean,
    default: false
  }
})

const videoRef = ref(null)
const wrapperRef = ref(null)
let player = null
let scrollPosition = 0

const handleFullscreenChange = () => {
  if (document.fullscreenElement) {
    scrollPosition = window.scrollY
    document.body.style.overflow = 'hidden'
  } else {
    document.body.style.overflow = ''
    window.scrollTo({ top: scrollPosition, behavior: 'instant' })
  }
}

const initPlayer = () => {
  if (!videoRef.value) return

  player = new Plyr(videoRef.value, {
    controls: [
      'play-large',
      'play',
      'progress',
      'current-time',
      'duration',
      'mute',
      'volume',
      'settings',
      'pip',
      'airplay',
      'fullscreen'
    ],
    settings: ['quality', 'speed'],
    speed: {
      selected: 1,
      options: [0.5, 0.75, 1, 1.25, 1.5, 2]
    },
    quality: {
      default: 720,
      options: [720, 480, 360]
    },
    autoplay: props.autoplay,
    loop: {
      active: props.loop
    },
    muted: props.muted,
    ratio: '16:9',
    hideControls: false,
    resetOnEnd: false,
    disableContextMenu: true,
    loadSprite: true,
    iconUrl: 'https://cdn.plyr.io/3.7.8/plyr.svg',
    fullscreen: {
      enabled: true,
      fallback: true,
      iosNative: false
    }
  })

  player.on('enterfullscreen', handleFullscreenChange)
  player.on('exitfullscreen', handleFullscreenChange)
}

const setupVideoSource = () => {
  if (!videoRef.value || !props.src) return
  videoRef.value.src = props.src
  if (props.autoplay) {
    videoRef.value.play().catch(() => {})
  }
}

watch(() => props.src, (newSrc) => {
  if (player && newSrc) {
    nextTick(() => {
      setupVideoSource()
    })
  }
})

watch(() => props.poster, (newPoster) => {
  if (player && newPoster) {
    player.poster = newPoster
  }
})

onMounted(() => {
  initPlayer()
  nextTick(() => {
    setupVideoSource()
  })
  document.addEventListener('fullscreenchange', handleFullscreenChange)
})

onBeforeUnmount(() => {
  document.removeEventListener('fullscreenchange', handleFullscreenChange)
  if (player) {
    player.destroy()
    player = null
  }
})

defineExpose({
  getInstance: () => player,
  play: () => player?.play(),
  pause: () => player?.pause()
})
</script>

<style scoped>
.video-player-wrapper {
  width: 100%;
  max-width: 800px;
  margin: 0 auto;
  border-radius: 8px;
  overflow: hidden;
  background: #000;
}

:deep(.plyr) {
  border-radius: 8px;
}

:deep(.plyr--video) {
  background: #000;
}

:deep(.plyr__control--overlaid) {
  background: rgba(var(--v-theme-primary), 0.9);
  border-radius: 50%;
}

:deep(.plyr__control--overlaid:hover) {
  background: rgba(var(--v-theme-primary), 1);
}

:deep(.plyr__controls) {
  background: linear-gradient(transparent, rgba(0, 0, 0, 0.7));
  padding: 10px;
}

:deep(.plyr__control) {
  opacity: 0.9;
  transition: opacity 0.2s;
}

:deep(.plyr__control:hover) {
  opacity: 1;
}

:deep(.plyr__progress__container) {
  flex: 1;
}

:deep(.plyr__progress) {
  height: 4px;
}

:deep(.plyr__progress--played) {
  background: rgb(var(--v-theme-primary));
}

:deep(.plyr__progress__buffer) {
  background: rgba(255, 255, 255, 0.3);
}

:deep(.plyr__volume) {
  max-width: 80px;
}

:deep(.plyr__time) {
  font-size: 12px;
}

:deep(.plyr__menu) {
  background: rgba(0, 0, 0, 0.9);
  border-radius: 8px;
}

:deep(.plyr__menu__item) {
  padding: 8px 16px;
}

:deep(.plyr__menu__item:hover) {
  background: rgba(255, 255, 255, 0.1);
}

:deep(.plyr--full-ui.plyr--fullscreen) {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  width: 100vw;
  height: 100vh;
  z-index: 9999;
}

:deep(.plyr--fullscreen .plyr__control--overlaid) {
  background: rgba(var(--v-theme-primary), 0.9);
}

:deep(.plyr--fullscreen .plyr__controls) {
  background: linear-gradient(transparent, rgba(0, 0, 0, 0.8));
}

@media (max-width: 768px) {
  .video-player-wrapper {
    max-width: 100%;
    border-radius: 0;
  }

  :deep(.plyr) {
    border-radius: 0;
  }
}

@media (max-width: 480px) {
  :deep(.plyr__controls) {
    padding: 8px;
  }

  :deep(.plyr__control) {
    padding: 8px;
  }

  :deep(.plyr__time) {
    font-size: 10px;
  }
}
</style>