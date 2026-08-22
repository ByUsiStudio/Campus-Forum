<template>
  <div class="page-container leaderboard-page">
    <div class="leaderboard-header">
      <div class="leaderboard-header-icon">
        <el-icon :size="26"><Trophy /></el-icon>
      </div>
      <div class="leaderboard-header-text">
        <span class="leaderboard-title">排行榜</span>
        <span class="leaderboard-subtitle">看看谁的积分更高、更活跃</span>
      </div>
    </div>

    <div class="leaderboard-card">
      <!-- 排行榜类型/周期选择 -->
      <div class="leaderboard-controls">
        <div class="leaderboard-control-row">
          <span class="control-label">排行榜类型</span>
          <el-segmented
            v-model="selectedType"
            :options="typeSegmentedOptions"
            @change="loadLeaderboard"
          />
        </div>
        <div class="leaderboard-control-row">
          <span class="control-label">统计周期</span>
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

      <!-- 颁奖台：前三名 -->
      <div v-if="leaderboard.length" class="podium">
        <template v-for="(item, index) in leaderboard.slice(0, 3)" :key="item.id">
          <div class="podium-item" :class="`podium-rank-${index + 1}`">
            <div class="podium-medal">
              <span class="podium-medal-num">{{ item.rank }}</span>
            </div>
            <el-avatar
              :size="index === 0 ? 72 : 60"
              :src="item.user?.avatar || ''"
              class="podium-avatar"
              @click="goProfile(item.user.id)"
            >
              {{ (item.user.display_name || item.user.username || 'U')[0] }}
            </el-avatar>
            <div class="podium-name">{{ item.user.display_name || item.user.username }}</div>
            <div class="podium-score">{{ getScoreText(item.score) }}</div>
          </div>
        </template>
      </div>

      <!-- 排行榜列表 -->
      <div class="leaderboard-list">
        <div
          v-for="(item, index) in leaderboard"
          :key="item.id"
          class="leaderboard-item"
          :class="{ 'is-top': index === 0, 'is-second': index === 1, 'is-third': index === 2 }"
          @click="goProfile(item.user.id)"
        >
          <span class="rank-badge" :style="{ backgroundColor: getRankColor(index + 1) }">
            <el-icon v-if="index < 3" class="rank-medal"><Medal /></el-icon>
            <span v-else class="rank-num">{{ item.rank || index + 1 }}</span>
          </span>

          <el-avatar
            :size="42"
            :src="item.user?.avatar || ''"
            class="user-avatar"
          >
            {{ (item.user.display_name || item.user.username || 'U')[0] }}
          </el-avatar>

          <div class="leaderboard-item-user">
            <div class="leaderboard-item-name">
              {{ item.user.display_name || item.user.username }}
            </div>
            <div class="leaderboard-item-sub">
              <span class="level-chip">Lv.{{ item.user.level || 1 }}</span>
              <span class="text-secondary">{{ getScoreText(item.score) }}</span>
            </div>
          </div>

          <el-button
            link
            type="primary"
            class="view-btn"
            @click.stop="goProfile(item.user.id)"
          >
            查看
            <el-icon class="view-icon"><ArrowRight /></el-icon>
          </el-button>
        </div>
      </div>

      <!-- 我的排名 -->
      <div class="my-rank-card">
        <div class="my-rank-head">
          <el-icon color="var(--campus-primary)" :size="18"><User /></el-icon>
          <span class="my-rank-label">我的排名</span>
        </div>
        <div class="my-rank-value">
          <span class="my-rank-tag">
            <el-icon class="mr-1"><Medal /></el-icon>
            第 {{ myRank.rank || '未上榜' }} 名
          </span>
          <span class="text-secondary my-rank-score">{{ getScoreText(myRank.score || 0) }}</span>
        </div>
      </div>
    </div>
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
      return '#cbd2e1'
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

/* 顶部标题 */
.leaderboard-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 18px;
}

