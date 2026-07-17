<template>
  <div class="video-page-container">
    <div class="layui-card video-card">
      <div class="layui-card-header text-center">
        <h3 style="font-size: 18px;">视频播放</h3>
      </div>
      
      <div class="layui-card-body text-center">
        <VideoPlayer :src="videoUrl" :poster="posterUrl" />
      </div>
      
      <div class="layui-card-footer text-center">
        <button class="layui-btn" @click="goBack">
          <i class="fa-solid fa-arrow-left mr-1"></i>
          返回文章
        </button>
      </div>
    </div>
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

.mr-1 {
  margin-right: 5px;
}
</style>