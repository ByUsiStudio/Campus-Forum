<template>
  <div class="admin-statistics">
    <div class="layui-card mb-4">
      <div class="layui-card-header">
        <i class="fa-solid fa-chart-column mr-2"></i>
        系统概览
      </div>
      <div class="layui-card-body">
        <div class="layui-row">
          <div class="layui-col-xs12 layui-col-sm6 layui-col-md3">
            <div class="stat-card stat-card-primary">
              <div class="stat-content text-center">
                <i class="fa-solid fa-users" style="font-size: 40px; color: #1E9FFF;"></i>
                <div style="font-size: 28px; font-weight: 700; margin-top: 10px;">{{ overview.total_users }}</div>
                <div style="color: #666; font-size: 14px;">总用户数</div>
              </div>
            </div>
          </div>
          <div class="layui-col-xs12 layui-col-sm6 layui-col-md3">
            <div class="stat-card stat-card-success">
              <div class="stat-content text-center">
                <i class="fa-solid fa-file-lines" style="font-size: 40px; color: #52C41A;"></i>
                <div style="font-size: 28px; font-weight: 700; margin-top: 10px;">{{ overview.total_articles }}</div>
                <div style="color: #666; font-size: 14px;">总文章数</div>
              </div>
            </div>
          </div>
          <div class="layui-col-xs12 layui-col-sm6 layui-col-md3">
            <div class="stat-card stat-card-info">
              <div class="stat-content text-center">
                <i class="fa-solid fa-comment" style="font-size: 40px; color: #3B82F6;"></i>
                <div style="font-size: 28px; font-weight: 700; margin-top: 10px;">{{ overview.total_comments }}</div>
                <div style="color: #666; font-size: 14px;">总评论数</div>
              </div>
            </div>
          </div>
          <div class="layui-col-xs12 layui-col-sm6 layui-col-md3">
            <div class="stat-card stat-card-warning">
              <div class="stat-content text-center">
                <i class="fa-solid fa-user-circle" style="font-size: 40px; color: #FAAD14;"></i>
                <div style="font-size: 28px; font-weight: 700; margin-top: 10px;">{{ overview.online_users }}</div>
                <div style="color: #666; font-size: 14px;">在线用户</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="layui-row">
      <div class="layui-col-xs12 layui-col-md6">
        <div class="layui-card mb-4">
          <div class="layui-card-header">
            <i class="fa-solid fa-chart-line mr-2"></i>
            最近7天活跃度
          </div>
          <div class="layui-card-body">
            <table class="layui-table">
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
            </table>
          </div>
        </div>
      </div>

      <div class="layui-col-xs12 layui-col-md6">
        <div class="layui-card mb-4">
          <div class="layui-card-header">
            <i class="fa-solid fa-flame mr-2"></i>
            热门文章
          </div>
          <div class="layui-card-body">
            <div v-for="article in hotArticles" :key="article.id" class="article-item border-b last:border-b-0" style="padding: 10px 0;">
              <div class="flex items-center justify-between">
                <div>
                  <div style="font-weight: 500;">{{ article.title }}</div>
                  <div style="color: #999; font-size: 12px; margin-top: 4px;">
                    <i class="fa-solid fa-eye mr-2"></i>{{ article.view_count }}
                    <i class="fa-solid fa-heart mr-2 ml-3"></i>{{ article.like_count }}
                    <i class="fa-solid fa-comment mr-2 ml-3"></i>{{ article.comment_count }}
                  </div>
                </div>
                <a :href="`/article/${article.id}`" class="layui-btn layui-btn-sm layui-btn-primary">查看</a>
              </div>
            </div>
            <div v-if="hotArticles.length === 0" class="text-center py-8 text-muted">暂无热门文章</div>
          </div>
        </div>
      </div>
    </div>

    <div class="layui-card">
      <div class="layui-card-header">
        <i class="fa-solid fa-star mr-2"></i>
        活跃用户
      </div>
      <div class="layui-card-body">
        <div class="layui-row">
          <div class="layui-col-xs12 layui-col-sm6 layui-col-md4 layui-col-lg3" v-for="user in activeUsers" :key="user.id">
            <div class="user-card" style="text-align: center; padding: 15px; border: 1px solid #E8E8E8; border-radius: 8px; margin-bottom: 15px;">
              <div class="avatar" style="width: 60px; height: 60px; border-radius: 50%; background: #F5F7FA; margin: 0 auto; overflow: hidden;">
                <img :src="user.avatar || '/default-avatar.png'" alt="avatar" style="width: 100%; height: 100%; object-fit: cover;" />
              </div>
              <div style="margin-top: 10px;">{{ user.display_name || user.username }}</div>
              <span class="layui-badge layui-bg-green" style="margin-top: 5px;">在线</span>
            </div>
          </div>
        </div>
        <div v-if="activeUsers.length === 0" class="text-center py-8 text-muted">暂无活跃用户</div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { statisticsApi } from '../../api'

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
.admin-statistics {
  padding: 24px;
}

.stat-card {
  padding: 15px;
  border: 1px solid #E8E8E8;
  border-radius: 8px;
  background: #fff;
}

.stat-card-primary {
  border-left: 4px solid #1E9FFF;
}

.stat-card-success {
  border-left: 4px solid #52C41A;
}

.stat-card-info {
  border-left: 4px solid #3B82F6;
}

.stat-card-warning {
  border-left: 4px solid #FAAD14;
}

.border-b {
  border-bottom: 1px solid #E8E8E8;
}

.last\:border-b-0:last-child {
  border-bottom: none;
}
</style>