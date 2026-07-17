<template>
  <div class="article-list">
    <div v-if="loading" class="article-grid">
      <div v-for="i in 6" :key="i" class="article-skeleton"></div>
    </div>
    
    <div v-else-if="articles.length === 0" class="empty-state">
      <i class="fa-regular fa-file-lines"></i>
      <div class="empty-title">暂无文章</div>
      <div class="empty-hint">快来发布第一篇文章吧</div>
    </div>
    
    <div v-else class="article-grid">
      <div v-for="article in articles" :key="article.id" class="article-card">
        <router-link :to="'/article/' + article.id" class="article-link">
          <div class="article-header">
            <UserAvatar :user="article.user" :size="40" />
            <div class="article-meta">
              <div class="article-author">
                <span v-if="article.is_pinned" class="pin-tag">
                  <i class="fa-solid fa-thumbtack"></i>
                  置顶
                </span>
                <span>{{ article.user?.display_name || article.user?.username }}</span>
              </div>
              <div class="article-time">{{ formatDate(article.created_at) }}</div>
            </div>
            <span v-if="article.category?.name" class="category-tag" :style="{ background: article.category.color + '20', color: article.category.color }">
              {{ article.category.name }}
            </span>
          </div>

          <div class="article-title">{{ article.title }}</div>

          <div class="article-excerpt">{{ getExcerpt(article.content) }}</div>

          <div class="article-footer">
            <div class="article-stats">
              <span class="stat-item">
                <i class="fa-solid fa-heart text-error"></i>
                {{ article.like_count || 0 }}
              </span>
              <span class="stat-item">
                <i class="fa-regular fa-eye"></i>
                {{ article.view_count || 0 }}
              </span>
              <span class="stat-item">
                <i class="fa-regular fa-comment"></i>
                {{ article.comment_count || 0 }}
              </span>
            </div>
          </div>
        </router-link>
      </div>
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
  background: #fff;
  border: 1px solid #f0f0f0;
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.25s ease;

  &:hover {
    border-color: rgba(30, 159, 255, 0.3);
    box-shadow: 0 4px 20px rgba(30, 159, 255, 0.1);
    transform: translateY(-4px);
  }
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

.article-author {
  font-size: 14px;
  font-weight: 500;
  color: #333;
  display: flex;
  align-items: center;
  gap: 8px;
}

.pin-tag {
  font-size: 12px;
  color: #FAAD14;
  background: rgba(250, 173, 20, 0.1);
  padding: 2px 8px;
  border-radius: 4px;
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.article-time {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.category-tag {
  font-size: 12px;
  padding: 4px 10px;
  border-radius: 4px;
}

.article-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  line-height: 1.4;
  margin-bottom: 8px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  transition: color 0.2s;
}

.article-card:hover .article-title {
  color: #1E9FFF;
}

.article-excerpt {
  flex: 1;
  font-size: 14px;
  color: #666;
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
  border-top: 1px solid #f0f0f0;
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
  font-size: 13px;
  color: #666;
}

.text-error {
  color: #FF5722;
}

.article-skeleton {
  background: #f5f5f5;
  border-radius: 12px;
  overflow: hidden;
  height: 200px;
}

.empty-state {
  text-align: center;
  padding: 48px 24px;

  i {
    font-size: 80px;
    color: #e8e8e8;
    margin-bottom: 16px;
  }

  .empty-title {
    font-size: 18px;
    color: #666;
    margin-bottom: 8px;
  }

  .empty-hint {
    font-size: 14px;
    color: #999;
  }
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