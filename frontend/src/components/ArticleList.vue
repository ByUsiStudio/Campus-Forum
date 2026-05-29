<template>
  <div class="article-list">
    <div v-if="articles.length === 0" class="text-center pa-8 text-medium-emphasis">
      <v-icon size="48" class="mb-2">mdi-text-box-outline</v-icon>
      <div>暂无文章</div>
    </div>
    <v-card
      v-for="article in articles"
      :key="article.id"
      class="article-card mb-4"
      variant="elevated"
      hover
    >
      <v-card-text class="pa-4">
        <div class="d-flex gap-4">
          <UserAvatar :user="article.user" :size="48" class="flex-shrink-0" />
          <div class="flex-grow-1 min-width-0">
            <div class="d-flex align-start justify-space-between gap-2 mb-2">
              <router-link
                :to="'/article/' + article.id"
                class="article-title text-decoration-none"
              >
                {{ article.title }}
              </router-link>
            </div>

            <div class="d-flex flex-wrap align-center gap-3 mb-3 text-caption text-medium-emphasis">
              <span class="d-flex align-center gap-1">
                <v-icon size="small">mdi-folder-outline</v-icon>
                {{ article.category?.name || '未分类' }}
              </span>
              <span class="d-flex align-center gap-1">
                <v-icon size="small">mdi-clock-outline</v-icon>
                {{ formatDate(article.created_at) }}
              </span>
            </div>

            <div class="article-excerpt text-body-2 text-medium-emphasis mb-3">
              {{ getExcerpt(article.content) }}
            </div>

            <div class="d-flex align-center gap-4 text-caption">
              <span class="d-flex align-center gap-1">
                <v-icon size="small" color="pink">mdi-heart</v-icon>
                {{ article.like_count || 0 }}
              </span>
              <span class="d-flex align-center gap-1">
                <v-icon size="small">mdi-eye-outline</v-icon>
                {{ article.view_count || 0 }}
              </span>
              <span class="d-flex align-center gap-1">
                <v-icon size="small">mdi-comment-outline</v-icon>
                {{ article.comment_count || 0 }}
              </span>
              <span v-if="article.favorite_count" class="d-flex align-center gap-1">
                <v-icon size="small" color="amber">mdi-star</v-icon>
                {{ article.favorite_count || 0 }}
              </span>
            </div>
          </div>
        </div>
      </v-card-text>
    </v-card>
  </div>
</template>

<script>
import UserAvatar from './UserAvatar.vue'

export default {
  name: 'ArticleList',
  components: {
    UserAvatar
  },
  props: {
    articles: {
      type: Array,
      default: () => []
    }
  },
  setup() {
    const formatDate = (date) => {
      if (!date) return ''
      const d = new Date(date)
      const now = new Date()
      const diff = now - d
      const days = Math.floor(diff / (1000 * 60 * 60 * 24))

      if (days === 0) {
        const hours = Math.floor(diff / (1000 * 60 * 60))
        if (hours === 0) {
          const minutes = Math.floor(diff / (1000 * 60))
          return minutes <= 1 ? '刚刚' : `${minutes} 分钟前`
        }
        return `${hours} 小时前`
      } else if (days === 1) {
        return '昨天'
      } else if (days < 7) {
        return `${days} 天前`
      } else {
        return d.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' })
      }
    }

    const getExcerpt = (content) => {
      if (!content) return ''
      const text = content
        .replace(/#{1,6}\s+/g, '')
        .replace(/\*\*(.+?)\*\*/g, '$1')
        .replace(/\*(.+?)\*/g, '$1')
        .replace(/`(.+?)`/g, '$1')
        .replace(/\[(.+?)\]\(.+?\)/g, '$1')
        .replace(/!\[.*?\]\(.+?\)/g, '[图片]')
        .replace(/\n+/g, ' ')
        .trim()
      return text.length > 120 ? text.substring(0, 120) + '...' : text
    }

    return {
      formatDate,
      getExcerpt
    }
  }
}
</script>

<style scoped>
.article-list {
  display: flex;
  flex-direction: column;
}

.article-card {
  transition: transform 0.2s, box-shadow 0.2s;
}

.article-card:hover {
  transform: translateY(-2px);
}

.article-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: rgb(var(--v-theme-primary));
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  transition: color 0.2s;
}

.article-title:hover {
  color: rgb(var(--v-theme-secondary));
  text-decoration: underline;
}

.article-excerpt {
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.min-width-0 {
  min-width: 0;
}

.flex-shrink-0 {
  flex-shrink: 0;
}
</style>