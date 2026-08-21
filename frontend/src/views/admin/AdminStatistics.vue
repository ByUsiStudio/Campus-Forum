<template>
  <v-container fluid>
    <v-row>
      <!-- 系统概览 -->
      <v-col cols="12">
        <v-card>
          <v-card-title>
            <v-icon left>mdi-chart-box</v-icon>
            系统概览
          </v-card-title>
          <v-card-text>
            <v-row>
              <v-col cols="12" sm="6" md="3">
                <v-card outlined>
                  <v-card-text class="text-center">
                    <v-icon color="primary" size="40">mdi-account-group</v-icon>
                    <div class="headline mt-2">{{ overview.total_users }}</div>
                    <div class="subtitle-2">总用户数</div>
                  </v-card-text>
                </v-card>
              </v-col>
              <v-col cols="12" sm="6" md="3">
                <v-card outlined>
                  <v-card-text class="text-center">
                    <v-icon color="success" size="40">mdi-file-document</v-icon>
                    <div class="headline mt-2">{{ overview.total_articles }}</div>
                    <div class="subtitle-2">总文章数</div>
                  </v-card-text>
                </v-card>
              </v-col>
              <v-col cols="12" sm="6" md="3">
                <v-card outlined>
                  <v-card-text class="text-center">
                    <v-icon color="info" size="40">mdi-comment</v-icon>
                    <div class="headline mt-2">{{ overview.total_comments }}</div>
                    <div class="subtitle-2">总评论数</div>
                  </v-card-text>
                </v-card>
              </v-col>
              <v-col cols="12" sm="6" md="3">
                <v-card outlined>
                  <v-card-text class="text-center">
                    <v-icon color="warning" size="40">mdi-account-circle</v-icon>
                    <div class="headline mt-2">{{ overview.online_users }}</div>
                    <div class="subtitle-2">在线用户</div>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- 最近7天统计 -->
      <v-col cols="12" md="6">
        <v-card>
          <v-card-title>
            <v-icon left>mdi-chart-line</v-icon>
            最近7天活跃度
          </v-card-title>
          <v-card-text>
            <v-simple-table>
              <template v-slot:default>
                <thead>
                  <tr>
                    <th>日期</th>
                    <th>新增用户</th>
                    <th>活跃用户</th>
                    <th>新增文章</th>
                    <th>新增评论</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="stat in recentStats" :key="stat.date">
                    <td>{{ stat.date }}</td>
                    <td>{{ stat.new_users }}</td>
                    <td>{{ stat.active_users }}</td>
                    <td>{{ stat.new_articles }}</td>
                    <td>{{ stat.new_comments }}</td>
                  </tr>
                </tbody>
              </template>
            </v-simple-table>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- 热门文章 -->
      <v-col cols="12" md="6">
        <v-card>
          <v-card-title>
            <v-icon left>mdi-fire</v-icon>
            热门文章
          </v-card-title>
          <v-card-text>
            <v-list>
              <v-list-item v-for="article in hotArticles" :key="article.id">
                <v-list-item-content>
                  <v-list-item-title>{{ article.title }}</v-list-item-title>
                  <v-list-item-subtitle>
                    <v-icon small>mdi-eye</v-icon> {{ article.view_count }}
                    <v-icon small class="ml-2">mdi-heart</v-icon> {{ article.like_count }}
                    <v-icon small class="ml-2">mdi-comment</v-icon> {{ article.comment_count }}
                  </v-list-item-subtitle>
                </v-list-item-content>
                <v-list-item-action>
                  <v-btn text small color="primary" :to="`/article/${article.id}`">
                    查看
                  </v-btn>
                </v-list-item-action>
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- 活跃用户 -->
      <v-col cols="12">
        <v-card>
          <v-card-title>
            <v-icon left>mdi-account-star</v-icon>
            活跃用户
          </v-card-title>
          <v-card-text>
            <v-row>
              <v-col cols="12" sm="6" md="4" lg="3" v-for="user in activeUsers" :key="user.id">
                <v-card outlined>
                  <v-card-text class="text-center">
                    <v-avatar size="60">
                      <img :src="user.avatar || '/default-avatar.png'" alt="avatar">
                    </v-avatar>
                    <div class="mt-2">{{ user.display_name || user.username }}</div>
                    <v-chip small color="success" class="mt-1">
                      <v-icon small left>mdi-circle</v-icon>
                      在线
                    </v-chip>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { ref, onMounted } from 'vue'
import { statisticsApi } from '@/api'

export default {
  name: 'StatisticsDashboard',
  setup() {
    const overview = ref({
      total_users: 0,
      total_articles: 0,
      total_comments: 0,
      online_users: 0
    })
    const recentStats = ref([])
    const hotArticles = ref([])
    const activeUsers = ref([])

    const loadDashboard = async () => {
      try {
        const res = await statisticsApi.getStatisticsDashboard()
        if (res.data.success) {
          overview.value = res.data.data.overview
          recentStats.value = res.data.data.recent_stats
          hotArticles.value = res.data.data.hot_articles
          activeUsers.value = res.data.data.active_users
        }
      } catch (error) {
        console.error('加载仪表板数据失败:', error)
      }
    }

    onMounted(() => {
      loadDashboard()
    })

    return {
      overview,
      recentStats,
      hotArticles,
      activeUsers
    }
  }
}
</script>

<style scoped>
.v-card {
  margin-bottom: 16px;
}
</style>