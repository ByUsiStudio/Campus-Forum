<template>
  <div class="video-player-wrapper">
    <video
      ref="videoRef"
      class="video-js vjs-big-play-centered"
      :poster="poster"
      playsinline
    ></video>
  </div>
</template>

<script>
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import videojs from 'video.js'
import 'video.js/dist/video-js.css'

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
  },
  setup(props) {
    const videoRef = ref(null)
    let player = null
    
    const initPlayer = () => {
      if (!videoRef.value) return
      
      player = videojs(videoRef.value, {
        sources: [{
          src: props.src,
          type: 'video/mp4'
        }],
        autoplay: props.autoplay,
        loop: props.loop,
        muted: props.muted,
        fluid: true,
        preload: 'auto',
        controls: true,
        playbackRates: [0.5, 0.75, 1, 1.25, 1.5, 2],
        controlBar: {
          children: [
            'playToggle',
            'volumePanel',
            'currentTimeDisplay',
            'timeDivider',
            'durationDisplay',
            'progressControl',
            'playbackRateMenuButton',
            'fullscreenToggle'
          ]
        }
      })
      
      player.on('error', (e) => {
        console.error('视频播放错误:', e)
      })
    }
    
    watch(() => props.src, (newSrc) => {
      if (player && newSrc) {
        player.src({
          src: newSrc,
          type: 'video/mp4'
        })
      }
    })
    
    watch(() => props.poster, (newPoster) => {
      if (player && newPoster) {
        player.poster(newPoster)
      }
    })
    
    onMounted(() => {
      initPlayer()
    })
    
    onBeforeUnmount(() => {
      if (player) {
        player.dispose()
        player = null
      }
    })
    
    return {
      videoRef
    }
  }
}
</script>

<style scoped>
.video-player-wrapper {
  width: 100%;
  max-width: 900px;
  margin: 0 auto;
  border-radius: 8px;
  overflow: hidden;
}
</style>
