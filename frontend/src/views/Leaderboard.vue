<template>
  <v-container>
    <v-card>
      <v-card-title>
        <v-icon left>mdi-podium</v-icon>
        排行榜
      </v-card-title>
      
      <v-card-text>
        <!-- 排行榜类型选择 -->
        <v-row class="mb-4">
          <v-col cols="12" sm="6">
            <v-select
              v-model="selectedType"
              :items="typeOptions"
              label="排行榜类型"
              outlined
              dense
              @change="loadLeaderboard"
            ></v-select>
          </v-col>
          <v-col cols="12" sm="6">
            <v-select
              v-model="selectedPeriod"
              :items="periodOptions"
              label="统计周期"
              outlined
              dense
              @change="loadLeaderboard"
            ></v-select>
          </v-col>
        </v-row>

        <!-- 排行榜列表 -->
        <v-list three-line>
          <v-list-item v-for="(item, index) in leaderboard" :key="item.id">
            <v-list-item-avatar>
              <v-avatar :color="getRankColor(index + 1)" size="40">
                <span class="white--text headline">{{ item.rank }}</span>
              </v-avatar>
            </v-list-item-avatar>

            <v-list-item-content>
              <v-list-item-title>
                {{ item.user.display_name || item.user.username }}
              </v-list-item-title>
              <v-list-item-subtitle>
                <v-chip small color="primary" class="mr-2">
                  Lv.{{ item.user.level || 1 }}
                </v-chip>
                <span>{{ getScoreText(item.score) }}</span>
              </v-list-item-subtitle>
            </v-list-item-content>

            <v-list-item-action>
              <v-btn text small color="primary" :to="`/user/${item.user.id}`">
                查看
              </v-btn>
            </v-list-item-action>
          </v-list-item>
        </v-list>

        <!-- 我的排名 -->
        <v-divider class="my-4"></v-divider>
        <v-card outlined>
          <v-card-text>
            <div class="d-flex align-center">
              <v-icon left color="primary">mdi-account</v-icon>
              <span class="subtitle-1">我的排名</span>
            </div>
            <div class="mt-2">
              <v-chip color="primary" large>
                <v-avatar left>
                  <v-icon>mdi-medal</v-icon>
                </v-avatar>
                第 {{ myRank.rank || '未上榜' }} 名
              </v-chip>
              <span class="ml-2">{{ getScoreText(myRank.score || 0) }}</span>
            </div>
          </v-card-text>
        </v-card>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script>
import { ref, onMounted } from 'vue'
import { leaderboardApi } from '@/api'

export default {
  name: 'Leaderboard',
  setup() {
    const selectedType = ref('experience')
    const selectedPeriod = ref('all_time')
    const leaderboard = ref([])
    const myRank = ref({})

    const typeOptions = [
      { text: '经验值排行', value: 'experience' },
      { text: '文章数排行', value: 'articles' },
      { text: '点赞数排行', value: 'likes' },
      { text: '评论数排行', value: 'comments' },
      { text: '签到排行', value: 'sign_in' },
      { text: '活跃度排行', value: 'active' }
    ]

    const periodOptions = [
      { text: '总排行', value: 'all_time' },
      { text: '本月排行', value: 'monthly' },
      { text: '本周排行', value: 'weekly' },
      { text: '今日排行', value: 'daily' }
    ]

    const getRankColor = (rank) => {
      if (rank === 1) return 'gold'
      if (rank === 2) return 'silver'
      if (rank === 3) return 'bronze'
      return 'grey'
    }

    const getScoreText = (score) => {
      const texts = {
        experience: `${Math.floor(score)} 经验值`,
        articles: `${Math.floor(score)} 篇文章`,
        likes: `${Math.floor(score)} 个点赞`,
        comments: `${Math.floor(score)} 条评论`,
        sign_in: `${Math.floor(score)} 天签到`,
        active: `${score.toFixed(1)} 活跃度`
      }
      return texts[selectedType.value] || `${Math.floor(score)} 分`
    }

    const loadLeaderboard = async () => {
      try {
        const res = await leaderboardApi.getLeaderboard(selectedType.value, selectedPeriod.value)
        if (res.data.success) {
          leaderboard.value = res.data.data.leaderboard
        }

        const rankRes = await leaderboardApi.getUserRank(selectedType.value, selectedPeriod.value)
        if (rankRes.data.success) {
          myRank.value = rankRes.data.data
        }
      } catch (error) {
        console.error('加载排行榜失败:', error)
      }
    }

    onMounted(() => {
      loadLeaderboard()
    })

    return {
      selectedType,
      selectedPeriod,
      leaderboard,
      myRank,
      typeOptions,
      periodOptions,
      getRankColor,
      getScoreText,
      loadLeaderboard
    }
  }
}
</script>

<style scoped>
.v-avatar.gold {
  background-color: #FFD700 !important;
}

.v-avatar.silver {
  background-color: #C0C0C0 !important;
}

.v-avatar.bronze {
  background-color: #CD7F32 !important;
}
</style>