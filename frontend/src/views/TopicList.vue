<template>
  <div class="topic-list-page">
    <div class="topic-card">
      <div class="card-header">
        <div class="header-left">
          <i class="fa-solid fa-tags mr-2"></i>
          话题广场
        </div>
        <button class="layui-btn layui-btn-normal layui-btn-sm" @click="showHotTopics">
          <i class="fa-solid fa-flame mr-2"></i>
          热门话题
        </button>
      </div>

      <div class="card-body">
        <div class="topics-grid">
          <div 
            v-for="topic in topics" 
            :key="topic.id" 
            class="topic-item"
            @click="viewTopic(topic.id)"
          >
            <div class="topic-header">
              <div class="topic-avatar">
                <i :class="getIconClass(topic.icon)"></i>
              </div>
              <div class="topic-info">
                <div class="topic-name">{{ topic.display_name || topic.name }}</div>
                <div class="topic-desc">{{ topic.description }}</div>
              </div>
            </div>
            
            <div class="topic-stats">
              <span class="stat-tag">
                <i class="fa-solid fa-file-lines mr-1"></i>
                {{ topic.article_count }} 文章
              </span>
              <span class="stat-tag">
                <i class="fa-solid fa-user mr-1"></i>
                {{ topic.follow_count }} 关注
              </span>
            </div>

            <span v-if="topic.is_hot" class="hot-tag">
              <i class="fa-solid fa-flame mr-1"></i>
              热门
            </span>
          </div>
        </div>

        <div v-if="totalPages > 1" class="pagination">
          <button 
            class="page-btn" 
            :disabled="page <= 1"
            @click="page--; loadTopics()"
          >
            <i class="fa-solid fa-chevron-left"></i>
          </button>
          <span class="page-info">{{ page }} / {{ totalPages }}</span>
          <button 
            class="page-btn" 
            :disabled="page >= totalPages"
            @click="page++; loadTopics()"
          >
            <i class="fa-solid fa-chevron-right"></i>
          </button>
        </div>
      </div>
    </div>

    <div v-if="topicDialog" class="dialog-overlay" @click.self="topicDialog = false">
      <div class="topic-dialog">
        <div class="dialog-header">
          <div class="header-left">
            <div class="topic-avatar">
              <i :class="getIconClass(selectedTopic?.icon)"></i>
            </div>
            <span>{{ selectedTopic?.display_name || selectedTopic?.name }}</span>
          </div>
          <button 
            class="layui-btn layui-btn-normal layui-btn-sm"
            @click="followTopic"
          >
            <i :class="isFollowing ? 'fa-solid fa-heart' : 'fa-regular fa-heart'" class="mr-2"></i>
            {{ isFollowing ? '已关注' : '关注' }}
          </button>
        </div>

        <div class="dialog-body">
          <p class="topic-description">{{ selectedTopic?.description }}</p>
          
          <div class="divider"></div>
          
          <div class="related-articles">
            <h3 class="articles-title">相关文章</h3>
            <div class="articles-list">
              <div 
                v-for="articleTopic in topicArticles" 
                :key="articleTopic.id" 
                class="article-item"
              >
                <div class="article-info">
                  <div class="article-title">{{ articleTopic.article.title }}</div>
                  <div class="article-meta">
                    {{ articleTopic.article.user.display_name || articleTopic.article.user.username }}
                    · {{ formatDate(articleTopic.article.created_at) }}
                  </div>
                </div>
                <router-link :to="`/article/${articleTopic.article.id}`" class="view-link">
                  查看
                </router-link>
              </div>
            </div>
          </div>
        </div>

        <div class="dialog-footer">
          <button class="layui-btn layui-btn-primary" @click="topicDialog = false">关闭</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { topicApi } from '@/api'

