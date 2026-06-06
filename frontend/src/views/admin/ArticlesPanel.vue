<template>
  <div>
    <!-- 搜索和筛选栏 -->
    <v-card class="mb-4 pa-3" variant="flat" rounded="lg">
      <v-row dense align="center">
        <v-col cols="12" sm="6" md="4">
          <v-text-field
            v-model="searchQuery"
            placeholder="搜索文章标题..."
            prepend-inner-icon="mdi-magnify"
            variant="outlined"
            density="compact"
            hide-details
            clearable
          />
        </v-col>
        <v-col cols="12" sm="6" md="3">
          <v-select
            v-model="localFilter"
            :items="statusOptions"
            item-title="label"
            item-value="value"
            label="状态筛选"
            variant="outlined"
            density="compact"
            hide-details
            @update:model-value="$emit('update:filter', $event)"
          />
        </v-col>
        <v-col cols="12" sm="12" md="5" class="text-end">
          <v-btn-group variant="tonal" density="compact">
            <v-btn color="primary" @click="$emit('refresh')" :loading="loading">
              <v-icon start>mdi-refresh</v-icon>
              刷新
            </v-btn>
          </v-btn-group>
        </v-col>
      </v-row>
    </v-card>

    <!-- 文章列表 -->
    <v-card variant="flat" rounded="lg">
      <v-list lines="two" v-if="filteredArticles.length > 0">
        <v-list-item v-for="article in filteredArticles" :key="article.id" class="py-3">
          <template v-slot:prepend>
            <v-avatar size="48" color="primary" variant="tonal">
              <v-icon>mdi-file-document</v-icon>
            </v-avatar>
          </template>

          <v-list-item-title class="font-weight-medium mb-1">
            {{ article.title }}
            <v-chip size="x-small" :color="getStatusColor(article.status)" variant="tonal" class="ml-2">
              {{ getStatusText(article.status) }}
            </v-chip>
          </v-list-item-title>

          <v-list-item-subtitle>
            <div class="d-flex flex-wrap align-center ga-2">
              <span class="d-flex align-center">
                <v-icon size="14" class="mr-1">mdi-identifier</v-icon>
                ID: {{ article.id }}
              </span>
              <span class="d-flex align-center">
                <UserAvatar :user="article.User || {}" :size="20" class="mr-1" />
                {{ article.User?.display_name || '-' }}
              </span>
              <v-chip size="x-small" variant="tonal" color="primary" v-if="article.Category">
                {{ article.Category.name }}
              </v-chip>
            </div>
            <div class="d-flex flex-wrap align-center ga-3 mt-1">
              <span class="d-flex align-center text-caption">
                <v-icon size="14" color="pink" class="mr-1">mdi-heart</v-icon>
                {{ article.like_count || 0 }}
              </span>
              <span class="d-flex align-center text-caption">
                <v-icon size="14" color="blue" class="mr-1">mdi-eye</v-icon>
                {{ article.view_count || 0 }}
              </span>
              <span class="d-flex align-center text-caption">
                <v-icon size="14" color="grey" class="mr-1">mdi-clock-outline</v-icon>
                {{ formatDate(article.created_at) }}
              </span>
            </div>
          </v-list-item-subtitle>

          <template v-slot:append>
            <v-btn-group variant="text" density="compact" divided>
              <v-btn size="small" color="primary" @click="$emit('change-status', article)">
                <v-icon>mdi-state-machine</v-icon>
                <v-tooltip activator="parent">修改状态</v-tooltip>
              </v-btn>
              <v-btn size="small" color="info" :to="`/article/${article.id}`" target="_blank">
                <v-icon>mdi-eye</v-icon>
                <v-tooltip activator="parent">查看文章</v-tooltip>
              </v-btn>
              <v-btn size="small" color="error" @click="$emit('delete', article)" v-if="canDeleteArticle()">
                <v-icon>mdi-delete</v-icon>
                <v-tooltip activator="parent">删除</v-tooltip>
              </v-btn>
            </v-btn-group>
          </template>
        </v-list-item>
      </v-list>

      <v-card-text v-else class="text-center py-8">
        <v-icon size="48" color="grey-lighten-1">mdi-file-document-search</v-icon>
        <div class="text-body-1 text-medium-emphasis mt-2">
          {{ searchQuery ? '未找到匹配的文章' : '暂无文章数据' }}
        </div>
      </v-card-text>
    </v-card>

    <!-- 分页 -->
    <div class="d-flex align-center justify-center ga-4 mt-4" v-if="totalPages > 1">
      <v-pagination
        v-model="localPage"
        :length="totalPages"
        :total-visible="5"
        rounded="lg"
        @update:model-value="$emit('update:page', $event)"
      />
      <div class="text-caption text-medium-emphasis">
        第 {{ page }} / {{ totalPages }} 页
      </div>
    </div>
  </div>
</template>

<script>
import { ref, watch, computed } from 'vue'
import UserAvatar from '../../components/UserAvatar.vue'

export default {
  name: 'ArticlesPanel',
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
    },
    page: {
      type: Number,
      default: 1
    },
    totalPages: {
      type: Number,
      default: 1
    },
    filter: {
      type: String,
      default: ''
    },
    statusOptions: {
      type: Array,
      default: () => []
    },
    currentUserRole: {
      type: String,
      default: null
    }
  },
  emits: ['change-status', 'delete', 'refresh', 'update:page', 'update:filter'],
  setup(props, { emit }) {
    const localPage = ref(props.page)
    const localFilter = ref(props.filter)
    const searchQuery = ref('')

    watch(() => props.page, (val) => {
      localPage.value = val
    })

    watch(() => props.filter, (val) => {
      localFilter.value = val
    })

    const filteredArticles = computed(() => {
      if (!searchQuery.value) return props.articles
      const query = searchQuery.value.toLowerCase()
      return props.articles.filter(article => 
        article.title?.toLowerCase().includes(query) ||
        article.User?.display_name?.toLowerCase().includes(query)
      )
    })

    const getStatusColor = (status) => {
      const colors = {
        pending: 'warning',
        published: 'success',
        rejected: 'error'
      }
      return colors[status] || 'default'
    }

    const getStatusText = (status) => {
      const texts = {
        pending: '待审核',
        published: '已发布',
        rejected: '已拒绝'
      }
      return texts[status] || status
    }

    const formatDate = (dateString) => {
      if (!dateString) return '-'
      const date = new Date(dateString)
      return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit'
      })
    }

    const canDeleteArticle = () => {
      return props.currentUserRole === 'system' || props.currentUserRole === 'admin'
    }

    return {
      localPage,
      localFilter,
      searchQuery,
      filteredArticles,
      getStatusColor,
      getStatusText,
      formatDate,
      canDeleteArticle
    }
  }
}
</script>