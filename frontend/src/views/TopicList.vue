<template>
  <div class="page-container topic-list-page">
    <el-card class="mb-4 card-surface">
      <template #header>
        <div class="topic-header">
          <el-icon :size="22" class="topic-header-icon"><Collection /></el-icon>
          <span class="topic-header-title">话题广场</span>
          <el-button
            type="primary"
            plain
            class="topic-header-btn"
            @click="showHotTopics"
          >
            <el-icon class="mr-1"><TrendCharts /></el-icon>
            热门话题
          </el-button>
        </div>
      </template>

      <!-- 话题卡片网格 -->
      <div v-if="topics.length" class="topic-grid">
        <div
          v-for="topic in topics"
          :key="topic.id"
          class="topic-grid-item"
        >
          <el-card
            shadow="hover"
            class="topic-card cursor-pointer"
            :body-style="{ padding: '16px' }"
            @click="viewTopic(topic.id)"
          >
            <div class="topic-main">
              <el-avatar :size="40" class="topic-avatar">
                <el-icon :size="22"><PriceTag /></el-icon>
              </el-avatar>
              <div class="topic-info">
                <div class="topic-name">{{ topic.display_name || topic.name }}</div>
                <div class="topic-desc text-secondary">{{ topic.description }}</div>
              </div>
            </div>

            <div class="topic-stats mt-2">
              <el-tag size="small" effect="plain" class="topic-stat-tag">
                <el-icon class="mr-1"><Document /></el-icon>
                {{ topic.article_count }} 文章
              </el-tag>
              <el-tag size="small" effect="plain" class="topic-stat-tag">
                <el-icon class="mr-1"><User /></el-icon>
                {{ topic.follow_count }} 关注
              </el-tag>
            </div>

            <div v-if="topic.is_hot" class="mt-2">
              <el-tag size="small" type="warning" effect="light">
                <el-icon class="mr-1"><TrendCharts /></el-icon>
                热门
              </el-tag>
            </div>
          </el-card>
        </div>
      </div>

      <el-empty v-else description="暂无话题" />

      <!-- 分页 -->
      <div v-if="totalPages > 1" class="text-center mt-4">
        <el-pagination
          layout="prev, pager, next"
          v-model:current-page="page"
          :page-size="limit"
          :total="total"
          @current-change="loadTopics"
        />
      </div>
    </el-card>

    <!-- 话题详情弹窗 -->
    <el-dialog
      v-model="topicDialog"
      width="800px"
      :title="selectedTopic ? (selectedTopic.display_name || selectedTopic.name) : ''"
    >
      <template v-if="selectedTopic">
        <div class="topic-dialog-header mb-4">
          <el-avatar :size="40" class="topic-avatar mr-2">
            <el-icon :size="22"><PriceTag /></el-icon>
          </el-avatar>
          <el-button
            type="primary"
            :plain="!isFollowing"
            @click="followTopic"
          >
            <el-icon class="mr-1">
              <Star v-if="isFollowing" />
              <StarFilled v-else />
            </el-icon>
            {{ isFollowing ? '已关注' : '关注' }}
          </el-button>
        </div>

        <div class="mb-4 text-secondary">{{ selectedTopic.description }}</div>

        <el-divider class="my-4" />

        <div class="topic-articles-title mb-2">相关文章</div>
        <el-list v-if="topicArticles.length" class="topic-articles-list">
          <el-list-item
            v-for="articleTopic in topicArticles"
            :key="articleTopic.id"
            class="topic-article-item"
          >
            <div class="topic-article-main">
              <div class="topic-article-title">{{ articleTopic.article.title }}</div>
              <div class="topic-article-sub text-secondary">
                {{ articleTopic.article.user.display_name || articleTopic.article.user.username }}
                · {{ formatDate(articleTopic.article.created_at) }}
              </div>
            </div>
            <el-button
              link
              type="primary"
              @click="goArticle(articleTopic.article.id)"
            >
              查看
            </el-button>
          </el-list-item>
        </el-list>
        <el-empty v-else description="暂无相关文章" />
      </template>

      <template #footer>
        <el-button type="primary" plain @click="topicDialog = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  Collection,
  TrendCharts,
  Document,
  User,
  PriceTag,
  Star,
  StarFilled
} from '@element-plus/icons-vue'
import { topicApi } from '@/api'

