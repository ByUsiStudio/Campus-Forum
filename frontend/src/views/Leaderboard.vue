<template>
  <div class="page-container leaderboard-page">
    <el-card class="mb-4 leaderboard-card">
      <template #header>
        <div class="d-flex align-center leaderboard-header">
          <el-icon :size="24" class="leaderboard-header-icon"><Trophy /></el-icon>
          <span class="leaderboard-title">排行榜</span>
        </div>
      </template>

      <!-- 排行榜类型选择 -->
      <div class="leaderboard-controls mb-4">
        <div class="leaderboard-control-row">
          <span class="text-secondary control-label">排行榜类型</span>
          <el-segmented
            v-model="selectedType"
            :options="typeSegmentedOptions"
            @change="loadLeaderboard"
          />
        </div>
        <div class="leaderboard-control-row">
          <span class="text-secondary control-label">统计周期</span>
          <el-select
            v-model="selectedPeriod"
            style="width: 180px"
            @change="loadLeaderboard"
          >
            <el-option
              v-for="opt in periodOptions"
              :key="opt.value"
              :label="opt.text"
              :value="opt.value"
            />
          </el-select>
        </div>
      </div>

      <!-- 排行榜列表 -->
      <el-list class="leaderboard-list">
        <el-list-item v-for="(item, index) in leaderboard" :key="item.id" class="leaderboard-item">
          <div class="leaderboard-item-main">
            <el-avatar
              :size="40"
              :style="{ backgroundColor: getRankColor(index + 1) }"
              class="rank-avatar"
            >
              <span class="rank-num">{{ item.rank }}</span>
            </el-avatar>

            <div class="leaderboard-item-user">
              <div class="leaderboard-item-name">
                {{ item.user.display_name || item.user.username }}
              </div>
              <div class="leaderboard-item-sub">
                <el-tag size="small" type="primary" class="mr-2">
                  Lv.{{ item.user.level || 1 }}
                </el-tag>
                <span class="text-secondary">{{ getScoreText(item.score) }}</span>
              </div>
            </div>

            <el-button
              link
              type="primary"
              @click="goProfile(item.user.id)"
            >
              查看
            </el-button>
          </div>
        </el-list-item>
      </el-list>

      <!-- 我的排名 -->
      <el-divider class="mt-4" />
      <div class="card-surface my-rank-card">
        <div class="my-rank-title">
          <el-icon color="var(--campus-primary)"><User /></el-icon>
          <span class="text-secondary my-rank-label">我的排名</span>
        </div>
        <div class="my-rank-value">
          <el-tag type="primary" size="large" effect="plain" class="mr-2">
            <el-icon class="mr-1"><Medal /></el-icon>
            第 {{ myRank.rank || '未上榜' }} 名
          </el-tag>
          <span class="text-secondary">{{ getScoreText(myRank.score || 0) }}</span>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Trophy, User, Medal } from '@element-plus/icons-vue'
import { leaderboardApi } from '@/api'

export default {
  name: 'Leaderboard',
  setup() {
    const router = useRouter()
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

    // el-segmented options computed from typeOptions
    const typeSegmentedOptions = typeOptions.map(t => ({
      label: t.text,
      value: t.value
    }))

    const periodOptions = [
      { text: '总排行', value: 'all_time' },
      { text: '本月排行', value: 'monthly' },
      { text: '本周排行', value: 'weekly' },
      { text: '今日排行', value: 'daily' }
    ]

    const getRankColor = (rank) => {
      if (rank === 1) return '#FFD700'
      if (rank === 2) return '#C0C0C0'
      if (rank === 3) return '#CD7F32'
      return '#909399'
    }

    const getScoreText = (score) => {
      const value = Number(score) || 0
      const texts = {
        experience: `${Math.floor(value)} 经验值`,
        articles: `${Math.floor(value)} 篇文章`,
        likes: `${Math.floor(value)} 个点赞`,
        comments: `${Math.floor(value)} 条评论`,
        sign_in: `${Math.floor(value)} 天签到`,
        active: `${value.toFixed(1)} 活跃度`
      }
      return texts[selectedType.value] || `${Math.floor(value)} 分`
    }

    const goProfile = (userId) => {
      router.push(`/profile?id=${userId}`)
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
      typeSegmentedOptions,
      periodOptions,
      getRankColor,
      getScoreText,
      goProfile,
      loadLeaderboard
    }
  }
}
</script>

<style scoped>
.leaderboard-page {
  padding-top: 16px;
}

.leaderboard-card {
  border-radius: var(--campus-radius);
}

.leaderboard-header {
  font-size: 1.25rem;
  font-weight: 600;
}

.leaderboard-header-icon {
  color: var(--campus-primary);
}

.leaderboard-title {
  line-height: 1.3;
}

.leaderboard-controls {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.leaderboard-control-row {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.control-label {
  font-size: 0.9rem;
  flex-shrink: 0;
}

.leaderboard-list {
  --el-list-item-paddings: 12px 0;
}

.leaderboard-item-main {
  display: flex;
  align-items: center;
  width: 100%;
  gap: 16px;
}

.rank-avatar {
  flex-shrink: 0;
}

.rank-num {
  color: #fff;
  font-size: 1.05rem;
  font-weight: 600;
}

.leaderboard-item-user {
  flex: 1;
  min-width: 0;
}

.leaderboard-item-name {
  font-weight: 500;
  color: var(--campus-text);
}

.leaderboard-item-sub {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 4px;
  font-size: 0.85rem;
}

.my-rank-card {
  padding: 16px;
  margin-top: 12px;
}

.my-rank-title {
  display: flex;
  align-items: center;
  gap: 6px;
}

.my-rank-label {
  font-size: 1rem;
  font-weight: 500;
}

.my-rank-value {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 12px;
}
</style>
