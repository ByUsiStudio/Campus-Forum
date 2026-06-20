<template>
  <v-container>
    <v-card>
      <v-card-title class="d-flex align-center">
        <v-icon left>mdi-tag-multiple</v-icon>
        话题广场
        <v-btn text color="primary" class="ml-auto" @click="showHotTopics">
          <v-icon left>mdi-fire</v-icon>
          热门话题
        </v-btn>
      </v-card-title>

      <v-card-text>
        <!-- 话题列表 -->
        <v-row>
          <v-col cols="12" sm="6" md="4" v-for="topic in topics" :key="topic.id">
            <v-card outlined hover class="topic-card" @click="viewTopic(topic.id)">
              <v-card-text>
                <div class="d-flex align-center mb-2">
                  <v-avatar size="40" color="primary">
                    <v-icon dark>{{ topic.icon || 'mdi-tag' }}</v-icon>
                  </v-avatar>
                  <div class="ml-3">
                    <div class="title">{{ topic.display_name || topic.name }}</div>
                    <div class="caption">{{ topic.description }}</div>
                  </div>
                </div>
                
                <v-row class="mt-2">
                  <v-col cols="6">
                    <v-chip small>
                      <v-icon small left>mdi-file-document</v-icon>
                      {{ topic.article_count }} 文章
                    </v-chip>
                  </v-col>
                  <v-col cols="6">
                    <v-chip small>
                      <v-icon small left>mdi-account</v-icon>
                      {{ topic.follow_count }} 关注
                    </v-chip>
                  </v-col>
                </v-row>

                <div v-if="topic.is_hot" class="mt-2">
                  <v-chip small color="orange">
                    <v-icon small left>mdi-fire</v-icon>
                    热门
                  </v-chip>
                </div>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>

        <!-- 分页 -->
        <div class="text-center mt-4">
          <v-pagination
            v-model="page"
            :length="totalPages"
            :total-visible="7"
            @input="loadTopics"
          ></v-pagination>
        </div>
      </v-card-text>
    </v-card>

    <!-- 话题详情对话框 -->
    <v-dialog v-model="topicDialog" max-width="800">
      <v-card v-if="selectedTopic">
        <v-card-title>
          <v-avatar size="40" color="primary" class="mr-3">
            <v-icon dark>{{ selectedTopic.icon || 'mdi-tag' }}</v-icon>
          </v-avatar>
          {{ selectedTopic.display_name || selectedTopic.name }}
          <v-btn text color="primary" class="ml-auto" @click="followTopic">
            <v-icon left>{{ isFollowing ? 'mdi-heart' : 'mdi-heart-outline' }}</v-icon>
            {{ isFollowing ? '已关注' : '关注' }}
          </v-btn>
        </v-card-title>

        <v-card-text>
          <div class="mb-4">{{ selectedTopic.description }}</div>

          <v-divider class="my-4"></v-divider>

          <div class="subtitle-1 mb-2">相关文章</div>
          <v-list three-line>
            <v-list-item v-for="articleTopic in topicArticles" :key="articleTopic.id">
              <v-list-item-content>
                <v-list-item-title>{{ articleTopic.article.title }}</v-list-item-title>
                <v-list-item-subtitle>
                  {{ articleTopic.article.user.display_name || articleTopic.article.user.username }}
                  · {{ formatDate(articleTopic.article.created_at) }}
                </v-list-item-subtitle>
              </v-list-item-content>
              <v-list-item-action>
                <v-btn text small color="primary" :to="`/article/${articleTopic.article.id}`">
                  查看
                </v-btn>
              </v-list-item-action>
            </v-list-item>
          </v-list>
        </v-card-text>

        <v-card-actions>
          <v-btn text @click="topicDialog = false">关闭</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
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

          // 检查是否已关注
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
      formatDate
    }
  }
}
</script>

<style scoped>
.topic-card {
  cursor: pointer;
  transition: transform 0.2s;
}

.topic-card:hover {
  transform: translateY(-4px);
}
</style>