export default {
  name: 'TopicList',
  setup() {
    const router = useRouter()
    const topics = ref([])
    const page = ref(1)
    const limit = ref(20)
    const total = ref(0)
    const topicDialog = ref(false)
    const selectedTopic = ref(null)
    const topicArticles = ref([])
    const isFollowing = ref(false)

    const totalPages = computed(() => {
      return Math.ceil(total.value / limit.value)
    })

    const loadTopics = async () => {
      try {
        const res = await topicApi.getTopics(page.value, limit.value)
        if (res.data.success) {
          topics.value = res.data.data.topics || []
          total.value = res.data.data.total || 0
        }
      } catch (error) {
        console.error('加载话题失败:', error)
      }
    }

    const viewTopic = async (topicId) => {
      try {
        const res = await topicApi.getTopic(topicId)
        if (res.data.success) {
          selectedTopic.value = res.data.data.topic
          topicArticles.value = res.data.data.articles || []
          topicDialog.value = true

          const followedRes = await topicApi.getFollowedTopics()
          const followedList = (followedRes.data.success && followedRes.data.data) || []
          isFollowing.value = followedList.some(f => f.topic_id === topicId)
        }
      } catch (error) {
        console.error('加载话题详情失败:', error)
      }
    }

    const followTopic = async () => {
      const topicId = selectedTopic.value?.id
      if (!topicId) return

      try {
        if (isFollowing.value) {
          await topicApi.unfollowTopic(topicId)
          isFollowing.value = false
        } else {
          await topicApi.followTopic(topicId)
          isFollowing.value = true
        }
      } catch (error) {
        console.error('关注操作失败:', error)
      }
    }

    const showHotTopics = async () => {
      try {
        const res = await topicApi.getTopics(1, 20, true)
        if (res.data.success) {
          topics.value = res.data.data.topics || []
          total.value = res.data.data.total || 0
        }
      } catch (error) {
        console.error('加载热门话题失败:', error)
      }
    }

    const goArticle = (articleId) => {
      router.push(`/article/${articleId}`)
    }

    const formatDate = (date) => {
      return new Date(date).toLocaleDateString('zh-CN')
    }

    onMounted(() => {
      loadTopics()
    })

    return {
      topics,
      page,
      total,
      totalPages,
      topicDialog,
      selectedTopic,
      topicArticles,
      isFollowing,
      Collection,
      TrendCharts,
      Document,
      User,
      PriceTag,
      Star,
      StarFilled,
      loadTopics,
      viewTopic,
      followTopic,
      showHotTopics,
      goArticle,
      formatDate
    }
  }
}
</script>

<style scoped>
.topic-list-page {
  padding-top: 16px;
}

.topic-card {
  border-radius: var(--campus-radius);
  transition: transform 0.2s;
}

.topic-card:hover {
  transform: translateY(-4px);
}

.topic-header {
  display: flex;
  align-items: center;
  font-size: 1.25rem;
  font-weight: 600;
  gap: 8px;
}

.topic-header-icon {
  color: var(--campus-primary);
}

.topic-header-btn {
  margin-left: auto;
}

.topic-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}

.topic-main {
  display: flex;
  align-items: center;
  gap: 12px;
}

.topic-avatar {
  flex-shrink: 0;
  color: #fff;
  font-weight: 600;
}

.topic-info {
  flex: 1;
  min-width: 0;
}

.topic-name {
  font-weight: 600;
  font-size: 1.05rem;
  color: var(--campus-text);
}

.topic-desc {
  font-size: 0.85rem;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.topic-stats {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.topic-stat-tag {
  display: inline-flex;
  align-items: center;
}

.topic-dialog-header {
  display: flex;
  align-items: center;
  gap: 8px;
  justify-content: space-between;
}

.topic-articles-title {
  font-weight: 600;
  font-size: 1rem;
}

.topic-article-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.topic-article-main {
  flex: 1;
  min-width: 0;
}

.topic-article-title {
  font-weight: 500;
  color: var(--campus-text);
}

.topic-article-sub {
  font-size: 0.85rem;
}
</style>
