<template>
  <el-list class="reply-list">
    <el-list-item
      v-for="reply in replies"
      :key="reply.id"
      class="py-2 reply-item"
    >
      <div class="reply-body">
        <div class="reply-header">
          <UserAvatar
            :user="reply.user"
            :size="32"
            :showUsername="false"
            class="cursor-pointer reply-avatar"
            @click="$emit('goToUserProfile', reply.user.id)"
          />
          <div class="reply-meta">
            <span class="reply-username">
              {{ reply.user.display_name || reply.user.username || '匿名用户' }}
            </span>
            <span class="reply-time text-secondary">
              {{ formatDate(reply.created_at) }}
            </span>
          </div>
        </div>

        <div class="reply-content comment-text">
          {{ reply.content }}
        </div>

        <div class="reply-actions">
          <el-button
            text
            size="small"
            class="reply-action-btn"
            @click="$emit('toggleLike', reply)"
          >
            <el-icon class="action-icon"><ThumbUp /></el-icon>
            <span :class="commentLiked[reply.id] ? 'liked-text' : ''">
              {{ reply.like_count }}
            </span>
          </el-button>
          <el-button
            v-if="token"
            text
            size="small"
            class="reply-action-btn"
            @click="$emit('showReplyForm', reply.id)"
          >
            <el-icon class="action-icon"><ChatDotRound /></el-icon>
            回复
          </el-button>
          <el-button
            v-if="canDeleteComment(reply)"
            text
            size="small"
            class="reply-action-btn danger"
            @click="$emit('deleteComment', reply.id, reply)"
          >
            <el-icon class="action-icon"><Delete /></el-icon>
          </el-button>
        </div>

        <div v-if="replyingTo === reply.id" class="reply-form mt-2">
          <el-textarea
            v-model="localReplyContent"
            :placeholder="'回复 ' + (reply.user.display_name || reply.user.username || '匿名用户')"
            :rows="2"
            resize="none"
            class="reply-textarea"
          />
          <div class="d-flex align-center gap-2 mt-2 form-footer">
            <el-checkbox v-model="localReplyIsAnonymous" size="small">
              匿名
            </el-checkbox>
            <div class="flex-spacer"></div>
            <el-button size="small" type="primary" @click="submitReply(reply.id)">
              发送
            </el-button>
            <el-button size="small" text @click="cancelReply">取消</el-button>
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
    </el-list-item>
  </el-list>
</template>

<script>
import { ref, watch } from 'vue'
import { ThumbUp, ChatDotRound, Delete } from '@element-plus/icons-vue'
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
  padding-left: 16px;
  background: rgba(103, 80, 164, 0.04);
  border-radius: 8px;
  border-left: 2px solid var(--campus-border);
}

.reply-item {
  padding: 8px 0;
}

.reply-item + .reply-item {
  border-top: 1px solid var(--campus-border);
}

.reply-body {
  width: 100%;
  min-width: 0;
}

.reply-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.reply-avatar {
  flex-shrink: 0;
}

.reply-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.reply-username {
  font-size: 13px;
  font-weight: 500;
  color: var(--campus-text);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.reply-time {
  font-size: 12px;
}

.reply-content {
  font-size: 13px;
  color: var(--campus-text);
  margin-bottom: 6px;
}

.reply-actions {
  display: flex;
  gap: 4px;
}

.reply-action-btn {
  padding: 4px 6px;
  height: auto;
  min-height: 0;
  font-size: 12px;
}

.reply-action-btn.danger {
  color: #f56c6c;
}

.reply-action-btn.danger:hover {
  background-color: rgba(245, 108, 108, 0.1);
}

.action-icon {
  margin-right: 2px;
  font-size: 14px;
}

.liked-text {
  color: var(--campus-primary);
  font-weight: 500;
}

.reply-form {
  margin-top: 8px;
}

.reply-textarea {
  width: 100%;
}

.form-footer {
  margin-top: 8px;
}

.d-flex {
  display: flex;
}

.align-center {
  align-items: center;
}

.gap-2 {
  gap: 8px;
}

.mt-2 {
  margin-top: 8px;
}

.flex-spacer {
  flex: 1;
}

.comment-text {
  white-space: pre-wrap;
  word-break: break-word;
}
</style>
