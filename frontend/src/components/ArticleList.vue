<template>
  <div class="article-list">
    <!-- 加载状态 -->
    <div v-if="loading" class="article-grid">
      <div
        v-for="i in 6"
        :key="i"
        class="article-card article-skeleton"
      >
        <div class="skeleton-cover"></div>
        <div class="skeleton-body">
          <div class="skeleton-line w40"></div>
          <div class="skeleton-line w90"></div>
          <div class="skeleton-line w70"></div>
          <div class="skeleton-line w80"></div>
          <div class="skeleton-foot">
            <div class="skeleton-dot"></div>
            <div class="skeleton-dot"></div>
            <div class="skeleton-dot"></div>
          </div>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else-if="articles.length === 0" class="empty-state">
      <div class="empty-icon">
        <el-icon :size="56"><Document /></el-icon>
      </div>
      <div class="empty-title">暂无文章</div>
      <div class="text-secondary empty-subtitle">快来发布第一篇文章吧</div>
    </div>

    <!-- 文章列表 -->
    <div v-else class="article-grid">
      <router-link
        v-for="article in articles"
        :key="article.id"
        :to="'/article/' + article.id"
        class="article-card"
      >
        <div class="card-topbar">
          <div class="author-line">
            <UserAvatar
              :user="article.user"
              :size="34"
              :show-username="false"
              class="author-avatar"
            />
            <span class="author-name">
              {{ article.user?.display_name || article.user?.username }}
            </span>
          </div>
          <span v-if="article.category?.name" class="category-tag" :style="categoryStyle(article.category)">
            {{ article.category.name }}
          </span>
        </div>

        <div v-if="coverImage(article)" class="cover-wrap">
          <img :src="coverImage(article)" :alt="article.title" class="cover-img" loading="lazy" />
        </div>

        <div class="card-body">
          <div v-if="article.is_pinned" class="pinned-tag">
            <el-icon class="tag-icon"><Promotion /></el-icon>
            置顶
          </div>
          <h3 class="article-title" v-html="highlightText(article.title)"></h3>
          <p class="article-excerpt article-excerpt-multiline" v-html="DOMPurify.sanitize(highlightText(getExcerpt(article.content)))"></p>
        </div>

        <div class="card-footer">
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
          <span class="article-date">{{ formatDate(article.created_at) }}</span>
        </div>
      </router-link>
    </div>
  </div>
</template>

<script>
import { Document, Promotion, SemiSelect, View, ChatDotRound } from '@element-plus/icons-vue'
import UserAvatar from './UserAvatar.vue'
import DOMPurify from 'dompurify'

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
    },
    // 可选搜索关键词；提供时在标题/摘要中高亮匹配文本
    highlightQuery: {
      type: String,
      default: ''
    }
  },
  setup(props) {
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

    // 封面：优先取文章自带 cover_image，其次提取 markdown 第一张图片（纯展示，不引入接口）
    const coverImage = (article) => {
      if (article?.cover_image) return article.cover_image
      const content = article?.content || ''
      if (!content) return ''
      const m = content.match(/!\[.*?\]\((.+?)\)/)
      return m ? m[1] : ''
    }

    // 分类标签内联样式（保持原 el-tag 配色的视觉习惯）
    const categoryStyle = (category) => {
      const color = category.color || '#4f6ef7'
      return {
        backgroundColor: applyAlpha(color, 0.12),
        color
      }
    }

    const applyAlpha = (hex, alpha) => {
      const clean = String(hex || '').replace('#', '')
      if (!/^[0-9a-fA-F]{6}$/.test(clean)) return hex
      const r = parseInt(clean.slice(0, 2), 16)
      const g = parseInt(clean.slice(2, 4), 16)
      const b = parseInt(clean.slice(4, 6), 16)
      return `rgba(${r}, ${g}, ${b}, ${alpha})`
    }

    // 高亮关键词（转义正则，避免特殊字符破坏模式）
    const highlightText = (text) => {
      const raw = String(props.highlightQuery || '').trim()
      if (!raw || !text) return text
      const escaped = raw.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
      const re = new RegExp(`(${escaped})`, 'gi')
      return String(text).replace(re, '<mark class="hl">$1</mark>')
    }

    return {
      formatDate,
      getExcerpt,
      coverImage,
      categoryStyle,
      highlightText
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
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 18px;
}

.article-card {
  display: flex;
  flex-direction: column;
  background: var(--campus-surface);
  border: 1px solid var(--campus-border);
  border-radius: var(--campus-radius);
  overflow: hidden;
  text-decoration: none;
  color: inherit;
  transition: transform 0.22s ease, box-shadow 0.22s ease, border-color 0.22s ease;
}

.article-card:hover {
  border-color: rgba(79, 110, 247, 0.4);
  box-shadow: 0 10px 30px rgba(15, 23, 42, 0.08);
  transform: translateY(-4px);
}

/* 顶部：作者 + 分类 */
.card-topbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  padding: 14px 14px 0;
}

