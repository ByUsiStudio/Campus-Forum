<template>
  <div class="article-list">
    <!-- 加载状态 -->
    <div v-if="loading" class="article-grid">
      <v-skeleton-loader
        v-for="i in 6"
        :key="i"
        type="article"
        class="article-skeleton"
      />
    </div>
    
    <!-- 空状态 -->
    <div v-else-if="articles.length === 0" class="text-center pa-12">
      <v-icon size="80" color="grey-lighten-2" class="mb-4">mdi-text-box-outline</v-icon>
      <div class="text-h6 text-medium-emphasis mb-2">暂无文章</div>
      <div class="text-body-2 text-medium-emphasis">快来发布第一篇文章吧</div>
    </div>
    
    <!-- 文章列表 -->
    <div v-else class="article-grid">
      <v-card
        v-for="article in articles"
        :key="article.id"
        class="article-card"
        variant="flat"
        rounded="lg"
      >
        <router-link
          :to="'/article/' + article.id"
          class="article-link"
        >
          <div class="article-header">
            <UserAvatar :user="article.user" :size="40" />
            <div class="article-meta flex-grow-1">
              <div class="d-flex align-center gap-2">
                <v-chip
                  v-if="article.is_pinned"
                  size="x-small"
                  color="orange"
                  variant="tonal"
                  prepend-icon="mdi-pin"
                >
                  置顶
                </v-chip>
                <span class="text-body-2 font-weight-medium">{{ article.user?.display_name || article.user?.username }}</span>
              </div>
              <div class="text-caption text-medium-emphasis">
                {{ formatDate(article.created_at) }}
              </div>
            </div>
            <v-chip
              v-if="article.category?.name"
              size="x-small"
              :color="article.category.color"
              variant="tonal"
            >
              {{ article.category.name }}
            </v-chip>
          </div>

          <div class="article-title">{{ article.title }}</div>

          <div class="article-excerpt text-body-2 text-medium-emphasis">
            {{ getExcerpt(article.content) }}
          </div>

          <div class="article-footer">
            <div class="article-stats">
              <span class="stat-item">
                <v-icon size="small" color="pink">mdi-heart</v-icon>
                {{ article.like_count || 0 }}
              </span>
              <span class="stat-item">
                <v-icon size="small">mdi-eye-outline</v-icon>
                {{ article.view_count || 0 }}
              </span>
              <span class="stat-item">
                <v-icon size="small">mdi-comment-outline</v-icon>
                {{ article.comment_count || 0 }}
              </span>
            </div>
          </div>
        </router-link>
      </v-card>
    </div>
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
    },
    loading: {
      type: Boolean,
      default: false
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
  width: 100%;
}

.article-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 16px;
}

.article-card {
  background: rgb(var(--v-theme-surface));
  border: 1px solid rgba(var(--v-theme-on-surface), 0.08);
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.25s ease;
}

.article-card:hover {
  border-color: rgba(var(--v-theme-primary), 0.3);
  box-shadow: 0 4px 20px rgba(var(--v-theme-primary), 0.1);
  transform: translateY(-4px);
}

.article-link {
  display: flex;
  flex-direction: column;
  padding: 16px;
  text-decoration: none;
  color: inherit;
  height: 100%;
  box-sizing: border-box;
}

.article-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.article-meta {
  flex: 1;
  min-width: 0;
}

.author-name {
  font-weight: 600;
  font-size: 0.875rem;
  color: rgb(var(--v-theme-on-surface));
  line-height: 1.3;
}

.article-title {
  font-size: 1rem;
  font-weight: 600;
  color: rgb(var(--v-theme-primary));
  line-height: 1.4;
  margin-bottom: 8px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  transition: color 0.2s;
}

.article-card:hover .article-title {
  color: rgb(var(--v-theme-secondary));
}

.article-excerpt {
  flex: 1;
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  margin-bottom: 12px;
}

.article-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 12px;
  border-top: 1px solid rgba(var(--v-theme-on-surface), 0.06);
}

.article-stats {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 0.75rem;
  color: rgb(var(--v-theme-on-surface));
  opacity: 0.7;
}

.article-skeleton {
  border-radius: 12px;
  overflow: hidden;
}

@media (max-width: 600px) {
  .article-grid {
    grid-template-columns: 1fr;
  }

  .article-card {
    border-radius: 8px;
  }
}
</style>