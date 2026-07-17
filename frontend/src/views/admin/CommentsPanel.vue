<template>
  <div>
    <div class="layui-card mb-4">
      <div class="layui-card-body">
        <div class="layui-row align-items-center">
          <div class="layui-col-xs12 layui-col-sm8">
            <div class="layui-input-wrap">
              <i class="fa-solid fa-magnifying-glass layui-input-prefix" style="color: #999;"></i>
              <input type="text" v-model="searchQuery" placeholder="搜索评论内容..." class="layui-input" />
            </div>
          </div>
          <div class="layui-col-xs12 layui-col-sm4 text-right">
            <button class="layui-btn" @click="$emit('refresh')" :disabled="loading">
              <i class="fa-solid fa-refresh mr-2"></i>
              刷新
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="layui-card">
      <div v-if="filteredComments.length > 0">
        <div v-for="comment in filteredComments" :key="comment.id" class="layui-card-body border-b last:border-b-0" style="padding: 15px;">
          <div class="flex items-center gap-3">
            <UserAvatar :user="getUserInfo(comment)" :size="48" />
            <div class="flex-1">
              <div class="font-medium">
                {{ getUserName(comment) }}
                <span v-if="comment.is_anonymous" class="layui-badge layui-bg-gray ml-2">匿名</span>
              </div>
              <div class="text-muted mt-2">{{ comment.content }}</div>
              <div class="flex flex-wrap gap-3 mt-2 text-sm text-muted">
                <span v-if="getArticleInfo(comment)">
                  <i class="fa-solid fa-file-lines mr-1" style="color: #1E9FFF;"></i>
                  {{ getArticleInfo(comment).title }}
                </span>
                <span>
                  <i class="fa-solid fa-heart mr-1" style="color: #FF69B4;"></i>
                  {{ comment.like_count || 0 }}
                </span>
                <span>
                  <i class="fa-solid fa-reply mr-1" style="color: #3B82F6;"></i>
                  {{ comment.reply_count || 0 }}
                </span>
                <span>
                  <i class="fa-solid fa-clock mr-1" style="color: #999;"></i>
                  {{ formatDate(comment.created_at) }}
                </span>
              </div>
            </div>
            <div class="flex gap-2">
              <a v-if="comment.article_id" :href="`/article/${comment.article_id}`" target="_blank" class="layui-btn layui-btn-primary layui-btn-sm">
                <i class="fa-solid fa-eye"></i>
              </a>
              <button class="layui-btn layui-btn-danger layui-btn-sm" @click="$emit('delete', comment.id)">
                <i class="fa-solid fa-trash"></i>
              </button>
            </div>
          </div>
        </div>
      </div>
      <div v-else class="layui-card-body text-center py-8">
        <i class="fa-solid fa-message" style="font-size: 48px; color: #dcdcdc;"></i>
        <div class="text-muted mt-2">{{ searchQuery ? '未找到匹配的评论' : '暂无评论数据' }}</div>
      </div>
    </div>

    <div class="flex items-center justify-center gap-4 mt-4" v-if="pagination.totalPages > 1">
      <div class="layui-laypage">
        <button 
          v-for="page in visiblePages" 
          :key="page"
          class="layui-laypage-btn"
          :class="{ 'layui-laypage-curr': page === pagination.page }"
          @click="handlePageChange(page)"
        >
          {{ page }}
        </button>
      </div>
      <div class="text-sm text-muted">
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

    const visiblePages = computed(() => {
      const total = props.pagination.totalPages
      const current = props.pagination.page
      const pages = []
      const start = Math.max(1, current - 2)
      const end = Math.min(total, current + 2)
      for (let i = start; i <= end; i++) {
        pages.push(i)
      }
      return pages
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
      visiblePages,
      formatDate,
      getUserInfo,
      getUserName,
      getArticleInfo,
      handlePageChange
    }
  }
}
</script>

<style scoped>
.last\:border-b-0:last-child {
  border-bottom: none;
}
</style>