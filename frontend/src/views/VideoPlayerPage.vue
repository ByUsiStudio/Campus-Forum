<template>
  <div class="page-container">
    <div class="video-card card-surface">
      <div class="video-player-wrap">
        <VideoPlayer :src="videoUrl" :poster="posterUrl" />
      </div>

      <div class="video-info">
        <h1 class="video-title">视频播放</h1>
        <div class="video-meta">
          <span class="meta-item">
            <el-icon :size="14"><VideoPlay /></el-icon>
            高清在线播放
          </span>
          <span class="meta-item">
            <el-icon :size="14"><InfoFilled /></el-icon>
            建议全屏观看以获得最佳体验
          </span>
        </div>

        <div class="video-actions">
          <el-button type="primary" @click="goBack">
            <el-icon class="mr-8"><ArrowLeft /></el-icon>
            返回文章
          </el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, VideoPlay, InfoFilled } from '@element-plus/icons-vue'
import VideoPlayer from '../components/VideoPlayer.vue'

export default {
  name: 'VideoPlayerPage',
  components: {
    VideoPlayer,
    ArrowLeft,
    VideoPlay,
    InfoFilled
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
.video-card {
  max-width: 960px;
  margin: 0 auto;
  border-radius: var(--campus-radius);
  overflow: hidden;
}

.video-player-wrap {
  background: #000;
}

.video-info {
  padding: 20px 24px 24px;
}

.video-title {
  font-size: 20px;
  font-weight: 700;
  color: var(--campus-text);
  margin: 0 0 12px;
  letter-spacing: -0.01em;
}

.video-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  margin-bottom: 16px;
}

.meta-item {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: var(--campus-text-secondary);
}

.video-actions {
  display: flex;
  padding-top: 6px;
  border-top: 1px solid var(--campus-border);
}

.mr-8 {
  margin-right: 8px;
}

@media (max-width: 768px) {
  .video-info {
    padding: 16px;
  }
}
</style>
