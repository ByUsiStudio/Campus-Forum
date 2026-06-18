<template>
  <v-list class="ml-8 bg-grey-lighten-5 rounded">
    <v-list-item
      v-for="reply in replies"
      :key="reply.id"
      class="py-2"
    >
      <template v-slot:prepend>
        <UserAvatar
          :user="reply.user"
          :size="32"
          :showUsername="false"
          class="cursor-pointer"
          @click="$emit('goToUserProfile', reply.user.id)"
        />
      </template>

      <v-list-item-title class="d-flex align-center gap-2 mb-1">
        <span class="font-weight-medium text-body-2">
          {{ reply.user.display_name || reply.user.username || '匿名用户' }}
        </span>
        <span class="text-caption text-medium-emphasis">
          {{ formatDate(reply.created_at) }}
        </span>
      </v-list-item-title>

      <v-list-item-subtitle class="text-body-2 comment-text">
        {{ reply.content }}
      </v-list-item-subtitle>

      <template v-slot:append>
        <div class="d-flex gap-1">
          <v-btn
            variant="text"
            size="x-small"
            @click="$emit('toggleLike', reply)"
            :color="commentLiked[reply.id] ? 'primary' : 'default'"
          >
            <v-icon size="12">mdi-thumb-up</v-icon>
            {{ reply.like_count }}
          </v-btn>
          <v-btn
            variant="text"
            size="x-small"
            @click="$emit('showReplyForm', reply.id)"
            v-if="token"
          >
            <v-icon size="12">mdi-reply</v-icon>
            回复
          </v-btn>
          <v-btn
            variant="text"
            size="x-small"
            color="error"
            @click="$emit('deleteComment', reply.id, reply)"
            v-if="canDeleteComment(reply)"
          >
            <v-icon size="12">mdi-delete</v-icon>
          </v-btn>
        </div>

        <v-expand-transition>
          <div v-if="replyingTo === reply.id" class="reply-form mt-2">
            <v-textarea
              v-model="localReplyContent"
              :placeholder="'回复 ' + (reply.user.display_name || reply.user.username || '匿名用户')"
              variant="outlined"
              rows="2"
              density="compact"
              hide-details
            />
            <div class="d-flex align-center gap-2 mt-2">
              <v-checkbox
                v-model="localReplyIsAnonymous"
                label="匿名"
                color="primary"
                hide-details
                density="compact"
              />
              <v-spacer></v-spacer>
              <v-btn size="small" color="primary" @click="submitReply(reply.id)">发送</v-btn>
              <v-btn size="small" variant="text" @click="cancelReply">取消</v-btn>
            </div>
          </div>
        </v-expand-transition>
      </template>
    </v-list-item>

    <CommentReply
      v-if="reply.replies && reply.replies.length > 0"
      v-for="reply in replies"
      :key="'nested-' + reply.id"
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
  </v-list>
</template>

<script>
import { ref } from 'vue'
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

    const submitReply = (parentId) => {
      emit('submitReply', { parentId, content: props.localReplyContent, isAnonymous: props.localReplyIsAnonymous })
    }

    const cancelReply = () => {
      emit('cancelReply')
    }

    return {
      formatDate,
      canDeleteComment,
      submitReply,
      cancelReply
    }
  }
}
</script>

<style scoped>
.comment-text {
  white-space: pre-wrap;
  word-break: break-word;
}
</style>