.leaderboard-header-icon {
  width: 52px;
  height: 52px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 16px;
  background: linear-gradient(135deg, #ffb347, #ff6b6b);
  color: #fff;
  box-shadow: 0 6px 16px rgba(255, 150, 90, 0.35);
}

.leaderboard-header-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.leaderboard-title {
  font-size: 1.5rem;
  font-weight: 800;
  color: var(--campus-text);
  line-height: 1.2;
}

.leaderboard-subtitle {
  font-size: 0.875rem;
  color: var(--campus-text-secondary);
}

/* 主卡片 */
.leaderboard-card {
  background: var(--campus-surface);
  border: 1px solid var(--campus-border);
  border-radius: var(--campus-radius);
  box-shadow: var(--campus-shadow);
  padding: 20px;
}

/* 控制区 */
.leaderboard-controls {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 16px;
}

.leaderboard-control-row {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.control-label {
  font-size: 0.9rem;
  font-weight: 600;
  color: var(--campus-text);
  flex-shrink: 0;
  min-width: 76px;
}

/* 颁奖台 */
.podium {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
  align-items: end;
  padding: 16px 4px 8px;
  border-radius: 14px;
  background: linear-gradient(180deg, rgba(79, 110, 247, 0.06), rgba(79, 110, 247, 0))
}

.podium-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 14px 8px;
  text-align: center;
  position: relative;
}

.podium-rank-1 {
  order: 2;
}

.podium-rank-2 {
  order: 1;
}

.podium-rank-3 {
  order: 3;
}

.podium-medal {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  font-size: 0.9rem;
  font-weight: 800;
  color: #4a3b00;
  border: 2px solid rgba(0, 0, 0, 0.08);
}

.podium-rank-1 .podium-medal {
  background: #ffd700;
  box-shadow: 0 4px 12px rgba(255, 215, 0, 0.5);
}

.podium-rank-2 .podium-medal {
  background: #d7dbe4;
  box-shadow: 0 4px 12px rgba(192, 192, 192, 0.45);
}

.podium-rank-3 .podium-medal {
  background: #f4b27a;
  box-shadow: 0 4px 12px rgba(205, 127, 50, 0.4);
}

.podium-avatar {
  background: var(--campus-primary);
  color: #fff;
  font-weight: 700;
  border: 3px solid var(--campus-surface);
  box-shadow: 0 4px 12px rgba(15, 23, 42, 0.12);
  cursor: pointer;
}

.podium-name {
  font-weight: 700;
  font-size: 0.9rem;
  color: var(--campus-text);
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.podium-score {
  font-size: 0.8rem;
  color: var(--campus-primary);
  font-weight: 600;
}

/* 排行榜列表 */
.leaderboard-list {
  margin-top: 8px;
}

.leaderboard-item {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 12px 10px;
  border-radius: 12px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.leaderboard-item:hover {
  background: rgba(79, 110, 247, 0.06);
}

.leaderboard-item.is-top {
  background: rgba(255, 215, 0, 0.1);
}

.leaderboard-item.is-second {
  background: rgba(192, 192, 192, 0.1);
}

.leaderboard-item.is-third {
  background: rgba(205, 127, 50, 0.1);
}

.rank-badge {
  flex-shrink: 0;
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  color: #fff;
  font-weight: 700;
  font-size: 0.85rem;
}

.rank-medal {
  font-size: 15px;
}

.rank-num {
  font-size: 0.85rem;
}

.user-avatar {
  flex-shrink: 0;
  background: var(--campus-primary);
  color: #fff;
  font-weight: 700;
}

.leaderboard-item-user {
  flex: 1;
  min-width: 0;
}

.leaderboard-item-name {
  font-weight: 600;
  color: var(--campus-text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.leaderboard-item-sub {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
  font-size: 0.8rem;
}

.level-chip {
  font-size: 0.72rem;
  font-weight: 600;
  color: var(--campus-primary);
  background: rgba(79, 110, 247, 0.1);
  padding: 1px 8px;
  border-radius: 999px;
}

.view-btn {
  flex-shrink: 0;
}

.view-icon {
  margin-left: 2px;
  vertical-align: -2px;
}

/* 我的排名 */
.my-rank-card {
  margin-top: 16px;
  padding: 16px;
  border-radius: 14px;
  background: linear-gradient(135deg, #eef2ff, #f7f4ff);
  border: 1px solid #e3e4ff;
}

.my-rank-head {
  display: flex;
  align-items: center;
  gap: 6px;
}

.my-rank-label {
  font-size: 1rem;
  font-weight: 700;
  color: var(--campus-text);
}

.my-rank-value {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 12px;
}

.my-rank-tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 1rem;
  font-weight: 700;
  color: var(--campus-primary);
  background: var(--campus-surface);
  border-radius: 999px;
  padding: 6px 16px;
  box-shadow: 0 2px 8px rgba(79, 110, 247, 0.15);
}

.my-rank-score {
  font-size: 0.9rem;
}

.mr-1 {
  margin-right: 4px;
}

/* 移动端 */
@media (max-width: 575px) {
  .leaderboard-control-row .el-segmented {
    max-width: 100%;
    overflow-x: auto;
    -webkit-overflow-scrolling: touch;
    flex-shrink: 1;
  }

  .leaderboard-controls {
    gap: 8px;
  }

  .leaderboard-card {
    padding: 14px;
  }

  .leaderboard-item {
    padding: 10px 6px;
    gap: 10px;
  }

  .podium-name {
    font-size: 0.8rem;
  }
}
</style>
