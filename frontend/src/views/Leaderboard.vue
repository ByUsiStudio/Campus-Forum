<template>
  <div class="leaderboard-page">
    <div class="leaderboard-card">
      <div class="card-header">
        <i class="fa-solid fa-podium mr-2"></i>
        <span>排行榜</span>
      </div>

      <div class="card-body">
        <div class="filter-row">
          <div class="filter-item">
            <label class="filter-label">排行榜类型</label>
            <select v-model="selectedType" class="layui-select" @change="loadLeaderboard">
              <option v-for="option in typeOptions" :key="option.value" :value="option.value">
                {{ option.text }}
              </option>
            </select>
          </div>
          <div class="filter-item">
            <label class="filter-label">统计周期</label>
            <select v-model="selectedPeriod" class="layui-select" @change="loadLeaderboard">
              <option v-for="option in periodOptions" :key="option.value" :value="option.value">
                {{ option.text }}
              </option>
            </select>
          </div>
        </div>

        <div class="leaderboard-list">
          <div 
            v-for="(item, index) in leaderboard" 
            :key="item.id"
            class="leaderboard-item"
          >
            <div class="rank-avatar" :class="getRankClass(index)">
              {{ item.rank }}
            </div>

            <div class="user-info">
              <div class="user-name">{{ item.user.display_name || item.user.username }}</div>
              <div class="user-level">
                <span class="level-tag">Lv.{{ item.user.level || 1 }}</span>
                <span class="score-text">{{ getScoreText(item.score) }}</span>
              </div>
            </div>

            <router-link :to="`/user/${item.user.id}`" class="view-link">
              查看
            </router-link>
          </div>
        </div>

        <div class="divider"></div>

        <div class="my-rank-card">
          <div class="my-rank-header">
            <i class="fa-solid fa-user mr-2 text-primary"></i>
            <span class="subtitle-1">我的排名</span>
          </div>
          <div class="my-rank-content">
            <div class="rank-chip">
              <i class="fa-solid fa-medal mr-2"></i>
              第 {{ myRank.rank || '未上榜' }} 名
            </div>
            <span class="score-text">{{ getScoreText(myRank.score || 0) }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
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

    const getRankClass = (index) => {
      if (index === 0) return 'gold'
      if (index === 1) return 'silver'
      if (index === 2) return 'bronze'
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
      getRankClass,
      getScoreText,
      loadLeaderboard
    }
  }
}
</script>

<style scoped>
.leaderboard-page {
  padding: 24px;
  max-width: 800px;
  margin: 0 auto;
}

.leaderboard-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
  overflow: hidden;
}

.card-header {
  display: flex;
  align-items: center;
  padding: 16px 24px;
  border-bottom: 1px solid #f0f0f0;
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.card-body {
  padding: 24px;
}

.filter-row {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.filter-item {
  display: flex;
  flex-direction: column;
}

.filter-label {
  font-size: 12px;
  color: #999;
  margin-bottom: 8px;
}

.layui-select {
  padding: 10px 12px;
  border: 1px solid #e8e8e8;
  border-radius: 6px;
  font-size: 14px;
  color: #333;
  outline: none;
  
  &:focus {
    border-color: var(--primary);
  }
}

.leaderboard-list {
  display: flex;
  flex-direction: column;
}

.leaderboard-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 0;
  border-bottom: 1px solid #f0f0f0;
  
  &:last-child {
    border-bottom: none;
  }
}

.rank-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  font-weight: 600;
  color: white;
  
  &.gold { background: #FFD700; }
  &.silver { background: #C0C0C0; }
  &.bronze { background: #CD7F32; }
  &.grey { background: #e8e8e8; color: #666; }
}

.user-info {
  flex: 1;
}

.user-name {
  font-size: 15px;
  font-weight: 500;
  color: #333;
  margin-bottom: 4px;
}

.user-level {
  display: flex;
  align-items: center;
  gap: 8px;
}

.level-tag {
  font-size: 12px;
  color: white;
  background: var(--primary);
  padding: 2px 8px;
  border-radius: 4px;
}

.score-text {
  font-size: 13px;
  color: #666;
}

.view-link {
  font-size: 13px;
  color: var(--primary);
  padding: 6px 12px;
  border: 1px solid var(--primary);
  border-radius: 4px;
  
  &:hover {
    background: rgba(30, 159, 255, 0.1);
  }
}

.divider {
  height: 1px;
  background: #f0f0f0;
  margin: 24px 0;
}

.my-rank-card {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 16px;
}

.my-rank-header {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
}

.subtitle-1 {
  font-size: 15px;
  font-weight: 600;
  color: #333;
}

.my-rank-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.rank-chip {
  display: flex;
  align-items: center;
  font-size: 14px;
  font-weight: 600;
  color: white;
  background: var(--primary);
  padding: 8px 16px;
  border-radius: 20px;
}

.mr-2 {
  margin-right: 8px;
}
</style>