export default {
  name: 'TopicList',
  setup() {
    const topics = ref([])
    const page = ref(1)
    const limit = ref(20)
    const total = ref(0)
    const topicDialog = ref(false)
    const selectedTopic = ref(null)
    const topicArticles = ref([])
    const isFollowing = ref(false)

    const mdiToFa = {
      'mdi-tag': 'fa-solid fa-tag',
      'mdi-file-document': 'fa-solid fa-file-lines',
      'mdi-account': 'fa-solid fa-user',
      'mdi-fire': 'fa-solid fa-flame',
      'mdi-heart': 'fa-solid fa-heart',
      'mdi-heart-outline': 'fa-regular fa-heart'
    }

    const getIconClass = (icon) => {
      return mdiToFa[icon] || 'fa-solid fa-tag'
    }

    const totalPages = computed(() => {
      return Math.ceil(total.value / limit.value)
    })

    const loadTopics = async () => {
      try {
        const res = await topicApi.getTopics(page.value, limit.value)
        if (res.data.success) {
          topics.value = res.data.data.topics
          total.value = res.data.data.total
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
          topicArticles.value = res.data.data.articles
          topicDialog.value = true

          const followedRes = await topicApi.getFollowedTopics()
          if (followedRes.data.success) {
            isFollowing.value = followedRes.data.data.some(f => f.topic_id === topicId)
          }
        }
      } catch (error) {
        console.error('加载话题详情失败:', error)
      }
    }

    const followTopic = async () => {
      try {
        if (isFollowing.value) {
          await topicApi.unfollowTopic(selectedTopic.value.id)
          isFollowing.value = false
        } else {
          await topicApi.followTopic(selectedTopic.value.id)
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
          topics.value = res.data.data.topics
          total.value = res.data.data.total
        }
      } catch (error) {
        console.error('加载热门话题失败:', error)
      }
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
      loadTopics,
      viewTopic,
      followTopic,
      showHotTopics,
      formatDate,
      getIconClass
    }
  }
}
</script>

<style scoped>
.topic-list-page {
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;
}

.topic-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
  overflow: hidden;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 24px;
  border-bottom: 1px solid #f0f0f0;
  font-size: 18px;
  font-weight: 600;
}

.header-left {
  display: flex;
  align-items: center;
}

.card-body {
  padding: 24px;
}

.topics-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;
}

.topic-item {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 16px;
  cursor: pointer;
  transition: all 0.2s;
  
  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }
}

.topic-header {
  display: flex;
  gap: 12px;
  margin-bottom: 12px;
}

.topic-avatar {
  width: 40px;
  height: 40px;
  background: var(--primary);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
  
  i {
    font-size: 18px;
  }
}

.topic-info {
  flex: 1;
  min-width: 0;
}

.topic-name {
  font-size: 15px;
  font-weight: 600;
  color: #333;
  margin-bottom: 4px;
}

.topic-desc {
  font-size: 13px;
  color: #999;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.topic-stats {
  display: flex;
  gap: 12px;
  margin-bottom: 12px;
}

.stat-tag {
  font-size: 12px;
  color: #666;
  background: white;
  padding: 4px 8px;
  border-radius: 4px;
}

.hot-tag {
  font-size: 12px;
  color: white;
  background: #FF5722;
  padding: 4px 8px;
  border-radius: 4px;
}

.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  margin-top: 24px;
}

.page-btn {
  width: 36px;
  height: 36px;
  border: 1px solid #e8e8e8;
  border-radius: 6px;
  background: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  
  &:hover:not(:disabled) {
    border-color: var(--primary);
    color: var(--primary);
  }
  
  &:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }
}

.page-info {
  font-size: 14px;
  color: #666;
}

.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 24px;
}

.topic-dialog {
  background: white;
  border-radius: 12px;
  width: 100%;
  max-width: 800px;
  overflow: hidden;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
}

.dialog-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 24px;
  border-bottom: 1px solid #f0f0f0;
  font-size: 16px;
  font-weight: 600;
}

.dialog-body {
  padding: 24px;
  overflow-y: auto;
}

.topic-description {
  font-size: 14px;
  color: #666;
  margin-bottom: 16px;
}

.divider {
  height: 1px;
  background: #f0f0f0;
  margin: 16px 0;
}

.related-articles {
  margin-top: 16px;
}

.articles-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-bottom: 12px;
}

.articles-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.article-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px;
  background: #f8f9fa;
  border-radius: 6px;
}

.article-info {
  flex: 1;
}

.article-title {
  font-size: 14px;
  color: #333;
  margin-bottom: 4px;
}

.article-meta {
  font-size: 12px;
  color: #999;
}

.view-link {
  font-size: 12px;
  color: var(--primary);
}

.dialog-footer {
  padding: 16px 24px;
  border-top: 1px solid #f0f0f0;
  text-align: right;
}

.mr-1 {
  margin-right: 4px;
}

.mr-2 {
  margin-right: 8px;
}
</style>
