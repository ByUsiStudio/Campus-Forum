<template>
  <div>
    <!-- 搜索栏 -->
    <v-card class="mb-4 pa-3" variant="flat" rounded="lg">
      <v-row dense align="center">
        <v-col cols="12" sm="8">
          <v-text-field
            v-model="searchQuery"
            placeholder="搜索评论内容..."
            prepend-inner-icon="mdi-magnify"
            variant="outlined"
            density="compact"
            hide-details
            clearable
          />
        </v-col>
        <v-col cols="12" sm="4" class="text-end">
          <v-btn variant="tonal" color="primary" @click="$emit('refresh')" :loading="loading" prepend-icon="mdi-refresh">
            刷新
          </v-btn>
        </v-col>
      </v-row>
    </v-card>

    <!-- 评论列表 -->
    <v-card variant="flat" rounded="lg">
      <v-list lines="three" v-if="filteredComments.length > 0">
        <v-list-item v-for="comment in filteredComments" :key="comment.id" class="py-3">
          <template v-slot:prepend>
            <UserAvatar :user="getUserInfo(comment)" :size="48" />
          </template>

          <v-list-item-title class="font-weight-medium mb-1">
            {{ getUserName(comment) }}
            <v-chip size="x-small" color="grey" variant="tonal" class="ml-2" v-if="comment.is_anonymous">
              匿名
            </v-chip>
          </v-list-item-title>

          <v-list-item-subtitle>
            <div class="mb-2">{{ comment.content }}</div>
            <div class="d-flex flex-wrap align-center ga-2">
              <span class="d-flex align-center text-caption" v-if="getArticleInfo(comment)">
                <v-icon size="14" color="primary" class="mr-1">mdi-file-document</v-icon>
                {{ getArticleInfo(comment).title }}
              </span>
              <span class="d-flex align-center text-caption">
                <v-icon size="14" color="pink" class="mr-1">mdi-heart</v-icon>
                {{ comment.like_count || 0 }}
              </span>
              <span class="d-flex align-center text-caption">
                <v-icon size="14" color="blue" class="mr-1">mdi-reply</v-icon>
                {{ comment.reply_count || 0 }}
              </span>
              <span class="d-flex align-center text-caption">
                <v-icon size="14" color="grey" class="mr-1">mdi-clock-outline</v-icon>
                {{ formatDate(comment.created_at) }}
              </span>
            </div>
          </v-list-item-subtitle>

          <template v-slot:append>
            <v-btn-group variant="text" density="compact" divided>
              <v-btn size="small" color="info" :to="`/article/${comment.article_id}`" target="_blank" v-if="comment.article_id">
                <v-icon>mdi-eye</v-icon>
                <v-tooltip activator="parent">查看文章</v-tooltip>
              </v-btn>
              <v-btn size="small" color="error" @click="$emit('delete', comment.id)">
                <v-icon>mdi-delete</v-icon>
                <v-tooltip activator="parent">删除</v-tooltip>
              </v-btn>
            </v-btn-group>
          </template>
        </v-list-item>
      </v-list>

      <v-card-text v-else class="text-center py-8">
        <v-icon size="48" color="grey-lighten-1">mdi-comment-text-outline</v-icon>
        <div class="text-body-1 text-medium-emphasis mt-2">
          {{ searchQuery ? '未找到匹配的评论' : '暂无评论数据' }}
        </div>
      </v-card-text>
    </v-card>

    <!-- 分页 -->
    <div class="d-flex align-center justify-center ga-4 mt-4" v-if="pagination.totalPages > 1">
      <v-pagination
        v-model="currentPage"
        :length="pagination.totalPages"
        :total-visible="5"
        rounded="lg"
        @update:model-value="handlePageChange"
      />
      <div class="text-caption text-medium-emphasis">
        第 {{ pagination.page }} / {{ pagination.totalPages }} 页 (共 {{ pagination.total }} 条)
      </div>
    </div>
  </div>
</template>

<script>
import { ref, watch, computed } from 'vue'
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
    },
    pagination: {
      type: Object,
      default: () => ({
        page: 1,
        pageSize: 20,
        total: 0,
        totalPages: 0
      })
    }
  },
  emits: ['delete', 'refresh', 'page-change'],
  setup(props, { emit }) {
    const currentPage = ref(1)
    const searchQuery = ref('')

    watch(() => props.pagination.page, (newVal) => {
      currentPage.value = newVal
    })

    const filteredComments = computed(() => {
      if (!searchQuery.value) return props.comments
      const query = searchQuery.value.toLowerCase()
      return props.comments.filter(comment => 
        comment.content?.toLowerCase().includes(query) ||
        getUserName(comment).toLowerCase().includes(query)
      )
    })

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

    const getUserInfo = (comment) => {
      return comment.user || comment.User || {}
    }

    const getUserName = (comment) => {
      const user = getUserInfo(comment)
      return user.display_name || user.DisplayName || user.username || user.Username || '未知用户'
    }

    const getArticleInfo = (comment) => {
      return comment.article || comment.Article || null
    }

    const handlePageChange = (page) => {
      emit('page-change', page)
    }

    return {
      currentPage,
      searchQuery,
      filteredComments,
      formatDate,
      getUserInfo,
      getUserName,
      getArticleInfo,
      handlePageChange
    }
  }
}
</script>