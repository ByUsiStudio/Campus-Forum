<template>
  <div class="video-page-container page-container">
    <el-card class="video-card card-surface" shadow="always">
      <template #header>
        <div class="card-header">
          <span class="video-title">视频播放</span>
        </div>
      </template>

      <div class="video-body">
        <VideoPlayer :src="videoUrl" :poster="posterUrl" />
      </div>

      <div class="video-actions">
        <el-button type="primary" @click="goBack">
          <el-icon class="mr-8"><ArrowLeft /></el-icon>
          返回文章
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft } from '@element-plus/icons-vue'
import VideoPlayer from '../components/VideoPlayer.vue'

export default {
  name: 'VideoPlayerPage',
  components: {
    VideoPlayer,
    ArrowLeft
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

.card-header {
  display: flex;
  align-items: center;
  justify-content: center;
}

.video-title {
  font-size: 20px;
  font-weight: 600;
}

.video-body {
  display: flex;
  justify-content: center;
  padding: 0 0 12px;
}

.video-actions {
  display: flex;
  justify-content: center;
  padding-top: 8px;
}

.mr-8 {
  margin-right: 8px;
}
</style>
