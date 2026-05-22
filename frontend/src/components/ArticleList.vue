<template>
  <div class="article-list">
    <div v-if="articles.length === 0" class="text-center pa-8 text-medium-emphasis">
      暂无文章
    </div>
    <v-card
      v-for="article in articles"
      :key="article.id"
      class="mb-3"
      variant="outlined"
    >
      <v-card-text>
        <div class="d-flex gap-3 mb-3">
          <v-avatar size="48" color="primary">
            <v-img v-if="article.user.avatar" :src="article.user.avatar"></v-img>
            <span v-else class="text-h6">{{ article.user.display_name?.[0] || 'U' }}</span>
          </v-avatar>
          <div class="flex-grow-1">
            <router-link
              :to="'/article/' + article.id"
              class="text-decoration-none"
            >
              <div class="text-h6 text-primary font-weight-bold hover-underline">
                {{ article.title }}
              </div>
            </router-link>
            <div class="d-flex flex-wrap gap-2 mt-1 text-caption text-medium-emphasis">
              <span>
                <v-icon size="small">mdi-account</v-icon>
                {{ article.user.display_name }}
              </span>
              <span>
                <v-icon size="small">mdi-folder</v-icon>
                {{ article.category?.name || '未分类' }}
              </span>
              <span>
                <v-icon size="small">mdi-clock</v-icon>
                {{ formatDate(article.created_at) }}
              </span>
            </div>
          </div>
        </div>
        
        <div class="text-body-2 text-medium-emphasis mb-3 article-excerpt">
          {{ getExcerpt(article.content_html) }}
        </div>
        
        <div class="d-flex gap-4 text-caption text-medium-emphasis">
          <span>
            <v-icon size="small" color="pink">mdi-heart</v-icon>
            {{ article.like_count }}
          </span>
          <span>
            <v-icon size="small">mdi-eye</v-icon>
            {{ article.view_count }}
          </span>
          <span>
            <v-icon size="small">mdi-comment</v-icon>
            {{ article.comment_count || 0 }}
          </span>
        </div>
      </v-card-text>
    </v-card>
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
      if (!html) return ''
      const div = document.createElement('div')
      div.innerHTML = html
      const text = div.textContent || div.innerText || ''
      return text.substring(0, 150) + (text.length > 150 ? '...' : '')
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

.hover-underline:hover {
  text-decoration: underline;
}

.article-excerpt {
  line-height: 1.5;
}
</style>
