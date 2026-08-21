<template>
  <div class="video-page-container">
    <v-card class="video-card" elevation="4">
      <v-card-title class="text-center text-h5">视频播放</v-card-title>
      
      <v-card-text class="text-center">
        <VideoPlayer :src="videoUrl" :poster="posterUrl" />
      </v-card-text>
      
      <v-card-actions class="justify-center">
        <v-btn color="primary" @click="goBack">
          <v-icon start>mdi-arrow-left</v-icon>
          返回文章
        </v-btn>
      </v-card-actions>
    </v-card>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import VideoPlayer from '../components/VideoPlayer.vue'

export default {
  name: 'VideoPlayerPage',
  components: {
    VideoPlayer
  },
  setup() {
    const route = useRoute()
    const router = useRouter()
    const videoUrl = ref('')
    const posterUrl = ref('')
    const articleId = ref('')
    
    const goBack = () => {
      if (articleId.value) {
        router.push(`/article/${articleId.value}`)
      } else {
        router.push('/')
      }
    }
    
    onMounted(() => {
      videoUrl.value = route.query.src || ''
      posterUrl.value = route.query.poster || ''
      articleId.value = route.query.articleId || ''
    })
    
    return {
      videoUrl,
      posterUrl,
      articleId,
      goBack
    }
  }
}
</script>

<style scoped>
.video-page-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.video-card {
  max-width: 950px;
  width: 100%;
  border-radius: 16px;
}
</style>
