<template>
  <div class="article-list">
    <div v-if="articles.length === 0" class="empty">
      暂无文章
    </div>
    <div v-for="article in articles" :key="article.id" class="article-item">
      <div class="article-header">
        <h3 class="article-title">
          <router-link :to="'/article/' + article.id">{{ article.title }}</router-link>
        </h3>
        <div class="article-meta">
          <span>作者：{{ article.user.display_name }}</span>
          <span>分区：{{ article.category.name }}</span>
          <span>时间：{{ formatDate(article.created_at) }}</span>
        </div>
      </div>
      <div class="article-excerpt markdown-body" v-html="getExcerpt(article.content_html)"></div>
      <div class="article-stats">
        <span>{{ article.like_count }} 点赞</span>
        <span>{{ article.view_count }} 阅读</span>
        <span>{{ article.comment_count || 0 }} 评论</span>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ArticleList',
  props: {
    articles: {
      type: Array,
      default: () => []
    }
  },
  setup() {
    const formatDate = (date) => {
      return new Date(date).toLocaleString('zh-CN')
    }
    
    const getExcerpt = (html) => {
      // 简单提取纯文本作为摘要
      const div = document.createElement('div')
      div.innerHTML = html
      const text = div.textContent || div.innerText || ''
      return text.substring(0, 200) + (text.length > 200 ? '...' : '')
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
  gap: 20px;
}

.article-item {
  padding: 20px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  transition: background 0.2s;
}

.article-header {
  margin-bottom: 15px;
}

.article-title {
  margin-bottom: 10px;
  font-size: 18px;
}

.article-title a {
  text-decoration: none;
  color: #1e293b;
}

.article-meta {
  font-size: 13px;
  color: #6b7280;
  display: flex;
  gap: 15px;
  flex-wrap: wrap;
}

.article-excerpt {
  margin: 10px 0;
  color: #4b5563;
  line-height: 1.5;
  font-size: 14px;
}

.article-stats {
  display: flex;
  gap: 20px;
  font-size: 13px;
  color: #6b7280;
  margin-top: 10px;
}

.empty {
  text-align: center;
  padding: 40px;
  color: #6b7280;
}
</style>