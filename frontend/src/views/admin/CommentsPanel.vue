<template>
  <div class="comments-panel">
    <div class="panel-header">
      <div class="header-left">
        <h2 class="panel-title">评论管理</h2>
        <p class="panel-subtitle">审核与删除用户评论</p>
      </div>
      <v-btn variant="outlined" color="primary" @click="$emit('refresh')" :loading="loading">
        <v-icon start>mdi-refresh</v-icon>
        刷新
      </v-btn>
    </div>

    <div v-if="comments.length === 0" class="empty-state">
      <v-icon size="64" color="grey-lighten-1">mdi-comment-text-outline</v-icon>
      <div class="empty-text">暂无评论</div>
    </div>

    <div v-else class="comments-list">
      <v-card
        v-for="comment in comments"
        :key="comment.id"
        class="comment-card"
        variant="outlined"
      >
        <div class="comment-main">
          <UserAvatar :user="comment.User || {}" :size="44" />

          <div class="comment-content">
            <div class="comment-header">
              <span class="comment-author">{{ comment.User?.display_name || '未知用户' }}</span>
              <span class="comment-time">{{ formatDate(comment.created_at) }}</span>
            </div>

            <div class="comment-reply-info" v-if="comment.Article">
              <v-icon size="14">mdi-subdirectory-arrow-right</v-icon>
              回复文章：{{ comment.Article.title }}
            </div>

            <div class="comment-text">{{ comment.content }}</div>
          </div>

          <div class="comment-actions">
            <v-btn
              variant="text"
              size="small"
              color="error"
              @click="$emit('delete', comment.id)"
            >
              <v-icon start>mdi-delete</v-icon>
              删除
            </v-btn>
          </div>
        </div>
      </v-card>
    </div>
  </div>
</template>

<script>
import UserAvatar from '../../components/UserAvatar.vue'

export default {
  name: 'CommentsPanel',
  components: {
    UserAvatar
  },
  props: {
    comments: {
      type: Array,
      default: () => []
    },
    loading: {
      type: Boolean,
      default: false
    }
  },
  emits: ['delete', 'refresh'],
  setup() {
    const formatDate = (dateString) => {
      if (!dateString) return '-'
      const date = new Date(dateString)
      return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
      })
    }

    return {
      formatDate
    }
  }
}
</script>

<style scoped>
.comments-panel {
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 16px;
}

.panel-title {
  font-size: 1.5rem;
  font-weight: 700;
  color: #1a1a2e;
  margin: 0 0 4px 0;
}

.panel-subtitle {
  font-size: 0.9rem;
  color: #6b7280;
  margin: 0;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  background: white;
  border-radius: 16px;
}

.empty-text {
  margin-top: 16px;
  font-size: 1rem;
  color: #9ca3af;
}

.comments-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.comment-card {
  border-radius: 12px;
  padding: 20px;
  transition: box-shadow 0.2s;
}

.comment-card:hover {
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.comment-main {
  display: flex;
  gap: 16px;
  align-items: flex-start;
}

.comment-content {
  flex: 1;
  min-width: 0;
}

.comment-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
  flex-wrap: wrap;
}

.comment-author {
  font-weight: 600;
  color: #1a1a2e;
}

.comment-time {
  font-size: 0.8rem;
  color: #9ca3af;
}

.comment-reply-info {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 0.85rem;
  color: #6b7280;
  margin-bottom: 8px;
  padding: 6px 10px;
  background: #f8f9ff;
  border-radius: 6px;
  width: fit-content;
}

.comment-text {
  color: #4b5563;
  line-height: 1.6;
  word-break: break-word;
}

.comment-actions {
  flex-shrink: 0;
}

@media (max-width: 600px) {
  .comment-main {
    flex-direction: column;
  }

  .comment-actions {
    width: 100%;
    display: flex;
    justify-content: flex-end;
  }
}
</style>