.author-line {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.author-avatar {
  flex-shrink: 0;
}

.author-name {
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--campus-text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.category-tag {
  flex-shrink: 0;
  font-size: 0.72rem;
  font-weight: 600;
  padding: 3px 10px;
  border-radius: 999px;
}

/* 封面 */
.cover-wrap {
  margin: 12px 14px 0;
  border-radius: 12px;
  overflow: hidden;
  aspect-ratio: 16 / 9;
  background: var(--campus-bg);
}

.cover-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
  transition: transform 0.35s ease;
}

.article-card:hover .cover-img {
  transform: scale(1.04);
}

/* 正文 */
.card-body {
  padding: 12px 14px 8px;
  flex: 1;
  min-width: 0;
}

.pinned-tag {
  display: inline-flex;
  align-items: center;
  gap: 2px;
  font-size: 0.72rem;
  font-weight: 600;
  color: #e6a23c;
  background: rgba(255, 165, 0, 0.12);
  padding: 2px 8px;
  border-radius: 6px;
  margin-bottom: 6px;
}

.pinned-tag .tag-icon {
  vertical-align: -2px;
}

.article-title {
  font-size: 1.02rem;
  font-weight: 650;
  color: var(--campus-text);
  line-height: 1.4;
  margin: 0 0 6px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  transition: color 0.2s;
}

.article-card:hover .article-title {
  color: var(--campus-primary);
}

.article-excerpt {
  font-size: 0.84rem;
  line-height: 1.55;
  color: var(--campus-text-secondary);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  margin: 0;
}

/* 搜索关键词高亮 */
.hl {
  background: rgba(79, 110, 247, 0.18);
  color: var(--campus-primary);
  border-radius: 4px;
  padding: 0 3px;
  font-weight: 600;
}

/* 底部 */
.card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 14px 14px;
  border-top: 1px solid var(--campus-border);
  margin-top: 4px;
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
  font-size: 0.78rem;
  color: var(--campus-text-secondary);
}

.stat-icon {
  font-size: 14px;
}

.stat-like {
  color: #ec6b8f;
}

.article-date {
  font-size: 0.75rem;
  color: var(--campus-text-secondary);
  white-space: nowrap;
}

/* 骨架屏 */
.article-skeleton {
  cursor: default;
}

.article-skeleton:hover {
  transform: none;
  box-shadow: none;
  border-color: var(--campus-border);
}

.skeleton-cover {
  margin: 12px 14px 0;
  height: 150px;
  border-radius: 12px;
  background: linear-gradient(100deg, #eef1f7 40%, #f6f8fc 50%, #eef1f7 60%);
  background-size: 200% 100%;
  animation: shimmer 1.4s infinite;
}

.skeleton-body {
  padding: 12px 14px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.skeleton-line {
  height: 14px;
  border-radius: 6px;
  background: linear-gradient(100deg, #eef1f7 40%, #f6f8fc 50%, #eef1f7 60%);
  background-size: 200% 100%;
  animation: shimmer 1.4s infinite;
}

.skeleton-line.w40 { width: 40%; }
.skeleton-line.w70 { width: 70%; }
.skeleton-line.w80 { width: 80%; }
.skeleton-line.w90 { width: 90%; }

.skeleton-foot {
  display: flex;
  gap: 16px;
  margin-top: 6px;
}

.skeleton-dot {
  width: 54px;
  height: 14px;
  border-radius: 6px;
  background: linear-gradient(100deg, #eef1f7 40%, #f6f8fc 50%, #eef1f7 60%);
  background-size: 200% 100%;
  animation: shimmer 1.4s infinite;
}

@keyframes shimmer {
  to { background-position: -200% 0; }
}

/* 空状态 */
.empty-state {
  text-align: center;
  padding: 56px 16px;
}

.empty-icon {
  margin-bottom: 14px;
  color: var(--campus-text-secondary);
  opacity: 0.5;
}

.empty-title {
  font-size: 1.1rem;
  font-weight: 600;
  margin-bottom: 6px;
  color: var(--campus-text);
}

.empty-subtitle {
  font-size: 0.875rem;
}

@media (max-width: 640px) {
  .article-grid {
    grid-template-columns: 1fr;
    gap: 14px;
  }

  .article-card {
    border-radius: 12px;
  }
}
</style>
