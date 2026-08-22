<template>
  <div class="video-player-wrapper" ref="wrapperRef">
    <video
      ref="videoRef"
      class="plyr-video"
      :poster="poster"
      playsinline
      controls
    >
      <source :src="src" type="video/mp4" />
    </video>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
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
    window.scrollTo({ top: scrollPosition, behavior: 'auto' })
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

watch(() => props.src, (newSrc) => {
  if (player && newSrc) {
    player.source = {
      type: 'video',
      sources: [
        {
          src: newSrc,
          type: 'video/mp4'
        }
      ]
    }
  }
})

watch(() => props.poster, (newPoster) => {
  if (player && newPoster) {
    player.poster = newPoster
  }
})

onMounted(() => {
  initPlayer()
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
  border-radius: var(--campus-radius);
  overflow: hidden;
  background: #000;
  box-shadow: var(--campus-shadow);
}

:deep(.plyr) {
  border-radius: var(--campus-radius);
}

:deep(.plyr--video) {
  background: #000;
}

:deep(.plyr__control--overlaid) {
  background: var(--campus-primary);
  border-radius: 50%;
  box-shadow: 0 6px 16px rgba(79, 110, 247, 0.4);
  transition: transform 0.2s ease, background-color 0.2s ease;
}

:deep(.plyr__control--overlaid:hover) {
  background: var(--campus-primary-dark);
  transform: scale(1.06);
}

:deep(.plyr__controls) {
  background: linear-gradient(transparent, rgba(0, 0, 0, 0.78));
  padding: 12px;
}

:deep(.plyr--video .plyr__control) {
  color: #fff;
}

:deep(.plyr__control) {
  opacity: 0.92;
  transition: opacity 0.2s;
}

:deep(.plyr__control:hover) {
  opacity: 1;
}

:deep(.plyr__controls .plyr__control--pressed) {
  color: var(--campus-primary-light);
}

:deep(.plyr__progress__container) {
  flex: 1;
}

:deep(.plyr__progress) {
  height: 5px;
}

:deep(.plyr__progress input[type='range']::-webkit-slider-runnable-track) {
  border-radius: 999px;
}

:deep(.plyr__progress__buffer) {
  background: rgba(255, 255, 255, 0.25);
  border-radius: 999px;
}

:deep(.plyr__progress__played) {
  background: var(--campus-primary-light);
}

:deep(.plyr__volume) {
  max-width: 80px;
}

:deep(.plyr__time) {
  font-size: 12px;
}

:deep(.plyr__menu) {
  background: rgba(15, 23, 42, 0.95);
  border-radius: var(--campus-radius-sm);
  border: 1px solid rgba(255, 255, 255, 0.08);
}

:deep(.plyr__menu__item) {
  padding: 8px 16px;
}

:deep(.plyr__menu__item:hover) {
  background: rgba(255, 255, 255, 0.08);
}

:deep(.plyr__menu__item--selected) {
  color: var(--campus-primary-light);
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
  background: var(--campus-primary);
}

:deep(.plyr--fullscreen .plyr__controls) {
  background: linear-gradient(transparent, rgba(0, 0, 0, 0.85));
}

@media (max-width: 768px) {
  .video-player-wrapper {
    max-width: 100%;
    border-radius: var(--campus-radius-sm);
  }

  :deep(.plyr) {
    border-radius: var(--campus-radius-sm);
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
