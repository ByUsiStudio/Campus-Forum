<template>
  <div class="reply-list ml-8">
    <div
      v-for="reply in replies"
      :key="reply.id"
      class="reply-item"
    >
      <div class="reply-header">
        <div
          class="reply-avatar"
          @click="$emit('goToUserProfile', reply.user.id)"
        >
          <img
            :src="reply.user.avatar || `/api/users/${reply.user.id}/avatar`"
            :alt="reply.user.display_name"
            class="avatar-img"
          >
        </div>
        <div class="reply-info">
          <div class="reply-author">
            {{ reply.user.display_name || reply.user.username || '匿名用户' }}
          </div>
          <div class="reply-time">{{ formatDate(reply.created_at) }}</div>
        </div>
      </div>
      
      <div class="reply-content">{{ reply.content }}</div>
      
      <div class="reply-actions">
        <button
          class="reply-action-btn"
          :class="{ liked: commentLiked[reply.id] }"
          @click="$emit('toggleLike', reply)"
        >
          <i class="fa-solid fa-thumbs-up"></i>
          {{ reply.like_count }}
        </button>
        <button
          v-if="token"
          class="reply-action-btn"
          @click="$emit('showReplyForm', reply.id)"
        >
          <i class="fa-solid fa-reply"></i>
          回复
        </button>
        <button
          v-if="canDeleteComment(reply)"
          class="reply-action-btn reply-action-delete"
          @click="$emit('deleteComment', reply.id, reply)"
        >
          <i class="fa-solid fa-trash"></i>
        </button>
      </div>

      <div v-if="replyingTo === reply.id" class="reply-form">
        <textarea
          v-model="localReplyContent"
          :placeholder="'回复 ' + (reply.user.display_name || reply.user.username || '匿名用户')"
          class="reply-textarea"
          rows="2"
        ></textarea>
        <div class="reply-form-footer">
          <label class="checkbox-label">
            <input type="checkbox" v-model="localReplyIsAnonymous">
            <span>匿名</span>
          </label>
          <div class="reply-form-actions">
            <button class="layui-btn layui-btn-sm" @click="submitReply(reply.id)">发送</button>
            <button class="layui-btn layui-btn-sm layui-btn-primary" @click="cancelReply">取消</button>
          </div>
        </div>
      </div>

      <CommentReply
        v-if="reply.replies && reply.replies.length > 0"
        :replies="reply.replies"
        :commentLiked="commentLiked"
        :token="token"
        :currentUser="currentUser"
        :replyingTo="replyingTo"
        :localReplyContent="localReplyContent"
        :localReplyIsAnonymous="localReplyIsAnonymous"
        @toggleLike="$emit('toggleLike', $event)"
        @showReplyForm="$emit('showReplyForm', $event)"
        @deleteComment="$emit('deleteComment', $event.id, $event.comment)"
        @goToUserProfile="$emit('goToUserProfile', $event)"
        @submitReply="$emit('submitReply', $event)"
        @cancelReply="$emit('cancelReply')"
      />
    </div>
  </div>
</template>

<script>
import { ref, watch } from 'vue'
import UserAvatar from './UserAvatar.vue'

export default {
  name: 'CommentReply',
  components: {
    UserAvatar
  },
  props: {
    replies: {
      type: Array,
      required: true
    },
    commentLiked: {
      type: Object,
      default: () => ({})
    },
    token: {
      type: String,
      default: ''
    },
    currentUser: {
      type: Object,
      default: null
    },
    replyingTo: {
      type: [Number, null],
      default: null
    },
    localReplyContent: {
      type: String,
      default: ''
    },
    localReplyIsAnonymous: {
      type: Boolean,
      default: false
    }
  },
  setup(props, { emit }) {
    const formatDate = (date) => {
      return new Date(date).toLocaleString('zh-CN')
    }

    const canDeleteComment = (comment) => {
      if (!props.currentUser) return false
      const currentUserId = Number(props.currentUser.id)
      const commentUserId = Number(comment.user_id)
      return currentUserId === commentUserId || props.currentUser.role === 'admin'
    }

    const localReplyContent = ref(props.localReplyContent)
    const localReplyIsAnonymous = ref(props.localReplyIsAnonymous)

    watch(() => props.localReplyContent, (newVal) => {
      localReplyContent.value = newVal
    })

    watch(() => props.localReplyIsAnonymous, (newVal) => {
      localReplyIsAnonymous.value = newVal
    })

    const submitReply = (parentId) => {
      emit('submitReply', { parentId, content: localReplyContent.value, isAnonymous: localReplyIsAnonymous.value })
    }

    const cancelReply = () => {
      emit('cancelReply')
    }

    return {
      formatDate,
      canDeleteComment,
      submitReply,
      cancelReply,
      localReplyContent,
      localReplyIsAnonymous
    }
  }
}
</script>

<style scoped>
.reply-list {
  padding-left: 24px;
}

.reply-item {
  padding: 12px 0;
  border-bottom: 1px solid #f5f5f5;
}

.reply-header {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  margin-bottom: 8px;
}

.reply-avatar {
  flex-shrink: 0;
  cursor: pointer;
}

.avatar-img {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  object-fit: cover;
}

.reply-info {
  flex: 1;
}

.reply-author {
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

.reply-time {
  font-size: 12px;
  color: #999;
}

.reply-content {
  font-size: 14px;
  color: #333;
  line-height: 1.6;
  margin-bottom: 8px;
  padding-left: 48px;
}

.reply-actions {
  display: flex;
  gap: 16px;
  padding-left: 48px;
}

.reply-action-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  background: transparent;
  border: none;
  color: #999;
  font-size: 13px;
  cursor: pointer;
  transition: color 0.3s ease;

  &:hover {
    color: #1E9FFF;
  }

  &.liked {
    color: #1E9FFF;
  }

  &.reply-action-delete:hover {
    color: #FF5722;
  }

  i {
    font-size: 12px;
  }
}

.reply-form {
  margin-top: 12px;
  padding: 12px;
  background: #f8f9fa;
  border-radius: 8px;
  margin-left: 48px;
}

.reply-textarea {
  width: 100%;
  padding: 10px;
  border: 1px solid #e0e0e0;
  border-radius: 6px;
  font-size: 14px;
  resize: none;

  &:focus {
    outline: none;
    border-color: #1E9FFF;
  }
}

.reply-form-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 10px;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #666;
  cursor: pointer;
}

.reply-form-actions {
  display: flex;
  gap: 8px;
}
</style>