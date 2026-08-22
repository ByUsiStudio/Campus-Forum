<template>
  <div class="article-list">
    <!-- 加载状态 -->
    <div v-if="loading" class="article-grid">
      <el-card
        v-for="i in 6"
        :key="i"
        shadow="always"
        class="article-card article-skeleton"
      >
        <el-skeleton :rows="5" animated />
      </el-card>
    </div>

    <!-- 空状态 -->
    <div v-else-if="articles.length === 0" class="empty-state">
      <el-icon :size="80" class="text-secondary empty-icon">
        <Document />
      </el-icon>
      <div class="empty-title">暂无文章</div>
      <div class="text-secondary empty-subtitle">快来发布第一篇文章吧</div>
    </div>

    <!-- 文章列表 -->
    <div v-else class="article-grid">
      <el-card
        v-for="article in articles"
        :key="article.id"
        shadow="hover"
        class="article-card"
      >
        <router-link
          :to="'/article/' + article.id"
          class="article-link"
        >
          <div class="article-header">
            <UserAvatar :user="article.user" :size="40" />
            <div class="article-meta">
              <div class="meta-top">
                <el-tag
                  v-if="article.is_pinned"
                  size="small"
                  color="rgba(255,165,0,0.12)"
                  class="pinned-tag"
                  effect="plain"
                >
                  <el-icon class="tag-icon"><Promotion /></el-icon>
                  置顶
                </el-tag>
                <span class="author-name">{{ article.user?.display_name || article.user?.username }}</span>
              </div>
              <div class="text-secondary article-date">
                {{ formatDate(article.created_at) }}
              </div>
            </div>
            <el-tag
              v-if="article.category?.name"
              size="small"
              :color="article.category.color"
              effect="light"
              class="category-tag"
            >
              {{ article.category.name }}
            </el-tag>
          </div>

          <div class="article-title">{{ article.title }}</div>

          <div class="article-excerpt text-secondary">
            {{ getExcerpt(article.content) }}
          </div>

          <div class="article-footer">
            <div class="article-stats">
              <span class="stat-item">
                <el-icon class="stat-icon stat-like"><SemiSelect /></el-icon>
                {{ article.like_count || 0 }}
              </span>
              <span class="stat-item">
                <el-icon class="stat-icon"><View /></el-icon>
                {{ article.view_count || 0 }}
              </span>
              <span class="stat-item">
                <el-icon class="stat-icon"><ChatDotRound /></el-icon>
                {{ article.comment_count || 0 }}
              </span>
            </div>
          </div>
        </router-link>
      </el-card>
    </div>
  </div>
</template>

<script>
import { Document, Promotion, SemiSelect, View, ChatDotRound } from '@element-plus/icons-vue'
import UserAvatar from './UserAvatar.vue'

export default {
  name: 'ArticleList',
  components: {
    UserAvatar,
    Document,
    Promotion,
    SemiSelect,
    View,
    ChatDotRound
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
  background: var(--campus-surface);
  border: 1px solid var(--campus-border);
  border-radius: var(--campus-radius);
  overflow: hidden;
  transition: all 0.25s ease;
}

.article-card:hover {
  border-color: var(--campus-primary);
  box-shadow: 0 4px 20px rgba(103, 80, 164, 0.12);
  transform: translateY(-4px);
}

.article-card :deep(.el-card__body) {
  padding: 0;
}

.article-link {
  display: flex;
  flex-direction: column;
  padding: 10px 12px;
  text-decoration: none;
  color: inherit;
  height: 100%;
  box-sizing: border-box;
}

.article-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 6px;
}

.article-meta {
  flex: 1;
  min-width: 0;
}

.meta-top {
  display: flex;
  align-items: center;
  gap: 6px;
}

.pinned-tag {
  border: none;
  color: #e6a23c;
}

.pinned-tag .tag-icon {
  margin-right: 2px;
  vertical-align: -2px;
}

.category-tag {
  flex-shrink: 0;
}

.author-name {
  font-weight: 600;
  font-size: 0.875rem;
  color: var(--campus-text);
  line-height: 1.3;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.article-date {
  font-size: 0.75rem;
  line-height: 1.3;
  margin-top: 1px;
}

.article-title {
  font-size: 1rem;
  font-weight: 600;
  color: var(--campus-primary);
  line-height: 1.35;
  margin-bottom: 4px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  transition: color 0.2s;
}

.article-card:hover .article-title {
  color: var(--campus-primary-dark);
}

.article-excerpt {
  flex: 1;
  font-size: 0.875rem;
  line-height: 1.45;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  margin-bottom: 8px;
}

.article-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 8px;
  border-top: 1px solid var(--campus-border);
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
  color: var(--campus-text);
  opacity: 0.7;
}

.stat-icon {
  font-size: 14px;
}

.stat-like {
  color: #ec6b8f;
}

.article-skeleton {
  cursor: default;
}

.article-skeleton:hover {
  border-color: var(--campus-border);
  box-shadow: var(--campus-shadow);
  transform: none;
}

.article-skeleton :deep(.el-card__body) {
  padding: 16px;
}

.empty-state {
  text-align: center;
  padding: 48px 16px;
}

.empty-icon {
  margin-bottom: 16px;
  display: block;
}

.empty-title {
  font-size: 1.125rem;
  font-weight: 600;
  margin-bottom: 8px;
  color: var(--campus-text);
}

.empty-subtitle {
  font-size: 0.875rem;
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
